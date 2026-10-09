# Subtask 01 — benchmark-first-line (task 251, gitmap)

## File box
`cli/cmdagent/agent_collisions.go`, `agent_stats.go`, `agent_heatmap.go`,
`agent_ps.go` — timing wiring ONLY. No logic changes.

## Build (spec 01 §A2)
1. Each command measures its own work with `time.Now()` at entry and prints
   `[Nms]` as the VERY FIRST line before any other output (tables included).
2. `--json` mode: `elapsed_ms` is the first key in the payload.
3. Keep it simple: `start := time.Now()` at the top of each `Run*`, compute at
   each return point (or defer-print — but must be FIRST line, so print at the
   end won't work; restructure to print timing first, then the report).

## Rules
GitMap tools only; no rg/grep. `go build ./...` only, never `go test`.
Never `git add`/`git commit`.

## Deliverable
Diff + build exit code. Then stop.
