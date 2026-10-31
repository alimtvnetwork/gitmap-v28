# Subtask: TASK-06 - Compilation, Live Terminal Verification & Release Ceremony (v6.645.0)

- **Parent Task:** `370-ai-execution-and-error-handling-instructions`
- **Subtask ID:** `sub-121929-task-06`
- **Owned Files:** `cli/constants/constants.go`, `version.json`, `changelog.md`

## Objectives
1. Bump version to `v6.645.0` in `cli/constants/constants.go` and `version.json`.
2. Add detailed entry to `changelog.md`.
3. Compile with `go build -p 2 -ldflags=-s -o ../bin/gitmap.exe .` and deploy to local appdata.
4. Execute live terminal tests verifying `gitmap ps <dir>` intercepts cleanly with AI reflection card.
5. Move plan to completed and commit via `gitmap cpr "v6.645.0 ai execution and powershell path error handling instructions"`.
