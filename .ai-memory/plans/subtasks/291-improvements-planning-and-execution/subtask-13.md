# Subtask 13: Implement FormatSectionDivider and PrintSectionDivider

## Status: Completed
## Parent Plan: .ai-memory/plans/291-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/291-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/render.go`

---

## 1. Objective

Implement `FormatSectionDivider(title string, theme Theme) string` and `PrintSectionDivider(title string, theme Theme)` in `cli/termout/render.go` to provide uniform subtle horizontal section breaks across CLI help menus, diagnostics, and command summaries.

---

## 2. Scope & Implementation Details

1. **Signatures**:
   - `FormatSectionDivider(title string, theme Theme) string`
   - `PrintSectionDivider(title string, theme Theme)`

2. **Formatting Logic**:
   - Standard divider total width: 64 characters.
   - Divider character: UTF-8 light horizontal box-drawing glyph `─` (`\u2500`).
   - Line color: resolved from `theme.HeaderColor` or `RoleSGR(RoleFGSubtle)`.
   - Empty Title Case:
     - When `strings.TrimSpace(title) == ""`, return a solid 64-character horizontal line: `strings.Repeat("─", 64)` styled with line color.
   - Non-Empty Title Case:
     - Leading rule: `─── ` (3 glyphs + 1 space).
     - Title block: styled with `theme.HeaderColor` or `RoleSGRBold(RolePrimary)`.
     - Trailing rule: fill remaining width up to 64 columns with ` ─...` (1 space + trailing `─` glyphs).
   - Plain Fallback:
     - When `!colorAllowed()`, strip ANSI formatting and return pure text divider.
   - `PrintSectionDivider`:
     - Calls `fmt.Fprintln(output.UI(), FormatSectionDivider(title, theme))`.

3. **Decomposition**:
   - `buildDividerLine(title string, lineWidth int) (string, string, string)`
   - `FormatSectionDivider(title string, theme Theme) string`
   - `PrintSectionDivider(title string, theme Theme)`
   - Each function $\le 15$ lines.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines of executable code.
- **Zero Nested Conditionals**: Guard clauses for empty title and no-color checks.
- **Positive Booleans**: Use `colorAllowed()` and positive length checks.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `FormatSectionDivider` and `PrintSectionDivider` exported in `cli/termout/render.go`.
- [x] Handles empty title with continuous subtle rule.
- [x] Formats non-empty title with leading rule, padded title, and trailing rule.
- [x] Non-TTY environments produce clean plain character output without ANSI escapes.
- [x] All functions $\le 15$ lines.
