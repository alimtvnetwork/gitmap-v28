# 23 — Building Block Components & Visual Anatomy

> **/goal** Deconstruct modern interface components into fundamental building blocks, explaining visual hierarchy, spacing, and hover mechanics in clear, foundational terms with standalone SVGs and modular LESS styles.
> **/learn** Master the component anatomy of curriculum cards, 4-cell sub-grids, interactive rows, and pricing tables, understanding how micro-surfaces assemble into cohesive layouts.

**Version:** 1.0.0
**Updated:** 2026-09-24
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Explaining Building Blocks Simply (Foundational Guide)

To help an AI understand interface construction without confusion, think of a modern card component as building with children's toy blocks. Every card is an organized toy shelf with specific pieces placed in a thoughtful order:

```text
+--------------------------------------------------------------+
| 1. The Hat (Category Badge Pill)                             |
|    - Small colorful oval at the top                          |
|    - Tells you what kind of lesson or module is inside       |
+--------------------------------------------------------------+
| 2. The Big Name Tag (Heading 700 Weight)                     |
|    - Large, strong, high-contrast dark letters               |
|    - Makes the subject crystal clear in 2 seconds            |
+--------------------------------------------------------------+
| 3. The Story (Subtitle & Description)                        |
|    - Softer slate color text                                 |
|    - Explains what you will build and achieve                |
+--------------------------------------------------------------+
| 4. The Toy Cubbies (2x2 Micro Sub-Grid)                      |
|    - 4 smaller boxes resting neatly inside the big box       |
|    - Each has an index circle (01, 02), title, and details   |
+--------------------------------------------------------------+
| 5. The Magic Hover Touch (Darkish Lift & Shade)              |
|    - When your cursor hovers, the card gently lifts 4px      |
|    - A darkish soft shade blossoms underneath                |
|    - It never disappears into the white background           |
+--------------------------------------------------------------+
| 6. The Doorway (Action Button / Arrow Link)                  |
|    - The bottom path inviting the learner to step forward    |
+--------------------------------------------------------------+
```

---

## 2. Standalone SVG Diagrams

The visual architecture is detailed in standalone SVG files within this repository:

- **[Card Anatomy & Subgrid Diagram](01-svg/card-anatomy.svg):** Deconstructs the category badge, title, 4-block inner grid, and hairline borders.
- **[Hover Mechanics & Darkish Shade Diagram](01-svg/hover-elevation.svg):** Compares the resting container against the elevated hover container and illustrates non-white-blended interactive rows.
- **[Curriculum Card Architecture](01-svg/curriculum-card.svg):** Detailed multi-module syllabus card layout in dark obsidian navy mode.
- **[3-Tier Pricing Table Diagram](01-svg/pricing-tiers.svg):** Illustrates optical weight hierarchy, center card elevation, feature bullets, and gold badge.

---

## 3. Core Component Catalog & LESS Implementations

LESS is the preferred styling framework across this repository due to nested block scoping and modular variables.

### 3.1 Component 1: The Curriculum Card with 4-Cell Subgrid

```less
// ============================================================================
// CURRICULUM CARD WITH 4-CELL SUBGRID (LESS PREFERRED)
// ============================================================================

.curriculum-module-card {
  display: flex;
  flex-direction: column;
  padding: 2rem;
  background-color: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 1.25rem;
  box-shadow: 0 4px 6px -1px rgba(15, 23, 42, 0.05);
  transition: transform 300ms cubic-bezier(0.16, 1, 0.3, 1),
              box-shadow 300ms cubic-bezier(0.16, 1, 0.3, 1),
              border-color 200ms ease;

  &:hover {
    transform: translateY(-4px);
    border-color: rgba(99, 102, 241, 0.4);
    box-shadow: 0 20px 40px -15px rgba(15, 23, 42, 0.12);
  }

  // 1. The Hat
  .module-badge {
    align-self: flex-start;
    padding: 0.25rem 0.75rem;
    border-radius: 9999px;
    background-color: #e0e7ff;
    color: #4338ca;
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.05em;
    text-transform: uppercase;
    margin-bottom: 1rem;
  }

  // 2. The Big Name Tag
  .module-title {
    font-size: 1.5rem;
    font-weight: 700;
    color: #0f172a;
    line-height: 1.25;
    margin-bottom: 0.5rem;
  }

  // 3. The Story
  .module-description {
    font-size: 0.9375rem;
    color: #64748b;
    line-height: 1.5;
    margin-bottom: 1.75rem;
  }

  // 4. The Toy Cubbies (2x2 Grid)
  .subgrid-container {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
    margin-bottom: 1.75rem;

    @media (max-width: 640px) {
      grid-template-columns: 1fr;
    }

    .subgrid-cell {
      padding: 1rem;
      background-color: #f8fafc;
      border: 1px solid #f1f5f9;
      border-radius: 0.75rem;
      transition: background-color 150ms ease, border-color 150ms ease;

      &:hover {
        background-color: #f1f5f9;
        border-color: #cbd5e1;
      }

      .cell-number {
        display: inline-block;
        font-family: monospace;
        font-size: 0.75rem;
        font-weight: 700;
        color: #6366f1;
        margin-bottom: 0.25rem;
      }

      .cell-title {
        font-size: 0.875rem;
        font-weight: 600;
        color: #1e293b;
        margin-bottom: 0.25rem;
      }

      .cell-desc {
        font-size: 0.75rem;
        color: #64748b;
        line-height: 1.4;
      }
    }
  }

  // 6. The Doorway
  .module-footer-link {
    display: inline-flex;
    align-items: center;
    gap: 0.375rem;
    font-size: 0.875rem;
    font-weight: 600;
    color: #6366f1;
    text-decoration: none;
    margin-top: auto;
    transition: gap 200ms ease;

    &:hover {
      gap: 0.625rem;
      color: #4f46e5;
    }
  }
}
```

