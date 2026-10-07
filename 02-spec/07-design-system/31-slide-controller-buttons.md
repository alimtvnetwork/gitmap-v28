# 31 — Slide Controller Buttons, HUD Pill, Webcam Overlay & Dots

> **/goal** Specify the global floating presenter HUD pill, 8-position mounting system, navigation buttons, audio chimes, background music toggle, onboarding popup, webcam PIP overlay, timer, top progress bar, and bottom dot pagination.
> **/learn** These controls sit on the global viewport (`z-index: 50`–`60`), never inside the scaled `1920×1080` stage. Web Audio API synthesizers and the 8-layer slide skeleton operate in strict coordination with these controls.

**Version:** 4.2.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Floating Controller HUD Pill & 8-Position Mounting System

The controller is a single unified floating pill anchored to any of 8 viewport positions (default `'BottomCenter'` or `'TopRight'`). It renders via `createPortal` into `document.body` so it floats above the scaled stage and is never clipped by `overflow: hidden`:

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│  ‹   5 / 37   ›  │  ⌁ Share  │  🎨 Theme  │  📹 Cam  │  🎵 Music  │  ✎ Build  │  ⤢ Full │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
  ▲─── Height: 56px, Radius: 9999px, Backdrop Blur 12px, Hover-Reveal Pill HUD ───────────────▲
```

### 1.1 Geometry, Tokens & Surface Properties

| Property | Value | Notes |
|---|---|---|
| **Height** | `56px` | Constant across all 8 positions |
| **Padding** | `8px 12px` | Horizontal and vertical inner margin |
| **Border Radius** | `9999px` | Full capsule pill |
| **Background** | `hsl(240 8% 8% / 0.85)` + `backdrop-filter: blur(12px)` | Always near-black chrome on all themes |
| **Border** | `1px solid hsl(var(--border))` | Subtle hairline boundary |
| **Shadow** | `0 8px 24px hsl(0 0% 0% / 0.40)` | High-elevation lift |
| **Z-Index** | `60` | Above slides (`1–4`), below modals (`80`) |
| **Internal Dividers** | `1px` wide × `24px` high | In `hsl(var(--border))` separating action clusters |

### 1.2 The 8-Position Mounting Matrix (`ControllerPosition`)

Presenters can configure the pill anchor across 8 edge coordinates, persisted to `localStorage` (`ctrl.position.v1`):

```typescript
export type ControllerPosition =
  | "TopLeft" | "TopCenter" | "TopRight"
  | "BottomLeft" | "BottomCenter" | "BottomRight"
  | "LeftCenter" | "RightCenter";
```

| Anchor Position | Viewport Offset CSS (includes `env(safe-area-inset-*)`) | Tooltip Direction | Expand Axis |
|:---|:---|:---:|:---:|
| **`TopLeft`** | `top: calc(24px + env(safe-area-inset-top)); left: calc(24px + env(safe-area-inset-left));` | `bottom` | downward |
| **`TopCenter`** | `top: calc(24px + env(safe-area-inset-top)); left: 50%; transform: translateX(-50%);` | `bottom` | downward |
| **`TopRight`** | `top: calc(32px + env(safe-area-inset-top)); right: calc(32px + env(safe-area-inset-right));` | `bottom` | downward |
| **`BottomLeft`** | `bottom: calc(24px + env(safe-area-inset-bottom)); left: calc(24px + env(safe-area-inset-left));` | `top` | upward |
| **`BottomCenter`** | `bottom: calc(24px + env(safe-area-inset-bottom)); left: 50%; transform: translateX(-50%);` | `top` | upward |
| **`BottomRight`** | `bottom: calc(24px + env(safe-area-inset-bottom)); right: calc(24px + env(safe-area-inset-right));` | `top` | upward |
| **`LeftCenter`** | `top: 50%; left: calc(24px + env(safe-area-inset-left)); transform: translateY(-50%);` | `right` | inward |
| **`RightCenter`** | `top: 50%; right: calc(24px + env(safe-area-inset-right)); transform: translateY(-50%);` | `left` | inward |

---

## 2. Hover-Reveal & Auto-Hide Lifecycle

To eliminate visual distraction during executive presentations, the controller implements a 2.5s idle fade:
1. **Idle State:** After `2500ms` without mouse movement, pointer input, or touch interaction, the pill fades to `opacity: 0.15` (or `opacity: 0` in presenter mode).
2. **Awakening:** Any mousemove on window, keypress, or hovering over the pill boundary immediately transitions opacity back to `1.0` over `150ms ease-out`.
3. **Lock:** While a popover (Theme Menu, Audio Settings, Quick Jumper) is open, auto-hide is completely suppressed.

---

## 3. Controller Action Groups & Button Contracts

Every icon button inside the pill uses a `40×40px` circular hit target with hover background `hsl(0 0% 100% / 0.08)` and active `hsl(0 0% 100% / 0.14)`:

### 3.1 Navigation Group
- **Previous (`ChevronLeft`):** `40×40px`, disabled when `current === 1` (`opacity: 0.35`). Navigates `current - 1`, updates `#slide-{n}`.
- **Slide Counter:** Poppins `500`, `20px`, tabular numbers, `user-select: none`. Minimum width `64px`, format: `${current} / ${total}`. Double-clicking activates numeric Quick Jumper.
- **Next (`ChevronRight`):** `40×40px`, disabled when `current === total` (`opacity: 0.35`). Navigates `current + 1`, updates `#slide-{n}`.

