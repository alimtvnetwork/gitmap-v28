package termtable

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PrintTable formats and outputs a table directly to standard output.
func PrintTable(cfg TableConfig) {
	fmt.Print(RenderTable(cfg))
}

// RenderTable builds a complete formatted table string.
func RenderTable(cfg TableConfig) string {
	if len(cfg.Columns) == 0 {
		return ""
	}
	widths := calculateColumnWidths(cfg)
	var sb strings.Builder

	appendHeaderBlock(&sb, cfg, widths)
	appendRowBlocks(&sb, cfg, widths)

	return sb.String()
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
	width := len(col.Title)
	if col.MinWidth > width {
		return col.MinWidth
	}

	return width
}

func evaluateRowsForColWidth(rows []Row, colIdx, curWidth, maxCap int) int {
	res := curWidth
	for _, r := range rows {
		if colIdx < len(r.Cells) && len(r.Cells[colIdx]) > res {
			res = len(r.Cells[colIdx])
		}
	}
	if maxCap > 0 && res > maxCap {
		return maxCap
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
		cell := formatCell(col.Title, widths[i], col.Align, cfg.EllipsisText)
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

func getColAlign(cols []Column, idx int) AlignType {
	if idx < len(cols) {
		return cols[idx].Align
	}

	return AlignLeft
}

func formatCell(text string, width int, align AlignType, ellipsis string) string {
	truncated := TruncateMiddle(text, width, ellipsis)
	padLen := width - len(truncated)
	if padLen <= 0 {
		return truncated
	}
	if align == AlignRight {
		return strings.Repeat(" ", padLen) + truncated
	}
	if align == AlignCenter {
		left := padLen / 2
		right := padLen - left

		return strings.Repeat(" ", left) + truncated + strings.Repeat(" ", right)
	}

	return truncated + strings.Repeat(" ", padLen)
}
