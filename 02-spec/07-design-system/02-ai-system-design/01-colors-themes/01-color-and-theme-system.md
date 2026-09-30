# Color & Variable-Driven Theme Architecture

> **/goal** Provide a centralized, single-source-of-truth variable system that enables instantaneous theme swapping without modifying component-level styles.
> **/learn** Master the token definitions, CSS custom properties, LESS mixins, and theme conversion rules defined below.

## 🎯 Actionable CI/CD & AI Agent Checklist

- [ ] `/goal` Verify zero hardcoded hex codes inside component files; all values must reference CSS variables.
- [ ] `/learn` Verify that replacing the 5 core theme tokens converts the UI seamlessly from Pink/Red to Green or any custom palette.
- [ ] `/goal` Maintain WCAG 2.1 AA contrast ratio (>= 4.5:1 for normal text, >= 3.0:1 for large text) across all states.
- [ ] `/learn` Ensure LESS parametric mixins compile cleanly and mirror the CSS variable tokens.

---

## 1. Centralized Variable Architecture

At the root of the design system, all visual tokens are declared as CSS custom properties on `:root`. All child components, layouts, buttons, and modules MUST consume these variables.

### Global CSS Custom Properties (`theme.css`)

```css
:root {
  /* ------------------------------------------------------------- */
  /* Tier 1: Core Brand Accent Tokens (Default: Pink / Red)         */
  /* ------------------------------------------------------------- */
  --color-primary: #FF2D6F;          /* Strong pink-red brand accent */
  --color-primary-hover: #E02663;    /* 10% darkened interactive state */
  --color-primary-active: #C71F54;   /* Pressed feedback state */
  --color-primary-light: #FFE4EC;    /* 10% pastel wash for tags/chips */
  --color-primary-glow: rgba(255, 45, 111, 0.15); /* Focus shadow ring */

  /* ------------------------------------------------------------- */
  /* Tier 2: Hero Background Gradient Tokens                       */
  /* ------------------------------------------------------------- */
  --gradient-start: #F7DDE5;         /* Soft pastel rose tint */
  --gradient-end: #EADFD8;           /* Warm creamy beige tone */
  --hero-gradient: linear-gradient(135deg, var(--gradient-start) 0%, var(--gradient-end) 100%);

  /* ------------------------------------------------------------- */
  /* Tier 3: Neutral Surface & Border Tokens                       */
  /* ------------------------------------------------------------- */
  --color-bg-base: #F7F4F3;          /* Warm light gray page background */
  --color-surface: #FFFFFF;          /* Pure white card and container surface */
  --color-surface-translucent: rgba(255, 255, 255, 0.85); /* Sticky nav backdrop */
  --color-border: #E5E7EB;           /* Hairline neutral border (1px) */
  --color-border-subtle: #F3F4F6;    /* Divider and separator lines */

  /* ------------------------------------------------------------- */
  /* Tier 4: Typography & Content Tokens                           */
  /* ------------------------------------------------------------- */
  --color-text-primary: #111827;     /* High-contrast neutral dark text */
  --color-text-secondary: #6B7280;   /* Muted supporting body copy */
  --color-text-tertiary: #9CA3AF;    /* Form placeholders and inactive icons */
  --color-text-inverse: #FFFFFF;     /* Text on primary filled buttons */

  /* ------------------------------------------------------------- */
  /* Tier 5: Semantic State Tokens                                 */
  /* ------------------------------------------------------------- */
  --color-success: #10B981;          /* Active confirmation / online status */
  --color-warning: #F59E0B;          /* Caution / pending review badge */
  --color-danger: #EF4444;           /* Error border and alert text */
  --color-info: #3B82F6;             /* Supplemental system badges */

  /* ------------------------------------------------------------- */
  /* Tier 6: Elevation & Ambient Shadow Tokens                     */
  /* ------------------------------------------------------------- */
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  --shadow-card: 0 10px 30px rgba(0, 0, 0, 0.05);
  --shadow-card-hover: 0 18px 40px rgba(0, 0, 0, 0.08);
  --shadow-fab: 0 12px 32px rgba(255, 45, 111, 0.35);
  --shadow-fab-hover: 0 16px 40px rgba(255, 45, 111, 0.45);
}
```

---

## 2. Theme Conversion System: Pink/Red → Green Preset

To transform the complete website from the default **Pink/Red** theme to the **Green** theme, an AI or engineer simply substitutes the 5 core brand variables or applies the `[data-theme="green"]` attribute:

