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
	testStandardFileDescriptors(t)
	testPlatformSpecificPaths(t)
}

// testStandardFileDescriptors tests standard file descriptors (stdin, stdout, stderr)
func testStandardFileDescriptors(t *testing.T) {
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
}

// testPlatformSpecificPaths tests platform-specific code paths
func testPlatformSpecificPaths(t *testing.T) {
	switch runtime.GOOS {
	case "windows":
		testWindowsPaths(t)
	case "plan9":
		testPlan9Paths(t)
	case "js":
		if runtime.GOARCH == "wasm" {
			testWasmPaths(t)
		}
	}
}

// testWindowsPaths tests Windows-specific code paths
func testWindowsPaths(t *testing.T) {
	t.Run("cygwin-detection", func(t *testing.T) {
		fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
		for _, fd := range fds {
			probe.ClearCache()
			isCygwin := probe.IsCygwinTerminal(fd)
			t.Logf("fd %d is Cygwin terminal: %v", fd, isCygwin)
		}
	})
}

// testPlan9Paths tests Plan9-specific code paths
func testPlan9Paths(t *testing.T) {
	t.Run("plan9-paths", func(t *testing.T) {
		fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
		for _, fd := range fds {
			probe.ClearCache()
			probe.IsTerminal(fd)
		}
	})
}

// testWasmPaths tests WebAssembly-specific code paths
func testWasmPaths(t *testing.T) {
	t.Run("wasm-paths", func(t *testing.T) {
		fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
		for _, fd := range fds {
			probe.ClearCache()
			probe.IsTerminal(fd)
		}
	})
}
