package cmdautomation

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderGuardResult(res GuardResult) {
	printGuardHeader()
	renderOversizedList(res.OversizedFiles)
	renderBinaryList(res.BinaryFiles)
	renderLargeJsonList(res.ExcludedLargeJsons)
	printGuardSummary(res)
}

func printGuardHeader() {
	fmt.Printf("\n%s[Repository File Size & Binary Guard]%s\n", constants.ColorBold, constants.ColorReset)
}

func renderOversizedList(files []OversizedFile) {
	if len(files) == 0 {
		return
	}
	fmt.Printf("\n%s✖ Oversized Files (Exceeding Limit):%s\n", constants.ColorRed, constants.ColorReset)
	for _, f := range files {
		kb := f.SizeBytes / 1024
		maxKb := f.MaxAllowedBytes / 1024
		fmt.Printf("  %s%s%s (%d KB > %d KB limit)\n", constants.ColorRed, f.Path, constants.ColorReset, kb, maxKb)
	}
}

func renderBinaryList(files []OversizedFile) {
	if len(files) == 0 {
		return
	}
	fmt.Printf("\n%s● Discovered Common Binaries:%s\n", constants.ColorYellow, constants.ColorReset)
	for _, f := range files {
		kb := f.SizeBytes / 1024
		fmt.Printf("  %s%s%s (%d KB)\n", constants.ColorYellow, f.Path, constants.ColorReset, kb)
	}
}

func renderLargeJsonList(files []OversizedFile) {
	if len(files) == 0 {
		return
	}
	fmt.Printf("\n%s○ Large JSON Files Excluded from Search Stream:%s\n", constants.ColorCyan, constants.ColorReset)
	for _, f := range files {
		kb := f.SizeBytes / 1024
		fmt.Printf("  %s%s%s (%d KB)\n", constants.ColorCyan, f.Path, constants.ColorReset, kb)
	}
}

func printGuardSummary(res GuardResult) {
	hasIssues := len(res.OversizedFiles) > 0 || len(res.BinaryFiles) > 0
	if !hasIssues {
		fmt.Printf("\n%s✔ All passed.%s %d file(s) scanned in %s (0 oversized files).\n\n",
			constants.ColorGreen, constants.ColorReset, res.TotalFiles, res.Duration)
		return
	}

	fmt.Printf("\n%sAudit completed in %s:%s %d total files, %d oversized, %d binaries, %d large JSONs excluded.\n\n",
		constants.ColorBold, res.Duration, constants.ColorReset,
		res.TotalFiles, len(res.OversizedFiles), len(res.BinaryFiles), len(res.ExcludedLargeJsons))
}
