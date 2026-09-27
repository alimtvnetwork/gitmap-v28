package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

var (
	agyLsOnlyMissing bool
	agyLsOnlyActive  bool
	agyLsOnlyPinned  bool
	agyLsJSON        bool
	agyLsSSH         bool
	agyLsSortBy      string
	agyLsFilter      string
	agyLsFile        string
)

var agyLsCmd = &cobra.Command{
	Use:   "ls [N]",
	Short: "List projects in a status table",
	RunE: func(cmd *cobra.Command, args []string) error {
		n := 0
		n = parseArgCount(args, n)
		return runAgyLs(n)
	},
}

func init() {
	agyLsCmd.Flags().BoolVarP(&agyLsOnlyMissing, "missing", "m", false, "Show only missing projects")
	agyLsCmd.Flags().BoolVarP(&agyLsOnlyActive, "active", "a", false, "Show only active projects")
	agyLsCmd.Flags().BoolVarP(&agyLsOnlyPinned, "pinned", "p", false, "Show only pinned projects")
	agyLsCmd.Flags().BoolVar(&agyLsJSON, "json", false, "Output results as JSON")
	agyLsCmd.Flags().BoolVar(&agyLsSSH, "ssh", false, "Aggregate projects across remote SSH cluster nodes")
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
	projects, loadErr := DiscoverUnifiedAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load projects")
	}
	if agyLsSSH {
		remote := FetchClusterSSHProjects()
		projects = mergeRemoteProjects(projects, remote)
	}

	filtered := filterAndSortAgyProjects(projects)
	if n > 0 && len(filtered) > n {
		filtered = filtered[:n]
	}

	store.RecordAgyDecision("ls", fmt.Sprintf("%d projects", len(filtered)), dirPath, "", "listed projects across discovery sources", "success")
	return renderAgyLsResult(filtered, dirPath)
}

func mergeRemoteProjects(local, remote []AgyProject) []AgyProject {
	seen := make(map[string]bool)
	var merged []AgyProject
	for _, p := range local {
		ws := cleanProjectWorkspace(p.GetPath())
		if ws != "" {
			seen[ws] = true
		}
		merged = append(merged, p)
	}
	for _, p := range remote {
		ws := cleanProjectWorkspace(p.GetPath())
		if ws != "" && !seen[ws] {
			seen[ws] = true
			merged = append(merged, p)
		}
	}
	return merged
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
	entries := make([]CachedSequenceEntry, 0, len(filtered))
	for i, p := range filtered {
		entries = append(entries, CachedSequenceEntry{
			Seq:      i + 1,
			ID:       p.ID,
			Name:     p.Name,
			Path:     p.GetPath(),
			Category: "project",
		})
	}
	_ = SaveSequenceCache(entries)
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
