# 10 — CI/CD Pipeline: Known Issues, Root Causes & Fixes

> **Purpose**: This is a self-contained post-mortem catalog of every CI/CD failure encountered in this project. Any AI model or engineer reading this should be able to diagnose and fix the same class of issue instantly without external context.
>
> **Format for each entry**: Symptom → Root Cause → Why It Wasn't Caught → Fix → Prevention Rule → Related Files.

---

## Issue Index

| # | Issue | Pipeline Stage | Severity | Status |
|---|-------|---------------|----------|--------|
| 1 | `npm ci` lockfile drift | docs-site build | 🔴 Blocker | ✅ Fixed v2.82.0 |
| 2 | `go-winres` icon size > 256x256 | release: resource embed | 🔴 Blocker | ✅ Fixed v2.81.0 |
| 3 | `cd: dist: No such file or directory` | release: compress | 🔴 Blocker | ✅ Fixed v2.54.0 |
| 4 | Job-level `if` blocks required status checks | CI: SHA dedup | 🟠 High | ✅ Fixed |
| 5 | `cancel-in-progress` cancels mark-success job | CI: SHA dedup | 🟠 High | ✅ Fixed |
| 6 | `@latest` tool installs are non-reproducible | CI: setup | 🟡 Medium | ✅ Fixed |
| 7 | Release branch runs cancelled mid-pipeline | release | 🔴 Blocker | ✅ Fixed |

---

## Issue #1 — `npm ci` Lockfile Drift

### Symptom

```
npm error code EUSAGE
npm error `npm ci` can only install packages when your package.json and
npm error package-lock.json or npm-shrinkwrap.json are in sync.
npm error Missing: vitest@3.2.4 from lock file
npm error Missing: @testing-library/react@16.3.2 from lock file
... (100+ missing entries)
Error: Process completed with exit code 1.
```

### Root Cause

Dependencies were added to `package.json` (via `code--add_dependency` or manual edit) **without** regenerating `package-lock.json`. CI uses `npm ci` (clean install) which **refuses to run** if the two files are out of sync — by design, to guarantee reproducible builds.

`npm install` would silently update the lockfile on a developer machine, but CI must use `npm ci` for determinism.

### Why It Wasn't Caught

- Local development typically uses `npm install` or `bun install`, both of which auto-update the lockfile.
- The `package-lock.json` is in the read-only file list (the AI cannot modify it via normal edits) — it must be regenerated via the shell.
- No pre-commit hook validates lockfile sync.

### Fix (v2.82.0)

Regenerate the lockfile from scratch:

```bash
cd /dev-server
rm -f package-lock.json
npm install --package-lock-only --ignore-scripts
# Verify:
npm ci --ignore-scripts --dry-run
```

### Prevention Rules

1. **After ANY change to `package.json` dependencies, regenerate `package-lock.json` in the same commit.**
2. **Add an early CI step** that runs `npm ci --dry-run` before any heavy work, so lockfile drift fails fast (seconds) instead of after the build (minutes).
3. **Never hand-edit `package-lock.json`** — always regenerate.
4. When using AI tooling, after `code--add_dependency` runs, immediately run `npm install --package-lock-only --ignore-scripts` if the AI cannot do so itself.

### Related Files

| File | Purpose |
|------|---------|
| `package.json` | Source of truth for dependency versions |
| `package-lock.json` | Locked transitive dependency tree (read-only to AI) |
| `bun.lock` | Parallel lockfile for Bun runtime |
| `.github/workflows/*.yml` | Workflows that invoke `npm ci` |

---

## Issue #2 — `go-winres` Icon Size > 256x256

### Symptom

```
go: downloading github.com/cpuguy83/go-md2man/v2 v2.0.2
go: downloading github.com/russross/blackfriday/v2 v2.1.0
2026/04/16 16:26:46 image size too big, must fit in 256x256
Error: Process completed with exit code 1.
```

