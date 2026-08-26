# Consolidated: Design System — Complete Reference

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Purpose

This is the **standalone consolidated reference** for the AI-Adaptable Design System. An AI reading only this file must be able to reproduce, extend, and re-theme the entire visual language without consulting source spec files.

---

## Architecture — Variable-Driven

All colors, spacing, borders, and visual tokens are defined as CSS custom properties (HSL format) in `index.css`. Components never use hardcoded colors — they reference semantic tokens. Changing a token propagates to every component.

```
CSS Custom Properties (index.css :root / .dark)
  → Tailwind Config mapping (hsl(var(--token)))
    → Component semantic classes
      → Component states (:hover, :focus, .dark)
        → Page composition
```

### Variable Format

```css
--token-name: H S% L%;
```

Usage in components:

```css
color: hsl(var(--token-name));                /* full opacity */
background: hsl(var(--token-name) / 0.1);     /* 10% opacity */
```

Usage in Tailwind classes:

```html
<div class="bg-primary text-primary-foreground">
```

### Tailwind Config Mapping

Every CSS custom property is mapped to a Tailwind utility class in `tailwind.config.ts`:

```typescript
colors: {
  primary: {
    DEFAULT: "hsl(var(--primary))",
    foreground: "hsl(var(--primary-foreground))",
  },
  accent: {
    DEFAULT: "hsl(var(--accent))",
    foreground: "hsl(var(--accent-foreground))",
  },
  destructive: {
    DEFAULT: "hsl(var(--destructive))",
    foreground: "hsl(var(--destructive-foreground))",
  },
  warning: {
    DEFAULT: "hsl(var(--warning))",
    foreground: "hsl(var(--warning-foreground))",
  },
  success: {
    DEFAULT: "hsl(var(--success))",
    foreground: "hsl(var(--success-foreground))",
  },
  muted: {
    DEFAULT: "hsl(var(--muted))",
    foreground: "hsl(var(--muted-foreground))",
  },
  // ... all tokens follow same pattern
}
```

---

## Design Philosophy

| Principle | Rule |
|-----------|------|
| **Variable-First** | Every color, spacing, visual property from a CSS custom property. No hardcoded values. |
| **Semantic Tokens** | Named by purpose (`--primary`, `--accent`), not by value (`--purple`). |
| **HSL Color Model** | All colors in HSL for easy theme derivation (adjust lightness without changing hue). |
| **CSS3 Motion Only** | All transitions and animations use CSS transitions/transforms/keyframes. No JS animation libraries. |
| **Dark/Light Parity** | Every token has `:root` (light) and `.dark` (dark) values. Components never branch on theme. |
| **Progressive Enhancement** | Hover effects and animations enhance but never gate functionality. |
| **Portability** | Works with React, WordPress, static HTML, or any CSS-capable framework. |
| **Component Independence** | Components read from the global token layer, never from sibling/parent styles. |
| **Consistent State Language** | All interactive elements follow the same state pattern (default → hover → active → focus → disabled). |

### Forbidden Patterns

| Pattern | Why Forbidden | Correct Alternative |
|---------|--------------|---------------------|
| `color: #7c3aed` | Hardcoded hex — bypasses theme system | `hsl(var(--primary))` |
| `bg-purple-600` | Tailwind literal color — not theme-aware | `bg-primary` |
| `rgb(124, 58, 237)` | RGB format — can't use opacity modifier | `hsl(var(--primary))` |
| `hsl(252, 85%, 60%)` | Full hsl() — should use variable | `hsl(var(--primary))` |
| `text-white` | Literal color — bypasses tokens | `text-primary-foreground` |
| `bg-black` | Literal color — bypasses tokens | `bg-background` (dark mode) |

**Exception:** Code block backgrounds (`hsl(220, 14%, 11%)`) are intentionally fixed to maintain readability across themes.

---

## Complete Token Registry

### Core Surface Colors

