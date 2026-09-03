# 17 — Release Pipeline Issues: Root Cause Analysis & Prevention

**Version:** 2.0.0
**Created:** 2026-04-16
**Updated:** 2026-04-16
**Status:** Active reference
**Audience:** Any AI model or engineer maintaining the CI/CD pipeline
**Goal:** Document every release/CI failure encountered (this repo + sibling reference implementation lessons), with root cause and durable fix, so the same mistakes never recur.

---

## Issue Index

| # | Issue | Stage | Severity | Source |
|---|-------|-------|----------|--------|
| 1 | `npm ci` lockfile drift | release | 🔴 Blocker | This repo |
| 2 | `setup-python` pip cache requires manifest | CI | 🟠 High | This repo |
| 3 | `release.sh` depended on Node at runtime | release | 🟠 High | This repo |
| 4 | `go-winres` icon > 256×256 px | release: resource embed | 🔴 Blocker | sibling-ref |
| 5 | `cd: dist: No such file or directory` | release: compress | 🔴 Blocker | sibling-ref |
| 6 | Job-level `if` blocks required status checks | CI: SHA dedup | 🟠 High | sibling-ref |
| 7 | `cancel-in-progress` cancels marker job | CI: cache write | 🟠 High | sibling-ref |
| 8 | `@latest` tool/action installs non-reproducible | CI: setup | 🟡 Medium | sibling-ref |
| 9 | Release branch run cancelled by follow-up commit | release: concurrency | 🔴 Blocker | sibling-ref |
| 10 | Install-script placeholder unreplaced | release: script gen | 🟠 High | sibling-ref |
| 11 | Missing `GITHUB_TOKEN` silently skips upload | release: asset upload | 🟡 Medium | sibling-ref |
| 12 | Asset name mismatch between checksum and upload | release: packaging | 🟠 High | sibling-ref |

---

## Purpose

This document is a living **post-mortem ledger** for the project's GitHub Actions CI and Release pipelines. Each entry captures:

1. **Symptom** — exact error message as seen in CI logs
2. **Trigger** — what action/commit/state caused it
3. **Root cause** — the *why*, not just the *what*
4. **Fix applied** — the concrete change committed
5. **Prevention rule** — the durable guideline future contributors and AI models must follow

Read this file **before** modifying `.github/workflows/*.yml`, `release.sh`, `install.sh`, `install.ps1`, or `package.json` version-related fields.

---

## Pipeline Architecture (1-line summary)

- **CI** (`.github/workflows/ci.yml`): runs on push/PR to `main`. Validates linter scripts (Go + Python) and Axios pinning. **No Node/npm involvement.**
- **Release** (`.github/workflows/release.yml`): runs on `v*` tags. Executes `release.sh` to build artifacts, then publishes a GitHub Release. **No Node/npm involvement.**

> ⚠️ **CRITICAL INVARIANT:** Neither workflow may invoke `npm ci`, `npm install`, `node`, or any JS toolchain. The repository ships a `package.json` for *local* preview tooling only — the lockfile is **not** kept in sync with `package.json` and must never be relied on in CI.

---

## Issue Ledger

### Issue #1 — `npm ci` fails: lockfile out of sync

**Date observed:** 2026-04-16
**Workflow:** Release (`release.yml`)
**Tag:** `vX.Y.Z`

**Symptom:**

```
npm error `npm ci` can only install packages when your package.json and
package-lock.json or npm-shrinkwrap.json are in sync.
npm error Missing: @monaco-editor/react@4.7.0 from lock file
npm error Missing: monaco-editor@0.55.1 from lock file
... (40+ missing packages)
Error: Process completed with exit code 1.
```

**Trigger:**
The release workflow ran `npm ci` to bootstrap a Node environment so it could call `node -p "require('./package.json').version"` to read the version string.

**Root cause:**

