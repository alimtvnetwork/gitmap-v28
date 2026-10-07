# 43 — Presenter Webcam PIP Overlay & Auto-Frame Specification

> **/goal** Architect, construct, and calibrate the floating presenter webcam picture-in-picture (PIP) overlay, squircle gold-rim frame, draggable stage coordinates, auto-frame face tracking, and presentation key passthrough.
> **/learn** Master the `PresenterWebcamProvider` state machine, 4 stepped sizing presets (S, M, L, XL), squircle mask geometry, keyboard shortcuts (`C`, `F`, `M`, `+`, `-`, `O`, `H`, `1`), and the `deck:webcam-passthrough` event forwarder.

**Version:** 4.3.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Executive System Overview

The **Presenter Webcam Overlay** provides an integrated broadcast-grade camera bubble floating over the 1920×1080 slide stage:
1. **Camera Life Cycle (`getUserMedia`):** On-demand hardware activation with error recovery and permissions management.
2. **4 Sizing Scales & Free Pointer Drag:** Stepped presets (S/M/L/XL) plus 16:9 locked freehand pointer resize.
3. **Squircle Gold-Rim Frame & Vignette Halo:** Rounded squircle mask with metallic gold rim and optional ambient shadow halo.
4. **Auto-Frame Face Tracking:** Continuous face detection adjusting scale and translation to center the speaker.
5. **Keyboard Hotkeys & Event Forwarding:** Transparent event bridging so slide deck navigation keys continue functioning when the camera is active or fullscreen.

---

## 2. Architecture & State Machine (`PresenterWebcamProvider`)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ <PresenterWebcamProvider> (React Context & State Machine)                   │
│   • Manages phase: idle | starting | active | error                         │
│   • Holds stream, stage position {x, y}, scale preset, shape, and zoom      │
│   • Persists settings to localStorage ("deck.webcam.settings.v1")           │
│                                                                             │
│   ┌───────────────────────────────────────────────────────────────────────┐ │
│   │ <PresenterWebcamOverlay> (Viewport Render Layer - z-index: 55)        │ │
│   │   • Video surface with squircle/circle CSS mask-image                 │ │
│   │   • Draggable pointer handlers in 1920×1080 stage coordinates         │ │
│   │   • Keyboard shortcut listener (C, M, F, +, -, O, H, 1, Esc)          │ │
│   │   • useAutoFrame() video transform coordinator                        │ │
│   └───────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.1 State Model Interface

```typescript
export type WebcamPhase = "idle" | "starting" | "active" | "error";
export type WebcamShape = "squircle" | "circle" | "rectangle";
export type WebcamSizePreset = "sm" | "md" | "lg" | "xl";

export interface WebcamState {
  phase: WebcamPhase;
  stream: MediaStream | null;
  position: { x: number; y: number }; // In 1920×1080 stage pixels
  sizePreset: WebcamSizePreset;
  shape: WebcamShape;
  zoom: number;                       // 1.0 to 3.0
  isMirrored: boolean;                // Default true
  isAutoFrame: boolean;               // Face tracking
  isHalo: boolean;                    // Ambient vignette shadow
  isFullscreen: boolean;              // Stage cover
}
```

---

## 3. Geometric Sizing Presets & Stage Coordinates

Webcam dimensions are calibrated for the standard 1920×1080 reference canvas:

| Preset ID | 16:9 Locked Dimensions | Circular / Squircle Dimensions | Optimal Stage Location |
|:---|:---:|:---:|:---|
| **`sm`** | `240 × 135px` | `160 × 160px` | Top-right corner, dense data slides |
| **`md`** *(Default)* | `360 × 202px` | `220 × 220px` | Bottom-right corner, balanced presentation |
| **`lg`** | `480 × 270px` | `300 × 300px` | Right editorial column, persona introduction |
| **`xl`** | `640 × 360px` | `400 × 400px` | Narrative talk, full keynote segment |
| **`stage-fill`** (`1` key) | `1920 × 1080px` | Fullscreen Stage | Direct presenter audience address |

