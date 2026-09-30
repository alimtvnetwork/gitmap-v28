# White Blue Theme: Buttons, Motion Grammar & Micro-Interactions

> **/goal** Define the complete button system (`WhiteBlueButton`), CSS3 keyframe animations (`wb-shine`, `wb-text-sweep`, `wb-marquee`), pointer-tracking surfaces (`SpotlightCard`, `TiltCard`, `Magnetic`), scroll entrance primitives (`Reveal`, `StaggerGroup`, `MaskedHeading`), and interactive controls (`ToggleRow`, `DragRail`) for the White Blue Theme (`04-white-blue-theme`).
> **/learn** Master all 6 button variants, 4 button sizes, 3-format color states, hardware-accelerated sheen and bottom-up pointer fills, and reduced-motion safety guards.

---

## 1. Button & Interaction Color Matrix (Mandatory 3-Format Reference)

Every color used across `WhiteBlueButton`, `.shine-sweep`, `.pointer-fill`, `.gradient-text`, `.spotlight`, `.spotlight-light`, and `ToggleRow` is documented below in **HEX**, **RGB / RGBA**, **HSL / HSLA**, and **OKLCH**:

| Component / Variant | Token / CSS Property | Format 1: HEX | Format 2: RGB / RGBA | Format 3: HSL / HSLA | OKLCH Runtime Value |
|:---|:---|:---|:---|:---|:---|
| **`primary` Button Fill** | `--gradient-accent` | `#2563EB` -> `#822EE8` | `rgb(37, 99, 235)` -> `rgb(130, 46, 232)` | `hsl(221, 83%, 53%)` -> `hsl(267, 80%, 55%)` | `oklch(0.546 0.215 262.9)` -> `oklch(0.532 0.253 296.9)` |
| **`primary` / `solid` Hover Fill** | `--gradient-accent-hover` | `#4277EE` -> `#924AEB` | `rgb(66, 119, 238)` -> `rgb(146, 74, 235)` | `hsl(222, 83%, 60%)` -> `hsl(267, 80%, 61%)` | `88%` primary + `12%` white in oklab |
| **`solid` Button Resting Fill** | `--brand-primary` (`--blue-900`) | `#0D2975` | `rgb(13, 41, 117)` | `hsl(224, 80%, 25%)` | `oklch(0.317 0.135 264.5)` |
| **`outline` Button Resting Border** | `--border` (`--neutral-200`) | `#E2E6ED` | `rgb(226, 230, 237)` | `hsl(218, 23%, 91%)` | `oklch(0.924 0.011 265)` |
| **`outline` Button Pointer Fill** | `--pointer-fill: var(--blue-500)` | `#2563EB` | `rgb(37, 99, 235)` | `hsl(221, 83%, 53%)` | `oklch(0.546 0.215 262.9)` |
| **`glass` Button Surface** | `--gradient-card-dark` | `#FFFFFF12` -> `#FFFFFF05` | `rgba(255, 255, 255, 0.07)` -> `rgba(255, 255, 255, 0.02)` | `hsla(0, 0%, 100%, 0.07)` -> `hsla(0, 0%, 100%, 0.02)` | `oklch(1 0 0 / 7%)` -> `oklch(1 0 0 / 2%)` |
| **`ghost` Button Hover Fill** | `--muted` (`--paper-muted`) | `#F4F6FA` | `rgb(244, 246, 250)` | `hsl(220, 37%, 97%)` | `oklch(0.972 0.006 264.5)` |
| **`.shine-sweep` Specular Band** | `color-mix(in oklab, var(--brand-highlight) 45%, transparent)` | `#03D5E773` | `rgba(3, 213, 231, 0.45)` | `hsla(185, 97%, 46%, 0.45)` | `oklch(0.797 0.136 205.6 / 45%)` |
| **`.spotlight` Radial Core** | `color-mix(in oklab, var(--brand-secondary) 16%, transparent)` | `#2563EB29` | `rgba(37, 99, 235, 0.16)` | `hsla(221, 83%, 53%, 0.16)` | `oklch(0.546 0.215 262.9 / 16%)` |
| **`.spotlight-light` Radial Core** | `color-mix(in oklab, var(--brand-secondary) 9%, transparent)` | `#2563EB17` | `rgba(37, 99, 235, 0.09)` | `hsla(221, 83%, 53%, 0.09)` | `oklch(0.546 0.215 262.9 / 9%)` |
| **`ToggleRow` Active Light Surface** | `bg-primary/[0.06]` + `border-primary/40` | `#2563EB0F` / `#2563EB66` | `rgba(37, 99, 235, 0.06)` / `rgba(37, 99, 235, 0.40)` | `hsla(221, 83%, 53%, 0.06)` / `hsla(221, 83%, 53%, 0.40)` | `oklch(0.546 0.215 262.9 / 6%)` |

