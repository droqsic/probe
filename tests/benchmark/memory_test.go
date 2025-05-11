package benchmark

import (
    "os"
    "testing"

    "github.com/droqsic/probe"
    "github.com/mattn/go-isatty"
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

// BenchmarkIsattyMemoryAllocations measures memory allocations in isatty.IsTerminal
func BenchmarkIsattyMemoryAllocations(b *testing.B) {
    fd := os.Stdout.Fd()
    
    b.ReportAllocs()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        isatty.IsTerminal(fd)
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

// BenchmarkIsattyCygwinMemoryAllocations measures memory allocations in isatty.IsCygwinTerminal
func BenchmarkIsattyCygwinMemoryAllocations(b *testing.B) {
    if testing.Short() {
        b.Skip("Skipping in short mode")
    }
    
    fd := os.Stdout.Fd()
    
    b.ReportAllocs()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        isatty.IsCygwinTerminal(fd)
    }
}