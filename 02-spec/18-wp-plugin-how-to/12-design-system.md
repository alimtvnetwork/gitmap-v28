# Phase 12 — WordPress Plugin Design System

> **Created:** 2026-04-09
> **Status:** ✅ Active
> **Applies to:** All WordPress admin plugin pages

---

## 1. Overview

Every WordPress plugin admin page MUST follow a unified design system built on CSS custom properties, consistent typography, a shadow hierarchy, a shared animation library, and reusable component patterns. This specification extracts all conventions from the reference implementation and defines the rules an AI needs to produce visually consistent, production-quality plugin UI.

---

## 2. CSS Custom Properties (Design Tokens)

All colors, radii, shadows, and transitions are defined as CSS variables with `--riseup-` namespace prefix. Every property MUST include a hardcoded fallback value for environments where variables are not defined.

### 2.1 Token Registry

| Token | Fallback | Purpose |
|-------|----------|---------|
| `--riseup-primary` | `#1d4ed8` | Primary brand (buttons, links, active states) |
| `--riseup-primary-light` | `#3b82f6` | Lighter primary (gradients, hover accents) |
| `--riseup-primary-bg` | `#eff6ff` | Primary tinted background |
| `--riseup-primary-glow` | `rgba(59, 130, 246, 0.25)` | Focus ring / inset glow |
| `--riseup-success` | `#059669` | Success text / icon color |
| `--riseup-success-bg` | `#ecfdf5` | Success background tint |
| `--riseup-danger` | `#dc2626` | Error / destructive text |
| `--riseup-danger-bg` | `#fef2f2` | Error background tint |
| `--riseup-warning` | `#d97706` | Warning text |
| `--riseup-warning-bg` | `#fffbeb` | Warning background tint |
| `--riseup-text` | `#0f172a` | Primary text (headings, labels) |
| `--riseup-text-secondary` | `#475569` | Secondary text (descriptions, stats) |
| `--riseup-text-muted` | `#94a3b8` | Tertiary text (placeholders, legends) |
| `--riseup-bg` | `#f8fafc` | Surface background (cards, panels) |
| `--riseup-border` | `#e2e8f0` | Default border |
| `--riseup-border-strong` | `#cbd5e1` | Emphasized border |
| `--riseup-radius` | `8px` | Default border radius |
| `--riseup-radius-lg` | `12px` | Large border radius (cards, modals) |
| `--riseup-transition` | `0.2s` | Default transition duration |
| `--riseup-shadow-sm` | `0 1px 2px rgba(0, 0, 0, 0.05)` | Subtle shadow |
| `--riseup-shadow` | `0 1px 3px rgba(0, 0, 0, 0.08)` | Default shadow |
| `--riseup-shadow-md` | `0 4px 6px rgba(0, 0, 0, 0.08)` | Medium shadow |

### 2.2 Usage Pattern

```css
/* ✅ CORRECT — always include fallback */
.my-card {
    background: var(--riseup-bg, #f8fafc);
    border: 1px solid var(--riseup-border, #e2e8f0);
    border-radius: var(--riseup-radius, 8px);
}

/* ❌ WRONG — no fallback */
.my-card {
    background: var(--riseup-bg);
}

/* ❌ WRONG — hardcoded without variable */
.my-card {
    background: #f8fafc;
}
```

### 2.3 Naming Convention

- Prefix: `--riseup-` (replace with your plugin slug)
- Pattern: `--{slug}-{category}-{modifier}`
- Examples: `--riseup-primary`, `--riseup-text-muted`, `--riseup-shadow-md`

### 2.4 Slug Substitution Guide

All examples in this design system use the `riseup` prefix from the reference implementation. When building a new plugin, **replace every occurrence** with your plugin's kebab-case slug.

#### What to replace

| Category | Reference pattern | Your plugin (`my-tool`) |
|----------|------------------|------------------------|
| CSS custom properties | `--riseup-primary` | `--my-tool-primary` |
| CSS class prefixes | `.riseup-admin`, `.riseup-card` | `.my-tool-admin`, `.my-tool-card` |
| CSS file scoping | `.riseup-admin.riseup-agents` | `.my-tool-admin.my-tool-agents` |
| Keyframe names | `@keyframes riseupFadeIn` | `@keyframes myToolFadeIn` |
| JS localized objects | `window.RiseupErrors` | `window.MyToolErrors` |
| AJAX action prefixes | `riseup_dismiss_error_flash` | `my_tool_dismiss_error_flash` |

#### Derivation rules

