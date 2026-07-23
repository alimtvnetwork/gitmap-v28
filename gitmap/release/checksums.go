// Package release — checksums.go generates SHA256 checksums for release assets.
package release

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/verbose"
)

// GenerateChecksums computes SHA256 hashes for all assets and writes
// a checksums.txt file in the same directory as the first asset.
// Returns the path to the checksums file.
func GenerateChecksums(assets []string) (string, error) {
	if len(assets) == 0 {
		return "", nil
	}

	dir := filepath.Dir(assets[0])
	checksumPath := filepath.Join(dir, constants.ChecksumsFile)

	file, err := os.Create(checksumPath)
	if err != nil {
		return "", fmt.Errorf("create checksums file: %w", err)
	}
	defer file.Close()

	for _, asset := range assets {
		hash, hashErr := hashFile(asset)
		if hashErr != nil {
			fmt.Fprintf(os.Stderr, constants.ErrChecksumFailed, filepath.Base(asset), hashErr)

			continue
		}

		if verbose.IsEnabled() {
			verbose.Get().Log("checksum: %s  sha256:%s", filepath.Base(asset), hash)
		}

		fmt.Fprintf(file, "%s  %s\n", hash, filepath.Base(asset))
	}

	return checksumPath, nil
}

// hashFile computes the SHA256 hex digest of a file.
func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
