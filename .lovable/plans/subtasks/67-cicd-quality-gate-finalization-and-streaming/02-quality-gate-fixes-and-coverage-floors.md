# Subtask 02: Quality Gate Fixes, Test Hardening & Coverage Floor Protection

## Objective
Ensure all failing gates discovered during local CI/CD execution are fixed and covered:
1. `gitmap/logging`: Ensure unit tests in `jsonlog_test.go` achieve 100% statement coverage.
2. `gitmap/visibility`: Ensure tests across `pattern_test.go`, `exclude_test.go`, and `fuzzy_test.go` achieve >75% coverage.
3. `.github/scripts/coverage-floor.py`: Ensure `go tool cover` executes within `gitmap` module directory, resolves relative paths cleanly, and verifies all packages registered in `.github/coverage.floor`.
4. `06-cicd-local-runner.py`: Remove `-coverpkg=./...` from `Go Test Coverage Profile` to avoid profile corruption and memory bloat.

## Files
- `gitmap/logging/jsonlog_test.go`
- `gitmap/visibility/pattern_test.go`
- `gitmap/visibility/exclude_test.go`
- `gitmap/visibility/fuzzy_test.go`
- `.github/scripts/coverage-floor.py`
- `03-ai-scripts/06-cicd-local-runner.py`
- `.lovable/ai-fix-scripts/06-cicd-local-runner.py`

## Acceptance Criteria
- [ ] `gitmap/logging` tests pass and coverage is >= 80%.
- [ ] `gitmap/visibility` tests pass and coverage is >= 75%.
- [ ] `python .github/scripts/coverage-floor.py coverage.out` passes with exit code 0.
