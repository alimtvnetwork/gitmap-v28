# CI/CD Pipeline Architecture

## Overview

The project uses three GitHub Actions workflows to enforce code quality,
produce release artifacts, and scan for vulnerabilities. All workflows
implement concurrency controls to cancel superseded runs.

---

## Workflows

### 1. CI (`ci.yml`)

**Triggers:** Push to `main`, pull requests targeting `main`.

| Stage          | Purpose                                           |
|----------------|---------------------------------------------------|
| Lint           | `go vet` + `golangci-lint` (28 linters, 5m timeout) |
| Vulnerability  | `govulncheck` — fails only on third-party issues  |
| Test           | Parallel matrix: unit, store, integration, tui    |
| Test Summary   | Aggregates failures, generates coverage breakdown |
| Build          | Cross-compile 6 targets after tests pass          |
| Build Summary  | Lists all binaries with sizes                      |

**Binary builds on main.** After all tests pass, the CI pipeline
cross-compiles 6 binaries (windows/linux/darwin × amd64/arm64)
versioned as `dev-<sha>` and uploads them as artifacts (14-day
retention). This provides pre-built binaries for every green commit
without requiring a formal release.

### 2. Release (`release.yml`)

**Triggers:** Push to `release/**` branches or `v*` tags.

Produces 6 cross-compiled binaries (windows/linux/darwin ×
amd64/arm64) for both `gitmap-v28` and `gitmap-updater`. Generates
versioned artifacts, SHA256 checksums, and changelog excerpts.

#### Code Signing

