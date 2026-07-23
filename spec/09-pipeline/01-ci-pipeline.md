# CI Pipeline

## Overview

The CI pipeline validates every push and pull request to the `main` branch. It runs linting, vulnerability scanning, parallel test suites, and cross-compiled builds — then caches the result so identical commits are never re-validated.

---

## Trigger and Concurrency

### Trigger

```yaml
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  workflow_dispatch:
    # Manual escape hatches for the golangci-lint baseline cache. See the
    # "Job: Lint Baseline Diff" section below for full semantics.
    inputs:
      lint_baseline_cache_version:
        description: "Bump to invalidate the golangci-lint baseline cache (e.g. v2)"
        required: false
        default: "v1"
        type: string
      lint_baseline_disable:
        description: "Skip restoring/saving the lint baseline cache for this run"
        required: false
        default: false
        type: boolean
```

The two `workflow_dispatch` inputs let an operator manually intervene when the cached lint baseline becomes stale (linter version bump, configuration change, or a partial earlier run that seeded the cache with bogus findings) **without** having to delete the cache through the GitHub UI. Default behavior on `push` / `pull_request` is unchanged: `lint_baseline_cache_version` falls back to `"v1"` and `lint_baseline_disable` to `false`.

### Concurrency Control

Scope concurrent runs by branch reference. Cancel superseded runs on feature branches, but **never cancel release branches** (they must always run to completion):

```yaml
concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: ${{ !startsWith(github.ref, 'refs/heads/release/') }}
```

**Why**: If two pushes land on the same branch in quick succession, the first run is cancelled to save resources. Release branches are exempt because every release commit must produce complete artifacts.

---

## Job Graph

```
sha-check ──┬── lint ──────┬── test (matrix: 4 suites) ──┬── test-summary ── build (matrix: 6 targets) ── build-summary
             │              │                              │
             └── vulncheck ─┘                              └── (cache SHA on success)
```

All jobs depend on `sha-check`. If the SHA is already cached, every job prints a skip message and exits successfully — no work is repeated.

---

## Pattern: SHA-Based Build Deduplication (Passthrough Gate)

### Problem

Re-running the pipeline on an already-validated commit wastes compute and blocks PRs with slow re-checks.

### Solution

A gate job probes the GitHub Actions cache for a key tied to the commit SHA. Downstream jobs always run (never `if: false` at job level) but use **step-level conditionals** to skip expensive work when the SHA is cached.

### Why Not Job-Level `if`?

GitHub treats skipped jobs as neither success nor failure. Required status checks see "skipped" as not passing — this blocks merges and shows grey icons in the UI. The passthrough pattern ensures every job shows a green checkmark.

### Implementation

**Step 1 — Gate job** outputs whether the SHA was found:

```yaml
sha-check:
  runs-on: ubuntu-latest
  outputs:
    already-built: ${{ steps.cache-check.outputs.cache-hit }}
  steps:
    - name: Check SHA cache
      id: cache-check
      uses: actions/cache@v4
      with:
        path: /tmp/ci-passed
        key: ci-passed-${{ github.sha }}
        lookup-only: true
```

**Step 2 — Every downstream step** uses a conditional:

```yaml
- name: Already validated
  if: needs.sha-check.outputs.already-built == 'true'
  run: echo "✅ SHA ${{ github.sha }} already passed"

- uses: actions/checkout@v6
  if: needs.sha-check.outputs.already-built != 'true'

# ... all remaining steps guard with the same condition
```

**Step 3 — Cache write is inlined** into the last validation job (not a separate job):

```yaml
- name: Mark SHA as built
  if: success() && needs.sha-check.outputs.already-built != 'true'
  run: mkdir -p /tmp/ci-passed && echo "${{ github.sha }}" > /tmp/ci-passed/sha.txt

- name: Save SHA to cache
  if: success() && needs.sha-check.outputs.already-built != 'true'
  uses: actions/cache/save@v4
  with:
    path: /tmp/ci-passed
    key: ci-passed-${{ github.sha }}
```

