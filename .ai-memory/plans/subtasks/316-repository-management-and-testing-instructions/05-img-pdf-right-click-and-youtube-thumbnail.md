# Subtask Plan 05: img-pdf Right-Click and YouTube Thumbnail Investigation

- **Subtask Slug:** `05-img-pdf-right-click-and-youtube-thumbnail`
- **Parent Task:** `316-repository-management-and-testing-instructions`
- **Target Repository:** img-pdf-v2 (local folder `img-pdf`, TypeScript/Vite app)

---

## 1. Overview & Objectives

1. **Right-click button does not work:** locate the context-menu / right-click
   handling in `src/` — component registration, event wiring, conditional
   rendering, platform gating — and produce a root cause with file-level evidence.
2. **YouTube thumbnail download not done:** determine whether the feature is
   absent, stubbed, or broken. Cite the modules involved or prove their absence,
   and state precisely what finishing the feature requires (API surface, UI entry
   point, wiring).

## 2. Method

- Search via GitMap (`gitmap aum search`, `gitmap find`, `gitmap cat`) from inside
  the `img-pdf` folder; read source files directly where search is inconclusive.
  No rg/grep.
- Read-only investigation; the deliverable is a root-cause report. Fixes are a
  follow-up decision once the cause is proven.

## 3. Acceptance Criteria

- Right-click: root cause stated in one sentence, backed by relative paths and
  the failing code path.
- YouTube thumbnail: status (absent / stub / broken) proven with evidence, plus a
  concrete finishing plan.
- Report delivered into the dedicated `img-pdf-v2 repo` chat's parent program —
  i.e. to the lead in the main chat.
