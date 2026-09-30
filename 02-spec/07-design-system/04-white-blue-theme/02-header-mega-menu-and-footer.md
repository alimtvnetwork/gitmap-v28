# White Blue Theme: Sticky Header, 3D Flip Mega-Menu, Footer & Floating Chrome

> **/goal** Specify the complete navigation, mega-menu, mobile drawer, soft-band footer, and floating viewport chrome (`StickyCta`, `BackToTop`) for the White Blue Theme (`04-white-blue-theme`).
> **/learn** Master the `72px` glassmorphic sticky header, scroll-threshold shadow elevation (`window.scrollY > 12`), per-character `SlideSwapLabel` hover links, `14px` safe-region pointer tracking, the **3D Flip Promo Card** (`[perspective:1400px]`, `rotateY(180deg)`), and non-destructive mobile scroll locking.

---

## 1. Navigation & Chrome Color Tokens (Mandatory 3-Format Reference)

Every surface, border, and interactive highlight in the header, mega-menu, footer, and floating bars resolves to semantic tokens documented below in **HEX**, **RGB / RGBA**, **HSL / HSLA**, and **OKLCH**:

| Element / State | CSS Token / Expression | Format 1: HEX | Format 2: RGB / RGBA | Format 3: HSL / HSLA | OKLCH Runtime Value |
|:---|:---|:---|:---|:---|:---|
| **Header Translucent Ground** | `bg-background/90` (`--paper` @ 90%) | `#FFFFFFE6` | `rgba(255, 255, 255, 0.90)` | `hsla(0, 0%, 100%, 0.90)` | `oklch(1 0 0 / 90%)` |
| **Header Resting Hairline** | `border-border/40` (`--neutral-200` @ 40%) | `#E2E6ED66` | `rgba(226, 230, 237, 0.40)` | `hsla(218, 23%, 91%, 0.40)` | `oklch(0.924 0.011 265 / 40%)` |
| **Header Scrolled Border** | `border-border` (`--neutral-200`) | `#E2E6ED` | `rgb(226, 230, 237)` | `hsl(218, 23%, 91%)` | `oklch(0.924 0.011 265)` |
| **Header Scrolled Shadow** | `--shadow-card` (`--brand-primary` @ 18%) | `#0D29752E` | `rgba(13, 41, 117, 0.18)` | `hsla(224, 80%, 25%, 0.18)` | `oklch(0.317 0.135 264.5 / 18%)` |
| **Nav Link Resting Text** | `text-foreground/75` (`--ink` @ 75%) | `#090E18BF` | `rgba(9, 14, 24, 0.75)` | `hsla(220, 45%, 6%, 0.75)` | `oklch(0.163 0.023 265 / 75%)` |
| **Nav Link Active / Hover Text** | `text-foreground` (`--ink`) | `#090E18` | `rgb(9, 14, 24)` | `hsl(220, 45%, 6%)` | `oklch(0.163 0.023 265)` |
| **Nav Bottom Hairline Rule** | `bg-brand-secondary` (`--blue-500`) | `#2563EB` | `rgb(37, 99, 235)` | `hsl(221, 83%, 53%)` | `oklch(0.546 0.215 262.9)` |
| **Mega-Menu Item Hover Tint** | `color-mix(in oklab, var(--primary) 7%, transparent)` | `#2563EB12` | `rgba(37, 99, 235, 0.07)` | `hsla(221, 83%, 53%, 0.07)` | `oklch(0.546 0.215 262.9 / 7%)` |
| **3D Promo Front Gradient** | `--gradient-accent` (`--blue-500` -> `--violet-500`) | `#2563EB` -> `#822EE8` | `rgb(37, 99, 235)` -> `rgb(130, 46, 232)` | `hsl(221, 83%, 53%)` -> `hsl(267, 80%, 55%)` | `oklch(0.546 0.215 262.9)` -> `oklch(0.532 0.253 296.9)` |
| **Footer Canvas (`.band-soft`)** | `bg-surface-soft` (`--surface-soft`) | `#EEF3FA` | `rgb(238, 243, 250)` | `hsl(215, 55%, 96%)` | `oklch(0.962 0.011 258)` |
| **Messaging Accent Badge** | Brand messaging green | `#25D366` | `rgb(37, 211, 102)` | `hsl(142, 70%, 49%)` | `oklch(0.765 0.198 150.5)` |

---

## 2. Sticky `72px` Glassmorphic Header & Scroll Threshold

