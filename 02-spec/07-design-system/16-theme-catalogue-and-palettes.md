# 16 — Theme Catalogue & Multi-Theme Architecture

> **/goal** Provide a comprehensive, standardized catalogue of production-grade color themes (Navy Blue & Purple, VS Code themes, heatmaps, and warm editorial) with exact color formulas, semantic token mappings, contrast verification, and explicit guidance for AI models on selecting and applying the optimal visual identity.
> **/learn** Master the exact Hex, HSL, and semantic usage rules for each theme family, understand why and when each theme is deployed, and enforce the single-accent Von Restorff rule and 60/30/10 weight balance across all interface states.

**Version:** 4.0.0
**Updated:** 2026-09-24
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Executive Architecture & Theme Philosophy

A modern design system cannot rely on a single color palette. Different software domains require distinct emotional registers:
- **Cloud Infrastructure, AI, & Developer Tools:** Require deep navy, midnight slate, vibrant purple, and electric cyan to signal precision, concurrency, and intelligence.
- **Code Editors & Technical Dashboards:** Require syntax-informed VS Code palettes (Dark+, Tokyo Night, One Dark, Monokai Pro) that reduce eye strain during prolonged multi-hour focus.
- **Analytics, Metrics & Observability:** Require standardized sequential, diverging, and density heatmap scales (Plasma, Emerald-to-Gold, Red-Amber-Green) where color encodes mathematical values without cognitive ambiguity.
- **High-Trust Enterprise & Editorial Work:** Require warm-paper lights and warm-black darks with amber/bronze restraint that convey bespoke human craftsmanship rather than disposable SaaS templates.

### The Governing Rules for Every Theme

1. **60 / 30 / 10 Surface-to-Accent Ratio:**
   - **60% Dominant Neutral:** Canvas base, section planes, and outer wrappers.
   - **30% Structural Surfaces & Typography:** Card bodies, table rows, headers, hairlines, and muted content.
   - **10% Accent:** Primary call-to-action buttons, active indicators, and singular focal highlights.
2. **Single-Accent Von Restorff Rule:** At any given viewport or scroll position, **exactly ONE prominent accent-filled element** may be visible. If two buttons or cards both shout in primary saturation, the visual hierarchy collapses into noise.
3. **Semantic Decoupling:** Components NEVER declare literal color classes (e.g. `bg-blue-900`, `text-purple-400`). Components bind exclusively to semantic tokens (`--bg`, `--surface`, `--primary`, `--accent`, `--border`).

---

## 2. Theme Family 1: Navy Blue & Purple ("Midnight Cyber / Deep Oceanic")

### 2.1 Domain & Purpose
- **Recommended For:** AI platforms, distributed cloud platforms, cybersecurity consoles, developer tools, high-velocity SaaS.
- **Why Used:** Deep navy backgrounds provide a more comfortable, higher-contrast environment than pure black, while electric violet and ice cyan create visual energy and focal pop without triggering eye fatigue.

### 2.2 Token Registry

| Semantic Token | Hex | HSL (`H S% L%`) | Visual Role & Constraint |
|:---|:---|:---|:---|
| `--bg` (Base) | `#0A0F1D` | `224 49% 8%` | Deep midnight navy page ground (Plane 0) |
| `--bg-raised` | `#10172A` | `222 45% 11%` | Section alternation, container wells (Plane 1) |
| `--surface` | `#172036` | `223 40% 15%` | Card, table, and panel backgrounds (Plane 2) |
| `--surface-hover` | `#1E2945` | `223 39% 19%` | Interactive card and row hover states (Plane 3) |
| `--primary` | `#A855F7` | `271 91% 65%` | Electric Violet — primary CTAs, active pills |
| `--primary-hover` | `#C084FC` | `270 95% 75%` | Lightened violet for pointer hover |
| `--primary-fg` | `#0A0F1D` | `224 49% 8%` | Text on primary violet (High contrast: 11.2:1) |
| `--accent` | `#38BDF8` | `199 89% 60%` | Sky / Ice Cyan — data highlights, secondary cues |
| `--accent-fg` | `#0A0F1D` | `224 49% 8%` | Text on cyan accent (High contrast: 10.8:1) |
| `--depth-wash` | `#6366F1` | `239 84% 67%` | Royal Indigo — soft radial gradient glow (6–10% opacity) |
| `--border-hairline`| `rgba(248, 250, 252, 0.08)` | `210 40% 98% / 8%` | Structural container dividers and card edges |
| `--border-strong`  | `rgba(248, 250, 252, 0.16)` | `210 40% 98% / 16%`| Focus states and active card perimeter |
| `--fg` | `#F8FAFC` | `210 40% 98%` | Primary heading and body text (Contrast: 16.4:1) |
| `--fg-muted` | `#94A3B8` | `215 20% 65%` | Secondary descriptions and metadata (Contrast: 7.2:1)|
| `--fg-subtle` | `#64748B` | `215 16% 47%` | Labels, captions, table headers (Contrast: 4.6:1) |

