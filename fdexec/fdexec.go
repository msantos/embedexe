// Package fdexec runs a command by file descriptor.
package fdexec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

	"codeberg.org/msantos/embedexe"
	"codeberg.org/msantos/embedexe/internal/reexec"
)

// Cmd is a wrapper around the os/exec Cmd struct.
type Cmd struct {
	*exec.Cmd

	// The command name (proctitle) stored in /proc/self/cmdline.
	// Defaults to the command name of the current running process.
	Name string

	// Enable debug messages to stderr.
	Verbose bool

	fd *embedexe.FD // the executable file descriptor
}

// Command returns the Cmd struct to execute the program referenced by
// the file descriptor with the given arguments.
func Command(fd *embedexe.FD, arg ...string) *Cmd {
	cmd := exec.Command("/proc/self/exe", arg...)
	return &Cmd{
		Cmd: cmd,
		fd:  fd,
	}
}

// CommandContext returns a Cmd struct using the provided context.
func CommandContext(ctx context.Context, fd *embedexe.FD, arg ...string) *Cmd {
	cmd := exec.CommandContext(ctx, "/proc/self/exe", arg...)
	return &Cmd{
		Cmd: cmd,
		fd:  fd,
	}
}

// Run starts the specified command and waits for it to complete.
func (cmd *Cmd) Run() error {
	if err := cmd.fdsetenv(); err != nil {
		return err
	}
	return cmd.Cmd.Run()
}

// Start starts the specified command but does not wait for it to complete.
func (cmd *Cmd) Start() error {
	if err := cmd.fdsetenv(); err != nil {
		return err
	}
	return cmd.Cmd.Start()
}

// CombinedOutput runs the command and returns its combined standard
// output and standard error.
func (cmd *Cmd) CombinedOutput() ([]byte, error) {
	if err := cmd.fdsetenv(); err != nil {
		return nil, err
	}
	return cmd.Cmd.CombinedOutput()
}

// Output runs the command and returns its standard output.
func (cmd *Cmd) Output() ([]byte, error) {
	if err := cmd.fdsetenv(); err != nil {
		return nil, err
	}
	return cmd.Cmd.Output()
}

func (cmd *Cmd) fdsetenv() error {
	env := []string{fmt.Sprintf("%s=%d", reexec.EnvVar, int(cmd.fd.FD()))}

	if cmd.fd.CloseExec() {
		env = append(env, reexec.EnvFlags+"="+reexec.CLOEXEC)
		if err := cmd.fd.SetCloseExec(false); err != nil {
			return err
		}
	}

	if cmd.Verbose {
		env = append(env, reexec.EnvVerbose+"=1")
	}

	cmd.Env = append(cmd.Env, env...)

	if cmd.Name == "" {
		cmd.Args[0] = path.Base(os.Args[0])
	} else {
		cmd.Args[0] = cmd.Name
	}

	return nil
}
