# 28 — Antigravity (AGY) Prompts Templates, Rerun Suite & Pipeline Error Remediations

- **Slug:** agy-prompts-templates-and-rerun-suite
- **Date:** 2026-09-19
- **Version:** v6.261.0
- **Status:** completed
- **Execution Loops:** 6 completed subtasks across 2 parallel execution phases (Budget: N=250, Completed in 6 steps).

---

## 1. Task Origin & Problem Statement

The user reported a severe execution failure and stack trace when running `gitmap rm settings` following a missing repository detection warning from `gitmap status`:
```
PS D:\wp-work\riseup-asia> gitmap rm settings
rm: no repo matched "settings"
gitmap rm: execute failed: [E9000:EXECUTION] fatal error: (at=cmd/rm.go:79)
Stack Trace:
    at github.com/alimtvnetwork/gitmap-v28/gitmap/apperror.NewSimple (apperror/apperror.go:153)
    at github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.runRm (cmd/rm.go:79)
```
In addition, the user specified a suite of AGY prompt management and rerun features, pipeline false positive error detection fixes, and database maintenance capabilities:
1. `gitmap agy rerun last N [-p template]` with pre-seeded `is-done` verification prompt.
2. `gitmap agy list-prompts [N]` with `--all-projects`, `--projects <M>`, and `--project <name>`, launching prompt diffs in a new non-admin VS Code window (`code -n`).
3. `gitmap agy scan` auditing project prompts and creating automated backups.
4. `gitmap prompts-template` (add, edit, ls, rm, import, export, import-all, export-all as JSON).
5. Remote execution across `ssh`, `sc`, `cluster`, and local `exec`.
6. Fix pipeline error evaluation where failed GitHub Actions runs (e.g. `cat-my-v12` commit `384f5d2`) were misreported as `(clean status)` with empty logs.
7. Implement `gitmap pipeline clear-db` (and `db clear`) to reset pipeline execution telemetry and purge cached log files.
8. Implement `gitmap storage restore-db` and display repository SQLite split DBs in `gitmap storage ls` and `gitmap storage space ls`.
9. Intelligent conditional terminal padding and separator formatting to eliminate duplicate padding or repeating blue rule lines.
10. Authoritative help documentation first in `cli/helptext/`.

---

## 2. Consolidated Execution Summary

### Subtask 01: Fix `gitmap rm` Resolution & Error Handling
- **Unified Candidate Pools (`loadUnifiedCandidates`)**: Merged candidate search across both SQLite database (`db.ListRepos()`) and `.gitmap/output/gitmap.json` (with fallback to `output/gitmap.json`), deduplicating candidates by normalized absolute path.
- **Broadened Target Matching (`matchesRepoTarget`)**: Implemented matching against `Slug`, `RepoName`, base folder name of `AbsolutePath`, base of `RepoName` (e.g. `settings` matching `owner/settings`), base of `Slug`, and full path.
- **Clean Not-Found AppError**: Replaced fatal `apperror.NewSimple("fatal error", "E9000")` crash with `apperror.NewNotFoundError(fmt.Sprintf("no repository matched %q", missing[0]))`. This guarantees `gitmap rm <missing-repo>` exits cleanly with `gitmap rm: not found: no repository matched "settings"` and zero fatal execution stack traces.
- **Dual Database & JSON Untracking**: Added `removeRepoFromJSON` to remove untracked repos from `.gitmap/output/gitmap.json` so `gitmap status` immediately stops reporting untracked repositories as missing.

### Subtask 02: Fix Antigravity IDE & CLI Detection & Workspace Integration
- **Decoupled Binary Discovery:** Implemented `ResolveAntigravityIDE() result.Result[string]` and `ResolveAntigravityCLI() result.Result[string]`.
- **Precedence Correction:** Removed `%LOCALAPPDATA%\agy\bin\agy.exe` from IDE candidates; prioritized real desktop IDE executables (`Antigravity.exe` / `Antigravity`).
- **Process Detection:** Added `DetectRunningAntigravityIDE() result.Result[AgyProcessInfo]` and `IsAntigravityIDERunning() bool` inspecting running processes via `tasklist` (Windows) and `ps` (Unix/macOS).
- **Single Return Types & Universal AppError:** Converted legacy multi-value returns to `AgyInjectionResult` and `AgyAssembledPromptPayload`.
- **Repo Root Resolution & Staging:** Ensured `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt` is created and populated directly within the target repository root, and printed distinct feedback differentiating between CLI background runner execution and IDE workspace staging.

