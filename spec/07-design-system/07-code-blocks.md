# Code Blocks

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

Code blocks are the most complex visual component in the design system. They maintain a **fixed dark appearance** regardless of light/dark theme, feature per-language accent colors, interactive line selection, font controls, and a fullscreen mode.

---

## Structure

```
┌─────────────────────────────────────────────────┐
│ ● LANG_BADGE    N lines  A- A A+  Copy  ⇲  ⛶  │  ← Header
├────┬────────────────────────────────────────────┤
│  1 │ const x = 42;                              │  ← Body
│  2 │ console.log(x);                             │
│  3 │                                             │
├────┴────────────────────────────────────────────┤
│ Lines 1-2 selected          Copy Selected  ✕    │  ← Selection Bar
└─────────────────────────────────────────────────┘
```

### HTML Structure

```html
<div class="code-block-wrapper" style="--lang-accent: H S% L%" data-block-id="N">
  <!-- Header -->
  <div class="code-block-header">
    <div class="code-lang-badge">
      <span class="code-lang-dot"></span>
      LANGUAGE
    </div>
    <div class="code-header-right">
      <span class="code-line-count">N lines</span>
      <div class="code-font-controls">
        <button class="code-tool-btn">A-</button>
        <button class="code-tool-btn">A</button>
        <button class="code-tool-btn">A+</button>
      </div>
      <button class="code-tool-btn">Copy</button>
      <button class="code-tool-btn">Download</button>
      <button class="code-tool-btn">Fullscreen</button>
    </div>
  </div>

  <!-- Body -->
  <div class="code-block-body">
    <pre class="code-line-numbers">...</pre>
    <pre class="code-content"><code class="hljs language-X">...</code></pre>
  </div>

  <!-- Selection Bar (conditional) -->
  <div class="copy-selected-bar">...</div>
</div>
```

---

## Color System (Fixed, Not Themed)

Code blocks always use a dark background to ensure code readability. These values are **not** CSS custom properties — they are fixed HSL values.

| Element | HSL Value | Hex Equivalent |
|---------|-----------|---------------|
| Block background | `220, 14%, 11%` | `#181c24` |
| Header background | `220, 14%, 14%` | `#1f232b` |
| Header border | `220, 13%, 20%` | `#2c3038` |
| Block outer border | `220, 13%, 22%` | `#30353e` |
| Line number background | `220, 14%, 9%` | `#141820` |
| Line number border | `220, 13%, 18%` | `#272b33` |
| Line number text | `220, 10%, 35%` | `#535862` |
| Line hover background | `220, 15%, 16%` | `#232830` |
| Pinned line background | `var(--primary) / 0.12` | Themed |
| Tool button background | `220, 13%, 20%` | `#2c3038` |
| Tool button border | `220, 13%, 25%` | `#373c45` |
| Tool button hover | `220, 13%, 28%` | `#3e434c` |
| Font controls background | `220, 13%, 18%` | `#272b33` |

### Themed Elements Within Code Blocks

These elements reference CSS custom properties:

| Element | Token Used |
|---------|-----------|
| Syntax: keywords | `hsl(var(--primary))` |
| Syntax: strings/attributes | `hsl(var(--accent))` |
| Syntax: numbers | `hsl(var(--warning))` |
| Syntax: comments | `hsl(var(--muted-foreground))` |
| Pinned line background | `hsl(var(--primary) / 0.12)` |
| Pinned line number color | `hsl(var(--primary))` |
| Selection label | `hsl(var(--primary))` |
| Copy bar background | `hsl(var(--primary) / 0.08)` |
| Copy bar border | `hsl(var(--primary) / 0.2)` |

---

## Language Badge

The badge displays the detected language with a colored dot:

```css
.code-lang-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: hsl(var(--badge-color));
  box-shadow: 0 0 6px hsl(var(--badge-color) / 0.5);
}
```

