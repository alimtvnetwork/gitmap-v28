# Subtask 11: Implement TokyoTag Lipgloss Token in Theme System

## Status: Completed
## Parent Plan: .ai-memory/plans/291-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/291-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/theme_lipgloss.go`

---

## 1. Objective

Implement `TokyoTag(label, value string) string` in `cli/termout/theme_lipgloss.go` to provide a standardized, high-contrast, dual-pill tag badge formatter using Tokyo Night design tokens, with guaranteed non-TTY fallback and strict coding guideline compliance.

---

## 2. Scope & Implementation Details

1. **Function Signature**:
   - `TokyoTag(label, value string) string`

2. **Formatting Logic**:
   - Guard Clause: When `!colorAllowed()`:
     - If `value == ""`, return plain unstyled `fmt.Sprintf("[%s]", label)`.
     - Otherwise, return plain unstyled `fmt.Sprintf("[%s: %s]", label, value)`.
   - When color is supported:
     - Render `label` using Lipgloss with `TokyoFG` foreground, `TokyoSurface` background, bold font weight, and horizontal padding of 1 space.
     - If `value == ""`, return the rendered label tag.
     - Render `value` using Lipgloss with `TokyoPrimary` foreground, `TokyoBase` background, normal font weight, and horizontal padding of 1 space.
     - Return the combined string of both styled elements.

3. **Decomposition Strategy**:
   - Decompose into small helper functions if needed to keep each function $\le 15$ lines of executable code.
   - Example helper: `formatColoredTag(label, value string) string`.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines of executable code.
- **Zero Nested Conditionals**: Flatten all branches with guard clauses and early returns.
- **Positive Booleans**: Strict adherence to `colorAllowed()` or `isColorSupported`.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `TokyoTag(label, value string) string` exported in `cli/termout/theme_lipgloss.go`.
- [x] Non-TTY or no-color environments safely return unstyled fallback `[label: value]`.
- [x] Handles empty `value` gracefully by returning single label pill.
- [x] Styled version uses `TokyoFG`, `TokyoSurface`, `TokyoPrimary`, and `TokyoBase`.
- [x] All functions $\le 15$ lines without nested `if` statements.
