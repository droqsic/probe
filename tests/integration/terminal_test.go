package integration

import (
	"os"
	"runtime"
	"testing"

	"github.com/droqsic/probe"
	"github.com/droqsic/probe/platform"
	"github.com/mattn/go-isatty"
)

// TestTerminalDevices specifically tests terminal device detection
func TestTerminalDevices(t *testing.T) {
	// Skip on Windows as the concept of /dev/tty doesn't exist
	if runtime.GOOS == "windows" {
		t.Skip("Terminal device tests only run on Unix-like systems")
	}

	// Clear the cache to ensure a fresh test
	probe.ClearCache()

	// Try to open /dev/tty which should be a terminal device
	tty, err := os.Open("/dev/tty")
	if err != nil {
		t.Logf("Could not open /dev/tty: %v", err)
		t.Skip("Test requires access to /dev/tty")
	}
	defer tty.Close()

	fd := tty.Fd()

	// Test with all available methods
	isattyResult := isatty.IsTerminal(fd)
	probeResult := probe.IsTerminal(fd)
	platformResult := platform.IsTerminal(fd)

	t.Logf("/dev/tty detection results:")
	t.Logf("  isatty.IsTerminal: %v", isattyResult)
	t.Logf("  probe.IsTerminal: %v", probeResult)
	t.Logf("  platform.IsTerminal: %v", platformResult)

	// Both libraries should identify /dev/tty as a terminal
	if !isattyResult {
		t.Errorf("isatty.IsTerminal failed to identify /dev/tty as a terminal")
	}

	if !probeResult {
		t.Errorf("probe.IsTerminal failed to identify /dev/tty as a terminal")
	}

	if !platformResult {
		t.Errorf("platform.IsTerminal failed to identify /dev/tty as a terminal")
	}
}

// TestPtsDevices tests pseudo-terminal slave devices
func TestPtsDevices(t *testing.T) {
	// Skip on Windows as the concept of pts doesn't exist
	if runtime.GOOS == "windows" {
		t.Skip("PTS device tests only run on Unix-like systems")
	}

	// Clear the cache to ensure a fresh test
	probe.ClearCache()

	// Try common pts devices
	ptsDevices := []string{
		"/dev/pts/0",
		"/dev/pts/1",
		"/dev/pts/2",
	}

	for _, device := range ptsDevices {
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

		t.Logf("%s detection results:", device)
		t.Logf("  isatty.IsTerminal: %v", isattyResult)
		t.Logf("  probe.IsTerminal: %v", probeResult)
		t.Logf("  platform.IsTerminal: %v", platformResult)

		// Both libraries should have consistent results
		if isattyResult != probeResult {
			t.Errorf("Inconsistent results for %s: isatty=%v, probe=%v",
				device, isattyResult, probeResult)
		}
	}
}

// TestCachingBehavior tests if caching might be causing inconsistencies
func TestCachingBehavior(t *testing.T) {
	// Skip on Windows as the concept of /dev/tty doesn't exist
	if runtime.GOOS == "windows" {
		t.Skip("Terminal device tests only run on Unix-like systems")
	}

	// Clear the cache to ensure a fresh test
	probe.ClearCache()

	// Try to open /dev/tty which should be a terminal device
	tty, err := os.Open("/dev/tty")
	if err != nil {
		t.Logf("Could not open /dev/tty: %v", err)
		t.Skip("Test requires access to /dev/tty")
	}
	defer tty.Close()

	fd := tty.Fd()

	// First call to probe.IsTerminal
	firstResult := probe.IsTerminal(fd)

	// Direct call to platform.IsTerminal
	platformResult := platform.IsTerminal(fd)

	// Second call to probe.IsTerminal
	secondResult := probe.IsTerminal(fd)

	t.Logf("Caching behavior for /dev/tty (fd %d):", fd)
	t.Logf("  First probe.IsTerminal call: %v", firstResult)
	t.Logf("  platform.IsTerminal call: %v", platformResult)
	t.Logf("  Second probe.IsTerminal call: %v", secondResult)

	if firstResult != platformResult {
		t.Errorf("First probe.IsTerminal call (%v) doesn't match platform.IsTerminal (%v)",
			firstResult, platformResult)
	}

	if firstResult != secondResult {
		t.Errorf("First probe.IsTerminal call (%v) doesn't match second call (%v)",
			firstResult, secondResult)
		t.Log("This suggests the caching mechanism might be inconsistent")
	}
}
