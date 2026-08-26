# Feature: UI Theme, Colors & Animations

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

This document defines the complete color theme, animation system, and interactive styling used in the docs viewer UI. Any AI or developer can use this as a reference to replicate or extend the visual design.

---

## Color System

### Design Tokens (CSS Custom Properties)

All colors use **HSL format** stored as CSS custom properties in `src/index.css`. Components reference these via `hsl(var(--token))` — never hardcoded color values.

#### Light Mode

| Token | HSL Value | Usage |
|-------|-----------|-------|
| `--primary` | `252 85% 60%` | Primary brand purple — buttons, accents, badges |
| `--accent` | `330 85% 60%` | Accent pink — gradient endpoints, link hovers |
| `--background` | `0 0% 100%` | Page background |
| `--foreground` | `230 25% 15%` | Primary text |
| `--muted` | `220 20% 96%` | Subtle backgrounds, disabled states |
| `--muted-foreground` | `220 10% 46%` | Secondary text, metadata |
| `--border` | `220 20% 90%` | Borders, dividers |
| `--success` | `152 70% 42%` | Checked states, positive indicators |
| `--warning` | `38 92% 50%` | Warnings, caution states |
| `--destructive` | `0 84% 60%` | Error states |

#### Dark Mode

| Token | HSL Value | Usage |
|-------|-----------|-------|
| `--primary` | `252 85% 65%` | Brighter purple for dark backgrounds |
| `--accent` | `330 85% 65%` | Brighter pink for dark backgrounds |
| `--background` | `230 25% 8%` | Deep navy-charcoal page background |
| `--foreground` | `220 20% 92%` | Light text on dark |
| `--muted` | `230 15% 18%` | Subtle dark backgrounds |
| `--card` | `230 20% 12%` | Card/panel backgrounds |

#### Reading-Specific Tokens

| Token | Light | Dark | Usage |
|-------|-------|------|-------|
| `--heading-gradient-from` | `252 85% 60%` | `252 85% 70%` | Heading gradient start (purple) |
| `--heading-gradient-to` | `330 85% 60%` | `330 85% 70%` | Heading gradient end (pink) |
| `--link-color` | `252 85% 55%` | `252 85% 72%` | Link text color |
| `--code-bg` | `250 25% 95%` | `230 20% 15%` | Inline code background |
| `--code-text` | `330 85% 45%` | `330 85% 70%` | Inline code text |
| `--highlight-glow` | `252 85% 60%` | `252 85% 65%` | Hover glow effects |
| `--table-header-bg` | `252 85% 97%` | `230 20% 14%` | Table header background |

---

## Typography

| Element | Font Family | Weight | Size |
|---------|-------------|--------|------|
| Headings (h1–h6) | `Ubuntu` | 600–700 | 0.95–1.6rem |
| Body text | `Poppins` (via `font-body`) | 400 | 0.9rem |
| Code (inline) | `JetBrains Mono`, `Fira Code` | 500 | 0.85em |
| Code (blocks) | `Ubuntu Mono`, `JetBrains Mono` | 400 | 18px |

---

## Animation & Interaction Patterns

### Heading Hover Animation

Headings use a **gradient text fill** (purple → pink) with a brightness boost on hover:

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

### H3 Hover — Border Slide

H3 headings have a left accent border that brightens and slides right on hover:

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

### Paragraph Hover — Subtle Highlight

Body text gets a subtle primary-tinted background on hover:

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

Links use a CSS `::after` pseudo-element for an animated underline that sweeps from right to left:

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

Code blocks use a per-language accent color (`--lang-accent`) for hover glow:

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

---

## Code Block Language Colors

Each language has an HSL accent used for the badge, glow, and hover effects:

| Language | HSL Accent |
|----------|-----------|
| TypeScript/TSX | `99, 83%, 62%` |
| JavaScript | `53, 93%, 54%` |
| Go | `194, 66%, 55%` |
| PHP | `234, 45%, 60%` |
| CSS | `264, 55%, 58%` |
| JSON | `38, 92%, 50%` |
| Bash/Shell | `120, 40%, 55%` |
| SQL | `200, 70%, 55%` |
| Rust | `25, 85%, 55%` |
| HTML/XML | `12, 80%, 55%` |
| YAML | `0, 75%, 55%` |

---

## Checklist Rendering

Checklists are wrapped in a `.checklist-block` container with:
- A header bar showing "Checklist" label and a copy button
- Checked items get a green gradient checkbox (`--success` color)
- Unchecked items get a bordered empty box
- Hover slides item right with subtle highlight

---

## Implementation Rules

1. **Never use raw color values** — always reference CSS custom properties
2. **All transitions use `ease` or `cubic-bezier(0.4, 0, 0.2, 1)`** for smooth motion
3. **Hover effects combine 2–3 properties** (color + transform + shadow) for richness
4. **Duration scale**: micro-interactions 0.15s, standard 0.2–0.3s, emphasis 0.3–0.5s
5. **Transform hover patterns**: `translateX(3px)` for horizontal slide, `translateY(-2px)` for float, `scale(1.1–1.3)` for emphasis
6. **Gradient direction**: `135deg` for diagonal, `90deg` for horizontal sweeps
7. **Glow formula**: `box-shadow: 0 0 Npx hsl(var(--token) / 0.1–0.25)`

---

*UI theme and animations — updated: 2026-04-03*