---

### 3.2 Component 2: Interactive List Row with Darkish Shade

```less
// ============================================================================
// INTERACTIVE LIST ROW (LESS PREFERRED)
// ============================================================================

.interactive-syllabus-list {
  display: flex;
  flex-direction: column;
  border-radius: 0.75rem;
  overflow: hidden;
  border: 1px solid #e2e8f0;

  .list-row {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem 1.5rem;
    background-color: #ffffff;
    border-bottom: 1px solid #e2e8f0;
    transition: background-color 150ms ease, padding-left 250ms cubic-bezier(0.16, 1, 0.3, 1);

    &:last-child {
      border-bottom: none;
    }

    // Left accent pill sliding indicator
    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 25%;
      bottom: 25%;
      width: 3.5px;
      border-radius: 9999px;
      background-color: #f43f5e;
      opacity: 0;
      transform: scaleY(0.4);
      transition: opacity 150ms ease, transform 250ms cubic-bezier(0.16, 1, 0.3, 1);
    }

    // Darkish hover state preserves dark tone without washing out into white
    &:hover {
      background-color: rgba(15, 23, 42, 0.05);
      padding-left: 1.875rem;

      &::before {
        opacity: 1;
        transform: scaleY(1);
      }

      .row-title {
        color: #0f172a;
        font-weight: 600;
      }
    }

    .row-meta {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      font-size: 0.75rem;
      color: #64748b;
    }
  }
}
```

---

### 3.3 Component 3: The 3-Tier Modular Pricing Table

The pricing table uses optical weight differentiation: the recommended/popular plan floats slightly higher with an obsidian navy body and glowing accent, while surrounding plans rest on clean light surfaces.

