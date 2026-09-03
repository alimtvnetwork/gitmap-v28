# 01 — Code Block System

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 09 Code Block System.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Read and understand all numbered specifications under `09-code-block-system/`.
- [ ] `/learn` Adhere strictly to `.lovable/folder-structure.md` and `.lovable/strictly-avoid.md`.
- [ ] `/goal` Verify zero explicit `true` boolean evaluations and no mixed-polarity conditionals.
- [ ] `/learn` Run all local verification linters via `python 03-ai-scripts/06-cicd-local-runner.py`.


. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Purpose

This specification fully documents the code block rendering system — a rich, IDE-like experience for displaying fenced code blocks inside a markdown viewer. It covers every file, every CSS rule, every interaction, and every color so that **any AI agent can reproduce this system from scratch** in a new project.

---

## Folder Contents

| File | Description |
|------|-------------|
| `01-index.md` | This file — table of contents and architecture summary |
| `02-architecture.md` | Pipeline overview, file map, data flow |
| `03-html-03-structure.md` | Exact HTML output for code blocks, checklists, inline code |
| `04-syntax-highlighting.md` | highlight.js setup, language registration, tree detection |
| `05-interactions.md` | Copy, download, fullscreen, font controls, line selection, drag, keyboard |
| `06-styling.md` | Complete CSS reference — every class, color, transition, animation |
| `07-constants-and-maps.md` | Language labels, accent colors, file extensions, font sizes |
| `08-clipboard.md` | Clipboard utility — fallback strategy, toast feedback |
| `09-checklist-blocks.md` | Checklist extraction, rendering, copy/export |
| `10-tree-structure-rendering.md` | Tree/folder visualization with box-drawing characters |
| `11-selection-bar.md` | Selection bar UI, line pinning, drag-select, keyboard navigation |
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

_Auto-generated section — see `02-spec/09-code-block-system/97-acceptance-criteria.md` for the full criteria index._

### AC-CB-001: Code-block system conformance: Index

**Given** Render fenced code blocks (incl. nested 4-backtick fences) and checklist blocks from the spec tree.
**When** Run the verification command shown below.
**Then** Nested fences preserve backtick counts; clipboard copy returns exact source; tree rendering matches the constants map.

**Verification command:**

```bash
npm run test
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-08-30_
