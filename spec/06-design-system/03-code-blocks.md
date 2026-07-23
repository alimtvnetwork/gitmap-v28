# 03 — Code Blocks

> All colors referenced here are tokens from `01-colors-and-themes.md`.
> Syntax token colors are detailed in `05-syntax-token-colors.md`.

---

## Overview

Code blocks are rendered by `src/components/docs/CodeBlock.tsx`.
They provide syntax highlighting, line numbers, interactive features,
and a terminal-inspired look.

---

## Visual Structure

```
┌─────────────────────────────────────────────┐
│ ● ● ●   [language badge]     [actions bar]  │  ← Header
├─────────────────────────────────────────────┤
│  1 │ import fmt                              │  ← Code area
│  2 │ func main() {                           │
│  3 │     fmt.Println("hello")                │
│  4 │ }                                       │
└─────────────────────────────────────────────┘
```

### Header

- **Background:** `--terminal` token (`220 25% 8%` light / `220 25% 5%` dark).
- **Traffic light dots:** Three circles (red, yellow, green) — decorative.
- **Language badge:** Colored pill using `LANG_COLORS` map (per-language accent).
- **Actions:** Copy, Download, Font Size (S/M/L), Fullscreen toggle.

### Code Area

- **Background:** `--terminal` token.
- **Text color:** `hsl(220, 20%, 92%)` (base `.hljs` color).
- **Font:** `font-mono` (Ubuntu Mono).
- **Font sizes:** Three tiers — S (`0.75rem`), M (`0.85rem`), L (`0.95rem`).
- **Line numbers:** `--muted-foreground` at reduced opacity.

---

## Language Accent Colors

Each language has a CSS variable `--lang-accent` used for badge color and hover glow:

| Language | Accent HSL |
|----------|-----------|
| go | `171 68% 45%` (cyan/teal) |
| typescript / ts | `212 92% 45%` (blue) |
| javascript / js | `50 90% 50%` (yellow) |
| bash / shell | `142 71% 45%` (green — `--primary`) |
| json | `35 80% 50%` (orange) |
| sql | `280 60% 55%` (purple) |
| css | `330 70% 55%` (pink) |
| powershell | `212 70% 55%` (blue) |
| markdown | `220 15% 60%` (gray-blue) |

---

## Interactive Features

### Line Hover

```css
.code-line:hover {
  background: hsl(220, 15%, 16%);
  border-left-color: hsl(var(--primary) / 0.6);
}
```

### Line Pinning (Click-to-Pin)

- Click a line number to pin/unpin.
- Shift+click to pin a range from last pinned line.
- Pinned lines get:

```css
.code-line-pinned {
  background: hsl(var(--primary) / 0.22) !important;
  border-left-color: hsl(var(--primary)) !important;
}
```

- **Rule:** Pinned background opacity must be ≥ 0.2 (see `ui-issues.md` Issue 3).

### Copy Button

- Copies raw code (no line numbers) to clipboard.
- Shows checkmark for 2 seconds on success.

### Download Button

- Downloads as file with correct extension from `LANG_EXTENSIONS` map.
- Filename: `code.{ext}` or uses `title` prop if provided.

### Fullscreen

- Toggles a fullscreen overlay with the code block.
- Uses `fixed inset-0 z-50` positioning.

---

## Code Block Hover Glow

```css
.code-block-hover:hover {
  box-shadow: 0 8px 32px hsl(var(--lang-accent, 220 10% 50%) / 0.1),
              0 0 0 1px hsl(var(--lang-accent, 220 10% 50%) / 0.15);
}
```

Uses the `--lang-accent` CSS variable set per block, falling back to neutral.

---

## Text Selection in Code Blocks

```css
pre ::selection,
code ::selection,
.hljs ::selection {
  background: hsl(142 65% 40% / 0.35);
  color: hsl(0 0% 100%);
}
```

- Always white text on green background for readability.
- See `ui-issues.md` Issue 1 for rationale.

---

## Implementation Reference

```tsx
<CodeBlock
  code={`func main() {\n  fmt.Println("hello")\n}`}
  language="go"
  title="main.go"
/>
```

### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `code` | `string` | required | Raw code string |
| `language` | `string` | `"bash"` | Highlight.js language identifier |
| `title` | `string` | — | Optional title shown in header |

---

## Constraints

- Code block background always uses `--terminal` token.
- Never use `translateY` animations on code blocks.
- Line pin opacity ≥ 0.2 on dark backgrounds.
- Text selection must be white on colored background in dark mode.
- All accent colors are per-language, not hardcoded per-instance.
