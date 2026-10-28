# Subtask 05: E2E Verification and Release Ceremony
Parent Plan: 338-gitmap-pull-errors-and-ui-enhancement

## Scope
- Run `go test` across `cmdfix`, `cmdpull`, `cmdmail`, `cmdlogviewer`.
- Run CI script tests (`python .github/scripts/tests/test_ci_scripts.py`).
- Run `python .github/scripts/go-format-check.py --no-commit`.
- Execute minor bump to `v6.624.0` via `37-bump-version.py`.
- Rebuild binary and release via `gitmap cpr`.
