// Package release — compresstar.go handles tar.gz archive creation.
package release

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// createTarGz wraps a file into a .tar.gz archive.
func createTarGz(srcPath string) (string, error) {
	archivePath := srcPath + ".tar.gz"
	outFile, err := os.Create(archivePath)
	if err != nil {
		return "", fmt.Errorf("create tar.gz: %w", err)
	}
	defer outFile.Close()

	gw := gzip.NewWriter(outFile)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	err = addFileToTar(tw, srcPath)
	if err != nil {
		return "", err
	}

	tw.Close()
	gw.Close()
	outFile.Close()

	os.Remove(srcPath)

	return archivePath, nil
}

// addFileToTar adds a single file entry to a tar writer.
func addFileToTar(tw *tar.Writer, srcPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return fmt.Errorf("stat source: %w", err)
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("tar header: %w", err)
	}

	header.Name = filepath.Base(srcPath)

	err = tw.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	_, err = io.Copy(tw, src)

	return err
}
