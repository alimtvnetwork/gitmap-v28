package movemerge

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies src -> dst, preserving mode bits and parent dirs.
// Symlinks are recreated, not followed.
func CopyFile(src, dst string, info os.FileInfo) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(dst), err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return copySymlink(src, dst)
	}

	return copyRegular(src, dst, info.Mode())
}

// copySymlink replicates a symlink at dst.
func copySymlink(src, dst string) error {
	target, err := os.Readlink(src)
	if err != nil {
		return fmt.Errorf("readlink %s: %w", src, err)
	}
	_ = os.Remove(dst)

	return os.Symlink(target, dst)
}

// copyRegular streams the file content with mode preservation.
func copyRegular(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open %s: %w", src, err)
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
	if err != nil {
		return fmt.Errorf("create %s: %w", dst, err)
	}
	defer out.Close()
	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("copy %s -> %s: %w", src, dst, err)
	}

	return nil
}

// CopyTree copies every file from src into dst, honoring opts ignore list.
func CopyTree(src, dst string, opts Options) (int, error) {
	idx, err := IndexTree(src, opts)
	if err != nil {
		return 0, err
	}
	count := 0
	for rel, meta := range idx {
		srcPath := filepath.Join(src, filepath.FromSlash(rel))
		dstPath := filepath.Join(dst, filepath.FromSlash(rel))
		if copyErr := CopyFile(srcPath, dstPath, meta.Info); copyErr != nil {
			return count, copyErr
		}
		count++
	}

	return count, nil
}
