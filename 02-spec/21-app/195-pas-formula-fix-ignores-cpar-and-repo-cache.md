# 195 — GitMap PAS Formula, Ignore Management Suite, CPAR, and Repository Split-DB Cache Engine

## Status: Approved
- **Spec ID:** 195
- **Subsystems:** `cli/cmdpull`, `cli/cmdignore`, `cli/cmdcpar`, `cli/cmdsee`, `cli/cmdcache`, `cli/cmdssh`, `cli/cloner`
- **Date:** 2026-10-01

---

## 1. Architectural Overview & The GitMap PAS Formula

### 1.1 The GitMap PAS Formula
The **GitMap PAS Formula** defines the canonical orchestration standard for hybrid local/fleet execution across GitMap:
1. **Local Host Execution (Direct In-Process)**:
   - The current machine (`127.0.0.1` / localhost) executes directly in-process with real-time UI/streaming.
   - It is never wrapped in SSH or enqueued through asynchronous fleet workers.
2. **Remote Fleet Nodes (Bounded Async Dispatch)**:
   - Connected remote fleet nodes are queried for liveness (with differentiated cache TTLs and resilient probing).
   - Online nodes receive bounded asynchronous tasks via SSH JSON IPC with low worker concurrency (default `A=2`, `H=2`).
   - Offline nodes are skipped immediately with clear status reporting.
3. **Structured Telemetry & Audit**:
   - Every operation logs to the local Task Audit DB and Errors DB.
   - Any remote command can be audited via corresponding `see` and `ses` (see errors ssh) commands.

```
                      ┌─────────────────────────────────┐
                      │    User Command (e.g. pas)      │
                      └────────────────┬────────────────┘
                                       │
                ┌──────────────────────┴──────────────────────┐
                ▼                                             ▼
     ┌──────────────────────┐                     ┌──────────────────────┐
     │  Local Node Direct   │                     │ Remote Fleet Manager │
     │  In-Process Run      │                     │ (Bounded A=2, H=2)   │
     └──────────────────────┘                     └──────────┬───────────┘
                                                             │
                                              ┌──────────────┴──────────────┐
                                              ▼                             ▼
                                     ┌─────────────────┐           ┌─────────────────┐
                                     │ Node w1 (Online)│           │ Node w3 (Online)│
                                     │ Async SSH Task  │           │ Async SSH Task  │
                                     └─────────────────┘           └─────────────────┘
```

---

## 2. Pull-First Workflow & Windows Credential Safety

### 2.1 Pull-First Mandate
- `gitmap pa`, `gitmap pull-all`, `gitmap pas`, `gitmap pull-all-ssh` MUST execute pulls immediately.
- Pre-flight interactive prompts (such as `.gitignore` or resume task file questions) are strictly decoupled from the start of the pull.
- Ignore scans run asynchronously or report as an end-of-pull summary with an optional non-blocking resolution.

### 2.2 Windows Credential Store Safe Subprocess Environment
- All Git subprocess executions use an OS-aware environment:
  - `GCM_NO_PERSIST=1`
  - `GCM_INTERACTIVE=never`
  - `GIT_TERMINAL_PROMPT=0`
  - On Windows: Do NOT set `GCM_CREDENTIAL_STORE=cache`.
  - On non-Windows: `GCM_CREDENTIAL_STORE=cache` is permitted.

---

## 3. Ignore Management & Grouping Architecture

### 3.1 Command Suite
- **Fix All**:
  - `gitmap fix ignore all [-y]` (`fia`)
  - `gitmap fix-ignore-all [-y]`
  - `gitmap fix ignores all ssh [-y]`
  - `gitmap fix-ignores-all-ssh (fias) [-y]` (follows GitMap PAS Formula)
- **Ignore Subcommand (`gitmap ignore` / `gitmap ig`)**:
  - `add <entry>`: Add ignore entry.
  - `scan`: Scan repos for ignore issues, duplicates, and untracked artifacts.
  - `scan-ssh` (`ss`): Scan ignore issues across remote fleet nodes (PAS Formula).
  - `remove <entry>`: Remove ignore entry.
  - `edit`: Open ignore configuration in editor.
  - `action`: Perform ignore actions.
  - `ls`: List active ignore patterns.
  - `help`: Display ignore help.
  - `ui` / `app`: Open web/terminal UI for ignore management.
  - `add-group <name>`: Create named ignore group.
  - `remove-group` (`rm-grp`) <name>: Delete named ignore group.
  - `set-default-group <name>`: Set default group applied to all repos.
  - `add-grp-to-default` (`agtd`) <name>: Add group to default list.
  - `apply <group-name>`: Apply group patterns to repositories.
  - `connect-group-with-repo` (`cgwp`) <group-name> <path1,alias> [--add-with-default (`awd`)]: Connect group to repositories.
  - `export` / `import`: Export/import ignore configurations.

### 3.2 Idempotent GitIgnore Sanitizer
- Automatically includes `.gitmap/`, `.gitmap/backup/`, and resume task artifacts in default ignores.
- Deduplicates lines in `.gitignore` files preserving section comments and original structure.

---

## 4. Commit & Push All Repositories (`cpar`)

### 4.1 Commands & Flags
- `gitmap commit-push-all-repos` (`cpar`) `[-y]`: Stages all changes, commits with standard message, and pushes across all dirty repos.
- `gitmap commit-push-all-repos (cpar) --review` (`-r`): Displays pending commits across dirty repos for user review before committing and pushing.
- `gitmap commit-push-all-repos (cpar) --review --commit-only` (`-co`): Displays pending changes, reviews, commits without pushing.

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
- `gitmap cache create .`
- `gitmap cache create <relpath1>,<abspath2>`
- `gitmap cache create "a.json", "b.json"`
- `gitmap cache ls`, `add`, `create`, `remove` (`rm`), `help`
- `gitmap cache search "text search" "*.md" [--lines 10] [--limit 20]`
- `gitmap cache search "text search" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]`
- `gitmap cache search-multi "text search", "multi *" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]`
- `gitmap cache search-multi-grep "regex search", "multi *" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]`
- `gitmap cache recache` / `reconcile` / `sync`

### 6.2 Split-DB Architecture
- **Root Cache DB** (`gitmap.cache.db` / split index):
  - Repository URL, relative folder hierarchy, file metadata, modification timestamps.
  - File size threshold: files <= 200KB indexed; binaries, large images, and >200KB files skipped.
  - `is_keep` boolean flag to mark cached assets for persistence across refreshes.
- **Folder Slug DBs**:
  - Partitioned per top-level folder slug (e.g. `cli.db`, `docs.db`) to enable zero-latency parallel search queries.
