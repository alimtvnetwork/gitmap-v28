# 42 — Slide Quiz Preview Chrome & Default Shadow Pair

> **/goal** One copy-paste contract for default **text-shadow** and **box-shadow** tokens, quiz option cards, botanical-light colors, presenter shortcuts, and center-stage layout.
> **/learn** Faults and remediation: `02-spec/21-app/13-slide-quiz-preview-and-presenter-chrome/05-blind-agent-fault-register.md`. Verify colors: `node scripts/verify-botanical-light-contrast.mjs`.

**Version:** 2.0.0
**Status:** Active
**AI Confidence:** 100 within the closed scope in section 0.1

---

## 0. Anti-hallucination

Do not invent shadow rgba, shortcut keys, or green hex. Copy sections **1–2**, **3**, **5.1–5.3**, and **6** verbatim. If a key is not listed in section 5, do not bind it.

### 0.1 Closed scope (100% license)

| In scope at 100% | Out of scope (do not invent) |
|:---|:---|
| CSS in §1–2, §3, §4, §6 | Per-slide layouts in `34-slide-layout-catalog.md` |
| `PRESENTER_SHORTCUTS_CORE` + handlers §5.3–5.4 | Custom deck shortcut rows (§5.2 optional only) |
| HUD shortcuts button §5.5 + `31` §3.2.1 | New HUD geometry not in file 31 |
| Embed theme §9 | Ninth entry in `40-theme-switch.md` table |
| Contrast script exit 0 | Marketing pages, mega menu |

---

## 1. Semantic tokens (text + elevation)

| Token | Purpose |
|:---|:---|
| `--text-shadow-rest` | Body and option labels at rest |
| `--text-shadow-hover` | Option labels on card hover |
| `--elevation-rest` | Cards and images at rest |
| `--elevation-hover` | Cards and images on hover |
| `--elevation-selected` | Selected option glow |
| `--elevation-rest-quiz-plate` | Dark glass quiz shell (measured) |
| `--elevation-hover-quiz-plate` | Dark glass quiz shell hover |

### 1.1 Light surfaces (`:root`, `.light`, `[data-theme="botanical-light"]`, `[data-theme="green-choice"]`)

```css
:root,
.light,
[data-theme="botanical-light"],
[data-theme="green-choice"] {
  --text-shadow-rest: 0 1px 2px rgba(0, 0, 0, 0.06);
  --text-shadow-hover: rgba(0, 0, 0, 0.3) 1px 0.7px 0px;
  --elevation-rest: 0 1px 3px 0 rgba(0, 0, 0, 0.04), 0 1px 2px -1px rgba(0, 0, 0, 0.03);
  --elevation-hover: 0 10px 30px -4px hsl(var(--primary) / 0.24), 0 2px 8px -1px rgba(0, 0, 0, 0.35);
  --elevation-selected: 0 0 16px hsl(var(--primary) / 0.27);
}
```

**Source:** exam `theme.css` L62–73, L87–108 (light branch); hover primary alpha **0.24** is the spec improvement over **0.28**.

### 1.2 Dark text-shadow on options (measured)

```css
.dark,
[data-theme="vscode-dark"],
[data-theme="dracula"],
[data-theme="noir-gold"],
[data-theme="bright-gold-tech"] {
  --text-shadow-rest: 0 1px 4px rgba(0, 0, 0, 0.45), 0 2px 8px rgba(0, 0, 0, 0.25);
  --text-shadow-hover: #000000 1px 0.7px 0px;
  --elevation-rest: 0 1px 3px 0 rgba(0, 0, 0, 0.04), 0 1px 2px -1px rgba(0, 0, 0, 0.03);
  --elevation-hover: 0 10px 30px -4px hsl(var(--primary) / 0.28), 0 2px 8px -1px rgba(0, 0, 0, 0.35);
  --elevation-selected: 0 0 16px hsl(var(--primary) / 0.27);
}
```

