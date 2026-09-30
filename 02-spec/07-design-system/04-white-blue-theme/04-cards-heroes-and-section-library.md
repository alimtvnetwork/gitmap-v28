# White Blue Theme: Cards, Split Heroes & Enterprise Section Library

> **/goal** Document the complete card primitives (`SurfaceCard`, `GlassCard`, `NeuCard`, `card-premium`, `row-premium`), the Light-First Split Hero with Self-Assembling `CapabilityStack`, Flagship `SolutionsGrid`, Pinned Sticky-Note Workflow Board (`PushPin` + `NoteConnector`), Fluted Ribbed Glass `ScrollStack`, Tabbed Capabilities, Interactive Pricing, and Lead Form for the White Blue Theme (`04-white-blue-theme`).
> **/learn** Master every section pattern with exact 3-format color values (`HEX`, `RGB`/`RGBA`, `HSL` + `OKLCH`), CSS `@utility` definitions, and copy-pasteable React/TSX and LESS blueprints.

---

## 1. Card & Section Surface Color Matrix (Mandatory 3-Format Reference)

Every color and tint across the card primitives and section templates is documented in **HEX**, **RGB / RGBA**, **HSL / HSLA**, and **OKLCH**:

| Surface / Accent Element | CSS Token / Expression | Format 1: HEX | Format 2: RGB / RGBA | Format 3: HSL / HSLA | OKLCH Runtime Value |
|:---|:---|:---|:---|:---|:---|
| **`card-premium` Resting Ground** | `var(--card)` (`--paper`) | `#FFFFFF` | `rgb(255, 255, 255)` | `hsl(0, 0%, 100%)` | `oklch(1 0 0)` |
| **`card-premium` Hover Surface (4% Tint)** | `color-mix(in oklab, var(--primary) 4%, var(--card))` | `#F6F9FE` | `rgb(246, 249, 254)` | `hsl(218, 80%, 98%)` | `4%` `--blue-500` + `96%` `--paper` |
| **`card-premium` Hover Border (32% Tint)** | `color-mix(in oklab, var(--primary) 32%, var(--border))` | `#A6BEF0` | `rgb(166, 190, 240)` | `hsl(221, 71%, 80%)` | `32%` `--blue-500` + `68%` `--border` |
| **`row-premium` Hover Overlay (5% Tint)** | `color-mix(in oklab, var(--primary) 5%, var(--card))` | `#F4F7FE` | `rgb(244, 247, 254)` | `hsl(222, 83%, 98%)` | `5%` `--blue-500` + `95%` `--paper` |
| **Hero Top Wash Start** | `var(--blue-50)` | `#F0F6FF` | `rgb(240, 246, 255)` | `hsl(216, 100%, 97%)` | `oklch(0.972 0.014 259)` |
| **Workflow Ruled Backdrop Line** | `oklch(0.35 0.06 265 / 10%)` | `#2B38561A` | `rgba(43, 56, 86, 0.10)` | `hsla(222, 33%, 25%, 0.10)` | `oklch(0.35 0.06 265 / 10%)` |
| **Workflow Push-Pin #1 (Blue)** | `var(--blue-500)` | `#2563EB` | `rgb(37, 99, 235)` | `hsl(221, 83%, 53%)` | `oklch(0.546 0.215 262.9)` |
| **Workflow Push-Pin #2 (Violet)** | `var(--violet-500)` | `#822EE8` | `rgb(130, 46, 232)` | `hsl(267, 80%, 55%)` | `oklch(0.532 0.253 296.9)` |
| **Workflow Push-Pin #3 (Cyan)** | `var(--cyan-400)` | `#03D5E7` | `rgb(3, 213, 231)` | `hsl(185, 97%, 46%)` | `oklch(0.797 0.136 205.6)` |
| **Workflow Push-Pin #4 (Indigo)** | `var(--indigo-500)` | `#5856E9` | `rgb(88, 86, 233)` | `hsl(241, 77%, 63%)` | `oklch(0.545 0.216 277)` |
| **Workflow Push-Pin #5 (Teal)** | `var(--teal-500)` | `#009F9F` | `rgb(0, 159, 159)` | `hsl(180, 100%, 31%)` | `oklch(0.63 0.118 195)` |
| **Active Tab Pill Surface** | `color-mix(in oklab, var(--primary) 14%, transparent)` | `#2563EB24` | `rgba(37, 99, 235, 0.14)` | `hsla(221, 83%, 53%, 0.14)` | `oklch(0.546 0.215 262.9 / 14%)` |

