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
	fs := elements.NewFuzzySet("mypackage")
	ast.Walk(fs, parsed)

	in := createMetaData(fs)
	b, err := in.doTmpl()

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

func TestMe3(sliceType []int,
	slicePointerType []*int,
	arrayType [10]int,
	arrayPointerType [10]*int,
	arraySelectorType [10]bytes.Buffer,
	arraySelectorPointerType [10]*bytes.Buffer,
	sliceSelectorType []bytes.Buffer,
	sliceSelectorPointerType []*bytes.Buffer,
	ptrSlice *[]int,
	ptrSlicePtr *[]*int,
	ptrArray *[10]int,
	ptryArrayPtr *[10]*int,
	ptrSliceSelector *[]bytes.Buffer,
	ptrSlicePtrSelector *[]*bytes.Buffer,
	ptrArraySelector *[10]bytes.Buffer,
	ptrArrayPtrSelector *[10]*bytes.Buffer) {}
`
	parsed := setupAndParse(t, src)
	fs := elements.NewFuzzySet("main")
	ast.Walk(fs, parsed)

	md := createMetaData(fs)

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

	actualTest = md.TestsToGenerate[2]
	est = &singleTest{
		FunctionUnderTest: "TestMe3",
		PtrsToInit: []argDescriptor{
			{"sliceType", "[]int", false},
			{"slicePointerType", "[]*int", false},
			{"arrayType", "[10]int", false},
			{"arrayPointerType", "[10]*int", false},
			{"arraySelectorType", "[10]bytes.Buffer", false},
			{"arraySelectorPointerType", "[10]*bytes.Buffer", false},
			{"sliceSelectorType", "[]bytes.Buffer", false},
			{"sliceSelectorPointerType", "[]*bytes.Buffer", false},
			{"ptrSlice", "[]int", true},
			{"ptrSlicePtr", "[]*int", true},
			{"ptrArray", "[10]int", true},
			{"ptrArrayPtr", "[10]*int", true},
			{"ptrSliceSelector", "[]bytes.Buffer", true},
			{"ptrSlicePtrSelector", "[]*bytes.Buffer", true},
			{"ptrArraySelector", "[10]bytes.Buffer", true},
			{"ptrArrayPtrSelector", "[10]*bytes.Buffer", true},
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

	fs := elements.NewFuzzySet("main")
	ast.Walk(fs, parsed)

	md := createMetaData(fs)
	t.Log("finished CreateMetaData")

	if md.PackageName != "main" {
		t.Errorf("wrong PackageName, got %s", md.PackageName)
	}

	t.Logf("md.TestsToGenerate: %v", md.TestsToGenerate)

	assertEquivSingleTest(t, md.TestsToGenerate[0], est1)
	assertEquivSingleTest(t, md.TestsToGenerate[1], est2)
}
