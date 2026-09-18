# Plan 206: Pipeline Errors AGY Fix Hardening and Multi-Project Batch Orchestration

> **Task Type:** Hardening, Background Injection Detachment, Payload Formatting, Documentation & Help Parity  
> **Execution Context:** Root Task N-Step Loop (Steps 1 to 5)  
> **Outcome:** Successfully verified and deployed direct Antigravity injection with headless permissions bypass, process detachment (`DETACHED_PROCESS`), full absolute file path display, full CI/CD error log embedding, prompt queue ledger management, and multi-project parallel batching across single-repo and workspace scopes.

## Background & Starting Context
The user encountered issues when running `gitmap pipeline errors agy fix` in target repository `D:\work\Antigravity-Manager`:
1. `active-agy-pipeline-fix-prompt.txt` contained only the `# Verification Check: Is It Fixed?` text without error logs because earlier versions allowed the followup prompt stage to overwrite the primary payload.
2. In terminal output, paths were printed as relative paths (`.ai-memory/temp/...`).
3. Direct Antigravity IDE/CLI injection exited immediately because headless `agy.exe` was invoked without `--dangerously-skip-permissions` and without process detachment on Windows.
4. Total payload size in terminal readout matched the fix prompt size exactly instead of including the error report and commit history.
5. User requested explicit support and terminal help explanation for running parallel batch error fixes across multiple projects (default: 3 projects) with pagination cursor state.

## Core Implementations & Fixes Applied

### 1. Direct Antigravity Injection with Process Detachment & Auto-Approve Permissions
- **File:** `cli/cmdagy/agy_proc_windows.go` & `cli/cmdagy/agy_proc_other.go`
- Configured Windows process attributes using `DETACHED_PROCESS` (`0x00000008`) to ensure background processes launched by `gitmap` survive terminal exit.
- **File:** `cli/cmdagy/agy_fix_pipeline_inject.go`
  - Added `--dangerously-skip-permissions` to `agy.exe` invocations in headless background mode so autonomous tool calls are not rejected.
  - Redirected background process stdout and stderr to `.ai-memory/temp/agy-injection.log` inside the target repository directory.
  - Added `cmd.Process.Release()` to safely detach child process handles.

### 2. Full Error Logs Inside Text File Payload
- **File:** `cli/cmdagy/agy_fix_pipeline_assembly.go`
  - `AssembleRcaFixPayload` bundles the complete 4-part RCA directives, git commit history, and full failing CI/CD pipeline error report into `active-agy-pipeline-fix-prompt.txt`.
- **File:** `cli/cmdagy/agy_fix_pipeline_persist.go` & `cli/cmdagy/agy_fix_pipeline_queue.go`
  - `active-agy-pipeline-fix-prompt.txt` retains the primary fix payload.
  - `queued-agy-followup-prompt.txt` stores the secondary `# Verification Check: Is It Fixed?` verification prompt.
  - Both tasks are logged in `.ai-memory/temp/agy-prompt-queue.json` under `active` and `queued` fields.

### 3. 100% Full Absolute File Paths in Terminal
- **File:** `cli/cmdagy/agy_fix_pipeline_render.go`
  - All displayed paths (`Saved Payload`, `Follow-up Verification Prompt Queued`, `Queue Ledger`, `Fix Prompt`) are wrapped with `toAbsPath(filepath.Clean(...))`.

### 4. Multi-Project Parallel Batch Scanning & Pagination
- **File:** `cli/cmdagy/agy_fix_pipeline_batch_scan.go`, `batch_exec.go`, `batch_cursor.go`
  - Concurrently scans all candidate projects for pipeline error reports.
  - Batches at 3 projects by default (`--projects 3` / `--limit 3`).
  - Persists cursor in `.ai-memory/temp/pipeline-fix-batch-cursor.json`.
  - Subsequent invocations automatically process subsequent batches; `--reset-batch` resets cursor.

### 5. CLI Help, Web UI Documentation, and Antigravity Skill Parity
- **File:** `cli/helptext/pipeline.md`: Added comprehensive subcommands, flags (`--all`, `--projects`, `--limit`, `--reset-batch`, `--no-inject`), shortcuts (`aef`), and examples.
- **File:** `cli/helptext/agy-fix-pipeline.md`: Added direct injection, full paths, error log embedding, and batch examples.
- **File:** `cli/helptext/agy.md`: Added Section 5 for `agy fix-pipeline (aef)`.
- **File:** `src/data/commands.ts`: Documented `pipeline errors agy fix` in web UI commands dataset.
- **File:** `.agents/skills/pipeline-errors-agy-fix-injection-and-multi-project-batching/SKILL.md`: Documented full architecture.

### 6. Binary Synchronization
- Deployed freshly compiled binary to:
  - `C:\Users\Administrator\AppData\Local\gitmap-cli\gitmap.exe`
  - `d:\work\gitmap\bin\gitmap.exe`
- Verified live runs in both `d:\work\gitmap` and `D:\work\Antigravity-Manager`.
