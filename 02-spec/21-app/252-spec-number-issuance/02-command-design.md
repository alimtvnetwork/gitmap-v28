# Spec 252 — Command Design: `gitmap spec next`

## CLI Surface

```
gitmap spec next [--json]
```

- `gitmap spec` (bare) and `gitmap spec --help` print usage (from
  `cli/helpdoc/spec.md`) and exit 0.
- `gitmap spec next` issues the next spec number for the current repository
  and prints it.
- Flag list:
  - `--json` — print the issuance record as JSON instead of a bare number.
  - `-h`, `--help` — print usage, exit 0. No other flags exist.
- Exit codes:
  - `0` — number issued successfully (both plain and `--json` modes), or
    usage printed.
  - `2` — usage error (unknown flag, unknown subcommand under `spec`).
  - `1` — issuance failure (DB open failed, insert failed, directory scan
    failed). The error message goes to stderr; stdout must contain no number.

## Output Contracts

- Plain mode: stdout gets ONLY the number followed by a newline, e.g.
  `252\n`. Nothing else — no labels, no decoration, no trailing commentary.
  This keeps it script-friendly:
  `N=$(gitmap spec next) && mkdir "02-spec/21-app/${N}-my-spec"`.
- JSON mode (`--json`): stdout is exactly one JSON object, pretty-printed
  with two-space indent via `json.MarshalIndent(x, "", "  ")` (precedent:
  `cli/cmdstats/stats.go:52-58,150-160`), parsed with
  `flag.NewFlagSet`:
  ```json
  {
    "repo": "<absolute repo path>",
    "number": 252,
    "issued_at": "<RFC3339>"
  }
  ```
  - `repo` is the absolute path of the repository the number was issued for
    (the DB's repo).
  - `issued_at` is the issuance timestamp in RFC3339 (UTC).
- Errors: returned as `error` and printed to stderr; exit code nonzero.
  A failed issuance must never print a number to stdout.

## Issuance Algorithm

1. Resolve the Tier-1 DB path cwd-based via
   `store.ResolveMasterAgentDbPath("")`.
2. Open via `store.InitMasterAgentDB` (runs `CREATE TABLE IF NOT EXISTS
   spec_numbers` idempotently at open).
3. Seed/candidate: `candidate = max(highest number seen on disk under
   02-spec/21-app/, highest number recorded in spec_numbers) + 1`.
   (Max number on disk today is 251, so the first issuance is 252.)
4. Attempt `INSERT OR IGNORE` of `candidate` via `store.ExecWrapper`.
   - `RowsAffected == 1` → issued. Return `candidate`.
   - `RowsAffected == 0` → taken by a concurrent issuer. Increment and
     retry (bounded; fail with an error if the retry budget is exhausted).
5. The insert is the atomic claim. Two parallel callers cannot both get
   `RowsAffected == 1` for the same number — the loser retries to the next
   free number.

## Worked Examples

Basic issuance (repo at `/home/user/work/gitmap-v28`):

```
$ gitmap spec next
252
```

JSON mode:

```
$ gitmap spec next --json
{
  "repo": "/home/user/work/gitmap-v28",
  "number": 252,
  "issued_at": "2026-10-09T09:00:00Z"
}
```

Parallel callers: two agents run `gitmap spec next` at the same time.

```
agent A $ gitmap spec next     → 253
agent B $ gitmap spec next     → 254
```

Both get distinct numbers; the second to commit wins the lower number,
the loser retries and takes the next one. No collision, no filesystem
scan, no manual coordination.

Failure:

```
$ gitmap spec next
Error: open /repo/.ai-memory/temp-agents/ai_agents.db: permission denied
$ echo $?
1
```

The error is on stderr; stdout is empty; the exit code is nonzero.

## Registration Checklist

1. **Constants** — add `CmdSpec` / `CmdSpecNext` to
   `cli/constants/constants_cli.go`, following the `CmdSpace` /
   `CmdSpaceBackupBranch` block at lines 304-312.
2. **Package** — create `cli/cmdspec/`, modeled on `cli/cmdspace/`:
   a `DispatchSpec` router switching on `os.Args[2:]` (cf.
   `cli/cmdspace/space.go:87`), with the `next` implementation in its own
   file. Precedent for script-friendly single-value stdout:
   `cli/cmdspace/backup_branch.go:56`.
3. **root.go dispatch chain** — register via direct call in
   `cli/cmd/root.go:601`: `found, err = cmdspec.DispatchSpec(command)`,
   mirroring the `space` line. Alternative: a dispatch-table entry in
   `cli/cmd/rootcore.go` (choose one; keep the other out).
4. **Help** — hand-write `cli/helpdoc/spec.md`; it is auto-embedded via
   `//go:embed *.md` in `cli/helpdoc/print.go` (no Go wiring).
5. **Collision-net check** — `cli/cmddispatch/registry.go` parses dispatch
   tables via go/ast and fails on name collisions: verify `spec` and `next`
   are unclaimed before landing. Also confirm no existing `spec` alias or
   subcommand elsewhere claims the name.
6. **DB plumbing** — `spec_numbers` table created idempotently at open
   (precedent `ensureCollisionSchema`, `cli/cmdagent/agent_claim.go:17-52`);
   all writes through `store.ExecWrapper`; per-repo path via
   `store.ResolveMasterAgentDbPath("")` — no `--repo` flag (D4).
