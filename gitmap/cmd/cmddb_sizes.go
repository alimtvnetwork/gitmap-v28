package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runDBSizes(args []string) error {
	mainInfo, exists := collectMainDBInfo()
	splitDBs := collectSplitDBs()
	profileDBs := collectProfileDBs()

	var allItems []DBFileInfo
	if exists {
		allItems = append(allItems, mainInfo)
	}
	allItems = append(allItems, profileDBs...)
	allItems = append(allItems, splitDBs...)

	printDBSizesHeader()
	printDBSizesTable(allItems)
	printDBSizesTotal(allItems)
	return nil
}

func printDBSizesHeader() {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + "── Gitmap Database Disk Sizes ──" + constants.ColorReset)
	fmt.Println()
}

func printDBSizesTable(items []DBFileInfo) {
	fmt.Printf("  %-28s %-20s %-12s %s\n", "DATABASE FILE", "CATEGORY", "SIZE", "PATH")
	fmt.Println("  " + strings.Repeat("─", 90))
	for _, it := range items {
		catColor := constants.ColorCyan
		if it.Category == "Primary Master DB" {
			catColor = constants.ColorWhite
		} else if it.Category == "Split Repository DB" {
			catColor = constants.ColorDim
		}
		fmt.Printf("  %-28s %s%-20s%s %-12s %s\n",
			truncateStr(it.Name, 27),
			catColor, it.Category, constants.ColorReset,
			formatBytes(it.Size),
			it.Path,
		)
	}
}

func printDBSizesTotal(items []DBFileInfo) {
	var totalBytes int64
	for _, it := range items {
		totalBytes += it.Size
	}
	fmt.Println("  " + strings.Repeat("─", 90))
	fmt.Printf("  Total: %s%d database file(s)%s, combined size: %s%s%s\n\n",
		constants.ColorWhite, len(items), constants.ColorReset,
		constants.ColorGreen, formatBytes(totalBytes), constants.ColorReset)
}
