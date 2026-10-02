# Button System & Interactive Action Component Specification

> **/goal** Define the complete, production-grade button architecture across enterprise websites, web applications, and slide presentation systems.
> **/learn** Master the 6 standard button variants, 4 sizing scales, magnetic cursor physics, hardware-accelerated CSS3 interactions (`shine-sweep`, `pointer-fill`, `gradient-ring`), presentation controller action buttons, header shine pills, and capsule pill buttons.

**Version:** 4.1.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. System Overview

Buttons represent the primary interactive bridge between user intent and interface response. This specification codifies:
1. **6 Semantic Variants:** `primary`, `solid`, `outline`, `glass`, `ghost`, and `link`.
2. **4 Proportional Sizing Scales:** `sm` (36px), `md` (44px), `lg` (52px), and `icon` (44×44px).
3. **Magnetic Cursor Physics:** Optional cursor-tracking pull (`strength: 0.22`) for high-conversion primary CTAs.
4. **CSS3 Micro-Interactions:** Non-blocking hardware-accelerated visual feedback including shine sweeps, radial cursor fills, and gradient borders.
5. **Presentation Controller Action Buttons:** Circular 40×40px floating HUD action triggers with tabular numerals counter.
6. **Capsule & Status Pill Buttons:** High-contrast presentation telemetry chips, filter pills, and gradient badges.

---

## 2. Core Button Variants (Tokens & Utility Mapping)

| Variant | Visual Architecture | Resting Treatment | Hover & Focus State |
|:---|:---|:---|:---|
| **`primary`** | Gradient Accent Fill + Sheen | `bg-[image:var(--gradient-accent)] text-white shadow-[var(--shadow-card)]` | `.shine-sweep` active, `hover:bg-[image:var(--gradient-accent-hover)] hover:shadow-[var(--shadow-lift)]` |
| **`solid`** | Single Brand Ink Fill | `bg-brand-primary text-white shadow-[var(--shadow-card)]` | `.shine-sweep` active, `hover:bg-[image:var(--gradient-accent-hover)] hover:shadow-[var(--shadow-lift)]` |
| **`outline`** | Interactive Border + Radial Pointer Fill | `border border-border bg-transparent text-foreground` | `.pointer-fill` active, `hover:border-[var(--brand-secondary)] hover:text-white` |
| **`glass`** | Frosted Acrylic + Gradient Border | `bg-[image:var(--gradient-card-dark)] text-on-dark backdrop-blur-md` | `.shine-sweep gradient-ring` active, `hover:bg-white/10` |
| **`ghost`** | Zero-Chrome Utility | `bg-transparent text-foreground` | `hover:bg-muted text-foreground` |
| **`link`** | Text Anchor with Animated Underline | `text-brand-secondary px-0` | `hover:underline underline-offset-4` |

---

## 3. Sizing & Geometry Scale

| Size ID | Height | Horizontal Padding | Font Size | Icon Dimensions | Touch Target |
|:---|:---|:---|:---|:---|:---|
| **`sm`** | `36px` (`h-9`) | `16px` (`px-4`) | `13px` (`text-[13px]`) | `14px` (`size-3.5`) | `36px` |
| **`md`** | `44px` (`h-11`) | `20px` (`px-5`) | `14px` (`text-sm`) | `16px` (`size-4`) | `44px` |
| **`lg`** | `52px` (`h-13`) | `28px` (`px-7`) | `15px` (`text-[15px]`) | `18px` (`size-4.5`) | `52px` |
| **`icon`** | `44px` (`h-11`) | `0px` (`w-11`) | `14px` (`text-sm`) | `20px` (`size-5`) | `44px` |

---

## 4. Hardware-Accelerated CSS3 Button Interactions

### 4.1 The Continuous Sheen Sweep (`.shine-sweep`)
Primary and solid CTAs feature an angled 45-degree light sweep running across the surface on user interaction. To conserve device resources, it animates only when hovered or focused:

