//go:build darwin || freebsd || openbsd || netbsd || dragonfly || hurd || zos
// +build darwin freebsd openbsd netbsd dragonfly hurd zos

package platform

import (
	"golang.org/x/sys/unix"
)

// isTerminal returns true if the given file descriptor is a terminal on BSD systems.
// It uses the TIOCGETA ioctl call which is common across BSD variants.
func isTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TIOCGETA)
	return err == nil
}

// isCygwin always returns false on BSD systems.
// Cygwin terminals are only relevant on Windows.
func isCygwin(fd uintptr) bool {
	return false
}
