package tests

import (
	"os"
	"runtime"
	"testing"

	"github.com/droqsic/probe/platform"
)

// TestPlatformIsTerminal tests the platform-specific IsTerminal function
func TestPlatformIsTerminal(t *testing.T) {
	fd := os.Stdout.Fd()

	result := platform.IsTerminal(fd)
	t.Logf("platform.IsTerminal(%d) = %v on %s", fd, result, runtime.GOOS)
}

// TestPlatformIsCygwin tests the platform-specific IsCygwin function
func TestPlatformIsCygwin(t *testing.T) {
	// This function is only meaningful on Windows
	if runtime.GOOS != "windows" {
		t.Skip("IsCygwin is only relevant on Windows")
	}

	fd := os.Stdout.Fd()

	result := platform.IsCygwin(fd)
	t.Logf("platform.IsCygwin(%d) = %v", fd, result)
}

// TestPlatformConsistency ensures that platform functions are consistent with main API
func TestPlatformConsistency(t *testing.T) {
	fds := []struct {
		name string
		fd   uintptr
	}{
		{"stdin", os.Stdin.Fd()},
		{"stdout", os.Stdout.Fd()},
		{"stderr", os.Stderr.Fd()},
	}

	for _, f := range fds {
		// Import the main package to access its functions
		mainResult := platform.IsTerminal(f.fd)

		// We can't directly compare with probe.IsTerminal here because that would
		// create an import cycle. Instead, we just log the results.
		t.Logf("%s (fd %d): platform.IsTerminal = %v", f.name, f.fd, mainResult)
	}
}

// TestPlatformEdgeCases tests edge cases for platform-specific functions
func TestPlatformEdgeCases(t *testing.T) {
	// Test with an invalid file descriptor
	invalidFd := uintptr(999999)

	// Should return false and not panic
	result := platform.IsTerminal(invalidFd)
	if result {
		t.Errorf("Expected platform.IsTerminal to return false for invalid fd, got true")
	}

	// On Windows, also test IsCygwin
	if runtime.GOOS == "windows" {
		result = platform.IsCygwin(invalidFd)
		if result {
			t.Errorf("Expected platform.IsCygwin to return false for invalid fd, got true")
		}
	}
}
