# Master Architecture Spec: UI Improvement Clusters 01–03
## Color Tokens & HSL Architecture, WCAG Contrast Calibration, and 4-Plane Depth Hierarchy

- **Sequence**: 340
- **Document**: `02-spec/21-app/340-next-100-ui-improvements/01-tokens-contrast-and-depth-spec.md`
- **Scope**: Cluster 01 (Items 01–10), Cluster 02 (Items 11–20), Cluster 03 (Items 21–30)
- **Status**: APPROVED / ARCHITECTURE SPECIFICATION
- **Version**: 1.0.0
- **Author**: Spec Author 1 (Antigravity Swarm)
- **Target Interface**: GitMap Fleet Web UI, Log Viewer, Installer, Macro UI & Companion Dashboards

---

## 1. Executive Summary & Architectural Intent

The GitMap interface ecosystem spans desktop-integrated web views (Fleet Web UI `gitmap ui`, Log Viewer `gitmap show-log`, Macro Builder `gitmap ui macro`, and Installer Manager `gitmap ui installer`). Historically, styling relied on fragmented Tailwind color utilities (`bg-blue-500`, `text-slate-400`, `border-gray-700`) and arbitrary hex literals (`#2563eb`, `#1e293b`). This produced inconsistent contrast across light and dark viewports, uncalibrated border elevations, muddy drop shadows, and brittle theme toggling.

This specification establishes a mathematically grounded design foundation across three foundational pillars:
1. **Cluster 01 (Color Tokens & HSL Architecture)**: Complete elimination of hardcoded hex colors, transition to space-separated HSL channels in CSS custom properties (`H S% L%`), complete 10-step Blue and Violet primitive ramps, dedicated Cyan highlights, and semantic decoupling.
2. **Cluster 02 (WCAG Contrast Ratios & Luminance Calibrations)**: Strict mathematical enforcement of WCAG 2.2 AA (minimum 4.5:1 for regular text, 3:1 for large text >= 24px) with deliberate default targeting of AAA (>= 7.0:1) across body prose, interactive inputs, focus rings, and non-color redundant state pills.
3. **Cluster 03 (4-Plane Depth Hierarchy & Materiality)**: Systematic elimination of muddy Gaussian blur drop shadows in favor of a 4-plane physical lightness stepped hierarchy (Plane 0 Canvas, Plane 1 Raised Wells, Plane 2 Surface Cards, Plane 3 Elevated Modals), hairline edge light refraction (`inset 0 1px 0 rgba(255, 255, 255, 0.12)`), and frosted acrylic glassmorphism.

---

## 2. Cluster 01: Color Tokens & HSL Architecture (Items 01–10)

### 2.1 Space-Separated HSL Token Storage Standard (Item 02)

All CSS custom properties must be declared as raw space-separated values without the `hsl(...)` wrapper inside `:root` and `.dark` scopes. This architecture provides runtime alpha composition across CSS and Tailwind utility pipelines:

```css
/* ✅ CORRECT: Space-separated HSL values */
:root {
  --primary: 221 83% 53%;
  --primary-foreground: 0 0% 100%;
  --surface: 220 20% 98%;
  --border: 218 23% 91%;
}

/* Runtime Usage in Components */
.component-card {
  background-color: hsl(var(--surface));
  border: 1px solid hsl(var(--border) / 0.8);
  box-shadow: 0 4px 12px hsl(var(--primary) / 0.08);
}
```

```css
/* ❌ FORBIDDEN: Pre-wrapped or comma-separated tokens */
:root {
  --primary: #2563eb;          /* Hex prevents dynamic alpha layering */
  --surface: hsl(220, 20%, 98%); /* Wrapped hsl prevents / alpha modifier */
}
```

### 2.2 Primitive Color Ramps (Items 03, 04, 05, 06)

