// Package cmd — agy_find_cure_dups.go handles finding and curing duplicate Antigravity projects.
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var (
	agyFdpExcept string
	agyFdpJson   bool
)

var agyFindDupsCmd = &cobra.Command{
	Use:     "find-duplicate-projects",
	Aliases: []string{"fdp", "find-duplicate", "find-duplicates"},
	Short:   "Find duplicate Antigravity projects and show cure commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyFindDups()
	},
}

func init() {
	agyFindDupsCmd.Flags().StringVarP(&agyFdpExcept, "except", "e", "", "Exclude projects matching id, name, slug, or path")
	agyFindDupsCmd.Flags().BoolVar(&agyFdpJson, "json", false, "Emit structured JSON output")
	agyOptimizeCmd.Aliases = append(agyOptimizeCmd.Aliases, "cure-duplicate-projects", "cdp", "cure-duplicates", "cure-duplicate")
}

func runAgyFindDups() error {
	dirPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load projects")
	}
	dupGroups := buildAgyDuplicateGroups(projects, agyFdpExcept)
	if len(dupGroups) == 0 {
		fmt.Printf("%s No duplicate Antigravity projects found across %d projects.\n",
			constants.ColorGreen+"✓"+constants.ColorReset, len(projects))
		return nil
	}
	renderAgyDuplicateFindings(dupGroups)
	return nil
}

func buildAgyDuplicateGroups(projects []AgyProject, exceptStr string) map[string][]AgyProject {
	groups := groupProjectsByPath(projects)
	dupGroups := make(map[string][]AgyProject)
	for path, list := range groups {
		filtered := make([]AgyProject, 0)
		for _, p := range list {
			if !isAgyProjectExcepted(p, exceptStr) {
				filtered = append(filtered, p)
			}
		}
		if len(filtered) > 1 {
			dupGroups[path] = filtered
		}
	}
	return dupGroups
}

func renderAgyDuplicateFindings(dupGroups map[string][]AgyProject) {
	totalDups := 0
	for _, list := range dupGroups {
		totalDups += len(list) - 1
	}
	fmt.Printf("\n  %s── Antigravity (AGY) Duplicate Projects ──%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("  Found %d duplicate project group(s) (%d duplicate entries total):\n\n", len(dupGroups), totalDups)
	groupNum := 1
	for path, list := range dupGroups {
		sortGroupNewestFirst(list)
		printAgyDupGroup(groupNum, path, list)
		groupNum++
	}
	printAgyDupRemediation()
}

func printAgyDupGroup(num int, path string, list []AgyProject) {
	fmt.Printf("  Group %d: Path: %s%s%s (%d entries)\n",
		num, constants.ColorCyan, path, constants.ColorReset, len(list))
	fmt.Printf("    %-38s %-22s %-20s %s\n", "PROJECT ID", "NAME", "UPDATED", "STATUS")
	fmt.Printf("    %s\n", strings.Repeat("─", 90))
	for idx, p := range list {
		status := constants.ColorGreen + "KEEP (newest)" + constants.ColorReset
		if idx > 0 {
			status = constants.ColorRed + "DUPLICATE (older)" + constants.ColorReset
		}
		fmt.Printf("    %-38s %-22s %-20s %s\n", p.ID, p.Name, p.UpdatedAt, status)
	}
	fmt.Println()
}

func printAgyDupRemediation() {
	fmt.Printf("  %sCure & Remediation Commands for Antigravity:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s\n", strings.Repeat("─", 74))
	fmt.Println("  ● Cure All Duplicates (keep newest per path):")
	fmt.Println("    gitmap agy cure-duplicate-projects   (alias: gitmap agy cdp)")
	fmt.Println()
	fmt.Println("  ● Preview Cure without Deleting (dry-run):")
	fmt.Println("    gitmap agy cure-duplicate-projects --dry-run")
	fmt.Println()
	fmt.Println("  ● Cure with Exclusions:")
	fmt.Println("    gitmap agy cure-duplicate-projects --except \"project-id, name, path\"")
	fmt.Println()
}
