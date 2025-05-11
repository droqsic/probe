package probe

import (
	"runtime"
	"sync"

	"github.com/droqsic/probe/platform"
)

// Cache store the result of IsTerminal and IsCygwinTerminal calls for each file descriptor.
// It is used to avoid calling the underlying platform functions multiple times for the same file descriptor.
var (
	cache = struct {
		terminal map[uintptr]bool // Maps file descriptors to terminal status
		cygwin   map[uintptr]bool // Maps file descriptors to Cygwin status
		mutex    sync.RWMutex     // Protects concurrent access to the maps
	}{
		terminal: make(map[uintptr]bool),
		cygwin:   make(map[uintptr]bool),
	}
)

// getCache retrieves a cached result for a file descriptor.
// It returns the cached value and a boolean indicating if the value was found in the cache.
func getCache(m map[uintptr]bool, fd uintptr) (bool, bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	result, ok := m[fd]
	return result, ok
}

// setCache stores a result for a file descriptor in the cache.
func setCache(m map[uintptr]bool, fd uintptr, result bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	m[fd] = result
}

// IsTerminal returns true if the file descriptor is a terminal.
// It uses platform-specific implementations and cache results for performance.
// This function is thread-safe and can be called from multiple goroutines.
func IsTerminal(fd uintptr) bool {
	// Check cache first to avoid expensive platform calls.
	if result, ok := getCache(cache.terminal, fd); ok {
		return result
	}

	// Determine if the file descriptor is a terminal based on the platform.
	// The platform-specific implementation is selected at compile time.
	result := platform.IsTerminal(fd)

	// Cache the result for future use.
	setCache(cache.terminal, fd, result)
	return result
}

// IsCygwinTerminal returns true if the file descriptor is a Cygwin/MSYS2 terminal.
// This function is only relevant on Windows and always returns false on other platforms.
// This function is thread-safe and can be called from multiple goroutines.
func IsCygwinTerminal(fd uintptr) bool {
	// Early return for non-Windows platforms.
	if runtime.GOOS != "windows" {
		return false
	}

	// Check cache first to avoid expensive platform calls.
	if result, ok := getCache(cache.cygwin, fd); ok {
		return result
	}

	result := platform.IsCygwin(fd)    // Call the platform-specific implementation.
	setCache(cache.cygwin, fd, result) // Cache the result for future use.
	return result
}

// ClearCache clears the internal cache.
// This is mainly useful for testing purposes.
func ClearCache() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.terminal = make(map[uintptr]bool)
	cache.cygwin = make(map[uintptr]bool)
}
