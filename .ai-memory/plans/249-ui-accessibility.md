# Plan 249 — Dashboard UI keyboard accessibility + compact professional polish

## User request (verbatim)
"In digit map, setting UI and other UI, the UIs are not that good. Can you please work on it on all the UI settings UI and other parts of the UI which we have like SEO UI make it more professional. Nice and compact looking put into the UI UX concepts high level concept smaller stuff focus stuff also have the accessibility like tab and shortcuts to access to buttons pages and do this stuff from the keyboard from the UI section so make it more professional currently it's not can you please do that for me"

## Scope
- Repo: `~/workspace/repos/gitmap-v28`, main at 49c8631 (pulled 2026-10-09, tree has one unrelated dirty file: `.ai-memory/plans/246-parallel-fix-command.md` — do NOT touch).
- Backup branch `backup/249-ui-accessibility` created from 49c8631 and pushed to origin (via git direct; installed `gitmap` binary predates `space backup-branch`, dogfood command unavailable in it).
- Builds on program 246's dashboard overhaul (theme system, SVG icons, toasts, 4pt scale, responsive shell, settings page).
- Spec: `02-spec/21-app/249-ui-accessibility/` (to be written).

## Conflicts (logged, Top-Instruction Priority)
- V6 skill R1 bans `go build`; OWNER explicitly authorized builds + serving the dashboard for verification. Owner wins.
- V6 skill R8/R9 wants ONE atomic commit at Phase 3; parent task demands atomic commits per wave + immediate push. Parent wins.

## Wave plan
- Wave 0 — Research: 2 research subagents inventory dashboard assets, pages/tabs, keyboard/ARIA state, SEO UI.
- Wave 1 — Spec 249 + subtasks (lead writes after research synthesis).
- Wave 2 — Workers (A=2, H=2 = 4 workstreams, disjoint file boxes):
  - WS1: keyboard accessibility core — focus-visible styles, skip-links, tab order fixes, ARIA roles/labels (theme CSS + shell markup).
  - WS2: keyboard shortcuts engine + `?` shortcut-map overlay (JS).
  - WS3: compact polish — settings page + SEO-related UI.
  - WS4: compact polish — remaining 8 tabs.
- Wave 3 — Lead verification: `go build ./...`, serve dashboard, HTTP 200 every page, keyboard-walk settings page, screenshot review. Then atomic `gitmap cpf` + push per wave.

## Checkboxes
- [x] Step 0: git pull (49c8631), backup branch `backup/249-ui-accessibility` pushed.
- [x] Wave 0: Research 01 received (asset inventory: 15 files/2763 lines; 10 tabs; keyboard/ARIA gaps mapped; zero-overlap WS boxes proposed). Research 02 received (SEO UI: NOT FOUND — no such surface; settings control inventory + 10 dangling handlers; compactness audit with 10 file:line items; help tab is home for shortcut docs).
- [x] Wave 1: spec 249 written (`02-spec/21-app/249-ui-accessibility/` 01-overview/02-keyboard-accessibility/03-compact-polish) + 4 subtasks. Decisions: D1 no SEO tab invented; D2 dead handlers get toast stubs; D3 overlay built in JS; D4 keep Go-const asset model; D5 no new tabs/SSE.
- [x] Wave 2: WS3 settings + SEO compact polish DONE (Worker 02: savebar unified, catalog compacted, badge tokens, terminal focusable, 21/21 labels paired).
- [x] Wave 2: WS4 remaining tabs compact polish DONE (Worker 02: heading scale fixed, button-row double-spacing fixed, 11/11 labels paired, hex→tokens).
- [x] Wave 2: WS1 keyboard a11y core DONE (Worker 01: focus system, skip-link, tablist ARIA, toast/scrim keyboard, 10 toast stubs, labels paired).
- [x] Wave 2: WS2 shortcuts + overlay DONE (Worker 01: `?` overlay, `1-0`/`[`/`]`/arrows/Esc, showTab hardening, help-tab shortcuts card).
- [x] Wave 3: lead verification DONE — `go build ./...` exit 0 (lead-run); served dashboard `/` → 200, all 10 tab panes present, `/api/settings` → 200; a11y markup verified in served HTML (skip-link, tablist, overlay, stubs, help card); JS braces/parens balanced; secrets gate clean.
- [ ] Wave 3: lead verification (build exit 0, all pages 200, keyboard-walk, screenshots).
- [ ] Commits pushed per wave; final report to parent.

## Assumptions
- Dashboard = Go string-constant embedded assets in `cli/cmdui/` (zero-dependency single binary; no npm/build step) — keep that model.
- 10 tabs from 246: settings, commitin, ssh, macro, installer, prompts, import-export, schedules, editor, help.
- "SEO UI" location TBD by research (may be a settings pane or separate surface).

## Follow-ups (out of scope)
- SSE live telemetry (declined in 246, stays declined).
- New dashboard/cluster/templates/repos tabs (declined in 246).
