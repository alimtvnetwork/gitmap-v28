# Architecture — Code Block System

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Pipeline

The markdown-to-HTML pipeline uses a **placeholder extraction** pattern. Complex blocks are replaced with unique placeholders before inline formatting is applied, then restored at the end. This prevents regex collisions.

### Extraction Order (Critical)

1. **Code Blocks** — fenced triple-backtick blocks are replaced with CODEBLOCK placeholders
2. **Checklists** — checkbox list runs are replaced with CHECKLIST placeholders
3. **Inline Codes** — single-backtick spans are replaced with INLINECODE placeholders

> This order matters. Checklists can contain backticks, so code blocks must be extracted first. Inline codes are last because they're the simplest.

### Processing Steps

After extraction, the placeholder-free text goes through:

1. **convertTables** — pipe-delimited tables become HTML table elements
2. **convertInlineFormatting** — headings, bold, italic, links, blockquotes, HRs
3. **convertLists** — unordered and ordered lists
4. **wrapParagraphs** — remaining lines become paragraph tags

### Restoration

The **restorePlaceholders** function iterates all stores and replaces each placeholder with its stored HTML.

### Data Flow Diagram

```
markdownParser.ts
  ├── codeBlockExtractor.ts
  │     └── codeBlockBuilder.ts
  │           └── highlighter.ts
  │                 └── constants.ts
  ├── checklistBuilder.ts
  └── (inline extraction — internal)
```

---

## React Integration

### MarkdownRenderer Component

```
MarkdownRenderer({ content, allCollapsed, filePath })
  │
  ├── useMemo → renderMarkdown(content) → HTML string
  ├── useCodeBlockEvents(containerRef, html, setFullscreenBlock)
  ├── useEscapeFullscreen(!!fullscreenBlock, callback)
  ├── useSyncFullscreenClass(containerRef, fullscreenBlock, html)
  └── useCollapsibleSections(containerRef, html, allCollapsed, filePath)
```

The HTML is rendered via **dangerouslySetInnerHTML**. All interactivity is handled through event delegation on the container div.

### Event Delegation Pattern

Instead of attaching listeners to each button or line, a single set of listeners is attached to the container element. Each handler checks the event target against a selector to determine if it should act. This is efficient because:

- Code blocks are generated as raw HTML strings (not React components)
- The number of code blocks and lines can be very large
- Listeners are cleaned up and reattached when the HTML changes

---

## File Responsibilities

### codeBlockExtractor.ts

- Scans markdown line-by-line
- Detects opening fences (three or more backticks followed by optional language)
- Detects closing fences (matching backtick count)
- Handles unclosed blocks gracefully (appends raw text back)
- Delegates HTML building to codeBlockBuilder.ts

### codeBlockBuilder.ts

- Receives raw code string + language + block ID
- Calls highlighter.ts for syntax highlighting
- Builds complete HTML structure: wrapper → header → body → selection bar
- Generates all buttons with embedded data attributes

### highlighter.ts

- Wraps highlight.js with tree structure detection
- Registers 12 languages explicitly
- Falls back to auto-detection, then to tree rendering
- Tree renderer adds 📁/📄 icons, colorizes guides, ellipsis, comments

### useCodeBlockEvents.ts

- Single React hook that creates all handlers and attaches them
- Maintains refs for: active wrapper, anchor index, cursor index, pinned lines, drag state
- Splits handler creation into buildActionHandlers and buildLineHandlers

---

*Architecture — updated: 2026-04-08*
