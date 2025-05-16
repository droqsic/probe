package unit

import (
	"os"
	"runtime"
	"testing"

	"github.com/droqsic/probe"
)

// TestIsTerminalBasic tests the basic functionality of IsTerminal.
// This test checks that the cache is consistent for standard file descriptors.
// It does not check the correctness of the terminal detection, only that the cache is consistent.
func TestIsTerminalBasic(t *testing.T) {
	fds := []struct {
		name string
		fd   uintptr
	}{
		{"stdin", os.Stdin.Fd()},
		{"stdout", os.Stdout.Fd()},
		{"stderr", os.Stderr.Fd()},
	}

	for _, fd := range fds {
		t.Run(fd.name, func(t *testing.T) {
			probe.ClearCache()

			result := probe.IsTerminal(fd.fd)
			t.Logf("%s is terminal: %v", fd.name, result)

			cachedResult := probe.IsTerminal(fd.fd)
			if result != cachedResult {
				t.Errorf("Inconsistent results between initial and cached checks for %s", fd.name)
			}
		})
	}
}

// TestIsCygwinTerminalBasic tests the basic functionality of IsCygwinTerminal.
// This test only runs on Windows.
// It checks that the cache is consistent for standard file descriptors.
func TestIsCygwinTerminalBasic(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("IsCygwinTerminal is only relevant on Windows")
	}

	fds := []struct {
		name string
		fd   uintptr
	}{
		{"stdin", os.Stdin.Fd()},
		{"stdout", os.Stdout.Fd()},
		{"stderr", os.Stderr.Fd()},
	}

	for _, fd := range fds {
		t.Run(fd.name, func(t *testing.T) {
			probe.ClearCache()

			result := probe.IsCygwinTerminal(fd.fd)
			t.Logf("%s is Cygwin terminal: %v", fd.name, result)

			cachedResult := probe.IsCygwinTerminal(fd.fd)
			if result != cachedResult {
				t.Errorf("Inconsistent results between initial and cached Cygwin checks for %s", fd.name)
			}
		})
	}
}

// TestNonTerminalFiles tests that non-terminal files are correctly identified.
// This test creates a temporary file and checks that it is not detected as terminal.
// It also checks that non-terminal files are not detected as Cygwin terminals.
func TestNonTerminalFiles(t *testing.T) {
	f, err := os.CreateTemp("", "probe-test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	if probe.IsTerminal(f.Fd()) {
		t.Errorf("Regular file should not be detected as terminal")
	}

	if probe.IsCygwinTerminal(f.Fd()) {
		t.Errorf("Regular file should not be detected as Cygwin terminal")
	}
}
