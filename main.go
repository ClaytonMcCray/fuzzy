package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"
)

func main() {
	if err := Run(os.Stdin, os.Stdout, os.Stderr, os.Args); err != nil {
		log.Fatal(err)
	}
}

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
