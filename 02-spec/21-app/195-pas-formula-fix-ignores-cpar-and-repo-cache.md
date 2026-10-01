# 195 — GitMap PAS Formula, Ignore Management Suite, CPAR, and Repository Split-DB Cache Engine

## Status: Approved
- **Spec ID:** 195
- **Subsystems:** `cli/cmdpull`, `cli/cmdignore`, `cli/cmdcpar`, `cli/cmdsee`, `cli/cmdcache`, `cli/cmdssh`, `cli/cloner`
- **Date:** 2026-10-01

---

## 1. Architectural Overview & The GitMap PAS Formula

### 1.1 The GitMap PAS Formula
The **GitMap PAS Formula** (Pull-All-SSH Standard) is the canonical orchestration pattern in GitMap for coordinating hybrid operations across local workstations and distributed remote fleet nodes:

1. **Local Host Execution (Direct In-Process):**
   - The current machine (`127.0.0.1` / localhost) executes directly in-process with real-time UI/streaming output.
   - Local operations are never wrapped in SSH loops or queued behind asynchronous fleet worker tasks.
   - Standard stdout and stderr are preserved without escape leakages.

2. **Remote Fleet Nodes (Bounded Async Dispatch):**
   - Remote fleet nodes registered in the cluster database (`gitmap.db` / `installation.db`) are probed for liveness using resilient health checks (3000ms timeout with 1 retry, caching offline states for 5s and online states for 45s).
   - Online nodes receive bounded asynchronous tasks via SSH JSON IPC with low worker concurrency (default `A=2`, `H=2`).
   - Offline nodes are skipped immediately with structured reporting.

3. **Structured Telemetry & Audit:**
   - Every operation logs to the local Task Audit DB (`gitmap.db` / `TaskHistory` table) and Errors DB.
   - Any remote command can be audited via corresponding `see` and `ses` (`see errors ssh`) commands.

```
                      ┌───────────────────────────────────────────────┐
                      │            User Command (e.g. pas)            │
                      └───────────────────────┬───────────────────────┘
                                              │
                 ┌────────────────────────────┴────────────────────────────┐
                 ▼                                                         ▼
      ┌──────────────────────────────┐                         ┌──────────────────────┐
      │   Local Host (In-Process)    │                         │ Remote Fleet Manager │
      │   Direct Streaming Output    │                         │ (Bounded A=2, H=2)   │
      └──────────────────────────────┘                         └──────────┬───────────┘
                                                                          │
                                                      ┌───────────────────┴───────────────────┐
                                                      ▼                                       ▼
                                           ┌─────────────────────┐                 ┌─────────────────────┐
                                           │  Node w1 (Online)   │                 │  Node w3 (Online)   │
                                           │  Async SSH JSON Task│                 │  Async SSH JSON Task│
                                           └─────────────────────┘                 └─────────────────────┘
```

---

## 2. Pull-First Workflow & Windows Credential Safety

### 2.1 Pull-First Mandate
- Commands `gitmap pa`, `gitmap pull-all`, `gitmap pas`, `gitmap pull-all-ssh` MUST execute pulls immediately.
- Upfront interactive prompts (such as `.gitignore` or resume task questions) are strictly decoupled from the start of the pull.
- Ignore scans run asynchronously or report as an end-of-pull summary with an optional non-blocking resolution.

### 2.2 Windows Credential Store Safe Subprocess Environment
- All Git subprocess executions use an OS-aware environment constructed via `gitutil.BuildSafeGitEnv`:
  - `GCM_NO_PERSIST=1`
  - `GCM_INTERACTIVE=never`
  - `GIT_TERMINAL_PROMPT=0`
  - On Windows: Do NOT set `GCM_CREDENTIAL_STORE=cache`. Git for Windows lacks UNIX domain socket support for the credential cache daemon. Setting `cache` causes immediate fatal exit.
  - On non-Windows (Linux / macOS): `GCM_CREDENTIAL_STORE=cache` is permitted.

---

## 3. Ignore Management & Grouping Architecture

### 3.1 Command Suite
- **Fix All Repositories:**
  - `gitmap fix ignore all [-y]` (alias `fia`): Scans and remediates `.gitignore` issues across all repositories.
  - `gitmap fix-ignore-all [-y]`: Hyphenated form of fix ignore all.
  - `gitmap fix ignores all ssh [-y]`: Runs fix ignore across local host and all SSH fleet nodes.
  - `gitmap fix-ignores-all-ssh (fias) [-y]`: Canonical PAS Formula command for fleet-wide ignore remediation.

