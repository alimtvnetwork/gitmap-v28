package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runDBLs displays an architectural overview of all Gitmap SQLite databases.
func runDBLs(args []string) error {
	mainInfo, exists := collectMainDBInfo()
	splitDBs := collectSplitDBs()
	profileDBs := collectProfileDBs()

	printDBHeader()
	printMainDBSection(mainInfo, exists)
	printSplitDBSection(splitDBs)
	printProfileDBSection(profileDBs)
	printDBSummary(mainInfo, splitDBs, profileDBs)
	return nil
}

func printDBHeader() {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + "── Gitmap Database Architecture ──" + constants.ColorReset)
	fmt.Println()
}

func printMainDBSection(m DBFileInfo, exists bool) {
	fmt.Println("  " + constants.ColorCyan + "● Primary Master Database:" + constants.ColorReset)
	fmt.Printf("    %-13s %s%s%s\n", "Name:", constants.ColorWhite, m.Name, constants.ColorReset)
	fmt.Printf("    %-13s %s%s%s\n", "Type:", constants.ColorDim, m.Category, constants.ColorReset)
	fmt.Printf("    %-13s %s\n", "Location:", m.Path)
	if exists {
		fmt.Printf("    %-13s %s%s%s\n", "Size:", constants.ColorGreen, formatBytes(m.Size), constants.ColorReset)
	} else {
		fmt.Printf("    %-13s %s(not yet initialized - run 'gitmap scan')%s\n", "Size:", constants.ColorYellow, constants.ColorReset)
	}
	printPurposeWrapped("Purpose:", m.Purpose)
	fmt.Println()
}

func printPurposeWrapped(label, purpose string) {
	fmt.Printf("    %-13s %s\n", label, purpose)
}

func printSplitDBSection(splitDBs []DBFileInfo) {
	fmt.Println("  " + constants.ColorCyan + "● Split Repository Databases (repo_search):" + constants.ColorReset)
	dirs := findSplitDBDirs()
	dirStr := "(none)"
	if len(dirs) > 0 {
		dirStr = dirs[0]
	}
	fmt.Printf("    %-13s %s\n", "Directory:", dirStr)
	fmt.Printf("    %-13s %s%d split database(s)%s\n", "Total:", constants.ColorWhite, len(splitDBs), constants.ColorReset)
	if len(splitDBs) == 0 {
		fmt.Println("    " + constants.ColorDim + "(Split DBs are created automatically when searching or indexing a repo)" + constants.ColorReset)
	} else {
		printSplitDBRows(splitDBs)
	}
	printSplitDBWhy()
	fmt.Println()
}

func printSplitDBRows(splitDBs []DBFileInfo) {
	fmt.Println()
	fmt.Printf("    %-8s %-28s %-26s %-10s\n", "REPO ID", "REPO SLUG", "DB FILE", "SIZE")
	fmt.Println("    " + strings.Repeat("─", 76))
	for _, s := range splitDBs {
		idStr := "-"
		if s.RepoID > 0 {
			idStr = fmt.Sprintf("%d", s.RepoID)
		}
		slug := truncateStr(s.RepoSlug, 27)
		fmt.Printf("    %-8s %-28s %-26s %-10s\n", idStr, slug, truncateStr(s.Name, 25), formatBytes(s.Size))
	}
}

func truncateStr(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

func printSplitDBWhy() {
	fmt.Println()
	fmt.Println("    " + constants.ColorYellow + "Why Split DBs?" + constants.ColorReset)
	fmt.Println("    Split DBs isolate per-repository search caches, full-text file indexes, and sequence")
	fmt.Println("    operations into separate SQLite files under repo_search/. This eliminates SQLite write locks,")
	fmt.Println("    enables fast concurrent CLI operations, and guarantees zero cross-domain lock contention.")
}

func printProfileDBSection(profiles []DBFileInfo) {
	if len(profiles) == 0 {
		return
	}
	fmt.Println("  " + constants.ColorCyan + "● Profile Databases:" + constants.ColorReset)
	for _, p := range profiles {
		fmt.Printf("    %-13s %s (%s) - %s\n", p.Name, p.Path, formatBytes(p.Size), p.Purpose)
	}
	fmt.Println()
}

func printDBSummary(main DBFileInfo, splitDBs, profileDBs []DBFileInfo) {
	totalCount := 1 + len(splitDBs) + len(profileDBs)
	totalBytes := main.Size
	for _, s := range splitDBs {
		totalBytes += s.Size
	}
	for _, p := range profileDBs {
		totalBytes += p.Size
	}
	fmt.Println("  " + strings.Repeat("─", 78))
	fmt.Printf("  Total: %s%d database file(s)%s, combined size: %s%s%s on disk\n\n",
		constants.ColorWhite, totalCount, constants.ColorReset,
		constants.ColorGreen, formatBytes(totalBytes), constants.ColorReset)
}
