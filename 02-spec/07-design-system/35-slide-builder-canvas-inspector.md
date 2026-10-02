# 35 — Slide Builder Mode, Interactive Canvas & Inspector Specification

> **/goal** Master and enforce the dual-store architecture, 7 visual layers, floating draggable builder panel, slide-type switching engine, interactive selection overlays, bounding box overrides, and multi-format exports of the Slide Presentation Builder Engine.
> **/learn** Master the separation of `useDeckStore` (persisted) and `useEditStore` (ephemeral), the 7 canvas stacking layers, draggable `BuilderPanel` physics, `convertSlideType` text-preservation mechanics across all 20 slide layouts, builder hotkeys (`B`/`E`, `Tab`, `Cmd+Z`, `1`–`7`), and headless print-ready PDF/handout exports.

**Version:** 4.1.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. System Overview & Architectural Role

The **Slide Builder Engine** empowers authors, presenters, and AI agents to visually configure presentations in real-time directly on the scaled `1920×1080` virtual canvas:
- Move, drag, and resize text blocks, bullet cards, and media plates with live bounding boxes.
- Instantly switch layout models across all 20 slide types while preserving existing headlines and body text.
- Reassign 10-step theme gradients, pill presets, and 9-cell alignment coordinates.
- Maintain a non-destructive undo/redo history stack without polluting presentation playback timers.

---

## 2. Decoupled Dual-Store Architecture

To eliminate unnecessary re-renders and guarantee that editing artifacts never leak into audience presentation view, state is decoupled into two discrete stores:

```
┌─────────────────────────────────┐     ┌─────────────────────────────────┐
│   useDeckStore (Persisted)      │     │    useEditStore (Ephemeral)     │
│ ─────────────────────────────── │     │ ─────────────────────────────── │
│ • deck: DeckData                │     │ • isEditMode: boolean           │
│ • activeSlideIndex: number      │     │ • selectedElementId: string     │
│ • upsertSlide(slide)            │     │ • activePanel: PanelType        │
│ • updateSlideTheme(themeId)     │     │ • undoStack: HistoryAction[]    │
│ • reorderSlides(from, to)       │     │ • redoStack: HistoryAction[]    │
│ • LocalStorage Sync: "deck-v1"  │     │ • activeBoxHover: string | null │
└─────────────────────────────────┘     └─────────────────────────────────┘
                 ▲                                       ▲
                 └───────────────────┬───────────────────┘
                                     │
                             applyEdit() Gateway
                                     │
                        ┌────────────────────────┐
                        │ Single Mutator Gateway │
                        └────────────────────────┘
```

1. **`useDeckStore` (Persisted):** Holds the authoritative presentation JSON tree. Persists changes to `localStorage` under a versioned key (`deck-v1`). Synchronizes cross-window edits via browser `storage` events.
2. **`useEditStore` (Ephemeral):** Manages interactive UI states, selected element IDs, drag coordinates, undo/redo stacks, and toolbars. Cleared on page refresh.

---

## 3. Visual Canvas Layer Stack (7 Discrete Layers)

