// Package cmd — chromeprofile_undo.go handles undo and redo for Chrome profile mutations.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func getChromeBackupRoot() (string, *apperror.AppError) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", apperror.WrapSimple(homeErr, "user home dir")
	}

	root := filepath.Join(homeDir, ".gemini", "config", "backup", "chrome")
	if mkErr := os.MkdirAll(root, constants.DirPermission); mkErr != nil {
		return "", apperror.WrapSimple(mkErr, "mkdir chrome backup")
	}

	return root, nil
}

func snapshotChromeProfile(profileDir string, tag string) (string, *apperror.AppError) {
	backupRoot, rootErr := getChromeBackupRoot()
	if rootErr != nil {
		return "", rootErr
	}

	baseName := filepath.Base(profileDir)
	ts := time.Now().UTC().Format("20060102-150405")
	destDir := filepath.Join(backupRoot, fmt.Sprintf("%s-%s-%s", ts, baseName, tag))

	if copyErr := copyDirContents(profileDir, destDir); copyErr != nil {
		return "", copyErr
	}

	return destDir, nil
}

func listChromeSnapshots() ([]string, *apperror.AppError) {
	backupRoot, rootErr := getChromeBackupRoot()
	if rootErr != nil {
		return nil, rootErr
	}

	entries, readErr := os.ReadDir(backupRoot)
	if readErr != nil {
		return nil, nil
	}

	snapshots := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			snapshots = append(snapshots, filepath.Join(backupRoot, e.Name()))
		}
	}

	sort.Strings(snapshots)

	return snapshots, nil
}

func runChromeProfileUndo(profileName string) error {
	snapshots, err := listChromeSnapshots()
	if err != nil {
		return fmt.Errorf("list snapshots: %w", err.Unwrap())
	}

	chosen := pickLatestSnapshotWithTag(snapshots, "pre-import")
	if chosen == "" && len(snapshots) > 0 {
		chosen = snapshots[len(snapshots)-1]
	}

	if chosen == "" {
		fmt.Printf("%s No Chrome profile snapshots available to undo.\n", constants.ColorYellow+"ℹ"+constants.ColorReset)
		return nil
	}

	targetDir := filepath.Join(chromeUserDataDir(), profileName)
	if _, snapErr := snapshotChromeProfile(targetDir, "pre-undo"); snapErr != nil {
		return fmt.Errorf("snapshot pre-undo: %w", snapErr.Unwrap())
	}

	if copyErr := copyDirContents(chosen, targetDir); copyErr != nil {
		return fmt.Errorf("restore snapshot: %w", copyErr.Unwrap())
	}

	fmt.Printf("%s Restored Chrome profile %q from snapshot: %s\n", constants.ColorGreen+"✓"+constants.ColorReset, profileName, filepath.Base(chosen))
	return nil
}