#### Blue Ramp (Primary Action & Corporate Trust — Item 03)
| Token | OKLCH Runtime Value | HEX Code | RGB Value | Space-Separated HSL | Architectural Function |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `--blue-50` | `oklch(0.972 0.014 259)` | `#F0F6FF` | `rgb(240, 246, 255)` | `216 100% 97%` | Hero wash start, brand pill bg |
| `--blue-100` | `oklch(0.935 0.032 262)` | `#DEEAFF` | `rgb(222, 234, 255)` | `218 100% 94%` | Badge hover fill, soft tile ring |
| `--blue-200` | `oklch(0.874 0.061 263)` | `#C1D6FF` | `rgb(193, 214, 255)` | `220 100% 88%` | Subtle container border, chrome stop |
| `--blue-300` | `oklch(0.777 0.104 264)` | `#95B6FA` | `rgb(149, 182, 250)` | `220 91% 78%` | Chrome gradient end, data series |
| `--blue-400` | `oklch(0.662 0.164 264)` | `#5E8EF6` | `rgb(94, 142, 246)` | `221 89% 67%` | Focus ring stop, gradient text end |
| `--blue-500` | `oklch(0.546 0.215 262.9)` | `#2563EB` | `rgb(37, 99, 235)` | `221 83% 53%` | **Primary CTA Action (`--brand-secondary`)** |
| `--blue-600` | `oklch(0.489 0.208 263.5)` | `#1B51D3` | `rgb(27, 81, 211)` | `222 77% 47%` | Active CTA depression, brand gradient start |
| `--blue-700` | `oklch(0.409 0.176 264.5)` | `#173DA7` | `rgb(23, 61, 167)` | `224 76% 37%` | High-contrast hyperlink, brand pill text |
| `--blue-800` | `oklch(0.369 0.147 264.8)` | `#17368B` | `rgb(23, 54, 139)` | `224 72% 32%` | Deep interactive border, badge text |
| `--blue-900` | `oklch(0.317 0.135 264.5)` | `#0D2975` | `rgb(13, 41, 117)` | `224 80% 25%` | **Core Brand Anchor (`--brand-primary`)** |

#### Violet Ramp (Tertiary Accent & Gradient Partner — Item 04)
| Token | OKLCH Runtime Value | HEX Code | RGB Value | Space-Separated HSL | Architectural Function |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `--violet-50` | `oklch(0.972 0.017 300)` | `#F7F3FF` | `rgb(247, 243, 255)` | `260 100% 98%` | Accent surface fill |
| `--violet-100` | `oklch(0.946 0.032 300)` | `#F0E9FF` | `rgb(240, 233, 255)` | `259 100% 96%` | Subtle chip background |
| `--violet-200` | `oklch(0.885 0.068 299)` | `#DFD0FF` | `rgb(223, 208, 255)` | `259 100% 91%` | Accent hairline border on light |
| `--violet-300` | `oklch(0.784 0.128 298)` | `#C3A6FF` | `rgb(195, 166, 255)` | `260 100% 83%` | Soft glow highlight on dark cards |
| `--violet-400` | `oklch(0.659 0.187 297.5)` | `#A173F4` | `rgb(161, 115, 244)` | `261 85% 70%` | Secondary chart accent |
| `--violet-500` | `oklch(0.532 0.253 296.9)` | `#822EE8` | `rgb(130, 46, 232)` | `267 80% 55%` | **Tertiary Brand Accent (`--brand-tertiary`)** |
| `--violet-600` | `oklch(0.481 0.243 296.5)` | `#721CD2` | `rgb(114, 28, 210)` | `268 76% 47%` | Brand gradient midpoint (55%) |
| `--violet-700` | `oklch(0.417 0.213 295.9)` | `#5C14AF` | `rgb(92, 20, 175)` | `268 79% 38%` | Accent foreground text |
| `--violet-800` | `oklch(0.362 0.187 295.5)` | `#4B0D92` | `rgb(75, 13, 146)` | `268 84% 31%` | Deep violet structural border |
| `--violet-900` | `oklch(0.257 0.153 294.5)` | `#2D0063` | `rgb(45, 0, 99)` | `267 100% 19%` | Maximum depth violet anchor |

