# YouTube Thumbnail & Banner Design — Visual Identity & Typography Workflow

> **Prompt Version:** 2.0.0
> **Target Environment:** Lovable, Image Generation AI (Flux, Ideogram, Midjourney), Canva, Figma & Design AI Platforms
> **Synchronization:** Main Meta-Repo & Connected Workspaces

This prompt instructs design AI platforms on how to craft high-conversion, professional YouTube thumbnails and channel banners with strict text fidelity, zero hallucination, cinematic photographic shot direction, a 5-zone layout grid, and a high-contrast visual hierarchy.

---

## Strictly Avoid (Critical Negative Constraints)

> [!CAUTION]
> **ZERO TOLERANCE FOR TEXT HALLUCINATION, UNBLENDED EDGES, CLIP-ART & UNSEQUENCED FILES (AUTO-REJECT)**
> AI image models frequently produce waxy skin, harsh cutout edges, tacky clip-art stickers, and scrambled text. Accuracy, photographic realism, and spatial balance are the absolute top priorities.

1. **NO Text Hallucinations or Gibberish:** NEVER render garbled, pseudo-Latin, merged glyphs, or invented words. Every single character in titles, channel names, handles, and badges MUST match the user's exact spelling verbatim.
2. **NO Falsified Titles or Subtitles:** NEVER invent or embellish credentials, names, awards, statistics, testimonials, or quotes not explicitly provided or approved by the user.
3. **NO Unblended Subject Cutouts (MANDATORY BOTTOM & EDGE GRADIENT FADE):** NEVER paste a subject cutout with hard, abrupt bottom or side edges. The subject MUST blend naturally into the bottom edge of the frame using a soft linear gradient fade or dark feathered vignette so the torso/suit dissolves organically into the canvas.
4. **NO Tacky Clip-Art, Star Badges, or Fake Icons (TOTAL BAN):** NEVER add cheap clip-art shapes, cartoonish star stickers (e.g. `#01` star badges), arbitrary geometric badges, or fake icons. Do NOT invent or create avatars or icons unless explicitly requested by the user. Professional thumbnails rely on typography, authentic subject photography, and atmospheric lighting—NOT clip-art.
5. **NO Chaotic Overlapping Lines Across Subjects:** NEVER run background lines, graph strokes, or grid wires across the subject's face, neck, or body. Elements must stay behind the subject or maintain clean negative space.
6. **NO Unsequenced or Uppercase Filenames:** ALL generated files and images MUST use lowercase kebab-case preceded by a two-digit zero-padded sequence number (e.g. `01-thumbnail-1280x720.png`, `02-thumbnail-1920x1080.png`, `01-prompt.md`, `02-plan.md`). Unsequenced names like `thumbnail.png` are strictly BANNED.
7. **NO Low-Contrast Text:** NEVER place light text over bright, busy backgrounds or dark text over dark shadows without proper contrast treatment (e.g., dropshadows, dark vignettes, or solid badge backdrops).
8. **NO Empty Background Placeholders:** NEVER generate an empty background containing blank colored boxes or disconnected shapes while putting all text in an SVG with broken images. The composition must be fully realized, cohesive, and intentional.
9. **NO Uncontrolled Facial or Hand Distortions:** The human subject in the foreground must have anatomically correct eyes, glasses, hands, and fingers without AI melting or extra digits.
10. **NO HTML or Web Page Code:** Do NOT generate full HTML web pages or application templates. This is strictly a graphic design, thumbnail, and banner workflow.
11. **NO Low-Resolution Exports:** NEVER export or specify standard 72 DPI blurry images. Always specify Full HD (`1280x720 px` standard, `1920x1080 px` high-res for thumbnails, `2560x1440 px` for banners).
12. **NO Distorted Lighting or Neon Gradients:** Do not use purple, violet, cyan, fluorescent yellow, pastel colors, rainbow effects, or harsh saturated neon gradients unless explicitly requested.

---

## 1. Canvas Architecture & 5-Zone Layout Grid

All thumbnails are engineered on a `1280 × 720 px` canvas (16:9 aspect ratio) with strict spatial zone boundaries:

### Canvas Specifications

- **Dimensions:** `1280 × 720 pixels` (16:9 aspect ratio).
- **Edge Safe Area:** Keep all important faces, text, and focal graphics at least `55 pixels` away from every outside edge.
- **Central Core:** Primary content must remain inside the central 80% of the canvas.
- **Continuous Scene Blending:** Do NOT make sections look like separate rectangular floating cards. Blend them together as one continuous cinematic environment.

### 5-Zone Spatial Breakdown

```
+---------------------------------------------------------------------------------------------------+
|  [ ZONE A: LEFT CONTENT ]   [ ZONE B: PORTRAIT ]   [ ZONE C: NAME & CRED ]  [ ZONE D: ACHIEVEMENT]|
|  x: 25 to 405 px            x: 355 to 735 px       x: 655 to 965 px         x: 930 to 1245 px     |
|  - Handwritten Quote        - Seated Person        - Main Name (Stacked)    - "Top 1%" Badge      |
|  - Book Display / Table     - Natural Overlap      - Professional Titles    - Supporting Statement|
|  - Warm Studio Lights       - 54-58% Canvas Height - Values Line + Underline- Mic / YouTube Line  |
|                                                                                                   |
|                             [ ZONE E: LOWER INFORMATION STRIP ]                                   |
|                             x: 670 to 1280 px | y: 548 to 720 px                                  |
|                             - 4 Professional Role Groups behind lower body                        |
+---------------------------------------------------------------------------------------------------+
```

- **Zone A (Left Content Zone, x: 25 to 405 px, width ~380px):** Contains the handwritten quote, physical book display on a wood table, and warm studio lighting.
- **Zone B (Portrait Zone, x: 355 to 735 px, width ~380px):** Contains the seated featured subject. Overlaps the left and center zones slightly to connect the composition organically.
- **Zone C (Name & Credential Zone, x: 655 to 965 px, width ~310px):** Contains the stacked prominent name, professional credentials with vertical bar separators, and the values line with a hand-painted gold underline.
- **Zone D (Achievement Zone, x: 930 to 1245 px, width ~315px):** Contains the "Top 1%" achievement callout, supporting statement, YouTube identity line, and the podcast microphone entering from the right edge.
- **Zone E (Lower Information Strip, y: 548 to 720 px, x: 670 to 1280 px):** A dark translucent strip containing four professional role groups. It appears behind the subject's lower body where the two areas overlap.

---

## 2. Complete Color System & Distribution Formula

Apply a disciplined, luxury-authority color hierarchy across every element:

### Color Palette

- **Main Background Tones:**
  - Deep midnight navy: `#07111F`
  - Dark blue-black: `#0A1424`
  - Charcoal black: `#15171C`
  - Warm near-black brown: `#211612`
- **Primary Text Tones:**
  - Main white: `#F7F7F4`
  - Soft white: `#ECEBE6`
  - Secondary light gray: `#B8BBC2`
  - Muted gray: `#858B94`
