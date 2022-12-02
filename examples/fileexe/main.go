// Fileexe forks and runs an embedded executable.
//
//	cp /bin/echo .
//	go build
//	./filexe 3 hello world
package main

import (
	_ "embed"
	"log"
	"os"

	"codeberg.org/msantos/embedexe/exec"
)

//go:embed exe
var exe []byte

func main() {
	cmd := exec.Command(exe, os.Args)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		log.Fatalln("run:", cmd, err)
	}
}
