---
Slug: reclone-transport-and-vscode-open
Status: pending
Created: 2026-06-07
---

# Reclone honors stored transport + new `gitmap code` opens VS Code

**Slug:** reclone-transport-and-vscode-open
**Steps:** 5
**Status:** pending
**Created:** 2026-06-07

## Context
Two related gaps surfaced by the user:
1. `cfr` / `cfrp` / `clone-now` re-clone over HTTPS even when the working repo's `origin` is SSH, repeating the v6.19–v6.22 class bug in the reclone path and losing browser-auth-free behavior. Transport is never persisted on the `Repo` row, so even a manual fix is forgotten by the next reclone, and the event is invisible in `gitmap history`.
2. There is no `gitmap code` (aliases `vcode`, `vscode`) command to open the current/argument folder in VS Code and register it in the Project Manager extension's `projects.json` (the existing `vpm` only syncs an already-populated `projects.json`).

Captured artifacts:
- `.lovable/spec/commands/04-reclone-honors-stored-transport.md`
- `.lovable/spec/commands/05-gitmap-code-opens-vscode.md`
- `.lovable/issues/02-reclone-loses-ssh-transport.md`
- `.lovable/issues/03-no-gitmap-code-command.md`

Pulled in from existing pending plans (see "Appended from prior pending tasks").

## Steps
1. **Audit + classify** every reclone-class URL picker (`cfr`, `cfrp`, `clone-now`/`reclone`, direct-URL `clone`) and the shared `cloner.pickURL` to confirm which already honor identified transport vs which still hard-code HTTPS-first. Produce a one-line verdict per command with file:line refs. See `./subtasks/03-reclone-transport-and-vscode-open/01-audit-reclone-pickers.md`.
2. **Persist transport on the `Repo` row** via schema migration 007 adding `IdentifiedTransport TEXT NOT NULL DEFAULT ''`, extend `UpsertRepoByPath` + `SelectAll/BySlug/ByPath` + the `model.Repo` struct, and backfill from `origin` on first read after the migration. See `./subtasks/03-reclone-transport-and-vscode-open/02-db-identifiedtransport-column.md`.
3. **Wire reclone to read → store → reuse → log**: in every reclone entrypoint, resolve repo path (cwd / arg-folder / DB row by URL), parse `.git/config` `remote.origin.url`, classify transport, persist via step 2, pass the classified transport into the URL picker so the actual `git clone` uses it, and write a history event (`gitmap history` row) with `from`, `to`, `transport`, `command`, `exit`. See `./subtasks/03-reclone-transport-and-vscode-open/03-reclone-folder-and-url-flow.md`.
4. **Add `gitmap code` (aliases `vcode`, `vscode`)**: new `cmd/code.go` + dispatch entry + `CmdCode*` constants in `constants_cli.go` + marker comment for completion. Behavior: no-arg → cwd, with-arg → folder; resolve VS Code binary (PATH, then platform fallbacks); on success append/update the target into Project Manager's `projects.json` reusing `vscode-pm-sync` merge helpers; soft-fail on missing VS Code. Add focused unit test for the projects.json merge path. See `./subtasks/03-reclone-transport-and-vscode-open/04-gitmap-code-command.md`.
5. **Release closeout**: bump `gitmap/constants/constants.go` `Version` → `6.25.0`, mirror in `src/constants/index.ts`, repin `README.md` (every `v6.24.0` occurrence including the version-matrix table), and prepend a `## v6.25.0` `CHANGELOG.md` entry describing the reclone-transport fix + the new `code`/`vcode`/`vscode` command with file lists, root cause sentences, and verification notes — matching the format of v6.22–v6.24 entries.

## Verification
- Step 1: audit file checked in under the subtask folder with file:line evidence.
- Step 2: migration runs idempotently; `SelectAll` returns the new column; backfill smoke test.
- Step 3: in a fixture repo whose `origin` is SSH, `gitmap cfr` emits `git clone -b … git@…`, `Repo.IdentifiedTransport = 'ssh'` after one run, and `gitmap history` shows the event.
- Step 4: `gitmap code` opens cwd in VS Code on Windows + macOS smoke runs; `projects.json` gains the new entry without duplicating; `gitmap vcode` and `gitmap vscode` dispatch to the same handler; missing-VS-Code path prints a single stderr line and exits non-zero.
- Step 5: `gitmap --version` prints `6.25.0`; README and CHANGELOG diffs reviewed; preview confirms `src/constants/index.ts` version banner.

## Appended from prior pending tasks
- `01-bulk-visibility-mapub-mapri.md` — still pending; not blocking this plan but tracked in `.lovable/plans/pending/`.
- `02-ssh-aware-clone.md` — partially implemented across v6.19–v6.22 for scan/probe/cloner/direct-clone/report paths; the **reclone** path (this plan, step 3) and the **`IdentifiedTransport` persisted column** (this plan, step 2) are the remaining unresolved items from that plan.
