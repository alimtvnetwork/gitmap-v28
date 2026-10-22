// Package cmd — agy_undo_ops.go implements execute logic for AGY undo and redo.
package cmdagy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runAgyUndo() error {
	restoredConvs, hasCacheBackup, _ := restoreLatestCacheClearBackup()
	if hasCacheBackup && restoredConvs > 0 {
		fmt.Printf("%s Restored %d conversation(s) from temporary cache backup.\n",
			constants.ColorGreen+"✓"+constants.ColorReset, restoredConvs)
	}

	snapshots, err := listAgySnapshots()
	if err != nil {
		return apperror.WrapSimple(err.Unwrap(), "list snapshots")
	}

	chosen := pickLatestSnapshotWithTag(snapshots, "pre-clear")
	if chosen == "" && len(snapshots) > 0 {
		chosen = snapshots[len(snapshots)-1]
	}

	if chosen == "" && (!hasCacheBackup || restoredConvs == 0) {
		fmt.Printf("%s No Antigravity snapshots available to undo.\n", constants.ColorYellow+"ℹ"+constants.ColorReset)
		return nil
	}
	if chosen == "" {
		return nil
	}

	destDir, dirErr := getProjectsDirPath()
	if dirErr != nil {
		return apperror.WrapSimple(dirErr, "projects dir")
	}

	if _, snapErr := snapshotAgyProjects("pre-undo"); snapErr != nil {
		return apperror.WrapSimple(snapErr.Unwrap(), "snapshot pre-undo")
	}

	if copyErr := copyDirContents(chosen, destDir); copyErr != nil {
		return apperror.WrapSimple(copyErr.Unwrap(), "restore snapshot")
	}

	fmt.Printf("%s Restored Antigravity projects from snapshot: %s\n", constants.ColorGreen+"✓"+constants.ColorReset, filepath.Base(chosen))

	return nil
}

func runAgyRedo() error {
	snapshots, err := listAgySnapshots()
	if err != nil {
		return apperror.WrapSimple(err.Unwrap(), "list snapshots")
	}

	chosen := pickLatestSnapshotWithTag(snapshots, "pre-undo")
	if chosen == "" {
		fmt.Printf("%s No Antigravity redo snapshots available.\n", constants.ColorYellow+"ℹ"+constants.ColorReset)

		return nil
	}

	destDir, dirErr := getProjectsDirPath()
	if dirErr != nil {
		return apperror.WrapSimple(dirErr, "projects dir")
	}

	if copyErr := copyDirContents(chosen, destDir); copyErr != nil {
		return apperror.WrapSimple(copyErr.Unwrap(), "restore redo snapshot")
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
