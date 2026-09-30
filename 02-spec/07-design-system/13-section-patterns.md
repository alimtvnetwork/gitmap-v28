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

## Pattern 7: Asymmetric 5-Question Hero with Visual Anchor

An editorial 7+5 or 8+4 golden-ratio split answering all 5 above-the-fold questions with a tangible visual anchor (no blank gradients):

```text
┌──────────────────────────────────────────────────────────┐
│  [Proof Eyebrow: Badge / Verified Metrics]               │
│                                                          │
│  Display Headline (H1, <= 9 words)      ┌──────────────┐ │
│  ─── second clause at 50% opacity ───   │ Visual       │ │
│                                         │ Product      │ │
│  Lead paragraph (max-w-[54ch])          │ Artefact     │ │
│                                         │ (Scorecard,  │ │
│  [Primary CTA Pill]  [Secondary Ghost]  │ Diagram,     │ │
│  Micro reassurance line                 │ or Photo)    │ │
│                                         └──────────────┘ │
└──────────────────────────────────────────────────────────┘
```

| Property | Token / Rule |
|:---|:---|
| Layout Grid | 12-column grid: 7 cols (content) + 5 cols (visual) or 8 + 4 |
| Background | Plane 0 Canvas Base + subtle radial glow (6–10% opacity) |
| Visual Anchor | Occupies ≥ 25% of desktop frame; uses progressive blur or card chrome |
| Primary CTA | Exactly ONE solid accent-filled button (Von Restorff rule) |
| Secondary CTA | Transparent background with hairline border (`rgba(255,255,255,0.16)`) |
| Padding | Desktop `140px–160px` vertical; mobile `72px` |

---

## Pattern 8: Controlled Multi-Card Sliding Carousel Showcase

Interactive portfolio, case studies, or talent showcase with a controlled sliding track:

```text
┌──────────────────────────────────────────────────────────┐
│  Section Heading                     [◄ Prev] [Next ►]   │
│  Lead narrative                      Index: 01 / 06      │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │ Active Card  │  │ Card N+1     │  │ Card N+2     │    │
│  │ (Border glow)│  │ (Hairline)   │  │ (Hairline)   │    │
│  └──────────────┘  └──────────────┘  └──────────────┘    │
│  Filter Pills: [ All ][ Track A ][ Track B ][ Track C ]   │
└──────────────────────────────────────────────────────────┘
```

| Property | Token / Rule |
|:---|:---|
| Desktop Layout | 3 cards visible; advance by 1 card step per click/swipe |
| Mobile Layout | 1 active card + 24px glimpse of next card |
| Motion | `transform: translate3d(-Xpx, 0, 0)` with `--ease-out-editorial` |
| Hover Rule | Zero `translateY` hover lifts (use border glow and image scale `1.02`) |
| Autoplay | 5.2s interval with auto-pause on hover, focus, touch, and offscreen |
| Reduced Motion | Static horizontal scroll list without cloned elements |

---

## Pattern 9: Pinned Progression & Signature Funnel / Timeline

Full-bleed section where a central narrative or filtering process advances through stages as the user scrolls:

- Left column (4 cols): Sticky glass stage rail showing active milestone and live counters.
- Right column (8 cols): Evolving canvas, state-machine diagram, or stepped interactive cards.
- Mobile / Reduced-Motion Fallback: Unpinned vertical stacked cards with static diagrams.

---

## Composing New Sections

When creating a new section pattern:

1. **Choose background:** `--background` (Plane 0), `--bg-raised` (Plane 1), or surface cards (Plane 2)
2. **Alternate planes:** Never place three consecutive sections on the exact same plane
3. **Choose border:** `1px solid rgba(255,255,255,0.08)` for dark contained sections, top-highlight inset
4. **Choose heading:** Fluid clamp display scale with negative tracking on large sizes
5. **Enforce single accent:** Only one primary saturated CTA visible per viewport
6. **Maintain spacing rhythm:** Alternate between open (`160px`) and compressed (`48px`) bands

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Spacing Scale | [05-spacing-layout.md](./05-spacing-layout.md) |
| Theme Catalogue & Palettes | [16-theme-catalogue-and-palettes.md](./16-theme-catalogue-and-palettes.md) |
| Dark Mode & Materiality | [18-dark-mode-and-materiality.md](./18-dark-mode-and-materiality.md) |
| Modern Motion & Sliding Interactions | [19-modern-motion-and-sliding-interactions.md](./19-modern-motion-and-sliding-interactions.md) |
| AI Training & Anti-Slop Guide | [20-ai-training-and-checklist-guide.md](./20-ai-training-and-checklist-guide.md) |
| Button Variants | [11-button-system.md](./11-button-system.md) |
| Motion Rules | [08-motion-transitions.md](./08-motion-transitions.md) |