- **Gold & Yellow Accents:**
  - Primary warm gold: `#F5A817`
  - Bright highlight gold: `#FFB51B`
  - Dark gold shadow: `#B96E08`
  - Hand-drawn underline gold: `#E6A51D`
- **Supporting Accents:**
  - Deep burgundy clothing: `#651E2B` (shadow: `#361019`)
  - YouTube red: `#FF0000` (strictly for the official icon)
  - Dark wood table: `#4B281D` (highlight: `#92553B`, shadow: `#27130E`)
  - Icon white: `#F5F4F0`
  - Divider gray: `#A5A5A5` at 60% opacity

### Color Distribution Formula

- **65%:** Deep navy, charcoal, and black (atmospheric background and shadow separation).
- **15%:** Warm brown and burgundy (studio warmth, clothing, and table textures).
- **12%:** White and light gray (high-contrast typography and readable headlines).
- **7%:** Gold accents (focal names, "1%", and hand-painted underlines).
- **<= 1%:** YouTube red (official logo badge only).

---

## 3. 3-Zone Cinematic Background Transition

The background must never be a flat solid color or an artificial neon gradient. Construct a seamless left-to-right cinematic transition:

1. **Left Background (Warm Studio & Bookshelves):**
   - Base color: dark charcoal brown `#211612`.
   - Add heavily blurred bookshelves and warm office ambient lights in soft amber `#C47A32`.
   - Low-contrast, creamy bokeh with zero sharp edges.
2. **Center Background (Contrast Separation Void):**
   - Base color: dark blue-black `#0A1424`.
   - The darkest area of the background sits directly behind the person's face and shoulders to create maximum separation and pop.
   - Subtle cool edge light in desaturated blue-gray `#6F8193`.
   - No bright objects directly behind the head.
3. **Right Background (Distant Architectural Silhouette):**
   - Base color: midnight navy `#07111F`.
   - Add a subtle, distant professional city or architectural skyline near the lower-right in muted gray-blue `#536170` at 20–30% opacity.
   - Very subtle warm horizon glow in dusty amber `#B47646` near the bottom.
   - Background remains darker than all foreground text and subject layers.

---

## 4. Input Capture & Verification

Before generating any design concepts or prompts, capture and validate the following inputs:

1. `channel_name` / `host_name`: Exact name of the creator or brand (e.g., "MD ALIM UL KARIM").
2. `primary_title`: Main headline or topic (e.g., "Better Ideas, Bigger Impact").
3. `subtitles_and_credentials`: Professional titles or roles separated by vertical bars (e.g., `Author | Marketer | Trainer | Consultant | Podcaster`).
4. `core_pillars`: 3 key thematic words (e.g., `Ideas | Strategy | Impact`).
5. `achievement_callout`: Main authority badge or metric (e.g., `"Top 1%"`).
6. `supporting_statement`: 1–2 line hook (e.g., `"Real Stories, Real People, Real Growth."`).
7. `youtube_handle`: Exact channel handle (e.g., `@alimulkarim`).
8. `youtube_url`: Full channel URL (e.g., `https://youtube.com/@alimulkarim`).
9. `website_url`: Creator or brand website URL (e.g., `https://alimulkarim.com`).
10. `email_address`: Contact or business inquiry email address to embed (e.g., `contact@alimulkarim.com`).
11. `qr_code`: Optional scannable QR code destination URL or asset to embed (optional).
12. `channel_icon_or_avatar`: Profile picture, logo mark, or avatar image (optional — ask user if available, or skip).
13. `person_photo_or_avatar`: High-resolution photo of the person, visual description, or subject cutout (optional — ask user if available, or skip).
14. `font_family`: Primary font family (Defaults to **Ubuntu** across all typography, titles, and overlays, with brush-script style for quotes).
15. `color_palette`: Primary brand colors (defaults to the 65/15/12/7/1 system above).

> [!IMPORTANT]
> Always ask the user if they have a **person's photo, subject cutout, or avatar image** to include, or if they prefer to **skip** it (for a graphics/typography-focused layout). Remember: do NOT generate unsolicited avatars, clip-art icons, or star badges. Also ask if they want to embed an **email address** or a **scannable QR code**. If `channel_name`, `primary_title`, `youtube_handle`, or `website_url` are missing or not provided, **STOP and ask the user** before proceeding.

---

## 5. Photographic, Lighting & Subject Specification

Professional YouTube thumbnails require deliberate cinematic direction rather than generic AI generation:

### A. Camera Shot & Framing

- **Lens & Optics:** 85mm prime portrait lens equivalent with a wide aperture (`f/1.8` to `f/2.8`), creating an authentic shallow depth of field where the subject remains razor-sharp while the background dissolves into soft, creamy bokeh.
- **Shot Distance:** Seated medium bust or torso shot (from waist/hands to head). The subject occupies approximately `54–58%` of canvas height from head to seated hands, with lower body continuing to the bottom edge.
- **Placement (Zone B):** Center of face at approximately `x = 555, y = 150`, portrait width `360–390 px`. Overlaps Zone A and Zone C slightly.
- **Pose & Expression:** Seated, facing forward, direct eye contact with the camera. Calm, confident, approachable expression with hands naturally clasped. Body turned no more than 5 degrees away.

### B. 3-Point Cinematic Lighting

- **Warm Key Light:** Upper-left softbox in warm amber `#F0B078`, creating natural facial modeling with gentle tonal transitions.
- **Soft Cool Fill:** Front-right ambient fill in muted blue-gray `#7D91A8`, preserving shadow detail on cheekbones and clothing.
- **Thin Warm Rim Light (Crucial):** Rim light in warm gold `#D99145` tracing the left shoulder, hair, and jawline for clean separation from dark backgrounds.
- **Shadow Quality:** Soft natural contact shadow beneath the chin and behind the body. No harsh cutout edges, no glowing halo outlines, and no white sticker borders.

### C. Texture & Environmental Realism

- **Organic Skin Texture:** Visible skin pores, natural micro-contrast, realistic skin tones, and authentic specular highlights. Total ban on plastic, waxy, or airbrushed AI skin.
- **Crisp Textile Weave:** High-fidelity fabric texture on clothing (e.g., deep burgundy overshirt `#651E2B`, white inner T-shirt `#F1F0EC`, dark navy trousers `#12213A`).
- **Glasses & Eyes:** Clear eyes behind lenses with controlled, minimal specular reflections.
- **Mandatory Bottom Gradient Fade:** Apply a soft vertical linear gradient fade across the bottom 20% of the subject cutout, allowing the lower body to dissolve organically into the canvas.

---

## 6. Component & Typography Specifications

### A. Left-Side Handwritten Quote (Zone A)

- **Exact Text:** e.g., `“Better Ideas, Bigger Impact”`
- **Position:** `x = 90, y = 55`, max width `285 px`, max height `125 px`. Keep at least `30 px` from books.
- **Font Style:** Bold handwritten brush-script (e.g., Caveat Brush, Kalam Bold, or Ubuntu Bold Script). Feels personally painted with visible brush variation.
- **Color:** Soft white `#F7F7F4` with subtle shadow (`#05070A` at 55% opacity, 3px offset).
- **Size:** Approximately `48–54 px`.
- **Underline:** Curved hand-painted brush stroke beneath the second line in warm gold `#E6A51D` (darker edge `#B96E08`), thickness `7–12 px`, rising slightly toward the right.

