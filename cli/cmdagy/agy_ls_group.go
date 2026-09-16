// Package cmdagy — agy_ls_group.go groups projects by root folder for table display.
package cmdagy

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func resolveProjectParentFolder(p AgyProject) string {
	path := strings.TrimSpace(p.GetPath())
	if path == "" || path == "—" {
		return "[Global / Config]"
	}

	clean := filepath.Clean(path)
	parent := filepath.Dir(clean)
	if parent == "." || parent == "/" || parent == "\\" {
		return clean
	}

	return parent
}

func groupProjectsByRootFolder(projects []AgyProject) ([]string, map[string][]AgyProject) {
	groupMap := make(map[string][]AgyProject)
	for _, p := range projects {
		parent := resolveProjectParentFolder(p)
		groupMap[parent] = append(groupMap[parent], p)
	}

	folders := make([]string, 0, len(groupMap))
	for f := range groupMap {
		folders = append(folders, f)
	}
	sort.Strings(folders)

	return folders, groupMap
}

func printFolderHeader(folder string, count int) {
	fmt.Printf("  %s📁 %s%s %s(%d projects)%s\n",
		constants.ColorCyan, folder, constants.ColorReset,
		constants.ColorDim, count, constants.ColorReset)
}

func renderAgyFolderTable(list []AgyProject) (int, int) {
	ctx := newAgyTableContext()
	activeCount, missingCount := 0, 0

	for _, p := range list {
		row := buildAgyTableRow(p)
		if row.IsMissing {
			missingCount++
		} else {
			activeCount++
		}
		ctx.addRow(row)
	}

	printAgyTableHeader(ctx)
	for i, r := range ctx.Rows {
		printAgyTableRow(ctx, r, i)
	}
	fmt.Println()

	return activeCount, missingCount
}

func renderAgyProjectsGrouped(projects []AgyProject) (int, int) {
	folders, groupMap := groupProjectsByRootFolder(projects)
	totalActive, totalMissing := 0, 0

	for _, folder := range folders {
		list := groupMap[folder]
		printFolderHeader(folder, len(list))
		active, missing := renderAgyFolderTable(list)
		totalActive += active
		totalMissing += missing
	}

	return totalActive, totalMissing
}
