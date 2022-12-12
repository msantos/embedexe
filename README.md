# embedexe

[![Go Reference](https://pkg.go.dev/badge/codeberg.org/msantos/execve.svg)](https://pkg.go.dev/codeberg.org/msantos/embedexe)

Run a program stored in a byte array such as an executable or a directory
of executables embedded in a Go binary.

# LIMITATIONS

* the executable must either be statically linked or the linked libraries
  available on the filesystem

* only works on Linux (but not on ChromeOS/crostini where presumably
  kernel hardening measures disable executing memfds)

* an embedded executable cannot run another executable embedded in the
  same binary

* an embedded executable cannot use a library embedded in the same binary

# EXAMPLES

* [fsexe](examples/fsexe/main.go): embed a directory of executables
  in a Go binary and run from memory

## Run an embedded executable

```go
// Echotest forks and runs an embedded echo(1).
//
//	cp /bin/echo .
//	go build
//	./echotest hello world
package main

import (
	_ "embed"
	"log"
	"os"

	"codeberg.org/msantos/embedexe/exec"
)

//go:embed echo
var echo []byte

func main() {
	cmd := exec.Command(echo, os.Args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln("run:", cmd, err)
	}
}
```
