# Subtask 08 — Lead final verification (program 273, WS6)

## File box

- Verification only — NO new code. A dated verification note may be added
  under `.ai-memory/plans/subtasks/gitmap-help-text-and-ui-enhancement/`
  if findings need recording. NO version bump. NO release.

## Context

Lead-owned final pass after WS2–WS5 land. Acceptance criteria:
spec `97-acceptance-criteria.md` (33 criteria). Hard rules: `go build ./...`
from `cli/` MUST be green; every commit atomic with
`gitmap cpf "<module> - <summary>"` and pushed (`git ls-remote` verified);
no binaries; no deletions; all files ≤300 lines; relative paths only;
NEVER `go test` locally. Verify theme pipe-safety (`gitmap <cmd> | cat`
byte-identical across themes), `help example` coverage = dispatcher verbs,
`mr` E2E, proposals file = exactly 20 with full templates, RCA file present.

## Checklist

- [ ] `go build ./...` green from `cli/`.
- [ ] All 33 acceptance criteria checked off with evidence.
- [ ] Every wave commit pushed (spot-check `git ls-remote`).
- [ ] Update the program ledger
  (`.ai-memory/plans/273-gitmap-help-text-and-ui-enhancement.md`) wave
  checkboxes to reflect completion.
- [ ] Confirm: no version bump, no release cut in program 273.
- [ ] Final report to the user with commit SHAs.
