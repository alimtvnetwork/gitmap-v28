# Sweet Digs Cards & Sections Architecture

> **/goal** Master and enforce the thematic card patterns, ambient elevation hierarchies, and "Team Greenhouse" section DNA for the Sweet Digs Design System.
> **/learn** Master 16–20px rounded corners, compound card hover states, two-tone section headers, and staggered entry animations.

---

## 1. "Team Greenhouse" Section Layout DNA

The "Team Greenhouse" section is the foundational structural pattern for all thematic feature presentations across the Sweet Digs platform:

```html
<section class="section-greenhouse">
  <div class="section-container">
    <!-- Centered Two-Tone Header -->
    <header class="section-header">
      <span class="section-badge">Core Atmosphere Science</span>
      <h2 class="section-title">
        Understanding the Greenhouse <span class="text-primary">Effect</span>
      </h2>
      <p class="section-description">
        Key atmospheric gases trap solar thermal radiation, creating a habitable planetary climate while requiring balanced emission governance.
      </p>
    </header>

    <!-- Staggered Card Grid -->
    <div class="card-grid">
      <!-- Gas Card 1: Carbon Dioxide -->
      <article class="gas-card" style="--card-index: 1;">
        <div class="gas-card__header">
          <div class="gas-card__icon-box">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
            </svg>
          </div>
          <span class="gas-card__badge">76% Share</span>
        </div>
        <h3 class="gas-card__title">Carbon Dioxide (CO₂)</h3>
        <p class="gas-card__desc">Primary long-lived greenhouse gas produced by fossil fuel combustion, deforestation, and industrial processing.</p>
        <div class="gas-card__meta">
          <span class="meta-label">Atmospheric Lifetime</span>
          <span class="meta-value">100+ Years</span>
        </div>
      </article>

      <!-- Gas Card 2: Methane -->
      <article class="gas-card" style="--card-index: 2;">
        <div class="gas-card__header">
          <div class="gas-card__icon-box">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="12 2 2 7 12 12 22 7 12 2"></polygon>
              <polyline points="2 17 12 22 22 17"></polyline>
              <polyline points="2 12 12 17 22 12"></polyline>
            </svg>
          </div>
          <span class="gas-card__badge">16% Share</span>
        </div>
        <h3 class="gas-card__title">Methane (CH₄)</h3>
        <p class="gas-card__desc">Potent short-lived climate pollutant released during natural gas extraction, livestock digestion, and landfills.</p>
        <div class="gas-card__meta">
          <span class="meta-label">Global Warming Potential</span>
          <span class="meta-value">28x CO₂</span>
        </div>
      </article>

      <!-- Gas Card 3: Nitrous Oxide -->
      <article class="gas-card" style="--card-index: 3;">
        <div class="gas-card__header">
          <div class="gas-card__icon-box">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <path d="m4.93 4.93 4.24 4.24"></path>
              <path d="m14.83 9.17 4.24-4.24"></path>
              <path d="m14.83 14.83 4.24 4.24"></path>
              <path d="m9.17 14.83-4.24 4.24"></path>
            </svg>
          </div>
          <span class="gas-card__badge">6% Share</span>
        </div>
        <h3 class="gas-card__title">Nitrous Oxide (N₂O)</h3>
        <p class="gas-card__desc">Formed predominantly via synthetic agricultural nitrogen fertilizers and chemical manufacturing operations.</p>
        <div class="gas-card__meta">
          <span class="meta-label">Atmospheric Lifetime</span>
          <span class="meta-value">114 Years</span>
        </div>
      </article>
    </div>
  </div>
</section>
```

---

## 2. Section Structural DNA & Metrics