| Plugin name | Kebab slug | CSS/property prefix | Class prefix |
|-------------|-----------|---------------------|--------------|
| Riseup Asia Uploader | `riseup` | `--riseup-` | `.riseup-` |
| QUpload | `qupload` | `--qupload-` | `.qupload-` |
| My Custom Plugin | `my-custom` | `--my-custom-` | `.my-custom-` |

#### How to derive your slug

1. Take the value of `PluginConfigType::Slug` (e.g., `'my-custom-plugin'`)
2. For CSS: use the slug directly as the prefix → `--my-custom-plugin-primary`
3. For classes: use the slug directly → `.my-custom-plugin-admin`
4. Alternatively, use a shortened form if the slug is long — define it once in your shared CSS and use consistently

#### ❌ Common mistakes

```css
/* ❌ WRONG — using reference prefix in a different plugin */
.riseup-admin { background: var(--riseup-bg, #f8fafc); }

/* ❌ WRONG — mixing prefixes */
.my-tool-admin { background: var(--riseup-bg, #f8fafc); }

/* ✅ CORRECT — consistent prefix throughout */
.my-tool-admin { background: var(--my-tool-bg, #f8fafc); }
```

---

## 3. Color System

### 3.1 Semantic Color Groups

Colors are organized into semantic groups defined in `data/colors.json` and managed by a `ColorGroupType` enum:

| Group | Purpose | Example keys |
|-------|---------|-------------|
| `logLevel` | Error log severity badges | Error `#dc3545`, Warn `#fd7e14`, Info `#0d6efd`, Debug `#6c757d` |
| `status` | Operation result states | success `#46b450`, error `#dc3232`, warning `#dba617` |
| `wpAdmin` | WordPress admin palette | primary `#2271b1`, headerAccent `#667eea` |

### 3.2 Gradient Patterns

Gradients follow a consistent `135deg` angle for backgrounds and `90deg` for horizontal bars:

```css
/* Status badge gradient — 135° diagonal */
background: linear-gradient(135deg, var(--riseup-success-bg, #ecfdf5), #d1fae5);

/* Progress bar gradient — 90° horizontal */
background: linear-gradient(90deg, var(--riseup-primary, #1d4ed8), var(--riseup-primary-light, #3b82f6));

/* Chart bar gradient — 180° vertical */
background: linear-gradient(180deg, var(--riseup-primary-light, #3b82f6), var(--riseup-primary, #1d4ed8));
```

### 3.3 Semantic Status Colors (Complete Table)

| Semantic State | Background | Text | Border | Use Case |
|----------------|------------|------|--------|----------|
| Success / Connected / Active | `#ecfdf5 → #d1fae5` | `#059669` | `#a7f3d0` | Connected agents, success badges |
| Error / Failed / Expired | `#fef2f2 → #fecaca` | `#dc2626` | `#fca5a5` | Error badges, expired licenses |
| Warning / Pending | `#fffbeb → #fef3c7` | `#d97706` | `#fbbf24` | Warning banners, pending states |
| Neutral / Inactive | `#f1f5f9 → #e2e8f0` | `#475569` | `#cbd5e1` | Inactive states, disabled badges |
| Info / Primary | `#dbeafe → #bfdbfe` | `#1e40af` | `#93c5fd` | Info badges, primary accents |
| Purple / Incremental | `#f3e5f5` | `#7b1fa2` | `#ce93d8` | Cron triggers, incremental snapshots |

### 3.4 Dark Mode & WordPress Admin Color Schemes

WordPress ships with multiple admin color schemes (Default, Light, Modern, Blue, Coffee, Ectoplasm, Midnight, Ocean, Sunrise). Plugins MUST adapt gracefully to at least the dark schemes to avoid unreadable UI.

#### Strategy: Override tokens per color scheme

WordPress adds a `body.admin-color-{scheme}` class. Override your design tokens for dark schemes:

```css
/* ── Dark scheme overrides ─────────────────────────────── */
body.admin-color-midnight .pluginname-admin,
body.admin-color-ectoplasm .pluginname-admin,
body.admin-color-coffee .pluginname-admin,
body.admin-color-ocean .pluginname-admin,
body.admin-color-sunrise .pluginname-admin {
    --pluginname-bg: #1e1e2f;
    --pluginname-text: #e2e8f0;
    --pluginname-text-secondary: #94a3b8;
    --pluginname-text-muted: #64748b;
    --pluginname-border: #334155;
    --pluginname-border-strong: #475569;
    --pluginname-shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.3);
    --pluginname-shadow: 0 1px 3px rgba(0, 0, 0, 0.4);
    --pluginname-shadow-md: 0 4px 6px rgba(0, 0, 0, 0.4);
}

/* Card backgrounds in dark mode */
body.admin-color-midnight .pluginname-admin .pluginname-card,
body.admin-color-ectoplasm .pluginname-admin .pluginname-card,
body.admin-color-coffee .pluginname-admin .pluginname-card,
body.admin-color-ocean .pluginname-admin .pluginname-card,
body.admin-color-sunrise .pluginname-admin .pluginname-card {
    background: #252540;
    border-color: var(--pluginname-border, #334155);
}

/* Table rows in dark mode */
body.admin-color-midnight .pluginname-admin .wp-list-table tbody tr,
body.admin-color-ectoplasm .pluginname-admin .wp-list-table tbody tr {
    background: transparent;
    color: var(--pluginname-text, #e2e8f0);
}

body.admin-color-midnight .pluginname-admin .wp-list-table tbody tr:nth-child(even) {
    background: rgba(255, 255, 255, 0.03);
}

/* Status colors remain unchanged — they use self-contained backgrounds */
```

