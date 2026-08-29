// Package cmd — chromeprofile_copy.go: resilient Chrome profile tree copy
// helpers used by `gitmap chrome-profile-copy` / `gitmap cpc`.
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// copyChromeProfile copies the curated subset of entries from src to
// dst. Missing entries are skipped silently — Chrome regenerates them.
func copyChromeProfile(src, dst string) (int, error) {
	if err := os.MkdirAll(dst, constants.DirPermission); err != nil {
		return 0, newChromeProfileCopyError(src, dst, constants.ChromeProfileCopyOpMkdir, err)
	}
	return copyChromeProfileEntries(src, dst)
}

func copyChromeProfileEntries(src, dst string) (int, error) {
	total := 0
	for _, name := range constants.ChromeProfileCopyEntries {
		n, err := copyEntry(filepath.Join(src, name), filepath.Join(dst, name))
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

type chromeProfileCopyError struct {
	Source string
	Target string
	Op     string
	Err    error
}

func (e *chromeProfileCopyError) Error() string {
	return fmt.Sprintf("%s %s -> %s: %v", e.Op, e.Source, e.Target, e.Err)
}

func (e *chromeProfileCopyError) Unwrap() error { return e.Err }

func newChromeProfileCopyError(src, dst, op string, err error) error {
	return &chromeProfileCopyError{Source: src, Target: dst, Op: op, Err: err}
}

// copyEntry copies a single file or directory tree. Returns file count.
func copyEntry(src, dst string) (int, error) {
	info, err := os.Stat(src)
	if err != nil {
		return handleCopyStatError(src, dst, err)
	}
	if info.IsDir() {
		return copyDir(src, dst)
	}
	return copyRegularFile(src, dst)
}

func handleCopyStatError(src, dst string, err error) (int, error) {
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	return 0, newChromeProfileCopyError(src, dst, constants.ChromeProfileCopyOpStat, err)
}

func copyRegularFile(src, dst string) (int, error) {
	copied, err := chromeProfileCopyFile(src, dst)
	if !copied || err != nil {
		return 0, err
	}
	return 1, nil
}

// chromeProfileCopyFile copies a single file from src to dst preserving mode.
func chromeProfileCopyFile(src, dst string) (bool, error) {
	in, err := os.Open(src) //nolint:gosec // curated entry list
	if err != nil {
		return handleChromeFileOpenError(src, dst, err)
	}
	defer in.Close()
	return writeChromeProfileCopyFile(in, src, dst)
}

func handleChromeFileOpenError(src, dst string, err error) (bool, error) {
	if isChromeVolatileLockFile(src) {
		warnChromeProfileLockSkip(src, dst, err)
		return false, nil
	}
	return false, newChromeProfileCopyError(src, dst, constants.ChromeProfileCopyOpRead, err)
}

func writeChromeProfileCopyFile(in *os.File, src, dst string) (bool, error) {
	if err := os.MkdirAll(filepath.Dir(dst), constants.DirPermission); err != nil {
		return false, newChromeProfileCopyError(src, dst, constants.ChromeProfileCopyOpMkdir, err)
	}
	out, err := os.Create(dst) //nolint:gosec // curated entry list
	if err != nil {
		return false, newChromeProfileCopyError(src, dst, constants.ChromeProfileCopyOpWrite, err)
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return handleChromeFileCopyError(src, dst, err)
	}
	return true, nil
}

func handleChromeFileCopyError(src, dst string, err error) (bool, error) {
	if isChromeVolatileLockFile(src) {
		warnChromeProfileLockSkip(src, dst, err)
		return false, nil
	}
	return false, newChromeProfileCopyError(src, dst, constants.ChromeProfileCopyOpWrite, err)
}

// copyDir recursively copies a directory tree.
func copyDir(src, dst string) (int, error) {
	if err := os.MkdirAll(dst, constants.DirPermission); err != nil {
		return 0, newChromeProfileCopyError(src, dst, constants.ChromeProfileCopyOpMkdir, err)
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return 0, newChromeProfileCopyError(src, dst, constants.ChromeProfileCopyOpList, err)
	}
	return copyDirEntries(entries, src, dst)
}

func copyDirEntries(entries []os.DirEntry, src, dst string) (int, error) {
	total := 0
	for _, e := range entries {
		n, err := copyEntry(filepath.Join(src, e.Name()), filepath.Join(dst, e.Name()))
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func isChromeVolatileLockFile(path string) bool {
	return filepath.Base(path) == constants.ChromeProfileLockFileName
}

// chromeProfileLockSkipCount is reset at the start of each copy run and
// incremented for every volatile LOCK file the copier skips. The caller
// prints a single colorful summary line at the end instead of one banner
// per file (see runChromeProfileCopy).
var chromeProfileLockSkipCount int

func warnChromeProfileLockSkip(src, dst string, err error) {
	_ = dst
	_ = err
	chromeProfileLockSkipCount++
	fmt.Fprintf(os.Stderr, constants.WarnChromeProfileSkipLock, src)
}

func printChromeProfileCopyError(src, dst chromeProfileResolution, err error) {
	copyErr := unwrapChromeProfileCopyError(err)
	fmt.Fprintf(os.Stderr, constants.ErrChromeProfileCopyFailed,
		chromeProfileSummary(src), chromeProfileSummary(dst), src.Path, dst.Path,
		copyErr.Source, copyErr.Op, copyErr.Err)
}

func unwrapChromeProfileCopyError(err error) chromeProfileCopyError {
	var copyErr *chromeProfileCopyError
	if errors.As(err, &copyErr) {
		return *copyErr
	}
	return chromeProfileCopyError{Source: constants.ChromeProfileCopyUnknown, Op: constants.ChromeProfileCopyOpCopy, Err: err}
}
