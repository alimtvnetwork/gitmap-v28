package cmdagy

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type promptScanStats struct {
	totalPrompts  int
	recent24h     int
	recent7d      int
	totalArchives int
}

func computePromptScanStats() promptScanStats {
	allPrompts := CollectAllPrompts()
	now := time.Now()
	cutoff24h := now.Add(-24 * time.Hour)
	cutoff7d := now.Add(-7 * 24 * time.Hour)
	var stats promptScanStats
	stats.totalPrompts = len(allPrompts)
	for _, p := range allPrompts {
		tallyPromptTime(p.CreatedAt, cutoff24h, cutoff7d, &stats)
	}
	stats.totalArchives = countHistoricalArchives()

	return stats
}

func tallyPromptTime(createdAt, cutoff24h, cutoff7d time.Time, stats *promptScanStats) {
	if createdAt.After(cutoff24h) {
		stats.recent24h++
	}
	if createdAt.After(cutoff7d) {
		stats.recent7d++
	}
}

func countHistoricalArchives() int {
	brainDir, err := GetBrainLogsDirPath()
	if err != nil {
		return 0
	}
	entries, readErr := os.ReadDir(brainDir)
	if readErr != nil {
		return 0
	}

	return len(entries)
}

func printPromptScanSummary(stats promptScanStats) {
	fmt.Printf("  %s● AGY Prompts Activity:%s %d total prompts · %s%d last 24h%s · %s%d last 7d%s · %d archives\n",
		constants.ColorCyan, constants.ColorReset,
		stats.totalPrompts,
		constants.ColorGreen, stats.recent24h, constants.ColorReset,
		constants.ColorYellow, stats.recent7d, constants.ColorReset,
		stats.totalArchives,
	)
}