### B. Book Display on Wood Table (Zone A)

- **Position:** `x = 25 to 405 px`, `y = 225 to 475 px`, width `350–375 px`. Keep at least `25 px` from the subject's torso.
- **Four Upright Books:**
  1. Book 1: Mustard yellow cover `#F4BC24` with blue/white artwork.
  2. Book 2: Deep emerald green cover `#087158` with mint details.
  3. Book 3: Royal blue cover `#0758A4` with white details.
  4. Book 4: Warm off-white cover `#ECE8DF` with charcoal title and yellow accent.
- **Preserve Real Covers:** Preserve exact titles, typography, and Bengali or English lettering from supplied references. Do NOT invent fake book covers or garbled symbols.
- **Wood Table:** Dark wood `#4B281D`, highlight `#92553B`, shadow `#27130E` with realistic horizontal wood grain and soft warm reflections beneath the books.

### C. Main Name: Primary Title (Zone C)

- **Exact Stacked Wording:** e.g.,
  ```
  MD ALIM UL
  KARIM
  ```
- **Position:** Left edge `x = 660`, top `y = 85`, max width `300 px`. Keep `30–45 px` clear space from subject head, and `18–26 px` from Zone D. Both lines share the same left edge.
- **Font Family (Ubuntu Standard):** `Ubuntu Bold` or extra-bold geometric sans-serif (weight 800–900), tightly stacked with line height 0.82–0.9.
- **First Name Color:** Soft white `#F7F7F4`, lower shading `#D7D9DC`, dark shadow `#02060D` (65% opacity, 5px offset).
- **Last Name Color:** Warm golden yellow `#F5A817`, bright highlight `#FFB51B`, dark shadow `#B96E08`. Must be the strongest color accent in the center.

### D. Professional Credentials & Values Line (Zone C)

- **Credentials Block:** Directly beneath the name (`y = 305`, max width `320 px`).
  - Text: e.g., `Author | Marketer | Trainer` / `Consultant | Podcaster`
  - Font: `Ubuntu Medium` or clean brush script (size `26–31 px`, color `#ECEBE6`, separators `#C8C5BE`).
- **Values Line:** Below credentials (`y = 400`, max width `300 px`).
  - Text: `Ideas  |  Strategy  |  Impact`
  - Font: `Ubuntu Medium` (size `21–25 px`, color `#B8BBC2`, separators `#858B94`).
  - Underline: Thin hand-painted gold line `#E6A51D` (thickness `4–7 px`, rising 5–8 degrees toward the right).

### E. Achievement Area & Supporting Statement (Zone D)

- **"Top 1%" Callout:** Horizontal start `x = 945`, top `y = 165`.
  - "Top": Soft white `#F7F7F4`, size `78–90 px`.
  - "1%": Warm gold `#F5A817`, highlight `#FFB51B`, size `78–94 px`.
  - Both on the same line with precisely aligned baseline.
- **Supporting Statement:** Directly beneath `Top 1%` (`y = 305`).
  - Text: e.g., `Real Stories, Real People,` / `Real Growth.`
  - Font: `Ubuntu Medium` (size `22–27 px`, color `#B8BBC2`).
- **YouTube Channel Line:** Below supporting statement (`y = 385`).
  - Official YouTube play button icon (red `#FF0000`, white triangle, 30–34px wide).
  - Channel name or handle in `Ubuntu Bold` (`18–22 px`, color `#ECEBE6`).

### F. Podcast Microphone (Zone D)

- **Position:** Head centered near `x = 1220, y = 235`, entering from far-right edge at `25–35 degrees` downward-left angle.
- **Color & Finish:** Dark metal body `#171A20`, mesh highlights `#5B616B`, subtle burgundy reflection `#651E2B`, tiny warm highlight `#D99145`. No bright chrome or silver.
- **Depth:** Realistic shallow depth of field, slightly less sharp than the subject's face.

### G. Lower Professional Role Strip (Zone E)

- **Position:** `x = 675 to 1280 px`, `y = 550 to 720 px`. Positioned behind the subject where they overlap.
- **Background:** Base `#211612` at 88% opacity, upper edge `#4B281D` at 50%, lower edge `#0D0908` at 95%. Subtle 1–2px top separation line in `#7F5A45` at 35% opacity.
- **Four Role Groups (width ~140–150 px each, separated by 60% opacity dividers `#A5A5A5`):**
  1. Open-book line icon + `Bestselling Author`
  2. Graduation-cap line icon + `Hard-skill Trainer`
  3. Megaphone line icon + `Brand Marketing Professional`
  4. Briefcase line icon + `Business Consultant`
- **Icon & Label Typography:** Clean white line icons (`#F5F4F0`, 42–50px), labels in `Ubuntu Medium` (`18–23 px`, `#F5F4F0`, centered).

---

## 7. Spacing, Visual Separation & Required Reading Hierarchy

### Strict Spacing Rules

- Canvas edge to important text: minimum `55 pixels`
- Person's face to name: `30–45 pixels`
- Name to "Top 1%" area: `18–26 pixels`
- First name to last name gap: `0–8 pixels`
- Name to credential block: `18–25 pixels`
- Credentials to values line: `16–22 pixels`
- Values line to gold underline: `8–12 pixels`
- "Top 1%" to supporting statement: `18–26 pixels`
- Supporting statement to YouTube line: `15–20 pixels`
- Books to portrait: minimum `25 pixels`
- Quote to books: minimum `30 pixels`
- Microphone to achievement text: minimum `20 pixels`

### Required Reading Hierarchy (Mobile & Desktop)

1. The featured person's face & eye contact
2. The primary stacked name (white first name + gold last name)
3. The "Top 1%" achievement callout
4. The professional credentials & roles
5. The four colorful books on the table
6. The handwritten quote ("Better Ideas, Bigger Impact")
7. The lower professional role strip

---

## 8. Text Accuracy & Anti-Hallucination Protocol

When generating prompts for image engines (Flux, Midjourney v6, Ideogram) or compositing layers:

1. **Explicit Text Quoting:** Always specify text inside literal quotes in the generation prompt:
   - `with the exact text "MD ALIM UL" in bold white Ubuntu letters, and "KARIM" in bold warm-gold Ubuntu letters`
2. **Character Verification Gate:** Inspect the generated output. If even a single character is warped, merged, or misspelled, the image MUST be rejected or the text layer must be re-rendered as a clean vector overlay.
3. **Hybrid Compositing (Recommended):** For production-grade thumbnails:
   - Use AI to generate the photographic background, cinematic lighting, wood table, and subject portrait.
   - Render all typography, badges, book covers, and icons as crisp SVG vector layers over the background to guarantee 100% spelling precision.

---

## 9. Directory & File Hierarchy (Strict Lowercase & Two-Digit Sequence)

