package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/probe"
)

// BenchmarkCacheHitRate measures the cache hit rate
func BenchmarkCacheHitRate(b *testing.B) {
	fd := os.Stdout.Fd()

	b.Run("100-percent-hits", func(b *testing.B) {
		probe.IsTerminal(fd)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			probe.IsTerminal(fd)
		}
	})

	b.Run("50-percent-hits", func(b *testing.B) {
		probe.IsTerminal(fd)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if i%2 == 0 {
				probe.ClearCache()
			}
			probe.IsTerminal(fd)
		}
	})

	b.Run("0-percent-hits", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			probe.ClearCache()
			probe.IsTerminal(fd)
		}
	})
}

// BenchmarkMultipleDescriptorsCaching measures the performance of IsTerminal with multiple descriptors
func BenchmarkMultipleDescriptorsCaching(b *testing.B) {
	fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}

	b.Run("all-descriptors", func(b *testing.B) {
		for _, fd := range fds {
			probe.IsTerminal(fd)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, fd := range fds {
				probe.IsTerminal(fd)
			}
		}
	})

	b.Run("alternating-descriptors", func(b *testing.B) {
		for _, fd := range fds {
			probe.IsTerminal(fd)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			probe.IsTerminal(fds[i%len(fds)])
		}
	})
}

// BenchmarkCacheOverhead measures the overhead of caching
func BenchmarkCacheOverhead(b *testing.B) {
	fd := os.Stdout.Fd()

	b.Run("with-cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			probe.IsTerminal(fd)
		}
	})

	b.Run("direct-platform-call", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			probe.ClearCache()
			probe.IsTerminal(fd)
		}
	})
}