#### Token override table

| Token | Light default | Dark override | Rationale |
|-------|-------------|---------------|-----------|
| `--pluginname-bg` | `#f8fafc` | `#1e1e2f` | Dark surface |
| `--pluginname-text` | `#0f172a` | `#e2e8f0` | Light-on-dark text |
| `--pluginname-text-secondary` | `#475569` | `#94a3b8` | Slightly muted |
| `--pluginname-text-muted` | `#94a3b8` | `#64748b` | Placeholder text |
| `--pluginname-border` | `#e2e8f0` | `#334155` | Subtle dark border |
| `--pluginname-border-strong` | `#cbd5e1` | `#475569` | Emphasized border |
| `--pluginname-primary` | `#1d4ed8` | `#3b82f6` | Slightly lighter for dark bg contrast |
| `--pluginname-primary-bg` | `#eff6ff` | `#1e293b` | Tinted dark surface |
| `--pluginname-danger-bg` | `#fef2f2` | `#3b1c1c` | Dark red tint |
| `--pluginname-success-bg` | `#ecfdf5` | `#1a2e22` | Dark green tint |
| `--pluginname-warning-bg` | `#fffbeb` | `#2e2a1a` | Dark amber tint |

#### Which schemes are "dark"

| Scheme | Type | Needs override |
|--------|------|---------------|
| Default | Light | No |
| Light | Light | No |
| Modern | Light | No |
| Blue | Light | No |
| Coffee | **Dark** | **Yes** |
| Ectoplasm | **Dark** | **Yes** |
| Midnight | **Dark** | **Yes** |
| Ocean | **Dark** | **Yes** |
| Sunrise | **Dark** | **Yes** |

#### Detection in PHP (for conditional asset loading)

```php
$isDarkScheme = in_array(
    get_user_option('admin_color'),
    ['midnight', 'ectoplasm', 'coffee', 'ocean', 'sunrise'],
    true,
);
```

This can be used to conditionally enqueue a `admin-dark.css` file instead of embedding all overrides in the shared CSS:

```php
if ($isDarkScheme) {
    wp_enqueue_style(
        $pluginSlug . '-admin-dark',
        plugin_dir_url(__FILE__) . 'assets/css/admin-dark.css',
        [$pluginSlug . '-admin-shared'],
        PluginConfigType::Version->value,
    );
}
```

#### Rules

1. Semantic status colors (success/error/warning badges) do NOT change — they have self-contained background+text contrast
2. Only surface, text, border, and shadow tokens need dark overrides
3. The `.code-pre` terminal block (§18.1) is already dark — no override needed
4. Test with Midnight scheme at minimum — it's the most common dark scheme
5. Use `rgba()` for dark mode shadows with higher opacity (0.3–0.4 vs 0.05–0.08)

### 3.5 Anti-Patterns

1. ❌ Never use hardcoded hex without a CSS variable reference
2. ❌ Never mix color systems (Material + Tailwind palettes in same badge)
3. ❌ Never use `opacity` on badge backgrounds — use explicit `rgba()` or gradient endpoints
4. ❌ Never place light text on light backgrounds

---

## 4. Typography

### 4.1 Font Stack

```css
/* Body text */
font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', sans-serif;

/* Monospace (code, timestamps, values, keys) */
font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;

/* Feedback / special pages */
font-family: 'Poppins', -apple-system, BlinkMacSystemFont, sans-serif;
```

### 4.2 Font Size Scale

| Size | Use Case | Weight |
|------|----------|--------|
| `9px` | Bar chart labels | 400 |
| `10px` | Summary labels, chart axis | 700 |
| `11px` | Badges, legends, descriptions, source files | 500–700 |
| `12px` | Filter labels, timestamps, progress meta, code blocks | 500–600 |
| `13px` | Body text, table cells, inline status, license keys | 500–600 |
| `14px` | Section headings, submit buttons, card headings | 600–700 |
| `20px` | Stat values (hero numbers) | 700 |
| `24px` | Modal close button | — |