**Source:** exam `theme.css` L91–100 for text-shadow; option-card elevation on dark decks reuses the same hover stack as light cards with `--primary` from the active deck.

### 1.3 Dark glass quiz plate (measured, optional class)

```css
.modern-quiz-card {
  box-shadow: var(--elevation-rest-quiz-plate);
  transition: box-shadow 250ms cubic-bezier(0.4, 0, 0.2, 1);
}
.modern-quiz-card:hover {
  box-shadow: var(--elevation-hover-quiz-plate);
}
.dark,
[data-theme="bright-gold-tech"] {
  --elevation-rest-quiz-plate: 0 10px 25px -5px rgba(0, 0, 0, 0.4), 0 8px 10px -6px rgba(0, 0, 0, 0.3);
  --elevation-hover-quiz-plate: 0 20px 30px -10px rgba(0, 0, 0, 0.5), 0 10px 15px -5px rgba(0, 0, 0, 0.35);
}
```

**Source:** exam `theme.css` L125–136.

### 1.4 Legacy aliases (mandatory one release)

Implementers MUST add this block so older exam CSS keeps working:

```css
:root {
  --option-text-shadow-rest: var(--text-shadow-rest);
  --option-text-shadow-hover: var(--text-shadow-hover);
}
.option-text-shadow {
  text-shadow: var(--text-shadow-rest);
}
```

Do not delete `--option-text-shadow-*` until all call sites read `--text-shadow-*`.

### 1.5 Utility classes

```css
.text-shadow-default {
  text-shadow: var(--text-shadow-rest);
  transition: text-shadow 200ms ease, color 200ms ease;
}
.elevation-default {
  box-shadow: var(--elevation-rest);
  transition: transform 220ms cubic-bezier(0.16, 1, 0.3, 1),
    box-shadow 220ms ease,
    border-color 200ms ease,
    background-color 200ms ease;
}
.elevation-default:hover,
.elevation-default:focus-visible {
  box-shadow: var(--elevation-hover);
}
img.elevation-default,
.media-plate.elevation-default {
  border-radius: var(--radius, 0.75rem);
}
img.elevation-default:hover {
  transform: translate3d(0, -2px, 0);
}
```

---

## 2. Presentation option card (normative)

```css
.presentation-option-card {
  opacity: 0.82;
  box-shadow: var(--elevation-rest);
  transition: transform 220ms cubic-bezier(0.16, 1, 0.3, 1),
    box-shadow 220ms ease,
    border-color 200ms ease,
    background-color 200ms ease,
    background 200ms ease,
    opacity 200ms ease;
  will-change: transform, opacity, box-shadow;
}
.presentation-option-card:hover {
  opacity: 1 !important;
  transform: translate3d(6px, 0, 0) !important;
  border-color: hsl(var(--primary) / 0.75) !important;
  background: linear-gradient(
    90deg,
    hsl(var(--primary) / 0.14) 0%,
    hsl(var(--card) / 0.92) 100%
  ) !important;
  box-shadow: var(--elevation-hover) !important;
}
.presentation-option-card:hover .option-text,
.presentation-option-card:hover .option-text-shadow {
  color: hsl(var(--foreground)) !important;
  font-weight: 600 !important;
  text-shadow: var(--text-shadow-hover) !important;
}
.presentation-option-card:hover .option-badge {
  border-color: hsl(var(--primary) / 0.75) !important;
  background-color: hsl(var(--primary) / 0.22) !important;
  color: hsl(var(--primary)) !important;
  transform: scale(1.05);
}
.option-text,
.option-text-shadow {
  text-shadow: var(--text-shadow-rest);
  transition: text-shadow 200ms ease, color 200ms ease;
}
@media (prefers-reduced-motion: reduce) {
  .presentation-option-card {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
    transform: none !important;
  }
}
```

Selected state (TSX):

```typescript
const selectedStyle = {
  opacity: 1,
  backgroundColor: "hsl(var(--card-active-bg))",
  border: "1px solid hsl(var(--primary) / 0.85)",
  boxShadow: "var(--elevation-selected)",
};
```

