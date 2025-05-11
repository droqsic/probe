package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/probe"
)

// BenchmarkIsTerminal measures the performance of probe.IsTerminal
func BenchmarkIsTerminal(b *testing.B) {
	// Clear cache before benchmarking
	probe.ClearCache()

	fd := os.Stdout.Fd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		probe.IsTerminal(fd)
	}
}

// BenchmarkIsCygwinTerminal measures the performance of probe.IsCygwinTerminal
func BenchmarkIsCygwinTerminal(b *testing.B) {
	// Clear cache before benchmarking
	probe.ClearCache()

	fd := os.Stdout.Fd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		probe.IsCygwinTerminal(fd)
	}
}

// BenchmarkIsTerminalCached measures the performance of repeated calls to IsTerminal
// This demonstrates the effectiveness of the caching mechanism
func BenchmarkIsTerminalCached(b *testing.B) {
	fd := os.Stdout.Fd()

	// Warm up the cache with one call
	probe.IsTerminal(fd)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		probe.IsTerminal(fd)
	}
}

// BenchmarkIsCygwinTerminalCached measures the performance of repeated calls to IsCygwinTerminal
// This demonstrates the effectiveness of the caching mechanism
func BenchmarkIsCygwinTerminalCached(b *testing.B) {
	fd := os.Stdout.Fd()

	// Warm up the cache with one call
	probe.IsCygwinTerminal(fd)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		probe.IsCygwinTerminal(fd)
	}
}

// BenchmarkMultipleFileDescriptors measures the performance when checking multiple file descriptors
func BenchmarkMultipleFileDescriptors(b *testing.B) {
	// Clear cache before benchmarking
	probe.ClearCache()

	fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, fd := range fds {
			probe.IsTerminal(fd)
		}
	}
}
