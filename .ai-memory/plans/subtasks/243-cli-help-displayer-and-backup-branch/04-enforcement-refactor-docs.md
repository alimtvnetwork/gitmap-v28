# Subtask 04 — Enforcement, cmd split, enums, stale docs

## Objective
Per `02-spec/21-app/243-cli-help-displayer-and-backup-branch/04-enforcement-and-refactor.md`.

## Files (disjoint box — coordinate the move map with the coordinator first)
- The 10 files >1000 lines (split targets; one file per wave).
- `03-ai-scripts/` (file-size check script).
- `cli/cmd/` (dispatch only afterwards) + `cmdX/` packages (receivers).
- `cli/enums/` (new package).
- `.ai-memory/overview.md`, `docs/architecture/state-ownership.md` (new).

## Done criteria
- [ ] The 10 files >1000 lines split (junk-drawer clusters → own files); then the 500–1000 band, largest first.
- [ ] File-size check script (`gitmap py`) fails past 500 lines; wired into the repo's pre-commit path.
- [ ] `cli/cmd/` keeps only thin dispatch (`root*.go`, argv preprocessing, alias context); impl files moved to their `cmdX` packages; exact move map documented in the ledger BEFORE moving; `go build` verified after each wave.
- [ ] `cli/enums/` created per the Go coding-guideline enum pattern (cite the pattern); most self-contained group migrated first.
- [ ] `.ai-memory/overview.md` version corrected (from `version.json`); state-ownership doc written (`.gitmap/` vs `gitmap.json` vs SQLite).
- [ ] New-command checklist updated to the spec-03 help technique.
- [ ] `go build` + `go vet` clean after every wave; atomic commits per wave (`gitmap cpf`/`cpb`) + immediate push.
- [ ] Do NOT touch version cadence. Do NOT prune aliases.
