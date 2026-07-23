# Missing macOS Binaries and Lint Regression

## Ticket

v2.64.0 release shipped without macOS (darwin) gitmap-v27 binaries and with
5 unresolved golangci-lint errors that block CI.

## Symptoms

1. GitHub Release v2.64.0 shows only 8 assets instead of the expected ~14.
   Missing: `gitmap-v27.64.0-darwin-amd64.tar.gz`,
   `gitmap-v27.64.0-darwin-arm64.tar.gz`, and several other platform
   binaries. The `gitmap-updater-v2.64.0-darwin-amd64.tar.gz` IS present.
2. CI fails with 5 lint errors:
   - `paramTypeCombine` on `autocommit.go` (gocritic)
   - 3× `S1039` unnecessary `fmt.Sprintf` on `installtools.go` (gosimple)
   - `misspell`: "cancelled" → "canceled" on `constants_install.go`

## Root Cause

### Lint Errors

Introduced during the v2.65.0 install UX refactor. The `writeInstallErrorLog`
function was added with `fmt.Sprintf` wrapping bare string literals (no
format verbs), and the `AutoCommit` signature was not updated to use Go's
grouped-parameter syntax. The "cancelled" misspelling existed from initial
constant creation and was never caught because misspell was not enabled in
earlier linter configs.

### Missing Binaries

The release workflow (`release.yml`) correctly defines all 6 targets
(windows/linux/darwin × amd64/arm64). The most likely cause is that lint
failures in the CI pipeline (`ci.yml`) block the green status required
for release branch protection, causing the release job to run on a
partially-built or stale commit. Since the release workflow is a single
job, a build failure in any target would prevent ALL subsequent steps
(compress, checksum, upload) from completing.

## Fix

| File | Change |
|------|--------|
| `release/autocommit.go` | `func AutoCommit(version string, dryRun, yes bool)` — group same-type params |
| `cmd/installtools.go` | Replace `fmt.Sprintf("literal\n")` with `"literal\n"` (3 sites) |
| `constants/constants_install.go` | `cancelled` → `canceled` |

## Prevention

1. **Pre-commit lint gate**: Before any Go PR is merged, the CI must pass
   `golangci-lint run` with zero errors. This is already configured but
   the errors were introduced after the gate was last green.
2. **Refactor guardrail** (added to `.lovable/plan.md` and
   `spec/04-generic-cli/13-checklist.md`): After any Go file split or
   refactor, run `go test` and `go vet` on the affected package
   immediately — do not defer lint cleanup.
3. **Grouped parameters**: Always use Go's `func(a, b bool)` syntax when
   consecutive parameters share a type. The `gocritic/paramTypeCombine`
   linter enforces this.
4. **No `fmt.Sprintf` for bare strings**: Use `sb.WriteString("literal")`
   directly. The `gosimple/S1039` linter catches this.
5. **American English spelling**: Use `canceled`, `color`, `behavior` in
   all string constants. The `misspell` linter enforces US English.

## Verification

After merging the fix, re-run the release pipeline for the next version
tag and confirm:
- CI lint step passes green
- All 6 gitmap-v27 + 6 updater binaries appear in the release assets
- macOS binaries are downloadable and executable

## Related

- `.golangci.yml` — linter configuration
- `spec/02-app-issues/14-security-hardening-gosec-fixes.md` — prior lint fixes
- `.lovable/plan.md` — Go refactor validation guardrail