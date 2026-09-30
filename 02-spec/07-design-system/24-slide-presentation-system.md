# 24 — Slide Presentation System & Presenter Engine Architecture

> **/goal** Specify the comprehensive architecture for web-based slide presentations, 16:9 responsive canvas scaling, draggable webcam picture-in-picture (PIP) overlays, incremental step reveals, and multi-deck routing synthesized from local presentation repositories.
> **/learn** Master the responsive coordinate transforms (`1920x1080` canvas), `PresenterWebcamOverlay` webcam stream integration, `stepMotionOverride` state tracking, dual-screen presenter consoles, and modular LESS slide styling.

**Version:** 1.0.0
**Updated:** 2026-09-24
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. System Overview & Presentation Repositories Synthesis

Web-based presentation decks combine the interactivity and fluidity of the browser with the structured delivery of professional slide software. Based on architectural patterns synthesized across active presentation systems (e.g. `remix-of-presentation-riseup-asia`, `flat-slide-show`, and `slides-app`), modern presentation systems operate on five core pillars:

1. **Fixed-Aspect Ratio Virtual Canvas:** Standard `1920x1080` coordinate space scaled down dynamically using CSS `transform: scale(min(w/1920, h/1080))` so layouts never shift or distort across screens.
2. **Floating Draggable Webcam PIP (`PresenterWebcamOverlay`):** Video stream with circular boundary, glowing accent perimeter, draggable coordinate snapping, mirror toggle, and customizable opacity.
3. **Step-by-Step Progressive Reveals:** Sub-slide animations where individual bullets, code blocks, or diagram nodes reveal on subsequent arrow key presses before advancing slides.
4. **Presenter Dual-Screen Console:** Broadcast channel sync communicating between audience display and presenter control view (notes, elapsed timer, next slide preview).
5. **Standalone Vector Layouts:** System illustrated via standalone SVG diagrams in `01-svg/slide-layout.svg`.

---

## 2. 16:9 Canvas Virtual Coordinate Architecture

To guarantee that typography, diagrams, and cards render identically regardless of audience screen resolution, all presentation slides render inside an isolated virtual canvas:

```text
+-------------------------------------------------------------------------------+
| BROWSER VIEWPORT (Window W x H)                                               |
|                                                                               |
|   +-----------------------------------------------------------------------+   |
|   | VIRTUAL SLIDE CANVAS (Fixed 1920px x 1080px)                          |   |
|   |                                                                       |   |
|   | transform-origin: center center;                                      |   |
|   | transform: scale(min(viewportWidth / 1920, viewportHeight / 1080));  |   |
|   |                                                                       |   |
|   | [Slide Content: Header, Body Cards, Metrics, Diagrams]               |   |
|   |                                                                       |   |
|   |                             [Draggable Webcam PIP: 160px x 160px]     |   |
|   +-----------------------------------------------------------------------+   |
|                                                                               |
+-------------------------------------------------------------------------------+
```

### 2.1 Responsive Canvas Scaling (TypeScript)

```typescript
export interface ViewportScale {
  scale: number;
  offsetX: number;
  offsetY: number;
}

export function calculateSlideScale(
  containerWidth: number,
  containerHeight: number,
  baseWidth = 1920,
  baseHeight = 1080
): ViewportScale {
  const scale = Math.min(containerWidth / baseWidth, containerHeight / baseHeight);
  const offsetX = (containerWidth - baseWidth * scale) / 2;
  const offsetY = (containerHeight - baseHeight * scale) / 2;

  return { scale, offsetX, offsetY };
}
```

---

## 3. Draggable Webcam PIP Overlay (`PresenterWebcamOverlay`)

The presenter's camera stream floats above the presentation canvas to maintain eye contact with the audience without occluding content.

### 3.1 Component Features
- Circular or squircle geometry with customizable corner radius (`border-radius: 9999px` or `2rem`).
- Real-time drag positioning with bounded edge collision detection.
- Subtle violet/gold glowing shadow bloom.
- Video flip/mirror control (`transform: scaleX(-1)`).
- Quick keyboard shortcut (`C` key) to cycle size presets (Small 120px, Medium 180px, Large 240px, Hidden).

### 3.2 LESS Styling (Preferred)

```less
// ============================================================================
// PRESENTER WEBCAM OVERLAY (LESS PREFERRED)
// ============================================================================

.presenter-webcam-container {
  position: absolute;
  z-index: 9999;
  width: 180px;
  height: 180px;
  border-radius: 9999px;
  overflow: hidden;
  border: 3px solid #8b5cf6;
  box-shadow: 0 10px 25px -5px rgba(139, 92, 246, 0.5),
              0 0 0 1px rgba(255, 255, 255, 0.1);
  cursor: grab;
  user-select: none;
  touch-action: none;
  transition: box-shadow 200ms ease, border-color 200ms ease;

  &:active {
    cursor: grabbing;
    box-shadow: 0 15px 35px -5px rgba(139, 92, 246, 0.7);
  }

  video {
    width: 100%;
    height: 100%;
    object-fit: cover;
    pointer-events: none;

    &.mirrored {
      transform: scaleX(-1);
    }
  }

  .status-dot {
    position: absolute;
    top: 12px;
    right: 12px;
    width: 10px;
    height: 10px;
    border-radius: 9999px;
    background-color: #10b981;
    border: 2px solid #0b1329;
  }
}
```

---

## 4. Incremental Step Reveal Architecture

Slides frequently contain multiple sequential ideas that must reveal one at a time. The engine models this using a two-tier position tracker:

```typescript
export interface PresentationState {
  currentSlideIndex: number;
  currentStepIndex: number;
  totalSlides: number;
  totalStepsInCurrentSlide: number;
}
```

### 4.1 Step Motion Execution
Elements with `data-step="1"`, `data-step="2"` stay hidden or dimmed (`opacity: 0.15; filter: blur(2px);`) until the presenter reaches their respective step. Upon activation, they animate to full prominence using `@ease-emphasized` (`cubic-bezier(0.16, 1, 0.3, 1)`).

---

## 5. Deck Registry & URL Routing

Multi-deck presentation repositories organize decks modularly:

```text
slides-app/
├── src/
│   ├── decks/
│   │   ├── 01-agentic-ai/
│   │   │   ├── 01-intro.tsx
│   │   │   ├── 02-architecture.tsx
│   │   │   └── manifest.json
│   │   └── 02-coding-guidelines/
│   │       ├── 01-booleans.tsx
│   │       ├── 02-error-handling.tsx
│   │       └── manifest.json
│   ├── components/
│   │   ├── SlideViewport.tsx
│   │   ├── PresenterWebcamOverlay.tsx
│   │   └── SpeakerNotesView.tsx
│   └── main.tsx
```

Decks are accessible via hash or path routing:
- Audience View: `/deck/01-agentic-ai/#/4` (Slide 4)
- Presenter Console: `/deck/01-agentic-ai/presenter#/4` (Slide 4 with notes & timer)