All generated YouTube thumbnail, banner, and prompt assets must follow strict lowercase naming and two-digit zero-padded sequence numbers. The AI must persist the design plan into `prompts/02-plan.md`, the generation prompt into `prompts/01-prompt.md`, and save all high-resolution images inside `youtube-thumbnails/`:

```
/ (repo root)
└── 02-projects/
    ├── 01-{project-name}/
    │   ├── readme.md (project overview, visual previews, and layout specs)
    │   ├── prompts/
    │   │   ├── 01-prompt.md (the exact prompt, user inputs, and AI parameters used)
    │   │   └── 02-plan.md (design plan, composition strategy, safe zones, and checklist)
    │   └── youtube-thumbnails/
    │       ├── 01-thumbnail-1280x720.png (Standard 16:9 YouTube video thumbnail)
    │       ├── 02-thumbnail-1920x1080.png (Full HD 1080p high-resolution thumbnail)
    │       ├── 03-banner-2560x1440.png (Full YouTube channel banner / TV master)
    │       ├── 04-banner-safe-zone-1546x423.png (Desktop & mobile safe crop banner)
    │       └── 05-vector-overlay.svg (Crisp SVG vector typography, badges, and logo overlay)
    └── 02-{project-name}/
```

### Hierarchy Rules

1. **Root Projects Folder:** All projects live under `02-projects/`.
2. **Project Folder Naming:** `{sequence}-{project-name}` using two-digit zero-padding and kebab-case (e.g., `01-tech-podcast`, `02-coding-insights`).
3. **Prompt & Plan Preservation (Mandatory):** When formulating the thumbnail/banner design, every detail of the design plan (visual composition, color grading, typography pairing, bottom gradient fade strategy, safe zone mapping, and anti-pattern checklist) MUST be saved directly to the file system at `prompts/02-plan.md`. The exact prompt given to the generation engine, user inputs, and model parameters MUST be saved in `prompts/01-prompt.md` so designs can be reproduced, audited, and iterated on.
4. **Asset Organization:** All thumbnail and banner raster images MUST be stored inside `youtube-thumbnails/` with zero-padded sequence prefixes (`01-`, `02-`, etc.).
5. **Relative Paths:** All links and image embeds in `readme.md` must use relative paths (e.g., `![Thumbnail](youtube-thumbnails/01-thumbnail-1280x720.png)`).

---

## 10. Asset Specifications

### A. YouTube Thumbnails (`youtube-thumbnails/`)

- `01-thumbnail-1280x720.png`: Standard YouTube video thumbnail (16:9 aspect ratio, under 2MB).
- `02-thumbnail-1920x1080.png`: Full HD high-resolution thumbnail for pristine visual quality on high-DPI displays.
- `03-banner-2560x1440.png`: Full YouTube channel banner (TV master dimension).
- `04-banner-safe-zone-1546x423.png`: Centered safe zone banner crop ensuring logos and text are fully visible on desktop and mobile.
- `05-vector-overlay.svg`: Crisp SVG vector layer containing all Ubuntu typography, badges, URLs, and icons for hybrid compositing.

### B. Generation Prompt & Plan Archive (`prompts/`)

- `01-prompt.md`: Contains the full generation prompt, model parameters (aspect ratio, style, negative prompts), and exact text strings used for the generation run.
- `02-plan.md`: Contains the full design plan, composition strategy, color palette, typography specification, safe zone mapping, and anti-pattern quality verification checklist.

---

## 11. Master Production Sample Prompt & Templates

### Template 1: Master Production Personal-Brand Thumbnail Prompt (Full Unabridged Reference)

Use this complete reference prompt when generating a cinematic, authority-grade YouTube thumbnail. AI models must follow this exact structure, coordinate precision, color discipline, and component breakdown:

