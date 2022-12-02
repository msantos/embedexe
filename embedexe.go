// Package embedexe embeds an executable or a directory of executables
// in a Go binary and runs them from memory.
package embedexe

import (
	"codeberg.org/msantos/execve"
	"golang.org/x/sys/unix"
)

// Exec runs the executable referenced by the file descriptor, replacing
// the current running process image.
func Exec(fd uintptr, argv, env []string) error {
	return execve.Fexecve(fd, argv, env)
}

// CloseExec checks if the O_CLOEXEC flag is set on the file descriptor.
func CloseExec(fd uintptr) bool {
	flag, err := unix.FcntlInt(fd, unix.F_GETFD, 0)
	if err != nil {
		return false
	}

	return flag&unix.MFD_CLOEXEC != 0
}

// SetCloseExec enables or disables the O_CLOEXEC flag on the file
// descriptor.
func SetCloseExec(fd uintptr, b bool) error {
	flag, err := unix.FcntlInt(fd, unix.F_GETFD, 0)
	if err != nil {
		return err
	}

	if b {
		flag |= unix.MFD_CLOEXEC
	} else {
		flag &= ^unix.MFD_CLOEXEC
	}

	_, err = unix.FcntlInt(fd, unix.F_SETFD, flag)
	return err
}

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
func Open(exe []byte, arg0 string) (uintptr, error) {
	flag := unix.MFD_CLOEXEC

	if len(exe) > 1 && exe[0] == '#' && exe[1] == '!' {
		flag &= ^unix.MFD_CLOEXEC
	}

	fd, err := unix.MemfdCreate(arg0, flag)
	if err != nil {
		return 0, err
	}

	if err := write(fd, exe); err != nil {
		return 0, err
	}

	return uintptr(fd), nil
}
