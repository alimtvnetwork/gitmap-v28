# Consolidated: Code Block System — Complete Reference

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Purpose

This is the **standalone consolidated reference** for the code block rendering system — a rich, IDE-like experience for displaying fenced code blocks inside a markdown viewer. An AI reading only this file must be able to reproduce the full system without consulting source spec files.

---

## Pipeline

```
Markdown string
  │
  ▼
extractCodeBlocks()       → replaces ``` fences with placeholders
  │                         builds HTML via codeBlockBuilder.ts
  ▼
extractChecklistBlocks()  → replaces [ ]/[x] runs with placeholders
  │
  ▼
extractInlineCodes()      → replaces `code` with placeholders
  │
  ▼
convertTables → convertInlineFormatting → convertLists → wrapParagraphs
  │
  ▼
restorePlaceholders()     → re-inserts all stored HTML
  │
  ▼
Final HTML string         → dangerouslySetInnerHTML
  │
  ▼
useCodeBlockEvents()      → attaches click/drag/keyboard listeners
```

### Extraction Order (Critical)

1. **Code Blocks** — fenced backtick blocks replaced with CODEBLOCK placeholders
2. **Checklists** — checkbox list runs replaced with CHECKLIST placeholders
3. **Inline Codes** — single-backtick spans replaced with INLINECODE placeholders

> This order matters. Checklists can contain backticks, so code blocks must be extracted first.

### React Integration

The `MarkdownRenderer` component renders HTML via `dangerouslySetInnerHTML`. All interactivity uses **event delegation** on the container div — no per-button or per-line listeners. Handlers are created fresh on each `useEffect` cycle and cleaned up on unmount.

```
MarkdownRenderer({ content, allCollapsed, filePath })
  ├── useMemo → renderMarkdown(content) → HTML string
  ├── useCodeBlockEvents(containerRef, html, setFullscreenBlock)
  ├── useEscapeFullscreen(!!fullscreenBlock, callback)
  ├── useSyncFullscreenClass(containerRef, fullscreenBlock, html)
  └── useCollapsibleSections(containerRef, html, allCollapsed, filePath)
```

---

## Source File Map

| File | Role |
|------|------|
| `markdownParser.ts` | Pipeline orchestrator |
| `codeBlockExtractor.ts` | Fence detection & placeholder insertion |
| `codeBlockBuilder.ts` | HTML generation for each code block |
| `highlighter.ts` | highlight.js wrapper + tree rendering |
| `constants.ts` | Language maps, font size limits |
| `types.ts` | Shared `ExtractionResult` type |
| `useCodeBlockEvents.ts` | React hook — event listener orchestrator |
| `codeBlockActionHandlers.ts` | Copy, download, fullscreen, checklist handlers |
| `codeBlockLineHandlers.ts` | Line click, pin, range, keyboard navigation |
| `codeBlockDragHandlers.ts` | Drag-select, hover highlight, font size |
| `codeBlockDomHelpers.ts` | Selectors, constants, DOM utilities |
| `checklistBuilder.ts` | Checklist extraction & HTML |
| `MarkdownRenderer.tsx` | React component — renders HTML, manages fullscreen |
| `clipboard.ts` | Clipboard utility with fallback |
| `index.css` | All visual styles |

---

## 4-Backtick Fence Rendering (Nested Code Blocks)

When markdown contains code blocks inside code blocks (e.g., a template showing example fences), the outer fence **must** use 4+ backticks and the inner fences use 3 backticks.

### How the Parser Handles It

The `codeBlockExtractor.ts` captures the backtick count via `` /^(`{3,})/ ``:

1. **Opening fence**: Records the exact backtick count (e.g., 4)
2. **Inner 3-backtick fences**: Ignored because they don't match the opening count
3. **Closing fence**: Must match the exact same backtick count as opening

### Matching Algorithm

