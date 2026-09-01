# docs-site Not Bundled in Release & Swallowed Errors Audit

## Issue 1: `gitmap-v28 hd` Fails — docs-site Not Included in Release

### Symptom

```
PS D:\scripts-fixture> gitmap-v28 hd
  ✗ Docs site directory not found at D:\dev\GitMap\docs-site (operation: resolve, reason: directory does not exist)
```

### Root Cause

The release workflow in `release/workflowfinalize.go` (`pushAndFinalize`) collects
Go cross-compiled binaries, zip group archives, and ad-hoc assets — but never
bundles `docs-site/dist/` alongside the binary. The `help-dashboard` command
resolves `docs-site/` relative to the binary directory via `resolveBinaryDir()`,
so the folder must be co-located with the installed binary.

### Fix

1. Add a `buildDocsSiteAsset` function that copies `docs-site/dist/` into the
   release staging directory as a zip archive (`docs-site.zip`).
2. Call it from `pushAndFinalize` so the archive is uploaded as a release asset.
3. Update install scripts to extract `docs-site/` next to the binary after download.

---

## Issue 2: Directory Casing — `GitMap` vs `gitmap-v28`

### Symptom

The error path shows `D:\dev\GitMap\docs-site` — PascalCase `GitMap` directory.
`constants_doctor.go` correctly uses `GitMapSubdir = "gitmap-v28"` (lowercase), but
the deployment directory on Windows is created as `GitMap` by install scripts.

### Root Cause

The install/deploy scripts (PowerShell `install.ps1`) create the target directory
using PascalCase naming. The Go code itself is consistent with lowercase, but
the directory on disk doesn't match because it was created by external scripts.

### Fix

Ensure all references use lowercase `gitmap-v28` consistently. The `resolveBinaryDir()`
already resolves from the actual executable path, so this is primarily an install
script concern — but the Go code should also not assume casing.

---

## Issue 3: Swallowed Errors Audit

### Root Cause

Multiple call sites ignore error return values without logging or handling,
violating the project's Code Red error management rule.

### Instances Found & Fixes Applied

| File | Line | Call | Fix |
|------|------|------|-----|
| `cmd/amendexec.go` | 68-69 | `exec.Command(...).Output()` — git author lookup | Log warning on failure, use empty fallback |
| `cmd/amendexecprint.go` | 37, 42 | `exec.Command("git", "config", ...).Output()` | Log warning on failure |
| `cmd/amendlist.go` | 136 | `json.MarshalIndent(...)` | Check error, print stderr |
| `cmd/bookmarklist.go` | 64 | `json.MarshalIndent(...)` | Check error, print stderr |
| `cmd/clone.go` | 99, 138 | `filepath.Abs(...)` | Check error, print stderr |
| `cmd/diffprofiles.go` | 92 | `json.MarshalIndent(...)` | Check error, print stderr |
| `cmd/doctorchecks.go` | 37, 191 | `filepath.Abs(...)` | Check error, print stderr |
| `cmd/doctorfixpath.go` | 42, 56 | `filepath.Abs(...)` | Check error, print stderr |
| `cmd/doctorversion.go` | 36, 108-109 | `filepath.Abs(...)` | Check error, print stderr |
| `cmd/history.go` | 149 | `json.MarshalIndent(...)` | Check error, print stderr |
| `cmd/installnpp.go` | 95 | `filepath.Abs(...)` | Check error, use original path |
| `localdirs/migrate.go` | 56, 61, 69, 82-83, 85-86, 119 | WalkDir errors, Rel, MkdirAll, Read/WriteFile | Log errors instead of silent skip |
| `release/assets.go` | 185 | `os.RemoveAll(...)` in CleanupStagingDir | Log warning on failure |
| `release/selfrelease_resolve.go` | 145 | `db.SetSetting(...)` | Log warning on failure |
| `release/workflow.go` | 78, 181 | `CurrentCommitSHA()`, `CurrentBranchName()` | Log warning on failure |

### Classification

- **Must fix**: `json.MarshalIndent` failures, `filepath.Abs` in paths shown to user, file I/O in migration
- **Best-effort OK**: `git config` lookups (fallback to empty), `os.RemoveAll` cleanup, browser open

---

## Status

- [x] Issue documented
- [x] Fixes applied
- [ ] Verified with `golangci-lint run ./...`

> **Issue 1 user-visible failure mode superseded by**
> `spec/02-app-issues/34-hd-hosted-docs-fallback.md` (v6.0.x): `gitmap hd`
> no longer hard-exits when `docs-site/` cannot be materialized — it falls
> back to the hosted docs at `constants.DocsURL`. The bundling fix here
> remains the happy path; #33 is the safety net.

## Contributors

- AI-assisted audit and fix
