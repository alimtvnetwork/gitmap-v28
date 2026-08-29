# CI/CD Issue 01 — golangci-lint misspell: "labeled"

## Pipeline
- **Tool:** `golangci-lint v1.64.8` (`misspell` linter)
- **Command:** `golangci-lint run --path-prefix=gitmap --timeout=5m`
- **Runner:** GitHub Actions (`gitmap-v28` repo)

## Symptom
```
gitmap/cmd/scanbenchmark.go:27:25: `labeled` is a misspelling of `labeled` (misspell)
// benchPhase holds one labeled timing measurement.
```

## Root Cause
British-English spelling `labeled` used in a Go doc comment. `misspell` enforces US spelling.

## Fix
- Replaced `labeled` → `labeled` in `gitmap/cmd/scanbenchmark.go` (line 27).
- Repo-wide grep found a second occurrence in `gitmap/scripts/install.sh` ("one labeled line") — also fixed.
- Verified `aria-labelledby` in `Troubleshooting.tsx` is a standard ARIA attribute and must NOT be touched.

## Verification
- `grep -rn "labeled" gitmap/` (excluding `aria-labelledby`) → 0 matches.
- `grep -rni "\blabeled\b\|\bcanceled\b\|\bbehavior\b\|\bcolor\b\|\boccurred\b\|\breceive\b\|\bseparate\b" --include="*.go" gitmap/` → 0 matches.

## Status
✅ Resolved (session 2026-04-23)

## Prevention
- Prefer US spelling in all Go comments, identifiers, help text, and shell scripts.
- When adding/reviewing comments, treat `misspell`'s US dictionary as the source of truth.
- Common offenders to avoid: `labeled`, `canceled`, `behavior`, `color`, `occurred`, `receive`, `separate`, `traveling`, `modeling`.
