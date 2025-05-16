package unit

import (
	"os"
	"testing"

	"github.com/droqsic/probe"
	"github.com/droqsic/probe/platform"
)

// TestPlatformConsistency ensures that platform functions are consistent with main API.
// This test checks that the platform-specific functions return the same results as the main API.
// It does not check the correctness of the results, only that they are consistent.
func TestPlatformConsistency(t *testing.T) {
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

			platformResult := platform.IsTerminal(fd.fd)
			probeResult := probe.IsTerminal(fd.fd)

			if platformResult != probeResult {
				t.Errorf("Inconsistent results between platform.IsTerminal and probe.IsTerminal for %s", fd.name)
			}

			platformCygwin := platform.IsCygwin(fd.fd)
			probeCygwin := probe.IsCygwinTerminal(fd.fd)

			if platformCygwin != probeCygwin {
				t.Errorf("Inconsistent results between platform.IsCygwin and probe.IsCygwinTerminal for %s", fd.name)
			}
		})
	}
}
