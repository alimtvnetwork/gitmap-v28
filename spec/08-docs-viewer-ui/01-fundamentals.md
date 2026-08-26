# Docs Viewer UI — Fundamentals

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** 90%  
**Ambiguity:** 10%

---

## Overview

The Docs Viewer (`/docs`) is a React-based specification browser that renders markdown files from the `spec/` folder tree. This document covers the core architecture and shared design decisions for all UI enhancements.

---

## Technology Stack

| Layer | Technology |
|-------|-----------|
| Framework | React 18 + TypeScript 5 |
| Styling | Tailwind CSS v3 + shadcn/ui |
| Fonts | Google Fonts — Ubuntu (headings), Poppins (body) |
| Syntax Highlighting | highlight.js (lightweight, 40+ language support) |
| Layout | Sidebar (shadcn) + flex content area |
| State | React hooks (useState, useCallback, useEffect) |

---

## Typography System

### Font Assignments

| Element | Font Family | Weight | Fallback |
|---------|-------------|--------|----------|
| H1–H4 headings | Ubuntu | 700 (bold), 600 (semi) | `sans-serif` |
| Body text, paragraphs | Poppins | 400 (regular) | `sans-serif` |
| Code blocks / inline code | `monospace` (system) | 400 | `Courier New` |
| Navigation / sidebar | Poppins | 400, 500 | `sans-serif` |

### Loading Strategy

- Fonts loaded via Google Fonts `<link>` in `index.html`
- Weights: Ubuntu 400,500,600,700 — Poppins 300,400,500,600
- `font-display: swap` for performance

---

## Navigation Model

### Keyboard Shortcuts

| Key | Action | Scope |
|-----|--------|-------|
| `←` Left Arrow | Previous file in current folder | File-level |
| `→` Right Arrow | Next file in current folder | File-level |
| `↑` Up Arrow | Previous folder (jump to its first file) | Folder-level |
| `↓` Down Arrow | Next folder (jump to its first file) | Folder-level |
| `F` | Toggle fullscreen mode | Global |
| `?` | Toggle shortcuts help overlay | Global |
| `Escape` | Exit fullscreen / close overlay | Global |

### Navigation Logic

- Files are ordered by their numeric prefix within each folder
- Folder navigation follows the `specTree.json` order
- Wrap-around: last file → first file (within folder); last folder → first folder (at root level)
- Shortcuts only active when no input/textarea is focused

---

## Feature Summary

| # | Feature | Priority |
|---|---------|----------|
| 1 | Ubuntu + Poppins typography | Must Have |
| 2 | Syntax highlighting (highlight.js) | Must Have |
| 3 | Fullscreen toggle button | Must Have |
| 4 | Keyboard arrow navigation | Must Have |
| 5 | Shortcuts help overlay (`?` key) | Must Have |
| 6 | Copy markdown button | Must Have |

---

## Cross-References

- [Spec Authoring Guide](../01-spec-authoring-guide/00-overview.md) — Folder conventions
- [App Project Template](../01-spec-authoring-guide/05-app-project-template.md) — Template this spec follows

---

*Fundamentals — updated: 2026-04-03*
