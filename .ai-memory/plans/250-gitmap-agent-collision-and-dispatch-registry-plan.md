# Plan 250 — agent collision handling, dispatch registry, A08, version truth

## Scope (spec: `02-spec/21-app/250-gitmap-agent-collision-and-dispatch-registry-plan/`)
1. **Dispatch registry** — central name/alias aggregation + `gitmap doctor --check-dispatch`
   + collision test. No dispatch flow refactor.
2. **Agent collision feature** — file-box claims in the task DB (`OwnedFilesJson`),
   `gitmap agent collisions` report, look-ahead before dispatch, scoped staging in
   the commit flow (`--task <slug>` stages only claimed files; agents never `git add`).
3. **A08** — `writeLastErrorFile` discards become visible warnings (+ audit).
4. **Version truth** — single `Version` field; bump script sole writer; sync check in gate.
5. **Release 6.518.0** via `37-bump-version.py`; tag `v6.518.0`.

## Execution (owner lifted the review gate — implementation starts immediately)
- 5 workers, disjoint file boxes (subtasks 01–05). Lead: build verification, e2e,
  release (bump + tag + push).
- Mindset: overlaps reported, never hard-blocked; end-of-task verify-and-fix.
- `go build ./...` only; no `go test`. Targeted staging only.
- Subtask 05 (new): agent observability methods — `stats`, `heatmap`, `ps`,
  collision history (spec 01 §3.5).

## Acceptance
D01–D03, C01–C05, A01–A03, V01–V03 (spec 02). Tag cut only after functional IDs pass.

## Status
**AWAITING OWNER REVIEW** — no implementation until the plan is reviewed.