---

## 2. Card Primitives & No-Jump Premium Surface System

### 2.1 Base Card Primitives (`SurfaceCard`, `GlassCard`, `NeuCard`, `Pill`)

- **`SurfaceCard`:** Rounded `20px` (`rounded-[var(--radius-card)]`), `border border-border bg-card p-6 shadow-[var(--shadow-card)]`, with `hover:-translate-y-1 hover:shadow-[var(--shadow-lift)]`.
- **`GlassCard`:** For `.band-dark` sections — uses `.gradient-ring`, `bg-[image:var(--gradient-card-dark)]`, `backdrop-blur-xl`, `text-on-dark` (`#D0D7E5`).
- **`NeuCard`:** Uses `.neu-raised p-6 text-ink` on `--neu-surface` (`#ECEFF5` / `rgb(236, 239, 245)` / `hsl(220, 31%, 94%)`).
- **`Pill`:** Rounded full (`rounded-full px-3 py-1 font-mono text-[11px] font-medium uppercase tracking-[0.12em]`) in three tones:
  - `default`: `border border-border bg-muted text-muted-foreground` (`#F4F6FA` fill, `#48536B` text)
  - `brand`: `bg-blue-50 text-blue-700` (`#F0F6FF` fill, `#173DA7` text)
  - `glass`: `gradient-ring bg-white/5 text-on-dark backdrop-blur-md`

### 2.2 `.card-premium`, `.row-premium` & `.index-fill` CSS Utilities

Instead of bouncing cards up and down across dense enterprise grids, `.card-premium` and `.row-premium` keep geometry locked while growing a `2px` vertical gradient bar (`scaleY(0 -> 1)` from `transform-origin: top`) along the left edge and warming the background:

```css
@utility card-premium {
  position: relative;
  isolation: isolate;
  overflow: hidden;
  border-radius: var(--radius-card);
  border: 1px solid var(--border);
  background: var(--card);
  box-shadow: var(--shadow-card);
  transition:
    background-color var(--dur-base) var(--ease-out),
    box-shadow var(--dur-base) var(--ease-out),
    border-color var(--dur-base) var(--ease-out);

  &::before {
    content: "";
    position: absolute;
    inset-block: 0;
    left: 0;
    width: 2px;
    z-index: 1;
    border-radius: inherit;
    background: var(--gradient-brand);
    transform: scaleY(0);
    transform-origin: top;
    transition: transform var(--dur-base) var(--ease-out);
    pointer-events: none;
  }

  &:hover,
  &:focus-visible,
  &:focus-within {
    border-color: color-mix(in oklab, var(--primary) 32%, var(--border));
    background-color: color-mix(in oklab, var(--primary) 4%, var(--card));
  }

  &:hover::before,
  &:focus-visible::before,
  &:focus-within::before {
    transform: scaleY(1);
  }
}

@utility row-premium {
  --row-bleed: 1.25rem;
  position: relative;
  isolation: isolate;
  transition: border-color var(--dur-base) var(--ease-out);

  &::after {
    content: "";
    position: absolute;
    inset-block: 0.25rem;
    inset-inline: calc(-1 * var(--row-bleed));
    z-index: -1;
    border-radius: var(--radius, 0.75rem);
    background-color: color-mix(in oklab, var(--primary) 5%, var(--card));
    opacity: 0;
    transition: opacity var(--dur-base) var(--ease-out);
    pointer-events: none;
  }

  &::before {
    content: "";
    position: absolute;
    inset-block: 0.25rem;
    left: calc(-1 * var(--row-bleed));
    width: 2px;
    z-index: 1;
    border-radius: 999px;
    background: var(--gradient-accent);
    transform: scaleY(0);
    transform-origin: top;
    transition: transform var(--dur-base) var(--ease-out);
    pointer-events: none;
  }

  &:hover::before,
  &:focus-visible::before,
  &:focus-within::before {
    transform: scaleY(1);
  }

  &:hover::after,
  &:focus-visible::after,
  &:focus-within::after {
    opacity: 1;
  }
}

@utility index-fill {
  transition:
    background var(--dur-fast) var(--ease-out),
    color var(--dur-fast) var(--ease-out);

  .group:hover &,
  .group:focus-visible &,
  .group:focus-within & {
    background-image: var(--gradient-accent);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
  }
}
```

