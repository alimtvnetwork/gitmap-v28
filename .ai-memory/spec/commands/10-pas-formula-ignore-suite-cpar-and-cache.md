# Command Specification: GitMap PAS Formula, Ignore Management Suite, CPAR & Split-DB Cache Engine

Spec Reference: [02-spec/21-app/197-gitmap-pas-fix-and-repo-cache-commands.md](../../../02-spec/21-app/197-gitmap-pas-fix-and-repo-cache-commands.md) & [02-spec/21-app/181-gitmap-ignore-and-cache-engine/01-overview.md](../../../02-spec/21-app/181-gitmap-ignore-and-cache-engine/01-overview.md)

## 1. The GitMap PAS Formula

The **GitMap PAS Formula** defines the canonical orchestration standard for hybrid local/fleet execution across GitMap:
- **Local Host Execution (Direct In-Process):** The current machine (`127.0.0.1` / localhost) executes directly in-process with real-time streaming output. It is never wrapped in SSH or queued through async fleet worker pools.
- **Remote Fleet Nodes (Bounded Async Dispatch):** Connected remote fleet nodes are queried for liveness (with differentiated cache TTLs and resilient probing). Online nodes receive bounded asynchronous tasks via SSH JSON IPC with low worker concurrency (default `A=2`, `H=2`; `A=1`, `H=2` under high CPU pressure). Offline nodes are skipped immediately with clear status reporting.
- **Structured Telemetry & Audit:** Every operation logs to the local Task Audit DB (`gitmap.db`) and Errors DB. Any remote command can be audited via corresponding `see` and `ses` (see errors ssh) commands.

---

## 2. All Commands Added & Enhanced (Exhaustive Screenshot Parity)

### 2.1 Pull-First Suite (`gitmap pa` / `gitmap pas`)
- `gitmap pull all` (alias `pa`): Immediately pulls all repositories in current workspace without blocking on pre-flight ignore prompts. Asynchronous background workers inspect `.gitignore` files for duplicates/tracked files and summarize at completion.
- `gitmap pull all ssh` (alias `pas`): Follows GitMap PAS Formula — executes local pull in-process and dispatches async pull tasks across active SSH fleet nodes using bounded workers (`A=2`, `H=2`).
- `gitmap pull-all-ssh` (alias `pas`): Direct alias for `gitmap pull all ssh`.

### 2.2 Ignore Remediation Suite (`gitmap fia` / `gitmap fias`)
- `gitmap fix-ignores-all (fia)` or `gitmap fix-ignore-all (fia) [no prompt -y]`: Scans all repositories for `.gitignore` duplicates and tracked ignore issues; prompts with choices (all at once vs interactive repo-by-repo session) or auto-applies with `-y`.
- `gitmap fix-ignores-all-ssh (fias)` or `gitmap fix-ignore-all-ssh (fias)`: Follows GitMap PAS Formula — repairs ignore issues across fleet nodes (local runs in-process as is; remote fleet nodes run over SSH with `-y`).

### 2.3 Commit & Push All Repositories (`gitmap cpar`)
- `gitmap commit-push-all-repos` (alias `cpar`) `[-y]`: Stages changes, commits with standard message, and pushes across all dirty repos.
- `gitmap commit-push-all-repos (cpar) -r` (alias `--review`): Displays pending commits across dirty repos as folder trees for interactive review before committing and pushing.
- `gitmap commit-push-all-repos (cpar) --commit-only` (alias `--co`): Commits all staged changes without pushing to remote.
- `gitmap commit-push-all-repos (cpar) [-r] -y [--co]`: Non-interactive auto-confirm variant for batch automated scripts.

### 2.4 Observability & Inspection Suite (`gitmap see` / `gitmap c`)
- `gitmap see commit pending (gitmap c commit pending)`: Lists all repositories with pending commits, displaying tree views (tree limit default 10). Prompts the user: "Would you like to commit all?". If confirmed, immediately invokes `gitmap commit-push-all-repos -y`.
- `gitmap see git-ignore issues (gitmap c ig issues)`: Finds duplicate patterns and tracked ignore files across all repositories.
- `gitmap see git-ignore issues -d (gitmap c ig issues -d)`: Filters and shows duplicate ignore pattern issues only.
- `gitmap see git-ignore issues -t (gitmap c ig issues -t)`: Filters and shows tracked ignore file issues only.
- `gitmap see errors (gitmap c errors) [--limit 10]`: Displays local error journal from GitMap Errors DB.
- `gitmap see errors ssh (gitmap c errors ssh or ses) [--limit 10]`: Queries and displays error logs across remote SSH fleet nodes following the PAS Formula.
- `gitmap see history` / `gitmap history ssh` / `gitmap nodes history`: Inspects task dispatch and execution history across local and remote nodes.
- `gitmap repo manage ui` (alias `gitmap repo-manage ui`): Launches the repository management interface for Git, GitHub, and GitLab fleet cloning and status overview.

