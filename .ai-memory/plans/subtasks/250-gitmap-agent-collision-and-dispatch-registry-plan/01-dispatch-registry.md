# Subtask 01 — dispatch registry (task 250)

## File box
New package `cli/cmddispatch/` ONLY (+ test file). Read-only elsewhere. Do NOT
refactor the dispatch flow.

## Build (per spec 01 §2)
1. `cli/cmddispatch/registry.go`: `CollectNames() []DispatchName` — aggregate every
   claimed command/alias from each dispatch table (core entries `cli/cmd/rootcore.go`,
   `dispatchAgySubsystem` cases in `cli/cmd/root.go`, extra/general commands), each
   tagged with owner file:line. `FindCollisions() []Collision` — names claimed ≥2.
2. Wire `gitmap doctor --check-dispatch`: print all names with owners; exit 1
   listing both owners on any collision.
3. Go test `TestDispatchRegistryNoCollisions` (same assertion; synthetic duplicate
   injected in-test only, never in production tables).
4. Regression: the check must catch `"agm"` if re-added to gitignore aliases
   (prove via the synthetic-duplicate test shape).

## Rules
GitMap tools only; no rg/grep. Small files, ≤300 lines. `go build ./...` only,
never `go test` (the test file is written but NOT executed — standing rule).

## Deliverable
Diff + build exit code. Then stop.
