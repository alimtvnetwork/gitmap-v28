# Subtask 03 — Docs site + UI surface for new commands

**Parent:** 02-ssh-aware-clone
**Status:** pending
**Created:** 2026-06-07

## Work
- Verify `src/data/commands.ts` entries for `make-all-public` / `make-all-private` are wired and surface in:
  - command list page
  - search
  - sidebar (`DocsSidebar.tsx`)
  - dedicated pages `src/pages/MakeAllPublic.tsx` + `MakeAllPrivate.tsx` (already created — confirm routes in `App.tsx`).
- Add a "Transport awareness" section to the scan docs page explaining identified transport and the dual-URL persistence.
- Add release-notes blurb on the changelog page if it's data-driven.
