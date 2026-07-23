package clonepick

// promote.go: move the picker's metadata-only clone into the user's
// final destination directory so we don't pay the network cost of
// `git clone` twice.
//
// Strategy:
//  1. Try os.Rename(src, dest). Fast path -- works when src and dest
//     live on the same filesystem AND dest is empty or removable.
//  2. If rename fails (cross-filesystem on Linux returns EXDEV; on
//     Windows it returns various errors when dest exists), fall back
//     to copy-tree + remove-src so the optimisation degrades safely
//     instead of failing the whole clone.
//
// prepareDest already MkdirAll'd dest and verified emptiness, so the
// rename path needs to remove the empty dest first (rename refuses
// to clobber an existing dir on Windows even when it's empty).

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// promotePreClonedSrc moves src's contents into dest. dest is
// expected to exist and be empty (prepareDest's contract). On
// success src no longer exists. On error the caller should clean
// up src; dest may be in a partially-populated state and the
// caller must treat the whole clone as failed.
func promotePreClonedSrc(src, dest string) error {
	if err := os.Remove(dest); err == nil {
		if renameErr := os.Rename(src, dest); renameErr == nil {
			return nil
		}
		// Recreate dest so the copy fallback has somewhere to land.
		if mkErr := os.MkdirAll(dest, 0o755); mkErr != nil {
			return mkErr
		}
	}

	return copyTreeThenRemove(src, dest)
}

// copyTreeThenRemove is the cross-filesystem fallback. Walks src,
// recreates the directory layout under dest, copies regular files
// (preserving mode bits), and removes src on success. Symlinks are
// recreated as symlinks; anything else (devices, sockets) is
// skipped because git never produces them.
func copyTreeThenRemove(src, dest string) error {
	walkErr := filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}

		return copyOneEntry(path, filepath.Join(dest, rel), entry)
	})
	if walkErr != nil {
		return walkErr
	}

	return os.RemoveAll(src)
}

// copyOneEntry handles one filesystem entry during the recursive
// copy. Split out so copyTreeThenRemove stays under the
// function-length cap.
func copyOneEntry(src, dest string, entry fs.DirEntry) error {
	if entry.IsDir() {
		return os.MkdirAll(dest, 0o755)
	}
	if entry.Type()&os.ModeSymlink != 0 {
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}

		return os.Symlink(target, dest)
	}
	if !entry.Type().IsRegular() {
		return nil
	}

	return copyRegularFile(src, dest)
}

// copyRegularFile copies a regular file preserving mode bits. Used
// by copyOneEntry; kept separate so the symlink branch isn't
// shadowed by file-IO scaffolding.
func copyRegularFile(src, dest string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return err
	}
	if _, copyErr := io.Copy(out, in); copyErr != nil {
		out.Close()

		return copyErr
	}

	return out.Close()
}

// errPromoteUnsupported is reserved for future use by callers that
// want to distinguish "rename failed because of a known incompat"
// from generic IO errors. Not currently emitted -- kept so the
// fallback path can grow typed errors without an API churn.
var errPromoteUnsupported = errors.New("clone-pick: promote unsupported")

var _ = errPromoteUnsupported
