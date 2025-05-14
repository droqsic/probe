//go:build (js && wasm) || (tinygo && wasm)
// +build js,wasm tinygo,wasm

package platform

import "syscall/js"

// isTerminal determines if the file descriptor is a terminal in a WASM environment.
// For WebAssembly, it checks for Node.js terminal properties.
// If running in a non-Node.js environment, it returns false.
func isTerminal(fd uintptr) bool {
	global := js.Global()

	// Check if the environment is Node.js-like.
	if !global.Get("process").IsUndefined() {
		// Node.js environment assumed.
		if fd == 1 && !global.Get("process").Get("stdout").IsUndefined() &&
			!global.Get("process").Get("stdout").Get("isTTY").IsUndefined() {
			return global.Get("process").Get("stdout").Get("isTTY").Bool()
		}

		if fd == 2 && !global.Get("process").Get("stderr").IsUndefined() &&
			!global.Get("process").Get("stderr").Get("isTTY").IsUndefined() {
			return global.Get("process").Get("stderr").Get("isTTY").Bool()
		}

		if fd == 0 && !global.Get("process").Get("stdin").IsUndefined() &&
			!global.Get("process").Get("stdin").Get("isTTY").IsUndefined() {
			return global.Get("process").Get("stdin").Get("isTTY").Bool()
		}
	}

	return false
}

// isCygwin determines if the file descriptor is a Cygwin terminal in a WASM environment.
// WASM doesn't support Cygwin, so it always returns false.
func isCygwin(fd uintptr) bool {
	return false
}
