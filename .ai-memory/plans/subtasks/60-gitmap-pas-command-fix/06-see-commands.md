# Subtask 06: GitMap Unified See Inspection & Remote Telemetry

## Objectives
- Implement and verify `gitmap see` (`c`) suite and aliases:
  - `gitmap see commit pending` / `gitmap c commit pending`: display repos with pending commits.
  - `gitmap see git-ignore issues` / `gitmap c ig issues`: inspect ignore duplicates and tracked ignore issues.
  - `gitmap see errors` / `gitmap c errors`: inspect local GitMap Errors journal.
  - `gitmap see errors ssh` / `gitmap ses`: query remote fleet node errors following GitMap PAS Formula.
  - `gitmap history ssh` / `gitmap nodes history`: view remote execution task history.
  - `gitmap repo-manage ui`: launch interactive terminal dashboard.

## Target Files
- `cli/cmdsee/see.go`
- `cli/cmd/rootcore.go`
- `cli/cmd/rootdata.go`
