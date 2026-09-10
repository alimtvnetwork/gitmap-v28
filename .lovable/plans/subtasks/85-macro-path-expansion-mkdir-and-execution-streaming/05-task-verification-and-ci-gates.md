# Subtask 05: Verification, Quality Gates & Consolidation

## 1. Objective
Run comprehensive automated tests, enforce CI/CD linters (nested ifs, boolean guidelines, relative paths, error management), verify all interactive and CLI commands end-to-end, compile the production binary into canonical target locations (`bin/gitmap.exe` and `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`), and consolidate Plan 85 to completed.

## 2. Target Files
- `.lovable/plans/completed/85-macro-path-expansion-mkdir-and-execution-streaming.md`
- `.lovable/plans/01-index.md`
- `bin/gitmap.exe`

## 3. Requirements
- Automated test suites:
  - Run `go test -v ./macro/...`
  - Run `go test -v ./cmd/...`
  - Ensure 100% pass rate.
- Linter verification:
  - `python 03-ai-scripts/check-nested-ifs.py` (0 violations)
  - `python 03-ai-scripts/check-boolean-guidelines.py` (0 violations)
  - `python 03-ai-scripts/check-relative-paths.py` (0 violations)
  - `python 03-ai-scripts/check-error-management.py` (0 violations)
- Functional verification:
  - Build binary: `go build -o bin/gitmap.exe .`
  - Deploy binary to `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`
  - Test `gitmap mkdir -p test-deep//nested/folder`
  - Test `gitmap mkdir -f test-deep/nested/file.txt`
  - Test `gitmap macro add` environment expansion logic
- Plan Consolidation:
  - Move `.lovable/plans/pending/85-macro-path-expansion-mkdir-and-execution-streaming.md` to `.lovable/plans/completed/85-macro-path-expansion-mkdir-and-execution-streaming.md`.
  - Update `.lovable/plans/01-index.md`.
