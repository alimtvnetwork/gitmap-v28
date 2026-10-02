# 33 — Avant-Garde Mega Menu Components & Dropdown System

> **/goal** Master and enforce the component architecture, physics parameters, staggered entrance timings, left-border growth rules, SlideSwapLabel keyframes, and 3D promotional flip cards of the Avant-Garde Mega Menu.
> **/learn** Master the exact entrance easing `cubic-bezier(0.16, 1, 0.3, 1)`, group stagger delay formula (`0.05s + gi * 0.05s`), link entrance delay formula (`0.08s + gi * 0.05s + li * 0.03s`), left hairline growth (`scaleY(0) -> scaleY(1)` over 420ms), SlideSwapLabel CSS keyframes, and 3D flip card physics (`perspective: 1400px`, `rotateY(180deg)` over 820ms).

**Version:** 4.1.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Executive System Overview

The **Avant-Garde Mega Menu** is an ultra-polished, multi-column navigation surface deploying underneath the 72px sticky glass header:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ [Column 1: Enterprise]   [Column 2: Cloud]    [Column 3: AI]  [3D Flip Card]│
│  • ERP Core               • Lakehouse DDL      • LLM Agent     ┌───────────┐│
│  • Supply Chain           • Vector Storage     • Prompt Spec   │ Front     ││
│  • FinTech Audit          • Real-Time Streams  • Neural Flow   │ (Flip 3D) ││
│                                                                └───────────┘│
└─────────────────────────────────────────────────────────────────────────────┘
  ▲─── Absolute Left-0 Right-0 Top-Full Z-40, Pt-3, Container Enclosed ───────▲
```

---

## 2. Component Taxonomy

| Component ID | Visual Role | Key Layout & Behavior |
|:---|:---|:---|
| **`SiteHeader`** | Sticky Viewport Anchor | `72px` height, `top-0 z-50`, `bg-background/90`, `backdrop-blur-xl`. |
| **`MegaPanel`** | Full-Bleed Dropdown Surface | Mounts at `top-full pt-3 z-40`, animates from `y: -8, scale: 0.985`. |
| **`MegaGroup`** | Categorized Link Column | Uppercase mono eyebrow title with 1–6 nested `MegaLink` items. |
| **`MegaLink`** | Interactive Navigation Row | Left growing accent hairline, `SlideSwapLabel`, trailing arrow reveal. |
| **`PromoFlipCard`** | Interactive Right-Side Feature | 3D two-sided card with `perspective: 1400px` flipping on hover. |
| **`HeaderCtaPair`** | Right Header Actions | Outline button (`variant="outline", size="sm"`) + Primary button (`size="sm"`). |
| **`MobileDrawer`** | Responsive Mobile Surface | Fullscreen drawer with sticky-safe wheel/touch suppression. |

---

## 3. Dropdown Motion & Entrance Physics

All transitions are coordinated using hardware-accelerated transforms and explicit stagger mathematics:

```
Panel Open:   0.00s ───► 0.26s  (Opacity 0 -> 1, Y -8 -> 0, Scale 0.985 -> 1)
Group Stagger:0.05s ───► 0.33s  (Opacity 0 -> 1, Y 8 -> 0, Delay: 0.05s + gi * 0.05s)
Link Stagger: 0.08s ───► 0.34s  (Opacity 0 -> 1, X -6 -> 0, Delay: 0.08s + gi*0.05s + li*0.03s)
Promo Entrance:0.14s ──► 0.44s  (Opacity 0 -> 1, Y 10 -> 0, Duration: 0.30s)
```

| Phase | Target Element | Initial State | Animate State | Exit State | Timing & Easing |
|:---|:---|:---|:---|:---|:---|
| **Panel Surface** | `MegaPanel` | `opacity: 0, y: -8, scale: 0.985` | `opacity: 1, y: 0, scale: 1` | `opacity: 0, y: -6, scale: 0.99` | `260ms`, `[0.16, 1, 0.3, 1]` |
| **Column Group** | `MegaGroup` | `opacity: 0, y: 8` | `opacity: 1, y: 0` | — | `280ms`, delay: `0.05s + gi * 0.05s` |
| **Link Item** | `MegaLink` | `opacity: 0, x: -6` | `opacity: 1, x: 0` | — | `260ms`, delay: `0.08s + gi*0.05s + li*0.03s` |
| **Feature Promo** | `PromoFlipCard` | `opacity: 0, y: 10` | `opacity: 1, y: 0` | — | `300ms`, delay: `0.14s` |
| **Reduced Motion**| All Surfaces | `opacity: 0` | `opacity: 1` | `opacity: 0` | `120ms linear` |

---

## 4. Multi-Column Grid Responsive Templates

```typescript
const gridColumnClass =
  groups.length >= 3
    ? "lg:grid-cols-[1fr_1fr_1fr_0.9fr]"
    : groups.length === 2
      ? "lg:grid-cols-[1fr_1fr_1.1fr]"
      : "lg:grid-cols-[1.6fr_1fr]";
