package elements

import (
	"go/ast"
)

type fuzzable int

const (
	PlainFuncs fuzzable = iota
	PtrReceivers
)

type FuzzySet struct {
	PackageName         string
	CompletePackagePath string

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

func NewFuzzySet(packageName, completePackagePath string) *FuzzySet {
	return &FuzzySet{
		PackageName:         packageName,
		CompletePackagePath: completePackagePath,
	}
}

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
