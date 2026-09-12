//go:build windows

package cmdpurge

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

type shFileOpStruct struct {
	hwnd                  uintptr
	wFunc                 uint32
	pFrom                 *uint16
	pTo                   *uint16
	fFlags                uint16
	fAnyOperationsAborted int32
	hNameMappings         uintptr
	lpszProgressTitle     *uint16
}

func sendToRecycleBin(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil
	}

	return invokeShFileDelete(absPath)
}

func invokeShFileDelete(absPath string) error {
	pFrom, err := syscall.UTF16PtrFromString(absPath + "\x00")
	if err != nil {
		return err
	}

	shell32 := syscall.NewLazyDLL("shell32.dll")
	proc := shell32.NewProc("SHFileOperationW")
	op := shFileOpStruct{wFunc: 3, pFrom: pFrom, fFlags: 0x40 | 0x10 | 0x0400 | 0x0004}
	ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&op)))
	if ret != 0 {
		return fmt.Errorf("SHFileOperation failed with code %d", ret)
	}

	return nil
}
