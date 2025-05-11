package integration

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/droqsic/probe"
)

// TestIntegrationWithRedirection tests the behavior when stdout is redirected
func TestIntegrationWithRedirection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if we're not in a terminal to begin with
	if !probe.IsTerminal(os.Stdout.Fd()) {
		t.Skip("Test requires running in a terminal")
	}

	// Create a temporary file to redirect output to
	tmpfile, err := os.CreateTemp("", "probe-redirect-test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	// Run the current test binary with a specific test and redirect output
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	cmd.Stdout = tmpfile

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run helper process: %v", err)
	}

	// Read the output
	tmpfile.Seek(0, 0)
	output := make([]byte, 100)
	n, err := tmpfile.Read(output)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	// The helper process should report that stdout is not a terminal
	if !strings.Contains(string(output[:n]), "stdout: false") {
		t.Errorf("Expected redirected stdout to not be a terminal, got: %s", string(output[:n]))
	}
}

// TestHelperProcess is not a real test, it's used as a helper for TestIntegrationWithRedirection
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		t.Skip("Not a helper process")
	}

	// Check if stdout is a terminal and print the result
	isTerminal := probe.IsTerminal(os.Stdout.Fd())
	os.Stdout.WriteString(fmt.Sprintf("stdout: %v", isTerminal))

	// Exit immediately
	os.Exit(0)
}

// TestConcurrentAccess tests that the library is thread-safe
func TestConcurrentAccess(t *testing.T) {
	// Number of goroutines to spawn
	const numGoroutines = 100
	// Number of checks per goroutine
	const checksPerGoroutine = 1000

	// File descriptors to check
	fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}

	// Channel to collect errors
	errCh := make(chan error, numGoroutines)

	// Start goroutines
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < checksPerGoroutine; j++ {
				// Check each file descriptor
				for _, fd := range fds {
					// These calls should not panic
					probe.IsTerminal(fd)

					if runtime.GOOS == "windows" {
						probe.IsCygwinTerminal(fd)
					}
				}
			}
			errCh <- nil // Signal completion
		}(i)
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		if err := <-errCh; err != nil {
			t.Errorf("Goroutine reported error: %v", err)
		}
	}
}

// TestRedirection tests behavior when output is redirected
func TestRedirection(t *testing.T) {
	if os.Getenv("TEST_REDIRECTION") == "1" {
		// This code runs in the child process with redirected stdout
		fmt.Println(probe.IsTerminal(os.Stdout.Fd()))
		return
	}

	t.Log("Note: For a complete test, run with output redirected to a file")
	t.Log("Example: go test -run=TestRedirection ./tests > output.txt")
}

// TestFileModes tests with different file modes
func TestFileModes(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "probe-test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	// Open the same file in read-only mode
	roFile, err := os.OpenFile(tmpfile.Name(), os.O_RDONLY, 0)
	if err != nil {
		t.Fatalf("Failed to open file in read-only mode: %v", err)
	}
	defer roFile.Close()

	result := probe.IsTerminal(roFile.Fd())
	t.Logf("Read-only file: IsTerminal = %v", result)

	// A regular file should not be a terminal regardless of mode
	if result {
		t.Errorf("Expected IsTerminal to return false for a read-only file, got true")
	}
}
