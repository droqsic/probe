//go:build plan9
// +build plan9

package platform

import (
	"syscall"
)

// isTerminal returns true if the given file descriptor is a terminal on Plan9.
// In Plan9, terminals are represented by specific device paths.
func isTerminal(fd uintptr) bool {
	path, err := syscall.Fd2path(int(fd))
	if err != nil {
		return false
	}
	return path == "/dev/cons" || path == "/mnt/term/dev/cons"
}

// isCygwin always returns false on Plan9.
// Cygwin terminals are only relevant on Windows.
func isCygwin(fd uintptr) bool {
	return false
}
