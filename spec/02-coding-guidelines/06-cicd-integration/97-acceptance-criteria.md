# Acceptance Criteria — CI/CD Integration

> **Version:** 1.0.0
> **Updated:** 2026-04-19

---

## AC-CI-001 — Portability

Each check script in `linters-cicd/checks/` runs on a stock Ubuntu
runner with only `python3` (≥ 3.10) and `bash` available. **No** `pip
install` required for Phase 1.

## AC-CI-002 — SARIF compliance

`linters-cicd/scripts/validate-sarif.py` validates every emitted file
against the official SARIF 2.1.0 schema. CI run is green.

## AC-CI-003 — Self-test on this repo

Running `./linters-cicd/run-all.sh --path .` against this repository
produces a SARIF file with **zero** CODE RED findings. (We dogfood our
own rules.)

## AC-CI-004 — Composite Action one-liner

A consumer can add coding-guidelines linting to their GitHub workflow
with exactly one `uses:` line and no other config.

## AC-CI-005 — Versioned release artifact

Every `v*` tag attaches `coding-guidelines-linters-vX.Y.Z.zip` to the
GitHub Release with a SHA-256 entry in `checksums.txt`.

## AC-CI-006 — Plugin model unchanged when adding language

Adding a new language plugin requires zero edits to `run-all.sh`,
`action.yml`, or any check script for another language. Verified by
PR-template checklist.

## AC-CI-007 — Exit codes

Every check exits `0` on clean, `1` on findings, `2` on tool error.
Verified by `linters-cicd/checks/_tests/test_exit_codes.sh`.

---

*Part of [CI/CD Integration](./00-overview.md)*
