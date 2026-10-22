package termout

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/output"
)

// RenderMenu prints a formatted, auto-aligned HelpMenu to standard output.
func RenderMenu(menu HelpMenu) {
	theme := DefaultTheme()
	bannerColor := menu.BannerColor
	if bannerColor == "" {
		bannerColor = theme.BannerColor
	}

	PrintBanner(menu.Title, bannerColor)
	PrintUsage(menu.UsageLines)
	colWidth := CalculateMaxCommandWidth(menu.Sections, 26, 44)
	for _, sec := range menu.Sections {
		PrintSection(sec, colWidth, theme)
	}
	printFooter(menu, colWidth, theme)
}

func printFooter(menu HelpMenu, colWidth int, theme Theme) {
	if len(menu.FooterFlags) > 0 {
		PrintFlags(menu.FooterFlags, colWidth, theme)
	}
	for _, tip := range menu.Tips {
		fmt.Fprintf(output.UI(), "  %s💡 Tip: %s%s\n", theme.TipColor, tip, Reset())
	}
	fmt.Fprintln(output.UI())
}

// PrintBanner renders a decorative box banner around the title.
func PrintBanner(title, color string) {
	width := len(title) + 6
	if width < 50 {
		width = 50
	}
	border := strings.Repeat("═", width)
	pad := (width - len(title)) / 2
	leftSpace := strings.Repeat(" ", pad)
	rightSpace := strings.Repeat(" ", width-len(title)-pad)

	fmt.Fprintf(output.UI(), "\n  %s╔%s╗%s\n", color, border, Reset())
	fmt.Fprintf(output.UI(), "  %s║%s%s%s║%s\n", color, leftSpace, title, rightSpace, Reset())
	fmt.Fprintf(output.UI(), "  %s╚%s╝%s\n\n", color, border, Reset())
}

// PrintUsage prints the usage block for the command. Colors come from
// semantic roles (primary for the header, accent for the lines) — never
// literal color codes — so the block follows the active theme.
func PrintUsage(usageLines []string) {
	if len(usageLines) == 0 {
		return
	}
	fmt.Fprintf(output.UI(), "  %sUsage:%s\n", RoleSGRBold(RolePrimary), Reset())
	for _, line := range usageLines {
		fmt.Fprintf(output.UI(), "    %s%s%s\n", RoleSGR(RoleAccent), line, Reset())
	}
	fmt.Fprintln(output.UI())
}

// PrintSection renders a section with its header and entries.
func PrintSection(sec HelpSection, colWidth int, theme Theme) {
	headerColor := sec.Color
	if len(headerColor) == 0 {
		headerColor = theme.HeaderColor
	}
	badge := ""
	if sec.BadgeText != "" {
		badge = formatSectionBadge(sec, theme)
	}
	fmt.Fprintf(output.UI(), "  %s%s%s%s:\n", headerColor, sec.Title, Reset(), badge)
	for _, entry := range sec.Entries {
		printEntryLine(entry, colWidth, theme)
	}
	fmt.Fprintln(output.UI())
}

func formatSectionBadge(sec HelpSection, theme Theme) string {
	badgeColor := sec.BadgeColor
	if badgeColor == "" {
		badgeColor = RoleSGR(RoleFGSubtle)
	}

	return fmt.Sprintf(" %s[%s]%s", badgeColor, sec.BadgeText, Reset())
}

func printEntryLine(entry CommandEntry, colWidth int, theme Theme) {
	cmdStr := formatCommandToken(entry, theme)
	padLen := colWidth - len(entry.Command)
	if entry.HasSubcommands {
		padLen -= 2
	}
	if padLen < 2 {
		padLen = 2
	}
	padding := strings.Repeat(" ", padLen)
	fmt.Fprintf(output.UI(), "    %s%s%s%s%s\n",
		cmdStr, padding,
		theme.DescColor, entry.Description, Reset())
}

func resolveSubcommandIndicator(hint string) string {
	if len(hint) > 0 {
		return hint
	}

	return DefaultSubcommandIndicator
}

