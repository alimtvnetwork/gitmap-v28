# Design Principles

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

This document defines the visual philosophy and interaction feel that every component, page, and future extension must follow. These principles are non-negotiable — they define what makes the design system feel cohesive.

---

## 1. Variable-Driven, Never Hardcoded

**Rule:** Every visual property that could change with a theme (color, shadow color, border color, glow color) MUST reference a CSS custom property.

```css
/* ✅ CORRECT */
color: hsl(var(--primary));
border: 1px solid hsl(var(--border));
box-shadow: 0 0 8px hsl(var(--primary) / 0.15);

/* ❌ WRONG — hardcoded values */
color: #7c3aed;
border: 1px solid #333;
box-shadow: 0 0 8px rgba(124, 58, 237, 0.15);
```

**Exception:** Code block backgrounds (`hsl(220, 14%, 11%)`) are intentionally fixed to maintain readability across themes. These are documented in [09-code-blocks.md](./09-code-blocks.md).

---

## 2. HSL Color Model

All colors use **HSL (Hue, Saturation, Lightness)** format stored as space-separated values without the `hsl()` wrapper:

```css
--primary: 252 85% 60%;   /* NOT: hsl(252, 85%, 60%) */
```

This allows opacity modifiers via Tailwind/CSS:

```css
background: hsl(var(--primary) / 0.1);   /* 10% opacity primary */
```

**Why HSL:** Theme variants can be derived by adjusting lightness (dark mode = +5-10% lightness on primaries, -80% lightness on backgrounds) without changing hue.

---

## 3. Semantic Token Naming

Tokens are named by **purpose**, not by color:

| ✅ Semantic | ❌ Literal |
|------------|-----------|
| `--primary` | `--purple` |
| `--accent` | `--pink` |
| `--destructive` | `--red` |
| `--muted-foreground` | `--gray-text` |

This ensures re-theming works — changing `--primary` from purple to teal updates everything without renaming tokens.

---

## 4. Dark/Light Parity

Every token has two values:

```css
:root {
  --primary: 252 85% 60%;        /* Light mode */
}
.dark {
  --primary: 252 85% 65%;        /* Dark mode — slightly lighter */
}
```

Components NEVER check which theme is active. They always use `hsl(var(--primary))` — the cascade handles the rest.

**Dark mode adjustments:**

- Backgrounds: much darker (8-12% lightness)
- Foreground text: lighter (85-95% lightness)
- Primary/accent: +5% lightness for contrast on dark backgrounds
- Borders: darker (18-22% lightness)

---

## 5. CSS3 Motion Only

All animations and transitions use CSS3:

| Allowed | Forbidden |
|---------|-----------|
| `transition` property | `framer-motion` |
| `@keyframes` | `react-spring` |
| `transform` | `gsap` |
| `opacity` | `anime.js` |
| CSS `animation` property | JS `requestAnimationFrame` for visual effects |

**Why:** CSS3 animations are GPU-accelerated, don't block the main thread, work in any framework, and are portable to WordPress.

---

## 6. Interaction Feel

The design system aims for **subtle, responsive, professional** interactions:

| Property | Standard Value | Rationale |
|----------|---------------|-----------|
| Hover transition duration | `0.2s` | Fast enough to feel instant, slow enough to be noticed |
| Transform transitions | `0.15s` | Slightly faster for positional shifts |
| Easing | `ease` or `cubic-bezier(0.4, 0, 0.2, 1)` | Natural deceleration |
| Hover lift | `translateY(-1px)` to `translateY(-2px)` | Subtle depth, never dramatic |
| Hover scale | `scale(1.03)` to `scale(1.1)` | Gentle emphasis |
| Hover glow | `box-shadow` with `/ 0.1` to `/ 0.25` opacity | Soft halo, never harsh |

**Anti-patterns to avoid:**

- Bouncy or spring-like easing
- Hover effects that take longer than `0.3s`
- Transforms greater than `translateY(-4px)` or `scale(1.15)`
- Neon/harsh glows above `/ 0.3` opacity

---

## 7. Progressive Enhancement

Every interactive enhancement must degrade gracefully:

- **No hover?** Content remains fully readable and functional
- **No animation?** Respects `prefers-reduced-motion`
- **No JavaScript?** CSS handles all visual states (hover, focus, active)
- **No custom fonts?** System fonts provide acceptable fallback

