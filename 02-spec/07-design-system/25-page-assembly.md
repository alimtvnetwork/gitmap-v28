# 25 — Page Assembly (Sites and Blogs)

> **/goal** Compose a marketing page or a blog only from the White Blue theme. Do not invent a section, a color, or a motion value.
> **/learn** Shell order, band rhythm, generic menu groups, and which section file owns each pattern.

**Version:** 1.0.0
**Status:** Active

---

## 0. Anti-hallucination

If a length, color, duration, or section name is not in this file or in the files it cites, do not invent it. Marketing pages use the White Blue palette. They do not use the slide palette in `27-slide-canvas-and-themes.md`.

Read, in order, before drawing a page:

1. `02-spec/07-design-system/04-white-blue-theme/readme.md`
2. `02-spec/07-design-system/04-white-blue-theme/01-colors-typography-and-tokens.md`
3. `02-spec/07-design-system/04-white-blue-theme/02-header-mega-menu-and-footer.md`
4. `02-spec/07-design-system/04-white-blue-theme/03-buttons-motion-and-interactions.md`
5. `02-spec/07-design-system/04-white-blue-theme/04-cards-heroes-and-section-library.md`

Copy caps for a body section: eyebrow at most 3 words, H2 at most 8 words, lead at most 28 words, card body at most 32 words. The H1 is at most 9 words (`20-ai-training-and-checklist-guide.md`, `13-section-patterns.md`). Do not add a cap beyond these.

Menu components and their entrance timings are `33-mega-menu-components.md`. Buttons are `04-white-blue-theme/03-buttons-motion-and-interactions.md` plus `HeaderShineButton` in file 33.

---

## 1. Shell

Every public page uses this order. Do not skip the header or the footer. Do not place two heroes.

1. `SiteHeader` from `04-white-blue-theme/02-header-mega-menu-and-footer.md`
2. Exactly one hero. Default is the light-first split hero (`CapabilityStack`) from `04-cards-heroes-and-section-library.md` section 3.
3. Three to seven body sections, chosen only from section 3 of this file.
4. One closer. Default is the lead form from `04-cards-heroes-and-section-library.md` section 7.3, or a final band that reuses an allowed section.
5. Footer from `04-white-blue-theme/02-header-mega-menu-and-footer.md`

Legal, error, and search pages may omit the closer. They still use the header and the footer.

---

## 2. Band rhythm

Adjacent sections never share a tone. Alternate `light`, `soft`, and a rare `dark` band, matching the White Blue readme (`light` then `soft` then `light` then `dark`). A `dark` band uses `GlassCard` and `text-on-dark` (`#D0D7E5`). A light band uses `card-premium` or `row-premium`. Do not stack two `soft` bands.

---

## 3. Allowed body sections

Use only these patterns. The anatomy and colors stay in the cited file.

| Id | Use | Source |
|---|---|---|
| `capability-stack` | Split hero with a self-assembling stack | `04-cards-heroes-and-section-library.md` section 3 |
| `solutions-grid` | Flagship card grid under a sticky editorial header | same file, section 4 |
| `workflow-board` | Pinned sticky-note process | same file, section 5 |
| `scroll-stack` | Fluted glass sticky card deck | same file, section 6 |
| `capability-tabs` | Tabbed capability board on `.band-soft` | same file, section 7.1 |
| `pricing` | Scope-builder pricing | same file, section 7.2 |
| `lead-form` | Lead form card | same file, section 7.3 |
| `hero-page` | Breadcrumb, eyebrow, one accent word, lead, two CTAs, optional 3-item proof | same hero type scale |
| `hero-split` | Copy column plus one media card (`TiltCard`) | file 33, max tilt `8deg` |
| `hero-form` | Value copy plus a form card. Do not animate the form | `card-premium` |
| `hero-editorial` | Category pill, title, byline, cover. No mask on a long title | type scale in the color file |
| `rail-tech` | Mono label plus a bordered logo grid | soft, tight |
| `band-stats` | Copy, CTA, `CountUp` grid | file 33 |
| `rail-testimonial` | One quote, portrait, arrows | soft |
| `grid-clients` | Logo tiles, grayscale until hover | soft |
| `band-certs` | Mono label plus a badge row | light |
| `list-values` | Numbered rows, mono index | light |
| `grid-cards` | 2 or 3 columns, icon, title, body | `card-premium` |
| `grid-industry` | Dense tiles, hover reveals one blurb | soft |
| `accordion-faq` | Accordion plus one question card | light |
| `prose-body` | Long form on a `720px` measure | blog article |
| `cta-band` | Closing band, always `soft`. The band before it is `light` | one H2, proof, one primary and two secondary actions |

