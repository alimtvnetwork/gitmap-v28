# 07 — Release Pipeline: Known Issues, Root Causes & Fixes

> **Purpose**: Self-contained post-mortem catalog of every release-pipeline failure. Any AI model or engineer reading this should be able to diagnose and fix the same class of issue instantly.
>
> **Scope**: Issues specific to the **release** workflow (`.github/workflows/release.yml`) — version resolution, cross-compilation, packaging, install-script generation, GitHub Release creation, asset upload.
>
> For CI-pipeline issues (lint, test, build matrix), see [`spec/09-pipeline/10-known-issues-and-fixes.md`](../09-pipeline/10-known-issues-and-fixes.md).

---

## Issue Index

| # | Issue | Stage | Severity | Status |
|---|-------|-------|----------|--------|
| 1 | `go-winres` icon > 256x256 | resource embed | 🔴 Blocker | ✅ Fixed v2.81.0 |
| 2 | `cd: dist: No such file or directory` | compress | 🔴 Blocker | ✅ Fixed v2.54.0 |
| 3 | `npm ci` lockfile drift (docs-site sub-build) | docs-site | 🔴 Blocker | ✅ Fixed v2.82.0 |
| 4 | Release branch run cancelled by follow-up commit | concurrency | 🔴 Blocker | ✅ Fixed |
| 5 | Install script version-pinning placeholder unreplaced | script generation | 🟠 High | ✅ Fixed |
| 6 | Missing `GITHUB_TOKEN` silently skips upload | asset upload | 🟡 Medium | ✅ Mitigated |
| 7 | Asset name mismatch between checksum and upload | packaging | 🟠 High | ✅ Fixed |

---

## Issue #1 — `go-winres` Icon Size > 256x256

### Symptom

```
2026/04/16 16:26:46 image size too big, must fit in 256x256
Error: Process completed with exit code 1.
```

Fails during the `go-winres make` step before any cross-compilation begins.

### Root Cause

Windows `.ico` resource format hard-limits each image frame to **256x256 pixels**. `go-winres` reads `gitmap-v27/winres/winres.json`, opens the referenced PNG, and converts it to `.ico` — refusing any image larger than the limit.

The project shipped a 512x512 `icon.png` for web/docs use, then mistakenly referenced it from `winres.json`.

### Why It Wasn't Caught

- Local Windows builds (`run.ps1`) skip `go-winres` entirely — they produce a metadata-less binary.
- The `go build` step succeeds without resource embedding.
- No automated dimension check existed for files referenced by `winres.json`.

### Fix (v2.81.0)

```diff
  // gitmap-v27/winres/winres.json
- "Path": "../assets/icon.png"
+ "Path": "../assets/icon-256.png"
```

Created `gitmap-v27/assets/icon-256.png` (256x256). Kept original 512x512 `icon.png` for non-Windows uses.

### Prevention Rules

1. **Icons referenced from `winres.json` MUST be ≤ 256x256.**
2. **Maintain separate files** by purpose: `icon.png` (web/docs), `icon-256.png` (Windows resource).
3. **Add a pre-check in CI** before `go-winres make`:
   ```bash
   python3 -c "from PIL import Image; img=Image.open('gitmap-v27/assets/icon-256.png'); assert max(img.size)<=256"
   ```
4. **Document the constraint inline** in `gitmap-v27/winres/winres.json` via a `_comment` field.

### Related Files

- `gitmap-v27/winres/winres.json`
- `gitmap-v27/assets/icon-256.png` / `icon.png`
- `.github/workflows/release.yml`
- `spec/08-generic-update/09-winres-icon-constraint.md`

---

## Issue #2 — `cd: dist: No such file or directory`

### Symptom

```
cd: dist: No such file or directory
Error: Process completed with exit code 1.
```

Fails during the compress/checksum step.

### Root Cause

In a monorepo with `gitmap-v27/` and `gitmap-updater/`, each `run:` step in GitHub Actions starts at the repository root unless `working-directory:` is set. A `cd dist` command assumed the previous step's CWD persisted — it doesn't.

### Why It Wasn't Caught

Locally, `run.ps1` always executes from `gitmap-v27/`, so the relative path resolves. CI is the only environment with the root-relative behavior.

### Fix (v2.54.0)

