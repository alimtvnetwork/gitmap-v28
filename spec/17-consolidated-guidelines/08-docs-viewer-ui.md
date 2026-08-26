# Consolidated: Docs Viewer UI — Complete Reference

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Purpose

This is the **standalone consolidated reference** for the interactive documentation viewer at `/docs`. An AI reading only this file must be able to reproduce the complete UI — layout, typography, color system, keyboard navigation, animations, and all component rendering.

---

## Architecture

| Layer | Technology |
|-------|-----------|
| Framework | React 18 + TypeScript 5 |
| Styling | Tailwind CSS v3 + shadcn/ui |
| Fonts | Google Fonts — Ubuntu (headings), Poppins (body), JetBrains Mono (code) |
| Syntax Highlighting | highlight.js with custom theme |
| Layout | Sidebar (shadcn) + flex content area |
| State | React hooks (useState, useCallback, useEffect) |
| Data Source | `src/data/specTree.json` — tree of nodes with `name`, `type`, `path`, `content`, `children` |

---

## Typography System

| Element | Font Family | Weight | Size | Fallback |
|---------|-------------|--------|------|----------|
| H1 | Ubuntu | 700 | 1.6rem | sans-serif |
| H2 | Ubuntu | 700 | 1.25rem | sans-serif |
| H3 | Ubuntu | 600 | 1.05rem | sans-serif |
| H4 | Ubuntu | 600 | 0.95rem | sans-serif |
| Body text | Poppins | 400 | 0.9rem | sans-serif |
| Sidebar / nav | Poppins | 400, 500 | — | sans-serif |
| Inline code | JetBrains Mono, Fira Code | 500 | 0.85em | ui-monospace |
| Code blocks | Ubuntu Mono, JetBrains Mono | 400 | 18px | monospace |

Font loading via Google Fonts `<link>` in `index.html`:
- Ubuntu: 400, 500, 600, 700
- Poppins: 300, 400, 500, 600
- `font-display: swap` for performance

---

## Color System — CSS Custom Properties (HSL)

All colors stored as HSL channels in CSS custom properties. Components reference via `hsl(var(--token))`.

### Core Tokens

| Token | Light Mode | Dark Mode | Usage |
|-------|-----------|-----------|-------|
| `--primary` | `252 85% 60%` | `252 85% 65%` | Brand purple — buttons, accents, badges |
| `--accent` | `330 85% 60%` | `330 85% 65%` | Accent pink — gradient endpoints, link hovers |
| `--background` | `0 0% 100%` | `230 25% 8%` | Page background |
| `--foreground` | `230 25% 15%` | `220 20% 92%` | Primary text |
| `--muted` | `220 20% 96%` | `230 15% 18%` | Subtle backgrounds, disabled states |
| `--muted-foreground` | `220 10% 46%` | `220 10% 60%` | Secondary text, metadata |
| `--border` | `220 20% 90%` | `220 13% 22%` | Borders, dividers |
| `--card` | `0 0% 100%` | `230 20% 12%` | Card/panel backgrounds |
| `--success` | `152 70% 42%` | `152 70% 42%` | Checked states, positive indicators |
| `--warning` | `38 92% 50%` | `38 92% 50%` | Warnings |
| `--destructive` | `0 84% 60%` | `0 84% 60%` | Error states |

### Reading-Specific Tokens

| Token | Light | Dark | Usage |
|-------|-------|------|-------|
| `--heading-gradient-from` | `252 85% 60%` | `252 85% 70%` | Heading gradient start (purple) |
| `--heading-gradient-to` | `330 85% 60%` | `330 85% 70%` | Heading gradient end (pink) |
| `--link-color` | `252 85% 55%` | `252 85% 72%` | Link text color |
| `--code-bg` | `250 25% 95%` | `230 20% 15%` | Inline code background |
| `--code-text` | `330 85% 45%` | `330 85% 70%` | Inline code text |
| `--highlight-glow` | `252 85% 60%` | `252 85% 65%` | Hover glow effects |
| `--table-header-bg` | `252 85% 97%` | `230 20% 14%` | Table header background |
| `--table-row-hover` | — | — | Table row hover tint |

---

## Keyboard Navigation — Complete Reference

### Key Bindings

