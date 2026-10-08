# Spec 243 — CLI Help Displayer, Backup Branch, Enforcement & Refactor

## Mindset
This work happens inside the staging-repo mindset (`.ai-memory/what-to-read.md` → "The Mindset — READ FIRST"):
velocity over ceremony, backup-first, small files, aliases are a feature, DRY, enforcement beats documentation.

## Goals
1. `gitmap space backup-branch "<task>"` — one-command backup branch (`backup/<slug>`) from current code (spec 02).
2. CLI help displayer system — `Displayer` interface + `HelpDisplay` struct + command/sub-command helpers + groups + theme system + smart suggestions, bound per the coding-guideline binding principle (spec 03).
3. DRY help — single source of truth; stop hand-maintaining help in two formats (spec 03).
4. Enforcement — close the gaps verified 2026-10-08: file-size cap, stale docs (spec 04).
5. Refactor — split the 741-file `cli/cmd/` glue package into smaller packages; `enums` become their own package per the Go coding guidelines (spec 04).

## Non-goals
- Do NOT touch the version cadence — fast-forward versioning is intentional.
- Do NOT prune aliases — they are a kept feature.
- Do NOT big-bang migrate all 162 help builders at once — incremental migration, new technique first.

## Decisions (2026-10-08, recorded)
- D1: `space backup-branch` pushes to origin by default (`--no-push` keeps it local). A backup only counts off-machine.
- D2: Dirty working tree → the command refuses with a clear message (commit or stash first). No silent stash games.
- D3: `--force` on backup-branch = overwrite an existing `backup/<slug>` branch.
- D4: Help displayer builds on the existing `cli/theme/` and `cli/glyphs/` packages — no new color engine.
- D5: Help-content structs are the source of truth; `helptext/*.md` is generated from them, never hand-maintained.
- D6: `cmd/` refactor = move implementation files into their `cmdX` packages; `cmd/` keeps only thin dispatch (root tables + argv preprocessing).
