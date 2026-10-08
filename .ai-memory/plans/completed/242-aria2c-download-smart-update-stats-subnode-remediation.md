# Completed Plan 242: aria2c Download, Smart Update, Diffstat Isolation & Remediation Sub-Node Tree

**Plan Slug:** `242-aria2c-download-smart-update-stats-subnode-remediation`  
**Parent Task ID:** `Task-242`  
**Status:** `COMPLETED`  
**Completed Date:** `2026-10-08`  
**Specs:**
- Architecture Spec: `02-spec/21-app/242-aria2c-download-smart-update-stats-subnode-remediation/01-architecture-spec.md`
- Component Spec: `02-spec/21-app/242-aria2c-download-smart-update-stats-subnode-remediation/02-download-and-update-spec.md`

---

## 1. Executive Summary & Problem Resolution

This milestone resolves five core functional defects and architectural requirements:
1. **Diffstat & Commit Counting Isolation (Image 2)**: Eliminated identical metrics `+167988/-14766 (1095)` across multiple repositories during pull/status by adding strict empty-path and non-existent path validation in `cli/gitutil/last_sha.go` (`GetLastCommitSHA`), `cli/cmdpull/pull_worker.go` (`queryGitDiffStat`, `CaptureRepoPrePullState`, `CaptureRepoPostPullState`), and explicitly setting `cmd.Dir = cleanDir` in `cli/cmdpull/pull.go` (`executeGitPullCommand`).
2. **Remediation Sub-Node Tree UX Restructuring (Image 1)**: Converted flat batch pull remediation output (`→ Remediating <repo>...` and `✓ [<repo>] Successfully updated.`) into structured parent repository nodes (`  • <repo>`) with hierarchical tree sub-nodes (`    ├── Attempting auto-merge...`, `    └── ✓ Successfully updated.`, `    └── ✗ Merge conflict detected...`).
3. **`gitmap download` Subcommand**: Implemented high-speed download engine (`cli/cmddownload/`) with automated `aria2c` detection and multi-connection execution (`-s 16 -x 16 -k 1M`), 3-tier fallback to `curl` and Go `net/http`, and centered symmetric visual progress bar rendering.
4. **Smart `gitmap update` Engine**:
   - Explicit version target support via positional argument or `--version <v>`.
   - Fast short-circuit skip if current version matches target version (`Already updated`).
   - Up to 5-tag automatic fallback loop when candidate releases lack compiled binary assets.
   - 24-hour daily caching in Split SQLite `ReleaseCache` table (`ReleaseCacheId`, `Tag`, `Version`, `HasExecutables`, `AssetUrl`, `CheckedAt`, `ExpiresAt`).
   - Concise interactive output displaying numbers and progress, plus structured `--json` payload support.
   - Synchronized runner scripts (`install.ps1`, `gitmap.ps1`, `run.ps1`).
5. **Minor Version Bump**: Incremented repository to `v6.513.0`.

---

## 2. Completed Subtasks & Verified Evidence

### Subtask 01: Diffstat Isolation & Remediation Tree UI
- **Subtask Code:** `Task-01`
- **Owner:** `Worker 01`
- **Files Modified:**
  - `cli/gitutil/last_sha.go`
  - `cli/cmdpull/pull_worker.go`
  - `cli/cmdpull/pull.go`
  - `cli/cmdpull/pull_remediation.go`
  - `cli/cmdpull/pull_missing_remediate.go`
- **Verification:** `go vet -C cli ./gitutil ./cmdpull` passed (exit 0); `go test -C cli -v ./cmdpull` passed (exit 0).

### Subtask 02: aria2c Download Method & Smart Update Engine
- **Subtask Code:** `Task-02`
- **Owner:** `Worker 02`
- **Files Created/Modified:**
  - `cli/cmddownload/aria2c.go`
  - `cli/cmddownload/curl.go`
  - `cli/cmddownload/download.go`
  - `cli/cmddownload/download_test.go`
  - `cli/cmddownload/http_stream.go`
  - `cli/cmddownload/progress.go`
  - `cli/cmddownload/types.go`
  - `cli/cmd/rootutility.go`
  - `cli/cmdupdate/types.go`
  - `cli/cmdupdate/update_fallback_loop.go`
  - `cli/cmdupdate/update_installer_cmd.go`
  - `cli/cmdupdate/updateremoteinstall.go`
  - `cli/cmdupdate/update_test.go`
  - `cli/store/release_cache.go`
  - `cli/store/store.go`
  - `install.ps1`
  - `install.sh`
- **Verification:** `go vet -C cli ./cmddownload ./cmdupdate ./store ./cmd` passed (exit 0); `go test -C cli -v ./cmddownload ./cmdupdate` passed (exit 0).

---

## 3. Invariants & Acceptance Gates
- **Zero Raw Grep / Select-String**: All searches executed exclusively via `gitmap`.
- **Strict Relative Git Paths**: Zero absolute paths across documentation and manifests.
- **Positive Booleans**: All boolean flags (`HasExecutables`, `IsForce`, `IsQuiet`, `IsJSON`) use affirmative naming.
- **Zero Intermediate Commits**: All deliverables finalized for atomic final commit.
