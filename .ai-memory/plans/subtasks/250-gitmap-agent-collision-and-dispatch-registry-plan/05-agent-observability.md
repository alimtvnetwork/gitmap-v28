# Subtask 05 — agent observability methods (task 250)

## File box
NEW files only under `cli/cmdagent/`: `agent_stats.go`, `agent_heatmap.go`,
`agent_ps.go` (one per method; exact names — do not create other files). Read-only
elsewhere. Shares the §3.5 schema contract (`FileClaim`, `CollisionEvent` in the
central registry DB; confirm the central DB path — `ParentTaskRegistry` location
is the lead candidate).

## Build (per spec 01 §3.5)
1. `gitmap agent stats [--days N]`: recently completed tasks, subtasks completed,
   per-agent completion counts, total collisions (from `CollisionEvent`).
2. `gitmap agent heatmap [--parent <slug>]`: files ranked by write-claim count
   (`FileClaim` + `AgentActionLog.TargetFile`).
3. `gitmap agent ps`: running tasks — slug, status, active agent count, subtask
   progress rollup.
4. All read-only; `--json` supported like the other agent commands. Plain,
   professional terminal output (no raw dumps).

## Rules
GitMap tools only; no rg/grep. Small files, ≤300 lines each. `go build ./...`
only, never `go test`.

## Deliverable
Diff + build exit code. Then stop.
