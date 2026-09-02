package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runAgyLsEmptyConvs(args []string) error {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return apperror.WrapSimple(err, "path error")
	}
	projects, err := loadAllAgyProjects(dirPath)
	if err != nil {
		return apperror.WrapSimple(err, "load projects")
	}
	convs, err := scanAllConversations()
	if err != nil {
		return apperror.WrapSimple(err, "scan conversations")
	}
	mapped := mapProjectsToConversations(projects, convs)
	emptyList := filterEmptyProjectConvs(mapped)

	printEmptyConvsHeader(len(emptyList), len(mapped))
	if len(emptyList) == 0 {
		return nil
	}
	printEmptyConvsTable(emptyList)
	printEmptyConvsRemediations()
	return nil
}

func filterEmptyProjectConvs(mapped []AgyProjectConvs) []AgyProjectConvs {
	var out []AgyProjectConvs
	for _, m := range mapped {
		if !m.HasActive {
			out = append(out, m)
		}
	}
	return out
}

func printEmptyConvsHeader(emptyCount, totalCount int) {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + "── Antigravity Projects with Empty Conversations ──" + constants.ColorReset)
	fmt.Println()
	if emptyCount == 0 {
		fmt.Printf("  %s✓ All %d Antigravity projects have active conversations.%s\n\n",
			constants.ColorGreen, totalCount, constants.ColorReset)
		return
	}
	fmt.Printf("  Found %s%d%s project(s) with empty or zero conversations (out of %d total):\n\n",
		constants.ColorYellow, emptyCount, constants.ColorReset, totalCount)
}

func printEmptyConvsTable(items []AgyProjectConvs) {
	fmt.Printf("  %-38s %-24s %-42s %-6s %s\n",
		"PROJECT ID", "NAME", "WORKSPACE PATH", "CONVS", "STATUS")
	fmt.Println("  " + strings.Repeat("─", 120))
	for _, it := range items {
		status := "No Convs"
		if len(it.Convs) > 0 {
			status = "Empty/Abort"
		}
		pathStr := it.Project.GetPath()
		fmt.Printf("  %-38s %-24s %-42s %-6d %s%s%s\n",
			it.Project.ID,
			truncateStr(it.Project.Name, 23),
			truncateStr(pathStr, 41),
			len(it.Convs),
			constants.ColorDim, status, constants.ColorReset)
	}
	fmt.Println()
}

func printEmptyConvsRemediations() {
	fmt.Println("  " + constants.ColorCyan + "Remediation Commands:" + constants.ColorReset)
	fmt.Println("  " + strings.Repeat("─", 80))
	fmt.Println("  ● Remove all projects with empty conversations:")
	fmt.Println("    " + constants.ColorGreen + "gitmap agy remove-projects-with-empty-conversations" + constants.ColorReset)
	fmt.Println()
	fmt.Println("  ● Remove all except specified projects, paths, or CSV:")
	fmt.Println("    " + constants.ColorGreen + "gitmap agy remove-projects-with-empty-conversations --except \"id, path, alias\"" + constants.ColorReset)
	fmt.Println("    " + constants.ColorGreen + "gitmap agy remove-projects-with-empty-conversations --except preserved_projects.csv" + constants.ColorReset)
	fmt.Println()
}
