# CI Hardening Session — 2026-04-29

> Session focused on tightening the GitHub Actions `lint` + compile pipeline. All changes scoped to `.github/workflows/ci.yml` (no Go source changes).

## ✅ Done

1. **Compile gate** — Added a `go test ./...` step after code changes; any typecheck/build error now fails CI before the matrix tests run.
2. **Build/test caching** — Enabled `GOMODCACHE` + `GOCACHE` reuse via `actions/setup-go@v5` cache + explicit `actions/cache` keys for `go.sum` hash. Speeds up the compile gate and downstream jobs.
3. **gofmt check** — Added a step that runs `gofmt -l` on all `.go` files; fails (with file list) on any unformatted file.
4. **goimports check** — Added a `goimports` step pinned to `golang.org/x/tools/cmd/goimports@v0.24.0`. Reads `LOCAL_PREFIX` dynamically from `go.mod`, runs `-l` for detection then `-d` for diff output, prints copy-pasteable fix command. Positioned between `gofmt` and `go vet`.
5. **golangci-lint strict gate** — Updated `golangci/golangci-lint-action@v6` step args to include `--issues-exit-code=1`. Pinned version remains `v1.64.8`, timeout `5m`. Working dir `gitmap`.
6. **Lint guard semantics — final state (2026-04-29, two-pass)** — All five guarded linters now share ONE contract: baseline-diff via `.github/scripts/check-single-linter-diff.sh`.
   - **Pass 1:** Renamed `lint-regression-guard` → `lint-hard-floor` to make the split semantics honest.
   - **Pass 2:** User chose to unify everything as baseline-diff. Renamed → `lint-baseline-guard`. Replaced the `unused`+`G115` hard-floor step with two new baseline-diff sub-steps (each with rolling per-linter cache: `golangci-unused-baseline-main-` and `golangci-gosec-g115-baseline-main-`). Extended `check-single-linter-diff.sh` with `TEXT_FILTER` (regex on `.Text`, used to scope `gosec` to G115 — applied to BOTH current and baseline) and `LABEL` (banner/annotation override). Deleted `.github/scripts/check-lint-regressions.sh`. Added `lint-baseline-guard` to `test-summary` job's `needs:`.
   - Linters covered: `unused`, `gosec G115`, `misspell`, `gocritic`, `exhaustive`. Each has its own rolling cache slot keyed by SHA, restored via prefix.
   - Full JSON diff still in separate `lint-baseline-diff` job (`lint-diff.py`).

## ⏳ Pending / Open

- ⚠️ **Branch protection follow-up**: required-check name in branch protection (if configured) needs updating to `Lint Baseline Guard (unused, gosec G115, misspell, gocritic, exhaustive)`. Old check names (`Lint Regression Guard …`, `Lint Hard Floor …`) will no longer appear after this merges.
- ⏳ Verify CI green on next push (no live run inspected this session).
- ⏳ Consider adding `goimports` and `gofmt` to the local `hooks/pre-commit` for parity with CI.

## Key Snippets

### goimports check
```yaml
- name: goimports check (import grouping + formatting)
  run: |
    go install golang.org/x/tools/cmd/goimports@v0.24.0
    GOIMPORTS="$(go env GOPATH)/bin/goimports"
    LOCAL_PREFIX="$(awk '/^module /{print $2; exit}' go.mod)"
    unformatted=$("$GOIMPORTS" -l -local "$LOCAL_PREFIX" .)
    if [ -n "$unformatted" ]; then
      echo "::error::The following .go files have goimports issues..."
      "$GOIMPORTS" -d -local "$LOCAL_PREFIX" $unformatted
      exit 1
    fi
```

### golangci-lint strict gate
```yaml
- name: golangci-lint (strict, fail on any error)
  if: needs.sha-check.outputs.already-built != 'true'
  uses: golangci/golangci-lint-action@v6
  with:
    version: v1.64.8
    working-directory: gitmap
    args: --timeout=5m --issues-exit-code=1
```

## Files Touched

- `.github/workflows/ci.yml` (all six steps above)

## Lessons / Anti-patterns

- **Pin every tool**: `goimports@v0.24.0`, `golangci-lint@v1.64.8`. Never `@latest` in CI.
- **Compute `-local` from `go.mod`**: avoids hardcoding the module path in CI.
- **Don't conflate floor vs diff**: Final resolution 2026-04-29 — job renamed `lint-regression-guard` → `lint-hard-floor` → `lint-baseline-guard`, all guarded linters now use the same baseline-diff contract via `check-single-linter-diff.sh`. Hard-floor mechanism (`check-lint-regressions.sh`) deleted entirely. If a true hard-floor is ever needed again, give it its own dedicated job — never mix models inside one.

## Next AI Pickup Point

CI guard work is complete and uniform. Branch protection (if used) needs the required-check name updated to `Lint Baseline Guard (unused, gosec G115, misspell, gocritic, exhaustive)`. Otherwise next logical step is wiring `gofmt`/`goimports` into `hooks/pre-commit` for local parity, or moving on to other CI improvements (test coverage gates, release pipeline hardening, etc.).
