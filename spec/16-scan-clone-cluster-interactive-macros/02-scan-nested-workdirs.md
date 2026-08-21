# Specification 16 - Part 2: Nested Scanning, Work Directory Grouping & Missing Repo Remediation

## 1. Nested Directory Traversal & Discovery

### 1.1 Deep Recursive Walk
Current scanner implementations can terminate shallowly or miss repositories nested several levels deep within project workspaces (e.g. `D:\work\client-projects\team-a\repo1`).
- The directory walker MUST traverse child directories recursively up to a configurable max depth (default: 16 levels).
- Standard exclusion lists MUST be respected (`.git`, `node_modules`, `vendor`, `.terraform`, `dist`, `bin`, `.cache`).
- When a `.git` directory is found, the containing folder is recorded as a repository, and the walker DOES NOT descend further into that `.git` folder, but continues scanning sibling directories.

### 1.2 Multi-Work-Directory Scans
When multiple paths or scan-roots are supplied (e.g. `gitmap scan D:\work D:\personal C:\tools`):
- Each root directory is treated as an explicit work directory scope.
- The scanner records the originating root path against each repository record (`ScanFolderRoot`).

## 2. Dual Display Modes

Users can select between two visual table representations for scan and status results:

### 2.1 Flat Table View (`--view=flat` / default for single folder)
Displays a unified column table sorted alphabetically by repository slug:

```text
  REPO                      BRANCH    STATUS     SYNC    STASH   FILES
  ────────────────────────────────────────────────────────────────────
  ai-empathy-prompt-tuner   main      ✔ clean     ─       ─       ─
  cat-my-v12                main      ✔ clean     ─       ─       ─
  gitmap-v28                main      ✔ clean     ─       ─       ─
  prompt-architect-v2       ─         ✖ missing   ─       ─       ─
```

### 2.2 Work-Directory Grouped View (`--view=grouped` / default when multi-root)
Groups repositories under distinct work directory header sections.
- The **Default Work Directory** is prominently highlighted with a `[DEFAULT]` badge.
- Nested folders under the primary work directory are formatted hierarchically.

```text
  📂 Work Directory: D:\work [DEFAULT] (18 repos)
  ────────────────────────────────────────────────────────────────────
  REPO                      BRANCH    STATUS     SYNC    STASH   FILES
  alim-status-sample-v2     main      ✔ clean     ─       ─       ─
  gitmap-v28                main      ✔ clean     ─       ─       ─
  coding-guidelines-v24     main      ✔ clean     ─       ─       ─

  📂 Work Directory: D:\personal\experiments (6 repos)
  ────────────────────────────────────────────────────────────────────
  REPO                      BRANCH    STATUS     SYNC    STASH   FILES
  prototype-agent-v1        main      ✔ clean     ─       ─       ─
  ui-prompts-cat            develop   ✔ clean     ─       ─       ─
```

## 3. Default Work Directory Tracking

### 3.1 Persistence & Rule
- The first work directory scanned (or the directory from which `gitmap scan` is initially run without arguments) is recorded in the SQLite database table `scan_folders` with `is_default = 1`.
- If a default already exists, subsequent scans update metadata (`last_scanned_at`, `repo_count`) while preserving the primary default flag unless explicitly overridden with `--set-default`.
- Commands such as `gitmap cd <name>` and `gitmap clone` fall back to the default work directory when no relative path or target directory is specified.

## 4. Missing Repository Remediation Workflow

### 4.1 Missing Status Reporting
When repositories recorded in the database or `gitmap.json` cannot be located on disk at their expected paths:
- They are flagged with status `✖ missing` (or `✖ not found`) in red/yellow.
- The summary footer explicitly counts missing repos:
  `26 repos · 24 clean · 2 missing`

### 4.2 Remediation Commands Output
At the conclusion of `gitmap scan` and `gitmap status`, if missing repos exist, the CLI displays an actionable guidance block:

```text
  ────────────────────────────────────────────────────────────────────
  ⚠  2 missing repositories detected:
     • prompt-architect-v2  (expected at D:\work\prompt-architect-v2)
     • global-ppt-v1        (expected at D:\work\global-ppt-v1)

  To resolve missing repositories:
  1. Relocate to a new folder:
     $ gitmap scan-folder update <slug> <new-path>
  2. Untrack from database:
     $ gitmap rm prompt-architect-v2 global-ppt-v1
```

### 4.3 Interactive Resolution Prompt
When run interactively in a TTY (and `--yes` is not supplied):
```text
  Would you like to resolve missing repositories now? [y/N]: 
  1) Relocate paths interactively
  2) Remove missing repos from tracking (gitmap rm)
  3) Ignore and keep in database
  Select option [1-3]: 
```