### Root Cause

`go-winres` embeds icons as Windows `.ico` resources inside the compiled `.exe`. The Windows `.ico` format has a **hard limit of 256x256 pixels** per image frame. The project's `gitmap-v27/assets/icon.png` was **512x512**, which exceeds this limit.

The error originates from the `go-winres make` step in the release workflow (`.github/workflows/release.yml`), which reads `gitmap-v27/winres/winres.json` and converts the referenced PNG into an `.ico` resource.

### Why It Wasn't Caught

- Local Windows builds via `run.ps1` do **not** invoke `go-winres` — icon embedding only happens in CI release pipeline.
- Plain `go build` succeeds without `go-winres`; it just produces a binary without embedded Windows metadata (icon, manifest, version info).
- The `.png` source file passed all visual review — there is no automatic dimension validation.

### Fix (v2.81.0)

1. Created a 256x256 resized copy: `gitmap-v27/assets/icon-256.png`.
2. Updated `gitmap-v27/winres/winres.json` to reference `icon-256.png` instead of `icon.png`.
3. Kept the original 512x512 `icon.png` for web/docs use.

### Prevention Rules

1. **Any icon referenced in `winres.json` MUST be ≤ 256x256 pixels.**
2. **Maintain separate icon files** by purpose: `icon.png` (512x512+ for web/docs), `icon-256.png` (256x256 for `.exe` embedding).
3. **Add a pre-check step** in CI before `go-winres make`:
   ```bash
   python3 -c "from PIL import Image; img=Image.open('gitmap-v27/assets/icon-256.png'); assert max(img.size)<=256, f'Icon too large: {img.size}'"
   ```
4. **Document the constraint** in `gitmap-v27/winres/README.md` so future contributors don't replace the icon with a higher-resolution version.

### Related Files

| File | Purpose |
|------|---------|
| `gitmap-v27/winres/winres.json` | Windows resource manifest for `go-winres` |
| `gitmap-v27/assets/icon-256.png` | 256x256 icon for `.exe` embedding |
| `gitmap-v27/assets/icon.png` | 512x512 original icon (web/docs) |
| `.github/workflows/release.yml` | CI pipeline that runs `go-winres make` |
| `spec/08-generic-update/09-winres-icon-constraint.md` | Original post-mortem |

---

## Issue #3 — `cd: dist: No such file or directory`

### Symptom

```
cd: dist: No such file or directory
Error: Process completed with exit code 1.
```

Failed in the compress/checksum step of `release.yml`.

### Root Cause

The compress step ran inside `gitmap-updater/` (which has no `dist/` folder) instead of `gitmap-v27/dist/` where cross-compiled binaries are output. The script used `cd dist` which assumed the working directory was `gitmap-v27/`, but GitHub Actions defaults to the repository root for every `run:` step unless `working-directory` is explicitly set.

In a monorepo with multiple Go modules (`gitmap-v27/`, `gitmap-updater/`), the previous step's working directory does **not** carry over to the next step.

### Why It Wasn't Caught

- Locally, `run.ps1` always executes from `gitmap-v27/`, so `cd dist` works.
- Only the CI environment exhibits the root-relative behavior.
- No `test -d dist` guard was present to fail fast with a useful message.

### Fix (v2.54.0)

Replaced `cd dist` with an explicit `working-directory` directive:

```yaml
- name: Compress and checksum
  working-directory: gitmap-v27/dist
  run: |
    for f in gitmap-*; do
      ...
    done
```

### Prevention Rules

1. **NEVER use `cd` inside CI scripts.** Always use the YAML `working-directory:` field on the step.
2. **Validate output directories exist** before operating: `test -d "$DIR" || { echo "::error::$DIR missing"; exit 1; }`.
3. **In monorepos, use absolute or explicitly anchored paths** — never assume the prior step's CWD.
4. **Test pipeline changes on a `release/test-*` branch** before merging to `main`.

