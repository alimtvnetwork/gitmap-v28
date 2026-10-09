# Subtask 03 — Settings page compact polish (WS3, Worker 02)

## Owned files (ONLY these)
- `cli/cmdui/ui_assets_markup_settings_b.go` (panes 5–7 + savebar)
- Settings region of `cli/cmdui/ui_assets_markup_ops.go`: ONLY the
  `uiAssetsMarkupSettingsA` constant (settings panes 1–4, ~line 142 to end).
  Do NOT touch the ops-tabs markup above it (that's subtask 04, same worker —
  do this file's settings region first, then subtask 04's ops region).

## Objective
Compact, professional settings page. Implements spec
`02-spec/21-app/249-ui-accessibility/03-compact-polish.md` §2 (+ §1 density).

## Changes
1. **Savebar unification**: DELETE the inline `style="..."` on the savebar div
   (`ui_assets_markup_settings_b.go:190`). The `.settings-savebar` CSS class
   (owned by Worker 01) becomes the single source of truth. Keep the caption
   text + single Save button.
2. **Density**: rely on Worker 01's tightened `.card`/`.form-group` tokens;
   remove any pane-level inline padding/margin overrides that fight them.
3. **Sub-tab bar**: remove inline `margin-bottom:1.25rem; padding-bottom:0.75rem`
   (`ui_assets_markup_ops.go:143`); single token margin via CSS class.
4. **Catalog cards** (7 cards, `settings_b.go:39-137`): padding → compact,
   collapse internal stack margins; replace hardcoded badge hex colors
   (`#1e3a5f`, `#60a5fa`, …) with theme badge tokens (`--badge-ok-bg` etc. —
   check `ui_assets_theme.go` for exact token names first).
5. **Terminal pane**: `min-height:260px` → `180px`; output padding → compact;
   add `tabindex="0"` + `aria-label="Terminal output"` to `#term-output` so
   keyboard users can focus/scroll history.
6. **Labels**: pair every `<label>` with its control via `for`/`id` in all
   7 panes (generate stable ids: `setting-<key>` already exists on most).
7. **Instances pane**: keep the honest empty state ("No instances detected —
   click Scan Instances."); no demo data.

## Out of scope
- Sub-tab JS behavior (Worker 01). Ops tabs markup (subtask 04, do after).
- Do NOT touch JS files, theme CSS, `ui_assets_markup_misc.go`.

## Acceptance
- Zero inline style overrides on savebar/sub-tab bar; zero hardcoded hex in
  badges; zero bare labels; terminal output focusable; every touched file
  ≤300 lines.
