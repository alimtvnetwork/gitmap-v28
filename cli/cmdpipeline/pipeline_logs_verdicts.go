package cmdpipeline

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderRunsSummaryTable(p PipelineErrorLogsPayload) {
	runs := p.Runs
	if len(runs) == 0 {
		runs = resolveCachedRunsOrFetch(p.Repo)
	}
	RenderHistorySummaryTable(runs)
}

func renderCleanSuccessTerminal(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Pipeline Status: CLEAN (No errors found)%s\n",
		constants.ColorGreen, constants.ColorReset)
	renderPipelineMetaBlock(p)
	fmt.Printf("\n  %s✓ All recent pipeline runs for %s on branch %s are PASSING (100%% green).%s\n\n",
		constants.ColorGreen, p.Repo, p.LatestBranch, constants.ColorReset)
	if len(p.DbPath) > 0 {
		fmt.Printf("  • Pipeline DB:     %s\n", filepath.ToSlash(FormatRelativeDbPath(p.DbPath)))
		fmt.Printf("  • DB Size:         %s\n", ResolveDbFileSize(p.DbPath))
		fmt.Printf("  • Cleanup:         gitmap pipeline clear -y\n")
	}
	renderRunsSummaryTable(p)
	printRerunETA(p.RerunEtaSeconds)
}

func renderUnknownTerminal(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Pipeline Status: UNKNOWN (No pipeline data retrieved)%s\n",
		constants.ColorYellow, constants.ColorReset)
	renderPipelineMetaBlock(p)
	fmt.Printf("\n  %s⚠ The pipeline status for %s could not be determined. This is NOT a clean result.%s\n",
		constants.ColorYellow, p.Repo, constants.ColorReset)
	fmt.Printf("    %s\n\n", p.ErrorLogs)
	if len(p.DbPath) > 0 {
		fmt.Printf("  • Pipeline DB:     %s\n", filepath.ToSlash(FormatRelativeDbPath(p.DbPath)))
		fmt.Printf("  • DB Size:         %s\n", ResolveDbFileSize(p.DbPath))
		fmt.Printf("  • Cleanup:         gitmap pipeline clear -y\n")
	}
}