---

## 2. `WhiteBlueButton` Component System (6 Variants × 4 Sizes)

### 2.1 Architectural Rules

- **Geometry & Typography:** Uses `--radius-button` (`12px`), `font-display` (`Ubuntu`), `font-medium` (`500`), `inline-flex items-center justify-center gap-2`, and tactile `active:translate-y-[1px]` press feedback.
- **Child Layering (`[&>*]:relative [&>*]:z-10`):** Direct children and SVGs sit at `z-10` so pseudo-element backgrounds (`.shine-sweep::after` and `.pointer-fill::before` at `z-index: -1`) sweep cleanly behind the label without obscuring text.
- **Optional Magnetic Wrapper (`isMagnetic`):** Passing `isMagnetic` wraps the button in `<Magnetic strength={0.22}>` for subtle spring cursor attraction on desktop pointers.

| Variant | Visual Behavior & Hover Mechanics |
|:---|:---|
| `primary` (Default) | `.shine-sweep`, white text (`#FFFFFF`), `--gradient-accent` (`#2563EB` -> `#822EE8`), `shadow-[var(--shadow-card)]`. On hover: shifts to `--gradient-accent-hover`, elevates to `shadow-[var(--shadow-lift)]`, and runs `@keyframes wb-shine`. |
| `solid` | `.shine-sweep`, Deep Royal Navy fill (`bg-brand-primary` `#0D2975`), white text (`#FFFFFF`). On hover: transitions to `--gradient-accent-hover` with `shadow-[var(--shadow-lift)]`. |
| `outline` | `.pointer-fill`, `1px` border (`#E2E6ED`), transparent background, dark ink text (`#090E18`). On hover: `.pointer-fill::before` scales vertically from bottom (`scaleY(0 -> 1)`) in `--blue-500` (`#2563EB`) while text shifts to `#FFFFFF`. |
| `glass` | `.shine-sweep gradient-ring`, translucent dark surface (`--gradient-card-dark`), `backdrop-blur-md`, `text-on-dark` (`#D0D7E5`), hover brightens to `bg-white/10`. |
| `ghost` | Transparent background, `text-foreground` (`#090E18`), hover fills with `bg-muted` (`#F4F6FA`). |
| `link` | Zero horizontal padding (`px-0`), `text-brand-secondary` (`#2563EB`), `underline-offset-4 hover:underline`. |

| Size | Height & Padding | Font Size | Usage |
|:---|:---|:---|:---|
| `sm` | `h-9 px-4` (`36px` tall, `16px` horizontal padding) | `text-[13px]` | Sticky header CTAs, footer actions, compact card actions |
| `md` (Default) | `h-11 px-5` (`44px` tall, `20px` horizontal padding) | `text-sm` (`14px`) | Standard section CTAs, modal actions |
| `lg` | `h-13 px-7` (`52px` tall, `28px` horizontal padding) | `text-[15px]` | Hero primary/secondary CTAs, Pricing quote CTA, Lead Form submit |
| `icon` | `h-11 w-11 px-0` (`44px × 44px` square) | — | Standalone icon triggers |

### 2.2 Complete React/TSX Implementation (`WhiteBlueButton`)

