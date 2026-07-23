// Package localdirs handles migration of legacy repo-local directories.
package localdirs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// MigrateLegacyDirs moves old directories into .gitmap/ if found.
func MigrateLegacyDirs() {
	migrations := []struct{ old, sub string }{
		{constants.LegacyOutputDir, constants.OutputDirName},
		{constants.LegacyReleaseDir, constants.ReleaseDirName},
		{constants.LegacyDeployedDir, constants.DeployedDirName},
	}

	for _, m := range migrations {
		migrateSingleDir(m.old, m.sub)
	}
}

// migrateSingleDir moves a legacy directory to .gitmap/<sub> if it exists.
// If the target already exists, it merges files from the legacy directory
// (without overwriting) and removes the legacy directory.
func migrateSingleDir(oldDir, subDir string) {
	if !dirExists(oldDir) {
		return
	}

	target := filepath.Join(constants.GitMapDir, subDir)
	if dirExists(target) {
		mergeAndRemoveLegacy(oldDir, target)

		return
	}

	ensureGitMapDir()
	if err := os.Rename(oldDir, target); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrMigrationFailed, oldDir, err)

		return
	}

	fmt.Printf(constants.MsgMigrated, oldDir, target)
}

// mergeAndRemoveLegacy copies files from oldDir into target (skipping files
// that already exist in target), then removes the legacy directory.
func mergeAndRemoveLegacy(oldDir, target string) {
	var merged, skipped int

	walkErr := filepath.WalkDir(oldDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Walk error at %s: %v\n", path, err)

			return nil
		}

		rel, relErr := filepath.Rel(oldDir, path)
		if relErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not compute relative path for %s: %v\n", path, relErr)

			return nil
		}

		if rel == "." {
			return nil
		}

		dest := filepath.Join(target, rel)

		if d.IsDir() {
			if mkErr := os.MkdirAll(dest, constants.DirPermission); mkErr != nil {
				fmt.Fprintf(os.Stderr, "  ⚠ Could not create directory %s: %v\n", dest, mkErr)
			}

			return nil
		}

		if fileExists(dest) {
			skipped++

			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not read %s: %v\n", path, readErr)

			return nil
		}

		if writeErr := os.WriteFile(dest, data, constants.FilePermission); writeErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not write %s: %v\n", dest, writeErr)

			return nil
		}

		merged++

		return nil
	})

	if walkErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Walk failed for %s: %v\n", oldDir, walkErr)
	}

	if err := os.RemoveAll(oldDir); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrMigrationFailed, oldDir, err)

		return
	}

	fmt.Printf(constants.MsgMergedAndRemoved, oldDir, target, merged, skipped)
}

// dirExists checks if a directory exists at the given path.
func dirExists(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

// fileExists checks if a file exists at the given path.
func fileExists(path string) bool {
	info, err := os.Stat(path)

	return err == nil && !info.IsDir()
}

// ensureGitMapDir creates the .gitmap/ directory if it does not exist.
func ensureGitMapDir() {
	if err := os.MkdirAll(constants.GitMapDir, constants.DirPermission); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not create %s: %v\n", constants.GitMapDir, err)
	}
}