**Typography rule:** Body sentences use `color: hsl(var(--foreground))`. Use `hsl(var(--primary))` only for badges, icons, borders, and filled buttons—not paragraph text.

---

## 3. Theme `botanical-light` (improved green)

**Embed-only.** Does not appear in the 8-id deck switcher (`40-theme-switch.md` §7).

`[data-theme="green-choice"]` MUST mirror the same custom properties (legacy product id).

```css
[data-theme="botanical-light"],
[data-theme="green-choice"] {
  --primary: 142 70% 30%;
  --primary-foreground: 0 0% 100%;
  --background: 140 18% 97%;
  --foreground: 160 22% 10%;
  --card: 0 0% 100%;
  --card-foreground: 160 22% 10%;
  --secondary: 142 45% 93%;
  --secondary-foreground: 142 55% 28%;
  --muted: 150 14% 96%;
  --muted-foreground: 160 9% 42%;
  --border: 150 14% 90%;
  --ring: 142 70% 30%;
  --card-active-bg: 142 48% 94%;
  --radius: 0.75rem;
  color-scheme: light;
  background-color: hsl(var(--background));
  background-image:
    radial-gradient(ellipse 70% 50% at 50% -10%, hsl(142 70% 30% / 0.06), transparent 70%),
    radial-gradient(circle 500px at 100% 100%, hsl(160 35% 40% / 0.03), transparent 60%);
  background-attachment: fixed;
}
```

### 3.1 WCAG pairs (computed 2026-10-06)

Run `node scripts/verify-botanical-light-contrast.mjs` — exit **0** required before merge.

| Pair | Ratio | AA normal |
|:---|:---:|:---:|
| `--foreground` on `--background` | 15.96:1 | Pass |
| `--foreground` on `--card` | 16.91:1 | Pass |
| `--muted-foreground` on `--background` | 4.63:1 | Pass |
| `--primary` on `--background` (links/icons) | 4.61:1 | Pass |
| `--primary-foreground` on `--primary` (fills) | 4.88:1 | Pass |

Hex twins (review only, not for TSX):

| Role | Hex |
|:---|:---|
| Ground | `#F3F7F5` |
| Primary | `#166534` |
| Active fill | `#E8F5EC` |
| Body text | `#13201B` |
| Muted text | `#5F7369` |

---

## 4. Stagger entrance (optional)

```css
@keyframes slide-up-fade {
  from {
    opacity: 0;
    transform: translate3d(0, 12px, 0);
  }
  to {
    opacity: 1;
    transform: translate3d(0, 0, 0);
  }
}
.slide-up-anim {
  animation: slide-up-fade 420ms cubic-bezier(0.16, 1, 0.3, 1) both;
}
.stagger-1 { animation-delay: 0.06s; }
.stagger-2 { animation-delay: 0.12s; }
.stagger-3 { animation-delay: 0.18s; }
.stagger-4 { animation-delay: 0.24s; }
.stagger-5 { animation-delay: 0.3s; }
.stagger-6 { animation-delay: 0.36s; }
```

---

## 5. Presenter keyboard shortcuts

### 5.1 Types (copy first)

```typescript
export interface ShortcutItem {
  keys: string[];
  label: string;
}

export interface ShortcutGroup {
  group: string;
  items: ShortcutItem[];
}
```

### 5.2 Optional slide-specific groups

Append after the core export. **Not required for acceptance.** Replace labels when a deck has a profile grid slide:

```typescript
export const SLIDE_SPECIFIC_SHORTCUTS_TEMPLATE: ShortcutGroup[] = [
  {
    group: "Profile grid slide (example slot 22)",
    items: [
      { keys: ["1"], label: "Open profile A" },
      { keys: ["2"], label: "Open profile B" },
      { keys: ["Esc", "Enter", "→", "←"], label: "Return to grid" },
    ],
  },
];
```

