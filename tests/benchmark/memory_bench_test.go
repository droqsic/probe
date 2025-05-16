package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/probe"
)

// BenchmarkProbeMemoryAllocations measures memory allocations in probe.IsTerminal
func BenchmarkProbeMemoryAllocations(b *testing.B) {
	fd := os.Stdout.Fd()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		probe.IsTerminal(fd)
	}
}

// BenchmarkProbeCygwinMemoryAllocations measures memory allocations in probe.IsCygwinTerminal
func BenchmarkProbeCygwinMemoryAllocations(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping in short mode")
	}

	fd := os.Stdout.Fd()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		probe.IsCygwinTerminal(fd)
	}
}

// BenchmarkMemoryWithCaching measures memory allocations with caching
func BenchmarkMemoryWithCaching(b *testing.B) {
	fd := os.Stdout.Fd()

	b.Run("first-call", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			probe.ClearCache()
			probe.IsTerminal(fd)
		}
	})

	b.Run("cached-call", func(b *testing.B) {
		probe.IsTerminal(fd)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			probe.IsTerminal(fd)
		}
	})
}

// BenchmarkMemoryMultipleDescriptors measures memory allocations with multiple descriptors
func BenchmarkMemoryMultipleDescriptors(b *testing.B) {
	fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, fd := range fds {
			probe.IsTerminal(fd)
		}
	}
}