```css
@keyframes shine-sweep-action {
  from { background-position: 200% 0; }
  to { background-position: -200% 0; }
}

.shine-sweep { isolation: isolate; position: relative; overflow: hidden; }
.shine-sweep::after {
  content: ""; position: absolute; inset: 0; z-index: -1; border-radius: inherit; pointer-events: none;
  background-image: linear-gradient(45deg, transparent 25%, color-mix(in oklab, var(--brand-highlight, #38bdf8) 45%, transparent) 50%, transparent 75%, transparent 100%);
  background-size: 250% 250%; background-position: 200% 0;
}
.shine-sweep:hover::after, .shine-sweep:focus-visible::after, .group:hover .shine-sweep::after {
  animation: shine-sweep-action 3s linear infinite;
}
@media (prefers-reduced-motion: reduce) {
  .shine-sweep:hover::after, .shine-sweep:focus-visible::after { animation: none; }
}
```

### 4.2 Radial Cursor Pointer-Fill (`.pointer-fill`)
Outline buttons fill smoothly from the exact angle of cursor entry:

```css
.pointer-fill { position: relative; overflow: hidden; --pointer-fill: var(--brand-secondary, #2563eb); }
.pointer-fill::before {
  content: ""; position: absolute; inset: 0; z-index: 0; border-radius: inherit;
  background-color: var(--pointer-fill); opacity: 0; transform: scale(0.95);
  transition: opacity 240ms cubic-bezier(0.16, 1, 0.3, 1), transform 240ms cubic-bezier(0.16, 1, 0.3, 1); pointer-events: none;
}
.pointer-fill:hover::before { opacity: 1; transform: scale(1); }
```

### 4.3 Magnetic Cursor Pull Physics
Wrap primary hero CTAs in `<Magnetic strength={0.22}>`:

```typescript
export function Magnetic({ children, strength = 0.22 }: { children: React.ReactNode; strength?: number }) {
  const ref = useRef<HTMLDivElement>(null);
  const [pos, setPos] = useState({ x: 0, y: 0 });
  const handleMouseMove = (e: React.MouseEvent) => {
    const node = ref.current;
    if (!node) return;
    const { clientX, clientY } = e;
    const { left, top, width, height } = node.getBoundingClientRect();
    setPos({ x: (clientX - (left + width / 2)) * strength, y: (clientY - (top + height / 2)) * strength });
  };
  return (
    <motion.div ref={ref} onMouseMove={handleMouseMove} onMouseLeave={() => setPos({ x: 0, y: 0 })}
      animate={{ x: pos.x, y: pos.y }} transition={{ type: "spring", damping: 15, stiffness: 150, mass: 0.1 }}>
      {children}
    </motion.div>
  );
}
```

---

## 5. Header Shine Button & Floating Spark Particles (`ShineButton`)

The `ShineButton` is the primary high-conversion header affordance. It reads as the highest-contrast element on a light surface with a slow continuous 45-degree sheen sweep and an orbiting 8-particle floating spark constellation:

```
                  *  .   * (Spark Constellation)
┌───────────────────────────────────────────────┐
│  Got an idea?                        [✓ User] │ ──► /book-consultation
└───────────────────────────────────────────────┘
  ▲── Pill with 45-deg continuous sheen sweep ──▲
```

### 5.1 Architecture & Two-Level Span Hierarchy
The spark layer sits on the wrapper, NOT inside the button, because the button uses `overflow-hidden` to clip its sheen:
1. **Outer Positioning Wrapper:** `span.relative.inline-flex` (bounds relative coordinates).
2. **Link / Button Body:** Pill shape (`rounded-xl md:rounded-2xl`, `px-3 py-2 md:px-4 md:py-3`), `overflow-hidden`.
3. **Continuous Sheen Sweep (`::before`):** `linear-gradient(45deg, transparent 25%, color-mix(in oklch, var(--brand) 50%, transparent) 50%, transparent 75%, transparent 100%)`, background-size `250% 250%`, animated via `@keyframes shine` over 3s linear infinite.
4. **Spark Layer (`hidden md:block`):** Anchored `absolute right-5 top-0 -z-10`. Two-level span hierarchy preserves horizontal offsets:
   - Outer span owns seeded X/Y placement, opacity, and scale via inline style.
   - Inner span (`block rounded-full bg-brand animate-spark`) owns diameter and `@keyframes spark-float` (`translateY(5px) -> translateY(-5px)` over 2.4s ease-in-out infinite alternate).

