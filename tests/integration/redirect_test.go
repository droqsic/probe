package integration

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/droqsic/probe"
)

// TestRedirectedOutput tests the behavior of IsTerminal when output is redirected.
// This test uses a helper program to ensure that the behavior is consistent
// regardless of how the main test binary is built and run.
func TestRedirectedOutput(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping redirect test in CI environment")
	}

	helperCode := `
package main

import (
    "fmt"
    "os"

    "github.com/droqsic/probe"
)

func main() {
    fmt.Printf("%v", probe.IsTerminal(os.Stdout.Fd()))
}
`

	tempDir, err := os.MkdirTemp("", "probe-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	helperFile := filepath.Join(tempDir, "helper.go")
	if err := os.WriteFile(helperFile, []byte(helperCode), 0644); err != nil {
		t.Fatalf("Failed to write helper program: %v", err)
	}

	helperBin := filepath.Join(tempDir, "helper")
	if runtime.GOOS == "windows" {
		helperBin += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", helperBin, helperFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build helper program: %v", err)
	}

	var buf bytes.Buffer
	cmd = exec.Command(helperBin)
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run helper program: %v", err)
	}

	result := buf.String()
	if result != "false" {
		t.Errorf("Expected redirected stdout to not be a terminal, got %s", result)
	}
}

// TestPipeRedirection tests the behavior of IsTerminal with pipe redirection.
// This test uses a helper program to ensure that the behavior is consistent
// regardless of how the main test binary is built and run.
func TestPipeRedirection(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	os.Stdout = w
	result := probe.IsTerminal(os.Stdout.Fd())
	os.Stdout = oldStdout

	if result {
		t.Errorf("Expected piped stdout to not be a terminal")
	}
}

// TestFileRedirection tests the behavior of IsTerminal with file redirection.
// This test uses a helper program to ensure that the behavior is consistent
// regardless of how the main test binary is built and run.
func TestFileRedirection(t *testing.T) {
	f, err := os.CreateTemp("", "probe-test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	os.Stdout = f

	result := probe.IsTerminal(os.Stdout.Fd())
	os.Stdout = oldStdout

	if result {
		t.Errorf("Expected file-redirected stdout to not be a terminal")
	}
}

// TestMultipleRedirections tests multiple redirections in sequence.
// This test uses a helper program to ensure that the behavior is consistent
// regardless of how the main test binary is built and run.
func TestMultipleRedirections(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	f1, err := os.CreateTemp("", "probe-test-1")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := os.CreateTemp("", "probe-test-2")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	redirections := []struct {
		name string
		file *os.File
	}{
		{"file1", f1},
		{"pipe", w},
		{"file2", f2},
	}

	for _, redir := range redirections {
		os.Stdout = redir.file

		probe.ClearCache()

		result := probe.IsTerminal(os.Stdout.Fd())

		if result {
			t.Errorf("Expected %s-redirected stdout to not be a terminal", redir.name)
		}
	}
}

// TestRedirectionWithCaching tests that caching works correctly with redirections.
// This test uses a helper program to ensure that the behavior is consistent
// regardless of how the main test binary is built and run.
func TestRedirectionWithCaching(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	f, err := os.CreateTemp("", "probe-test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	originalResult := probe.IsTerminal(os.Stdout.Fd())
	t.Logf("Original stdout is terminal: %v", originalResult)

	os.Stdout = f
	redirectedResult1 := probe.IsTerminal(os.Stdout.Fd())

	probe.ClearCache()

	redirectedResult2 := probe.IsTerminal(os.Stdout.Fd())

	if redirectedResult2 {
		t.Errorf("Expected file-redirected stdout to not be a terminal after cache clear")
	}

	if redirectedResult1 == originalResult && redirectedResult1 != redirectedResult2 {
		t.Logf("Cache was used for redirected stdout (expected behavior)")
	}
}
