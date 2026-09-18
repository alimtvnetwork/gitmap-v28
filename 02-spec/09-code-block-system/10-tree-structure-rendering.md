# 10 — Tree-Structure Rendering

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

Code blocks containing folder/file hierarchies receive **special tree rendering** instead of standard syntax highlighting. The system auto-detects tree content, applies semantic coloring and emoji icons, and handles comments and ellipsis markers — producing a clean, navigable visual representation of directory structures.

---

## Auto-Detection Logic

**Source:** `highlighter.ts` → `looksLikeTree(code)`

A code block is classified as a tree structure if **any** of these three patterns match:

| Pattern | Regex | Matches |
|---------|-------|---------|
| Box-drawing characters | `/[├└│─]/` | `├── src/`, `└── index.ts` |
| Trailing-slash directory | `/^\s*[A-Za-z0-9{}._-]+\/$/m` | `src/`, `components/` |
| File with extension | `/^\s*[A-Za-z0-9{}._-]+\.[A-Za-z0-9_-]+\s*$/m` | `index.ts`, `styles.css` |

**Priority:** Tree detection runs **before** highlight.js. If `looksLikeTree()` returns `true`, the block bypasses language-based highlighting entirely and uses `highlightAsTree()`.

**Language badge:** Tree blocks display **"Structure"** as their language badge (mapped from the `tree` key in `LANGUAGE_NAMES`).

---

## Line-by-Line Rendering Pipeline

**Source:** `highlighter.ts` → `highlightTreeLine(line)`

Each line is processed through a 5-step pipeline, applied in strict order:

### Step 1 — Comment Extraction

```
const commentIndex = line.indexOf("#");
```

- If `#` is found, the line is split into **content** (before `#`) and **comment** (from `#` onward)
- Trailing whitespace on the content portion is trimmed
- The comment is rendered separately at the end

### Step 2 — HTML Escaping

The content portion is passed through `escapeHtml()` to convert `&`, `<`, `>` to entities. This runs **before** any regex replacements to prevent injection.

### Step 3 — Regex Replacement Chain

Applied in this exact order on the escaped content:

| Order | Target | Regex | Output |
|-------|--------|-------|--------|
| 1 | Guide characters | `/([├└│─┌┐┘┬┴┤┼]+)/g` | `<span class="tree-guide">$1</span>` |
| 2 | Ellipsis | `/(\.\.\.)/g` | `<span class="tree-ellipsis">$1</span>` |
| 3 | Directories | `/([A-Za-z0-9{}._-]+\/)/g` | `<span class="tree-dir">📁 $1</span>` |
| 4 | Files | `/([A-Za-z0-9{}._-]+\.[A-Za-z0-9._-]+)/g` | `<span class="tree-file">📄 $1</span>` |

> **Order matters:** Guide characters must be wrapped first so they are not consumed by the directory regex. Ellipsis must be wrapped before files to prevent `...` being matched as a filename.

### Step 4 — Comment Append

If a comment was extracted in Step 1, it is appended:

```html
<highlighted-content> <span class="tree-comment"># comment text</span>
```

The comment text is also HTML-escaped before wrapping.

### Step 5 — Line Assembly

All lines are joined with `\n` and returned as the full highlighted block:

```typescript
function highlightAsTree(code: string): string {
  return code.split("\n").map(highlightTreeLine).join("\n");
}
```

---

## Icon Rendering

### Directory Icon — 📁

- Prepended to any token matching the directory regex (trailing `/`)
- Rendered **inside** the `<span class="tree-dir">` wrapper
- Example: `📁 components/`

### File Icon — 📄

- Prepended to any token matching the file regex (contains `.` extension)
- Rendered **inside** the `<span class="tree-file">` wrapper
- Example: `📄 index.ts`

### Icon Placement

Icons are inserted as literal Unicode emoji characters in the HTML output, not as CSS pseudo-elements. This ensures they are:

- Visible in copy/paste operations
- Selectable with text
- Rendered consistently across platforms

---

## CSS Classes & Styling

**Source:** `src/index.css` — code block section

All tree classes are scoped under `.code-content`:

| Class | Property | Value | Purpose |
|-------|----------|-------|---------|
| `.tree-guide` | `color` | `hsl(var(--muted-foreground) / 0.5)` | Subdued box-drawing lines |
| `.tree-dir` | `color` | `hsl(var(--foreground))` | Full-brightness directory names |
| `.tree-dir` | `font-weight` | `600` | Bold directories for emphasis |
| `.tree-file` | `color` | `hsl(var(--foreground) / 0.85)` | Slightly muted file names |
| `.tree-ellipsis` | `color` | `hsl(var(--accent))` | Accent-colored `...` markers |
| `.tree-comment` | `color` | `hsl(var(--muted-foreground))` | Muted comment text |
| `.tree-comment` | `font-style` | `italic` | Italic for visual distinction |

### Color Design Rationale