The primary site header (`SiteHeader`) is a sticky `72px` (`h-[72px]`) bar anchored at `top-0 z-50`:

1. **Backdrop Materiality:** Uses `bg-background/90 backdrop-blur-xl` (`rgba(255, 255, 255, 0.90)` over `backdrop-filter: blur(24px)`).
2. **Scroll Threshold Elevation (`window.scrollY > 12`):**
   - When `window.scrollY <= 12`, the bottom border is a soft 40% hairline (`border-b border-border/40`, `rgba(226, 230, 237, 0.40)`) with zero shadow so the header blends seamlessly into the hero's `--blue-50` (`#F0F6FF`) top wash.
   - Once `window.scrollY > 12` (`isScrolled`), the header transitions over `420ms cubic-bezier(0.16, 1, 0.3, 1)` to a solid `1px` border (`#E2E6ED`) and casts `var(--shadow-card)` (`0 8px 24px -12px rgba(13, 41, 117, 0.18)`).
3. **Three-Zone Horizontal Architecture:**
   - **Left Cluster:** Brand logo mark + `1px` vertical divider (`h-7 w-px bg-border`, visible at `xl`) + 2-line monospace partner credential badge (`font-mono text-[10px] uppercase leading-[1.35] tracking-[0.16em] text-muted-foreground`).
   - **Center Primary Nav (`lg:flex`):** Top-level links rendered with `SlideSwapLabel` (`stagger={0.04}`), a `1px` bottom underline (`-bottom-1.5 h-px w-full origin-left scale-x-0 bg-brand-secondary transition-transform duration-[var(--dur-fast)] group-hover:scale-x-100`), and a `14px` (`size-3.5`) `ChevronDown` that rotates `180deg` when its Mega-Menu panel is active.
   - **Right CTA Pair (`lg:flex`):** Secondary `outline` `WhiteBlueButton` (`size="sm"`) paired with primary gradient `WhiteBlueButton` (`size="sm"`). Hovering the CTA cluster immediately closes any open Mega-Menu panel (`openPanelNow(null)`).

---

## 3. Safe-Region Hover Mega-Menu & 3D Flip Promo Card

### 3.1 Safe-Region Pointer Geometry & Close Delay

Standard hover dropdowns flicker when the user moves the pointer diagonally from a nav trigger toward the panel. The White Blue Theme solves this with two coordinated mechanisms:

1. **Close Delay Timer (`110ms`–`220ms`):** Leaving a trigger schedules a delayed close (`scheduleClose`) cancelled immediately if the pointer enters the Mega-Menu container (`cancelClose`).
2. **`14px` Padded Bounding-Rect Safe Region:** While a panel is open, a passive `pointermove` listener computes the union of `headerRef.getBoundingClientRect()` and `panelRef.getBoundingClientRect()`, expanded by `pad = 14` pixels on all four sides. As long as the cursor stays inside this padded corridor (`isInsideSafeRegion`), the menu remains locked open. Pressing `Escape` or scrolling the window closes the panel immediately.

### 3.2 Mega-Menu Grid & Capability Link Anatomy

- **Outer Container:** Positioned at `absolute left-0 right-0 top-full z-40 pt-3`, wrapping a `<Container>` with `rounded-[var(--radius-card)]` (`20px`), `border border-border` (`#E2E6ED`), `bg-card` (`#FFFFFF`), and `shadow-[var(--shadow-lift)]` (`0 24px 60px -24px rgba(13, 41, 117, 0.35)`).
- **Responsive Column Template:**
  - 3+ link groups: `lg:grid-cols-[1fr_1fr_1fr_0.9fr]`
  - 2 link groups: `lg:grid-cols-[1fr_1fr_1.1fr]`
  - 1 link group: `lg:grid-cols-[1.6fr_1fr]`
- **Capability Link Item Hover:** Each link (`px-3 py-2 rounded-[10px]`) warms on hover with `hover:bg-[color-mix(in_oklab,var(--primary)_7%,transparent)]` (`rgba(37, 99, 235, 0.07)`), grows a `1px` vertical accent rule along its left edge (`origin-top scale-y-0 bg-[image:var(--gradient-accent)] group-hover:scale-y-100`), animates its title via `SlideSwapLabel`, and slides a `14px` `ArrowRight` icon from `-translate-x-1 opacity-0` to `translate-x-0 opacity-100`.

### 3.3 Right-Hand 3D Flip Promo Card (`[perspective:1400px]`)

