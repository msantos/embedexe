package exec_test

import (
	"bytes"
	"os"
	"testing"

	"codeberg.org/msantos/embedexe/exec"
)

func TestCommand(t *testing.T) {
	exe, err := os.ReadFile("/bin/echo")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	stdout := bytes.Buffer{}
	cmd := exec.Command(exe, "test")
	cmd.Env = append(os.Environ(), "EMBEDEXE_VERBOSE=1")
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	if err := cmd.Run(); err != nil {
		t.Errorf("%v", err)
		return
	}

	if stdout.String() != "test\n" {
		t.Errorf("expected: test, got: %v", stdout.String())
		return
	}
}

func TestCommandNullArgv(t *testing.T) {
	exe, err := os.ReadFile("/bin/echo")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	cmd := exec.Command(exe)

	if err := cmd.Run(); err != nil {
		t.Errorf("%v", err)
		return
	}
}