When this template is mounted, **camera stage-fill MUST NOT bind bare `1`** on that slide index (use another key or disable stage-fill).

### 5.3 `PRESENTER_SHORTCUTS_CORE` (normative, paste entire file)

```typescript
export const PRESENTER_SHORTCUTS_CORE: ShortcutGroup[] = [
  {
    group: "Deck navigation",
    items: [
      { keys: ["→", "Space", "Enter"], label: "Next slide" },
      { keys: ["←", "Backspace"], label: "Previous slide" },
      { keys: ["F"], label: "Toggle fullscreen" },
      { keys: ["G"], label: "Toggle full-deck overview grid" },
      { keys: ["J"], label: "Toggle top slide jumper" },
      { keys: ["T"], label: "Toggle theme palette" },
      { keys: ["Esc"], label: "Close overlay / exit fullscreen" },
      { keys: ["/"], label: "Open this keyboard map" },
    ],
  },
  {
    group: "Deck builder",
    items: [
      { keys: ["E"], label: "Toggle slide builder" },
      { keys: ["B"], label: "Toggle slide builder (alias)" },
      { keys: ["S"], label: "Open settings panel" },
    ],
  },
  {
    group: "Quick jump",
    items: [
      { keys: ["2", "3", "4", "5", "6", "7", "8", "9"], label: "Type slide number digit" },
      { keys: ["Enter"], label: "Jump to typed number" },
      { keys: ["Backspace"], label: "Delete last digit" },
      { keys: ["Esc"], label: "Cancel pending jump" },
    ],
  },
  {
    group: "Sidebar",
    items: [
      { keys: ["Ctrl", "1"], label: "Toggle slide outline (⌘+1 on macOS)" },
      { keys: ["Esc"], label: "Close sidebar" },
    ],
  },
  {
    group: "Camera — power & surfaces",
    items: [
      { keys: ["I"], label: "Hard toggle camera" },
      { keys: ["M"], label: "Soft minimize / restore stream" },
      { keys: ["P"], label: "Enter camera fullscreen" },
      { keys: ["["], label: "Exit camera fullscreen" },
      { keys: ["]"], label: "Cinematic 3-state cycle" },
      { keys: ["1"], label: "Stage-fill toggle (when allowed — see §5.4)" },
      { keys: ["Esc"], label: "Exit fullscreen / stage" },
    ],
  },
  {
    group: "Camera — sizing & framing",
    items: [
      { keys: ["+"], label: "Step size up" },
      { keys: ["−"], label: "Step size down" },
      { keys: ["O"], label: "Toggle circle / rectangle frame" },
      { keys: ["H"], label: "Toggle vignette halo" },
    ],
  },
  {
    group: "Camera — fullscreen passthrough",
    items: [
      { keys: ["→", "↓", "Enter", "Space"], label: "Next slide while camera fullscreen" },
      { keys: ["←"], label: "Previous slide" },
      { keys: ["PageUp", "PageDown"], label: "Prev / next slide" },
    ],
  },
];

export const isFormFocus = (target: EventTarget | null): boolean => {
  if (!(target instanceof HTMLElement)) {
    return false;
  }

  const tagName = target.tagName;

  return tagName === "INPUT" || tagName === "TEXTAREA" || target.isContentEditable;
};
```

Export for UI:

```typescript
export const SHORTCUTS: ShortcutGroup[] = [
  ...PRESENTER_SHORTCUTS_CORE,
  ...SLIDE_SPECIFIC_SHORTCUTS_TEMPLATE,
];
```

Remove the template spread when the deck has no slide-specific binds.

### 5.4 Key dispatch order (mandatory)

Evaluate in order; **first match wins**:

