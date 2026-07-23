# CI/CD Issues Index

Tracks every CI/CD pipeline failure or hardening decision encountered, its root cause, and resolution. New entries go in `.lovable/cicd-issues/XX-short-name.md` with sequential numeric prefixes.

## Conventions
- File naming: `XX-kebab-case-name.md` (XX = zero-padded sequence starting at `01`).
- One file per distinct issue. Do **not** duplicate — if the same root cause recurs, append a "Recurrence" section to the existing file.
- Status values: `✅ Resolved`, `🔄 In Progress`, `⏳ Pending`, `🚫 Blocked`.

## Issues

| # | Title | Tool / Stage | Status | File |
|---|-------|--------------|--------|------|
| 01 | misspell: `labelled` → `labeled` | golangci-lint (misspell) | ✅ Resolved | [01-misspell-labelled.md](cicd-issues/01-misspell-labelled.md) |
| 02 | `lint-regression-guard` → `lint-hard-floor` → `lint-baseline-guard` (now uniformly baseline-diff for all 5 linters) | golangci-lint (baseline-guard job) | ✅ Resolved | [02-lint-regression-guard-semantics.md](cicd-issues/02-lint-regression-guard-semantics.md) |

## Patterns Learned
- **US-English everywhere in Go**: `misspell` flags British spellings in comments and identifiers. Avoid `labelled`, `cancelled`, `behaviour`, `colour`, `occured`, `recieve`, `seperate`.
- **Pinned linter versions**: golangci-lint is pinned to `v1.64.8`; `goimports` to `v0.24.0`; `govulncheck` to `v1.1.4`. Do not assume newer rules.
- **ARIA attributes are exempt**: `aria-labelledby` is a standard HTML/ARIA token and must never be "corrected".
- **Decide hard-floor vs baseline-diff per linter** before adding to CI; document choice in the script header. Don't conflate them under one job.
- **Compute `-local` prefix from `go.mod`** — never hardcode the module path in CI scripts.
- **Cache `GOMODCACHE`/`GOCACHE`** keyed on `go.sum` hash for compile-gate + matrix speed.