Each code block sets `--lang-accent` inline, which drives:
- Badge dot color and glow
- Hover glow on the entire block
- Fullscreen shadow color
- Tool button hover border

---

## Interactions

### Line Click (Pin)

- Single click on a line number → toggles **pinned** state
- Pinned lines get `line-pinned` class → primary-tinted background
- Line number gets primary color + left border accent

### Shift-Click (Range)

- Shift-click extends selection from last pinned line to clicked line
- All lines in range get `line-pinned` class

### Click-Drag (Multi-Select)

- Click and drag across line numbers selects a range
- Selection bar appears at bottom with "Lines X–Y" label

### Selection Bar

Appears when lines are pinned:

```
┌────────────────────────────────────────────┐
│ Lines 4–9          Copy Selected      ✕    │
└────────────────────────────────────────────┘
```

- Animated in via `@keyframes slideUpBar`
- Copy button copies only selected lines
- Clear button (✕) removes all pins

### Font Size Controls

Three buttons: `A-`, `A`, `A+`

| Button | Action |
|--------|--------|
| A- | Decrease `--code-font-size` by 2px (min: 12px) |
| A | Reset to default 18px |
| A+ | Increase `--code-font-size` by 2px (max: 32px) |

Line numbers scale proportionally: `calc(var(--code-font-size) * 0.7)`.

### Copy Button

- Copies escaped code content to clipboard
- Button state changes to `.copied` class:
  - Background: `hsl(152, 60%, 18%)` (success tint)
  - Color: `hsl(152, 70%, 60%)`
  - Text: "Copied ✓"
- Reverts after 2 seconds

### Download Button

- Downloads code as a file with language-appropriate extension
- Filename: `code-{id}.{ext}`

### Fullscreen Button

- Expands block to `position: fixed; inset: 2rem`
- Overlay: `hsl(0 0% 0% / 0.7)` with `backdrop-filter: blur(4px)`
- Body receives overflow lock
- Escape key or overlay click exits fullscreen
- Shadow: `0 25px 80px hsl(var(--lang-accent) / 0.25)`

---

## Hover Effects

### Block Wrapper Hover

```css
box-shadow: 0 8px 32px hsl(var(--lang-accent) / 0.1),
            0 0 0 1px hsl(var(--lang-accent) / 0.15);
transform: translateY(-2px);
transition: box-shadow 0.3s ease, transform 0.2s ease;
```

### Tool Button Hover

```css
background: hsl(220, 13%, 28%);
color: hsl(0, 0%, 95%);
border-color: hsl(var(--lang-accent) / 0.4);
box-shadow: 0 0 8px hsl(var(--lang-accent) / 0.15);
transition: all 0.2s ease;
```

### Line Hover

```css
background: hsl(220, 15%, 16%);
transition: background 0.15s ease;
```

---

## Tree / Structure Rendering

When code content is detected as a file/folder tree (box-drawing characters, directory patterns), special rendering applies:

| Element | Rendering |
|---------|-----------|
| Directories (`name/`) | Prefixed with 📁, bold, `--foreground` color |
| Files (`name.ext`) | Prefixed with 📄, `--foreground / 0.85` color |
| Tree guides (`├└│─`) | `--muted-foreground / 0.5` — very subtle |
| Ellipsis (`...`) | `--accent` color |
| Comments (`# text`) | `--muted-foreground`, italic |

Tree blocks use the language label "STRUCTURE" and bypass syntax highlighting.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Theme Variables | [02-theme-variable-architecture.md](./02-theme-variable-architecture.md) |
| Motion Transitions | [06-motion-transitions.md](./06-motion-transitions.md) |
| Visual Rendering Guide | `../08-docs-viewer-ui/02-features/07-visual-rendering-guide.md` |
| Source: codeBlockBuilder | `src/components/markdown/codeBlockBuilder.ts` |
| Source: highlighter | `src/components/markdown/highlighter.ts` |
| Source: CSS styles | `src/index.css` (lines 264–605) |