### 4.3 Text Transform Conventions

| Element | Transform | Letter Spacing |
|---------|-----------|---------------|
| Filter labels | `uppercase` | `0.5px` |
| Stat labels | `uppercase` | `0.5px` |
| Section sub-headings | `uppercase` | `0.5px` |
| Badge text | `uppercase` (trigger/source badges) | `0.3px` |
| Date group headers | `none` | `0.3px` |

### 4.4 Font Weight Scale

| Weight | Use Case |
|--------|----------|
| `400` | Body text, descriptions |
| `500` | Secondary emphasis, filter labels, timestamps |
| `600` | Labels, inline status, badge text, tab text |
| `700` | Headings, stat values, badge counts, required markers |

---

## 5. Shadow Hierarchy

### 5.1 Elevation Levels

| Level | CSS | Use Case |
|-------|-----|----------|
| **sm** | `0 1px 2px rgba(0, 0, 0, 0.05)` | Subtle depth (slider values) |
| **default** | `0 1px 3px rgba(0, 0, 0, 0.08)` | Cards at rest, stat cards |
| **md** | `0 4px 6px rgba(0, 0, 0, 0.08)` | Hover elevation, storage cards |
| **lg** | `0 2px 12px rgba(0, 0, 0, 0.15)` | Flash banners, alert cards |
| **xl** | `0 20px 50px -12px rgba(0, 0, 0, 0.25)` | Modals |
| **2xl** | `0 25px 60px -12px rgba(0, 0, 0, 0.3)` | Fullscreen modals |
| **inset** | `inset 0 1px 3px rgba(0, 0, 0, 0.08)` | Progress bar tracks |
| **colored** | `0 2px 6px rgba(102, 126, 234, 0.3)` | Colored badge shadows |

### 5.2 Hover Shadow Pattern

Elements that elevate on hover MUST combine `translateY(-Npx)` with a shadow upgrade:

```css
.card {
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.card:hover {
    transform: translateY(-2px);
    box-shadow: var(--riseup-shadow-md, 0 4px 6px rgba(0, 0, 0, 0.08));
}
```

### 5.3 Focus Ring

All focusable elements use a blue ring for accessibility:

```css
.input:focus {
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
    outline: none;
}
```

---

## 6. Animation Library

### 6.1 Keyframe Registry

| Name | Effect | Duration | Easing | Use Case |
|------|--------|----------|--------|----------|
| `riseupFadeIn` | `opacity 0→1` | `0.25–0.3s` | `ease-out` | Overlays, tab content |
| `riseupFadeInUp` | `opacity 0→1 + translateY(12px→0)` | `0.3–0.5s` | `ease-out` | Table rows, cards, sections |
| `riseupScaleIn` | `opacity 0→1 + scale(0.95→1)` | `0.3–0.35s` | `cubic-bezier(0.34, 1.56, 0.64, 1)` | Modals, dialogs |
| `riseupPulse` | `scale(1→1.05→1)` | `1.5–3s` | `ease-in-out` | Live indicators, alert badges |
| `riseupSpin` | `rotate(0→360deg)` | `0.8–1.2s` | `linear` | Loading spinners |
| `riseupShimmer` | `background-position slide` | `1.5s` | `linear` | Progress bar shine effect |
| `livePulse` | `opacity + box-shadow pulse` | `1.5s` | — | Live dot indicator |

### 6.2 Keyframe Definitions

```css
@keyframes riseupFadeIn {
    from { opacity: 0; }
    to   { opacity: 1; }
}

@keyframes riseupFadeInUp {
    from { opacity: 0; transform: translateY(12px); }
    to   { opacity: 1; transform: translateY(0); }
}

@keyframes riseupScaleIn {
    from { opacity: 0; transform: scale(0.95); }
    to   { opacity: 1; transform: scale(1); }
}

@keyframes riseupPulse {
    0%, 100% { transform: scale(1); }
    50%      { transform: scale(1.05); }
}

@keyframes riseupSpin {
    0%   { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
}

@keyframes riseupShimmer {
    0%   { background-position: -200% 0; }
    100% { background-position: 200% 0; }
}

@keyframes livePulse {
    0%, 100% { opacity: 1; box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.4); }
    50%      { opacity: 0.7; box-shadow: 0 0 0 4px rgba(34, 197, 94, 0); }
}
```

### 6.3 Staggered Row Animation

Table rows MUST use staggered delays for entrance:

