---
name: Clone → VS Code Project Manager sync
description: Every clone variant (clone, clone-next, clone-from, clone-now, clone-pick, clone-multi, cfr/cfrp) must update alefragnani.project-manager projects.json after a successful clone. Honors --no-vscode-sync. v4.16.0+.
type: feature
---

# Feature: Clone-Time Sync to VS Code Project Manager (v4.16.0+)

## Rule
Every gitmap command that lands a new repo on disk MUST update the
`alefragnani.project-manager` `projects.json` file in the same
invocation. Equivalent to running `gitmap code <new-folder>` after
each clone — except automatic.

## Scope (every clone variant)
- `clone` (direct URL, manifest, multi-URL)
- `clone-next` / `cn` (single + batch)
- `clone-from` (manifest)
- `clone-now` / `reclone` (manifest)
- `clone-pick` (sparse subset)
- `clone-fix-repo` (cfr) and `clone-fix-repo-pub` (cfrp) — inherit via `executeDirectClone`

## Implementation
- Shared helper: `cmd/clonepmsync.go::syncClonedReposToVSCodePM(pairs, skip)`
- Calls `vscodepm.Sync` once per command (NOT per repo) — single atomic write to projects.json.
- Auto-tags via `vscodepm.DetectTags(absPath)` (mirrors `gitmap code`).
- Honors `constants.FlagNoVSCodeSync` (`--no-vscode-sync`) on every clone command.
- Soft-fails via `reportVSCodePMSoftError` — sync failures NEVER fail a successful clone.
- Only successful clones are synced (status=ok). Skipped/failed rows are excluded.

## Why
Original report: "when I do clone or CFRP, the new repo never appears
in Project Manager." Root cause: `syncCodeEntry` was private to
`gitmap code`; the clone family had no equivalent call. Fix factors
the sync into a shared helper so every entry-point that clones a repo
also updates Project Manager.

## Spec
`spec/01-vscode-project-manager-sync/02-clone-sync.md`
