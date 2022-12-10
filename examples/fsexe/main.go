package main

import (
	"embed"
	"log"
	"os"

	"codeberg.org/msantos/embedexe/exec"
)

//go:embed exe/*
var exe embed.FS

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("usage:", os.Args[0], "<command>")
	}

	b, err := exe.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command(b, os.Args[2:])

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}
