//go:build zos
// +build zos

package platform

import (
	"golang.org/x/sys/unix"
)

// isTerminal returns true if the given file descriptor is a terminal on z/OS.
// It uses the TCGETS ioctl call which is appropriate for z/OS.
func isTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	return err == nil
}

// isCygwin always returns false on z/OS.
// Cygwin terminals are only relevant on Windows.
func isCygwin(fd uintptr) bool {
	return false
}
