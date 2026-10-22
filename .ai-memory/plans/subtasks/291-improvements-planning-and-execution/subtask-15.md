# Subtask 15: Implement Column Alignment Rendering in Table Engine

## Status: Completed
## Parent Plan: .ai-memory/plans/291-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/291-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/table.go`

---

## 1. Objective

Enhance `cli/termout/table.go` to support robust, visual-width-aware column text alignment (`AlignLeft`, `AlignRight`, `AlignCenter`) across both header and row cells using `AlignmentType` and `VisualStringWidth`.

---

## 2. Scope & Implementation Details

1. **Alignment Resolution Helper**:
   - Update `getColAlign(cols []Column, idx int) AlignmentType`:
     - Checks if `idx < len(cols)`.
     - Returns `cols[idx].Alignment` if set, else `cols[idx].Align` if set, defaulting to `AlignLeft`.

2. **Visual Cell Alignment in `formatCell`**:
   - Signature: `formatCell(text string, width int, align AlignmentType, ellipsis string) string`.
   - Middle-truncate text if exceeding `width`.
   - Calculate visual width via `VisualStringWidth(truncated)`.
   - Calculate required space padding: `padLen := width - visWidth`.
   - If `padLen <= 0`, return `truncated`.
   - When `align == AlignRight`:
     - Return `SpacePool(padLen) + truncated`.
   - When `align == AlignCenter`:
     - `left := padLen / 2`
     - `right := padLen - left`
     - Return `SpacePool(left) + truncated + SpacePool(right)`.
   - When `align == AlignLeft` (default):
     - Return `truncated + SpacePool(padLen)`.

3. **Header Alignment**:
   - Update `formatHeader` and `formatFramedHeader` to align column title cells according to the column's configured `AlignmentType`.

4. **Decomposition**:
   - Keep `formatCell`, `alignCenterCell`, and alignment resolution helpers $\le 15$ lines each.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines of executable code.
- **Zero Nested Conditionals**: Flatten all branches with guard clauses and early returns.
- **Positive Booleans**: Strict adherence to positive logic (`padLen <= 0`).
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `getColAlign` correctly reads `col.Alignment` and `col.Align`.
- [x] `formatCell` handles `AlignRight`, `AlignCenter`, and `AlignLeft` with exact visual spacing.
- [x] ANSI escape codes stripped during width calculation via `VisualStringWidth`.
- [x] Headers respect column alignment rules.
- [x] All functions $\le 15$ lines.
