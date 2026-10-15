# Subtask 06 — 20 new command proposals (program 273, WS4)

## File box

- ONE file: `.ai-memory/plans/gitmap-273-twenty-command-proposals.md` (≤300
  lines). Proposals ONLY — ZERO code, zero dispatch edits, zero version bump.

## Context

Spec: `02-spec/21-app/gitmap-help-text-and-ui-enhancement/06-twenty-new-command-proposals.md`.
Each proposal fills the mandatory template: PROBLEM, INVOCATION, EXPECTED
OUTPUT, ALIASES, HOST/TARGET, EFFORT (S/M/L), DEPENDS ON. Exactly 20. Seed
directions: `pr` family beyond `mr`, repo export/import round-trip,
`help search`/`help topic`, `theme list`/`theme preview`, fleet ergonomics,
pipeline triage shortcuts — plus subtask-04 audit gaps. No proposal may
duplicate an existing verb/alias. Rules: gitmap search tools only; relative
paths; ≤300 lines.

## Checklist

- [ ] Exactly 20 proposals, all template fields filled.
- [ ] Each has concrete invocation + normalized expected output.
- [ ] Each references an audit gap ID where applicable.
- [ ] Zero implementation artifacts (no .go changes in this wave).
- [ ] Commit: `gitmap cpf "proposals - twenty new command proposals"` + push
  (verify `git ls-remote`).
