# Section Patterns

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

Reusable section patterns provide consistent visual templates for building pages. Each pattern defines layout, spacing, content grouping, and interaction style. Future pages should compose from these patterns.

---

## Pattern 1: Hero Section

Large introductory section with icon, title, description, and CTA.

```
┌────────────────────────────────────────────┐
│                                            │
│              [Icon / Logo]                 │
│                                            │
│          Gradient Heading (H1)             │
│                                            │
│        Muted description paragraph         │
│                                            │
│         [Primary CTA Button]               │
│           (slide text hover)               │
│                                            │
└────────────────────────────────────────────┘
```

| Property | Token |
|----------|-------|
| Background | `hsl(var(--background))` |
| Icon color | `hsl(var(--primary))` |
| Heading | Gradient from `--heading-gradient-from` to `--heading-gradient-to` |
| Description | `hsl(var(--muted-foreground))` |
| CTA Button | `.btn-primary.slide-btn` |
| Padding | `4rem 2rem` vertical, centered |
| Max width | `600px` for text content |

---

## Pattern 2: Feature Cards Grid

Grid of cards highlighting features or capabilities.

```
┌──────────┐  ┌──────────┐  ┌──────────┐
│  [Icon]  │  │  [Icon]  │  │  [Icon]  │
│  Title   │  │  Title   │  │  Title   │
│  Desc    │  │  Desc    │  │  Desc    │
└──────────┘  └──────────┘  └──────────┘
```

| Property | Token |
|----------|-------|
| Card bg | `hsl(var(--card))` |
| Card border | `1px solid hsl(var(--border))` |
| Card radius | `var(--radius)` |
| Hover | `transform: scale(1.03)`, `box-shadow` elevation, icon scale(1.1) |
| Icon | `hsl(var(--primary))`, `1.5rem` size |
| Title | `font-weight: 600`, `hsl(var(--foreground))` |
| Description | `hsl(var(--muted-foreground))`, `0.85rem` |
| Grid | `grid-template-columns: repeat(auto-fit, minmax(280px, 1fr))` |
| Gap | `1.5rem` |
| Transition | `transform 0.2s ease, box-shadow 0.3s ease` |

---

## Pattern 3: Team / Content Section ("Team Greenhouse" Style)

A thematic content section with a heading, descriptive content, and visual grouping.

```
┌────────────────────────────────────────────┐
│  ┌─ Section Heading (gradient) ─────────┐  │
│  │                                       │  │
│  │  Content paragraph with context       │  │
│  │                                       │  │
│  │  ┌──────┐  ┌──────┐  ┌──────┐       │  │
│  │  │Member│  │Member│  │Member│       │  │
│  │  └──────┘  └──────┘  └──────┘       │  │
│  └───────────────────────────────────────┘  │
└────────────────────────────────────────────┘
```

| Property | Token |
|----------|-------|
| Section bg | `hsl(var(--muted) / 0.3)` |
| Section border | `1px solid hsl(var(--border))` |
| Section radius | `0.75rem` |
| Heading | Gradient text, `font-weight: 700` |
| Content | `hsl(var(--foreground) / 0.9)`, `line-height: 1.65` |
| Member cards | Same as Feature Cards but smaller |
| Padding | `2rem` |
| Hover | Subtle `border-color` shift toward primary |

---

## Pattern 4: CTA Banner

Full-width call-to-action strip.

```
┌────────────────────────────────────────────┐
│    Heading text         [Action Button]    │
└────────────────────────────────────────────┘
```

| Property | Token |
|----------|-------|
| Background | `linear-gradient(135deg, hsl(var(--primary)), hsl(var(--accent)))` |
| Text | `hsl(var(--primary-foreground))` |
| Button | Ghost/outline variant on gradient bg |
| Padding | `2rem 3rem` |
| Radius | `var(--radius)` |

---

## Pattern 5: Data Table Section

Content section with a prominent table.

| Property | Token |
|----------|-------|
| Table wrapper | `border: 1px solid hsl(var(--border))`, `radius: 0.5rem` |
| Header bg | `hsl(var(--table-header-bg))` |
| Row hover | `hsl(var(--table-row-hover))` with left accent shadow |
| Alternating rows | `hsl(var(--muted) / 0.15)` |

---

## Pattern 6: Checklist / Audit Section

Interactive checklist with progress tracking.

| Property | Token |
|----------|-------|
| Block bg | `hsl(var(--card))` |
| Block border | `1px solid hsl(var(--border))` |
| Hover | `border-color: hsl(var(--primary) / 0.3)`, subtle shadow |
| Checked items | Success gradient checkbox |
| Header bg | `hsl(var(--muted) / 0.4)` |

---

## Composing New Sections

When creating a new section pattern:

1. **Choose background:** `--background` (default), `--muted / 0.3` (highlighted), or gradient (CTA)
2. **Choose border:** `1px solid hsl(var(--border))` for contained sections, none for full-bleed
3. **Choose heading:** Gradient for primary sections, solid for secondary
4. **Add hover:** Always include a transition on the container or child elements
5. **Maintain spacing rhythm:** Use the spacing scale from [04-spacing-layout.md](./04-spacing-layout.md)

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Spacing Scale | [04-spacing-layout.md](./04-spacing-layout.md) |
| Button Variants | [09-button-system.md](./09-button-system.md) |
| Motion Rules | [06-motion-transitions.md](./06-motion-transitions.md) |
