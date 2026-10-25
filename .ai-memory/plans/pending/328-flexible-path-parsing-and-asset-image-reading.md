# Plan: Flexible Path Parsing and Asset Image Reading

- **Task Sequence:** 328
- **Task Slug:** `flexible-path-parsing-and-asset-image-reading`
- **Spec Path:** `02-spec/21-app/328-flexible-path-parsing-and-asset-image-reading/`
- **Status:** `COMPLETED`
- **Created At:** 2026-10-25

---

## Executive Summary

This plan orchestrates the design and implementation of flexible path parsing and safe image reading in GitMap:

1. **Flexible Path Parser Engine (`cli/pathparse`):**
   - Normalizes conversational, markdown, and list-formatted file references into clean, canonical relative or absolute paths.
   - Fully covers the four required path formats:
     - **Format 1:** Bracketed Windows absolute paths (`@[C:\Users\Administrator\...\media_1792894009084_45b4893f.png]`).
     - **Format 2:** Markdown image syntax with URI schemes (`![Attached Image 1](file:///<drive>:/path/to/media_1792894009084_45b4893f.png)`).
     - **Format 3:** Numbered list, repo-relative paths with trailing comments (`1. @[assets/media_1792894009084_45b4893f.png] # consider the file is relative from the repo`).
     - **Format 4:** Numbered list, parent-relative paths with trailing comments (`1. @[../repo-cache/assets/media_1792894009084_45b4893f.png] # consider the file is relative from the repo outside of the repo`).
2. **Safe Image Asset Cat Subsystem (`cli/cmdmacro`):**
   - Upgrades `gitmap cat` to prevent raw binary streaming to stdout in interactive terminal sessions.
   - Detects image files via extension matching (`.png`, `.jpg`, `.jpeg`, `.gif`, `.webp`, `.bmp`, `.ico`, `.tiff`) and magic byte header sniffing.
   - Renders a styled, high-contrast **Safe Image Metadata Card** displaying dimensions, MIME type, byte size, color space, SHA-256 hash, and a terminal safety notice.
   - Provides an explicit `--raw` flag to stream raw binary bytes to stdout when required by pipelines or redirection.
3. **Skills & Guidance Synchronization:**
   - Updates `.agents/skills/` to instruct autonomous swarms and AI agents on using `gitmap cat` with rich path inputs and inspecting image attachments safely.

---

## Subtasks

- [x] `01-subtask-spec-and-planning.md` — Author comprehensive Task 328 specifications, audit ledger, and subtask implementation plans.
- [x] `02-subtask-flexible-path-parser.md` — Implement the core `cli/pathparse` package handling Formats 1, 2, 3, and 4 with positive booleans and $\le 15$ line functions.
- [x] `03-subtask-safe-image-cat-and-skills-sync.md` — Integrate `cli/pathparse` into `gitmap cat`, implement safe image detection and metadata card rendering, and synchronize skills.
