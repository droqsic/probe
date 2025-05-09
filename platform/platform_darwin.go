//go:build darwin
// +build darwin

package platform

import (
	"golang.org/x/sys/unix"
)

// isTerminal returns true if the given file descriptor is a terminal on Darwin (macOS).
// It uses the TIOCGETA ioctl call which is specific to Darwin systems.
func isTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TIOCGETA)
	return err == nil
}

// isCygwin always returns false on Darwin.
// Cygwin terminals are only relevant on Windows.
func isCygwin(fd uintptr) bool {
	return false
}