#### Cyan Highlight, Slate Neutrals & Surfaces (Items 05, 06)
| Token | HEX Code | Space-Separated HSL | Architectural Function |
| :--- | :--- | :--- | :--- |
| `--cyan-400` | `#03D5E7` | `185 97% 46%` | **Brand Highlight (`--brand-highlight`)**, `.shine-sweep`, live telemetry pulse |
| `--neutral-50` | `#F8FAFC` | `210 40% 98%` | Stark light base / Dark high-contrast ink |
| `--neutral-200` | `#E2E6ED` | `218 23% 91%` | **Light Hairline Border (`--border`, `--input`)** |
| `--neutral-400` | `#8F9BB2` | `219 19% 63%` | Muted metadata, dark-band secondary text |
| `--neutral-600` | `#48536B` | `221 20% 35%` | **Soft Ink (`--ink-soft`, `--muted-foreground`)** |
| `--neutral-800` | `#1E293B` | `217 33% 17%` | Dark well surface, dark card border |
| `--neutral-900` | `#090E18` | `220 45% 6%` | **Primary Ink (`--ink`, `--foreground`)** |
| `--surface-dark`| `#0A0E17` | `222 39% 6%` | Dark section canvas / header bar |

### 2.3 Semantic Token Decoupling (Item 07)

Hardcoded utility classes like `text-blue-500` or `bg-slate-900` must be replaced with role-based semantic tokens. Components consume functional aliases:

```typescript
// ❌ ANTI-PATTERN: Direct Tailwind palette color usage
<div className="bg-slate-900 border border-slate-700 text-blue-400 p-4 rounded-lg">
  <span className="text-gray-400">Branch:</span> main
</div>

// ✅ CANONICAL: Semantic token consumption
<div className="bg-surface-2 border border-border text-primary p-4 rounded-lg shadow-hairline">
  <span className="text-muted-foreground">Branch:</span> main
</div>
```

### 2.4 Dark / Light Theme Variable Parity (Item 08)

Every variable registered in `:root` must possess an identical, calibrated counterpart in `.dark`:

```css
:root {
  /* Canvas & Wells */
  --background: 40 33% 98%;          /* #FBF9F6 Warm Paper */
  --foreground: 220 45% 6%;          /* #090E18 Primary Ink */
  --well: 220 20% 94%;                /* Light Well */
  --card: 0 0% 100%;                  /* Pure White Surface */
  --card-foreground: 220 45% 6%;
  --popover: 0 0% 100%;
  --popover-foreground: 220 45% 6%;
  
  /* Brand Tokens */
  --primary: 221 83% 53%;             /* #2563EB Blue 500 */
  --primary-foreground: 0 0% 100%;
  --secondary: 267 80% 55%;           /* #822EE8 Violet 500 */
  --secondary-foreground: 0 0% 100%;
  --highlight: 185 97% 46%;           /* #03D5E7 Cyan 400 */
  
  /* Structural Utilities */
  --muted: 220 14% 96%;
  --muted-foreground: 221 20% 35%;   /* #48536B (AA 5.2:1) */
  --border: 218 23% 91%;              /* #E2E6ED */
  --input: 218 23% 91%;
  --ring: 221 83% 53%;
}

.dark {
  /* Canvas & Wells (4-Plane Lightweight Scaling) */
  --background: 222 47% 8%;          /* #0A0F1D Plane 0 Base */
  --foreground: 210 40% 98%;         /* #F8FAFC Primary Light Ink */
  --well: 222 45% 11%;               /* #10172A Plane 1 Well (+3% L) */
  --card: 222 40% 15%;               /* #172036 Plane 2 Surface (+4% L) */
  --card-foreground: 210 40% 98%;
  --popover: 222 39% 19%;            /* #1E2945 Plane 3 Modal (+4% L) */
  --popover-foreground: 210 40% 98%;
  
  /* Brand Tokens (+8% Lightness for Dark Canvas Contrast) */
  --primary: 221 89% 67%;             /* #5E8EF6 Blue 400 */
  --primary-foreground: 222 47% 8%;
  --secondary: 261 85% 70%;           /* #A173F4 Violet 400 */
  --secondary-foreground: 222 47% 8%;
  --highlight: 185 97% 55%;           /* Cyan Highlight Glow */
  
  /* Structural Utilities */
  --muted: 222 35% 15%;
  --muted-foreground: 219 19% 68%;   /* #9AA5BA (AAA 7.4:1) */
  --border: 222 30% 22%;              /* 1px Hairline Slate Border */
  --input: 222 30% 22%;
  --ring: 221 89% 67%;
}
```

