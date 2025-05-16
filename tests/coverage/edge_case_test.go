package coverage

import (
	"os"
	"runtime"
	"testing"

	"github.com/droqsic/probe"
)

// TestEdgeCases tests various edge cases to ensure maximum code coverage, and to catch any potential issues.
// This test is intended to cover edge cases that might not be covered by other tests.
func TestEdgeCases(t *testing.T) {
	testInvalidFileDescriptors(t)
	testPipeFileDescriptors(t)
	testPlatformSpecificEdgeCases(t)
}

// testInvalidFileDescriptors tests behavior with invalid file descriptors
func testInvalidFileDescriptors(t *testing.T) {
	t.Run("invalid-file-descriptors", func(t *testing.T) {
		invalidFDs := []uintptr{
			uintptr(999999),
			uintptr(1 << 30),
			uintptr(^uint(0) - 1),
		}

		for _, fd := range invalidFDs {
			probe.ClearCache()

			result := probe.IsTerminal(fd)
			if result {
				t.Errorf("Invalid file descriptor should not be detected as terminal")
			}

			if runtime.GOOS == "windows" {
				cygwinResult := probe.IsCygwinTerminal(fd)
				if cygwinResult {
					t.Errorf("Invalid file descriptor should not be detected as Cygwin terminal")
				}
			}
		}
	})
}

// testPipeFileDescriptors tests behavior with pipe file descriptors
func testPipeFileDescriptors(t *testing.T) {
	t.Run("pipe-file-descriptors", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("Failed to create pipe: %v", err)
		}
		defer r.Close()
		defer w.Close()

		probe.ClearCache()
		if probe.IsTerminal(r.Fd()) {
			t.Errorf("Pipe read end should not be detected as terminal")
		}

		probe.ClearCache()
		if probe.IsTerminal(w.Fd()) {
			t.Errorf("Pipe write end should not be detected as terminal")
		}

		// Test named pipe if supported on this platform
		if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
			testNamedPipe(t)
		}
	})
}

// testNamedPipe tests behavior with named pipes
func testNamedPipe(t *testing.T) {
	// Implementation depends on platform
	// This is a placeholder for the actual implementation
	t.Log("Named pipe tests would go here")
}

// testPlatformSpecificEdgeCases tests platform-specific edge cases
func testPlatformSpecificEdgeCases(t *testing.T) {
	t.Run("platform-specific", func(t *testing.T) {
		switch runtime.GOOS {
		case "windows":
			testWindowsEdgeCases(t)
		case "plan9":
			testPlan9EdgeCases(t)
		case "js":
			if runtime.GOARCH == "wasm" {
				testWasmEdgeCases(t)
			}
		}
	})
}

// testWindowsEdgeCases tests Windows-specific edge cases
func testWindowsEdgeCases(t *testing.T) {
	stdHandles := []uintptr{0, 1, 2} // STD_INPUT_HANDLE, STD_OUTPUT_HANDLE, STD_ERROR_HANDLE

	for _, handle := range stdHandles {
		probe.ClearCache()
		probe.IsTerminal(handle)
		probe.IsCygwinTerminal(handle)
	}
}

// testPlan9EdgeCases tests Plan9-specific edge cases
func testPlan9EdgeCases(t *testing.T) {
	stdFDs := []uintptr{0, 1, 2}
	for _, fd := range stdFDs {
		probe.ClearCache()
		probe.IsTerminal(fd)
	}
}

// testWasmEdgeCases tests WebAssembly-specific edge cases
func testWasmEdgeCases(t *testing.T) {
	stdFDs := []uintptr{0, 1, 2}
	for _, fd := range stdFDs {
		probe.ClearCache()
		probe.IsTerminal(fd)
	}
}
