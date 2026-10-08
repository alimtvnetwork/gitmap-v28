# Subtask 01 — `gitmap space backup-branch`

## Objective
Implement `gitmap space backup-branch "<task string>" [--no-push] [--force]` per
`02-spec/21-app/243-cli-help-displayer-and-backup-branch/02-space-backup-branch.md`.

## Files (disjoint box — no other worker touches these)
- `cli/cmdspace/` (new package; locate the existing `space` command implementation first
  via `gitmap aum search` — it currently only hosts `space common`)
- `cli/constants/` — add `CmdSpaceBackupBranch` constant (+ test registration)

## Done criteria
- [ ] Slug generation matches the spec (lowercase, hyphenated, 60-char cap, empty → error).
- [ ] Branch `backup/<slug>` created from HEAD; dirty tree refused with the spec'd message.
- [ ] Pushes to origin by default; `--no-push` keeps it local; `--force` recreates existing.
- [ ] `--help` renders via the existing help builder (migration to the new displayer happens in subtask 02/03).
- [ ] `go build` + `go vet` clean; new files small.
- [ ] Atomic commit via `gitmap cpf "space - backup-branch command"` + immediate push.
- [ ] Evidence: command output showing branch creation + `git branch -r` listing, pasted in the ledger.