```text
Create a highly polished, cinematic personal-brand YouTube thumbnail at exactly 1280 × 720 pixels, 16:9 aspect ratio, based on the supplied reference image.

IMPORTANT: Recreate the reference image’s overall composition, visual hierarchy, spacing, lighting, and professional atmosphere, but change the featured name to:

MD ALIM UL
KARIM

Use the supplied portrait, books, logo, and other real assets as references whenever available. Preserve the person’s actual identity, face, complexion, hairstyle, beard, glasses, clothing, and body shape. Do not invent a different person.

==================================================
1. OVERALL CANVAS AND LAYOUT GRID
==================================================

Canvas:
- Width: 1280 pixels
- Height: 720 pixels
- Aspect ratio: 16:9
- Safe area: Keep all important faces and text at least 55 pixels away from every outside edge.
- Primary content must remain inside the central 80% of the canvas.

Divide the canvas into these visual zones:

A. LEFT CONTENT ZONE:
- Horizontal position: x = 25 to 405 pixels
- Contains the handwritten quote, books, and supporting visual details.
- Approximate width: 380 pixels.

B. PORTRAIT ZONE:
- Horizontal position: x = 355 to 735 pixels
- Contains the seated featured person.
- Approximate width: 380 pixels.
- The person overlaps the left and center zones slightly to connect the composition.

C. NAME AND CREDENTIAL ZONE:
- Horizontal position: x = 655 to 965 pixels
- Contains “MD ALIM UL KARIM,” professional titles, and the short values line.
- Approximate width: 310 pixels.

D. ACHIEVEMENT ZONE:
- Horizontal position: x = 930 to 1245 pixels
- Contains “Top 1%,” the supporting statement, the YouTube identity, and the microphone.
- Approximate width: 315 pixels.

E. LOWER INFORMATION STRIP:
- Vertical position: y = 548 to 720 pixels
- Runs from approximately x = 670 to 1280 pixels.
- Contains four professional role groups.
- It must appear behind the lower body where the two areas overlap.

Do not make the sections look like separate rectangular cards. Blend them together as one continuous cinematic scene.

==================================================
2. COMPLETE COLOR SYSTEM
==================================================

Use the following exact color hierarchy throughout the design:

MAIN BACKGROUND COLORS:
- Deep midnight navy: #07111F
- Dark blue-black: #0A1424
- Charcoal black: #15171C
- Warm near-black brown: #211612

PRIMARY TEXT COLORS:
- Main white: #F7F7F4
- Soft white: #ECEBE6
- Secondary light gray: #B8BBC2
- Muted gray: #858B94

GOLD AND YELLOW ACCENTS:
- Primary warm gold: #F5A817
- Bright highlight gold: #FFB51B
- Dark gold shadow: #B96E08
- Hand-drawn underline gold: #E6A51D

ADDITIONAL ACCENTS:
- Deep burgundy clothing accent: #651E2B
- Burgundy shadow: #361019
- YouTube red: #FF0000
- Dark wood: #4B281D
- Warm wood highlight: #92553B
- Icon white: #F5F4F0
- Divider gray: #A5A5A5 at approximately 60% opacity

COLOR DISTRIBUTION:
- Approximately 65% deep navy, charcoal, and black
- Approximately 15% warm brown and burgundy
- Approximately 12% white and light gray
- Approximately 7% gold
- No more than 1% bright red

Do not use purple, violet, cyan, fluorescent yellow, pastel colors, rainbow effects, or colorful gradients.

==================================================
3. BACKGROUND DESIGN
==================================================

Create a dark, cinematic professional studio environment.

Use a left-to-right background color transition:

LEFT BACKGROUND:
- Base color: dark charcoal brown #211612
- Add blurred shelves and warm office lights.
- Use subtle amber light sources in #C47A32.
- Keep the lights heavily blurred and low-contrast.

CENTER BACKGROUND:
- Base color: dark blue-black #0A1424
- Place the darkest area directly behind the person’s face and shoulders to create strong separation.
- Add a subtle cool edge light in desaturated blue-gray #6F8193.
- Do not place bright objects directly behind the face.

RIGHT BACKGROUND:
- Base color: midnight navy #07111F
- Add a subtle, distant professional city or architectural silhouette near the lower-right.
- Color the architecture in muted gray-blue #536170.
- Keep it at only 20–30% visual opacity.
- Add a very subtle warm horizon glow in dusty amber #B47646 near the bottom.
- The background must remain darker than every important foreground element.

Apply a soft photographic depth-of-field effect:
- Background shelves: strongly blurred
- Distant architecture: moderately blurred
- Person, books, name, and microphone: sharp
- No visible artificial bokeh circles
- No decorative glowing orbs
- No radial neon glow

==================================================
4. FEATURED PERSON
==================================================

Place the featured person slightly left of center.

POSITION:
- Center of face: approximately x = 555, y = 150
- Top of head: approximately y = 25
- Bottom of visible body: approximately y = 720
- Portrait width: approximately 360–390 pixels
- The person should occupy approximately 54–58% of the canvas height when measured from head to seated hands, with the lower legs continuing toward the bottom edge.

POSE:
- Seated and facing forward
- Direct eye contact with the camera
- Calm, confident, approachable expression
- Hands naturally clasped around the lower center
- Shoulders relaxed
- Body turned no more than 5 degrees away from the camera
- Do not make the pose theatrical or aggressive

CLOTHING:
- Deep burgundy overshirt: main color #651E2B
- Burgundy shadows: #361019
- White inner T-shirt: #F1F0EC
- Dark navy trousers: #12213A
- Black glasses: #171717
- Dark watch with restrained metallic highlights

LIGHTING:
- Warm key light from upper-left: #F0B078
- Soft cool fill from front-right: #7D91A8
- Thin warm rim light around left shoulder: #D99145
- Soft shadow beneath the chin
- Natural skin texture
- Controlled highlights on forehead and cheeks
- Clear eyes behind the glasses
- Minimal reflection on the lenses
- No artificial beauty-filter skin
- No glowing outline around the body

CUTOUT QUALITY:
- Hair, beard, shoulders, and glasses must have clean natural edges.
- Use a subtle realistic contact shadow behind the body.
- Do not create a sticker-like white border.
- Do not blur the face.

==================================================
5. LEFT-SIDE HANDWRITTEN QUOTE
==================================================

Place the quote in the upper-left area.

EXACT WORDING:

“Better Ideas
Bigger Impact”

POSITION:
- Quote block begins approximately at x = 90, y = 55
- Maximum width: 285 pixels
- Maximum height: 125 pixels
- Keep at least 30 pixels between the quote and the top of the books.
- Align the two lines slightly diagonally rather than perfectly mechanical.
- Do not overlap the person’s head.

FONT:
- Use a bold handwritten brush-script style similar to “Caveat Brush,” “Kalam Bold,” or a premium hand-painted marker font.
- The letters should feel personally written with visible brush variation.
- Do not use a formal calligraphy font.
- Do not use thin cursive writing.
- Do not use an elegant wedding-style script.

FONT COLOR:
- Main quote color: soft white #F7F7F4
- Very subtle shadow: #05070A at 55% opacity
- Shadow offset: approximately 3 pixels downward and 3 pixels right
- No outline or glowing effect

FONT SIZE:
- “Better Ideas”: approximately 48–54 pixels
- “Bigger Impact”: approximately 48–54 pixels
- Line gap: approximately 3–7 pixels
- Letter spacing: 0
- Slight natural variation in baseline is acceptable

QUOTATION MARKS:
- Add one opening quotation mark before “Better.”
- Add one closing quotation mark after “Impact.”
- Use the same soft white #F7F7F4.
- Keep quotation marks smaller than the main lettering.

UNDERLINE:
- Add one curved hand-painted brush stroke beneath “Bigger Impact.”
- Color: warm gold #E6A51D
- Darker edge: #B96E08
- Width: approximately 210–240 pixels
- Thickness: 7–12 pixels
- Place it 5–10 pixels beneath the text.
- Make the stroke rise slightly toward the right.
- It must look painted, not digitally glowing.

==================================================
6. BOOK DISPLAY
==================================================

Place four upright books on a realistic wooden table in the lower-left area.

BOOK AREA:
- Horizontal position: x = 25 to 405 pixels
- Vertical position: y = 225 to 475 pixels
- Books should occupy approximately 350–375 pixels in total width.
- Maintain 5–12 pixels of natural spacing between books.
- The books may slightly overlap one another because of perspective.
- Keep at least 25 pixels between the books and the person’s torso.

BOOK COLORS:
1. First book:
   - Cover color: mustard yellow #F4BC24
   - Shadow color: #B67B0C
   - Optional megaphone artwork in blue #164A8A and white #F7F7F4

2. Second book:
   - Cover color: deep emerald green #087158
   - Shadow color: #034838
   - Supporting details in muted mint #A4D0B5

3. Third book:
   - Cover color: professional royal blue #0758A4
   - Shadow color: #06396D
   - Supporting details in white #F7F7F4

4. Fourth book:
   - Cover color: warm off-white #ECE8DF
   - Shadow color: #C8C1B3
   - Main title color: charcoal #242321
   - Small accent color: warm yellow #EBAE24

Use the original book covers exactly if supplied. Preserve their real Bengali or English text, typography, and imagery. Do not invent book names. Do not replace Bengali lettering with meaningless symbols.

TABLE:
- Main dark wood color: #4B281D
- Highlight color: #92553B
- Shadow color: #27130E
- Add realistic horizontal wood grain.
- Use a warm reflection directly beneath the books.
- Keep the surface glossy enough to show a soft reflection, but not mirror-like.

==================================================
7. MAIN NAME: “MD ALIM UL KARIM”
==================================================

Place the name immediately to the right of the person’s head and upper torso.

EXACT WORDING AND LINE BREAK:

MD ALIM UL
KARIM

Use uppercase lettering for stronger thumbnail readability.

POSITION:
- Left edge of name block: approximately x = 660
- Top edge: approximately y = 85
- Maximum width: approximately 300 pixels
- The name block must not overlap the person’s face, hair, glasses, or beard.
- Keep approximately 30–45 pixels of clear space between the person’s head and the first letter.
- Keep approximately 18–26 pixels between the name block and the “Top 1%” area.
- Both name lines must share the same left edge.

FONT:
- Use an extra-bold geometric sans-serif similar to “Ubuntu Bold,” “Archivo Black,” “Montserrat ExtraBold,” “League Spartan Black,” or “Anton.”
- Choose a wide, substantial typeface rather than a narrow condensed font.
- Weight: 800–900
- Letter spacing: 0
- No italics
- No serif styling
- No handwritten styling
- No outline

FIRST NAME COLOR — MD ALIM UL:
- Main fill: soft white #F7F7F4
- Subtle lower shading: #D7D9DC
- Dark shadow: #02060D at 65% opacity
- Shadow offset: 5 pixels down and 5 pixels right
- Very subtle depth only; do not make it look metallic.

LAST NAME COLOR — KARIM:
- Main fill: warm golden yellow #F5A817
- Bright top highlight: #FFB51B
- Dark lower shadow: #B96E08
- External dark shadow: #02060D at 70% opacity
- Shadow offset: 5 pixels down and 5 pixels right
- “KARIM” must be the strongest color accent in the central section.
- Do not color “KARIM” orange-red.
- Do not use a gold metallic texture.
- Do not apply an outline.

FONT SIZE AND SPACING:
- “MD ALIM UL”: approximately 82–94 pixels
- “KARIM”: approximately 82–94 pixels
- Line height: approximately 0.82–0.9 of the font size
- Gap between the visual bottoms and tops of the two lines: approximately 0–8 pixels
- The two lines should appear tightly stacked, as in a professional poster.
- “KARIM” should begin directly beneath “MD ALIM UL,” not indented.

READING HIERARCHY:
- The person’s face is the first visual focus.
- “MD ALIM UL KARIM” is the second focus.
- “Top 1%” is the third focus.
- The white first name and gold last name must be instantly distinguishable.

==================================================
8. PROFESSIONAL CREDENTIALS
==================================================

Place the credentials directly beneath “KARIM.”

EXACT WORDING:

Author | Marketer | Trainer
Consultant | Podcaster

POSITION:
- Top of credential block: approximately y = 305
- Center the block beneath the name.
- Maximum width: approximately 320 pixels.
- Keep 18–25 pixels between the bottom of “KARIM” and the first credential line.
- Keep 4–8 pixels between the two credential lines.
- Keep at least 18 pixels between this block and the values line beneath it.

FONT:
- Use the same handwritten brush-script family used for the left quote (or clean Ubuntu Medium).
- Recommended style: “Kalam Bold,” “Caveat Brush,” or an equivalent legible brush script.
- The lettering should feel personal but remain clear.
- Do not use tiny, thin cursive writing.

COLOR:
- Main color: soft white #ECEBE6
- Separators: warm light gray #C8C5BE
- Shadow: #02060D at 60% opacity
- No gold coloring in this section.
- No outline or glow.

FONT SIZE:
- Approximately 26–31 pixels
- Line height: approximately 32–36 pixels
- Letter spacing: 0

SEPARATORS:
- Use simple vertical bars: |
- Keep approximately one normal character space on each side.
- Do not replace separators with dots, commas, or slashes.

==================================================
9. VALUES LINE
==================================================

Place this line beneath the professional credentials:

Ideas  |  Strategy  |  Impact

POSITION:
- Approximate top position: y = 400
- Align it horizontally with the center of the name block.
- Keep 16–22 pixels between this line and the credentials.
- Maximum width: approximately 300 pixels.

FONT:
- Clean medium-weight sans-serif similar to “Ubuntu Medium,” “Montserrat SemiBold,” or “Archivo SemiBold.”
- Font size: approximately 21–25 pixels
- Weight: 600
- Letter spacing: 0

COLOR:
- Words: muted light gray #B8BBC2
- Separators: #858B94
- Do not color any individual word gold.

UNDERLINE:
- Add a thin hand-painted gold line below this phrase.
- Color: #E6A51D
- Dark edge: #B96E08
- Width: approximately 185–220 pixels
- Thickness: 4–7 pixels
- Vertical gap beneath text: 8–12 pixels
- Angle: rise approximately 5–8 degrees toward the right

==================================================
10. “TOP 1%” ACHIEVEMENT AREA
==================================================

Place the main achievement on the right side.

EXACT WORDING:

Top 1%

POSITION:
- Horizontal starting point: approximately x = 945
- Top position: approximately y = 165
- Keep at least 55 pixels from the right canvas edge.
- Keep at least 20 pixels between this phrase and the microphone.
- Do not overlap the name block.

FONT:
- Use the same heavy geometric sans-serif used for “MD ALIM UL KARIM” (e.g. “Ubuntu Bold”).
- Weight: 800–900
- Letter spacing: 0
- No italics
- No outline

“TOP” COLOR:
- Main fill: #F7F7F4
- Lower soft shading: #D7D9DC
- Dark shadow: #02060D at 70% opacity
- Shadow offset: 5 pixels down and 5 pixels right

“1%” COLOR:
- Main fill: #F5A817
- Bright upper highlight: #FFB51B
- Lower shadow: #B96E08
- External dark shadow: #02060D at 70% opacity
- No metallic texture
- No orange-red coloring

FONT SIZE:
- “Top”: approximately 78–90 pixels
- “1%”: approximately 78–94 pixels
- Keep “Top” and “1%” on the same line.
- Gap between “Top” and “1%”: approximately 16–24 pixels.
- The baseline must align precisely.

==================================================
11. SUPPORTING ACHIEVEMENT STATEMENT
==================================================

Place this directly beneath “Top 1%”:

Real Stories, Real People,
Real Growth.

POSITION:
- Top: approximately y = 305
- Left-align it with the beginning of “Top.”
- Keep 18–26 pixels below the large title.
- Maximum width: approximately 265 pixels.
- Keep 4–7 pixels between the two lines.

FONT:
- Clean sans-serif similar to “Ubuntu Medium,” “Montserrat Medium,” or “Archivo Medium.”
- Font size: approximately 22–27 pixels
- Weight: 500–600
- Letter spacing: 0
- Line height: approximately 28–32 pixels

COLOR:
- Main color: #B8BBC2
- Important punctuation stays the same color.
- Soft shadow: #02060D at 55% opacity
- No gold words in this section.
- No italics.

==================================================
12. YOUTUBE CHANNEL LINE
==================================================

Place a small YouTube identification line below the supporting statement.

POSITION:
- Top: approximately y = 385
- Left-align it with the supporting statement.
- Keep 15–20 pixels below “Real Growth.”
- Keep the entire line within approximately 265 pixels.

ICON:
- Use the official simple YouTube play-button shape.
- Red background: #FF0000
- White play triangle: #FFFFFF
- Icon size: approximately 30–34 pixels wide
- Do not distort the icon.

CHANNEL TEXT:
- Place the confirmed channel handle or channel name to the right of the icon.
- Gap between icon and text: approximately 9–12 pixels.
- Font: clean sans-serif, weight 600–700 (e.g. “Ubuntu Bold”).
- Font size: approximately 18–22 pixels.
- Color: #ECEBE6.
- Do not invent a channel name or handle.
- If no confirmed channel name is supplied, omit this complete line.

==================================================
13. PODCAST MICROPHONE
==================================================

Place a professional podcast microphone entering from the far-right edge.

POSITION:
- Microphone head centered approximately near x = 1220, y = 235
- Angle the microphone downward-left by approximately 25–35 degrees.
- Allow part of the microphone and support arm to crop naturally outside the right edge.
- Keep at least 20 pixels between the microphone and “1%.”
- Do not cover any important text.

COLORS:
- Main metal body: #171A20
- Dark shadow: #050608
- Mesh highlights: #5B616B
- Very subtle burgundy reflection: #651E2B
- Tiny warm highlight: #D99145
- No bright chrome or silver appearance.

DEPTH:
- Keep the microphone sharp enough to identify.
- Apply slightly less sharpness than the person’s face.
- Use realistic shallow depth of field.
- Do not make the microphone larger than the person’s head.

==================================================
14. LOWER PROFESSIONAL ROLE STRIP
==================================================

Create one continuous dark translucent information strip across the bottom-right section.

POSITION:
- Begin around x = 675
- Extend to x = 1280
- Begin around y = 550
- Extend to the bottom edge at y = 720
- The strip should remain behind the person where they overlap.

BACKGROUND COLOR:
- Base: #211612 at approximately 88% opacity
- Upper edge: subtle #4B281D at approximately 50% opacity
- Lower edge: #0D0908 at approximately 95% opacity
- Use a very subtle vertical tonal transition.
- Do not use a bright border.
- Do not make it appear as a floating rounded card.
- Corner radius: 0

TOP SEPARATION:
- Add a very subtle 1–2 pixel line in #7F5A45 at 35% opacity.
- Keep this line understated.

Create four evenly spaced role groups inside the strip.

GROUP WIDTH:
- Approximately 140–150 pixels each.
- Vertical dividers should sit midway between groups.
- Keep approximately 18–25 pixels of interior horizontal space around each icon and label.

DIVIDERS:
- Color: #A5A5A5
- Opacity: approximately 60%
- Width: 1–2 pixels
- Height: approximately 95 pixels
- Start around y = 575
- End around y = 670

GROUP 1:
- White open-book icon
- Label:
  Bestselling
  Author

GROUP 2:
- White graduation-cap icon
- Label:
  Hard-skill
  Trainer

GROUP 3:
- White megaphone icon
- Label:
  Brand Marketing
  Professional

GROUP 4:
- White briefcase icon
- Label:
  Business
  Consultant

ICON STYLE:
- Simple professional line icons.
- Icon color: #F5F4F0
- Stroke width should be visually consistent across all icons.
- Icon size: approximately 42–50 pixels.
- Keep icons centered over their labels.
- No colored icon backgrounds.
- No circular icon containers.

LABEL TYPOGRAPHY:
- Font: clean sans-serif similar to “Ubuntu Medium,” “Archivo SemiBold,” or “Montserrat SemiBold.”
- Font size: approximately 18–23 pixels.
- Weight: 600.
- Line height: approximately 22–26 pixels.
- Text color: #F5F4F0.
- Text alignment: centered.
- Gap between icon and label: approximately 10–14 pixels.
- Keep every label inside its assigned group.

==================================================
15. OPTIONAL LOWER-LEFT SUPPORTING PHOTOGRAPHS
==================================================

If the reference photographs are available, place two subtle monochrome supporting images behind the person near the bottom-left:

- One image showing the person speaking into a microphone.
- One image showing the person presenting to an audience.

POSITION:
- Vertical range: y = 505 to 720
- Keep them behind the main portrait.
- Do not cover the books.
- Do not compete with the main face.

COLOR TREATMENT:
- Convert to near-monochrome.
- Main tone: #A7A9AD
- Shadow tone: #171A20
- Opacity: approximately 45–60%
- Add a subtle cool blue-gray tint #536170.
- Keep contrast lower than the main portrait.
- Use a thin muted gray divider between the photographs if necessary.

Omit this section entirely if real supporting photographs are unavailable. Do not invent fake event photographs.

==================================================
16. SPACING AND VISUAL SEPARATION
==================================================

Maintain these approximate spacing rules:

- Canvas edge to important text: minimum 55 pixels
- Person’s face to name: 30–45 pixels
- Name to “Top 1%” area: 18–26 pixels
- “MD ALIM UL” to “KARIM”: 0–8 pixels
- Name to credential block: 18–25 pixels
- Credential line 1 to line 2: 4–8 pixels
- Credentials to values line: 16–22 pixels
- Values line to gold underline: 8–12 pixels
- “Top 1%” to supporting statement: 18–26 pixels
- Supporting statement to YouTube line: 15–20 pixels
- Icons to role labels: 10–14 pixels
- Books to portrait: minimum 25 pixels
- Quote to books: minimum 30 pixels
- Microphone to achievement text: minimum 20 pixels

No text should touch another section. No element should look accidentally squeezed into the remaining space.

==================================================
17. REQUIRED READING ORDER
==================================================

At small mobile-thumbnail size, the viewer must notice elements in this order:

1. The featured person’s face
2. “MD ALIM UL KARIM”
3. “Top 1%”
4. “Author | Marketer | Trainer / Consultant | Podcaster”
5. The four colorful books
6. “Better Ideas, Bigger Impact”
7. The lower professional role strip

Use scale, contrast, lighting, and whitespace to enforce this exact hierarchy.

==================================================
18. TEXT ACCURACY REQUIREMENTS
==================================================

Render only these approved phrases:

“Better Ideas
Bigger Impact”

“MD ALIM UL KARIM”

“Author | Marketer | Trainer
Consultant | Podcaster”

“Ideas | Strategy | Impact”

“Top 1%”

“Real Stories, Real People,
Real Growth.”

“Bestselling
Author”

“Hard-skill
Trainer”

“Brand Marketing
Professional”

“Business
Consultant”

Every word must be spelled exactly as written.

Do not:
- Write “ALIM UL KARIM” (without MD) or misspelling
- Add another person’s name
- Invent awards, statistics, testimonials, prices, or claims
- Invent a YouTube handle
- Add random Bengali or English text
- Produce garbled lettering
- Repeat any phrase
- Merge words
- Replace vertical separators with unrelated symbols

If the image model cannot reproduce long text accurately, leave those specific areas visually clean and reserved for later typography placement rather than generating misspelled text.

==================================================
19. FINAL PHOTOGRAPHIC AND DESIGN QUALITY
==================================================

The final result must feel like:
- A premium business and marketing personality thumbnail
- A cinematic podcast or documentary cover
- Confident, credible, intelligent, and approachable
- Professionally photographed and manually designed
- Highly readable on both desktop and mobile

Use:
- Realistic photography
- Natural skin
- Crisp facial detail
- Controlled shadows
- High text contrast
- Clean subject separation
- Restrained cinematic color grading
- Professional spacing
- Consistent typography
- Realistic material textures

Avoid:
- Generic AI-generated face
- Changed identity
- Different beard or hairstyle
- Distorted glasses
- Extra fingers
- Malformed hands
- Plastic skin
- Excessive sharpening
- Sticker-like portrait outline
- Strong glowing effects
- Purple or neon gradients
- Bright blue lighting
- Random decorative shapes
- Floating cards
- Excessive rounded corners
- Metallic gold text
- Thin unreadable script
- Crowded text
- Tiny captions
- Fake books
- Garbled Bengali writing
- Misspelled English writing
- Text over the face
- Microphone covering “Top 1%”
- Important content outside the safe area
- Any additional wording not explicitly approved
```

