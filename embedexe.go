// Package embedexe executes a program stored in a byte array such as
// an executable or a directory of executables embedded in a Go binary.
package embedexe

import (
	"fmt"
	"os"

	"codeberg.org/msantos/execve"

	"golang.org/x/sys/unix"
)

type FD uintptr

func write(fd int, p []byte) error {
	for i := 0; i < len(p); {
		n, err := unix.Write(fd, p[i:])
		if err != nil {
			return err
		}
		// check if Write will return io.EOF
		if n <= 0 {
			return nil
		}
		i += n
	}
	return nil
}

// Open returns a file descriptor to an executable stored in memory.
func Open(exe []byte, arg0 string) (FD, error) {
	flag := unix.MFD_CLOEXEC

	if len(exe) > 1 && exe[0] == '#' && exe[1] == '!' {
		flag &= ^unix.MFD_CLOEXEC
	}

	fd, err := unix.MemfdCreate(arg0, flag)
	if err != nil {
		return FD(0), err
	}

	if err := write(fd, exe); err != nil {
		return FD(0), err
	}

	return FD(fd), nil
}

// Close closes the executable file descriptor.
func (fd FD) Close() error {
	return unix.Close(int(fd))
}

// Path returns the path to the executable file descriptor. Running the
// executable using the file descriptor path directly is an alternative to
// running by file descriptor in Exec.
func (fd FD) Path() string {
	return fmt.Sprintf("/proc/%d/fd/%d", os.Getpid(), int(fd))
}

// Exec runs the executable referenced by the file descriptor, replacing
// the current running process image.
func (fd FD) Exec(argv, env []string) error {
	return execve.Fexecve(uintptr(fd), argv, env)
}

// CloseExec checks if the O_CLOEXEC flag is set on the file descriptor.
func (fd FD) CloseExec() bool {
	flag, err := unix.FcntlInt(uintptr(fd), unix.F_GETFD, 0)
	if err != nil {
		return false
	}

	return flag&unix.MFD_CLOEXEC != 0
}

// SetCloseExec enables or disables the O_CLOEXEC flag on the file
// descriptor.
func (fd FD) SetCloseExec(b bool) error {
	flag, err := unix.FcntlInt(uintptr(fd), unix.F_GETFD, 0)
	if err != nil {
		return err
	}

	if b {
		flag |= unix.MFD_CLOEXEC
	} else {
		flag &= ^unix.MFD_CLOEXEC
	}

	_, err = unix.FcntlInt(uintptr(fd), unix.F_SETFD, flag)
	return err
}
