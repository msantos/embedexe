// Fileexe forks and runs an embedded executable.
//
// Copy an executable into the build directory:
//
//	cp /bin/echo bin
//	go build
//	./filexe hello world
package main

import (
	_ "embed"
	"log"
	"os"

	"codeberg.org/msantos/embedexe/exec"
)

//go:embed bin
var bin []byte

func main() {
	cmd := exec.Command(bin, os.Args[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln("run:", cmd, err)
	}
}
