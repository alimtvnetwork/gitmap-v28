import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { Keyboard, ChevronDown, RefreshCw } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

/**
 * TabOrderMap — collapsible accessibility helper that derives the live
 * keyboard tab sequence from the DOM and renders it as a numbered list.
 *
 * The order is computed from the actual document at the moment the panel
 * opens (and on subsequent DOM mutations). Positive `tabindex` values
 * jump ahead of `0`/auto elements, matching real browser behavior.
 */

interface FocusEntry {
  step: number;
  label: string;
  /** Optional secondary text — usually from aria-describedby / aria-description. */
  sublabel?: string;
  tag: string;
  tabIndex: number;
  section: string;
  isSelf: boolean;
}

// ─────────────────────────────────────────────────────────────────────────────
// DOM derivation
// ─────────────────────────────────────────────────────────────────────────────

const FOCUSABLE_SELECTOR = [
  "a[href]",
  "button:not([disabled])",
  "input:not([disabled]):not([type='hidden'])",
  "select:not([disabled])",
  "textarea:not([disabled])",
  "[tabindex]:not([tabindex='-1'])",
  "audio[controls]",
  "video[controls]",
  "[contenteditable]:not([contenteditable='false'])",
].join(",");

/**
 * True when the element is rendered AND interactable. Filters out:
 *   - `disabled` / `aria-hidden` / hidden subtrees
 *   - zero-size boxes
 *   - `display:none` / `visibility:hidden` (self or ancestor)
 *   - `opacity: 0` anywhere in the ancestor chain (Radix overlay pattern)
 *   - elements covered by another overlay (modal backdrop, popover, etc.)
 */
const isVisible = (el: HTMLElement): boolean => {
  if (el.hasAttribute("disabled")) return false;
  if (el.getAttribute("aria-hidden") === "true") return false;
  if (el.closest("[aria-hidden='true']")) return false;

  const rects = el.getClientRects();
  if (rects.length === 0) return false;

  // Walk ancestors checking display/visibility/opacity. An ancestor with
  // `opacity:0` or `visibility:hidden` makes the descendant non-interactable
  // even though the descendant's own computed style may still report 1.
  let node: HTMLElement | null = el;
  while (node && node !== document.body) {
    const s = window.getComputedStyle(node);
    if (s.display === "none") return false;
    if (s.visibility === "hidden" || s.visibility === "collapse") return false;
    if (parseFloat(s.opacity) === 0) return false;
    if (s.pointerEvents === "none" && node === el) {
      // Only block when the element itself opts out; ancestors with
      // pointer-events:none + a child that re-enables them are valid.
      return false;
    }
    node = node.parentElement;
  }

  // Overlay/cover detection — sample a few points within the element's
  // first rect and require at least one to hit the element (or a
  // descendant of it). If every sample resolves to an unrelated node
  // higher in the paint stack, treat it as covered.
  const r = rects[0];
  if (r.width < 1 || r.height < 1) return false;
  const samples: [number, number][] = [
    [r.left + r.width / 2, r.top + r.height / 2], // center
    [r.left + 2, r.top + 2],                       // top-left
    [r.right - 2, r.bottom - 2],                   // bottom-right
  ];
  const inViewport = samples.some(
    ([x, y]) => x >= 0 && y >= 0 && x <= window.innerWidth && y <= window.innerHeight,
  );
  if (!inViewport) {
    // Off-screen elements (e.g. below the fold) are still focusable; skip
    // the overlay test rather than wrongly excluding them.
    return true;
  }
  const isHitOrDescendant = samples.some(([x, y]) => {
    if (x < 0 || y < 0 || x > window.innerWidth || y > window.innerHeight) return false;
    const top = document.elementFromPoint(x, y);
    return !!top && (top === el || el.contains(top));
  });
  return isHitOrDescendant;
};


/** Resolve a space-separated id-list reference (aria-labelledby / aria-describedby). */
const resolveIdRefs = (ids: string | null): string => {
  if (!ids) return "";
  return ids
    .split(/\s+/)
    .map((id) => document.getElementById(id)?.textContent?.trim() ?? "")
    .filter(Boolean)
    .join(" ");
};

const truncate = (s: string, max = 80): string =>
  s.length > max ? `${s.slice(0, max - 3)}…` : s;

