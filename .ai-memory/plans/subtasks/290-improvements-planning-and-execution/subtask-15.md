# Subtask 15: UTF-8 Aware Right-Truncation Engine

## Status: completed
## Parent Plan: .ai-memory/plans/290-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/290-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/truncate.go`

---

## 1. Objective

Implement `TruncateRightWithEllipsis(text string, maxWidth int) string` in `cli/termout/truncate.go` to provide a display-width aware, rune-safe right truncation routine that handles multi-byte UTF-8 sequences and East Asian double-width characters without string corruption or slicing panics.

---

## 2. Scope & Implementation Details

1. **`TruncateRightWithEllipsis` Contract**:
   - Signature: `TruncateRightWithEllipsis(text string, maxWidth int) string`
   - Edge Cases:
     - If `maxWidth <= 0`, return `""`.
     - Compute visual width via `VisualStringWidth(text)`.
     - If `VisualStringWidth(text) <= maxWidth`, return `text` untouched.
   - Ellipsis Strategy:
     - Ellipsis rune: `…` (visual display width 1).
     - If `maxWidth <= 1`, return `"…"`.
     - Target character width: `avail := maxWidth - 1`.
   - Rune Iteration:
     - Iterate through `text` runes, tracking accumulated visual display width.
     - While `curWidth + runeWidth <= avail`, append rune to output builder.
     - Terminate scan early once `avail` is reached.
     - Append `"…"` and return string.

2. **Helper Functions & Decomposition**:
   - Extract rune visual width calculation helper $\le 10$ lines:
     `resolveRuneVisualWidth(r rune) int`
   - Extract string assembly helper $\le 12$ lines:
     `assembleRightTruncation(text string, avail int) string`

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Max 15 lines per function.
- **Zero Nested Ifs**: Guard clauses and early returns only.
- **Rune Safety**: Never perform byte slicing (`text[:n]`) on UTF-8 strings.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `TruncateRightWithEllipsis` exported and callable.
- [x] Returns `""` on `maxWidth <= 0`.
- [x] Returns unmodified string if visual width fits within `maxWidth`.
- [x] Accurately truncates multi-byte UTF-8 runes without panicking or creating broken glyphs.
- [x] Appends `…` cleanly within the exact `maxWidth` visual boundary.
- [x] Functions conform to $\le 15$ lines limit.