```css
.wp-list-table tbody tr {
    animation: riseupFadeInUp 0.4s ease-out both;
}
.wp-list-table tbody tr:nth-child(1) { animation-delay: 0s; }
.wp-list-table tbody tr:nth-child(2) { animation-delay: 0.06s; }
.wp-list-table tbody tr:nth-child(3) { animation-delay: 0.12s; }
.wp-list-table tbody tr:nth-child(4) { animation-delay: 0.18s; }
.wp-list-table tbody tr:nth-child(5) { animation-delay: 0.24s; }
```

### 6.4 Transition Standards

| Property | Duration | Easing | Use Case |
|----------|----------|--------|----------|
| `all` | `0.2s` | `cubic-bezier(0.4, 0, 0.2, 1)` | Default for interactive elements |
| `color` | `0.2s` | `ease` | Status text changes |
| `width` | `0.5s` | `cubic-bezier(0.4, 0, 0.2, 1)` | Progress bar fill |
| `transform` | `0.2s` | `cubic-bezier(0.34, 1.56, 0.64, 1)` | Bouncy icon scale (dashicons) |
| `background` | `0.2s` | `cubic-bezier(0.4, 0, 0.2, 1)` | Table row hover |

### 6.5 Animation Rules

1. Every page MUST re-declare the keyframes it uses (pages load CSS independently)
2. `infinite` animations are reserved for: spinners, live indicators, and alert pulses
3. Entrance animations use `both` fill mode for staggered delays
4. Hover animations are NEVER `infinite` — they are state transitions only

---

## 7. Badge System

### 7.1 Base Badge Anatomy

```css
.badge {
    display: inline-block;           /* or inline-flex with gap */
    padding: 3px 10px;
    border-radius: 20px;             /* pill shape, status badges */
    /* OR */
    border-radius: 4px;              /* tag shape, action/trigger badges */
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.3px;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
```

### 7.2 Badge Variants

| Variant | Shape | Background | Text | Border | Shadow |
|---------|-------|------------|------|--------|--------|
| **Status (pill)** | `border-radius: 20px` | Gradient `135deg` | Semantic color | `1px solid` lighter | Scale on hover |
| **Action (tag)** | `border-radius: 4px` | Translucent `rgba()` | Semantic dark | `1px solid rgba()` | `0 2px 4px rgba(0,0,0,0.1)` |
| **Trigger (tag)** | `border-radius: 4px` | Solid tint | Semantic dark | `1px solid` | `0 2px 4px rgba(0,0,0,0.12)` |
| **Method (tag)** | `border-radius: 6px` | Gradient `135deg` | Dark semantic | `1px solid` | None |
| **Count (pill)** | `border-radius: 10–12px` | Red gradient | White | None | Colored shadow |
| **Source (tag)** | `border-radius: 4px` | Dark `#1a1a2e` | White | None | `0 2px 6px rgba(0,0,0,0.2)` |

### 7.3 HTTP Method Badge Colors

| Method | Background | Text | Border |
|--------|------------|------|--------|
| `GET` | `#dcfce7 → #bbf7d0` | `#166534` | `#86efac` |
| `POST` | `#dbeafe → #bfdbfe` | `#1e40af` | `#93c5fd` |
| `PUT` | `#fef3c7 → #fde68a` | `#92400e` | `#fcd34d` |
| `DELETE` | `#fee2e2 → #fecaca` | `#991b1b` | `#fca5a5` |

### 7.4 Badge Hover Behavior

```css
.badge:hover {
    transform: scale(1.05);         /* status badges */
    /* OR */
    transform: scale(1.08);         /* level badges */
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}
```

---

## 8. Card Patterns

### 8.1 Stat Card

```css
.stat-card {
    flex: 1;
    min-width: 90px;
    background: var(--riseup-bg, #f8fafc);
    border: 1px solid var(--riseup-border, #e2e8f0);
    border-radius: var(--riseup-radius, 8px);
    padding: 14px 16px;
    text-align: center;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.stat-card:hover {
    border-color: var(--riseup-primary-light, #3b82f6);
    box-shadow: var(--riseup-shadow, 0 1px 3px rgba(0, 0, 0, 0.08));
    transform: translateY(-1px);
}
```

**Internal structure:**
- `.stat-value` — `font-size: 20px; font-weight: 700; font-family: monospace; color: primary`
- `.stat-label` — `font-size: 11px; uppercase; letter-spacing: 0.5px; font-weight: 600; color: muted`

### 8.2 Storage/Selection Card

