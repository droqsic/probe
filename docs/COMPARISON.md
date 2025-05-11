# Comparison: Probe vs. go-isatty

This document provides a detailed comparison between Probe and go-isatty, the most widely used terminal detection library in the Go ecosystem.

## Platform Support

| Platform        | Probe | go-isatty | Notes                                                 |
| --------------- | :---: | :-------: | ----------------------------------------------------- |
| Windows         |  ✅   |    ✅     | Both use `GetConsoleMode`                             |
| macOS           |  ✅   |    ✅     | Both use `TIOCGETA` ioctl                             |
| Linux           |  ✅   |    ✅     | Both use `TCGETS` ioctl                               |
| FreeBSD         |  ✅   |    ✅     | Both use `TIOCGETA` ioctl                             |
| OpenBSD         |  ✅   |    ✅     | Both use `TIOCGETA` ioctl                             |
| NetBSD          |  ✅   |    ✅     | Both use `TIOCGETA` ioctl                             |
| DragonFly BSD   |  ✅   |    ✅     | Both use `TIOCGETA` ioctl                             |
| Solaris/Illumos |  ✅   |    ✅     | Both use `TCGETA` ioctl                               |
| AIX             |  ✅   |    ✅     | Probe uses `IoctlGetTermios`, go-isatty uses `TCGETA` |
| z/OS            |  ✅   |    ❌     | Probe has explicit support                            |
| Plan9           |  ✅   |    ❓     | Probe uses `Fd2path`, go-isatty has limited support   |
| Android         |  ✅   |    ✅     | Both use Linux implementation (TCGETS)                |
| iOS             |  ✅   |    ✅     | Both use Darwin implementation (TIOCGETA)             |
| Hurd            |  ✅   |    ❌     | Probe has explicit support                            |
| Cygwin/MSYS2    |  ✅   |    ✅     | Both detect special named pipes                       |

## Performance

| Metric             | Probe  | go-isatty |           Improvement           |
| ------------------ | ------ | --------- | :-----------------------------: |
| First call         | ~600ns | ~600ns    |             Similar             |
| Cached calls       | ~13ns  | ~600ns    |         **~44x faster**         |
| Cygwin detection   | ~13ns  | ~887ns    |         **~64x faster**         |
| Memory allocations | 0      | 0-2       | **Better** for Cygwin detection |

### Benchmark Results

```
BenchmarkProbeIsTerminal-12           81792342        13.68 ns/op        0 B/op        0 allocs/op
BenchmarkIsattyIsTerminal-12           2012518       608.60 ns/op        0 B/op        0 allocs/op
BenchmarkProbeCachedIsTerminal-12     89733715        13.57 ns/op        0 B/op        0 allocs/op
BenchmarkIsattyCachedIsTerminal-12     1990125       588.90 ns/op        0 B/op        0 allocs/op
BenchmarkProbeIsCygwinTerminal-12     88964013        13.73 ns/op        0 B/op        0 allocs/op
BenchmarkIsattyIsCygwinTerminal-12     1340991       887.20 ns/op       64 B/op        2 allocs/op
```

### Performance Analysis

1. **Initial Call Performance**: Both libraries have similar performance for the first call to check a file descriptor.

2. **Cached Performance**: Probe's caching mechanism provides a dramatic performance improvement for repeated checks of the same file descriptor, making it approximately 44 times faster than go-isatty.

3. **Cygwin Detection**: Probe's optimized Cygwin detection is approximately 64 times faster than go-isatty's implementation and makes no memory allocations.

4. **Memory Efficiency**: Probe makes zero memory allocations for all operations, while go-isatty makes allocations for Cygwin detection.

5. **Real-world Impact**: In applications that frequently check terminal status (such as interactive CLI tools), Probe can significantly reduce CPU usage and improve responsiveness.

## Features

| Feature                | Probe | go-isatty | Notes                                     |
| ---------------------- | :---: | :-------: | ----------------------------------------- |
| Terminal detection     |  ✅   |    ✅     | Core functionality                        |
| Cygwin detection       |  ✅   |    ✅     | Both support                              |
| Result caching         |  ✅   |    ❌     | Major performance advantage for Probe     |
| Thread safety          |  ✅   |    ✅     | Both are thread-safe                      |
| Cache clearing         |  ✅   |    N/A    | Useful for testing                        |
| Zero allocations       |  ✅   |    ❌     | Probe makes no allocations                |
| Direct platform access |  ✅   |    ❌     | Probe exposes platform-specific functions |

## API Comparison

| API                      | Probe                     | go-isatty              | Notes                             |
| ------------------------ | ------------------------- | ---------------------- | --------------------------------- |
| Terminal detection       | `IsTerminal(fd)`          | `IsTerminal(fd)`       | Same interface                    |
| Cygwin detection         | `IsCygwinTerminal(fd)`    | `IsCygwinTerminal(fd)` | Same interface                    |
| Cache control            | `ClearCache()`            | N/A                    | Additional functionality in Probe |
| Platform-specific access | `platform.IsTerminal(fd)` | N/A                    | Additional functionality in Probe |

## Implementation Details

| Aspect                 | Probe                           | go-isatty                       | Notes                              |
| ---------------------- | ------------------------------- | ------------------------------- | ---------------------------------- |
| Code organization      | Modular with platform package   | Single package                  | Probe has better separation        |
| Build tags             | Modern `//go:build`             | Legacy `// +build`              | Probe uses newer syntax            |
| Error handling         | Consistent                      | Varies by platform              | Probe has more consistent approach |
| Windows implementation | GetConsoleMode + pipe detection | GetConsoleMode + pipe detection | Similar approach                   |
| Cygwin implementation  | Optimized pipe name check       | Standard pipe name check        | Probe is more efficient            |
| Cache implementation   | RWMutex + map                   | N/A                             | Probe has sophisticated caching    |

## Documentation and Maintenance

| Aspect      |       Probe       |   go-isatty    | Notes                        |
| ----------- | :---------------: | :------------: | ---------------------------- |
| README      |   Comprehensive   |     Basic      | Probe has more detailed docs |
| Examples    | Multiple examples | Basic examples | Probe has more examples      |
| Benchmarks  |     Extensive     |    Limited     | Probe has more benchmarks    |
| Tests       |   Comprehensive   |     Basic      | Probe has more test coverage |
| Last update |      Recent       |     Active     | Both maintained              |
| Community   |        New        |  Established   | go-isatty has longer history |

## Summary

Probe offers significant improvements over go-isatty, particularly in:

1. **Performance**: The caching mechanism provides dramatic speed improvements for repeated checks
2. **Memory efficiency**: Zero allocations for all operations
3. **Platform support**: Explicit support for more platforms
4. **Code organization**: Better separation of platform-specific code
5. **Documentation**: More comprehensive documentation and examples
6. **Testing**: More extensive test suite
7. **API**: Additional functionality with platform-specific access

The main advantage of go-isatty is its established position in the Go ecosystem and longer history, but technically, Probe is superior in almost every measurable aspect.

For applications that check terminal status frequently, Probe can provide significant performance benefits while maintaining complete compatibility with the go-isatty API.