---

### Adaptation Guide: Modifying the Master Prompt for Any Subject or Client

When adapting this master prompt for different creators, topics, or brand archetypes, modify the 19 sections as follows:

| Section | What to Keep | What to Modify / Customize |
|---|---|---|
| **1. Canvas & Grid** | 1280x720 px, 16:9, 55px safe margin, 5-zone structure | Adjust zone widths if the subject is on the right instead of the left. |
| **2. Color System** | 65/15/12/7/1 distribution formula | Swap burgundy/warm gold for the client's primary/accent brand colors (e.g., deep navy + electric cyan). |
| **3. Background** | 3-zone left-to-right contrast transition | Match the background theme to the industry (e.g., studio shelves for authors, server racks/mesh for tech, trading floor for finance). |
| **4. Featured Person** | 85mm `f/1.8` lens, 3-point lighting, bottom gradient fade, natural skin | Subject photo reference, clothing, glasses, facial hair, and pose. |
| **5. Handwritten Quote** | Brush-script font style, gold hand-painted underline | Replace with the creator's signature hook, manifesto, or episode title. |
| **6. Props / Display** | Realistic table with horizontal grain, soft reflection | Replace books with client-relevant props (e.g., laptop with code, award statuette, podcast audio console, product packaging). |
| **7. Main Name** | Two-line stacked layout, white first line, gold second line | The creator's actual full name (`{first_name}` + `{last_name}`). |
| **8. Credentials** | Vertical bar separators `\|`, soft white font | Creator's professional titles (e.g., `Software Architect \| Founder \| Speaker`). |
| **9. Values Line** | 3 thematic pillars + thin hand-painted underline | Creator's 3 core pillars (e.g., `Code \| Scale \| Security`). |
| **10. Achievement** | Large stacked callout (white + gold accent) | Custom metric (e.g., `"Top 1%"`, `"$10M ARR"`, `"100k Subs"`, `"EPISODE 42"`). |
| **11. Supporting Statement** | 1–2 line medium-weight hook | Custom tagline or episode premise. |
| **12. YouTube Line** | Official YouTube icon + handle | Creator's verified `@handle`. |
| **13. Accent Prop** | Shallow depth of field, entering from edge | Podcast mic, camera lens, stylized stylus, or omit if not applicable. |
| **14. Role Strip** | Dark translucent strip, 4 groups, white line icons, dividers | 4 core expertise domains matching the creator's portfolio. |
| **15. Supporting Photos** | Low-contrast monochrome photos behind portrait | Real stage/speaking photos if available, otherwise omit entirely. |
| **16. Spacing** | Strict spacing minimums (canvas edge 55px, etc.) | Preserve exact pixel spacing rules. |
| **17. Reading Order** | 1 to 7 visual priority sequence | Maintain exact visual hierarchy. |
| **18. Text Accuracy** | Verbatim text quoting, anti-hallucination gate | List only confirmed text strings for the project. |
| **19. Quality** | Negative constraints, ban on waxy skin, ban on clip-art | Preserve all negative constraints verbatim. |

