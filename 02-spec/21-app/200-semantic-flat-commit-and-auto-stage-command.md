# 200 — Semantic Flat Commit & Auto-Stage Command (`gitmap commit`, `cm`, `commit-all`, `ca`) and Macro Execution Resilience

## Status: Active
- **Spec ID:** 200
- **Scope:** Application CLI, Git Automation, Shell Dispatch, Macro Runner
- **Created At:** 2026-10-01

---

## 1. Executive Summary & Problem Statement

GitMap provides workflow commands such as `cpf` (`commit-push-feature`) and `cpb` (`commit-push-bug`), but developers frequently require a flat, flexible 1-step auto-staging commit command that:
1. Automatically stages all tracked modifications, untracked new files, and deletions (`git add -A`) without requiring manual two-step `git add` and `git commit` commands.
2. Accepts arbitrary flat commit messages without forcing semantic `feat:` or `fix:` prefixes.
3. Supports unquoted multi-word positional arguments (e.g. `gitmap cm update readme with install instructions`) or quoted strings (`gitmap cm "refactor: clean up split-db"`).
4. Provides optional `--push` / `-p` flags to push immediately to the remote branch.
5. Provides optional `--dry-run` / `-n` flags to inspect changes (`git status --short`) without committing.
6. Intercepts `--help` / `-h` to render a styled, boxed terminal help menu.

Additionally, this specification addresses the macro execution regression where Windows PowerShell's native non-idempotent `rm` cmdlet (`Remove-Item`) threw terminating `ItemNotFoundException` (`exit status 1`) when encountering absent paths, mandating the active runtime deployment of GitMap's cross-platform `safe-rm` adapter and synchronized binary currency.

---

## 2. Command Architecture & Routing

### 2.1 Aliases and Precedence
The command is bound to four top-level aliases defined in `cli/constants/constants_cli.go`:
- `gitmap commit "<msg>"`
- `gitmap cm "<msg>"`
- `gitmap commit-all "<msg>"`
- `gitmap ca "<msg>"`

### 2.2 Dispatch Collision Resolution
Previously, `cli/cmd/rootcore.go` contained a legacy passthrough handler:
```go
{[]string{"commit", "cm"}, func() error { return runCommitCLI(argsTail()) }},
```
which shadowed new tooling registrations in `cli/cmd/roottooling.go`. 

The architecture consolidates commit dispatching exclusively in `cli/cmd/rootcore.go` using the centralized constant references:
```go
{[]string{constants.CmdCommit, constants.CmdCommitAlias, constants.CmdCommitAlias2, constants.CmdCommitAlias3}, func() error { return runCommit(argsTail()) }},
```
The obsolete `cli/cmd/commit_cli.go` file is purged, eliminating shadow routing.

---

## 3. Flag Specification & Execution Flow

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--push` | `-p` | Boolean | Push committed changes to remote repository immediately |
| `--dry-run` | `-n` | Boolean | Preview changes to be committed via `git status --short` |
| `-m` | | String | Explicit commit message (optional; positional words accepted) |
| `--help` | `-h` | Boolean | Display the boxed terminal help menu |

### 3.1 Execution Pipeline
1. **Help Check:** If help flags are detected (`-h`, `--help`), render `RenderCommitHelp()` and return nil.
2. **Flag Parsing:** `parseCommitFlags(args)` extracts clean message tokens, `hasPush`, and `isDryRun`.
3. **Dry Run Preview:** If `isDryRun == true`, invoke `execGitInheritCP("status", "--short")` and exit without mutation.
4. **Auto-Stage:** Execute `git add -A`. If staging fails, return `apperror.WrapSimple(err, "git add failed:")`.
5. **Commit Dispatch:** Execute `git commit -m "<msg>"` (or open editor if no message was provided).
6. **Optional Remote Push:** If `hasPush == true`, execute `git push`. If push fails, return `apperror.WrapSimple(err, "git push failed:")`.

---

## 4. Macro Execution Resilience & Safe Removal

When macro steps execute deletion operations on Windows (`rm <target>`, `Remove-Item <target>`, `del <target>`), GitMap intercepts the command in `cli/macro/safe_rm.go` and transforms it into an idempotent PowerShell loop:
```powershell
foreach ($__target in @(<targets>)) {
    if (Test-Path -LiteralPath $__target) {
        Remove-Item -Recurse -Force -LiteralPath $__target
    } elseif (Test-Path -Path $__target) {
        Remove-Item -Recurse -Force -Path $__target
    }
}
```
Furthermore, the compiled binary at `./gitmap.exe` and `%LOCALAPPDATA%\gitmap-cli\gitmap.exe` must remain synchronized with the repository code to prevent stale binary execution.

---

## 5. Verification & Acceptance Criteria

- [x] `gitmap commit "<msg>"` stages all changes and creates a flat git commit without forced semantic prefixes.
- [x] Short alias `gitmap cm "<msg>"` and alternate aliases `commit-all` and `ca` execute identical logic.
- [x] Positional multi-word arguments without quotes are correctly concatenated into a single commit message.
- [x] `--push` / `-p` pushes to remote after successful commit.
- [x] `--dry-run` / `-n` previews status without staging or committing.
- [x] Unit tests in `cli/cmd/commit_cmd_test.go` pass 100%.
- [x] Macro step execution tests in `cli/macro/safe_rm_test.go` pass 100%.
- [x] Binary synchronization updates both workspace and local app data binaries.
