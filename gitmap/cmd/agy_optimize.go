// Package cmd — agy_optimize.go implements project deduplication for Antigravity.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var (
	agyOptDryRun bool
	agyOptYes    bool
	agyOptExcept string
)

var agyOptimizeCmd = &cobra.Command{
	Use:     "optimize-projects",
	Aliases: []string{"optimize", "dedupe", "dedup", "--repeat-fix"},
	Short:   "Deduplicate and clear repeated Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyOptimize()
	},
}

func init() {
	agyOptimizeCmd.Flags().BoolVarP(&agyOptDryRun, "dry-run", "d", false, "Preview duplicates without deleting")
	agyOptimizeCmd.Flags().BoolVarP(&agyOptYes, "yes", "y", false, "Confirm removal without prompt")
	agyOptimizeCmd.Flags().StringVarP(&agyOptExcept, "except", "e", "", "Exclude projects matching id, name, slug, or path")
	agyOptimizeCmd.Flags().Bool("repeat-fix", false, "Alias flag for optimize-projects")
}

func runAgyOptimize() error {
	dirPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load projects")
	}
	duplicates := findAgyDuplicates(projects, agyOptExcept)
	if len(duplicates) == 0 {
		fmt.Printf("%s No duplicate Antigravity projects found. Total active: %d\n",
			constants.ColorGreen+"✓"+constants.ColorReset, len(projects))
		return nil
	}
	return executeAgyOptimize(dirPath, duplicates, len(projects)-len(duplicates))
}

func findAgyDuplicates(projects []AgyProject, exceptStr string) []AgyProject {
	groups := groupProjectsByPath(projects)
	duplicates := make([]AgyProject, 0)
	for _, group := range groups {
		if len(group) <= 1 {
			continue
		}
		sortGroupNewestFirst(group)
		for _, dup := range group[1:] {
			if !isAgyProjectExcepted(dup, exceptStr) {
				duplicates = append(duplicates, dup)
			}
		}
	}
	return duplicates
}

func groupProjectsByPath(projects []AgyProject) map[string][]AgyProject {
	groups := make(map[string][]AgyProject)
	for _, p := range projects {
		if p.ID == "outside-of-project" {
			continue
		}
		path := strings.ToLower(filepath.Clean(p.GetPath()))
		if path == "" {
			continue
		}
		groups[path] = append(groups[path], p)
	}
	return groups
}

func sortGroupNewestFirst(group []AgyProject) {
	sort.Slice(group, func(i, j int) bool {
		return group[i].UpdatedAt > group[j].UpdatedAt
	})
}

func executeAgyOptimize(dirPath string, duplicates []AgyProject, remainingCount int) error {
	printClearPreview(duplicates)
	if agyOptDryRun {
		fmt.Printf("\n%s [dry-run] %d duplicate project(s) would be removed. Remaining: %d\n",
			constants.ColorYellow+"ℹ"+constants.ColorReset, len(duplicates), remainingCount)
		return nil
	}
	if !agyOptYes && !askClearConfirmation(len(duplicates)) {
		fmt.Println("Optimization cancelled. No changes made.")
		return nil
	}
	return deleteAgyDuplicates(dirPath, duplicates, remainingCount)
}

func deleteAgyDuplicates(dirPath string, duplicates []AgyProject, remainingCount int) error {
	deleted := 0
	for _, p := range duplicates {
		filePath := filepath.Join(dirPath, p.ID+".json")
		if err := os.Remove(filePath); err == nil {
			deleted++
		}
	}
	fmt.Printf("\n%s Successfully removed %d duplicate Antigravity project(s). Remaining: %d\n",
		constants.ColorGreen+"✓"+constants.ColorReset, deleted, remainingCount)
	return nil
}
