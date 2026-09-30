# 21 — CSS3 Animations, Easing Curves & Professional Hover Effects

> **/goal** Specify high-craft CSS3 keyframe animations, transition curves, line hover highlights, and non-white-blended darkish hover elevations for web applications, design systems, and presentation slides.
> **/learn** Master the cubic-bezier physics curves, micro-interactions, darkish shade contrasts on light backgrounds, glow pulses, and accessible motion controls with LESS as the preferred stylesheet format.

**Version:** 1.0.0
**Updated:** 2026-09-24
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Principles of High-Craft Motion & Hover States

Web interfaces often suffer from washed-out hover effects where gray backgrounds blend into pure white, creating muddy or imperceptible visual feedback. Professional design systems solve this with three foundational rules:

1. **Non-White-Blended Darkish Tint:** Light surface hover tones must use a crisp slate tint (`rgba(15, 23, 42, 0.05)`) paired with high-contrast text (`#0f172a`), preserving immediate visibility.
2. **Directional Accent Indicators:** Interactive list rows and table items present a vertical accent pill (`3px` width, `border-radius: 9999px`) that expands or slides in on hover, signaling interactive readiness.
3. **Physics-Based Cubic-Bezier Curves:** Linear transitions feel robotic. Interfaces must employ emphasized deceleration curves (`cubic-bezier(0.16, 1, 0.3, 1)`) so elements snap smoothly into place.

---

## 2. Timing Functions & Easing Tokens

| Token Name | Cubic-Bezier Value | Duration | Purpose & Feel |
|:---|:---|:---|:---|
| `@ease-emphasized` | `cubic-bezier(0.16, 1, 0.3, 1)` | 300ms–400ms | Snappy entrance, card lift, modal reveal |
| `@ease-standard` | `cubic-bezier(0.4, 0, 0.2, 1)` | 200ms–250ms | Neutral opacity, color shifts |
| `@ease-exit` | `cubic-bezier(0.4, 0, 1, 1)` | 150ms–200ms | Dismissals, dropdown closing |

---

## 3. Darkish Line & Card Hover Mechanics (LESS Preferred)

LESS is the preferred styling framework across this repository due to parametric mixins and scoped color manipulation.

```less
// ============================================================================
// LESS IMPLEMENTATION (PREFERRED)
// ============================================================================

@color-ground-light: #ffffff;
@color-slate-900:    #0f172a;
@color-slate-500:    #64748b;
@color-accent-rose:  #f43f5e;
@color-accent-violet:#6366f1;
@ease-snap:          cubic-bezier(0.16, 1, 0.3, 1);

// Interactive list item with non-white-blended darkish tint
.interactive-curriculum-row {
  position: relative;
  display: flex;
  align-items: center;
  padding: 1rem 1.25rem;
  background-color: transparent;
  border-bottom: 1px solid #e2e8f0;
  transition: background-color 150ms ease,
              padding-left 250ms @ease-snap;

  // Sliding left indicator pill
  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 20%;
    bottom: 20%;
    width: 3.5px;
    border-radius: 9999px;
    background-color: @color-accent-rose;
    opacity: 0;
    transform: scaleY(0.3);
    transition: opacity 150ms ease,
                transform 250ms @ease-snap;
  }

  // Hover state: Crisp darkish shade, not muddy gray
  &:hover {
    background-color: rgba(15, 23, 42, 0.05);
    padding-left: 1.75rem;

    &::before {
      opacity: 1;
      transform: scaleY(1);
    }

    .row-title {
      color: @color-slate-900;
      font-weight: 600;
    }
  }
}

// Elevated card hover with shadow bloom
.interactive-card {
  background-color: @color-ground-light;
  border: 1px solid #e2e8f0;
  border-radius: 1rem;
  box-shadow: 0 4px 6px -1px rgba(15, 23, 42, 0.04);
  transition: transform 300ms @ease-snap,
              box-shadow 300ms @ease-snap,
              border-color 200ms ease;

  &:hover {
    transform: translateY(-4px);
    border-color: fade(@color-accent-violet, 40%);
    box-shadow: 0 20px 40px -15px rgba(15, 23, 42, 0.12);
  }
}
```

```css
/* ============================================================================
   RAW CSS EQUIVALENT
   ============================================================================ */

.interactive-curriculum-row {
  position: relative;
  display: flex;
  align-items: center;
  padding: 1rem 1.25rem;
  background-color: transparent;
  border-bottom: 1px solid #e2e8f0;
  transition: background-color 150ms ease, padding-left 250ms cubic-bezier(0.16, 1, 0.3, 1);
}

.interactive-curriculum-row::before {
  content: '';
  position: absolute;
  left: 0;
  top: 20%;
  bottom: 20%;
  width: 3.5px;
  border-radius: 9999px;
  background-color: #f43f5e;
  opacity: 0;
  transform: scaleY(0.3);
  transition: opacity 150ms ease, transform 250ms cubic-bezier(0.16, 1, 0.3, 1);
}

.interactive-curriculum-row:hover {
  background-color: rgba(15, 23, 42, 0.05);
  padding-left: 1.75rem;
}

.interactive-curriculum-row:hover::before {
  opacity: 1;
  transform: scaleY(1);
}
```