```css
.selection-card {
    cursor: pointer;
    border: 2px solid var(--riseup-border, #e2e8f0);
    border-radius: var(--riseup-radius-lg, 12px);
    flex: 1;
    min-width: 160px;
    max-width: 220px;
    overflow: hidden;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.selection-card:hover {
    border-color: var(--riseup-primary-light, #3b82f6);
    box-shadow: var(--riseup-shadow-md);
    transform: translateY(-2px);
}
.selection-card.active {
    border-color: var(--riseup-primary, #1d4ed8);
    box-shadow: 0 0 0 1px var(--riseup-primary), var(--riseup-shadow-md);
    background: var(--riseup-primary-bg, #eff6ff);
}
```

### 8.3 File Viewer Card

- White background with `border-radius: 10px`
- Header: `background: #f8fafc; border-bottom: 1px solid #e2e8f0; padding: 14px 18px`
- Body: dark code panel `background: #0f172a; color: #e2e8f0; font-family: monospace`

### 8.4 Warning/Flash Banner Card

```css
.flash-banner {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 14px 18px;
    background: linear-gradient(135deg, #fffbeb, #fef3c7);
    border: 1px solid #fbbf24;
    border-left: 4px solid #f59e0b;
    border-radius: 10px;
    box-shadow: 0 2px 12px rgba(245, 158, 11, 0.15);
    animation: riseupFadeInUp 0.5s ease-out;
}
```

---

## 9. Button Variants

### 9.1 WordPress `.button` Enhancements

All WordPress `.button` elements inside `.riseup-admin` receive:

```css
.riseup-admin .button .dashicons {
    vertical-align: middle;
    margin-top: -2px;
    margin-right: 3px;
    transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.riseup-admin .button:hover .dashicons {
    transform: scale(1.15);
}
```

### 9.2 Primary Action Button

```css
.primary-action-btn {
    font-size: 14px;
    padding: 8px 24px;
    border-radius: 8px;
    background: linear-gradient(135deg, #1d4ed8, #3b82f6);
    border-color: #1d4ed8;
    box-shadow: 0 1px 3px rgba(29, 78, 216, 0.2);
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;
}
.primary-action-btn:hover {
    background: linear-gradient(135deg, #1e40af, #1d4ed8);
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(29, 78, 216, 0.3);
}
.primary-action-btn:active {
    transform: translateY(0) scale(0.98);
    box-shadow: 0 1px 2px rgba(29, 78, 216, 0.2);
}
```

### 9.3 Action Button (Table Row)

```css
.action-btn {
    margin-right: 4px;
    padding: 4px 10px;
    border-radius: 6px;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.action-btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}
```

### 9.4 Report / Pill Button

```css
.pill-btn {
    border-radius: 20px !important;
    padding: 4px 14px !important;
    font-weight: 500 !important;
}
.pill-btn:hover {
    background: #1d4ed8 !important;
    color: #fff !important;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(29, 78, 216, 0.25);
}
```

### 9.5 Dashicon Hover Behavior

Interactive dashicons inside buttons MUST scale on hover:

```css
.btn:hover .dashicons {
    transform: scale(1.15);              /* Standard */
    /* OR */
    transform: scale(1.2) rotate(-5deg); /* Playful (submit, report) */
}
```

---

## 10. Form Input Styling

### 10.1 Text Inputs & Textareas

```css
input[type="text"],
textarea {
    border-radius: 8px;
    border: 1px solid #e2e8f0;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
input[type="text"]:focus,
textarea:focus {
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
    outline: none;
}
```

### 10.2 Checkboxes

```css
input[type="checkbox"] {
    width: 18px;
    height: 18px;
    border-radius: 4px;
    border: 2px solid #e2e8f0;
    transition: all 0.2s;
    cursor: pointer;
}
input[type="checkbox"]:checked {
    background: #1d4ed8;
    border-color: #1d4ed8;
}
input[type="checkbox"]:hover {
    border-color: #3b82f6;
}
```

### 10.3 File Input

```css
input[type="file"] {
    padding: 8px;
    border: 2px dashed #e2e8f0;
    border-radius: 10px;
    background: #f8fafc;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    width: 100%;
    max-width: 500px;
}
input[type="file"]:hover {
    border-color: #3b82f6;
    background: #eff6ff;
}
```

### 10.4 Range Slider

```css
input[type="range"] {
    flex: 1;
    max-width: 300px;
    accent-color: var(--riseup-primary, #1d4ed8);
    height: 6px;
}
```

### 10.5 Settings Form Table

```css
.form-table tr {
    transition: background 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.form-table tr:hover {
    background: #f8fafc;
}
.form-table th {
    color: #0f172a;
    font-weight: 600;
}
.form-table .description {
    color: #64748b;
    font-style: normal;
    margin-top: 4px;
    font-size: 12px;
}
```

### 10.6 Required Field Marker

```css
.required {
    color: var(--riseup-danger, #dc2626);
    font-weight: 700;
}
```