1. `package.json` is updated by Lovable's preview environment whenever new deps are added.
2. `package-lock.json` is **read-only** in this project and is *not* regenerated on every dep change.
3. Therefore `package.json` and `package-lock.json` drift out of sync continuously.
4. `npm ci` enforces strict sync — it fails immediately.
5. **Deeper cause:** the release pipeline had no business depending on Node at all. Reading a version string does not require `npm ci`.

**Fix applied:**

- Removed `actions/setup-node` and `npm ci` from `release.yml`.
- Replaced `node -p` with a shell `resolve_version` function in `release.sh` that:
  1. Prefers `$RELEASE_VERSION` env var (set by workflow from `GITHUB_REF_NAME`).
  2. Falls back to `sed`-based extraction from `package.json`.
- Workflow now passes `RELEASE_VERSION: ${{ steps.version.outputs.version }}` into `release.sh`.

**Prevention rule:**
🔴 **NEVER add `npm ci`, `npm install`, or `actions/setup-node` to any CI/CD workflow in this repo.** The lockfile is not maintained. If a workflow needs a value from `package.json`, extract it with `sed`/`grep`/`jq` (jq is preinstalled on `ubuntu-latest`). If a workflow genuinely needs Node tooling, it must first regenerate the lockfile with `npm install --package-lock-only` and accept the resulting drift — but this is **discouraged**.

---

### Issue #2 — `setup-python` cache fails: no `requirements.txt`

**Date observed:** 2026-04-16
**Workflow:** CI (`ci.yml`)

**Symptom:**

```
Error: No file in /home/runner/work/coding-guidelines-v24/coding-guidelines-v24
matched to [**/requirements.txt or **/pyproject.toml], make sure you have
checked out the target repository
```

**Trigger:**
`actions/setup-python@v5` was configured with `cache: 'pip'`. The built-in pip cache requires a dependency manifest (`requirements.txt` or `pyproject.toml`) to compute its cache key. The repo has neither — the Python validator (`linter-scripts/validate-guidelines.py`) uses only the standard library.

**Root cause:**
Built-in caching in `setup-python` and `setup-go` assumes the project has a canonical dependency manifest at a discoverable path. This repo is *not* a Python or Go project — it merely *contains* validator scripts in those languages with zero external deps. Built-in caching is the wrong tool.

**Fix applied:**

- Removed `cache: 'pip'` from `setup-python` and `cache: true` from `setup-go`.
- Added explicit `actions/cache@v4` steps with cache keys derived from `hashFiles('linter-scripts/**/*.go')` and `hashFiles('linter-scripts/**/*.py')`.
- Cache paths: `~/.cache/go-build`, `~/go/pkg/mod` for Go; `~/.cache/pip` for pip.

**Prevention rule:**
🔴 **Do not enable built-in `cache:` on `setup-python` or `setup-go` in this repo.** Always use explicit `actions/cache@v4` with `hashFiles()` keyed on the script files themselves. If a future contributor adds a real `requirements.txt` or `go.mod`, only then may built-in caching be reconsidered.

---

### Issue #3 — Release script depended on Node at runtime