### 3.2 Share, Theme & Media Group
- **Share (`Share2`):** Deep links to current slide `#slide-{current}`. Uses `navigator.share` or copies URL to clipboard with 2-second Sonner toast.
- **Theme Palette (`Palette`):** Opens upward-expanding Theme Menu popover with live WCAG contrast validator.
- **Presenter Webcam (`Video` / `VideoOff`):** Toggles floating PIP overlay (`z-index: 55`).
- **Background Music (`Music` / `Volume2` / `VolumeX`):** Toggles ambient audio track loop on/off with volume slider.
- **Builder Mode (`Pencil`):** Activates in-canvas drag-and-drop editing and floating `BuilderPanel`. Hotkey: `B` or `E`.
- **Fullscreen (`Maximize2` / `Minimize2`):** Toggles HTML5 Fullscreen API. Hotkey: `F`.

### 3.2.1 Keyboard shortcuts map (`Keyboard`)

| Property | Value |
|---|---|
| **Icon** | Lucide `Keyboard` |
| **Element id** | `ctrl-shortcuts` |
| **Hit target** | `40×40px` circle (same as §3) |
| **Hover / active** | Same as §3 (`hsl(0 0% 100% / 0.08)` / `0.14`) |
| **Tooltip** | `Keyboard shortcuts (/)` |
| **Action** | `setShortcutsOpen(true)` — same state as pressing **`/`** |
| **Data source** | `PRESENTER_SHORTCUTS_CORE` in `42-slide-quiz-preview-chrome-and-default-shadows.md` §5.3 |
| **Placement** | Immediately after **Theme Palette**, before **Presenter Webcam** |

While the shortcuts dialog is open, auto-hide (§2) is suppressed. **`Esc`** closes the dialog.

---

## 4. Audio Chimes & Sound Synthesis (Web Audio API)

All UI feedback sounds are synthesized via the native Web Audio API (`AudioContext`) with zero external sound asset downloads:

```typescript
class DeckSoundEngine {
  private ctx: AudioContext | null = null;
  private init() { if (!this.ctx) this.ctx = new (window.AudioContext || (window as any).webkitAudioContext)(); }

  playTick() { // Navigation step: 800Hz sine decay over 30ms
    this.init(); if (!this.ctx) return;
    const osc = this.ctx.createOscillator(); const gain = this.ctx.createGain();
    osc.frequency.setValueAtTime(800, this.ctx.currentTime);
    gain.gain.setValueAtTime(0.08, this.ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.001, this.ctx.currentTime + 0.03);
    osc.connect(gain); gain.connect(this.ctx.destination);
    osc.start(); osc.stop(this.ctx.currentTime + 0.03);
  }

  playChord() { // Slide jump chord: dual 523Hz (C5) + 659Hz (E5) over 120ms
    this.init(); if (!this.ctx) return;
    [523.25, 659.25].forEach((freq) => {
      const osc = this.ctx!.createOscillator(); const gain = this.ctx!.createGain();
      osc.frequency.setValueAtTime(freq, this.ctx!.currentTime);
      gain.gain.setValueAtTime(0.06, this.ctx!.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.001, this.ctx!.currentTime + 0.12);
      osc.connect(gain); gain.connect(this.ctx!.destination);
      osc.start(); osc.stop(this.ctx!.currentTime + 0.12);
    });
  }
}
```

