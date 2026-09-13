package cmdinstall

import (
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/archive"
	"github.com/alimtvnetwork/gitmap-v28/cli/tempdir"
)

func extractArchiveStaging(ctx context.Context, srcPath, appName string) (string, error) {
	stagingBase := tempdir.RepoTempDir("archive-install")
	if isStandaloneGz(srcPath) {
		return extractSingleGz(srcPath, stagingBase, appName)
	}

	res, err := archive.CompactExtract(ctx, srcPath, stagingBase)
	if err != nil {
		return "", apperror.Wrap(err, "archive.extract", map[string]any{"source": srcPath})
	}
	return res.OutputDir, nil
}

func isStandaloneGz(path string) bool {
	low := strings.ToLower(path)
	return strings.HasSuffix(low, ".gz") && !strings.HasSuffix(low, ".tar.gz")
}

func extractSingleGz(srcPath, stagingBase, appName string) (string, error) {
	destDir := filepath.Join(stagingBase, appName)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}

	destFile := filepath.Join(destDir, appName)
	if err := decompressGzToFile(srcPath, destFile); err != nil {
		return "", err
	}
	_ = os.Chmod(destFile, 0755)
	return destDir, nil
}

func decompressGzToFile(srcPath, destPath string) error {
	in, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer in.Close()

	gzReader, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, gzReader)
	return err
}

func deriveArchiveAppName(archivePath string) string {
	base := filepath.Base(archivePath)
	for _, ext := range archiveExts {
		if strings.HasSuffix(strings.ToLower(base), ext) {
			return base[:len(base)-len(ext)]
		}
	}
	return base
}
