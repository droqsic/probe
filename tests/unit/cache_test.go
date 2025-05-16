package unit

import (
	"os"
	"testing"

	"github.com/droqsic/probe"
)

// TestCacheConsistency tests that the cache returns consistent results.
// This test checks that the cache returns consistent results for standard file descriptors.
// It does not check the correctness of the terminal detection, only that the cache is consistent.
func TestCacheConsistency(t *testing.T) {
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
			firstResult := probe.IsTerminal(fd.fd)
			secondResult := probe.IsTerminal(fd.fd)

			if firstResult != secondResult {
				t.Errorf("Inconsistent results between uncached and cached calls for %s", fd.name)
			}

			probe.ClearCache()
			thirdResult := probe.IsTerminal(fd.fd)

			if firstResult != thirdResult {
				t.Errorf("Inconsistent results after clearing cache for %s", fd.name)
			}
		})
	}
}

// TestCacheClear tests that ClearCache properly clears the cache.
// This test uses a pipe to ensure we have a non-terminal file descriptor.
// It checks that the pipe is not detected as a terminal before and after clearing the cache.
func TestCacheClear(t *testing.T) {
	fd := os.Stdout.Fd()

	probe.IsTerminal(fd)
	probe.ClearCache()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	if probe.IsTerminal(r.Fd()) {
		t.Errorf("Pipe reader should not be detected as terminal")
	}

	if probe.IsTerminal(w.Fd()) {
		t.Errorf("Pipe writer should not be detected as terminal")
	}

	probe.ClearCache()

	if probe.IsTerminal(r.Fd()) {
		t.Errorf("Pipe reader should not be detected as terminal after clearing cache")
	}
}

// TestMultipleFileDescriptors tests cache behavior with multiple file descriptors.
// This test creates multiple temporary files and checks that they are not detected as terminal.
// It also checks that standard file descriptors are consistent with the cache.
func TestMultipleFileDescriptors(t *testing.T) {
	numFiles := 5
	files := make([]*os.File, numFiles)
	for i := range files {
		f, err := os.CreateTemp("", "probe-test")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		defer f.Close()
		files[i] = f
	}

	probe.ClearCache()

	for i, f := range files {
		result := probe.IsTerminal(f.Fd())
		if result {
			t.Errorf("File %d should not be detected as terminal", i)
		}
	}

	stdFDs := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
	for _, fd := range stdFDs {
		first := probe.IsTerminal(fd)
		second := probe.IsTerminal(fd)

		if first != second {
			t.Errorf("Inconsistent results for fd %d", fd)
		}
	}
}