```tsx
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import type { ButtonHTMLAttributes } from "react";
import { Magnetic } from "@/components/primitives/motion";
import { cn } from "@/lib/utils";

export const whiteBlueButtonVariants = cva(
  "relative inline-flex select-none items-center justify-center gap-2 whitespace-nowrap rounded-[var(--radius-button)] font-display font-medium transition-[transform,box-shadow,background-color,color] duration-[var(--dur-fast)] ease-[var(--ease-out)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:pointer-events-none disabled:opacity-50 active:translate-y-[1px] [&_svg]:relative [&_svg]:z-10 [&_svg]:size-4 [&_svg]:shrink-0 [&>*]:relative [&>*]:z-10",
  {
    variants: {
      variant: {
        primary:
          "shine-sweep overflow-hidden text-white shadow-[var(--shadow-card)] bg-[image:var(--gradient-accent)] hover:bg-[image:var(--gradient-accent-hover)] hover:shadow-[var(--shadow-lift)]",
        solid:
          "shine-sweep overflow-hidden bg-brand-primary text-white hover:bg-[image:var(--gradient-accent-hover)] hover:shadow-[var(--shadow-lift)]",
        outline:
          "pointer-fill overflow-hidden border border-border bg-transparent text-foreground hover:border-[var(--blue-500)] hover:text-white [--pointer-fill:var(--blue-500)]",
        glass:
          "shine-sweep gradient-ring overflow-hidden bg-[image:var(--gradient-card-dark)] text-on-dark backdrop-blur-md hover:bg-white/10",
        ghost: "text-foreground hover:bg-muted",
        link: "text-brand-secondary underline-offset-4 hover:underline px-0",
      },
      size: {
        sm: "h-9 px-4 text-[13px]",
        md: "h-11 px-5 text-sm",
        lg: "h-13 px-7 text-[15px]",
        icon: "h-11 w-11 px-0",
      },
    },
    defaultVariants: { variant: "primary", size: "md" },
  },
);

export type WhiteBlueButtonProps = ButtonHTMLAttributes<HTMLButtonElement> &
  VariantProps<typeof whiteBlueButtonVariants> & {
    asChild?: boolean;
    isMagnetic?: boolean;
  };

export function WhiteBlueButton({
  className,
  variant,
  size,
  asChild = false,
  isMagnetic = false,
  ...props
}: WhiteBlueButtonProps) {
  const Comp = asChild ? Slot : "button";
  const buttonNode = (
    <Comp className={cn(whiteBlueButtonVariants({ variant, size }), className)} {...props} />
  );

  if (isMagnetic) {
    return <Magnetic strength={0.22}>{buttonNode}</Magnetic>;
  }

  return buttonNode;
}
```

---

## 3. CSS3 Kinetic Utilities: `.shine-sweep`, `.pointer-fill`, `.gradient-text` & `.slide-swap`

```css
/* 1. 45-Degree Cyan Specular Button Sweep (runs strictly on hover/focus) */
@keyframes wb-shine {
  from {
    background-position: 200% 0;
  }
  to {
    background-position: -200% 0;
  }
}

.shine-sweep {
  isolation: isolate;
}

.shine-sweep::after {
  content: "";
  position: absolute;
  inset: 0;
  z-index: -1;
  border-radius: inherit;
  pointer-events: none;
  background-image: linear-gradient(
    45deg,
    transparent 25%,
    color-mix(in oklab, var(--brand-highlight) 45%, transparent) 50%,
    transparent 75%,
    transparent 100%
  );
  background-size: 250% 250%;
  background-position: 200% 0;
}

.shine-sweep:hover::after,
.shine-sweep:focus-visible::after,
.group:hover .shine-sweep::after {
  animation: wb-shine 3s linear infinite;
}

/* 2. Bottom-Origin ScaleY Fill for Outline Buttons */
.pointer-fill {
  isolation: isolate;
}

.pointer-fill::before {
  content: "";
  position: absolute;
  inset: 0;
  z-index: -1;
  border-radius: inherit;
  pointer-events: none;
  transform: scaleY(0);
  transform-origin: bottom;
  transition: transform var(--dur-fast) var(--ease-out);
  background: var(--pointer-fill, var(--gradient-brand));
}

.pointer-fill:hover::before {
  transform: scaleY(1);
  transform-origin: bottom;
}

/* 3. Gradient Accent Text with Optical Depth & Hover Sweep */
@utility gradient-text {
  position: relative;
  display: inline;
  background-image: var(--gradient-text);
  background-size: 200% 100%;
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  background-position: 0% 0;
  transition: background-position 0.4s ease;
  filter: drop-shadow(1px 0.7px 0 color-mix(in oklab, var(--ink) 22%, transparent));
}

@keyframes wb-text-sweep {
  from {
    background-position: 200% 0;
  }
  to {
    background-position: -200% 0;
  }
}

.gradient-text:hover,
h1:hover .gradient-text,
h2:hover .gradient-text,
h3:hover .gradient-text,
a:hover .gradient-text,
button:hover .gradient-text {
  animation: wb-text-sweep 3s linear infinite;
}

/* 4. Per-Character Vertical Label Swap */
.slide-swap {
  line-height: 1.15;
}

.slide-swap-top,
.slide-swap-bottom {
  transition: transform 520ms var(--ease-out);
}

.slide-swap-bottom {
  transform: translateY(100%);
}

.slide-swap:hover .slide-swap-top,
.group:hover .slide-swap-top,
a:hover > .slide-swap .slide-swap-top,
button:hover .slide-swap-top {
  transform: translateY(-110%);
}

.slide-swap:hover .slide-swap-bottom,
.group:hover .slide-swap-bottom,
a:hover > .slide-swap .slide-swap-bottom,
button:hover .slide-swap-bottom {
  transform: translateY(0);
}
```

---

