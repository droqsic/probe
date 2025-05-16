package coverage

import (
	"os"
	"runtime"
	"testing"

	"github.com/droqsic/probe"
	"github.com/droqsic/probe/platform"
)

// TestPlatformSpecificCoverage tests platform-specific code paths to ensure maximum code coverage.
// This test covers various scenarios to ensure that all platform-specific code paths are executed.
// It does not check the correctness of the results, only that all code paths are executed.
func TestPlatformSpecificCoverage(t *testing.T) {
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

			probeResult := probe.IsTerminal(fd.fd)
			platformResult := platform.IsTerminal(fd.fd)

			if probeResult != platformResult {
				t.Errorf("Inconsistent results for %s: probe=%v, platform=%v",
					fd.name, probeResult, platformResult)
			}

			t.Logf("%s is terminal: %v", fd.name, probeResult)
		})
	}

	if runtime.GOOS == "windows" {
		t.Run("cygwin-detection", func(t *testing.T) {
			for _, fd := range fds {
				probe.ClearCache()
				isCygwin := probe.IsCygwinTerminal(fd.fd)
				t.Logf("%s is Cygwin terminal: %v", fd.name, isCygwin)
			}
		})
	}

	t.Run("different-file-types", func(t *testing.T) {
		f, err := os.CreateTemp("", "probe-coverage")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		defer f.Close()

		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("Failed to create pipe: %v", err)
		}
		defer r.Close()
		defer w.Close()

		fileTypes := []struct {
			name string
			fd   uintptr
		}{
			{"regular-file", f.Fd()},
			{"pipe-reader", r.Fd()},
			{"pipe-writer", w.Fd()},
		}

		for _, ft := range fileTypes {
			probe.ClearCache()

			probeResult := probe.IsTerminal(ft.fd)
			platformResult := platform.IsTerminal(ft.fd)

			if probeResult != platformResult {
				t.Errorf("Inconsistent results for %s: probe=%v, platform=%v",
					ft.name, probeResult, platformResult)
			}

			if probeResult {
				t.Errorf("%s should not be detected as terminal", ft.name)
			}

			t.Logf("%s is terminal: %v", ft.name, probeResult)
		}
	})

	t.Run("invalid-file-descriptors", func(t *testing.T) {
		invalidFDs := []uintptr{
			uintptr(999999),
			uintptr(1 << 30),
			uintptr(^uint(0) - 1),
		}

		for _, fd := range invalidFDs {
			probe.ClearCache()

			probeResult := probe.IsTerminal(fd)
			platformResult := platform.IsTerminal(fd)

			if probeResult {
				t.Errorf("Invalid file descriptor should not be detected as terminal")
			}

			if probeResult != platformResult {
				t.Errorf("Inconsistent results for invalid fd: probe=%v, platform=%v",
					probeResult, platformResult)
			}

			if runtime.GOOS == "windows" {
				cygwinResult := probe.IsCygwinTerminal(fd)
				if cygwinResult {
					t.Errorf("Invalid file descriptor should not be detected as Cygwin terminal")
				}
			}
		}
	})

	t.Run("cache-behavior", func(t *testing.T) {
		fd := os.Stdout.Fd()

		probe.ClearCache()

		result1 := probe.IsTerminal(fd)
		result2 := probe.IsTerminal(fd)

		if result1 != result2 {
			t.Errorf("Inconsistent results between uncached and cached calls")
		}

		probe.ClearCache()
		result3 := probe.IsTerminal(fd)

		if result1 != result3 {
			t.Errorf("Inconsistent results after clearing cache")
		}
	})
}
