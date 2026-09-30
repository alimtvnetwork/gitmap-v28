# 19 — Modern Motion & Sliding Interactions

> **/goal** Master the construction of modern, production-grade sliding carousels, scroll-linked choreography, entrance grammars, and section pacing without layout thrashing, maintaining strict 60fps performance and full accessibility.
> **/learn** Implement controlled multi-card carousels with invisible wrap resets, direction-aware entrance choreography (`maskUp`, `rise`, `drawRule`), section density alternation, and instant zero-movement fallbacks under `prefers-reduced-motion: reduce`.

**Version:** 4.0.0
**Updated:** 2026-09-24
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Principles of High-Craft Motion

1. **Motion Explains, Never Decorates:** Every transition communicates a state change, a spatial relationship, or a causal action. Motion that exists purely for flash is distraction and must be deleted.
2. **Transform & Opacity Only:** Animating layout properties (`width`, `height`, `top`, `left`, `margin`, `padding`) triggers expensive CPU browser re-flows and causes layout thrashing. Only GPU-composited properties (`transform: translate3d/scale/rotate`, `opacity`, and `filter: blur`) are permitted.
3. **One Easing Family:** Visual cohesion requires a shared kinetic personality. Bouncy spring physics mixed with linear easing makes an interface look stitched together from random templates.
4. **Reading First, Motion Second:** The critical LCP element (hero headline) must be rendered in the very first paint. Motion must never gate or delay the user's ability to read critical content.
5. **Reduced Motion is a First-Class Citizen:** Under `prefers-reduced-motion: reduce`, all animations resolve instantly with zero transform travel. Content and interactive controls remain 100% accessible.

---

## 2. Kinetic Tokens & Easing Curves

| Token Name | Value | Purpose & Application |
|:---|:---|:---|
| `--ease-out-editorial` | `cubic-bezier(0.16, 1, 0.3, 1)` | The signature entrance easing; smooth rapid acceleration with a long, calm deceleration tail. |
| `--ease-state` | `cubic-bezier(0.4, 0, 0.2, 1)` | Micro-interactions, button hovers, dropdown expansions. |
| `--ease-linear` | `linear` | Strictly for continuous marquees and infinite ticker tracks. |
| `--dur-instant` | `120ms` | Icon toggles, selection checkmarks. |
| `--dur-fast` | `200ms` | Hover states, color shifts, card lift. |
| `--dur-base` | `320ms` | Modal appearances, drawer slides, accordion reveals. |
| `--dur-slow` | `520ms` | Section entrance choreography, card grid stagger. |
| `--dur-extended` | `800ms` | Hero headline reveal, large signature diagram builds. |
| `--stagger-step` | `60ms – 80ms` | Delay between consecutive children in a staggered grid. |

---

## 3. The Controlled Multi-Card Sliding Carousel

Modern product showcases, customer stories, and engineering portfolios require **controlled multi-card sliding loops** that provide effortless navigation without page-level scroll hijacking.

```text
DESKTOP (≥ 1024px): 3 visible cards + 1 card step per advance
┌────────────────────────────────────────────────────────────────────────┐
│  [ Card N-1 ]         [ Card N (Active) ]         [ Card N+1 ]         │
│  Border: Hairline     Border: Strong Glow         Border: Hairline     │
└────────────────────────────────────────────────────────────────────────┘
          [◄ Previous]    ●   ●   ●   ○   ○    [Next ►]

MOBILE (≤ 768px): 1 active card + peek glimpse of next card (24px)
┌──────────────────────────────────────────────┐
│  ┌─────────────────────────┐  ┌───────────┐  │
│  │   Card N (Active)       │  │ Card N+1  │  │ (controlled horizontal stage)
│  └─────────────────────────┘  └───────────┘  │
└──────────────────────────────────────────────┘
```

### 3.1 Architectural Rules for Sliding Carousels

1. **Discrete One-Card Advance:** Clicking Next/Previous or pressing arrow keys advances the visible stage by exactly **one card width plus gutter**, never a jarring multi-card screen wipe.
2. **Controlled Stage Transform:** The track container uses `transform: translate3d(-Xpx, 0, 0)` driven by an active index state. Native scroll snapping is disabled for animated carousels to avoid conflict between browser scroll engines and user gestures.
3. **Seamless Invisible Wrap Reset:**
   - To create an infinite loop without jumping, append clone elements of the first and last cards to the edges of the data array.
   - When the active index reaches the clone boundary, wait for the `transitionend` event, then silently disable transitions (`transition: none`), reset the index to the authentic card, force a reflow, and restore the transition.
4. **Comprehensive Autoplay Safeguards:**
   - Standard autoplay interval: **5.2s**.
   - Autoplay MUST automatically pause when:
     - The user hovers the pointer over any card or control (`:hover`).
     - Any interactive element inside the carousel gains keyboard focus (`:focus-within`).
     - A touch gesture is active (`touchstart`).
     - The carousel scrolls outside the visible viewport (`IntersectionObserver`).
     - The system has `prefers-reduced-motion: reduce` active.
5. **No Hover Lift on Carousel Cards:** Cards in an active carousel must NOT apply `translateY` hover lifts. Vertical transforms fight horizontal track transforms and cause visual jitter. Use border brightening and subtle image scale (`scale(1.02)`) instead.

---

## 4. Scroll Choreography & Entrance Grammar

When a user scrolls down the page, content must enter with a disciplined, predictable grammatical vocabulary:

### 4.1 Entrance Grammar Matrix

| Pattern Name | Visual Mechanic | Application |
|:---|:---|:---|
| `rise` | `opacity: 0 → 1`, `transform: translate3d(0, 18px, 0) → (0, 0, 0)` over `520ms` | Default for headings, body copy, and grid cards. |
| `riseStagger` | Children execute `rise` sequentially, delayed by `70ms` per item | Card grids, metric shelves, feature rows. |
| `maskUp` | CSS `clip-path: inset(0 0 100% 0) → inset(0 0 0 0)` per text line | Display headlines, section H2 titles. |
| `drawRule` | `transform: scaleX(0) → scaleX(1)` from `transform-origin: left` | Structural dividers, eyebrow rules, timeline rails. |
| `countUp` | Numerical figures increment from 0 to target over `800ms` | Metric strips, stats blocks (requires `tabular-nums`). |
| `depthReveal` | `opacity: 0 → 1`, `scale: 0.96 → 1.0` over `650ms` | Focal media cards, closing call-to-action blocks. |
| `pairedEditorial` | Left text enters `x: -28px → 0`; Right visual enters `x: +28px → 0` | Two-column asymmetric problem/solution sections. |

### 4.2 Section-by-Section Choreography Blueprint

```text
1. HERO:
   - Eyebrow & Headline: maskUp line-by-line (stagger 80ms)
   - Lead Paragraph & Buttons: rise at +300ms
   - Right Visual Anchor: depthReveal at +200ms
   - Background Living Canvas: gentle fade-in (opacity 0 → 1 over 1.2s)

2. PROOF STRIP / LOGO MARQUEE:
   - Numerals countUp in tabular figures
   - Hairline divider drawRule across 100% width
   - Logo track continuous linear drift (40s loop, pause on hover)

3. PROBLEM / SOLUTION:
   - Problem statement enters from left; Metric diagram enters from right
   - Diagram bars expand upward from baseline (scaleY 0 → 1)

4. SERVICES CHAPTERS:
   - Left index numerals enter first (01 → 04)
   - Right content cards follow with 60ms stagger
   - Card hover warms border and highlights internal icon

5. SIGNATURE FILTER / FUNNEL (Pinned Canvas):
   - Stage rail pins for 200vh while user scrolls
   - Canvas dots eliminate stage-by-stage (1,000 → 20 → final shortlist)
   - Live counter tracks survivor count in real time

6. PROCESS TIMELINE:
   - Timeline rule draws horizontally across desktop grid
   - Milestone nodes pop in sequence with subtle spring scale

7. SELECTED WORK CAROUSEL:
   - Controlled multi-card carousel slides one card per click/swipe
   - Progress bar tracks index

8. FAQ DISCLOSURE:
   - Accordion height animates smoothly via grid-template-rows (0fr → 1fr)
   - Plus icon rotates 180° into a minus
```

---

## 5. Section Rhythm & Spatial Pacing

Pages that maintain the exact same section height and padding feel robotic and monotonous. A high-craft page uses **dynamic rhythm**:

```text
[ HERO: Open, Spacious ]        Padding: 160px top / 120px bottom   (Plane 0 Base)
          │
[ PROOF STRIP: Compressed ]     Padding: 48px top / 48px bottom     (Plane 1 Raised)
          │
[ PROBLEM: Dense Editorial ]    Padding: 120px top / 120px bottom   (Plane 0 Base)
          │
[ SERVICES: Open Card Grid ]    Padding: 160px top / 140px bottom   (Plane 1 Raised)
          │
[ FUNNEL: Full-Bleed Pinned ]   Height: 200vh Pinned Stage          (Plane 0 Base)
          │
[ CTA BAND: Focused Closing ]   Padding: 160px top / 120px bottom   (Plane 0 Base)
```

- **Rhythm Rule:** Alternate between open breathing sections (`140px–160px`) and compressed structural bands (`48px–96px`).
- **Container Break Rule:** At least two sections per page should break the standard `1280px` container (e.g. a full-bleed marquee and a full-bleed visual showcase).

---

## 6. Strict Reduced Motion Fallbacks

When the user specifies `prefers-reduced-motion: reduce`:

```css
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

### Component Fallback Manifest:
- **Carousels:** Autoplay is disabled. Cloned wrap elements are removed. The track renders as a standard accessible horizontal scroll list with native scrollbar and arrow buttons.
- **Marquees:** Infinite translation stops. Items arrange into an accessible multi-row flex grid.
- **Stat Counters:** Numeral counting is bypassed; final target numbers render immediately.
- **Entrances:** All elements render at `opacity: 1` and `transform: none` on first paint without delay.
- **Pinned Funnel:** Pinned scrolling is unpinned. Renders as a vertical stacked sequence of stages with a static diagram.

---

## 7. AI Verification Checklist for Motion & Sliding

- [ ] `/goal` Verify only `transform` and `opacity` are animated (zero animated `width`, `height`, `top`, `left`).
- [ ] `/learn` Verify carousel uses one-card-step navigation and pauses on hover, focus, touch, and offscreen.
- [ ] `/goal` Verify carousel cards have zero vertical `translateY` hover lifts (use border and image scale instead).
- [ ] `/learn` Verify all entrance animations fire once in reading order and do not replay when scrolling upward.
- [ ] `/goal` Verify `prefers-reduced-motion: reduce` completely removes transforms and renders complete static UI.
- [ ] `/learn` Verify 60fps frame rate without layout thrashing on throttled devices.