### 5.2 Deterministic Sparks Particle Array
The 8 particles use deterministic mathematical offsets rather than runtime randomness:

| # | Diameter | translateX | translateY | Opacity | Scale | Delay |
|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| **1** | `6px` | `13px` | `-13px` | `0.80` | `1.00` | `0s` |
| **2** | `5px` | `16px` | `-8px` | `0.71` | `0.89` | `0.4s` |
| **3** | `8px` | `9px` | `-11px` | `0.52` | `0.65` | `0.8s` |
| **4** | `5px` | `18px` | `-11px` | `0.62` | `0.78` | `1.2s` |
| **5** | `4px` | `3px` | `-5px` | `0.27` | `0.34` | `1.6s` |
| **6** | `6px` | `9px` | `-9px` | `0.40` | `0.50` | `2.0s` |
| **7** | `4px` | `17px` | `-7px` | `0.77` | `0.96` | `2.4s` |
| **8** | `5px` | `4px` | `-2px` | `0.12` | `0.15` | `2.8s` |

### 5.3 Variants: Dark vs Inverse
- **Dark (Default Header):** `border border-white/10 bg-night text-white shadow-[0_15px_30px_-5px_var(--night)]` with active sheen and sparks.
- **Inverse (In-Card Callout):** `bg-white text-ink shadow-[0_15px_35px_-18px_oklch(0_0_0/0.6)]`, static without sheen or spark layer.

---

## 6. Slide Presentation Controller Buttons

Presentation decks use floating 40×40px round action buttons inside a fixed viewport pill HUD (`top: 32px, right: 32px`, height: 56px, radius: 9999px):

| Control | Icon (Lucide) | Target Dimension | State Rules | Action Payload |
|:---|:---|:---|:---|:---|
| **Previous** | `ChevronLeft` (20px) | `40×40px` round | Disabled at slide 1 (`opacity: 0.35`) | Decrements slide, updates hash `#slide-{n-1}` |
| **Next** | `ChevronRight` (20px) | `40×40px` round | Disabled at last slide (`opacity: 0.35`) | Increments slide, updates hash `#slide-{n+1}` |
| **Counter** | Live DOM Text | Min width `64px` | Tabular numbers, `user-select: none` | Displays `{current} / {total}` in Poppins 500, 20px |
| **Share** | `Share2` (20px) | `40×40px` round | Deep link to current slide | Invokes `navigator.share` or clipboard copy + Sonner toast |
| **Theme** | `Palette` (20px) | `40×40px` round | Hover background `0.08` alpha | Toggles theme picker popover with live WCAG contrast |
| **Webcam** | `Video`/`VideoOff` | `40×40px` round | Cycles webcam PIP presets | Toggles floating presenter webcam overlay |
| **Builder** | `Pencil` (20px) | `40×40px` round | Hotkey: `B` or `E` | Activates in-canvas slide builder mode |
| **Fullscreen**| `Maximize2`/`Minimize2` | `40×40px` round | Hotkey: `F` | Executes `requestFullscreen()` or `exitFullscreen()` |

---

## 7. Capsule & Status Pill Buttons

Presentation decks and hardware telemetry showcases use tokenized capsule pills for status chips, filter selectors, and telemetry readouts:

| Token Name | Background Fill | Foreground Text | Border Stroke | Application |
|:---|:---|:---|:---|:---|
| `--capsule-gold` | `40 96% 48% / 0.20` | `42 100% 78%` | `40 96% 48% / 0.55` | Brand highlights, verified tags |
| `--capsule-ember` | `14 80% 57% / 0.18` | `14 90% 72%` | `14 80% 57% / 0.45` | High-energy callouts, warnings |
| `--capsule-cream` | `var(--cream)` | `var(--ink)` | `var(--cream)` | Solid light badge, executive tag |
| `--capsule-ink` | `var(--surface-2)` | `var(--cream)` | `var(--gold) / 0.55` | Elevated dark contrast chip |
| `--capsule-outline`| `transparent` | `var(--cream)` | `var(--gold) / 0.70` | Wireframe tag, secondary filter |
| **Gradient Violet**| `linear(265 75% 66% -> 280 80% 72%)` | `0 0% 4%` (Near-black) | `none` | Creative tags, AI models (WCAG AA) |
| **Gradient Rose** | `linear(340 80% 60% -> 352 85% 70%)` | `0 0% 4%` (Near-black) | `none` | High-priority warnings (WCAG AA) |
| **Gradient Sky** | `linear(205 85% 55% -> 195 90% 65%)` | `0 0% 4%` (Near-black) | `none` | Telemetry readouts, metrics (WCAG AA) |

