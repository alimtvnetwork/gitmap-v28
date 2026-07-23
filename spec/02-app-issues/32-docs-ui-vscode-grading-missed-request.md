# 32 — Docs UI VS Code grading request was missed repeatedly

## Summary

The user asked multiple times for the docs UI to use a VS Code-style color grading, but the implementation kept focusing on backend/update RCA work instead of the requested frontend visual change.

## Root cause

1. The request contained strong frustration plus references to earlier backend failures, and the workflow over-weighted the debugging context instead of the explicit new ask: **change the docs UI color grading**.
2. The existing frontend already had a partial VS Code Dark+ token set in `src/index.css`, which created a false sense that the design request was already done. In reality the UI still mixed default light surfaces, green brand accents, aurora hero styling, and generic card treatments.
3. Theme behavior defaulted to stored/system preference instead of making the requested VS-style dark grading the primary docs experience, so even the existing dark palette was not reliably what the user saw first.
4. Release hygiene drift compounded trust loss: prior responses discussed fixes without consistently syncing version + both changelog sources for the user-visible UI request.

## Solution

1. Make the docs shell visually consistent with a VS Code-like information architecture: darker explorer/sidebar, restrained borders, flat header, panel-style content framing, and blue selection accents.
2. Promote the VS-style dark grade to the default theme path while preserving the light toggle.
3. Remove decorative hero treatment that clashes with the requested editor-style look.
4. Record the change in both changelog sources and the pending-issues tracker so the request is auditable.

## Files involved

- `src/index.css`
- `src/lib/theme.ts`
- `src/components/docs/DocsLayout.tsx`
- `src/components/docs/DocsSidebar.tsx`
- `src/components/docs/FeatureCard.tsx`
- `src/components/docs/InstallBlock.tsx`
- `src/components/docs/CodeBlock.tsx`
- `src/pages/Index.tsx`

## Prevention

1. When a user asks for a UI style change, handle the frontend request directly unless they explicitly ask to debug backend logic too.
2. Partial design similarity in tokens does **not** count as completion; verify the visible shell, hero, cards, and code surfaces all match the requested direction.
3. Any user-visible fix must ship with synced version metadata in `CHANGELOG.md` and `src/data/changelog.ts`.