| Token | Light Value | Dark Value | Purpose |
|-------|-------------|------------|---------|
| `--background` | `0 0% 100%` | `230 25% 8%` | Page background |
| `--foreground` | `230 25% 15%` | `220 20% 92%` | Primary text color |
| `--card` | `0 0% 100%` | `230 20% 12%` | Card/panel background |
| `--card-foreground` | `230 25% 15%` | `220 20% 92%` | Card text color |
| `--popover` | `0 0% 100%` | `230 20% 12%` | Dropdown/popover background |
| `--popover-foreground` | `230 25% 15%` | `220 20% 92%` | Popover text color |

### Brand Colors

| Token | Light Value | Dark Value | Purpose |
|-------|-------------|------------|---------|
| `--primary` | `252 85% 60%` | `252 85% 65%` | Primary brand — buttons, links, accents |
| `--primary-foreground` | `0 0% 100%` | `0 0% 100%` | Text on primary backgrounds |
| `--accent` | `330 85% 60%` | `330 85% 65%` | Secondary accent — gradients, highlights |
| `--accent-foreground` | `0 0% 100%` | `0 0% 100%` | Text on accent backgrounds |

### Neutral & Muted

| Token | Light Value | Dark Value | Purpose |
|-------|-------------|------------|---------|
| `--secondary` | `210 40% 96%` | `230 20% 16%` | Secondary surface (subtle backgrounds) |
| `--secondary-foreground` | `230 25% 15%` | `220 20% 92%` | Text on secondary surfaces |
| `--muted` | `220 20% 96%` | `230 15% 18%` | Muted backgrounds (headers, wells) |
| `--muted-foreground` | `220 10% 46%` | `220 10% 60%` | De-emphasized text (labels, captions) |

### Feedback / State Colors

| Token | Light Value | Dark Value | Purpose |
|-------|-------------|------------|---------|
| `--destructive` | `0 84% 60%` | `0 72% 51%` | Error, delete, danger |
| `--destructive-foreground` | `0 0% 100%` | `0 0% 100%` | Text on destructive |
| `--warning` | `38 92% 50%` | `48 96% 53%` | Warnings, caution |
| `--warning-foreground` | `48 96% 89%` | `26 83% 14%` | Text on warning |
| `--success` | `152 70% 42%` | `152 65% 48%` | Success, checked, valid |
| `--success-foreground` | `150 85% 96%` | `150 80% 10%` | Text on success |

### Border & Input

| Token | Light Value | Dark Value | Purpose |
|-------|-------------|------------|---------|
| `--border` | `220 20% 90%` | `230 15% 20%` | Default border color |
| `--input` | `220 20% 90%` | `230 15% 20%` | Input field borders |
| `--ring` | `252 85% 60%` | `252 85% 65%` | Focus ring color |
| `--radius` | `0.5rem` | `0.5rem` | Base border radius |

### Sidebar Tokens

| Token | Light Value | Dark Value | Purpose |
|-------|-------------|------------|---------|
| `--sidebar-background` | `240 20% 98%` | `230 22% 10%` | Sidebar panel background |
| `--sidebar-foreground` | `230 15% 30%` | `220 15% 80%` | Sidebar text |
| `--sidebar-primary` | `252 85% 60%` | `252 85% 65%` | Active item highlight |
| `--sidebar-primary-foreground` | `0 0% 100%` | `0 0% 100%` | Text on active item |
| `--sidebar-accent` | `252 85% 96%` | `252 40% 18%` | Hover/selected background |
| `--sidebar-accent-foreground` | `252 85% 40%` | `252 85% 80%` | Hover/selected text |
| `--sidebar-border` | `220 20% 92%` | `230 15% 18%` | Sidebar dividers |
| `--sidebar-ring` | `252 85% 60%` | `252 85% 65%` | Sidebar focus ring |

### Reading / Prose Theme Tokens

