// Package cmdagy — agy_pin_projects_table.go renders table formatting for pinned projects.
package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderPinnedProjectsTable(projects []PinnedProject) {
	if len(projects) == 0 {
		printPinnedEmptyState()
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
	fmt.Printf("\n  %s╔══════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║       pinned antigravity projects    ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s%d pinned projects%s\n", constants.ColorDim, count, constants.ColorReset)
	fmt.Printf("  %s%s%s\n\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
}

type pinnedTableContext struct {
	MaxSeq, MaxProject, MaxID, MaxBranch, MaxStatus, MaxPinned int
}

func buildPinnedTableContext(projects []PinnedProject) *pinnedTableContext {
	ctx := &pinnedTableContext{MaxSeq: 3, MaxProject: 18, MaxID: 10, MaxBranch: 8, MaxStatus: 14, MaxPinned: 10}
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
	fmt.Printf("  %s%-*s%s%-*s%s%-*s%s%-*s%s%-*s%s%-*s%s%s%s\n",
		constants.ColorWhite,
		c.MaxSeq, "SEQ", gap,
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
	seqCol := fmt.Sprintf("%s%03d%s", constants.ColorYellow, index+1, constants.ColorReset)
	statusCol := formatAgyStatus("active", checkPathMissing(p.Path), c.MaxStatus)
	branch := p.Branch
	if branch == "" {
		branch = "—"
	}

	branchCol := fmt.Sprintf("%s%-*s%s", constants.ColorCyan, c.MaxBranch, branch, constants.ColorReset)
	projectCol := fmt.Sprintf("%s%-*s%s", color, c.MaxProject, p.Name, constants.ColorReset)
	idCol := fmt.Sprintf("%s%-*s%s", constants.ColorDim, c.MaxID, shortProjectId(p.ID), constants.ColorReset)
	pinnedCol := fmt.Sprintf("%s%-*s%s", constants.ColorYellow, c.MaxPinned, formatRelativeTime(p.PinnedAt), constants.ColorReset)

	fmt.Printf("  %s%s%s%s%s%s%s%s%s%s%s%s%s\n",
		seqCol, gap, projectCol, gap, idCol, gap, branchCol, gap, statusCol, gap, pinnedCol, gap, p.Path)
}

func printPinnedEmptyState() {
	fmt.Printf("\n  %sNo pinned Antigravity projects found.%s\n", constants.ColorDim, constants.ColorReset)
	fmt.Printf("  Pin a project with: %sgitmap agy pins add <seq|id|slug>%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %sCommands:%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy pins add <target...>%s   Pin projects\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy pins rm <target...>%s    Unpin projects\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy pins edit <target> <name>%s Rename a pin\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy pins help%s                Show help guide\n\n", constants.ColorCyan, constants.ColorReset)
}