/**
 * Best-effort accessible name. Mirrors the ARIA name calculation order:
 * aria-label → aria-labelledby → associated <label> → title → text content.
 */
const labelFor = (el: HTMLElement): string => {
  const aria = el.getAttribute("aria-label");
  if (aria) return truncate(aria.trim());

  const labelledby = resolveIdRefs(el.getAttribute("aria-labelledby"));
  if (labelledby) return truncate(labelledby.replace(/\s+/g, " "));

  // <label for="..."> or wrapping <label>
  if (
    el instanceof HTMLInputElement ||
    el instanceof HTMLSelectElement ||
    el instanceof HTMLTextAreaElement
  ) {
    if (el.labels && el.labels.length > 0) {
      const text = Array.from(el.labels)
        .map((l) => l.textContent?.trim() ?? "")
        .filter(Boolean)
        .join(" ");
      if (text) return truncate(text.replace(/\s+/g, " "));
    }
  }

  const title = el.getAttribute("title");
  if (title) return truncate(title.trim());

  const text = (el.textContent ?? "").replace(/\s+/g, " ").trim();
  if (text) return truncate(text);

  if (el instanceof HTMLInputElement) return `${el.type} input`;
  return el.tagName.toLowerCase();
};

/**
 * Accessible description — what aria-describedby / aria-description would
 * make a screen reader announce *after* the name. Returned separately so
 * callers can render it as a sublabel without colliding with the name.
 */
const descriptionFor = (el: HTMLElement): string | undefined => {
  const direct = el.getAttribute("aria-description");
  if (direct?.trim()) return truncate(direct.trim(), 120);

  const referenced = resolveIdRefs(el.getAttribute("aria-describedby"));
  if (referenced) return truncate(referenced.replace(/\s+/g, " "), 120);

  return undefined;
};

/** Closest meaningful landmark for grouping in the list. */
const sectionFor = (el: HTMLElement): string => {
  const landmark = el.closest<HTMLElement>(
    "header, nav, main, aside, footer, [role='banner'], [role='navigation'], [role='main'], [role='complementary'], [role='contentinfo']",
  );
  if (landmark) {
    const role =
      landmark.getAttribute("role") ?? landmark.tagName.toLowerCase();
    return role.charAt(0).toUpperCase() + role.slice(1);
  }
  const section = el.closest<HTMLElement>("section[aria-labelledby], section[aria-label]");
  if (section) {
    const labelId = section.getAttribute("aria-labelledby");
    if (labelId) {
      const ref = document.getElementById(labelId);
      if (ref?.textContent) return ref.textContent.trim();
    }
    const label = section.getAttribute("aria-label");
    if (label) return label.trim();
  }
  return "Page";
};

/**
 * Collect focusable elements in real tab order.
 *
 * Browsers visit positive tabindex values (1, 2, …) first in ascending
 * order, then everything with tabindex 0 or no tabindex attribute in
 * document order. We replicate that with a stable sort.
 */
export const getTabOrder = (root: ParentNode = document.body): HTMLElement[] => {
  const all = Array.from(root.querySelectorAll<HTMLElement>(FOCUSABLE_SELECTOR));
  const visible = all.filter(isVisible);

  // Stable index for tie-breaking (document order).
  const indexed = visible.map((el, i) => ({
    el,
    i,
    ti: Number(el.getAttribute("tabindex") ?? "0"),
  }));

  indexed.sort((a, b) => {
    const aPositive = a.ti > 0;
    const bPositive = b.ti > 0;
    if (aPositive && !bPositive) return -1;
    if (!aPositive && bPositive) return 1;
    if (aPositive && bPositive && a.ti !== b.ti) return a.ti - b.ti;
    return a.i - b.i;
  });

  return indexed.map((x) => x.el);
};

// ─────────────────────────────────────────────────────────────────────────────
// Component
// ─────────────────────────────────────────────────────────────────────────────