- **Neutral palette:** Tree structures use foreground/muted tones only — no red, pink, or purple syntax colors
- **Visual hierarchy:** Directories (bold, full opacity) → Files (85% opacity) → Guides (50% opacity) → Comments (muted italic)
- **Accent reserved for `...`:** The ellipsis marker uses the accent color to draw attention to omitted content

---

## Nesting Rules

### Indentation

Tree nesting is conveyed through **whitespace indentation** in the source markdown. The renderer does not interpret or modify indentation — it preserves the original spacing exactly as authored.

### Box-Drawing Characters

Standard Unicode box-drawing characters define the tree structure:

| Character | Unicode | Name | Usage |
|-----------|---------|------|-------|
| `├` | U+251C | Box Drawings Light Vertical and Right | Branch with siblings below |
| `└` | U+2514 | Box Drawings Light Up and Right | Last branch in a group |
| `│` | U+2502 | Box Drawings Light Vertical | Continuation line |
| `─` | U+2500 | Box Drawings Light Horizontal | Horizontal connector |
| `┌` | U+250C | Box Drawings Light Down and Right | Top-left corner |
| `┐` | U+2510 | Box Drawings Light Down and Left | Top-right corner |
| `┘` | U+2518 | Box Drawings Light Up and Left | Bottom-right corner |
| `┬` | U+252C | Box Drawings Light Down and Horizontal | T-junction down |
| `┴` | U+2534 | Box Drawings Light Up and Horizontal | T-junction up |
| `┤` | U+2524 | Box Drawings Light Vertical and Left | T-junction left |
| `┼` | U+253C | Box Drawings Light Vertical and Horizontal | Cross junction |

All of these are matched by the guide regex and wrapped in `.tree-guide`.

### Typical Nesting Pattern

```
project/
├── src/
│   ├── components/
│   │   ├── Header.tsx
│   │   └── Footer.tsx
│   └── index.ts
├── package.json
└── readme.md
```

---

## Edge Cases

### 1 — Mixed Content (Tree + Non-Tree)

If a code block contains both tree characters and non-tree text (e.g., explanatory paragraphs), the **entire block** is rendered as a tree. The detection is binary — there is no partial tree mode.

### 2 — Filenames with Multiple Dots

Files like `app.config.ts` or `styles.module.css` are correctly matched by the file regex since it matches the **last** occurrence of `.` in a continuous token.

### 3 — Curly Braces in Names

The regex supports `{}` in filenames/directories for placeholder syntax:

- `{id}/` → rendered as directory `📁 {id}/`
- `{slug}.tsx` → rendered as file `📄 {slug}.tsx`

### 4 — Lines with Only Whitespace

Blank or whitespace-only lines pass through the pipeline unchanged — no icons or spans are added.

### 5 — Hash in File/Directory Names

Since `#` is used as the comment delimiter, a filename containing `#` (e.g., `C#-notes.md`) would have everything after the `#` treated as a comment. **Workaround:** This is a known limitation; avoid `#` in tree-rendered filenames.

### 6 — Ellipsis as Omission Marker

The `...` pattern is commonly used to indicate omitted content:

```
src/
├── components/
│   ├── ...
│   └── Footer.tsx
```

The `...` receives accent coloring to visually signal "more items exist here."

### 7 — Explicit `tree` Language Tag

Authors can force tree rendering by using ` ```tree ` as the language tag. This bypasses auto-detection and directly invokes `highlightAsTree()`.

### 8 — Comments After File/Directory Names

Inline comments are supported and common:

```
├── constants.ts    # Language maps, accent colors
└── types.ts        # Shared interfaces
```

The content before `#` is rendered with icons; the comment after is rendered in `.tree-comment` style.

### 9 — Root Directory on First Line

A standalone directory name on the first line (e.g., `project/`) is rendered as a directory with 📁 icon and bold styling, serving as the visual root of the tree.

### 10 — No Box-Drawing Characters

A block with only trailing-slash directories or file extensions (no `├└│─`) still triggers tree detection. Example:

```
src/
  components/
    Header.tsx
    Footer.tsx
  index.ts
```

This renders with icons but without guide-line styling (since there are no box-drawing characters to wrap).

---

## Interaction Behavior

Tree blocks retain all standard code block interactions:

- **Copy:** Copies raw text including box-drawing characters but **without** emoji icons (icons are HTML-injected, not in source)
- **Download:** Downloads as `.txt` file with original source text
- **Fullscreen:** Expands tree with same styling
- **Font controls:** A-/A/A+ adjust tree text size identically to code
- **Line selection:** Click/Shift+click/drag on line numbers works as with any code block

---

## Source Files

| File | Role |
|------|------|
| `src/components/markdown/highlighter.ts` | Detection (`looksLikeTree`), rendering (`highlightTreeLine`, `highlightAsTree`) |
| `src/components/markdown/constants.ts` | `"tree": "Structure"` language name mapping |
| `src/index.css` | `.tree-guide`, `.tree-dir`, `.tree-file`, `.tree-ellipsis`, `.tree-comment` styles |

---

*Tree-Structure Rendering — updated: 2026-04-08*
