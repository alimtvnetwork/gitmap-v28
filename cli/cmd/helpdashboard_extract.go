package cmd

import (
	"archive/zip"
	"errors"
	"fmt"
t"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// extractDocsSiteZip extracts docs-site.zip into the target directory.
// Validates paths to prevent traversal (G305) and limits total size (G110).
func extractDocsSiteZip(zipPath, targetDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return apperror.WrapSimple(err, "open zip")
	}

	defer r.Close()

	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return apperror.WrapSimple(err, "resolve target dir")
	}

	var totalSize int64

	for _, f := range r.File {
		written, entryErr := extractDocsZipEntry(f, absTarget, totalSize)
		if entryErr != nil {
			return entryErr
		}

		totalSize += written
		if totalSize >= maxDocsSiteSize {
			return apperror.NewValidationError(fmt.Sprintf("archive exceeds maximum extraction size (%d bytes)", maxDocsSiteSize))
		}
	}

	return nil
}

// extractDocsZipEntry writes a single zip entry to absTarget, validating the path
// against traversal and respecting the remaining size budget.
func extractDocsZipEntry(f *zip.File, absTarget string, totalSize int64) (int64, error) {
	destPath := filepath.Join(absTarget, f.Name) // #nosec G305 — validated below
	absDestPath, absErr := filepath.Abs(destPath)
	if absErr != nil {
		return 0, apperror.WrapSimple(absErr, fmt.Sprintf("illegal file path in zip %s", f.Name))
	}
	if !strings.HasPrefix(absDestPath, absTarget+string(os.PathSeparator)) {
		return 0, apperror.NewValidationError(fmt.Sprintf("illegal file path in zip: %s", f.Name))
	}

	if f.FileInfo().IsDir() {
		return handleExtractDir(absDestPath)
	}

	if mkErr := os.MkdirAll(filepath.Dir(absDestPath), constants.DirPermission); mkErr != nil {
		return 0, apperror.WrapSimple(mkErr, "create parent dir")
	}

	rc, openErr := f.Open()
	if openErr != nil {
		return 0, apperror.WrapSimple(openErr, fmt.Sprintf("open entry %s", f.Name))
	}

	defer rc.Close()

	outFile, createErr := os.Create(absDestPath) // #nosec G304 — absDestPath validated above
	if createErr != nil {
		return 0, apperror.WrapSimple(createErr, fmt.Sprintf("create file %s", absDestPath))
	}

	written, copyErr := io.CopyN(outFile, rc, maxDocsSiteSize-totalSize) // #nosec G110 — size-limited
	outFile.Close()

	if copyErr != nil && !errors.Is(copyErr, io.EOF) {
		return written, apperror.WrapSimple(copyErr, fmt.Sprintf("write file %s", absDestPath))
	}

	return written, nil
}

// handleExtractDir creates a directory and returns an error if it fails.
func handleExtractDir(absDestPath string) (int64, error) {
	mkErr := os.MkdirAll(absDestPath, constants.DirPermission)
	if mkErr != nil {
		return 0, apperror.WrapSimple(mkErr, fmt.Sprintf("create dir %s", absDestPath))
	}

	return 0, nil
}
