# Header & Precision Navigation System Specification

> **/goal** Specify the enterprise sticky glass header, shrinking pill navbar, safe-region pointer physics, character-staggered `SlideSwapLabel` navigation links, and full responsive behavior.
> **/learn** Master the 72px sticky header height, 12px scroll threshold, pointer safe-region collision math (`pad = 14px`, 220ms close debounce), character-staggered vertical slide swap (`stagger: 0.04`), and background scroll lock that preserves `position: sticky` without body overflow tampering.

**Version:** 4.1.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Executive System Overview

The **Header & Precision Navigation System** anchors the site's information architecture. It balances effortless discoverability with editorial visual luxury:
1. **72px Glass Header Bar:** Fixed sticky bar with `backdrop-blur-xl` and 90% opacity surface fill.
2. **Dynamic 12px Scroll Threshold:** Elevates from a subtle hairline border into a distinct card shadow once scrolled past 12px.
3. **Alternative Shrinking Pill Header:** Transitions from a 1240px transparent container into a 900px floating white pill after scroll.
4. **`SlideSwapLabel` Nav Links:** Per-character staggered vertical slide swap on hover (`stagger: 0.04`).
5. **Pointer Safe-Region Geometry:** Prevents accidental dropdown dismissals using a 14px quad collision buffer and 220ms debounce.
6. **Sticky-Preserving Mobile Drawer:** Prevents background scroll leaks via window event suppression **without** modifying `document.body.style.overflow` (which breaks CSS sticky positioning).

---

## 2. Header Geometry & Elevation Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ [Logo + Partner Badge]         [Primary Nav Links]        [CTA Button Pair] │
└─────────────────────────────────────────────────────────────────────────────┘
  ◄── 72px Height, Sticky Top-0, Z-Index 50, Backdrop-Blur-XL (12px Blur) ────►
```

### 2.1 Geometry Metrics

| Dimension / Property | Standard Desktop (≥1024px) | Mobile / Tablet (<1024px) |
|:---|:---|:---|
| **Height** | `72px` (`h-[72px]`) | `72px` (`h-[72px]`) |
| **Positioning** | `sticky top-0 z-50 w-full` | `sticky top-0 z-50 w-full` |
| **Max Content Width** | `1280px` (`max-w-[1280px] mx-auto`) | `100%` with gutters (`px-6`) |
| **Surface Background** | `bg-background/90 text-foreground` | `bg-background/90 text-foreground` |
| **Backdrop Filter** | `backdrop-blur-xl` | `backdrop-blur-xl` |
| **Transition Properties**| `background-color, box-shadow, backdrop-filter` | `background-color, box-shadow, backdrop-filter` |
| **Transition Timing** | `duration-[var(--dur-base,420ms)] ease-[var(--ease-out)]` | `duration-[var(--dur-base,420ms)] ease-[var(--ease-out)]` |

### 2.2 Elevation & Scroll State Transition

The header listens to viewport scroll via passive event listener:
```typescript
const [scrolled, setScrolled] = useState(false);

useEffect(() => {
  const onScroll = () => setScrolled(window.scrollY > 12);
  onScroll();
  window.addEventListener("scroll", onScroll, { passive: true });
  return () => window.removeEventListener("scroll", onScroll);
}, []);
```

- **Resting State (`scrollY <= 12px`):**
  - Class: `border-b border-border/40 shadow-none`
- **Scrolled State (`scrollY > 12px`):**
  - Class: `border-b border-border shadow-[var(--shadow-card)]`
  - Easing: `420ms cubic-bezier(0.16, 1, 0.3, 1)`

---

## 3. Nav Link Micro-Interactions

Primary navigation links incorporate three concurrent micro-interactions on hover:

```
          ┌─ SlideSwapLabel (Per-character upward roll, 520ms)
          ▼
   [ Solutions ▾ ]
   └───────────────┘
     ▲          ▲
     │          └─ ChevronDown (Rotates 180deg over 240ms)
     └─ Underline Grow (Scale-X: 0 -> 1 from left, 240ms)
