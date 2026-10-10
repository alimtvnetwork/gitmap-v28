# 45 — AI Visual Whitespace Discipline for Generated Images

> **/goal** Give AI agents a non-negotiable whitespace-first discipline for rendering
> standalone images (social cards, infographics), so generated visuals stop coming
> back cramped, text-heavy, and overlapping.
> **/learn** The text-zone budget, spacing minimums, delete-don't-shrink rule, and
> the pre-delivery self-check below are mandatory before any AI-rendered image is
> shown to the user.

**Version:** 1.0.0
**Status:** Active
**Updated:** 2026-10-09

---

## 0. Why this file exists

Three rounds of user feedback on AI-rendered SEO infographics ("really really bad" →
"too many texts … no spacing" → "spacing is alright") converged on one root cause:
the AI filled the canvas with text and left no air. Whitespace is not empty space —
whitespace IS the design. This file embeds that as an instruction no agent may skip.

---

## 1. The whitespace mandate (non-negotiable)

1. **Empty space is a design element.** A 1200×630 card with four well-spaced blocks
   beats a card with seven cramped blocks every time.
2. **If it doesn't add meaning, delete it.** Never shrink a redundant line to make
   it fit — remove it. Duplicate kickers, keyword footers that restate the eyebrow,
   and decorative rules that touch text are the first things to go.
3. **One hero per image.** A single numeral or headline dominates (Ubuntu 700);
   everything else supports it. Two competing large elements = instant reject.

---

## 2. Text-zone budget (1200×630 `social-og-card`)

A *text zone* is one visually separated block of copy. Hard cap: **5 zones**,
preferred: **4**.

Canonical zone order:

| # | Zone | Type | Example |
|---|------|------|---------|
| 1 | Eyebrow kicker | Poppins 600, uppercase, tracked `+0.14em` | `SOFTWARE FRANCHISE · KUALA LUMPUR` |
| 2 | Hero numeral | Ubuntu 700, 190–220px, tabular numerals | `$2M`, `6%` |
| 3 | Headline | Ubuntu 700, 54–62px, max two lines | `franchise investment` |
| 4 | Support line (optional) | Poppins 400, 34–38px, muted token | `88% use AI · only 39% see earnings impact` |
| 5 | Source / brand (tiny) | Poppins 400, 22–24px, dim token | `McKinsey · State of AI 2025` |

Rules: zones 1–3 and 5 are the default set. Zone 4 is allowed only when it carries
the card's core contrast (e.g. the adoption-vs-impact gap); otherwise delete it.
Never add a sixth zone — restructure instead.

**Label-first variant (preferred for stat cards):** zone 1 becomes the metric label
(`AI HIGH PERFORMERS`, Poppins 600 tracked, accent color), zone 2 the hero numeral
below it. No company/brand kicker at the top — brand attribution never headlines
a stat card.

**Data-proof rule:** every claim on the card must be proven by real figures shown
on the card itself. The proof is rendered as a highlighted label — a pill or card
(surface fill, accent outline) with the key figures picked out in the accent color.
Plain muted running text does not count as proof display.

**Small-text budget:** at most one tiny line (source *or* brand, never both);
prefer neither when the data speaks for itself.

---

## 3. Spacing minimums

| Rule | Minimum | Notes |
|------|---------|-------|
| Between blocks | 60px | 80–90px preferred between hero and headline |
| Side margins | 80px | No text starts left of x=80 or ends right of x=1120 |
| Perimeter safe zone | 55px | Per `37-image-specifications.md`; no text inside it |
| Footer clearance | 24px | Smallest gap allowed anywhere (source line only) |

Measure gaps edge-to-edge of rendered glyphs (descenders count), not baselines.

---

## 4. Separated segments for multi-fact cards

When one card must carry 2–4 related facts, never stack them as a dense bullet
list. Render each fact as its own **segment**: a rounded surface card
(`radius 16–18px`, surface token fill, hairline outline) with generous internal
padding and a real gap (32–40px) between cards. One short label per card
(`50% PROFIT MARGIN`), value-first where possible. Three segments maximum per row.

---

## 5. Type and color discipline

- Headings/numerals: **Ubuntu 700**; body/labels: **Poppins** (400 body, 600 kickers).
  Per `04-typography.md` v3.2.0. Large Ubuntu sizes use negative tracking.
- Colors come from `17-theme-tokens.json` only: `warm-editorial` for brand/craft
  narratives, `navy-purple` for technical/franchise subjects. Do not invent palettes.
- Canvas: `social-og-card` (1200×630px) from the closed set in
  `37-image-specifications.md`. Do not invent dimensions.

---

## 6. Mistake log (what the bad versions did)

Recorded so no agent repeats them:

1. Seven text zones on one card; footer rule drawn through body copy.
2. Bullets hugging text with 12px gaps; three facts crammed in one column.
3. Keyword footer duplicating the eyebrow kicker word-for-word.
4. 300px hero numeral with 24px supporting text nobody could read at thumbnail size.
5. Decorative divider lines placed with no clearance check, visually striking
   through words.

---

## 7. AI pre-delivery self-check (mandatory)

Before showing any rendered image to the user:

- [ ] Count text zones: ≤ 5. If 6+, delete or restructure before delivery.
- [ ] Render and VIEW the actual PNG. Do not trust coordinates alone.
- [ ] No two glyph bounding boxes touch or overlap. Overlap = instant reject, re-render.
- [ ] Every inter-block gap ≥ 60px; side margins ≥ 80px.
- [ ] Read the card at thumbnail scale (≈300px wide): hero + headline legible.
- [ ] No invented colors, fonts, or canvas sizes — tokens only.

---

## 8. Sibling references

- Canvas geometries, safe zones, closed image palette: [`37-image-specifications.md`](./37-image-specifications.md)
- Typography tokens: [`04-typography.md`](./04-typography.md)
- Theme tokens (machine-readable): [`17-theme-tokens.json`](./17-theme-tokens.json)
- Spacing scale: [`05-spacing-layout.md`](./05-spacing-layout.md)
