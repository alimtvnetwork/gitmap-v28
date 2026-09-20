package osclean

import (
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type categoryCleaner struct {
	Name    string
	Aliases []string
	Clean   func(bool) CategoryCleanStats
}

func CleanDevCaches(opts DevCleanOptions) result.Result[DevCleanSummary] {
	start := time.Now()
	var summary DevCleanSummary
	summary.IsDryRun = opts.IsDryRun

	cleaners := getCategoryCleaners()
	for _, c := range cleaners {
		if !isCategorySelected(c, opts.OnlyCategories) {
			continue
		}
		stats := c.Clean(opts.IsDryRun)
		accumulateSummary(&summary, stats)
	}

	summary.DurationMs = time.Since(start).Milliseconds()

	return result.Ok(summary)
}

func accumulateSummary(summary *DevCleanSummary, stats CategoryCleanStats) {
	summary.TotalItemsRemoved += stats.ItemsRemoved
	summary.TotalDirsRemoved += stats.DirsRemoved
	summary.TotalBytesFreed += stats.BytesFreed
	summary.Categories = append(summary.Categories, stats)
}

func getCategoryCleaners() []categoryCleaner {
	return []categoryCleaner{
		{"go-buildcache", []string{"go", "gobuild", "gocache"}, cleanGoCache},
		{"pnpm-store", []string{"pnpm", "pnpmstore"}, cleanPnpmCache},
		{"npm-cache", []string{"npm", "npmcache"}, cleanNpmCache},
		{"choco-cache", []string{"choco", "chocolatey"}, cleanChocoCache},
		{"yarn-cache", []string{"yarn", "yarncache"}, cleanYarnCache},
		{"bun-cache", []string{"bun", "buncache"}, cleanBunCache},
		{"pip-cache", []string{"pip", "python", "pipcache"}, cleanPipCache},
		{"cargo-cache", []string{"cargo", "rust", "cargocache"}, cleanCargoCache},
		{"nuget-cache", []string{"nuget", "dotnet", "nugetcache"}, cleanNugetCache},
		{"gradle-maven-cache", []string{"gradle", "maven", "m2"}, cleanGradleMavenCache},
	}
}

func isCategorySelected(c categoryCleaner, only []string) bool {
	if len(only) == 0 {
		return true
	}
	for _, raw := range only {
		needle := strings.ToLower(strings.TrimSpace(raw))
		if needle == c.Name || containsAlias(c.Aliases, needle) {
			return true
		}
	}

	return false
}

func containsAlias(aliases []string, needle string) bool {
	for _, a := range aliases {
		if a == needle {
			return true
		}
	}

	return false
}
