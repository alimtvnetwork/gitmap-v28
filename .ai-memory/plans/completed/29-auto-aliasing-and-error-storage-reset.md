# 29 — Auto-Aliasing Engine, Scan Bracket/Tree Views & Error Storage Reset Suite

- **Slug:** auto-aliasing-and-error-storage-reset
- **Date:** 2026-09-19
- **Version:** v6.262.0
- **Status:** completed
- **Execution Loops:** 4 completed subtasks across 2 parallel execution phases (Budget: N=250, Completed in 4 steps).

---

## User Request (Verbatim)

```text
When, okay, so new, new, uh, let's say feature in, uh, in Git map, it would be auto aliasing. Uh, we can see the aliasing when you do the scan or see the table that show as a bracket. Uh, the default aliasing, if it has too much, then it would be as a tree view of multiple aliasing would, would show up. So this is how I want. So usually how it would, would go is that if the border has a space, every, uh, space first, I mean, uh, first words, first character, words first character would be considered as the, uh, I mean, characters to combine as the alias. For example, anti-gravity manager. So if it's written as hyphen or space, A-A and M would be considered as the, um, considered as the alias. Okay. Similar could go for the Git map. So Git map does not have two words, so it would have GM. Um, so in, in this way, it would try to find what would be the next word that two words combined, and it would try to find that, uh, aliasing. Do you understand this feature? Then we integrate this. And this would happen when we do the scan, and when we, let's say, uh, update the installation, the latest one, we would check the existing ones, and it would create those, um, let's say, aliasing. In some cases, the aliasing cannot be done. For example, the core package, the packages which are, let's say, very one word, uh, cannot have the other understanding. Any package that has a WP, then WP would be prefix. Okay. That would be, uh, not, not relevant or just W. So WP-GL would be Git log. WP-HA would be HTML automated. Uh, things like that. So this type of, uh, code you need to write. Um, also for the presentations, it would have like prep hyphen the short version of the name. So you know that. Mm. So these are the ways. So when it updates, first time it checks if this aliasing is there. Aliasing should have a separate table collection compared to connection with this. And also remember the error storage. We should be able to reset. That command needs to be show up when we do the reverse and see the database. We should have a reset command showing up there so that we can clear that error as well. Remember these few things. Is it clear? Do you have any question, confusion?
```

## Extracted Actionable Task List

1. **Auto-Aliasing Word-Boundary Tokenization & Acronym Generation**: Tokenize on `-`, `_`, space; combine first characters of word tokens (e.g. `anti-gravity-manager` $\rightarrow$ `agm`).
2. **Compound Word Splitting & Acronym Generation**: Sub-word dictionary heuristic for delimiter-less names (e.g. `gitmap` $\rightarrow$ `gm`, `scriptsfixer` $\rightarrow$ `sf`, `gitsync` $\rightarrow$ `gs`).
3. **WordPress Package Prefix Rules**: Names starting with `wp-` preserve `wp-` prefix with abbreviated acronym remainder (e.g. `wp-git-log` $\rightarrow$ `wp-gl`, `wp-html-automated` $\rightarrow$ `wp-ha`).
4. **Presentation Repository Prefix Rules**: Presentation repositories formatted as `prep-` + short name (e.g. `presentations-repos/bsrm-presentation-hiltrax` $\rightarrow$ `prep-bsrm`).
5. **Fallback & Collision Handling**: Gracefully fallback on single-word core packages; resolve collisions against existing DB aliases with numeric suffixes (`agm2`, `gm2`, etc.).
6. **Dedicated Table Collection & Auto-Population**: Store aliases in dedicated `Alias` / `repo_aliases` collection tracking `IsPrimary` and `Source`; auto-populate on `gitmap scan` and binary update/installation migrations.
7. **Scan Table Bracket Display**: Render primary alias in brackets (e.g. `[gm]`, `[wp-gl]`) in scan output headers and block renders.
8. **Multi-Alias Branch Tree View**: When a repository has $>1$ aliases, display a clean branch tree view (`├── ... (primary)`, `└── ... (secondary)`).
9. **Error Storage Reset Command**: Implement `gitmap storage reset-errors` (`error-reset`, `reset`, `clear-errors`) to purge pipeline and scan error logs, error reports, and database error tables.
10. **Database & Storage Reset Guidance**: Prominently display error storage reset guidance whenever inspecting the database (`gitmap storage ls`, `gitmap db`, `gitmap db ls`, `gitmap db status`, and `gitmap repo db status`).

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
