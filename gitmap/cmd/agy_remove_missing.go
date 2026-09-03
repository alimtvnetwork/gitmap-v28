// Package cmd — agy_remove_missing.go removes Antigravity projects whose paths are missing on disk.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var (
	agyRmMissingDryRun bool
	agyRmMissingYes    bool
	agyRmMissingExcept string
)

var agyRemoveMissingCmd = &cobra.Command{
	Use:     "remove-missing-projects",
	Aliases: []string{"remove-misisng-projects", "rm-missing-projects", "rm-missing", "clean-missing"},
	Short:   "Remove Antigravity projects whose target directories no longer exist on disk",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyRemoveMissing()
	},
}

func init() {
	agyRemoveMissingCmd.Flags().BoolVarP(&agyRmMissingDryRun, "dry-run", "d", false, "Preview projects to be removed without deleting")
	agyRemoveMissingCmd.Flags().BoolVarP(&agyRmMissingYes, "yes", "y", false, "Confirm removal without prompt")
	agyRemoveMissingCmd.Flags().StringVarP(&agyRmMissingExcept, "except", "e", "", "Exclude projects matching id, name, path, or from CSV file")
}

func runAgyRemoveMissing() error {
	dirPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load projects")
	}
	missingProjects := findMissingAgyProjects(projects, agyRmMissingExcept)
	if len(missingProjects) == 0 {
		fmt.Printf("%s No missing Antigravity projects found to remove. All paths exist.\n",
			constants.ColorGreen+"✓"+constants.ColorReset)
		return nil
	}
	return executeRemoveMissing(dirPath, missingProjects, len(projects))
}

func findMissingAgyProjects(projects []AgyProject, exceptStr string) []AgyProject {
	tokens := parseAgyExceptTokens(exceptStr)
	missing := make([]AgyProject, 0)
	for _, p := range projects {
		if p.ID == "outside-of-project" || isAgyProjectExceptedWithTokens(p, tokens) {
			continue
		}
		path := p.GetPath()
		if path != "" && !checkDirExists(path) {
			missing = append(missing, p)
		}
	}
	return missing
}

func executeRemoveMissing(dirPath string, targets []AgyProject, totalCount int) error {
	printMissingPreview(targets)
	if agyRmMissingDryRun {
		fmt.Printf("\n%s [dry-run] %d missing project(s) would be removed. Remaining active: %d\n",
			constants.ColorYellow+"ℹ"+constants.ColorReset, len(targets), totalCount-len(targets))
		return nil
	}
	if !agyRmMissingYes && !askMissingConfirmation(len(targets)) {
		fmt.Println("Removal canceled. No changes made.")
		return nil
	}
	return deleteMissingFiles(dirPath, targets, totalCount-len(targets))
}

func printMissingPreview(targets []AgyProject) {
	fmt.Printf("\n  %sTargeting %d missing Antigravity project(s) for removal:%s\n\n",
		constants.ColorYellow, len(targets), constants.ColorReset)
	fmt.Printf("    %-12s %-22s %-20s %s\n", "ID", "NAME", "SLUG", "MISSING PATH")
	fmt.Printf("    %s\n", strings.Repeat("─", 88))
	for _, p := range targets {
		slug := filepath.Base(p.GetPath())
		fmt.Printf("    %-12s %-22s %-20s %s\n", shortProjectId(p.ID), p.Name, slug, p.GetPath())
	}
	fmt.Printf("\n    %sTip: Exclude items using: --except \"<id, name, slug, or starts-with text>\"%s\n",
		constants.ColorDim, constants.ColorReset)
}

func askMissingConfirmation(count int) bool {
	fmt.Printf("\n  %sAre you sure you want to permanently remove %d missing project(s)? [y/N]: %s",
		constants.ColorYellow, count, constants.ColorReset)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(strings.ToLower(text))
	return text == "y" || text == "yes"
}

func deleteMissingFiles(dirPath string, targets []AgyProject, remainingCount int) error {
	deleted := 0
	for _, p := range targets {
		filePath := filepath.Join(dirPath, p.ID+".json")
		if err := os.Remove(filePath); err == nil {
			deleted++
		}
	}
	fmt.Printf("\n%s Successfully removed %d missing Antigravity project(s). Remaining: %d\n",
		constants.ColorGreen+"✓"+constants.ColorReset, deleted, remainingCount)
	return nil
}
