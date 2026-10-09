# Plan 251 — collision benchmark, agent ls/show, Antigravity review

## Scope (spec: `02-spec/21-app/251-collision-detection-and-antigravity-manager-review/`)
1. gitmap: benchmark-first-line on `agent collisions|stats|heatmap|ps` (+`elapsed_ms` in JSON).
2. gitmap: `agent ls` (last N tasks, IDs, statuses, subtask tree) + `agent show <id>` (detail).
3. AGM: openai.rs attribution (Jeik = upstream lbjlaq or fork?), 7.5k-line rescue plan.
4. AGM: PR #6 "Sync v5" review — extract ideas, record commits, never break our code.
5. AGM: scope tier map + `unwrap()`/`expect()` reduction in proxy handlers.
6. Release gitmap 6.519.0; honest `gitmap update` report.

## Execution
- 4 workers, disjoint boxes. W1+W2: gitmap (`cli/cmdagent/` new files only).
  W3: AGM investigation (read-only, `~/workspace/repos/Antigravity-Manager`).
  W4: AGM scope/improvements (AGM repo).
- Lead: build verify, e2e, release.
- `go build ./...` only; no `go test`. Targeted staging via `gitmap cpf --task`.

## Status
**IN PROGRESS** — workers dispatched.
