# Subtask 02 — agent ls + agent show (task 251, gitmap)

## File box
NEW files only: `cli/cmdagent/agent_ls.go`, `cli/cmdagent/agent_show.go`.
Read-only elsewhere.

## Build (spec 01 §A3–A4)
1. `gitmap agent ls [--limit N] [--json]` (default 12, max 50): last N parent
   tasks from Tier-1 `ParentTaskRegistry` ordered by UpdatedAt DESC. Columns:
   ID | SLUG | STATUS | AGENTS | SUBTASKS done/total | UPDATED.
2. Subtask subtree: under each parent, indented rows from its Tier-2 DB
   (`sub-<id> | <code> | <status> | <agent>`). Missing/broken Tier-2 → show
   `BROKEN` marker, never crash.
3. `gitmap agent show <id-or-slug> [--json]`: full parent detail (ID, slug, name,
   status, run dir, budget, created/completed at) + subtask subtree with file-box
   counts and evidence excerpt.
4. Canonical ID-or-slug resolution (follow existing `resolveSubtaskParent` patterns).
5. All SQLite, no network. Target <300ms for 15 tasks (the B01 benchmark line
   on `ls` proves it — add the timing first-line here too).

## Rules
GitMap tools only; no rg/grep. ≤300 lines per file. `go build ./...` only,
never `go test`. Never `git add`/`git commit`.

## Deliverable
Diff + build exit code. Then stop.
