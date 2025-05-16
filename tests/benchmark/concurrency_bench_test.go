package benchmark

import (
	"os"
	"sync"
	"testing"

	"github.com/droqsic/probe"
)

// BenchmarkProbeConcurrent measures the performance of probe.IsTerminal under concurrent access
func BenchmarkProbeConcurrent(b *testing.B) {
	fd := os.Stdout.Fd()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			probe.IsTerminal(fd)
		}
	})
}

// BenchmarkProbeHighConcurrency tests with a fixed number of goroutines
func BenchmarkProbeHighConcurrency(b *testing.B) {
	const goroutines = 100
	fd := os.Stdout.Fd()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for j := 0; j < goroutines; j++ {
			go func() {
				defer wg.Done()
				probe.IsTerminal(fd)
			}()
		}

		wg.Wait()
	}
}

// BenchmarkProbeConcurrentMixedFDs measures performance with concurrent access to different FDs for probe
func BenchmarkProbeConcurrentMixedFDs(b *testing.B) {
	fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		localCounter := 0
		for pb.Next() {
			fdToTest := fds[localCounter%len(fds)]
			probe.IsTerminal(fdToTest)
			localCounter++
		}
	})
}

// BenchmarkConcurrentCacheClearing measures performance with concurrent cache clearing
func BenchmarkConcurrentCacheClearing(b *testing.B) {
	fd := os.Stdout.Fd()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		localCounter := 0
		for pb.Next() {
			if localCounter%10 == 0 {
				probe.ClearCache()
			}
			probe.IsTerminal(fd)
			localCounter++
		}
	})
}
