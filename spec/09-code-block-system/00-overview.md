# 08 — Code Block System

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Purpose

This specification fully documents the code block rendering system — a rich, IDE-like experience for displaying fenced code blocks inside a markdown viewer. It covers every file, every CSS rule, every interaction, and every color so that **any AI agent can reproduce this system from scratch** in a new project.

---

## Folder Contents

| File | Description |
|------|-------------|
| `00-overview.md` | This file — table of contents and architecture summary |
| `01-architecture.md` | Pipeline overview, file map, data flow |
| `02-html-structure.md` | Exact HTML output for code blocks, checklists, inline code |
| `03-syntax-highlighting.md` | highlight.js setup, language registration, tree detection |
| `04-interactions.md` | Copy, download, fullscreen, font controls, line selection, drag, keyboard |
| `05-styling.md` | Complete CSS reference — every class, color, transition, animation |
| `06-constants-and-maps.md` | Language labels, accent colors, file extensions, font sizes |
| `07-clipboard.md` | Clipboard utility — fallback strategy, toast feedback |
| `08-checklist-blocks.md` | Checklist extraction, rendering, copy/export |
| `09-tree-structure-rendering.md` | Tree/folder visualization with box-drawing characters |
| `10-selection-bar.md` | Selection bar UI, line pinning, drag-select, keyboard navigation |
| `99-consistency-report.md` | Cross-reference validation |

---

## Architecture at a Glance

```
Markdown string
  │
  ▼
extractCodeBlocks() ─── replaces ``` fences with placeholders
  │                      builds HTML via codeBlockBuilder.ts
  ▼
extractChecklistBlocks() ─── replaces [ ]/[x] runs with placeholders
  │
  ▼
extractInlineCodes() ─── replaces `code` with placeholders
  │
  ▼
convertTables → convertInlineFormatting → convertLists → wrapParagraphs
  │
  ▼
restorePlaceholders() ─── re-inserts all stored HTML
  │
  ▼
Final HTML string → dangerouslySetInnerHTML
  │
  ▼
useCodeBlockEvents() ─── attaches all click/drag/keyboard listeners
```

---

## Source File Map

| File | Role |
|------|------|
| `src/components/markdown/markdownParser.ts` | Pipeline orchestrator |
| `src/components/markdown/codeBlockExtractor.ts` | Fence detection & placeholder insertion |
| `src/components/markdown/codeBlockBuilder.ts` | HTML generation for each code block |
| `src/components/markdown/highlighter.ts` | highlight.js wrapper + tree rendering |
| `src/components/markdown/constants.ts` | Language maps, font size limits |
| `src/components/markdown/types.ts` | Shared `ExtractionResult` type |
| `src/components/markdown/useCodeBlockEvents.ts` | React hook — event listener orchestrator |
| `src/components/markdown/codeBlockActionHandlers.ts` | Copy, download, fullscreen, checklist handlers |
| `src/components/markdown/codeBlockLineHandlers.ts` | Line click, pin, range, keyboard navigation |
| `src/components/markdown/codeBlockDragHandlers.ts` | Drag-select, hover highlight, font size |
| `src/components/markdown/codeBlockDomHelpers.ts` | Selectors, constants, DOM utilities |
| `src/components/markdown/checklistBuilder.ts` | Checklist extraction & HTML |
| `src/components/MarkdownRenderer.tsx` | React component — renders HTML, manages fullscreen |
| `src/lib/clipboard.ts` | Clipboard utility with fallback |
| `src/index.css` | All visual styles |

---

*Code Block System — overview — updated: 2026-04-08*

---

## Verification

_Auto-generated section — see `spec/09-code-block-system/97-acceptance-criteria.md` for the full criteria index._

### AC-CB-000: Code-block system conformance: Overview

**Given** Render fenced code blocks (incl. nested 4-backtick fences) and checklist blocks from the spec tree.  
**When** Run the verification command shown below.  
**Then** Nested fences preserve backtick counts; clipboard copy returns exact source; tree rendering matches the constants map.

**Verification command:**

```bash
npm run test
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
