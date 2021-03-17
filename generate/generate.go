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

// TODO:
// 1) There is no special handling of interface arguments right now.
// 	* They will probably need to be handled with a directive like "//fuzzy:concrete type Concrete Interface"
// 		which will do type aliasing to the concrete type.
// 2) It should probably use goimports instead of go/format.
// 3) I think arguments are missed if they are declared with `x, y int` syntax instead of `x int, y int`

//go:embed tmpl.go.txt
var headerTmpl string

type argDescriptor struct {
	Name, Type string

	// FunctionArgIsPointer only matters for function arguments.
	// All receivers are treated as pointers.
	FunctionArgIsPointer bool
}

type singleTest struct {
	FunctionUnderTest string
	PtrsToInit        []argDescriptor
	ReceiverName      argDescriptor
}

type metaData struct {
	PackageName         string
	CompletePackagePath string
	TestsToGenerate     []*singleTest
}

func (st *singleTest) convertFields(fields []*ast.Field) {
	for _, arg := range fields {
		if len(arg.Names) < 1 {
			log.Fatalf("failing in CreateMetaData, arg.Names: %v, too short", arg.Names)
		}

		switch arg.Type.(type) {
		case *ast.StarExpr:
			switch arg.Type.(*ast.StarExpr).X.(type) {
			case *ast.Ident:
				st.PtrsToInit = append(st.PtrsToInit, argDescriptor{
					Name:                 arg.Names[0].Name,
					Type:                 arg.Type.(*ast.StarExpr).X.(*ast.Ident).Name,
					FunctionArgIsPointer: true,
				})

			case *ast.SelectorExpr:
				st.PtrsToInit = append(st.PtrsToInit, argDescriptor{
					Name: arg.Names[0].Name,
					Type: fmt.Sprintf("%s.%s",
						arg.Type.(*ast.StarExpr).X.(*ast.SelectorExpr).X.(*ast.Ident).Name,
						arg.Type.(*ast.StarExpr).X.(*ast.SelectorExpr).Sel.Name),

					FunctionArgIsPointer: true,
				})
			}

		case *ast.SelectorExpr:
			st.PtrsToInit = append(st.PtrsToInit, argDescriptor{
				Name: arg.Names[0].Name,
				Type: fmt.Sprintf("%s.%s", arg.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name,
					arg.Type.(*ast.SelectorExpr).Sel.Name),
				FunctionArgIsPointer: false,
			})

		default:
			st.PtrsToInit = append(st.PtrsToInit, argDescriptor{
				Name: arg.Names[0].Name,
				Type: arg.Type.(*ast.Ident).Name,
			})
		}
	}
}

// createMetaData takes in a FuzzySet and converts its data into the structure metaData,
// needed to execute the test function templates.
func createMetaData(fs *elements.FuzzySet) *metaData {
	md := &metaData{
		PackageName:         fs.PackageName,
		CompletePackagePath: fs.CompletePackagePath,
	}

	fs.Inspect(elements.PlainFuncs, func(fd *ast.FuncDecl) {
		curTest := &singleTest{
			FunctionUnderTest: fd.Name.Name,
		}

		curTest.convertFields(fd.Type.Params.List)
		md.TestsToGenerate = append(md.TestsToGenerate, curTest)
	})

	fs.Inspect(elements.PtrReceivers, func(fd *ast.FuncDecl) {
		curTest := &singleTest{
			FunctionUnderTest: fd.Name.Name,
			ReceiverName: argDescriptor{
				Name: fmt.Sprintf("%s", fd.Recv.List[0].Names[0]),
			},
		}

		switch fd.Recv.List[0].Type.(type) {
		case *ast.Ident:
			curTest.ReceiverName.Type = fmt.Sprintf("%v", fd.Recv.List[0].Type.(*ast.Ident).Name)
		case *ast.StarExpr:
			curTest.ReceiverName.Type =
				fmt.Sprintf("%v", fd.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name)
		}

		curTest.convertFields(fd.Type.Params.List)
		md.TestsToGenerate = append(md.TestsToGenerate, curTest)
	})

	return md
}

func (md *metaData) doTmpl() ([]byte, error) {
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

// TODO: This is the entry point of the package. All work should be
// 	done here.
func Gen(stdout, stderr io.Writer, input *ast.Package) error {
	return errors.New("not implemented")
}
