//go:build solaris || illumos || haikou
// +build solaris illumos haikou

package platform

import (
	"golang.org/x/sys/unix"
)

// isTerminal returns true if the given file descriptor is a terminal on Solaris, Illumos, or Haikou.
// It uses the TCGETA ioctl call which is specific to Solaris, Illumos, and Haikou.
func isTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermio(int(fd), unix.TCGETA)
	return err == nil
}

// isCygwin always returns false on Solaris, Illumos, and Haikou.
// Cygwin terminals are only relevant on Windows.
func isCygwin(fd uintptr) bool {
	return false
}