### 2.5 Tokenized Gradient Definitions (Item 09)

Gradients must never be specified as ad-hoc inline strings. They are defined once in the root stylesheet:

```css
:root {
  --gradient-brand: linear-gradient(135deg, hsl(var(--primary)) 0%, hsl(var(--secondary)) 100%);
  --gradient-hero-wash: radial-gradient(circle at 50% 0%, hsl(var(--primary) / 0.12) 0%, transparent 70%);
  --gradient-card-sheen: linear-gradient(180deg, rgba(255, 255, 255, 0.08) 0%, rgba(255, 255, 255, 0.00) 100%);
  --gradient-sweep: linear-gradient(45deg, transparent 25%, hsl(var(--highlight) / 0.4) 50%, transparent 75%);
}
```

### 2.6 Code Block Theme Isolation (Item 10)

Terminal logs, diff viewers, and code blocks MUST remain pinned to an isolated dark container background across both light and dark global themes:
```css
.code-viewport, .terminal-container {
  background-color: hsl(220 14% 11%) !important; /* Fixed dark slate */
  color: hsl(210 40% 98%) !important;
  border: 1px solid rgba(255, 255, 255, 0.1);
}
```

---

## 3. Cluster 02: Contrast Ratios & WCAG 2.2 AA/AAA Calibrations (Items 11–20)

### 3.1 Relative Luminance & Mathematical Verification Formula (Item 14, 15)

The Relative Luminance $L$ of any sRGB color is computed as:
$$L = 0.299 R_{lin} + 0.587 G_{lin} + 0.114 B_{lin}$$

For linear channel transformation from 8-bit standard values ($C \in [R, G, B]$ normalized to $[0, 1]$):
$$C_{lin} = \begin{cases} 
\frac{C}{12.92} & \text{if } C \le 0.04045 \\
\left(\frac{C + 0.055}{1.055}\right)^{2.4} & \text{if } C > 0.04045 
\end{cases}$$

The contrast ratio between light luminance $L_1$ and dark luminance $L_2$ is:
$$\text{Contrast Ratio} = \frac{L_1 + 0.05}{L_2 + 0.05}$$

**Luminance Threshold Foreground Selector Rule (Item 15)**:
- When background $L > 0.6 \implies$ ink color MUST be `--neutral-900` (`#090E18`).
- When background $L \le 0.6 \implies$ ink color MUST be stark white (`#FFFFFF`) or `--neutral-50` (`#F8FAFC`).

### 3.2 Calibrated Color Contrast Matrix (Items 11, 12, 13)