### 3.1 Positioning Boundaries
- **Brand Inset Boundary:** Presenter webcam automatically respects `--brand-inset-x` (218px) when anchored to top corners to prevent colliding with brand logos.
- **Header & Footer Clearance:** Must never obscure the top 80px brand band or bottom 80px controller pill zone.

---

## 4. Squircle Gold Rim & Mask Geometry

```css
/* Squircle Mask with Super-Ellipse Curvature */
.webcam-squircle {
  mask-image: radial-gradient(circle at center, black 65%, transparent 70%);
  border-radius: 28% / 28%;
}

/* Metallic Gold Rim Plate */
.webcam-gold-rim {
  border: 2px solid transparent;
  background: linear-gradient(135deg, #fdd072, #bb8e4e, #f3c258) border-box;
  box-shadow: 0 12px 36px rgba(0, 0, 0, 0.45), 0 0 20px rgba(243, 194, 88, 0.25);
}

/* Ambient Vignette Halo */
.webcam-halo {
  box-shadow: 0 0 0 8px rgba(243, 194, 88, 0.12), 0 20px 50px rgba(0, 0, 0, 0.6);
}
```

---

## 5. Keyboard Shortcuts & Deck Event Passthrough

The webcam overlay registers window keyboard shortcuts while avoiding intercepting typing in text inputs:

| Key | Action | Behavior |
|:---:|:---|:---|
| **`C`** | Toggle Camera | Starts camera stream or shuts down hardware tracks |
| **`M`** | Toggle Mirror | Flips video horizontal axis (`scaleX(-1)`) |
| **`F`** | Toggle Auto-Frame | Engages face tracking zoom & center algorithm |
| **`+` / `=`** | Zoom In | Increments digital zoom by `+0.2` (max `3.0`) |
| **`-` / `_`** | Zoom Out | Decrements digital zoom by `-0.2` (min `1.0`) |
| **`O`** | Toggle Shape | Cycles `squircle` -> `circle` -> `rectangle` |
| **`H`** | Toggle Halo | Toggles ambient golden glow aura |
| **`1`** | Stage Fill | Expands camera to fill 1920×1080 stage; press again to restore |
| **`Esc`** | Exit Full/Dismiss | Collapses stage-fill or closes active camera panel |

### 5.1 Slide Navigation Event Passthrough
When the camera is active, keyboard navigation keys (`ArrowRight`, `ArrowLeft`, `Space`) dispatch a custom synthetic event `deck:webcam-passthrough` so the presenter can advance slides without clicking back into the stage canvas:

```typescript
function handleWebcamKey(e: KeyboardEvent) {
  if (["ArrowRight", "ArrowLeft", "Space"].includes(e.code)) {
    window.dispatchEvent(new CustomEvent("deck:webcam-passthrough", { detail: { code: e.code } }));
  }
}
```

---

## 6. Auto-Frame Face Tracker Engine (`useAutoFrame`)

The face tracking hook observes video frames via an offscreen `HTMLCanvasElement` or browser `FaceDetector` API:
1. Detects face bounding box coordinates `(fx, fy, fw, fh)`.
2. Calculates center offset relative to the video viewport center:
   $$\Delta x = \text{targetX} - \text{currentX}, \quad \Delta y = \text{targetY} - \text{currentY}$$
3. Applies smooth damping interpolation ($k = 0.12$) to transform CSS:
   ```css
   transform: scale(var(--cam-zoom)) translate3d(var(--cam-pan-x), var(--cam-pan-y), 0);
   transition: transform 300ms cubic-bezier(0.16, 1, 0.3, 1);
   ```

---

## 7. Anti-Hallucination & Quality Verification Checklist

- [ ] Webcam overlay mounts at `z-index: 55` floating strictly above slides and below modal dialogs.
- [ ] 4 Stepped presets (S, M, L, XL) maintain 16:9 aspect ratio or 1:1 squircle proportions.
- [ ] Camera hardware stream is completely stopped (`track.stop()`) on unmount or toggle off.
- [ ] Face auto-framing applies smooth damping to prevent rapid visual jitter.
- [ ] Keyboard events for slide navigation are passed through via `deck:webcam-passthrough`.
- [ ] User preferences (size, shape, position, mirror) are persisted in `localStorage` under `deck.webcam.settings.v1`.
