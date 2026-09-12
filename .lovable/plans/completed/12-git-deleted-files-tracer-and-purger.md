# 12-git-deleted-files-tracer-and-purger.md

**Title:** Git Historical Deleted Files Tracer, Restorer & Deep Purger
**Status:** COMPLETED
**Parent Task:** User Request (Git History Recovery & Deep Purging Engine)
**Target Script:** `03-ai-scripts/33-git-history-tracer-and-purger.py`
**Shared Engine:** `03-ai-scripts/02-shared-engine.py`

---

## 1. Executive Summary & Problem Statement

As repositories evolve and undergo major file consolidations (such as reducing hundreds of micro-plans and subtasks into high-density milestones), historical commits retain copies of every deleted file in Git's object database (`.git/objects/`).

Users face two complementary challenges:
1. **Tracing & Recovery:** Identifying exactly what files were removed across historical commits and recovering specific lost files back to the working directory without having to manually search through `git log` commit hashes.
2. **Deep History Purge:** Permanently expelling sensitive, obsolete, or bloated deleted files from the entire Git database (including all commits, branches, releases, and tags) to reduce packfile size and ensure eliminated files do not linger in repository history.

This specification designs an optimized, interactive Python CLI tool (`03-ai-scripts/33-git-history-tracer-and-purger.py`) backed by `02-shared-engine.py` and `git-filter-repo` that provides pre-flight indexing, selective exclusion, granular file restoration, and full Git history purging.

---

## 2. Core Functional Requirements

### Requirement 1: Historical Deletion Tracer & Pre-flight Indexing
- Traverse Git commit history using `git log --diff-filter=D --name-only` plumbing to catalog all files that were deleted in history and do not exist in current `HEAD`.
- Extract metadata for each deleted file:
  - Relative repository path.
  - Deletion commit hash (`%H`).
  - Commit author (`%an`).
  - Commit date (`%ad`).
  - Commit summary subject (`%s`).
  - Size of blob prior to deletion (`git cat-file -s <commit>~1:<path>`).
- Render a structured, numbered pre-flight candidate table:
  ```text
  [Index]  Deleted File Path                             Commit Date   Author             Commit Subject
  ---------------------------------------------------------------------------------------------------------
  [ 1]     .lovable/plans/subtasks/01-task.md           2026-09-10    MD ALIM UL KARIM   chore: consolidate...
  [ 2]     .lovable/plans/subtasks/02-task.md           2026-09-10    MD ALIM UL KARIM   chore: consolidate...
  [ 3]     spec/old-feature/01-spec.md                  2026-08-20    MD ALIM UL KARIM   refactor: remove...
  ```

### Requirement 2: Preset Shortcuts & Flexible Filtering
Provide built-in shortcuts for high-frequency operations alongside generic CLI flags:
- `[target]`: Positional argument for arbitrary directory or file path from root (e.g. `spec/21-app/25-app-spec-audit`, `spec/19-main-worker-service/audit`, `.lovable`, or `.` for entire repository).
- `--spec-25-audit`: Exact folder shortcut targeting `spec/21-app/25-app-spec-audit` (Folder 25 app spec audit documents).
- `--spec-audit`: Exact folder shortcut targeting `spec/19-main-worker-service/audit` (audit documents removed after review).
- `--audit`: Repository-wide shortcut targeting all deleted audit files (`*audit*`).
- `--lovable-subtasks`: Scopes tracing to `.lovable/plans/subtasks/`.
- `--lovable-audits`: Scopes tracing to `.lovable/audits/`.
- `--lovable`: Scopes tracing to all deleted files under `.lovable/`.
- `--lovable-md`: Scopes tracing to all deleted `.md` files under `.lovable/`.
- `--spec`: Scopes tracing to all deleted files under `spec/`.
- `--spec-md`: Scopes tracing to all deleted `.md` files under `spec/`.
- `--path <folder>`: Filter by arbitrary folder path prefix.
- `--ext <extension>`: Filter by file extension (e.g. `.md`, `.go`, `.json`, `.sql`).
- `--pattern <glob/regex>`: Safe filter by path glob or regex substring (protected against invalid regex compile crashes).

### Requirement 3: Interactive & Flag-Based Selective Exclusion
- Pre-flight candidate selection allows users to inspect the numbered list and exclude files:
  - CLI flag: `--exclude 1,3,5-8,12`
  - Interactive mode: prompts user `Enter numbers or ranges to exclude (e.g. 1, 4-7) [or press Enter to keep all]:`
  - Pattern exclusion: `--exclude-pattern "<glob_or_regex>"`
- Re-renders the filtered list showing remaining target candidates and count.

