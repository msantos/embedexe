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
		os.Exit(1)
	}
	defer fd.Close()

	if err := fd.Exec([]string{"example", "test"}, os.Environ()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ExampleFD_Path() {
	b, err := os.ReadFile("/bin/echo")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fd, err := embedexe.Open(b, "echo")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := exec.Command(fd.Path(), "-n", "test", "abc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Output: test abc
}
