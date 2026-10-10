# Subtask 04: Workspace Heatmap TreeView & CI/CD Telemetry Fusion
Parent Task: gitmap-command-enhancements
Status: COMPLETED

## Objective
Implement `gitmap full summary [N]`, `gitmap full status [N]`, `gitmap fs [N]`, `gitmap fs+pe`, and `gitmap pe all` with 48h active filter, dirty worktree breakdown, master sanitize commit command, and pipeline failure stack traces.

## Target Files
- `cli/cmdsummary/full_summary_core.go`
- `cli/cmdsummary/full_summary_pe.go`
- `cli/cmdsummary/summary_pe_all.go`

## Verification
Unit tests pass (`TestParseFullSummaryOptions`).
Live CLI smoke test: `.\cli\gitmap.exe fs 2` renders TreeView with pending changes, releases, and master sanitize command.
`.\cli\gitmap.exe pe all --json` emits structured JSON telemetry.
