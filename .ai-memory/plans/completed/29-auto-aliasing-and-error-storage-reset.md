# 29 — Auto-Aliasing Engine, Scan Bracket/Tree Views & Error Storage Reset Suite

- **Slug:** auto-aliasing-and-error-storage-reset
- **Date:** 2026-09-19
- **Version:** v6.262.0
- **Status:** completed
- **Execution Loops:** 4 completed subtasks across 2 parallel execution phases (Budget: N=250, Completed in 4 steps).

---

## 1. Task Origin & Problem Statement

The user requested an auto-aliasing feature and error storage reset capabilities for GitMap:
1. **Auto-Aliasing Generation**:
   - Word boundary splitting: hyphens (`-`), underscores (`_`), and spaces (` `). First character of each word combined (e.g. `anti-gravity-manager` $\rightarrow$ `agm`).
   - Compound words: Single words like `gitmap` decomposed by prefix/suffix roots $\rightarrow$ `gm`.
   - Special prefix `wp-`: Packages prefixed with `wp` retain prefix (e.g. `wp-git-log` $\rightarrow$ `wp-gl`, `wp-html-automated` $\rightarrow$ `wp-ha`).
   - Special prefix presentations: Presentation repositories formatted as `prep-` + short name (e.g. `presentations-repos/bsrm-presentation-hiltrax` $\rightarrow$ `prep-bsrm`).
   - Single-word/core packages: Graceful fallback where abbreviations cannot be derived.
2. **Dedicated Table Collection**:
   - Stored in a separate table collection (`Alias` / `repo_aliases`) in SQLite, tracking `IsPrimary` and `Source`.
   - Automatically checked and populated during `gitmap scan` and binary update/installation.
3. **Scan Table Bracket & Multi-Alias Tree Display**:
   - Display primary alias in brackets (e.g. `[gm]`, `[wp-gl]`) in scan outputs.
   - Display a clean branch tree view when a repository has multiple registered aliases.
4. **Error Storage Reset Command & Guidance**:
   - Implement `gitmap storage reset-errors` (`error-reset`, `reset`, `clear-errors`) to purge error caches.
   - Display reset command guidance in the storage inventory footer (`storage ls` and `storage space ls`).

---

## 2. Consolidated Execution Summary

### Subtask 01: Auto-Aliasing Engine & Heuristics
- Implemented `cli/cmd/auto_alias.go` with deterministic alias generation (`GenerateAutoAlias`).
- Supported word boundary splitting across `-`, `_`, and spaces.
- Built compound root matching (`knownCompoundPrefixes` and `knownCompoundSuffixes`), mapping `gitmap` to `gm`, `gitlab` to `gl`, etc.
- Implemented prefix rules for WordPress (`wp-*` $\rightarrow$ `wp-<acronym>`) and presentations (`prep-<token>`).
- Handled collision avoidance with incremental numeric suffixes (`ResolveUniqueAlias`).
- Authored unit test suite in `cli/cmd/auto_alias_test.go` covering all permutations.

### Subtask 02: Dedicated Database Table & Auto-Population on Scan/Install
- Extended `cli/constants/constants_alias.go` and `cli/store/alias.go` to support `IsPrimary` and `Source` columns.
- Added `CreateAliasWithDetails` and `FindAliasesByRepoID`.
- Updated `cli/cmdinstall/install_aliases.go` (`EnsureTrackedRepoAliases`) to leverage `GenerateAutoAlias`.
- Integrated `autoPopulateScanAliases` in `cli/cmdscan/scan.go` to automatically populate aliases for newly discovered repositories.

### Subtask 03: Scan Table Bracket & Multi-Alias Tree View Rendering
- Updated `cli/formatter/terminal.go` to display primary aliases in brackets (e.g. `■ my-package [mp] (main)`).
- Implemented `printRepoAliasesBranchIfMulti` in `terminal.go` to render multi-alias tree views (`├── mp (primary)`, `└── my-pkg (secondary)`).
- Updated `cli/render/repotermblock.go` and `cli/render/adapters.go` with `formatBlockHeader` and `appendBlockAliasesTree`.
- Added unit tests in `cli/render/repotermblock_alias_test.go`.

### Subtask 04: Error Storage Reset Command & Help Documentation
- Implemented `runStorageResetErrors` in `cli/cmd/storage_reset.go`, supporting `--dry-run`, `--verbose`, and `--force`.
- Purges pipeline `.log` and `.json` cache files, scan error reports, and database error records.
- Added `reset-errors`, `error-reset`, `reset`, and `clear-errors` aliases in `cli/cmd/storage_cmd.go`.
- Added prominent command guidance in `printStorageResetGuidance` in `cli/cmd/storage_ls.go`.
- Documented commands and flags in `cli/helptext/storage.md` and `cli/helptext/scan.md`.

---

## 3. Verification & Coding Guidelines Compliance

- **Function Sizing:** All functions refactored to <= 8–15 lines.
- **Single Return Types:** Multi-value return tuples eliminated.
- **Universal AppError:** Replaced generic error returns with strongly-typed `*apperror.AppError`.
- **Boolean Conventions:** All boolean variables and parameters use affirmative `is*` and `has*` prefixes.
- **Line Endings:** Strictly Unix LF (`\n`) across all modified and new files.
