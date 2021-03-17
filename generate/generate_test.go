package generate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/ClaytonMcCray/fuzzy/elements"
)

// TestTemplate doesn't make any strong assertions about the output. It needs to
// be run in verbose mode (go test -v ./generate) so that the output can be inspected.
func TestTemplate(t *testing.T) {
	src := `
package mypackage

type MyRecvrType struct {}
func (mrt *MyRecvrType) MyFunc(one int, two int, three *bytes.Buffer) error {
	return nil
}

func PlainFunc(x *ast.Ident) {}
`
	parsed := setupAndParse(t, src)
	fs := elements.NewFuzzySet("mypackage", "git.com/me/mine/mypackage")
	ast.Walk(fs, parsed)

	in := CreateMetaData(fs)
	b, err := in.DoTmpl()

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", b)
}

func setupAndParse(t *testing.T, toParse string) *ast.File {
	expr, err := parser.ParseFile(token.NewFileSet(), "src.go", toParse, 0)

	if err != nil {
		t.Fatalf("got %s trying to parse toParse; expected nil", err)
	}

	return expr
}

func assertEquivSingleTest(t *testing.T, actualTest, est *singleTest) {
	if actualTest.FunctionUnderTest != est.FunctionUnderTest {
		t.Errorf("expected %s, got %s", est.FunctionUnderTest, actualTest.FunctionUnderTest)
	}

	for i, actual := range actualTest.PtrsToInit {
		if actual != est.PtrsToInit[i] {
			t.Errorf("expected %v, got %v", est.PtrsToInit[i], actual)
		}
	}

	if actualTest.ReceiverName != est.ReceiverName {
		t.Errorf("expected %v, got %v", est.ReceiverName, actualTest.ReceiverName)
	}
}

func TestCreateMetaDataPlainFunction(t *testing.T) {
	src := `
package main

func TestMe(x int, y *float64, z bytes.Buffer, a *ast.Ident) error {
	return nil
}

func TestMe2(b, c int) {}
`
	parsed := setupAndParse(t, src)
	fs := elements.NewFuzzySet("main", "my/main")
	ast.Walk(fs, parsed)

	md := CreateMetaData(fs)

	if md.CompletePackagePath != "my/main" {
		t.Errorf("wrong CompletePackagePath, got %s", md.CompletePackagePath)
	}

	if md.PackageName != "main" {
		t.Errorf("wrong PackageName, got %s", md.PackageName)
	}

	actualTest := md.TestsToGenerate[0]

	est := &singleTest{
		FunctionUnderTest: "TestMe",
		PtrsToInit: []argDescriptor{
			{"x", "int", false},
			{"y", "float64", true},
			{"z", "bytes.Buffer", false},
			{"a", "ast.Ident", true},
		},
	}

	assertEquivSingleTest(t, actualTest, est)

	actualTest = md.TestsToGenerate[1]
	est = &singleTest{
		FunctionUnderTest: "TestMe2",
		PtrsToInit: []argDescriptor{
			{"b", "int", false},
			{"c", "int", false},
		},
	}

	assertEquivSingleTest(t, actualTest, est)
}

func TestCreateMetaDataRecvr(t *testing.T) {
	src := `
package main

type X struct{}
func (x X) XFuncOnePtr(a int, b *int) {}
func (x *X) XFuncTwoNotPtr(a, b string) {}
`
	parsed := setupAndParse(t, src)
	est1 := &singleTest{
		FunctionUnderTest: "XFuncOnePtr",
		PtrsToInit: []argDescriptor{
			{"a", "int", false},
			{"b", "int", true},
		},

		ReceiverName: argDescriptor{
			Name: "x",
			Type: "X",
		},
	}

	est2 := &singleTest{
		FunctionUnderTest: "XFuncTwoNotPtr",
		PtrsToInit: []argDescriptor{
			{"a", "string", false},
			{"b", "string", false},
		},

		ReceiverName: argDescriptor{
			Name: "x",
			Type: "X",
		},
	}

	fs := elements.NewFuzzySet("main", "my/main")
	ast.Walk(fs, parsed)

	md := CreateMetaData(fs)
	t.Log("finished CreateMetaData")

	if md.CompletePackagePath != "my/main" {
		t.Errorf("wrong CompletePackagePath, got %s", md.CompletePackagePath)
	}

	if md.PackageName != "main" {
		t.Errorf("wrong PackageName, got %s", md.PackageName)
	}

	t.Logf("md.TestsToGenerate: %v", md.TestsToGenerate)

	assertEquivSingleTest(t, md.TestsToGenerate[0], est1)
	assertEquivSingleTest(t, md.TestsToGenerate[1], est2)
}