### 2.5 Ignore Management Engine (`gitmap ignore` / `gitmap ig`)
- `gitmap ignore connect-group-with-repo (gitmap ig cgwp) <group-name> [--add-with-default] [--awd]`: Binds an ignore group to a repository path or alias. Links repo-specific system rules (`system-rules/<alias>`).
- `gitmap ignore add-grp-to-default (gitmap ig agtd) <group-name>`: Chains the specified group execution to run automatically after the default ignore group.
- `gitmap ignore remove-group-from-default (gitmap ig rgfd) <group-name>`: Removes a group from the default chained list.
- `gitmap ignore remove-group-from-repo (gitmap ig rgfr) <repo-name> <group-name>`: Unlinks an ignore group from a specific repository.
- `gitmap ignore list-repo-groups`: Lists all connected ignore groups and repository bindings.
- `gitmap ignore add <entry>` (alias `ig add`): Adds ignore entry to current repository `.gitignore`.
- `gitmap ignore scan`: Scans tracked repositories for duplicate lines, syntax errors, and missing ignores.
- `gitmap ignore scan-ssh` (alias `ig ss`): Scans ignore issues across remote fleet nodes using the PAS Formula.
- `gitmap ignore remove <entry>` (alias `ig remove`): Removes ignore entry from repository `.gitignore`.
- `gitmap ignore edit` (alias `ig edit`): Opens repository `.gitignore` in default system editor.
- `gitmap ignore action` (alias `ig action`): Interactive menu to triage and resolve ignore anomalies.
- `gitmap ignore ls` (alias `ig ls`): Lists active ignore patterns in current repository.
- `gitmap ignore help` (alias `ig help`): Displays comprehensive ignore help.
- `gitmap ignore ui` / `app` (alias `ig ui` / `ig app`): Launches terminal or web UI for ignore management.
- `gitmap ignore add-group <name>`: Creates a named group of ignore patterns.
- `gitmap ignore remove-group <name>` (alias `rm-grp`): Deletes a named ignore group.
- `gitmap ignore set-default-group <name>`: Sets the default ignore group automatically applied to repositories.
- `gitmap ignore apply [folder|repo]`: Applies default and connected ignore rules to the target directory.
- `gitmap ignore export` / `import`: Exports and imports ignore configuration definitions in JSON/YAML.

### 2.6 Repository Split-DB Cache Engine (`gitmap cache`)
- `gitmap cache create`: Scans and indexes current repository into multi-tier split SQLite databases (`sql.db` root + `<folder-slug>.db` per top-level folder).
- `gitmap cache create <relative_path>`: Creates cache for a specific relative repository folder.
- `gitmap cache create --ignore ".git,node_modules"`: Creates cache with custom comma-separated ignore folders.
- `gitmap cache create search-text <text-to-find>`: Creates cache and immediately executes text search.
- `gitmap cache create search-text <text-to-find> -ext md`: Creates cache and searches within markdown files.
- `gitmap cache search-text <text-to-find> -ext md,ts`: Searches indexed text filtered by extension (default limit 20 matches, 10 lines context window from-to).
- `gitmap cache search-text <text-to-find> -format/ext md -s/all "all"`: Searches across all indexed text matching format/extension.
- `gitmap cache search-multi "word-word", "word..." (filesystem fly) "s/all" "c/ext" -limit 20`: Multi-term search across indexed files with filesystem-fly reconciliation.
- `gitmap cache search-grep "regex-search" "s/all" -format/ext [ts/md] --limit 20`: Ultra-fast SQLite-native regular expression search with automatic file `mtime` reconciliation.
- `gitmap cache reconcile` / `recache` / `sync`: Checks filesystem modification times (`mtime`) against `sql.db` and updates modified files without hashing.
- `gitmap cache ls`: Lists indexed repositories, folder slugs, and cached file counts.
- `gitmap cache add <path...>`: Adds individual files or directories to the split cache.
- `gitmap cache remove <path...>` (alias `rm`): Removes cached paths from the database.
- `gitmap cache help`: Displays comprehensive help documentation and performance guidelines.