---

## 4. Keyframe Animations

### 4.1 Ambient Glow Pulse (Hero Cards & Live Status)

```less
@keyframes pulse-subtle-glow {
  0%, 100% {
    box-shadow: 0 0 0 0 rgba(139, 92, 246, 0.4);
  }
  50% {
    box-shadow: 0 0 0 8px rgba(139, 92, 246, 0);
  }
}

.status-badge-live {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 9999px;
  background-color: #10b981;
  animation: pulse-subtle-glow 2.5s infinite cubic-bezier(0.4, 0, 0.6, 1);
}
```

### 4.2 Smooth Shimmer Skeleton Loader

```less
@keyframes skeleton-shimmer {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}

.skeleton-box {
  background: linear-gradient(
    90deg,
    rgba(241, 245, 249, 1) 0%,
    rgba(226, 232, 240, 0.8) 50%,
    rgba(241, 245, 249, 1) 100%
  );
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.8s infinite linear;
}
```

### 4.3 Infinite Marquee Ribbon (Infinite Horizontal Ticker)

```less
// ============================================================================
// INFINITE MARQUEE TICKER (LESS PREFERRED)
// ============================================================================

@keyframes marquee-drift {
  0% {
    transform: translate3d(0, 0, 0);
  }
  100% {
    transform: translate3d(-50%, 0, 0);
  }
}

.marquee-container {
  overflow: hidden;
  display: flex;
  user-select: none;
  gap: 2rem;
  mask-image: linear-gradient(
    90deg,
    transparent 0%,
    rgba(0, 0, 0, 1) 10%,
    rgba(0, 0, 0, 1) 90%,
    transparent 100%
  );

  .marquee-track {
    display: flex;
    flex-shrink: 0;
    gap: 2rem;
    animation: marquee-drift 35s infinite linear;
    will-change: transform;

    &:hover {
      animation-play-state: paused;
    }
  }
}
```

### 4.4 Fluid Accordion Height Animation (Zero-JS)

Traditional web accordions required calculating pixel heights with JavaScript. Modern CSS allows pure declarative height animations through CSS Grid fractions or `interpolate-size`:

```less
// ============================================================================
// ZERO-JS FLUID ACCORDION (LESS PREFERRED)
// ============================================================================

.accordion-item {
  border-bottom: 1px solid #e2e8f0;

  .accordion-trigger {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.25rem 0;
    background: none;
    border: none;
    font-size: 1.125rem;
    font-weight: 600;
    color: #0f172a;
    cursor: pointer;

    .chevron-icon {
      width: 0.5rem;
      height: 0.5rem;
      border-right: 2px solid #64748b;
      border-bottom: 2px solid #64748b;
      transform: rotate(45deg);
      transition: transform 250ms cubic-bezier(0.16, 1, 0.3, 1);
    }
  }

  // Smooth height transition container
  .accordion-collapse {
    display: grid;
    grid-template-rows: 0fr;
    transition: grid-template-rows 300ms cubic-bezier(0.16, 1, 0.3, 1);

    .accordion-inner {
      overflow: hidden;
    }
  }

  // Expanded State
  &[aria-expanded="true"],
  &.is-open {
    .chevron-icon {
      transform: rotate(-135deg);
    }

    .accordion-collapse {
      grid-template-rows: 1fr;
    }
  }
}
```

### 4.5 Rotating Neon Border Gradient Sweep

```less
// ============================================================================
// CONTINUOUS ROTATING BORDER GRADIENT (LESS PREFERRED)
// ============================================================================

@property --gradient-angle {
  syntax: "<angle>";
  initial-value: 0deg;
  inherits: false;
}

@keyframes rotate-border {
  to {
    --gradient-angle: 360deg;
  }
}

.neon-sweep-card {
  position: relative;
  border-radius: 1.25rem;
  background: #0b1329;
  padding: 2rem;

  &::before {
    content: '';
    position: absolute;
    inset: -2px;
    border-radius: inherit;
    background: conic-gradient(
      from var(--gradient-angle),
      #8b5cf6 0%,
      #38bdf8 25%,
      transparent 50%,
      #8b5cf6 100%
    );
    z-index: -1;
    animation: rotate-border 6s linear infinite;
  }
}
```

### 4.6 Masked Text Reveal Grammar

```less
// ============================================================================
// MASKED TEXT ENTRANCE (LESS PREFERRED)
// ============================================================================

@keyframes text-reveal-up {
  0% {
    opacity: 0;
    transform: translateY(100%);
    clip-path: inset(0 0 100% 0);
  }
  100% {
    opacity: 1;
    transform: translateY(0%);
    clip-path: inset(0 0 0% 0);
  }
}

.reveal-heading {
  display: inline-block;
  overflow: hidden;

  .reveal-line {
    display: inline-block;
    animation: text-reveal-up 700ms cubic-bezier(0.16, 1, 0.3, 1) forwards;
  }
}
```

---

## 5. Accessibility & Reduced Motion

All transitions and animations must strictly respect the user's OS-level motion preference:

```less
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
```
