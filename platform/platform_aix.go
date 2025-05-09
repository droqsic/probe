//go:build aix
// +build aix

package platform

import (
	"golang.org/x/sys/unix"
)

// isTerminal returns true if the given file descriptor is a terminal on AIX.
// It uses the TCGETA ioctl call which is specific to AIX.
func isTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TCGETA)
	return err == nil
}

// isCygwin always returns false on AIX.
// Cygwin terminals are only relevant on Windows.
func isCygwin(fd uintptr) bool {
	return false
}