| Key | Action | Scope |
|-----|--------|-------|
| `→` Right Arrow | Next file within current folder | File-level |
| `←` Left Arrow | Previous file within current folder | File-level |
| `↓` Down Arrow | Next folder (jump to first file) | Folder-level |
| `↑` Up Arrow | Previous folder (jump to first file) | Folder-level |
| `Enter` | Open selected file | Sidebar focus |
| `F` | Toggle fullscreen mode | Global |
| `?` | Toggle shortcuts help overlay | Global |
| `Escape` | Exit fullscreen / close overlay / close search | Global |
| `Ctrl+K` | Focus search input | Global |

### Navigation Logic

- Files ordered by numeric prefix within each folder
- Folder navigation follows `specTree.json` order
- **Wrap-around**: last file → first file (within folder); last folder → first folder (at root level)
- Only active when no `<input>`, `<textarea>`, or `[contenteditable]` is focused
- Only active when a file is currently selected (no-op on welcome screen)

### Implementation Guard

```typescript
const isInputFocused = () => {
  const el = document.activeElement;
  return el?.tagName === 'INPUT' || el?.tagName === 'TEXTAREA' || 
         el?.getAttribute('contenteditable') === 'true';
};
```

---

## Shortcuts Help Overlay

### Trigger

- `?` key (when no input focused)
- Keyboard icon button in the header

### Content — Grid Layout

| Category | Key | Action |
|----------|-----|--------|
| Navigation | `←` / `→` | Previous / Next file |
| Navigation | `↑` / `↓` | Previous / Next folder |
| Actions | `Ctrl+K` | Focus search |
| View | `F` | Toggle fullscreen |
| View | `?` | Toggle this help |
| View | `Escape` | Close overlay / exit fullscreen |

### Dismiss

- `Escape` key
- Clicking outside the overlay
- `?` key again (toggle)

### Visual

- Modal with backdrop blur (`backdrop-filter: blur(4px)`)
- Grouped by category (Navigation, Actions, View)
- Key badges styled with `--muted` background and `--border` border
- `z-index: 1000` above all content

---

## Fullscreen Mode

- Toggle button on doc content area
- Hides sidebar and header for distraction-free reading
- `Escape` key exits fullscreen
- `F` key toggles
- Preserves scroll position on toggle
- `position: fixed; inset: 0; z-index: 50`

---

## Copy Markdown

- Button to copy raw markdown source to clipboard
- Uses Clipboard API (`navigator.clipboard.writeText()`) with `document.execCommand('copy')` fallback
- Toast notification on success ("Copied!") or failure ("Failed to copy")
- Button shows checkmark (✓) for 2 seconds after successful copy

---

## Sidebar Tree Structure

```
Sidebar
├── Search Input (Ctrl+K)
├── Folder (expandable)
│   ├── 00-overview.md (file)
│   ├── 01-fundamentals.md (file)
│   └── 02-features/ (subfolder)
│       ├── 00-overview.md
│       └── 01-feature.md
└── Folder 2 (expandable)
    └── ...
```

- Folders show file count badge
- Active file highlighted with accent color left border + background tint
- Collapse state persisted in `localStorage`
- Files sorted by numeric prefix
- Expand/collapse: chevron rotation `200ms ease-out`

### Data Source — `specTree.json`

Each node requires:
- `name`: display name (string)
- `type`: `"file"` or `"folder"` (string)
- `path`: relative path to `.md` file (string)
- `content`: full markdown text (string, required for rendering)
- `children`: array of child nodes (folders only)

---

## Animation & Interaction Patterns

### Heading Hover — Gradient Brightness (H1/H2)

```css
.spec-h1, .spec-h2 {
  background: linear-gradient(135deg, hsl(var(--heading-gradient-from)), hsl(var(--heading-gradient-to)));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  transition: filter 0.3s ease;
}
.spec-h1:hover, .spec-h2:hover {
  filter: brightness(1.2) saturate(1.1);
}
```

### H3 — Left Border Slide

```css
.spec-h3 {
  padding-left: 0.65rem;
  border-left: 3px solid hsl(var(--primary) / 0.5);
  transition: color 0.2s ease, border-color 0.2s ease, padding-left 0.2s ease;
}
.spec-h3:hover {
  color: hsl(var(--primary));
  border-color: hsl(var(--primary));
  padding-left: 0.85rem;
}
```

### H4 — Subtle Color Shift

```css
.spec-h4:hover { color: hsl(var(--foreground)); }
```

### Paragraph Hover — Subtle Highlight

