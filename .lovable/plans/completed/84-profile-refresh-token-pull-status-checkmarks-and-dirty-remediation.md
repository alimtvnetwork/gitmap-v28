# Master Plan: Google Profile Refresh Token Export, Pull/Status Table Check Marks, Accurate Dirty Numbers & Direct Remediation Commands

## 1. Overview & Root Cause Analysis

### Problem Statement
1. **Google Profile Export & Refresh Tokens**:
   - The user asked whether exporting Google profiles includes the refresh token and if it is done properly.
   - When Google profiles are imported, Chrome was prompting for sign-in again.
   - While `gitmap/cmd/chromeprofile_tokens.go` extracts tokens from `Web Data` (`token_service` table) during single-profile exports (`writeChromeExport`), `chromeprofile_export_all.go` (`loadSingleChromeProfileExport`) **omitted** calling `readChromeTokenService(srcPath)`!
   - Multi-profile JSON, YAML, and SQLite exports therefore dropped all OAuth refresh tokens.
   - In addition, Chrome's refresh tokens are OS-key encrypted (DPAPI on Windows, Secret Service/GNOME Keyring on Linux). Full cross-machine recovery without re-login requires capturing both token service blobs and profile preferences.

2. **Pull & Status Table Check Mark Display**:
   - In `gitmap/cmd/pull_table_row.go`, the table renders raw uppercase string statuses (`UP_TO_DATE`, `synced`, `DIRTY`) rather than an elegant check mark (`✓` / `✔`).
   - `gitmap/cmd/pull_table_style.go` contained a dead-coded `formatPullStatus` with `//nolint:unused` that was bypassed.
   - In status and project listings, users prefer seeing green checkmarks (`✔ active` or `✔ up-to-date`) rather than plain text `active` or `UP_TO_DATE`.

3. **Dirty Numbers Accuracy**:
   - There are conflicting dirty status parsers (`gitutil.Status` in `gitutil.go` vs `gitutil.InspectDirtyState` in `dirty_inspect.go`).
   - In `dirty_inspect.go`:
     - `classifyDirtyFile`: Any line with `M` (such as `M ` staged modification) is misclassified as modified instead of staged because of `strings.Contains(prefix, "M")`.
     - `collectReasonParts`: `StagedCount` is completely omitted from the reason breakdown!
     - In `gitutil.go`: `parsePortcelainStatus` double-counts files with both staged and unstaged modifications (`MM`).
   - This causes `gitmap pull` and `gitmap status` to report incorrect or contradictory file counts.

4. **Direct Remediation Commands & Prompting**:
   - In `gitmap/cmd/remediation_box.go` (`printPendingReposList`), dirty repositories are listed, but no direct copy-pasteable runnable command (e.g. `git -C <absPath> stash -u` or `gitmap fix <slug> 1`) is provided per repository.
   - No interactive prompting exists in `gitmap pull` to ask whether the user wants to apply the remediation immediately or skip.

---

## 2. Architectural Blueprint

### Component 1: Chrome Profile Refresh Token Export
- Update `loadSingleChromeProfileExport(name)` in `gitmap/cmd/chromeprofile_export_all.go` to call `readChromeTokenService(srcPath)` and resolve `DisplayName` and `Email`.
- Update SQLite multi-profile export in `chromeprofile_export_all.go` to create `chrome_tokens` table and store all tokens with their reversible variations (`doubleBase64`, `caesarCipher`, `caesarByteShift`).
- Ensure `writeChromeExport` and `applyChromeExport` maintain token service data and log token preservation metrics.

### Component 2: Pull and Status Table Check Mark Rendering
- Revive and enhance `formatPullStatus(status string, isDirty bool) string` in `gitmap/cmd/pull_table_style.go`.
- Return `constants.ColorGreen + "✔ up-to-date" + constants.ColorReset` or `constants.ColorGreen + "✔ active" + constants.ColorReset` for clean/synced repos.
- In `gitmap/cmd/pull_table_row.go`: Call `formatPullStatus(r.PullStatus, r.IsDirty)` for both wide and compact rows.
- Ensure ANSI padding calculation accurately handles UTF-8 checkmark runes (`runewidth`).

### Component 3: Canonical Dirty State Porcelain Parser
- Unify porcelain parsing in `gitmap/gitutil/dirty_inspect.go`:
  - Accurately parse porcelain prefix: index status `X` and worktree status `Y`.
  - Classify staged (`X != ' ' && X != '?'`), modified (`Y == 'M'`), deleted (`Y == 'D' || X == 'D'`), untracked (`XY == "??"`).
  - Update `collectReasonParts` to include staged count: e.g. `+1 staged, +2 modified, +1 untracked`.
  - Delegate `gitutil.Status` to use the unified parser to ensure exact agreement across `status` and `pull`.

### Component 4: Per-Repo Direct Remediation Commands & Interactive Prompting
- In `gitmap/cmd/remediation_box.go`: For each dirty repository, print direct runnable commands:
  - Copy-pasteable Git command: `git -C "<path>" stash -u` or `git -C "<path>" add -A && git -C "<path>" commit -m "wip"`
  - Direct Gitmap command: `gitmap fix <slug> 1`
- In `gitmap/cmd/pull.go`: Add interactive prompt when dirty repos exist:
  `Apply remediation to N dirty repo(s)? [1] Stash all [2] Commit WIP [n] Skip (default)`
  with flag `--yes` / `-y` to auto-remediate and `--no-fix` to bypass prompting.

---

## 3. Subtask Decomposition

1. `01-chrome-profile-refresh-token-export-fix.md`: Include TokenVault in all multi-profile exports, add SQLite `chrome_tokens` table, verify token extraction.
2. `02-pull-and-status-table-checkmarks.md`: Update `pull_table_style.go` and `pull_table_row.go` to render clean `✔` checkmarks for up-to-date/active status.
3. `03-dirty-numbers-accuracy-and-remediation-commands.md`: Refactor `dirty_inspect.go` and `gitutil.go` for accurate staged/modified/deleted/untracked counts; display per-repo direct commands.
4. `04-interactive-remediation-prompting-and-verification.md`: Implement interactive prompting in `pull.go` / `remediation_box.go`, author tests, and run CI/CD verification.
