# [V6] Plan 60: GitMap PAS Formula, Ignore Grouping Engine, CPAR Suite & Split-DB Repo Cache (Consolidated)

- **Status:** `COMPLETED`
- **Completed At:** 2026-10-01
- **Canonical Spec:** [02-spec/21-app/197-gitmap-pas-fix-and-repo-cache-commands.md](../../../02-spec/21-app/197-gitmap-pas-fix-and-repo-cache-commands.md) & [02-spec/21-app/181-gitmap-ignore-and-cache-engine/01-overview.md](../../../02-spec/21-app/181-gitmap-ignore-and-cache-engine/01-overview.md)
- **Execution Lifecycle:** 2-Wave Parallel Multi-Agent Subagent Execution (`A = 2`, `H = 2`) via `invoke_subagent` (`TypeName: "self"`).
- **Total Loops / Waves:** 2 Waves (Wave 1: Worker 01 & Worker 02; Wave 2: Worker 03 & Worker 04).

---

## 1. Executive Summary & Verification Ledger

All 7 subtasks have been fully implemented, verified with file-level quality linters, and consolidated:

1. **Subtask 01 (`01-pull-all-optimization.md`): Decoupled Asynchronous Pull & PAS Formula Engine**
   - Prioritizes git pull operations first across repositories.
   - Non-blocking asynchronous ignore scanning executes via `startAsyncIgnoreScan` (worker pool: 5 repos per worker).
   - Results collected and presented at the end of pull batch execution with clean per-repo summaries.
   - Interactive prompt supports resolving all at once (`[y/all]`), one-by-one (`[s/single]`), or skip (`[n]`).
   - Non-interactive `-y` / auto-yes fully supported.

2. **Subtask 02 (`02-fix-ignore-all.md`): GitMap Fix Ignore All (Local & SSH Fleet)**
   - `gitmap fix-ignore-all` (`fia`) & `gitmap fix-ignores-all-ssh` (`fias`).
   - Scans repositories for duplicate lines in `.gitignore` and tracked ignored files.
   - Deduplicates patterns preserving comments and blank line structure.
   - Git index verification via `filterTrackedFilesInIndex` (`git ls-files --error-unmatch`) ensures untracked files are never falsely untracked or prompted.
   - Registered in `TasksSplitDB` with status transitions (`running`, `completed`, `failed`) and errors logged to GitMap Errors DB.
   - `fias` enforces the canonical **Gitmap PAS Formula** (2 workers, 2 async operations max).

3. **Subtask 03 (`03-commit-push-all.md`): GitMap Commit-Push All Repos Suite (CPAR)**
   - `gitmap commit-push-all-repos` (`cpar`).
   - Discovers all dirty repositories across registered workspaces.
   - `--review` (`-r`) displays structured change tree (modified, untracked, staged counts, individual files).
   - Interactive review flow supports committing as feature, bug, single-repo sessions, or chore.
   - `--commit-only` (`-co`) mode stages and commits locally without pushing upstream.
   - Registered in `TasksSplitDB` with error logging via `store.LogInternalError("CPAR", ...)`.

4. **Subtask 04 (`04-ignore-engine.md`): GitMap Ignore Engine & Hierarchical Groups**
   - `gitmap ignore` (`ig`) suite:
     - `add <pattern>`: adds to default group in SQLite DB and `.gitmap/ignores/default.txt`. Duplicate check prints `Already added. Don't need to add.`
     - `ls` / `list`: lists groups, pattern counts, default tags, chained groups, and repo bindings.
     - `add-group <name>`: creates named ignore group.
     - `remove-group` / `rm-grp <name>`: deletes named group.
     - `set-default-group <name>`: sets default group.
     - `add-grp-to-default` (`agtd`) `<name>`: chains named group after default group.
     - `connect-group-with-repo` (`cgwp`) `<group> <repo> [--add-with-default | --awd]`: binds group to repo.
     - `apply [path]`: applies active rules to `.gitignore`, cleanly deduplicating patterns.
     - `export` / `import`: JSON configuration backup and restore.
     - `help`: comprehensive terminal help text with syntax and examples.
     - Default patterns guarantee `.gitmap/` and `.gitmap/backup/` persistence.

