# Subtask 19: Output Standardization & Status Badges in Pending Commits Render

## Status: completed
## Parent Plan: .ai-memory/plans/290-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/290-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmdpending/pending_commits_render.go`

---

## 1. Objective

Standardize table rendering and box alignment in `cli/cmdpending/pending_commits_render.go` by deleting the private `padRight` helper, rebinding box framing to `termout.PadRightMax`, and replacing hardcoded status color strings with `termout.TokyoOKPill` and `termout.TokyoWarnPill`.

---

## 2. Scope & Implementation Details

1. **Delete Redundant Private Function**:
   - Delete `padRight(s string, width int) string` (lines 141–147).

2. **Rebind Box Padding to `termout.PadRightMax`**:
   - In `renderSummaryLine`: Replace `padRight(summaryText, 78)` with `termout.PadRightMax(summaryText, 78)`.
   - In `renderHeaderLine`: Replace `padRight(headerText, 78)` with `termout.PadRightMax(headerText, 78)`.
   - In `renderFooterLine`: Replace `padRight(footerMsg, 78)` with `termout.PadRightMax(footerMsg, 78)`.
   - In row formatting: Replace column padding with `termout.PadRightMax(repoDisplay, 22)`, `termout.PadRightMax(branchDisplay, 16)`.

3. **Modernize Status Indicators with Tokyo Pills**:
   - In `renderStatusPill(isDirty bool) string` (or `resolveStatusToken`):
     - Replace `constants.ColorYellow + "● PEND   " + reset` with `termout.TokyoWarnPill("PEND")`.
     - Replace `constants.ColorGreen + "○ CLEAN  " + reset` with `termout.TokyoOKPill("CLEAN")`.

4. **Preserve Fixed 80-Column Boundary**:
   - Ensure the outer table container (`│ ... 78 cols ... │`) maintains mathematical precision without horizontal overflow.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines.
- **Zero Nested Ifs**: Guard clauses and early returns only.
- **Affirmative Booleans**: `isDirty`, `isClean`, `hasUncommitted`.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] Private `padRight` function deleted from `cli/cmdpending/pending_commits_render.go`.
- [x] All padding call sites use `termout.PadRightMax`.
- [x] Status indicators formatted with `termout.TokyoWarnPill` and `termout.TokyoOKPill`.
- [x] Box boundaries remain aligned at 80 terminal columns.
- [x] All functions $\le 15$ lines.
