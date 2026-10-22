# Subtask 16: Implement Rune-Safe Visual TruncateMiddleWithEllipsis

## Status: Completed
## Parent Plan: .ai-memory/plans/291-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/291-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/truncate.go`

---

## 1. Objective

Implement `TruncateMiddleWithEllipsis(text string, maxWidth int) string` in `cli/termout/truncate.go` to deliver visual-width-aware, rune-safe middle truncation with `"..."` ellipsis, preventing slice index panics and UTF-8 multi-byte character corruption.

---

## 2. Scope & Implementation Details

1. **Function Signature**:
   - `TruncateMiddleWithEllipsis(text string, maxWidth int) string`

2. **Algorithm & Edge Cases**:
   - Guard Conditions:
     - If `maxWidth <= 0`, return `""`.
     - If `VisualStringWidth(text) <= maxWidth`, return `text`.
     - If `maxWidth <= 3`, return `"..."[:maxWidth]`.
   - Budget Calculations:
     - Ellipsis string: `"..."` (visual width 3).
     - Available visual budget: `avail := maxWidth - 3`.
     - Prefix visual width: `prefixBudget := (avail + 1) / 2`.
     - Suffix visual width: `suffixBudget := avail - prefixBudget`.
   - Rune Slicing Strategy:
     - Convert text to `[]rune`.
     - `slicePrefixRunes(runes []rune, maxVisualWidth int) string`: accumulates runes from index 0 while total visual width does not exceed `maxVisualWidth`.
     - `sliceSuffixRunes(runes []rune, maxVisualWidth int) string`: accumulates runes from the end while total visual width does not exceed `maxVisualWidth`.
     - Return `prefix + "..." + suffix`.

3. **Decomposition**:
   - `TruncateMiddleWithEllipsis(text string, maxWidth int) string`
   - `slicePrefixRunes(runes []rune, budget int) string`
   - `sliceSuffixRunes(runes []rune, budget int) string`
   - Each helper $\le 15$ lines of executable code.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines of executable code.
- **Zero Nested Conditionals**: Invert conditions and use early returns.
- **Rune-Safety**: Never index raw byte slices (`text[i]`) on strings that may contain UTF-8 runes.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `TruncateMiddleWithEllipsis(text string, maxWidth int) string` exported in `cli/termout/truncate.go`.
- [x] Safe handling for `maxWidth <= 0`, `maxWidth <= 3`, and `VisualStringWidth(text) <= maxWidth`.
- [x] Preserves multi-byte UTF-8 characters without splitting runes.
- [x] Symmetrically distributes prefix and suffix widths around `"..."`.
- [x] All functions $\le 15$ lines.
