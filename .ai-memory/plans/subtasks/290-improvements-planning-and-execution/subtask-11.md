# Subtask 11: Tokyo Night Lipgloss Styling Extensions & Progress Bar

## Status: completed
## Parent Plan: .ai-memory/plans/290-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/290-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/theme_lipgloss.go`

---

## 1. Objective

Enhance `cli/termout/theme_lipgloss.go` with standardized Tokyo Night Lipgloss styling components, status badge pill renderers, and an inline truecolor progress bar renderer with complete TTY and `colorAllowed()` safety checks.

---

## 2. Scope & Implementation Details

1. **Tokyo Framing Styles**:
   - `TokyoBannerStyle`: Lipgloss style with `TokyoPrimary` foreground, bold text, horizontal padding of 1 space.
   - `TokyoBoxStyle`: Lipgloss style with rounded border (`lipgloss.RoundedBorder()`), `TokyoMuted` border color, horizontal padding of 1 space.

2. **Status Badge Pill Helpers**:
   - `TokyoPill(text, fgHex, bgHex string) string`:
     - If `!colorAllowed()`, returns fallback plain string `fmt.Sprintf("[%s]", text)`.
     - Otherwise, applies Lipgloss foreground `fgHex` and background `bgHex` with padding (0, 1).
   - `TokyoOKPill(text string) string`:
     - Returns pill formatted with `TokyoOK` foreground and `TokyoSurface` background.
   - `TokyoWarnPill(text string) string`:
     - Returns pill formatted with `TokyoWarn` foreground and `TokyoSurface` background.
   - `TokyoErrorPill(text string) string`:
     - Returns pill formatted with `TokyoError` foreground and `TokyoSurface` background.

3. **Inline Progress Bar (`RenderTokyoProgress`)**:
   - Signature: `RenderTokyoProgress(current, total, barWidth int) string`
   - Guard conditions:
     - If `total <= 0` or `barWidth <= 0`, return unformatted zero bar: `[░░░░] 0%`.
   - Clamping:
     - If `current > total`, clamp `current = total`.
     - If `current < 0`, clamp `current = 0`.
   - Calculation:
     - `filled := (current * barWidth) / total`
     - `pct := (current * 100) / total`
   - Blocks:
     - Filled block: `█` styled with `TokyoOK` (or `TokyoPrimary`).
     - Remaining block: `░` styled with `TokyoMuted`.
   - Returns string: `[████░░░░] 50%`.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines of executable code.
- **Zero Nested Ifs**: Flatten all branches with guard clauses and early returns.
- **Positive Booleans**: Use `isColorSupported` or `colorAllowed()`.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `TokyoBannerStyle` and `TokyoBoxStyle` exported and correctly configured.
- [x] `TokyoPill`, `TokyoOKPill`, `TokyoWarnPill`, `TokyoErrorPill` exported and check `colorAllowed()`.
- [x] `RenderTokyoProgress` safely handles division by zero (`total <= 0`) and clamps inputs.
- [x] All functions $\le 15$ lines without nested `if` statements.