```yaml
- name: Compress and checksum
  working-directory: gitmap-v27/dist
  run: |
    for f in gitmap-*; do
      [ -f "$f" ] || continue
      ...
    done
```

### Prevention Rules

1. **NEVER use `cd` inside CI scripts.** Use the YAML `working-directory:` field.
2. **Guard with `test -d`** before operating on directories: `test -d gitmap-v27/dist || { echo "::error::dist missing"; exit 1; }`.
3. **Test on `release/test-*` branch** before promoting.

### Related Files

- `.github/workflows/release.yml`
- `spec/02-app-issues/13-release-pipeline-dist-directory.md`

---

## Issue #3 — `npm ci` Lockfile Drift in Docs-Site Build

### Symptom

```
Building docs-site...
npm error code EUSAGE
npm error `npm ci` can only install packages when your package.json and
npm error package-lock.json are in sync.
npm error Missing: vitest@3.2.4 from lock file
... (100+ entries)
```

### Root Cause

The release pipeline builds the documentation site as part of the release artifact bundle. New devDependencies (`vitest`, `@testing-library/*`, `axios`, `framer-motion`, etc.) were added to `package.json` without regenerating `package-lock.json`. `npm ci` refuses to proceed when the two files diverge — by design.

### Why It Wasn't Caught

- Local development uses `npm install` / `bun install`, which auto-update lockfiles.
- `package-lock.json` is read-only to AI tooling — must be regenerated via shell.
- No early lockfile-validation step ran before the docs-site build.

### Fix (v2.82.0)

```bash
rm -f package-lock.json
npm install --package-lock-only --ignore-scripts
npm ci --ignore-scripts --dry-run  # verification
```

Commit the regenerated `package-lock.json` alongside the `package.json` change.

### Prevention Rules

1. **Every `package.json` dependency change MUST be committed with a regenerated `package-lock.json`.**
2. **Add an early CI gate** that runs `npm ci --dry-run` to fail in seconds rather than minutes.
3. **Never hand-edit `package-lock.json`.** Always regenerate.
4. **Mirror the rule for Go**: `go.sum` must be committed alongside `go.mod` changes.

### Related Files

- `package.json`, `package-lock.json`, `bun.lock`
- `.github/workflows/release.yml` (docs-site step)

---

## Issue #4 — Release Branch Run Cancelled by Follow-Up Commit

### Symptom

A push to `release/v2.5x.0` started the release workflow. A follow-up commit (changelog fix, typo) on the same branch cancelled the in-progress run. Result: GitHub Release created with partial assets, or no release at all.

### Root Cause

`concurrency.cancel-in-progress: true` is appropriate for PRs and feature branches, but catastrophic for release branches where every commit must produce complete artifacts.

### Fix

```yaml
# release.yml
concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: false  # NEVER cancel release runs
```

For shared CI workflows that also run on release branches, use a conditional:

```yaml
cancel-in-progress: ${{ !startsWith(github.ref, 'refs/heads/release/') }}
```

### Prevention Rules

1. **Release workflow `cancel-in-progress` MUST be `false`.**
2. **Shared workflows** must use a conditional that exempts release branches.
3. **Document this rule** in any new workflow file's header comment.

### Related Files

- `.github/workflows/release.yml`
- `.github/workflows/ci.yml`
- `.lovable/memory/tech/ci-release-automation.md`

---

## Issue #5 — Install Script Version-Pinning Placeholder Unreplaced

### Symptom

Released `install.ps1` / `install.sh` contained literal `VERSION_PLACEHOLDER` or `REPO_PLACEHOLDER` strings, causing user installs to attempt downloading a non-existent release.

### Root Cause

The release workflow generates install scripts inline via `sed` substitution:

```bash
sed -i "s|VERSION_PLACEHOLDER|${VERSION}|g; s|REPO_PLACEHOLDER|${GITHUB_REPOSITORY}|g" install.sh
```

If `VERSION` was unset (e.g., version-resolution step failed silently) or `sed` ran on the wrong file, placeholders survived into the published asset.

### Fix

1. **Fail fast** if `VERSION` is empty:
   ```bash
   : "${VERSION:?VERSION must be set before generating install scripts}"
   ```
2. **Verify substitution** before upload:
   ```bash
   ! grep -q "PLACEHOLDER" install.sh install.ps1 || { echo "::error::Unreplaced placeholder"; exit 1; }
   ```

### Prevention Rules

