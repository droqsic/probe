package integration

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/droqsic/probe"
	"github.com/mattn/go-isatty"
)

// TestStressConcurrency tests with high concurrency
func TestStressConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	const goroutines = 1000
	const iterations = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	// Launch goroutines for isatty
	start := time.Now()
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				isatty.IsTerminal(os.Stdout.Fd())
			}
		}()
	}

	// Launch goroutines for probe
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				probe.IsTerminal(os.Stdout.Fd())
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	t.Logf("Stress test completed in %v", elapsed)
}

// TestRapidSequentialCalls tests behavior with rapid sequential calls
func TestRapidSequentialCalls(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping rapid sequential test in short mode")
	}

	const iterations = 100000

	// Test isatty
	startIsatty := time.Now()
	for i := 0; i < iterations; i++ {
		isatty.IsTerminal(os.Stdout.Fd())
	}
	isattyElapsed := time.Since(startIsatty)

	// Test probe
	startProbe := time.Now()
	for i := 0; i < iterations; i++ {
		probe.IsTerminal(os.Stdout.Fd())
	}
	probeElapsed := time.Since(startProbe)

	t.Logf("Rapid sequential calls (%d iterations):", iterations)
	t.Logf("  isatty: %v", isattyElapsed)
	t.Logf("  probe: %v", probeElapsed)
	t.Logf("  probe is %.2fx faster", float64(isattyElapsed)/float64(probeElapsed))
}

// TestAlternatingFileDescriptors tests with alternating file descriptors
func TestAlternatingFileDescriptors(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping alternating FD test in short mode")
	}

	const iterations = 10000

	// File descriptors to alternate between
	fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}

	// Test isatty with alternating FDs
	startIsatty := time.Now()
	for i := 0; i < iterations; i++ {
		isatty.IsTerminal(fds[i%len(fds)])
	}
	isattyElapsed := time.Since(startIsatty)

	// Test probe with alternating FDs
	startProbe := time.Now()
	for i := 0; i < iterations; i++ {
		probe.IsTerminal(fds[i%len(fds)])
	}
	probeElapsed := time.Since(startProbe)

	t.Logf("Alternating file descriptors (%d iterations):", iterations)
	t.Logf("  isatty: %v", isattyElapsed)
	t.Logf("  probe: %v", probeElapsed)
	t.Logf("  probe is %.2fx faster", float64(isattyElapsed)/float64(probeElapsed))
}

// TestCachingPerformance tests the performance impact of caching
func TestCachingPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping caching performance test in short mode")
	}

	const iterations = 100000

	// First call to prime the cache
	probe.IsTerminal(os.Stdout.Fd())

	// Measure cached performance
	start := time.Now()
	for i := 0; i < iterations; i++ {
		probe.IsTerminal(os.Stdout.Fd())
	}
	elapsed := time.Since(start)

	t.Logf("Cached performance (%d iterations): %v", iterations, elapsed)
	t.Logf("Average time per call: %v", elapsed/time.Duration(iterations))
}
