package cmd

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// extractDocsSiteZip extracts docs-site.zip into the target directory.
// Validates paths to prevent traversal (G305) and limits total size (G110).
func extractDocsSiteZip(zipPath, targetDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("resolve target dir: %w", err)
	}

	var totalSize int64

	for _, f := range r.File {
		written, entryErr := extractDocsZipEntry(f, absTarget, totalSize)
		if entryErr != nil {
			return entryErr
		}

		totalSize += written
		if totalSize >= maxDocsSiteSize {
			return fmt.Errorf("archive exceeds maximum extraction size (%d bytes)", maxDocsSiteSize)
		}
	}

	return nil
}

// extractDocsZipEntry writes a single zip entry to absTarget, validating the path
// against traversal and respecting the remaining size budget.
func extractDocsZipEntry(f *zip.File, absTarget string, totalSize int64) (int64, error) {
	destPath := filepath.Join(absTarget, f.Name) // #nosec G305 — validated below
	absDestPath, absErr := filepath.Abs(destPath)
	if absErr != nil || !strings.HasPrefix(absDestPath, absTarget+string(os.PathSeparator)) {
		return 0, fmt.Errorf("illegal file path in zip: %s", f.Name)
	}

	if f.FileInfo().IsDir() {
		if mkErr := os.MkdirAll(absDestPath, constants.DirPermission); mkErr != nil {
			return 0, fmt.Errorf("create dir %s: %w", absDestPath, mkErr)
		}
		return 0, nil
	}

	if mkErr := os.MkdirAll(filepath.Dir(absDestPath), constants.DirPermission); mkErr != nil {
		return 0, fmt.Errorf("create parent dir: %w", mkErr)
	}

	rc, openErr := f.Open()
	if openErr != nil {
		return 0, fmt.Errorf("open entry %s: %w", f.Name, openErr)
	}
	defer rc.Close()

	outFile, createErr := os.Create(absDestPath) // #nosec G304 — absDestPath validated above
	if createErr != nil {
		return 0, fmt.Errorf("create file %s: %w", absDestPath, createErr)
	}

	written, copyErr := io.CopyN(outFile, rc, maxDocsSiteSize-totalSize) // #nosec G110 — size-limited
	outFile.Close()

	if copyErr != nil && !errors.Is(copyErr, io.EOF) {
		return written, fmt.Errorf("write file %s: %w", absDestPath, copyErr)
	}

	return written, nil
}
