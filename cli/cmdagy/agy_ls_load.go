// Package cmdagy — agy_ls_load.go handles reading, sorting, and JSON output of Antigravity projects.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func loadAllAgyProjects(dirPath string) ([]AgyProject, error) {
	_ = os.MkdirAll(dirPath, 0755)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	return collectAgyProjectsFromEntries(dirPath, entries), nil
}

func collectAgyProjectsFromEntries(dirPath string, entries []os.DirEntry) []AgyProject {
	projects := make([]AgyProject, 0, len(entries))
	for _, entry := range entries {
		p, isLoaded := tryLoadAgyProjectEntry(dirPath, entry)
		if isLoaded {
			projects = append(projects, p)
		}
	}

	return projects
}

func tryLoadAgyProjectEntry(dirPath string, entry os.DirEntry) (AgyProject, bool) {
	if filepath.Ext(entry.Name()) != ".json" {
		return AgyProject{}, false
	}

	p, err := readAgyProject(filepath.Join(dirPath, entry.Name()))

	return p, err == nil
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

func sortAgyProjects(projects []AgyProject, sortBy string) {
	pinnedMap := getPinnedProjectsSet()
	sort.Slice(projects, func(i, j int) bool {
		pi := pinnedMap[projects[i].ID] || pinnedMap[projects[i].Name]
		pj := pinnedMap[projects[j].ID] || pinnedMap[projects[j].Name]
		if pi != pj {
			return pi
		}
		if sortBy == "time" || sortBy == "recent" {
			return projects[i].UpdatedAt > projects[j].UpdatedAt
		}
		return strings.ToLower(projects[i].Name) < strings.ToLower(projects[j].Name)
	})
}

func outputAgyProjectsJSON(projects []AgyProject) error {
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "json marshal")
	}

	fmt.Println(string(data))

	return nil
}
