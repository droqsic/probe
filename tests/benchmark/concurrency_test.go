package benchmark

import (
	"os"
	"sync"
	"testing"

	"github.com/droqsic/probe"
	"github.com/mattn/go-isatty"
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

// BenchmarkIsattyConcurrent measures the performance of isatty.IsTerminal under concurrent access
func BenchmarkIsattyConcurrent(b *testing.B) {
	fd := os.Stdout.Fd()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			isatty.IsTerminal(fd)
		}
	})
}

// BenchmarkProbeHighConcurrency tests with a fixed number of goroutines
func BenchmarkProbeHighConcurrency(b *testing.B) {
	const goroutines = 100
	fd := os.Stdout.Fd()

	b.ResetTimer()

	b.Run("probe", func(b *testing.B) {
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
	})

	b.Run("isatty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			wg.Add(goroutines)

			for j := 0; j < goroutines; j++ {
				go func() {
					defer wg.Done()
					isatty.IsTerminal(fd)
				}()
			}

			wg.Wait()
		}
	})
}

// BenchmarkProbeConcurrentMixedFDs measures performance with concurrent access to different FDs for probe
func BenchmarkProbeConcurrentMixedFDs(b *testing.B) {
	fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
	// For more variety, you could add temporary file FDs here as in cache_test.go

	probe.ClearCache() // Ensure a clean slate for the benchmark run

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine will cycle through the FDs.
		// Using a simple counter per goroutine to select FDs.
		// A shared atomic counter or random selection could also be used
		// but might introduce other contention points not directly related to probe's cache.
		// This approach ensures mixed FD access under concurrency.
		localCounter := 0
		for pb.Next() {
			fdToTest := fds[localCounter%len(fds)]
			probe.IsTerminal(fdToTest)
			localCounter++
		}
	})
}

// BenchmarkIsattyConcurrentMixedFDs measures performance with concurrent access to different FDs for isatty
func BenchmarkIsattyConcurrentMixedFDs(b *testing.B) {
	fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
	// For more variety, you could add temporary file FDs here

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine will cycle through the FDs.
		localCounter := 0
		for pb.Next() {
			fdToTest := fds[localCounter%len(fds)]
			isatty.IsTerminal(fdToTest)
			localCounter++
		}
	})
}
