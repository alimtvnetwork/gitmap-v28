# Subtask 05: CI Linter Integration & End-to-End Verification

> **Parent Plan:** `17-boolean-and-naming-audit.md`  
> **Scope:** `linter-scripts/check-enum-and-boolean.py`, `.lovable/ai-fix-scripts/03-cicd-local-runner.py`, `.github/workflows/ci.yml`

## Objective

Connect `check-enum-and-boolean.py` to the local CI runner and GitHub Actions CI workflow, and verify 100% green pass.

## Action Steps

1. Register `check-enum-and-boolean.py` in `03-cicd-local-runner.py`.
2. Verify GitHub Actions workflow step in `.github/workflows/ci.yml`.
3. Run `python linter-scripts/check-enum-and-boolean.py` to confirm zero violations.
4. Run `python .lovable/ai-fix-scripts/03-cicd-local-runner.py` and `go test ./...`.
