# Subtask 02 — agent collision feature (task 250)

## File box
`cli/cmdagent/` (claim/collisions CLI) + `cli/cmdcommit/commit_push.go`
(`executeStageAllStep` scoping ONLY). Nothing else.

## Build (per spec 01 §3)
1. **Look-ahead claim**: `gitmap agent subtask claim --task-id <slug> --files <rel,paths…>`
   (also `--files` on `subtask add`). Persist to `Subtask.OwnedFilesJson` (column
   exists in `cli/cmdagent/agent_config.go`). Extend `subtask ls` output with files.
2. **Collisions report**: `gitmap agent collisions [--parent <slug>]` — overlapping
   write claims across active subtasks, with owners and statuses. Overlaps are
   REPORTED, never blocked.
3. **Scoped staging**: `executeStageAllStep()` — when `--task <slug>` is passed to
   `cpf/cpb/cpc/cpr` (or `GITMAP_TASK` env set), stage ONLY the union of that task's
   `OwnedFilesJson` paths (`git add -- <paths>`) instead of `git add -A`. Without
   `--task`, behavior unchanged. Agents must never run `git add`/`git commit` directly.
4. Communication stays DB-based (`AgentActionLog.TargetFile` already exists).

## Rules
GitMap tools only; no rg/grep. Small files. `go build ./...` only, never `go test`.

## Deliverable
Diff + build exit code. Then stop.