The rightmost column of every Mega-Menu houses a two-sided 3D card (`min-h-[220px] [perspective:1400px]`) that flips `180deg` around the Y-axis over `820ms cubic-bezier(0.16, 1, 0.3, 1)` when hovered:

- **Front Face (`[backface-visibility:hidden]`):** Full `--gradient-accent` background (`#2563EB` -> `#822EE8`), bold white headline (`#FFFFFF`), 85% white supporting copy (`rgba(255, 255, 255, 0.85)`), and a bottom `"HOVER TO FLIP ->"` monospace prompt.
- **Back Face (`[backface-visibility:hidden] [transform:rotateY(180deg)]`):** Crisp white card surface (`#FFFFFF` / `rgb(255, 255, 255)` / `hsl(0, 0%, 100%)`), `1px` border (`#E2E6ED`), dark ink headline (`#090E18`), muted body copy (`#48536B`), and a pill CTA button (`rounded-full bg-[image:var(--gradient-accent)] px-4 py-2.5 text-sm font-semibold text-white shadow-[var(--shadow-lift)]`).

### 3.4 Complete React/TSX Implementation (`SiteHeader` & `MegaPanel`)

```tsx
import { useCallback, useEffect, useRef, useState, type RefObject } from "react";
import { AnimatePresence, motion, useReducedMotion } from "motion/react";
import { ArrowRight, ChevronDown, Menu, X } from "lucide-react";
import { Container } from "@/components/primitives/layout";
import { WhiteBlueButton } from "@/components/primitives/white-blue-button";
import { SlideSwapLabel } from "@/components/primitives/motion";
import { cn } from "@/lib/utils";

export interface NavLinkItem {
  label: string;
  href: string;
  description?: string;
}

export interface MegaPromoCard {
  front: { title: string; body: string };
  back: { title: string; body: string; cta: { label: string; href: string } };
}

export function MegaPanel({
  panelKey,
  groups,
  promo,
  panelRef,
  onNavigate,
  onMouseEnter,
  onMouseLeave,
}: {
  panelKey: string;
  groups: { title: string; links: NavLinkItem[] }[];
  promo: MegaPromoCard;
  panelRef: RefObject<HTMLDivElement | null>;
  onNavigate: () => void;
  onMouseEnter: () => void;
  onMouseLeave: () => void;
}) {
  const isReducedMotion = useReducedMotion();

  const initialAnim = isReducedMotion
    ? { opacity: 0 }
    : { opacity: 0, y: -8, scale: 0.985 };

  const exitAnim = isReducedMotion
    ? { opacity: 0 }
    : { opacity: 0, y: -6, scale: 0.99 };

  return (
    <motion.div
      ref={panelRef}
      key={panelKey}
      initial={initialAnim}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      exit={exitAnim}
      transition={{ duration: isReducedMotion ? 0.12 : 0.26, ease: [0.16, 1, 0.3, 1] }}
      className="absolute left-0 right-0 top-full z-40 pt-3"
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
    >
      <Container>
        <div className="overflow-hidden rounded-[var(--radius-card)] border border-border bg-card shadow-[var(--shadow-lift)]">
          <div className="grid gap-8 p-8 lg:grid-cols-[1fr_1fr_1.1fr]">
            {groups.map((group) => (
              <div key={group.title} className="flex flex-col gap-3">
                <p className="text-eyebrow text-muted-foreground">{group.title}</p>
                <ul className="flex flex-col gap-1">
                  {group.links.map((link) => (
                    <li key={link.href}>
                      <a
                        href={link.href}
                        onClick={onNavigate}
                        className="group relative block overflow-hidden rounded-[var(--radius-sm,10px)] px-3 py-2 transition-colors hover:bg-[color-mix(in_oklab,var(--primary)_7%,transparent)]"
                      >
                        <span
                          aria-hidden
                          className="absolute inset-y-1 left-0 w-px origin-top scale-y-0 bg-[image:var(--gradient-accent)] transition-transform duration-[var(--dur-base)] ease-[var(--ease-out)] group-hover:scale-y-100"
                        />
                        <span className="flex items-center gap-1.5 font-display text-sm font-medium text-foreground">
                          <SlideSwapLabel>{link.label}</SlideSwapLabel>
                          <ArrowRight className="size-3.5 -translate-x-1 opacity-0 transition-all group-hover:translate-x-0 group-hover:opacity-100" />
                        </span>
                        {link.description ? (
                          <span className="mt-0.5 block text-xs leading-relaxed text-muted-foreground">
                            {link.description}
                          </span>
                        ) : null}
                      </a>
                    </li>
                  ))}
                </ul>
              </div>
            ))}

            {/* Right-Hand 3D Flip Promo Card */}
            <div className="group/promo relative min-h-[220px] [perspective:1400px]">
              <div className="relative h-full w-full transition-transform duration-[820ms] ease-[var(--ease-out)] [transform-style:preserve-3d] group-hover/promo:[transform:rotateY(180deg)]">
                {/* Front Face */}
                <div className="absolute inset-0 flex flex-col justify-between gap-4 overflow-hidden rounded-[var(--radius-card)] bg-[image:var(--gradient-accent)] p-6 text-white [backface-visibility:hidden]">
                  <div>
                    <p className="text-base font-bold leading-snug">{promo.front.title}</p>
                    <p className="mt-2 text-sm text-white/85">{promo.front.body}</p>
                  </div>
                  <span className="inline-flex items-center gap-1.5 text-xs font-semibold uppercase tracking-[0.14em] text-white/80">
                    Hover to flip
                    <ArrowRight className="size-3.5" />
                  </span>
                </div>

                {/* Back Face */}
                <div className="absolute inset-0 flex flex-col justify-between gap-4 overflow-hidden rounded-[var(--radius-card)] border border-border bg-card p-6 text-foreground [backface-visibility:hidden] [transform:rotateY(180deg)]">
                  <div>
                    <p className="text-base font-bold leading-snug">{promo.back.title}</p>
                    <p className="mt-2 text-sm text-muted-foreground">{promo.back.body}</p>
                  </div>
                  <a
                    href={promo.back.cta.href}
                    onClick={onNavigate}
                    className="group/cta inline-flex items-center justify-center gap-2 rounded-full bg-[image:var(--gradient-accent)] px-4 py-2.5 text-sm font-semibold text-white shadow-[var(--shadow-lift)] transition-transform duration-[var(--dur-fast)] hover:scale-[1.02]"
                  >
                    {promo.back.cta.label}
                    <ArrowRight className="size-4 transition-transform duration-[var(--dur-base)] group-hover/cta:translate-x-1" />
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Container>
    </motion.div>
  );
}
```

