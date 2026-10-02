# 28 — Slide Layouts & AI Blind-Authoring Contracts

> **/goal** Codify the closed set of slide layouts, slot definitions, and the 3-axis blind-authoring mental model so any AI can create or fix layouts without rendering errors.
> **/learn** Master the 3-axis layout rule (vertical section placement vs internal sibling spacing vs runtime brand insets), CSS Grid track packing (`min-content` vs `1fr`), compact card mechanics, and canonical slide templates.

**Version:** 4.2.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. The 3-Axis Layout Mental Model (Non-Negotiable)

Every slide section is a single flex column owning three completely independent axes. **Never confuse them:**

| Axis | Lever (CSS Property) | Applied Element | Architectural Job |
|---|---|---|---|
| **A. Vertical Placement** | `justify-content` (`justify-center`, `justify-start`) | Outer `<section>` | Positions the entire eyebrow→title→content→caption block on canvas (top / middle / bottom). |
| **B. Internal Sibling Spacing** | `mb-*` / `mt-*` / `gap-*` | Child components (header, grid, caption) | Controls tightness of the header→content stack. Never tweak to move the whole group. |
| **C. Horizontal Brand Alignment** | `paddingLeft` / `paddingRight` | Outer `<section>` as **inline style** | Uses `var(--brand-inset-x)` so headlines align with the brand logo. Never use Tailwind `px-*`. |

### 1.1 Lever Selection Rules
- **Do NOT** move the content group up/down by tweaking `pt-*` / `pb-*` / `mt-*` on children (Axis B leaking into Axis A). **Use `justify-content` on the section.**
- **Do NOT** tighten header-to-grid spacing by adding `flex-1` or `content-center` to the grid (Axis A leaking into Axis B). **Use `mb-6` on the header.**
- **Do NOT** align with the logo using `px-24` or any Tailwind utility. Brand inset is a runtime token that drifts across screen sizes. **Use `style={{ paddingLeft: 'var(--brand-inset-x)', paddingRight: 'var(--brand-inset-x)' }}`.**

---

## 2. CSS Grid Track Model & Compact Card Mechanics

When layout slots contain cards, two independent properties determine card height and spacing:

### 2a. Track Sizing (`grid-auto-rows`)
- `grid-auto-rows: 1fr`: Splits available vertical height evenly among rows. Two compact cards in one column become tall cards taking 50% height each.
- `grid-auto-rows: min-content`: Each row is **strictly the height of its content**. Two compact cards pack tightly at the top with `row-gap`.
- Auto-packing rule in `index.css`:
  ```css
  .slide-grid-2-equal:has(.slide-card.is-compact),
  .slide-grid-5-7:has(.slide-card.is-compact),
  .slide-grid-4-8:has(.slide-card.is-compact),
  .slide-grid-3-9:has(.slide-card.is-compact) {
    grid-auto-rows: min-content;
    align-content: start;
    row-gap: 1rem;
  }
  ```

### 2b. Item Alignment (`align-self`)
- `align-self: stretch` (default): Card expands to fill its row track. Padding adjustments have zero visible effect on card height.
- `align-self: start` (applied by `.slide-card.is-compact`): Card hugs content, anchored to top of the row.

### 2c. Decision Matrix: Symptom to Lever

| Symptom in Preview | Diagnosis | Surgical Lever |
|---|---|---|
| Compact card is as tall as a hero card | `align-self: stretch` winning | Add `"compact": true` in JSON; remove any `h-*` or `min-h-*` overrides. |
| Two compact cards spread to top and bottom | `grid-auto-rows: 1fr` splitting space | Ensure container selector `:has(.slide-card.is-compact)` triggers `grid-auto-rows: min-content`. |
| Content group sits too close to top bar | Outer section `justify-start` without padding | Set `<section className="... justify-center pt-24 pb-40">`. |
| Headline misaligned with top-left logo | Tailwind `px-*` drifting from brand inset | Apply `style={{ paddingLeft: 'var(--brand-inset-x)', paddingRight: 'var(--brand-inset-x)' }}`. |

