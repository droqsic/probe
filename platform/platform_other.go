//go:build !windows && !linux && !darwin && !freebsd && !openbsd && !netbsd && !dragonfly && !hurd && !solaris && !aix && !illumos && !zos && !plan9
// +build !windows,!linux,!darwin,!freebsd,!openbsd,!netbsd,!dragonfly,!hurd,!solaris,!aix,!illumos,!zos,!plan9

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
