// Package exec runs a command stored as bytes in memory.
package exec

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"codeberg.org/msantos/embedexe"
	"codeberg.org/msantos/embedexe/internal/reexec"
)

// Cmd is a wrapper around the os/exec Cmd struct.
type Cmd struct {
	*exec.Cmd

	// Exe holds the executable as a byte array.
	Exe []byte

	// The command name (proctitle) stored in /proc/self/comm.
	// Defaults to the command name of the current running process.
	Name string
}

// Command returns the Cmd struct to execute the program held in exe
// with the given arguments.
func Command(exe []byte, arg ...string) *Cmd {
	cmd := exec.Command("/proc/self/exe", arg...)
	return &Cmd{
		Cmd: cmd,
		Exe: exe,
	}
}

// CommandContext returns a Cmd struct using the provided context.
func CommandContext(ctx context.Context, exe []byte, arg ...string) *Cmd {
	cmd := exec.CommandContext(ctx, "/proc/self/exe", arg...)
	return &Cmd{
		Cmd: cmd,
		Exe: exe,
	}
}

// Run starts the specified command and waits for it to complete.
func (cmd *Cmd) Run() error {
	fd, err := cmd.fdopen()
	if err != nil {
		return err
	}
	defer fd.Close()
	return cmd.Cmd.Run()
}

// Start starts the specified command but does not wait for it to complete.
func (cmd *Cmd) Start() error {
	fd, err := cmd.fdopen()
	if err != nil {
		return err
	}
	defer fd.Close()
	return cmd.Cmd.Start()
}

// CombinedOutput runs the command and returns its combined standard
// output and standard error.
func (cmd *Cmd) CombinedOutput() ([]byte, error) {
	fd, err := cmd.fdopen()
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return cmd.Cmd.CombinedOutput()
}

// Output runs the command and returns its standard output.
func (cmd *Cmd) Output() ([]byte, error) {
	fd, err := cmd.fdopen()
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return cmd.Cmd.Output()
}

func (cmd *Cmd) fdopen() (embedexe.FD, error) {
	if cmd.Name == "" {
		cmd.Args[0] = os.Args[0]
	} else {
		cmd.Args[0] = cmd.Name
	}

	fd, err := embedexe.Open(cmd.Exe, cmd.Args[0])
	if err != nil {
		return 0, err
	}

	environ, err := fdset(fd)
	if err != nil {
		_ = fd.Close()
		return 0, err
	}

	cmd.Env = append(cmd.Env, environ...)

	return fd, nil
}

func fdset(fd embedexe.FD) ([]string, error) {
	env := make([]string, 0)
	if fd.CloseExec() {
		env = append(env, reexec.EnvFlags+"="+reexec.CLOEXEC)
		if err := fd.SetCloseExec(false); err != nil {
			return env, err
		}
	}
	return append(env, fmt.Sprintf("%s=%d", reexec.EnvVar, int(fd))), nil
}