```css
.section-greenhouse {
  width: 100%;
  padding-block: 5rem;         /* 80px top & bottom */
  background-color: hsl(var(--background));
}

.section-container {
  max-width: 1152px;           /* 72rem / max-w-6xl */
  margin-inline: auto;
  padding-inline: 1.5rem;
}

.section-header {
  text-align: center;
  max-width: 680px;
  margin-inline: auto;
  margin-bottom: 3.5rem;       /* 56px spacing to card grid */
  animation: section-fade-in-up 700ms ease-out both;
}

.section-badge {
  display: inline-block;
  padding: 0.25rem 0.875rem;
  border-radius: 9999px;
  background-color: hsl(var(--secondary));
  color: hsl(var(--secondary-foreground));
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.75rem;
}

.section-title {
  font-size: 2.25rem;          /* 36px mobile, 44px md+ */
  font-weight: 800;
  color: hsl(var(--foreground));
  line-height: 1.2;
  letter-spacing: -0.02em;
  margin-bottom: 1rem;
}

.section-title .text-primary {
  color: hsl(var(--primary));
}

.section-description {
  font-size: 1.0625rem;        /* 17px */
  color: hsl(var(--muted-foreground));
  line-height: 1.6;
}
```

---

## 3. Gas Card Component Pattern (16–20px Rounded Geometry)

The Gas Card encapsulates atmospheric metrics using soft ambient shadows, a 16–20px border radius, and compound hover transitions:

```css
.card-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.75rem;
}

@media (min-width: 640px) {
  .card-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .card-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

.gas-card {
  position: relative;
  background-color: hsl(var(--card));
  border: 1px solid hsl(var(--border));
  border-radius: 1.25rem;      /* 20px curvature */
  padding: 2rem;
  box-shadow: 0 4px 12px -2px rgba(0, 0, 0, 0.04);
  transition: transform 300ms ease, box-shadow 300ms ease, border-color 300ms ease;
  animation: card-entry 700ms ease-out both;
  animation-delay: calc(var(--card-index, 1) * 120ms);
  overflow: hidden;
}

/* Subtle Gradient Overlay on Hover */
.gas-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    135deg,
    hsl(var(--primary) / 0.05) 0%,
    hsl(var(--secondary) / 0.3) 100%
  );
  opacity: 0;
  transition: opacity 300ms ease;
  pointer-events: none;
}

/* Compound Hover State */
.gas-card:hover {
  transform: translateY(-6px);
  border-color: hsl(var(--primary) / 0.4);
  box-shadow: 0 16px 32px -6px hsl(var(--primary) / 0.12),
              0 6px 12px -4px rgba(0, 0, 0, 0.04);
}

.gas-card:hover::before {
  opacity: 1;
}

/* Header & Micro-Elements */
.gas-card__header {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
}

.gas-card__icon-box {
  width: 48px;
  height: 48px;
  border-radius: 0.75rem;     /* 12px */
  background-color: hsl(var(--primary) / 0.12);
  color: hsl(var(--primary));
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 300ms ease;
}

.gas-card:hover .gas-card__icon-box {
  transform: scale(1.08) rotate(4deg);
}

.gas-card__badge {
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  background-color: hsl(var(--primary) / 0.12);
  color: hsl(var(--primary));
  font-size: 0.8125rem;
  font-weight: 600;
}

.gas-card__title {
  position: relative;
  z-index: 1;
  font-size: 1.25rem;
  font-weight: 700;
  color: hsl(var(--card-foreground));
  margin-bottom: 0.625rem;
}

.gas-card__desc {
  position: relative;
  z-index: 1;
  font-size: 0.875rem;
  color: hsl(var(--muted-foreground));
  line-height: 1.55;
  margin-bottom: 1.5rem;
}

.gas-card__meta {
  position: relative;
  z-index: 1;
  padding-top: 1rem;
  border-top: 1px solid hsl(var(--border) / 0.8);
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8125rem;
}

.meta-label {
  color: hsl(var(--muted-foreground));
}

.meta-value {
  font-weight: 600;
  color: hsl(var(--foreground));
}
```

---

## 4. Staggered Keyframe Animations

```css
@keyframes section-fade-in-up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes card-entry {
  from {
    opacity: 0;
    transform: translateY(24px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
```

---

## 5. Composition Checklist for AI Agents

1. **Card Curvature:** Always use `16px` to `20px` (`rounded-2xl` or `1.25rem`). Never use sharp 2px or 4px corners for primary cards.
2. **Elevated Shadows:** Shadows must incorporate low-alpha primary tint on hover rather than pure harsh black.
3. **Heading Split:** Ensure every thematic section header splits the concept into neutral text + primary keyword.
4. **Content Separation:** Spacing between header block and card grid must maintain 48px to 56px (`3rem` to `3.5rem`).
