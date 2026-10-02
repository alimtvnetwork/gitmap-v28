# 38 — Card, Pricing & High-Conversion Showcase Sections

> **/goal** Specify 3-tier modular pricing tables, interactive split pricing blocks, scroll-driven services showcases, tactile ruled-paper workflow stacks, and enterprise hero sections.
> **/learn** Master the exact geometry, tokens, state machines, and motion physics across high-conversion marketing surfaces. Curriculum cards and interactive rows live in `23-building-block-components.md`.

**Version:** 4.2.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Scope & System Overview

This specification codifies the primary editorial cards and high-conversion sections:
1. **3-Tier Subscription Pricing Matrix (`PricingTierCard`)**
2. **Interactive Split Pricing Blocks (`InteractiveSplitPricing`)**
3. **Scroll-Driven Services Showcase (`WhatWeDoSection`)**
4. **Ruled-Paper Process Reassurance Stack (`WorkflowProcessSection`)**
5. **Split Enterprise Hero with Depth-Scaled Cards (`EnterpriseHeroSection`)**

---

## 2. Tokens & Variables

| Token | Hex / Value | Role |
|---|---|---|
| `--color-ground-light` | `#FFFFFF` | Pure white resting background |
| `--color-border-subtle`| `#E2E8F0` | Standard hairline border (1px) |
| `--color-popular-bg` | `#0B1329` | Dark obsidian ground for popular tier |
| `--color-popular-border`| `#8B5CF6` | Electric violet border for highlighted tier |
| `--color-badge-gold` | `#FACC15` | Most popular pill badge background |
| `--color-badge-ink` | `#0F172A` | Contrast dark text on gold badge |
| `--night` | `oklch(0.19 0 0)` | Near-black surface for dark showcase shells |
| `--ruled-paper-line` | `rgba(0, 0, 0, 0.06)` | 1px horizontal rule across process sheets |

---

## 3. High-Conversion Section Components

### 3.1 3-Tier Subscription Pricing Matrix (`PricingTierCard`)
- **Job:** Standard 3-column SaaS pricing grid with feature checklists.
- **Geometry:** Resting card: padding `2.25rem 1.75rem`, radius `20px`, border `1px solid #e2e8f0`. Highlighted tier card: `border: 2px solid #8b5cf6`, `transform: scale(1.04)`, `z-index: 2`. Badge: `top: -12px`, padding `0.25rem 0.875rem`, radius `9999px`.
- **States:** Resting: shadow `0 4px 6px -1px rgba(15, 23, 42, 0.05)`. Hover: `transform: translateY(-4px)`, shadow `0 20px 40px -15px rgba(15, 23, 42, 0.12)`. Easing: `300ms cubic-bezier(0.16, 1, 0.3, 1)`.

### 3.2 Interactive Split Pricing Blocks (`InteractiveSplitPricing`)
- **Job:** Two stacked plan blocks pairing interactive option toggles on the left with a cream price card on the right. Toggles serve as a confidence device showing inclusions.
- **Layout Grid:** `grid items-start gap-10 md:grid-cols-2 md:gap-14`.
- **Left Column:** Plan name (`h3`), "Choose what works for you" subhead, lead paragraph, and 3 selectable option buttons with active radio dots.
- **Right Column (Cream Price Card):**
  - Surface: `bg-[#FDFBF7]` (warm cream), border `1px solid #EAE5D9`, radius `24px`, padding `32px`.
  - Content: Label row, large numerical price (`$2,400`), billing frequency, perk checklist with checkmark icons, and `ShineButton` or booking CTA.

### 3.3 Scroll-Driven Services Showcase (`WhatWeDoSection`)
- **Job:** Dark, rounded services panel pinning service copy while paired project screenshots scroll past. When a new image pair enters view, copy swaps via filmstrip slide.
- **Shell:** Dark container, radius `32px`, overflow clip, background `#0A0A0A` (or `--night`), padding `80px 24px`.
- **Layout:** Two-column flex (`gap: 32px`, `align-items: flex-start`).
  - **Left Sticky Column:** Max width `396px`, height `580px`, `position: sticky; top: 20vh;`. Contains an overflow-hidden mask with a vertical track translating `translateY(-active * panelHeight)`. Each panel has category title, gradient divider (`1px solid linear-gradient(90deg, #43009D, #180037)`), body copy, and "See More" link.
  - **Right Scrolling Column:** Flex 1, row gap `50px`. Contains 4 image blocks (`height: 680px`, flex row, gap `32px`). Image box 1 (`396×556px`, radius `8px`), Image box 2 (`396×556px`, radius `8px`, `margin-top: 120px` for staggered rhythm).

