# 39 Logo Construction

> **/goal** Design a reusable mark for any product once the name, the idea, and the colors are supplied.
> **/learn** This file is the spec. The prompts `01-logo-create.md` and `03-svg-logo.md` follow it. A missing input stops the job.

**Version:** 4.3.0
**Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** A clean mark, yes. A finished brand system, only after the three inputs exist.
**Ambiguity:** Without `product_name`, `product_idea`, and `sample_colors`, do not draw.

---

## 1. Inputs

Stop and ask when any of these is missing:

| Input | Meaning |
|---|---|
| `product_name` | The words the wordmark may set |
| `product_idea` | What the mark has to suggest |
| `sample_colors` | Hex or a theme id such as `bright-gold-tech` |

Optional: `inspiration_logos` (silhouette only, do not trace), `needs_dark_white_variants` (default yes), `is_animated` (default no).

For `bright-gold-tech`, the mark uses `--pres-accent` on a transparent canvas. The dark slide already supplies the black. Do not paint a black rectangle behind the mark.

---

## 2. Construction

Pick one archetype. Do not mix two in one mark.

| Archetype | Construction |
|---|---|
| `silhouette` | One filled shape and one cutout |
| `modular` | Two to four shapes, equal gap, no overlap |
| `ribbon` | One continuous path, rotational |
| `stroke` | One curve, thick to thin, no crossing |

Rules:

1. Readable at `16px`, comfortable at `48px`, balanced at `1024px`.
2. No tangled overlaps. Gaps stay visible at `16px`.
3. No background rectangle in the SVG. No solid plate behind a PNG. Alpha is `0`.
4. No embedded bitmap and no editor metadata.
5. Filenames are lowercase and numbered: `01-logo.svg`, `02-logo-dark.svg`, `03-logo-white.svg`.
6. Monochrome strokes use `currentColor`. A two-color mark uses `var(--pres-accent)` and `var(--pres-text)` or the supplied hex.
7. Root SVG has `viewBox` and no fixed `width` or `height`.
8. Include `<title>` and `<desc>`, `role="img"`, and `aria-labelledby`.

PNG sizes, only when asked: `52`, `128`, `256`, `512`, `1024`. Animation only when `is_animated` is yes: one SVG or one GIF, same geometry, no new shapes.

---

## 3. Example

This is a construction sample, not a brand. Replace the path when the inputs arrive.

```xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48" fill="none" role="img" aria-labelledby="mark-title mark-desc">
  <title id="mark-title">{PRODUCT_NAME}</title>
  <desc id="mark-desc">{PRODUCT_IDEA}</desc>
  <path fill="currentColor" d="M8 36 L24 8 L40 36 H32 L24 20 L16 36 Z"/>
</svg>
```

On a bright-gold slide the parent sets `color: hsl(var(--pres-accent))` and the image box is `height: 48px`.

---

## 4. Refuse

Do not output a web page, a React tree, or a stylesheet when the task is the mark. The palette note, when needed, is a separate `01-palette.md` that lists the theme id and the hex twins. It does not restate this file.