---

## 3. Theme Family 2: VS Code Theme Ecosystem

### 3.1 VS Code Dark+ (Standard IDE Dark)
- **Domain:** Developer documentation, API portals, technical blogs, code-heavy apps.
- **Why Used:** Instantly recognizable to millions of engineers; established visual syntax creates immediate familiarity and comfort.

| Semantic Role | Token Name | Hex | HSL | Usage |
|:---|:---|:---|:---|:---|
| Canvas Ground | `--bg` | `#1E1E1E` | `0 0% 12%` | Editor canvas |
| Sidebar / Wells | `--bg-raised` | `#252526` | `240 1% 15%` | Navigation, file tree, toolbars |
| Component Panels| `--surface` | `#2D2D2D` | `0 0% 18%` | Modal dialogs, floating cards |
| Surface Hover | `--surface-hover` | `#37373D` | `240 5% 23%` | Active row hover |
| Primary Accent | `--primary` | `#007ACC` | `204 100% 40%` | Status blue, primary actions |
| Syntax Cyan | `--accent` | `#4EC9B0` | `168 53% 55%` | Type names, interfaces, success |
| Syntax Yellow | `--accent-warn`| `#DCDCAA` | `60 41% 77%` | Function calls, warnings |
| Syntax Orange | `--accent-str` | `#CE9178` | `17 48% 64%` | Strings, parameters, highlights |
| Syntax Purple | `--accent-kw`  | `#C586C0` | `305 34% 65%` | Keywords, control flow |
| Border | `--border` | `#3C3C3C` | `0 0% 24%` | Hairline dividers |
| Text Primary | `--fg` | `#D4D4D4` | `0 0% 83%` | Default code and prose |
| Text Muted | `--fg-muted` | `#858585` | `0 0% 52%` | Comments, line numbers, subtle text |

### 3.2 Tokyo Night (Storm Edition)
- **Domain:** Modern CLI tools, terminal dashboards, engineering community hubs.
- **Why Used:** Subtle deep blue-grey canvas prevents glare, paired with high-saturation neon pastel accents for maximum readability.

| Token Name | Hex | HSL | Role |
|:---|:---|:---|:---|
| `--bg` | `#1F2335` | `230 26% 17%` | Storm canvas ground |
| `--bg-raised` | `#24283B` | `229 24% 19%` | Section backgrounds and sidebars |
| `--surface` | `#292E42` | `229 23% 21%` | Card surfaces |
| `--primary` | `#7AA2F7` | `221 89% 72%` | Electric cornflower blue CTA |
| `--accent` | `#BB9AF7` | `261 88% 79%` | Soft lilac purple highlight |
| `--cyan` | `#7DCFFF` | `202 100% 75%` | Ice cyan data points |
| `--green` | `#9ECE6A` | `88 48% 65%` | Fresh spring green for positive states |
| `--orange` | `#FF9E64` | `23 100% 69%` | Warm peach orange for alerts |
| `--fg` | `#C0CAF5` | `228 75% 86%` | Pale violet white primary text |
| `--fg-muted` | `#787C99` | `232 13% 54%` | Slate purple secondary text |

### 3.3 One Dark Pro (Atom Classic)
- **Domain:** Documentation viewers, technical specifications, markdown readers.
- **Why Used:** Chalky slate base with soft, warm pastel syntax accents that have been battle-tested over a decade of developer usage.

| Token Name | Hex | HSL | Role |
|:---|:---|:---|:---|
| `--bg` | `#282C34` | `220 13% 18%` | Signature dark slate base |
| `--bg-raised` | `#21252B` | `220 13% 15%` | Sunken toolbars and tree views |
| `--surface` | `#2E3440` | `220 16% 22%` | Raised cards and popovers |
| `--primary` | `#61AFEF` | `207 82% 66%` | Calming sky blue primary action |
| `--accent` | `#C678DD` | `286 60% 67%` | Lavender violet emphasis |
| `--coral` | `#E06C75` | `355 65% 65%` | Soft red-coral alert / deletion |
| `--sage` | `#98C379` | `95 38% 62%` | Sage green success / insertion |
| `--amber` | `#E5C07B` | `39 67% 69%` | Soft amber warning / constants |
| `--fg` | `#ABB2BF` | `219 14% 71%` | Neutral warm slate text |

