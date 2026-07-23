package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

// renderChangelogEntry pretty-prints a single changelog entry.
//
// When pretty is true, **bold** / `code` / "double quotes" markdown is
// rendered with ANSI and the header/marker colors are applied. When
// false, the same layout is emitted with every ANSI escape sequence
// suppressed — output stays terminal-safe for redirects and pipes.
func renderChangelogEntry(entry release.ChangelogEntry, pretty bool) {
	bullets := selectChangelogBullets(entry)
	printChangelogHeader(entry, pretty)
	printChangelogBullets(bullets, pretty)
	fmt.Println()
}

// selectChangelogBullets prefers the structured Bullets slice, falling
// back to the flat Notes slice for legacy callers / tests.
func selectChangelogBullets(entry release.ChangelogEntry) []release.ChangelogBullet {
	if len(entry.Bullets) > 0 {
		return entry.Bullets
	}

	out := make([]release.ChangelogBullet, 0, len(entry.Notes))
	for _, note := range entry.Notes {
		out = append(out, release.ChangelogBullet{Depth: 0, Marker: "-", Text: note})
	}

	return out
}

// printChangelogHeader prints the rule + version + title block. When
// pretty is false, every ANSI color slot collapses to "" so the layout
// is preserved without escape codes.
func printChangelogHeader(entry release.ChangelogEntry, pretty bool) {
	dim := colorOrEmpty(constants.ColorDim, pretty)
	cyan := colorOrEmpty(constants.ColorCyan, pretty)
	white := colorOrEmpty(constants.ColorWhite, pretty)
	reset := colorOrEmpty(constants.ColorReset, pretty)

	fmt.Println()
	fmt.Printf("  %s%s%s\n", dim, constants.ChangelogPrettyRule, reset)
	if len(entry.Title) > 0 {
		fmt.Printf(constants.ChangelogPrettyHeaderFmt,
			cyan, entry.Version, reset,
			dim+"  •  "+reset,
			white, entry.Title, reset)
	} else {
		fmt.Printf(constants.ChangelogPrettyHeaderBare,
			cyan, entry.Version, reset)
	}
	fmt.Printf("  %s%s%s\n", dim, constants.ChangelogPrettyRule, reset)
}

// printChangelogBullets renders each bullet with depth-aware styling.
func printChangelogBullets(bullets []release.ChangelogBullet, pretty bool) {
	width := changelogWrapWidth()
	for i := range bullets {
		printChangelogBullet(bullets[i], width, pretty)
	}
}

// printChangelogBullet renders a single bullet with hanging indent.
func printChangelogBullet(bullet release.ChangelogBullet, wrapWidth int, pretty bool) {
	indent := changelogIndent(bullet.Depth)
	marker := changelogMarker(bullet)
	color := colorOrEmpty(changelogMarkerColor(bullet.Depth), pretty)
	reset := colorOrEmpty(constants.ColorReset, pretty)
	prefix := fmt.Sprintf("  %s%s%s%s ", indent, color, marker, reset)
	hanging := "  " + indent + repeatSpace(visibleLen(marker)+1)

	body := renderInlineMarkdown(bullet.Text, bullet.Depth, pretty)
	wrapped := wrapWithHangingIndent(body, prefix, hanging, wrapWidth)
	fmt.Print(wrapped)
}

// colorOrEmpty returns the ANSI color string when pretty is true,
// otherwise the empty string. Centralizes the "strip ANSI" toggle so
// every header / marker site stays consistent — adding a new color slot
// only needs this helper, not a new boolean threading.
func colorOrEmpty(c string, pretty bool) string {
	if pretty {
		return c
	}

	return ""
}

// changelogIndent returns the leading indent string for the given depth.
func changelogIndent(depth int) string {
	out := ""
	for i := 0; i < depth; i++ {
		out += constants.ChangelogPrettyIndentUnit
	}

	return out
}

// changelogMarker returns the bullet glyph or ordered-list marker.
func changelogMarker(bullet release.ChangelogBullet) string {
	if bullet.Ordered {
		return bullet.Marker
	}
	if bullet.Depth == 0 {
		return constants.ChangelogPrettyMarkerL0
	}
	if bullet.Depth == 1 {
		return constants.ChangelogPrettyMarkerL1
	}

	return constants.ChangelogPrettyMarkerLN
}

// changelogMarkerColor selects a color based on bullet depth.
func changelogMarkerColor(depth int) string {
	if depth == 0 {
		return constants.ColorGreen
	}
	if depth == 1 {
		return constants.ColorCyan
	}

	return constants.ColorDim
}
