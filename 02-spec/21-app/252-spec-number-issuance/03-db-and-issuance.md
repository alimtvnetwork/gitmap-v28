# Spec 252 — Spec-number issuance — 03 — DB and issuance

Program: 252 — spec number issuance (concurrency-safe, serialized via
Tier-1 DB). This file owns the DB layer contract and the issuance
algorithm. Command-layer policy lives in the command layer owned by
`cli/cmdspec/` (see subtask 02-command-layer.md); the dumb DB primitive
lives in `cli/store/spec_numbers.go` (see subtask 01-db-layer.md).

## 3.1 Exact DDL — `spec_numbers`

```sql
CREATE TABLE IF NOT EXISTS spec_numbers (
  number     INTEGER NOT NULL,
  repo_root  TEXT    NOT NULL DEFAULT '',
  task_slug  TEXT    NOT NULL DEFAULT '',
  issued_at  TEXT    NOT NULL,
  UNIQUE(number)
);
```

Column semantics:

- `number` — the issued spec sequence number. `UNIQUE(number)` is the
  race-safety anchor (see §3.4).
- `repo_root` — absolute path of the repo that issued the number, stored
  for auditability. The Tier-1 DB file itself lives at
  `<repo>/.ai-memory/temp-agents/ai_agents.db`, so per-repo scoping is
  already inherent; `repo_root` exists so a shared or copied DB file can
  still be traced back to its origin. Never used in uniqueness logic.
- `task_slug` — the program slug that claimed the number (e.g.
  `252-spec-number-issuance`), for traceability.
- `issued_at` — RFC3339 UTC timestamp of the successful claim.

Conventions (matching the repo's existing schema code):

- Package-level DDL const, e.g. `sqlCreateSpecNumbersTable`, executed via
  `store.ExecWrapper` inside a loop over `[]string{...}` at every DB open
  — copy the `ensureCollisionSchema` shape
  (`cli/cmdagent/agent_claim.go:17-52`).
- Table name is lowercase (`spec_numbers`), consistent with the snake_case
  collision tables (`FileClaim`, `CollisionEvent` are CamelCase legacy;
  new tables follow the lowercased `spec_numbers` form).

## 3.2 The `cli/store/spec_numbers.go` contract

The DB primitive is deliberately dumb — no policy, no disk access, no
spec-directory knowledge. Three exported functions, mirroring existing
store helpers:

```go
// EnsureSpecNumbersSchema creates spec_numbers when missing (idempotent,
// runs at every open via the DDL loop).
func EnsureSpecNumbersSchema(conn *sql.DB) *appfault.AppError

// ClaimSpecNumber inserts (number, repoRoot, taskSlug, now) with
// INSERT OR IGNORE. Returns (true, nil) when this caller won the claim,
// (false, nil) when the number was already taken, and an error only on a
// genuine DB failure.
func ClaimSpecNumber(conn *sql.DB, number int, repoRoot, taskSlug string) (bool, error)

// MaxIssuedSpecNumber returns the highest number ever issued through this
// DB, or 0 when the table is empty (fresh repo or pre-issuance state).
func MaxIssuedSpecNumber(conn *sql.DB) (int, error)
```

Error model: `ClaimSpecNumber` distinguishes "number already taken"
(normal, returns `false`) from "database broken" (error propagates).
Callers must never treat `RowsAffected == 0` as an error — that is the
normal "I lost the race, try the next number" signal.

## 3.3 Issuance algorithm (policy — lives in `cli/cmdspec/`)

```
1. diskScan = max leading-NNN prefix found under 02-spec/21-app/ (§3.5),
              0 when the directory is empty or missing.
2. dbMax    = store.MaxIssuedSpecNumber(conn), 0 when the table is empty.
3. candidate = max(diskScan, dbMax) + 1
4. loop:
     claimed, err = store.ClaimSpecNumber(conn, candidate, repoRoot, slug)
     if err        → abort with the DB error (something is genuinely broken)
     if claimed    → DONE; print candidate
     candidate     → candidate + 1, retry (bound: 1000 attempts, then abort)
5. The printed number is the caller's number; its directory is created by
   the caller afterwards (DB claim precedes mkdir, so a failed mkdir never
   leaves a phantom hole — the number simply stays claimed).
```

Why the retry loop exists: two agents can compute the same candidate
from the same `max+1` inputs before either inserts. The loser of that race
does NOT receive a new number handed to it — it simply increments and
tries again. There are no gaps by design: consecutive callers get
consecutive numbers unless a raced loser skipped one, which is the
correct, unavoidable outcome of serialization.

## 3.4 Race-window analysis — why INSERT OR IGNORE + retry is correct

- SQLite in WAL mode serializes writers at the page level; with
  `busy_timeout = 5000` and `MaxOpenConns(1)` (per `cli/store/sqlite_config.go`
  and `InitMasterAgentDB`), a second concurrent writer blocks briefly and
  then executes — it does NOT get a hard `SQLITE_BUSY` for the normal
  window. The `UNIQUE(number)` constraint is the single arbiter of
  ownership.
- The only way two claims for the same number can overlap is for both
  writers to pass through the constraint check — impossible: the second
  writer's `INSERT OR IGNORE` evaluates the unique index after the first
  writer's commit is visible, finds the conflict, and silently inserts
  zero rows. No error, no exception path, no lock table of our own.
- `RowsAffected == 0` is therefore an exact, trustworthy "already taken"
  signal. There is no TOCTOU window: no separate SELECT-then-INSERT, no
  advisory lock that could drift out of sync with the table.
- Failure modes: a crash between INSERT and the caller printing the
  number burns that number (it stays claimed). This is acceptable — holes
  are harmless; duplicates are not. A crash before INSERT leaves no trace.
- This is the same primitive the codebase already relies on:
  `INSERT OR IGNORE` at `constants/constants_bookmark.go:37` and
  `cli/cmdpipeline/pipeline_recorder.go:154`.

## 3.5 Disk-scan rule

- Root: `<repo>/02-spec/21-app/`.
- For each directory entry, parse the leading numeric prefix with the
  pattern `^(\d+)-`. Entries not matching (e.g. `01-cli-architecture`,
  wait — `01-...` DOES match; entries like `10-github-desktop.md` files,
  `.md` files, non-numeric names) — match only when the entry is a
  DIRECTORY whose name starts with digits followed by `-`. Skip files.
- `diskScan = max` of all parsed prefixes. Current reality check
  (2026-10-09): highest on disk is `251`, duplicates exist (`240-*` twice,
  `250-*` twice — manual numbering collisions this program exists to
  eliminate), which is why the algorithm takes the max and claims via the
  DB rather than trusting disk order.
- Empty or missing directory → `diskScan = 0`, so the first issuance on a
  fresh repo is `1`. Verified against repo convention: the earliest spec
  dirs in this repo start at low numbers (`01-`, `02-`, ...), and `252`
  is the current max — the sequence began at 1 and the formula
  `max(diskScan, dbMax) + 1` reproduces it exactly.

## 3.6 Future-proofing (note only — do NOT implement now)

If a later program needs extra columns on `spec_numbers` (e.g. a
`released` flag for number reclamation), follow the
`ensureSubtaskSlugColumnConn` precedent (`cli/store/agent_store.go:628`):
`PRAGMA table_info(spec_numbers)` check, then `ALTER TABLE` only when the
column is missing. Add a migration step inside `EnsureSpecNumbersSchema`,
keeping schema creation and column migration in the same function so
every open self-heals.
