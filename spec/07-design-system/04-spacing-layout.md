# Spacing & Layout

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

Consistent spacing creates visual rhythm. All spacing uses a base-4 scale. Layout follows container-constrained, flex/grid patterns with responsive breakpoints.

---

## Spacing Scale

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
| `5rem` | 80px | Large section margins (`.my-5`) |

---

## Container

```typescript
container: {
  center: true,
  padding: "2rem",
  screens: { "2xl": "1400px" }
}
```

Max content width: **1400px**, centered with **2rem** horizontal padding.

---

## Layout Patterns

### Full-Width Sidebar Layout

```
┌──────────┬─────────────────────────────────┐
│ Sidebar  │         Main Content            │
│ (fixed)  │   ┌─────────────┬─────────┐     │
│          │   │   Prose     │   TOC   │     │
│          │   │   Content   │ (sticky)│     │
│          │   └─────────────┴─────────┘     │
└──────────┴─────────────────────────────────┘
```

- Sidebar: collapsible via `SidebarProvider`, toggle with Ctrl+B
- Main content: fills remaining width
- TOC: sticky right panel, hidden in split mode

### Split View Layout

```
┌──────────┬────────────────┬──┬──────────────┐
│ Sidebar  │    Editor      │÷ │   Preview    │
│          │  (Monaco)      │  │  (Markdown)  │
└──────────┴────────────────┴──┴──────────────┘
```

- Divider: `6px` wide, draggable
- Split ratio clamped: 20%-80%
- Divider hover: primary color accent, handle grows from 32px to 48px

---

## Responsive Breakpoints

| Breakpoint | Width | Behavior |
|-----------|-------|----------|
| Mobile | `< 768px` | Sidebar becomes sheet overlay |
| Tablet | `768px - 1024px` | Sidebar can be collapsed |
| Desktop | `> 1024px` | Full sidebar + content + TOC |
| Wide | `> 1400px` | Container max-width caps |

---

## Component Spacing Patterns

### Code Block

| Area | Spacing |
|------|---------|
| Wrapper margin | `my-5` (1.25rem vertical) |
| Header padding | `0.5rem 1rem` |
| Body line numbers padding | `1rem 0` (vertical), `0.75rem` (right) |
| Body content padding | `1rem 1.25rem 1rem 0.5rem` |

### Tables

| Area | Spacing |
|------|---------|
| Wrapper margin | `0.5rem 0` |
| Header cells | `0.45rem 0.75rem` |
| Body cells | `0.35rem 0.75rem` |

### Lists

| Area | Spacing |
|------|---------|
| UL left padding | `1.5rem` |
| LI vertical padding | `0.1rem 0` |
| LI left padding | `0.4rem` |
| Bullet offset | `-0.75rem` from left |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Border Radius | [05-borders-shapes.md](./05-borders-shapes.md) |
| Code Block Layout | [07-code-blocks.md](./07-code-blocks.md) |
