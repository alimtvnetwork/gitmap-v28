package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runFindDuplicatesAgy() error {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		fmt.Println("  " + constants.ColorYellow + "Antigravity config directory not found." + constants.ColorReset)
		return nil
	}
	projects, err := loadAllAgyProjects(dirPath)
	if err != nil || len(projects) == 0 {
		fmt.Println("  " + constants.ColorDim + "No Antigravity projects configured." + constants.ColorReset)
		return nil
	}
	dupGroups := groupAgyDuplicates(projects)
	if len(dupGroups) == 0 {
		fmt.Printf("  %s✓ Antigravity (AGY): No duplicate projects found. Total active: %d%s\n\n",
			constants.ColorGreen, len(projects), constants.ColorReset)
		return nil
	}
	printAgyDupFindings(dupGroups)
	printAgyRemediations(dupGroups)
	return nil
}

func groupAgyDuplicates(projects []AgyProject) map[string][]AgyProject {
	groups := groupProjectsByPath(projects)
	dupGroups := make(map[string][]AgyProject)
	for path, list := range groups {
		if len(list) > 1 {
			dupGroups[path] = list
		}
	}
	return dupGroups
}

func printAgyDupFindings(dupGroups map[string][]AgyProject) {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + "── Antigravity (AGY) Duplicate Projects ──" + constants.ColorReset)
	totalDups := 0
	for _, list := range dupGroups {
		totalDups += len(list) - 1
	}
	fmt.Printf("  Found %s%d%s duplicate project group(s) (%s%d%s duplicate entries total):\n\n",
		constants.ColorWhite, len(dupGroups), constants.ColorReset,
		constants.ColorYellow, totalDups, constants.ColorReset)

	groupNum := 1
	for path, list := range dupGroups {
		fmt.Printf("  Group %d: Path: %s%s%s (%d entries)\n", groupNum, constants.ColorCyan, path, constants.ColorReset, len(list))
		fmt.Printf("    %-38s %-22s %s\n", "PROJECT ID", "NAME", "UPDATED")
		fmt.Println("    " + strings.Repeat("─", 74))
		for _, p := range list {
			fmt.Printf("    %-38s %-22s %s\n", p.ID, truncateStr(p.Name, 21), truncateStr(p.UpdatedAt, 20))
		}
		fmt.Println()
		groupNum++
	}
}

func printAgyRemediations(dupGroups map[string][]AgyProject) {
	var sampleID, samplePath string
	for path, list := range dupGroups {
		if len(list) > 1 {
			sampleID = list[1].ID
			samplePath = path
			break
		}
	}
	fmt.Println("  " + constants.ColorCyan + "Remediation & Fix Commands for Antigravity:" + constants.ColorReset)
	fmt.Println("  " + strings.Repeat("─", 74))
	fmt.Printf("  ● Fix Single (Delete specific duplicate project ID):\n")
	fmt.Printf("    %sgitmap agy rm %s%s\n\n", constants.ColorGreen, sampleID, constants.ColorReset)
	fmt.Printf("  ● Fix All Together (Deduplicate & keep newest per path):\n")
	fmt.Printf("    %sgitmap agy optimize-projects%s   (alias: %sgitmap agy --repeat-fix%s)\n\n",
		constants.ColorGreen, constants.ColorReset, constants.ColorWhite, constants.ColorReset)
	fmt.Printf("  ● Remap / Re-add Cleanly:\n")
	fmt.Printf("    %sgitmap agy add \"%s\"%s\n\n", constants.ColorGreen, samplePath, constants.ColorReset)
}
