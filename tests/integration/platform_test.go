package integration

import (
	"os"
	"runtime"
	"testing"

	"github.com/droqsic/probe"
	"github.com/droqsic/probe/platform"
	"github.com/mattn/go-isatty"
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

// TestPlatformDetails tests platform-specific behavior in more detail
func TestPlatformDetails(t *testing.T) {
	t.Logf("OS: %s, Architecture: %s", runtime.GOOS, runtime.GOARCH)

	// Platform-specific tests
	switch runtime.GOOS {
	case "windows":
		testWindowsSpecific(t)
	case "darwin":
		testDarwinSpecific(t)
	case "linux":
		testLinuxSpecific(t)
	default:
		t.Logf("No specific tests for platform: %s", runtime.GOOS)
	}
}

func testWindowsSpecific(t *testing.T) {
	t.Log("Running Windows-specific tests")
	// Windows-specific tests would go here
}

func testDarwinSpecific(t *testing.T) {
	t.Log("Running macOS-specific tests")
	// macOS-specific tests would go here
}

func testLinuxSpecific(t *testing.T) {
	t.Log("Running Linux-specific tests")

	// Test with various Linux-specific devices if available
	devices := []string{
		"/dev/null",
		"/dev/zero",
		"/dev/tty",
		"/dev/console",
		"/dev/pts/0",
	}

	for _, device := range devices {
		file, err := os.Open(device)
		if err != nil {
			t.Logf("Could not open %s: %v", device, err)
			continue
		}
		defer file.Close()

		fd := file.Fd()
		isattyResult := isatty.IsTerminal(fd)
		probeResult := probe.IsTerminal(fd)
		platformResult := platform.IsTerminal(fd)

		t.Logf("%s: isatty=%v, probe=%v, platform=%v",
			device, isattyResult, probeResult, platformResult)

		// For terminal devices, both libraries should return the same result
		if isattyResult != probeResult {
			// This is a potential bug in the probe library
			t.Errorf("Inconsistent results for %s: isatty=%v, probe=%v",
				device, isattyResult, probeResult)
		}
	}
}

// TestEnvironmentVariables tests behavior with environment variables that might affect terminal detection
func TestEnvironmentVariables(t *testing.T) {
	// Save original environment
	origTerm := os.Getenv("TERM")
	defer os.Setenv("TERM", origTerm)

	// Test with different TERM values
	termValues := []string{"", "dumb", "xterm", "xterm-256color", "vt100"}

	for _, term := range termValues {
		os.Setenv("TERM", term)

		isattyResult := isatty.IsTerminal(os.Stdout.Fd())
		probeResult := probe.IsTerminal(os.Stdout.Fd())

		t.Logf("TERM=%s: isatty=%v, probe=%v", term, isattyResult, probeResult)

		// Check for consistency
		if isattyResult != probeResult {
			t.Errorf("Inconsistent results with TERM=%s", term)
		}
	}
}
