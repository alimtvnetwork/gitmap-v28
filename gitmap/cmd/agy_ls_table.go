// Package cmd — agy_ls_table.go renders table formatting for antigravity projects.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type agyTableRow struct {
	ID        string
	Name      string
	Branch    string
	Status    string
	Updated   string
	Path      string
	IsMissing bool
}

type agyTableContext struct {
	Rows       []agyTableRow
	MaxProject int
	MaxID      int
	MaxBranch  int
	MaxStatus  int
	MaxUpdated int
}

func newAgyTableContext() *agyTableContext {
	return &agyTableContext{
		Rows:       make([]agyTableRow, 0),
		MaxProject: 18,
		MaxID:      10,
		MaxBranch:  8,
		MaxStatus:  10,
		MaxUpdated: 9,
	}
}

func (c *agyTableContext) addRow(r agyTableRow) {
	c.Rows = append(c.Rows, r)
	if l := len(r.Name); l > c.MaxProject {
		c.MaxProject = l
	}
	if l := len(r.ID); l > c.MaxID {
		c.MaxID = l
	}
	if l := len(r.Branch); l > c.MaxBranch {
		c.MaxBranch = l
	}
	if l := len(r.Updated); l > c.MaxUpdated {
		c.MaxUpdated = l
	}
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
	branchColor := constants.ColorCyan
	if r.Branch == "—" {
		branchColor = constants.ColorDim
	}
	projectCol := fmt.Sprintf("%s%-*s%s", color, c.MaxProject, r.Name, constants.ColorReset)
	idCol := fmt.Sprintf("%s%-*s%s", constants.ColorDim, c.MaxID, r.ID, constants.ColorReset)
	branchCol := fmt.Sprintf("%s%-*s%s", branchColor, c.MaxBranch, r.Branch, constants.ColorReset)
	statusCol := formatAgyStatus(r.Status, r.IsMissing, c.MaxStatus)
	updatedCol := fmt.Sprintf("%s%-*s%s", constants.ColorDim, c.MaxUpdated, r.Updated, constants.ColorReset)

	fmt.Printf("  %s%s%s%s%s%s%s%s%s%s%s\n",
		projectCol, gap,
		idCol, gap,
		branchCol, gap,
		statusCol, gap,
		updatedCol, gap,
		r.Path,
	)
}

func formatAgyStatus(status string, isMissing bool, width int) string {
	if isMissing {
		return fmt.Sprintf("%s%-*s%s", constants.ColorRed, width, "✖ missing", constants.ColorReset)
	}
	if status == "global" {
		return fmt.Sprintf("%s%-*s%s", constants.ColorDim, width, "— global", constants.ColorReset)
	}
	return fmt.Sprintf("%s%-*s%s", constants.ColorGreen, width, "✔ active", constants.ColorReset)
}

func printAgySummary(total, active, missing int) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	fmt.Printf("  %d projects · %s%d active%s · %s%d missing%s\n",
		total,
		constants.ColorGreen, active, constants.ColorReset,
		constants.ColorRed, missing, constants.ColorReset)
	if missing > 0 {
		fmt.Printf("  %sTip: Run 'gitmap agy clear --missing' to remove stale projects.%s\n",
			constants.ColorDim, constants.ColorReset)
	}
	fmt.Println()
}
