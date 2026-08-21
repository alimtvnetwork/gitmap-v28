package fsutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// SafeRemoveAll removes path using Windows long-path prefix if needed.
func SafeRemoveAll(path string) error {
	if len(path) == 0 {
		return nil
	}
	longPath := EnsureLongPath(path)
	err := os.RemoveAll(longPath)
	if err != nil {
		return os.RemoveAll(path)
	}
	return nil
}

// SafeRename renames src to dst, falling back to copy+delete.
func SafeRename(src, dst string) error {
	err := os.Rename(EnsureLongPath(src), EnsureLongPath(dst))
	if err == nil {
		return nil
	}
	if errCopy := CopyDirectory(src, dst); errCopy != nil {
		return fmt.Errorf("rename fallback failed: %w", errCopy)
	}
	return SafeRemoveAll(src)
}

// CopyDirectory recursively copies src tree to dst.
func CopyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
