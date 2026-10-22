package termout

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/output"
)

// PrintTable formats and outputs a table directly to standard output.
func PrintTable(cfg TableConfig) {
	fmt.Fprint(output.UI(), RenderTable(cfg))
}

func isFramedTable(cfg TableConfig) bool {
	if cfg.BorderStyle == BorderNone {
		return false
	}

	return cfg.BorderStyle != "" || cfg.HasBorders
}

// RenderTable builds a complete formatted table string.
func RenderTable(cfg TableConfig) string {
	if len(cfg.Columns) == 0 {
		return ""
	}
	widths := calculateColumnWidths(cfg)
	if isFramedTable(cfg) {
		return renderFramedTable(cfg, widths)
	}
	var sb strings.Builder

	appendHeaderBlock(&sb, cfg, widths)
	appendRowBlocks(&sb, cfg, widths)

	return sb.String()
}

type tableBorderGlyphs struct {
	topL     string
	topT     string
	topR     string
	midL     string
	cross    string
	midR     string
	botL     string
	botT     string
	botR     string
	vert     string
	horiz    string
	hasOuter bool
}

func resolveBorderGlyphs(style BorderStyleType) tableBorderGlyphs {
	if style == BorderSharp {
		return tableBorderGlyphs{
			topL: "┌", topT: "┬", topR: "┐",
			midL: "├", cross: "┼", midR: "┤",
			botL: "└", botT: "┴", botR: "┘",
			vert: "│", horiz: "─", hasOuter: true,
		}
	}
	if style == BorderMinimal {
		return tableBorderGlyphs{
			midL: "─", cross: "──", midR: "─",
			botL: "─", botT: "──", botR: "─",
			vert: "  ", horiz: "─", hasOuter: false,
		}
	}

	return tableBorderGlyphs{
		topL: "╭", topT: "┬", topR: "╮",
		midL: "├", cross: "┼", midR: "┤",
		botL: "╰", botT: "┴", botR: "╯",
		vert: "│", horiz: "─", hasOuter: true,
	}
}

func resolveBorderColor(cfg TableConfig) string {
	if cfg.BorderColor != "" {
		return cfg.BorderColor
	}
	if colorAllowed() {
		return RoleSGR(RoleFGSubtle)
	}

	return ""
}

func buildBorderLine(widths []int, left, mid, right, horiz, color string) string {
	var parts []string
	for _, w := range widths {
		parts = append(parts, strings.Repeat(horiz, w+2))
	}
	line := "  " + left + strings.Join(parts, mid) + right
	if color != "" {
		return fmt.Sprintf("%s%s%s", color, line, constants.ColorReset)
	}

	return line
}

func formatFramedHeader(cfg TableConfig, widths []int, g tableBorderGlyphs, bColor string) string {
	var parts []string
	color := cfg.HeaderColor
	if len(color) == 0 {
		color = constants.ColorYellow
	}
	for i, col := range cfg.Columns {
		cell := formatCell(col.Title, widths[i], getColAlign(cfg.Columns, i), cfg.EllipsisText)
		parts = append(parts, fmt.Sprintf("%s%s%s", color, cell, constants.ColorReset))
	}
	vPipe := g.vert
	if bColor != "" {
		vPipe = fmt.Sprintf("%s%s%s", bColor, g.vert, constants.ColorReset)
	}
	if !g.hasOuter {
		return "  " + strings.Join(parts, "  ")
	}

	return "  " + vPipe + " " + strings.Join(parts, " "+vPipe+" ") + " " + vPipe
}

func formatFramedRow(row Row, widths []int, cfg TableConfig, g tableBorderGlyphs, bColor string) string {
	var parts []string
	rowColor := row.Color
	if len(rowColor) == 0 {
		rowColor = constants.ColorWhite
	}
	for i, w := range widths {
		raw := getCellText(row.Cells, i)
		cell := formatCell(raw, w, getColAlign(cfg.Columns, i), cfg.EllipsisText)
		parts = append(parts, fmt.Sprintf("%s%s%s", rowColor, cell, constants.ColorReset))
	}
	vPipe := g.vert
	if bColor != "" {
		vPipe = fmt.Sprintf("%s%s%s", bColor, g.vert, constants.ColorReset)
	}
	if !g.hasOuter {
		return "  " + strings.Join(parts, "  ")
	}

	return "  " + vPipe + " " + strings.Join(parts, " "+vPipe+" ") + " " + vPipe
}

func renderFramedTable(cfg TableConfig, widths []int) string {
	g := resolveBorderGlyphs(cfg.BorderStyle)
	bColor := resolveBorderColor(cfg)
	var sb strings.Builder

	if g.hasOuter {
		sb.WriteString(buildBorderLine(widths, g.topL, g.topT, g.topR, g.horiz, bColor))
		sb.WriteString("\n")
	}
	sb.WriteString(formatFramedHeader(cfg, widths, g, bColor))
	sb.WriteString("\n")
	sb.WriteString(buildBorderLine(widths, g.midL, g.cross, g.midR, g.horiz, bColor))
	sb.WriteString("\n")
	appendFramedRows(&sb, cfg, widths, g, bColor)
	if g.hasOuter {
		sb.WriteString(buildBorderLine(widths, g.botL, g.botT, g.botR, g.horiz, bColor))
		sb.WriteString("\n")
	}

	return sb.String()
}

