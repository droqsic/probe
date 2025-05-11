package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/probe"
)

// BenchmarkProbeAlternatingFDs measures performance when alternating between file descriptors
func BenchmarkProbeAlternatingFDs(b *testing.B) {
	fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		probe.IsTerminal(fds[i%len(fds)])
	}
}

// BenchmarkProbeCacheHitRate measures the impact of cache hit rate
func BenchmarkProbeCacheHitRate(b *testing.B) {
	// Create multiple file descriptors
	files := make([]*os.File, 0, 10)
	fds := make([]uintptr, 0, 10)

	// Create temporary files
	for i := 0; i < 10; i++ {
		f, err := os.CreateTemp("", "probe-benchmark")
		if err != nil {
			b.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		defer f.Close()

		files = append(files, f)
		fds = append(fds, f.Fd())
	}

	// Test with different cache hit rates
	scenarios := []struct {
		name   string
		numFDs int // Number of different FDs to use
	}{
		{"100%", 1}, // Always the same FD (100% hit rate)
		{"50%", 2},  // Alternate between 2 FDs (50% hit rate)
		{"20%", 5},  // Alternate between 5 FDs (20% hit rate)
		{"10%", 10}, // Alternate between 10 FDs (10% hit rate)
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				idx := i % scenario.numFDs
				probe.IsTerminal(fds[idx])
			}
		})
	}
}

// BenchmarkProbeConcurrentCacheHitRate measures the impact of cache hit rate under concurrent access
func BenchmarkProbeConcurrentCacheHitRate(b *testing.B) {
	// Create multiple file descriptors
	files := make([]*os.File, 0, 10)
	fds := make([]uintptr, 0, 10)

	// Create temporary files
	for i := 0; i < 10; i++ {
		f, err := os.CreateTemp("", "probe-benchmark-concurrent-cache")
		if err != nil {
			b.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		defer f.Close()

		files = append(files, f)
		fds = append(fds, f.Fd())
	}

	scenarios := []struct {
		name   string
		numFDs int // Number of different FDs to use
	}{
		{"100%", 1}, // Always the same FD (100% hit rate)
		{"50%", 2},  // Alternate between 2 FDs (50% hit rate)
		{"20%", 5},  // Alternate between 5 FDs (20% hit rate)
		{"10%", 10}, // Alternate between 10 FDs (10% hit rate)
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			probe.ClearCache() // Clear cache for each scenario run
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for i := 0; pb.Next(); i++ { // Use a local counter for each goroutine
					idx := i % scenario.numFDs
					probe.IsTerminal(fds[idx])
				}
			})
		})
	}
}
