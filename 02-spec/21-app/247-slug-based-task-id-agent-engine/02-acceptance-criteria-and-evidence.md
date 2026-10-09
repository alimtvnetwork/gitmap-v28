# 02 — Acceptance Criteria and Evidence (task 247: slug-based task ID agent engine)

## Commands under test

All run against a temp `XDG`/data dir (isolated from the user's real agent DBs).

### AC1 — `gitmap agent task enqueue` (get-or-create by slug)
- `enqueue --slug "SEO Writing Task" --name "SEO Writing Task"` creates a parent
  task; stored `TaskSlug` = `seo-writing-task` (sanitized), `TaskName` keeps
  Title Case.
- Running the same `enqueue` again does NOT create a duplicate: second run takes
  the "check" path and prints the existing task + progress. `task ls` shows one row.
- `--json` emits machine-readable output in both paths.

### AC2 — `gitmap agent task progress --slug <slug>`
- Prints parent status + subtask rollup (pending/in-progress/done/failed,
  done/total) + up to 5 most recent *other* parent tasks.
- Unknown slug → clean "not found" error, exit non-zero, no stack trace.

### AC3 — `gitmap agent task pending`
- With 2 pending + 1 done subtask across tasks: `pending --count` prints `2`;
  without `--count` lists the 2 pending subtasks with their parent slugs.
- `--task-id <slug>` scopes to one task.

### AC4 — `gitmap agent task recent` / `completed`
- `recent --limit 2` returns the 2 most recently created parent tasks, newest first.
- After completing one task, `completed` lists it (root-level completions).

### AC5 — subtask slugs
- `gitmap agent subtask add --parent <slug> --slug "SEO Writing Task - Task 01" --code task-01 --title "…"`
  stores sanitized `TaskSlug` (`seo-writing-task-task-01`); re-running `add`
  against an older DB without the column migrates it (no error, column present).
- Omitting `--slug` derives it from parent slug + code.

### AC6 — no regressions
- `go build ./...` exit 0; existing `agent task|subtask` commands (`init`, `ls`,
  `status`, subtask `add/claim/start/complete/fail/ls`) behave as before
  (spot-check `init` → `ls` → `status`).

## Evidence to record on completion

- Build log (`go build ./...` exit code).
- E2E transcript (temp-dir run covering AC1–AC5).
- `git log` SHA of the implementing commit.
