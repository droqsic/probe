package integration

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/droqsic/probe"
	"github.com/mattn/go-isatty"
)

// TestIsTerminal verifies that IsTerminal correctly identifies terminals
func TestIsTerminal(t *testing.T) {
	// Standard file descriptors to test
	fds := []struct {
		name string
		fd   uintptr
	}{
		{"stdin", os.Stdin.Fd()},
		{"stdout", os.Stdout.Fd()},
		{"stderr", os.Stderr.Fd()},
	}

	for _, f := range fds {
		result := probe.IsTerminal(f.fd)
		t.Logf("%s (fd %d): IsTerminal = %v", f.name, f.fd, result)

		// We can't assert the exact result as it depends on how the test is run
		// (terminal vs. CI environment), but we can ensure the function runs without errors
	}
}

// TestIsCygwinTerminal verifies that IsCygwinTerminal works correctly
func TestIsCygwinTerminal(t *testing.T) {
	// Skip on non-Windows platforms
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

	for _, f := range fds {
		result := probe.IsCygwinTerminal(f.fd)
		t.Logf("%s (fd %d): IsCygwinTerminal = %v", f.name, f.fd, result)

		// We can't assert the exact result as it depends on the environment
	}
}

// TestCaching verifies that the caching mechanism works correctly
func TestCaching(t *testing.T) {
	fd := os.Stdout.Fd()

	// First call should determine the actual status
	firstResult := probe.IsTerminal(fd)

	// Second call should use the cached result
	secondResult := probe.IsTerminal(fd)

	// Results should be consistent
	if firstResult != secondResult {
		t.Errorf("Inconsistent results: first call returned %v, second call returned %v",
			firstResult, secondResult)
	}
}

// TestNonExistentFileDescriptor tests behavior with an invalid file descriptor
func TestNonExistentFileDescriptor(t *testing.T) {
	// Use a likely invalid file descriptor
	invalidFd := uintptr(999999)

	// Should return false and not panic
	result := probe.IsTerminal(invalidFd)
	if result {
		t.Errorf("Expected IsTerminal to return false for invalid fd, got true")
	}

	// On Windows, also test IsCygwinTerminal
	if runtime.GOOS == "windows" {
		result = probe.IsCygwinTerminal(invalidFd)
		if result {
			t.Errorf("Expected IsCygwinTerminal to return false for invalid fd, got true")
		}
	}
}

// TestFileDescriptorFromFile tests with a file descriptor from a regular file
func TestFileDescriptorFromFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "probe-test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	// Get the file descriptor
	fd := tmpfile.Fd()

	// A regular file should not be a terminal
	if probe.IsTerminal(fd) {
		t.Errorf("Expected IsTerminal to return false for a regular file, got true")
	}

	// On Windows, also test IsCygwinTerminal
	if runtime.GOOS == "windows" {
		if probe.IsCygwinTerminal(fd) {
			t.Errorf("Expected IsCygwinTerminal to return false for a regular file, got true")
		}
	}
}

// TestConsistency verifies that both libraries return the same results
func TestConsistency(t *testing.T) {
	fds := []struct {
		name string
		fd   uintptr
	}{
		{"Stdout", os.Stdout.Fd()},
		{"Stdin", os.Stdin.Fd()},
		{"Stderr", os.Stderr.Fd()},
	}

	for _, f := range fds {
		isattyResult := isatty.IsTerminal(f.fd)
		probeResult := probe.IsTerminal(f.fd)

		if isattyResult == probeResult {
			t.Logf("Consistent results for %s: both returned %v", f.name, isattyResult)
		} else {
			t.Errorf("Inconsistent results for %s: isatty=%v, probe=%v",
				f.name, isattyResult, probeResult)
		}
	}
}

// TestPipe tests with pipe file descriptors
func TestPipe(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	isattyR := isatty.IsTerminal(r.Fd())
	probeR := probe.IsTerminal(r.Fd())

	isattyW := isatty.IsTerminal(w.Fd())
	probeW := probe.IsTerminal(w.Fd())

	t.Logf("Pipe read end: isatty=%v, probe=%v", isattyR, probeR)
	t.Logf("Pipe write end: isatty=%v, probe=%v", isattyW, probeW)

	// Both should return false for pipes
	if isattyR || isattyW || probeR || probeW {
		t.Errorf("Unexpected true result for pipe test")
	}
}

// TestSubprocess tests behavior in a subprocess
func TestSubprocess(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		// This code runs in the subprocess
		fmt.Printf("isatty=%v,probe=%v\n",
			isatty.IsTerminal(os.Stdout.Fd()),
			probe.IsTerminal(os.Stdout.Fd()))
		return
	}

	// Run the test in a subprocess with various configurations
	cmd := exec.Command(os.Args[0], "-test.run=TestSubprocess")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=1")

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run subprocess: %v", err)
	}

	t.Logf("Subprocess output: %s", output)
}
