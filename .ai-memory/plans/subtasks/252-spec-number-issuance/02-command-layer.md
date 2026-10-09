# Subtask 02 — command layer (task 252, gitmap)

## File box
- `cli/cmdspec/spec.go` — new package `cmdspec`: `DispatchSpec(command string) (bool, error)`,
  `runSpecNext`, disk scan, retry loop, output.
- `cli/cmdspec/doc.go` — package doc comment (keep each file ≤ ~300 lines).
- `cli/constants/constants_cli.go` — add `CmdSpec = "spec"` (+ `CmdSpecNext = "next"`).
- `cli/cmd/root.go` — register `cmdspec.DispatchSpec(command)` in the dispatch
  chain near line 601 (`cmdspace`).
- `cli/helpdoc/spec.md` — hand-written help (auto-embedded by the helpdoc
  mechanism); covers `gitmap spec next [--json] [--slug <slug>]`.

## Build (spec 03 §3.3–3.5)
1. `DispatchSpec` routes `next` off `os.Args[2:]`; handles `--help`;
   returns `(false, nil)` for anything it does not own.
2. `runSpecNext`:
   - Resolve DB: `dbPath := store.ResolveMasterAgentDbPath("")`,
     `conn, appErr := store.InitMasterAgentDB(dbPath)`;
     call `store.EnsureSpecNumbersSchema(conn)` at open.
   - Disk scan: read `<repo>/02-spec/21-app/` (repo root = the directory
     gitmap is running from), take max of leading `^(\d+)-` directory
     prefixes; ignore non-matching entries and files; missing/empty dir → 0.
   - `candidate := max(diskScan, store.MaxIssuedSpecNumber(conn)) + 1`;
     loop `store.ClaimSpecNumber(conn, candidate, repoRoot, slug)`:
     on `true` → print, on `false` → `candidate++`, bound 1000 attempts.
3. Output: plain `<number>` on stdout (e.g. `252`); with `--json`,
   `{"number":253,"slug":"...","disk_scan_max":251,"db_max":252}`.
   `--slug <slug>` records the claiming program; default `""`.
4. Verify the collision-net: `cli/cmddispatch/registry.go` must route or
   list `spec` without conflicting with any existing top-level command.

## Rules
GitMap tools only (`gitmap aum search`); no rg/grep. `go build ./...`
only, never `go test`. Relative paths only, small functions, zero nesting
where possible, positive boolean names. Never `git add`/`git commit`.

## Acceptance criteria
- `go build ./...` exits 0.
- `gitmap spec next` prints a plain number on stdout (first run on this
  repo should print `252`: current disk max is 251, DB max is 0).
- Second immediate run prints `253` (each caller gets a fresh number).
- `gitmap spec next --json` emits the documented JSON shape.
- `gitmap spec --help` renders the `cli/helpdoc/spec.md` content.

## Deliverable
Diff + build exit code. Then stop.