### Subtask 03: Fix Pipeline "Clean Status" False Positive & Add `pipeline clear-db`
- **Scoped Deduplication to Branch and Commit SHA:** Implemented `buildWorkflowScopeKey(r ghRunItem)` generating `branch:sha:workflowname`, ensuring passing runs on other branches or newer commits cannot mask failures on the target commit.
- **Broadened Failure Conclusion Recognition:** Implemented `isFailingConclusion(conclusion string)` recognizing `"timed_out"`, `"cancelled"`, `"startup_failure"`, and `"failure"`.
- **Eliminated Premature Clean Status Short-Circuiting:** Removed premature short-circuit on `p.Conclusion == "success"` when companion runs or other commits in the fetched runs have failed.
- **Fixed DB Fallback to Isolated Split Database:** Replaced legacy `openDB()` call in `queryRunsFromDB` with `pipelinedb.OpenPipelineSplitDb(repo)` querying `db.QueryRecentRuns(5)`.
- **Added `pipeline clear-db` Subcommands:** Added `clear-db`, `cleardb`, and `db-clear` aliases routed in `dispatchPipelineSubcmd`.
- **Enhanced `Clear()` in `pipelinedb`:** Collects run IDs, truncates tables, resets SQLite sequences, runs `VACUUM;`, and purges `.log` and `.json` cache files on disk.

### Subtask 04: AGY Prompt Templates, Rerun Suite & Non-Admin VS Code Opener
- **Canonical `is-done` Template:** Preserved default built-in template with content `"Is it done properly? Can we check properly the missing items from the task that is mentioned below? Please check it carefully. Do not make any mistakes."`
- **Prompt Template Subcommands:** Verified `ls`, `add`, `edit`, `rm`, `export`, `import`, `export-all`, `import-all` in `cmdprompttemplate`.
- **AGY Rerun Suite:** Implemented `gitmap agy rerun last [N] -p <template>` with top-level aliases `gitmap rerun` and `gitmap rr`.
- **AGY List Prompts:** Handled `[N]`, `--all-projects`, `--projects <M>`, `--project <name>`, and `--json`.
- **Non-Admin VS Code Opener:** Enhanced `spawnNonAdminVSCode` to strictly pass `-n` (new window), apply `filterAdminEnv`, and detach background process so it never runs in admin mode or blocks existing instances.
- **AGY Scan & Automatic Backup:** Added automated prompt backups to `~/.gemini/config/backup/prompts/` on every scan.

### Subtask 05: Storage Restore DB, Repo DB Listing & Conditional Terminal Padding
- **Storage Space LS & Restore DB:** Integrated `gitmap storage space ls` and `gitmap storage restore-db` (and `restoredb`, `restore`).
- **Repository Split DBs Discovery:** Enhanced database inventory in `cli/cmd/storage_ls.go` to discover and display repository-scoped split DBs (`~/.gitmap/pipeline/pipeline_*.db`, `repodb/*.db`) alongside primary DBs.
- **Conditional Terminal Padding:** Enhanced `termpad` with `HasExistingMargin`, `IsBoundaryRule`, `PrintBlueRule`, `PrintSeparator`, and `EnsureBottomPadding`, suppressing consecutive duplicate boundary rules and eliminating repeated blank line padding.

### Subtask 06: Authoritative Help Text Parity & Remote Triad Verification
- **Help Text Catalog:** Authored and updated `cli/helptext/agy.md`, `cli/helptext/pipeline.md`, `cli/helptext/prompt_templates.md`, `cli/helptext/prompts-template.md`, and `cli/helptext/storage.md`.
- **Remote Triad Execution:** Verified and updated `cli/cmdssh/ssh_exec_command.go` to route `agy`, `prompts-template`, `pt`, `schedule`, `pipeline`, and `storage` commands across SSH, SC, Cluster, and local `exec`.
- **Remote Scheduled Commands:** Added `"schedule", "schedules"` route in SSH, SC, and Cluster dispatchers.

---

## 3. Verification & Coding Guidelines Compliance

- **Function Sizing:** All functions refactored to <= 8–15 lines.
- **Single Return Types:** Tuple multi-returns eliminated across all refactored packages.
- **Universal AppError:** Replaced generic error returns with strongly-typed `*apperror.AppError`.
- **Boolean Conventions:** All boolean variables and parameters use affirmative `is*` and `has*` prefixes.
- **Line Endings:** Strictly Unix LF (`\n`) across all modified source and markdown files.