---

## 8. Implementation: Universal Button Primitive (`AppButton.tsx`)

```tsx
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import type { ButtonHTMLAttributes } from "react";
import { Magnetic } from "./magnetic";
import { cn } from "@/lib/utils";

export const appButtonVariants = cva(
  "relative inline-flex select-none items-center justify-center gap-2 whitespace-nowrap rounded-[var(--radius-button,12px)] font-display font-medium transition-[transform,box-shadow,background-color,color] duration-[var(--dur-fast,240ms)] ease-[var(--ease-out,cubic-bezier(0.16,1,0.3,1))] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:pointer-events-none disabled:opacity-50 active:translate-y-[1px] [&_svg]:relative [&_svg]:z-10 [&_svg]:size-4 [&_svg]:shrink-0 [&>*]:relative [&>*]:z-10",
  {
    variants: {
      variant: {
        primary: "shine-sweep overflow-hidden text-white shadow-[var(--shadow-card)] bg-[image:var(--gradient-accent)] hover:bg-[image:var(--gradient-accent-hover)] hover:shadow-[var(--shadow-lift)]",
        solid: "shine-sweep overflow-hidden bg-brand-primary text-white hover:bg-[image:var(--gradient-accent-hover)] hover:shadow-[var(--shadow-lift)]",
        outline: "pointer-fill overflow-hidden border border-border bg-transparent text-foreground hover:border-[var(--brand-secondary)] hover:text-white [--pointer-fill:var(--brand-secondary)]",
        glass: "shine-sweep gradient-ring overflow-hidden bg-[image:var(--gradient-card-dark)] text-on-dark backdrop-blur-md hover:bg-white/10",
        ghost: "text-foreground hover:bg-muted",
        link: "text-brand-secondary underline-offset-4 hover:underline px-0",
      },
      size: { sm: "h-9 px-4 text-[13px]", md: "h-11 px-5 text-sm", lg: "h-13 px-7 text-[15px]", icon: "size-11 px-0" },
    },
    defaultVariants: { variant: "primary", size: "md" },
  },
);

export type AppButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & VariantProps<typeof appButtonVariants> & { asChild?: boolean; magnetic?: boolean };
export function AppButton({ className, variant, size, asChild = false, magnetic = false, ...props }: AppButtonProps) {
  const Comp = asChild ? Slot : "button";
  const node = <Comp className={cn(appButtonVariants({ variant, size }), className)} {...props} />;
  return magnetic ? <Magnetic strength={0.22}>{node}</Magnetic> : node;
}
```

---

## 9. Anti-Hallucination & Quality Verification Checklist

- [ ] Every CTA variant uses only the 6 specified names (`primary`, `solid`, `outline`, `glass`, `ghost`, `link`).
- [ ] No button hardcodes inline hex codes; all colors pull from semantic CSS variables.
- [ ] `ShineButton` separates outer placement wrapper from inner floating dot animation so horizontal offsets survive.
- [ ] Sizing strictly follows `sm: 36px`, `md: 44px`, `lg: 52px`, `icon: 44px`.
- [ ] Active state depresses subtly by `translate-y-[1px]` without jarring layout shift.
- [ ] Focus state displays `ring-2 ring-ring ring-offset-2` for accessibility.
- [ ] Controller HUD buttons have `40×40px` circular hit targets with tabular numerals counter.
- [ ] Capsule gradient buttons enforce near-black ink (`0 0% 4%`) to pass WCAG 2.2 AA contrast.
