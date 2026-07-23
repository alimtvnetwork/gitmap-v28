# 02 — Typography

> All colors referenced here are tokens from `01-colors-and-themes.md`.
> Never use raw color values.

---

## Font Families

Loaded via Google Fonts in `src/index.css`:

| Role | Font | Tailwind Class | CSS Fallback |
|------|------|----------------|--------------|
| Body text | Poppins | `font-sans` | `system-ui, sans-serif` |
| Headings | Ubuntu | `font-heading` | `system-ui, sans-serif` |
| Code / monospace | Ubuntu Mono | `font-mono` | `monospace` |

**Import URL:**
```
https://fonts.googleapis.com/css2?family=Ubuntu+Mono:wght@400;700&family=Poppins:wght@400;500;600;700&family=Ubuntu:wght@400;500;700&display=swap
```

---

## Heading Styles

All headings (`h1`–`h6`) use `font-family: 'Ubuntu'` and `font-weight: 700`
via the base layer in `index.css`.

### Gradient Headings (`.docs-h1`, `.docs-h2`)

```css
/* Light mode */
background: linear-gradient(135deg, hsl(142 65% 32%), hsl(142 55% 42%));
-webkit-background-clip: text;
-webkit-text-fill-color: transparent;

/* Dark mode (.dark) */
background: linear-gradient(135deg, hsl(var(--primary)), hsl(142 71% 60%));
```

- Uses **darker green on light backgrounds** for legibility.
- Uses **brighter green on dark backgrounds** for vibrancy.
- Hover: `filter: brightness(1.15) saturate(1.1)`.

### Sub-headings (`.docs-h3`)

```css
padding-left: 0.65rem;
border-left: 3px solid hsl(142 55% 35% / 0.5);
color: hsl(142 55% 30%); /* light mode */
```

- Dark mode inherits `color` and uses `hsl(var(--primary) / 0.5)` border.
- Hover: shifts `padding-left` to `0.85rem` and solidifies border.

---

## Body Text

Applied via base layer:

```css
p, li, td, th, label {
  font-family: 'Poppins', system-ui, sans-serif;
}
```

Color: `text-foreground` (token from `01`).

---

## Inline Code (`.docs-inline-code`)

```css
@apply bg-[hsl(var(--code-bg))] text-primary font-mono text-[0.85em]
       font-medium px-[0.45em] py-[0.2em] rounded-[5px]
       border border-border/50;
```

- Background: `--code-bg` token.
- Text: `--primary` token.
- Hover: `box-shadow: 0 0 0 2px hsl(var(--primary) / 0.15)` + `translateY(-1px)`.

---

## Constraints

- Headings must always use separate light/dark color values (see `ui-issues.md` Issue 2).
- Never apply `translateY` or `scale` animations to text elements.
- Body font weights: 400 (normal), 500 (medium), 600 (semibold), 700 (bold).
