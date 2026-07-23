package release

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/verbose"
)

// buildDocsSiteAsset bundles docs-site/dist/ into a zip archive for release.
// Returns the archive path if successful, or empty string if docs-site is not found.
func buildDocsSiteAsset(stagingDir string) string {
	docsDir := constants.HDDocsDir
	distDir := filepath.Join(docsDir, constants.HDDistDir)

	info, err := os.Stat(distDir)
	if err != nil || !info.IsDir() {
		if verbose.IsEnabled() {
			verbose.Get().Log("docs-site: no dist/ directory found at %s, skipping", distDir)
		}

		return ""
	}

	archivePath := filepath.Join(stagingDir, constants.DocsSiteArchive)
	fmt.Printf(constants.MsgDocsSiteBundling, distDir)

	items, err := collectDocsSiteItems(distDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDocsSiteBundle, err)

		return ""
	}

	if err := createDocsSiteZip(archivePath, distDir, items); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDocsSiteBundle, err)

		return ""
	}

	if verbose.IsEnabled() {
		archiveInfo, statErr := os.Stat(archivePath)
		if statErr == nil {
			verbose.Get().Log("docs-site: bundled %d file(s) into %s (%d bytes)",
				len(items), filepath.Base(archivePath), archiveInfo.Size())
		}
	}

	fmt.Printf(constants.MsgDocsSiteBundled, filepath.Base(archivePath))

	return archivePath
}

// collectDocsSiteItems walks the dist directory and returns relative paths.
func collectDocsSiteItems(distDir string) ([]string, error) {
	var items []string

	err := filepath.Walk(distDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error at %s: %w", path, err)
		}

		if info.IsDir() {
			return nil
		}

		items = append(items, path)

		return nil
	})

	return items, err
}

// createDocsSiteZip creates a zip with docs-site/ prefix so it extracts correctly.
func createDocsSiteZip(archivePath, distDir string, items []string) error {
	outFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("create archive %s: %w", archivePath, err)
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	defer w.Close()

	for _, itemPath := range items {
		relPath, relErr := filepath.Rel(distDir, itemPath)
		if relErr != nil {
			return fmt.Errorf("relative path for %s: %w", itemPath, relErr)
		}

		// Prefix with docs-site/ so it extracts into the right directory.
		entryName := filepath.ToSlash(filepath.Join(constants.HDDocsDir, constants.HDDistDir, relPath))

		if err := addSingleFileToZip(w, itemPath, entryName); err != nil {
			return fmt.Errorf("add %s to zip: %w", itemPath, err)
		}
	}

	return nil
}
