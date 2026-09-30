# Learned 35: In-Process Fleet Pull JSON Delegation and Wincredman Concurrency Isolation

## Context
When running `gitmap pa --ssh`, the host machine was running `exec.Command("gitmap", "pa", "--json")` out-of-process and multiple repos failed on Windows with `fatal: Unable to persist credentials with the 'wincredman' credential store.`.

## Patterns & Conventions
1. **Host Fleet Operations In-Process:**
   - Any fleet command executed on the local host machine (`127.0.0.1 - localhost`) must be executed directly in-process via package function calls or injected callback hooks (e.g. `cmdssh.RunLocalPullAllJSONFn = cmdpull.RunPullAllJSON`).
   - Output must be routed into in-memory writers (`SetPullBatchJSONWriter(&buf)`) rather than spawned subprocesses or OS pipe redirects.

2. **Git Credential Manager Headless Safety:**
   - In automated, batch, or parallel Git operations on Windows, Git Credential Manager (`wincredman`) must be instructed not to persist credentials to the disk vault to avoid concurrent `CredWriteW` lock collisions:
     ```go
     "GCM_NO_PERSIST=1",
     "GCM_CREDENTIAL_STORE=cache",
     ```
   - Must be present on `cmd.Env` across `safe_pull.go`, `cloner.go`, `safe_push.go`, and `pull.go`.