const TabOrderMap = () => {
  const [open, setOpen] = useState(false);
  const [entries, setEntries] = useState<FocusEntry[]>([]);
  const [focusedStep, setFocusedStep] = useState<number | null>(null);
  const selfRef = useRef<HTMLElement | null>(null);
  const elementsRef = useRef<HTMLElement[]>([]);

  const refresh = useCallback(() => {
    const els = getTabOrder(document.body);
    elementsRef.current = els;
    const list: FocusEntry[] = els.map((el, idx) => {
      const label = labelFor(el);
      const desc = descriptionFor(el);
      // Skip sublabel if it duplicates the name (case/whitespace insensitive).
      const norm = (s: string) => s.toLowerCase().replace(/\s+/g, " ").trim();
      const sublabel = desc && norm(desc) !== norm(label) ? desc : undefined;
      return {
        step: idx + 1,
        label,
        sublabel,
        tag: el.tagName.toLowerCase(),
        tabIndex: Number(el.getAttribute("tabindex") ?? "0"),
        section: sectionFor(el),
        isSelf: !!selfRef.current && selfRef.current.contains(el),
      };
    });
    setEntries(list);
    // Re-resolve focused step against the newly-collected element list.
    const active = document.activeElement as HTMLElement | null;
    const idx = active ? els.indexOf(active) : -1;
    setFocusedStep(idx >= 0 ? idx + 1 : null);
  }, []);

  // Recompute on open + on DOM mutations + on resize while open.
  useEffect(() => {
    if (!open) return;
    refresh();

    let raf = 0;
    const schedule = () => {
      cancelAnimationFrame(raf);
      raf = requestAnimationFrame(refresh);
    };

    const mo = new MutationObserver(schedule);
    mo.observe(document.body, {
      subtree: true,
      childList: true,
      attributes: true,
      attributeFilter: ["disabled", "tabindex", "aria-hidden", "hidden"],
    });
    window.addEventListener("resize", schedule);
    return () => {
      cancelAnimationFrame(raf);
      mo.disconnect();
      window.removeEventListener("resize", schedule);
    };
  }, [open, refresh]);

  // Track focus globally while the panel is open. Uses focusin/focusout
  // (which bubble, unlike focus/blur) so we catch every change.
  useEffect(() => {
    if (!open) return;
    const onFocusIn = (e: FocusEvent) => {
      const target = e.target as HTMLElement | null;
      if (!target) {
        setFocusedStep(null);
        return;
      }
      const idx = elementsRef.current.indexOf(target);
      setFocusedStep(idx >= 0 ? idx + 1 : null);
    };
    const onFocusOut = () => {
      // Defer so the next focusin (if any) wins this frame.
      requestAnimationFrame(() => {
        const active = document.activeElement as HTMLElement | null;
        if (!active || active === document.body) {
          setFocusedStep(null);
        }
      });
    };
    document.addEventListener("focusin", onFocusIn);
    document.addEventListener("focusout", onFocusOut);
    // Seed with whatever currently has focus.
    onFocusIn({ target: document.activeElement } as unknown as FocusEvent);
    return () => {
      document.removeEventListener("focusin", onFocusIn);
      document.removeEventListener("focusout", onFocusOut);
    };
  }, [open]);

  // Group entries by section for readability while keeping global numbering.
  const grouped = useMemo(() => {
    const groups: { section: string; items: FocusEntry[] }[] = [];
    for (const e of entries) {
      const last = groups[groups.length - 1];
      if (last && last.section === e.section) last.items.push(e);
      else groups.push({ section: e.section, items: [e] });
    }
    return groups;
  }, [entries]);

  return (
    <section ref={(el) => { selfRef.current = el; }} className="py-6">
      <div className="flex items-center gap-2">
        <button
          type="button"
          onClick={() => setOpen((o) => !o)}
          aria-expanded={open}
          aria-controls="tab-order-map-panel"
          className="inline-flex items-center gap-2 rounded-md border border-border bg-card px-3 py-1.5 text-xs font-sans text-muted-foreground hover:text-foreground hover:border-primary/40 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
        >
          <Keyboard className="h-3.5 w-3.5" />
          <span>{open ? "Hide" : "Show"} tab order map</span>
          <motion.span
            animate={{ rotate: open ? 180 : 0 }}
            transition={{ duration: 0.2 }}
            className="inline-flex"
          >
            <ChevronDown className="h-3.5 w-3.5" />
          </motion.span>
        </button>
        {open && (
          <button
            type="button"
            onClick={refresh}
            className="inline-flex items-center gap-1.5 rounded-md border border-border bg-card px-2.5 py-1.5 text-xs font-sans text-muted-foreground hover:text-foreground hover:border-primary/40 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
            aria-label="Recompute tab order from the current DOM"
          >
            <RefreshCw className="h-3 w-3" />
            Refresh
          </button>
        )}
      </div>

      <AnimatePresence initial={false}>
        {open && (
          <motion.div
            id="tab-order-map-panel"
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.25, ease: "easeInOut" }}
            className="overflow-hidden"
          >
            <div className="mt-4 rounded-lg border border-border bg-muted/20 p-5">
              <p className="mb-4 font-sans text-xs text-muted-foreground">
                Derived live from the DOM ({entries.length} focusable element
                {entries.length === 1 ? "" : "s"}
                {focusedStep !== null && (
                  <>
                    {" · currently focused: "}
                    <span className="font-mono font-semibold text-primary">
                      #{focusedStep}
                    </span>
                  </>
                )}
                ). Use{" "}
                <kbd className="rounded border border-border bg-card px-1.5 py-0.5 font-mono text-[10px]">Tab</kbd>{" "}
                /{" "}
                <kbd className="rounded border border-border bg-card px-1.5 py-0.5 font-mono text-[10px]">Shift+Tab</kbd>{" "}
                to walk the order below.
              </p>

              {entries.length === 0 ? (
                <p className="font-sans text-xs text-muted-foreground">
                  No focusable elements detected.
                </p>
              ) : (
                <div className="space-y-4">
                  {grouped.map((group) => (
                    <div
                      key={`${group.section}-${group.items[0].step}`}
                      className="rounded-lg border border-dashed border-border bg-card/30 p-3"
                    >
                      <div className="mb-2 font-heading text-[11px] uppercase tracking-[0.18em] text-muted-foreground">
                        {group.section}
                      </div>
                      <ol className="space-y-1.5 list-none p-0 m-0">
                        {group.items.map((e) => {
                          const isActive = focusedStep === e.step;
                          return (
                            <li
                              key={e.step}
                              aria-current={isActive ? "true" : undefined}
                              className={[
                                "flex items-center gap-3 rounded-md border px-3 py-2 transition-colors",
                                isActive
                                  ? "border-primary bg-primary/10 ring-2 ring-primary/40 shadow-[0_0_0_3px_hsl(var(--primary)/0.15)]"
                                  : "border-border bg-card/60",
                              ].join(" ")}
                            >
                              <span
                                aria-hidden="true"
                                className={[
                                  "inline-flex h-5 w-5 shrink-0 items-center justify-center rounded-full text-[10px] font-mono font-bold shadow-sm",
                                  isActive
                                    ? "bg-primary text-primary-foreground ring-2 ring-primary/50 ring-offset-1 ring-offset-card scale-110"
                                    : "bg-primary text-primary-foreground",
                                ].join(" ")}
                              >
                                {e.step}
                              </span>
                              <div className="min-w-0 flex-1">
                                <div className="font-sans text-sm text-foreground truncate">
                                  {e.label}
                                  {e.isSelf && (
                                    <span className="ml-2 rounded border border-border bg-muted px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
                                      self
                                    </span>
                                  )}
                                  {isActive && (
                                    <span className="ml-2 rounded bg-primary/20 px-1.5 py-0.5 font-mono text-[10px] font-semibold text-primary">
                                      focused
                                    </span>
                                  )}
                                </div>
                                {e.sublabel && (
                                  <div className="font-sans text-xs text-muted-foreground italic mt-0.5 line-clamp-2">
                                    {e.sublabel}
                                  </div>
                                )}
                                <div className="font-mono text-[11px] text-muted-foreground/70 mt-0.5">
                                  &lt;{e.tag}&gt; · tabindex={e.tabIndex}
                                </div>
                              </div>
                            </li>
                          );
                        })}
                      </ol>
                    </div>
                  ))}
                </div>
              )}

              <p className="mt-4 font-sans text-[11px] text-muted-foreground">
                Order updates automatically when the DOM changes. Positive
                tabindex values jump ahead of <code className="font-mono">tabindex=&quot;0&quot;</code> /
                auto elements, matching browser behavior.
              </p>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </section>
  );
};

export default TabOrderMap;