func appendFramedRows(sb *strings.Builder, cfg TableConfig, widths []int, g tableBorderGlyphs, bColor string) {
	for _, row := range cfg.Rows {
		sb.WriteString(formatFramedRow(row, widths, cfg, g, bColor))
		sb.WriteString("\n")
	}
}

func appendHeaderBlock(sb *strings.Builder, cfg TableConfig, widths []int) {
	sb.WriteString(formatHeader(cfg, widths))
	sb.WriteString("\n")
	sb.WriteString(formatSeparator(widths, cfg.BorderColor))
	sb.WriteString("\n")
}

func appendRowBlocks(sb *strings.Builder, cfg TableConfig, widths []int) {
	for _, row := range cfg.Rows {
		sb.WriteString(formatRow(row, widths, cfg))
		sb.WriteString("\n")
	}
}

func calculateColumnWidths(cfg TableConfig) []int {
	widths := make([]int, len(cfg.Columns))
	for i, col := range cfg.Columns {
		widths[i] = resolveInitialColWidth(col)
		widths[i] = evaluateRowsForColWidth(cfg.Rows, i, widths[i], col.MaxWidth)
	}

	return widths
}

func resolveInitialColWidth(col Column) int {
	width := VisualStringWidth(col.Title)
	if col.MinWidth > width {
		return col.MinWidth
	}

	return width
}

func evaluateRowsForColWidth(rows []Row, colIdx, curWidth, maxCap int) int {
	res := maxRowWidth(rows, colIdx, curWidth)
	if maxCap > 0 && res > maxCap {
		return maxCap
	}

	return res
}

func maxRowWidth(rows []Row, colIdx, curWidth int) int {
	res := curWidth
	for _, r := range rows {
		if colIdx < len(r.Cells) && VisualStringWidth(r.Cells[colIdx]) > res {
			res = VisualStringWidth(r.Cells[colIdx])
		}
	}

	return res
}

func formatHeader(cfg TableConfig, widths []int) string {
	var parts []string
	color := cfg.HeaderColor
	if len(color) == 0 {
		color = constants.ColorYellow
	}
	for i, col := range cfg.Columns {
		cell := formatCell(col.Title, widths[i], getColAlign(cfg.Columns, i), cfg.EllipsisText)
		parts = append(parts, fmt.Sprintf("%s%s%s", color, cell, constants.ColorReset))
	}

	return "  " + strings.Join(parts, "  ")
}

func formatSeparator(widths []int, borderColor string) string {
	var parts []string
	if len(borderColor) == 0 {
		borderColor = constants.ColorDim
	}
	for _, w := range widths {
		parts = append(parts, strings.Repeat("─", w))
	}

	return fmt.Sprintf("  %s%s%s", borderColor, strings.Join(parts, "──"), constants.ColorReset)
}

func formatRow(row Row, widths []int, cfg TableConfig) string {
	var parts []string
	rowColor := row.Color
	if len(rowColor) == 0 {
		rowColor = constants.ColorWhite
	}
	for i, w := range widths {
		raw := getCellText(row.Cells, i)
		cell := formatCell(raw, w, getColAlign(cfg.Columns, i), cfg.EllipsisText)
		parts = append(parts, fmt.Sprintf("%s%s%s", rowColor, cell, constants.ColorReset))
	}

	return "  " + strings.Join(parts, "  ")
}

func getCellText(cells []string, idx int) string {
	if idx < len(cells) {
		return cells[idx]
	}

	return ""
}

func getColAlign(cols []Column, idx int) AlignmentType {
	if idx >= len(cols) {
		return AlignLeft
	}
	if cols[idx].Alignment != "" {
		return cols[idx].Alignment
	}
	if cols[idx].Align != "" {
		return cols[idx].Align
	}

	return AlignLeft
}

func formatCell(text string, width int, align AlignmentType, ellipsis string) string {
	truncated := TruncateMiddle(text, width, ellipsis)
	padLen := width - VisualStringWidth(truncated)
	if padLen <= 0 {
		return truncated
	}

	return alignPaddedCell(truncated, padLen, align)
}

func alignPaddedCell(truncated string, padLen int, align AlignmentType) string {
	if align == AlignRight {
		return SpacePool(padLen) + truncated
	}
	if align == AlignCenter {
		return alignCenterCell(truncated, padLen)
	}

	return truncated + SpacePool(padLen)
}

func alignCenterCell(truncated string, padLen int) string {
	left := padLen / 2
	right := padLen - left

	return SpacePool(left) + truncated + SpacePool(right)
}

// VisualStringWidth calculates the printable terminal column width of string s,
// stripping ANSI escape codes and accounting for double-width runes.
func VisualStringWidth(s string) int {
	plain := StripAnsi(s)
	width := 0
	for _, r := range plain {
		if isDoubleWidthRune(r) {
			width += 2
			continue
		}
		width++
	}

	return width
}

func isDoubleWidthRune(r rune) bool {
	return (r >= 0x1100 && r <= 0x115F) ||
		(r >= 0x2E80 && r <= 0xA4CF) ||
		(r >= 0xAC00 && r <= 0xD7A3) ||
		(r >= 0xF900 && r <= 0xFAFF) ||
		(r >= 0xFE10 && r <= 0xFE19) ||
		(r >= 0xFE30 && r <= 0xFE6F) ||
		(r >= 0xFF00 && r <= 0xFF60) ||
		(r >= 0xFFE0 && r <= 0xFFE6) ||
		(r >= 0x1F300 && r <= 0x1F64F) ||
		(r >= 0x1F680 && r <= 0x1F6FF)
}