**Date observed:** 2026-04-16 (same incident as Issue #1)
**Workflow:** Release (`release.sh`)

**Symptom:**
`release.sh` invoked `node -p "require('./package.json').version"`, which forced the workflow to install Node and run `npm ci`, cascading into Issue #1.

**Root cause:**
The release script was written assuming a JS-tooled environment. A version-bump workflow does not need a JS runtime to read a JSON field.

**Fix applied:**
`resolve_version()` shell function in `release.sh`:
```bash
resolve_version() {
  if [ -n "${RELEASE_VERSION:-}" ]; then
    echo "${RELEASE_VERSION#v}"
    return
  fi
  sed -n 's/.*"version": *"\([^"]*\)".*/\1/p' package.json | head -1
}
```

**Prevention rule:**
🔴 **`release.sh` and any future CI shell scripts must remain language-agnostic.** Never invoke `node`, `npm`, `python`, or `go` from `release.sh` unless the script's *purpose* is to build that language's artifacts. Use POSIX shell utilities (`sed`, `awk`, `grep`, `jq`, `cut`) for parsing.

---

## Issues Imported from Sibling Reference Implementation

The following lessons were captured during the development of the a sibling reference implementation (its `02-spec/12-cicd-pipeline-workflows/10-known-issues-and-fixes.md` and `02-spec/16-generic-release/08-known-issues-and-fixes.md`). They describe failure modes that **will eventually occur in this repo too** if/when we add: a Go binary release, Windows resource embedding, SHA-deduplication, install-script generation, or multi-asset uploads. They are documented preemptively so the same diagnosis cycle does not repeat.

> ⚠️ Issues #4–#12 below describe pipeline patterns this repo does not yet ship but is **likely to adopt**. Treat them as canonical guardrails for any future workflow extension.

---

### Issue #4 — `go-winres` Icon Size > 256×256 px

**Symptom:**

```
2026/04/16 16:26:46 image size too big, must fit in 256x256
Error: Process completed with exit code 1.
```

**Trigger:** A release workflow step running `go-winres make` references a PNG larger than 256×256 from `winres.json`.

**Root cause:** The Windows `.ico` resource format hard-limits each frame to 256×256 px. `go-winres` refuses any larger source image. Local `go build` succeeds without resource embedding, so the constraint is invisible until CI runs.

**Fix applied (in sibling reference implementation):**

- Created a 256×256 copy `assets/icon-256.png`.
- Updated `winres/winres.json` to reference the smaller file.
- Kept the original 512×512 `icon.png` for web/docs.

**Prevention rule:**
🔴 **Any PNG referenced from `winres.json` MUST be ≤ 256×256.** Maintain two separate files (`icon.png` for docs, `icon-256.png` for `.exe` embedding). Add a CI pre-check:
```bash
python3 -c "from PIL import Image; img=Image.open('assets/icon-256.png'); assert max(img.size)<=256, img.size"
```

---

### Issue #5 — `cd: dist: No such file or directory`

**Symptom:**

```
cd: dist: No such file or directory
Error: Process completed with exit code 1.
```

**Trigger:** A `run:` step in `release.yml` uses `cd dist && ...` assuming the previous step's working directory persists.

**Root cause:** Every `run:` step in GitHub Actions starts at the repository root unless `working-directory:` is set explicitly. CWD does **not** carry over between steps. In a monorepo with multiple modules, this causes silent path drift.

**Fix applied (in sibling reference implementation):**
```yaml
- name: Compress and checksum
  working-directory: <module>/dist
  run: |
    test -d . || { echo "::error::dist missing"; exit 1; }
    for f in <binary>-*; do ...; done
```

**Prevention rule:**
🔴 **NEVER use `cd` inside CI `run:` blocks.** Always use the YAML `working-directory:` field. Guard every output directory with `test -d` before operating on it.

---

### Issue #6 — Job-Level `if` Blocks Required Status Checks

**Symptom:** SHA-deduplication added with job-level `if: needs.sha-check.outputs.already-built != 'true'`. GitHub UI shows grey "skipped" icons; branch protection treats skipped jobs as neither success nor failure → PRs permanently blocked.

**Root cause:** GitHub Actions distinguishes `success`, `failure`, and `skipped` conclusions. Required status checks only accept `success`. A job-level `if` that evaluates false produces `skipped`, which never resolves the gate.

**Fix applied (in sibling reference implementation):** **Passthrough gate pattern** — jobs always run; step-level `if:` skips the actual work; an unconditional first step echoes "✅ Already validated (SHA cached)" so the job always concludes `success`.

**Prevention rule:**
🟠 **Never use job-level `if` for cache/dedup gating** when the job is a required status check. Use **step-level** conditionals and always include at least one unconditional step.

---

### Issue #7 — `cancel-in-progress` Cancels the Marker Job

**Symptom:** A trailing `mark-success` job that writes the `ci-passed-<SHA>` cache entry was intermittently cancelled by `cancel-in-progress: true`, leaving the SHA uncached and forcing full re-runs on the next push.

**Root cause:** When all validation jobs finish and `mark-success` is queued, a new push to the same ref cancels the entire workflow run — including the still-pending marker job.

**Fix applied (in sibling reference implementation):** Inline the cache write as the **final step of the last validation job** (`test-summary`). Guard with `if: success()`.

**Prevention rule:**
🟠 **Side-effects that must persist after success (cache writes, telemetry, deployment markers) belong in the LAST validation job, not in a separate trailing job** when `cancel-in-progress: true` is set.

---

### Issue #8 — `@latest` Tool/Action Installs Are Non-Reproducible

**Symptom:** Builds that passed yesterday fail today with no code changes — a tool installed via `go install foo@latest` or an action pinned to `@v1` (floating major) introduced a breaking change overnight.

**Root cause:** `@latest`, `@main`, and floating major tags resolve to whatever is current at install time. A new upstream release between two CI runs introduces a new lint rule, removed flag, or behavior change.

**Fix applied (in sibling reference implementation):** Pin every tool and action to an exact tag.

| Tool / Action | Pinned Version |
|---|---|
| `golangci-lint` | `vX.Y.Z` |
| `govulncheck` | `vX.Y.Z` |
| `actions/checkout` | `@v6` |
| `actions/setup-go` | `@v6` |
| `actions/cache` | `@v4` |
| `softprops/action-gh-release` | `@v2` |

**Prevention rule:**
🟡 **`@latest` and `@main` are PROHIBITED in any workflow or `setup.sh`.** Pin every action and CLI tool to an exact tag. Use Dependabot/Renovate to propose pinned bumps.

---

### Issue #9 — Release Branch Run Cancelled by Follow-Up Commit

**Symptom:** A push to `release/vX.Y.Z` started the release workflow. A follow-up commit (changelog typo fix) on the same branch cancelled the in-progress run, leaving artifacts half-built and the GitHub Release inconsistent.

**Root cause:** `concurrency.cancel-in-progress: true` is appropriate for PRs but catastrophic for release branches where every commit must produce complete artifacts.

**Fix applied (in sibling reference implementation):**
```yaml

# release.yml

concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: false   # NEVER cancel release runs
```
For shared CI workflows that also run on release branches:
```yaml
cancel-in-progress: ${{ !startsWith(github.ref, 'refs/heads/release/') }}
```

**Prevention rule:**
🔴 **Release workflow `cancel-in-progress` MUST be `false`.** Shared workflows must use a conditional that exempts release branches.

---

### Issue #10 — Install-Script Placeholder Unreplaced

**Symptom:** Released `install.ps1` / `install.sh` shipped with literal `VERSION_PLACEHOLDER` or `REPO_PLACEHOLDER` strings, causing user installs to download a non-existent release.

**Root cause:** The release workflow generated install scripts with `sed` substitution. When `VERSION` was unset (e.g., the version-resolution step failed silently) or `sed` ran on the wrong file, placeholders survived into the published asset.

**Fix applied (in sibling reference implementation):**
```bash
: "${VERSION:?VERSION must be set before generating install scripts}"
sed -i "s|VERSION_PLACEHOLDER|${VERSION}|g; s|REPO_PLACEHOLDER|${GITHUB_REPOSITORY}|g" install.sh install.ps1
! grep -q "PLACEHOLDER" install.sh install.ps1 \
  || { echo "::error::Unreplaced placeholder"; exit 1; }
```

**Prevention rule:**
🟠 **Always validate generated scripts for residual placeholders before upload.** Use `:?` parameter expansion for required env vars. Test install scripts end-to-end on a `release/test-*` branch.

---

### Issue #11 — Missing `GITHUB_TOKEN` Silently Skips Upload

**Symptom:** Local invocations of a release helper produced binaries but no GitHub Release. No error printed — silent skip.

**Root cause:** Upload code checked for `GITHUB_TOKEN` and **returned early without erroring** when absent. Intentional for local dev, but the silence misled both users and CI debuggers.

**Fix applied (in sibling reference implementation):** Print an explicit warning to stderr when assets exist but the token is missing. In CI, treat a missing `GITHUB_TOKEN` as a **hard failure**, not a skip.

**Prevention rule:**
🟡 **Silent skips are a bug.** Always log to stderr when skipping an expected operation. In CI, missing required secrets MUST fail the job, not warn.

---

### Issue #12 — Asset Name Mismatch Between Checksum and Upload

**Symptom:** Users running install scripts saw checksum-verification failures even when the binary was intact.

**Root cause:** The compress step produced `app-vX.Y.Z-windows-amd64.zip`, but the checksum step (running in a different working directory) generated `checksums.txt` listing `app-windows-amd64.zip` (no version). Install scripts looked up the versioned name in a non-versioned manifest → mismatch.

**Fix applied (in sibling reference implementation):**

1. Centralize asset naming as a shell function:
   ```bash
   asset_name() { echo "app-${VERSION}-${OS}-${ARCH}.${EXT}"; }
   ```
2. Generate checksums in the **same `working-directory`** as the artifacts.
3. Round-trip test: extract `checksums.txt`, verify each listed file exists in `dist/`.

**Prevention rule:**
🟠 **Asset naming must be defined once and reused.** Checksum generation MUST run in the artifact directory. Add a pre-publish round-trip test that validates checksums against actual files.

---



---

### Issue #13 — Release Skew & Out-of-Band Tagging

**Date observed:** 2026-04-16
**Workflow:** Release (`release.yml`)
**Tag:** `vX.Y.Z`

**Symptom:**

```
Installer Smoke Test failed: Version mismatch. expected: vX.Y.Z, actual: vX.Y.Z
Error: Process completed with exit code 1.
```

**Trigger:**

A `git tag vX.Y.Z` was pushed out-of-band to trigger a release before `version.json`, `constants.go`, and `changelog.md` had been updated and committed on `main`.

**Root cause:**

1. In Go, `constants.Version` declared as `const` cannot be overridden via `-ldflags` during `go build`.
2. Pushing a tag before synchronizing the repository source creates a race condition where the compiled binary reflects stale code.
3. The release workflow parsed `changelog.md` for ````[vX.Y.Z]````, found none, and emitted "No changelog entry found".

**Fix applied:**

1. Switched `Version` in Go to `var Version = "dev"` so `-ldflags -X` injection works reliably.
2. Synchronized `version.json`, `package.json`, and `changelog.md` on `main` before creating the git tag.
3. Created the AI Release Synchronization standard in `05-release-pipeline.md`.

**Prevention rule:**

🔴 **Source control dictates the tag; tags NEVER dictate source control.** Always run `npm run sync` and push to the `main` branch before creating any release tag.


## Standing Rules (apply to every CI/CD change)

| # | Rule | Rationale | Source |
|---|------|-----------|--------|
| 1 | No `npm ci` / `npm install` in any workflow | Lockfile is not kept in sync; will always fail | Issue #1 |
| 2 | No `actions/setup-node` in any workflow | No Node-based build step exists in this repo | Issue #1 |
| 3 | No built-in `cache:` on language setup actions | No canonical manifests exist for those languages | Issue #2 |
| 4 | Version reads from `package.json` use `sed`, not `node` | Avoid Node toolchain dependency | Issue #3 |
| 5 | All `actions/*` versions pinned to exact tag (no `@latest`/`@main`) | Reproducibility | Issue #8 |
| 6 | Tool versions pinned exactly (e.g. `golangci-lint@vX.Y.Z`) | Reproducibility | Issue #8 |
| 7 | Every code change bumps at least the minor version | Per `.lovable/user-preferences` | — |
| 8 | Touching `release-artifacts/` outside of `release.sh` is forbidden | Generated content; do not hand-edit | — |
| 9 | NEVER `cd` inside `run:` blocks — use `working-directory:` | Steps reset CWD to repo root | Issue #5 |
| 10 | Every output directory used in CI must be guarded with `test -d` | Fail fast with actionable error | Issue #5 |
| 11 | Job-level `if` MUST NOT gate required status checks (use step-level) | `skipped` ≠ `success` in branch protection | Issue #6 |
| 12 | Cache-write side effects belong in the LAST validation job, not a trailing marker job | `cancel-in-progress` kills trailing jobs | Issue #7 |
| 13 | Release workflow `cancel-in-progress` MUST be `false` | Half-built releases corrupt the artifact set | Issue #9 |
| 14 | Shared CI workflows must exempt `refs/heads/release/*` from cancellation | Same reason as #13 | Issue #9 |
| 15 | Required env vars in shell use `:?` parameter expansion | Fails fast instead of producing garbage output | Issue #10 |
| 16 | Generated scripts grep-validated for residual `*PLACEHOLDER` strings before upload | Prevents shipping broken installers | Issue #10 |
| 17 | Silent skips are forbidden; always log to stderr; in CI fail hard on missing required secrets | Surface invisible problems | Issue #11 |
| 18 | Asset names centralized in one function/variable; checksums generated in artifact dir | Prevents name drift between produce/verify | Issue #12 |
| 19 | Any PNG referenced from `winres.json` MUST be ≤ 256×256 px (when Go/Windows release added) | Hard limit of `.ico` format | Issue #4 |

---

## Pre-flight Checklist for Workflow Edits

Before committing changes to `.github/workflows/*.yml`, `release.sh`, `install.sh`, `install.ps1`:

**Repo-specific (always required):**

- [ ] No `npm`, `node`, `setup-node`, or JS dep introduced? → else **STOP**, see Issue #1.
- [ ] No built-in `cache:` on `setup-python` / `setup-go`? → else **STOP**, see Issue #2.
- [ ] All action versions pinned to exact tag (`@v6`, not `@latest`)? — Issue #8
- [ ] All tool versions pinned (e.g. `@vX.Y.Z`)? — Issue #8
- [ ] `release.sh` resolves version without Node? (`grep -E 'node |require\(' release.sh` returns empty) — Issue #3
- [ ] `package.json` version bumped (minor or major)?
- [ ] No `cd` in `run:` blocks? Use `working-directory:` instead. — Issue #5
- [ ] Every directory used has a `test -d` guard? — Issue #5

**If pipeline grows (apply when relevant):**

- [ ] Required-status-check jobs have NO job-level `if`? — Issue #6
- [ ] Cache writes are inlined into the last validation job, not a trailing job? — Issue #7
- [ ] `release.yml` has `cancel-in-progress: false`? — Issue #9
- [ ] Shared CI workflows exempt `refs/heads/release/*`? — Issue #9
- [ ] Required env vars use `${VAR:?msg}` syntax? — Issue #10
- [ ] Generated install scripts grep-checked for residual `PLACEHOLDER`? — Issue #10
- [ ] Missing required secrets fail loudly (no silent skip)? — Issue #11
- [ ] Asset naming centralized; checksum generated in artifact `working-directory`? — Issue #12
- [ ] Any `winres.json`-referenced PNG is ≤ 256×256? — Issue #4

- [ ] This ledger updated if a new failure mode was discovered.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| CI pipeline spec | [`./02-ci-pipeline.md`](./02-ci-pipeline.md) |
| Release pipeline spec | [`./05-release-pipeline.md`](./05-release-pipeline.md) |
| Shared conventions | [`./03-shared-conventions.md`](./03-shared-conventions.md) |
| Current CI workflow | `.github/workflows/ci.yml` |
| Current Release workflow | `.github/workflows/release.yml` |
| Current release script | `release.sh` |
| Coding standards | `02-spec/02-coding-guidelines/01-index.md` |

---

## Update Protocol

When a new CI/CD failure occurs:

1. Capture the **exact** error message (do not paraphrase).
2. Add a new `### Issue #N` section using the template structure above.
3. Add a corresponding rule to **Standing Rules** if the failure mode is generalizable.
4. Update the **Pre-flight Checklist** if a new check is warranted.
5. Bump the version of this document in the front-matter.
6. Cross-reference the issue from the affected spec (e.g., `05-release-pipeline.md`).