func formatCommandToken(entry CommandEntry, theme Theme) string {
	if !entry.HasSubcommands {
		return fmt.Sprintf("%s%s%s", theme.CommandColor, entry.Command, Reset())
	}
	sym := resolveSubcommandIndicator(entry.SubcommandHint)
	indicator := fmt.Sprintf(" %s%s%s", theme.HintColor, sym, Reset())

	return fmt.Sprintf("%s%s%s%s", theme.CommandColor, entry.Command, Reset(), indicator)
}

// PrintFlags renders standard flags section.
func PrintFlags(flags []CommandEntry, colWidth int, theme Theme) {
	fmt.Fprintf(output.UI(), "  %sFlags:%s\n", RoleSGRBold(RolePrimary), Reset())
	for _, entry := range flags {
		printEntryLine(entry, colWidth, theme)
	}
	fmt.Fprintln(output.UI())
}

// FormatKeyValPadded formats a key-value pair padded to keyWidth with theme styling.
func FormatKeyValPadded(key, val string, keyWidth int, theme Theme) string {
	paddedKey := PadRightMax(key+":", keyWidth)
	if !colorAllowed() {
		return fmt.Sprintf("  %s %s", paddedKey, val)
	}

	keyColor := theme.HeaderColor
	if keyColor == "" {
		keyColor = RoleSGRBold(RolePrimary)
	}

	return fmt.Sprintf("  %s%s%s %s", keyColor, paddedKey, Reset(), val)
}

// PrintKeyValPadded prints a formatted key-value pair directly to output.UI().
func PrintKeyValPadded(key, val string, keyWidth int, theme Theme) {
	fmt.Fprintln(output.UI(), FormatKeyValPadded(key, val, keyWidth, theme))
}

func resolveKeyColor(theme Theme) string {
	if theme.HeaderColor != "" {
		return theme.HeaderColor
	}

	return RoleSGRBold(RolePrimary)
}

func buildPaddedKeyLabel(key string, keyWidth int) string {
	keyLabel := key + ":"
	padLen := keyWidth - VisualStringWidth(keyLabel)
	if padLen < 1 {
		padLen = 1
	}

	return keyLabel + SpacePool(padLen)
}

// FormatKeyValuePair formats a key-value pair aligned by printable terminal visual width.
func FormatKeyValuePair(key, val string, keyWidth int, theme Theme) string {
	paddedKey := buildPaddedKeyLabel(key, keyWidth)
	if !colorAllowed() {
		return fmt.Sprintf("  %s %s", paddedKey, val)
	}
	keyColor := resolveKeyColor(theme)

	return fmt.Sprintf("  %s%s%s %s", keyColor, paddedKey, Reset(), val)
}

// PrintKeyValuePair prints a visual-width aligned key-value pair directly to output.UI().
func PrintKeyValuePair(key, val string, keyWidth int, theme Theme) {
	fmt.Fprintln(output.UI(), FormatKeyValuePair(key, val, keyWidth, theme))
}


// FormatSectionDivider builds a horizontal subtle line divider with an optional title.
func FormatSectionDivider(title string, theme Theme) string {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return formatEmptyDivider()
	}

	return formatTitledDivider(trimmed, theme)
}

func formatEmptyDivider() string {
	rule := strings.Repeat("─", 64)
	if !colorAllowed() {
		return rule
	}

	return fmt.Sprintf("%s%s%s", RoleSGR(RoleFGSubtle), rule, Reset())
}

func formatTitledDivider(title string, theme Theme) string {
	trailingCount := 64 - 4 - len(title) - 1
	if trailingCount < 3 {
		trailingCount = 3
	}
	trail := strings.Repeat("─", trailingCount)
	if !colorAllowed() {
		return fmt.Sprintf("─── %s %s", title, trail)
	}

	titleColor := resolveDividerTitleColor(theme)
	subtle := RoleSGR(RoleFGSubtle)

	return fmt.Sprintf("%s─── %s%s%s %s%s%s", subtle, titleColor, title, Reset(), subtle, trail, Reset())
}

func resolveDividerTitleColor(theme Theme) string {
	if theme.HeaderColor != "" {
		return theme.HeaderColor
	}

	return RoleSGRBold(RolePrimary)
}

// PrintSectionDivider outputs FormatSectionDivider directly to output.UI().
func PrintSectionDivider(title string, theme Theme) {
	fmt.Fprintln(output.UI(), FormatSectionDivider(title, theme))
}

