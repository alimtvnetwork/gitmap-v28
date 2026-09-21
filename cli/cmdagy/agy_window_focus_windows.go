//go:build windows

package cmdagy

import (
	"syscall"
	"unsafe"
)

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procSetForegroundWindow      = user32.NewProc("SetForegroundWindow")
	procShowWindow               = user32.NewProc("ShowWindow")
)

const swRestore = 9

func checkWindowMatch(hwnd uintptr, targetPID uint32, matched *uintptr) uintptr {
	var pid uint32
	procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	isSamePID := pid == targetPID
	if isSamePID == false {
		return 1
	}
	vis, _, _ := procIsWindowVisible.Call(hwnd)
	isVisible := vis != 0
	if isVisible {
		*matched = hwnd
		return 0
	}

	return 1
}

func findHWNDForPID(targetPID uint32) uintptr {
	var matchedHWND uintptr
	cb := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		return checkWindowMatch(hwnd, targetPID, &matchedHWND)
	})
	procEnumWindows.Call(cb, 0)

	return matchedHWND
}

// FocusAntigravityWindow activates and restores the Antigravity window for targetPID.
func FocusAntigravityWindow(targetPID int) bool {
	hasTarget := targetPID > 0
	if hasTarget == false {
		return false
	}
	hwnd := findHWNDForPID(uint32(targetPID))
	hasHWND := hwnd != 0
	if hasHWND == false {
		return false
	}
	procShowWindow.Call(hwnd, swRestore)
	res, _, _ := procSetForegroundWindow.Call(hwnd)

	return res != 0
}
