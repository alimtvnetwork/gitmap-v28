package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// printMissingRepoRemediation renders guidance when repos are missing on disk.
func printMissingRepoRemediation(c *statusTableContext) {
	var missingSlugs []string
	for _, r := range c.Rows {
		if r.Missing {
			missingSlugs = append(missingSlugs, r.RepoName)
		}
	}
	if len(missingSlugs) == 0 {
		return
	}

	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	fmt.Printf("  %s▲ %d missing repository(ies) detected:%s\n",
		constants.ColorYellow, len(missingSlugs), constants.ColorReset)
	for _, slug := range missingSlugs {
		fmt.Printf("     • %s%s%s\n", constants.ColorRed, slug, constants.ColorReset)
	}
	fmt.Println()
	fmt.Println("  To resolve missing repositories:")
	fmt.Println("  1. Relocate to a new folder:")
	fmt.Println("     $ gitmap scan-folder update <slug> <new-path>")
	fmt.Println("  2. Untrack from database:")
	fmt.Printf("     $ gitmap rm %s\n", slugList(missingSlugs))
	fmt.Println()
}

func slugList(slugs []string) string {
	if len(slugs) > 3 {
		return fmt.Sprintf("%s %s ... (+%d more)", slugs[0], slugs[1], len(slugs)-2)
	}
	out := ""
	for i, s := range slugs {
		if i > 0 {
			out += " "
		}
		out += s
	}
	return out
}
