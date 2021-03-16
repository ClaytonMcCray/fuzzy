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
	functionsToTest := []*singleTest{
		{
			FunctionUnderTest: "MyFunction",
			PtrsToInit: []argDescriptor{
				{Type: "myType", Name: "one"},
				{Type: "int", Name: "one"},
				{Type: "float64", Name: "two"},
				{Type: "int", Name: "two"},
				{Type: "float64", Name: "three", IsPointer: true},
			},
			ReceiverName:        argDescriptor{Type: "myRecvrType", Name: "recvr"},
			ReceiverIsPointer:   true,
			ReceiverIsComposite: true,
		},
		{
			FunctionUnderTest: "PlainFunction",
			PtrsToInit: []argDescriptor{
				{Type: "CompAuto", Name: "three"},
				{Type: "bytes.Buffer", Name: "four", IsPointer: true},
			},
		},
	}

	in := &metaData{
		PackageName:         "mypackage",
		CompletePackagePath: "gitter.com/me/mypackage",
		TestsToGenerate:     functionsToTest,
	}

	b, err := in.Do()

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

	if actualTest.ReceiverIsPointer {
		t.Errorf("got true for ReceiverIsPointer, should be false")
	}

	if actualTest.ReceiverIsComposite {
		t.Errorf("got true for ReceiverIsComposite, should be false")
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
