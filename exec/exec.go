// Package exec runs a command stored in memory.
package exec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"codeberg.org/msantos/embedexe"
	"golang.org/x/sys/unix"
)

const (
	EnvVar     = "EMBEDEXE"
	EnvFlags   = "EMBEDEXE_FLAGS"
	EnvVerbose = "EMBEDEXE_VERBOSE" // enable debug error messages
)

type Cmd struct {
	*exec.Cmd
	Exe []byte
}

func errexit(status int, err error) {
	if os.Getenv(EnvVerbose) != "" {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[1], err)
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

	fd, err := strconv.Atoi(v)
	if err != nil {
		errexit(127, err)
	}

	if os.Getenv(EnvFlags) == "CLOEXEC" {
		if err := embedexe.SetCloseExec(uintptr(fd), true); err != nil {
			errexit(128, err)
		}

		if err := os.Unsetenv(EnvFlags); err != nil {
			errexit(128, err)
		}
	}

	err = embedexe.Exec(uintptr(fd), os.Args[1:], os.Environ())

	errexit(126, err)
}

func Command(exe []byte, argv []string) *Cmd {
	cmd := exec.Command("/proc/self/exe", argv...)
	return &Cmd{
		Cmd: cmd,
		Exe: exe,
	}
}

func CommandContext(ctx context.Context, exe []byte, argv []string) *Cmd {
	cmd := exec.CommandContext(ctx, "/proc/self/exe", argv...)
	return &Cmd{
		Cmd: cmd,
		Exe: exe,
	}
}

func (cmd *Cmd) Run() error {
	if err := cmd.fdopen(); err != nil {
		return err
	}
	return cmd.Cmd.Run()
}

func (cmd *Cmd) Start() error {
	if err := cmd.fdopen(); err != nil {
		return err
	}
	return cmd.Cmd.Start()
}

func (cmd *Cmd) CombinedOutput() ([]byte, error) {
	if err := cmd.fdopen(); err != nil {
		return nil, err
	}
	return cmd.Cmd.CombinedOutput()
}

func (cmd *Cmd) Output() ([]byte, error) {
	if err := cmd.fdopen(); err != nil {
		return nil, err
	}
	return cmd.Cmd.Output()
}

func (cmd *Cmd) fdopen() error {
	// 0: /proc/self/exe
	// 1: os.Args[0]
	if len(cmd.Args) < 2 {
		return unix.EINVAL
	}

	fd, err := embedexe.Open(cmd.Exe, os.Args[1])
	if err != nil {
		return err
	}

	environ, err := fdset(fd)
	if err != nil {
		return err
	}

	cmd.Env = append(cmd.Env, environ...)

	return nil
}

func fdset(fd uintptr) ([]string, error) {
	env := make([]string, 0)
	if embedexe.CloseExec(fd) {
		env = append(env, EnvFlags+"=CLOEXEC")
		if err := embedexe.SetCloseExec(fd, false); err != nil {
			return env, err
		}
	}
	return append(env, fmt.Sprintf("%s=%d", EnvVar, int(fd))), nil
}
