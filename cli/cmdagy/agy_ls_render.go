// Package cmdagy — agy_ls_render.go renders formatted tables and banners for Antigravity projects.
package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderAgyProjectsTable(projects []AgyProject, dirPath string) {
	printAgyBanner(len(projects), dirPath)
	activeCount, missingCount := renderAgyProjectsGrouped(projects)
	printAgySummary(len(projects), activeCount, missingCount)
	printExternalAgyDuplicates(projects)
}

func printAgyBanner(count int, dirPath string) {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║         antigravity projects         ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════╝%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println()
	fmt.Printf("  %s%d projects from %s%s\n", constants.ColorDim, count, dirPath, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	fmt.Println()
}

func printAgyTableHeader(c *agyTableContext) {
	const gap = "   "
	fmt.Printf("  %s%-*s%s%-*s%s%-*s%s%-*s%s%-*s%s%s%s\n",
		constants.ColorWhite,
		c.MaxProject, "PROJECT", gap,
		c.MaxID, "ID", gap,
		c.MaxBranch, "BRANCH", gap,
		c.MaxStatus, "STATUS", gap,
		c.MaxUpdated, "UPDATED", gap,
		"PATH",
		constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
}

func printAgyTableRow(c *agyTableContext, r agyTableRow, index int) {
	const gap = "   "
	color := constants.ColorCycle[index%len(constants.ColorCycle)]
	branchColor := resolveBranchColor(r.Branch)

	projectCol := fmt.Sprintf("%s%-*s%s", color, c.MaxProject, r.Name, constants.ColorReset)
	idCol := fmt.Sprintf("%s%-*s%s", constants.ColorDim, c.MaxID, r.ID, constants.ColorReset)
	branchCol := fmt.Sprintf("%s%-*s%s", branchColor, c.MaxBranch, r.Branch, constants.ColorReset)
	statusCol := formatAgyStatus(r.Status, r.IsMissing, c.MaxStatus)
	updatedCol := fmt.Sprintf("%s%-*s%s", constants.ColorDim, c.MaxUpdated, r.Updated, constants.ColorReset)

	fmt.Printf("  %s%s%s%s%s%s%s%s%s%s%s\n",
		projectCol, gap, idCol, gap, branchCol, gap, statusCol, gap, updatedCol, gap, r.Path)
}

func resolveBranchColor(branch string) string {
	if branch == "—" {
		return constants.ColorDim
	}

	return constants.ColorCyan
}

func formatAgyStatus(status string, isMissing bool, width int) string {
	if isMissing {
		return fmt.Sprintf("%s%-*s%s", constants.ColorRed, width, "✖   missing", constants.ColorReset)
	}

	if status == "global" {
		return fmt.Sprintf("%s%-*s%s", constants.ColorDim, width, "—   global", constants.ColorReset)
	}

	return fmt.Sprintf("%s%-*s%s", constants.ColorGreen, width, "✔   active", constants.ColorReset)
}
