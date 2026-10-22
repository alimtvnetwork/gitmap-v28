# Subtask 14: Structured Table Border Framing Engine

## Status: completed
## Parent Plan: .ai-memory/plans/290-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/290-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/table.go`

---

## 1. Objective

Enhance `cli/termout/table.go` to support complete terminal table framing (top border, header separator, row divider, and bottom border) using `BorderStyleType` (`BorderRounded`, `BorderSharp`, `BorderMinimal`), defaulting border colors to subtle Tokyo Night palette tokens (`TokyoMuted` / `TokyoSurface`).

---

## 2. Scope & Implementation Details

1. **Border Glyph Mapping**:
   - Internal struct `tableBorderGlyphs` storing characters for:
     - Top line: `topL`, `topR`, `topT`, `horiz`
     - Mid line: `midL`, `midR`, `cross`, `midT`
     - Bottom line: `botL`, `botR`, `botT`
     - Vertical pipe: `vert`
   - Mapping:
     - `BorderRounded`: `╭`, `╮`, `╰`, `╯`, `├`, `┤`, `┬`, `┴`, `┼`, `─`, `│`
     - `BorderSharp`: `┌`, `┐`, `└`, `┘`, `├`, `┤`, `┬`, `┴`, `┼`, `─`, `│`
     - `BorderMinimal`: separator lines `─` without outer vertical caps.

2. **Framed Table Assembly**:
   - Update `RenderTable(cfg TableConfig) string`:
     - If `cfg.BorderStyle != ""` or `cfg.HasBorders`: invoke framed table rendering pipeline.
     - Top border row: `buildTopBorder(cfg, widths, glyphs)`.
     - Header row with outer vertical borders: `│ Col1 │ Col2 │`.
     - Header separator line with T-crossings: `├──────┼──────┤`.
     - Data rows with vertical column borders: `│ Val1 │ Val2 │`.
     - Bottom border row: `buildBottomBorder(cfg, widths, glyphs)`.

3. **Subtle Tokyo Border Palette**:
   - If `cfg.BorderColor == ""`, default to Tokyo Night subtle border color:
     `termout.SafeTokyoStyle(TokyoMuted).Render("")` or role `RoleFGSubtle`.

4. **Function Decomposition & Sizing**:
   - Decompose into discrete helper functions $\le 15$ lines:
     - `resolveBorderGlyphs(style BorderStyleType) tableBorderGlyphs`
     - `formatFramedHeader(cfg TableConfig, widths []int, g tableBorderGlyphs) string`
     - `formatFramedRow(row Row, widths []int, cfg TableConfig, g tableBorderGlyphs) string`
     - `buildBorderLine(widths []int, left, mid, right, horiz, color string) string`

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines.
- **Zero Nested Ifs**: Flatten branching using guard clauses and early returns.
- **Visual Width Accuracy**: Continue using `VisualStringWidth` to ensure proper column alignment.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `BorderRounded` renders `╭`, `─`, `┬`, `╮` top border and `╰`, `─`, `┴`, `╯` bottom border.
- [x] `BorderSharp` renders `┌`, `─`, `┬`, `┐` top border and `└`, `─`, `┴`, `┘` bottom border.
- [x] `BorderMinimal` renders subtle separators without outer vertical pipes.
- [x] Default border color uses Tokyo Night subtle border tokens when unspecified.
- [x] All functions $\le 15$ lines.
