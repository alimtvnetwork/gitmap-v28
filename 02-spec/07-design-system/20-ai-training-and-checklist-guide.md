# 20 — AI Training, Anti-Slop & Design Execution Guide

> **/goal** Train any AI model (including lower-capability, blind, or prompt-following agents) to construct world-class, human-grade web interfaces without falling into generic AI template tropes.
> **/learn** Execute the non-negotiable 5-question above-the-fold contract, apply golden ratio asymmetry, enforce the single-accent Von Restorff law, eliminate all 10 AI-slop failure modes, and follow a deterministic step-by-step design sequence.

**Version:** 4.0.0
**Updated:** 2026-09-24
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. The Core Problem: Why AI Code Looks Like "AI Slop"

When an unguided AI assistant is asked to build a modern landing page or interface, it defaults to statistical averages:
- A purple-to-blue gradient on white.
- A centered hero headline with zero visual anchor.
- A generic 3-card grid with small circular outline icons and one-sentence descriptions.
- Adjectives like "revolutionary", "cutting-edge", and "seamless" instead of hard numbers.
- Identical padding and identical background color on every single section.

This specification exists to **retrain and constrain AI behavior** through deterministic rules, checklists, and step-by-step algorithms.

---

## 2. The 10-Point "Anti-AI-Slop" Rejection Rubric

If an AI generates any of the following 10 anti-patterns, the output is considered an **immediate failure and must be rejected and refactored**:

| # | Anti-Pattern Failure | What Amateur AI Does | Mandatory Professional Correction |
|---|:---|:---|:---|
| **1** | **Generic 3-Card Grid** | Three equal cards with a round icon and a vague sentence. | Replace with an asymmetric layout (e.g. 7+5 split, numbered chapters, or a controlled sliding carousel). |
| **2** | **Purple-Blue Gradient Soup** | Full-screen purple/blue gradient backgrounds. | Use a rich dark neutral base (`#0A0F1D` or `#0B0A09`) with hairlines and subtle 6–10% radial depth washes. |
| **3** | **Unanchored Hero** | Large centered text floating over a blank or gradient background. | Hero MUST carry a real visual anchor (e.g. interactive scorecard, product artefact, or progressive-blur photo) occupying ≥ 25% of desktop frame. |
| **4** | **Adjective-Laden Copy** | Using "cutting-edge", "world-class", "seamless", "synergy". | Replace adjectives with verifiable facts, hard numbers, timeframes, and concrete architectural guarantees. |
| **5** | **Uniform Section Rhythm** | Every section using identical padding (e.g. `py-16`) and identical plane. | Alternate between open sections (`160px`) and compressed bands (`48px`); alternate Plane 0 and Plane 1. |
| **6** | **Symmetrical Everything** | Stacking centered block after centered block down the page. | Require asymmetry in at least 4 sections (e.g. 7+5, 8+4, 5+7 column splits). Max 2 centered blocks per page. |
| **7** | **Stock Photography / AI People** | Photos of smiling people shaking hands at laptops. | Real authentic team photography, or replace with designed technical diagrams, scorecards, or UI artefacts. |
| **8** | **Accent Overuse** | Multiple primary colored buttons, colored headings, and icons all shouting at once. | Enforce 60/30/10 and single-accent Von Restorff: exactly ONE primary accent button per screen. |
| **9** | **Drop Shadows on Dark Grounds** | Heavy black `box-shadow` on dark backgrounds. | Replace shadows with 1px hairlines at 8–16% white opacity and top-edge chamfer highlights. |
| **10**| **Layout-Thrashing Motion** | Animating `width`, `height`, `top`, or `margin` on hover. | Animate ONLY GPU-composited `transform` and `opacity`. Hover lifts capped at `translateY(-2px)`. |

---

## 3. The 5-Question Above-the-Fold Contract

The very first viewport (before the user scrolls a single pixel) must answer **five fundamental questions**:

```text
┌────────────────────────────────────────────────────────────────────────┐
│ FIRST VIEWPORT CONTRACT (Desktop 1440×900 & Mobile 390×844)            │
├────────────────────────────────────────────────────────────────────────┤
│ 1. WHAT IS IT?                                                         │
│    Answered in the H1 headline (≤ 9 words, display scale).             │
├────────────────────────────────────────────────────────────────────────┤
│ 2. WHO IS IT FOR?                                                      │
│    Answered in the first clause of the subhead.                        │
├────────────────────────────────────────────────────────────────────────┤
│ 3. WHAT IS THE OUTCOME?                                                │
│    Answered in the second clause of the subhead (concrete metric).     │
├────────────────────────────────────────────────────────────────────────┤
│ 4. WHY SHOULD I BELIEVE IT?                                            │
│    Proof anchor in the first screen (hard stats, client logos, badge). │
├────────────────────────────────────────────────────────────────────────┤
│ 5. WHAT SHOULD I DO NOW?                                               │
│    Exactly ONE primary CTA button + ONE secondary ghost/outline action.│
└────────────────────────────────────────────────────────────────────────┘
```