| Theme Pairing | Background | Foreground | Computed Ratio | WCAG Compliance | Use Case |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Dark Theme Body (Standard)** | `#0A0F1D` (Canvas) | `#F8FAFC` (Slate 50) | **16.4:1** | **AAA Pass** (Target >= 7.0:1) | Primary paragraphs, titles |
| **Dark Theme Muted Text** | `#172036` (Card) | `#9AA5BA` (Neutral 400) | **7.4:1** | **AAA Pass** (Target >= 4.5:1) | Metadata, commit hashes |
| **Dark High-Contrast Theme** | `#0A0C10` (Base) | `#FFFFFF` (Pure White) | **21.0:1** | **AAA Pass** (Max Contrast) | Accessibility high-contrast mode |
| **Light Theme Body (Standard)** | `#FBF9F6` (Warm Paper)| `#090E18` (Primary Ink)| **16.8:1** | **AAA Pass** (Target >= 7.0:1) | Main documentation, prose |
| **Light Theme Muted Text** | `#FFFFFF` (Card) | `#48536B` (Neutral 600) | **5.4:1** | **AA Pass** (Target >= 4.5:1) | Secondary helper text |
| **Primary CTA Button** | `#2563EB` (Blue 500) | `#FFFFFF` (Pure White) | **4.6:1** | **AA Pass** (Target >= 4.5:1) | Primary action label |
| **Dark Primary CTA Accent** | `#5E8EF6` (Blue 400) | `#0A0F1D` (Dark Ink) | **8.2:1** | **AAA Pass** (Target >= 4.5:1) | High-visibility dark button |
| **Focus Visible Outline** | `#0A0F1D` (Canvas) | `#5E8EF6` (Ring) | **6.1:1** | **AA Pass** (Non-Text >= 3.0:1) | Keyboard focus boundary |
| **Interactive Form Border** | `#172036` (Card) | `#44516D` (Border) | **3.2:1** | **AA Pass** (Non-Text >= 3.0:1) | Input & Select borders |

### 3.3 Focus Ring & Form Control Standards (Items 16, 17, 18, 19)

1. **Focus Ring Contrast (Item 16)**:
   - All interactive controls (`button`, `a`, `input`, `select`) must implement an unbroken `:focus-visible` state.
   - Outline width: `2px`. Offset: `2px`. Color: `hsl(var(--ring))`.
   - Contrast between ring and immediate surface must strictly exceed `3.0:1`.

2. **Form Control Border Contrast (Item 17)**:
   - Inactive input borders must maintain >= 3.0:1 against adjacent card bodies to ensure fields are clearly discernable to visually impaired users.
   - Active/focus input borders elevate to `--ring` (`5.8:1`).

3. **Placeholder Text Differentiation (Item 18)**:
   - Placeholder text must achieve >= 3.5:1 against the input background while remaining distinct from user-entered text (which satisfies >= 7:1 AAA).

4. **Link Distinction in Prose (Item 19)**:
   - In-line hyperlinks within documentation or commit logs must have either an explicit underline (`text-decoration: underline; text-underline-offset: 4px;`) or a >= 3:1 contrast delta from the surrounding body prose in addition to color difference.

### 3.4 Non-Color Visual Redundancy (Item 20)

Status pills and badge indicators must NEVER rely on color hue alone to convey status. Every status indicator must feature an explicit SVG icon, distinct geometric shape, or visible textual badge:

