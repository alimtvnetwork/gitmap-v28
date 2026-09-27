# Completed Plan: Antigravity Running Prompts Backup, Restore, and Automation

Spec Reference: [02-spec/21-app/167-agy-running-prompts-backup-restore-and-green-automation/01-overview.md](../../../02-spec/21-app/167-agy-running-prompts-backup-restore-and-green-automation/01-overview.md)

## Summary of Accomplishments
1. **Prompts Split-DB & Operations (`cli/store/backup_prompts_split_*.go`)**:
   - Implemented `OpenBackupPromptsSplitDB(customPath)` defaulting to `filepath.Join(store.BinaryDataDir(), "backup-prompts", "sql.db")`.
   - Created tables `PromptBackupBatch`, `PromptBackupItem`, and `PromptRestoreLedger`.
   - Built operations for batch insertions, listing, restoration tracking, 24h retention auto-pruning with `"old data has been removed"` notice, force cleanup, and storage metrics reporting.
2. **CLI Commands, Rendering & I/O (`cli/cmdagy/agy_running_prompts_*.go`)**:
   - Implemented `running-prompts` (aliases `running-prompt`, `rp-prompts`), `backup-running-prompts` (aliases `backup-running-prompt`, `brp`), and `restore-running-prompts` (aliases `restore-running-prompt`, `rrp`).
   - Implemented `ls` with flags `--limit`/`-l` (default 8), `--wordcount`/`--wc` (default 100), `--full`, and `--json`.
   - Implemented `export` and `import` for SQLite `.db` and canonical `.json` formats.
   - Word count truncation via `TruncateWords(text, n)` preserving whitespace boundaries.
3. **Running Projects Discovery & Cluster SSH Aggregation (`cli/cmdagy/agy_running_projects*.go`)**:
   - Implemented `gitmap agy running-projects [ls] [--json] [-file/-f <path>] [--ssh]`.
   - Active and queued prompt discovery across local Antigravity workspaces.
   - Cluster SSH aggregation using decoupled `SSHConnectionsFetcher` callback avoiding import cycles.
4. **Finish Prompts Until Green (`cli/cmdagy/agy_fpug.go`)**:
   - Implemented `gitmap agy finish-prompts-until-green` (alias `fpug`) `<target...> [-t <duration>]`.
   - Target keyword `running-projects` dynamically expands to all active projects.
   - Configurable polling interval with minimum 30s enforcement and non-intrusive `gitmap pe` checks.
5. **Shutdown Until Green (`cli/cmdagy/agy_sug*.go`)**:
   - Implemented `gitmap agy shutdown-until-green` (alias `sug`) `ls/help/run/add-projects/rm/agy-running-projects [-t <duration>]`.
   - Persisted watch configuration in `~/.gemini/antigravity/sug_watch_list.json`.
   - Cross-platform OS shutdown execution across Windows, Ubuntu/Linux, and macOS with mockable executor `OSShutdownExecutorFn` ensuring zero unwanted shutdowns during automated tests.
6. **Pipeline Error Fast DB Cache & Decoupled Fix Dispatch (`cli/cmdpipeline/`)**:
   - In `cli/cmdpipeline/pipeline_cache_eval.go`: added `checkHeadShaCacheHit` so if current local HEAD commit SHA is recorded in SQLite DB, telemetry is served immediately without remote network calls.
   - In `cli/cmdpipeline/pipeline_logs.go`: added `params.WantFix` check so `gitmap pe` only dispatches automatic AGY fix when `--fix` is explicitly provided.
7. **Top-Level Root Aliases (`cli/cmd/root.go`)**:
   - Added `backup-running-prompts`, `restore-running-prompts`, `running-prompts`, `running-projects`, `fpug`, and `sug` directly to root command routing table so commands can be called with or without `agy` prefix.
