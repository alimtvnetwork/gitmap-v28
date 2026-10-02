# Bright Gold Tech: Tokens, Background, Chrome

> **/goal** Fix the color, the dot field, the card, and the footer so a second agent cannot drift.
> **/learn** These values are copied from the bright-gold preset and from `.slide-dark-bg`. Hex in the table is the computed twin of the HSL. Components still consume the HSL token.

**Version:** 4.3.0
**Updated:** 2026-10-02
**Status:** Active

---

## 1. Tokens

Theme id: `bright-gold-tech`. Appearance: dark only. Do not invent a light inversion.

| Token | HSL | 8-bit twin | Use |
|---|---|---|---|
| `--pres-bg` | `240 20% 4%` | `#08080C` | Canvas |
| `--pres-bg-surface` | `240 17% 8%` | `#111118` | Raised surface |
| `--pres-bg-card` | `240 15% 11%` | `#181820` | Card fill |
| `--pres-accent` | `41 100% 50%` | `#FFAE00` | Gold |
| `--pres-accent-glow` | `41 100% 65%` | `#FFC64D` | Glow and gradient end |
| `--pres-accent-dark` | `33 100% 40%` | `#CC7000` | Chart bar foot |
| `--pres-text` | `0 0% 100%` | `#FFFFFF` | Primary type |
| `--pres-text-muted` | `218 11% 75%` | `#B8BDC6` | Eyebrow index, captions |
| `--pres-text-dim` | `220 9% 46%` | `#6B7280` | Page index |

Implement the HSL column. A preset swatch also writes `#0A0A14` and `#FFAD01` for this same pair. Do not put those swatch strings in new CSS.

`--background`, `--foreground`, `--primary`, and `--accent` on this theme copy `--pres-bg`, `--pres-text`, `--pres-accent`, and `--pres-accent`. Primary text on a gold fill is `--pres-bg`, not white.

Display font: Ubuntu. Body font: Poppins. Headings use the display font. Slide body uses the body font.

Do not use `#F5A623`, `#0B0B0E`, or the noir-gold pair `#0D0D0D` / `#C9A84C` on this theme. Those belong to other files.

---

## 2. Canvas

The slide canvas is `1920px` by `1080px`. The stage centers it and scales with:

```css
.slide-scaler {
  position: absolute;
  width: 1920px;
  height: 1080px;
  left: 50%;
  top: 50%;
  margin-left: -960px;
  margin-top: -540px;
  transform-origin: center center;
}
```

A website section that quotes this theme uses the same tokens. It does not use the `1920×1080` stage unless it is a slide.

---

## 3. Background

Every bright-gold page uses this field. The dot grid is gold at low opacity, masked so the center is visible and the corners fall off. Do not replace it with a line grid, a purple glow, or a flat black fill.

```css
.slide-dark-bg {
  background: hsl(var(--pres-bg));
  position: relative;
}
.slide-dark-bg::before {
  content: "";
  position: absolute;
  inset: 0;
  background-image:
    radial-gradient(ellipse 600px 600px at 55% 50%, hsl(41 100% 50% / 0.5), transparent);
  opacity: 0.22;
  pointer-events: none;
}
.slide-dark-bg::after {
  content: "";
  position: absolute;
  inset: 0;
  background-image:
    radial-gradient(circle 2px at center, hsl(41 80% 50% / 0.08) 2px, transparent 2px);
  background-size: 16px 16px;
  mask-image:
    radial-gradient(ellipse 700px 550px at 55% 50%, #000 0%, rgba(0, 0, 0, 0.35) 55%, rgba(0, 0, 0, 0.1) 100%);
  pointer-events: none;
}
```

Content sits above both layers (`z-index` greater than the pseudo elements).

---

## 4. Card, rule, pill, emphasis

```css
.glass-card {
  background: linear-gradient(135deg, hsl(var(--pres-bg-surface) / 0.8), hsl(var(--pres-bg-card) / 0.6));
  border: 1px solid hsl(var(--pres-accent) / 0.08);
  backdrop-filter: blur(20px);
  border-radius: 24px;
}
.accent-bar {
  height: 4px;
  width: 72px;
  background: linear-gradient(90deg, hsl(var(--pres-accent)), hsl(var(--pres-accent-glow)));
  border-radius: 2px;
}
.gold-emphasis {
  color: hsl(var(--pres-accent));
}
.gold-pill {
  display: inline-flex;
  align-items: center;
  border: 1px solid hsl(var(--pres-accent));
  color: hsl(var(--pres-accent));
  background: transparent;
  border-radius: 999px;
  padding: 8px 16px;
  font-size: 14px;
  line-height: 1;
}
```

The `72px` rule width and the `8px 16px` pill padding are the chrome minimums for this theme. Do not grow the pill into a filled gold button on a definition page.

Hover on a recommendation card, when that page type is used:

```css
.rec-card {
  transition: transform 0.3s ease, border-color 0.3s ease, box-shadow 0.3s ease;
}
.rec-card:hover {
  transform: scale(1.06) translateY(-8px);
  border-color: hsl(var(--pres-accent) / 0.6);
  box-shadow:
    0 16px 52px hsl(var(--pres-accent) / 0.18),
    0 0 36px hsl(var(--pres-accent) / 0.1);
}
```

Entrance, when motion is allowed:

```css
@keyframes gold-rise {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}
.gold-rise {
  animation: gold-rise 0.75s cubic-bezier(0.25, 0.46, 0.45, 0.94) both;
}
```

Stagger sibling blocks by `0.22s`. Under `prefers-reduced-motion: reduce`, set animation and transition duration to `0.01ms`.

---

## 5. Footer chrome

A full-width `.accent-bar` sits above the footer row. Set that instance to `width: 100%` and keep `height: 4px`.

The footer row holds three slots, left to right:

| Slot | Align | Color | Content |
|---|---|---|---|
| Mark | left | the logo file | Wordmark from `39-logo-construction.md`. Height `48px`. |
| Address | right | `--pres-text-muted` | Optional. The deck's own address only. |
| Index | right, under the address | `--pres-text-dim` | `{current} / {total}`, two digits each. |

Header logo on pages that show a mark at the top uses the same asset at height `48px`, inset `60px` from the left and `36px` from the top.