5. **Subtask 05 (`05-cache-engine.md`): GitMap Split-DB Repository Cache & Search Engine**
   - Split-DB cache architecture:
     - Root `sql.db` (`.gitmap/cache/repos/<slug>/sql.db`) storing `RepoMetadata`, `FolderTree`, and `Files` index with `mtime` (zero slow hashing).
     - Top-level folder databases: separate SQLite DB `<folder-slug>.db` per folder (no subfolder nesting); root files in `root.db`.
     - Exclusion gates: excludes `.git/`, `node_modules/`, `.vscode/`, `.idea/`, `.gitmap/`, `vendor/`, files > 200KB, binary files, large JSON (> 150KB).
     - `is_keep` flag (`--keep`, `-k`) supported.
   - `search "<text>" ["<glob>"]`: fast text search with 10 context lines and limit 20.
   - `search "<text>" -file-pattern (fp) "a*.md", "b*.md"`: glob path filtering.
   - `search-multi "t1", "t2"`: multi-term search.
   - `search-multi-grep "<regex>"`: SQLite compiled regex search.
   - `recache` / `reconcile` / `sync`: checks `mtime` against filesystem, updates outdated files, prunes deleted files.
   - `ls` & `rm`: cache management commands.

6. **Subtask 06 (`06-see-commands.md`): GitMap Unified See Inspection & Remote Telemetry**
   - `gitmap see commit pending` / `gitmap c commit pending` / `gitmap c cp`: display dirty repos.
   - `gitmap see git-ignore issues` / `gitmap see ignore issues` / `gitmap see ig issues` / `gitmap c ig issues`: inspect ignore issues.
   - `gitmap see errors` / `gitmap c errors`: inspect local error journal.
   - `gitmap see history` / `gitmap c history`: inspect task execution history.
   - `gitmap see errors ssh` / `gitmap ses`: query remote fleet node errors following Gitmap PAS Formula.
   - `gitmap history ssh` / `gitmap nodes history` / `gitmap nodes histories`: view remote task history.
   - `gitmap repo-manage ui`: launch interactive terminal dashboard.
   - Root aliases `c`, `ses`, `nodes-history`, `repo-manage-ui` registered in `cli/cmd/rootcore.go`.

7. **Subtask 07 (`07-fix-screenshot-bug.md`): Pull Engine Remediation & Screenshot Bug Fix**
   - In `cli/cmdpull/pull_remediation_hint.go`:
     - Replaced raw git advice (`git -C ...`, `git pull --rebase`, `git config pull.rebase`) with native GitMap actionable commands:
       - Dirty: `gitmap fix <repo>`, `gitmap cpar`, or `gitmap stash`.
       - Merge conflicts: `gitmap fix <repo>` or `gitmap stash`.
       - Diverged branches: `gitmap pull <repo>`.
       - Missing repositories: `gitmap clone <repo>`.
       - Auth / general: `gitmap status <repo>` or `gitmap fix <repo>`.
     - Zero raw git command leakage.
   - In `cli/gitignoreagm/cli_prompt.go`:
     - Verified git index tracking status via `git ls-files --error-unmatch` so untracked files are never falsely prompted.

---

## 2. Modified Files Audit

- [cli/cmd/rootcore.go](../../../cli/cmd/rootcore.go)
- [cli/cmdcache/cache_cli.go](../../../cli/cmdcache/cache_cli.go)
- [cli/cmdcache/cache_create.go](../../../cli/cmdcache/cache_create.go)
- [cli/cmdcache/cache_search.go](../../../cli/cmdcache/cache_search.go)
- [cli/cmdcache/cache_types.go](../../../cli/cmdcache/cache_types.go)
- [cli/cmdcpar/cpar.go](../../../cli/cmdcpar/cpar.go)
- [cli/cmdcpar/types.go](../../../cli/cmdcpar/types.go)
- [cli/cmdignore/fias.go](../../../cli/cmdignore/fias.go)
- [cli/cmdignore/fix_ignore.go](../../../cli/cmdignore/fix_ignore.go)
- [cli/cmdignore/ignore_cli.go](../../../cli/cmdignore/ignore_cli.go)
- [cli/cmdignore/ignore_groups.go](../../../cli/cmdignore/ignore_groups.go)
- [cli/cmdpull/pull.go](../../../cli/cmdpull/pull.go)
- [cli/cmdpull/pull_efficient.go](../../../cli/cmdpull/pull_efficient.go)
- [cli/cmdpull/pull_remediation_hint.go](../../../cli/cmdpull/pull_remediation_hint.go)
- [cli/cmdsee/see.go](../../../cli/cmdsee/see.go)
- [cli/gitignoreagm/cli_prompt.go](../../../cli/gitignoreagm/cli_prompt.go)
- [cli/store/split_db_cache.go](../../../cli/store/split_db_cache.go)
- [cli/store/split_db_ignore.go](../../../cli/store/split_db_ignore.go)