---

## 3. Canonical Slide Templates

### 3.1 Cover Slide (Section Opener)
- **Title:** `display-hero` (Ubuntu 700, `128px`), `hsl(var(--primary))`, centered on both axes.
- **Subtitle / Kicker:** `body-md` (Poppins 400, `22px`), `hsl(var(--foreground-muted))`, centered, `24px` below title baseline.
- **Spotlight:** Radial aura centered exactly behind the title text.
- **Exclusion Zone:** Decorative icons stay outside an inner `900×400px` rectangle around title.

### 3.2 Content Slide (Two-Column Text & Bar Chart)
- **Geometry:** 1920×1080 stage. Left column starts at `x: 96`, max width `1100px`. Right column chart container at `x: 1180, y: 280, w: 640, h: 520`.
- **Left Column:** Headline (`display-xl`, 88px) + mini-header (`display-sm`, 28px) + 9 emoji-led bullet items with `40px` vertical rhythm.
- **Right Column (Growth Chart):** 5 vertical bars (Month 1 → 18%, Month 3 → 38%, Month 6 → 58%, Month 9 → 78%, Month 12 → 100%), widths `72px`, gap `48px`, `border-radius: 8px 8px 0 0`, vertical gradient base to +12% lightness, drop shadow `0 4px 16px rgba(0,0,0,0.35)`.

### 3.3 Stat Slide (Big Number)
- Centered 200px+ tabular figures number in Ubuntu 700 with label below, radial accent glow centered under digits.

---

## 4. Closed Layout Union & Slot Specifications

```text
left | center | steps | timeline | process | quote | bullets | image
poll | qa | embed | reveal-grid | counter-stat | typewriter
compare | priority | depth-stack
```

| Layout Type | Required Slots | Optional Slots | Layout Architecture |
|---|---|---|---|
| **`left`** | `heading` (RichText) | `kicker`, `body`, `media` (`src`, `alt`) | Text left half, media right half |
| **`center`** | `heading` (RichText) | `subhead`, `display` | Centered text, no side column |
| **`steps`** | `heading`, `steps[]` | `media` (`src`, `alt`, `caption`, `fit`) | Primary active step, remaining dimmed |
| **`timeline`** | `items[]` | `heading` | Single vertical rail, accent nodes |
| **`process`** | `stages[]` | `heading`, `subhead` | One row or column of circles with icons |
| **`quote`** | `quote` (RichText) | `attribution` | Centered editorial quote |
| **`bullets`** | `heading`, `bullets[]` | `kicker` | Single bullet list with sequential step reveals |
| **`counter-stat`** | `stats[]` (1–4) | `heading` | Numerical counter with `value`, `prefix`, `suffix`, tabular nums |
| **`typewriter`** | `lines[]` (1–6) | `heading`, `richLines` | Incremental character typing with sound option |
| **`compare`** | `before`, `after` | `heading` | Before/after image wipe with `durationMs` |
| **`depth-stack`** | `sentences[]` (1–6) | `heading` | 3D perspective card stack (`perspective: 1200`) |

---

## 5. Anti-Hallucination & Quality Verification Checklist

- [ ] Vertical placement uses `justify-content` on outer `<section>`, never margins on child elements.
- [ ] Horizontal inset strictly binds `style={{ paddingLeft: 'var(--brand-inset-x)', paddingRight: 'var(--brand-inset-x)' }}`.
- [ ] Compact cards specify `"compact": true` and have `align-self: start` in CSS.
- [ ] Grid containers with compact cards trigger `grid-auto-rows: min-content` and `align-content: start`.
- [ ] Layout type strictly belongs to the closed union of 17 layout identifiers.
- [ ] Zero private company or client names exist in slide content.
