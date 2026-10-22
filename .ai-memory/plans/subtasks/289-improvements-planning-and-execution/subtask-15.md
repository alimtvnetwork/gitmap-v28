# Subtask 15: Table Visual String Width & Ellipsis Truncation

> **Subtask ID:** Subtask-15  
> **Parent Plan:** `.ai-memory/plans/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/02-component-spec.md`  
> **Target Path:** `cli/termout/table.go`  
> **Status:** COMPLETED  

---

## 1. Objectives & Scope
1. Enhance `cli/termout/table.go` to compute column widths and alignment using visual string width rather than raw byte counts (`len(s)`).
2. Account for ANSI escape sequences (ignoring escape bytes) and wide unicode characters (emojis, CJK glyphs).
3. Implement `VisualStringWidth` helper using `StripAnsi` and rune categorization.
4. Update `resolveInitialColWidth`, `evaluateRowsForColWidth`, and `formatCell` to use visual width.
5. Retain robust middle and end ellipsis truncation.

---

## 2. Detailed Technical Plan

### 2.1 Implement `VisualStringWidth` Helper
- Place in `cli/termout/table.go` or reference shared width helper:
  ```go
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
  ```

### 2.2 Update Column Width Calculation
- Modify `resolveInitialColWidth`:
  ```go
  func resolveInitialColWidth(col Column) int {
  	width := VisualStringWidth(col.Title)
  	if col.MinWidth > width {
  		return col.MinWidth
  	}

  	return width
  }
  ```
- Modify `evaluateRowsForColWidth`:
  ```go
  func evaluateRowsForColWidth(rows []Row, colIdx, curWidth, maxCap int) int {
  	res := curWidth
  	for _, r := range rows {
  		if colIdx < len(r.Cells) {
  			w := VisualStringWidth(r.Cells[colIdx])
  			if w > res {
  				res = w
  			}
  		}
  	}
  	if maxCap > 0 && res > maxCap {
  		return maxCap
  	}

  	return res
  }
  ```

### 2.3 Update Cell Padding & Alignment in `formatCell`
- Use `VisualStringWidth(truncated)` to determine required space padding:
  ```go
  func formatCell(text string, width int, align AlignType, ellipsis string) string {
  	truncated := TruncateMiddle(text, width, ellipsis)
  	visWidth := VisualStringWidth(truncated)
  	padLen := width - visWidth
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
  ```

---

## 3. Strict Coding Guidelines Compliance
- **Function Limit:** All modified functions remain <= 15 lines.
- **Zero Nested Ifs:** Branches flattened with early returns and guard checks.
- **Rune-safe Traversal:** Ranging over strings iterates by unicode runes, not raw bytes.
- **Return Spacing:** Single blank line before each return statement.

---

## 4. Verification Criteria
- [x] Cells with emojis (e.g., `✔`, `💡`, `📋`) do not push subsequent columns out of vertical alignment.
- [x] ANSI escape codes inside cells do not distort column width calculations.
- [x] Middle-truncation preserves string aesthetics when cells exceed `MaxWidth`.
- [x] Left, right, and center alignment remain visually accurate.
