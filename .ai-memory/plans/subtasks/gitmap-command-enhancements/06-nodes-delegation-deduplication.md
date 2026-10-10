# Subtask 06: Fleet Nodes Delegation with Local Precedence
Parent Task: gitmap-command-enhancements
Status: COMPLETED

## Objective
Implement `gitmap nodes fs`, `nodes fspe`, `nodes summary $repo`, and `nodes pe all` with two-phase handshake and local-machine precedence deduplication (if local machine has the repo, do not query remote nodes for it).

## Target Files
- `cli/cmdnodes/nodes_summary.go`
- `cli/cmdnodes/nodes_cmd.go`

## Verification
Unit tests pass: `go test ./cmdnodes/...` (4.315s).
Live CLI smoke test: `.\cli\gitmap.exe nodes fs` falls back gracefully when no remote SSH nodes exist and renders local summary.
