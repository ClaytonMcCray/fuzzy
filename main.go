package main

import (
	"errors"
	"flag"
	"io"
	"log"
)

func main() {
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
