# Specification 281: Pull Summary Transparency, Remediation Menu Guidance, Safe-Pull Conflict Hardening & LFS Dirty Status Remediation

> **Specification Slug:** `281-pull-summary-remediation-help-and-alim-seo-fix`  
> **Status:** `Active / Authoritative`  
> **Date:** `2026-10-17`  
> **Target Release:** `v6.565.0`  

---

## 1. Executive Summary & Problem Statement

During multi-repository operations (`gitmap pa` / `gitmap pull-all`, `gitmap fix --all`), four operational friction points were identified:

1. **Pull Progress & Summary Opacity:**
   - In multi-worker parallel pull (`-p 4` up to `-p 20`), the progress line hard-capped active repo names to 2 (`(+N more)`), and standard 80-character terminals truncated active repo names.
   - The post-pull summary in `cli/cmdpull/pull_efficient_render.go` only printed `Updated Repositories`, `Dirty Repositories`, and `Failed Repositories`. Up-to-date repositories were omitted from the categorized list, leaving users without explicit confirmation of which repositories completed successfully.

2. **Remediation Prompt Confusion:**
   - In `cli/cmdpull/pull_remediation.go`, the prompt `[?] Remediate failed repositories?` presented options `[1/a] Auto-merge`, `[2/s] Step through`, `[3/r] Rebase` without any explanatory guidance. Users lacked information on what actions and git commands these options would execute, what safety guarantees were provided, and what risk of conflict existed.

3. **Safe-Pull Conflict Detection & Destructive Abort:**
   - When auto-merge encountered conflicts in `cli/cloner/safe_pull.go`, `cloner.SafeAbortMerge()` was called immediately, executing `git merge --abort`. This wiped out the merge state, leaving the working tree seemingly clean and causing `gitmap conflicts` / `gitmap resolve` to report 0 conflicts.
   - `isMergeConflict` in `safe_pull.go` checked only `err.Error()` instead of command output, failing to detect `"unmerged files"` / `"unresolved conflict"` and causing 4 blind retry attempts on exit status 128.

4. **`alim-seo-writing` Persistent Dirty Status:**
   - In `alim-seo-writing`, 11 avatar image files (`assets/avatars/*.jpg`) were committed via the GitHub web UI as 562KB raw blobs rather than Git LFS pointers. Because `.gitattributes` mandates `filter=lfs`, local Git runs `git-lfs clean` which hashes them to 131-byte LFS pointers. Git identifies `Bin 562829 -> 131 bytes` and permanently marks the repository dirty (+11 modified), blocking automated pulls.

---

## 2. Architectural Scope & Target Modules

1. **`cli/cmdpull/` (Progress & Post-Pull Reporting):**
   - Enhance `cli/cmdpull/pull_progress_bar_workers_summary.go` to show up to 4 active workers dynamically based on terminal width.
   - Update `cli/cmdpull/pull_efficient_render.go` to categorize and render `Up-to-Date Repositories (%d)` using compact wrapping so users have complete visibility into all pulled repositories.
2. **`cli/cmdpull/pull_remediation.go` (Remediation Menu Guidance):**
   - Add `[h/?/help]` choice to `printRemediationPromptMenu()`.
   - Implement `renderRemediationHelpGuide()` displaying detailed explanations of Auto-Merge, Step-Through, and Rebase operations before execution.
   - Wrap prompt in a loop so viewing help returns cleanly to the menu.
3. **`cli/cloner/safe_pull.go` & `cli/cloner/safe_pull_conflicts.go` (Conflict Detection Hardening):**
   - Extend `isMergeConflict` to inspect both error envelope and git output for keywords (`"unmerged files"`, `"unresolved conflict"`, `"automatic merge failed"`, `"conflict ("`).
   - Prevent 4 blind retries when an unmerged state or conflict is already active.
4. **`cli/cmdpending/nc_catalog_gen.go` & `cli/cmdpending/ncgen/nc_annotations.json` (Drift Fix):**
   - Sync `nc_catalog_gen.go` generated entries (`Alias: "rec"` for reclone) and verify `go generate ./...` produces 0 drift.
5. **`alim-seo-writing` (LFS Pointer Normalization):**
   - Stage and commit the clean 131-byte LFS pointer files for the 11 avatar images in `alim-seo-writing` so working tree matches `.gitattributes` and HEAD.
