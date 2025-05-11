//go:build windows
// +build windows

package platform

import (
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	fileTypePipe   = 3 // Pipe file type constant for GetFileType
	fileTypeChar   = 2 // Character file type constant for GetFileType
	fileNameInfo   = 2 // File name information class constant for GetFileInformationByHandleEx
	objectNameInfo = 1 // Object name information class constant for NtQueryObject
)

// Windows API function pointers and flags
var (
	kernel32                         = syscall.NewLazyDLL("kernel32.dll")
	ntdll                            = syscall.NewLazyDLL("ntdll.dll")
	procGetConsoleMode               = kernel32.NewProc("GetConsoleMode")
	procGetFileInformationByHandleEx = kernel32.NewProc("GetFileInformationByHandleEx")
	procGetFileType                  = kernel32.NewProc("GetFileType")
	procNtQueryObject                = ntdll.NewProc("NtQueryObject")
	hasGetFileInfoByHandleEx         = procGetFileInformationByHandleEx.Find() == nil
)

// isTerminal checks if the file descriptor is a Windows console.
// It uses the GetConsoleMode function, which is available on all Windows versions.
func isTerminal(fd uintptr) bool {
	ft, _, _ := syscall.Syscall(procGetFileType.Addr(), 1, fd, 0, 0)
	if ft == fileTypePipe {
		return false
	}

	var mode uint32
	r, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, fd, uintptr(unsafe.Pointer(&mode)), 0)
	return r != 0 && e == 0
}

// isCygwin checks if the file descriptor is a Cygwin/MSYS2 terminal.
// Cygwin/MSYS2 terminals are implemented as named pipes with specific naming conventions.
func isCygwin(fd uintptr) bool {
	// Check for pipe type first. Cygwin/MSYS2 terminals are always pipes.
	ft, _, _ := syscall.Syscall(procGetFileType.Addr(), 1, fd, 0, 0)
	if ft != fileTypePipe {
		return false
	}

	// Get pipe name using the appropriate method based on Windows version.
	var pipeName string
	var err error

	if hasGetFileInfoByHandleEx {
		// Use GetFileInformationByHandleEx on newer Windows versions.
		var buf [2 + syscall.MAX_PATH]uint16
		r, _, _ := syscall.Syscall6(procGetFileInformationByHandleEx.Addr(),
			4, fd, fileNameInfo, uintptr(unsafe.Pointer(&buf)), uintptr(len(buf)*2), 0, 0)
		if r == 0 {
			return false
		}
		l := *(*uint32)(unsafe.Pointer(&buf))
		pipeName = string(utf16.Decode(buf[2 : 2+l/2]))
	} else {
		// Fallback to NtQueryObject on older Windows versions (XP, Server 2003).
		pipeName, err = getFileNameByHandle(fd)
		if err != nil {
			return false
		}
	}

	// Check if the pipe name matches the Cygwin/MSYS2 naming pattern.
	return isCygwinPipeName(pipeName)
}

// isCygwinPipeName checks if a pipe name matches the Cygwin/MSYS2 naming pattern.
// Cygwin/MSYS2 PTY has a name like: \{cygwin,msys}-XXXXXXXXXXXXXXXX-ptyN-{from,to}-master
func isCygwinPipeName(name string) bool {
	tokens := strings.Split(name, "-")

	// Check the number of tokens and the prefix.
	if len(tokens) < 5 {
		return false
	}

	// Check the prefix, pty name, and direction.
	if !strings.Contains(tokens[0], "msys") && !strings.Contains(tokens[0], "cygwin") {
		return false
	}

	// Check the pty name and direction.
	if !strings.HasPrefix(tokens[2], "pty") {
		return false
	}

	// Check the direction and the "master" suffix.
	if tokens[3] != "from" && tokens[3] != "to" {
		return false
	}

	// Check the "master" suffix.
	return tokens[4] == "master"
}

// getFileNameByHandle retrieves the file name for a given file descriptor using NtQueryObject.
// This function is used as a fallback on older Windows versions where GetFileInformationByHandleEx is not available.
func getFileNameByHandle(fd uintptr) (string, error) {
	if procNtQueryObject == nil {
		return "", syscall.EWINDOWS
	}
	var buf [4 + syscall.MAX_PATH]uint16
	var result int
	r, _, _ := syscall.Syscall6(procNtQueryObject.Addr(), 5,
		fd, objectNameInfo, uintptr(unsafe.Pointer(&buf)), uintptr(2*len(buf)), uintptr(unsafe.Pointer(&result)), 0)
	if r != 0 {
		return "", syscall.EINVAL
	}
	return string(utf16.Decode(buf[4 : 4+buf[0]/2])), nil
}
