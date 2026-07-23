//go:build windows

package cmd

import (
	"syscall"
	"unsafe"
)

// Windows console API constants for UTF-8 output and ANSI escape support.
const (
	consoleCodePageUTF8                = 65001
	consoleEnableVirtualTerminalOutput = 0x0004
	consoleStdOutHandle                = ^uintptr(10) // -11 as unsigned
	consoleStdErrHandle                = ^uintptr(11) // -12 as unsigned
)

// initConsole switches the active Windows console to UTF-8 and enables
// Virtual Terminal Processing so multi-byte glyphs (✓, →, ⚠, ▸, etc.)
// emitted by gitmap render correctly instead of as cp437/cp1252 mojibake
// (e.g. "Γ£ô" for "✓").
//
// Safe to call from any process: failures (e.g. when stdout is redirected
// to a pipe or file) are silently ignored — the caller still receives raw
// UTF-8 bytes, which is exactly what file consumers want.
func initConsole() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleOutputCP := kernel32.NewProc("SetConsoleOutputCP")
	setConsoleCP := kernel32.NewProc("SetConsoleCP")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	getStdHandle := kernel32.NewProc("GetStdHandle")

	_, _, _ = setConsoleOutputCP.Call(uintptr(consoleCodePageUTF8))
	_, _, _ = setConsoleCP.Call(uintptr(consoleCodePageUTF8))

	enableVTOnHandle(getStdHandle, getConsoleMode, setConsoleMode, consoleStdOutHandle)
	enableVTOnHandle(getStdHandle, getConsoleMode, setConsoleMode, consoleStdErrHandle)
}

// enableVTOnHandle turns on ENABLE_VIRTUAL_TERMINAL_PROCESSING for one
// standard handle. Errors are intentionally swallowed — the handle may
// not be a console (e.g. piped output) and that's fine.
func enableVTOnHandle(getStdHandle, getConsoleMode, setConsoleMode *syscall.LazyProc, stdHandleID uintptr) {
	handle, _, _ := getStdHandle.Call(stdHandleID)
	if handle == 0 || handle == ^uintptr(0) {
		return
	}

	var mode uint32
	ret, _, _ := getConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	if ret == 0 {
		return
	}

	mode |= consoleEnableVirtualTerminalOutput
	_, _, _ = setConsoleMode.Call(handle, uintptr(mode))
}
