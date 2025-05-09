package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/droqsic/probe"
)

// main is the entry point of the program.
// It checks if standard input, output, and error streams are connected to terminals.
func main() {
	checkFd("stdin", os.Stdin.Fd())   // Check standard input
	checkFd("stdout", os.Stdout.Fd()) // Check standard output
	checkFd("stderr", os.Stderr.Fd()) // Check standard error
}

// checkFd prints the name and file descriptor of a file, and whether it is connected to a terminal.
// It takes a name string for identification and a file descriptor to check.
// On Windows, it also checks if the terminal is a Cygwin/MSYS2 terminal.
// Parameters:
//   - name: A human-readable name for the file descriptor (e.g., "stdin")
//   - fd: The file descriptor to check
//
// The function prints:
//   - The name and numeric value of the file descriptor
//   - Whether it's connected to a terminal
//   - On Windows, if it's a Cygwin/MSYS2 terminal or a standard Windows console
func checkFd(name string, fd uintptr) {

	// Check if the file descriptor is connected to a terminal
	isTerminal := probe.IsTerminal(fd)
	fmt.Printf("%s (fd %d): %t", name, fd, isTerminal)

	// On Windows, differentiate between Cygwin/MSYS2 terminals and standard Windows consoles
	if runtime.GOOS == "windows" && isTerminal {
		isCygwin := probe.IsCygwinTerminal(fd)
		if isCygwin {
			fmt.Print(" (Cygwin/MSYS2 terminal)")
		} else {
			fmt.Print(" (Windows console)")
		}
	}
	fmt.Println()
}