### Related Files

| File | Purpose |
|------|---------|
| `.github/workflows/release.yml` | Compress/checksum step |
| `spec/02-app-issues/13-release-pipeline-dist-directory.md` | Original post-mortem |

---

## Issue #4 — Job-Level `if` Blocks Required Status Checks

### Symptom

When SHA-deduplication was implemented with job-level `if: needs.sha-check.outputs.already-built != 'true'`, the GitHub UI showed grey "skipped" icons for all jobs. Branch protection rules treated skipped jobs as **neither success nor failure**, blocking PR merges that required passing checks.

### Root Cause

GitHub Actions distinguishes three job conclusion states: `success`, `failure`, and `skipped`. Required status checks configured in branch protection only accept `success`. A skipped job never resolves the check, so the PR is permanently blocked.

### Fix

Adopted the **passthrough gate pattern**:
- Jobs always run (no job-level `if`).
- Each step inside the job uses `if: needs.sha-check.outputs.already-built != 'true'` to skip the actual work.
- A first step always echoes "✅ Already validated (SHA cached)" so the job always concludes as `success` with a green checkmark.

### Prevention Rules

1. **Never use job-level `if` for cache/dedup gating** if any required status check depends on that job.
2. **Use step-level conditionals** for skip logic.
3. **Always include at least one unconditional step** (even just an `echo`) so the job concludes successfully.

### Related Files

| File | Purpose |
|------|---------|
| `.github/workflows/ci.yml` | `sha-check` gate + downstream jobs |
| `spec/05-coding-guidelines/29-ci-sha-deduplication.md` | Pattern documentation |

---

## Issue #5 — `cancel-in-progress` Cancels Mark-Success Job

### Symptom

A separate `mark-success` job was added at the end of the CI pipeline to write the `ci-passed-<SHA>` cache entry. Intermittently, this job was cancelled before completing, leaving the SHA uncached. Re-runs of the same SHA then re-executed the entire pipeline.

### Root Cause

The workflow used `concurrency.cancel-in-progress: true`. When all validation jobs completed and `mark-success` was queued, a new push to the same ref would cancel the **entire workflow run** — including the still-pending `mark-success` job. The validation work had succeeded, but the cache save was lost.

### Fix

**Inlined the cache write as the final step of the `test-summary` job** instead of using a separate job. Because `test-summary` is itself a validation job that downstream success depends on, it cannot be cancelled after success without also affecting the visible status.

### Prevention Rules

1. **Side-effects that must persist after pipeline success (cache writes, telemetry, deployment markers) belong in the LAST validation job, not a separate downstream job.**
2. **Avoid trailing "marker" jobs** when `cancel-in-progress: true` is set.
3. **Failed pipelines must NEVER write success markers** — guard cache writes with `if: success()`.

### Related Files

| File | Purpose |
|------|---------|
| `.github/workflows/ci.yml` | `test-summary` final cache write step |
| `spec/05-coding-guidelines/29-ci-sha-deduplication.md` | Pattern documentation |

---

## Issue #6 — `@latest` Tool Installs Are Non-Reproducible

### Symptom

CI builds that passed yesterday started failing today with no code changes, due to a breaking change in a Go tool installed via `go install foo@latest`.

### Root Cause

`@latest` resolves to whatever version is current at install time. A new tool release between two CI runs introduced an incompatibility (new lint rule, removed flag, behavior change).

### Fix

Pinned every tool to an exact version tag:

| Tool | Pinned Version |
|------|---------------|
| `golangci-lint` | `v1.64.8` |
| `govulncheck` | `v1.1.4` |
| `actions/checkout` | `@v6` |
| `actions/setup-go` | `@v6` |
| `actions/cache` | `@v4` |
| `softprops/action-gh-release` | `@v2` |

### Prevention Rules

