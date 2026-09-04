// Package cmd — agy_ls.go handles listing Antigravity projects in table format.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
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
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load projects")
	}
	filtered := filterAgyProjects(projects)
	sortAgyProjects(filtered, agyLsSortBy)

	if agyLsJSON {
		return outputAgyProjectsJSON(filtered)
	}
	renderAgyProjectsTable(filtered, dirPath)
	return nil
}

func loadAllAgyProjects(dirPath string) ([]AgyProject, error) {
	_ = os.MkdirAll(dirPath, 0755)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	projects := make([]AgyProject, 0, len(entries))
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		p, parseErr := readAgyProject(filepath.Join(dirPath, entry.Name()))
		if parseErr == nil {
			projects = append(projects, p)
		}
	}
	return projects, nil
}

func readAgyProject(filePath string) (AgyProject, error) {
	var p AgyProject
	data, err := os.ReadFile(filePath)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(data, &p)
	return p, err
}

func filterAgyProjects(projects []AgyProject) []AgyProject {
	pinnedMap := getPinnedMapIfRequested()
	out := make([]AgyProject, 0, len(projects))

	for _, p := range projects {
		if agyLsOnlyPinned && !pinnedMap[p.ID] && !pinnedMap[filepath.Clean(p.GetPath())] {
			continue
		}

		if !matchesAgyFilter(p) {
			continue
		}

		out = append(out, p)
	}

	return out
}

func getPinnedMapIfRequested() map[string]bool {
	if !agyLsOnlyPinned {
		return make(map[string]bool)
	}

	return getPinnedProjectsMap()
}

func matchesAgyFilter(p AgyProject) bool {
	path := p.GetPath()
	isMissing := path != "" && !checkDirExists(path)

	if agyLsOnlyMissing && !isMissing {
		return false
	}
	if agyLsOnlyActive && isMissing {
		return false
	}
	if agyLsFilter != "" && !matchesFilter(p.Name, path, agyLsFilter) {
		return false
	}
	return true
}

func matchesFilter(name, path, filter string) bool {
	term := strings.ToLower(filter)
	matchesName := strings.Contains(strings.ToLower(name), term)
	matchesPath := strings.Contains(strings.ToLower(path), term)
	return matchesName || matchesPath
}

func sortAgyProjects(projects []AgyProject, sortBy string) {
	if sortBy == "time" || sortBy == "recent" {
		sort.Slice(projects, func(i, j int) bool {
			return projects[i].UpdatedAt > projects[j].UpdatedAt
		})
		return
	}
	sort.Slice(projects, func(i, j int) bool {
		return strings.ToLower(projects[i].Name) < strings.ToLower(projects[j].Name)
	})
}

func renderAgyProjectsTable(projects []AgyProject, dirPath string) {
	ctx := newAgyTableContext()
	activeCount, missingCount := 0, 0

	for _, p := range projects {
		row := buildAgyTableRow(p)
		if row.IsMissing {
			missingCount++
		} else {
			activeCount++
		}
		ctx.addRow(row)
	}

	printAgyBanner(len(projects), dirPath)
	printAgyTableHeader(ctx)
	for i, r := range ctx.Rows {
		printAgyTableRow(ctx, r, i)
	}
	printAgySummary(len(projects), activeCount, missingCount)
}

func buildAgyTableRow(p AgyProject) agyTableRow {
	path := p.GetPath()
	branch := p.GetBranch()
	if branch == "" {
		branch = "—"
	}
	isMissing := path != "" && !checkDirExists(path)
	status := "active"
	if path == "" {
		status = "global"
		path = "—"
	}
	return agyTableRow{
		ID:        shortProjectId(p.ID),
		Name:      p.Name,
		Branch:    branch,
		Status:    status,
		Updated:   formatRelativeTime(p.UpdatedAt),
		Path:      path,
		IsMissing: isMissing,
	}
}

func checkDirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func outputAgyProjectsJSON(projects []AgyProject) error {
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "json marshal")
	}
	fmt.Println(string(data))
	return nil
}