### 3.4 Ruled-Paper Process Reassurance Stack (`WorkflowProcessSection`)
- **Job:** Warm, tactile process reassurance block answering "What happens after I book?" with 5 numbered sticky notes pinned to a sheet of ruled paper.
- **Shell:** `section.relative.overflow-hidden`, rounded top corners only (`rounded-t-[32px]`), flat tinted background `#F7F5F0`.
- **Ruled Paper Overlay:** `aria-hidden absolute inset-0 opacity-70 pointer-events-none` with repeating linear gradient:
  ```css
  background-image: repeating-linear-gradient(
    180deg,
    transparent 0px,
    transparent 31px,
    var(--ruled-paper-line) 31px,
    var(--ruled-paper-line) 32px
  );
  ```
- **Notes Column (`flex flex-col`):** 5 alternating notes. Each note is an `<article>` wrapped in a `Reveal` container with `md:w-[42%]` alternating `md:mr-auto` (steps 1, 3, 5) and `md:ml-auto` (steps 2, 4).
- **Note Internals:**
  - Tilt: Subtly rotated (`-1.5deg` to `+1.8deg`), hovering straightens to `0deg` with hover lift `translateY(-4px)` and shadow lift.
  - Glowing Pin Circle: Absolutely positioned, overlapping top edge of card (`size-6 rounded-full bg-brand shadow-glow`).
  - Content: Script numeral (`font-script text-2xl text-brand`), step title, and descriptive body.

### 3.5 Split Enterprise Hero with Depth-Scaled Cards (`EnterpriseHeroSection`)
- **Job:** Above-the-fold developer platform landing hero with action cluster and depth-scaled feature toggle cards.
- **Layout:** Container max width `1240px`, grid `lg:grid-cols-[minmax(0,520px)_1fr]`, gap `64px`, `items-center`.
- **Left Column:** Eyebrow link with Lucide `ChevronRight`, two-line headline (line 1 dark ink, line 2 gradient span), two-line lead, and action row (`ShineButton` + ghost link).
- **Right Column (Depth Stack):**
  - Background: Striped violet gradient glow, blurred.
  - Floating Card Stack: 7 rows of white cards (`radius-12 border border-border/60`), depth-scaled: middle cards render at `scale(1.0)` and `opacity: 1`, while top and bottom cards scale to `0.92` and fade to `opacity: 0.4` to suggest an infinite vertical list.
  - Card Switch: `40×22px` toggle track with `18px` circular knob.

---

## 4. Closed Sets & Layout Rules

1. Allowed pricing matrices: `three-column-grid` (desktop) and `single-column-stack` (mobile).
2. Allowed tier roles: `standard` (resting white) and `featured` (scaled 1.04, obsidian ground, gold badge). Only one featured tier per table.
3. Allowed workflow notes: Exactly 5 alternating steps with alternating horizontal alignment.

---

## 5. Anti-Hallucination & Quality Verification Checklist

- [ ] Interactive split pricing cards render option toggles on left and cream card (`#FDFBF7`) on right.
- [ ] "What We Do" services section uses sticky left column with `translateY(-active * panelHeight)` copy slide and right-side staggered image blocks (`mt-120px` on second image).
- [ ] Process workflow section paints ruled paper overlay using `repeating-linear-gradient` and alternates note alignment (`md:mr-auto` vs `md:ml-auto`).
- [ ] Sticky notes include glowing pin circle overlapping top card edge with subtle tilt (-1.5deg to +1.8deg).
- [ ] Enterprise hero right column implements 7-row depth-scaled card stack with blurred violet halo.
- [ ] Zero private or client names exist in component copy.
