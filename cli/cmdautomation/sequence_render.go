package cmdautomation

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderSequenceAuditorResult(res SequenceAuditorResult) {
	printSequenceHeader()
	for _, dirRes := range res.DirResults {
		renderSingleDirResult(dirRes)
	}
	printSequenceSummary(res)
}

func printSequenceHeader() {
	fmt.Printf("\n%s[Sequence, Numbering & Title Header Auditor]%s\n\n", constants.ColorBold, constants.ColorReset)
}

func renderSingleDirResult(res SequenceDirResult) {
	if res.IsClean {
		fmt.Printf("  %s✔%s %s (%d files)\n", constants.ColorGreen, constants.ColorReset, res.Dir, res.NumberedFilesCount)
		return
	}

	fmt.Printf("  %s✖ %s%s (%d files)\n", constants.ColorYellow, res.Dir, constants.ColorReset, res.NumberedFilesCount)
	for _, gap := range res.SequenceGaps {
		fmt.Printf("    %s::error::Gap:%s %s\n", constants.ColorRed, constants.ColorReset, gap)
	}
	for _, mismatch := range res.TitleMismatches {
		fmt.Printf("    %s::warning::Title:%s %s\n", constants.ColorYellow, constants.ColorReset, mismatch)
	}
	for _, fixed := range res.FixedTitles {
		fmt.Printf("    %s::fixed::%s %s\n", constants.ColorGreen, constants.ColorReset, fixed)
	}
}

func printSequenceSummary(res SequenceAuditorResult) {
	hasIssues := res.TotalGaps > 0 || (res.TotalMismatches > 0 && res.TotalFixed == 0)
	if !hasIssues {
		fmt.Printf("\n%s✔ All sequences passed.%s %d file(s) across %d folder(s) in %s.\n\n",
			constants.ColorGreen, constants.ColorReset, res.TotalNumberedFiles, res.ScannedDirs, res.Duration)
		return
	}

	fmt.Printf("\n%sAudit completed in %s:%s %d gaps, %d title mismatches, %d fixed across %d folder(s).\n\n",
		constants.ColorBold, res.Duration, constants.ColorReset,
		res.TotalGaps, res.TotalMismatches, res.TotalFixed, res.ScannedDirs)
}
