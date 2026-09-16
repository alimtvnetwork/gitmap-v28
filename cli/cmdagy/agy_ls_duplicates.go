// Package cmdagy — agy_ls_duplicates.go surfaces duplicate projects outside the status table.
package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func printExternalAgyDuplicates(projects []AgyProject) {
	dupGroups := groupAgyDuplicates(projects)
	if len(dupGroups) == 0 {
		fmt.Printf("  %s✓ Antigravity (AGY): No duplicate project paths found.%s\n\n",
			constants.ColorGreen, constants.ColorReset)

		return
	}

	printAgyDupFindings(dupGroups)
	printAgyRemediations(dupGroups)
}