---

## 3. Section Pattern 1: Light-First Split Hero & Self-Assembling `CapabilityStack`

### 3.1 Architectural Anatomy

- **Full-Bleed Top Underlap (`-mt-[72px] pt-[72px]`):** Pulls the hero canvas underneath the sticky `72px` header so the top wash (`linear-gradient(180deg, var(--blue-50), var(--paper) 72%)`, `#F0F6FF` -> `#FFFFFF`), `GridBackdrop` (`opacity-[0.6]`), and `Halo` glow shine through the translucent header.
- **2-Column Asymmetric Split (`lg:grid-cols-[1.35fr_0.95fr] xl:gap-x-16`):**
  - **Left Column:** Multi-line `MaskedHeading` (`as="h1"`, `clamp(38px, 3.6vw, 52px)` on desktop) with a `.gradient-text` accent word, `52ch` lead paragraph (`#48536B`), and `lg` primary + outline `WhiteBlueButton` pair.
  - **Right Column (`HeroCard` -> `TiltCard max={3}` -> `CapabilityStack`):**
    - **Header Bar:** Displays `"YOUR PLATFORM, ASSEMBLED"` in `font-mono text-[11px] uppercase tracking-[0.16em]` on the left and a live animated counter (`3/5 enabled`) on the right, underscored by a `1px` `--gradient-accent` progress bar (`width: ${(activeCount / totalCount) * 100}%`).
    - **Auto-Cycling `ToggleRow` Stack (`CYCLE_MS = 2200`, `RESUME_MS = 8000`):** Every `2200ms` while in view, the next capability row toggles its switch and sweeps a subtle specular sheen across the row. Clicking any `ToggleRow` pauses the auto-cycle for `8000ms` (`isPaused`) so the visitor has full manual control.
    - **Detached ROI Summary Strip:** A floating card below the stack (`mt-4 rounded-[var(--radius-card)] border border-border bg-card px-5 py-3.5 shadow-[var(--shadow-card)]`) pairing `"RETURN ON INVESTMENT (ROI)"` with `"200%"` in bold `font-display text-lg`.

### 3.2 Complete React/TSX Blueprint (`CapabilityStack`)

```tsx
import { useEffect, useRef, useState } from "react";
import { motion, useReducedMotion } from "motion/react";
import { ToggleRow } from "@/components/primitives/toggle-row";

const CYCLE_MS = 2200;
const RESUME_MS = 8000;

export interface CapabilityItem {
  id: string;
  title: string;
  sub: string;
  isDefaultEnabled: boolean;
}

export function CapabilityStack({ items }: { items: CapabilityItem[] }) {
  const isReducedMotion = useReducedMotion();
  const [enabledList, setEnabledList] = useState(() => items.map((item) => item.isDefaultEnabled));
  const [cursorIndex, setCursorIndex] = useState(0);
  const [isPaused, setIsPaused] = useState(false);
  const resumeTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const activeCount = enabledList.filter(Boolean).length;

  useEffect(() => {
    if (isReducedMotion || isPaused) {
      return;
    }

    const intervalId = setInterval(() => {
      setCursorIndex((currentCursor) => {
        setEnabledList((prev) =>
          prev.map((isItemOn, idx) => (idx === currentCursor ? isItemOn === false : isItemOn)),
        );

        return (currentCursor + 1) % items.length;
      });
    }, CYCLE_MS);

    return () => clearInterval(intervalId);
  }, [isReducedMotion, isPaused, items.length]);

  const handleUserToggle = (index: number) => {
    setEnabledList((prev) =>
      prev.map((isItemOn, idx) => (idx === index ? isItemOn === false : isItemOn)),
    );
    setIsPaused(true);

    if (resumeTimerRef.current) {
      clearTimeout(resumeTimerRef.current);
    }

    resumeTimerRef.current = setTimeout(() => setIsPaused(false), RESUME_MS);
  };

  return (
    <div className="relative w-full">
      <div
        aria-hidden
        className="pointer-events-none absolute -inset-8 -z-10 rounded-[var(--radius-media)] bg-[image:var(--gradient-halo)] opacity-70 blur-3xl"
      />
      <div className="overflow-hidden rounded-[var(--radius-media)] border border-border bg-card shadow-[var(--shadow-lift)]">
        <div className="border-b border-border">
          <div className="flex items-center justify-between px-5 py-3.5">
            <p className="font-mono text-[11px] uppercase tracking-[0.16em] text-muted-foreground">
              Your platform, assembled
            </p>
            <span className="font-mono text-[11px] tabular-nums text-muted-foreground">
              <strong className="text-foreground">{activeCount}</strong>/{enabledList.length} enabled
            </span>
          </div>
          <div className="h-px w-full bg-[color:var(--hairline)]">
            <motion.div
              className="h-px bg-[image:var(--gradient-accent)]"
              animate={{ width: `${(activeCount / enabledList.length) * 100}%` }}
              transition={{ duration: isReducedMotion ? 0 : 0.6, ease: [0.16, 1, 0.3, 1] }}
            />
          </div>
        </div>
        <div className="flex flex-col gap-2.5 p-4 sm:p-5">
          {items.map((cap, idx) => (
            <ToggleRow
              key={cap.id}
              title={cap.title}
              sub={cap.sub}
              isEnabled={enabledList[idx] ?? false}
              onToggle={() => handleUserToggle(idx)}
              tone="light"
            />
          ))}
        </div>
      </div>
    </div>
  );
}
```

