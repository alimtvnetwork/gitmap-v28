# Architecture Specification: GitMap Pull Errors and UI Enhancement
# Spec Sequence: 338
# Slug: gitmap-pull-errors-and-ui-enhancement

## 1. Overview & Problem Definition

During high-speed multi-repository pull operations (`gitmap pa` / `gitmap pull-all` / `gitmap fix --all`), three systemic issues recurred:
1. **Repeated Uncommitted Casing Normalizations**:
   - `reconcileAllCasingDiscrepancies` in `cli/cmdfix/fix_cmd.go` detects casing collisions (e.g. `README.md` -> `readme.md`) across tracked repositories and applies git/filesystem renames (`ExecuteCasingNormalization`).
   - However, the normalization was staged via `git add -A` but **never committed** to the local branch.
   - During subsequent batch runs (or when batch stash-autostash runs in parallel across dirty repositories), git stash or merge cycles reset or restaged the files, repeating the warning: `• Normalized casing: README.md -> readme.md` indefinitely.
2. **Missing Repository Failures & Stale Database Records**:
   - Registered repositories missing from disk (`lara-publishing-v1`, `laravel-automation-v1`, etc.) failed `git pull` with:
     `directory exists but is not a Git repository (missing .git)`.
   - Batch remediation attempted `RemediateMissingRepo`, printing:
     `✗ Cannot clone: no remote URL recorded in database.`
   - While `purgeInvalidRepoEntry(f.RepoPath)` deleted the dead row from `gitmap.db`, it failed to clear the recorded error from the Pull Split-DB (`store.OpenPullSplitDB()`).
   - Consequently, every subsequent run of `gitmap fix --all` re-read the stale error from `PullSplitDB.QueryAllLatestPullErrors`, prompting or failing repeatedly with `Found 7 failed repository(ies) from recent runs`.
3. **UI Crispness, Icons & Mail Composer Enhancements**:
   - The favicon and logo in `cli/cmdlogviewer/favicon.go` and UI assets featured low-contrast borders, solid background cards, or fuzzy details.
   - The mail composer Web UI (`gitmap mail ui`) lacked multi-mailbox status inspection, account switching headers, recipient pills, and mail viewing capabilities inspired by `go-email-reader` (without any warm-up dependencies).

## 2. Root Cause Analysis (RCA) Matrix

| Issue | Root Cause | Architectural Solution |
| :--- | :--- | :--- |
| **Casing Re-normalization Loop** | `ExecuteCasingNormalization` in `cli/cmdfix/fix_casing.go` runs `git mv` and `git add -A`, but does not commit the normalization to git. Uncommitted changes are stashed, dirty, or reverted across pulls. | Create `CommitCasingNormalization(repoPath, d)` with dedicated semantic commit `cpc "chore: normalize %s to %s"` if working tree has staged changes, ensuring changes persist permanently. |
| **Ghost Missing Repo Loop** | `RemediateMissingRepo` in `cli/cmdpull/pull_missing_remediate.go` calls `purgeInvalidRepoEntry(f.RepoPath)` to delete from `gitmap.db`, but does NOT evict the error record from `PullSplitDB` via `ClearPullErrorsForRepo`. | Add `_ = db.ClearPullErrorsForRepo(f.RepoName, f.RepoPath)` to `purgeInvalidRepoEntry` and `RemediateMissingRepo` upon unrecoverable clone failure so stale errors are evicted immediately. |
| **Fuzzy / Subpar UI Icons** | Favicons and SVGs in `favicon.go` had solid background cards and tight margins violating the `logo-and-icon-design` spec (requires 100% transparent canvas and 8-10% quiet zone). | Upgrade `GitMapLogoSVG`, embedded favicon ICO/SVG, and UI icon assets with sharp paths, Catppuccin color tokens (`#89b4fa`, `#a6e3a1`, `#cdd6f4`), transparent background, and precision vector paths. |
| **Mail Composer & Inbox UI** | `gitmap mail ui` only provided basic sending without mailbox overview, sender switching, or message previewing. | Modernize `cli/cmdmail/mail_ui_html.go` with Mailbox Inspector pane, sender status badges, rich message viewer, and clean Catppuccin Mocha styling without warm-up bloat. |

## 3. Strict Architectural Boundary Rules
- Zero warm-up features or imports (`internal/warmup` strictly banned).
- All relative paths throughout code, specs, and release documents.
- 8-15 line function limits with positive booleans (`is*`, `has*`).
