# Subtask 01: Windows Git Subprocess Credential Store Fix

## Scope
- Files:
  - `cli/cloner/cloner.go`
  - `cli/cloner/safe_pull.go`
  - `cli/cloner/safe_push.go`
  - `cli/cmdclone/clonepretty.go`
  - `cli/cmdpull/pull.go`
- Objective:
  - Strip `GCM_CREDENTIAL_STORE=cache` on Windows across all Git subprocess builders.
  - Implement OS-aware GCM environment builder or helper:
    - Retain `GCM_NO_PERSIST=1`, `GCM_INTERACTIVE=never`, `GIT_TERMINAL_PROMPT=0`.
    - If `runtime.GOOS != "windows"`, allow `GCM_CREDENTIAL_STORE=cache`.
    - If `runtime.GOOS == "windows"`, do NOT inject `GCM_CREDENTIAL_STORE=cache`.
