# CI Pipeline Architecture

The CI pipeline (GitHub Actions) uses a parallel matrix strategy ('fail-fast: false') to execute four distinct test suites: unit, store, integration, and tui. Test output and coverage profiles ('-covermode=atomic') are collected as artifacts and consolidated by a final 'test-summary' job. This summary job aggregates failures into a single report, calculates project-wide coverage using 'go tool cover', and generates a per-package breakdown. To ensure visibility, the test stage uses 'set +e' and 'grep' to filter for specific Go failure patterns (e.g., '--- FAIL', 'build failed', 'undefined') before exiting with the original code.

## Failure Report Script

The `.github/scripts/test-summary.sh` script parses each suite's `test-output.txt` to extract failing test names and their specific failure reasons (assertion errors, expected/got mismatches, panics, undefined references). It produces a **"FAILURE REPORT (copy-paste ready)"** block at the end — a self-contained summary that can be shared directly without scrolling through full logs. The script uses `awk` to capture lines between `=== RUN` and `--- FAIL` markers, filtering for `.go:<line>:` patterns and error keywords.

## Binary Builds on Main

After all tests pass, a `build` matrix job cross-compiles 6 binaries (windows/linux/darwin × amd64/arm64) versioned as `dev-<sha>` using `CGO_ENABLED=0` and uploads them as artifacts with 14-day retention. A subsequent `build-summary` job downloads all artifacts and prints a formatted table of binary names and human-readable file sizes.

## SHA-Based Build Deduplication (Passthrough Gate Pattern)

A 'sha-check' gate job runs before all other jobs. It probes the GitHub Actions cache for key 'ci-passed-<SHA>' using 'lookup-only: true'. Downstream jobs always run (no job-level `if` skipping) but use **step-level conditionals**: when the SHA is already cached, each job executes only an "Already validated" echo step and exits with ✅ Success. This ensures the GitHub UI always shows green checkmarks — never grey "skipped" icons that look like failures and block required status checks. When the cache misses, steps guarded by `if: needs.sha-check.outputs.already-built != 'true'` execute normally. The cache write is **inlined as the final step of `test-summary`** (not a separate `mark-success` job) to prevent `cancel-in-progress` from cancelling the cache save while all validation jobs have already passed. Failed pipelines never cache, so re-runs of the same SHA execute fully. Documented in spec/05-coding-guidelines/29-ci-sha-deduplication.md.

## Concurrency Control

All workflows use 'concurrency: group: ci-${{ github.ref }}' to scope runs by branch. For non-release branches, 'cancel-in-progress: true' cancels superseded runs. Release branches ('release/**') are **never cancelled** — they always run to completion to ensure every release commit produces complete artifacts and metadata. The CI workflow uses a conditional expression: 'cancel-in-progress: ${{ !startsWith(github.ref, 'refs/heads/release/') }}'. The release workflow uses 'cancel-in-progress: false' unconditionally.

## Lessons Learned

1. **Never use `cd` in CI scripts** — use `working-directory` in the workflow step definition. The v2.54.0 release pipeline failed with `cd: dist: No such file or directory` because the compress step ran in `gitmap-updater/` instead of `gitmap/`. Fixed by setting explicit `working-directory: gitmap/dist`. See `spec/02-app-issues/13-release-pipeline-dist-directory.md`.
2. **Pin Go tool versions** — `go install tool@latest` is non-reproducible. All tools (e.g., `golangci-lint@v1.64.8`, `govulncheck@v1.1.4`) must use exact version tags. Documented in `setup.sh` and `spec/05-coding-guidelines/17-cicd-patterns.md`.
3. **Validate build output directories** before operating on them: `test -d "$DIR" || exit 1`.
4. **Never use job-level `if` for SHA deduplication** — GitHub treats skipped jobs as neither success nor failure, blocking required status checks. Use the passthrough gate pattern with step-level conditionals instead.
5. **Inline cache writes into the last validation job** — a separate `mark-success` job can be cancelled by `cancel-in-progress` after all validation passes, leaving the SHA uncached. Inlining the cache save as the final step of `test-summary` prevents this.
6. **Compile gate before matrix** (2026-04-29) — `go test ./...` typecheck step runs before the per-suite matrix to fail fast on build errors. Backed by `actions/setup-go@v5` cache + explicit `actions/cache` for `GOMODCACHE`/`GOCACHE` keyed on `go.sum`.
7. **Format/imports/lint hard gate** (2026-04-29) — Lint job runs (in order): `gofmt -l` → `goimports -l/-d` (pinned `v0.24.0`, `-local` derived from `go.mod`) → `go vet` → `golangci-lint v1.64.8 --issues-exit-code=1`. All four are hard failures.
8. **Unified baseline-diff guard** (2026-04-29, final state) — `lint-baseline-guard` job (history: `lint-regression-guard` → `lint-hard-floor` → `lint-baseline-guard`) runs five analyzers all under the SAME contract: only NEW findings vs the cached main-branch baseline fail the build. Linters: `unused`, `gosec G115` (via `TEXT_FILTER=G115`), `misspell`, `gocritic`, `exhaustive` — each with its own rolling per-linter cache slot keyed by SHA. All driven by `.github/scripts/check-single-linter-diff.sh` (now supports `TEXT_FILTER` regex + `LABEL` override). The hard-floor script `check-lint-regressions.sh` was deleted. Full JSON diff still lives in the separate `lint-baseline-diff` job (`lint-diff.py`). When adding a linter, copy an existing sub-step block (restore → diff → save → persist). See `.lovable/cicd-issues/02-lint-regression-guard-semantics.md`.