### 3.4 GitHub Dark High Contrast
- **Domain:** Accessibility-first interfaces, mission-critical operations, financial ledgers.
- **Why Used:** Exceeds WCAG AAA requirements with zero visual ambiguity and crisp boundaries.

| Token Name | Hex | HSL | Role |
|:---|:---|:---|:---|
| `--bg` | `#0A0C10` | `220 23% 5%` | Ultra-deep black-blue base |
| `--bg-raised` | `#12151C` | `223 22% 9%` | High-contrast sunken containers |
| `--surface` | `#1C212C` | `220 22% 14%` | Solid elevated cards |
| `--primary` | `#79C0FF` | `209 100% 74%` | High-chroma bright blue |
| `--primary-fg` | `#010409` | `218 80% 2%` | Black text on blue (Contrast: 13.5:1) |
| `--border` | `#7A828E` | `217 9% 52%` | Explicit visible 1px borders |
| `--fg` | `#FFFFFF` | `0 0% 100%` | Stark pure white text |
| `--fg-muted` | `#C9D1D9` | `210 17% 82%` | High-contrast muted text |

### 3.5 Monokai Pro
- **Domain:** Gaming platforms, creative tooling, telemetry feeds, high-impact landing pages.
- **Why Used:** Warm organic slate ground with electric fluorescent accents for maximum personality and visual snap.

| Token Name | Hex | HSL | Role |
|:---|:---|:---|:---|
| `--bg` | `#2D2A2E` | `285 5% 18%` | Warm espresso-slate ground |
| `--bg-raised` | `#221F22` | `300 5% 13%` | Deep inset panels |
| `--surface` | `#3A363B` | `288 4% 22%` | Elevated cards |
| `--primary` | `#FF6188` | `345 100% 69%` | Fluorescent rose magenta |
| `--accent` | `#FFD866` | `45 100% 70%` | Radiant sunny gold |
| `--mint` | `#A9DC76` | `90 60% 66%` | Lime mint green |
| `--cyan` | `#78DCE8` | `187 71% 69%` | Clean turquoise cyan |
| `--violet` | `#AB9DF2` | `250 78% 78%` | Electric lavender violet |
| `--fg` | `#FCFCFA` | `60 33% 98%` | Warm porcelain white |

---

## 4. Theme Family 3: Heat Map & Density Color Scales

When building data visualizations, activity matrices, risk grids, or performance heatmaps, colors represent **quantitative values**, not branding. AI models must apply these standardized scales strictly:

### 4.1 Sequential Plasma Scale (Thermal / Intensity)
Used for continuous metrics: latency distributions, CPU temperatures, request volume, error densities.

```text
Low Intensity (0%) ──────────────────────────────────────────► High Intensity (100%)
[#0D0887]      [#4B03A1]      [#7D03A8]      [#A82296]      [#CB4679]      [#E56B5D]      [#F89441]      [#FDC328]      [#F0F921]
Deep Midnight  Royal Violet   Magenta        Ruby Rose      Coral Thermal  Amber Flame    Gold Sun       Pure Radiant
```

| Step | Hex | HSL | Semantic Meaning |
|:---|:---|:---|:---|
| Step 0 (0–10%) | `#0D0887` | `242 83% 29%` | Baseline / Idle / Zero traffic |
| Step 1 (11–25%) | `#5A01A5` | `273 99% 33%` | Nominal low activity |
| Step 2 (26–40%) | `#8F0DA4` | `292 85% 35%` | Moderate load |
| Step 3 (41–55%) | `#BC3754` | `347 55% 48%` | Elevated threshold |
| Step 4 (56–70%) | `#DE5F39` | `14 73% 55%` | Substantial activity |
| Step 5 (71–85%) | `#F58C46` | `24 90% 62%` | High pressure / Hotspot |
| Step 6 (86–100%)| `#F0F921` | `63 93% 55%` | Peak capacity / Critical saturation |

### 4.2 Diverging Red-Amber-Green Scale (Quality, Risk, SLA)
Used for directional data centered on a neutral zero or target threshold: latency delta, test results, code coverage deviation, budget margins.