**Why inline the cache write?** A separate `mark-success` job could be cancelled by `cancel-in-progress` after all validation jobs pass but before the cache is saved. Inlining it into the final validation step prevents this race condition.

---

## Job: Lint

Runs static analysis using `go vet` and `golangci-lint`.

### Steps

1. Checkout repository
2. Setup Go toolchain (version from `go.mod`)
3. Run `go vet ./...`
4. Run `golangci-lint` (pinned version, 5-minute timeout)

### Tool Configuration

- `golangci-lint` is configured via `.golangci.yml` in the project root
- Version is pinned (e.g., `v1.64.8`) — never use `@latest`
- The `golangci-lint-action` handles caching automatically
- Use `working-directory` to point at the Go module root

```yaml
- uses: golangci/golangci-lint-action@v6
  with:
    version: v1.64.8
    working-directory: <module-root>
    args: --timeout=5m
```

---

## Job: Lint Baseline Diff

Soft-gate companion to **Job: Lint**. Runs `golangci-lint` with `--issues-exit-code=0 --out-format=json` so the raw report is collected even when findings exist, then diffs the current findings against a baseline cached from the last successful `main` run. **Fails only on NEW findings.** Pre-existing findings are surfaced in the PR sticky comment as `UNCHANGED` but never block the build — contributors can land incremental improvements without being forced to fix every historical lint issue at once.

### Cache strategy

| Aspect | Behavior |
|---|---|
| Cache key | `golangci-baseline-main-${cache_version}-${github.sha}` |
| Restore key (fallback) | `golangci-baseline-main-${cache_version}-` (rolling, single slot) |
| Save trigger | Only on `push` to `main` (or `workflow_dispatch` from `main`) after a green diff |
| PR behavior | Restore-only — PRs read the baseline but never advance it |
| Miss behavior | **Seeding mode** — diff script exits 0 and emits all current findings as warnings so the next run has a baseline to compare against |

### `workflow_dispatch` cache controls

Both inputs are optional; defaults preserve normal CI behavior.

#### `lint_baseline_cache_version` *(string, default `"v1"`)*

Bumps the cache key suffix. Old caches are abandoned (GitHub evicts them after 7 days of inactivity); the next `main`-branch run reseeds under the new version. Free-form string — common values are `"v2"`, `"v3"`, or a date stamp like `"2026-04-21"`. The `restore-keys` fallback also carries this version, so a pre-bump baseline is **never** accidentally restored.

**When to bump:** golangci-lint version upgrade, `.golangci.yml` rule changes that broaden coverage, or recovery from a poisoned cache that contains findings the diff script can't normalize.

#### `lint_baseline_disable` *(boolean, default `false`)*

When `true`, **skips both the restore and save steps for this run.** The diff job enters seeding mode (exits 0, surfaces all current findings as warnings) without touching the cache. Use this to diagnose suspected stale-cache issues without losing baseline history — the next normal `main` run will continue using the existing cached baseline.

### Example usage

**Bump the cache version (recommended after a `golangci-lint` upgrade):**

```bash
gh workflow run ci.yml \
  --ref main \
  -f lint_baseline_cache_version=v2
```

**One-off run that ignores the baseline entirely (diagnostic):**

```bash
gh workflow run ci.yml \
  --ref main \
  -f lint_baseline_disable=true
```

**Combine both — bump the version AND skip restore/save for the dispatch run** (useful when you suspect the most recent `main` run poisoned the cache and you want to inspect raw findings before committing to a reseed):

```bash
gh workflow run ci.yml \
  --ref main \
  -f lint_baseline_cache_version=v2 \
  -f lint_baseline_disable=true
```

After verifying the dispatch output looks correct, push a no-op commit to `main` (or re-dispatch with `lint_baseline_disable=false`) to seed the new `v2` cache.

### Sticky PR comment

