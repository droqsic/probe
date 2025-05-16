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

	t.Run("unusual-file-types", func(t *testing.T) {
		f, err := os.CreateTemp("", "probe-edge-case")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		defer f.Close()

		var namedPipeFd uintptr
		if runtime.GOOS != "windows" {
			pipeName := "/tmp/probe-test-pipe"
			err := os.Remove(pipeName)
			if err != nil && !os.IsNotExist(err) {
				t.Logf("Failed to remove existing named pipe: %v", err)
			}

			r, w, err := os.Pipe()
			if err == nil {
				defer r.Close()
				defer w.Close()
				namedPipeFd = r.Fd()
			} else {
				t.Logf("Failed to create pipe: %v", err)
			}
		}

		probe.ClearCache()
		if probe.IsTerminal(f.Fd()) {
			t.Errorf("Regular file should not be detected as terminal")
		}

		if namedPipeFd != 0 {
			probe.ClearCache()
			if probe.IsTerminal(namedPipeFd) {
				t.Errorf("Pipe should not be detected as terminal")
			}
		}
	})

	t.Run("platform-specific", func(t *testing.T) {
		switch runtime.GOOS {
		case "windows":
			t.Run("windows", func(t *testing.T) {
				stdHandles := []uintptr{0, 1, 2} // STD_INPUT_HANDLE, STD_OUTPUT_HANDLE, STD_ERROR_HANDLE

				for _, handle := range stdHandles {
					probe.ClearCache()
					probe.IsTerminal(handle)
					probe.IsCygwinTerminal(handle)
				}
			})

		case "plan9":
			t.Run("plan9", func(t *testing.T) {
				probe.ClearCache()
				probe.IsTerminal(0)
				probe.IsTerminal(1)
				probe.IsTerminal(2)
			})

		case "js":
			if runtime.GOARCH == "wasm" {
				t.Run("wasm", func(t *testing.T) {
					probe.ClearCache()
					probe.IsTerminal(0)
					probe.IsTerminal(1)
					probe.IsTerminal(2)
				})
			}
		}
	})

	t.Run("cache-edge-cases", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			probe.IsTerminal(uintptr(i + 10000))
		}

		probe.ClearCache()

		fd := os.Stdout.Fd()
		firstResult := probe.IsTerminal(fd)
		secondResult := probe.IsTerminal(fd)

		if firstResult != secondResult {
			t.Errorf("Cache inconsistency after clearing and refilling")
		}
	})
}