| Token | Light Value | Dark Value | Purpose |
|-------|-------------|------------|---------|
| `--heading-gradient-from` | `252 85% 60%` | `252 85% 70%` | Heading gradient start |
| `--heading-gradient-to` | `330 85% 60%` | `330 85% 70%` | Heading gradient end |
| `--link-color` | `252 85% 55%` | `252 85% 72%` | Link text color |
| `--code-bg` | `250 25% 95%` | `230 20% 15%` | Inline code background |
| `--code-text` | `330 85% 45%` | `330 85% 70%` | Inline code text color |
| `--blockquote-border` | `252 85% 70%` | `252 60% 50%` | Blockquote left border |
| `--table-header-bg` | `252 85% 97%` | `230 20% 14%` | Table header background |
| `--table-row-hover` | `252 85% 97%` | `252 40% 15%` | Table row hover background |
| `--highlight-glow` | `252 85% 60%` | `252 85% 65%` | Inline code hover glow |

### Code Block Component Tokens (Fixed, Not Themed)

Code blocks maintain a consistent dark appearance in both light and dark modes. These are **NOT** CSS custom properties — they are fixed values.

| Property | HSL Value | Hex | Purpose |
|----------|-----------|-----|---------|
| Block background | `220, 14%, 11%` | `#181c24` | Always-dark code background |
| Header background | `220, 14%, 14%` | `#1f232b` | Code header bar |
| Header border | `220, 13%, 20%` | `#2c3038` | Header bottom border |
| Block outer border | `220, 13%, 22%` | `#30353e` | Outer border |
| Line number background | `220, 14%, 9%` | `#141820` | Line number gutter |
| Line number border | `220, 13%, 18%` | `#272b33` | Gutter right border |
| Line number text | `220, 10%, 35%` | `#535862` | Line number text |
| Line hover background | `220, 15%, 16%` | `#232830` | Line hover bg |
| Tool button background | `220, 13%, 20%` | `#2c3038` | Header action buttons |
| Tool button border | `220, 13%, 25%` | `#373c45` | Button borders |
| Tool button hover bg | `220, 13%, 28%` | `#3e434c` | Button hover |
| Font controls bg | `220, 13%, 18%` | `#272b33` | Font size control group |
| `--code-font-size` | — | — | Default `18px` |
| `--code-line-height` | — | — | Default `1.6` |

#### Themed Elements Within Code Blocks

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

### Language Accent Colors

Each language gets a unique HSL accent (`--lang-accent`) for badge dot, hover glow, and fullscreen shadow:

| Language | HSL Value |
|----------|-----------|
| TypeScript | `210 80% 60%` |
| JavaScript | `50 90% 55%` |
| Go | `190 80% 50%` |
| Rust | `20 85% 55%` |
| CSS | `280 70% 60%` |
| JSON | `45 85% 55%` |
| Bash/Shell | `120 50% 50%` |
| SQL | `200 70% 55%` |
| Markdown | `252 60% 60%` |
| YAML | `340 60% 55%` |
| PHP | `240 55% 60%` |
| HTML/XML | `15 80% 55%` |
| Plain Text | `220 10% 65%` |
| Tree/Structure | `220 10% 65%` |

---

## Typography

### Font Stacks

| Role | Primary Font | Fallbacks | CSS Variable |
|------|-------------|-----------|-------------|
| **Headings** | Ubuntu | sans-serif | `font-heading` |
| **Body Text** | Poppins | sans-serif | `font-body` |
| **Code (Inline)** | JetBrains Mono | Fira Code, ui-monospace, monospace | — |
| **Code (Blocks)** | Ubuntu Mono | JetBrains Mono, ui-monospace, monospace | — |
| **Keyboard** | Ubuntu | monospace | — |

Loading via Google Fonts in `<head>`:

```html
<link href="https://fonts.googleapis.com/css2?family=Ubuntu:wght@400;500;700&family=Poppins:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;500;700&family=Ubuntu+Mono:wght@400;700&display=swap" rel="stylesheet">
```

### Size Hierarchy — Prose Content

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

### Fullscreen Mode Sizes

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
| Code content | `var(--code-font-size)` (default `18px`) | 400 |
| Line numbers | `calc(var(--code-font-size) * 0.7)` | 400 |

