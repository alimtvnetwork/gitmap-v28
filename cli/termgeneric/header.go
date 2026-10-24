package termgeneric

import (
	"strings"

	"github.com/mattn/go-runewidth"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/termgeneric/enums"

// BorderStyle defines the boundary characters used when rendering box headers.
type BorderStyle struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
}

// Predefined box border styles.
var (
	BorderRounded = BorderStyle{
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
		Horizontal:  "─",
		Vertical:    "│",
	}

	BorderSharp = BorderStyle{
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
		Horizontal:  "─",
		Vertical:    "│",
	}

	BorderDouble = BorderStyle{
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
		Horizontal:  "═",
		Vertical:    "║",
	}

	BorderAscii = BorderStyle{
		TopLeft:     "+",
		TopRight:    "+",
		BottomLeft:  "+",
		BottomRight: "+",
		Horizontal:  "-",
		Vertical:    "|",
	}
)

// AlignmentType defines the horizontal text alignment within the box.
// AlignmentType aliases the strongly-typed enum from the enums package.
type AlignmentType = enums.AlignmentType

const (
	AlignLeft   = enums.AlignLeft
	AlignCenter = enums.AlignCenter
	AlignRight  = enums.AlignRight
)

// Alignment is a deprecated legacy alias for AlignmentType.
type Alignment = enums.AlignmentType

// HeaderConfig specifies the parameters for box header rendering.
type HeaderConfig struct {
	Title       string
	Subtitle    string
	Width       int
	Padding     int
	Alignment   Alignment
	Border      BorderStyle
	BorderColor string
	TitleColor  string
}

// RenderHeader renders a complete declarative box header string.
func RenderHeader(cfg HeaderConfig) string {
	cfg = normalizeConfig(cfg)
	innerWidth := calculateInnerWidth(cfg.Width)

	var sb strings.Builder
	sb.WriteString(renderTopBorder(cfg, innerWidth) + "\n")
	sb.WriteString(renderContentLine(cfg.Title, cfg, innerWidth, cfg.TitleColor))
	appendSubtitleSection(&sb, cfg, innerWidth)
	sb.WriteString("\n" + renderBottomBorder(cfg, innerWidth))

	return sb.String()
}

func calculateInnerWidth(width int) int {
	innerWidth := width - 2
	if innerWidth < 2 {
		return 2
	}
	return innerWidth
}

func appendSubtitleSection(sb *strings.Builder, cfg HeaderConfig, innerWidth int) {
	hasSubtitle := len(cfg.Subtitle) > 0
	if hasSubtitle {
		sb.WriteString("\n")
		sb.WriteString(renderContentLine(cfg.Subtitle, cfg, innerWidth, ""))
	}
}

// RenderBox renders a standard rounded box around title with given width.
func RenderBox(title string, width int) string {
	cfg := HeaderConfig{
		Title:     title,
		Width:     width,
		Border:    BorderRounded,
		Alignment: AlignLeft,
		Padding:   2,
	}

	return RenderHeader(cfg)
}

// RenderBoxWithSubtitle renders a box containing title and subtitle.
func RenderBoxWithSubtitle(title string, subtitle string, width int) string {
	cfg := HeaderConfig{
		Title:     title,
		Subtitle:  subtitle,
		Width:     width,
		Border:    BorderRounded,
		Alignment: AlignLeft,
		Padding:   2,
	}

	return RenderHeader(cfg)
}

func normalizeConfig(cfg HeaderConfig) HeaderConfig {
	isDefaultBorder := len(cfg.Border.TopLeft) == 0
	if isDefaultBorder {
		cfg.Border = BorderRounded
	}

	isZeroPadding := cfg.Padding <= 0
	if isZeroPadding {
		cfg.Padding = 2
	}

	isAutoWidth := cfg.Width <= 0
	if isAutoWidth {
		cfg.Width = calculateAutoWidth(cfg)
	}

	return cfg
}

