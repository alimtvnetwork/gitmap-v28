# RCA-56: SSH Fleet Pull Subprocess Bypass and Wincredman Credential Store Failure

**Status:** Resolved
**Date:** 2026-10-01
**Affected Subsystem:** `cli/cmdssh/` (`ssh_pull_fleet.go`), `cli/cmdpull/` (`pullall.go`, `pull_batch_json.go`, `pull.go`), `cli/cloner/` (`safe_pull.go`, `cloner.go`, `safe_push.go`)
**Associated Spec:** [02-spec/21-app/169-pull-all-concise-summary-and-ssh-fleet-json.md](../21-app/169-pull-all-concise-summary-and-ssh-fleet-json.md)

---

## 1. Symptom

When running `gitmap pa --ssh` (or `gitmap pull-all --ssh`) from a repository or workspace on Windows:
```text
PS C:\Users\Alim\awansoft-v10> gitmap pa --ssh

  Enqueuing 'pull-all' across SSH fleet:
    • Remote Node [w3] (node-w3): Offline (skipped, no task enqueued)
    • Remote Node [w1] (node-w1): Online → Enqueued (async)
    • Remote Node [w2] (node-w2): Online → Enqueued (async)
    • Remote Node [w4] (node-w4): Online → Enqueued (async)
    • Remote Node [u1] (node-u1): Offline (skipped, no task enqueued)
    • Remote Node [main] (node-main): Offline (skipped, no task enqueued)
    • Current Machine [Alim-Desktop (127.0.0.1)]: Running locally (direct execution, not enqueued)

  ▶ Local VM (127.0.0.1 - localhost): 45 pulled (22 active, 23 up-to-date)
      • ...
      • seo-packages-v2                        failed
        ↳ Reason: fatal: Unable to persist credentials with the 'wincredman' credential store.
        ↳ Next Step: Inspect repository: run 'git -C seo-packages-v2 status'
      • slides-spec                            failed
        ↳ Reason: fatal: Unable to persist credentials with the 'wincredman' credential store.
        ↳ Next Step: Inspect repository: run 'git -C slides-spec status'
```

Problems identified:
1. `executeLocalVMPull` in `cli/cmdssh/ssh_pull_fleet.go` bypassed in-process execution, invoking `exec.Command(resolveLocalGitmapExecutable(), "pa", "--json")` as an external OS subprocess. This re-invoked whatever outdated binary was registered on the system PATH, failing to run inside GitMap in-process.
2. When pulling repositories, Git Credential Manager attempted concurrent writes to Windows Credential Manager (`wincredman`), failing multiple repositories with `fatal: Unable to persist credentials with the 'wincredman' credential store.`.

---

## 2. Root Cause

`executeLocalVMPull` spawned an external OS subprocess rather than delegating the call in-process to `cmdpull.RunPullAllJSON` to receive structured JSON output directly, while concurrent git operations on Windows lacked `GCM_NO_PERSIST=1` and `GCM_CREDENTIAL_STORE=cache` environment protections to prevent Git Credential Manager lock contention against the Windows Credential Manager.

---

## 3. Resolution

1. **In-Process JSON Delegation for Local Fleet Pull:**
   - Added `cmdpull.RunPullAllJSON(args []string) (string, error)` which runs `runPull` in-process with `--all`, `--json`, and `--parallel 2`.
   - Introduced `SetPullBatchJSONWriter` and `ResetPullBatchJSONWriter` in `cli/cmdpull/pull_batch_json.go` (and `cli/cmdpull/pull_efficient_json.go`) allowing in-process callers to capture the structured JSON output into an in-memory buffer without redirecting OS pipes.
   - Added `cmdssh.RunLocalPullAllJSONFn` callback hook in `cli/cmdssh/ssh_pull_fleet.go` wired in `cli/cmd/clihelpers.go`.
   - Updated `executeLocalVMPull` to invoke `RunLocalPullAllJSONFn` in-process, parse the JSON payload, and only fall back to subprocess execution if uninitialized.

2. **Git Credential Manager Concurrency Hardening:**
   - Configured `GCM_NO_PERSIST=1` and `GCM_CREDENTIAL_STORE=cache` across all automated git environments in `cli/cloner/safe_pull.go`, `cli/cloner/cloner.go`, `cli/cloner/safe_push.go`, and `cli/cmdpull/pull.go`.
   - This disables disk-persisted credential writes during automated and batch operations, eliminating `wincredman` lock contention.

3. **Core Command Alias Recognition:**
   - Extended `isGitmapCoreCommand` in `cli/cmdssh/ssh_exec_command.go` to recognize `pa`, `pae`, `paet`, `ca`, `cfr`, `cfrp`, and `sa` as first-class GitMap commands for remote SSH delegation.

---

## 4. Prevention & Learnings

1. **Rule: No Out-of-Process Self-Invocation:** Fleet commands running tasks locally on the host machine MUST ALWAYS delegate in-process via package function calls or injected callback hooks; never spawn external `exec.Command("gitmap", ...)` subprocesses.
2. **Rule: GCM Non-Persist Environment:** All headless, batch, parallel, and automated Git invocations MUST inject `GCM_NO_PERSIST=1` and `GCM_CREDENTIAL_STORE=cache` into `cmd.Env` to guarantee Windows Credential Manager stability.