### Weight Usage

| Weight | Name | Usage |
|--------|------|-------|
| 300 | Light | Rarely — only large decorative text |
| 400 | Regular | Body text, paragraphs, table cells, code |
| 500 | Medium | Links, inline code, subtle emphasis |
| 600 | SemiBold | H3, H4, labels, table headers, badges |
| 700 | Bold | H1, H2, strong text, font controls |

### Text Spacing

| Property | Value | Applied To |
|----------|-------|------------|
| Letter-spacing `-0.02em` | Tighter | H1 headings only |
| Letter-spacing `0.03em` | Wider | Table headers (uppercase) |
| Letter-spacing `0.05em` | Widest | Badge labels, checklist titles |
| Line-height `1.65` | Relaxed | Body paragraphs |
| Line-height `1.55` | Standard | List items |
| Line-height `var(--code-line-height)` (`1.6`) | Code | Code blocks |

### Heading Gradient Effect

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

Hover: `filter: brightness(1.2) saturate(1.1)`.

---

## Spacing & Layout

### Spacing Scale

| Token | Value | Usage |
|-------|-------|-------|
| `0.1rem` | 1.6px | Micro padding (paragraph vertical) |
| `0.25rem` | 4px | Code line padding, minimal gaps |
| `0.35rem` | 5.6px | Table cell padding, checklist items |
| `0.4rem` | 6.4px | Badge gaps, header/checklist padding |
| `0.5rem` | 8px | Code block header padding, base `--radius` |
| `0.75rem` | 12px | Blockquote/checklist margins, code block radius |
| `1rem` | 16px | Standard section spacing, code body padding |
| `1.25rem` | 20px | H2 bottom margin, code content right padding |
| `1.5rem` | 24px | List left padding |
| `1.8rem` | 28.8px | H2 top margin |
| `2rem` | 32px | Container padding, fullscreen inset |
| `5rem` | 80px | Large section margins |

### Container

```typescript
container: {
  center: true,
  padding: "2rem",
  screens: { "2xl": "1400px" }
}
```

Max content width: **1400px**, centered with **2rem** horizontal padding.

### Layout Patterns

#### Full-Width Sidebar Layout

```
┌──────────┬─────────────────────────────────┐
│ Sidebar  │         Main Content            │
│ (fixed)  │   ┌─────────────┬─────────┐     │
│          │   │   Prose     │   TOC   │     │
│          │   │   Content   │ (sticky)│     │
│          │   └─────────────┴─────────┘     │
└──────────┴─────────────────────────────────┘
```

- Sidebar: collapsible via `SidebarProvider`, toggle with `Ctrl+B`
- Main content: fills remaining width
- TOC: sticky right panel, hidden in split mode

#### Split View Layout

```
┌──────────┬────────────────┬──┬──────────────┐
│ Sidebar  │    Editor      │÷ │   Preview    │
│          │  (Monaco)      │  │  (Markdown)  │
└──────────┴────────────────┴──┴──────────────┘
```

- Divider: `6px` wide, draggable, split ratio clamped: 20%-80%
- Divider hover: primary color accent, handle grows from 32px to 48px

---

## Border Radius Scale

| Token | Value | Usage |
|-------|-------|-------|
| `--radius` | `0.5rem` (8px) | Base radius — cards, inputs |
| `lg` | `var(--radius)` | Large radius |
| `md` | `calc(var(--radius) - 2px)` | Medium radius (6px) |
| `sm` | `calc(var(--radius) - 4px)` | Small radius (4px) |
| `0.75rem` | 12px | Code block wrapper |
| `1rem` | 16px | Fullscreen code block |
| `0.375rem` | 6px | Tool buttons, font controls |
| `0.35rem` | 5.6px | Checklist copy button, checkboxes |
| `0.65rem` | 10.4px | Checklist block |
| `50%` | Circle | Language badge dot |
| `3px` | Subtle | Paragraph hover, list items |
| `2px` | Minimal | Code line hover |
| `1px` | Hairline | Link underline |