After every diff, `lint-suggest.py` writes an actionable suggestions block to `/tmp/lint-suggestions/comment.md`. On `pull_request` events, `peter-evans/find-comment` locates the previous comment via the `<!-- gitmap-lint-suggestions -->` sentinel and `peter-evans/create-or-update-comment` edits it in place — preventing the PR from accumulating stale lint comments on every push. On `push` / `workflow_dispatch` events the comment is mirrored to `GITHUB_STEP_SUMMARY` instead.

---

## Job: Vulnerability Scan (In-CI)

Runs `govulncheck` to detect known vulnerabilities in dependencies.

### Stdlib vs. Third-Party Vulnerability Handling

The scanner differentiates between:
- **Third-party vulnerabilities** (packages you import) → **fail the build**
- **Standard library vulnerabilities** (unfixable until Go upgrades) → **warn only**

### Implementation

```bash
set +e
govulncheck ./... 2>&1 | tee /tmp/vulncheck.out
rc=$?
if [ $rc -ne 0 ]; then
  if grep "packages you import" /tmp/vulncheck.out | grep -qv "0 vulnerabilities in packages"; then
    echo "::error::Third-party vulnerabilities detected"
    exit 1
  fi
  echo "::warning::Only stdlib vulnerabilities found (no fix available in current Go version)"
fi
```

**Why**: Go's standard library vulnerabilities cannot be patched without upgrading the Go toolchain itself. Failing the build on these would block all development until a new Go release is available.

---

## Job: Test (Matrix)

Runs tests in parallel across multiple suites using a matrix strategy.

### Matrix Definition

```yaml
strategy:
  fail-fast: false
  matrix:
    include:
      - name: unit
        packages: ./pkg1/... ./pkg2/...
      - name: store
        packages: ./tests/store_test/...
      - name: integration
        packages: ./tests/integration_test/...
      - name: tui
        packages: ./tui/...
```

**Why `fail-fast: false`?** All suites run to completion even if one fails. This gives a complete picture of what's broken, not just the first failure.

### Test Execution

```bash
set -uo pipefail
set +e
go test ${{ matrix.packages }} -v -count=1 \
  -coverprofile=coverage-${{ matrix.name }}.out \
  -covermode=atomic 2>&1 | tee test-output.txt
exit_code=${PIPESTATUS[0]}
set -e

if [ "$exit_code" -ne 0 ]; then
  echo "========================================="
  echo "  FAILED TESTS (${{ matrix.name }})"
  echo "========================================="
  grep -E -A 5 '^--- FAIL:' test-output.txt || true
  echo "========================================="
fi

exit "$exit_code"
```

Key details:
- `set +e` prevents the step from exiting on the first test failure — capture the exit code instead
- `tee` writes output to both stdout and a file for later artifact upload
- `PIPESTATUS[0]` captures the exit code of `go test`, not `tee`
- `-count=1` disables test caching for reliable CI results
- `-covermode=atomic` enables safe concurrent coverage collection

### Artifact Upload

Upload test output and coverage profiles for aggregation:

```yaml
- uses: actions/upload-artifact@v4
  if: always()  # upload even on failure
  with:
    name: test-results-${{ matrix.name }}
    path: |
      test-output.txt
      coverage-${{ matrix.name }}.out
    retention-days: 7
```

---

## Job: Test Summary

Aggregates results from all test matrix jobs into a single report.

### Failure Report Script

The script (`.github/scripts/test-summary.sh`) parses each suite's `test-output.txt` to extract:
- Failing test names
- Specific failure reasons (assertion errors, expected/got mismatches, panics, undefined references)

It produces a **"FAILURE REPORT (copy-paste ready)"** block — a self-contained summary that can be shared directly without scrolling through full CI logs.

