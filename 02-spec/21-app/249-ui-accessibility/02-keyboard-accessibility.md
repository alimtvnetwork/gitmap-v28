# Spec 249.2 — Keyboard accessibility

## 1. Shortcut map (required)
All shortcuts are inert while focus is in `<input>`, `<textarea>`, `<select>`,
or `contenteditable` (except `Esc`, which always blurs/closes).

| Keys | Action |
|---|---|
| `?` (Shift+/) | Toggle shortcut-help overlay |
| `1`–`9`, `0` | Jump to tab 1–10 (tab order = nav order) |
| `[` / `]` | Previous / next tab |
| `Esc` | Close overlay, or blur focused control |
| `←` / `→` | Move between tabs when focus is inside the tablist (roving tabindex); same inside the settings sub-tab bar |

## 2. Focus system
- Keep `:focus-visible { outline: 2px solid var(--primary); outline-offset: 2px }`;
  add a `:focus` fallback with the same ring (Safari/older engines).
- Remove `#editor-area { outline: none }` (`ui_assets_theme.go:117`) — the editor
  is the worst place to hide focus.
- Toasts: `tabindex="0"`, `Enter`/`Esc` dismisses the focused toast.
- Sidebar scrim: keyboard-operable (Enter/Esc toggles).

## 3. ARIA contract
- Main nav: `role="tablist"` on `#nav-list`; each nav button `role="tab"`,
  `aria-selected`, `aria-controls="tab-<id>"`; each pane `role="tabpanel"`,
  `aria-labelledby` pointing at its tab button. `showTab`/`navGo` keep
  `aria-selected` in sync.
- Settings sub-tab bar: same tablist/tab/tabpanel pattern (8 sub-tabs).
- Hamburger: `aria-expanded` synced with sidebar state.
- Shortcut overlay: `role="dialog"`, `aria-modal="true"`,
  `aria-label="Keyboard shortcuts"`; focus moves to its close button on open,
  returns to the opener on close.
- Toast stack keeps `aria-live="polite"`.

## 4. Skip link
First element in `<body>`: `<a class="skip-link" href="#main-content">Skip to
main content</a>`; visually hidden until focused. The active content wrapper
gets `id="main-content"` + `tabindex="-1"`.

## 5. Labels
Every `<label>` gets a `for` paired with its control's `id` (settings panes,
editor tab, help tab). No bare labels remain.

## 6. Overlay behavior
- Built and injected by JS (no new Go markup asset).
- Focus the close button on open; `Esc` closes; focus returns to opener.
- Simple focus trap: `Tab`/`Shift+Tab` cycle within the dialog.
- Lists the §1 shortcut map in a two-column table; also linked from the Help
  tab as a third card ("Dashboard Keyboard Shortcuts").

## 7. Dead-handler stubs (honesty)
Define in global JS: `deployMacroFleetModal`, `loadMacros`, `saveInstaller`,
`testInstaller`, `formatPromptJSON`, `importPrompt`, `exportFullConfig`,
`importFullConfig`, `addSchedule`, `testNode` — each calls
`showToast('Not implemented yet', 'info')`. Zero ReferenceErrors.
