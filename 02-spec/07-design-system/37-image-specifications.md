# 37 — Image Specifications and Canvas Geometries

> **/goal** Specify exact pixel geometries, safe zones, text rules, and color tokens for standalone marketing images, social cards, banners, and thumbnails.
> **/learn** Canvas sizes, collision zones, text caps, and image palettes live here. Web layout lives in `25-page-assembly.md` and slide stages live in `27-slide-canvas-and-themes.md`.

**Version:** 1.0.0
**Status:** Active

---

## 0. Anti-hallucination

If a dimension, coordinate, color, or text limit is not in this file or a file it cites, do not invent it. Image surfaces use their own dedicated palette. They do not borrow tokens from marketing pages or slide decks.

---

## 1. Scope

Use this file when generating or rendering standalone raster images (`.png`, `.webp`), social cards, YouTube thumbnails, channel banners, LinkedIn banners, and SVG vector marks.

Do not use this file for live DOM layout. Do not flatten selectable web page text into raster images.

---

## 2. Tokens (Image Palette)

Image surfaces use this closed palette. Do not borrow marketing blue `#2563EB` or slide stage `#0a0a14`.

| Token | Hex | HSL | Role |
|---|---|---|---|
| `--img-canvas-dark` | `#090D16` | `hsl(222, 43%, 6%)` | Dark comparison card background |
| `--img-card-surface` | `#111827` | `hsl(222, 47%, 11%)` | Card containers on dark images |
| `--img-badge-gold` | `#FFAA00` | `hsl(40, 100%, 50%)` | Badge fill, divider line, left symbol |
| `--img-badge-ink` | `#0F172A` | `hsl(222, 47%, 11%)` | Badge text on gold background |
| `--img-left-border` | `#F59E0B` | `hsl(38, 92%, 50%)` | Left comparison card outline |
| `--img-right-border` | `#38BDF8` | `hsl(199, 89%, 60%)` | Right comparison card outline |
| `--img-text-pure` | `#FFFFFF` | `hsl(0, 0%, 100%)` | High-contrast headline text |
| `--img-text-muted` | `#9CA3AF` | `hsl(218, 11%, 65%)` | Description text on cards |
| `--img-text-subtle` | `#94A3B8` | `hsl(215, 16%, 65%)` | Right symbol and footer links |
| `--img-dot-grid` | `rgba(255, 255, 255, 0.055)` | `hsl(0, 0%, 100% / 0.055)` | Background halftone dot grid |

---

## 3. Components

### 3.1 InfographicComparisonCard

- **Name and job:** `InfographicComparisonCard`. Renders a 1080×1080 square comparison graphic for social posts and documentation.
- **Anatomy:** Root canvas, 32px dot grid, header title, badge pill, subtitle, two comparison cards (left amber, right blue), central divider line, footer brand, and URL.
- **Geometry:** Canvas `1080×1080px`. Left card: x `80px`, y `280px`, w `430px`, h `480px`, radius `32px`. Right card: x `570px`, y `280px`, w `430px`, h `480px`, radius `32px`. Divider line: `(80, 830)` to `(1000, 830)`, width `3px`.
- **Tokens:** Canvas `--img-canvas-dark`, card fill `--img-card-surface`, left outline `--img-left-border`, right outline `--img-right-border`, badge fill `--img-badge-gold`, badge text `--img-badge-ink`.
- **States:** Static export. No interactive states.
- **Motion:** Not specified in source. Do not invent.
- **Reduced motion:** Not applicable to static raster export.
- **Accessibility:** Must include descriptive alt text matching headline and card points when embedded in HTML.
- **Closed options:** Size `1080×1080px`. Anything else does not exist.
- **Do not:** Do not render at arbitrary aspect ratios. Do not omit the footer URL or brand row.

### 3.2 YouTubeThumbnail

- **Name and job:** `YouTubeThumbnail`. Renders high-conversion 16:9 thumbnail artwork with strict safe zones.
- **Anatomy:** Background scene, left content zone (Zone A), subject portrait zone (Zone B), name and credential zone (Zone C), achievement zone (Zone D), lower information strip (Zone E).
- **Geometry:** Standard canvas `1280×720px`, high-res canvas `1920×1080px`. Edge safe margin `55px`. Zone A: x `25–405px`. Zone B: x `355–735px`, height 54%–58% of canvas. Zone C: x `655–965px`. Zone D: x `930–1245px`. Zone E: x `670–1280px`, y `548–720px`.
- **Tokens:** High-contrast text `#FFFFFF` and `#FFAA00`, background feathered linear gradient fade at bottom edge.
- **States:** Static export. No interactive states.
- **Motion:** Not specified in source. Do not invent.
- **Reduced motion:** Not applicable to static raster export.
- **Accessibility:** Title text rendered in image must match video title metadata.
- **Closed options:** `1280×720px` and `1920×1080px`. Anything else does not exist.
- **Do not:** Do not place text within 55px of the perimeter. Do not use clip-art, star badges, or unblended hard cutout edges on the subject.