```less
// ============================================================================
// 3-TIER MODULAR PRICING TABLE (LESS PREFERRED)
// ============================================================================

.pricing-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
  align-items: center;

  @media (max-width: 1024px) {
    grid-template-columns: 1fr;
    gap: 2rem;
  }

  // Standard Resting Tier Card
  .pricing-tier-card {
    position: relative;
    padding: 2.25rem 1.75rem;
    background-color: #ffffff;
    border: 1px solid #e2e8f0;
    border-radius: 1.25rem;
    box-shadow: 0 4px 6px -1px rgba(15, 23, 42, 0.05);
    transition: transform 300ms cubic-bezier(0.16, 1, 0.3, 1),
                box-shadow 300ms cubic-bezier(0.16, 1, 0.3, 1);

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 20px 40px -15px rgba(15, 23, 42, 0.12);
    }

    .tier-title {
      font-size: 1.25rem;
      font-weight: 800;
      color: #0f172a;
      margin-bottom: 0.375rem;
    }

    .tier-desc {
      font-size: 0.8125rem;
      color: #64748b;
      margin-bottom: 1.5rem;
    }

    .tier-price {
      font-size: 2.25rem;
      font-weight: 900;
      color: #0f172a;
      line-height: 1;
      margin-bottom: 1.5rem;

      .period {
        font-size: 0.875rem;
        font-weight: 500;
        color: #64748b;
      }
    }

    .feature-list {
      list-style: none;
      padding: 0;
      margin: 0 0 2rem 0;

      li {
        display: flex;
        align-items: center;
        gap: 0.625rem;
        font-size: 0.875rem;
        color: #334155;
        margin-bottom: 0.75rem;

        &::before {
          content: '✓';
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 1.125rem;
          height: 1.125rem;
          border-radius: 9999px;
          background-color: rgba(16, 185, 129, 0.15);
          color: #10b981;
          font-weight: 700;
          font-size: 0.6875rem;
        }
      }
    }

    .tier-btn {
      display: block;
      width: 100%;
      padding: 0.75rem 1rem;
      border-radius: 0.75rem;
      text-align: center;
      font-weight: 700;
      font-size: 0.875rem;
      text-decoration: none;
      background-color: #f1f5f9;
      color: #0f172a;
      border: 1px solid #cbd5e1;
      transition: background-color 150ms ease;

      &:hover {
        background-color: #e2e8f0;
      }
    }

    // Highlighted / Most Popular Tier
    &.tier-popular {
      background-color: #0b1329;
      border: 2px solid #8b5cf6;
      box-shadow: 0 20px 45px -10px rgba(139, 92, 246, 0.35);
      transform: scale(1.04);
      z-index: 2;

      &:hover {
        transform: scale(1.04) translateY(-4px);
        box-shadow: 0 25px 55px -10px rgba(139, 92, 246, 0.50);
      }

      .popular-badge {
        position: absolute;
        top: -12px;
        left: 50%;
        transform: translateX(-50%);
        padding: 0.25rem 0.875rem;
        border-radius: 9999px;
        background-color: #facc15;
        color: #0f172a;
        font-size: 0.6875rem;
        font-weight: 800;
        letter-spacing: 0.05em;
        text-transform: uppercase;
      }

      .tier-title,
      .tier-price {
        color: #ffffff;
      }

      .tier-desc,
      .tier-price .period {
        color: #94a3b8;
      }

      .feature-list li {
        color: #e2e8f0;

        &::before {
          background-color: rgba(139, 92, 246, 0.25);
          color: #a78bfa;
        }
      }

      .tier-btn {
        background-color: #8b5cf6;
        color: #ffffff;
        border: none;

        &:hover {
          background-color: #7c3aed;
        }
      }
    }
  }
}
```

---

### 3.4 Component 4: High-Trust Editorial Whiteness & Floating Containers

Inspired by clean, minimalist digital portals, this pattern pairs pure white canvases (`#ffffff`) with subtle hairline boundaries (`#f1f5f9`), delicate rose focus underlines (`#f43f5e`), and floating elevation bloom:

```less
// ============================================================================
// HIGH-TRUST EDITORIAL WHITENESS (LESS PREFERRED)
// ============================================================================

.editorial-section {
  background-color: #ffffff;
  padding: 5rem 1.5rem;

  .editorial-header {
    max-width: 680px;
    margin: 0 auto 3.5rem auto;
    text-align: center;

    .editorial-badge {
      display: inline-block;
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      background-color: #fff1f2;
      border: 1px solid #ffe4e6;
      color: #e11d48;
      font-size: 0.75rem;
      font-weight: 700;
      letter-spacing: 0.05em;
      text-transform: uppercase;
      margin-bottom: 1rem;
    }

    .editorial-title {
      font-size: clamp(2rem, 4vw, 3rem);
      font-weight: 800;
      color: #0f172a;
      line-height: 1.15;
      letter-spacing: -0.02em;

      .highlight-underline {
        position: relative;
        white-space: nowrap;

        &::after {
          content: '';
          position: absolute;
          left: 0;
          bottom: 2px;
          width: 100%;
          height: 3px;
          background: linear-gradient(90deg, #f43f5e 0%, #ec4899 100%);
          border-radius: 9999px;
        }
      }
    }
  }

  .floating-card-row {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;

    @media (max-width: 768px) {
      grid-template-columns: 1fr;
    }

    .floating-card {
      padding: 2.25rem;
      background-color: #ffffff;
      border: 1px solid #f1f5f9;
      border-radius: 1.5rem;
      box-shadow: 0 10px 30px -10px rgba(15, 23, 42, 0.04);
      transition: transform 300ms cubic-bezier(0.16, 1, 0.3, 1),
                  box-shadow 300ms cubic-bezier(0.16, 1, 0.3, 1);

      &:hover {
        transform: translateY(-6px);
        box-shadow: 0 25px 50px -12px rgba(15, 23, 42, 0.10);
      }
    }
  }
}
```

---

## 4. Guidelines for AI Construction

1. **Hierarchy First:** Always establish the container boundary before nesting internal content items.
2. **Predictable Ratios:** Card padding should scale proportionally (`1.5rem` on mobile, `2rem` to `2.5rem` on desktop).
3. **Darkish Shade Preservation:** Never use light gray `#fafafa` alone for light surface hover states; combine with low-alpha slate `rgba(15, 23, 42, 0.05)` to maintain visible dark contrast.
4. **Standalone SVGs:** Reference local SVGs (`01-svg/*.svg`) for visual documentation rather than external image URLs.
