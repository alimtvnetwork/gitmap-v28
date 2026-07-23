# 01 — Colors and Themes (Root Reference)

All color values live as CSS custom properties in `src/index.css`.
Every component, spec, and style rule references these variables —
never raw hex/hsl literals.

---

## Core Tokens

All values are **HSL without the `hsl()` wrapper** so they can be used
with Tailwind's `hsl(var(--token))` pattern.

### Light Mode (`:root`)

| Token | HSL Value | Usage |
|-------|-----------|-------|
| `--background` | `220 20% 97%` | Page background |
| `--foreground` | `220 25% 10%` | Default text |
| `--card` | `0 0% 100%` | Card surfaces |
| `--card-foreground` | `220 25% 10%` | Card text |
| `--popover` | `0 0% 100%` | Popover/dropdown bg |
| `--popover-foreground` | `220 25% 10%` | Popover text |
| `--primary` | `142 71% 45%` | Brand green — buttons, links, accents |
| `--primary-foreground` | `220 25% 5%` | Text on primary bg |
| `--secondary` | `220 14% 92%` | Secondary surfaces |
| `--secondary-foreground` | `220 25% 10%` | Secondary text |
| `--muted` | `220 14% 92%` | Muted backgrounds |
| `--muted-foreground` | `220 10% 46%` | Muted/hint text |
| `--accent` | `142 71% 45%` | Accent (matches primary) |
| `--accent-foreground` | `220 25% 5%` | Text on accent bg |
| `--destructive` | `0 84% 60%` | Error/danger actions |
| `--destructive-foreground` | `0 0% 100%` | Text on destructive bg |
| `--border` | `220 13% 87%` | Borders and dividers |
| `--input` | `220 13% 87%` | Input field borders |
| `--ring` | `142 71% 45%` | Focus ring color |
| `--radius` | `0.5rem` | Default border radius |
| `--terminal` | `220 25% 8%` | Terminal block bg |
| `--terminal-foreground` | `142 71% 55%` | Terminal text (green) |
| `--code-bg` | `220 20% 94%` | Inline code bg |

### Dark Mode (`.dark`)

| Token | HSL Value | Usage |
|-------|-----------|-------|
| `--background` | `220 25% 6%` | Page background |
| `--foreground` | `220 10% 90%` | Default text |
| `--card` | `220 25% 9%` | Card surfaces |
| `--card-foreground` | `220 10% 90%` | Card text |
| `--popover` | `220 25% 9%` | Popover bg |
| `--popover-foreground` | `220 10% 90%` | Popover text |
| `--primary` | `142 71% 45%` | Brand green (same both modes) |
| `--primary-foreground` | `220 25% 5%` | Text on primary |
| `--secondary` | `220 20% 14%` | Secondary surfaces |
| `--secondary-foreground` | `220 10% 90%` | Secondary text |
| `--muted` | `220 20% 14%` | Muted backgrounds |
| `--muted-foreground` | `220 10% 55%` | Muted text |
| `--accent` | `142 71% 45%` | Accent |
| `--accent-foreground` | `220 25% 5%` | Text on accent |
| `--destructive` | `0 62% 30%` | Error (darker in dark mode) |
| `--destructive-foreground` | `0 0% 100%` | Text on destructive |
| `--border` | `220 20% 16%` | Borders |
| `--input` | `220 20% 16%` | Input borders |
| `--ring` | `142 71% 45%` | Focus ring |
| `--terminal` | `220 25% 5%` | Terminal bg (deeper) |
| `--terminal-foreground` | `142 71% 55%` | Terminal green text |
| `--code-bg` | `220 25% 10%` | Inline code bg |

### Sidebar Tokens

| Token | Light | Dark |
|-------|-------|------|
| `--sidebar-background` | `220 20% 95%` | `220 25% 8%` |
| `--sidebar-foreground` | `220 10% 30%` | `220 10% 75%` |
| `--sidebar-primary` | `142 71% 45%` | `142 71% 45%` |
| `--sidebar-primary-foreground` | `220 25% 5%` | `220 25% 5%` |
| `--sidebar-accent` | `220 14% 90%` | `220 20% 12%` |
| `--sidebar-accent-foreground` | `220 25% 10%` | `220 10% 90%` |
| `--sidebar-border` | `220 13% 87%` | `220 20% 14%` |
| `--sidebar-ring` | `142 71% 45%` | `142 71% 45%` |

---

## Tailwind Config Mapping

In `tailwind.config.ts`, every token is mapped to a Tailwind class:

```ts
colors: {
  background: "hsl(var(--background))",
  foreground: "hsl(var(--foreground))",
  primary: {
    DEFAULT: "hsl(var(--primary))",
    foreground: "hsl(var(--primary-foreground))",
  },
  // ... same pattern for all tokens
}
```

**Usage in components:**
```tsx
// ✅ Correct — uses design tokens
<div className="bg-background text-foreground border-border" />
<button className="bg-primary text-primary-foreground" />

// ❌ Wrong — hardcoded colors
<div className="bg-white text-gray-900 border-gray-200" />
<button className="bg-green-500 text-black" />
```

---

## How to Change the Theme

1. Open `src/index.css`.
2. Modify the HSL values under `:root` (light) and `.dark` (dark).
3. All components automatically pick up the new values.
4. No component files need editing.

### Example: Change brand color from green to blue

```css
/* Before */
--primary: 142 71% 45%;

/* After */
--primary: 217 91% 60%;
```

Every button, link, heading gradient, focus ring, and accent
updates automatically.

---

## Constraints

- **Never use raw colors** in components (`text-white`, `bg-black`, `text-green-500`).
- **Always use semantic tokens** (`text-foreground`, `bg-primary`, `text-muted-foreground`).
- **All HSL values** must omit the `hsl()` wrapper in CSS variables.
- **`--primary` is identical** in light and dark mode by design.
- **Selection colors** must ensure contrast (see `spec/05-coding-guidelines/ui-issues.md`).