```css
.spec-p {
  padding: 0.1rem 0.35rem;
  border-radius: 3px;
  transition: color 0.15s ease, background 0.2s ease;
}
.spec-p:hover {
  color: hsl(var(--foreground));
  background: hsl(var(--primary) / 0.04);
}
```

### Link Hover — Underline Sweep

```css
.spec-link::after {
  content: '';
  position: absolute;
  width: 100%;
  height: 2px;
  bottom: -2px;
  left: 0;
  background: linear-gradient(90deg, hsl(var(--heading-gradient-from)), hsl(var(--heading-gradient-to)));
  transform: scaleX(0);
  transform-origin: bottom right;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.spec-link:hover::after {
  transform: scaleX(1);
  transform-origin: bottom left;
}
.spec-link:hover {
  color: hsl(var(--accent));
}
```

### Inline Code Hover — Lift + Glow

```css
.inline-code:hover {
  box-shadow: 0 0 0 2px hsl(var(--highlight-glow) / 0.15);
  transform: translateY(-1px);
}
```

### List Item Hover — Slide + Bullet Grow

```css
.spec-li:hover {
  transform: translateX(3px);
  background: hsl(var(--primary) / 0.04);
}
.spec-li:hover::before {
  background: hsl(var(--primary));
  transform: scale(1.3);
}
```

### Table Row Hover — Highlight + Left Bar

```css
tbody tr:hover {
  background: hsl(var(--table-row-hover));
  box-shadow: inset 3px 0 0 hsl(var(--primary) / 0.5);
}
```

### Code Block Hover — Float + Glow

```css
.code-block-wrapper:hover {
  box-shadow: 0 8px 32px hsl(var(--lang-accent) / 0.1), 0 0 0 1px hsl(var(--lang-accent) / 0.15);
  transform: translateY(-2px);
}
```

### Blockquote Hover — Slide + Shadow

```css
.spec-blockquote:hover {
  background: hsl(var(--muted) / 0.5);
  transform: translateX(3px);
  box-shadow: -4px 0 12px hsl(var(--heading-gradient-from) / 0.1);
}
```

### Sidebar Item Hover

- Background: subtle highlight with `150ms` transition
- Active file: accent color left border + background tint

### Content Transition

- Fade-in on file switch: `200ms opacity` transition

---

## Heading Properties

| Property | H1 | H2 | H3 | H4 |
|----------|----|----|----|----|
| Font size | 1.6rem | 1.25rem | 1.05rem | 0.95rem |
| Weight | 700 | 700 | 600 | 600 |
| Margin | `1rem 0 0.6rem` | `1.8rem 0 0.5rem` | — | — |
| Bottom border | None | `1px solid hsl(var(--border))` | None | None |
| Left border | None | None | `3px solid hsl(var(--primary) / 0.5)` | None |
| Hover | brightness(1.2) | brightness(1.2) | border slide | color shift |

---

## Visual Rendering — Component Guide

### Inline Code

```css
.inline-code {
  background: hsl(var(--code-bg));
  color: hsl(var(--code-text));
  padding: 0.2em 0.45em;
  border-radius: 5px;
  font-size: 0.85em;
  font-family: 'JetBrains Mono', 'Fira Code', ui-monospace, monospace;
  font-weight: 500;
  border: 1px solid hsl(var(--border) / 0.5);
  transition: background 0.2s ease, box-shadow 0.2s ease, transform 0.15s ease;
}
```

### Bold & Italic

```css
strong { color: hsl(var(--foreground)); font-weight: 700; }
em     { color: hsl(var(--muted-foreground)); }
```

### Paragraphs

- `line-height: 1.65`
- `color: hsl(var(--foreground) / 0.9)`

### Tables

- Bordered, striped rows (even rows get muted background)
- Responsive horizontal scroll wrapper
- Header: `--table-header-bg` background, `font-weight: 600`
- Row hover: highlight + 3px left bar in primary color

### Blockquotes

- Left accent border: `4px solid hsl(var(--primary) / 0.4)`
- Background: `hsl(var(--muted) / 0.3)`
- Hover: slides right 3px with shadow

### Checklists

- `[ ]` → bordered empty checkbox
- `[x]` → green gradient checkbox using `--success` color
- Wrapped in `.checklist-block` with header bar ("Checklist" label + copy button)
- Hover slides item right with subtle highlight

### Horizontal Rules

- `border-top: 1px solid hsl(var(--border))`
- `margin: 1.5rem 0`

### Links

