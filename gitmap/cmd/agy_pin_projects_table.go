// Package cmd — agy_pin_projects_table.go renders table formatting for pinned projects.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func renderPinnedProjectsTable(projects []PinnedProject) {
	if len(projects) == 0 {
		fmt.Printf("\n  %sNo pinned Antigravity projects found.%s\n", constants.ColorDim, constants.ColorReset)
		fmt.Printf("  Pin a project with: %sgitmap agy pin-projects add <project-id-or-path>%s\n\n", constants.ColorCyan, constants.ColorReset)
		return
	}

	printPinnedBanner(len(projects))
	ctx := buildPinnedTableContext(projects)
	printPinnedTableHeader(ctx)

	for i, p := range projects {
		printPinnedTableRow(ctx, p, i)
	}

	printPinnedSummary(projects)
}

func printPinnedBanner(count int) {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║       pinned antigravity projects    ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════╝%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println()
	fmt.Printf("  %s%d pinned projects%s\n", constants.ColorDim, count, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	fmt.Println()
}

type pinnedTableContext struct {
	MaxProject int
	MaxID      int
	MaxBranch  int
	MaxStatus  int
	MaxPinned  int
}

func buildPinnedTableContext(projects []PinnedProject) *pinnedTableContext {
	ctx := &pinnedTableContext{
		MaxProject: 18,
		MaxID:      10,
		MaxBranch:  8,
		MaxStatus:  14,
		MaxPinned:  10,
	}

	for _, p := range projects {
		if len(p.Name) > ctx.MaxProject {
			ctx.MaxProject = len(p.Name)
		}
		if len(shortProjectId(p.ID)) > ctx.MaxID {
			ctx.MaxID = len(shortProjectId(p.ID))
		}
	}

	return ctx
}

func printPinnedTableHeader(c *pinnedTableContext) {
	const gap = "   "
	fmt.Printf("  %s%-*s%s%-*s%s%-*s%s%-*s%s%-*s%s%s%s\n",
		constants.ColorWhite,
		c.MaxProject, "PROJECT", gap,
		c.MaxID, "ID", gap,
		c.MaxBranch, "BRANCH", gap,
		c.MaxStatus, "STATUS", gap,
		c.MaxPinned, "PINNED", gap,
		"PATH",
		constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
}

func printPinnedTableRow(c *pinnedTableContext, p PinnedProject, index int) {
	const gap = "   "
	color := constants.ColorCycle[index%len(constants.ColorCycle)]
	isMissing := checkPathMissing(p.Path)
	statusCol := formatAgyStatus("active", isMissing, c.MaxStatus)
	branch := p.Branch

	if branch == "" {
		branch = "—"
	}

	branchCol := fmt.Sprintf("%s%-*s%s", constants.ColorCyan, c.MaxBranch, branch, constants.ColorReset)
	projectCol := fmt.Sprintf("%s%-*s%s", color, c.MaxProject, p.Name, constants.ColorReset)
	idCol := fmt.Sprintf("%s%-*s%s", constants.ColorDim, c.MaxID, shortProjectId(p.ID), constants.ColorReset)
	pinnedCol := fmt.Sprintf("%s%-*s%s", constants.ColorYellow, c.MaxPinned, formatRelativeTime(p.PinnedAt), constants.ColorReset)

	fmt.Printf("  %s%s%s%s%s%s%s%s%s%s%s\n",
		projectCol, gap,
		idCol, gap,
		branchCol, gap,
		statusCol, gap,
		pinnedCol, gap,
		p.Path,
	)
}

func checkPathMissing(path string) bool {
	if path == "" {
		return true
	}

	_, err := os.Stat(path)

	return err != nil
}

func printPinnedSummary(projects []PinnedProject) {
	missing := 0

	for _, p := range projects {
		if checkPathMissing(p.Path) {
			missing++
		}
	}

	active := len(projects) - missing
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	fmt.Printf("  %d pinned · %s%d active%s · %s%d missing%s\n\n",
		len(projects),
		constants.ColorGreen, active, constants.ColorReset,
		constants.ColorRed, missing, constants.ColorReset)
}

func outputPinnedProjectsJSON(projects []PinnedProject) error {
	data, err := json.MarshalIndent(projects, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
