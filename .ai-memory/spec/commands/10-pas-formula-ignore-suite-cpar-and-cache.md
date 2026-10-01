# Command Specification: GitMap PAS Formula, Ignore Management Suite, CPAR & Split-DB Cache Engine

Spec Reference: [02-spec/21-app/195-pas-formula-fix-ignores-cpar-and-repo-cache.md](../../../02-spec/21-app/195-pas-formula-fix-ignores-cpar-and-repo-cache.md)

## 1. The GitMap PAS Formula

The **GitMap PAS Formula** defines the canonical orchestration standard for hybrid local/fleet execution across GitMap:
- **Local Host Execution (Direct In-Process):** The current machine (`127.0.0.1` / localhost) executes directly in-process with real-time streaming output. It is never wrapped in SSH or queued through async fleet worker pools.
- **Remote Fleet Nodes (Bounded Async Dispatch):** Connected remote fleet nodes are queried for liveness (with differentiated cache TTLs and resilient probing). Online nodes receive bounded asynchronous tasks via SSH JSON IPC with low worker concurrency (default `A=2`, `H=2`). Offline nodes are skipped immediately with clear status reporting.
- **Structured Telemetry & Audit:** Every operation logs to the local Task Audit DB and Errors DB. Any remote command can be audited via corresponding `see` and `ses` (see errors ssh) commands.

---

## 2. All Commands Added & Enhanced

### 2.1 Pull-First Suite (`gitmap pa` / `gitmap pas`)
- `gitmap pull all` (alias `pa`): Immediately pulls all repositories in current workspace without blocking on pre-flight ignore prompts.
- `gitmap pull all ssh` (alias `pas`): Follows GitMap PAS Formula — executes local pull in-process and dispatches async pull tasks across active SSH fleet nodes.
- `gitmap pull-all-ssh` (alias `pas`): Alias for `gitmap pull all ssh`.

### 2.2 Ignore Remediation Suite
- `gitmap fix ignore all [-y]`: Scans all repositories for `.gitignore` issues, duplicates, and missing rules; prompts or auto-fixes with `-y`.
- `gitmap fix-ignore-all (fia) [-y]`: Direct hyphenated alias for `gitmap fix ignore all`.
- `gitmap fix ignores all ssh [-y]`: Fleet ignore fix executing locally in-process and fanning out to remote nodes over SSH.
- `gitmap fix-ignores-all-ssh (fias) [-y]`: Follows GitMap PAS Formula — fixes `.gitignore` issues across all nodes (current node runs in-process as is; remote fleet nodes run over SSH with `-y`).

### 2.3 Ignore Management CLI (`gitmap ignore` / `gitmap ig`)
- `gitmap ignore add <entry>` (alias `ig add`): Adds ignore entry to current repository `.gitignore`.
- `gitmap ignore scan` (alias `ig scan`): Scans tracked repositories for duplicate lines, syntax errors, and missing `.gitmap/` or resume task ignores.
- `gitmap ignore scan-ssh` (alias `ig ss`): Scans ignore issues across remote fleet nodes using the PAS Formula.
- `gitmap ignore remove <entry>` (alias `ig remove`): Removes ignore entry from repository `.gitignore`.
- `gitmap ignore edit` (alias `ig edit`): Opens repository `.gitignore` in default system editor.
- `gitmap ignore action` (alias `ig action`): Interactive menu to triage and resolve ignore anomalies.
- `gitmap ignore ls` (alias `ig ls`): Lists active ignore patterns in current repository.
- `gitmap ignore help` (alias `ig help`): Displays comprehensive ignore help.
- `gitmap ignore ui` / `app` (alias `ig ui` / `ig app`): Launches terminal or web UI for ignore management.
- `gitmap ignore add-group <name>`: Creates a named group of ignore patterns.
- `gitmap ignore remove-group <name>` (alias `rm-grp`): Deletes a named ignore group.
- `gitmap ignore set-default-group <name>`: Sets the default ignore group automatically applied to new or existing repos.
- `gitmap ignore add-grp-to-default <name>` (alias `agtd`): Appends a named group to the default groups list.
- `gitmap ignore apply <group-name>`: Applies a named ignore group to repositories.
- `gitmap ignore connect-group-with-repo <group-name> <path1,alias> [--add-with-default]` (alias `cgwp` / `--add-with-default` alias `awd`): Links an ignore group with target repositories.
- `gitmap ignore export` / `import`: Exports and imports ignore configuration definitions.

### 2.4 Commit & Push All Repositories (`gitmap cpar`)
- `gitmap commit-push-all-repos` (alias `cpar`) `[-y]`: Stages changes, commits with standard message, and pushes across all dirty repos.
- `gitmap commit-push-all-repos (cpar) --review` (alias `-r`): Displays pending commits across dirty repos for interactive review before committing and pushing.
- `gitmap commit-push-all-repos (cpar) --review --commit-only` (alias `-co`): Displays pending changes, reviews, and commits without pushing to remote.

### 2.5 Observability & Inspection Suite (`gitmap see`)
- `gitmap see commit pending`: Lists all repositories with uncommitted or unstaged changes.
- `gitmap see git-ignore issues` (aliases `see ignore issues`, `see ig issues`): Lists repositories with `.gitignore` duplicates or issues.
- `gitmap see errors`: Displays local error history (alias for `gitmap errors`).
- `gitmap see history`: Displays recent execution history (alias for `gitmap history`).
- `gitmap see errors ssh` (alias `ses`): Queries and aggregates error logs from remote fleet nodes over SSH following the PAS Formula.
- `gitmap repo-manage ui`: Opens interactive repository management UI.

### 2.6 Repository Split-DB Cache Engine (`gitmap cache`)
- `gitmap cache create .`: Scans and indexes the current repository into split SQLite databases.
- `gitmap cache create <relpath1>,<abspath2>`: Indexes specified paths into split cache databases.
- `gitmap cache create "a.json", "b.json"`: Indexes specific files into cache databases.
- `gitmap cache ls`: Lists indexed files and folder slugs in the cache.
- `gitmap cache add <path...>`: Adds individual files or directories to the cache database.
- `gitmap cache remove <path...>` (alias `rm`): Removes files or paths from the cache.
- `gitmap cache help`: Displays help text and usage examples for cache subcommands.
- `gitmap cache search "text search" "*.md" [--lines 10] [--limit 20]`: Full-text search across indexed files with optional glob pattern, line context, and match limit.
- `gitmap cache search "text search" -file-pattern (fp) "a*.md", "b*.md" [--lines 10] [--limit 20]`: Full-text search filtered by file patterns.
- `gitmap cache search-multi "text search", "multi *"`: Multi-term search across the indexed repository.
- `gitmap cache search-multi-grep "regex search", "multi *"`: Multi-pattern regular expression search across indexed content.
- `gitmap cache recache` / `reconcile` / `sync`: Reconciles cached file hashes and updates modified files.
