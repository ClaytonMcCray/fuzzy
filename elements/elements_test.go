package elements

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func setupAndParse(t *testing.T, toParse string) *ast.File {
	expr, err := parser.ParseFile(token.NewFileSet(), "src.go", toParse, 0)

	if err != nil {
		t.Fatalf("got %s trying to parse toParse; expected nil", err)
	}

	return expr
}

func TestFuzzySetRegister(t *testing.T) {
	toParse := `package main

var (
	unexported int
	ExportedString string
	ExportedInt int
)

func F(i int, s string) error { return nil }
func Y() error { return nil }

type Z map[int]string
func (z Z) DoNothing(arg ...int) (string, error) {
	return "", nil
}

type X struct {}
func (x X) String() string {
	return "str"
}
func (x *X) PtrString() string {
	return "ptrstring"
}
func (x *X) Noop() {
	return
}
`
	expr := setupAndParse(t, toParse)
	fs := NewFuzzySet("main", "my/main")

	ast.Walk(fs, expr)

	expectedFuncs := 2
	expectedPtrs := 4

	if len(fs.funcs) != expectedFuncs {
		t.Fatalf("wrong numbder of fs.funcs, expected %d, got %d", expectedFuncs, len(fs.funcs))
	}

	if len(fs.ptrRecvrs) != expectedPtrs {
		t.Errorf("wrong number of fs.ptrRecvrs, expected %d got %d", expectedPtrs, len(fs.ptrRecvrs))
	}
}
