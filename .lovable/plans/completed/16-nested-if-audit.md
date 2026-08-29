# Plan 16: Nested `if` Elimination & Guard Clause Audit

## Objective

Autonomously scan, plan, refactor, and fix all nested `if` statements across the codebase, directly modifying source files to flatten conditional branching with guard clauses, early returns, inverted conditions, and function decomposition (<= 8–15 lines) until 100% green without stopping.

## Background & Rationale

1. **Cognitive Load Explosion:** Nested conditionals create multiple execution paths that make code harder to reason about and maintain.
2. **Hidden Invariant Bugs:** Deeply nested branches frequently conceal edge case regressions and unhandled errors.
3. **Canonical Size Tier Compliance:** Flattening conditionals enables functions to satisfy the mandatory <= 8-line preferred (<= 15 lines max) function limit.

## Target Violations Inventory

| Id | File | Line | Snippet | Planned Fix | Status |
|----|------|------|---------|-------------|--------|
| 1 | `gitmap/cmd/comp_053_test.go` | 26 | `if !ok` inside `if err != nil` | Guard clause / flat assertion | Pending |
| 2 | `gitmap/cmd/comp_102_test.go` | 24 | `if appErr.Code !=` inside `if err != nil` | Guard clause / flat assertion | Pending |
| 3 | `gitmap/cmd/comp_109_test.go` | 25 | `if appErr.Code !=` inside `if err != nil` | Guard clause / flat assertion | Pending |
| 4 | `gitmap/cmd/fix_cmd.go` | 20 | `if len(args) == 0` inside `if isHelp` | Invert condition & return early | Pending |
| 5 | `gitmap/cmd/ip_cmd_test.go` | 62 | `if len(outStr) == 0` inside `if !hasIP` | Flatten conditional | Pending |
| 6 | `gitmap/cmd/rootusage.go` | 260 | `if curLen+1+wLen > width` inside `if len(cur) > 0` | Decompose / early return | Pending |
| 7 | `gitmap/cmd/ssh_client_test.go` | 48 | `if _, ok := err.(...); !ok` inside `if err != nil` | Flat error type assertion | Pending |
| 8 | `gitmap/scanner/progress_test.go` | 143 | `if p.IsFinal` inside loop conditional | Invert and continue early | Pending |
| 9 | `gitmap/store/ssh_repo_test.go` | 74 | `if appErr.Code !=` inside `if err != nil` | Flat error assertion | Pending |
| 10 | `gitmap/store/ssh_repo_test.go` | 117 | `if appErr.Code !=` inside `if err != nil` | Flat error assertion | Pending |
| 11 | `gitmap/store/ssh_repo_test.go` | 120 | `if !errors.Is(appErr.Cause...)` inside `if err != nil` | Flat error assertion | Pending |
| 12 | `.lovable/ai-fix-scripts/01-file-manipulator.py` | 96 | `if '=' in pair:` inside loop | Flatten with helper / continue | Pending |
| 13 | `.lovable/ai-fix-scripts/03-cicd-local-runner.py` | 193 | `if out.strip():` inside failed_jobs loop | Extract helper function | Pending |
| 14 | `.lovable/ai-fix-scripts/newline_fixer.py` | 47 | `if fix_trailing_newline(filepath):` inside check | Flatten with guard | Pending |
| 15 | `scripts/audit_codebase.py` | 20, 40, 46, 48, 79 | Nested conditions in audit loop | Guard clauses / early continue | Pending |

## Subtasks

- [ ] Subtask 01: Flatten Go test files (`comp_053_test.go`, `comp_102_test.go`, `comp_109_test.go`, `ip_cmd_test.go`, `ssh_client_test.go`, `progress_test.go`, `ssh_repo_test.go`).
- [ ] Subtask 02: Flatten Go command implementations (`gitmap/cmd/fix_cmd.go`, `gitmap/cmd/rootusage.go`).
- [ ] Subtask 03: Flatten Python scripts (`01-file-manipulator.py`, `03-cicd-local-runner.py`, `newline_fixer.py`, `scripts/audit_codebase.py`).
- [ ] Subtask 04: Connect `linter-scripts/check-nested-ifs.py` to `.github/workflows/ci.yml` and `03-cicd-local-runner.py`.
- [ ] Subtask 05: Verify 100% green CI test suite and complete plan lifecycle.

## Acceptance Criteria

- [ ] `python linter-scripts/check-nested-ifs.py` reports ✅ PASS with 0 violations.
- [ ] `python .lovable/ai-fix-scripts/03-cicd-local-runner.py` passes all quality gates with exit code 0.