---

## 11. Modal System

### 11.1 Standard Modal

Shared modal styles are defined in `admin-shared.css` and reused across all pages:

```
┌─────────────────────────────────────────────┐
│  .riseup-modal (fixed, z-index: 100000)     │
│  ┌─────────────────────────────────────────┐ │
│  │  .riseup-modal-overlay                  │ │
│  │  (blur backdrop, fadeIn 0.25s)          │ │
│  └─────────────────────────────────────────┘ │
│  ┌─────────────────────────────────────────┐ │
│  │  .riseup-modal-content                  │ │
│  │  (max-width: 600px, scaleIn 0.3s)       │ │
│  │  ┌───────────────────────────────────┐   │ │
│  │  │ .riseup-modal-header (bg: #f8fafc)│   │ │
│  │  │ ┌──────────┐  ┌─────────────────┐ │   │ │
│  │  │ │ h3 title │  │ × close button  │ │   │ │
│  │  │ └──────────┘  └─────────────────┘ │   │ │
│  │  ├───────────────────────────────────┤   │ │
│  │  │ .riseup-modal-body (scrollable)   │   │ │
│  │  ├───────────────────────────────────┤   │ │
│  │  │ .riseup-modal-footer (bg: #f8fafc)│   │ │
│  │  └───────────────────────────────────┘   │ │
│  └─────────────────────────────────────────┘ │
└─────────────────────────────────────────────┘
```

**Key properties:**
- Overlay: `rgba(15, 23, 42, 0.6)` + `backdrop-filter: blur(4px)`
- Content: `border-radius: 12px`, shadow xl, `max-height: 80vh`
- Close button: `32×32px`, rotates `90deg` on hover, turns red

### 11.2 Fullscreen Modal (Error Detail)

- `max-width: 1000px`, `max-height: 85vh`
- `border-radius: 14px`, shadow 2xl
- Overlay: `rgba(15, 23, 42, 0.65)` + `backdrop-filter: blur(6px)`
- `z-index: 100001` (above standard modal)
- Contains: summary bar → modal tabs → tab panes

### 11.3 Page-Specific Modal Overrides

Pages MAY override `max-width` on `.riseup-modal-content`:
- Agents page: `max-width: 800px`
- Error page: uses fullscreen variant

---

## 12. Tab System

### 12.1 Page-Level Tabs

```css
.nav-tab {
    border-radius: 8px 8px 0 0;
    padding: 8px 16px;
    font-weight: 500;
    position: relative;
    overflow: hidden;
}
/* Underline indicator */
.nav-tab::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 50%;
    width: 0;
    height: 2px;
    background: #1d4ed8;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    transform: translateX(-50%);
}
.nav-tab:hover::after { width: 80%; }
.nav-tab-active::after {
    width: 100%;
    background: linear-gradient(90deg, #1d4ed8, #3b82f6);
}
```

### 12.2 Modal Tabs

```css
.modal-tab {
    padding: 10px 18px;
    font-size: 13px;
    color: #64748b;
    border-bottom: 2px solid transparent;
    margin-bottom: -2px;
    font-weight: 500;
}
.modal-tab.active {
    color: #1d4ed8;
    border-bottom-color: #1d4ed8;
    font-weight: 600;
}
```

---

## 13. Table Patterns

### 13.1 Table Row Hover

```css
/* Standard hover */
tbody tr {
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
tbody tr:hover {
    background: #eff6ff !important;
}

/* Clickable row hover (logs) */
tbody tr.has-details:hover {
    background: #e8eeff !important;
    box-shadow: inset 3px 0 0 #667eea;
    cursor: pointer;
}
```

### 13.2 Date Group Headers

Tables with chronological data use date group separators:

```css
.date-group-header td {
    background: linear-gradient(135deg, #f8f9fb, #eef1f6) !important;
    border-top: 2px solid #667eea;
    padding: 10px 12px !important;
    font-size: 0;  /* hide td text, show label only */
}
.date-group-label {
    font-size: 13px;
    font-weight: 700;
    color: #1e2a4a;
}
.date-group-label::before {
    content: '📅';
}
```

### 13.3 Endpoint Group Headers

```css
.endpoint-group-header td {
    background: #f1f5f9 !important;
    font-size: 12px;
    color: #64748b;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    border-top: 2px solid #e2e8f0;
}
```

### 13.4 Nested Rows

```css
.nested-row td {
    background: var(--riseup-bg, #f8fafc) !important;
}
.nested-row:hover td {
    background: var(--riseup-primary-bg, #eff6ff) !important;
}
.nested-row td:first-child {
    border-left: 3px solid #7b1fa2;  /* purple accent */
}
```

