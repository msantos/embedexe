package embedexe_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"codeberg.org/msantos/embedexe"
)

var errInvalidOutput = errors.New("unexpected output")

func run(cmd *exec.Cmd, output string) error {
	var buf bytes.Buffer

	cmd.Stdout = &buf
	cmd.Stderr = &buf
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}

	if !strings.HasPrefix(buf.String(), output) {
		return fmt.Errorf("Expected: %s\nOutput: %s\nError: %w",
			output,
			buf.String(),
			errInvalidOutput,
		)
	}

	return nil
}

func TestOpen(t *testing.T) {
	if os.Getenv("TESTING_EMBEDEXE_TESTOPEN") == "1" {
		ExampleOpen()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestOpen")
	cmd.Env = append(os.Environ(), "TESTING_EMBEDEXE_TESTOPEN=1")

	if err := run(cmd, "test"); err != nil {
		t.Errorf("%v", err)
		return
	}
}

func ExampleOpen() {
	b := []byte("#!/bin/sh\necho $@")

	fd, err := embedexe.Open(b, "example")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fd.Close()

	if err := fd.Exec([]string{"example", "test"}, os.Environ()); err != nil {
		fmt.Println(err)
		return
	}
}

func ExampleFD_Path() {
	b, err := os.ReadFile("/bin/echo")
	if err != nil {
		fmt.Println(err)
		return
	}

	fd, err := embedexe.Open(b, "echo")
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd := exec.Command(fd.Path(), "-n", "test", "abc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return
	}

	// Output: test abc
}

// An example of running a script contained in a memfd without execute
// permissions.
func ExampleFD_Path_sh() {
	b := []byte("#!/bin/sh\necho $@")

	fd, err := embedexe.Open(b, "sh")
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd := exec.Command("/bin/sh", fd.Path(), "-n", "test", "abc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return
	}

	// Output: test abc
}
