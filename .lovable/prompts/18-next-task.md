# Next 17 Task — Audit reclone URL pickers (plan 03 step 1)

Executed step 1 of `.lovable/plans/pending/03-reclone-transport-and-vscode-open.md`:

- Wrote `.lovable/audits/2026-06-07-reclone-pickers.md` with per-command verdict + file:line evidence.
- Flipped `.lovable/plans/subtasks/03-…/01-audit-reclone-pickers.md` `Status:` → `completed`.
- Bumped to v6.25.0 across `gitmap/constants/constants.go`, `src/constants/index.ts`, `README.md`, `CHANGELOG.md`.

Verdict (one-liner): `cfr` / `cfrp` are the only reclone-class commands that do NOT consult the destination folder's existing `remote.origin.url` before cloning, so an HTTPS positional URL silently downgrades an SSH-origin folder. Fix scheduled for plan 03 step 3.
