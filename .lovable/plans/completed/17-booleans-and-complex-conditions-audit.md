# Plan 17: Boolean Principles, Negative Elimination & Complex Conditions Audit

## Objective

Autonomously scan, plan, refactor, and fix all boolean naming, double negatives, mixed polarity, and complex condition violations across the codebase, modifying source files directly to enforce affirmative prefixes (`is`, `has`, `can`, `should`), implicit evaluation (no `== true`), positive framing (no `!isSuccess`, no `hasNo*`, no `isNot*`), and discrete condition decomposition until 100% green without stopping.

## Background & Rationale

1. **Cognitive Clarity:** Negative naming (`hasNoGroup`, `isNotSingleBlock`) forces engineers to reason about double negatives (`!hasNoGroup`).
2. **Standardization:** All boolean flags and variables must represent positive states (`hasGroup`, `isSingleBlock`, `hasColors`, `hasPayload`, `hasError`).
3. **Implicit Checks:** Zero explicit `== true` / `== false` checks.

## Target Violations Inventory

| Id | File | Line | Violation | Refactoring Plan | Status |
|----|------|------|-----------|------------------|--------|
| 1 | `gitmap/cmd/cdops.go` | 181, 183 | `hasNoGroup` | Rename to `hasGroup` with positive logic | Pending |
| 2 | `gitmap/cmd/installctx_linux_e2e_test.go` | 244, 245 | `isNotSingleBlock` | Rename to `isSingleBlock` with inverted guard | Pending |
| 3 | `gitmap/cmd/visibilityresolveowner.go` | 114, 116, 119, 122 | `hasNoSlash` | Rename to `hasSlash` with affirmative check | Pending |
| 4 | `src/components/docs/TabOrderMap.tsx` | 93, 94 | `hasNoRects` | Rename to `hasRects` and guard `if (!hasRects)` | Pending |
| 5 | `src/components/ui/chart.tsx` | 80, 133, 159, 254 | `hasNoColors`, `hasNoPayload` | Rename to `hasColors`, `hasPayload` with positive checks | Pending |
| 6 | `src/components/ui/form.tsx` | 88, 89 | `hasNoError` | Rename to `hasError` with affirmative ternary | Pending |

## Subtasks

- [ ] Subtask 01: Refactor Go files (`gitmap/cmd/cdops.go`, `gitmap/cmd/installctx_linux_e2e_test.go`, `gitmap/cmd/visibilityresolveowner.go`).
- [ ] Subtask 02: Refactor TypeScript / React UI components (`TabOrderMap.tsx`, `chart.tsx`, `form.tsx`).
- [ ] Subtask 03: Connect `check-boolean-guidelines.py` to `.github/workflows/ci.yml` and `03-cicd-local-runner.py`.
- [ ] Subtask 04: Verify 100% green test suite across all 23 quality gates.

## Acceptance Criteria

- [ ] `python linter-scripts/check-boolean-guidelines.py` passes with 0 violations.
- [ ] `python .lovable/ai-fix-scripts/03-cicd-local-runner.py` passes with exit code 0.
