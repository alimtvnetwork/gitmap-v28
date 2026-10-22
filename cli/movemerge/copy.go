package movemerge

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// CopyFile copies src -> dst, preserving mode bits and parent dirs.
// Symlinks are recreated, not followed.
func CopyFile(src, dst string, info os.FileInfo) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return apperror.WrapSimple(err, "mkdir "+filepath.Dir(dst))
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
		return apperror.WrapSimple(err, "readlink "+src)
	}

	_ = os.Remove(dst)

	return os.Symlink(target, dst)
}

// copyRegular streams the file content with mode preservation.
func copyRegular(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return apperror.WrapSimple(err, "open "+src)
	}

	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
	if err != nil {
		return apperror.WrapSimple(err, "create "+dst)
	}

	defer out.Close()
	if _, err = io.Copy(out, in); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("copy %s -> %s", src, dst))
	}

	return nil
}

// CopyTree copies every file from src into dst, honoring opts ignore list.
func CopyTree(src, dst string, opts Options) (int, error) {
	idxRes := IndexTree(src, opts)
	if idxRes.IsFailure() {
		return 0, idxRes.AppError()
	}

	count := 0
	for rel, meta := range idxRes.Data {
		srcPath := filepath.Join(src, filepath.FromSlash(rel))
		dstPath := filepath.Join(dst, filepath.FromSlash(rel))
		if copyErr := CopyFile(srcPath, dstPath, meta.Info); copyErr != nil {
			return count, copyErr
		}

		count++
	}

	return count, nil
}
