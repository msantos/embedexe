package exec_test

import (
	"bytes"
	"os"
	"testing"

	"codeberg.org/msantos/embedexe/exec"
	"golang.org/x/sys/unix"
)

func TestCommand(t *testing.T) {
	exe, err := os.ReadFile("/bin/echo")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	stdout := bytes.Buffer{}
	cmd := exec.Command(exe, []string{"procname", "test"})
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

	cmd := exec.Command(exe, []string{})

	if err := cmd.Run(); err != unix.EINVAL {
		t.Errorf("%v", err)
		return
	}
}