**AI Verification:** Take a virtual viewport screenshot at 1440×900. If any of these 5 answers are missing or pushed below the fold, fix the spacing and layout immediately.

---

## 4. Composition: The Golden Ratio & Required Asymmetry

1. **The Golden Split (61.8% / 38.2%):**
   - On a 12-column grid, use asymmetric splits: `7 cols (content) + 5 cols (visual anchor)` or `8 cols + 4 cols`.
   - Never default to a basic 6+6 split for marketing heroes.
2. **F-Pattern for High-Density Reading:**
   - In comparison tables, feature chapters, and FAQs, align content along a strong left anchor edge.
   - Restrict prose measure to `max-w-[62ch]` (45–75 characters per line). Long lines cause tracking fatigue.
3. **Z-Pattern for High-Conversion Hero & CTA Bands:**
   - Top-Left: Brand mark & proof badge.
   - Top-Right: Primary navigation CTA.
   - Bottom-Left: Value proposition headline & action buttons.
   - Bottom-Right: Visual product anchor / diagram.

---

## 5. Step-by-Step AI Decision Tree: Designing a New Screen

When an AI receives a prompt to create a new page, it must execute this deterministic algorithm:

```text
STEP 1: SELECT THEME IDENTITY (from 16-theme-catalogue-and-palettes.md)
  │  ├─ Developer / CLI Tool -> Select VS Code Dark+ or Tokyo Night
  │  ├─ Cloud / AI Platform   -> Select Navy Blue & Purple (#0A0F1D + Electric Violet)
  │  ├─ Data / Observability  -> Select Plasma Sequential or RAG Diverging
  │  └─ Enterprise / Bespoke  -> Select Warm Editorial (#0B0A09 + Brand Amber)
  │
STEP 2: ESTABLISH THE 4-PLANE DEPTH MATRIX
  │  Plane 0: Canvas Base (L: 6-8%)
  │  Plane 1: Raised Wells (L: 10-12%)
  │  Plane 2: Surface Cards (L: 14-16%)
  │  Plane 3: Elevated Modals (L: 18-20%)
  │
STEP 3: FORMULATE THE 5-QUESTION HERO
  │  Draft H1 (≤ 9 words)
  │  Draft Subhead (who + outcome)
  │  Place Visual Anchor in right 5 columns (occupying ≥ 25% of frame)
  │  Configure 1 Primary CTA + 1 Secondary CTA + micro reassurance line
  │
STEP 4: SEQUENCE SECTION PACING & PLANES
  │  Section 1 (Hero): Plane 0, Open (160px pad)
  │  Section 2 (Proof Strip): Plane 1, Compressed (48px pad)
  │  Section 3 (Problem/Solution): Plane 0, Asymmetric 5+7 (120px pad)
  │  Section 4 (Interactive Showcase): Plane 1, Controlled Sliding Carousel
  │  Section 5 (CTA Band): Plane 0, Focused Closing
  │
STEP 5: APPLY KINETIC TOKENS & SAFENETS
  │  Entrances: maskUp on display titles, rise on cards (stagger 70ms)
  │  Transitions: transform: translate3d and opacity ONLY
  │  Reduced Motion: instant resolution, static fallback for carousels
  │
STEP 6: RUN CONTRAST & ACCESSIBILITY AUDIT
     Text on Base ≥ 4.5:1 (AAA requires ≥ 7:1)
     Visible keyboard focus ring on all interactive elements
     Zero literal color classes in components
```

---

## 6. Train-and-Learn Self-Verification Checklist for AI Models

Before emitting final code, every AI agent must mentally verify each item:

- [ ] `/goal` Did I select a recognized theme from `16-theme-catalogue-and-palettes.md`?
- [ ] `/learn` Did I verify that the canvas background is NOT pure black `#000000`?
- [ ] `/goal` Did I eliminate the generic 3-card icon grid in favor of an asymmetric or carousel layout?
- [ ] `/learn` Does the first viewport answer What, Who for, What outcome, Why believe, and What now?
- [ ] `/goal` Does the hero carry a tangible visual anchor occupying at least 25% of the desktop frame?
- [ ] `/learn` Is there exactly ONE primary accent button visible in the viewport (Von Restorff rule)?
- [ ] `/goal` Are all transitions restricted to `transform` and `opacity` (zero layout property animations)?
- [ ] `/learn` Does the layout support `prefers-reduced-motion: reduce` with instant static rendering?
- [ ] `/goal` Are all colors referenced through semantic tokens (zero raw hex or literal Tailwind colors)?
- [ ] `/learn` Is the text line measure capped at 62ch to prevent visual fatigue?