```
State: NOT_IN_BLOCK

For each line:
  1. If NOT_IN_BLOCK:
     - Match /^(`{3,})(\w*)$/
     - If match: record backtickCount = match[1].length, lang = match[2]
     - Enter IN_BLOCK state

  2. If IN_BLOCK:
     - Match /^(`{N,})$/ where N = recorded backtickCount
     - If exact match: close block, build HTML, enter NOT_IN_BLOCK
     - If no match: append line to current block content

  3. At end of input:
     - If still IN_BLOCK: treat as unclosed — append raw text back
```

### Critical Data Integrity Rule

> In `specTree.json`, the `content` field **must preserve the exact backtick count** from the source `.md` file. If a source file uses 4-backtick fences, the JSON must contain 4 backticks — not 3. Truncating to 3 causes data corruption where inner fences are misinterpreted as closers.

### Example

````markdown
````markdown

# Template

```
NN-module-name/
├── 01-index.md
└── 99-consistency-report.md
```
````
````

The outer ```````` ```` ```````` fence (4 backticks) wraps content containing inner ```` ``` ```` fences (3 backticks). The parser correctly treats only the 4-backtick line as the closer.

### Nesting Depth Support

| Outer Fence | Inner Fence | Inner-Inner | Status |
|-------------|-------------|-------------|--------|
| 3 backticks | — | — | Standard code block |
| 4 backticks | 3 backticks | — | One level of nesting |
| 5 backticks | 4 backticks | 3 backticks | Two levels (rare) |

Each level adds one backtick. The parser's `/^(`{3,})/` regex supports arbitrary depth.

---

## HTML Structure

Each code block renders as:

```html
<div class="code-block-wrapper" style="--lang-accent: H S% L%; --badge-color: H S% L%"
     data-block-id="N" data-language="typescript">
  <div class="code-block-header">
    <div class="code-lang-badge">
      <span class="code-lang-dot"></span>
      TYPESCRIPT
    </div>
    <div class="code-header-right">
      <span class="code-line-count">N lines</span>
      <span class="code-selection-label" style="display:none"></span>
      <div class="code-font-controls">
        <button class="code-tool-btn font-decrease-btn">A-</button>
        <button class="code-tool-btn font-reset-btn">A</button>
        <button class="code-tool-btn font-increase-btn">A+</button>
      </div>
      <button class="code-tool-btn copy-code-btn" data-code="...">Copy</button>
      <button class="code-tool-btn download-code-btn" data-code="..." data-ext="ts">Download</button>
      <button class="code-tool-btn fullscreen-code-btn" data-block-id="N">Fullscreen</button>
    </div>
  </div>
  <div class="code-block-body">
    <pre class="code-line-numbers">
      <span class="code-line-number" data-line="1">1</span>
      <span class="code-line-number" data-line="2">2</span>
    </pre>
    <pre class="code-content"><code class="hljs language-typescript">
      <span class="code-line" data-line="1">...</span>
      <span class="code-line" data-line="2">...</span>
    </code></pre>
  </div>
  <div class="copy-selected-bar" style="display:none">
    <span class="copy-selected-label"></span>
    <button class="code-tool-btn copy-selected-btn">Copy selected</button>
    <button class="code-tool-btn clear-selected-btn">✕</button>
  </div>
</div>
```

---

## Interactions

All interactions use **event delegation** on the markdown container. Handlers check `.closest(selector)` on the event target.

### 1. Copy Code

**Trigger:** Click `.copy-code-btn`

1. Read `data-code` attribute (HTML-escaped raw code)
2. Decode: `&#10;` → `\n`, `&#39;` → `'`, `&amp;` → `&`, `&lt;` → `<`, `&gt;` → `>`, `&quot;` → `"`
3. Copy to clipboard via `copyTextToClipboard()` (Clipboard API with `execCommand('copy')` fallback)
4. Show feedback: hide copy icon, show check icon, label → "Copied!", add `.copied` class
5. After **2000ms**, revert all changes

### 2. Download Code

**Trigger:** Click `.download-code-btn`

1. Read `data-code` + `data-ext` from button
2. Decode escaped characters
3. Create `Blob` (`type: "text/plain"`), temporary `<a>` with `download="code.{ext}"`, trigger click

### 3. Fullscreen Toggle

**Trigger:** Click `.fullscreen-code-btn`

1. Read `data-block-id`, toggle via React state
2. `.code-fullscreen` class added/removed on wrapper
3. Overlay `<div class="code-fullscreen-overlay">` rendered behind

**Exit:** Click fullscreen button again, click overlay, or press **Escape**.

**Layout:** Block → `position: fixed; inset: 2rem; z-index: 999`. Overlay → `position: fixed; inset: 0; z-index: 998; background: hsl(0 0% 0% / 0.7); backdrop-filter: blur(4px)`.

### 4. Font Size Controls

**Trigger:** Click `.font-increase-btn`, `.font-decrease-btn`, or `.font-reset-btn`

| Action | Formula | Limits |
|--------|---------|--------|
| Increase | `current + 2` | Max **32px** |
| Decrease | `current - 2` | Min **12px** |
| Reset | `18` | Default |

Sets `--code-font-size` on the wrapper. Line numbers scale: `calc(var(--code-font-size) * 0.7)`.

### 5. Line Selection — Click

**Single click** on `.code-line` or `.code-line-number` (excluding tool buttons):
1. Clear all `.line-pinned` in wrapper
2. Add `.line-pinned` to clicked line + corresponding line number
3. Update selection bar: "Line {N}"

**Shift+click:** Pin range from previous anchor to clicked line. Selection bar: "Lines {from}–{to}".

### 6. Line Selection — Keyboard

Arrow Up/Down when a wrapper is active:
- **Without Shift:** Move single pin to adjacent line, scroll into view
- **With Shift:** Extend/contract range from anchor
- Guards: Ignored if `document.activeElement` is `INPUT`, `TEXTAREA`, or `contentEditable`

### 7. Line Selection — Drag

**Trigger:** `mousedown` on `.code-line-number`

1. `mousedown` — Record anchor, clear pins, pin anchor line
2. `mousemove` (document) — `document.elementFromPoint()` → resolve line → `pinRange()`
3. `mouseup` (document) — Clear drag state

### 8. Line Hover

`mouseover` → add `.line-highlight` to line + line number. `mouseout` → remove all `.line-highlight`.

### 9. Copy Selected / Clear Selection

- **Copy selected** (`.copy-selected-btn`): Extract `textContent` from all `.code-line.line-pinned`, join with `\n`, copy
- **Clear** (`.clear-selected-btn` ✕): Remove all `.line-pinned`, hide selection bar

### Selection Bar

Appears at bottom of code block when lines are pinned:
- Shows "Line {N}" or "Lines {from}–{to}"
- Animated in via `@keyframes slideUpBar` (opacity 0→1, translateY 4px→0)
- Header also shows selection label badge next to line count

### State Management (Ref-Based)

| Ref | Type | Purpose |
|-----|------|---------|
| `activeWrapperRef` | `Element \| null` | Wrapper that owns keyboard focus |
| `anchorIdxRef` | `number` | Anchor line for range selections |
| `cursorIdxRef` | `number` | Current cursor position |
| `lastPinnedRef` | `Map<Element, number>` | Last-pinned line per wrapper |
| `dragStateRef` | `DragState \| null` | Active drag session |

### Event Delegation Wiring

| Event | Target | Handler |
|-------|--------|---------|
| `click` | container | lineClick, copySelected, clearSelected |
| `mousedown` | container | dragStart |
| `mouseover` | container | lineHover |
| `mouseout` | container | lineLeave |
| `mousemove` | `document` | dragMove |
| `mouseup` | `document` | dragEnd |
| `keydown` | `document` | lineKey |

---

## Syntax Highlighting

### Library

- **highlight.js** v11+ (core only — tree-shakeable)
- Theme CSS: `highlight.js/styles/github-dark.css`
- Custom token colors override the theme

### Registered Languages

| Registration Name(s) | highlight.js Module |
|-----------------------|---------------------|
| `typescript`, `ts`, `tsx`, `javascript`, `js` | `typescript` |
| `go`, `golang` | `go` |
| `php` | `php` |
| `css` | `css` |
| `json` | `json` |
| `bash`, `sh`, `shell` | `bash` |
| `sql` | `sql` |
| `rust` | `rust` |
| `html`, `xml` | `xml` |
| `yaml`, `yml` | `yaml` |
| `markdown`, `md` | `markdown` |

### Language Resolution Flow

```
1. normalizeLang(lang) — trim, lowercase, check known groups
2. resolveDisplayLang(code, lang) — if no lang AND tree-like → "tree"
3. highlightCode(code, lang):
   - If no lang AND tree-like → highlightAsTree()
   - If lang registered → hljs.highlight(code, { language })
   - Otherwise → hljs.highlightAuto(code)
     - If auto returns plaintext AND tree-like → highlightAsTree()
   - All paths have try/catch → fallback to escapeHtml()
```

### Syntax Token Colors (CSS Custom Properties)

| Token Type | CSS Class(es) | Color Variable |
|------------|---------------|----------------|
| Keywords, types, built-ins | `.hljs-keyword`, `.hljs-type`, `.hljs-built_in` | `--primary` (purple) |
| Function/class names | `.hljs-title`, `.hljs-section` | `--foreground / 0.85` |
| Strings, attributes | `.hljs-string`, `.hljs-attr`, `.hljs-property` | `--accent` (pink) |
| Numbers, symbols | `.hljs-number`, `.hljs-symbol`, `.hljs-regexp` | `--warning` (amber) |
| Comments | `.hljs-comment`, `.hljs-quote` | `--muted-foreground` (italic) |

---

## Constants & Maps

### Language Labels (`LANG_LABELS`)

| Key(s) | Display Label |
|--------|---------------|
| `ts`, `tsx`, `typescript` | TypeScript |
| `js`, `javascript` | JavaScript |
| `go`, `golang` | Go |
| `php` | PHP |
| `css` | CSS |
| `json` | JSON |
| `bash` | Bash |
| `sh`, `shell` | Shell |
| `sql` | SQL |
| `rust` | Rust |
| `html` | HTML |
| `xml` | XML |
| `yaml`, `yml` | YAML |
| `md`, `markdown` | Markdown |
| `tree` | Structure |
| `text`, `""` | Plain Text |

### Language Accent Colors (`LANG_COLORS`)

HSL values for `--lang-accent` and `--badge-color`:

| Language | HSL |
|----------|-----|
| TypeScript / TS / TSX | `99 83% 62%` |
| JavaScript / JS | `53 93% 54%` |
| Go / Golang | `194 66% 55%` |
| PHP | `234 45% 60%` |
| CSS | `264 55% 58%` |
| JSON | `38 92% 50%` |
| Bash / SH / Shell | `120 40% 55%` |
| SQL | `200 70% 55%` |
| Rust | `25 85% 55%` |
| HTML / XML | `12 80% 55%` |
| YAML / YML | `0 75% 55%` |
| Markdown / MD / Tree | `252 85% 60%` |
| **Default** (unlisted) | `220 10% 50%` |

### Language Extensions (`LANG_EXTENSIONS`)

| Language | Extension |
|----------|-----------|
| TypeScript | `ts` |
| TSX | `tsx` |
| JavaScript / JS | `js` |
| Go | `go` |
| PHP | `php` |
| CSS | `css` |
| JSON | `json` |
| Bash / SH / Shell | `sh` |
| SQL | `sql` |
| Rust | `rs` |
| HTML | `html` |
| XML | `xml` |
| YAML / YML | `yaml` / `yml` |
| Markdown / MD | `md` |
| Tree / Text / "" | `txt` |

### Language Group Constants

| Constant | Values |
|----------|--------|
| `TYPESCRIPT_LANGS` | `["typescript", "ts", "tsx"]` |
| `JAVASCRIPT_LANGS` | `["javascript", "js"]` |
| `GO_LANGS` | `["go", "golang"]` |
| `PLAINTEXT_LANGS` | `["text", "plaintext", "plain", "tree"]` |
| `ALL_SUPPORTED_LANGS` | `["php", "css", "json", "bash", "sh", "shell", "sql", "rust", "html", "xml", "yaml", "yml", "markdown", "md"]` |

### Timing & Size Constants

| Constant | Value | Usage |
|----------|-------|-------|
| `DEFAULT_FONT_SIZE` | `18` (px) | Default code font |
| `MIN_FONT_SIZE` | `12` (px) | Minimum via A- button |
| `MAX_FONT_SIZE` | `32` (px) | Maximum via A+ button |
| `FONT_SIZE_STEP` | `2` (px) | Step per click |
| `COPY_FEEDBACK_DELAY` | `2000` (ms) | "Copied!" feedback duration |

### DOM Selectors

| Constant | Selector |
|----------|----------|
| `SELECTOR_COPY_BTN` | `.copy-code-btn` |
| `SELECTOR_DOWNLOAD_BTN` | `.download-code-btn` |
| `SELECTOR_FULLSCREEN_BTN` | `.fullscreen-code-btn` |
| `SELECTOR_CHECKLIST_COPY_BTN` | `.checklist-copy-btn` |
| `SELECTOR_CHECKLIST_EXPORT_BTN` | `.checklist-export-btn` |
| `SELECTOR_COPY_SELECTED_BTN` | `.copy-selected-btn` |
| `SELECTOR_CLEAR_SELECTED_BTN` | `.clear-selected-btn` |
| `SELECTOR_CODE_WRAPPER` | `.code-block-wrapper` |
| `SELECTOR_TOOL_EXCLUSIONS` | `.code-tool-btn, .code-font-controls, .copy-selected-bar` |
| `SELECTOR_LINE_NUMBER` | `.code-line-number` |
| `SELECTOR_CODE_LINE` | `.code-line` |

### CSS Class Constants

| Constant | Class |
|----------|-------|
| `LINE_PINNED_CLASS` | `line-pinned` |
| `LINE_HIGHLIGHT_CLASS` | `line-highlight` |
| `COPIED_FEEDBACK_COLOR` | `hsl(152 70% 50%)` |

---

## Styling

### Design Philosophy

Code blocks use a **permanently dark theme** regardless of light/dark mode. Background is always `hsl(220, 14%, 11%)`. Each language gets a unique accent color (HSL) that tints the badge, hover shadow, and glow effects.

### CSS Custom Properties on `.code-block-wrapper`

| Property | Default | Modified By |
|----------|---------|-------------|
| `--code-font-size` | `18px` | Font control buttons |
| `--code-line-height` | `1.6` | Fixed |
| `--lang-accent` | Per-language HSL | Set at build time |
| `--badge-color` | Same as `--lang-accent` | Set on badge |

### Hardcoded Color Constants (Not Themed)

| Color | HSL | Usage |
|-------|-----|-------|
| Block background | `220, 14%, 11%` | `.code-block-wrapper` |
| Header background | `220, 14%, 14%` | `.code-block-header` |
| Line numbers bg | `220, 14%, 9%` | `.code-line-numbers` |
| Header border | `220, 13%, 20%` | Header bottom |
| Wrapper border | `220, 13%, 22%` | Outer border |
| Button background | `220, 13%, 20%` | `.code-tool-btn` |
| Button border | `220, 13%, 25%` | `.code-tool-btn` border |
| Button hover bg | `220, 13%, 28%` | Hover state |
| Font controls bg | `220, 13%, 18%` | `.code-font-controls` |
| Line number color | `220, 10%, 35%` | Muted line numbers |
| Line count color | `220, 10%, 45%` | "N lines" text |
| Button text color | `220, 10%, 65%` | Default button text |
| Hover line bg | `220, 15%, 16%` | `.line-highlight` |
| Hover line-num bg | `220, 15%, 12%` | Line number highlight |
| Copied bg | `152, 60%, 18%` | Copied button bg |
| Copied text | `152, 70%, 60%` | Copied button text |
| Copied border | `152, 50%, 30%` | Copied button border |

### Key Component Styles

**`.code-block-wrapper`:** `border-radius: 0.75rem`, `border: 1px solid hsl(220, 13%, 22%)`, `background: hsl(220, 14%, 11%)`, `font-family: 'Ubuntu Mono', 'JetBrains Mono', ui-monospace, monospace`.

**Hover:** `box-shadow: 0 8px 32px hsl(var(--lang-accent) / 0.1), 0 0 0 1px hsl(var(--lang-accent) / 0.15)`, `transform: translateY(-2px)`.

**`.code-lang-dot`:** `width: 7px`, `height: 7px`, `border-radius: 50%`, `background: hsl(var(--badge-color))`, `box-shadow: 0 0 6px hsl(var(--badge-color) / 0.5)`.

**`.line-pinned`:** `background: hsl(var(--primary) / 0.12)`. Line number: `color: hsl(var(--primary))`, `border-right: 2px solid hsl(var(--primary) / 0.6)`.

**Fullscreen:** `position: fixed; inset: 2rem; z-index: 999; border-radius: 1rem`, `box-shadow: 0 25px 80px hsl(var(--lang-accent) / 0.25), 0 0 0 1px hsl(var(--lang-accent) / 0.3)`.

**Inline code (`.inline-code`):** `background: hsl(var(--code-bg))`, `color: hsl(var(--code-text))`, `padding: 0.2em 0.45em`, `border-radius: 5px`, `font-size: 0.85em`, `font-family: 'JetBrains Mono'`. Hover: `box-shadow: 0 0 0 2px hsl(var(--highlight-glow) / 0.15)`, `transform: translateY(-1px)`.

### Font Stacks

| Context | Fonts |
|---------|-------|
| Code blocks | `'Ubuntu Mono', 'JetBrains Mono', ui-monospace, monospace` |
| Tool buttons | `'Ubuntu Mono', ui-monospace, monospace` |
| Inline code | `'JetBrains Mono', 'Fira Code', ui-monospace, monospace` |

---

## Tree Structure Rendering

### Auto-Detection

A code block is classified as a tree if **any** pattern matches:

| Pattern | Regex | Example |
|---------|-------|---------|
| Box-drawing characters | `/[├└│─]/` | `├── src/` |
| Trailing-slash directory | `/^\s*[A-Za-z0-9{}._-]+\/$/m` | `components/` |
| File with extension | `/^\s*[A-Za-z0-9{}._-]+\.[A-Za-z0-9_-]+\s*$/m` | `index.ts` |

Priority: Tree detection runs **before** highlight.js. Force tree rendering with ` ```tree ` language tag.

### Line-by-Line Pipeline

1. **Comment extraction:** Split on `#` — content before, comment after
2. **HTML escaping:** `escapeHtml()` on content portion
3. **Regex chain** (order matters):
   - Guide chars (`├└│─┌┐┘┬┴┤┼`) → `<span class="tree-guide">`
   - Ellipsis (`...`) → `<span class="tree-ellipsis">`
   - Directories (`name/`) → `<span class="tree-dir">📁 name/</span>`
   - Files (`name.ext`) → `<span class="tree-file">📄 name.ext</span>`
4. **Comment append:** `<span class="tree-comment"># text</span>`

### Tree CSS Classes

| Class | Color | Style |
|-------|-------|-------|
| `.tree-guide` | `--muted-foreground / 0.5` | Subdued guides |
| `.tree-dir` | `--foreground`, `font-weight: 600` | Bold directories |
| `.tree-file` | `--foreground / 0.85` | Slightly muted files |
| `.tree-ellipsis` | `--accent` | Accent-colored omission |
| `.tree-comment` | `--muted-foreground`, `italic` | Muted comments |

### Edge Cases

- Mixed tree + non-tree: entire block renders as tree (binary detection)
- Multi-dot files (`app.config.ts`): correctly matched
- `{}` in names: supported for placeholder syntax
- `#` in filenames: treated as comment delimiter (known limitation)
- Blank lines: pass through unchanged
- No box-drawing chars: still detected via directory/file patterns, renders with icons but no guide styling

---

## Checklist Blocks

### Detection Pattern

```regex
/^(\s*)([-*+]|\d+\.)\s+\[([ xX])\]\s+(.+)$/
```

Consecutive matching lines are grouped into a single checklist block.

### HTML Structure

```html
<div class="checklist-block" data-checklist-id="{ID}">
  <div class="checklist-header">
    <span class="checklist-title">Checklist</span>
    <div class="checklist-actions">
      <button class="checklist-copy-btn" data-checklist="{ENCODED_MD}">Copy</button>
      <button class="checklist-export-btn" data-checklist="{ENCODED_MD}">Export</button>
    </div>
  </div>
  <ul class="checklist-items">
    <li class="spec-checkbox checked" style="margin-left: {LEVEL}rem">
      <span class="checkbox-box checked-box">✓</span>
      <span class="checklist-item-content">{CONTENT}</span>
    </li>
  </ul>
</div>
```

### Behaviors

- **Nesting:** `Math.floor(indent / 2)` levels, each adding `1rem` left margin
- **Copy:** Copies raw markdown, HTML-encoded in `data-checklist`. Feedback: "Copied!" for 2000ms
- **Export:** Downloads as `checklist.md` via `Blob` + temporary `<a>`
- **Inline formatting:** Supports `` `code` ``, `**bold**`, `*italic*`, `[links](url)` in item content
- **Styling:** `border: 1px solid hsl(var(--border))`, `background: hsl(var(--card))`, hover: primary-tinted border + shadow
- **Checkbox hover:** `transform: translateX(3px)`, box `scale(1.1)`

---

*Consolidated code block system — v3.2.0 — 2026-04-16*
