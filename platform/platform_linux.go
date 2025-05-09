//go:build linux
// +build linux

package platform

import (
	"golang.org/x/sys/unix"
)

// isTerminal returns true if the given file descriptor is a terminal on Linux.
// It uses the TCGETS ioctl call which is specific to Linux.
func isTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	return err == nil
}

// isCygwin always returns false on Linux.
// Cygwin terminals are only relevant on Windows.
func isCygwin(fd uintptr) bool {
	return false
}