### Border Thickness

| Thickness | Usage |
|-----------|-------|
| `1px` | Default borders — cards, tables, inputs, code blocks |
| `1.5px` | Unchecked checkbox border |
| `2px` | Table header bottom, pinned line gutter, link underline |
| `3px` | H3 left accent border |
| `4px` | Blockquote left border (gradient) |
| `6px` | Split divider width |

### Special Border Patterns

**Gradient Blockquote Border:**

```css
border-left: 4px solid transparent;
border-image: linear-gradient(
  to bottom,
  hsl(var(--heading-gradient-from)),
  hsl(var(--heading-gradient-to))
) 1;
```

**H3 Left Accent:** `border-left: 3px solid hsl(var(--primary) / 0.5)` → hover: full opacity.

**Horizontal Rule:** Gradient background, not a border:

```css
height: 1px;
background: linear-gradient(
  90deg, transparent,
  hsl(var(--heading-gradient-from) / 0.4),
  hsl(var(--heading-gradient-to) / 0.4),
  transparent
);
```

---

## Motion & Transitions

### Core Transition Tokens

| Duration | Easing | Usage |
|----------|--------|-------|
| `0.15s` | `ease` | Transform shifts (translateX, translateY, scale) |
| `0.2s` | `ease` | Color changes, background changes, border changes, opacity |
| `0.3s` | `ease` | Box-shadow, filter, complex multi-property transitions |
| `0.3s` | `cubic-bezier(0.4, 0, 0.2, 1)` | Link underline sweep |

**Rule:** Never exceed `0.3s` for hover transitions. Users perceive > 300ms as sluggish.

### Transition Patterns

| Pattern | CSS | Used On |
|---------|-----|---------|
| **Color** | `transition: color 0.2s ease, background 0.2s ease` | Paragraphs, list items, H4, tool buttons |
| **Lift + Glow** | `transition: box-shadow 0.3s ease, transform 0.2s ease` | Code block wrapper (`translateY(-2px)`), inline code (`translateY(-1px)`) |
| **Slide / Nudge** | `transition: transform 0.15s ease` | List items, checkboxes, blockquotes (`translateX(3px)`) |
| **Scale** | `transition: transform 0.2s ease` | Checkbox box (`scale(1.1)`), list bullet (`scale(1.3)`) |
| **Link Underline Sweep** | `transform: scaleX(0)` → hover: `scaleX(1)`, origin swap | `.spec-link::after` pseudo-element |
| **H3 Border Slide** | `transition: padding-left 0.2s ease, border-color 0.2s ease` | H3 headings |
| **Brightness Boost** | `filter: brightness(1.2) saturate(1.1)` | H1, H2 gradient headings |

### Keyframe Animations

**slideUpBar** (copy selection bar):

```css
@keyframes slideUpBar {
  from { opacity: 0; transform: translateY(4px); }
  to   { opacity: 1; transform: translateY(0); }
}
```

**Accordion** (sidebar expand/collapse):

```css
@keyframes accordion-down {
  from { height: 0; }
  to   { height: var(--radix-accordion-content-height); }
}
```

### Button Slide Text Animation (CSS3)

CTA buttons use a CSS3 text slide — default text slides up and out, hover text slides up from below:

```css
.slide-btn { position: relative; overflow: hidden; }
.slide-btn .btn-text-default { transform: translateY(0); }
.slide-btn .btn-text-hover { position: absolute; transform: translateY(100%); }
.slide-btn:hover .btn-text-default { transform: translateY(-100%); }
.slide-btn:hover .btn-text-hover { transform: translateY(0); }
```

Transition: `0.3s cubic-bezier(0.4, 0, 0.2, 1)`.

### Reduced Motion Support

```css
@media (prefers-reduced-motion: reduce) {
  * {
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
  }
}
```

### Theme Transition

```css
body { transition: background-color 0.3s ease, color 0.3s ease; }
```

