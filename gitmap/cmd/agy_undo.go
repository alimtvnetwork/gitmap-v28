// Package cmd — agy_undo.go handles undo and redo for Antigravity workspace mutations.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var agyUndoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Revert the last Antigravity clear or project mutation",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyUndo()
	},
}

var agyRedoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Reapply the undone Antigravity workspace action",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyRedo()
	},
}

func getAgyBackupRoot() (string, *apperror.AppError) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", apperror.WrapSimple(homeErr, "user home dir")
	}

	root := filepath.Join(homeDir, ".gemini", "config", "backup", "agy")
	if mkErr := os.MkdirAll(root, constants.DirPermission); mkErr != nil {
		return "", apperror.WrapSimple(mkErr, "mkdir agy backup")
	}

	return root, nil
}

func snapshotAgyProjects(tag string) (string, *apperror.AppError) {
	backupRoot, rootErr := getAgyBackupRoot()
	if rootErr != nil {
		return "", rootErr
	}

	srcDir, dirErr := getProjectsDirPath()
	if dirErr != nil {
		return "", apperror.WrapSimple(dirErr, "projects dir")
	}

	ts := time.Now().UTC().Format("20060102-150405")
	destDir := filepath.Join(backupRoot, fmt.Sprintf("%s-%s", ts, tag))

	if copyErr := copyDirContents(srcDir, destDir); copyErr != nil {
		return "", copyErr
	}

	return destDir, nil
}

func copyDirContents(src string, dest string) *apperror.AppError {
	if mkErr := os.MkdirAll(dest, constants.DirPermission); mkErr != nil {
		return apperror.WrapSimple(mkErr, "mkdir dest")
	}

	entries, readErr := os.ReadDir(src)
	if readErr != nil {
		return nil
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		sFile := filepath.Join(src, e.Name())
		dFile := filepath.Join(dest, e.Name())

		content, rErr := os.ReadFile(sFile)
		if rErr != nil {
			continue
		}

		_ = os.WriteFile(dFile, content, constants.FilePermission)
	}

	return nil
}

func listAgySnapshots() ([]string, *apperror.AppError) {
	backupRoot, rootErr := getAgyBackupRoot()
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