---

## 4. Section Pattern 2: Sticky Editorial Header + Flagship `SolutionsGrid`

- **Grid Layout (`lg:grid-cols-[0.8fr_1.2fr] lg:gap-20`):**
  - **Left Column (`lg:sticky lg:top-28 lg:self-start`):** Remains pinned while the right column scrolls, displaying `Eyebrow`, `MaskedHeading`, `text-lead`, and a monospace footer kicker.
  - **Right Column Top — Flagship `SpotlightCard`:** Uses `group card-premium spotlight-light flex flex-col gap-4 p-8 sm:p-10`. Includes a top-right radial halo (`-right-20 -top-24 size-64 rounded-full bg-[image:var(--gradient-halo)] opacity-0 group-hover:opacity-100`), brand icon + `"MOST DEPLOYED"` pill badge, `text-h3` title, description, and `"Explore the platform ->"` link.
  - **Right Column Bottom — Numbered Typographic Rows (`focus-list`):** Each secondary solution renders as a `group row-premium grid grid-cols-[auto_1fr_auto] items-start gap-5 border-b border-border px-3 py-7`. The left column displays a `02`..`06` monospace index with `.index-fill` (filling with `--gradient-accent` on hover), while the vendor icon transitions from `opacity-70 grayscale` to `opacity-100 grayscale-0`.

---

## 5. Section Pattern 3: Pinned Sticky-Note Workflow Board (`Workflow`)

### 5.1 Architectural Anatomy

1. **Ruled Paper Backdrop (`ruled-backdrop-light`):** Renders horizontal 1px blue-slate rule lines every `48px` (`oklch(0.35 0.06 265 / 10%)`, `rgba(43, 56, 86, 0.10)`) masked by an elliptical radial fade (`mask-image: radial-gradient(ellipse at center, black, transparent 78%)`).
2. **Alternating Tilted Sticky Notes (`NOTE_STYLE`):**
   - Each note alternates horizontal alignment (`md:ml-auto md:mr-[4%] lg:mr-[8%]` vs. `md:mr-auto md:ml-[4%] lg:ml-[8%]`), initial tilt angle (`-4.5deg`, `5deg`, `-5deg`, `4deg`, `-4deg`), and token pin hue (`--blue-500` `#2563EB`, `--violet-500` `#822EE8`, `--cyan-400` `#03D5E7`, `--indigo-500` `#5856E9`, `--teal-500` `#009F9F`).
   - As each note scrolls toward the viewport center (`useCenterProgress`), a spring (`stiffness: 140, damping: 28`) untwists `rotate` from `tilt -> 0deg` and lifts `y` from `14px -> -6px`.
   - The inner note surface (`rounded-[18px]`) is tinted with `7%` of its pin color: `background: color-mix(in oklab, var(--pin) 7%, var(--card))`.
