package fdexec_test

import (
	"os"
	"testing"

	"codeberg.org/msantos/embedexe"
	"codeberg.org/msantos/embedexe/fdexec"
)

func TestCommand(t *testing.T) {
	exe, err := os.ReadFile("/bin/echo")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	fd, err := embedexe.Open(exe, "echo")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	cmd := fdexec.Command(fd, "test")
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

	fd, err := embedexe.Open(exe, "echo")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	cmd := fdexec.Command(fd)

	if err := cmd.Run(); err != nil {
		t.Errorf("%v", err)
		return
	}
}