### Motion Anti-Patterns

| ❌ Forbidden | ✅ Alternative |
|-------------|---------------|
| Bounce easing | `ease` or `cubic-bezier(0.4, 0, 0.2, 1)` |
| Duration > 0.3s for hover | 0.2s–0.3s max |
| `translateY` > `-4px` | `-1px` to `-2px` |
| `scale` > `1.15` | `1.03` to `1.1` |
| Glow opacity > `0.3` | `0.1` to `0.25` |
| JS animation libraries | CSS3 transitions/keyframes |
| Animating `width`/`height` | `transform: scale()` or `max-height` |
| Animating `top`/`left` | `transform: translate()` |

---

## Component Specifications

### Header/Navigation

- Fixed top header: `sticky`, `z-index: 50`, blur backdrop
- Layout: `flex`, `space-between` — Logo | Nav Items | Action Icons
- Nav item hover: underline sweep (CSS `::after`, `scaleX`, gradient)
- Nav item active: `color: hsl(var(--primary))`, `bg: hsl(var(--primary) / 0.08)`
- Icon hover: `scale(1.05)`, active: `scale(0.95)` — no bounce, no spring
- Mobile: hamburger menu, sidebar becomes sheet overlay

### Button System

| Variant | Background | Text | Hover |
|---------|-----------|------|-------|
| `primary` | `hsl(var(--primary))` | `hsl(var(--primary-foreground))` | `translateY(-1px)` + glow |
| `ghost` | transparent | `hsl(var(--foreground))` | `hsl(var(--muted) / 0.5)` bg |
| `highlight` | `linear-gradient(135deg, --primary, --accent)` | `hsl(var(--primary-foreground))` | `translateY(-2px)` + glow |
| `destructive` | `hsl(var(--destructive))` | `hsl(var(--destructive-foreground))` | Darker bg |

All buttons: `focus-visible: 2px solid hsl(var(--ring))`, disabled: `opacity: 0.5`.

### Sidebar

- Width: `16rem` expanded, icon-only collapsed
- Toggle: `Ctrl+B` / `Cmd+B`
- Tree items: `0.8rem`, `0.25rem 0.5rem` padding
- Default: `hsl(var(--sidebar-foreground))`, hover: `hsl(var(--sidebar-accent))`, active: `hsl(var(--sidebar-primary))`
- Folder chevrons rotate `0° → 90°` via `transform 0.2s ease`
- Search: focus ring via `--sidebar-ring`
- Mobile (< 768px): `Sheet` overlay, slides from left

### Section Patterns

| Pattern | Key Tokens |
|---------|-----------|
| **Hero** | `--background` bg, gradient heading, `.btn-primary.slide-btn` CTA |
| **Feature Cards Grid** | `--card` bg, `--border` border, hover: `scale(1.03)` + shadow |
| **Team Section** | `--muted / 0.3` bg, gradient heading, member cards |
| **CTA Banner** | `linear-gradient(135deg, --primary, --accent)` bg, `--primary-foreground` text |
| **Data Table** | `--table-header-bg`, `--table-row-hover`, left accent shadow |
| **Checklist** | `--card` bg, `--border` border, success gradient checkbox |

---

## Code Block System

### Structure

```
┌─────────────────────────────────────────────────┐
│ ● LANG_BADGE    N lines  A- A A+  Copy  ⇲  ⛶  │  ← Header
├────┬────────────────────────────────────────────┤
│  1 │ const x = 42;                              │  ← Body
│  2 │ console.log(x);                             │
├────┴────────────────────────────────────────────┤
│ Lines 1-2 selected          Copy Selected  ✕    │  ← Selection Bar
└─────────────────────────────────────────────────┘
```

### Interactions

- **Line click:** Toggle pinned state → primary-tinted background
- **Shift-click:** Extend selection range
- **Click-drag:** Multi-select across line numbers
- **Font controls:** `A-` (−2px, min 12px), `A` (reset 18px), `A+` (+2px, max 32px)
- **Copy:** Success state: `hsl(152, 60%, 18%)` bg, `hsl(152, 70%, 60%)` text, reverts after 2s
- **Fullscreen:** `position: fixed; inset: 2rem`, overlay `blur(4px)`, Escape to exit

