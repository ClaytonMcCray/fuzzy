package main

import (
	"errors"
	"flag"
	"io"
	"log"
)

func main() {
}

// TODO: when parsing, make sure to filter out *Test.go files, and note in help dialogue
//	 that fuzz tests should be written to a *Test.go file.
func Run(stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	log.SetOutput(stdout)
	flags := flag.NewFlagSet(args[0], flag.PanicOnError)
	flags.Parse(args[1:])

	if flags.NArg() < 2 {
		log.Printf("%d", flags.NArg())
		return errors.New("Usage: fuzzy DIRECTORY PACKAGE [-o OUTPUT]")
	}
	return nil
}
