// Package exec runs a command stored as bytes in memory.
package exec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"codeberg.org/msantos/embedexe"
)

const (
	EnvVar     = "EMBEDEXE"
	EnvFlags   = "EMBEDEXE_FLAGS"
	EnvVerbose = "EMBEDEXE_VERBOSE" // enable debug error messages
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

func errexit(status int, err error) {
	if os.Getenv(EnvVerbose) != "" {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
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

	err = embedexe.Exec(uintptr(fd), os.Args, os.Environ())

	errexit(126, err)
}

// Command returns the Cmd struct to execute the program held in exe
// with the given arguments.
func Command(exe []byte, argv []string) *Cmd {
	cmd := exec.Command("/proc/self/exe", argv...)
	return &Cmd{
		Cmd: cmd,
		Exe: exe,
	}
}

// CommandContext returns a Cmd struct using the provided context.
func CommandContext(ctx context.Context, exe []byte, argv []string) *Cmd {
	cmd := exec.CommandContext(ctx, "/proc/self/exe", argv...)
	return &Cmd{
		Cmd: cmd,
		Exe: exe,
	}
}

// Run starts the specified command and waits for it to complete.
func (cmd *Cmd) Run() error {
	if err := cmd.fdopen(); err != nil {
		return err
	}
	return cmd.Cmd.Run()
}

// Start starts the specified command but does not wait for it to complete.
func (cmd *Cmd) Start() error {
	if err := cmd.fdopen(); err != nil {
		return err
	}
	return cmd.Cmd.Start()
}

// CombinedOutput runs the command and returns its combined standard
// output and standard error.
func (cmd *Cmd) CombinedOutput() ([]byte, error) {
	if err := cmd.fdopen(); err != nil {
		return nil, err
	}
	return cmd.Cmd.CombinedOutput()
}

// Output runs the command and returns its standard output.
func (cmd *Cmd) Output() ([]byte, error) {
	if err := cmd.fdopen(); err != nil {
		return nil, err
	}
	return cmd.Cmd.Output()
}

func (cmd *Cmd) fdopen() error {
	if cmd.Name == "" {
		cmd.Args[0] = os.Args[0]
	} else {
		cmd.Args[0] = cmd.Name
	}

	fd, err := embedexe.Open(cmd.Exe, cmd.Args[0])
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
