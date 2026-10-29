# Subtask 01: Non-Git Error Handling & Web UI Dashboard Fallback

- **Status**: completed

## Objectives
- Update `RunShowLog` in `cli/cmdlogviewer/server.go`:
  - Extract helper `handleNonGitShowLog(args []string, isNoBrowser bool) error`.
  - Deepen child repository discovery to probe level 2 directories when level 1 has no git repos.
  - Skip noise/system directories (`node_modules`, `vendor`, `.git`, `AppData`, `Windows`, etc.).
  - Always launch `FallbackDashboardUIFn("search", 8080)` when `!isNoBrowser && FallbackDashboardUIFn != nil`.
- Update `runSinglePathDiff` in `cli/cmddiff/single_diff.go`:
  - When target is not a Git repo, emit friendly notice and tip, and return `nil` rather than an execution error.

## Verification
- Invoke `gitmap show-log` from `D:\` and verify it discovers child repos and launches dashboard.
- Verified: Detected 10 child repos in `work/` and exited cleanly with code 0.
