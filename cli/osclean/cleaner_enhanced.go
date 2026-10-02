package osclean

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// CleanEnhancedDevCaches runs the 10-category deep developer cache cleaner in parallel.
func CleanEnhancedDevCaches(opts DevCleanOptions) DevCleanSummary {
	start := time.Now()
	all := getEnhancedCategoryDefs()
	selected := filterEnhancedCategories(all, opts.OnlyCategories)

	stats := executeCategoryWorkers(selected, opts.IsDryRun)
	duration := time.Since(start).Milliseconds()

	return buildSummaryTotals(stats, opts.IsDryRun, duration)
}

func executeCategoryWorkers(categories []EnhancedCategoryDef, isDryRun bool) []CategoryCleanStats {
	ch := make(chan CategoryCleanStats, len(categories))
	var wg sync.WaitGroup

	for _, cat := range categories {
		wg.Add(1)
		go func(c EnhancedCategoryDef) {
			defer wg.Done()
			ch <- c.CleanFunc(isDryRun)
		}(cat)
	}

	wg.Wait()
	close(ch)
	return collectCategoryResults(ch)
}

func collectCategoryResults(ch <-chan CategoryCleanStats) []CategoryCleanStats {
	var results []CategoryCleanStats
	for s := range ch {
		results = append(results, s)
	}
	return results
}

func buildSummaryTotals(stats []CategoryCleanStats, isDryRun bool, durationMs int64) DevCleanSummary {
	summary := DevCleanSummary{
		IsDryRun:   isDryRun,
		DurationMs: durationMs,
		Categories: stats,
	}
	for _, s := range stats {
		summary.TotalItemsRemoved += s.ItemsRemoved
		summary.TotalDirsRemoved += s.DirsRemoved
		summary.TotalBytesFreed += s.BytesFreed
	}
	return summary
}

func filterEnhancedCategories(all []EnhancedCategoryDef, only []string) []EnhancedCategoryDef {
	if len(only) == 0 {
		return all
	}
	var filtered []EnhancedCategoryDef
	for _, c := range all {
		if isEnhancedCategoryMatched(c, only) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func isEnhancedCategoryMatched(c EnhancedCategoryDef, only []string) bool {
	for _, o := range only {
		norm := strings.TrimSpace(o)
		if strings.EqualFold(norm, c.ID) || containsAliasFold(c.Aliases, norm) {
			return true
		}
	}
	return false
}

func containsAliasFold(aliases []string, needle string) bool {
	for _, a := range aliases {
		if strings.EqualFold(a, needle) {
			return true
		}
	}
	return false
}

// FormatCleanSize formats byte counts as readable MB or GB strings.
func FormatCleanSize(bytes int64) string {
	const gb = 1024 * 1024 * 1024
	const mb = 1024 * 1024
	if bytes >= gb {
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	}
	return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
}

// RenderEnhancedSummaryTable outputs an aligned terminal summary table.
func RenderEnhancedSummaryTable(summary DevCleanSummary) {
	fmt.Println()
	fmt.Printf("  %-34s %-10s %8s %8s %12s\n", "Category", "Status", "Files", "Dirs", "Reclaimed")
	fmt.Println("  " + strings.Repeat("─", 78))
	for _, cat := range summary.Categories {
		renderEnhancedTableRow(cat, summary.IsDryRun)
	}
	fmt.Println("  " + strings.Repeat("─", 78))
	renderEnhancedTableTotals(summary)
}

func renderEnhancedTableRow(cat CategoryCleanStats, isDryRun bool) {
	status := "Cleaned"
	if isDryRun {
		status = "Reclaimable"
	}
	sizeStr := FormatCleanSize(cat.BytesFreed)
	fmt.Printf("  • %-32s %-10s %8d %8d %12s\n",
		cat.Label, status, cat.ItemsRemoved, cat.DirsRemoved, sizeStr)
}

func renderEnhancedTableTotals(summary DevCleanSummary) {
	mode := "Cleaned"
	if summary.IsDryRun {
		mode = "Dry-Run Reclaimable"
	}
	totalSize := FormatCleanSize(summary.TotalBytesFreed)
	fmt.Printf("  ✔ Total %s: %s across %d files and %d dirs (%dms)\n\n",
		mode, totalSize, summary.TotalItemsRemoved, summary.TotalDirsRemoved, summary.DurationMs)
}
