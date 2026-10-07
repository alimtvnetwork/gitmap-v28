# 36 — Website Content Builder Mode & In-Page Visual Editor Specification

> **/goal** Specify the production-grade, zero-server, review-only visual builder overlay that allows clients and AI agents to edit copy, media, menu items, and section order directly on the live website and export deterministic change manifests for developers.
> **/learn** Master the activation gate (`?builder=1&email={OWNER_EMAIL}`), 5 operational modes (Text, Images, Menu, Layout, Off), in-place `contentEditable` lifecycle, gradient text flattening, DOM sanitization allowlist, non-destructive React icon preservation, floating side panel, change tracking diffs, and deterministic ZIP exports.

**Version:** 4.4.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. System Overview & Core Philosophy

The **Website Content Builder Mode** is a zero-backend, client-side visual review overlay. It eliminates ambiguity between clients, designers, and developers:
1. **Review-Only Overlay (Total Ban on Live Auto-Publishing):** Edits are stored in client `localStorage` and rendered only when the builder session is active. The production website always renders clean source content. Live sites change ONLY when a developer applies the exported changes.
2. **Obscurity Gate (Zero Server Leakage):** Activated strictly via URL query parameters. No login form, no server-side auth endpoints, and no builder code bundled into public production visitors.
3. **Deterministic Change Exports:** Produces a standardized ZIP package containing human-readable Markdown change summaries, machine-applicable JSON AST diff manifests, and replaced image assets.

---

## 2. Activation & Session Lifecycle

### 2.1 URL Activation Gate
The editor bundle initializes only when both query parameters are present on any page URL:

```text
https://{SITE_DOMAIN}/{ANY_PAGE}?builder=1&email={OWNER_EMAIL}
```

- **Authentication via Key:** The owner email address **is** the key. If `builder=1` is missing, or the email fails validation against the client allowlist, the runtime bails out immediately without mounting UI, setting attributes, or downloading builder assets.
- **Query Parameter Preservation:** To maintain editing context during single-page app (SPA) navigation, the builder intercepts history events (`pushState`, `replaceState`, `popstate`) and re-appends `?builder=1&email={OWNER_EMAIL}` to the address bar.
- **Storage Sanitization:** Query parameters are stripped from storage keys and exported URLs.

### 2.2 Hydration Wait & Mutation Observer
1. **Hydration Quiet Window:** The orchestrator waits for `domQuiet(700ms, 4000ms)` via `MutationObserver` on `document.body` before tagging the DOM, preventing hydration mismatch warnings and lost React fiber references.
2. **Dynamic Rescan:** A debounced (120ms) observer catches lazy-loaded sections and route transitions, assigning IDs and re-applying draft modifications without infinite loops (all builder UI carries `data-bm-root` and is ignored by mutations).

---

## 3. Operational Modes & Interaction Taxonomy

A mode selector in the floating side panel switches the entire viewport between 5 mutually exclusive behaviors:

| Mode ID | Visual Trigger | Target Elements | Editing Behavior |
|:---|:---|:---|:---|
| **`Text`** *(Default)* | Double-click | Headings (`h1`–`h6`), paragraphs, list items, button labels, stat numerals, captions | In-place `contentEditable` with visible caret, accent outline, and highlight toolbar |
| **`Images`** | Single-click | Photographic images (`img`), picture sources, and inline SVG icons | Opens modal: image file upload (PNG/JPG/WebP/SVG max 2MB), required alt-text, SEO filename |
| **`Menu`** | Single-click | Primary nav items, header links, and footer links | Modal editing of Link Label, Description, and Link Target (`/`, `#`, `http(s)://`, `mailto:`, `tel:`) |
| **`Layout`** | Hover controls | Repeated card decks, feature grids, list rows | Reorder items up/down within sibling container, drag order handles |
| **`Off`** | Panel toggle | Viewport canvas | Hides all builder outlines, chips, and panel; keeps draft modifications visible on page |

### 3.2 DOM Auto-Tagging Pipeline (`assignIds()`) & Content Index Mapping

To reliably map live DOM nodes back to source code files without manual JSX tagging:
1. **Content Index (At Start):** The build module flattens all marketing data files (`src/content/*.ts`) to a depth cap of 6 into a normalized string-to-source map:
   `index.set(normalize(value), { file: "src/content/home.ts", key: "home.hero.title" })`.
2. **Deterministic 4-Pass DOM Tagging Order:**
   - **Pass 1 — Menu Links:** `header a, footer nav a` tagged as `data-builder-type="menu"`, with id `global.nav.<slug>` or `global.footer.<slug>`.
   - **Pass 2 — Images & Logo:** Every `<img>` tagged as `data-builder-type="image"`; logo tagged `global.header.logo`.
   - **Pass 3 — Animated / Masked Headings:** Containers with `[data-bm-masked="words"|"lines"]` tagged as single editable units.
   - **Pass 4 — Generic Text:** Evaluates remaining nodes against the text tag allowlist (`H1-H6, P, LI, SPAN, A, BUTTON, STRONG, LABEL, BLOCKQUOTE, TD, TH`).
