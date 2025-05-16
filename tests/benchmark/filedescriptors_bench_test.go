package benchmark

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/droqsic/probe"
)

// BenchmarkProbeWithDifferentFDs measures performance with different file descriptor types
func BenchmarkProbeWithDifferentFDs(b *testing.B) {
	stdFDs := []struct {
		name string
		fd   uintptr
	}{
		{"stdin", os.Stdin.Fd()},
		{"stdout", os.Stdout.Fd()},
		{"stderr", os.Stderr.Fd()},
	}

	for _, fd := range stdFDs {
		b.Run("std-"+fd.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				probe.IsTerminal(fd.fd)
			}
		})
	}

	tmpfile, err := ioutil.TempFile("", "probe-benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	b.Run("regular-file", func(b *testing.B) {
		fd := tmpfile.Fd()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			probe.IsTerminal(fd)
		}
	})

	r, w, err := os.Pipe()
	if err != nil {
		b.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	b.Run("pipe-read", func(b *testing.B) {
		fd := r.Fd()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			probe.IsTerminal(fd)
		}
	})

	b.Run("pipe-write", func(b *testing.B) {
		fd := w.Fd()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			probe.IsTerminal(fd)
		}
	})
}

// BenchmarkCacheWithDifferentFDTypes measures performance with different file descriptor types
func BenchmarkCacheWithDifferentFDTypes(b *testing.B) {
	tmpfile, err := ioutil.TempFile("", "probe-benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	r, w, err := os.Pipe()
	if err != nil {
		b.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	fds := []uintptr{
		os.Stdout.Fd(),
		tmpfile.Fd(),
		r.Fd(),
		w.Fd(),
	}

	b.Run("mixed-fd-types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			probe.IsTerminal(fds[i%len(fds)])
		}
	})
}