```css
[data-theme="green"] {
  /* Green Theme Preset */
  --color-primary: #22C55E;
  --color-primary-hover: #16A34A;
  --color-primary-active: #15803D;
  --color-primary-light: #DCFCE7;
  --color-primary-glow: rgba(34, 197, 94, 0.20);

  /* Green Hero Gradient */
  --gradient-start: #D1FAE5;
  --gradient-end: #ECFDF5;

  /* Floating Button Shadow Adaptation */
  --shadow-fab: 0 12px 32px rgba(34, 197, 94, 0.35);
  --shadow-fab-hover: 0 16px 40px rgba(34, 197, 94, 0.45);
}
```

### Visual Token Comparison Matrix

| Token Name | Default (Pink/Red) | Green Preset | Role Across All Pages |
|:---|:---|:---|:---|
| `--color-primary` | `#FF2D6F` | `#22C55E` | Main buttons, search action, active chips, FAB |
| `--color-primary-hover`| `#E02663` | `#16A34A` | Button hover state, icon hover accents |
| `--color-primary-light`| `#FFE4EC` | `#DCFCE7` | Active filter chip background, tag washes |
| `--gradient-start` | `#F7DDE5` | `#D1FAE5` | Hero background radial/linear origin |
| `--gradient-end` | `#EADFD8` | `#ECFDF5` | Hero background blend target |
| `--color-bg-base` | `#F7F4F3` | `#F7F4F3` | Universal page baseline background |
| `--color-surface` | `#FFFFFF` | `#FFFFFF` | Card, modal, and input background |
| `--color-text-primary` | `#111827` | `#111827` | Primary headings and high-contrast labels |

---

## 3. Parametric LESS Mixins (`theme-palette.less`)

Because LESS is explicitly supported and preferred for structured design tokens, parametric mixins mirror these CSS variables for preprocessed workflows:

```less
// -------------------------------------------------------------
// Parametric Palette Mixin
// -------------------------------------------------------------
.apply-theme-palette(@primary; @hover; @light; @grad-start; @grad-end) {
  --color-primary: @primary;
  --color-primary-hover: @hover;
  --color-primary-light: @light;
  --gradient-start: @grad-start;
  --gradient-end: @grad-end;
  --hero-gradient: linear-gradient(135deg, @grad-start 0%, @grad-end 100%);
}

// Default Theme Initialization
:root {
  .apply-theme-palette(
    #FF2D6F,
    #E02663,
    #FFE4EC,
    #F7DDE5,
    #EADFD8
  );
}

// Green Theme Override
[data-theme="green"] {
  .apply-theme-palette(
    #22C55E,
    #16A34A,
    #DCFCE7,
    #D1FAE5,
    #ECFDF5
  );
}
```

---

## 4. Typography Scale & Font Pairing Guidance

### Typography Tokens

```css
:root {
  /* Font Family Stacks */
  --font-heading: "Ubuntu", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  --font-body: "Ubuntu", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  --font-mono: "JetBrains Mono", ui-monospace, SFMono-Regular, monospace;

  /* Typography Scale */
  --text-hero-title: clamp(2.5rem, 5vw, 4rem);   /* 48px – 64px, weight 700 */
  --text-section-title: clamp(1.75rem, 3vw, 2.5rem); /* 28px – 40px, weight 700 */
  --text-card-title: 1.25rem;                     /* 20px, weight 700 */
  --text-body: 1rem;                              /* 16px, weight 400 */
  --text-subtext: 1.125rem;                       /* 18px, weight 400 */
  --text-btn: 0.9375rem;                          /* 15px, weight 600 */
  --text-chip: 0.875rem;                          /* 14px, weight 500 */
  --text-label: 0.75rem;                          /* 12px, weight 700, letter-spacing: 0.05em */
}
```

### Font Pairing Alternatives

When "Ubuntu" is not loaded in the client environment, fallback cleanly to:
- **Inter:** Ultra-clean, neutral, and precise for enterprise SaaS.
- **Poppins:** Geometric and rounded, matching the friendly aesthetic.
- **Manrope:** Modern semi-geometric, providing a refined, premium feel.

---

## 5. Rules for AI Generative Compliance

1. **Strict Variable Reference:** When generating any HTML, JSX, CSS, or Tailwind utility class, never write `#FF2D6F` or `#22C55E` directly. Use `var(--color-primary)` or Tailwind config alias `bg-brand-primary`.
2. **Predictable State Transitions:** Any element styled with `--color-primary` MUST use `--color-primary-hover` on `:hover` and `--color-primary-active` on `:active`.
3. **Contrast Integrity:** White text (`--color-text-inverse`) is reserved exclusively for solid `--color-primary` surfaces. On neutral surfaces, always use `--color-text-primary`.