- **Ignore Subcommand (`gitmap ignore` / `gitmap ig`):**
  - `add <entry>`: Adds an ignore entry to `.gitignore`.
  - `scan`: Scans repositories for ignore issues, duplicate lines, and untracked artifacts.
  - `scan-ssh` (`ss`): Scans ignore issues across remote fleet nodes (PAS Formula).
  - `remove <entry>`: Removes an ignore entry.
  - `edit`: Opens `.gitignore` in the default editor.
  - `action`: Interactive triage action for ignore issues.
  - `ls`: Lists active ignore patterns.
  - `help`: Displays comprehensive ignore command help.
  - `ui` / `app`: Opens web/terminal UI for ignore management.
  - `add-group <name>`: Creates a named ignore group.
  - `remove-group` (`rm-grp`) <name>: Deletes a named ignore group.
  - `set-default-group <name>`: Sets the default ignore group applied to all repositories.
  - `add-grp-to-default` (`agtd`) <name>: Adds an ignore group to the default list.
  - `apply <group-name>`: Applies group patterns to repositories.
  - `connect-group-with-repo` (`cgwp`) <group-name> <path1,alias> [--add-with-default (`awd`)]: Connects a group to specific repositories.
  - `export` / `import`: Exports and imports ignore configuration definitions.

### 3.2 Idempotent GitIgnore Sanitizer
- Automatically includes `.gitmap/backup/`, `.antigravity_resume_task.json`, `antigravity-resume_task.json` in default ignores.
- Deduplicates lines in `.gitignore` preserving section comments (`# Section`) and existing structure.
- Normalizes trailing newlines and removes accidental empty duplicate blocks.

---

## 4. Commit & Push All Repositories (`cpar`)

### 4.1 Commands & Flags
- `gitmap commit-push-all-repos` (`cpar`) `[-y]`: Stages changes (`git add -A`), creates commit, and pushes across all dirty repos.
- `gitmap commit-push-all-repos (cpar) --review` (`-r`): Displays pending commits across dirty repos for review before committing and pushing.
- `gitmap commit-push-all-repos (cpar) --review --commit-only` (`-co`): Displays pending changes, reviews, and commits without pushing.

---

## 5. Observability & Inspection Suite (`see`)

### 5.1 Commands
- `gitmap see commit pending`: Lists all repositories with uncommitted or unstaged changes.
- `gitmap see git-ignore` / `ignore` / `ig issues`: Lists repositories with ignore issues or duplicate patterns.
- `gitmap see errors`: Alias for `gitmap errors`.
- `gitmap see history`: Alias for `gitmap history`.
- `gitmap see errors ssh` (`ses`): Queries and aggregates error logs from remote fleet nodes over SSH following the PAS Formula.
- `gitmap history ssh`: Aggregated fleet execution history.
- `gitmap nodes histories` / `history`: Fleet nodes execution histories.
- `gitmap repo-manage ui`: Interactive repository management UI.

---

## 6. Repository Split-DB Cache Engine (`gitmap cache`)

### 6.1 Commands
- `gitmap cache create .`: Indexes current repository into split SQLite databases.
- `gitmap cache create <relpath1>,<abspath2>`: Indexes specified paths into split cache databases.
- `gitmap cache create "a.json", "b.json"`: Indexes specific files.
- `gitmap cache ls`, `add`, `create`, `remove` (`rm`), `help`.
- `gitmap cache search "text search" "*.md" [--lines 10] [--limit 20]`: Full-text search across indexed files with optional glob, line context, and match limit.
- `gitmap cache search "text search" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]`: File-pattern filtered search.
- `gitmap cache search-multi "text search", "multi *" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]`: Multi-term search across indexed cache.
- `gitmap cache search-multi-grep "regex search", "multi *" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]`: Regex search across cached content.
- `gitmap cache recache` / `reconcile` / `sync`: Reconciles cached file hashes and updates modified files.

### 6.2 Split-DB Architecture
- **Root Cache DB** (`gitmap.cache.db`):
  ```sql
  CREATE TABLE IF NOT EXISTS Files (
      FileId INTEGER PRIMARY KEY AUTOINCREMENT,
      RepoUrl TEXT,
      RelativePath TEXT UNIQUE,
      AbsolutePath TEXT,
      FileSize INTEGER,
      ModifiedTime INTEGER,
      FolderSlug TEXT,
      IsKeep INTEGER DEFAULT 0
  );
  CREATE INDEX IF NOT EXISTS Idx_Files_Slug ON Files(FolderSlug);
  CREATE INDEX IF NOT EXISTS Idx_Files_Rel ON Files(RelativePath);
  ```
  - File size threshold: files <= 200KB indexed; binaries, large images, and >200KB files skipped.
  - `IsKeep` boolean flag to mark cached assets for persistence across refreshes.

- **Folder Slug DBs** (`.gitmap/cache/<slug>.db`):
  ```sql
  CREATE TABLE IF NOT EXISTS Lines (
      LineId INTEGER PRIMARY KEY AUTOINCREMENT,
      RelativePath TEXT,
      LineNumber INTEGER,
      Content TEXT
  );
  CREATE INDEX IF NOT EXISTS Idx_Lines_Path ON Lines(RelativePath);
  ```
  - Partitioned per top-level folder slug (e.g. `cli.db`, `docs.db`) to enable zero-latency parallel search queries across partitioned SQLite stores.
