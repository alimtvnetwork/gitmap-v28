package termhelp

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RenderMenu prints a formatted, auto-aligned HelpMenu to standard output.
func RenderMenu(menu HelpMenu) {
	theme := DefaultTheme()
	PrintBanner(menu.Title, theme.BannerColor)
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
		fmt.Printf("  %s💡 Tip: %s%s\n", theme.TipColor, tip, constants.ColorReset)
	}
	fmt.Println()
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

	fmt.Printf("\n  %s╔%s╗%s\n", color, border, constants.ColorReset)
	fmt.Printf("  %s║%s%s%s║%s\n", color, leftSpace, title, rightSpace, constants.ColorReset)
	fmt.Printf("  %s╚%s╝%s\n\n", color, border, constants.ColorReset)
}

// PrintUsage prints the usage block for the command.
func PrintUsage(usageLines []string) {
	if len(usageLines) == 0 {
		return
	}
	fmt.Printf("  %sUsage:%s\n", constants.ColorWhite, constants.ColorReset)
	for _, line := range usageLines {
		fmt.Printf("    %s%s%s\n", constants.ColorCyan, line, constants.ColorReset)
	}
	fmt.Println()
}

// PrintSection renders a section with its header and entries.
func PrintSection(sec HelpSection, colWidth int, theme Theme) {
	headerColor := sec.Color
	if len(headerColor) == 0 {
		headerColor = theme.HeaderColor
	}
	fmt.Printf("  %s%s:%s\n", headerColor, sec.Title, constants.ColorReset)
	for _, entry := range sec.Entries {
		printEntryLine(entry, colWidth, theme)
	}
	fmt.Println()
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
	fmt.Printf("    %s%s%s%s%s\n",
		cmdStr, padding,
		theme.DescColor, entry.Description, constants.ColorReset)
}

func resolveSubcommandIndicator(hint string) string {
	if len(hint) > 0 {
		return hint
	}

	return DefaultSubcommandIndicator
}

func formatCommandToken(entry CommandEntry, theme Theme) string {
	if !entry.HasSubcommands {
		return fmt.Sprintf("%s%s%s", theme.CommandColor, entry.Command, constants.ColorReset)
	}
	sym := resolveSubcommandIndicator(entry.SubcommandHint)
	indicator := fmt.Sprintf(" %s%s%s", theme.HintColor, sym, constants.ColorReset)

	return fmt.Sprintf("%s%s%s%s", theme.CommandColor, entry.Command, constants.ColorReset, indicator)
}

// PrintFlags renders standard flags section.
func PrintFlags(flags []CommandEntry, colWidth int, theme Theme) {
	fmt.Printf("  %sFlags:%s\n", constants.ColorWhite, constants.ColorReset)
	for _, entry := range flags {
		printEntryLine(entry, colWidth, theme)
	}
	fmt.Println()
}
