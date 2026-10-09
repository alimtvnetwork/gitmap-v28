# Subtask 04 — version truth unification (task 250)

## File box
`version.json`, its readers, release-gate wiring. Nothing else.

## Build (per spec 01 §5)
1. `version.json`: exactly ONE version field (`Version`). Remove the lowercase
   duplicate. Enumerate EVERY reader of both fields (gitmap aum search) and migrate
   them to the canonical field. List the migrated readers in the report.
2. Write rule: human edits ONLY `version.json` → `Version`; `37-bump-version.py`
   remains the sole writer for the rest. Document in the spec-adjacent code comment.
3. Enforcement: wire `03-ai-scripts/14-version-sync-checker.py` into the pre-tag
   release gate path so a divergent copy refuses the tag. (If the gate lives in
   `29-release-orchestrator.py`, add the check there; minimal diff.)
4. `cli/constants/constants.go`: keep ldflags override; default stays script-synced
   (no hand edits).

## Rules
GitMap tools only; no rg/grep. `go build ./...` (+ `gitmap py` for the checker),
never `go test`.

## Deliverable
Diff + reader migration list + build exit code. Then stop.
