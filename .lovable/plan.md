# Plan (post v6.80.0)

Version pinned to **v6.80.0** across `gitmap/constants/constants.go`,
`src/constants/index.ts`, `.gitmap/release/latest.json`, and
`.gitmap/release/v6.80.0.json`.

## Shipped in v6.80.0 — Update Awareness

Specification-side artifacts landed; Go implementation is the follow-up work
(see acceptance criteria in `spec/01-app/117-update-awareness.md`).

- `gitmap list --update` (alias `lu`, also `gitmap update list`)
- `gitmap update apply <repo>` (alias `ua`)
- `gitmap update all` (alias `uall`)
- `gitmap hd` (alias for `gitmap help --dashboard`)
- `gitmap stats` Upgrades block (text + JSON)
- Post-scan hint linking to `list --update`
- New `TaskType` seed `Upgrade`; `do-pending` replays via `update apply`
- Helptext: `list-update.md`, `update-apply.md`, `update-all.md`, `hd.md`;
  cross-links added to `scan.md`, `stats.md`, `pending.md`, `do-pending.md`
- JSON schemas: `list-update.schema.json`, `update-apply.schema.json`,
  `hd.schema.json`
- Docs site: entries added to `src/data/commands.ts` under `scanning` and
  `tools`

## Remaining backlog (from the Rejog v1 report)

Priority order preserved from `.lovable/memory/reports/20260723-rejog-reliability.md`.

### P0

1. Duplicate spec-file prefixes in `spec/01-app/` (multiple `89-*`, `90-*`,
   `95-*`, `96-*`). Rename or consolidate before large refactors.
2. Function-length audit against the ≤ 15 lines rule in
   `spec/05-coding-guidelines/02-go-code-style.md`.

### P1

1. SSH transport bug: `.lovable/issues/01-ssh-repo-cloned-as-https.md` and
   `.lovable/issues/02-reclone-loses-ssh-transport.md`.
2. `gitmap code` command spec in `.lovable/issues/03-no-gitmap-code-command.md`.

### P2

1. Test backfill for the new update-awareness commands once the Go
   implementation lands.

## Next task selection

Pick from P0 or P1 above. The next natural task is P1-1 (SSH transport) since
it is a live user-facing correctness bug, or P0-1 (spec-file dedup) if a
low-risk cleanup pass is preferred.
