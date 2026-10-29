# Subtask 03: Navigation Physics, Motion & Hygiene (Clusters 07, 08, 09, 10)

**Parent:** `346-next-100-ui-improvements`  
**Scope:** Items 61–100 from master audit ledger.
**Status:** COMPLETED

## Execution Checklist

- [x] Incorporate safe-area padding and 72px min-height glassmorphic physics on navigation headers (`ui_assets_theme.go`, `ui_assets_markup_misc.go`, `Settings.tsx`).
- [x] Add `transform-gpu` to off-canvas sidebars, modal dialogs, and animated view transitions (`sidebar.tsx`, `dialog.tsx`, `ProjectDetailDialog.tsx`).
- [x] Deconstruct `src/pages/Settings.tsx` into modular cards <= 90 lines each (`NodeResourceCard.tsx`, `UiUxRepoCatalogCard.tsx`, `referenceReposData.ts`).
- [x] Deconstruct `src/pages/CommitIn.tsx` into modular tables <= 95 lines each (`CommitInFlagsTable.tsx`, `CommitInAutoInitTable.tsx`, `CommitInSponsorTemplates.tsx`, `CommitInExitCodesTable.tsx`).
- [x] Deconstruct `src/pages/Bookmarks.tsx` into modular cards <= 65 lines each (`BookmarkTerminalPreview.tsx`, `BookmarkSchemaTable.tsx`, `BookmarkFileLayoutTable.tsx`).
- [x] Deconstruct `src/components/projects/ProjectDetailDialog.tsx` into modular sections <= 65 lines each (`ProjectDetailSections.tsx`).
- [x] Harmonize container padding and child band spacing to uniform 8pt grid scales (`CloneFixRepo.tsx`, `CommitIn.tsx`, `MakePublic.tsx`, `MergeLeft.tsx`, `Replace.tsx`, `Bookmarks.tsx`).
