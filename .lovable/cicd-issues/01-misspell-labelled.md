# CI/CD Issue 01 — golangci-lint misspell: "labelled"

## Pipeline
- **Tool:** `golangci-lint v1.64.8` (`misspell` linter)
- **Command:** `golangci-lint run --path-prefix=gitmap --timeout=5m`
- **Runner:** GitHub Actions (`gitmap-v28` repo)

## Symptom
```
gitmap/cmd/scanbenchmark.go:27:25: `labelled` is a misspelling of `labeled` (misspell)
// benchPhase holds one labelled timing measurement.
```

## Root Cause
British-English spelling `labelled` used in a Go doc comment. `misspell` enforces US spelling.

## Fix
- Replaced `labelled` → `labeled` in `gitmap/cmd/scanbenchmark.go` (line 27).
- Repo-wide grep found a second occurrence in `gitmap/scripts/install.sh` ("one labelled line") — also fixed.
- Verified `aria-labelledby` in `Troubleshooting.tsx` is a standard ARIA attribute and must NOT be touched.

## Verification
- `grep -rn "labelled" gitmap/` (excluding `aria-labelledby`) → 0 matches.
- `grep -rni "\blabelled\b\|\bcancelled\b\|\bbehaviour\b\|\bcolour\b\|\boccured\b\|\brecieve\b\|\bseperate\b" --include="*.go" gitmap/` → 0 matches.

## Status
✅ Resolved (session 2026-04-23)

## Prevention
- Prefer US spelling in all Go comments, identifiers, help text, and shell scripts.
- When adding/reviewing comments, treat `misspell`'s US dictionary as the source of truth.
- Common offenders to avoid: `labelled`, `cancelled`, `behaviour`, `colour`, `occured`, `recieve`, `seperate`, `travelling`, `modelling`.
