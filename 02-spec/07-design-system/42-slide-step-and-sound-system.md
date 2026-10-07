# 42 — Slide Step System & Sound Synthesis Engine Specification

> **/goal** Architect, construct, and calibrate the declarative step presentation models (`StepTimelineSlide` vs `AdvanceStepSlide`) and the Web Audio API sound synthesis and playback engine with zero runtime hallucinations.
> **/learn** Master the step timeline layout with active focus row, right description column, capsule pills, sound playback triggers (`whoosh`, `click`, `fadeClick`, `zoom`, `fadeZoom`), runtime audio envelopes (attack/decay/release), procedural oscillator fallbacks, and ducking mechanics.

**Version:** 4.3.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Executive System Overview

Slide presentations require synchronized visual progression and auditory feedback to guide executive attention without cognitive fatigue. This specification codifies:
1. **Two Core Step Models:** `StepTimelineSlide` (holistic overview with focused row detail) vs `AdvanceStepSlide` (full-stage camera dolly per step).
2. **Declarative Step JSON Schema:** Bounded step structures containing labels, titles, subtitles, descriptions, capsule pills, and reveal CTAs.
3. **Web Audio API Singleton Sound Engine:** Audio asset table, procedural oscillator fallbacks, runtime gain envelopes, debounce gates, and sound ducking.
4. **Click-Reveal & Hotspot Contracts:** Stepwise disclosure via keyboard or pointer interaction.

---

## 2. Step Presentation Models (`StepTimelineSlide` vs `AdvanceStepSlide`)

| Dimension | `StepTimelineSlide` (Showcase Default) | `AdvanceStepSlide` (Camera Dolly) |
|:---|:---|:---|
| **Use When** | Presenting 3–6 sequential phases where audience must see the **complete chain at once** while focusing on one active phase. | Presenting a 3–7 step methodology where each step requires the entire stage. |
| **Stage Geometry** | Left: vertical chain of step cards (360px width). Right: large active description card (720px width). | Full stage canvas: camera translates along horizontal/vertical track. |
| **Focus Interaction** | Hover or Arrow keys focus step row; right description card morphs smoothly via `layoutId`. | Next/Prev navigates from step frame to step frame with camera zoom/pan. |
| **Audio Trigger** | `whoosh` on focus change; `click` on capsule or CTA activation. | `zoom` or `whoosh` on frame transition; `click` on navigation. |

---

## 3. Declarative Step JSON Schema

```jsonc
{
  "slideNumber": 3,
  "slideName": "engagement-process",
  "slideType": "StepTimelineSlide",
  "transition": "SlideIn",
  "textAnimation": "SlideUp",
  "enabled": true,
  "isClickReveal": false,
  "showBrandHeader": true,
  "showPresenterChip": true,
  "brandStrip": false,
  "titleStyle": "white",
  "titleShimmer": false,
  "sound": { "on": "focus", "kind": "whoosh", "volume": 0.5 },
  "content": {
    "eyebrow": "How We Work",
    "title": "Engagement Process",
    "steps": [
      {
        "label": "Step 1",
        "title": "Discovery & Intake",
        "subtitle": "Listen, audit, and align",
        "description": "Two-week intake covering stakeholder interviews, architecture audit, and scope alignment.",
        "capsule": { "text": "Week 1", "color": "gold" }
      },
      {
        "label": "Step 2",
        "title": "Architecture Strategy",
        "subtitle": "Frame the primary objective",
        "description": "We narrow to a single measurable objective: one page, one team, one core metric.",
        "capsule": { "text": "Weeks 2-3", "color": "ember" },
        "cta": { "text": "Inspect Architecture", "revealSlide": 7, "variant": "gold" }
      },
      {
        "label": "Step 3",
        "title": "Production Sprint",
        "subtitle": "Build and verify",
        "description": "Rapid component assembly, guideline compliance, and automated quality gates.",
        "capsule": { "text": "Weeks 4-8", "color": "gold" }
      }
    ]
  }
}
```

---

## 4. Visual Anatomy of `StepTimelineSlide`

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│ Eyebrow: HOW WE WORK              Header Title: Engagement Process                          │
│                                                                                             │
│ ┌──────────────────────────────────────┐  ┌───────────────────────────────────────────────┐ │
│ │ Step 1: Discovery         [Week 1]   │  │ Step 2 Focus Detail Card                      │ │
│ ├──────────────────────────────────────┤  │                                               │ │
│ │ ► Step 2: Strategy        [Wk 2-3]   │  │ "We narrow to a single measurable objective.  │ │
│ ├──────────────────────────────────────┤  │  One page, one team, one core metric."        │ │
│ │ Step 3: Production        [Wk 4-8]   │  │                                               │ │
│ └──────────────────────────────────────┘  │ [Inspect Architecture CTA →]                  │ │
│   ▲── Left Step List (360px) ──────────▲  └───────────────────────────────────────────────┘ │
│                                             ▲── Right Rich Description Column (720px) ────▲ │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

1. **Step Card Dimensions:** Height `84px`, border-radius `16px`, padding `16px 20px`.
2. **Resting State:** `bg-card/40 border border-border/40 opacity-70 text-foreground`.
3. **Focused State:** `bg-card border-2 border-primary shadow-lg scale-[1.02] opacity-100`.
4. **Connecting Track:** Vertical `2px` hairline running behind step numeral badges in `var(--border)`.
5. **Capsule Badges:** Small pill badges (`text-[11px] font-mono uppercase px-2.5 py-0.5 rounded-full`):
   - `gold`: `bg-amber-500/15 text-amber-300 border border-amber-500/30`.
   - `ember`: `bg-orange-500/15 text-orange-300 border border-orange-500/30`.

