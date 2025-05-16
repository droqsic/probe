package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/probe"
)

// BenchmarkIsTerminal measures the performance of IsTerminal
func BenchmarkIsTerminal(b *testing.B) {
	fd := os.Stdout.Fd()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		probe.IsTerminal(fd)
	}
}

// BenchmarkIsCygwinTerminal measures the performance of IsCygwinTerminal
func BenchmarkIsCygwinTerminal(b *testing.B) {
	fd := os.Stdout.Fd()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		probe.IsCygwinTerminal(fd)
	}
}

// BenchmarkCacheBenefit measures the benefit of caching
func BenchmarkCacheBenefit(b *testing.B) {
	fd := os.Stdout.Fd()

	b.Run("with-cache", func(b *testing.B) {
		probe.IsTerminal(fd)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			probe.IsTerminal(fd)
		}
	})

	b.Run("no-cache", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			probe.ClearCache()
			probe.IsTerminal(fd)
		}
	})
}

// BenchmarkMultipleFileDescriptors measures the performance of IsTerminal with multiple file descriptors
func BenchmarkMultipleFileDescriptors(b *testing.B) {
	fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}

	b.Run("sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, fd := range fds {
				probe.IsTerminal(fd)
			}
		}
	})

	b.Run("alternating", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			probe.IsTerminal(fds[i%len(fds)])
		}
	})
}
