// Package reexec reexecs the process image using a file descriptor.
package reexec

import (
	"fmt"
	"os"
	"strconv"

	"codeberg.org/msantos/embedexe"
)

const (
	EnvVar     = "EMBEDEXE"
	EnvFlags   = "EMBEDEXE_FLAGS"
	EnvVerbose = "EMBEDEXE_VERBOSE" // enable debug error messages

	CLOEXEC = "CLOEXEC"
)

func errexit(status int, err error) {
	if os.Getenv(EnvVerbose) != "" {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
	}
	os.Exit(status)
}

func init() {
	v := os.Getenv(EnvVar)
	if v == "" {
		return
	}

	if err := os.Unsetenv(EnvVar); err != nil {
		errexit(128, err)
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		errexit(127, err)
	}

	fd := embedexe.FD(uintptr(i))

	if os.Getenv(EnvFlags) == CLOEXEC {
		if err := fd.SetCloseExec(true); err != nil {
			errexit(128, err)
		}

		if err := os.Unsetenv(EnvFlags); err != nil {
			errexit(128, err)
		}
	}

	err = fd.Exec(os.Args, os.Environ())

	errexit(126, err)
}
