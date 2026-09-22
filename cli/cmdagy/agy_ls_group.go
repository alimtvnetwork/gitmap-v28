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

func renderAgyFolderTable(list []AgyProject, convMap map[string]AgyLatestConv, seqMap map[string]int) (int, int) {
	ctx := newAgyTableContext()
	activeCount, missingCount := 0, 0

	for _, p := range list {
		seq := resolveProjectSeq(p, seqMap)
		row := buildAgyTableRow(p, convMap, seq)
		if row.IsMissing {
			missingCount++
			continue
		}
		activeCount++
		ctx.addRow(row)
	}

	if len(ctx.Rows) == 0 {
		return activeCount, missingCount
	}

	printAgyTableHeader(ctx)
	for i, r := range ctx.Rows {
		printAgyTableRow(ctx, r, i)
	}
	fmt.Println()

	return activeCount, missingCount
}

func resolveProjectSeq(p AgyProject, seqMap map[string]int) int {
	if seq, ok := seqMap[p.ID]; ok && seq > 0 {
		return seq
	}
	seq, _ := GetOrAssignProjectSequence(p.ID, p.Name, p.GetPath())
	if seq > 0 {
		seqMap[p.ID] = seq
		return seq
	}

	return 1
}

func loadAllProjectSequences(projects []AgyProject) map[string]int {
	m, err := GetAllProjectSequences()
	if err != nil {
		m = make(map[string]int)
	}
	for _, p := range projects {
		if _, ok := m[p.ID]; ok {
			continue
		}
		seq, _ := GetOrAssignProjectSequence(p.ID, p.Name, p.GetPath())
		if seq > 0 {
			m[p.ID] = seq
		}
	}

	return m
}

func filterActiveProjects(list []AgyProject) []AgyProject {
	var active []AgyProject
	for _, p := range list {
		path := p.GetPath()
		isActive := path == "" || checkDirExists(path)
		if isActive {
			active = append(active, p)
		}
	}

	return active
}

func renderAgyProjectsGrouped(projects []AgyProject) (int, int) {
	convMap := loadLatestConversationsMap()
	seqMap := loadAllProjectSequences(projects)
	folders, groupMap := groupProjectsByRootFolder(projects)
	totalActive, totalMissing := 0, 0

	for _, folder := range folders {
		list := groupMap[folder]
		activeList := filterActiveProjects(list)
		if len(activeList) == 0 {
			totalMissing += len(list)
			continue
		}
		printFolderHeader(folder, len(activeList))
		active, missing := renderAgyFolderTable(list, convMap, seqMap)
		totalActive += active
		totalMissing += missing
	}

	return totalActive, totalMissing
}
