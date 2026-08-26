# Typography

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

Typography rules define font families, size hierarchy, weight usage, and text spacing. All typographic choices serve readability and visual hierarchy.

---

## Font Stacks

| Role | Primary Font | Fallbacks | CSS Variable |
|------|-------------|-----------|-------------|
| **Headings** | Ubuntu | sans-serif | `font-heading` |
| **Body Text** | Poppins | sans-serif | `font-body` |
| **Code (Inline)** | JetBrains Mono | Fira Code, ui-monospace, monospace | — |
| **Code (Blocks)** | Ubuntu Mono | JetBrains Mono, ui-monospace, monospace | — |
| **Keyboard** | Ubuntu | monospace | — |

### Loading Strategy

Fonts are loaded via Google Fonts in the HTML `<head>`:

```html
<link href="https://fonts.googleapis.com/css2?family=Ubuntu:wght@400;500;700&family=Poppins:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;500;700&family=Ubuntu+Mono:wght@400;700&display=swap" rel="stylesheet">
```

---

## Size Hierarchy

### Prose Content (Markdown Rendering)

| Element | Size | Weight | Additional |
|---------|------|--------|------------|
| H1 (`.spec-h1`) | `1.6rem` | 700 | Gradient text, letter-spacing: `-0.02em` |
| H2 (`.spec-h2`) | `1.25rem` | 700 | Gradient text, bottom border |
| H3 (`.spec-h3`) | `1.05rem` | 600 | Left border accent, color transition |
| H4 (`.spec-h4`) | `0.95rem` | 600 | Muted color, hover brightens |
| Body paragraph | `0.9rem` | 400 | Line-height: `1.65` |
| Inline code | `0.85em` | 500 | Monospace, background pill |
| Table text | `0.82rem` | 400 | — |
| Table headers | `0.75rem` | 600 | Uppercase, letter-spacing: `0.03em` |

### Fullscreen Mode

When entering fullscreen reading mode, sizes scale up:

| Element | Fullscreen Size |
|---------|----------------|
| H1 | `2rem` |
| H2 | `1.5rem` |
| H3 | `1.2rem` |
| Body | `1.05rem` |
| Code content | `20px` |

### Code Block Typography

| Element | Size | Weight |
|---------|------|--------|
| Language badge | `0.7rem` | 600 |
| Line count | `0.65rem` | 400 |
| Tool buttons | `0.65rem` | 500 |
| Font controls | `0.6rem` | 700 |
| Selection label | `0.6rem` | 600 |
| Copy label | `0.7rem` | 600 |
| Code content | `var(--code-font-size)` (default `18px`) | 400 |
| Line numbers | `calc(var(--code-font-size) * 0.7)` | 400 |

---

## Weight Usage

| Weight | Name | Usage |
|--------|------|-------|
| 300 | Light | Rarely used — only for very large decorative text |
| 400 | Regular | Body text, paragraphs, table cells, code |
| 500 | Medium | Links, inline code, subtle emphasis |
| 600 | SemiBold | H3, H4, labels, table headers, badges |
| 700 | Bold | H1, H2, strong text, font controls |

---

## Text Spacing

| Property | Value | Applied To |
|----------|-------|------------|
| Letter-spacing `-0.02em` | Tighter | H1 headings only |
| Letter-spacing `0.03em` | Wider | Table headers (uppercase) |
| Letter-spacing `0.05em` | Widest | Badge labels, checklist titles |
| Letter-spacing `0.02em` | Slight | Selection labels |
| Line-height `1.65` | Relaxed | Body paragraphs |
| Line-height `1.55` | Standard | List items |
| Line-height `var(--code-line-height)` (`1.6`) | Code | Code blocks |

---

## Text Transform Rules

| Pattern | Transform | Where Used |
|---------|-----------|------------|
| Language badges | `text-transform: uppercase` | Code block headers |
| Table headers | `text-transform: uppercase` | Prose tables |
| Checklist titles | `text-transform: uppercase` | Checklist block headers |
| Everything else | None | Default — preserve original case |

---

## Heading Gradient Effect

H1 and H2 headings use a gradient text effect:

```css
background: linear-gradient(
  135deg,
  hsl(var(--heading-gradient-from)),
  hsl(var(--heading-gradient-to))
);
-webkit-background-clip: text;
background-clip: text;
-webkit-text-fill-color: transparent;
```

The gradient flows from `--heading-gradient-from` (primary hue) to `--heading-gradient-to` (accent hue) at a 135° angle. Changing these two tokens re-themes all heading gradients.

**Hover behavior:** `filter: brightness(1.2) saturate(1.1)` — subtle brightening.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Theme Variables | [02-theme-variable-architecture.md](./02-theme-variable-architecture.md) |
| Code Block Specifics | [07-code-blocks.md](./07-code-blocks.md) |
