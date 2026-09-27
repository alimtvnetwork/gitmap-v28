# Consolidated Plan 178: Smart File and Folder Deployment with Bidirectional Sync Modes

## 1. Specification Reference & Execution Metadata
- **Canonical Spec Reference:** [02-spec/21-app/178-smart-deploy-files-and-folders-with-sync-modes.md](../../../02-spec/21-app/178-smart-deploy-files-and-folders-with-sync-modes.md)
- **Status:** Completed
- **Execution Date:** 2026-09-27
- **Total Steps / Loops Executed:** 1 Phase 1 Planning Turn, 1 Dual-Subagent Parallel Turn (A = 2, H = 2)

---

## 2. User Request (Verbatim)

```text
these commands should work

 gitmap deploy <alias, ip, seq, id> <from or current releative file or folder> <abs path or path from default work dir into that machine> # clear???

 gitmap deploy <alias, ip, seq, id> <from or current releative file or folder> <abs path or path from default work dir into that machine> --json # clear???, will show summary response in json mode, if same file match or found then will ask for permission if we want to replace to skip

gitmap deploy <alias, ip, seq, id> <from or current releative file or folder> <abs path or path from default work dir into that machine> --overrite(o)/skip/sync/sync-right/sync-left  # don't seek for prompt, copy and replace , sync means it will only take if left one is latest by date or copy from right to left if the right has the latest date, for folder can do parallel replace ?? can???, sync right means don't pass any file to left, sync left means keep the left intact right on conflicts 

also create prompts for deploy-right, deply left
```

---

## 3. Consolidated Implementation & Subtask Outcomes

### Subtask 01: Deploy Command Dispatcher & Flag Parsing (`cli/cmd/roottooling.go`, `cli/cmd/clihelpers.go`, `cli/cmdssh/ssh.go`, `cli/cmdssh/ssh_deploy_cmd.go`)
- **Outcome:**
  - Added root dispatch table entries in `cli/cmd/roottooling.go` for `deploy`, `deploy-right`, and `deploy-left` delegating to `runSSHDeploy`.
  - Added helper `runSSHDeploy` in `cli/cmd/clihelpers.go` delegating to `cmdssh.RunSSHDeployCLI`.
  - Routed SSH subcommands `case "deploy", "deploy-right", "deploy-left":` in `cli/cmdssh/ssh.go`.
  - Implemented `RunSSHDeployCLI` in `cli/cmdssh/ssh_deploy_cmd.go` with positional parsing (`<target> <src> <dest>`), multi-identifier resolution (`<alias, ip, seq, id>`), and flag parsing (`--overwrite`/`-o`, `--skip`/`-s`, `--sync`, `--sync-right`, `--sync-left`, `--json`/`-j`, `--parallel`/`-p`, `--dry-run`/`-n`).

### Subtask 02: Sync Modes & Parallel Folder Replacement Engine (`cli/cmdssh/ssh_deploy_sync.go`, `cli/cmdssh/ssh_deploy_folder.go`, `cli/cmdssh/ssh_deploy_engine.go`)
- **Outcome:**
  - Implemented remote file metadata probing (`probeRemoteFileInfo`) via cross-platform PowerShell and POSIX stat commands.
  - Implemented conflict evaluation (`evaluateFileConflict`) supporting `--overwrite`, `--skip`, `--sync` (bidirectional mtime comparison), `--sync-right` (push latest to remote only), `--sync-left` (pull latest to local while keeping local intact on conflict), and interactive TTY fallback.
  - Implemented directory tree walker (`scanLocalFolder`) and channel-based parallel folder deployment worker pool (`executeParallelFolderDeploy`) with bounded concurrency (up to 16 workers).
  - Implemented `ExecuteDeploy` lifecycle managing local path verification, TCP liveness pre-flight checks, remote default workdir resolution (`D:/work` or `~`), single-file and folder deployment execution, and telemetry formatting.

### Subtask 03: JSON Telemetry & Canonical Prompts (`cli/cmdssh/ssh_deploy_json.go`, `cli/helptext/deploy.md`, `01-prompts/deploy-right.md`, `01-prompts/deploy-left.md`)
- **Outcome:**
  - Implemented structured JSON summary renderer `RenderDeployJSONSummary` and non-interactive conflict prompt payload `RenderConflictJSONPrompt` (with exit code 3).
  - Authored comprehensive terminal help and markdown documentation in `cli/helptext/deploy.md`.
  - Authored reusable canonical AI prompts in `01-prompts/deploy-right.md` and `01-prompts/deploy-left.md`.

---

## 4. Verification & Quality Gates
- **Linters:** 0 nested-if violations, 0 boolean violations across polyglot codebase.
- **Function Sizing:** All functions <= 8-15 lines.
- **Guideline Compliance:** Positive booleans only, structured `*apperror.AppError`, Unix LF line endings.
