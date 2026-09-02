package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/vscodepm"
)

func runFindDuplicatesVSCode() error {
	path, err := vscodepm.ProjectsJSONPath()
	if err != nil {
		fmt.Println("  " + constants.ColorYellow + "VS Code projects.json path not resolvable." + constants.ColorReset)
		return nil
	}
	entries, err := vscodepm.ReadEntries(path)
	if err != nil || len(entries) == 0 {
		fmt.Println("  " + constants.ColorDim + "No VS Code projects found in projects.json." + constants.ColorReset)
		return nil
	}
	dupGroups := groupVSCodeDuplicates(entries)
	if len(dupGroups) == 0 {
		fmt.Printf("  %s✓ VS Code: No duplicate projects found. Total active: %d%s\n\n",
			constants.ColorGreen, len(entries), constants.ColorReset)
		return nil
	}
	printVSCodeDupFindings(dupGroups)
	printVSCodeRemediations(dupGroups)
	return nil
}

func groupVSCodeDuplicates(entries []vscodepm.Entry) map[string][]vscodepm.Entry {
	groups := make(map[string][]vscodepm.Entry)
	for _, e := range entries {
		norm := strings.ToLower(filepath.Clean(e.RootPath))
		if norm == "" {
			continue
		}
		groups[norm] = append(groups[norm], e)
	}
	dupGroups := make(map[string][]vscodepm.Entry)
	for k, list := range groups {
		if len(list) > 1 {
			dupGroups[k] = list
		}
	}
	return dupGroups
}

func printVSCodeDupFindings(dupGroups map[string][]vscodepm.Entry) {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + "── VS Code Duplicate Projects ──" + constants.ColorReset)
	totalDups := 0
	for _, list := range dupGroups {
		totalDups += len(list) - 1
	}
	fmt.Printf("  Found %s%d%s duplicate group(s) (%s%d%s duplicate entries total):\n\n",
		constants.ColorWhite, len(dupGroups), constants.ColorReset,
		constants.ColorYellow, totalDups, constants.ColorReset)

	groupNum := 1
	for path, list := range dupGroups {
		fmt.Printf("  Group %d: RootPath: %s%s%s (%d entries)\n", groupNum, constants.ColorCyan, path, constants.ColorReset, len(list))
		fmt.Printf("    %-24s %s\n", "NAME", "PATHS / TAGS")
		fmt.Println("    " + strings.Repeat("─", 68))
		for _, e := range list {
			extra := strings.Join(e.Paths, ", ")
			fmt.Printf("    %-24s %s\n", truncateStr(e.Name, 23), truncateStr(extra, 42))
		}
		fmt.Println()
		groupNum++
	}
}

func printVSCodeRemediations(dupGroups map[string][]vscodepm.Entry) {
	var samplePath, sampleName string
	for path, list := range dupGroups {
		samplePath = path
		if len(list) > 0 {
			sampleName = list[0].Name
		}
		break
	}
	fmt.Println("  " + constants.ColorCyan + "Remediation & Fix Commands for VS Code:" + constants.ColorReset)
	fmt.Println("  " + strings.Repeat("─", 68))
	fmt.Printf("  ● Fix Single (Remove specific duplicate project by path or name):\n")
	fmt.Printf("    %sgitmap vscode rm \"%s\"%s\n", constants.ColorGreen, sampleName, constants.ColorReset)
	fmt.Printf("    %sgitmap vscode rm \"%s\"%s\n\n", constants.ColorGreen, samplePath, constants.ColorReset)
	fmt.Printf("  ● Fix All Together (Merge duplicates & clean repeated projects):\n")
	fmt.Printf("    %sgitmap vscode optimize-projects%s   (alias: %sgitmap vscode --repeat-fix%s)\n\n",
		constants.ColorGreen, constants.ColorReset, constants.ColorWhite, constants.ColorReset)
	fmt.Printf("  ● Remap / Re-add Cleanly:\n")
	fmt.Printf("    %sgitmap vscode add \"%s\"%s\n\n", constants.ColorGreen, samplePath, constants.ColorReset)
}