```bash
#!/usr/bin/env bash
set -uo pipefail

RESULTS_DIR="${1:?Usage: test-summary.sh <results-dir>}"
overall=0
all_failures=""

for dir in "$RESULTS_DIR"/test-results-*; do
  suite=$(basename "$dir" | sed 's/test-results-//')
  file="$dir/test-output.txt"
  [ ! -f "$file" ] && continue

  pass=$(grep -c '^--- PASS:' "$file" || true)
  fail=$(grep -c '^--- FAIL:' "$file" || true)

  if [ "$fail" -gt 0 ]; then
    overall=1
    # Extract failure details using awk between "=== RUN" and "--- FAIL" markers
    # Filter for .go:<line>: patterns and error keywords
    # Append to all_failures
  else
    echo "✅ $suite: $pass passed"
  fi
done

if [ "$overall" -ne 0 ]; then
  echo "========================================="
  echo "  FAILURE REPORT (copy-paste ready)"
  echo "========================================="
  echo "$all_failures"
  echo "========================================="
  exit 1
fi

echo "All test suites passed."
```

### Coverage Aggregation

After failure reporting, merge all coverage profiles and produce a per-package breakdown:

1. Concatenate all `coverage-*.out` files (stripping duplicate `mode:` headers)
2. Run `go tool cover -func=combined-coverage.out`
3. Use `awk` to calculate per-package average coverage
4. Print a formatted table with package paths and percentages

### Cache Write

The SHA cache is saved as the **final step** of `test-summary` (only on success). See the SHA deduplication pattern above for why this is inlined.

---

## Job: Build (Matrix)

Cross-compiles binaries for all target platforms after tests pass.

### Matrix Definition

```yaml
strategy:
  fail-fast: false
  matrix:
    include:
      - { os: windows, arch: amd64, ext: .exe }
      - { os: windows, arch: arm64, ext: .exe }
      - { os: linux,   arch: amd64, ext: ""   }
      - { os: linux,   arch: arm64, ext: ""   }
      - { os: darwin,  arch: amd64, ext: ""   }
      - { os: darwin,  arch: arm64, ext: ""   }
```

### Build Command

```bash
VERSION="dev-${GITHUB_SHA::10}"
LDFLAGS="-s -w -X '<module>/constants.Version=$VERSION'"
OUTPUT="<binary>-${VERSION}-${{ matrix.os }}-${{ matrix.arch }}${{ matrix.ext }}"
CGO_ENABLED=0 GOOS=${{ matrix.os }} GOARCH=${{ matrix.arch }} go build -ldflags "$LDFLAGS" -o "$OUTPUT" .
```

Key details:
- `CGO_ENABLED=0` produces static binaries with no C dependencies
- `-s -w` strips debug symbols for smaller binaries
- `-X` embeds the version string at compile time
- CI builds use `dev-<sha>` versioning; release builds use semantic versions

### Artifact Upload

Each binary is uploaded with 14-day retention:

```yaml
- uses: actions/upload-artifact@v4
  with:
    name: <binary>-${{ matrix.os }}-${{ matrix.arch }}
    path: <binary>-*
    retention-days: 14
```

---

## Job: Build Summary

Downloads all build artifacts and prints a formatted table of binary names and file sizes.

```bash
printf "  %-45s %s\n" "Binary" "Size"
for dir in binaries/*; do
  for file in "$dir"/*; do
    name=$(basename "$file")
    size=$(stat --format="%s" "$file")
    human=$(numfmt --to=iec --suffix=B "$size")
    printf "  %-45s %s\n" "$name" "$human"
  done
done
```

---

## Constraints

- No `cd` in CI steps — use `working-directory` in the step definition
- All tool installs use exact version tags — `@latest` is prohibited
- Validate build output directories before operating on them: `test -d "$DIR" || exit 1`
- Never use job-level `if` for SHA deduplication — use step-level conditionals
- Inline cache writes into the last validation job — never a separate job
- No notification steps (email, Slack, etc.) in workflows
- All scripts use `set -euo pipefail` (or `set -uo pipefail` when capturing exit codes)