3. **Split on Block-Level Children:** A hero heading with stacked `<span class="block">` lines is edited line-by-line, preventing accidental text collapse on save. Decorative elements (SVGs, `aria-hidden="true"`) do not disqualify parent buttons or labels.
4. **Attributes Written:**
   - `data-builder-id`: Content key or generated `<page-slug>.auto.<n>`.
   - `data-builder-type`: `"text" | "image" | "menu" | "layout"`.
   - `data-bm-section-name`: Nearest section heading text or semantic role.
   - `data-bm-order`: Document index for sorting exports top-to-bottom.
   - `data-bm-source`: File path and dotted object path.

---

## 4. Text Editing Lifecycle & DOM Sanitization

### 4.1 In-Place Text Editing Workflow
1. **Trigger:** Double-click on any `[data-builder-type="text"]`.
2. **Active Outline:** Sets `data-bm-editing="true"` with a `2px solid var(--primary)` focus ring and `border-radius: 4px`.
3. **Caret Placement:** Employs `document.caretPositionFromPoint(x, y)` (standards) or `document.caretRangeFromPoint(x, y)` (WebKit/Blink fallback) so the cursor blinks precisely where clicked.
4. **Gradient Text Flattening:** Brand headings using `background-clip: text` with transparent text are flattened to solid ink while editing so the caret remains legible:
   ```css
   html[data-bm-on="true"] [data-bm-editing="true"],
   html[data-bm-on="true"] [data-bm-editing="true"] .gradient-text {
     -webkit-text-fill-color: currentColor !important;
     background: none !important;
     color: var(--ink, #0f172a) !important;
     filter: none !important;
   }
   ```
5. **Floating Keyboard Chip:** Displays floating helper pill above element: *"Editing — select word to highlight · Ctrl/Cmd+Enter saves, Esc cancels"*.
6. **Selection Highlight Toolbar:** Selecting one or more words inside the active element displays a floating mini-toolbar with:
   - **Highlight:** Wraps selection in `<span class="gradient-text" data-text="{WORD}">` applying brand gradient.
   - **Remove Highlight:** Unwraps gradient span back to plain text.
7. **Commit & Cancel:**
   - Commit on `blur` or `Cmd/Ctrl + Enter`.
   - Cancel on `Escape` (restores original `innerHTML` verbatim).
   - Empty strings are rejected with an inline warning: *"Text cannot be empty"*.

### 4.2 HTML Sanitizer & Plain-Text Paste Allowlist
To prevent XSS and style regressions, rich clipboard formatting is intercepted. Plain text is inserted via `document.execCommand('insertText')`. Allowed tags on save:

```text
ALLOWED_TAGS = ['STRONG', 'B', 'EM', 'I', 'BR', 'A', 'SPAN']
```

- `SPAN` is allowed **strictly** when `class="gradient-text"`. All other spans are unwrapped.
- All attributes are stripped except `href` on `<a>`, and `class` plus `data-text` on `<span class="gradient-text">`.
- `href` values beginning with `javascript:` or `data:` are purged.
- Preserves child SVG icons by marking them `contenteditable="false"` prior to edit start.

---

## 5. Media & Navigation Editors

### 5.1 Images & Icon Replacement Modal
- **Accepted Formats:** `image/png`, `image/jpeg`, `image/webp`, `image/svg+xml`. Hard reject on files `> 2 MB`.
- **Required Alt Text:** Form validation blocks save if alt-text field is empty.
- **Deterministic SEO Naming:** Automatically names uploaded image `{page-slug}-{element-slug}-{short-hash}.{ext}`.
- **Inline SVG Support:** Custom SVG icons parsed via `DOMParser(..., "image/svg+xml")` ensuring `viewBox` preservation and script stripping.

### 5.2 Menu Link Editor
- **Inputs:** Navigation Label, Optional Subtitle / Description, and Target Destination.
- **Protocol Allowlist:** Target must start with `/` (relative route), `#` (anchor jump), `https://`, `http://`, `mailto:`, or `tel:`.
- **Navigation Suppression:** While builder mode is active, normal clicks on `<a>` links outside the builder UI are intercepted (`e.preventDefault()`) so links can be clicked to edit. Opening links for actual navigation requires holding `Cmd`, `Ctrl`, or `Alt`.

---

## 6. Floating Side Panel & Change History

