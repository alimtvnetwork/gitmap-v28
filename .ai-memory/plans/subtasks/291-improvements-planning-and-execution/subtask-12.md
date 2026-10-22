# Subtask 12: Implement TokyoMetricPill Count Formatter

## Status: Completed
## Parent Plan: .ai-memory/plans/291-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/291-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/theme_lipgloss.go`

---

## 1. Objective

Implement `TokyoMetricPill(label string, count int) string` in `cli/termout/theme_lipgloss.go` to provide a standardized badge pill displaying a metric label and numeric counter with semantic Tokyo Night coloring.

---

## 2. Scope & Implementation Details

1. **Function Signature**:
   - `TokyoMetricPill(label string, count int) string`

2. **Formatting Logic**:
   - Guard Clause: When `!colorAllowed()`:
     - Return plain fallback: `fmt.Sprintf("[%s: %d]", label, count)`.
   - Semantic Color Resolution (`resolveMetricColor(label string, count int) string`):
     - If label matches `"healthy"`, `"pass"`, or `"ok"` (case-insensitive): return `TokyoOK`.
     - If `count > 0` and label matches `"issues"`, `"errors"`, or `"fail"` (case-insensitive): return `TokyoWarn`.
     - Otherwise: return `TokyoPrimary`.
   - Lipgloss Styling:
     - Foreground: resolved semantic hex color.
     - Background: `TokyoSurface`.
     - Horizontal padding: 1 space.
     - Text content: `fmt.Sprintf("%s: %d", label, count)`.

3. **Sub-15-Line Function Decomposition**:
   - Decompose into `resolveMetricColor(label string, count int) string` and `TokyoMetricPill(label string, count int) string`.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines of executable code.
- **Zero Nested Conditionals**: Invert conditions and use early returns.
- **Positive Booleans**: Use positive checks (`colorAllowed()`, `isWarning`).
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `TokyoMetricPill(label string, count int) string` exported in `cli/termout/theme_lipgloss.go`.
- [x] Safe plain fallback `[label: count]` returned when colors are disabled.
- [x] Semantic coloring applied for success, warning, and default metrics.
- [x] Functions strictly $\le 15$ lines with zero nested `if` statements.
