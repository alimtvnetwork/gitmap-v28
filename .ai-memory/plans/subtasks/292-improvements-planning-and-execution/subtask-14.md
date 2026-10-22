# Subtask 14: Implement Table Cell Wrap Modes and Row Striping

## Status: Pending
## Parent Plan: .ai-memory/plans/292-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/292-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/table.go`

---

## 1. Objective

Enhance the table rendering engine in `cli/termout/table.go` to support configurable cell wrap modes (`WrapClip`, `WrapEllipsis`, `WrapWord`) and alternating row striping (`RowStripe bool`) with subtle Tokyo surface background tinting, while decomposing row formatting functions to satisfy the $\le 15$-line function limit.

---

## 2. Scope & Implementation Details

1. **Cell Wrap Strategy Resolution**:
   - Implement `formatCellWithWrap(text string, width int, align AlignmentType, wrap WrapModeType, ellipsis string) string`.
   - Wrapping Logic:
     - If `wrap == WrapClip`: perform hard truncation to `width` visual characters without appending ellipsis.
     - If `wrap == WrapEllipsis` or `wrap == ""`: invoke `TruncateMiddle(text, width, ellipsis)`.
     - Align the resulting truncated string using `alignPaddedCell`.
   - Update `formatCell` to route through wrap mode resolution.

2. **Wrap Mode Precedence**:
   - Helper `resolveColWrap(col ColumnConfig, defaultWrap WrapModeType) WrapModeType`:
     - If `col.WrapMode != ""`, return `col.WrapMode`.
     - If `defaultWrap != ""`, return `defaultWrap`.
     - Return `WrapEllipsis`.

3. **Alternating Row Striping**:
   - In `appendFramedRows` and `appendRowBlocks`, track row index `rowIdx`.
   - Helper `resolveRowBackground(rowIdx int, cfg TableConfig) string`:
     - If `!cfg.RowStripe || !colorAllowed() || rowIdx%2 == 0`: return `""`.
     - Return ANSI background code for subtle Tokyo surface tint (`\033[48;2;31;35;53m` or `SafeTokyoStyle(TokyoSurface)`).
   - In `formatFramedRow` and `formatRow`, apply row background prefix and reset suffix when row striping is active.

4. **Decomposition & Refactoring**:
   - Ensure all new and modified functions remain strictly $\le 15$ lines.
   - Separate row background wrapping from row cell concatenation.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Line Limit**: Every function must be $\le 15$ lines of executable code ($\le 8$ lines preferred).
- **Zero Nested Ifs**: Invert conditions to eliminate all nested conditional branches.
- **Positive Booleans**: Use `isStriped`, `hasOuter`, `colorAllowed()`.
- **Vertical Spacing**: Blank line after closing brace `}` and blank line before every `return`.
- **Zero Builds & Tests**: Strict prohibition on running `go test` or `go build`.

---

## 4. Verification Checklist

- [ ] `formatCellWithWrap` correctly implements `WrapClip` and `WrapEllipsis`.
- [ ] Column-level `WrapMode` takes precedence over table-level `WrapMode`.
- [ ] `RowStripe: true` renders subtle background on odd rows when color is allowed.
- [ ] Even rows and non-TTY outputs remain unstriped without ANSI leakage.
- [ ] All functions decomposed to $\le 15$ lines.
