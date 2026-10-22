# Subtask 11: Implement TokyoBadge and TokyoHeader in Lipgloss Theme System

## Status: Pending
## Parent Plan: .ai-memory/plans/292-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/292-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/theme_lipgloss.go`

---

## 1. Objective

Implement `TokyoBadge(label, bgHex, fgHex string) string` and `TokyoHeader(title string) string` in `cli/termout/theme_lipgloss.go`. These functions provide reusable Tokyo Night styled badge pills with custom color pairing and prominent bold section header banners, with guaranteed safe unstyled fallback for non-TTY environments and strict adherence to repository coding guidelines.

---

## 2. Scope & Implementation Details

1. **`TokyoBadge(label, bgHex, fgHex string) string`**:
   - Signature: `TokyoBadge(label, bgHex, fgHex string) string`
   - Non-TTY / No-Color Fallback:
     - Guard: If `!colorAllowed()`, return `fmt.Sprintf("[%s]", label)`.
   - Empty input: If `label == ""`, return `""`.
   - Styled Rendering:
     - Construct a Lipgloss style with `Foreground(lipgloss.Color(fgHex))` and `Background(lipgloss.Color(bgHex))`.
     - Apply `Bold(true)` and horizontal padding `Padding(0, 1)`.
     - Render `label` and return the formatted string.
   - Decomposition: Decompose into `formatPlainBadge(label string) string` and `renderStyledBadge(label, bgHex, fgHex string) string` to enforce the $\le 15$-line rule.

2. **`TokyoHeader(title string) string`**:
   - Signature: `TokyoHeader(title string) string`
   - Non-TTY / No-Color Fallback:
     - Guard: If `!colorAllowed()`, return `strings.ToUpper(title)`.
   - Empty input: If `strings.TrimSpace(title) == ""`, return `""`.
   - Styled Rendering:
     - Construct a Lipgloss style with `Bold(true)`, `Foreground(lipgloss.Color(TokyoPrimary))`.
     - Add `MarginBottom(1)`.
     - Return the rendered uppercase title string.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Cap**: Every function must be $\le 15$ lines of executable code ($\le 8$ lines preferred).
- **Zero Nested Conditionals**: All logic must use guard clauses and early returns.
- **Positive Booleans**: Strict use of positive booleans (`colorAllowed()`, `isColorSupported`).
- **Whitespace Rules**: Mandatory blank line after closing brace `}` and blank line before every `return`.
- **Zero Builds & Tests**: Spec and plan execution strictly prohibits running `go test` or `go build`.

---

## 4. Verification Checklist

- [ ] `TokyoBadge(label, bgHex, fgHex string) string` exported in `cli/termout/theme_lipgloss.go`.
- [ ] `TokyoHeader(title string) string` exported in `cli/termout/theme_lipgloss.go`.
- [ ] Safe fallback `[label]` returned when `!colorAllowed()`.
- [ ] Safe uppercase fallback returned for `TokyoHeader` when `!colorAllowed()`.
- [ ] All functions decomposed to $\le 15$ lines with zero nested `if` statements.
- [ ] Blank line after `}` and blank line before `return` strictly preserved.
