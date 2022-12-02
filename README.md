# embedexe

[![Go Reference](https://pkg.go.dev/badge/codeberg.org/msantos/execve.svg)](https://pkg.go.dev/codeberg.org/msantos/embedexe)

Run an executable embedded in a Go binary.

# LIMITATIONS

* the executable must either be statically linked or the linked libraries
  available on the filesytem

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
//    cp /bin/echo .
//    go build
//    ./echotest 3 hello world
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
	cmd := exec.Command(echo, os.Args)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		log.Fatalln("run:", cmd, err)
	}
}
```
