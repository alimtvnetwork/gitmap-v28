# Spec 249.3 — Compact professional polish

## 1. Density targets (4pt scale only, no magic numbers)
- `.card` padding `var(--sp-5)` → `var(--sp-4)`; card bottom margin → `var(--sp-4)`.
- `.form-group` bottom margin 16px → `var(--sp-3)`.
- `.content-area` padding `var(--sp-6)` → `var(--sp-5)`.
- Sub-tab bar: remove inline `margin-bottom:1.25rem; padding-bottom:0.75rem`;
  single `margin-bottom: var(--sp-4)` via the token scale.

## 2. Settings page (all 7 panes + savebar)
- Unify savebar: delete the inline style override
  (`ui_assets_markup_settings_b.go:190`); `.settings-savebar` CSS is the single
  source of truth (sticky, `z-index:10`, flush to content bottom).
- Catalog cards: cut padding to `var(--sp-3)`, collapse internal stack margins;
  target ≤120px per card.
- Terminal pane: `min-height` 260px → 180px; output box padding → `var(--sp-3)`.
- Badge colors: replace hardcoded hex (`#1e3a5f` etc.) with theme badge tokens.
- Label `for`/`id` pairing on every control (see 249.2 §5).

## 3. Ops tabs (commitin, ssh, macro, installer, prompts, import-export, schedules)
- Normalize heading scale: card `h3` = `--fs-md`; demote oversized `h4`s
  (ssh info box) to `--fs-sm`; catalog-style `<strong>` headers to `--fs-sm`.
- Fix `.btn-secondary { margin-left: var(--sp-2) }` double-spacing where a
  flex `gap` already exists (ssh + import-export button rows): scope the
  margin to non-flex contexts or drop it.
- Label pairing on every control; compact `.form-group` rhythm (§1).

## 4. Editor + help tabs
- `#editor-container`: fixed 500px → `clamp(320px, 60vh, 640px)`.
- Help tab: add third card "Dashboard Keyboard Shortcuts" (mirrors the `?`
  overlay content, static table).
- Label pairing for editor controls.

## 5. Hover/focus refinement
- Buttons/inputs/tabs: 150ms ease transitions already exist; ensure hover
  states use `--primary-hover` / raised surfaces consistently, no dead hovers.
- Focus ring (§249.2.2) must be visible on every interactive element in both
  themes — spot-check light theme contrast.

## 6. Out of scope
- New tabs, new features behind dead handlers (stubs only), SSE, re-theming
  (246's token system stays).