### 6.1 Panel Geometry & Layout
The panel is fixed to the right viewport edge (`top: 0, right: 0, bottom: 0, width: 320px, z-index: 9999`), collapsible into an unbranded `"Editor"` tab:
1. **Master Toggle:** Enable/disable builder mode.
2. **"Show Original" Button:** Temporarily suppresses stored draft to compare against baseline.
3. **Mode Radio Bar:** Quick-switch between Text, Images, Menu, Layout, and Off.
4. **Change Counter:** Live badge displaying `"X changes on this page · Y total"`.
5. **Action Row:**
   - **Undo (`Cmd+Z`):** Reverts single most recent change across the session.
   - **Reset Page:** Dialog confirmation restoring current page to baseline.
   - **History Button:** Opens modal listing all changes.
   - **Export Button:** Opens export dialog.

### 6.2 History Diff Ledger
Every save registers a structured record into `localStorage["builder-draft-v1"]`:

```typescript
export interface ChangeRecord {
  id: string;              // Deterministic element key (e.g. "home-hero-title")
  page: string;            // Page pathname (e.g. "/solutions")
  section: string;         // Section name (e.g. "Hero")
  type: "text" | "image" | "menu" | "layout";
  before: string;          // Baseline markup or image URL
  after: string;           // Modified markup or base64 data URI
  timestamp: number;       // Unix epoch ms
}
```

The History Modal renders changes grouped by page, showing timestamps and git-style `−` before / `+` after diff lines.

---

## 7. Deterministic ZIP Export Specification

Clicking **"Export (.zip)"** generates an archive named `content-changes--all-pages--YYYY-MM-DD-HHmm.zip` structured as:

```text
content-changes--all-pages--2026-10-02-1200.zip
├── SUMMARY.md              (Executive overview, page counts, change statistics)
├── manifest.json           (Machine-applicable JSON diff array for automated codegen)
├── pages/
│   ├── home.md             (Section-by-section markdown change report for /)
│   └── solutions.md        (Section-by-section markdown change report for /solutions)
├── menu-changes.md         (Navigation link updates and target diffs)
└── assets/
    ├── home-hero-bg-a1b2c3.webp
    └── solutions-icon-d4e5f6.svg
```

### 7.1 Machine-Readable Manifest (`manifest.json`)
```json
[
  {
    "id": "home-hero-headline",
    "page": "/",
    "section": "Hero",
    "type": "text",
    "sourceFile": "src/content/home.ts",
    "sourceKey": "home.hero.title",
    "before": "Transforming Enterprise Operations",
    "after": "Accelerating Modern Cloud Operations",
    "timestamp": 1790942400000
  }
]
```

### 7.2 Automated Before / After Section Screenshot Pairs (`shots.ts`)

During the export routine, the builder captures side-by-side photographic visual verification for each modified section using `html-to-image`:
1. **After State:** Captures the current DOM rendering with the edited element highlighted in a red accent outline (`#ef4444`, `2px solid`). Saved as `screenshots/{page-slug}/{NN}-after.png`.
2. **Before State:** Temporarily paints the original record value back onto the DOM element via `paintSide(node, rec, "before")`, highlights it in a muted slate outline (`#64748b`, `2px solid`), captures `screenshots/{page-slug}/{NN}-before.png`, and immediately restores the active draft.
3. **Capture Parameters:** Section bounds clamped to `1440 × 2400px`, `skipFonts: true`, `pixelRatio <= 1.5`, filtering out builder overlay chrome (`[data-bm-root]`, `.bm-chip`).
4. **Markdown Side-by-Side Tables:** Generates GitHub-flavored markdown tables presenting side-by-side comparisons of the before and after screenshots for human code review.

---

## 8. Anti-Hallucination & Quality Verification Checklist

- [ ] Builder activation strictly requires both `?builder=1` and `email={OWNER_EMAIL}`.
- [ ] No server endpoints or database publishing exists; overlay is strictly review-only.
- [ ] Query parameters are re-attached on SPA navigation and stripped from exported URLs.
- [ ] In-place text editing flattens gradient text to solid ink while editing for caret legibility.
- [ ] Caret is explicitly positioned via `caretPositionFromPoint` or `caretRangeFromPoint`.
- [ ] Plain-text pasting enforced; allowed tags restricted strictly to `STRONG`, `B`, `EM`, `I`, `BR`, `A`, and `<span class="gradient-text">`.
- [ ] Navigation on links is suppressed during editing unless modifier keys (`Cmd`, `Ctrl`, `Alt`) are held.
- [ ] Images mode restricts uploads to PNG, JPG, WebP, SVG with a 2 MB hard cap and mandatory alt-text.
- [ ] Menu mode validates target links against approved protocols (`/`, `#`, `http(s)://`, `mailto:`, `tel:`).
- [ ] Export routine generates before/after screenshot pairs with red (`#ef4444`) and slate (`#64748b`) outlines.
- [ ] Export produces deterministic ZIP with `SUMMARY.md`, `manifest.json`, per-page Markdown diffs, screenshots, and assets.