func calculateAutoWidth(cfg HeaderConfig) int {
	titleVis := VisibleWidth(cfg.Title)
	subVis := VisibleWidth(cfg.Subtitle)
	maxVis := titleVis
	if subVis > maxVis {
		maxVis = subVis
	}

	calcWidth := maxVis + (cfg.Padding * 2) + 2
	if calcWidth < 40 {
		return 40
	}

	return calcWidth
}

func renderTopBorder(cfg HeaderConfig, innerWidth int) string {
	raw := cfg.Border.TopLeft + strings.Repeat(cfg.Border.Horizontal, innerWidth) + cfg.Border.TopRight
	return colorize(raw, cfg.BorderColor)
}

func renderBottomBorder(cfg HeaderConfig, innerWidth int) string {
	raw := cfg.Border.BottomLeft + strings.Repeat(cfg.Border.Horizontal, innerWidth) + cfg.Border.BottomRight
	return colorize(raw, cfg.BorderColor)
}

func renderContentLine(text string, cfg HeaderConfig, innerWidth int, colorCode string) string {
	availWidth := innerWidth - (cfg.Padding * 2)
	if availWidth < 1 {
		availWidth = 1
	}

	truncated := TruncateTitle(text, availWidth)
	visWidth := VisibleWidth(truncated)
	leftPad, rightPad := computePadding(visWidth, innerWidth, cfg)

	vert := colorize(cfg.Border.Vertical, cfg.BorderColor)
	content := colorize(truncated, colorCode)

	return vert + strings.Repeat(" ", leftPad) + content + strings.Repeat(" ", rightPad) + vert
}

func computePadding(visWidth int, innerWidth int, cfg HeaderConfig) (int, int) {
	remaining := innerWidth - visWidth
	if remaining < 0 {
		remaining = 0
	}

	if cfg.Alignment == AlignCenter {
		return computeCenterPadding(remaining)
	}

	if cfg.Alignment == AlignRight {
		return computeRightPadding(remaining, cfg.Padding)
	}

	return computeLeftPadding(remaining, cfg.Padding)
}

func computeCenterPadding(remaining int) (int, int) {
	leftPad := remaining / 2
	return leftPad, remaining - leftPad
}

func computeRightPadding(remaining int, pad int) (int, int) {
	rightPad := pad
	leftPad := remaining - rightPad
	if leftPad < 0 {
		return 0, remaining
	}
	return leftPad, rightPad
}

func computeLeftPadding(remaining int, pad int) (int, int) {
	leftPad := pad
	rightPad := remaining - leftPad
	if rightPad < 0 {
		return remaining, 0
	}
	return leftPad, rightPad
}

// TruncateTitle shortens title to maxWidth terminal columns if it overflows.
func TruncateTitle(title string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	visWidth := VisibleWidth(title)
	if visWidth <= maxWidth {
		return title
	}

	const ellipsis = "..."
	ellipsisWidth := VisibleWidth(ellipsis)
	if maxWidth <= ellipsisWidth {
		return ellipsis[:maxWidth]
	}

	return truncateWithEllipsis(title, maxWidth-ellipsisWidth, ellipsis)
}

func truncateWithEllipsis(title string, targetWidth int, ellipsis string) string {
	var sb strings.Builder
	currWidth := 0
	clean := StripAnsi(title)

	for _, r := range clean {
		rw := runewidth.RuneWidth(r)
		if currWidth+rw > targetWidth {
			break
		}
		sb.WriteRune(r)
		currWidth += rw
	}

	sb.WriteString(ellipsis)
	return sb.String()
}

// VisibleWidth calculates the display column width of s excluding ANSI escapes.
func VisibleWidth(s string) int {
	clean := StripAnsi(s)
	return runewidth.StringWidth(clean)
}

func colorize(s string, colorCode string) string {
	hasColorCode := len(colorCode) > 0
	if !hasColorCode {
		return s
	}

	hasColor := IsColorEnabled()
	if !hasColor {
		return s
	}

	return colorCode + s + CodeReset
}
