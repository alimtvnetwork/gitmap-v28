# Spec 250 — Acceptance criteria and evidence

## Dispatch registry

- **D01**: `gitmap doctor --check-dispatch` exits 0 on the current tree and lists
  every registered command/alias with its owner. Evidence: command output.
- **D02**: a synthetic duplicate alias (test-only, not in production tables) makes
  the check exit 1 naming both owners. Evidence: test output.
- **D03**: re-adding `"agm"` to the gitignore alias list is caught (regression for
  the 247 hijack). Evidence: check output naming `rootcore.go` and `root.go`.

## Agent collision feature

- **C01**: `gitmap agent subtask claim --task-id <slug> --files a.go,b.go` persists
  the box; `gitmap agent subtask ls` shows it. Evidence: CLI transcript.
- **C02**: two active subtasks claiming the same file → `gitmap agent collisions`
  reports the overlap with both owners. Evidence: CLI transcript.
- **C03**: disjoint boxes → `collisions` reports clean. Evidence: CLI transcript.
- **C04**: `gitmap cpf --task <slug>` stages ONLY that task's claimed files —
  e2e with two tasks (overlapping + disjoint files): each commit contains exactly
  its task's files, sibling files stay uncommitted. Evidence: `git show --stat`
  of both commits.
- **C05**: without `--task`, `cpf` behavior is unchanged (full `git add -A`).
  Evidence: e2e transcript.

## A08 swallowed errors

- **A01**: with CWD read-only, a failing command prints the concise persistence
  warning (no silence). Evidence: e2e transcript.
- **A02**: normal path unchanged — `.gitmap/last_error.log` still written on errors.
  Evidence: e2e transcript.
- **A03**: `persistToErrorsDB` audit — any discards found are fixed the same way;
  none remain on the error path. Evidence: `aum search` for `_ =` on the path.

## Version truth

- **V01**: `version.json` has exactly one version field (`Version`); no reader
  references the removed duplicate. Evidence: `aum search` for `"version"` (lowercase
  key) returns zero code hits.
- **V02**: bumping via the script keeps every copy in sync; the sync checker passes.
  Evidence: dry-run + checker output.
- **V03**: release 6.518.0 — `version.json` = 6.518.0, annotated tag `v6.518.0`
  (`Release v6.518.0`) on origin, changelog entry, `gitmap version` = 6.518.0.
  Evidence: tag listing + version output.

## Evidence locations

E2E scripts + logs OUTSIDE the repo (`~/workspace/gitmap-e2e-250/`). Nothing
test-related committed. `go test` never runs (standing rule).

## Failure policy

Any FAIL → rejected, root-caused, re-verified. The 6.518.0 tag is cut ONLY after
all functional IDs pass.

## Agent observability methods

- **S01**: `gitmap agent stats` shows recently completed tasks, per-agent completion
  counts, and total collisions. Evidence: CLI transcript.
- **S02**: `gitmap agent heatmap` ranks files by write-claim count. Evidence: CLI
  transcript.
- **S03**: `gitmap agent ps` lists running tasks with slug, status, and active agent
  count. Evidence: CLI transcript.
- **S04**: overlapping claims write `CollisionEvent` rows; `stats` reflects them.
  Evidence: CLI transcript + row count.
