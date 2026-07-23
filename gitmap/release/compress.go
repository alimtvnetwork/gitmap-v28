// Package release — compress.go wraps release assets in archives.
// Windows assets → .zip, Linux/macOS assets → .tar.gz.
package release

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/verbose"
)

// CompressAssets wraps each file in assets into an archive.
// Windows binaries (.exe) → .zip, others → .tar.gz.
// Returns the list of archive paths (originals are removed).
func CompressAssets(assets []string) ([]string, error) {
	archives := make([]string, 0, len(assets))

	for _, path := range assets {
		archive, err := compressSingle(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrCompressFailed, filepath.Base(path), err)

			continue
		}

		if verbose.IsEnabled() {
			logCompressedArchive(archive)
		}

		archives = append(archives, archive)
	}

	return archives, nil
}

// logCompressedArchive logs the size and SHA-1 of a compressed archive.
func logCompressedArchive(archivePath string) {
	info, err := os.Stat(archivePath)
	if err != nil {
		verbose.Get().Log("compress: %s (stat error: %v)", filepath.Base(archivePath), err)

		return
	}

	hash, err := sha1File(archivePath)
	if err != nil {
		hash = "error"
	}

	verbose.Get().Log("compress: %s — %d bytes, sha1:%s",
		filepath.Base(archivePath), info.Size(), hash)
}

// compressSingle compresses a single file and removes the original.
func compressSingle(path string) (string, error) {
	if isWindowsBinary(path) {
		return createZip(path)
	}

	return createTarGz(path)
}

// isWindowsBinary returns true for .exe files.
func isWindowsBinary(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".exe")
}

// createZip wraps a file into a .zip archive.
func createZip(srcPath string) (string, error) {
	archivePath := srcPath + ".zip"
	outFile, err := os.Create(archivePath)
	if err != nil {
		return "", fmt.Errorf("create zip: %w", err)
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	defer w.Close()

	err = addFileToZip(w, srcPath)
	if err != nil {
		return "", err
	}

	w.Close()
	outFile.Close()

	os.Remove(srcPath)

	return archivePath, nil
}

// addFileToZip adds a single file entry to a zip writer.
func addFileToZip(w *zip.Writer, srcPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return fmt.Errorf("stat source: %w", err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("zip header: %w", err)
	}

	header.Name = filepath.Base(srcPath)
	header.Method = zip.Deflate

	writer, err := w.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("create entry: %w", err)
	}

	_, err = io.Copy(writer, src)

	return err
}

// DescribeCompression returns human-readable archive names for dry-run.
func DescribeCompression(assets []string) []string {
	var names []string

	for _, path := range assets {
		base := filepath.Base(path)
		if isWindowsBinary(path) {
			names = append(names, base+".zip")
		} else {
			names = append(names, base+".tar.gz")
		}
	}

	return names
}
