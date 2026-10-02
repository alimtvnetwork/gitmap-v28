# 24 — Slide Presentation System & Presenter Engine Architecture

> **/goal** Architect and enforce the complete, production-grade presentation presenter engine, dual-window presenter console, draggable webcam PIP overlay, incremental step reveals, overview thumbnail grid, and keyboard orchestration.
> **/learn** Master the responsive 1920×1080 canvas scaling, BroadcastChannel dual-screen console synchronization, `PresenterWebcamOverlay` video stream mechanics, progressive step-motion reveals, `G` key overview grid, `J` key top slide jumper, and the complete deck keyboard map.

**Version:** 4.1.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. System Overview & The 5 Presenter Pillars

The **Slide Presentation System** delivers high-authority interactive slide presentations running in web browsers:
1. **Fixed-Aspect Ratio Virtual Canvas:** Standard `1920×1080` (16:9) coordinate space scaled uniformly using CSS vector scaling so diagrams, cards, and typography never wrap or break.
2. **Dual-Screen Presenter Console:** BroadcastChannel synchronized console running at `/present` in a separate `1280×800` window for speaker notes, pacing clock, and next-slide previews.
3. **Floating Draggable Webcam PIP (`PresenterWebcamOverlay`):** Hardware-accelerated webcam feed with mirrored preview, circular boundary, glowing accent perimeter, and hotkey sizing.
4. **Incremental Step Reveals:** Sub-slide sequential reveals where individual points, code blocks, or cards activate before the deck advances.
5. **Presenter Tools & Hotkeys:** `G` key slide grid, `J` key quick jumper, `C` key camera toggle, `F` key fullscreen, `B`/`E` builder mode, and `/` keyboard shortcut modal.

---

## 2. 16:9 Virtual Canvas Scaling Architecture

```text
+-------------------------------------------------------------------------------+
| BROWSER VIEWPORT (Window W x H)                                               |
|                                                                               |
|   +-----------------------------------------------------------------------+   |
|   | VIRTUAL SLIDE STAGE (Fixed 1920px x 1080px)                           |   |
|   | transform-origin: center center;                                      |   |
|   | transform: scale(min(viewportWidth / 1920, viewportHeight / 1080));  |   |
|   | contain: layout paint; isolation: isolate; will-change: transform;    |   |
|   |                                                                       |   |
|   | [Slide Content: 20 Master Enterprise Layout Models]                  |   |
|   |                                                                       |   |
|   |                          [Draggable Webcam PIP: 120px / 180px / 240px]|   |
|   +-----------------------------------------------------------------------+   |
|                                                                               |
+-------------------------------------------------------------------------------+
```

All slide geometry is authored strictly in pixels on `1920 × 1080`. The container detects resizing via `ResizeObserver` and updates `--stage-scale`:

```typescript
export function calculateSlideScale(viewportWidth: number, viewportHeight: number): number {
  return Math.min(viewportWidth / 1920, viewportHeight / 1080);
}
```

---

## 3. Dual-Screen Presenter Console (`/present`)

