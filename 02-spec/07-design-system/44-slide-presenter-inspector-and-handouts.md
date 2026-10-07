# 44 — Slide Presenter Inspector & Handout Export Specification

> **/goal** Architect, construct, and calibrate the dedicated speaker dual-slide Presenter Inspector console, the 3-up printable speaker handout engine, and the 1-up landscape presentation print mode.
> **/learn** Master the 60/40 split inspector layout, elapsed speech timers, speaker notes parsing, 3-up page chunking with faint ruled writing lines, final-state step thumbnail rendering, and automated print triggering via `@page` geometry.

**Version:** 4.3.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Executive System Overview

High-stakes keynote presentations demand specialized dual-screen visibility and physical collateral:
1. **Presenter Inspector Route (`/slides/inspector/$slideId`):** Standalone speaker cockpit displaying the live current slide (60%), upcoming slide preview (40%), scrollable speaker notes, step counter, and pacing timer.
2. **3-Up Speaker Handout Mode (`/slides/handout-3up`):** Printable landscape PDF format grouping 3 slides per page with rendered visual thumbnails, notes, and ruled writing lines for physical note-taking.
3. **1-Up Landscape Print Mode (`/slides/print`):** Full-bleed 1920×1080 per-page PDF generator stripping all interactive HUD chrome.

---

## 2. Dedicated Presenter Inspector Cockpit

The Presenter Inspector mounts at `/slides/inspector/$slideId` and `/slides/inspector/$slideId/$step`. It exists as an isolated route branch outside the standard slide shell to prevent mounting duplicate HUD pills or audio engines:

```
┌───────────────────────────────────────────────────┬─────────────────────────────────────────┐
│ Current Slide (60% Width, Live ScaledSlide)       │ Next Slide Preview (40% Width, Scaled)  │
│                                                   │                                         │
│ • Full interactive stage                          │ • Non-interactive thumbnail             │
│ • Live animation and step reveals                 │ • Preview of immediate upcoming content │
├───────────────────────────────────────────────────┼─────────────────────────────────────────┤
│ Speaker Notes Panel (Scrollable)                  │ Session Telemetry & Pacing Clock        │
│                                                   │                                         │
│ • Formatted markdown & delivery cues              │ • Step Counter: Step 2 of 4             │
│ • White-space preservation                        │ • Elapsed Speech Timer: 14:32           │
│ • Key bullet points & transition reminders        │ • Pacing Status: On Target (±0:45)      │
└───────────────────────────────────────────────────┴─────────────────────────────────────────┘
```

### 2.1 Inspector Keyboard Contracts & Timers

| Key | Action | Behavior |
|:---:|:---|:---|
| **`I`** | Launch Inspector | Opens `/slides/inspector/{current}` in a dedicated second window or tab |
| **`Space` / `→`** | Advance Step/Slide | Synchronously advances both speaker console and audience screen |
| **`←`** | Previous Step/Slide | Reverses active step or moves to preceding slide |
| **`R`** | Reset Timer | Resets elapsed speech stopwatch back to `00:00` |
| **`P`** | Pause Timer | Freezes or resumes elapsed speech stopwatch |
| **`Esc`** | Exit Inspector | Navigates back to the main presentation view (`/slides/{current}`) |

### 2.2 Speech Stopwatch & Pacing Clock
- **Timer Persistence:** Session start timestamp is stored in `localStorage` under `deck.inspector.startedAt`. Refreshing the inspector tab preserves the running talk duration.
- **Pacing Alerts:** When target slide duration is configured in slide metadata (`slide.targetDurationSec`), the timer chip transitions:
  - Normal: `text-foreground`
  - Approaching limit (within 30s): `text-amber-400 font-semibold`
  - Over time: `text-red-400 font-bold animate-pulse`

---

## 3. 3-Up Speaker Handout Export Engine

