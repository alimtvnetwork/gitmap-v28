# UI Issues & Fixes Log

Tracking resolved UI issues to prevent regressions.

---

## Issue 1: Code Block Text Selection Contrast

**Problem:** Selected text in dark-themed code blocks used dark background + dark text, making it unreadable.

**Root Cause:** Global `::selection` used `hsl(var(--primary) / 0.2)` which was too subtle in dark mode, and inherited dark foreground colors clashed with syntax-highlighted text.

**Fix:** Added dedicated `pre ::selection, code ::selection, .hljs ::selection` rules with lighter background (`hsl(142 65% 40% / 0.35)`) and forced white text (`hsl(0 0% 100%)`). Dark mode global selection also uses lighter green (`hsl(142 71% 55% / 0.3)`) with near-white text.

**Rule:** Always ensure text selection has sufficient contrast — light text on colored backgrounds in dark mode, dark text on colored backgrounds in light mode.

---

## Issue 2: Heading Colors Too Bright on Light Backgrounds

**Problem:** `docs-h1`, `docs-h2`, `docs-h3` used bright green gradients that looked washed out on white backgrounds.

**Root Cause:** Same `hsl(var(--primary))` (142 71% 45%) was used for both light and dark themes.

**Fix:** Light mode headings now use darker green (`hsl(142 65% 32%)` to `hsl(142 55% 42%)`). Dark mode retains the original bright green via `.dark` overrides.

**Rule:** Always define separate light/dark color values for headings — dark green on light backgrounds, bright green on dark backgrounds.

---

## Issue 3: Pinned Line Selection Too Dim in Dark Mode

**Problem:** Pinned/selected lines in code blocks were barely visible in dark mode.

**Root Cause:** `.code-line-pinned` used `hsl(var(--primary) / 0.12)` background — too low opacity on dark backgrounds.

**Fix:** Increased to `hsl(var(--primary) / 0.22)` with full-opacity border color.

**Rule:** Pinned/highlighted states in code blocks must use ≥0.2 opacity for background tints.

---

## Issue 4: Vertical Translation Animations Causing Layout Shift

**Problem:** Feature cards and fade-in animations used `translateY` causing text to "jump" on page load.

**Root Cause:** `@keyframes fade-in` included `from { transform: translateY(10px) }` and cards had `hover:-translate-y-1`.

**Fix:** Removed all `translateY` from fade-in keyframes. Removed `hover:-translate-y-1` and `scale-110` from feature cards. Only horizontal slides or opacity transitions allowed.

**Rule:** Never use `translateY` or scale animations on text-containing elements. Allowed: opacity fades, horizontal slides, border/color transitions.

---

## General Contrast Rules

1. **Text selection:** Selected text must always be readable — use white/near-white text on colored selection backgrounds in dark mode.
2. **Headings:** Use darker color variants on light backgrounds; brighter variants on dark backgrounds.
3. **Interactive states (hover, pinned, active):** Must have ≥0.2 opacity difference from resting state on dark backgrounds.
4. **Animations:** No vertical position changes on text. No zoom/scale on content containers.
