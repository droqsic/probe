# Probe

[![Go Reference](https://pkg.go.dev/badge/github.com/droqsic/probe.svg)](https://pkg.go.dev/github.com/droqsic/probe)
[![Go Report Card](https://goreportcard.com/badge/github.com/droqsic/probe)](https://goreportcard.com/report/github.com/droqsic/probe)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A lightweight, cross-platform Go library to detect if a file descriptor is connected to a terminal.

## Features

- Cross-platform terminal detection (Windows, macOS, Linux, BSD)
- Cygwin/MSYS2 terminal detection on Windows
- Thread-safe implementation
- High-performance caching mechanism
- Zero dependencies (except for standard library and x/sys)
- Simple, clean API

## Installation

```bash
go get github.com/droqsic/probe
```

## Usage

```go
package main

import (
    "fmt"
    "os"

    "github.com/droqsic/probe"
)

func main() {
    // Check if standard streams are terminals
    if probe.IsTerminal(os.Stdout.Fd()) {
        fmt.Println("Stdout is a terminal")
    } else {
        fmt.Println("Stdout is not a terminal (redirected to a file or pipe)")
    }

    // On Windows, check for Cygwin/MSYS2 terminals
    if probe.IsCygwinTerminal(os.Stdout.Fd()) {
        fmt.Println("Running in a Cygwin/MSYS2 terminal")
    }
}
```

## How It Works

Probe uses platform-specific mechanisms to detect terminals:

- **Unix-like systems** (Linux, macOS, BSD): Uses the `TIOCGETA` ioctl call
- **Windows**: Uses the Win32 `GetConsoleMode` function
- **Cygwin/MSYS2 on Windows**: Detects special named pipes used by Cygwin terminals

Results are cached for performance, making repeated checks on the same file descriptor extremely fast.

## Performance

Terminal detection can be an expensive operation, especially on Windows with Cygwin detection. Probe implements an efficient caching mechanism that makes subsequent checks on the same file descriptor nearly instantaneous.

### Benchmark Comparison

Probe significantly outperforms other terminal detection libraries. Here's a comparison with the popular `go-isatty` library:

```
BenchmarkProbeIsTerminal-12           81792342        13.68 ns/op        0 B/op        0 allocs/op
BenchmarkIsattyIsTerminal-12           2012518       608.60 ns/op        0 B/op        0 allocs/op
BenchmarkProbeCachedIsTerminal-12     89733715        13.57 ns/op        0 B/op        0 allocs/op
BenchmarkIsattyCachedIsTerminal-12     1990125       588.90 ns/op        0 B/op        0 allocs/op
BenchmarkProbeIsCygwinTerminal-12     88964013        13.73 ns/op        0 B/op        0 allocs/op
BenchmarkIsattyIsCygwinTerminal-12     1340991       887.20 ns/op       64 B/op        2 allocs/op
```

Key observations (values are approximate and may vary based on hardware and Go version):

- Probe's terminal detection is **44x faster** than go-isatty
- Probe's Cygwin terminal detection is **64x faster** than go-isatty
- Probe makes **zero memory allocations** for all operations

## Supported Platforms

- Windows
- macOS
- Linux
- FreeBSD
- OpenBSD
- NetBSD
- Solaris/Illumos
- AIX
- z/OS

Other platforms will compile but terminal detection will always return `false`.

## Thread Safety

Probe is designed with concurrency in mind, ensuring that all functions can be safely called from multiple goroutines simultaneously. The library implements a sophisticated synchronization mechanism using read-write mutexes to protect its internal cache, allowing for high-throughput concurrent access patterns common in modern Go applications.

The caching layer is particularly optimized for the read-heavy workloads typical of terminal detection scenarios, where the same file descriptor is checked repeatedly. By using a read-write mutex instead of a standard mutex, Probe allows multiple goroutines to read from the cache simultaneously, only blocking when the cache needs to be updated with new information.

This thread-safe design ensures that Probe can be used in high-concurrency environments without performance degradation or race conditions, making it suitable for server applications, concurrent CLI tools, and other multi-threaded software.

## Contributing

Contributions to Probe are warmly welcomed and greatly appreciated. Whether you're fixing a bug, adding a feature, improving documentation, or suggesting enhancements, your help makes this project better for everyone.

Please see our [Contributing Guidelines](guides/docs/CONTRIBUTING.md) for details on how to contribute.

All contributors are expected to adhere to our [Code of Conduct](guides/docs/CODE_OF_CONDUCT.md), fostering an open and welcoming environment for everyone.

## License

Probe is released under the MIT License, a permissive open source license that places minimal restrictions on reuse and distribution. This means you can use Probe in your own projects, whether they're open source or commercial, with very few limitations.

The MIT License grants permissions to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the software, provided that the original copyright notice and permission notice appear in all copies or substantial portions of the software.

For the full license text, please see the [LICENSE](LICENSE) file included in this repository.

## Acknowledgements

Probe was inspired by and builds upon the foundation laid by [go-isatty](https://github.com/mattn/go-isatty), created by Yasuhiro Matsumoto. While reimplementing the core functionality, Probe introduces significant performance optimizations and a more efficient caching mechanism.

Special thanks to:

- The Go team for creating such an excellent programming language
- The maintainers of the x/sys package for providing the low-level system interfaces
- All contributors who have helped improve this project
- The open source community for their continuous support and feedback

The benchmark comparisons and performance improvements would not have been possible without the original work done by the go-isatty project, which has been a valuable resource for terminal detection in the Go ecosystem for many years.
