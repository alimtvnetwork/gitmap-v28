# Subtask 01 — DB layer (task 252, gitmap)

## File box
`cli/store/spec_numbers.go` — new file, DB primitive only. No disk
access, no spec-directory knowledge, no policy. Everything policy-shaped
(disk scan, candidate computation, retry loop) belongs in subtask 02's
`cli/cmdspec/`.

## Build (spec 03 §3.1–3.2)
1. Package-level DDL const, table lowercase `spec_numbers`:
   ```sql
   CREATE TABLE IF NOT EXISTS spec_numbers (
     number     INTEGER NOT NULL,
     repo_root  TEXT    NOT NULL DEFAULT '',
     task_slug  TEXT    NOT NULL DEFAULT '',
     issued_at  TEXT    NOT NULL,
     UNIQUE(number)
   );
   ```
2. `EnsureSpecNumbersSchema(conn *sql.DB) *appfault.AppError` — loop over
   `[]string{sqlCreateSpecNumbersTable}` through `store.ExecWrapper`
   (copy the `ensureCollisionSchema` shape in
   `cli/cmdagent/agent_claim.go:17-52`).
3. `ClaimSpecNumber(conn *sql.DB, number int, repoRoot, taskSlug string) (bool, error)` —
   `INSERT OR IGNORE INTO spec_numbers (number, repo_root, task_slug, issued_at)
   VALUES (?, ?, ?, <RFC3339 UTC now>)`; return `RowsAffected() == 1`.
   `RowsAffected == 0` is the normal "already taken" signal — return
   `false, nil`, never an error.
4. `MaxIssuedSpecNumber(conn *sql.DB) (int, error)` —
   `SELECT COALESCE(MAX(number), 0) FROM spec_numbers`.
5. Keep every function small, zero nesting where possible, positive
   boolean names (`claimed`, not `notFailed`).

## Rules
GitMap tools only (`gitmap aum search`); no rg/grep. `go build ./...`
only, never `go test`. Relative paths only. Never `git add`/`git commit`
(the lead commits). Do not touch any file outside the file box.

## Acceptance criteria
- `go build ./...` exits 0.
- Opening the master agent DB (`store.InitMasterAgentDB`) creates
  `spec_numbers` (verify with `sqlite3 .ai-memory/temp-agents/ai_agents.db
  ".schema spec_numbers"` — sqlite3 CLI is a read-only inspection tool,
  not a test suite).
- `ClaimSpecNumber(conn, 9999, ...)`: first call returns `true`, second
  call for the same number returns `false` (clean up the probe row
  afterwards with a single `DELETE`).

## Deliverable
Diff + build exit code. Then stop.
