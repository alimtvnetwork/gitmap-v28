package cmd

import (
	"fmt"
	"io"
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

func installUserBin(name, source string) (string, error) {
	binDir := getUserBinDir()
	if binDir == "" {
		return "", fmt.Errorf("could not determine user home directory")
	}

	dest := filepath.Join(binDir, name)
	if err := copyExecutable(source, dest); err != nil {
		return "", err
	}

	_ = os.Chmod(dest, 0755)

	return dest, nil
}

func copyExecutable(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}

	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	defer out.Close()
	_, err = io.Copy(out, in)

	return err
}

func postPnpmLinuxSetup() {
	target := findBinaryInFallbackPaths("pnpm")
	if target == "" {
		return
	}

	ensureDirInPath(filepath.Dir(target))
	symlinkUserBin("pnpm", target)
}