---

## 4. Non-Destructive Mobile Drawer Lock Rule

> [!CAUTION]
> **NEVER set `overflow: hidden` on `<body>` when opening the mobile drawer.**
> Mutating `body { overflow: hidden }` changes the document scrolling container and immediately breaks `position: sticky` on the `72px` header, causing the header to jump off-screen whenever the user opens the menu partway down the page.

Instead, intercept `wheel` and `touchmove` events at the `window` level with `{ passive: false }` and call `event.preventDefault()` only when the event target lies outside `mobileMenuRef.current`:

```tsx
useEffect(() => {
  if (isMobileMenuClosed) {
    return;
  }

  const isTargetOutsideMenu = (target: EventTarget | null) => {
    if (target instanceof Node) {
      return mobileMenuRef.current?.contains(target) === false;
    }

    return true;
  };

  const blockBackgroundScroll = (event: Event) => {
    if (isTargetOutsideMenu(event.target)) {
      event.preventDefault();
    }
  };

  window.addEventListener("wheel", blockBackgroundScroll, { passive: false });
  window.addEventListener("touchmove", blockBackgroundScroll, { passive: false });

  return () => {
    window.removeEventListener("wheel", blockBackgroundScroll);
    window.removeEventListener("touchmove", blockBackgroundScroll);
  };
}, [isMobileMenuClosed]);
```

---

## 5. Soft-Band Site Footer (`SiteFooter`)

The footer is rendered on the cool-mist `--surface-soft` canvas (`#EEF3FA` / `rgb(238, 243, 250)` / `hsl(215, 55%, 96%)` / `oklch(0.962 0.011 258)`) with four structural layers:

1. **Top Ambient Aura:** A blurred radial pill (`-top-52 left-1/2 h-[420px] w-[720px] -translate-x-1/2 rounded-full opacity-25 blur-3xl`) filled with `var(--gradient-accent)` (`#2563EB` -> `#822EE8`).
2. **Asymmetric Brand + 5-Column Sitemap Grid (`lg:grid-cols-[minmax(0,0.85fr)_minmax(0,3.15fr)]`):**
   - Left column: Brand mark, positioning statement (`max-w-xs text-sm text-muted-foreground`), and primary/outline CTA button pair.
   - Right 5-column grid: Each column heading uses `font-mono text-[10px] uppercase tracking-[0.2em] text-muted-foreground` and shifts to `text-primary` (`#2563EB`) when any link inside the column is hovered or focused (`group-hover:text-primary group-focus-within:text-primary`).
