// Package cmd — chromeprofile_redo.go implements redo and group operations for Chrome profiles.
package cmdchromeprofile

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/ecosystemgroup"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
)

func runChromeProfileRedo(profileName string) error {
	snapshots, err := listChromeSnapshots()
	if err != nil {
		return fmt.Errorf("list snapshots: %w", err.Unwrap())
	}

	chosen := fsutil.PickLatestSnapshotWithTag(snapshots, "pre-undo")
	if chosen == "" {
		fmt.Printf("%s No Chrome profile redo snapshots available.\n", constants.ColorYellow+"ℹ"+constants.ColorReset)

		return nil
	}

	targetDir := filepath.Join(chromeUserDataDir(), profileName)
	if copyErr := fsutil.CopyDirContents(chosen, targetDir); copyErr != nil {
		return fmt.Errorf("restore redo snapshot: %w", copyErr.Unwrap())
	}

	fmt.Printf("%s Redone Chrome profile %q from snapshot: %s\n", constants.ColorGreen+"✓"+constants.ColorReset, profileName, filepath.Base(chosen))

	return nil
}

func runChromeGroupDispatch(args []string) error {
	if len(args) == 0 || args[0] == "ls" || args[0] == "list" {
		return ecosystemgroup.PrintGroupList("chrome", "Chrome Profile")
	}

	if args[0] == "add" && len(args) >= 3 {
		return handleChromeGroupAdd(args[1], args[2:])
	}

	if (args[0] == "rm" || args[0] == "delete") && len(args) >= 2 {
		return handleChromeGroupRm(args[1:])
	}

	return fmt.Errorf("usage: gitmap chrome group [ls|add <group> <profiles...>|rm <group> [profile]]")
}

func handleChromeGroupAdd(group string, profiles []string) error {
	appErr := ecosystemgroup.AddGroup("chrome", group, "", profiles)
	if appErr != nil {
		return fmt.Errorf("add chrome group: %w", appErr.Unwrap())
	}

	fmt.Printf("%s Added %d profile(s) to Chrome group %q.\n", constants.ColorGreen+"✓"+constants.ColorReset, len(profiles), group)

	return nil
}

func handleChromeGroupRm(args []string) error {
	if len(args) == 1 {
		return handleChromeGroupDelete(args[0])
	}

	if err := ecosystemgroup.RemoveTarget("chrome", args[0], args[1]); err != nil {
		return fmt.Errorf("remove chrome target: %w", err.Unwrap())
	}

	fmt.Printf("%s Removed %q from Chrome group %q.\n", constants.ColorGreen+"✓"+constants.ColorReset, args[1], args[0])

	return nil
}

func handleChromeGroupDelete(group string) error {
	if err := ecosystemgroup.DeleteGroup("chrome", group); err != nil {
		return fmt.Errorf("delete chrome group: %w", err.Unwrap())
	}

	fmt.Printf("%s Deleted Chrome group %q.\n", constants.ColorGreen+"✓"+constants.ColorReset, group)

	return nil
}
