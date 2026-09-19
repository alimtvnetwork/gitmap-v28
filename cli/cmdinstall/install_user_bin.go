package cmdinstall

import (
	"os"
	"path/filepath"
)

func getUserBinDir() string {
	home := homeDir()
	if home == "" {
		return ""
	}

	binDir := filepath.Join(home, ".local", "bin")
	_ = os.MkdirAll(binDir, 0755)
	ensureDirInPath(binDir)

	return binDir
}

func symlinkUserBin(name, source string) {
	binDir := getUserBinDir()
	if binDir == "" {
		return
	}

	dest := filepath.Join(binDir, name)
	if _, err := os.Lstat(dest); err == nil {
		return
	}

	_ = os.Symlink(source, dest)
}

func postPnpmLinuxSetup() {
	target := findBinaryInFallbackPaths("pnpm")
	if target == "" {
		return
	}

	ensureDirInPath(filepath.Dir(target))
	symlinkUserBin("pnpm", target)
}
