package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func resolveBatchLimit(opts AgyFixOptions) int {
	if opts.Limit > 0 {
		return opts.Limit
	}
	if opts.ProjectsCount > 0 {
		return opts.ProjectsCount
	}

	return 3
}

func printBatchHeader(start, end, total int) {
	fmt.Printf("\n  %s● Multi-Project Batch Pipeline Fix%s (Projects %d-%d of %d failing)\n",
		constants.ColorCyan, constants.ColorReset, start+1, end, total)
	fmt.Printf("  ────────────────────────────────────────────────────────────────────────────\n")
}

func printBatchProgress(end, total int) {
	remaining := total - end
	if remaining > 0 {
		fmt.Printf("\n  %s✔ Batch complete: %d of %d projects processed.%s\n",
			constants.ColorGreen, end, total, constants.ColorReset)
		fmt.Printf("    Run again to process the next batch of %d projects.\n\n", remaining)
		return
	}

	fmt.Printf("\n  %s✔ All %d failing projects have been processed across all batches!%s\n\n",
		constants.ColorGreen, total, constants.ColorReset)
}
