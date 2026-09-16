// Package cmdagy — agy_ls.go handles listing Antigravity projects in table format.
package cmdagy

import (
	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var (
	agyLsOnlyMissing bool
	agyLsOnlyActive  bool
	agyLsOnlyPinned  bool
	agyLsJSON        bool
	agyLsSortBy      string
	agyLsFilter      string
)

var agyLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List projects in a status table",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyLs()
	},
}

func init() {
	agyLsCmd.Flags().BoolVarP(&agyLsOnlyMissing, "missing", "m", false, "Show only missing projects")
	agyLsCmd.Flags().BoolVarP(&agyLsOnlyActive, "active", "a", false, "Show only active projects")
	agyLsCmd.Flags().BoolVarP(&agyLsOnlyPinned, "pinned", "p", false, "Show only pinned projects")
	agyLsCmd.Flags().BoolVar(&agyLsJSON, "json", false, "Output results as JSON")
	agyLsCmd.Flags().StringVarP(&agyLsSortBy, "sort", "s", "name", "Sort by 'name' or 'time'")
	agyLsCmd.Flags().StringVarP(&agyLsFilter, "filter", "f", "", "Filter projects by name or path")
}

func runAgyLs() error {
	dirPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}

	return processAgyLsProjects(dirPath)
}

func processAgyLsProjects(dirPath string) error {
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load projects")
	}

	filtered := filterAndSortAgyProjects(projects)

	return renderAgyLsResult(filtered, dirPath)
}

func filterAndSortAgyProjects(projects []AgyProject) []AgyProject {
	filtered := filterAgyProjects(projects)
	sortAgyProjects(filtered, agyLsSortBy)

	return filtered
}

func renderAgyLsResult(filtered []AgyProject, dirPath string) error {
	if agyLsJSON {
		return outputAgyProjectsJSON(filtered)
	}

	renderAgyProjectsTable(filtered, dirPath)

	return nil
}
