# Subtask 01 — Keyboard accessibility core (WS1, Worker 01)

## Owned files (ONLY these)
- `cli/cmdui/ui_assets_theme.go`
- `cli/cmdui/ui_assets_markup_misc.go`
- `cli/cmdui/ui_assets_js_nodes.go`

## Objective
ARIA roles, skip-link, focus system, toast/scrim keyboard support, label pairing
for editor/help tabs, dead-handler stubs. Implements spec `02-spec/21-app/249-ui-accessibility/02-keyboard-accessibility.md` §§2–5, 7.

## Changes
1. **CSS** (`ui_assets_theme.go`): add `:focus` fallback matching the existing
   `:focus-visible` ring; DELETE `#editor-area { outline: none }`; add
   `.skip-link` (visually hidden until `:focus`); toast `:focus` style;
   keep everything on the 4pt/token scale.
2. **Shell markup** (`ui_assets_markup_misc.go`): skip-link as first `<body>`
   element → `#main-content`; add `id="main-content" tabindex="-1"` to the
   content wrapper; `role="tablist"` on `#nav-list`; hamburger gets
   `aria-expanded` (JS syncs it); scrim keyboard: `tabindex="0"` +
   keydown Enter/Esc → toggleSidebar.
3. **Nav ARIA** (`ui_assets_js_nodes.go`, `renderNav`): each generated button
   gets `role="tab"`, `aria-selected`, `aria-controls="tab-<id>"`,
   `id="tabbtn-<id>"`; each `.content-area` pane gets `role="tabpanel"` +
   `aria-labelledby="tabbtn-<id>"` (set once at render; `showTab`/`navGo` sync
   `aria-selected`).
4. **Sub-tab behavior** (`ui_assets_js_nodes.go`, `showSettingsSubTab`): add
   roving tabindex + ArrowLeft/ArrowRight between the 8 sub-tab buttons;
   toggle `aria-selected` (buttons already exist in markup with
   `class="sub-tab-btn"` — add `role="tab"` via JS to avoid touching Worker 02's
   markup file; panes get `role="tabpanel"` via JS too).
5. **Toasts** (`ui_assets_js_nodes.go`, `showToast`): `tabindex="0"`,
   keydown Enter/Esc dismisses the focused toast.
6. **Dead-handler stubs** (global JS, same file): `deployMacroFleetModal`,
   `loadMacros`, `saveInstaller`, `testInstaller`, `formatPromptJSON`,
   `importPrompt`, `exportFullConfig`, `importFullConfig`, `addSchedule`,
   `testNode` → each `showToast('Not implemented yet', 'info')`.
7. **Editor/help labels** (`ui_assets_markup_misc.go`): pair every `<label>`
   with its control via `for`/`id`.

## Out of scope
- Shortcut engine + `?` overlay (WS2). Settings/ops tab markup (Worker 02).
- Do NOT touch `ui_assets_js_fleet.go`, `ui_assets_js_terminal.go`.

## Acceptance
- No `outline: none` remains on any focusable element; skip-link present;
  `role=tablist/tab/tabpanel` + `aria-selected` wired; zero bare labels in
  editor/help; 10 stubs defined; every touched file ≤300 lines.
