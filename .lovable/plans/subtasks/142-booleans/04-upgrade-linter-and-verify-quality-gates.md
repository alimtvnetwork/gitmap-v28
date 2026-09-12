# Subtask 04: Upgrade Linter, Verify Quality Gates & Consolidate Plan 142

Parent Plan: [142-boolean-principles-negatives-and-complex-conditions.md](../../pending/142-boolean-principles-negatives-and-complex-conditions.md)

## Goals
1. Upgrade `linter-scripts/check-boolean-guidelines.py` to check banned prefixes (`can`, `should`, `was`, `will`, `did`, `must`) and ensure negative naming flags are detected.
2. Run `go vet ./...` in `cli/` and verify 0 errors.
3. Run `python linter-scripts/check-boolean-guidelines.py`.
4. Run `python linter-scripts/check-nested-ifs.py`.
5. Run `python linter-scripts/check-relative-paths.py`.
6. Run `python linter-scripts/check-enum-guidelines.py`.
7. Run `python 03-ai-scripts/35-result-wrapper-auditor.py`.
8. Record modified files under lock via `python 03-ai-scripts/33-test-inventory-generator.py --record`.
9. Consolidate Plan 142 into `.lovable/plans/completed/142-boolean-principles-negatives-and-complex-conditions.md`.
10. Update `.lovable/plans/01-index.md`.
11. Atomic git commit & push via SSH.

## Acceptance Criteria
- [x] All linters and quality gates pass with exit code 0.
- [x] Plan 142 consolidated and committed atomically.
