// Package cmdagy — agy_pins_runners.go handles execution logic for pinned projects commands.
package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runAgyPinProjectsLs() error {
	dirPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}

	if err := processAgyLsProjects(dirPath, 5); err != nil {
		return err
	}
	printPinsTableAdvisory()
	return nil
}

func printPinsTableAdvisory() {
	store, _ := loadPinnedProjectsStore()
	if store == nil || len(store.Projects) == 0 {
		fmt.Printf("  %sTip: Pin projects with: gitmap agy pins add <seq|id|slug>%s\n\n",
			constants.ColorDim, constants.ColorReset)
		return
	}
	fmt.Printf("  %s%d pinned project(s) shown first. Manage pins: gitmap agy pins <add|rm|edit|help>%s\n\n",
		constants.ColorDim, len(store.Projects), constants.ColorReset)
}

func runAgyPinProjectsAdd(args []string) error {
	targets := expandCommaSeparatedTokens(args)
	if len(targets) == 0 {
		targets = []string{"."}
	}

	for _, t := range targets {
		pinned, addErr := addPinnedProjectTarget(t)
		if addErr != nil {
			return fmt.Errorf("pin project %q: %w", t, addErr.Unwrap())
		}

		fmt.Printf("%s Pinned project: %s%s%s (%s)\n",
			constants.ColorGreen+"✓"+constants.ColorReset,
			constants.ColorCyan, pinned.Name, constants.ColorReset,
			pinned.Path)
	}

	return nil
}

func runAgyPinProjectsRm(args []string) error {
	if agyPinProjectsRmAll {
		return handleAgyPinProjectsClearAll()
	}

	tokens := expandCommaSeparatedTokens(args)
	if len(tokens) == 0 {
		return fmt.Errorf("requires project sequence (001), ID, slug, or --all")
	}

	return handleAgyPinProjectsRemoveTargets(tokens)
}

func handleAgyPinProjectsClearAll() error {
	count, clearErr := clearAllPinnedProjects()
	if clearErr != nil {
		return fmt.Errorf("clear pinned projects: %w", clearErr.Unwrap())
	}

	fmt.Printf("%s Unpinned all (%d) projects.\n", constants.ColorGreen+"✓"+constants.ColorReset, count)

	return nil
}

func handleAgyPinProjectsRemoveTargets(targets []string) error {
	for _, t := range targets {
		removed, rmErr := removePinnedProjectTarget(t)
		if rmErr != nil {
			return fmt.Errorf("unpin project %q: %w", t, rmErr.Unwrap())
		}

		fmt.Printf("%s Unpinned project: %s%s%s (%s)\n",
			constants.ColorGreen+"✓"+constants.ColorReset,
			constants.ColorCyan, removed.Name, constants.ColorReset,
			removed.Path)
	}

	return nil
}
