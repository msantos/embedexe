// Fsexe forks and runs an executable from an embedded directory of
// executables.
//
// To setup the executable directory:
//
//	mkdir bin
//	cp /bin/echo bin
//	cp /bin/sh bin
//	cp /bin/ls bin
//	go build
//	./fsexe bin/ls -al
package main

import (
	"embed"
	"flag"
	"log"
	"os"

	"codeberg.org/msantos/embedexe/exec"
)

//go:embed bin/*
var bin embed.FS

func main() {
	verbose := flag.Bool("verbose", false, "Enable debug messages")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	b, err := bin.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command(b, flag.Args()[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Verbose = *verbose

	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}
