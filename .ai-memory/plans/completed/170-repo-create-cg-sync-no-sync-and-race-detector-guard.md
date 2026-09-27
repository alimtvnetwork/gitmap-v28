# Plan 170: Repo Create Coding Guideline Auto-Sync, No-Sync Flag & Race Detector Guard

## Status
`completed`

## Summary of Completed Implementation
- **Coding Guideline Auto-Sync (`--cg`):**
  - Updated `resolveCGBaseDir()` in `cli/cmd/repo_create_init.go` to search canonical coding guidelines across candidate repositories (`coding-guidelines`, `gitmap`), environment variables, parent directory walks from cwd, and drive roots.
  - Implemented `tryCopyCGSpecFiles(baseDir, cgDst)` to copy the complete `02-spec/02-coding-guidelines` folder into the target repository.
  - Added remote installer fallback via `cmdcg.RunCgScriptInRepo(absDir)`.
- **Selective Sync (`--no-sync`):**
  - Added `IsNoSync bool` to `createRepoParams` in `cli/cmd/repo_create_params.go`.
  - Added `workspacesync.SyncAgyOnly` in `cli/workspacesync/workspacesync.go` to skip GitHub Desktop and VS Code while registering Antigravity.
  - Wired `applyWorkspaceSync` in `cli/cmd/create_ops.go`.
- **Race Condition & Linter Remediation:**
  - Synchronized `cloneTargets` slice in `cli/cmd/commitin/workspace/workspace_test.go` with `sync.Mutex`.
  - Ran `gofmt -w` on `cli/cmdpull/pull_flags_test.go`.
  - Verified with `test_ci_scripts.py` (18/18 tests pass).
- **Release Ceremony:**
  - Published minor release `v6.349.0` via `python 03-ai-scripts/29-release-orchestrator.py -t minor --skip-tests`.
  - Bumped `constants.Version` to `6.349.0` and compiled binary to local AppData installation.
