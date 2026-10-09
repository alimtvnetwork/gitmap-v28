# 247 — Slug-based task ID for the agent engine (scope and plan)

## 1. Problem statement

The agent task engine (`gitmap agent task|subtask`, cobra, 3-tier SQLite) already
stores a `TaskSlug` on parent tasks, but the slug is not usable as a task ID:

- Re-running `task init` with the same name creates a **duplicate** row
  (`TaskSlug` has no UNIQUE constraint; no existence check).
- There is **no get-or-create**: no single command that takes a slug and either
  enqueues the task (new) or reports progress (existing).
- Subtasks have **no slug column** — they cannot be addressed or verified by slug.
- There are **no query commands** for: pending-task count, recent tasks,
  completed root tasks.

The Letterly/Muse prompt pipeline now emits a `## slug` (Title Case task ID) on
every formatted task and a `Slug:` sub-item on every Task-NN. GitMap must accept
those slugs as IDs: enqueue-or-check by slug, create subtasks by slug, and answer
"how many pending / recent / completed" queries.

## 2. Target engine

The cobra engine in `cli/cmdagent/` (NOT the legacy `gitmap task` flag engine in
`agent_task_sqlite.go`, NOT the legacy queue in `cli/cmdtask/task.go`):

- `cli/cmdagent/agent_task.go` — `initTaskCommands()` (`agent_task.go:73`)
- `cli/cmdagent/agent_subtask.go` — `initSubtaskCommands()` (`agent_subtask.go:130`)
- `cli/cmdagent/agent_config.go` — schema DDL (`InitTier1Schema`, `InitTier2Schema`)
- `cli/store/agent_store.go` — store layer
- `cli/store/split_db_path.go` — `SanitizeSlug` (lowercase, non-alnum → `-`)

Verified 2026-10-09: `SanitizeSlug` exists as described; `TaskCmd.AddCommand` /
`SubtaskCmd.AddCommand` registration points confirmed.

## 3. Design

### 3.1 New commands (`gitmap agent task`)

| Command | Behavior |
|---|---|
| `enqueue --slug <slug> --name <title> [--budget N] [--json]` | Sanitize slug via `store.SanitizeSlug`. Look up Tier 1 `ParentTaskRegistry` by sanitized slug. **Found** → print existing task + progress rollup (the "check" path; exit 0). **Not found** → init new parent task exactly as `task init` does (the "enqueue" path). Idempotent: same slug twice never duplicates. |
| `progress --slug <slug> [--json]` | Parent row + subtask rollup (pending/in-progress/done/failed counts, done/total) + the 5 most recent *other* parent tasks ("related previous tasks"). Slug-first alternative to `task status --task-id`. |
| `pending [--count] [--task-id <id-or-slug>] [--json]` | Pending subtasks across all agent tasks (or one task). `--count` prints just the number. Answers "how many pending tasks are there in the agents". |
| `recent [--limit N] [--json]` | Parent tasks ordered by `CreatedAt` desc. |
| `completed [--limit N] [--json]` | Parent tasks with `Status='COMPLETED'` (root-level completions). Thin, discoverable alternative to `task ls --status COMPLETED`. |

### 3.2 Subtask slugs

- `gitmap agent subtask add` gains `--slug <slug>`: accepts the Title Case form
  from the prompt (e.g. `SEO Writing and Folder Structure Instructions - Task 01`),
  stores `SanitizeSlug(slug)` in a new `Subtask.TaskSlug` column, keeps `Title`
  as given.
- Schema migration: `ALTER TABLE Subtask ADD COLUMN TaskSlug TEXT NOT NULL
  DEFAULT ''` applied when missing (check `PRAGMA table_info(Subtask)` first —
  SQLite has no `ADD COLUMN IF NOT EXISTS`). New-DB DDL in `agent_config.go`
  includes the column from the start.
- If `--slug` omitted, derive from parent slug + code:
  `SanitizeSlug(parentSlug + "-" + code)`.

### 3.3 Store layer (`cli/store/agent_store.go`)

- `FindParentTaskBySlug(masterDB, slug) (*ParentTask, error)` — full-row lookup
  by sanitized slug (existing `scanRootDbByTask` returns only the DB path).
- `EnsureSubtaskSlugColumn(db) error` — idempotent migration.
- Extend subtask insert to persist `TaskSlug`.

### 3.4 Explicitly out of scope

- UNIQUE constraint migration on existing `ParentTaskRegistry.TaskSlug`
  (code-level existence check in `enqueue` covers idempotency on all DBs).
- The legacy `gitmap task` flag engine and the legacy queue (`create/list/run/…`).
- Changing `SanitizeSlug` behavior (Title Case input → sanitized storage is the
  documented contract; `TaskName`/`Title` preserve original case).

## 4. Implementation plan

1. Spec (this dir) + plan task entry.
2. Store: `FindParentTaskBySlug`, `EnsureSubtaskSlugColumn`, subtask-insert slug.
3. Schema DDL: `TaskSlug` in `Subtask` CREATE TABLE (new DBs).
4. Commands: new file `cli/cmdagent/agent_task_slug.go` (<300 lines) with
   `enqueue`, `progress`, `pending`, `recent`, `completed`; extend
   `subtask add` with `--slug`; register in `initTaskCommands` /
   `initSubtaskCommands`.
5. Build (`go build ./...`), then E2E in a temp dir: init → enqueue (new) →
   enqueue (existing → check path) → subtask add with slug → pending --count →
   progress → recent → completed.
6. Skills: document new commands in `.agents/skills/gitmap/SKILL.md`
   (agent-task-engine section) and any dedicated agent skill file.
7. Commit + push (atomic, hyphen-format).
