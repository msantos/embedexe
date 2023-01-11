package fdexec_test

import (
	"os"
	"testing"

	"codeberg.org/msantos/embedexe/fdexec"
)

func TestCommand(t *testing.T) {
	exe, err := os.ReadFile("/bin/echo")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	fd, err := fdexec.Open(exe)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	cmd := fdexec.Command(fd, "test")
	cmd.Env = append(os.Environ(), "EMBEDEXE_VERBOSE=1")

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

	fd, err := fdexec.Open(exe)
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
