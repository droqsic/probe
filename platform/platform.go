package platform

// IsTerminal returns true if the given file descriptor is a terminal.
// This function is implemented differently for each platform.
func IsTerminal(fd uintptr) bool {
	return isTerminal(fd)
}

// IsCygwin returns true if the given file descriptor is a Cygwin/MSYS2 terminal.
// This function is only relevant on Windows and always returns false on other platforms.
func IsCygwin(fd uintptr) bool {
	return isCygwin(fd)
}
