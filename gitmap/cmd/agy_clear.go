// Package cmd — agy_clear.go handles clearing stale Antigravity projects and caches safely.
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
	agyClearMissing bool
	agyClearAll     bool
	agyClearDryRun  bool
	agyClearYes     bool
	agyClearExcept  string
)

var agyClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear stale or missing Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyClear()
	},
}

func init() {
	agyClearCmd.Flags().BoolVarP(&agyClearMissing, "missing", "m", false, "Clear only projects whose target directories no longer exist")
	agyClearCmd.Flags().BoolVarP(&agyClearAll, "all", "a", false, "Clear all projects (requires confirmation)")
	agyClearCmd.Flags().BoolVarP(&agyClearDryRun, "dry-run", "d", false, "Preview items that would be cleared without deleting")
	agyClearCmd.Flags().BoolVarP(&agyClearYes, "yes", "y", false, "Confirm deletion without interactive prompt")
	agyClearCmd.Flags().StringVarP(&agyClearExcept, "except", "e", "", "Exclude projects matching id, name, slug, or path")
}

func runAgyClear() error {
	dirPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load projects")
	}
	targets := selectClearTargets(projects)
	if len(targets) == 0 {
		fmt.Printf("%s No stale or missing Antigravity projects found to clear.\n", constants.ColorGreen+"✓"+constants.ColorReset)
		return nil
	}
	return executeClearTargets(dirPath, targets)
}

func selectClearTargets(projects []AgyProject) []AgyProject {
	targets := make([]AgyProject, 0)
	for _, p := range projects {
		if p.ID == "outside-of-project" || isAgyProjectExcepted(p, agyClearExcept) {
			continue
		}
		path := p.GetPath()
		isMissing := path != "" && !checkDirExists(path)
		shouldClear := agyClearAll || isMissing
		if shouldClear {
			targets = append(targets, p)
		}
	}
	return targets
}

func isAgyProjectExcepted(p AgyProject, exceptStr string) bool {
	if exceptStr == "" {
		return false
	}
	for _, ex := range strings.Split(exceptStr, ",") {
		ex = strings.TrimSpace(ex)
		if isMatchExcept(p, ex) {
			return true
		}
	}
	return false
}

func isMatchExcept(p AgyProject, ex string) bool {
	if ex == "" {
		return false
	}
	normPath := filepath.ToSlash(filepath.Clean(p.GetPath()))
	normEx := filepath.ToSlash(filepath.Clean(ex))
	if strings.EqualFold(p.ID, ex) || strings.HasPrefix(strings.ToLower(p.ID), strings.ToLower(ex)) || strings.EqualFold(p.Name, ex) || strings.EqualFold(normPath, normEx) {
		return true
	}
	return strings.Contains(strings.ToLower(normPath), strings.ToLower(normEx)) || strings.Contains(strings.ToLower(p.Name), strings.ToLower(ex))
}

func executeClearTargets(dirPath string, targets []AgyProject) error {
	printClearPreview(targets)
	if agyClearDryRun {
		fmt.Printf("\n%s [dry-run] %d project(s) would be removed. No changes made.\n", constants.ColorYellow+"ℹ"+constants.ColorReset, len(targets))
		return nil
	}
	if !agyClearYes && !askClearConfirmation(len(targets)) {
		fmt.Println("Aborted. No projects were removed.")
		return nil
	}
	deletedCount := deleteTargetProjects(dirPath, targets)
	fmt.Printf("\n%s Successfully removed %d stale project(s).\n", constants.ColorGreen+"✓"+constants.ColorReset, deletedCount)
	return nil
}

func printClearPreview(targets []AgyProject) {
	fmt.Printf("\nFound %d project(s) to remove:\n", len(targets))
	for i, p := range targets {
		path := p.GetPath()
		if path == "" {
			path = "(no path)"
		}
		fmt.Printf("  %2d. %-20s (%s) -> %s\n", i+1, p.Name, shortProjectId(p.ID), path)
	}
}

func askClearConfirmation(count int) bool {
	fmt.Printf("\nAre you sure you want to remove these %d project(s)? [y/N]: ", count)
	reader := bufio.NewReader(os.Stdin)
	ans, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	ans = strings.TrimSpace(strings.ToLower(ans))
	return ans == "y" || ans == "yes"
}

func deleteTargetProjects(dirPath string, targets []AgyProject) int {
	count := 0
	for _, p := range targets {
		fileName := p.ID + ".json"
		filePath := filepath.Join(dirPath, fileName)
		if err := os.Remove(filePath); err == nil {
			count++
		}
	}
	return count
}