## 4. Scroll Entrance Grammar (`Reveal`, `StaggerGroup`, `MaskedHeading`) & Motion Variants

### 4.1 Canonical Motion Variants (`src/lib/motion.ts`)

| Variant Name | Hidden State | Visible State | Duration & Easing | Purpose |
|:---|:---|:---|:---|:---|
| `fadeUp` | `{ opacity: 0, y: 28 }` | `{ opacity: 1, y: 0 }` | `700ms` `[0.16, 1, 0.3, 1]` | Standard desktop scroll reveal |
| `fadeUpCompact` | `{ opacity: 0, y: 16 }` | `{ opacity: 1, y: 0 }` | `420ms` `[0.16, 1, 0.3, 1]` | Mobile / touch viewport scroll reveal |
| `scaleIn` | `{ opacity: 0, scale: 0.94 }` | `{ opacity: 1, scale: 1 }` | `700ms` `[0.16, 1, 0.3, 1]` | Modal and media card entrance |
| `maskUp` | `{ y: "110%" }` | `{ y: "0%" }` | `1100ms` `[0.16, 1, 0.3, 1]` | Per-word / per-line `MaskedHeading` reveal |
| `viewportOnce` | — | `{ once: true, amount: 0.15, margin: "0px 0px -12% 0px" }` | — | Shared IntersectionObserver trigger |

### 4.2 Static-First Hydration & Flick-Scroll Safety Guard

1. **Static-First Above-the-Fold (`useDeferredReveal`):** On initial SSR and first client paint, blocks already inside the viewport (`rect.top < window.innerHeight * 0.95`) remain static plain HTML (`opacity: 1`). Only below-the-fold elements opt into scroll-triggered reveal animations, guaranteeing zero LCP delay.
2. **900ms Flick-Scroll Safety Guard (`useRevealGuard`):** If a fast flick-scroll skips past a `<Reveal>` element and its computed `opacity` is still `< 0.99` after `900ms`, `isStuck` flips to `true` and drops the motion wrapper so content never remains invisible.
3. **`MaskedHeading` Word/Line Splitter:** Splits headlines by spaces (`words`) or explicit `lines[]`, wrapping each unit in `<span className="inline-block overflow-hidden align-bottom pb-[0.08em]">` and applying `.gradient-text` to any word matching the `accent` prop.

---

## 5. Interactive Pointer Primitives (`Magnetic`, `TiltCard`, `SpotlightCard`, `ToggleRow`, `DragRail`)

### 5.1 `Magnetic`, `TiltCard` & `SpotlightCard`

```tsx
import { motion, useMotionValue, useReducedMotion, useSpring } from "motion/react";
import { useRef, type ElementType, type ReactNode } from "react";
import type React from "react";
import { cn } from "@/lib/utils";

/** Pulls element toward cursor within `radius` px and springs back on leave. */
export function Magnetic({
  children,
  className,
  strength = 0.28,
  radius = 90,
}: {
  children: ReactNode;
  className?: string;
  strength?: number;
  radius?: number;
}) {
  const isReducedMotion = useReducedMotion();
  const x = useSpring(0, { stiffness: 260, damping: 18, mass: 0.4 });
  const y = useSpring(0, { stiffness: 260, damping: 18, mass: 0.4 });

  if (isReducedMotion) {
    return <span className={cn("inline-flex", className)}>{children}</span>;
  }

  return (
    <motion.span
      className={cn("inline-flex", className)}
      style={{ x, y }}
      onPointerMove={(event) => {
        if (event.pointerType !== "mouse") {
          return;
        }

        const rect = event.currentTarget.getBoundingClientRect();
        const dx = event.clientX - (rect.left + rect.width / 2);
        const dy = event.clientY - (rect.top + rect.height / 2);
        x.set(Math.max(-radius, Math.min(radius, dx)) * strength);
        y.set(Math.max(-radius, Math.min(radius, dy)) * strength);
      }}
      onPointerLeave={() => {
        x.set(0);
        y.set(0);
      }}
    >
      {children}
    </motion.span>
  );
}

/** rAF-throttled cursor-tracking radial spotlight card. */
export function SpotlightCard({
  children,
  className,
  as: Tag = "div",
  ...rest
}: {
  children: ReactNode;
  className?: string;
  as?: ElementType;
} & Record<string, unknown>) {
  const isReducedMotion = useReducedMotion();
  const frameRef = useRef(0);
  const Comp = Tag as ElementType;

  const handlePointerMove = (event: React.PointerEvent<HTMLDivElement>) => {
    if (event.pointerType !== "mouse") {
      return;
    }

    const element = event.currentTarget;
    const clientX = event.clientX;
    const clientY = event.clientY;

    if (frameRef.current !== 0) {
      return;
    }

    frameRef.current = requestAnimationFrame(() => {
      frameRef.current = 0;
      const rect = element.getBoundingClientRect();
      element.style.setProperty("--mx", `${clientX - rect.left}px`);
      element.style.setProperty("--my", `${clientY - rect.top}px`);
    });
  };

  return (
    <Comp
      className={cn("spotlight relative", className)}
      onPointerMove={isReducedMotion ? undefined : handlePointerMove}
      {...rest}
    >
      {children}
    </Comp>
  );
}
```