### 3.3 YouTubeChannelBanner

- **Name and job:** `YouTubeChannelBanner`. Header canvas for channel homepage across TV, desktop, and mobile.
- **Anatomy:** Panoramic background, central safe-zone text and badge container.
- **Geometry:** Canvas `2560×1440px`. Central safe area `1232×338px` centered at `(1280, 720)`.
- **Tokens:** Base ground `--img-canvas-dark`, accent `--img-badge-gold`.
- **States:** Static export. No interactive states.
- **Motion:** Not specified in source. Do not invent.
- **Reduced motion:** Not applicable to static raster export.
- **Accessibility:** All readable channel copy must remain within the `1232×338px` safe zone.
- **Closed options:** Size `2560×1440px`. Anything else does not exist.
- **Do not:** Do not place text outside the 1232×338px center window.

### 3.4 LinkedInProfileBanner

- **Name and job:** `LinkedInProfileBanner`. Personal profile header graphic with avatar collision clearance.
- **Anatomy:** Panoramic canvas, right-aligned content area, left-bottom clear area.
- **Geometry:** Canvas `1584×396px` (4:1 aspect ratio). Safe area `1350×350px`. Avatar exclusion zone: left `40–300px`, bottom `0–140px` kept clear of critical text.
- **Tokens:** Ground `--img-canvas-dark`, accent `--img-badge-gold`, text `--img-text-pure`.
- **States:** Static export. No interactive states.
- **Motion:** Not specified in source. Do not invent.
- **Reduced motion:** Not applicable to static raster export.
- **Accessibility:** High contrast ratio (>= 4.5:1) for all typography over background.
- **Closed options:** Size `1584×396px`. Anything else does not exist.
- **Do not:** Do not put names, logos, or URLs in the bottom-left avatar exclusion zone.

### 3.5 LinkedInCompanyBanner

- **Name and job:** `LinkedInCompanyBanner`. Enterprise page header graphic.
- **Anatomy:** Panoramic canvas, right-aligned value proposition, logo clear area.
- **Geometry:** Canvas `1128×191px` (5.9:1 aspect ratio). Central safe area `1050×150px`.
- **Tokens:** Ground `--img-canvas-dark`, text `--img-text-pure`, accent `--img-badge-gold`.
- **States:** Static export. No interactive states.
- **Motion:** Not specified in source. Do not invent.
- **Reduced motion:** Not applicable to static raster export.
- **Accessibility:** Text must remain legible at 50% viewport scale on mobile screens.
- **Closed options:** Size `1128×191px`. Anything else does not exist.
- **Do not:** Do not exceed 8 words in the banner proposition headline.

---

## 4. Closed Sets

Allowed image canvas identifiers:

1. `card-infographic-square` (`1080×1080px`)
2. `youtube-thumbnail-std` (`1280×720px`)
3. `youtube-thumbnail-hd` (`1920×1080px`)
4. `youtube-banner` (`2560×1440px`)
5. `linkedin-profile-banner` (`1584×396px`)
6. `linkedin-company-banner` (`1128×191px`)
7. `social-og-card` (`1200×630px`)
8. `brand-logo-icon` (`512×512px`)

anything else does not exist.

Allowed image file naming:

- Format: `NN-kebab-name.{png|webp|svg}`
- Where `NN` is a two-digit zero-padded integer (`01`, `02`).
- Example: `01-comparison-card.png`, `02-channel-banner.png`.

anything else does not exist.

---

## 5. Refusal List

- Do not use marketing blue (`#2563EB`) or slide stage amber (`#ffae00`, `#0a0a14`) on image canvas surfaces. Use the image palette in section 2.
- Do not export unsequenced file names like `image.png` or `banner.jpg`.
- Do not invent canvas dimensions outside the closed sets in section 4.
- Do not render text without verifying safe zone clearances.
- Do not place text over the bottom-left avatar collision area in LinkedIn profile banners.

---

## 6. Sibling References

- Marketing web page shell and grid layouts: [`25-page-assembly.md`](./25-page-assembly.md)
- Presentation virtual canvas (1920×1080) and theme tokens: [`27-slide-canvas-and-themes.md`](./27-slide-canvas-and-themes.md)
- Upstream creation prompts:
  - SVG Logo: `01-prompts/10-ui-and-design/01-svg-logo.md`
  - Logo Design: `01-prompts/10-ui-and-design/02-logo-create.md`
  - LinkedIn Profile Banner: `01-prompts/10-ui-and-design/03-linkedin-profile-banner.md`
  - YouTube Thumbnail: `01-prompts/10-ui-and-design/04-youtube-thumbnail-create.md`
  - LinkedIn Company Banner: `01-prompts/10-ui-and-design/05-linkedin-company-banner.md`
  - YouTube Channel Banner: `01-prompts/10-ui-and-design/06-youtube-channel-banner.md`