Mounts at `/slides/handout-3up`. Formats the entire slide deck into a physical or PDF handout suitable for executive workshops:

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│ ┌────────────────────────┐  Slide 01 — Executive Architecture Overview                     │
│ │ Slide 1 Thumbnail      │  Speaker Notes: Frame the core business problem before presenting│
│ │ (16:9 Scaled Preview)  │  the technical architecture. Emphasize multi-region redundancy.  │
│ └────────────────────────┘  ─────────────────────────────────────────────────────────────── │
│                             ─────────────────────────────────────────────────────────────── │
│ ┌────────────────────────┐  Slide 02 — Data Ingestion & Transformation                      │
│ │ Slide 2 Thumbnail      │  Speaker Notes: Walk through streaming batch pipelines.         │
│ │ (16:9 Scaled Preview)  │  ─────────────────────────────────────────────────────────────── │
│ └────────────────────────┘  ─────────────────────────────────────────────────────────────── │
│ ┌────────────────────────┐  Slide 03 — Deployment & Verification                            │
│ │ Slide 3 Thumbnail      │  Speaker Notes: Highlight automated zero-failure completion gates│
│ │ (16:9 Scaled Preview)  │  ─────────────────────────────────────────────────────────────── │
│ └────────────────────────┘  ─────────────────────────────────────────────────────────────── │
│ Page 1 of 8 ──────────────────────────────────────────────────────────── Confidential Handout│
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

### 3.1 Handout Architectural Rules
1. **Source of Truth:** Filters enabled slides (`enabled !== false`) from the active deck store in sequential order.
2. **Page Chunking:** Slides are partitioned into groups of exactly three. Each chunk renders a discrete `.handout-threeup-page`.
3. **Print Geometry:** Locked to 1920×1080 landscape:
   ```css
   @page {
     size: 1920px 1080px landscape;
     margin: 0;
   }
   .handout-threeup-page {
     width: 1920px;
     height: 1080px;
     page-break-after: always;
     display: grid;
     grid-template-rows: repeat(3, 1fr);
     padding: 48px 64px;
     gap: 24px;
   }
   ```
4. **Final-State Step Thumbnails:** In step-based or timeline slides, thumbnails render the final completed reveal state (`step = totalSteps - 1`) so printed handouts show all content.
5. **Ruled Writing Lines:** Below the speaker notes, faint ruled lines are rendered via CSS for handwritten attendee annotations:
   ```css
   .ruled-lines {
     background-image: repeating-linear-gradient(transparent, transparent 27px, hsl(var(--border) / 0.5) 28px);
     height: 84px;
   }
   ```
6. **Padding Incomplete Pages:** If the total slide count is not divisible by 3, the final page renders empty row placeholders (`.handout-threeup-row.is-empty`) to maintain stable row heights.
7. **Automated Print Trigger:** Appending `?auto=1` to the URL triggers `window.print()` after an initial 800ms layout hydration window.

---

## 4. 1-Up Presentation Print Mode (`/slides/print`)

Mounts at `/slides/print`:
1. Renders every enabled slide at full `1920 × 1080` resolution on its own dedicated PDF page.
2. Automatically suppresses all navigation bars, controller HUD pills, webcam bubbles, and floating chips using `[data-print-hide]` attributes:
   ```css
   @media print {
     [data-print-hide], .controller-hud, .webcam-overlay, .pagination-dots {
       display: none !important;
     }
     .slide-stage {
       page-break-after: always;
     }
   }
   ```

---

## 5. Anti-Hallucination & Quality Verification Checklist

- [ ] Presenter Inspector layout strictly partitions into 60% current slide, 40% next preview, notes, and telemetry.
- [ ] Speech stopwatch start time persists in `localStorage` under `deck.inspector.startedAt`.
- [ ] 3-Up Handout places exactly 3 slide rows per 1920×1080 page.
- [ ] Incomplete final handout pages use empty row placeholders to prevent vertical stretching.
- [ ] Handout thumbnails on step slides always render the final reveal state (`totalSteps - 1`).
- [ ] Ruled writing lines provide 28px baseline spacing for physical note-taking.
- [ ] All print routes automatically suppress HUD chrome, controller pills, and webcam bubbles.
