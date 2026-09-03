// Package cmd — agy_undo_ops.go implements execute logic for AGY undo and redo.
package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runAgyUndo() error {
	snapshots, err := listAgySnapshots()
	if err != nil {
		return fmt.Errorf("list snapshots: %w", err.Unwrap())
	}

	chosen := pickLatestSnapshotWithTag(snapshots, "pre-clear")
	if chosen == "" && len(snapshots) > 0 {
		chosen = snapshots[len(snapshots)-1]
	}

	if chosen == "" {
		fmt.Printf("%s No Antigravity snapshots available to undo.\n", constants.ColorYellow+"ℹ"+constants.ColorReset)
		return nil
	}

	destDir, dirErr := getProjectsDirPath()
	if dirErr != nil {
		return fmt.Errorf("projects dir: %w", dirErr)
	}

	if _, snapErr := snapshotAgyProjects("pre-undo"); snapErr != nil {
		return fmt.Errorf("snapshot pre-undo: %w", snapErr.Unwrap())
	}

	if copyErr := copyDirContents(chosen, destDir); copyErr != nil {
		return fmt.Errorf("restore snapshot: %w", copyErr.Unwrap())
	}

	fmt.Printf("%s Restored Antigravity projects from snapshot: %s\n", constants.ColorGreen+"✓"+constants.ColorReset, filepath.Base(chosen))
	return nil
}

func runAgyRedo() error {
	snapshots, err := listAgySnapshots()
	if err != nil {
		return fmt.Errorf("list snapshots: %w", err.Unwrap())
	}

	chosen := pickLatestSnapshotWithTag(snapshots, "pre-undo")
	if chosen == "" {
		fmt.Printf("%s No Antigravity redo snapshots available.\n", constants.ColorYellow+"ℹ"+constants.ColorReset)
		return nil
	}

	destDir, dirErr := getProjectsDirPath()
	if dirErr != nil {
		return fmt.Errorf("projects dir: %w", dirErr)
	}

	if copyErr := copyDirContents(chosen, destDir); copyErr != nil {
		return fmt.Errorf("restore redo snapshot: %w", copyErr.Unwrap())
	}

	fmt.Printf("%s Redone Antigravity projects from snapshot: %s\n", constants.ColorGreen+"✓"+constants.ColorReset, filepath.Base(chosen))
	return nil
}

func pickLatestSnapshotWithTag(snapshots []string, tag string) string {
	for i := len(snapshots) - 1; i >= 0; i-- {
		if strings.HasSuffix(snapshots[i], tag) {
			return snapshots[i]
		}
	}

	return ""
}
