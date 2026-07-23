# 04 — Component Patterns

> All colors referenced here are tokens from `01-colors-and-themes.md`.

---

## Selection and Focus

### Text Selection

```css
/* Global */
::selection {
  background: hsl(var(--primary) / 0.25);
  color: hsl(var(--foreground));
}

/* Dark mode */
.dark ::selection {
  background: hsl(142 71% 55% / 0.3);
  color: hsl(220 10% 95%);
}
```

### Focus Ring

All interactive elements use `--ring` (`142 71% 45%`) for focus indicators.

---

## Tables (`.docs-table`)

```css
thead    → bg-muted/50
th       → font-heading, uppercase, text-xs, tracking-wide, text-muted-foreground
tbody tr → hover: bg-muted/30, box-shadow: inset 3px 0 0 hsl(var(--primary) / 0.5)
```

---

## Blockquotes (`.docs-blockquote`)

```css
border-left: 4px solid hsl(var(--primary) / 0.5);
@apply bg-muted/30 py-2 px-4 rounded-r-lg italic text-muted-foreground;
```

Hover: `bg-muted/50` + `translateX(3px)`.

---

## Horizontal Rules (`.docs-hr`)

```css
height: 1px;
background: linear-gradient(90deg,
  transparent,
  hsl(var(--primary) / 0.4),
  hsl(142 71% 60% / 0.4),
  transparent
);
```

---

## Cards

Use `bg-card text-card-foreground` tokens. Never `bg-white` or `bg-gray-900`.

```tsx
<div className="bg-card text-card-foreground border border-border rounded-lg p-6">
  {children}
</div>
```

---

## Buttons

Use shadcn button variants. Primary buttons use `bg-primary text-primary-foreground`.

---

## Sidebar

Uses dedicated sidebar tokens (`--sidebar-background`, `--sidebar-foreground`, etc.)
to allow independent theming from the main content area.

---

## Transitions

### Allowed

- `opacity` fades
- `color` / `background-color` transitions
- `border-color` changes
- Horizontal `translateX` (small, ≤ 5px)
- `box-shadow` glow effects

### Prohibited

- `translateY` on text-containing elements
- `scale` / `zoom` on content containers
- Layout-shifting animations on page load

**Global transition:** All elements have `transition-colors duration-300` via base layer.

---

## Spacing System

Uses Tailwind's default spacing scale. Common patterns:

| Context | Spacing |
|---------|---------|
| Page padding | `p-6` or `px-6 py-8` |
| Card padding | `p-6` |
| Section gaps | `space-y-6` or `gap-6` |
| Inline spacing | `gap-2` or `space-x-2` |

---

## Border Radius

From `tailwind.config.ts`:

```ts
borderRadius: {
  lg: "var(--radius)",        // 0.5rem
  md: "calc(var(--radius) - 2px)",
  sm: "calc(var(--radius) - 4px)",
}
```

---

## Constraints

- No raw color classes — always semantic tokens.
- No vertical translation on text elements.
- All interactive states need ≥ 0.2 opacity difference from rest state in dark mode.
- Sidebar tokens are independent from main content tokens.