3. **3D Glossy `PushPin` (`pin-dome`, `pin-dome-gloss`, `pin-shaft`, `pin-contact`, `pin-halo`):**
   - Pure CSS 3D push-pin head pressed into the top center of each note (`-top-3 left-1/2 size-7 -translate-x-1/2 sm:-top-4 sm:size-8`), casting a colored halo bleed (`pin-halo`) and blurred elliptical contact shadow (`pin-contact`) onto the card.
4. **Scroll-Linked Dashed SVG `NoteConnector`:** Between adjacent notes on desktop (`lg:block`), a `230×64` SVG cubic Bezier curve (`M214 6 C 186 44, 96 14, 22 54`) animates its `pathLength` from `0.04 -> 1` using an SVG `<linearGradient>` bridging the current note's pin color (`from`) to the next note's pin color (`to`).

---

## 6. Section Pattern 4: Fluted Ribbed Glass `ScrollStack`

### 6.1 Sticky Card Deck & Ribbed Glass Media Panel

Each card in `ScrollStack` is wrapped in `.stack-slot` (`position: sticky; top: calc(var(--stack-top) + var(--i, 0) * var(--stack-peek))`) so subsequent cards slide up and stack cleanly over earlier cards with a `26px` desktop peek (`14px` on mobile) and progressive `4px` horizontal depth inset (`--inset: (total - 1 - index) * 4px`):

- **Dynamic Hue & Tint Ramp (`.stack-card`):** Reads `--hue` (e.g. `var(--blue-500)`) and `--tint` (`0`–`100`) to synthesize a `152deg` three-stop surface wash. Early cards use pale tints (`tint: 10`–`24`, dark ink text `#090E18`), while final cards ramp into deep saturated tones (`tint: 75`–`92`, `isDark: true`, crisp white text `#FFFFFF`).
- **Fluted Ribbed Glass Panel (`.stack-panel`):** The right-hand media frame layers a `30px` repeating vertical fluted glass shader (`repeating-linear-gradient(90deg, ...)` controlled by `--flute: 7` on light cards and `--flute: 11` on dark cards) behind a gently floating (`animate={{ y: [0, -10, 0] }}`, `6s` infinite) 3D isometric illustration.

```css
@utility stack-slot {
  --stack-peek: 14px;
  --stack-top: 64px;
  position: sticky;
  top: calc(var(--stack-top) + var(--i, 0) * var(--stack-peek));

  @media (width >= 64rem) {
    --stack-peek: 26px;
    --stack-top: 88px;
    margin-inline: var(--inset, 0px);
  }
}

@utility stack-panel {
  position: relative;
  overflow: hidden;
  border-radius: var(--radius-media);
  background-image:
    radial-gradient(62% 55% at 50% 44%, color-mix(in oklab, white 24%, transparent) 0%, transparent 72%),
    radial-gradient(120% 100% at 50% 50%, transparent 55%, color-mix(in oklab, black 14%, transparent) 100%),
    repeating-linear-gradient(
      90deg,
      color-mix(in oklab, white calc(var(--flute, 9) * 1%), transparent) 0px,
      transparent 12px,
      color-mix(in oklab, black calc(var(--flute, 9) * 1%), transparent) 30px,
      color-mix(in oklab, white calc(var(--flute, 9) * 1%), transparent) 30px
    ),
    linear-gradient(
      180deg,
      color-mix(in oklab, var(--hue, var(--blue-500)) calc(min(var(--tint, 12) + 14, 96) * 1%), white) 0%,
      color-mix(in oklab, var(--hue, var(--blue-500)) calc(min(var(--tint, 12) + 44, 100) * 1%), white) 58%,
      color-mix(in oklab, var(--hue, var(--blue-500)) calc(min(var(--tint, 12) + 68, 100) * 1%), black 10%) 100%
    );
}
```

---

## 7. Section Patterns 5, 6 & 7: Tabbed Capabilities, Interactive Pricing & Lead Form

### 7.1 Tabbed Capability Board (`Capabilities` on `.band-soft`)

