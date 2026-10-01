package cmdos

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type DevRenderCard struct {
	Category   string
	Label      string
	Badge      string
	BadgeColor string
	Paths      []DiscoveredCachePath
	TotalBytes int64
	TotalFiles int
}

var ecosystemDisplayNames = map[string]string{
	"go":     "Go (Golang)",
	"pnpm":   "pnpm Store",
	"npm":    "npm Cache",
	"yarn":   "Yarn Cache",
	"bun":    "Bun Cache",
	"python": "Python Pip/UV",
	"cargo":  "Rust Cargo",
	"nuget":  ".NET NuGet",
	"gradle": "Java Gradle",
	"maven":  "Java Maven",
}

func resolveEcosystemDisplayName(eco string) string {
	if name, ok := ecosystemDisplayNames[strings.ToLower(eco)]; ok {
		return name
	}
	return eco
}

func getPathEcosystemKey(p DiscoveredCachePath) string {
	if len(p.Ecosystem) > 0 {
		return p.Ecosystem
	}
	return p.Category
}

func groupPathsByEcosystem(paths []DiscoveredCachePath) ([]string, map[string][]DiscoveredCachePath) {
	var order []string
	grouped := make(map[string][]DiscoveredCachePath)
	for _, p := range paths {
		eco := getPathEcosystemKey(p)
		if _, exists := grouped[eco]; !exists {
			order = append(order, eco)
		}
		grouped[eco] = append(grouped[eco], p)
	}
	return order, grouped
}

func RenderDevCleanAnsiCards(paths []DiscoveredCachePath, isDryRun bool) {
	order, grouped := groupPathsByEcosystem(paths)
	for _, eco := range order {
		card := buildDevRenderCard(eco, grouped[eco], isDryRun)
		renderSingleDevCard(card)
	}
	renderDevCleanTotals(paths, isDryRun)
}

func sumCardMetrics(paths []DiscoveredCachePath) (int64, int) {
	var totalBytes int64
	var totalFiles int
	for _, p := range paths {
		totalBytes += p.SizeBytes
		totalFiles += p.FilesCount
	}
	return totalBytes, totalFiles
}

func buildDevRenderCard(eco string, paths []DiscoveredCachePath, isDryRun bool) DevRenderCard {
	bytes, files := sumCardMetrics(paths)
	return DevRenderCard{
		Category:   eco,
		Label:      resolveEcosystemDisplayName(eco),
		Badge:      resolveBadgeText(isDryRun),
		BadgeColor: resolveBadgeColor(isDryRun),
		Paths:      paths,
		TotalBytes: bytes,
		TotalFiles: files,
	}
}

func renderSingleDevCard(c DevRenderCard) {
	border := constants.ColorDim
	reset := constants.ColorReset
	fmt.Printf("\n  %s┌── 📦 %s %s%s%s %s%s\n", border, c.Label, c.BadgeColor, c.Badge, reset, border, strings.Repeat("─", 36))
	fmt.Printf("  %s│%s  Total: %s across %d files in %d paths\n", border, reset, formatAnsiSize(c.TotalBytes), c.TotalFiles, len(c.Paths))
	fmt.Printf("  %s│%s\n", border, reset)
	for i, p := range c.Paths {
		renderCardPathRow(p, i == len(c.Paths)-1, border)
	}
	fmt.Printf("  %s└──%s%s\n", border, strings.Repeat("─", 68), reset)
}

func renderCardPathRow(p DiscoveredCachePath, isLast bool, border string) {
	connector := "├──"
	if isLast {
		connector = "└──"
	}
	src := p.DiscoverySource
	if len(src) == 0 {
		src = p.Source
	}
	reset := constants.ColorReset
	fmt.Printf("  %s│%s  %s 📁 %s%s%s\n", border, reset, connector, constants.ColorCyan, p.Path, reset)
	fmt.Printf("  %s│%s      ├── Size: %s │ Files: %d │ [%s%s%s]\n",
		border, reset, formatAnsiSize(p.SizeBytes), p.FilesCount, constants.ColorDim, src, reset)
}

func renderDevCleanTotals(paths []DiscoveredCachePath, isDryRun bool) {
	var totalBytes int64
	var totalFiles int
	for _, p := range paths {
		totalBytes += p.SizeBytes
		totalFiles += p.FilesCount
	}
	actionLabel := "Cleaned"
	if isDryRun {
		actionLabel = "Dry-Run Reclaimable"
	}
	fmt.Printf("\n  %s✔ %s: %s across %d files and %d directories%s\n\n",
		constants.ColorGreen, actionLabel, formatAnsiSize(totalBytes), totalFiles, len(paths), constants.ColorReset)
}

func resolveBadgeText(isDryRun bool) string {
	if isDryRun {
		return "[⚡ Reclaimable]"
	}
	return "[✔ Cleaned]"
}

func resolveBadgeColor(isDryRun bool) string {
	if isDryRun {
		return constants.ColorYellow
	}
	return constants.ColorGreen
}

func FormatCleanSize(bytes int64) string {
	const kb = 1024
	const mb = 1024 * kb
	const gb = 1024 * mb
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func formatAnsiSize(bytes int64) string {
	const gb = 1024 * 1024 * 1024
	const mb = 50 * 1024 * 1024
	sizeStr := FormatCleanSize(bytes)
	if bytes >= gb {
		return fmt.Sprintf("%s%s%s", constants.ColorGreen, sizeStr, constants.ColorReset)
	}
	if bytes >= mb {
		return fmt.Sprintf("%s%s%s", constants.ColorYellow, sizeStr, constants.ColorReset)
	}
	return fmt.Sprintf("%s%s%s", constants.ColorCyan, sizeStr, constants.ColorReset)
}
