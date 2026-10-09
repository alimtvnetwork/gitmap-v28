# Subtask 04 — Remaining tabs compact polish (WS4, Worker 02)

## Owned files (ONLY these)
- Ops-tabs region of `cli/cmdui/ui_assets_markup_ops.go`: ONLY lines ~3–141
  (the 6 tab panes: commitin, ssh, macro, installer, prompts, import-export,
  schedules). Do NOT touch the `uiAssetsMarkupSettingsA` constant below
  (subtask 03's region, same worker — do subtask 03 first).

## Objective
Compact, professional ops tabs. Implements spec
`02-spec/21-app/249-ui-accessibility/03-compact-polish.md` §3 (+ §1 density).

## Changes
1. **Heading scale**: card `h3` stays `--fs-md`; demote the ssh info-box `h4`
   (browser-default 16px, larger than card h3) to `--fs-sm`; any
   `<strong style="font-size:0.95rem">` headers → `--fs-sm` class.
2. **Button-row double spacing**: where a flex `gap` already exists (ssh row,
   import-export row), remove the extra `margin-left: var(--sp-2)` effect —
   do it with a local utility class or inline-style removal, NOT by editing
   the global `.btn-secondary` rule (that's Worker 01's theme file).
3. **Labels**: pair every `<label>` with its control via `for`/`id` across all
   6 panes (checkboxes, inputs, selects, textareas).
4. **Density**: remove pane-level inline padding/margin overrides that fight
   the tightened card/form-group tokens; keep markup changes minimal.
5. **Ping buttons** (ssh nodes table is JS-rendered — markup here is the ssh
   pane shell only): no change needed here; `testNode()` stub is WS1's.

## Out of scope
- Settings panes (subtask 03, do first). JS files, theme CSS, shell markup.
- Do NOT touch `ui_assets_markup_settings_b.go`, `ui_assets_markup_misc.go`.

## Acceptance
- Consistent heading hierarchy (no h4 larger than h3); no double-spaced
  button rows; zero bare labels; file stays ≤300 lines.
