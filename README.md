# Probe

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/droqsic/probe.svg)](https://pkg.go.dev/github.com/droqsic/probe)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Workflow](https://github.com/droqsic/probe/actions/workflows/go.yml/badge.svg)](https://github.com/droqsic/probe/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/droqsic/probe?nocache=1)](https://goreportcard.com/report/github.com/droqsic/probe)
[![Latest Release](https://img.shields.io/github/v/release/droqsic/probe)](https://github.com/droqsic/probe/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/droqsic/probe)](https://golang.org/)

</div>

A lightweight, cross-platform Go library to detect if a file descriptor is connected to a terminal. Probe offers superior performance with a sophisticated caching mechanism, making it up to 44x faster than alternatives.

## Features

- **High Performance**: Optimized caching mechanism makes repeated checks nearly instantaneous
- **Cross-Platform**: Works on all major platforms including Windows, macOS, Linux, BSD, and more
- **Zero Allocations**: Makes no memory allocations for any operations
- **Thread-Safe**: Designed for concurrent access from multiple goroutines
- **Simple API**: Clean, intuitive interface that's easy to integrate
- **Minimal Dependencies**: Only depends on standard library and x/sys

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

- **Unix-like systems**: Uses appropriate ioctl calls (`TIOCGETA`, `TCGETS`, or `TCGETA`)
- **Windows**: Uses the Win32 `GetConsoleMode` function
- **Cygwin/MSYS2**: Detects special named pipes used by Cygwin terminals

Results are cached for performance, making repeated checks on the same file descriptor extremely fast.

## Performance

Terminal detection can be an expensive operation. Probe implements an efficient caching mechanism that makes subsequent checks nearly instantaneous.

### Benchmark Comparison

```
BenchmarkProbeIsTerminal-12           81792342        13.68 ns/op        0 B/op        0 allocs/op
BenchmarkIsattyIsTerminal-12           2012518       608.60 ns/op        0 B/op        0 allocs/op
BenchmarkProbeCachedIsTerminal-12     89733715        13.57 ns/op        0 B/op        0 allocs/op
BenchmarkIsattyCachedIsTerminal-12     1990125       588.90 ns/op        0 B/op        0 allocs/op
BenchmarkProbeIsCygwinTerminal-12     88964013        13.73 ns/op        0 B/op        0 allocs/op
BenchmarkIsattyIsCygwinTerminal-12     1340991       887.20 ns/op       64 B/op        2 allocs/op
```

Key observations:

- Probe's terminal detection is **44x faster** than alternatives
- Probe's Cygwin terminal detection is **64x faster** than alternatives
- Probe makes **zero memory allocations** for all operations

For a detailed comparison with other libraries, see [COMPARISON.md](docs/COMPARISON.md).

## Supported Platforms

| Platform        | Support | Implementation   |
| --------------- | ------- | ---------------- |
| Windows         | ✅      | `GetConsoleMode` |
| macOS           | ✅      | `TIOCGETA` ioctl |
| Linux           | ✅      | `TCGETS` ioctl   |
| FreeBSD         | ✅      | `TIOCGETA` ioctl |
| OpenBSD         | ✅      | `TIOCGETA` ioctl |
| NetBSD          | ✅      | `TIOCGETA` ioctl |
| Solaris/Illumos | ✅      | `TCGETA` ioctl   |
| AIX             | ✅      | `TCGETA` ioctl   |
| z/OS            | ✅      | `TIOCGETA` ioctl |
| Plan9           | ✅      | `Fd2path`        |
| Android         | ✅      | `TCGETS` ioctl   |
| iOS             | ✅      | `TIOCGETA` ioctl |
| Hurd            | ✅      | `TIOCGETA` ioctl |

Other platforms will compile but terminal detection will always return `false`.

## Thread Safety

Probe is designed with concurrency in mind. The library implements a synchronization mechanism using read-write mutexes to protect its internal cache, allowing for high-throughput concurrent access patterns common in modern Go applications.

The caching layer is optimized for read-heavy workloads typical of terminal detection scenarios. By using a read-write mutex, Probe allows multiple goroutines to read from the cache simultaneously, only blocking when the cache needs to be updated.

## Contributing

Contributions to Probe are warmly welcomed. Whether you're fixing a bug, adding a feature, or improving documentation, your help makes this project better for everyone.

Please see our [Contributing Guidelines](docs/CONTRIBUTING.md) for details on how to contribute.

All contributors are expected to adhere to our [Code of Conduct](docs/CODE_OF_CONDUCT.md).

## License

Probe is released under the MIT License. For the full license text, please see the [LICENSE](LICENSE) file.

## Acknowledgements

Probe was inspired by [go-isatty](https://github.com/mattn/go-isatty), created by Yasuhiro Matsumoto. While reimplementing the core functionality, Probe introduces significant performance optimizations and a more efficient caching mechanism.

Special thanks to:

- The Go team for creating such an excellent programming language
- The maintainers of the x/sys package for providing the low-level system interfaces
- All contributors who have helped improve this project
