# Subtask 18: Eliminate Local Formatting Routines in History Subsystem

## Status: completed
## Parent Plan: .ai-memory/plans/290-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/290-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmdhistory/history.go`

---

## 1. Objective

Purge redundant private string formatting and colorizing helpers (`padRight`, `truncateHist`, `colorize`) from `cli/cmdhistory/history.go`, migrating all table and row alignment to `termout.PadRightMax`, `termout.TruncateRightWithEllipsis`, and `termout.RoleSGR`.

---

## 2. Scope & Implementation Details

1. **Delete Obsolete Functions**:
   - Delete `padRight(s string, w int) string` (lines 350–356).
   - Delete `truncateHist(s string, w int) string` (lines 358–368).
   - Delete `colorize(color, s string) string` (lines 370–376).

2. **Migrate Call Sites**:
   - In `formatStandardRow(r model.CommandHistoryRecord)` (lines 152–158):
     - Replace `padRight(r.Command, 16)` with `termout.PadRightMax(r.Command, 16)`.
     - Replace `padRight(truncateHist(r.Flags, 22), 22)` with `termout.PadRightMax(termout.TruncateRightWithEllipsis(r.Flags, 22), 22)`.
     - Replace `colorize(constants.ColorYellow, ...)` with `termout.SafeTokyoStyle(termout.TokyoWarn).Render(...)` or `termout.RoleSGR(termout.RoleAccent)`.
   - In `formatFullRow(r model.CommandHistoryRecord)`:
     - Replace `truncateHist(r.Args, 18)` with `termout.TruncateRightWithEllipsis(r.Args, 18)`.
     - Replace `padRight(truncateHist(r.Summary, 30), 30)` with `termout.PadRightMax(termout.TruncateRightWithEllipsis(r.Summary, 30), 30)`.
   - In `colorizedStatus(code int) string` (lines 189–196):
     - Replace manual color string concatenation with `termout.TokyoOKPill("✓ OK")` and `termout.TokyoErrorPill("✗ FAIL")`.

3. **Function Sizing & Coding Rules**:
   - Ensure all refactored row builder functions are $\le 15$ lines.
   - Guard against `nil` or empty history records.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Max 15 lines per function.
- **Zero Nested Ifs**: Guard clauses and early returns only.
- **Leaf Error Sovereignty**: Preserve error returns from `RunHistory`.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] Local functions `padRight`, `truncateHist`, and `colorize` completely deleted.
- [x] Call sites migrated to `termout.PadRightMax` and `termout.TruncateRightWithEllipsis`.
- [x] `colorizedStatus` migrated to `termout.TokyoOKPill` / `TokyoErrorPill`.
- [x] No functions exceed 15 lines.
