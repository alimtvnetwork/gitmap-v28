# Borders & Shapes

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

Borders and corner radii create visual containment and hierarchy. All border colors reference theme tokens. Radius uses a base variable with computed variants.

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

---

## Border Thickness

| Thickness | Usage |
|-----------|-------|
| `1px` | Default borders — cards, tables, inputs, code blocks |
| `1.5px` | Unchecked checkbox border |
| `2px` | Table header bottom, pinned line gutter, link underline |
| `3px` | H3 left accent border |
| `4px` | Blockquote left border (gradient) |
| `6px` | Split divider width |

---

## Border Color Behavior by State

### Default State

```css
border: 1px solid hsl(var(--border));
```

### Hover State

| Component | Border Change |
|-----------|--------------|
| Code block wrapper | `box-shadow` glow with `--lang-accent`, no border change |
| Checklist block | `border-color: hsl(var(--primary) / 0.3)` |
| Tool buttons | `border-color: hsl(var(--lang-accent) / 0.4)` |
| Checklist copy button | `border-color: hsl(var(--primary) / 0.3)` |

### Active/Selected State

| Component | Border Change |
|-----------|--------------|
| Pinned code line | `border-right: 2px solid hsl(var(--primary) / 0.6)` on line number |
| Table row hover | `box-shadow: inset 3px 0 0 hsl(var(--primary) / 0.5)` |

### Fullscreen State

```css
box-shadow: 0 25px 80px hsl(var(--lang-accent) / 0.25),
            0 0 0 1px hsl(var(--lang-accent) / 0.3);
```

---

## Special Border Patterns

### Gradient Border (Blockquote)

```css
border-left: 4px solid transparent;
border-image: linear-gradient(
  to bottom,
  hsl(var(--heading-gradient-from)),
  hsl(var(--heading-gradient-to))
) 1;
```

### H2 Bottom Border

```css
border-bottom: 1px solid hsl(var(--border));
```

### H3 Left Accent

```css
border-left: 3px solid hsl(var(--primary) / 0.5);
/* Hover: */
border-color: hsl(var(--primary));
```

### Horizontal Rule

Not a border — uses a gradient background:

```css
height: 1px;
background: linear-gradient(
  90deg,
  transparent,
  hsl(var(--heading-gradient-from) / 0.4),
  hsl(var(--heading-gradient-to) / 0.4),
  transparent
);
```

---

## Shape Language Summary

| Visual Feel | How Achieved |
|------------|-------------|
| **Rounded, friendly** | `0.5rem`–`0.75rem` radius on containers |
| **Subtle containment** | 1px borders at low opacity |
| **Depth on hover** | `box-shadow` glow, not heavier borders |
| **Accent hierarchy** | Left borders for H3, blockquotes, table row highlight |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Theme Variables | [02-theme-variable-architecture.md](./02-theme-variable-architecture.md) |
| Code Block Borders | [07-code-blocks.md](./07-code-blocks.md) |
| Motion System | [06-motion-transitions.md](./06-motion-transitions.md) |
