// Package cmdagy — agy_rm_ops.go implements execution logic for agy rm.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runAgyRmEnhanced(args []string, folderFlag string) error {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return apperror.WrapSimple(err, "get projects dir path")
	}

	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load agy projects")
	}
	sortAgyProjects(projects, "name")

	targets, resolveErr := resolveRemovalTargets(args, folderFlag, projects)
	if resolveErr != nil {
		return resolveErr
	}

	return executeAgyProjectRemoval(targets)
}

func resolveRemovalTargets(args []string, folderFlag string, projects []AgyProject) ([]AgyProject, error) {
	if folderFlag != "" {
		return ResolveAgyFolderTargets(folderFlag, projects)
	}

	if len(args) >= 2 && (strings.ToLower(args[0]) == "folder" || strings.ToLower(args[0]) == "dir") {
		return ResolveAgyFolderTargets(args[1], projects)
	}

	if len(args) == 0 {
		return nil, apperror.NewSimple("requires target sequence, id, or slug (or --folder)", "E9000")
	}

	return ResolveAgyProjectTargets(args, projects)
}

func executeAgyProjectRemoval(targets []AgyProject) error {
	_, _ = snapshotAgyProjects("pre-rm")
	removedCount, lastErr := removeTargetProjects(targets)
	printAgyUndoGuidance()
	if removedCount == 0 && lastErr != nil {
		return apperror.WrapSimple(lastErr, "execute agy project removal")
	}
	return nil
}

func removeTargetProjects(targets []AgyProject) (int, error) {
	count := 0
	var lastErr error
	for _, p := range targets {
		if err := removeSingleAgyProject(p); err != nil {
			lastErr = err
			continue
		}
		count++
	}
	return count, lastErr
}

func removeSingleAgyProject(p AgyProject) error {
	inv, _ := json.Marshal(p)
	_ = recordAgyTask("rm", p.Name, p.ID, string(inv))
	if err := deleteProjectFile(p.ID); err != nil {
		fmt.Printf("  %s✗ Failed to remove project '%s': %v%s\n",
			constants.ColorRed, p.Name, err, constants.ColorReset)
		return err
	}
	fmt.Printf("  %s✓ Removed project '%s' (%s)%s\n",
		constants.ColorGreen, p.Name, shortProjectId(p.ID), constants.ColorReset)
	return nil
}

func completeAgyRmArgs(toComplete string) ([]string, cobra.ShellCompDirective) {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	sortAgyProjects(projects, "name")

	return CompleteAgyProjectSuggestions(toComplete, projects), cobra.ShellCompDirectiveNoFileComp
}