### Tree / Structure Rendering

| Element | Rendering |
|---------|-----------|
| Directories (`name/`) | 📁 prefix, bold, `--foreground` |
| Files (`name.ext`) | 📄 prefix, `--foreground / 0.85` |
| Tree guides (`├└│─`) | `--muted-foreground / 0.5` |
| Ellipsis (`...`) | `--accent` color |
| Comments (`# text`) | `--muted-foreground`, italic |

---

## Re-Theming

Change HSL values in `:root {}` and `.dark {}` — every component updates automatically. No component files need editing.

```css
/* Switch from purple/pink to teal/amber */
:root {
  --primary: 175 85% 40%;            /* was: 252 85% 60% */
  --accent: 38 92% 50%;              /* was: 330 85% 60% */
  --heading-gradient-from: 175 85% 45%;
  --heading-gradient-to: 38 92% 55%;
}
```

### How to Create a New Theme Variant

| Adjustment | Technique |
|-----------|-----------|
| Change brand identity | Alter hue on `--primary` and `--accent` |
| Make warmer | Shift hues toward 30-60° range |
| Make cooler | Shift hues toward 200-240° range |
| Increase contrast | Increase saturation, widen lightness gaps |
| Soften design | Decrease saturation by 10-20% |
| Dark mode | Backgrounds: 5-12% lightness; Foregrounds: 85-95% lightness |

**Validation rules:**
- Heading gradients must remain visually distinct
- Primary and accent should have minimum 60° hue difference
- Text on colored backgrounds must meet WCAG AA contrast (4.5:1)

---

## Responsive Breakpoints

| Breakpoint | Width | Behavior |
|-----------|-------|----------|
| Mobile | `< 768px` | Sidebar becomes sheet overlay, single column |
| Tablet | `768px – 1024px` | Sidebar can be collapsed |
| Desktop | `> 1024px` | Full sidebar + content + TOC |
| Wide | `> 1400px` | Container max-width caps |

---

## Rules for New Pages

1. **No page-specific colors** — use only semantic tokens. If a new color is needed, add it to the token registry first.
2. **Consistent spacing rhythm** — use the spacing scale. Section-to-section: `4rem–6rem`. Heading-to-content: `1rem–1.5rem`.
3. **Heading hierarchy** — H1 → H2 → H3 → H4, never skip levels. One H1 per page.
4. **Responsive by default** — mobile-first with breakpoint overrides. Touch targets min 44px.
5. **Dark mode automatic** — tokens handle it, no conditional logic.
6. **CSS3 motion only** — no JavaScript animation libraries. Max 3 unique durations per page.
7. **Compose from section patterns** — `Page = Header + [Hero] + N × [Section Pattern] + [CTA] + [Footer]`.
8. **Font consistency** — Headings: Ubuntu, Body: Poppins, Code: JetBrains Mono / Ubuntu Mono.

### New Page Checklist

- [ ] All colors from theme tokens
- [ ] Heading hierarchy H1 → H2 → H3 (no skips)
- [ ] Sections use patterns from section patterns spec
- [ ] Transitions use shared durations and easings
- [ ] Works in both light and dark mode
- [ ] Responsive at mobile/tablet/desktop breakpoints
- [ ] Hover effects follow state language
- [ ] Fonts match design system stacks
- [ ] Spacing uses the defined scale

---

## WordPress Migration Notes

- CSS custom properties work natively in WordPress block themes
- Map Tailwind classes to WordPress block styles
- Admin theming via custom admin CSS loading the same tokens
- No component logic changes needed — only template mapping
- `theme.json` maps tokens to WordPress palette, spacing, and typography settings
- Dark mode via body class toggle (JavaScript or plugin)

---

*Consolidated design system — v3.2.0 — 2026-04-16*