- Color: `hsl(var(--link-color))` — purple
- Hover: shifts to `hsl(var(--accent))` — pink
- Underline sweep animation (see animations section)

### Images

- `max-width: 100%`
- `border-radius: 0.5rem`
- Optional caption below

---

## Code Block Language Colors

Each language has a unique HSL accent used for badge, glow, and hover:

| Language | Badge Label | HSL Accent |
|----------|-------------|-----------|
| TypeScript/TSX | `TYPESCRIPT` | `99 83% 62%` |
| JavaScript | `JAVASCRIPT` | `53 93% 54%` |
| Go | `GO` | `194 66% 55%` |
| PHP | `PHP` | `234 45% 60%` |
| CSS | `CSS` | `264 55% 58%` |
| JSON | `JSON` | `38 92% 50%` |
| Bash/Shell | `BASH` | `120 40% 55%` |
| SQL | `SQL` | `200 70% 55%` |
| Rust | `RUST` | `25 85% 55%` |
| HTML/XML | `HTML` | `12 80% 55%` |
| YAML | `YAML` | `0 75% 55%` |
| Markdown | `MARKDOWN` | `252 85% 60%` |
| Plain Text | `PLAIN TEXT` | `220 10% 50%` |

---

## Tree Structure Rendering

### Detection Logic

Code blocks auto-detected as tree when content matches:
- Unicode box-drawing chars: `/[├└│─]/`
- Lines ending with `/` (directories): `/^\s*[A-Za-z0-9{}._-]+\/$/m`
- Lines with file extensions: `/^\s*[A-Za-z0-9{}._-]+\.[A-Za-z0-9_-]+\s*$/m`
- Explicit fence labels: ` ```tree ` or ` ```structure `

### Line Rendering Rules

| Pattern | CSS Class | Visual |
|---------|-----------|--------|
| Box-drawing (`├ └ │ ─`) | `.tree-guide` | Muted 50% opacity |
| `...` (ellipsis) | `.tree-ellipsis` | Accent pink |
| `name/` (directory) | `.tree-dir` | 📁 prefix, bold, `font-weight: 600` |
| `name.ext` (file) | `.tree-file` | 📄 prefix, 85% opacity |
| `# comment` | `.tree-comment` | Italic, muted |

### Key Design Decision

Tree blocks use a **neutral/white color scheme**, avoiding syntax highlighting colors that conflict with directory/file readability.

---

## Implementation Rules

1. **Never use raw color values** — always reference CSS custom properties
2. **All transitions use `ease` or `cubic-bezier(0.4, 0, 0.2, 1)`**
3. **Hover effects combine 2–3 properties** (color + transform + shadow)
4. **Duration scale**: micro 0.15s, standard 0.2–0.3s, emphasis 0.3–0.5s
5. **Transform patterns**: `translateX(3px)` horizontal slide, `translateY(-2px)` float, `scale(1.1–1.3)` emphasis
6. **Gradient direction**: `135deg` diagonal, `90deg` horizontal sweeps
7. **Glow formula**: `box-shadow: 0 0 Npx hsl(var(--token) / 0.1–0.25)`

---

## Syntax Highlighting

- Engine: **highlight.js** with custom theme matching design system tokens
- Supported languages: TypeScript, Go, PHP, Rust, SQL, JSON, YAML, Bash, PowerShell, CSS, HTML, Markdown, INI
- Language detection: explicit fence label (```` ```go ````) or auto-detect
- Code blocks always use fixed dark background regardless of app theme

### Token Colors (highlight.js overrides)

| Token Type | CSS Selector | Color |
|------------|-------------|-------|
| Keywords, types, built-ins | `.hljs-keyword`, `.hljs-type`, `.hljs-built_in` | `hsl(var(--primary))` — purple |
| Strings, attributes | `.hljs-string`, `.hljs-attr`, `.hljs-property` | `hsl(var(--accent))` — pink |
| Numbers, variables | `.hljs-number`, `.hljs-variable`, `.hljs-regexp` | `hsl(var(--warning))` — amber |
| Comments | `.hljs-comment`, `.hljs-quote` | `hsl(var(--muted-foreground))` italic |
| Functions, classes, tags | `.hljs-title`, `.hljs-section`, `.hljs-tag` | `hsl(var(--foreground) / 0.85)` |
| Default text | `code` | `hsl(var(--foreground))` |

---

*Consolidated docs viewer UI — v3.2.0 — 2026-04-16*