---

## 5. Web Audio API Sound Engine Playbook

The sound engine is a browser-native singleton initialized on the first user interaction (`pointerdown` or `keydown`), bypassing autoplay restrictions:

### 5.1 Sound Asset Table & Audio Parameters

| Sound Kind | Primary Asset URL | Duration | Vol | Attack | Release | Ducks Prev | Architectural Role |
|:---|:---|:---:|:---:|:---:|:---:|:---:|:---|
| **`whoosh`** | `/sounds/fade_swoosh_v2.mp3` | `350ms` | `0.45` | `0.06s` | `0.12s` | Yes | Step focus changes, timeline transitions |
| **`click`** | `/sounds/click.mp3` | `180ms` | `0.18` | `0.005s`| `0.06s` | No | Slide navigation, dot clicks, menu select |
| **`fadeClick`**| `/sounds/click.mp3` (re-enveloped) | `180ms` | `0.09` | `0.05s` | `0.18s` | No | Soft tap precursor before major reveal |
| **`zoom`** | `/sounds/zoom.mp3` | `500ms` | `0.55` | `0.04s` | `0.18s` | Yes | AdvanceStep camera zoom, card growth |
| **`fadeZoom`** | `/sounds/fade_zoom.mp3` | `450ms` | `0.40` | `0.04s` | `0.15s` | Yes | Capsule expand morph, modal emergence |
| **`pop`** | Procedural Oscillator Synth | `120ms` | `0.45` | `0.002s`| `0.08s` | No | Tiny badge reveals, toggle switches |

### 5.2 Sound Engine Singleton Architecture

```typescript
export class DeckSoundManager {
  private ctx: AudioContext | null = null;
  private buffers = new Map<string, AudioBuffer>();
  private activeSource: AudioBufferSourceNode | null = null;
  private lastPlayTime = 0;

  private initCtx() {
    if (!this.ctx && typeof window !== "undefined") {
      const AudioCtx = window.AudioContext || (window as unknown as { webkitAudioContext: typeof AudioContext }).webkitAudioContext;
      this.ctx = new AudioCtx();
    }
    if (this.ctx && this.ctx.state === "suspended") {
      this.ctx.resume();
    }
  }

  async play(kind: "whoosh" | "click" | "fadeClick" | "zoom" | "fadeZoom" | "pop", customVol?: number) {
    if (typeof window === "undefined" || document.hidden) return;
    if (localStorage.getItem("deck.sound.muted") === "true") return;

    const now = performance.now();
    if (now - this.lastPlayTime < 60) return; // 60ms debounce gate
    this.lastPlayTime = now;

    this.initCtx();
    if (!this.ctx) return;

    if (kind === "pop" || !this.buffers.has(kind)) {
      this.playProceduralSynth(kind, customVol);
      return;
    }

    const buffer = this.buffers.get(kind)!;
    const source = this.ctx.createBufferSource();
    const gain = this.ctx.createGain();
    source.buffer = buffer;

    // Apply envelope
    const vol = customVol ?? 0.45;
    gain.gain.setValueAtTime(0.001, this.ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(vol, this.ctx.currentTime + 0.05);
    gain.gain.exponentialRampToValueAtTime(0.001, this.ctx.currentTime + buffer.duration);

    source.connect(gain);
    gain.connect(this.ctx.destination);
    source.start();
  }

  private playProceduralSynth(kind: string, vol = 0.3) {
    if (!this.ctx) return;
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();
    const freq = kind === "click" ? 800 : kind === "whoosh" ? 220 : 540;
    osc.frequency.setValueAtTime(freq, this.ctx.currentTime);
    if (kind === "whoosh") {
      osc.frequency.exponentialRampToValueAtTime(80, this.ctx.currentTime + 0.25);
    }
    gain.gain.setValueAtTime(vol, this.ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.001, this.ctx.currentTime + 0.15);
    osc.connect(gain);
    gain.connect(this.ctx.destination);
    osc.start();
    osc.stop(this.ctx.currentTime + 0.15);
  }
}

export const deckSound = new DeckSoundManager();
```

---

## 6. Click-Reveal & Interactive Hotspot Contracts

1. **Step-Aware Slide Navigation:**
   - In step-based slides, pressing `ArrowRight` or `Space` advances the internal step index (`step = step + 1`) until `step === totalSteps - 1`. Only then does navigation proceed to the next slide.
   - Pressing `ArrowLeft` reverses the internal step index before navigating back.
2. **Hotspot Pins on Diagrams:**
   - Diameter: `32px` circular pin button.
   - Surface: `bg-primary/90 text-primary-foreground shadow-glow`.
   - Micro-Animation: Continuous breathing ring (`@keyframes pulse-ring`) with `2.4s` period.
   - Interaction: Clicking or hovering expands an anchored callout card (`280px` width) displaying technical specifications.

---

## 7. Anti-Hallucination & Quality Verification Checklist

- [ ] All step timeline layouts place step cards on the left (360px) and active description on the right (720px).
- [ ] Sound engine enforces 60ms debounce to prevent audio clipping during rapid keypresses.
- [ ] Mute state is persisted strictly in `localStorage` under `deck.sound.muted`.
- [ ] Procedural oscillator fallbacks activate cleanly when audio files are offline or unbuffered.
- [ ] Capsule badges strictly support `gold` and `ember` variants.
- [ ] Closed set: Sound kinds are `whoosh`, `click`, `fadeClick`, `zoom`, `fadeZoom`, and `pop`. *Anything else does not exist.*
