//go:build !windows && !linux && !android && !darwin && !freebsd && !openbsd && !netbsd && !dragonfly && !hurd && !solaris && !aix && !illumos && !zos && !ios && !plan9 && !js && !tinygo && !haikou && !appengine && !nacl
// +build !windows,!linux,!android,!darwin,!freebsd,!openbsd,!netbsd,!dragonfly,!hurd,!solaris,!aix,!illumos,!zos,!ios,!plan9,!js,!tinygo,!haikou,!appengine,!nacl

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