| Grade | Hex | HSL | Condition |
|:---|:---|:---|:---|
| Critical Breach | `#DC2626` | `0 72% 51%` | Severe failure, outage, negative deviation > 20% |
| High Risk | `#EF4444` | `0 84% 60%` | Warning threshold breached, test failure |
| Warning / Caution| `#F59E0B` | `38 92% 50%` | Approaching quota, flaky build, medium latency |
| Neutral / Unchanged| `#64748B` | `215 16% 47%` | 0% change, baseline maintained |
| Acceptable | `#84CC16` | `84 81% 44%` | Nominal pass, SLA met |
| Optimal / Leader | `#10B981` | `160 84% 39%` | Superior performance, 100% test pass, zero defect |

### 4.3 Activity Density Grid (GitHub Style)
Used for commit frequency, test run history, uptime streaks, audit logs.

| Activity Level | Dark Mode Hex | Light Mode Hex | Threshold Rule |
|:---|:---|:---|:---|
| Level 0 (None) | `#161B22` | `#EBEDF0` | Exactly 0 events |
| Level 1 (Low) | `#0E4429` | `#9BE9A8` | 1–3 events |
| Level 2 (Medium) | `#006D32` | `#40C463` | 4–9 events |
| Level 3 (High) | `#26A641` | `#30A14E` | 10–19 events |
| Level 4 (Max) | `#39D353` | `#216E39` | 20+ events |

---

## 5. Theme Family 4: Warm Editorial & Craft Theme

### 5.1 Domain & Purpose
- **Recommended For:** Executive presentations, agency showcases, portfolio platforms, architecture reviews.
- **Why Used:** Avoids cold, disposable tech-startup aesthetics by combining warm-paper lights and warm-black darks with an ember/amber accent.

| Token | Dark Mode (`#0B0A09` base) | Light Mode (`#FBF9F6` base) | Role |
|:---|:---|:---|:---|
| `--bg` | `#0B0A09` (Warm-Black) | `#FBF9F6` (Warm Paper) | Page foundation |
| `--bg-raised` | `#151312` | `#F4F0EA` | Alternate section bands |
| `--surface` | `#1D1A18` | `#FFFFFF` | Card surfaces |
| `--primary` | `#FFAD01` (Brand Amber) | `#FFAD01` fill / `#8A5A00` text | Primary CTA / ink accent |
| `--depth` | `#791AB0` (Royal Violet) | `#791AB0` (depth only) | Ambient glow, funnel depth |
| `--fg` | `#F7F5F2` (Warm White) | `#141210` (Deep Ink) | Primary typography |
| `--fg-muted` | `#ADA8A2` | `#57514B` | Secondary body text |

---

## 6. How AI Models Must Choose and Apply Themes

When generating or refactoring UI, an AI agent must follow this decision tree:

```text
[Incoming Project Domain]
   │
   ├─► Developer CLI / IDE / Technical Docs?
   │     └─► Select VS Code Theme Family (Dark+, Tokyo Night, or One Dark)
   │
   ├─► Cloud Infrastructure / AI Platform / Cybersecurity?
   │     └─► Select Navy Blue & Purple Theme (#0A0F1D + Electric Violet #A855F7)
   │
   ├─► Telemetry / Metrics Dashboard / Activity Ledger?
   │     └─► Base: VS Code Dark+ or Navy; Heatmaps: Sequential Plasma or Density Grid
   │
   ├─► High-Trust Enterprise / Consulting / Portfolio?
   │     └─► Select Warm Editorial Theme (#0B0A09 / #FBF9F6 + Amber #FFAD01)
   │
   └─► Accessibility-Mandated Government / Financial Portal?
         └─► Select GitHub Dark High Contrast (WCAG AAA certified)
```

---

## 7. AI Verification Checklist for Themes & Palettes

- [ ] `/goal` Verify all colors are declared as CSS custom properties in `:root` and `.dark`.
- [ ] `/learn` Verify zero literal color classes (`text-white`, `bg-black`, `bg-purple-600`) in component code.
- [ ] `/goal` Verify 60/30/10 visual balance: 60% base plane, 30% structural surfaces/type, 10% accent.
- [ ] `/learn` Verify single-accent Von Restorff: exactly one prominent accent button/pill visible per viewport.
- [ ] `/goal` Verify WCAG 2.2 AA contrast: body text on background ≥ 4.5:1, large text ≥ 3:1.
- [ ] `/learn` Verify data heatmaps use standardized perceptual scales (Plasma, Diverging RAG, or Activity Grid) rather than random decorative colors.
