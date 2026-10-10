# Subtask 07: CLI Root Routing, Constants & Multi-Word Rewrites
Parent Task: gitmap-command-enhancements
Status: COMPLETED

## Objective
Register CLI constants, dispatch routes, and multi-word command rewrites so that `gitmap full summary`, `gitmap full status`, `gitmap fs`, `gitmap fs+pe`, `gitmap fspe`, `gitmap pe all`, `gitmap merge-ai`, and `gitmap ma` dispatch cleanly without errors.

## Target Files
- `cli/constants/constants_cli.go`
- `cli/cmd/roottooling.go`
- `cli/cmd/rootutility.go`
- `cli/cmd/root.go`

## Verification
Binary compiles cleanly (`go build -o gitmap.exe .`).
Root command tests pass: `go test ./cmd/...` (PASS).
