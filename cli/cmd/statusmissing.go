package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// printMissingRepoRemediation renders guidance when repos are missing on disk.
func printMissingRepoRemediation(c *statusTableContext) {
	missingSlugs := collectMissingSlugs(c)
	hasMissing := len(missingSlugs) > 0
	if !hasMissing {
		return
	}

	printMissingBanner(len(missingSlugs))
	printMissingItems(missingSlugs)
	printMissingRemediationSteps(missingSlugs)
}

func collectMissingSlugs(c *statusTableContext) []string {
	var missingSlugs []string
	if c == nil {
		return nil
	}

	for _, r := range c.Rows {
		if r.Missing {
			missingSlugs = append(missingSlugs, r.RepoName)
		}
	}

	return missingSlugs
}

func printMissingBanner(count int) {
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	fmt.Printf("  %s▲ %d missing repository(ies) detected:%s\n",
		constants.ColorYellow, count, constants.ColorReset)
}

func printMissingItems(slugs []string) {
	for _, slug := range slugs {
		fmt.Printf("     • %s%s%s\n", constants.ColorRed, slug, constants.ColorReset)
	}

	fmt.Println()
}

func printMissingRemediationSteps(slugs []string) {
	fmt.Println("  To resolve missing repositories:")
	fmt.Println("  1. Relocate to a new folder:")
	fmt.Println("     $ gitmap scan-folder update <slug> <new-path>")
	fmt.Println("  2. Untrack from database:")
	fmt.Printf("     $ gitmap rm %s\n", slugList(slugs))
	fmt.Println("  3. Clear or reset tracking database:")
	fmt.Println("     $ gitmap db-reset  (or: $ gitmap storage reset-errors)")
	fmt.Println()
}

func slugList(slugs []string) string {
	if len(slugs) > 3 {
		return fmt.Sprintf("%s %s ... (+%d more)", slugs[0], slugs[1], len(slugs)-2)
	}

	return strings.Join(slugs, " ")
}