---

## 5. First-Run Onboarding Popup ("Story")

On a user's first visit (`localStorage.getItem("ctrl.onboarded.v1") !== "1"`), a floating onboarding card mounts above the controller pill:
- **Badge:** `[✨ PRO TIP]` in `--capsule-gold`.
- **Title:** "Keyboard-First Presentation Engine"
- **Shortcuts Grid:**
  - `←` / `→` or `Space` — Navigate slides & reveal steps
  - `F` — Fullscreen presentation mode
  - `G` — 37-slide thumbnail gallery
  - `J` — Direct numeric jump
  - `C` — Live presenter webcam PIP
- **Dismiss:** Button "Got it!" writes `"1"` to `ctrl.onboarded.v1` and unmounts the card. Pressing any navigation key also dismisses automatically.

---

## 6. The 8-Layer Z-Stack Slide Anatomy Skeleton

Every slide is rendered inside a shared 8-layer shell (`<DeckLayout>`). Individual slide components only author **Layer 4**:

```
Layer 8 ─ Controller HUD Pill       fixed top:32 right:32 (z-index: 60)
Layer 7 ─ Bottom Dot Pagination    fixed bottom:32 (z-index: 50)
Layer 6 ─ Top Progress Bar         4px fixed top:0 (z-index: 50)
Layer 5 ─ Brand Logo (Top-Left)    48px tall SVG/PNG (z-index: 45)
Layer 4 ─ Slide Content Component   1920×1080 stage content (z-index: 10)
Layer 3 ─ Radial Spotlight Aura    centered behind focal point (z-index: 3)
Layer 2 ─ Decorative Outline Icons 10–14 Lucide icons, 10% opacity (z-index: 2)
Layer 1 ─ Cross-Hatch Grid Mesh    48px cells, 2.5% opacity (z-index: 1)
Layer 0 ─ Slide Canvas Ground      hsl(var(--background)) (z-index: 0)
```

---

## 7. Bottom Pagination Dot Row

Fixed at `bottom: 32px`, centered horizontally across the viewport:
- **Container:** `<nav aria-label="Slide pagination">` at `fixed bottom-8 left-1/2 -translate-x-1/2 z-50 flex items-center gap-3`.
- **Max Width:** `min(1600px, calc(100vw - 96px))`.

| State | Dot Shape | Dimensions | Fill & Border |
|---|---|---|---|
| **Inactive** | Circle | `8×8px` | `hsl(var(--foreground-subtle) / 0.6)` |
| **Hover** | Circle | `8×8px` | `hsl(var(--foreground))` |
| **Active** | Expanded Pill | `28×8px` | `hsl(var(--pres-accent))` |
| **Visited** | Circle | `8×8px` | `hsl(var(--foreground-subtle) / 0.8)` |

- **Animation:** `all 200ms cubic-bezier(0.2, 0.8, 0.2, 1)`. Under `prefers-reduced-motion: reduce`, dimensions snap without interpolation.
- **Accessibility:** `<button aria-label="Go to slide {n}: {title}" aria-current={current === n ? "true" : undefined} />`.

---

## 8. Anti-Hallucination & Quality Verification Checklist

- [ ] Controller supports all 8 anchor positions with safe area `env(safe-area-inset-*)` offsets.
- [ ] Hover-reveal fades to `opacity: 0.15` after 2.5s idle; any mousemove restores `1.0`.
- [ ] Tooltip expansion direction is derived from position edge (bottom anchor opens top, top opens bottom).
- [ ] Sound synthesis uses native Web Audio API oscillators without external asset requests.
- [ ] First-run onboarding card mounts once and dismisses cleanly via click or keypress.
- [ ] Slide anatomy strictly preserves the 8-layer Z-stack with individual slide content residing exclusively in Layer 4.
- [ ] Deep link share formats `#slide-{current}` and falls back to clipboard copy + Sonner toast.
