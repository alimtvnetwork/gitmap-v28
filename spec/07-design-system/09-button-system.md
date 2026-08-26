# Button System

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

The button system defines default, hover, active, focus, and disabled states for all interactive buttons. All buttons use CSS3 transitions. The preferred CTA hover pattern is a **slide text animation**, not a simple color change.

---

## Button Variants

### Primary Button

```css
.btn-primary {
  background: hsl(var(--primary));
  color: hsl(var(--primary-foreground));
  font-family: 'Ubuntu', sans-serif;
  font-weight: 600;
  font-size: 0.875rem;
  padding: 0.5rem 1.25rem;
  border-radius: var(--radius);
  border: none;
  cursor: pointer;
  transition: background 0.2s ease, box-shadow 0.2s ease, transform 0.15s ease;
}

.btn-primary:hover {
  background: hsl(var(--primary) / 0.9);
  box-shadow: 0 4px 16px hsl(var(--primary) / 0.25);
  transform: translateY(-1px);
}

.btn-primary:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px hsl(var(--primary) / 0.15);
}

.btn-primary:focus-visible {
  outline: 2px solid hsl(var(--ring));
  outline-offset: 2px;
}

.btn-primary:disabled {
  opacity: 0.5;
  pointer-events: none;
}
```

### Secondary / Ghost Button

```css
.btn-ghost {
  background: transparent;
  color: hsl(var(--foreground));
  border: 1px solid hsl(var(--border));
  padding: 0.5rem 1.25rem;
  border-radius: var(--radius);
  transition: background 0.2s ease, border-color 0.2s ease, color 0.2s ease;
}

.btn-ghost:hover {
  background: hsl(var(--muted) / 0.5);
  border-color: hsl(var(--primary) / 0.3);
  color: hsl(var(--primary));
}
```

### Highlight Button (e.g., "Climate AI" Style)

A visually prominent button with gradient background and glow:

```css
.btn-highlight {
  background: linear-gradient(135deg, hsl(var(--primary)), hsl(var(--accent)));
  color: hsl(var(--primary-foreground));
  font-weight: 700;
  padding: 0.6rem 1.5rem;
  border-radius: var(--radius);
  border: none;
  position: relative;
  overflow: hidden;
  transition: box-shadow 0.3s ease, transform 0.15s ease;
}

.btn-highlight:hover {
  box-shadow: 0 6px 24px hsl(var(--primary) / 0.3),
              0 0 0 1px hsl(var(--accent) / 0.2);
  transform: translateY(-2px);
}

.btn-highlight::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, hsl(var(--accent)), hsl(var(--primary)));
  opacity: 0;
  transition: opacity 0.3s ease;
}

.btn-highlight:hover::before {
  opacity: 1;
}

.btn-highlight span {
  position: relative;
  z-index: 1;
}
```

**Usage:** Feature callout buttons, campaign CTAs, technology badges.

---

## Slide Text Animation (CTA Buttons)

**Replaces** simple hover color change on buttons like "Join Us". Uses CSS3 only.

### HTML Structure

```html
<button class="btn-primary slide-btn">
  <span class="btn-text-default">Join Us</span>
  <span class="btn-text-hover">Let's Go →</span>
</button>
```

### CSS

```css
.slide-btn {
  position: relative;
  overflow: hidden;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 8rem;
}

.slide-btn .btn-text-default,
.slide-btn .btn-text-hover {
  display: block;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1),
              opacity 0.2s ease;
}

.slide-btn .btn-text-default {
  transform: translateY(0);
  opacity: 1;
}

.slide-btn .btn-text-hover {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transform: translateY(100%);
  opacity: 0;
}

.slide-btn:hover .btn-text-default {
  transform: translateY(-100%);
  opacity: 0;
}

.slide-btn:hover .btn-text-hover {
  transform: translateY(0);
  opacity: 1;
}
```

**Visual flow:**
```
Default:  "Join Us" visible
Hover:    "Join Us" slides up, "Let's Go →" slides up from below
Un-hover: Reverse animation
```

---

## Tool Buttons (Code Block)

Smaller, low-emphasis buttons used in code block headers:

| Property | Value |
|----------|-------|
| Font size | `0.65rem` |
| Padding | `0.25rem 0.5rem` |
| Radius | `0.375rem` |
| Background | `hsl(220, 13%, 20%)` |
| Hover bg | `hsl(220, 13%, 28%)` |
| Hover border | `hsl(var(--lang-accent) / 0.4)` |
| Hover glow | `0 0 8px hsl(var(--lang-accent) / 0.15)` |
| Copied state bg | `hsl(152, 60%, 18%)` |
| Copied state color | `hsl(152, 70%, 60%)` |

---

## State Summary

| State | Visual Change | Transition |
|-------|--------------|-----------|
| Default | Base token values | — |
| Hover | Lighten bg, subtle lift, glow | `0.2s ease` |
| Active | Slightly darker, no lift | Instant |
| Focus | Ring outline (`--ring` token) | Instant |
| Disabled | 50% opacity, no pointer | — |
| Copied | Success color tint | `0.2s ease` |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Motion Transitions | [06-motion-transitions.md](./06-motion-transitions.md) |
| Header Navigation | [08-header-navigation.md](./08-header-navigation.md) |
| Theme Variables | [02-theme-variable-architecture.md](./02-theme-variable-architecture.md) |
