# 29 — Slide Navigation and Slide Builder

> **/goal** Specify the presenter HUD, the two transition families, the keys, and the slide builder.
> **/learn** Scripted decks cycle five transitions with fixed numbers. JSON decks use fade and camera-zoom. Colors switch in `40-theme-switch.md`.

**Version:** 1.0.0
**Status:** Active

---

## 0. Anti-hallucination

Do not invent a key, a transition, or a chrome size. Website editing (wording, images, menu href) is `26-visual-builder.md` and must not be mounted on a slide stage.

---

## 1. HUD

Button geometry, the `40×40` hit targets, share, fullscreen, the `4px` progress bar, and the dot row are `31-slide-controller-buttons.md`. Do not restyle them here.

Position shared by that file: `top: 32px`, `right: 32px`, height `56px`, radius `9999px`, `z-index: 50`. A class `top-6 right-6` is `24px` on a default `4px` scale. New decks use `32px`.

Active dot is `28×8px`. Inactive dot is `8×8px`. The older `32×8` mark is not this controller.

Mount the HUD on the viewport, not inside the `1920×1080` transform.

---

## 2. Transitions

Pick one family per deck. Do not mix families on adjacent slides.

### 2.1 Scripted deck

Cycle by slide index, in this order only. Index `0` is `Slide`. Then `Fade`, `Zoom`, `Flip`, `Rise`, then repeat. Do not add a sixth name.

Every member uses the same timing:

| Name | Value |
|---|---|
| Duration | `0.45s` |
| Ease | `cubic-bezier(0.4, 0, 0.2, 1)` |
| Perspective, `Flip` only | `1200px` |

| Member | Enter, next slide | Exit, next slide |
|---|---|---|
| `Slide` | `x: 100%`, `opacity: 0` | `x: -100%`, `opacity: 0` |
| `Fade` | `opacity: 0` | `opacity: 0` |
| `Zoom` | `scale: 0.85`, `opacity: 0` | `scale: 1.15`, `opacity: 0` |
| `Flip` | `rotateY: 45deg`, `scale: 0.92`, `opacity: 0` | `rotateY: -45deg`, `scale: 0.92`, `opacity: 0` |
| `Rise` | `y: 80px`, `scale: 0.96`, `opacity: 0` | `y: -80px`, `scale: 0.96`, `opacity: 0` |

Going backward flips the sign of `x`, `y`, and `rotateY`. Zoom backward swaps `0.85` and `1.15`. The resting state is `x: 0`, `y: 0`, `opacity: 1`, `scale: 1`, `rotateY: 0`.

```css
.slide-motion {
  transition-duration: 0.45s;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
}
.slide-motion.is-flip {
  perspective: 1200px;
}
@media (prefers-reduced-motion: reduce) {
  .slide-motion {
    transition-duration: 0.01ms;
  }
}
```

A horizontal swipe starts a next or previous slide only after the finger travels `50px`. A shorter drag does nothing.

### 2.2 JSON deck

Allowed `TransitionKind` values:

```text
fade | camera-zoom | slide
```

The settings UI exposes `fade` and `camera-zoom` only. `slide` may exist on the type and must not appear as a third settings choice unless this file is revised. Default is `fade`.

---

## 3. Keys

Minimal deck subset. Camera groups, quick jump, sidebar, and the `/` shortcut map are **`42-slide-quiz-preview-chrome-and-default-shadows.md` section 5**. Do not duplicate that table here.

| Key | Action |
|---|---|
| `ArrowRight` or `Space` | Next step, then next slide |
| `ArrowLeft` | Previous step, then previous slide |
| `G` | Slide grid |
| `S` | Settings |
| `E` | Toggle the slide builder |
| `M` | Mute or unmute the deck click, only when a sound engine is already mounted. Do not invent a sound file. |
| `/` | Open keyboard shortcuts dialog (full matrix in file 42) |

Double-click a counter numeral to jump to that slide index. Do not bind `E` to the website visual builder. Bare **`1`** is reserved for presenter camera stage-fill when camera shortcuts are enabled (file 42).

---

## 4. Slide builder

`E` toggles builder chrome on the deck. The audience route does not load it.

- Reorder slides from the pagination pill.
- Edit the current slide in a `SlideOptionsPanel`: layout id from `28-slide-layouts.md`, theme id from `40-theme-switch.md`, and the slots that layout allows.
- Persist on the deck JSON. Do not write a second store for colors or fonts.
- Reject a layout id outside the closed union.
- Reject a theme id outside `40-theme-switch.md`.

The slide builder does not double-click text on the marketing page, does not upload a `2 MB` image through the website overlay, and does not edit a mega-menu href.

---

## 5. Still capture

Capture uses the active theme canvas from `40-theme-switch.md`. Do not paint the capture with a fixed hex.

| Name | Value |
|---|---|
| Standard scale | `1` |
| High-res scale | `2` |
| JPEG quality | `0.92` |
| PNG quality | `0.95` |
| Blank sample | `64` by `36` |
| Blank mean at or under | `14` |
| Blank standard deviation at or under | `8` |
| Settle before capture | `950ms` |
| Poll | `50ms` |
| Give-up | `5000ms` |
| Retries | `3` |
| First retry wait | `350ms` |
| Extra wait each retry | `250ms` |

A frame at or under both blank thresholds is a failed capture. Retry. Do not ship the blank frame.

Theme color changes are `40-theme-switch.md`. This file does not list palettes.

---

## 6. Checklist

- [ ] The deck uses one transition family. Scripted decks use section 2.1. JSON decks use section 2.2.
- [ ] The five scripted names stay in that order. Duration stays `0.45s`. Ease stays `cubic-bezier(0.4, 0, 0.2, 1)`.
- [ ] `Zoom` uses `0.85` and `1.15`. `Flip` uses `45deg`, `0.92`, and `1200px`. `Rise` uses `80px` and `0.96`.
- [ ] Reduced motion sets the duration to `0.01ms`.
- [ ] Swipe distance under `50px` does not change the slide.
- [ ] `M` only toggles an existing sound engine.
- [ ] Capture reads `var(--canvas)` from the active theme. A blank frame under the two thresholds is retried, then dropped.
- [ ] The slide builder still does not edit the marketing page.

---

## 7. Presenter extras

Webcam overlay size, mirror, and drag behavior stay in `24-slide-presentation-system.md`. Do not copy those pixel values into this file. Dual-screen notes and the next-slide preview stay there too.