1. **`@latest` and `@main` are PROHIBITED** in any CI workflow or `setup.sh`.
2. **Pin every action and CLI tool to an exact tag** (`@v1.2.3`, not `@v1`).
3. **Document version bumps** in `CHANGELOG.md` so regressions can be bisected.
4. **Use Dependabot or Renovate** to propose pinned-version bumps via PR.

### Related Files

| File | Purpose |
|------|---------|
| `setup.sh` | Local dev tool install (must match CI versions) |
| `.github/workflows/*.yml` | All workflows |
| `spec/05-coding-guidelines/17-cicd-patterns.md` | CI/CD patterns spec |

---

## Issue #7 — Release Branch Runs Cancelled Mid-Pipeline

### Symptom

A push to `release/v2.5x.0` triggered the release workflow, but a follow-up commit (e.g., changelog typo fix) cancelled the in-progress release run, leaving artifacts half-built and the GitHub Release in an inconsistent state.

### Root Cause

The workflow used `cancel-in-progress: true` unconditionally. This is correct behavior for feature branches and PRs (avoid wasted compute) but **incorrect for release branches** where every commit must produce complete artifacts and metadata.

### Fix

Use a **conditional `cancel-in-progress`** that protects release branches:

```yaml
concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: ${{ !startsWith(github.ref, 'refs/heads/release/') }}
```

The release workflow itself uses `cancel-in-progress: false` unconditionally.

### Prevention Rules

1. **Release branch runs MUST run to completion.** Never cancel them.
2. **Use conditional `cancel-in-progress`** based on `github.ref` prefix.
3. **Release workflow should never have `cancel-in-progress: true`.**

### Related Files

| File | Purpose |
|------|---------|
| `.github/workflows/ci.yml` | Conditional `cancel-in-progress` |
| `.github/workflows/release.yml` | `cancel-in-progress: false` |

---

## Universal CI/CD Prevention Checklist

Apply to **every** CI/CD change before merging:

- [ ] No `cd` commands in `run:` blocks — use `working-directory:`
- [ ] No `@latest` or `@main` tool/action references — pin to exact tags
- [ ] No job-level `if` for cache/dedup gating that affects required status checks
- [ ] No trailing "marker" jobs for cache writes — inline into last validation job
- [ ] No `cancel-in-progress: true` for release branches
- [ ] Lockfile (`package-lock.json`, `go.sum`) regenerated alongside any dependency change
- [ ] All file paths used in CI either absolute, anchored to `working-directory`, or guarded by `test -d` / `test -f`
- [ ] All Windows resource icons ≤ 256x256 pixels
- [ ] Pipeline tested on a `release/test-*` or feature branch before merging to `main`
- [ ] Failure modes produce actionable error messages (use `::error::` annotation)

---

## Cross-References

- [01-ci-pipeline.md](./01-ci-pipeline.md) — CI pipeline architecture
- [02-release-pipeline.md](./02-release-pipeline.md) — Release pipeline architecture
- [09-binary-icon-branding.md](./09-binary-icon-branding.md) — Windows icon embedding
- [spec/07-generic-release/07-known-issues-and-fixes.md](../07-generic-release/07-known-issues-and-fixes.md) — Release-specific issue catalog
- [spec/05-coding-guidelines/17-cicd-patterns.md](../05-coding-guidelines/17-cicd-patterns.md) — CI/CD patterns
- [spec/05-coding-guidelines/29-ci-sha-deduplication.md](../05-coding-guidelines/29-ci-sha-deduplication.md) — SHA dedup pattern
- [spec/02-app-issues/13-release-pipeline-dist-directory.md](../02-app-issues/13-release-pipeline-dist-directory.md) — `dist` directory post-mortem
- [spec/08-generic-update/09-winres-icon-constraint.md](../08-generic-update/09-winres-icon-constraint.md) — winres icon post-mortem
- [spec/12-consolidated-guidelines/16-cicd.md](../12-consolidated-guidelines/16-cicd.md) — Consolidated CI/CD guidelines