```

### 3.1 Character-Staggered `SlideSwapLabel`
Nav links duplicate their text into an accessible hidden element and split visible characters into individual inline blocks with progressive transition delays:

- **Nav Link Stagger Rate:** `0.04s` per character.
- **Top Glyph (`.slide-swap-top`):** Starts at `translateY(0)`, transitions to `translateY(-110%)` over `520ms var(--ease-out)`.
- **Bottom Glyph (`.slide-swap-bottom`):** Starts at `translateY(100%)`, transitions to `translateY(0)` over `520ms var(--ease-out)`.
- **Accessibility:** Screen readers parse `<span class="sr-only">{children}</span>`, while visual glyphs are marked `aria-hidden`.
- **Reduced Motion:** When `prefers-reduced-motion` is active, the component renders a static `<span>` with zero transform or delay.

### 3.2 Underline Indicator Rule
- Height: `1px` (`h-px`).
- Position: `absolute -bottom-1.5 left-0 w-full`.
- Transform: `origin-left scale-x-0 group-hover:scale-x-100`.
- Color: `bg-brand-secondary` (or theme primary).
- Transition: `duration-[var(--dur-fast,240ms)] ease-[var(--ease-out)]`.

### 3.3 Rotating Chevron Indicator
- Icon: Lucide `ChevronDown` (`size-3.5` / `14px`).
- Rotation: `rotate-0` resting, `rotate-180` when panel is open.
- Transition: `duration-[var(--dur-fast,240ms)] ease-[var(--ease-out)]`.

---

## 4. Safe Pointer Region & Debounced Close Geometry

To prevent frustrating accidental menu collapses when users move diagonally from the header item into the dropdown panel, the system implements a geometric safe-buffer check:

```typescript
useEffect(() => {
  if (!openPanel) return;

  const onKey = (e: KeyboardEvent) => {
    if (e.key === "Escape") setOpenPanel(null);
  };

  const onPointerMove = (e: PointerEvent) => {
    const rects = [headerRef.current, panelRef.current]
      .map((element) => element?.getBoundingClientRect())
      .filter((rect): rect is DOMRect => Boolean(rect));

    const pad = 14; // 14px safe region padding around containers
    const insideSafeRegion = rects.some(
      (rect) =>
        e.clientX >= rect.left - pad &&
        e.clientX <= rect.right + pad &&
        e.clientY >= rect.top - pad &&
        e.clientY <= rect.bottom + pad,
    );

    if (insideSafeRegion) {
      cancelClose();
    } else {
      scheduleClose(); // 220ms debounced dismissal
    }
  };

  const onScroll = () => setOpenPanel(null);

  window.addEventListener("keydown", onKey);
  window.addEventListener("pointermove", onPointerMove, { passive: true });
  window.addEventListener("scroll", onScroll, { passive: true });

  return () => {
    window.removeEventListener("keydown", onKey);
    window.removeEventListener("pointermove", onPointerMove);
    window.removeEventListener("scroll", onScroll);
  };
}, [cancelClose, openPanel, scheduleClose]);
```

- **Safe Region Buffer:** `pad = 14px` around both the header rect and the mega-panel rect.
- **Close Debounce Delay:** `220ms` timeout before state resets to `null`.
- **Keyboard Dismissal:** `Escape` key immediately closes panel (`0ms`).
- **Scroll Dismissal:** Any page scroll immediately dismisses panel (`0ms`).

---

## 5. Sticky-Safe Mobile Drawer Scroll Locking

Traditional mobile drawers set `document.body.style.overflow = "hidden"` to prevent background scrolling. **This breaks `position: sticky` on the header in modern browsers**, shifting the header off-screen.

### 5.1 The Sticky-Safe Interception Algorithm
The mobile drawer intercepts `wheel` and `touchmove` events on the window, cancelling them **only** if the target is outside the mobile menu drawer:

```typescript
useEffect(() => {
  if (!mobileOpen) return;

  const isInsideMenu = (target: EventTarget | null) =>
    target instanceof Node && Boolean(mobileMenuRef.current?.contains(target));

  const block = (e: Event) => {
    if (!isInsideMenu(e.target)) e.preventDefault();
  };

  const onKey = (e: KeyboardEvent) => {
    if (e.key === "Escape") setMobileOpen(false);
  };

  window.addEventListener("wheel", block, { passive: false });
  window.addEventListener("touchmove", block, { passive: false });
  window.addEventListener("keydown", onKey);

  return () => {
    window.removeEventListener("wheel", block);
    window.removeEventListener("touchmove", block);
    window.removeEventListener("keydown", onKey);
  };
}, [mobileOpen]);
```

---

## 6. Anti-Hallucination & Quality Verification Checklist

- [ ] Header height is strictly `72px` with `top-0 z-50` sticky positioning.
- [ ] Scroll threshold is `12px` and transitions with `420ms` easing.
- [ ] Nav links utilize `SlideSwapLabel` with `stagger: 0.04s` and `520ms` roll transition.
- [ ] Underline expands from left over `240ms`.
- [ ] Chevron indicator is `14px` (`size-3.5`) and rotates `180deg` on open.
- [ ] Safe-region padding is `14px` with a `220ms` debounce timer.
- [ ] Mobile drawer NEVER mutates `document.body.style.overflow`; it intercepts wheel/touch events.
- [ ] Header action pair consists of one outline button (`variant="outline", size="sm"`) and one primary button (`variant="primary", size="sm"`).
