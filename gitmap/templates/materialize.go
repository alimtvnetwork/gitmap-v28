package templates

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Materialize copies every embedded asset into the user-overlay directory,
// SKIPPING any file that already exists. This makes the call idempotent and
// preserves user edits.
//
// Returns the overlay directory and the list of files actually written.
func Materialize() (string, []string, error) {
	dir, err := EnsureUserDir()
	if err != nil {
		return "", nil, err
	}

	var written []string
	walkErr := fs.WalkDir(FS, embedAssetsRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		w, copyErr := materializeOne(dir, p)
		if copyErr != nil {
			return copyErr
		}
		if len(w) > 0 {
			written = append(written, w)
		}

		return nil
	})
	if walkErr != nil {
		return dir, written, fmt.Errorf(errTemplateMaterialize, dir, walkErr)
	}

	return dir, written, nil
}

// materializeOne copies a single embedded file to the overlay if missing.
// Returns the destination path on a fresh write, "" when skipped.
func materializeOne(overlayDir, embedPath string) (string, error) {
	rel := strings.TrimPrefix(embedPath, embedAssetsRoot+"/")
	dst := filepath.Join(overlayDir, filepath.FromSlash(rel))

	if _, err := os.Stat(dst); err == nil {
		return "", nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return "", fmt.Errorf(errTemplateMaterialize, dst, err)
	}

	data, err := FS.ReadFile(embedPath)
	if err != nil {
		return "", fmt.Errorf(errTemplateRead, embedPath, err)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return "", fmt.Errorf(errTemplateMaterialize, dst, err)
	}

	return dst, nil
}
