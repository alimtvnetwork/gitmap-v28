# Subtask 02: Context Menu Icons & Stale Verb Purging

- **Status**: completed

## Objectives
- In `cli/constants/constants_installctx.go`:
  - Define distinct, high-definition Windows shell icon constants (`shell32.dll` and `imageres.dll`).
- In `cli/cmdinstall/installctxentries.go`:
  - Assign distinct icons to direct entries, background entries, categories, tail entries, and file menu items.
  - Maintain the intentional 2× duplication of `90_terminal` / `91_docs` required by parity tests.
- In `cli/cmdinstall/installctx.go`:
  - In `runInstallCtxWindows`, purge existing `\shell` subkeys under `Directory\Background`, `Directory`, and `*` before writing fresh ones, eliminating all orphaned verbs.
- In `cli/cmdinstall/installctx_icon_test.go`:
  - Update tests to reflect new icon constants (`constants.CtxIconHelp`).

## Verification
- Run `gitmap context-menu install` and inspect registry keys with `reg query`.
- Verified: Stale verbs purged from `Directory\Background`, all entries assigned high-definition icons.
