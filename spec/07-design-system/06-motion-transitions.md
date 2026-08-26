# Motion & Transitions

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

All motion in the design system uses **CSS3 only** — `transition`, `transform`, `@keyframes`, and `animation`. No JavaScript animation libraries. This ensures portability, GPU acceleration, and CMS compatibility.

---

## Core Transition Tokens

| Duration | Easing | Usage |
|----------|--------|-------|
| `0.15s` | `ease` | Transform shifts (translateX, translateY, scale) |
| `0.2s` | `ease` | Color changes, background changes, border changes, opacity |
| `0.3s` | `ease` | Box-shadow, filter, complex multi-property transitions |
| `0.3s` | `cubic-bezier(0.4, 0, 0.2, 1)` | Link underline sweep (Material-style decelerate) |

**Rule:** Never exceed `0.3s` for hover transitions. Users perceive > 300ms as sluggish.

---

## Transition Patterns

### 1. Color Transition

```css
transition: color 0.2s ease, background 0.2s ease;
```

Used on: paragraphs, list items, H4, tool buttons, checklist items.

**State flow:**
```
Default color → Hover: shift toward primary/foreground → Release: return
```

### 2. Lift + Glow

```css
transition: box-shadow 0.3s ease, transform 0.2s ease;
```

Used on: code block wrapper, inline code.

**State flow:**
```
Default (flat) → Hover: translateY(-2px) + box-shadow glow → Release: settle back
```

| Component | Lift | Glow |
|-----------|------|------|
| Code block wrapper | `translateY(-2px)` | `0 8px 32px hsl(--lang-accent / 0.1)` |
| Inline code | `translateY(-1px)` | `0 0 0 2px hsl(--highlight-glow / 0.15)` |

### 3. Slide / Nudge

```css
transition: transform 0.15s ease;
```

Used on: list items, checkboxes, blockquotes.

**State flow:**
```
Default → Hover: translateX(3px) → Release: return
```

### 4. Scale

```css
transition: transform 0.2s ease;
```

Used on: checkbox box, list bullet.

| Component | Scale |
|-----------|-------|
| Checkbox box | `scale(1.1)` on parent hover |
| List bullet | `scale(1.3)` on item hover |

### 5. Link Underline Sweep

```css
/* Default state */
transform: scaleX(0);
transform-origin: bottom right;
transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);

/* Hover state */
transform: scaleX(1);
transform-origin: bottom left;
```

Used on: `.spec-link::after` pseudo-element.

**Visual effect:** Underline sweeps from right to left on hover, retracts from left to right on un-hover.

The underline is a gradient:
```css
background: linear-gradient(
  90deg,
  hsl(var(--heading-gradient-from)),
  hsl(var(--heading-gradient-to))
);
height: 2px;
```

### 6. H3 Border Slide

```css
transition: padding-left 0.2s ease, border-color 0.2s ease, color 0.2s ease;
```

**State flow:**
```
Default: padding-left 0.65rem, border 50% opacity
Hover: padding-left 0.85rem, border full opacity, text becomes primary
```

### 7. Brightness/Saturation Boost

```css
transition: filter 0.3s ease;
/* Hover: */
filter: brightness(1.2) saturate(1.1);
```

Used on: H1 and H2 gradient headings.

---

## Keyframe Animations

### slideUpBar

```css
@keyframes slideUpBar {
  from { opacity: 0; transform: translateY(4px); }
  to   { opacity: 1; transform: translateY(0); }
}
```

Duration: `0.2s ease`. Used for: copy selection bar, selection label appearance.

### Accordion

```css
@keyframes accordion-down {
  from { height: 0; }
  to   { height: var(--radix-accordion-content-height); }
}
@keyframes accordion-up {
  from { height: var(--radix-accordion-content-height); }
  to   { height: 0; }
}
```

Duration: `0.2s ease-out`. Used for: sidebar expand/collapse.

---

## Button Slide Text Animation (CSS3)

For CTA buttons like "Join Us", use a **CSS3 text slide** instead of simple color hover:

```css
.slide-btn {
  position: relative;
  overflow: hidden;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.slide-btn .btn-text-default,
.slide-btn .btn-text-hover {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-btn .btn-text-default {
  transform: translateY(0);
}

.slide-btn .btn-text-hover {
  position: absolute;
  transform: translateY(100%);
}

.slide-btn:hover .btn-text-default {
  transform: translateY(-100%);
}

.slide-btn:hover .btn-text-hover {
  transform: translateY(0);
}
```

**Visual effect:** Default text slides up and out, hover text slides up from below. Pure CSS3, no JavaScript.

---

## Theme Transition

When switching light ↔ dark:

```css
body {
  transition: background-color 0.3s ease, color 0.3s ease;
}
```

All other elements inherit via CSS custom property changes — they update instantly (CSS variables don't animate, but the body transition provides a smooth background shift).

---

## Reduced Motion Support

```css
@media (prefers-reduced-motion: reduce) {
  * {
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
  }
}
```

---

## Fullscreen Transition

Code block fullscreen uses `position: fixed` with `inset: 2rem`:

```css
.code-block-wrapper.code-fullscreen {
  position: fixed !important;
  inset: 2rem;
  z-index: 999;
  border-radius: 1rem;
  transform: none !important;
  box-shadow: 0 25px 80px hsl(var(--lang-accent) / 0.25),
              0 0 0 1px hsl(var(--lang-accent) / 0.3);
}
```

The overlay backdrop uses `backdrop-filter: blur(4px)` with `z-index: 998`.

---

## Motion Anti-Patterns

| ❌ Forbidden | ✅ Alternative |
|-------------|---------------|
| Bounce easing | Use `ease` or `cubic-bezier(0.4, 0, 0.2, 1)` |
| Duration > 0.3s for hover | Keep at 0.2s-0.3s max |
| `translateY` > `-4px` | Use `-1px` to `-2px` for subtle lift |
| `scale` > `1.15` | Use `1.03` to `1.1` |
| Glow opacity > `0.3` | Use `0.1` to `0.25` |
| JavaScript animation for CSS-achievable effects | Use CSS3 transitions/keyframes |
| Animating `width`/`height` directly | Use `transform: scale()` or `max-height` |
| Animating `top`/`left` | Use `transform: translate()` |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Design Principles (Motion Rule) | [01-design-principles.md](./01-design-principles.md) |
| Button System (Slide Animation) | [09-button-system.md](./09-button-system.md) |
| Code Block Interactions | [07-code-blocks.md](./07-code-blocks.md) |