### 5.2 `ToggleRow` Spring Switch (`cubic-bezier(0.34, 1.4, 0.5, 1)`)

Used inside the Hero `CapabilityStack` and interactive `Pricing` quote configurator:

```tsx
import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

export function ToggleRow({
  title,
  sub,
  icon,
  isEnabled,
  onToggle,
  className,
  tone = "light",
}: {
  title: string;
  sub: string;
  icon?: ReactNode;
  isEnabled: boolean;
  onToggle: () => void;
  className?: string;
  tone?: "dark" | "light";
}) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={isEnabled}
      onClick={onToggle}
      className={cn(
        "group flex w-full items-center gap-4 rounded-[var(--radius-card)] border px-4 py-3.5 text-left transition-[background-color,border-color,transform] duration-[var(--dur-base)] ease-[var(--ease-out)]",
        tone === "dark"
          ? isEnabled
            ? "border-primary/40 bg-white/[0.06] text-on-dark"
            : "border-hairline bg-white/[0.02] text-on-dark/80 hover:bg-white/[0.04]"
          : isEnabled
            ? "border-primary/40 bg-primary/[0.06] text-foreground"
            : "border-border bg-card text-foreground hover:bg-muted",
        className,
      )}
    >
      {icon ? (
        <span
          className={cn(
            "grid size-10 shrink-0 place-items-center rounded-xl border transition-colors",
             isEnabled ? "border-primary/30 bg-primary/10" : "border-border bg-muted",
          )}
        >
          {icon}
        </span>
      ) : null}

      <span className="min-w-0 flex-1">
        <span className="block font-display text-sm font-semibold leading-tight">{title}</span>
        <span className="mt-1 block text-xs leading-snug text-muted-foreground">{sub}</span>
      </span>

      <span
        aria-hidden
        className={cn(
          "flex h-8 w-14 shrink-0 items-center rounded-full p-1 transition-colors duration-[var(--dur-base)]",
          isEnabled ? "bg-primary" : "bg-border",
        )}
      >
        <span
          className={cn(
            "size-6 rounded-full bg-white shadow transition-transform duration-[520ms] ease-[cubic-bezier(0.34,1.4,0.5,1)]",
            isEnabled ? "translate-x-6" : "translate-x-0",
          )}
        />
      </span>
    </button>
  );
}
```

---

## 6. Parametric LESS Mixins for Buttons & Micro-Interactions

```less
// ============================================================================
// WHITE BLUE THEME — BUTTONS & INTERACTION LESS MIXINS
// ============================================================================

.wb-button-base() {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  border-radius: 12px;
  font-family: "Ubuntu", system-ui, sans-serif;
  font-weight: 500;
  isolation: isolate;
  overflow: hidden;
  transition: transform 240ms cubic-bezier(0.16, 1, 0.3, 1),
              box-shadow 240ms cubic-bezier(0.16, 1, 0.3, 1),
              background-color 240ms cubic-bezier(0.16, 1, 0.3, 1),
              color 240ms cubic-bezier(0.16, 1, 0.3, 1);

  &:active {
    transform: translateY(1px);
  }
}

.wb-button-primary() {
  .wb-button-base();
  color: #FFFFFF; // rgb(255, 255, 255) | hsl(0, 0%, 100%)
  background-image: linear-gradient(90deg, #2563EB 0%, #822EE8 100%);
  box-shadow: 0 8px 24px -12px rgba(13, 41, 117, 0.18);

  &::after {
    content: "";
    position: absolute;
    inset: 0;
    z-index: -1;
    background-image: linear-gradient(
      45deg,
      transparent 25%,
      rgba(3, 213, 231, 0.45) 50%,
      transparent 75%
    );
    background-size: 250% 250%;
    background-position: 200% 0;
  }

  &:hover {
    background-image: linear-gradient(90deg, #4277EE 0%, #924AEB 100%);
    box-shadow: 0 24px 60px -24px rgba(13, 41, 117, 0.35);

    &::after {
      animation: wb-shine 3s linear infinite;
    }
  }
}
```
