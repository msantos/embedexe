package exec_test

import (
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

	cmd := exec.Command(exe, "test")
	cmd.Verbose = true

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	if string(out) != "test\n" {
		t.Errorf("expected: test, got: %v", out)
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