```

- **Panel Card Shell:** `rounded-[var(--radius-card,20px)] border border-border bg-card shadow-[var(--shadow-lift)] overflow-hidden`.
- **Internal Padding:** `p-8` (`32px`), gap: `gap-8` (`32px`).
- **Eyebrow Header:** `font-mono text-[12px] uppercase tracking-[0.16em] text-muted-foreground font-semibold`.

---

## 5. `MegaLink` Anatomy, `SlideSwapLabel` & Hover Mechanics

1. **Outer Boundary:** `rounded-[10px] px-3 py-2 block relative overflow-hidden transition-colors hover:bg-[color-mix(in_oklab,var(--primary)_7%,transparent)]`.
2. **Growing Left Hairline:** `absolute inset-y-1 left-0 w-px origin-top scale-y-0 bg-[image:var(--gradient-accent)] transition-transform duration-[var(--dur-base,420ms)] ease-[var(--ease-out)] group-hover:scale-y-100`.
3. **Trailing Arrow:** Lucide `ArrowRight` (`size-3.5`), `-translate-x-1 opacity-0 group-hover:translate-x-0 group-hover:opacity-100 transition-all duration-[var(--dur-fast,240ms)]`.

### 5.1 `SlideSwapLabel` Component & CSS Keyframes

```tsx
export function SlideSwapLabel({ children, className, stagger = 0.018 }: { children: string; className?: string; stagger?: number }) {
  const reduced = useReducedMotion();
  if (reduced) return <span className={className}>{children}</span>;
  const chars = children.split("");
  return (
    <span className={cn("slide-swap relative inline-flex overflow-hidden align-bottom", className)}>
      <span className="sr-only">{children}</span>
      <span aria-hidden className="inline-flex">
        {chars.map((c, i) => (
          <span key={`${c}-${i}`} className="relative inline-block overflow-hidden">
            <span className="slide-swap-top inline-block whitespace-pre" style={{ transitionDelay: `${i * stagger}s` }}>{c}</span>
            <span aria-hidden className="slide-swap-bottom absolute left-0 top-0 inline-block whitespace-pre" style={{ transitionDelay: `${i * stagger}s` }}>{c}</span>
          </span>
        ))}
      </span>
    </span>
  );
}
```

```css
.slide-swap { line-height: 1.15; }
.slide-swap-top, .slide-swap-bottom { transition: transform 520ms cubic-bezier(0.16, 1, 0.3, 1); }
.slide-swap-bottom { transform: translateY(100%); }
.slide-swap:hover .slide-swap-top, .group:hover .slide-swap-top, a:hover > .slide-swap .slide-swap-top, button:hover .slide-swap-top { transform: translateY(-110%); }
.slide-swap:hover .slide-swap-bottom, .group:hover .slide-swap-bottom, a:hover > .slide-swap .slide-swap-bottom, button:hover .slide-swap-bottom { transform: translateY(0); }
@media (prefers-reduced-motion: reduce) { .slide-swap-top, .slide-swap-bottom { transition: none; } }
```

---

## 6. 3D Promotional Flip Card (`PromoFlipCard`)

The right-hand column showcases announcements using pure 3D hardware-accelerated card rotation:

- **Perspective Container:** `group/promo relative min-h-[220px] [perspective:1400px]`.
- **Card Core:** `relative h-full w-full transition-transform duration-[820ms] ease-[var(--ease-out)] [transform-style:preserve-3d] group-hover/promo:[transform:rotateY(180deg)]`.
- **Front Face:** `absolute inset-0 flex flex-col justify-between gap-4 overflow-hidden rounded-[var(--radius-card,20px)] bg-[image:var(--gradient-accent)] p-6 text-white [backface-visibility:hidden]`.
- **Back Face:** `absolute inset-0 flex flex-col justify-between gap-4 overflow-hidden rounded-[var(--radius-card,20px)] border border-border bg-card p-6 text-foreground [backface-visibility:hidden] [transform:rotateY(180deg)]`.
- **CTA Button:** `inline-flex items-center justify-center gap-2 rounded-full bg-[image:var(--gradient-accent)] px-4 py-2.5 text-sm font-semibold text-white shadow-[var(--shadow-lift)] transition-transform duration-[var(--dur-fast)] hover:scale-[1.02]`.

---

## 7. Complete Reference Implementation (`MegaMenu.tsx`)

```tsx
export function MegaPanel({ panelKey, groups, promo, panelRef, onNavigate, onMouseEnter, onMouseLeave }: MegaPanelProps) {
  const reduced = useReducedMotion();
  return (
    <motion.div
      ref={panelRef} key={panelKey}
      initial={reduced ? { opacity: 0 } : { opacity: 0, y: -8, scale: 0.985 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      exit={reduced ? { opacity: 0 } : { opacity: 0, y: -6, scale: 0.99 }}
      transition={{ duration: reduced ? 0.12 : 0.26, ease: [0.16, 1, 0.3, 1] }}
      className="absolute left-0 right-0 top-full z-40 pt-3"
      onMouseEnter={onMouseEnter} onMouseLeave={onMouseLeave}
    >
      <div className="mx-auto max-w-[1280px] px-6">
        <div className="overflow-hidden rounded-[var(--radius-card,20px)] border border-border bg-card shadow-[var(--shadow-lift)]">
          <div className={cn("grid gap-8 p-8", groups.length >= 3 ? "lg:grid-cols-[1fr_1fr_1fr_0.9fr]" : groups.length === 2 ? "lg:grid-cols-[1fr_1fr_1.1fr]" : "lg:grid-cols-[1.6fr_1fr]")}>
            {groups.map((group, gi) => (
              <motion.div key={group.title} initial={reduced ? false : { opacity: 0, y: 8 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.28, delay: 0.05 + gi * 0.05, ease: [0.16, 1, 0.3, 1] }} className="flex flex-col gap-3">
                <p className="font-mono text-xs uppercase tracking-[0.16em] text-muted-foreground font-semibold">{group.title}</p>
                <ul className="flex flex-col gap-1">
                  {group.links.map((link, li) => (
                    <motion.li key={link.href} initial={reduced ? false : { opacity: 0, x: -6 }} animate={{ opacity: 1, x: 0 }} transition={{ duration: 0.26, delay: 0.08 + gi * 0.05 + li * 0.03 }}>
                      <a href={link.href} onClick={onNavigate} className="group relative block overflow-hidden rounded-[10px] px-3 py-2 transition-colors hover:bg-[color-mix(in_oklab,var(--primary)_7%,transparent)]">
                        <span aria-hidden className="absolute inset-y-1 left-0 w-px origin-top scale-y-0 bg-[image:var(--gradient-accent)] transition-transform duration-[var(--dur-base,420ms)] ease-[var(--ease-out)] group-hover:scale-y-100" />
                        <span className="flex items-center gap-1.5 font-display text-sm font-medium text-foreground">
                          <SlideSwapLabel stagger={0.018}>{link.label}</SlideSwapLabel>
                          <ArrowRight className="size-3.5 -translate-x-1 opacity-0 transition-all group-hover:translate-x-0 group-hover:opacity-100" />
                        </span>
                        {link.description ? <span className="mt-0.5 block text-xs leading-relaxed text-muted-foreground">{link.description}</span> : null}
                      </a>
                    </motion.li>
                  ))}
                </ul>
              </motion.div>
            ))}
            {promo ? (
              <motion.div initial={reduced ? false : { opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.3, delay: 0.14, ease: [0.16, 1, 0.3, 1] }} className="group/promo relative min-h-[220px] [perspective:1400px]">
                <div className="relative h-full w-full transition-transform duration-[820ms] ease-[var(--ease-out)] [transform-style:preserve-3d] group-hover/promo:[transform:rotateY(180deg)]">
                  <div className="absolute inset-0 flex flex-col justify-between gap-4 overflow-hidden rounded-[var(--radius-card,20px)] bg-[image:var(--gradient-accent)] p-6 text-white [backface-visibility:hidden]">
                    <div>
                      <p className="text-base font-bold leading-snug">{promo.front.title}</p>
                      <p className="mt-2 text-sm text-white/85">{promo.front.body}</p>
                    </div>
                    <span className="inline-flex items-center gap-1.5 text-xs font-semibold uppercase tracking-[0.14em] text-white/80">Hover to flip <ArrowRight className="size-3.5" /></span>
                  </div>
                  <div className="absolute inset-0 flex flex-col justify-between gap-4 overflow-hidden rounded-[var(--radius-card,20px)] border border-border bg-card p-6 text-foreground [backface-visibility:hidden] [transform:rotateY(180deg)]">
                    <div>
                      <p className="text-base font-bold leading-snug">{promo.back.title}</p>
                      <p className="mt-2 text-sm text-muted-foreground">{promo.back.body}</p>
                    </div>
                    <a href={promo.back.cta.href} onClick={onNavigate} className="group/cta inline-flex items-center justify-center gap-2 rounded-full bg-[image:var(--gradient-accent)] px-4 py-2.5 text-sm font-semibold text-white shadow-[var(--shadow-lift)] transition-transform duration-[var(--dur-fast,240ms)] hover:scale-[1.02]">
                      {promo.back.cta.label}
                      <ArrowRight className="size-4 transition-transform duration-[var(--dur-base,420ms)] group-hover/cta:translate-x-1" />
                    </a>
                  </div>
                </div>
              </motion.div>
            ) : null}
          </div>
        </div>
      </div>
    </motion.div>
  );
}
```

---

## 8. Anti-Hallucination & Quality Verification Checklist

- [ ] Outer dropdown container mounts at `top-full pt-3 z-40`.
- [ ] Panel entrance duration is strictly `260ms` with ease `cubic-bezier(0.16, 1, 0.3, 1)`.
- [ ] Group stagger entrance follows `0.05s + gi * 0.05s` with `y: 8 -> 0`.
- [ ] Link entrance follows `0.08s + gi * 0.05s + li * 0.03s` with `x: -6 -> 0`.
- [ ] Growing left hairline uses `origin-top scale-y-0 -> scale-y-1` over `420ms`.
- [ ] Trailing arrow is `14px` (`size-3.5`) and reveals from `-translate-x-1 opacity-0` to `translate-x-0 opacity-100`.
- [ ] 3D promo flip card utilizes `perspective: 1400px` and rotates `180deg` over `820ms`.
- [ ] Reduced motion suppresses all 3D rotations, translates, and scales into an instant `120ms` opacity fade.
