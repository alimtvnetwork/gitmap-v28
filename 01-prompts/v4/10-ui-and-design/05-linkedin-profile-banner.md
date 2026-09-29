# LinkedIn Profile Banner Design — Visual Identity & Authority Workflow

> **Prompt Version:** 1.0.0
> **Target Environment:** Lovable, Figma, Canva, Image Generation AI (Flux, Ideogram, Midjourney) & Design AI Platforms
> **Synchronization:** Main Meta-Repo & Connected Workspaces

This prompt instructs design AI platforms on how to craft high-authority, conversion-focused personal LinkedIn profile banners with exact dimensional standards, profile picture collision safe zones, and zero text hallucination.

---

## Strictly Avoid (Critical Negative Constraints)

> [!CAUTION]
> **ZERO TOLERANCE FOR AVATAR COLLISION, TEXT HALLUCINATION & WEB CODE (AUTO-REJECT)**
> The AI MUST strictly avoid generating HTML, placing text in the avatar collision zone, or distorting typography.

1. **NO HTML Files or Web Code:** NEVER generate `.html` files, React components (`.tsx`, `.jsx`), or CSS stylesheets. This is strictly a banner design and vector layout workflow.
2. **NO Text in the Avatar Collision Zone:** NEVER place critical text, headlines, contact info, faces, or logos in the bottom-left area (x: 0 to 380px). On both desktop and mobile, LinkedIn's circular profile picture will obscure anything placed there.
3. **NO Unblended Subject Cutouts (MANDATORY BOTTOM & EDGE BLEND):** When a photo or avatar cutout is provided, it MUST blend naturally into the bottom edge and surrounding background using soft gradient feathering or dark vignettes. NEVER paste a cutout with hard, abrupt edges.
4. **NO Tacky Clip-Art, Star Stickers, or Fake Icons (TOTAL BAN):** NEVER generate cartoonish star stickers (e.g. `#01` star badges), clip-art badges, or arbitrary icons. Do NOT invent or create avatars or icons unless explicitly requested by the user. Professional banners rely on typography, authentic photography, and clean minimalist layouts.
5. **NO Chaotic Overlapping Lines Across Subjects:** NEVER run background lines, graph strokes, or accent geometry across the subject's face, neck, or body. Elements must stay behind the subject or maintain clean negative space.
6. **NO Text Hallucinations or Gibberish:** NEVER render garbled glyphs, pseudo-Latin, or misspelled words. Every character in the headline, website URL, and credentials must match the user's input verbatim.
7. **NO Low-Resolution Exports:** NEVER export or specify standard 72 DPI 1x images that blur upon upload. Always specify **2x or 3x Retina resolution** (`3168 x 792 px` or `4752 x 1188 px`).
8. **NO Low-Contrast Text:** NEVER place light text over bright, busy backgrounds or dark text over dark shadows without proper contrast treatment (e.g. solid badge backdrops, dark vignettes, or subtle dropshadows).
9. **NO Unsequenced or Uppercase Filenames:** ALL generated files and images MUST use lowercase kebab-case preceded by a two-digit zero-padded sequence number (e.g. `01-banner-1584x396.png`, `02-banner-3168x792.png`, `01-prompt.md`, `02-plan.md`). Unsequenced names like `banner.png` or uppercase paths like `Banners/` are strictly BANNED.

---

## 1. Dimensional Standards & Safe Zones

LinkedIn applies aggressive cropping and compression to profile banners across devices. All banners must be designed against these precise dimensional specifications:

### Canvas Dimensions