Card primitives allowed inside those sections: `SurfaceCard`, `GlassCard`, `NeuCard`, `Pill`, `card-premium`, `row-premium`. Hover on `card-premium` grows a `2px` left rule with `scaleY(0)` to `scaleY(1)` over `--dur-base` and `--ease-out`. It does not translate the card. `SurfaceCard` is the only card that uses `hover:-translate-y-1`.

---

## 4. Page kinds

### 4.1 Marketing page

`SiteHeader` built from `33-mega-menu-components.md`, one hero (`capability-stack`, `hero-page`, `hero-split`, or `hero-form`), then three to seven rows from section 3. Close with `cta-band` or `lead-form`, then the footer. The band before `cta-band` is `light`.

### 4.2 Blog index

`SiteHeader`, a short hero that reuses the split-hero type scale (one `text-h1`, one lead), then a `solutions-grid` of article cards (title, one-line summary, meta in `text-eyebrow`). No pricing block. No workflow board. Footer follows. Filters, if present, sit in one band and do not invent a new control style: use the native select rules in `02-spec/07-design-system/22-native-css-select-and-border-shapes.md`.

### 4.3 Blog article

`SiteHeader`, one editorial hero (title at most 9 words, eyebrow in `text-eyebrow` at `12px` mono), then a prose column on `--paper` using body type from `01-colors-typography-and-tokens.md` (`--font-body`, Poppins). End with one `lead-form` or a related-cards row that reuses `solutions-grid`. Footer follows. Do not add a second H1.

---

## 5. Menu groups

Top-level groups, in this order, unless the product has fewer:

1. Product
2. Industries
3. Services
4. About
5. Resources

Contact and legal links live in the footer and in the header CTA pair, not as extra mega columns. Each mega panel follows `02-header-mega-menu-and-footer.md`: `SlideSwapLabel` with `stagger={0.04}`, underline `scale-x-0` to `scale-x-100` over `--dur-fast`, safe region `pad = 14`, close delay `110ms` to `220ms`, chevron `180deg`, promo flip `820ms` at `perspective: 1400px`. Do not retune those numbers.

Column templates stay as specified there:

- 3 or more groups: `lg:grid-cols-[1fr_1fr_1fr_0.9fr]`
- 2 groups: `lg:grid-cols-[1fr_1fr_1.1fr]`
- 1 group: `lg:grid-cols-[1.6fr_1fr]`

---

## 6. Sibling References

- Standalone marketing images, banners, and thumbnails: [`37-image-specifications.md`](./37-image-specifications.md)
- Detailed 3-tier pricing table and floating editorial card components: [`38-card-and-pricing-components.md`](./38-card-and-pricing-components.md)

---

## 7. Refusal list

- Do not add a section id that is not in section 3.
- Menu parts that are not in `33-mega-menu-components.md` do not exist.
- Do not introduce Inter, a third accent, or a slide amber (`#ffae00`, `#F5A623`, `#FFD83A`) on a marketing page.
- Do not change header height (`72px`), scroll threshold (`12px`), or menu motion.
- Do not name a client, a vendor, or a private repository in generated copy or comments.
