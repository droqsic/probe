package benchmark

import (
	"os"
	"runtime"
	"testing"

	"github.com/droqsic/probe"
)

// BenchmarkPlatformSpecific runs benchmarks specific to the current platform
func BenchmarkPlatformSpecific(b *testing.B) {
	fd := os.Stdout.Fd()

	switch runtime.GOOS {
	case "windows":
		b.Run("windows-terminal", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				probe.IsTerminal(fd)
			}
		})

		b.Run("windows-cygwin", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				probe.IsCygwinTerminal(fd)
			}
		})

	case "linux", "android":
		b.Run("linux-tcgets", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				probe.IsTerminal(fd)
			}
		})

	case "darwin", "freebsd", "openbsd", "netbsd", "dragonfly":
		b.Run("bsd-tiocgeta", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				probe.IsTerminal(fd)
			}
		})

	case "solaris", "illumos", "haikou":
		b.Run("solaris-tcgeta", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				probe.IsTerminal(fd)
			}
		})

	case "plan9":
		b.Run("plan9-fd2path", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				probe.IsTerminal(fd)
			}
		})

	case "js":
		if runtime.GOARCH == "wasm" {
			b.Run("wasm", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					probe.IsTerminal(fd)
				}
			})
		}

	default:
		b.Run("generic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				probe.IsTerminal(fd)
			}
		})
	}
}

// BenchmarkPlatformComparison compares performance across different operations on the current platform
func BenchmarkPlatformComparison(b *testing.B) {
	fds := []struct {
		name string
		fd   uintptr
	}{
		{"stdin", os.Stdin.Fd()},
		{"stdout", os.Stdout.Fd()},
		{"stderr", os.Stderr.Fd()},
	}

	for _, fd := range fds {
		b.Run(fd.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				probe.IsTerminal(fd.fd)
			}
		})
	}

	f, err := os.CreateTemp("", "probe-bench")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	b.Run("regular-file", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			probe.IsTerminal(f.Fd())
		}
	})

	r, w, err := os.Pipe()
	if err != nil {
		b.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	b.Run("pipe-reader", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			probe.IsTerminal(r.Fd())
		}
	})

	b.Run("pipe-writer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			probe.IsTerminal(w.Fd())
		}
	})
}