3. **Full-Width 4-Cell Contact Hairline Grid:**
   - Uses `grid gap-px overflow-hidden rounded-[var(--radius-card)] border border-border bg-border sm:grid-cols-2 lg:grid-cols-4` where each cell (`bg-card p-6 hover:bg-muted`) displays `Office`, `Email`, `Phone` (with inline `#25D366` / `rgb(37, 211, 102)` / `hsl(142, 70%, 49%)` messaging badge), and `Hours`.
4. **Legal & Social Bottom Bar:** Bordered by `border-t border-border py-6 text-xs text-muted-foreground`, pairing copyright links on the left with policy links and `36px` (`size-9`) circular social icon buttons on the right.

---

## 6. Floating Bottom-Center CTA Pill (`StickyCta`) & `BackToTop`

1. **Bottom-Center `StickyCta` (`fixed inset-x-0 bottom-6 z-40`):**
   - **Visibility Window:** Appears once the user scrolls past 90% of the initial viewport (`window.scrollY > window.innerHeight * 0.9`) and hides automatically within `900px` of the page bottom (`y + window.innerHeight <= doc.scrollHeight - 900`) so it never collides with the footer CTA.
   - **Composition:** Pairs a `48px` (`size-12`) circular messaging button (`#25D366` / `rgb(37, 211, 102)` / `hsl(142, 70%, 49%)`) with a gradient pill CTA (`rounded-full bg-[image:var(--gradient-accent)] py-2 pl-2 pr-5 text-white shadow-[var(--shadow-lift)]`) containing a `36px` (`size-9`) `bg-white/15` icon badge (`CalendarCheck`) and sliding `ArrowRight`.
2. **Floating `BackToTop` Button (`fixed bottom-24 right-[var(--gutter)] lg:bottom-6 z-40`):**
   - Appears past 60% total scroll depth (`window.scrollY / maxScroll > 0.6`).
   - Uses `size-11 rounded-full border border-hairline bg-card/95 text-foreground shadow-[var(--shadow-lift)] backdrop-blur-xl`, sitting at `bottom-24` on mobile (`<1024px`) so it never overlaps `StickyCta`.

---

## 7. Parametric LESS Mixins for Header & Mega-Menu

```less
// ============================================================================
// WHITE BLUE THEME — HEADER & 3D FLIP MEGA-MENU LESS MIXINS
// ============================================================================

.wb-sticky-header() {
  position: sticky;
  top: 0;
  z-index: 50;
  width: 100%;
  height: 72px;
  background-color: rgba(255, 255, 255, 0.90); // #FFFFFFE6 | hsla(0, 0%, 100%, 0.90)
  color: #090E18;                              // rgb(9, 14, 24) | hsl(220, 45%, 6%)
  backdrop-filter: blur(24px);
  border-bottom: 1px solid rgba(226, 230, 237, 0.40);
  transition: background-color 420ms cubic-bezier(0.16, 1, 0.3, 1),
              box-shadow 420ms cubic-bezier(0.16, 1, 0.3, 1),
              border-color 420ms cubic-bezier(0.16, 1, 0.3, 1);

  &.is-scrolled {
    border-bottom-color: #E2E6ED;              // rgb(226, 230, 237) | hsl(218, 23%, 91%)
    box-shadow: 0 8px 24px -12px rgba(13, 41, 117, 0.18);
  }
}

.wb-promo-flip-card() {
  position: relative;
  min-height: 220px;
  perspective: 1400px;

  &__inner {
    position: relative;
    width: 100%;
    height: 100%;
    transform-style: preserve-3d;
    transition: transform 820ms cubic-bezier(0.16, 1, 0.3, 1);
  }

  &:hover &__inner {
    transform: rotateY(180deg);
  }

  &__face {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding: 1.5rem;
    border-radius: 20px;
    backface-visibility: hidden;
  }

  &__face--front {
    background-image: linear-gradient(90deg, #2563EB 0%, #822EE8 100%);
    color: #FFFFFF;
  }

  &__face--back {
    background-color: #FFFFFF;
    border: 1px solid #E2E6ED;
    color: #090E18;
    transform: rotateY(180deg);
  }
}
```
