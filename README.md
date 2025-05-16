# Probe

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/droqsic/probe.svg)](https://pkg.go.dev/github.com/droqsic/probe)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Workflow](https://github.com/droqsic/probe/actions/workflows/go.yml/badge.svg)](https://github.com/droqsic/probe/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/droqsic/probe?nocache=1)](https://goreportcard.com/report/github.com/droqsic/probe)
[![codecov](https://codecov.io/gh/droqsic/probe/branch/main/graph/badge.svg)](https://codecov.io/gh/droqsic/probe)
[![Latest Release](https://img.shields.io/github/v/release/droqsic/probe)](https://github.com/droqsic/probe/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/droqsic/probe)](https://golang.org/)

</div>

Probe is a lightweight, cross-platform Go library for detecting whether a file descriptor is connected to a terminal. It is designed for high performance, thread safety, and ease of use in modern Go applications.

## Features

- ⚡ **High Performance** — Fast terminal detection with intelligent result caching
- 🌍 **Cross-Platform** — Works seamlessly on Windows, macOS, Linux, BSD, and more
- 🧵 **Thread-Safe** — Fully safe for concurrent use across goroutines
- 🧠 **Zero Allocations** — Designed to avoid heap allocations entirely
- 🧩 **Minimal Dependencies** — Built using only the Go standard library and x/sys
- 🧼 **Simple API** — Clean and intuitive interface for fast integration

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

## Performance

Probe is engineered for speed. Its caching layer makes repeated checks on the same file descriptor nearly instantaneous. Here are benchmark results under typical usage:

```
BenchmarkIsTerminal                       89545555	        13.61 ns/op	       0 B/op	      0 allocs/op
BenchmarkIsCygwinTerminal                 89867444	        13.68 ns/op	       0 B/op	      0 allocs/op
```

These results demonstrate Probe's ultra-low overhead and suitability for high-throughput applications.

## How It Works

Probe uses platform-specific mechanisms for terminal detection:

- **Unix-like systems**: Uses appropriate ioctl calls (`TIOCGETA`, `TCGETS`, or `TCGETA`)
- **Windows**: Uses the Win32 `GetConsoleMode` function
- **Cygwin/MSYS2**: Detects special named pipes used by Cygwin terminals
- **WebAssembly**: Detects terminals based on Node.js environment variables
- **Other platforms**: Always returns `false`

All results are cached after the first check per file descriptor to avoid repeated syscalls.

## Supported Platforms

| Platform       | Support | Implementation    |
| -------------- | ------- | ----------------- |
| Windows        | ✅      | `GetConsoleMode`  |
| Linux          | ✅      | `TCGETS` ioctl    |
| Android        | ✅      | `TCGETS` ioctl    |
| macOS (Darwin) | ✅      | `TIOCGETA` ioctl  |
| iOS            | ✅      | `TIOCGETA` ioctl  |
| FreeBSD        | ✅      | `TIOCGETA` ioctl  |
| OpenBSD        | ✅      | `TIOCGETA` ioctl  |
| NetBSD         | ✅      | `TIOCGETA` ioctl  |
| DragonFly BSD  | ✅      | `TIOCGETA` ioctl  |
| Hurd           | ✅      | `TIOCGETA` ioctl  |
| z/OS           | ✅      | `TIOCGETA` ioctl  |
| Solaris        | ✅      | `TCGETA` ioctl    |
| Illumos        | ✅      | `TCGETA` ioctl    |
| Haikou         | ✅      | `TCGETA` ioctl    |
| AIX            | ✅      | `TCGETA` ioctl    |
| Plan9          | ✅      | `Fd2path`         |
| WebAssembly    | ✅      | Node.js detection |

Other platforms will compile but terminal detection will always return `false`.

## Thread Safety

Probe is built with concurrency in mind. It uses a read-write mutex to protect its internal cache, allowing many goroutines to read in parallel while safely handling updates. The design ensures scalability and efficiency for read-heavy workloads.

## Contributing

Contributions are welcome! Whether you're fixing bugs, adding features, or improving documentation, your input helps make Probe better for everyone.

- Read the [Contributing Guidelines](docs/CONTRIBUTING.md)
- Follow the [Code of Conduct](docs/CODE_OF_CONDUCT.md)

## License

Probe is released under the MIT License. For the full license text, please see the [LICENSE](LICENSE) file.

## Acknowledgements

Probe is inspired by [go-isatty](https://github.com/mattn/go-isatty) by Yasuhiro Matsumoto. While conceptually similar, Probe provides an enhanced implementation focused on performance and scalability.

Special thanks to:

- The Go team for their exceptional language and tooling
- The maintainers of x/sys for low-level system access
- All contributors who help make Probe better
