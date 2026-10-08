# Subtask 03 — DRY help: generator + pilot migration

## Objective
Make the help structs the single source of truth; generate `helptext/*.md` from them;
migrate 2–3 pilot commands (`space`, `scan`, `pull-all`) to the new displayer.

## Files (disjoint box)
- Generator: `03-ai-scripts/` (new script, run via `gitmap py`) OR a gitmap subcommand — decide and document in the ledger before building.
- Pilot command help builders under their `cmdX/` packages + their `helptext/*.md` topics.

## Done criteria
- [ ] Generator produces `helptext/*.md` from `HelpDisplay` structs; hand-editing those topics is no longer needed (document this in the spec).
- [ ] Pilot commands render through the displayer; terminal output diffed against today's rendering as evidence (modulo intentional improvements).
- [ ] `--help <sub>` auto-renders sub-displays where `subHelpers` exist.
- [ ] `go build` + `go vet` clean.
- [ ] Atomic commit via `gitmap cpf "helpdisplay - dry generator and pilot migration"` + immediate push.
- [ ] Migration backlog: list of remaining builders ordered by duplication/staleness, recorded in the ledger for opportunistic follow-ups (no big bang).
