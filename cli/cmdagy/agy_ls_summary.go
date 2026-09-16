// Package cmdagy — agy_ls_summary.go prints project count summaries and missing project tips.
package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func printAgySummary(total, active, missing int) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	fmt.Printf("  %d projects · %s%d active%s · %s%d missing%s\n",
		total,
		constants.ColorGreen, active, constants.ColorReset,
		constants.ColorRed, missing, constants.ColorReset)
	if missing > 0 {
		printAgyMissingTips()
	}

	fmt.Println()
}

func printAgyMissingTips() {
	fmt.Printf("  %sReconciliation & Cleanup for Missing Projects:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    %s●%s Reconcile & re-link: %sgitmap agy reconcile%s (alias: recon)\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %s●%s Scan & discover:    %sgitmap agy scan [path]%s\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %s●%s Remove missing:     %sgitmap agy remove-missing-projects%s (alias: rm-missing)\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
}
