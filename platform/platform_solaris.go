//go:build solaris || illumos
// +build solaris illumos

package platform

import (
	"golang.org/x/sys/unix"
)

// isTerminal returns true if the given file descriptor is a terminal on Solaris or Illumos.
// It uses the TCGETA ioctl call which is specific to Solaris or Illumos.
func isTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermio(int(fd), unix.TCGETA)
	return err == nil
}

// isCygwin always returns false on Solaris or Illumos.
// Cygwin terminals are only relevant on Windows.
func isCygwin(fd uintptr) bool {
	return false
}