---

### Template 2: Flux / Midjourney Photorealistic Scene Generation Prompt

```text
Cinematic 16:9 YouTube thumbnail scene. Seated professional subject captured in a medium bust portrait with an 85mm portrait lens at f/1.8, razor-sharp focus on face and eyes, natural skin texture with visible pores, authentic confident expression, wearing a modern deep burgundy overshirt over a white inner shirt. Dramatic 3-point studio lighting with a warm soft key light from the upper-left and vibrant gold rim light tracing the hair and shoulder. The lower torso dissolves seamlessly into a soft dark gradient fade at the bottom edge. On the left, four colorful published books on a polished dark mahogany table with soft reflections. Dark atmospheric studio background with deep charcoal tones on the left and midnight navy on the right with subtle distant architectural bokeh. Ample clean negative space in the center and right for typography overlay. High commercial photography quality, 8k resolution --ar 16:9 --style raw
```

### Template 3: Figma / SVG Typography & Vector Overlay Specification (Ubuntu Font Family)

```markdown
# Canvas: 1280 x 720 px (16:9, Ubuntu Font Family)

## Zone A: Quote (x: 90, y: 55)

- Line 1: "“Better Ideas" (Font: Caveat Brush / Ubuntu Bold Script, Color: #F7F7F4, Size: 50pt)
- Line 2: "Bigger Impact”" (Font: Caveat Brush / Ubuntu Bold Script, Color: #F7F7F4, Size: 50pt)
- Underline: Hand-painted curve (Color: #E6A51D, Width: 220px, Height: 8px, y: 165px)

## Zone C: Main Name (x: 660, y: 85)

- First Line: "MD ALIM UL" (Font: Ubuntu Bold, Color: #F7F7F4, Size: 86pt, Shadow: #02060D 65% 5px down-right)
- Second Line: "KARIM" (Font: Ubuntu Bold, Color: #F5A817, Size: 86pt, Shadow: #02060D 70% 5px down-right)
- Credentials (y: 305): "Author | Marketer | Trainer | Consultant | Podcaster" (Font: Ubuntu Medium, Color: #ECEBE6, Size: 28pt)
- Values (y: 400): "Ideas  |  Strategy  |  Impact" (Font: Ubuntu Medium, Color: #B8BBC2, Size: 23pt)
- Values Underline: Hand-painted gold stroke (Color: #E6A51D, Width: 200px, Height: 5px, y: 435px)

## Zone D: Achievement (x: 945, y: 165)

- Achievement: "Top" (#F7F7F4) + " 1%" (#F5A817) (Font: Ubuntu Bold, Size: 84pt)
- Hook (y: 305): "Real Stories, Real People,\nReal Growth." (Font: Ubuntu Medium, Color: #B8BBC2, Size: 24pt)
- YouTube Line (y: 385): YouTube Play Icon (32px, #FF0000) + "{youtube_handle}" (Font: Ubuntu Bold, Color: #ECEBE6, Size: 20pt)

## Zone E: Lower Role Strip (x: 675, y: 550 to 720)

- Background: #211612 at 88% opacity, top line #7F5A45 at 35%
- Group 1: Book Icon + "Bestselling Author" (Font: Ubuntu Medium, Color: #F5F4F0, Size: 19pt)
- Group 2: Cap Icon + "Hard-skill Trainer" (Font: Ubuntu Medium, Color: #F5F4F0, Size: 19pt)
- Group 3: Megaphone Icon + "Brand Marketing Professional" (Font: Ubuntu Medium, Color: #F5F4F0, Size: 19pt)
- Group 4: Briefcase Icon + "Business Consultant" (Font: Ubuntu Medium, Color: #F5F4F0, Size: 19pt)
```
