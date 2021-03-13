package elements

import (
	"go/ast"
)

type fuzzable int

const (
	PlainFuncs fuzzable = iota
	PtrReceivers
	AutoReceivers
)

type FuzzySet struct {
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

	// autoRecvrs require:
	// 1) Know that the receiver is local
	// 2) Know the arguments
	// 3) Fuzz the type, which requires taking its address
	// 4) Fuzz the inputs
	// 5) Call the functions
	autoRecvrs []*ast.FuncDecl
}

func NewFuzzySet() *FuzzySet {
	return &FuzzySet{}
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

	case AutoReceivers:
		for _, val := range fs.autoRecvrs {
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
			for _, field := range concrete.Recv.List {
				switch field.Type.(type) {
				case *ast.StarExpr:
					fs.ptrRecvrs = append(fs.ptrRecvrs, concrete)
				case *ast.Ident:
					fs.autoRecvrs = append(fs.autoRecvrs, concrete)
				}
			}
		}
	}

	return fs
}
