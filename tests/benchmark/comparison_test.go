package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/probe"
	"github.com/mattn/go-isatty"
)

// BenchmarkProbeIsTerminal measures the performance of probe.IsTerminal
func BenchmarkProbeIsTerminal(b *testing.B) {
	// Clear cache before benchmarking
	probe.ClearCache()

	fd := os.Stdout.Fd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		probe.IsTerminal(fd)
	}
}

// BenchmarkIsattyIsTerminal measures the performance of isatty.IsTerminal
func BenchmarkIsattyIsTerminal(b *testing.B) {
	fd := os.Stdout.Fd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		isatty.IsTerminal(fd)
	}
}

// BenchmarkProbeCachedIsTerminal measures the performance of probe.IsTerminal with caching
func BenchmarkProbeCachedIsTerminal(b *testing.B) {
	fd := os.Stdout.Fd()

	// Warm up the cache with one call
	probe.IsTerminal(fd)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		probe.IsTerminal(fd)
	}
}

// BenchmarkIsattyCachedIsTerminal measures the performance of isatty.IsTerminal with repeated calls
// Note: isatty doesn't have built-in caching, so this shows the difference in approach
func BenchmarkIsattyCachedIsTerminal(b *testing.B) {
	fd := os.Stdout.Fd()

	// One call to simulate the same setup as the probe test
	isatty.IsTerminal(fd)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		isatty.IsTerminal(fd)
	}
}

// BenchmarkProbeIsCygwinTerminal measures the performance of probe.IsCygwinTerminal
func BenchmarkProbeIsCygwinTerminal(b *testing.B) {
	// Clear cache before benchmarking
	probe.ClearCache()

	fd := os.Stdout.Fd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		probe.IsCygwinTerminal(fd)
	}
}

// BenchmarkIsattyIsCygwinTerminal measures the performance of isatty.IsCygwinTerminal
func BenchmarkIsattyIsCygwinTerminal(b *testing.B) {
	fd := os.Stdout.Fd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		isatty.IsCygwinTerminal(fd)
	}
}

// TestCompareResults verifies that both libraries return the same results
func TestCompareResults(t *testing.T) {
	// Clear cache before testing
	probe.ClearCache()

	fds := []struct {
		name string
		fd   uintptr
	}{
		{"stdin", os.Stdin.Fd()},
		{"stdout", os.Stdout.Fd()},
		{"stderr", os.Stderr.Fd()},
	}

	for _, f := range fds {
		probeResult := probe.IsTerminal(f.fd)
		isattyResult := isatty.IsTerminal(f.fd)

		if probeResult != isattyResult {
			t.Errorf("%s: probe.IsTerminal returned %v, but isatty.IsTerminal returned %v",
				f.name, probeResult, isattyResult)
		}

		// Only compare Cygwin detection on Windows
		if os.Getenv("GOOS") == "windows" {
			probeCygwin := probe.IsCygwinTerminal(f.fd)
			isattyCygwin := isatty.IsCygwinTerminal(f.fd)

			if probeCygwin != isattyCygwin {
				t.Errorf("%s: probe.IsCygwinTerminal returned %v, but isatty.IsCygwinTerminal returned %v",
					f.name, probeCygwin, isattyCygwin)
			}
		}
	}
}
