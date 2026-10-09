# Subtask 02 — Keyboard shortcuts engine + `?` overlay (WS2, Worker 01)

## Owned files (ONLY these)
- `cli/cmdui/ui_assets_js_fleet.go`
- `cli/cmdui/ui_assets_js_terminal.go`

## Objective
Global shortcut map + help overlay + arrow-key tab navigation. Implements spec
`02-spec/21-app/249-ui-accessibility/02-keyboard-accessibility.md` §§1, 6.

## Changes
1. **Global keydown** (in `ui_assets_js_fleet.go`, near the load block): ignore
   when `event.target` is input/textarea/select/contenteditable (except Esc).
   Map: `?` → toggle overlay; `1`–`9`,`0` → jump to tab N (nav order);
   `[`/`]` → prev/next tab; `Esc` → close overlay or blur.
2. **Overlay**: built + injected purely in JS (no Go markup change).
   `role="dialog"`, `aria-modal="true"`, `aria-label="Keyboard shortcuts"`;
   two-column table of the §1 map; close button focused on open; `Esc` closes;
   focus returns to opener; simple Tab/Shift+Tab trap inside the dialog.
3. **Arrow-key tab nav**: when focus is inside `#nav-list`, ArrowLeft/Right
   move between tabs (roving tabindex: active tab `tabindex="0"`, others `-1`).
   Coordinate with WS1's `role="tab"` wiring (attributes set by WS1's renderNav;
   this workstream only moves focus + calls `showTab`).
4. **showTab hardening** (`ui_assets_js_fleet.go:3`): stop relying on the
   implicit global `event?.target` for active-classing — resolve the button
   via `document.querySelector` on the tab id instead.
5. **Terminal** (`ui_assets_js_terminal.go`): ensure the terminal input's
   existing Enter/Up/Down handler doesn't swallow `?`/`[`/`]`/`Esc` —
   terminal keeps its bindings; global handler already skips inputs, but
   verify no double-handling.

## Out of scope
- Overlay CSS lives in `ui_assets_theme.go` (same worker owns it — do WS1's
  CSS work first, then add `.shortcuts-overlay` styles there using theme vars).
- Do NOT touch `ui_assets_js_nodes.go` (WS1's file — coordinate via the ARIA
  contract: WS1 sets roles/attrs, WS2 only reads them).

## Acceptance
- `?` opens/closes overlay on every tab; `1`–`0`/`[`/`]` switch tabs;
  arrows move within tablist; Esc behavior correct; no JS errors on load;
  every touched file ≤300 lines.
