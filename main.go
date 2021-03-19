package main

import (
	"bufio"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/ClaytonMcCray/fuzzy/generate"
)

func main() {
	dirPath := flag.String("i", "", "The path to the directory you wish to fuzz. Required.")
	packageToFuzz := flag.String("package", "", "The package to fuzz in -i. Required.")
	outputFile := flag.String("o", "fuzzy_fuzz_test.go", "The path to the output file. It should end in _test.go")
	flag.Parse()

	if *dirPath == "" || *packageToFuzz == "" {
		log.Fatal("See usage. -i must be set. -package must be set.")
	}

	f, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}

	output := bufio.NewWriter(f)
	fset := token.NewFileSet()
	parsed, err := parser.ParseDir(fset, *dirPath, func(fi fs.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	toFuzz, ok := parsed[*packageToFuzz]
	if !ok {
		log.Fatalf("Failed to parse package %s", *packageToFuzz)
	}

	if err = Run(toFuzz, output, os.Stderr); err != nil {
		log.Fatal(err)
	}
}

func Run(input *ast.Package, stdout, stderr io.Writer) error {
	err := generate.Gen(stdout, stderr, input)
	if err != nil {
		return err
	}

	return nil
}
