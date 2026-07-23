package release

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/verbose"
)

// createMaxCompressZip creates a ZIP archive with Deflate level 9.
func createMaxCompressZip(archivePath string, items []model.ZipGroupItem) error {
	outFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("create zip: %w", err)
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	defer w.Close()

	// Register a custom compressor with max compression.
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return newMaxDeflateWriter(out), nil
	})

	fileCount := 0

	for _, item := range items {
		itemPath := item.FullPath
		if len(itemPath) == 0 {
			itemPath = item.Path
		}

		if item.IsFolder {
			err = addFolderToZip(w, itemPath)
		} else {
			err = addSingleFileToZip(w, itemPath, filepath.Base(itemPath))
			fileCount++
		}

		if err != nil {
			return err
		}
	}

	// Close writer before reading size/hash.
	w.Close()
	outFile.Close()

	if verbose.IsEnabled() {
		logArchiveSummary(archivePath, fileCount)
	}

	return nil
}

// logArchiveSummary logs the final archive size, file count, and SHA-1 hash.
func logArchiveSummary(archivePath string, fileCount int) {
	info, err := os.Stat(archivePath)
	if err != nil {
		verbose.Get().Log("zip-summary: %s (stat error: %v)", filepath.Base(archivePath), err)

		return
	}

	hash, err := sha1File(archivePath)
	if err != nil {
		hash = "error"
	}

	verbose.Get().Log("zip-summary: %s — %d file(s), %d bytes, sha1:%s",
		filepath.Base(archivePath), fileCount, info.Size(), hash)
}

// sha1File computes the SHA-1 hex digest of a file.
func sha1File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()

	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// addSingleFileToZip adds one file entry to a zip writer.
func addSingleFileToZip(w *zip.Writer, srcPath, entryName string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", srcPath, err)
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return fmt.Errorf("stat %s: %w", srcPath, err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("header %s: %w", srcPath, err)
	}

	header.Name = entryName
	header.Method = zip.Deflate

	writer, err := w.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("create entry %s: %w", entryName, err)
	}

	_, err = io.Copy(writer, src)

	return err
}

// addFolderToZip recursively adds a directory's contents to the archive.
func addFolderToZip(w *zip.Writer, folderPath string) error {
	return filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, relErr := filepath.Rel(filepath.Dir(folderPath), path)
		if relErr != nil {
			relPath = path
		}

		if verbose.IsEnabled() {
			verbose.Get().Log("  zip-add: %s → %s", path, filepath.ToSlash(relPath))
		}

		return addSingleFileToZip(w, path, filepath.ToSlash(relPath))
	})
}
