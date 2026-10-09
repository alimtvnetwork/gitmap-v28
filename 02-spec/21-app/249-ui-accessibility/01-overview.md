# Spec 249.1 — Overview: dashboard keyboard accessibility + compact polish

## 1. Goal
Make `gitmap ui` fully keyboard-operable and visually compact/professional.
Builds on program 246 (theme system, SVG icons, toasts, 4pt scale, responsive
shell). Asset model stays: Go string constants in `cli/cmdui/`, zero-dependency
single binary, no npm/build step.

## 2. Scope
- All 10 tabs: settings, commitin, ssh, macro, installer, prompts, import-export,
  schedules, editor, help.
- Keyboard: tab order, visible focus, shortcuts + `?` overlay, skip-link, ARIA.
- Compact polish: density, hierarchy, hover/focus states, no placeholder content.

## 3. Non-goals / decisions
- **D1 — No SEO tab.** Verified: no SEO UI exists in the dashboard (zero hits for
  "seo" across all 15 `cli/cmdui/` files; server registers no SEO route).
  Closest existing surfaces: `seowrite` CLI (`cli/cmdseowrite/`) and the settings
  "UI/UX Repo Index" catalog. Inventing a tab to match a misunderstanding is
  backwards (same principle as 246 §6.2). An SEO dashboard tab is a new feature
  → follow-up, not this program.
- **D2 — Dead handlers get toast stubs, not implementations.** 10 dangling
  `onclick` handlers throw ReferenceError today (`deployMacroFleetModal`,
  `loadMacros`, `saveInstaller`, `testInstaller`, `formatPromptJSON`,
  `importPrompt`, `exportFullConfig`, `importFullConfig`, `addSchedule`,
  `testNode`). Each gets a stub that shows an honest "Not implemented yet"
  toast. Full features are follow-ups.
- **D3 — Shortcut overlay built in JS**, injected at runtime. No new markup
  asset file; no new overlay component in Go.
- **D4 — Keep the Go-const asset model.** No framework, no build step.
- **D5 — No new tabs, no SSE** (reaffirms 246 §6.2/§6.3).

## 4. Acceptance (program-level)
- Every interactive element reachable and operable by keyboard alone.
- `?` opens a shortcut map on every tab; `Esc` closes.
- axe-style manual audit: zero missing labels, zero focus-less interactive divs.
- Settings page fits ~40% more content above the fold vs baseline (measured by
  scroll height at 900px viewport).
- `go build ./...` exit 0; every tab loads HTTP 200; no console errors on load.
