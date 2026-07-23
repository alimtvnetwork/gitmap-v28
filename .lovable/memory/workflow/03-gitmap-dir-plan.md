# Plan: Unified .gitmap/ Directory Migration

## Goal

Consolidate `.release/` and `gitmap-output/` under a single `.gitmap/` folder.

See: `spec/01-app/56-unified-gitmap-dir.md`

## Steps

### Step 1 — Update constants (constants.go)
Change `DefaultReleaseDir` from `.release` to `.gitmap/release`.
Change `DefaultOutputDir` from `./gitmap-output` to `.gitmap/output`.
Change `DefaultOutputFolder` from `gitmap-output` to `.gitmap/output`.
Change `DefaultVerboseLogDir` from `gitmap-output` to `.gitmap/output`.
Add `GitMapDir = ".gitmap"` root constant.

### Step 2 — Update hardcoded display strings
Fix `constants_terminal.go` — repo count formats, clone step hints.
Fix `constants_messages.go` — `MsgNoOutputDir` error text.
Fix `constants_cli.go` — `HelpOutputPath` flag description.

### Step 3 — Update tests
Fix `config/config_test.go` expected `OutputDir`.

### Step 4 — Update help text files
Search all `helptext/*.md` for `.release/` and `gitmap-output/` references.

### Step 5 — Update spec documents
Search all `spec/` files for `.release/` and `gitmap-output/` references.

### Step 6 — Update docs site
Update `src/data/` and `src/pages/` references.

### Step 7 — Update memory files
Update `.lovable/memory/` references.

### Step 8 — Doctor check simplified
Legacy directory warnings removed from doctor; migration handles cleanup automatically.

### Step 9 — Bump version and changelog (v2.36.3)

### Step 10 — Automatic legacy directory migration
Add shared migration logic and call it at CLI startup plus after `release`
returns to the original branch.
Moves `gitmap-output/` → `.gitmap/output/`, `.release/` → `.gitmap/release/`,
`.deployed/` → `.gitmap/deployed/` automatically when detected.
When target already exists, merges files and removes legacy directory.
Skipped for `version` command to keep stdout clean.

## Status

Steps 1–10 complete. Release flow cleanup verified in design: migration reruns after checkout back to original branch.
