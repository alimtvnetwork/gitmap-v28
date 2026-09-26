package cmdagy

import (
	"encoding/json"
	"os"
	"strconv"
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
	agyLsFile        string
)

var agyLsCmd = &cobra.Command{
	Use:   "ls [N]",
	Short: "List projects in a status table",
	RunE: func(cmd *cobra.Command, args []string) error {
		n := 8
		n = parseArgCount(args, n)
		return runAgyLs(n)
	},
}

func init() {
	agyLsCmd.Flags().BoolVarP(&agyLsOnlyMissing, "missing", "m", false, "Show only missing projects")
	agyLsCmd.Flags().BoolVarP(&agyLsOnlyActive, "active", "a", false, "Show only active projects")
	agyLsCmd.Flags().BoolVarP(&agyLsOnlyPinned, "pinned", "p", false, "Show only pinned projects")
	agyLsCmd.Flags().BoolVar(&agyLsJSON, "json", false, "Output results as JSON")
	agyLsCmd.Flags().StringVarP(&agyLsSortBy, "sort", "s", "name", "Sort by 'name' or 'time'")
	agyLsCmd.Flags().StringVarP(&agyLsFilter, "filter", "", "", "Filter projects by name or path")
	agyLsCmd.Flags().StringVarP(&agyLsFile, "file", "f", "", "Export JSON to file")
}

func runAgyLs(n int) error {
	dirPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}

	return processAgyLsProjects(dirPath, n)
}

func processAgyLsProjects(dirPath string, n int) error {
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load projects")
	}

	filtered := filterAndSortAgyProjects(projects)
	if len(filtered) > n {
		filtered = filtered[:n]
	}

	return renderAgyLsResult(filtered, dirPath)
}

func filterAndSortAgyProjects(projects []AgyProject) []AgyProject {
	filtered := filterAgyProjects(projects)
	sortAgyProjects(filtered, agyLsSortBy)

	return filtered
}

func renderAgyLsResult(filtered []AgyProject, dirPath string) error {
	if agyLsFile != "" {
		return outputAgyProjectsJSONFile(filtered, agyLsFile)
	}
	if agyLsJSON {
		return outputAgyProjectsJSON(filtered)
	}

	renderAgyProjectsTable(filtered, dirPath)
	return nil
}

func parseArgCount(args []string, defaultN int) int {
	if len(args) == 0 {
		return defaultN
	}
	val, err := strconv.Atoi(args[0])
	if err == nil && val > 0 {
		return val
	}
	return defaultN
}

func outputAgyProjectsJSONFile(projects []AgyProject, file string) error {
	data, _ := json.MarshalIndent(projects, "", "  ")
	return os.WriteFile(file, data, 0644)
}
