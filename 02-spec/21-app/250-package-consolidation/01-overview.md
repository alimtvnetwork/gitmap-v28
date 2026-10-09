# 250-package-consolidation — Overview

Spec program 250. Repo: `~/workspace/repos/gitmap-v28`.
Module: `github.com/alimtvnetwork/gitmap-v28/cli` (all paths below relative to repo root; packages live under `cli/`).

## User Request (Verbatim)

- Reduce the package count by grouping logically-together packages.
- Keep packages small (~300-line files) but consolidated.
- Mine coding-guidelines Go packages for reuse (before rebuilding helpers that may already exist there).
- Make stack-trace output a user setting (default ON).
- Kill the global stdout pipe hack.

## Decisions

- D1 — Merge, never delete logic. Packages fold into a surviving package; exported identifiers are re-homed (or re-exported via aliases during transition). No behavior is dropped without an explicit owner decision.
- D2 — Orphan triage before merges. The 20 orphan packages (no importers) are triaged first (see `01-orphan-triage.md`): merge-or-keep is the default verdict; nothing is deleted without the owner's explicit approval.
- D3 — Foundation-first merge order. Merge dependency-foundation packages (errors, JSON helpers, terminal output, secrets, path/temp/lock, docs) before command-layer leaves, so import rewrites happen bottom-up and each merge lands on a compiling tree.
- D4 — Hubs untouched this program. Giant hubs (`constants` 218 importers, `apperror` 191, `store` 113, `cliexit` 84, `model` 78, `cmdssh` 179 files, `cmdagy` 213 files) are not split here. Merging small packages INTO a hub is fine; splitting a hub is a separate task.
- D5 — Behavior identical, verified by lead-run build. Every merge group lands only with a clean `go build ./...`; no behavioral change beyond the specified stack-trace setting and the removal of the stdout pipe hack.
- D6 — Aliases preserved. Renamed/moved exported identifiers keep package-level type aliases for at least one release cycle where importers exist outside the tree (hubs excluded from this obligation).

## Scope

- In scope: concrete merge groups in `02-package-merge-map.md` (groups 1–10), the 20 orphan packages triage, stack-trace-as-setting, removal of the global stdout pipe hack, reuse scan of coding-guidelines Go packages.
- Import rewrites use `gitmap aum` batch operations; never `rg`/`grep`/`git grep` (repo rule).
- Verification per merge group: `go build ./...` clean + `go list` shows no new import cycles (Tarjan SCC confirmed ZERO cycles today — every merge is cycle-safe by construction).

## Non-goals

- No splitting of giant hubs (constants/store/cmdssh/cmdagy).
- No new features, no refactoring of internals beyond the mechanical move.
- No deletion of any package without the owner's explicit ask (hard rule).
- No changes to `cmd` dispatcher semantics.