Windows `.exe` binaries are optionally signed via [SignPath.io](https://signpath.io)
after compilation but before compression/checksumming.

| Variable / Secret | Type | Purpose |
|-------------------|------|---------|
| `SIGNPATH_SIGNING_ENABLED` | Repository **variable** | Set to `true` to enable signing; omit or set to any other value to skip |
| `SIGNPATH_API_TOKEN` | Repository **secret** | API token from SignPath dashboard |
| `SIGNPATH_ORGANIZATION_ID` | Repository **secret** | Organization ID from SignPath |
| `SIGNPATH_PROJECT_SLUG` | Repository **secret** | Project slug configured in SignPath |
| `SIGNPATH_SIGNING_POLICY_SLUG` | Repository **secret** | Signing policy slug (e.g., `release-signing`) |

See [05-code-signing.md](05-code-signing.md) for full setup instructions.

### 3. Vulnerability Scan (`vulncheck.yml`)

**Triggers:** Weekly schedule (Mondays 09:00 UTC), manual dispatch.

Runs `govulncheck` independently of the CI pipeline for proactive
dependency monitoring.

---

## SHA-Based Build Deduplication

Before any work begins, a `sha-check` gate job probes the GitHub
Actions cache for a key `ci-passed-<SHA>`. If the commit has already
passed CI, all downstream jobs print "Already validated" and exit
with ✅ Success (passthrough gate pattern). On full success, the
`test-summary` job writes a marker to the cache as its final step,
so future runs for the same SHA short-circuit immediately.

The cache write is **inlined into `test-summary`** (not a separate
job) to prevent `cancel-in-progress` from cancelling the marker
write while validation jobs have already passed.

This eliminates redundant builds when the same commit is pushed
multiple times (e.g., re-tagging, merge commits, or manual re-runs
after transient infrastructure failures).

---

## Concurrency Control

All three workflows use GitHub Actions `concurrency` groups to
manage parallel runs. When a new commit is pushed while a previous
run is still executing, the behavior depends on the branch type.

### Group Keys

| Workflow   | Concurrency Group             | Cancel In-Progress |
|------------|-------------------------------|--------------------|
| CI         | `ci-${{ github.ref }}`        | Yes (non-release)  |
| Release    | `release-${{ github.ref }}`   | No (never)         |
| Vulncheck  | `vulncheck-${{ github.ref }}` | Yes                |

### Release Branch Protection

Release branches (`release/**`) are **never cancelled**, even when
multiple commits are pushed in quick succession. This ensures that
every release commit runs the full CI and release pipeline to
completion — partial builds or missed artifacts are unacceptable
for release branches.

For the CI workflow, `cancel-in-progress` uses a conditional
expression:

```yaml
cancel-in-progress: ${{ !startsWith(github.ref, 'refs/heads/release/') }}
```

This evaluates to `true` (cancel) for `main` and feature branches,
but `false` (never cancel) for release branches. The release
workflow uses `cancel-in-progress: false` unconditionally since it
only triggers on release branches and tags.

### Branch Behavior

The `github.ref` suffix ensures that:

- Different branches run independently (a push to `feature/a` does
  not cancel a run on `main`).
- Multiple pushes to the **same** non-release branch cancel each
  other, keeping only the latest.
- Pull request runs are scoped to the PR ref, so updating a PR
  cancels its previous CI run without affecting `main`.
- **Release branches always run to completion** — no cancellation.

### Why cancel-in-progress (for non-release branches)?

| Problem                                  | Solution                          |
|------------------------------------------|-----------------------------------|
| Wasted CI minutes on outdated commits    | Auto-cancel superseded runs       |
| Queue buildup during rapid iteration     | Only latest commit matters        |
| Stale results reported on merged PRs     | Cancelled runs produce no output  |

### Why NOT cancel release branches?

| Problem                                       | Solution                               |
|------------------------------------------------|----------------------------------------|
| Partial release artifacts from cancelled builds| Never cancel — run to completion       |
| Missed binaries or checksums                   | Every push produces complete artifacts |
| Incomplete metadata writes                     | Full pipeline guarantees consistency   |

---

## Test Architecture

Tests run in a parallel matrix (`fail-fast: false`) across four
suites. Each suite produces:

- `test-output.txt` — full verbose output for failure analysis.
- `coverage-<suite>.out` — atomic coverage profile.

The `test-summary` job downloads all artifacts and delegates
failure analysis to `.github/scripts/test-summary.sh`. This script:

1. Iterates each suite's `test-output.txt`.
2. Counts pass/fail per suite.
3. For each failing test, extracts the test name and the actual
   failure reason (assertion errors, expected/got mismatches,
   panics, undefined references) from the verbose output.
4. Produces a **"FAILURE REPORT (copy-paste ready)"** block at the
   end — a self-contained summary that can be shared directly
   without scrolling through full logs.

Coverage profiles are merged separately via `go tool cover`.

---

## Build Strategy

| Context         | Binary Production | Rationale                        |
|-----------------|-------------------|----------------------------------|
| `main` branch   | 6 targets         | Dev binaries for every green SHA |
| Pull requests   | 6 targets         | Same as main                     |
| `release/**`    | 6 targets         | Official artifacts for release   |
| `v*` tags       | 6 targets         | Tagged release artifacts         |

CI builds produce `dev-<sha>` artifacts for testing and validation.
Release builds produce official versioned artifacts attached to
GitHub Releases with SHA256 checksums.

---

## File Layout

| File                              | Purpose                              |
|-----------------------------------|--------------------------------------|
| `.github/workflows/ci.yml`        | Lint, test, build, coverage on main  |
| `.github/workflows/release.yml`   | Cross-compile and publish releases   |
| `.github/workflows/vulncheck.yml` | Weekly vulnerability scan            |
| `.github/scripts/test-summary.sh` | Failure report aggregation script    |
| `.golangci.yml`                   | Linter configuration (28 rules)      |

## References

For detailed, portable pipeline implementation guides (suitable for sharing with any AI or engineer), see:

- [spec/pipeline/README.md](../09-pipeline/README.md) — Index and quick reference
- [spec/pipeline/01-ci-pipeline.md](../09-pipeline/01-ci-pipeline.md) — CI patterns (SHA dedup, test matrix, builds)
- [spec/pipeline/02-release-pipeline.md](../09-pipeline/02-release-pipeline.md) — Release automation (versioning, install scripts)
- [spec/pipeline/03-vulnerability-scanning.md](../09-pipeline/03-vulnerability-scanning.md) — Vulnerability scanning

---

## Acceptance Criteria

1. Pushing two commits in quick succession to the same branch
   cancels the first CI run and only completes the second.
2. Pushing to `main` runs lint, vulncheck, tests, and builds
   6 cross-compiled binaries uploaded as artifacts.
3. Pushing to `release/**` or a `v*` tag produces 6 binaries
   and uploads them as GitHub Release assets.
4. The weekly vulncheck runs independently and does not block
   or interfere with CI or release workflows.
5. Pull request CI runs are scoped to the PR ref and do not
   cancel runs on `main` or other PRs.
6. Test failures in one suite do not prevent other suites from
    completing (`fail-fast: false`).
7. When any test suite fails, the test summary produces a
   **"FAILURE REPORT (copy-paste ready)"** block listing each
   failing test name and its specific failure reason (assertion
   errors, expected/got mismatches, panics, undefined references).
8. The build summary job prints a formatted table of all 6
   binaries with their human-readable file sizes.