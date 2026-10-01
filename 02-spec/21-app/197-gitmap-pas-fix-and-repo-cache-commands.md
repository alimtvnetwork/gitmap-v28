# 197. Gitmap PAS Fix, Ignore Grouping, CPAR, and Split-DB Repo Cache

## Overview
This specification consolidates the feature requests from the multi-part prompt provided in `02-spec/21-app/pas-fix-parts/`. It covers sweeping enhancements to GitMap's repository management capabilities, async task handling, ignoring logic, and a new high-performance SQLite-based repository cache.

## 1. PAS Formula & Async Execution
- **`gitmap pull all (pa)`**: Modifies behavior to prioritize the actual pull operation over `gitignore` checks. Ignore checks should run asynchronously (e.g., 5 projects per worker) to find duplicates or tracked ignored files. 
- **`gitmap pull all ssh (pas)`**: Applies the "Gitmap PAS formula". Runs `pa` on the current machine and delegates to SSH nodes using GitMap's task servers. SSH delegation must be highly efficient: 1 or 2 workers handling 2 async operations maximum to preserve delicate SSH resource constraints.
- All actions enqueue tasks on the task server and record results in the history/errors DB.

## 2. Ignore System (Groups & Commands)
- **`gitmap ignore (ig)`**: Subcommands:
  - `add <string>`: Adds to the default group.
  - `ls`: Lists groups, count of ignores, and application rules.
  - `add-group / remove-group (rm-grp)`: Group management.
  - `set-default-group / add-grp-to-default (agtd)`: Sets default execution groups.
  - `connect-group-with-repo (cgwp)`: Maps a group to a repo alias/path.
  - `apply`: Applies ignore rules to the current folder/repo.
  - `export / import`: JSON-based backup/restore of settings.
- Default ignores must include `.gitmap/` and `.gitmap/backup/`.
- **Validation**: Reject duplicates. Only flag tracked files if they are in the current commit state.

## 3. Global Operations (Fix, Commit-Push, See)
- **`gitmap fix ignore all [-y] (fia)`**: Scans all repos for ignore inconsistencies, duplicates, or tracked ignored files. Prompts to fix all (`-y`) or sequentially step through fixes.
- **`gitmap fix ignores all ssh [-y] (fias)`**: SSH fleet variant following the PAS formula.
- **`gitmap commit-push-all-repos (cpar) [-y] [--review|-r] [--commit-only|-co]`**: Summarizes repos with changes, prompts to commit as features/bugs, enqueues as tasks, and executes.
- **`gitmap see (c)`**: View pending issues across the fleet.
  - `c commit pending`: Changes not committed.
  - `c gitignore issues (ig)`: Repo ignore inconsistencies.
  - `c errors (ses for SSH)`: Display GitMap errors DB.

## 4. Split-DB Repo Cache System
- **`gitmap cache create <path>`**: Builds a split-DB architecture close to the CLI root.
  - `sql.db` (Root): Contains repo URL, root files, folder paths, and last modified times (no hashes). Excludes files > 200KB, images, and binaries.
  - `[slug].db` (Folder): Each subfolder gets its own SQLite DB containing its exact relative path and files.
  - Avoids nesting DBs for sub-subfolders.
- **`gitmap cache search / multi-search / search-multi-grep`**: Performs regex/text searches directly against SQLite for massive performance gains. Displays 10 lines of context by default (limit 20 matches).
  - `gitmap cache search "text search" "*.md" [--lines 10] [--limit 20]`
  - `gitmap cache search "text search" -file-pattern (fp) "a*.md", "b*.md" [--lines 10] [--limit 20]`
  - `gitmap cache search-multi "text search", "multi *" -file-pattern (fp) "a*.md" [--lines 10] [--limit 20]`
  - `gitmap cache search-multi-grep "regex search", "multi *" -file-pattern (fp) "a*.md" [--lines 10] [--limit 20]`
  - `gitmap cache recache/reconcile/sync`
- **History Commands**:
  - `gitmap history ssh`
  - `gitmap nodes histories/history`
- **Auto-Reconciliation**: If search finds outdated `last_modified` times, async workers seamlessly update the SQLite DB from the filesystem.

## 5. UI and Help
- `gitmap repo-manage ui`: Visual management of cloned repos, execution, and discovery.
- All commands must have extensive help text and documentation reflecting these workflows.

## Execution Sequence (Parent Task)
This spec is to be executed in waves using `A=2`, `H=2` disjoint file boxes, strictly following the N-step orchestrator loop to ensure the code remains DRY, elegant, and system-design focused.
