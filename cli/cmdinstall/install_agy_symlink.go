package cmdinstall

import (
	"os"
	"path/filepath"
	"runtime"
)

func ensureAgySymlinksAndPath() {
	if runtime.GOOS == "windows" {
		ensureAgyWindowsWrapper()
		return
	}
	ensureAgyUnixSymlink()
}

func ensureAgyUnixSymlink() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	localBin := filepath.Join(home, ".local", "bin")
	agyBin := filepath.Join(localBin, "agy")
	linkUnixAntigravity(localBin, agyBin)
	ensureDirInPath(localBin)
}

func linkUnixAntigravity(localBin, agyBin string) {
	info, errStat := os.Stat(agyBin)
	if errStat != nil || info.IsDir() {
		return
	}
	_ = os.Chmod(agyBin, 0755)
	antiBin := filepath.Join(localBin, "antigravity")
	if _, errAnti := os.Stat(antiBin); os.IsNotExist(errAnti) {
		_ = os.Symlink(agyBin, antiBin)
	}
}

func ensureAgyWindowsWrapper() {
	localApp := os.Getenv("LOCALAPPDATA")
	if localApp == "" {
		return
	}
	binDir := filepath.Join(localApp, "agy", "bin")
	cmdWrapper := filepath.Join(binDir, "antigravity.cmd")
	if _, err := os.Stat(cmdWrapper); os.IsNotExist(err) {
		content := "@echo off\r\n\"%~dp0agy.exe\" %*\r\n"
		_ = os.WriteFile(cmdWrapper, []byte(content), 0755)
	}
	ensureDirInPath(binDir)
}
