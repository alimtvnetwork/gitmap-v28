# Issue 22: SSH Fleet Pull Subprocess Bypass and Wincredman Credential Store Failure

**Status:** Resolved
**Date:** 2026-10-01
**Related Spec:** [02-spec/22-app-issues/56-ssh-fleet-pull-subprocess-bypass-and-wincredman-credential-failure-rca.md](../../02-spec/22-app-issues/56-ssh-fleet-pull-subprocess-bypass-and-wincredman-credential-failure-rca.md)

## 1. Symptom
Running `gitmap pa --ssh` locally invoked an external `gitmap.exe` subprocess rather than delegating in-process to the pull engine, and multiple concurrent git pulls failed with `fatal: Unable to persist credentials with the 'wincredman' credential store.`.

## 2. Root Cause
1. `executeLocalVMPull` in `cli/cmdssh/ssh_pull_fleet.go` executed `exec.Command(resolveLocalGitmapExecutable(), "pa", "--json")` instead of delegating in-process to `cmdpull.RunPullAllJSON`.
2. Windows Credential Manager (`wincredman`) experienced lock contention during parallel pulls without `GCM_NO_PERSIST=1` and `GCM_CREDENTIAL_STORE=cache`.

## 3. Resolution
1. Created `cmdpull.RunPullAllJSON` with in-memory buffer routing via `SetPullBatchJSONWriter`.
2. Wired `cmdssh.RunLocalPullAllJSONFn = cmdpull.RunPullAllJSON` in `cli/cmd/clihelpers.go`.
3. Injected `GCM_NO_PERSIST=1` and `GCM_CREDENTIAL_STORE=cache` across all automated git pull/push/clone commands.
4. Added aliases `pa`, `pae`, `paet`, `ca`, `cfr`, `cfrp`, `sa` to `isGitmapCoreCommand`.

## 4. Prevention
Always delegate host execution in-process inside GitMap, and mandate `GCM_NO_PERSIST=1` across all batch Git operations.
