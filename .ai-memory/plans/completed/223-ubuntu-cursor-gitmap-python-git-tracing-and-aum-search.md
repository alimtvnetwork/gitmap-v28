# Completed Plan: 223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search

> **Plan ID:** 223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search  
> **Status:** COMPLETED  
> **Release:** v6.484.0  
> **Completion Date:** 2026-10-05  

---

## 1. Summary of Completed Deliverables

### A. Cursor Ubuntu Fleet Deployment & repo-secrets Tracking
- Implemented `repo-secrets/05-scripts/setup-cursor-ubuntu.py` and `setup-cursor-ubuntu.sh`:
  - Dynamically fetches official upstream Cursor metadata from `https://cursor.com/api/download?platform=linux-x64&releaseTrack=stable`.
  - Downloads AppImage to `/opt/cursor/Cursor.AppImage` (SHA256 `6ba0b06a8e9087688f312b8484e61ce28c185b609c3731b8906e17a6aa98a238`).
  - Generates sandboxed launcher wrapper at `/usr/local/bin/cursor` and `/home/a/.local/bin/cursor`.
  - Configures Dracula theme and development invariants in `/home/a/.config/Cursor/User/settings.json`.
  - Persists fleet health state in `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json` (`status: HEALTHY`).
- Updated `cli/cmdcursor/cursor_install.go` to delegate `--node <alias>` to this setup script via SSH.

### B. Native GitMap Python Runner (`gitmap py` / `gitmap python`)
- Implemented `cli/cmdpy/py_cmd.go`, `py_exec.go`, `py_record.go`, and `py_cmd_test.go`:
  - Directly executes Python code snippets via `gitmap py -c "<code>"`.
  - Executes Python script files with arguments via `gitmap py <script.py> [args...]`.
  - Transparently pipes `stdin`, `stdout`, and `stderr`, and propagates exit codes.
  - Automatically records telemetry into `CommandHistorySplitDB` and `RecordAiExecution`.
  - Registered `[]string{"py", "python"}` under `rootutility.go` and `rootcore.go`.

### C. Git Command Split-DB Heatmap Tracing
- Updated `cli/cmd/rootgit.go`:
  - Instrumented `runGitPassthrough` to measure monotonic execution duration and extract child exit code.
  - Records every raw Git invocation (`git clone`, `git status`, `git branch`, etc.) to `store.CommandHistorySplitDB` (`commands.db`).
- Updated `cli/cmd/commit_push.go`:
  - Instrumented `performCommitPush`, `executePullCommitPush`, `stageAndCheckChanges`, and `pushUnpushedCommitsCP` to record commit and push events to `commands.db` for developer heatmap analytics.
- Enhanced `cli/store/command_history_split_db.go`:
  - Configured `PRAGMA busy_timeout = 3000;`.
  - Added `RecordCommandSafely` to ensure non-blocking, non-fatal telemetry writes during high-concurrency Git operations.

### D. AUM Search Quote Unquoting, Regex Normalization & Path Resolution
- Updated `cli/cmdautomation/search.go` and `cli/cmdautomation/automation_cmd.go`:
  - Added `cleanSearchPattern` to strip shell-escaped quotes (`\"...\"`), single quotes, and backticks.
- Added `cli/cmdautomation/search_regex_helper.go`:
  - Implemented `normalizeRegexPattern` to normalize escaped alternation pipes (`\|` -> `|`) so regex patterns like `"timer\|Clock\|8:47"` and `"Candidate Response\|CANDIDATE RESPONSE"` match correctly in Go RE2.
  - Implemented `autoPromoteRegex` to automatically promote patterns with regex syntax (`.*`, `.+`, `\|`, `|`, `[0-9]`, etc.) to regex mode even when `-r` is omitted by AI agents.
- Added `cli/cmdautomation/search_path_resolve.go`:
  - Tolerant file path resolution: handles exact file paths, auto-resolves missing extensions (`FormRunner` -> `FormRunner.tsx`), and resolves prefix matches (`FormRunner.t` -> `FormRunner.tsx`).
- Created unit tests in `cli/cmdautomation/search_quote_test.go` and `search_regex_test.go` verifying all scenarios pass.

### E. Quality Gates & Release Ceremony
- All linters passed with 0 violations:
  - `check-relative-paths.py`: 0 violations across 8,718 files.
  - `check-nested-ifs.py`: 0 violations.
  - `check-boolean-guidelines.py`: 0 violations.
- Bumped minor release to `v6.484.0` via `03-ai-scripts/37-bump-version.py -t minor`.
- Verified release notes in `.ai-memory/release/release-notes-v6.484.0.md` strictly contain relative paths and zero IPs.