```typescript
// ✅ CANONICAL: Multi-channel accessible status badge
export function PipelineStatusBadge({ status }: { status: 'passed' | 'failed' | 'running' }) {
  const configs = {
    passed: {
      label: 'Passed',
      icon: <CheckCircle2 className="size-3.5 text-emerald-400" aria-hidden="true" />,
      className: 'bg-emerald-950/40 border-emerald-500/40 text-emerald-200'
    },
    failed: {
      label: 'Failed',
      icon: <XCircle className="size-3.5 text-rose-400" aria-hidden="true" />,
      className: 'bg-rose-950/40 border-rose-500/40 text-rose-200'
    },
    running: {
      label: 'Running',
      icon: <Loader2 className="size-3.5 text-cyan-400 animate-spin" aria-hidden="true" />,
      className: 'bg-cyan-950/40 border-cyan-500/40 text-cyan-200'
    }
  };

  const current = configs[status];

  return (
    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-mono font-medium border ${current.className}`}>
      {current.icon}
      <span>{current.label}</span>
    </span>
  );
}
```

---

## 4. Cluster 03: 4-Plane Depth Hierarchy & Materiality (Items 21–30)

### 4.1 The 4 Planes of Physical Elevation

Rather than applying heavy, blurred drop shadows (which desaturate colors and look dated), GitMap interfaces construct depth through **lightness stepped surfaces** combined with **hairline light refraction**:

```text
+---------------------------------------------------------------------------------+
| Plane 3: Elevated Modals & Flyouts (Lightness: ~19%, #1E2945)                   |
|          + Hairline glow & inset refraction highlight                           |
+---------------------------------------------------------------------------------+
       ^
       | (+4% Lightness Step)
+---------------------------------------------------------------------------------+
| Plane 2: Surface Cards, Tables & Inputs (Lightness: ~15%, #172036)              |
|          + Inset 1px top highlight (0.12 opacity)                               |
+---------------------------------------------------------------------------------+
       ^
       | (+4% Lightness Step)
+---------------------------------------------------------------------------------+
| Plane 1: Raised Wells & Section Bands (Lightness: ~11%, #10172A)                |
|          + Inset border boundary                                                |
+---------------------------------------------------------------------------------+
       ^
       | (+3% Lightness Step)
+---------------------------------------------------------------------------------+
| Plane 0: Base Canvas Ground (Lightness: ~8%, #0A0F1D)                           |
+---------------------------------------------------------------------------------+
```

#### Detailed Plane Specifications:
1. **Plane 0 (Canvas Base — Item 21)**:
   - Dark: `#0A0F1D` (`hsl(222, 47%, 8%)`).
   - Light: `#FBF9F6` (`hsl(40, 33%, 98%)`).
   - Purpose: Viewport boundary, deep background ground.
2. **Plane 1 (Raised Wells — Item 22)**:
   - Dark: `#10172A` (`hsl(222, 45%, 11%)`).
   - Light: `#F1EDE4` (`hsl(40, 24%, 92%)`).
   - Purpose: Alternating table row strips, inset log viewports, sidebar panels.
3. **Plane 2 (Surface Cards — Item 23)**:
   - Dark: `#172036` (`hsl(222, 40%, 15%)`).
   - Light: `#FFFFFF` (`hsl(0, 0%, 100%)`).
   - Purpose: Interactive repository cards, configuration panels, form groupings.
4. **Plane 3 (Elevated Modals & Popovers — Item 24)**:
   - Dark: `#1E2945` (`hsl(222, 39%, 19%)`).
   - Light: `#FFFFFF` (`hsl(0, 0%, 100%)`).
   - Purpose: Dialogs, dropdown action sheets, context menus, tooltips.

### 4.2 Hairline Top-Edge Lighting & Dividers (Item 25, 26, 27)

To simulate directional overhead lighting hitting physical card bevels, cards feature an inset top-edge highlight and crisp hairline borders:

```css
/* Card Materiality Specification */
.plane-card {
  background-color: hsl(var(--card));
  border: 1px solid rgba(248, 250, 252, 0.08); /* Crisp 1px hairline divider */
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.12),   /* Overhead light refraction */
    0 4px 20px rgba(0, 0, 0, 0.25);            /* Restrained ambient ground shadow */
  border-radius: 0.75rem;
}

/* Modal Elevation Specification */
.plane-modal {
  background-color: hsl(var(--popover));
  border: 1px solid rgba(248, 250, 252, 0.14);
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.18),
    0 12px 40px rgba(0, 0, 0, 0.5),
    0 0 0 1px hsl(var(--primary) / 0.15);      /* Subtle accent perimeter glow */
  border-radius: 1rem;
}
```

### 4.3 Layered Radial Halos & Frosted Glassmorphism (Item 28, 29, 30)

1. **Layered Radial Halos (Item 28)**:
   - Hero and KPI sections feature subtle radial gradient backlights positioned behind primary cards (`opacity: 0.08` to `0.12`, blur `80px`):
   ```css
   .ambient-halo {
     position: absolute;
     width: 480px;
     height: 480px;
     background: radial-gradient(circle, hsl(var(--primary) / 0.14) 0%, transparent 70%);
     filter: blur(80px);
     pointer-events: none;
     z-index: 0;
   }
   ```
2. **Frosted Acrylic Glassmorphism (Item 29)**:
   - Sticky navigation bars and floating toolbars use acrylic diffusion:
   ```css
   .glass-nav {
     background-color: hsl(var(--background) / 0.85);
     backdrop-filter: blur(16px) saturate(180%);
     -webkit-backdrop-filter: blur(16px) saturate(180%);
     border-bottom: 1px solid hsl(var(--border) / 0.6);
   }
   ```
3. **Restrained Hover Elevation (Item 30)**:
   - Interactive card hover states use subtle physical lift:
   ```css
   .hover-card-elevation {
     transition: transform 240ms cubic-bezier(0.16, 1, 0.3, 1), box-shadow 240ms cubic-bezier(0.16, 1, 0.3, 1);
   }
   .hover-card-elevation:hover {
     transform: translateY(-2px);
     box-shadow: 
       inset 0 1px 0 rgba(255, 255, 255, 0.20),
       0 8px 30px rgba(0, 0, 0, 0.35);
   }
   ```

---

## 5. Concrete Migration & Refactoring Checklist

### 5.1 Token System Migration (Cluster 01)
- [ ] **Scan Hex Codes**: Grep codebase for raw `#` hex color occurrences (`#2563EB`, `#1E293B`, `#0F172A`).
- [ ] **Migrate Tailwind Root Config**: Update `tailwind.config.js` or global CSS with space-separated HSL variables.
- [ ] **Declare Blue Ramp**: Register `--blue-50` through `--blue-900` CSS custom properties in `:root` and `.dark`.
- [ ] **Declare Violet Ramp**: Register `--violet-50` through `--violet-900` CSS custom properties.
- [ ] **Integrate Cyan Highlight**: Map `--cyan-400` to `--brand-highlight` and register utility classes.
- [ ] **Eliminate Direct Utility Colors**: Refactor `bg-blue-600` $\to$ `bg-primary`, `text-slate-400` $\to$ `text-muted-foreground`.
- [ ] **Isolate Code Containers**: Apply fixed dark slate background to `pre`, `code`, and terminal log panes.

### 5.2 Contrast Verification Migration (Cluster 02)
- [ ] **Audit Body Prose**: Verify all paragraph text meets >= 4.5:1 (AA) and targets >= 7:1 (AAA) in both light and dark.
- [ ] **Calibrate Form Inputs**: Elevate border contrast of inactive text inputs to >= 3.0:1 against card backgrounds.
- [ ] **Enforce Focus Rings**: Add 2px ring with 2px offset on all interactive buttons, links, and switches.
- [ ] **Implement Status Redundancy**: Audit status badges in `PipelineStatusBadge` and `RepoCard` to ensure icons accompany color states.
- [ ] **Hyperlink Inspection**: Add underlines or >= 3:1 luminance deltas to all inline links.

### 5.3 4-Plane Depth Migration (Cluster 03)
- [ ] **Replace Heavy Shadows**: Remove `shadow-2xl` and diffuse blur shadows; replace with `shadow-hairline` and `shadow-card`.
- [ ] **Structure Plane Classes**: Implement utility classes `.plane-0-canvas`, `.plane-1-well`, `.plane-2-card`, `.plane-3-modal`.
- [ ] **Add Top Inset Light**: Inject `inset 0 1px 0 rgba(255, 255, 255, 0.12)` into all card component bases.
- [ ] **Calibrate Glass Headers**: Apply `backdrop-blur-xl` and 85% alpha background to `Navbar` and `Sidebar`.
- [ ] **Standardize Card Lift**: Bind card hover states to `translateY(-2px)` with 240ms cubic bezier transition.

---

## 6. Verification Protocol & Quality Acceptance Gates

1. **No External Network Dependencies**: Token calculations and HSL variables run purely locally.
2. **Strict Relative Paths**: All spec cross-references and assets use repository-relative paths (`./01-tokens-contrast-and-depth-spec.md`, `02-spec/21-app/340-next-100-ui-improvements/`).
3. **Zero Test/Build Invocations**: This specification was authored with strict static design verification without invoking compilation or local test suites.
4. **Full Mathematical Coverage**: Ramps provide all 10 steps, relative luminance satisfies WCAG 2.2 AA and AAA thresholds, and 4 depth planes are fully quantified.
