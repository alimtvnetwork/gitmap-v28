# Sweet Digs Button & Animation System

> **/goal** Master and enforce the button architecture and CSS3 interaction models for the Sweet Digs Design System, emphasizing smooth color-switching, sliding text reveals, and glowing action elements without zoom distortion.
> **/learn** Master the ban on zoom-in hover scaling, the fluid water/color-switching gradient animation, the "Join Us" text slide mechanism, and the "Climate AI" pulse-glow floating action button.

---

## 1. Architectural Mandate: Buttons Without Zoom-In

> [!IMPORTANT]
> **TOTAL BAN ON ZOOM-IN (`transform: scale(1.05)`):**
> Standard buttons across the Sweet Digs system MUST NOT enlarge or zoom in on hover. Zoom-in animations create visual distortion, displace neighbouring elements, and trigger undesirable micro-reflows.
>
> Instead, hover emphasis is conveyed strictly through:
> 1. **Soft ambient shadow bloom** (`box-shadow: 0 6px 20px -2px hsl(var(--primary) / 0.40)`).
> 2. **CSS3 fluid water / color-switching animation** across the surface.
> 3. **Tactile micro-depression on active press** (`transform: scale(0.98)`).

---

## 2. Standard Primary Button with Fluid Water Animation

The standard primary button combines high-contrast typography, full pill curvature, and a dynamic fluid gradient sweep:

```html
<button class="btn btn-primary btn-water">
  <span class="btn-text">Get Started</span>
</button>
```

```css
.btn {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.625rem 1.5rem;    /* 10px vertical, 24px horizontal */
  font-size: 0.875rem;         /* 14px */
  font-weight: 600;
  border-radius: 9999px;       /* Pill shape */
  border: none;
  cursor: pointer;
  text-decoration: none;
  overflow: hidden;
  user-select: none;
  transition: box-shadow 300ms ease, background-position 400ms ease;
}

/* Primary Button Styling with Ambient Soft Shadow */
.btn-primary {
  background: hsl(var(--primary));
  color: hsl(var(--primary-foreground));
  box-shadow: 0 4px 14px 0 hsl(var(--primary) / 0.30);
}

/* CSS3 Fluid Water / Color-Switching Animation */
.btn-water {
  background-image: linear-gradient(
    135deg,
    hsl(var(--primary)) 0%,
    hsl(var(--primary) / 0.85) 25%,
    hsl(var(--secondary-foreground) / 0.9) 50%,
    hsl(var(--primary)) 75%,
    hsl(var(--primary) / 0.9) 100%
  );
  background-size: 250% 250%;
  background-position: 0% 50%;
}

.btn-water:hover {
  /* Smooth fluid color flow across the surface */
  background-position: 100% 50%;
  box-shadow: 0 8px 24px -2px hsl(var(--primary) / 0.45);
}

.btn-water:active {
  transform: scale(0.98);
  box-shadow: 0 2px 8px 0 hsl(var(--primary) / 0.30);
}
```

---

## 3. "Join Us" Button: CSS3 Sliding Text Reveal

The "Join Us" action button features an editorial vertical text reveal where an alternate phrase slides smoothly into position on hover:

```html
<button class="btn-slide" aria-label="Join our community">
  <span class="btn-slide__track">
    <span class="btn-slide__label btn-slide__label--default">Join Us</span>
    <span class="btn-slide__label btn-slide__label--hover" aria-hidden="true">Let's Go!</span>
  </span>
</button>
```

```css
.btn-slide {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 42px;
  padding-inline: 1.75rem;
  border-radius: 9999px;
  background-color: hsl(var(--primary));
  color: hsl(var(--primary-foreground));
  font-size: 0.875rem;
  font-weight: 600;
  border: none;
  cursor: pointer;
  overflow: hidden;
  box-shadow: 0 4px 14px 0 hsl(var(--primary) / 0.30);
  transition: box-shadow 300ms ease;
}

.btn-slide__track {
  display: flex;
  flex-direction: column;
  height: 100%;
  transition: transform 350ms cubic-bezier(0.4, 0, 0.2, 1);
  will-change: transform;
}

.btn-slide__label {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 42px;
  white-space: nowrap;
}

/* Hover: Inner track shifts upward exactly one label height */
.btn-slide:hover {
  box-shadow: 0 8px 22px -2px hsl(var(--primary) / 0.45);
}

.btn-slide:hover .btn-slide__track {
  transform: translateY(-100%);
}

.btn-slide:active {
  transform: scale(0.98);
}
```

### Sliding Interaction Kinetics

| Property | Value | Rationale |
|:---|:---|:---|
| **Height Boundary** | `42px` fixed | Prevents vertical reflows during translate |
| **Duration** | `350ms` | Calibrated for fluid text reveal velocity |
| **Easing Function** | `cubic-bezier(0.4, 0, 0.2, 1)` | Material deceleration curve for organic motion |
| **Transform Axis** | `translateY(-100%)` | Hardware-accelerated GPU translation |

---

## 4. "Climate AI" Glowing Floating Action Button (FAB)

The "Climate AI" FAB pins to the bottom-right corner of the interface, providing quick access to intelligent assistant workflows with a calming pulse-glow halo:

```html
<button class="fab-climate-ai" aria-label="Open Climate AI Assistant">
  <span class="fab-icon">
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
    </svg>
  </span>
  <span class="fab-text">Climate AI</span>
  <span class="fab-status-dot"></span>
</button>
```

```css
.fab-climate-ai {
  position: fixed;
  bottom: 1.5rem;              /* 24px from bottom */
  right: 1.5rem;               /* 24px from right */
  z-index: 50;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.25rem;
  border-radius: 9999px;
  background-color: hsl(var(--primary));
  color: hsl(var(--primary-foreground));
  font-size: 0.875rem;
  font-weight: 600;
  border: 1px solid hsl(var(--primary-foreground) / 0.2);
  cursor: pointer;
  animation: fab-pulse-glow 2.5s ease-in-out infinite, fab-fade-in 500ms ease-out both;
  transition: transform 200ms ease, box-shadow 200ms ease;
}

.fab-status-dot {
  width: 8px;
  height: 8px;
  border-radius: 9999px;
  background-color: #ffffff;
  box-shadow: 0 0 8px #ffffff;
}

/* Ambient Pulse Glow Keyframes */
@keyframes fab-pulse-glow {
  0%, 100% {
    box-shadow: 0 4px 14px 0 hsl(var(--primary) / 0.40),
                0 0 0 0 hsl(var(--primary) / 0.40);
  }
  50% {
    box-shadow: 0 8px 24px 0 hsl(var(--primary) / 0.50),
                0 0 24px 6px hsl(var(--primary) / 0.25);
  }
}

/* Entrance Fade */
@keyframes fab-fade-in {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fab-climate-ai:hover {
  /* No scale enlargement! Enhanced glow and light lift */
  transform: translateY(-2px);
  box-shadow: 0 10px 28px 0 hsl(var(--primary) / 0.55),
              0 0 28px 8px hsl(var(--primary) / 0.30);
}

.fab-climate-ai:active {
  transform: translateY(0) scale(0.98);
}
```

---

## 5. Reduced Motion Accessibility Contract

When users request reduced motion, all keyframe loops and transforms gracefully resolve to static high-contrast states:

```css
@media (prefers-reduced-motion: reduce) {
  .btn,
  .btn-water,
  .btn-slide,
  .btn-slide__track,
  .fab-climate-ai {
    animation: none !important;
    transition: none !important;
    transform: none !important;
  }
}
```