### Requirement 4: File Restoration / Recovery Engine (`--restore`)
- Recover any or all selected deleted files back into the active workspace.
- Retrieves the file contents from the commit immediately preceding deletion: `git show <commit>~1:<path>`.
- Writes the recovered file to disk at its canonical path (or `--restore-to <dir>` staging directory).
- Re-creates missing parent directories automatically.

### Requirement 5: Workspace Recycle Bin Deletion Engine (`--delete`)
- Deletes active target files from the workspace using the native OS **Recycle Bin** (`send2trash` / Windows Shell `SHFileOperationW`).
- Prior to recycling, copies each file to a timestamped OS temp directory (`%TEMP%` / `/tmp`).
- Prints exact temp backup path and recovery command upon completion.

### Requirement 6: Deep Git History Purge Engine (`--purge`)
- Permanently rewrites Git history to completely eradicate the target files from all commits, branches, and tags.
- Uses `git filter-repo --paths-from-file <tmp_paths> --invert-paths --force`.
- **Mandatory Safety Protocol:**
  1. Verifies working tree is clean (`git status --porcelain`).
  2. Creates timestamped OS temp directory backup of all historical blobs and workspace files.
  3. Saves the origin remote URL (`git remote get-url origin`).
  4. Creates a timestamped safety backup branch: `backup/history-purge-YYYYMMDD-HHMMSS`.
  5. Executes history rewrite.
  6. Restores origin remote URL.
  7. Reports local OS temp backup path and dual rollback recipes (Git reset & file copy).

---

## 3. Architecture & Technical Invariants

### 1. Reusable Tooling & Shared Engine
- Imports from `03-ai-scripts/02-shared-engine.py`:
  - `ExitCodeType`, `RegexPatternType`
  - `normalize_rel_path`, `read_file_lf`, `write_file_lf`
  - `run_command_safe`, `format_comma_separated`
- Fast caching: reuses `17-fast-file-reader.py` / caching structures for blob resolution.

### 2. Coding Guidelines Compliance
- Every function <= 15 lines.
- Affirmative booleans (`is_*`, `has_*`) with implicit evaluation (`if is_valid:`).
- Zero nested `if` statements (guard-clause early return pattern).
- Zero swallowed errors: every failure logs context and returns standard `ExitCodeType`.
- Strictly relative Git paths (zero absolute paths, zero `file:///` URIs).
- Unix LF line endings, UTF-8 (no BOM), single terminating newline at EOF.

---

## 4. CLI Interface & Help Synopsis

```text
usage: 33-git-history-tracer-and-purger.py [-h] [--path PATH] [--ext EXT]
                                          [--pattern PATTERN] [--lovable-md]
                                          [--spec-md] [--list] [--restore]
                                          [--restore-to DIR] [--purge]
                                          [--exclude EXCLUDE]
                                          [--exclude-pattern PATTERN]
                                          [-y] [--backup-only]

Git Historical Deleted Files Tracer, Restorer & Deep Purger.

actions:
  --list, -l            List and preview all historically deleted files (default).
  --restore, -r         Restore selected deleted files to the filesystem.
  --purge, -p           Deeply purge selected files from all Git commits, tags, and branches.

presets:
  --lovable-md          Target all deleted .md files under .lovable/
  --spec-md             Target all deleted .md files under spec/

filters:
  --path PATH           Filter by directory path prefix.
  --ext EXT             Filter by file extension (e.g. .md, .go, .sql).
  --pattern PATTERN     Filter by glob/regex pattern.
  --exclude EXCLUDE     Comma-separated indices or ranges to exclude (e.g. 1,3,5-8).
  --exclude-pattern PAT Exclude files matching regex/glob.

options:
  --restore-to DIR      Directory to restore files into (default: original paths).
  -y, --confirm         Bypass interactive confirmation for purge operations.
  --backup-only         Create timestamped backup branch without modifying history.
```

---

## 5. Implementation Steps (Phased) & Outcomes

- [x] **Step 1:** Create `03-ai-scripts/33-git-history-tracer-and-purger.py` with CLI argument parser, preset handlers, and shared engine integration.
- [x] **Step 2:** Implement `trace_deleted_files()` to query git log diff-filter plumbing, parse commit metadata, and filter current HEAD tracked files.
- [x] **Step 3:** Implement pre-flight table renderer and interactive/flag-based exclusion parser (`parse_exclusion_indices()`).
- [x] **Step 4:** Implement `restore_selected_files()` using `git show <commit>~1:<path>` with directory creation and recovery logging.
- [x] **Step 5:** Implement `purge_selected_files()` with safety backup branch, `git filter-repo --paths-from-file`, remote URL preservation, and confirmation gates.
- [x] **Step 6:** Register in `03-ai-scripts/01-index.md` and test across `--list`, `--restore`, `--exclude`, `--lovable-md`, and `--spec-md`.
- [x] **Step 7:** Verify linters and CI quality gates.