For rehearsals and live stage presentations, the presenter launches the dedicated console via the Deck Menu or by navigating to `/present`:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ [● LIVE 14:22] PACING: ON TRACK   Slide 5 / 37          [⤢ Full] [× Close]  │
├──────────────────────────────────────┬──────────────────────────────────────┤
│ CURRENT SLIDE (Scaled 16:9 Preview)  │ SPEAKER NOTES (20px Ubuntu/Poppins)  │
│ ┌──────────────────────────────────┐ │ • Emphasize 12-minute release cycle  │
│ │                                  │ │ • Direct audience to ROI metrics     │
│ │      Live Audience View          │ │ • Pause 3 seconds for slide ingest   │
│ │                                  │ │ • Call out security compliance       │
│ └──────────────────────────────────┘ ├──────────────────────────────────────┤
│ NEXT SLIDE PREVIEW (Upcoming: #6)    │ PRESENTATION PACING & TIMERS         │
│ ┌──────────────────────────────────┐ │ Elapsed: 14m 22s / Budget: 25m 00s   │
│ │  Next: SaaS Pricing Matrix       │ │ Progress: [████████░░░░░░░░░░░░] 57% │
│ └──────────────────────────────────┘ │ Next Milestone: Product Deep Dive    │
└──────────────────────────────────────┴──────────────────────────────────────┘
```

### 3.1 BroadcastChannel Synchronization
The audience window and the presenter console communicate via `BroadcastChannel("deck-sync")` with fallback to `localStorage` events:

```typescript
export interface DeckSyncMessage {
  type: "NAVIGATE" | "STEP" | "THEME_CHANGE" | "RELOAD";
  slideIndex: number;
  stepIndex: number;
  themeId?: string;
  timestamp: number;
}
```

- When the presenter presses Next in `/present`, the audience display navigates instantly (`< 5ms` latency).
- Window parameters: Opens at `1280×800` pixels (`window.open('/present', 'PresenterConsole', 'width=1280,height=800')`).

---

## 4. Draggable Webcam PIP Overlay (`PresenterWebcamOverlay`)

The presenter webcam floats in a dedicated layer above the slide canvas (`z-index: 55`):

### 4.1 Component Specifications
- **Hardware Integration:** Captured via `navigator.mediaDevices.getUserMedia({ video: { width: 1280, height: 720 }, audio: false })`.
- **Render Element:** `<video autoPlay playsInline muted className="mirrored" />` with `transform: scaleX(-1)` for natural mirrored feedback.
- **Sizing Presets (Cycled via `C` Key):**
  1. `Small`: `120×120px` circle (`rounded-full`).
  2. `Medium` *(Default)*: `180×180px` circle (`rounded-full`).
  3. `Large`: `240×160px` rounded squircle (`rounded-2xl`).
  4. `Hidden`: Component unmounts, camera stream stops to conserve battery.
- **Interactive Dragging:** Framer Motion `drag` with drag constraints bounded to the viewport edges and a `24px` padding snap.
- **Visual Styling:** `2.5px solid hsl(var(--primary))` with ambient glow `0 12px 32px rgba(0,0,0,0.4)` and a green `status-dot` (`10×10px`, `#10B981`) in the top-right.

---

## 5. Incremental Step Reveals Architecture

Slides with complex multi-step diagrams, code progressions, or bullet points reveal items sequentially:

```typescript
export interface StepPresentationState {
  currentSlideIndex: number;
  currentStepIndex: number;
  totalStepsInCurrentSlide: number;
}
```

- **Resting (Unrevealed) State:** Elements tagged with `data-step="N"` have `opacity: 0.15; filter: blur(2px); transform: translateY(8px);` when `currentStepIndex < N`.
- **Activated State:** When `currentStepIndex >= N`, elements transition to `opacity: 1; filter: blur(0px); transform: translateY(0);` over `300ms cubic-bezier(0.16, 1, 0.3, 1)`.
- **Advance Contract:** Pressing `ArrowRight` or `Space` increments `currentStepIndex` until `totalStepsInCurrentSlide` is reached, then navigates to the next slide.

---

## 6. Overview Thumbnail Grid (`G` Key)

Pressing `G` opens a full-screen interactive slide gallery:
- **Layout:** Responsive thumbnail grid (`grid-cols-2 md:grid-cols-3 lg:grid-cols-4`, gap `24px`, padding `48px`).
- **Thumbnail Item:** Renders an accurately scaled miniature of each slide (`280×157px`), slide number badge (`05`), and slide title.
- **Navigation:** Arrow keys move focus across thumbnails; `Enter` or click activates that slide and dismisses the grid.
- **Search Filter:** Top input bar allows instant typing to filter slides by title keywords.

---

## 7. Top Slide Jumper (`J` Key)

Pressing `J` or double-clicking the slide counter in the controller HUD opens the Quick Jumper:
- A compact floating input chip anchored at `top: 24px, center`: `[ Jump to slide: (1–37) ]`.
- Typing a number and pressing `Enter` navigates immediately and updates `#slide-{n}`. `Escape` dismisses.

---

## 8. Complete Keyboard Shortcut Map

| Key / Shortcut | Action | Description |
|:---|:---|:---|
| `ArrowRight` / `Space` | **Next** | Advances step reveal, or moves to next slide |
| `ArrowLeft` / `Shift+Space` | **Previous** | Regresses step reveal, or moves to previous slide |
| `PageDown` / `PageUp` | **Slide Jump** | Direct previous/next slide navigation |
| `F` | **Fullscreen** | Toggles HTML5 Fullscreen API |
| `B` or `E` | **Slide Builder** | Toggles live in-canvas slide builder mode |
| `C` | **Camera PIP** | Cycles webcam PIP presets (Hidden ➔ Small ➔ Medium ➔ Large) |
| `G` | **Overview Grid** | Opens thumbnail overview gallery |
| `J` | **Slide Jumper** | Opens quick numeric jump input |
| `P` | **Presenter Console** | Opens `/present` dual-screen console window |
| `1` through `0` | **Theme Switch** | Instantly switches between the 10 built-in themes |
| `/` or `?` | **Keyboard Help** | Opens keyboard shortcuts dialog |
| `Escape` | **Dismiss** | Closes active modal, popover, or exits fullscreen |

---

## 9. LLM Guide Download & Copy

Inside the Deck Menu (`FileJson` icon in HUD), presenters and AI agents can export the complete presentation deck manifest:
- **Download LLM Guide:** Generates a structured `.md` document detailing every slide's type, title, body content, metrics, and tokens.
- **Copy LLM Guide:** Writes the markdown deck manifest directly to clipboard for AI inspection.

---

## 10. Anti-Hallucination & Quality Verification Checklist

- [ ] All coordinates author strictly on the virtual `1920×1080` reference grid.
- [ ] Dual-screen presenter console syncs via BroadcastChannel with sub-5ms latency.
- [ ] Webcam PIP overlay includes mirror preview, edge clamping, and sizing hotkeys.
- [ ] Incremental step reveals preserve dimmed blur state (`opacity: 0.15; filter: blur(2px)`) until triggered.
- [ ] Complete keyboard map is bound to window event listeners with input tag suppression.
- [ ] All 10 themes and 20 layout models operate seamlessly through the presenter engine.
