package elements

import (
	"go/ast"
)

type fuzzable int

const (
	PlainFuncs fuzzable = iota
	PtrReceivers
)

// FuzzySet holds the set of functions for which fuzz tests are generated.
type FuzzySet struct {
	// PackageName is the name used in selector statements into the package.
	PackageName string

	// funcs only require:
	// 1) Know the function name and arguments
	// 2) Fuzz the inputs
	// 	* identify whether or not they are pointers!
	// 3) Call the functions
	funcs []*ast.FuncDecl

	// ptrRecvrs require:
	// 1) Know that the receiver is a pointer
	// 2) Know the arguments
	// 3) Fuzz the type
	// 4) Fuzz the inputs
	// 5) Call the functions
	ptrRecvrs []*ast.FuncDecl
}

// NewFuzzySet creates a FuzzySet where packageName is the name of
// the package and completePackagePath is the import path other packages
// would use.
func NewFuzzySet(packageName string) *FuzzySet {
	return &FuzzySet{
		PackageName: packageName,
	}
}

// Inspect calls do on each function declaration of the requested type.
func (fs *FuzzySet) Inspect(which fuzzable, do func(*ast.FuncDecl)) {
	switch which {
	case PlainFuncs:
		for _, val := range fs.funcs {
			do(val)
		}

	case PtrReceivers:
		for _, val := range fs.ptrRecvrs {
			do(val)
		}
	}
}

// Visit is used to add decls to  the FuzzySet. It satisfies
// the ast.Visitor interface.
func (fs *FuzzySet) Visit(n ast.Node) ast.Visitor {
	switch concrete := n.(type) {
	case *ast.FuncDecl:
		if concrete.Recv == nil {
			fs.funcs = append(fs.funcs, concrete)
		} else {
			for range concrete.Recv.List {
				fs.ptrRecvrs = append(fs.ptrRecvrs, concrete)
			}
		}
	}

	return fs
}
