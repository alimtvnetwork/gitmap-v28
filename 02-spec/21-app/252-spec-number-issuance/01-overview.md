# Spec 252 — Spec Number Issuance

## User Request (Verbatim)

The owner hand-picked spec numbers by scanning the filesystem. Two parallel
streams both created spec "250" — the directories `250-package-consolidation`
and `250-gitmap-agent-collision-and-dispatch-registry-plan` both exist under
`02-spec/21-app/`.

Do not renumber anything retroactively. Build a command that issues spec
numbers from a central authority. The command should give us the number; the
database tells us where it is coming from for whichever repository we are
working on. Then: final build and release.

## Problem Statement

- Spec numbers are currently hand-picked by agents scanning `02-spec/21-app/`
  for the highest existing number and adding one.
- Parallel agent streams race this scan: both streams in program 250 read the
  same filesystem state and both issued themselves "250", producing two
  directories for one number (`250-package-consolidation` and
  `250-gitmap-agent-collision-and-dispatch-registry-plan`).
- Filesystem listing is not an authority — it is a read with no mutual
  exclusion. Two writers can observe the same "max" and both write the same
  next number.
- The fix is a single issuing authority: `gitmap spec next` returns the next
  spec number from a central store (SQLite), so parallel agents never collide
  on a number again.

## Non-goals

- No retroactive renumbering of anything on disk — the 250 collision stays as
  it is (owner explicit).
- No spec-management suite. `spec next` is the only required surface. No
  `spec list`, no `spec rename`, no `spec status`, no directory scaffolding.
  If those ever become wanted, they are a later spec.
- No global numbering across machines or repos. Numbering is per repository
  (per-repo DB), which matches how spec folders live inside one repo.
- No migration or repair tool for pre-existing collisions.

## Decisions (2026-10-09, recorded)

- D1: No retroactive renumbering. The 250 collision is historical fact; the
  spec defines issuance going forward only.
- D2: Output contract is a plain number on stdout (e.g. `252`), script-
  friendly, plus `--json` for structured consumers. Errors go to stderr with
  a nonzero exit — never a number on stdout when issuance failed.
- D3: Issuance = `max(disk, db) + 1`, claimed atomically via
  `INSERT OR IGNORE` through `store.ExecWrapper`; `RowsAffected == 0` means
  the number was taken, so the issuer retries with the next candidate.
  Idempotent `CREATE TABLE IF NOT EXISTS` for `spec_numbers` at DB open,
  following the `ensureCollisionSchema` precedent in
  `cli/cmdagent/agent_claim.go:17-52`. Canonical DB package is `cli/store`
  (`InitMasterAgentDB`, `ExecWrapper`).
- D4: Per-repo DB means no `--repo` flag. The Tier-1 DB at
  `<repo>/.ai-memory/temp-agents/ai_agents.db` (WAL) is resolved cwd-based
  via `store.ResolveMasterAgentDbPath("")` — it is per-repo by construction,
  so the issuer never needs to be told which repository it is working on.
- D5: Help text lives in hand-written `cli/helpdoc/spec.md`, auto-embedded
  via the existing `//go:embed *.md` in `cli/helpdoc/print.go` — zero Go
  wiring needed for help.
- D6: Registration follows the `space` precedent: command constants in
  `cli/constants/constants_cli.go` (cf. `CmdSpace`/`CmdSpaceBackupBranch` at
  lines 304-312), a new `cli/cmdspec/` package modeled on `cli/cmdspace/`
  with a subcommand router off `os.Args[2:]` (cf. `DispatchSpace` in
  `cli/cmdspace/space.go:87`), registered via direct call in
  `cli/cmd/root.go:601` (cf. `found, err = cmdspace.DispatchSpace(command)`),
  with the dispatch-table entry in `cli/cmd/rootcore.go` as the alternative.
