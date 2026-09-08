# Subtask 06: Quality Gate Verification and CI/CD

**Parent Plan:** [77-scripts-fixer-installation-split-db-and-tooling-engine.md](../../pending/77-scripts-fixer-installation-split-db-and-tooling-engine.md)
**Status:** Pending
**Target Files:**
- Whole repository validation

## Objectives
1. Ensure all new functions <= 15 lines, blank lines before return statements, Unix LF line endings, and proper error wrapping.
2. Run Go unit tests: `go test -v ./gitmap/store/... ./gitmap/cmd/...`.
3. Run `python linter-scripts/check-nested-ifs.py`.
4. Run `python linter-scripts/check-enum-and-boolean.py`.
5. Run `python .github/scripts/go-format-check.py`.
6. Run full CI runner: `python 03-ai-scripts/06-cicd-local-runner.py`.