1. **`isFormFocus(target)`** → ignore single-key shortcuts (modifiers may still run sidebar).
2. **Shortcut dialog open** → only **`Esc`** closes (plus dialog default).
3. **`isEditMode`** (slide builder ON, file 35) → **`B`** / **`E`** toggle builder; **`1`–`7`** apply builder theme quick-switch from file 35 §7; do **not** run camera stage-fill.
4. **`isSlideJumpPending`** → digits **`2`–`9`**, **`Enter`**, **`Backspace`**, **`Esc`** only.
5. **`isSlideSpecificShortcutActive`** (profile grid slide focused) → route **`1`–`4`** to slide handlers; camera stage-fill disabled.
6. **Camera enabled** → bare **`1`** toggles stage-fill; **`I`**, **`M`**, **`P`**, **`[`**, **`]`**, sizing keys per §5.3.
7. **Default** → deck navigation keys from §5.3.

### 5.5 `/` listener and HUD button

```typescript
useEffect(() => {
  const handler = (event: KeyboardEvent): void => {
    if (event.key !== "/" || event.shiftKey || event.ctrlKey || event.metaKey || event.altKey) {
      return;
    }

    if (isFormFocus(event.target)) {
      return;
    }

    event.preventDefault();
    setShortcutsOpen(true);
  };

  window.addEventListener("keydown", handler);

  return () => window.removeEventListener("keydown", handler);
}, []);
```

HUD control: **`31-slide-controller-buttons.md` §3.2.1** — same `setShortcutsOpen(true)`.

---

## 6. Center-stage layout

```css
.slide-center-stage {
  box-sizing: border-box;
  width: 1920px;
  height: 1080px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 96px 120px;
  text-align: center;
  gap: 24px;
}
.slide-center-stage .headline {
  max-width: 920px;
  font-weight: 700;
  font-size: clamp(40px, 4.2vw, 72px);
  line-height: 1.08;
  text-shadow: var(--text-shadow-rest);
}
.slide-center-stage .kicker {
  font-weight: 600;
  font-size: clamp(14px, 1.1vw, 18px);
  letter-spacing: 0.12em;
  text-transform: uppercase;
  text-shadow: var(--text-shadow-rest);
}
.slide-center-stage .subcopy {
  max-width: 720px;
  font-size: clamp(18px, 1.6vw, 28px);
  line-height: 1.45;
  color: hsl(var(--muted-foreground));
}
```

---

## 7. Cross-references

| Topic | File |
|:---|:---|
| Five transition names | `29-slide-navigation-and-builder.md` |
| HUD pill geometry | `31-slide-controller-buttons.md` |
| Builder stores | `35-slide-builder-canvas-inspector.md` |
| Eight deck themes | `40-theme-switch.md` |
| Embed quiz themes | §9 below + `40-theme-switch.md` §7 |
| Motion curves | `21-css3-animations-and-interactions.md` |

---

## 8. Blind-agent checklist

- [ ] Paste §1.1–1.5 and §2 CSS; add §1.4 legacy aliases.
- [ ] Paste §3 on `[data-theme="botanical-light"]` and `[data-theme="green-choice"]`.
- [ ] Run `node scripts/verify-botanical-light-contrast.mjs` (exit 0).
- [ ] Copy §5.3 to `shortcuts.ts`; wire §5.4 dispatch + §5.5 dialog.
- [ ] Mount HUD shortcuts button per file 31 §3.2.1.
- [ ] Pass AC-SQZ-001 through AC-SQZ-012 in app folder `04-acceptance-criteria.md`.

---

## 9. Quiz embed themes (not deck switcher)

### 9.1 DOM contract

```html
<div class="slide-stage" data-theme="bright-gold-tech">
  <div class="quiz-embed" data-theme="botanical-light">
    <!-- option cards, runner -->
  </div>
</div>
```

- **`data-theme` on `.slide-stage`:** one of the **8** ids in file 40.
- **`data-theme` on `.quiz-embed`:** `botanical-light` or `green-choice` only.
- Never set `botanical-light` on `.slide-stage`. Never increment the count in file 40 §1.

### 9.2 Inner root must re-declare shadows

Copy §1.1 custom properties onto `.quiz-embed[data-theme="botanical-light"]` so `--elevation-*` and `--text-shadow-*` exist inside the embed.
