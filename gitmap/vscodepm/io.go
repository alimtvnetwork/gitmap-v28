package vscodepm

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

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
// Converts backslashes to forward slashes so cross-platform path representations
// (e.g. D:\work\my-app vs D:/work/my-app) compare consistently across all operating systems.
// Case-insensitive on Windows or for Windows-style drive paths, case-sensitive elsewhere.
func normalizePath(p string) string {
	cleaned := filepath.Clean(filepath.ToSlash(p))
	if runtime.GOOS == "windows" || isWindowsDrivePath(cleaned) {
		return strings.ToLower(cleaned)
	}

	return cleaned
}

func isWindowsDrivePath(p string) bool {
	return len(p) >= 2 && p[1] == ':' && ((p[0] >= 'a' && p[0] <= 'z') || (p[0] >= 'A' && p[0] <= 'Z'))
}

// pathsEqual compares two paths using normalizePath.
func pathsEqual(a, b string) bool {
	return normalizePath(a) == normalizePath(b)
}