- **Shared-Layout Sliding Pill Tablist:** Renders horizontal pill buttons (`rounded-full border border-hairline px-5 py-2.5 font-display text-sm font-semibold`). The active tab renders a `<motion.span layoutId="cap-pill" />` with `border border-primary/40 bg-[color-mix(in_oklab,var(--primary)_14%,transparent)]` (`rgba(37, 99, 235, 0.14)`) and spring physics (`stiffness: 380, damping: 34`).
- **Split Panel Content (`lg:grid-cols-[minmax(0,0.92fr)_minmax(0,1.08fr)]`):**
  - Left featured card (`p-8 md:p-10`): Top radial primary glow (`opacity-60`), monospace category pill, `clamp(24px, 2.4vw, 34px)` headline, bullet list, and CTA link.
  - Right `2×2` sub-card grid (`sm:grid-cols-2`): Each card (`p-6`) pairs a `44px` (`size-11 rounded-[14px]`) primary-tinted icon box on the left with a `01`..`04` monospace index on the right.

### 7.2 Interactive Scope-Builder Pricing (`Pricing`)

- Instead of static 3-column pricing cards, each engagement tier (`PlanBlock`) is a 2-column interactive configurator (`lg:grid-cols-[1.05fr_0.95fr] border-t border-border py-12`):
  - **Left Column:** Tier name (`text-h3`), description, core checkmark perks, and interactive `ToggleRow` add-ons (`"ADD TO THIS ENGAGEMENT"`).
  - **Right Quote Card (`rounded-[var(--radius-media)] border border-border bg-card p-8 shadow-[var(--shadow-card)]`):** Displays the price note, large `clamp(38px, 4.4vw, 58px)` display price, a live `"SCOPE IN THIS QUOTE"` manifest where toggled-off items appear with `line-through text-muted-foreground/60` (`"Not added"`) and toggled-on items read `"Included"`, followed by a full-width `lg` `WhiteBlueButton`.

### 7.3 Enterprise Lead Form Card (`LeadForm`)

- **Container:** `rounded-[var(--radius-card)] border border-hairline bg-card p-[clamp(20px,3vw,36px)] shadow-[var(--shadow-card)]`.
- **Field Labels:** `text-xs uppercase tracking-[0.14em] text-muted-foreground` with required asterisks in `text-brand-secondary` (`#2563EB`).
- **Inputs & Native Selects:** `h-11 w-full rounded-[var(--radius-button)]` (`12px` radius), `border border-input` (`#E2E6ED`), `bg-background` (`#FFFFFF`), `px-3 text-sm text-foreground`, with `focus-visible:ring-2 focus-visible:ring-ring` (`#2563EB`).

---

## 8. Parametric LESS Mixins for Premium Cards & Sticky Stack

```less
// ============================================================================
// WHITE BLUE THEME — CARDS & SECTION LIBRARY LESS MIXINS
// ============================================================================

.wb-card-premium() {
  position: relative;
  isolation: isolate;
  overflow: hidden;
  border-radius: 20px;
  border: 1px solid #E2E6ED;      // rgb(226, 230, 237) | hsl(218, 23%, 91%)
  background-color: #FFFFFF;      // rgb(255, 255, 255) | hsl(0, 0%, 100%)
  box-shadow: 0 8px 24px -12px rgba(13, 41, 117, 0.18);
  transition: background-color 420ms cubic-bezier(0.16, 1, 0.3, 1),
              border-color 420ms cubic-bezier(0.16, 1, 0.3, 1),
              box-shadow 420ms cubic-bezier(0.16, 1, 0.3, 1);

  &::before {
    content: "";
    position: absolute;
    inset-block: 0;
    left: 0;
    width: 2px;
    z-index: 1;
    background-image: linear-gradient(120deg, #1B51D3 0%, #721CD2 55%, #2563EB 100%);
    transform: scaleY(0);
    transform-origin: top;
    transition: transform 420ms cubic-bezier(0.16, 1, 0.3, 1);
  }

  &:hover,
  &:focus-within {
    border-color: #A6BEF0;        // rgb(166, 190, 240) | hsl(221, 71%, 80%)
    background-color: #F6F9FE;    // rgb(246, 249, 254) | hsl(218, 80%, 98%)

    &::before {
      transform: scaleY(1);
    }
  }
}
```
