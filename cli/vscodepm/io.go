package vscodepm

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
)

// WriteEntries writes entries atomically to path.
func WriteEntries(path string, entries []Entry) error {
	return writeEntriesAtomic(path, entries)
}

// writeEntriesAtomic encodes entries to a sibling .tmp then renames.
// On Windows, os.Rename overwrites the destination; on Unix it does too,
// but we explicitly remove the destination first if rename fails to keep
// behavior consistent.
func writeEntriesAtomic(path string, entries []Entry) error {
	tmpPath := path + constants.VSCodePMProjectsTempSuffix

	if err := writeEntriesToFile(tmpPath, entries); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)

		return fmt.Errorf(constants.ErrVSCodePMRenameFailed,
			filepath.Base(path), err)
	}

	return nil
}

// writeEntriesToFile serializes entries with tab indent + trailing newline.
func writeEntriesToFile(path string, entries []Entry) error {
	if err := os.MkdirAll(filepath.Dir(path), constants.DirPermission); err != nil {
		return fmt.Errorf(constants.ErrVSCodePMWriteTempFailed, path, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf(constants.ErrVSCodePMWriteTempFailed, path, err)
	}

	if err := encodeEntries(file, entries); err != nil {
		_ = file.Close()
		_ = os.Remove(path)

		return fmt.Errorf(constants.ErrVSCodePMWriteTempFailed, path, err)
	}

	return file.Close()
}

// encodeEntries writes entries as pretty JSON with a trailing newline.
func encodeEntries(w io.Writer, entries []Entry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", constants.VSCodePMJSONIndent)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(entries); err != nil {
		return err
	}

	return nil
}

// normalizePath returns the canonical key used for rootPath comparisons.
// Delegates to fsutil.CanonicalPathKey for OS-aware canonicalization.
func normalizePath(p string) string {
	return fsutil.CanonicalPathKey(p)
}

// pathsEqual compares two paths using fsutil.EqualPaths.
func pathsEqual(a, b string) bool {
	return fsutil.EqualPaths(a, b)
}