1. **Always validate generated scripts** for residual placeholders before upload.
2. **Use `:?` parameter expansion** for required environment variables.
3. **Test install scripts end-to-end** on a `release/test-*` branch.

### Related Files

- `.github/workflows/release.yml`
- `.lovable/memory/tech/ci-release-automation.md`

---

## Issue #6 — Missing `GITHUB_TOKEN` Silently Skips Upload

### Symptom

Local invocations of `gitmap-v27 release --bin` produced binaries but no GitHub Release. No error printed, just silence.

### Root Cause

`gitmap-v27/release/workflowgithub.go::uploadToGitHub` checks for `GITHUB_TOKEN` and **returns early without erroring** if absent. This is intentional for non-CI use, but the silence misled users.

### Fix

Print an explicit warning to stderr when assets exist but token is missing:

```go
if len(token) == 0 {
    if len(assets) > 0 {
        fmt.Fprint(os.Stderr, constants.ErrAssetNoToken)
    }
    return
}
```

(See `gitmap-v27/release/workflowgithub.go` lines 13–20.)

### Prevention Rules

1. **Silent skips are a bug.** Always log to stderr when skipping an expected operation.
2. **Use `constants.Err*` for all error messages** — never inline strings.
3. **In CI**, treat a missing `GITHUB_TOKEN` as a hard failure (not a skip).

### Related Files

- `gitmap-v27/release/workflowgithub.go`
- `gitmap-v27/constants/constants_release.go` (Err* messages)

---

## Issue #7 — Asset Name Mismatch Between Checksum and Upload

### Symptom

Users running install scripts saw checksum-verification failures even when the binary was intact.

### Root Cause

The compress step produced `gitmap-v27.56.0-windows-amd64.zip`, but the checksum step (which ran in a different working directory) generated `checksums.txt` listing `gitmap-windows-amd64.zip` (without the version). Install scripts looked up the versioned name in a non-versioned manifest → mismatch → "tampered binary" error.

### Fix

1. **Single source of truth for asset naming** — define a shell function:
   ```bash
   asset_name() { echo "gitmap-${VERSION}-${OS}-${ARCH}.${EXT}"; }
   ```
2. **Generate checksums in the same `working-directory`** as the compressed artifacts.
3. **Round-trip test**: extract checksums.txt, verify each listed file exists in dist/.

### Prevention Rules

1. **Asset naming is centralized** — define once, reuse everywhere.
2. **Checksum generation runs in the artifact directory** with the same naming convention.
3. **Pre-publish round-trip test** validates checksums against actual files.

### Related Files

- `.github/workflows/release.yml`
- `spec/07-generic-release/04-checksums-verification.md`

---

## Universal Release-Pipeline Prevention Checklist

Apply before merging any change to `release.yml`:

- [ ] No `cd` in `run:` blocks — use `working-directory:`
- [ ] All action versions pinned to exact tags
- [ ] `cancel-in-progress: false` (or release-branch conditional)
- [ ] Lockfiles regenerated for every dependency change
- [ ] `winres.json`-referenced icons ≤ 256x256
- [ ] Required env vars guarded with `:?` parameter expansion
- [ ] Generated scripts validated for residual placeholders
- [ ] Asset naming convention centralized
- [ ] Checksum generation in same directory as artifacts
- [ ] Tested on `release/test-*` branch before merge

---

## Cross-References

- [01-cross-compilation.md](./01-cross-compilation.md) — Cross-compilation matrix
- [02-release-pipeline.md](./02-release-pipeline.md) — Release pipeline overview
- [03-install-scripts.md](./03-install-scripts.md) — Install script generation
- [04-checksums-verification.md](./04-checksums-verification.md) — SHA-256 manifest format
- [spec/09-pipeline/10-known-issues-and-fixes.md](../09-pipeline/10-known-issues-and-fixes.md) — CI-pipeline issue catalog
- [spec/02-app-issues/13-release-pipeline-dist-directory.md](../02-app-issues/13-release-pipeline-dist-directory.md) — `dist` post-mortem
- [spec/08-generic-update/09-winres-icon-constraint.md](../08-generic-update/09-winres-icon-constraint.md) — winres icon post-mortem
- [spec/12-consolidated-guidelines/16-cicd.md](../12-consolidated-guidelines/16-cicd.md) — Consolidated CI/CD guidelines
