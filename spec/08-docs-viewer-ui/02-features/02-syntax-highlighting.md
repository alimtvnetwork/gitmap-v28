# Feature: Syntax Highlighting

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Specification

### Library

- **highlight.js** — lightweight, 40+ languages, auto-detection
- Theme: `github-dark` for dark mode, `github` for light mode
- Only register needed languages: `typescript`, `go`, `php`, `css`, `json`, `bash`, `sql`, `rust`

### Integration Point

In `MarkdownRenderer.tsx`, replace the plain `<pre><code>` output with highlight.js-processed HTML.

### Copy Button on Code Blocks

Each code block gets a floating "Copy" button (top-right corner) that copies the raw code to clipboard.

---

*Syntax highlighting — updated: 2026-04-03*
