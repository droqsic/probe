//go:build !windows && !linux && !darwin && !freebsd && !openbsd && !netbsd && !dragonfly && !solaris && !aix && !illumos && !zos
// +build !windows,!linux,!darwin,!freebsd,!openbsd,!netbsd,!dragonfly,!solaris,!aix,!illumos,!zos

package platform

// isTerminal is a stub implementation for unsupported platforms.
// It always returns false.
func isTerminal(fd uintptr) bool {
	return false
}

// isCygwin is a stub implementation for unsupported platforms.
// It always returns false.
func isCygwin(fd uintptr) bool {
	return false
}