---

## 14. Filter Bar

```css
.filters {
    margin-bottom: 16px;
    padding: 14px 16px;
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 10px;
}
.filters:focus-within {
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}
.filter-row {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    align-items: flex-end;
}
.filters label span {
    display: block;
    font-size: 11px;
    color: #475569;
    margin-bottom: 3px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}
```

---

## 15. Progress Bar

```css
.progress-bar-wrap {
    background: var(--riseup-border, #e2e8f0);
    border-radius: 10px;
    height: 22px;
    overflow: hidden;
    box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.08);
}
.progress-bar {
    background: linear-gradient(90deg, var(--riseup-primary), var(--riseup-primary-light));
    height: 100%;
    border-radius: 10px;
    transition: width 0.5s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;
}
/* Shimmer overlay */
.progress-bar::after {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
    background-size: 200% 100%;
    animation: riseupShimmer 1.5s linear infinite;
}
```

---

## 16. Live Indicator

```css
.live-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #cbd5e1;
    display: inline-block;
    transition: all 0.3s;
}
.live-dot.active {
    background: #22c55e;
    animation: livePulse 1.5s infinite;
}
```

---

## 17. Inline Status Text

```css
.inline-status {
    font-weight: 600;
    margin-left: 10px;
    font-size: 13px;
    transition: color 0.2s;
}
.inline-status.success { color: var(--riseup-success, #059669); }
.inline-status.error   { color: var(--riseup-danger, #dc2626); }
```

---

## 18. Code / Pre Blocks

### 18.1 Dark Terminal Style

```css
.code-pre {
    margin: 0;
    padding: 18px;
    background: #0f172a;
    color: #e2e8f0;
    font-size: 12px;
    line-height: 1.6;
    border-radius: 8px;
    white-space: pre-wrap;
    word-wrap: break-word;
    max-height: 400px;
    overflow: auto;
    font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
    border: 1px solid #1e293b;
}
```

### 18.2 Light Context Style (Stack Trace)

```css
.stack-trace {
    background: linear-gradient(135deg, #faf5ff, #f5f3ff);
    border: 1px solid #d8b4fe;
    border-radius: var(--riseup-radius, 8px);
    padding: 14px 16px;
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-size: 11px;
    line-height: 1.6;
    max-height: 240px;
    overflow: auto;
    white-space: pre-wrap;
    color: #581c87;
}
```

---

## 19. CSS File Organization

### 19.1 File Structure

```
assets/css/
├── admin-shared.css       # Modal system, shared keyframes
├── admin-settings.css     # Settings page, endpoint table
├── admin-snapshots.css    # Snapshot dashboard, charts, calendar
├── admin-logs.css         # Log table, badge colors
├── admin-errors.css       # Error log, file viewer, detail modal
├── admin-agents.css       # Agent management
├── admin-license.css      # License page
└── admin-feedback.css     # Feedback form
```

### 19.2 Loading Rules

1. `admin-shared.css` is loaded on ALL admin pages (contains modal, keyframes)
2. Page-specific CSS is loaded only on its respective page
3. Each page-specific CSS MUST re-declare any keyframes it uses (independent loading)
4. CSS class names use `riseup-` prefix for plugin-specific styles
5. WordPress native classes (`.button`, `.wp-list-table`, `.form-table`) are enhanced, not replaced

### 19.3 Specificity Rules

1. Plugin styles scope to `.riseup-admin` wrapper class
2. Page-specific overrides scope to `.riseup-admin.riseup-{page}` (e.g., `.riseup-admin.riseup-agents`)
3. Avoid `!important` except when overriding WordPress core table row backgrounds
4. Use class-based selectors, not ID selectors

---

## 20. Responsive Breakpoints

```css
/* Collapse grid layouts */
@media (max-width: 1100px) {
    .analytics-row { grid-template-columns: 1fr; }
}

/* General mobile adjustments handled by WordPress admin viewport */
```

---

## 21. Anti-Patterns (NEVER DO)

1. ❌ Use raw hex colors without a CSS variable reference
2. ❌ Use `opacity` on backgrounds for status colors — use explicit gradient endpoints
3. ❌ Use `!important` for anything other than WordPress core overrides
4. ❌ Create animations without the `riseup` prefix
5. ❌ Use `box-shadow` with large spread/blur for "glow" effects
6. ❌ Place light text on light backgrounds (contrast violation)
7. ❌ Use inline `<style>` blocks in templates — always use external CSS files
8. ❌ Skip fallback values in `var()` declarations
9. ❌ Use more than 5 staggered animation delays (performance)
10. ❌ Use `transition: all` with duration > `0.3s` (feels sluggish)
