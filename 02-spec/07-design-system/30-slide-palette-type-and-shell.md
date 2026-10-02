# 30 — Slide Palette, Type, and Shell

> **/goal** Fix the dark amber deck: every color, type size, and shell layer.
> **/learn** This file is the scripted dark deck. JSON theme ids stay in `27-slide-canvas-and-themes.md`. Pill chips stay in `32-slide-color-options.md`. Marketing pages do not use this file.

**Version:** 1.0.0
**Status:** Active

---

## 0. Anti-hallucination

Store colors as HSL triplets without an `hsl()` wrapper, then consume them as `hsl(var(--token))`. Do not write a raw hex inside a slide component. If a size is not in this file, do not invent it.

This palette is one theme: dark only. Do not invert it into a light mode.

---

## 1. Core palette

| Token | HSL | Hex | Use |
|---|---|---|---|
| `--background` | `240 9% 5%` | `#0B0B0E` | Canvas |
| `--background-elevated` | `240 8% 8%` | `#121217` | Cards, tooltips, controller pill |
| `--foreground` | `0 0% 100%` | `#FFFFFF` | Primary text |
| `--foreground-muted` | `240 5% 75%` | `#BFBFC7` | Sub-headlines, captions |
| `--foreground-subtle` | `240 4% 55%` | `#888892` | Inactive dots, axis labels |
| `--primary` | `38 91% 55%` | `#F5A623` | Accent, active dot, titles that shout |
| `--primary-hover` | `38 95% 62%` | `#F8B647` | Hover on amber controls |
| `--primary-glow` | `35 100% 50% / 0.35` | `#FF9500` | Radial glow, low opacity |
| `--border` | `240 6% 18%` | `#2A2A30` | Hairline inside the controller |
| `--border-subtle` | `240 6% 14%` | `#1F1F24` | Faint grid |

Chart bars, cool to warm, one sequence only:

| Step | HSL | Hex |
|---|---|---|
| 1 | `200 45% 55%` | `#5A9FBF` |
| 2 | `265 35% 60%` | `#8E7BB8` |
| 3 | `345 55% 60%` | `#C76B85` |
| 4 | `30 65% 55%` | `#D49147` |
| 5 | `42 90% 55%` | `#F0B23A` |

`--pres-bg` `240 20% 4%` (`#0a0a14`) and `--pres-accent` `41 100% 50%` (`#ffae00`) belong to a different scripted deck. Do not mix those two hex values into this palette on the same slide.

---

## 2. Spotlight, grid, surfaces

Spotlight, centered on the focal point:

```css
background:
  radial-gradient(
    ellipse 60% 55% at 50% 50%,
    hsl(35 100% 50% / 0.18) 0%,
    hsl(35 100% 50% / 0.06) 35%,
    transparent 70%
  ),
  hsl(var(--background));
```

On a `1920×1080` canvas that ellipse is about `1150×600`.

Grid:

```css
background-image:
  linear-gradient(hsl(0 0% 100% / 0.025) 1px, transparent 1px),
  linear-gradient(90deg, hsl(0 0% 100% / 0.025) 1px, transparent 1px);
background-size: 48px 48px;
```

| Surface | Background | Border | Shadow |
|---|---|---|---|
| Canvas | `--background` plus the spotlight | none | none |
| Controller pill | `hsl(240 8% 8% / 0.85)` and `backdrop-filter: blur(12px)` | `1px solid hsl(var(--border))` | `0 8px 24px hsl(0 0% 0% / 0.4)` |
| Tooltip | `--background-elevated` | `1px solid hsl(var(--border))` | `0 4px 12px hsl(0 0% 0% / 0.5)` |
| Chart bar | Linear gradient of the bar color, 12% lighter at the top | none | `0 4px 16px hsl(0 0% 0% / 0.35)` |

Focus ring on any control: `0 0 0 2px hsl(var(--background)), 0 0 0 4px hsl(var(--primary))`.

Top progress fill: `linear-gradient(90deg, hsl(210 90% 55%), hsl(var(--primary)))`. Track is `--border-subtle`. Height `4px`. Width transition `240ms ease-out`.

---

## 3. Type

Display is Ubuntu `400`, `500`, `700`. Body and UI are Poppins `300`, `400`, `500`, `600`. Inter is the body fallback on this dark deck only. Sizes are authoring pixels on `1920×1080`.

| Token | Size / line-height | Weight | Family | Use |
|---|---|---|---|---|
| `display-hero` | `128 / 1.05` | 700 | Ubuntu | Cover title. Tracking `-0.02em` |
| `display-xl` | `88 / 1.1` | 700 | Ubuntu | Content title |
| `display-lg` | `56 / 1.15` | 700 | Ubuntu | Section divider |
| `display-md` | `36 / 1.2` | 700 | Ubuntu | Chart title |
| `display-sm` | `28 / 1.3` | 700 | Ubuntu | Mini header |
| `body-lg` | `28 / 1.5` | 400 | Poppins | Bullet |
| `body-md` | `22 / 1.5` | 400 | Poppins | Subtitle under a hero |
| `body-sm` | `18 / 1.5` | 400 | Poppins | Footnote |
| `ui-md` | `20 / 1` | 500 | Poppins | Controller counter |
| `ui-sm` | `16 / 1` | 500 | Poppins | Tooltip title |
| `axis-label` | `16 / 1` | 400 | Poppins | Chart axis |

Other Ubuntu display tokens use tracking `-0.01em`. Poppins at `20px` and below uses tracking `0`.

Amber (`--primary`) is for titles that shout, stat numbers, and the tooltip's number prefix. Other text is `--foreground`. Captions and axis labels are `--foreground-muted`. Footnotes are `--foreground-subtle` at `80%` opacity.

---

## 4. Shell layers

Back to front. The slide component draws layer 4 only. The deck shell owns the rest.

| Layer | What | Rule |
|---|---|---|
| 0 | Canvas | `hsl(var(--background))` |
| 1 | Grid | `1px` lines, `48px` cells, `2.5%` white |
| 2 | Outline icons | 10 to 14 Lucide outlines, amber, opacity `0.08` to `0.12`, stroke `1.5px`, size `56px` to `96px` |
| 3 | Spotlight | Section 2 |
| 4 | Content | The layout in file 28 |
| 5 | Logo | Top-left, `48px` tall |
| 6 | Progress | `4px`, `top: 0` |
| 7 | Dots | File 31 |
| 8 | Controller | File 31 |

Safe area for content: `96px` left and right, `80px` top and bottom. Logo, controller, and dots sit outside that inset.

Decorative icons do not sit on a grid. Keep them out of a `900×400` ellipse at the center of a cover, and out of a `700×400` zone around a content title. They never overlap text. Allowed names: `Database`, `Cloud`, `Box`, `Code`, `Code2`, `Brain`, `Cpu`, `Palette`, `Layers`, `Globe`, `LayoutDashboard`, `Server`, `Workflow`. No duplicates within about `400px`.

Emoji bullets replace `•` in content lists. Size matches `body-lg` (`28px`). One emoji, then one space. `vertical-align: -2px`. Background icons stay faint. Emoji stay sharp. Do not swap those roles.