```css
@media (prefers-reduced-motion: reduce) {
  * {
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
  }
}
```

---

## 8. Component Independence

Components must be self-contained visual units:

- No component should depend on a parent component's color or spacing
- Every component reads from the global token layer, not from sibling styles
- Components can be composed together without visual conflicts

---

## 9. Consistent State Language

All interactive elements follow the same state pattern:

| State | Visual Change |
|-------|--------------|
| **Default** | Base token values |
| **Hover** | Lightened background, subtle lift/glow, color shift toward primary |
| **Active/Pressed** | Slightly darker than hover, no lift |
| **Focus** | Ring outline using `--ring` token |
| **Disabled** | 50% opacity, no pointer events |
| **Selected/Active** | Primary accent color, stronger border or background |

---

## 10. Portability Mandate

The design system must work across:

| Platform | How |
|----------|-----|
| React + Tailwind | Native — current implementation |
| Static HTML | CSS custom properties + vanilla CSS classes |
| WordPress Classic Theme | CSS file import + PHP template classes |
| WordPress Block Theme | `theme.json` mapping + `style.css` |
| Any CSS framework | Custom properties are framework-agnostic |

No design decision should depend on React, Tailwind, or any specific build tool.

---

## 11. 60 / 30 / 10 Visual Balance & Von Restorff Rule

To prevent visual clutter and maintain clear conversion intent:

- **60% Dominant Neutral Base:** The deepest ground (Plane 0 or Plane 1) provides breathing room and contrast.
- **30% Structural Surfaces & Typography:** Content cards, tables, headers, and secondary text define the informational architecture.
- **10% Purposeful Accent:** The saturated brand accent is reserved strictly for primary actions, active states, and focal metrics.

### Single-Accent Von Restorff Principle
At any given viewport position, **exactly ONE prominent accent-filled element** may be visible (e.g. a solid Primary CTA). Secondary actions must use outline, ghost, or hairline styling. If two solid primary buttons are placed side-by-side, the visual hierarchy collapses into noise.

---

## 12. 4-Plane Depth Hierarchy (Materiality Over Shadows)

In modern dark and light interfaces, depth is created by stepped lightness planes and top-edge hairline highlights rather than dark drop shadows:

- **Plane 0 (Canvas Base):** Deepest page ground (`#0A0F1D` or `#0B0A09` in dark; `#FBF9F6` in light).
- **Plane 1 (Raised Wells):** Inset containers, alternating section bands, and sidebars (+3% to 4% lightness).
- **Plane 2 (Surface):** Cards, panels, input fields, and code blocks (+3% to 4% lightness).
- **Plane 3 (Elevated):** Modals, dropdowns, and hovered cards with subtle edge glow.

Adjacent sections alternate between Plane 0 and Plane 1 to establish rhythm.

---

## 13. Anti-AI-Slop Governance

All AI models generating UI code against this design system are strictly forbidden from producing:
1. **Generic 3-card grids** with circular icons and 1-line descriptions.
2. **Purple-blue gradient soup** across large background canvases.
3. **Unanchored heroes** lacking a tangible visual anchor (diagram, photo, or interactive scorecard).
4. **Adjective proof** ("cutting-edge", "world-class") in place of concrete numbers and metrics.
5. **Uniform section padding** that deprives the page of visual rhythm.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Theme Variable Registry | [03-theme-variable-architecture.md](./03-theme-variable-architecture.md) |
| Theme Catalogue & Palettes | [16-theme-catalogue-and-palettes.md](./16-theme-catalogue-and-palettes.md) |
| Machine-Readable Theme Tokens | [17-theme-tokens.json](./17-theme-tokens.json) |
| Dark Mode & Materiality | [18-dark-mode-and-materiality.md](./18-dark-mode-and-materiality.md) |
| Motion & Sliding Interactions | [19-modern-motion-and-sliding-interactions.md](./19-modern-motion-and-sliding-interactions.md) |
| AI Training & Anti-Slop Guide | [20-ai-training-and-checklist-guide.md](./20-ai-training-and-checklist-guide.md) |
| Motion System | [08-motion-transitions.md](./08-motion-transitions.md) |
| WordPress Migration | [15-wordpress-migration.md](./15-wordpress-migration.md) |
