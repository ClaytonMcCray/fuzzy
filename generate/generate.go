package generate

import (
	"bytes"
	_ "embed"
	"errors"
	"go/ast"
	"go/format"
	"io"
	"text/template"

	"github.com/ClaytonMcCray/fuzzy/elements"
)

//go:embed tmpl.go.txt
var headerTmpl string

type argPair struct {
	Name, Type string
}

type singleTest struct {
	FunctionUnderTest            string
	CompositePtrsToInit          []argPair
	CompositeAutosToInit         []argPair
	SimplePtrsToInit             []argPair
	SimpleAutosToInit            []argPair
	FunctionArgumentNamesInOrder []string
	ReceiverName                 argPair
	ReceiverIsPointer            bool
	ReceiverIsComposite          bool
}

type metaData struct {
	PackageName         string
	CompletePackagePath string
	TestsToGenerate     []singleTest
}

func Gen(stdout, stderr io.Writer, input *ast.Package) error {
	fs := elements.NewFuzzySet()
	ast.Walk(fs, input)

	return errors.New("not implemented")
}

func (md metaData) Do() ([]byte, error) {
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
