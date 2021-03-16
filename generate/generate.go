package generate

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"io"
	"log"
	"text/template"

	"github.com/ClaytonMcCray/fuzzy/elements"
)

// TODO: Everything is created as a pointer. If the ast.Expr is not a star expression, we should
//	dereference the value before using it instead of trying to make it an automatic local value.
//	This means the template needs to be updated, and the tests.

// TODO:
// 1) There is no special handling of interface arguments right now.
// 	* They will probably need to be handled with a directive like "//fuzzy:concrete type Concrete Interface"
// 		which will do type aliasing to the concrete type.
// 2) It should probably use goimports instead of go/format.

//go:embed tmpl.go.txt
var headerTmpl string

type argDescriptor struct {
	Name, Type string
	IsPointer  bool
}

type singleTest struct {
	FunctionUnderTest   string
	PtrsToInit          []argDescriptor
	ReceiverName        argDescriptor
	ReceiverIsPointer   bool
	ReceiverIsComposite bool
}

type metaData struct {
	PackageName         string
	CompletePackagePath string
	TestsToGenerate     []*singleTest
}

func (st *singleTest) convertFields(fields []*ast.Field) {
	for _, arg := range fields {
		log.Printf("Working on %s", arg.Names[0].Name)
		if len(arg.Names) < 1 {
			log.Fatalf("failing in CreateMetaData, arg.Names: %v, too short", arg.Names)
		}

		switch arg.Type.(type) {
		case *ast.StarExpr:
			switch arg.Type.(*ast.StarExpr).X.(type) {
			case *ast.Ident:
				log.Printf("%s is *ast.StarExpr", arg.Names[0].Name)
				st.PtrsToInit = append(st.PtrsToInit, argDescriptor{
					Name:      arg.Names[0].Name,
					Type:      arg.Type.(*ast.StarExpr).X.(*ast.Ident).Name,
					IsPointer: true,
				})

			case *ast.SelectorExpr:
				log.Printf("%s is *ast.StarExpr (selector)", arg.Names[0].Name)
				st.PtrsToInit = append(st.PtrsToInit, argDescriptor{
					Name: arg.Names[0].Name,
					Type: fmt.Sprintf("%s.%s",
						arg.Type.(*ast.StarExpr).X.(*ast.SelectorExpr).X.(*ast.Ident).Name,
						arg.Type.(*ast.StarExpr).X.(*ast.SelectorExpr).Sel.Name),

					IsPointer: true,
				})
			}

		case *ast.SelectorExpr:
			log.Printf("%s is *ast.SelectorExpr", arg.Names[0].Name)
			log.Printf("%s is not *ast.StarExpr", arg.Names[0].Name)
			st.PtrsToInit = append(st.PtrsToInit, argDescriptor{
				Name: arg.Names[0].Name,
				Type: fmt.Sprintf("%s.%s", arg.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name,
					arg.Type.(*ast.SelectorExpr).Sel.Name),
				IsPointer: false,
			})

		default:
			log.Printf("%s hit default", arg.Names[0].Name)
			st.PtrsToInit = append(st.PtrsToInit, argDescriptor{
				Name: arg.Names[0].Name,
				Type: arg.Type.(*ast.Ident).Name,
			})
		}
	}
}

// CreateMetaData takes in a FuzzySet and converts its data into the structure metaData,
// needed to execute the test function templates.
func CreateMetaData(fs *elements.FuzzySet) *metaData {
	md := &metaData{
		PackageName:         fs.PackageName,
		CompletePackagePath: fs.CompletePackagePath,
	}

	// Create a singleTest for each plain function under test.
	fs.Inspect(elements.PlainFuncs, func(fd *ast.FuncDecl) {
		curTest := &singleTest{
			FunctionUnderTest: fd.Name.Name,
		}

		curTest.convertFields(fd.Type.Params.List)

		md.TestsToGenerate = append(md.TestsToGenerate, curTest)
	})

	// TODO: handle the receiver functions

	return md
}

func (md *metaData) Do() ([]byte, error) {
	b := &bytes.Buffer{}
	tmpl, err := template.New("package_fuzzy_test").Parse(headerTmpl)
	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(b, md)
	if err != nil {
		return nil, err
	}

	return format.Source(b.Bytes())
}

func Gen(stdout, stderr io.Writer, input *ast.Package) error {
	return errors.New("not implemented")
}
