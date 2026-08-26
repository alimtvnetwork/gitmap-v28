# Theme & Variable Architecture

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

This is the **single source of truth** for every color, spacing, and visual token in the design system. All components reference these variables — never hardcoded values. To re-theme the entire system, change values in this registry only.

Variables are defined as **CSS custom properties** in `src/index.css` within `:root` (light) and `.dark` (dark) blocks. Values use **HSL space-separated format** (e.g., `252 85% 60%`) to allow opacity modifiers.

---

## Variable Format

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

### Sidebar

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

### Code Block Component Tokens

These are **NOT** CSS custom properties — they are fixed values used in the code block component. They are intentionally NOT themed because code blocks maintain a consistent dark appearance in both light and dark modes.

| Property | Value | Purpose |
|----------|-------|---------|
| Block background | `hsl(220, 14%, 11%)` | Always-dark code background |
| Header background | `hsl(220, 14%, 14%)` | Code header bar |
| Header border | `hsl(220, 13%, 20%)` | Header bottom border |
| Block border | `hsl(220, 13%, 22%)` | Outer border |
| Line number background | `hsl(220, 14%, 9%)` | Line number gutter |
| Line number border | `hsl(220, 13%, 18%)` | Gutter right border |
| Line number text | `hsl(220, 10%, 35%)` | Line number text |
| Tool button background | `hsl(220, 13%, 20%)` | Header action buttons |
| Tool button border | `hsl(220, 13%, 25%)` | Button borders |
| Tool button hover bg | `hsl(220, 13%, 28%)` | Button hover |
| Font controls bg | `hsl(220, 13%, 18%)` | Font size control group |
| `--lang-accent` | Per-language HSL | Language-specific glow color |
| `--code-font-size` | `18px` | Default code font size |
| `--code-line-height` | `1.6` | Code line height |

### Language Accent Colors

Each language gets a unique HSL accent used for the badge dot, hover glow, and fullscreen shadow:

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

## Tailwind Config Mapping

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
  // ... all tokens follow same pattern
}
```

This enables usage like:

```html
<button class="bg-primary text-primary-foreground hover:bg-primary/90">
```

---

## How to Create a New Theme Variant

### Step 1: Copy the `:root` and `.dark` blocks

### Step 2: Adjust HSL values

| Adjustment | Technique |
|-----------|-----------|
| Change brand identity | Alter hue on `--primary` and `--accent` |
| Make warmer | Shift hues toward 30-60° range |
| Make cooler | Shift hues toward 200-240° range |
| Increase contrast | Increase saturation, widen lightness gaps |
| Soften design | Decrease saturation by 10-20% |
| Dark mode | Backgrounds: 5-12% lightness; Foregrounds: 85-95% lightness |

### Step 3: Validate

- Heading gradients should remain visually distinct
- Primary and accent should have minimum 60° hue difference
- Text on colored backgrounds must meet WCAG AA contrast (4.5:1)

---

## Forbidden Patterns

| Pattern | Why Forbidden |
|---------|--------------|
| `color: #7c3aed` | Hardcoded hex — bypasses theme system |
| `bg-purple-600` | Tailwind literal color — not theme-aware |
| `rgb(124, 58, 237)` | RGB format — can't use opacity modifier |
| `hsl(252, 85%, 60%)` | Full hsl() — should be `hsl(var(--primary))` |
| `text-white` | Literal color — use `text-primary-foreground` |
| `bg-black` | Literal color — use `bg-background` in dark mode |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| CSS Source File | `src/index.css` |
| Tailwind Config | `tailwind.config.ts` |
| Design Principles | [01-design-principles.md](./01-design-principles.md) |
| Code Block Tokens | [07-code-blocks.md](./07-code-blocks.md) |