- **Standard Dimensions (1x):** `1584 x 396 px` (4:1 aspect ratio).
- **High-Resolution Retina (2x — Recommended):** `3168 x 792 px` (guarantees razor-sharp text after LinkedIn's upload compression).
- **Ultra High-Resolution (3x):** `4752 x 1188 px`.
- **Target File Size & Format:** Under 8MB; export as uncompressed PNG or high-bitrate WebP.

### Safe Zone Architecture

```
+-----------------------------------------------------------------------------+
|                                    TOP SAFE BAND                            |
|                                                                             |
|   [ AVATAR COLLISION ]            CENTER / RIGHT CONTENT ZONE               |
|   [  DEAD ZONE       ]     - Headline & Value Proposition                   |
|   [  (Profile Pic)   ]     - Core Pillars / Credential Badges               |
|   [                  ]     - Website URL & Call-to-Action                   |
|                                                                             |
+-----------------------------------------------------------------------------+
|<--- 0 to 380px ------>|<----------------- 380px to 1584px ----------------->|
```

- **Dead Zone (Bottom-Left):** From `x: 0` to `x: 380px` and `y: 180px` to `y: 396px`, the circular avatar covers the canvas on desktop and shifts further inward on mobile. Keep this area free of text.
- **Mobile Crop Margin:** Mobile devices crop approximately `100px` from the left and right edges. Keep all critical messaging within the central safe zone (`x: 400px` to `x: 1450px`).

---

## 2. Input Capture & Verification

Before generating any design concepts or layout prompts, capture and validate the following inputs:

1. `full_name`: Person's exact name (e.g., "{full_name}").
2. `professional_headline`: Core value proposition (e.g., "Fractional CMO | Scaling B2B SaaS from $1M to $10M ARR").
3. `current_company` / `venture`: Current organization, agency, or brand.
4. `core_pillars`: 3–4 key focus areas (e.g., `GTM Strategy • Demand Generation • Team Leadership`).
5. `website_url`: Personal or business website URL (e.g., `https://example.com`).
6. `email_address`: Professional contact email address to embed on the banner (e.g., `contact@example.com`).
7. `qr_code`: Optional scannable QR code destination URL or asset to embed (e.g., linking directly to a calendar booking page, digital vCard, or portfolio).
8. `contact_or_cta`: Explicit call to action (e.g., "Book a strategy audit at example.com", "DM for speaking inquiries").
9. `portrait_photo_or_avatar`: High-resolution photo cutout of the person, visual description, or avatar image (optional — ask the user if they have an image to provide, or skip if they prefer a typography/minimalist layout).
10. `font_family`: Primary font family (Defaults to **Ubuntu** across all typography, headlines, subtitles, and overlays).
11. `brand_colors`: Primary brand palette (e.g., Midnight Navy `#0b0f19`, Vibrant Cyan `#06b6d4`, Pure White `#ffffff`).

> [!IMPORTANT]
> Always ask the user if they have a **person's photo, portrait cutout, or avatar image** to include on the banner, or if they prefer to **skip** it (using a clean typography or minimalist layout like Version 2). Remember: do NOT generate unsolicited avatars, clip-art icons, or star badges. Also ask if they want to embed an **email address** or a **scannable QR code** (e.g. for booking or vCard). If `full_name`, `professional_headline`, `website_url`, or `brand_colors` are missing or not provided, **STOP and ask the user** before proceeding.

---

## 3. Photographic, Lighting & Compositional Specification

When a subject portrait or human photo is included on the banner (typically on the right flank), apply deliberate cinematic direction:

### A. Camera Shot & Framing

- **Lens & Optics:** 85mm prime portrait lens equivalent with a wide aperture (`f/1.8` to `f/2.8`), creating an authentic shallow depth of field where the subject remains razor-sharp while the background dissolves into soft, creamy bokeh.
- **Shot Distance:** Eye-level medium close-up or bust shot (from mid-chest to head). Never extreme wide-angle or fisheye distortion that distorts facial proportions.
- **Placement:** Anchor the subject on the right flank of the panoramic canvas (`x: 1000px` to `x: 1500px` at 1x, or `x: 2000px` to `x: 3000px` at 2x). This preserves the entire center area for typography and leaves the left flank clear of avatar collisions.
- **Gaze & Expression:** Direct or slightly off-axis confident gaze. Natural, authentic executive expression (calm authority, approachable expertise).

### B. 3-Point Cinematic Lighting

- **Key Light:** Large diffused softbox at 45 degrees, providing soft, flattering facial modeling with gentle tonal transitions.
- **Rim / Edge Light (Crucial):** Distinctive colored rim light matching the brand accent (e.g. vibrant cyan, warm gold, or electric blue) tracing the hair, shoulders, and jawline from behind. This creates clean separation from dark backdrops.
- **Fill Light:** Low-intensity ambient fill preserving fine fabric weaves and subtle shadow details.

### C. Texture & Environmental Realism

- **Organic Skin Texture:** Visible skin pores, natural micro-contrast, and authentic specular highlights. Total ban on plastic, waxy, or airbrushed AI skin.
- **Crisp Textile Weave:** High-fidelity fabric texture on suits, blazers, or collars with realistic stitching.
- **Catchlights & Reflections:** Crisp, natural catchlights in the pupils of the eyes and subtle realistic reflections on glasses if worn.

### D. Mandatory Bottom & Edge Gradient Fade

- **Natural Blending:** The subject's torso/suit must NEVER be abruptly sliced with a hard pixel edge at the bottom or side of the banner.
- **Gradient Dissolve:** Apply a soft vertical linear gradient fade (or a dark feathered vignette) across the bottom 25% of the subject cutout, allowing the clothing to melt organically into the dark lower border of the canvas.

---

## 4. Visual Layout Versions (Multiple Modes)

### Version 1: Authority & Social Proof Layout (Default)

Ideal for consultants, founders, and executives:
- **Background:** Deep navy or charcoal gradient with subtle geometric or architectural depth of field.
- **Left Flank:** Clean, ambient negative space accommodating the profile picture.
- **Center:** Bold two-line headline in high-contrast white and gold/cyan accent, followed by 3 core pillars separated by bullet dots.
- **Right Flank:** High-resolution subject portrait cutout (with mandatory soft gradient bottom fade / dark feathered vignette), or authority proof (e.g. "Featured in Forbes, TechCrunch" or book cover mockup).
- **Lower Band:** Clear website URL and rounded pill call-to-action badge.

### Version 2: Minimalist Executive Layout

Ideal for enterprise leaders, investors, and board members, or when the person's photo is skipped:
- **Background:** Monochromatic dark slate or textured matte finish with subtle linear accent lighting.
- **Left Flank:** Clean dark negative space.
- **Center-Right:** Single, powerful positioning statement in elegant modern typography.
- **Right Flank:** Minimalist brand logo mark and official corporate website URL.

### Version 3: Visual Showcase / Product Showcase Layout

Ideal for creators, authors, and product builders:
- **Background:** Studio gradient with warm ambient lighting.
- **Center:** Creator name and hook headline.
- **Right Flank:** 3D product mockup (e.g. SaaS dashboard screen, published book, or podcast badge).
- **Lower Third:** Credential badges with clean line-art icons and exact website URL.

---

## 5. Typography, Contrast & Anti-Hallucination Rules

- **Primary Font Standard (Ubuntu):** All typography MUST default to the **Ubuntu** font family:
  - *Main Headline:* `Ubuntu Bold` at large scale for immediate authority and readability.
  - *Subtitles & Roles:* `Ubuntu Medium` or `Ubuntu Regular`.
  - *Pill Badges & CTAs:* `Ubuntu Bold` inside high-contrast rounded containers.
  - *URLs & Contact:* `Ubuntu Regular` or `Ubuntu Medium`.
- **Verbatim Text Enforcement:** All names, titles, and URLs must be explicitly quoted and verified character-by-character.
- **Contrast Ratios:** Text must meet WCAG AAA standards against the banner backdrop. Use darkened vignettes behind text layers or solid badge containers for URLs.

---

## 6. Directory & File Hierarchy (Strict Lowercase & Two-Digit Sequence)

All generated LinkedIn profile banner assets must follow strict lowercase naming and two-digit zero-padded sequence numbers. The AI must persist the design plan into `prompts/02-plan.md`, the generation prompt into `prompts/01-prompt.md`, and save all high-resolution banner images inside `linkedin-banners/`:

```
/ (repo root)
└── 02-projects/
    ├── 01-{project-name}/
    │   ├── readme.md (project overview, banner preview, and layout specs)
    │   ├── prompts/
    │   │   ├── 01-prompt.md (the exact prompt, user inputs, and AI parameters used)
    │   │   └── 02-plan.md (design plan, positioning strategy, safe zones, and checklist)
    │   └── linkedin-banners/
    │       ├── 01-banner-1584x396.png (Standard 1x profile banner)
    │       ├── 02-banner-3168x792.png (High-Resolution 2x Retina banner — Recommended)
    │       ├── 03-banner-4752x1188.png (Ultra High-Resolution 3x banner)
    │       └── 04-vector-overlay.svg (Crisp SVG vector typography, badges, and logo overlay)
    └── 02-{project-name}/
```

### Hierarchy Rules

1. **Root Projects Folder:** All projects live under `02-projects/`.
2. **Project Folder Naming:** `{sequence}-{project-name}` using two-digit zero-padding and kebab-case (e.g., `01-executive-brand`, `02-saas-advisor`).
3. **Prompt & Plan Preservation (Mandatory):** When formulating the profile banner design, every detail of the design plan (positioning statement, visual hierarchy, color palette, bottom gradient fade strategy, safe zone mapping, and verification checklist) MUST be saved directly to the file system at `prompts/02-plan.md`. The exact prompt given to the generation engine, user inputs, and model parameters MUST be saved in `prompts/01-prompt.md` so designs can be reproduced, audited, and iterated on.
4. **Asset Organization:** All banner raster images and vector overlays MUST be stored inside `linkedin-banners/` with zero-padded sequence prefixes (`01-`, `02-`, etc.).
5. **Relative Paths:** All links and image embeds in `readme.md` must use relative paths (e.g., `![Banner](linkedin-banners/02-banner-3168x792.png)`).

---

## 7. Asset Specifications

### A. Profile Banners (`linkedin-banners/`)

- `01-banner-1584x396.png`: Standard LinkedIn personal profile banner (4:1 aspect ratio, under 8MB).
- `02-banner-3168x792.png`: 2x Retina high-resolution banner (guarantees crystal-clear typography against compression).
- `03-banner-4752x1188.png`: 3x Ultra high-resolution banner for maximum fidelity.
- `04-vector-overlay.svg`: Crisp SVG vector layer containing all Ubuntu typography, badges, and icon elements for hybrid compositing.

### B. Generation Prompt & Plan Archive (`prompts/`)

- `01-prompt.md`: Contains the full generation prompt, model parameters (aspect ratio, style, negative prompts), and exact text strings used for the generation run.
- `02-plan.md`: Contains the full design plan, positioning strategy, color palette, typography specification, safe zone mapping, and verification checklist.

---

## 8. Ready-to-Use Templates

### Template 1: Midjourney / Flux Background & Portrait Generation Prompt

```text
Wide panoramic 4:1 banner background for LinkedIn profile, dimensions 3168x792. Deep navy blue and charcoal dark modern architectural interior with subtle amber bokeh lights. On the right third of the frame, a professional subject captured in a medium bust portrait with an 85mm portrait lens at f/1.8, razor-sharp focus on face and eyes, natural skin texture with visible pores, authentic expression, modern dark blazer. Soft diffused key light and vibrant electric cyan rim light tracing the hair and shoulders. The lower torso dissolves seamlessly into a soft dark gradient fade at the bottom edge. Clean minimalist negative space on the left and center for typography overlay, commercial photography, premium executive aesthetic, 8k resolution --ar 4:1 --style raw
```

### Template 2: Figma / SVG Typography Overlay Specification (Ubuntu Font Family)

```markdown
# Canvas: 3168 x 792 px (2x Retina, Ubuntu Font Family)

- Avatar Dead Zone: x: 0 to 760 px (Keep empty)
- Primary Headline: "{headline_text}" (Font: Ubuntu Bold, Color: #FFFFFF, Size: 64pt, x: 820px, y: 260px)
- Pillars: "{core_pillars}" (Font: Ubuntu Medium, Color: #06B6D4, Size: 32pt, x: 820px, y: 380px)
- Website URL: "{website_url}" (Font: Ubuntu Regular, Color: #E2E8F0, Size: 28pt, x: 820px, y: 500px)
- Call-to-Action Pill: "{cta_text}" (Font: Ubuntu Bold, Background: #06B6D4, Text: #0B0F19, Size: 24pt, x: 820px, y: 580px)
```