Every element on the `1920×1080` canvas is assigned to one of 7 isolated stacking layers (`z-index` and CSS `isolation`):

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Layer 6: Floating Inspector & HUD (Controls, Theme Bar, Builder Sidebar)   │
├─────────────────────────────────────────────────────────────────────────────┤
│ Layer 5: Builder Selection Overlays (Bounding boxes, blue focus outlines)  │
├─────────────────────────────────────────────────────────────────────────────┤
│ Layer 4: Ink Annotation Layer (Live presenter drawing & highlighter paths) │
├─────────────────────────────────────────────────────────────────────────────┤
│ Layer 3: Live DOM Typography & Cards (Headings, bullet lists, pricing)     │
├─────────────────────────────────────────────────────────────────────────────┤
│ Layer 2: Media Plates (Photographic hero silhouettes, feathered masks)     │
├─────────────────────────────────────────────────────────────────────────────┤
│ Layer 1: Brand Watermarks & Ribbons (Concentric arcs, organic SVG waves)   │
├─────────────────────────────────────────────────────────────────────────────┤
│ Layer 0: Canvas Base Background (Pure white #FFFFFF, 10-step gradient stop)│
└─────────────────────────────────────────────────────────────────────────────┘
```

| Layer Index | Name | CSS Stacking | Architectural Scope |
|:---:|:---|:---|:---|
| **Layer 0** | Canvas Base | `z-index: 0` | Pure background fill or base theme gradient ($S_0$–$S_2$). |
| **Layer 1** | Watermarks | `z-index: 10` | Non-interactive SVG organic waves, concentric circles. |
| **Layer 2** | Media Plates | `z-index: 20` | Hero photos, team avatars with feathered gradient masks. |
| **Layer 3** | DOM Typography | `z-index: 30` | Live HTML text, bullet clusters, pricing cards, pills. |
| **Layer 4** | Ink Annotations| `z-index: 40` | Real-time freehand SVG pen/marker canvas paths. |
| **Layer 5** | Selection Box | `z-index: 50` | 2px solid cyan/blue bounding box, corner resize handles. |
| **Layer 6** | Inspector HUD | `z-index: 60` | Fixed floating toolbars, slide reordering list, modal menus. |

---

## 4. Floating Draggable Builder Panel (`BuilderPanel`)

The builder settings panel floats above the viewport, draggable via a top grip handle and minimizable:

```
┌───────────────────────────────────────────────────┐
│ [:: Grip] Builder Settings   [↶ Undo] [↷ Redo] [−]│
├───────────────────────────────────────────────────┤
│ Slide Type: [ 3-Point Bullets                 ▼ ] │
├───────────────────────────────────────────────────┤
│ ┌───────────────┐ ┌───────────────┐ ┌───────────┐ │
│ │ Layout & Opts │ │ Typography    │ │ Gradients │ │
│ └───────────────┘ └───────────────┘ └───────────┘ │
│ ┌───────────────┐ ┌───────────────┐ ┌───────────┐ │
│ │ Insert Icon   │ │ Insert Image  │ │ Camera XY │ │
│ └───────────────┘ └───────────────┘ └───────────┘ │
├───────────────────────────────────────────────────┤
│ Active Selection: "slide-headline-title"          │
│ Position: X: 140px, Y: 210px, W: 900px            │
└───────────────────────────────────────────────────┘
```

- **Panel Dimensions:** Fixed width `264px`, `maxHeight: calc(100vh - 32px)`, default position `left: 16px, top: 16px`.
- **Grip Drag Physics:** `pointerdown` on header grip records offsets; `pointermove` on window updates coordinates clamped to viewport boundaries (`0 <= x <= window.innerWidth - 264`). Coordinates persist in `useBuilderPanelStore`.
- **Sub-Panel Modules:**
  1. **`SlideTypeMenu`:** Dropdown selector to switch between all 20 slide types.
  2. **`SlideOptionsPanel`:** Manages 9-cell text alignment, slide dark/light theme override, and speaker notes.
  3. **`FontPanel`:** Visual sliders for font scale, letter tracking, line height, and heading weight.
  4. **`GradientPanel`:** Toggle gradient fill, switch `linear` / `radial`, angle slider (0–360°), and add/remove color stops.
  5. **`InsertIconPanel`:** Modal browser containing the entire Lucide vector icon catalog with search and color tinting.
  6. **`InsertImagePanel` & `ImageLayer`:** File upload, URL input, crop rectangle, scale slider, and alignment grid.
  7. **`StepCameraPanel` & `XYPad`:** Interactive 2D pad configuring camera focal point pan (`x, y`) and zoom level (`1.0` to `2.5×`) for multi-step reveals.
  8. **`HistoryPanel`:** Visual audit ledger showing timestamped changes with instant rollback buttons.

---

## 5. Slide-Type Switching Engine (`convertSlideType`)

When an author changes the slide type in the builder, `convertSlideType` seamlessly maps existing content into the target structure without loss of critical narrative copy:

```typescript
export function convertSlideType(slide: Slide, targetType: SlideType): Slide {
  const base: BaseSlide = {
    id: slide.id,
    type: targetType,
    title: slide.title,
    theme: slide.theme,
    align: slide.align ?? "center-left",
    notes: slide.notes,
  };

  const primaryText = slide.title || "Headline";
  const secondaryText = extractSecondaryText(slide);

  switch (targetType) {
    case "bullets":
      return {
        ...base,
        subtitle: secondaryText || "Key strategic initiatives",
        bullets: [
          { icon: "Check", title: "Primary Objective", body: "Direct impact milestone" },
          { icon: "Users", title: "Team Alignment", body: "Cross-functional execution" },
          { icon: "Zap", title: "Velocity Metric", body: "Sub-millisecond performance" },
        ],
      };
    case "center":
      return {
        ...base,
        subhead: secondaryText || "Declarative mission statement",
      };
    case "counter-stat":
      return {
        ...base,
        statNumber: "99.9%",
        statLabel: primaryText,
        statDelta: "+45% YoY",
      };
    default:
      return { ...base, description: secondaryText };
  }
}
```

---

## 6. Interactive Selection Overlay & Overrides

When an element is selected on the live canvas:
1. **Overlay Geometry:** `SelectionOverlay` renders a bounding box with `border: 2px solid #0ea5e9`, background tint `rgba(14, 165, 233, 0.05)`.
2. **Transform Handles:** 4 corner nodes (`8×8px`, white fill, blue border) and 4 edge midpoint handles.
3. **Floating Metric Tag:** Position pill attached to top-left of the bounding box displaying `x: 140 | y: 210 | w: 900`.
4. **Coordinate Persistence:** Dragging writes explicit pixel overrides to `slide.boxes[elementId]` (`{ x, y, width, height }`).

---

## 7. Key Actions & Hotkey Navigation Matrix

| Hotkey | Action Payload | Context Guard |
|:---|:---|:---|
| **`B`** or **`E`** | Toggle Builder Mode ON / OFF | Suppressed when typing in text fields |
| **`Tab`** | Focus next editable element on slide | Active when builder mode is ON |
| **`Shift + Tab`** | Focus previous editable element | Active when builder mode is ON |
| **`Cmd/Ctrl + Z`** | Undo last canvas modification | Pops from `undoStack`, pushes to `redoStack` |
| **`Cmd/Ctrl + Shift + Z`** | Redo canvas modification | Pops from `redoStack`, pushes to `undoStack` |
| **`Escape`** | Deselect active element / Close modal | Blurs active focus ring |
| **`1`** | Quick-switch to Royal Violet theme | Applies to current slide or deck |
| **`2`** | Quick-switch to Corporate Gold theme | Applies to current slide or deck |
| **`3`** | Quick-switch to Enterprise Blue theme | Applies to current slide or deck |
| **`4`** | Quick-switch to Clinical Emerald theme | Applies to current slide or deck |
| **`5`** | Quick-switch to Sunset Crimson Orange theme | Applies to current slide or deck |
| **`6`** | Quick-switch to Monochrome Paper theme | Applies to current slide or deck |
| **`7`** | Quick-switch to Cyan Tech theme | Applies to current slide or deck |
| **`F`** | Toggle Fullscreen Mode | Invokes `requestFullscreen()` |

---

## 8. Multi-Format Presentation Exports

1. **Headless High-Resolution PDF Print (`slides.print.tsx`):**
   - Headless Chromium navigates to `/slides/print?theme=light&reducedMotion=true`.
   - Renders 1 slide per page in landscape 16:9 (`1920×1080`) format with zero animation delay.
2. **3-Up Executive Handout (`slides.handout-3up.tsx`):**
   - Formats 3 consecutive slides on the left column with structured blank ruled lines on the right column for executive note-taking.
3. **Deck Manifest JSON Export:**
   - One-click export downloading the entire deck JSON schema including all slide contents, themes, gradient stops, and bounding box overrides.

---

## 9. Anti-Hallucination & Quality Verification Checklist

- [ ] State architecture strictly maintains dual-store separation (`useDeckStore` vs `useEditStore`).
- [ ] Stacking context adheres strictly to the 7 defined visual layers.
- [ ] Floating builder panel is draggable, clamped to viewport, and minimizable.
- [ ] `convertSlideType` successfully maps primary and secondary text across all 20 slide types without null values.
- [ ] Hotkeys `B` and `E` toggle builder mode; `Escape` clears selection.
- [ ] PDF export forces 1920×1080 dimensions with background graphics enabled and animations bypassed.
