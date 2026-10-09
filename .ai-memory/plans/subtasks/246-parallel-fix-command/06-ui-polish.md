# Subtask 06 — UI Polish: Icons, Toasts, Responsive, Honesty Fixes (program 246)

## Objective

Implement spec `02-spec/21-app/246-parallel-fix-command/04-ui-ux-overhaul.md`
§5–§8: SVG iconography, toast notifications, card tiers, responsive sidebar,
demo-data removal, `ui` flag parsing, honest help menu, and the single sticky
Settings save bar. (Theme tokens themselves are subtask 05.)

## Read first

- `02-spec/21-app/246-parallel-fix-command/04-ui-ux-overhaul.md` (§5–§8)
- `.ai-memory/what-to-read.md` — "The Mindset" (§5: ~300 lines max per file)

## Owned files (ONLY these — do not touch any other file)

- `cli/cmdui/ui_assets_markup_misc.go` — sidebar shell: icon+label nav buttons,
  hamburger, `.table-wrap`, header title hook (style const lives in subtask 05's
  new file — do not duplicate it here)
- `cli/cmdui/ui_assets_markup_ops.go` — icon+label sub-tab pills
- `cli/cmdui/ui_assets_markup_settings_b.go` — remove demo data; single sticky save bar
- `cli/cmdui/ui_assets_js_fleet.go` — `alert(` → `showToast(` swaps
- `cli/cmdui/ui_assets_js_nodes.go` — `alert(` → `showToast(` swaps; delete phantom
  `u1` injection; `showTab` syncs `#header-title` + `document.title`
- `cli/cmdui/ui_assets_js_terminal.go` — `alert(` → `showToast(` swaps (if any)
- `cli/cmdui/ui_cmd.go` — parse `-p/--port`, `--host`, `--no-browser` (cmdui-local FlagSet)
- `cli/cmdui/ui_help_menu.go` — honest 10-tab page list; drop the SSE claim

## Implementation checklist

**Icons** (inline SVG via the `icon()` helper from subtask 05; Lucide path data):

- [ ] Sidebar: settings→gear, commitin→package, ssh→server, macro→scroll-text,
      installer→rocket, prompts→message-square, import-export→arrow-left-right,
      schedules→clock, editor→pen-line, help→book-open. Labels hidden in rail mode
      (`title` attr keeps names). Zero emoji left in nav chrome/buttons/headers.

**Toasts:**

- [ ] `showToast(message, kind)` with `#toast-stack` styles (styles live in
      subtask 05's theme file — wire the JS + container markup here).
- [ ] Replace ALL 19 `alert(` sites: fleet deploy/export/import (6),
      nodes file read/save (4), settings save (3), terminal (remainder).
- [ ] After the swap: zero `alert(` remains in the assets (verify with
      `gitmap aum search "alert(" cli/cmdui --ext .go`).

**Cards & header:**

- [ ] `.card-hero` (elevated + amber top accent) on each tab's primary panel;
      `.card-flat` (border-only) for nested groups.
- [ ] `showTab` sets `#header-title` text and `document.title = "GitMap — <Tab>"`.
- [ ] Status indicator → live dot + theme-aware badge (drop "Connected: localhost").

**Responsive:**

- [ ] ≥1200px full 240px sidebar; 900–1200px icon rail 68px; <900px off-canvas +
      header hamburger toggling `.sidebar-open`.
- [ ] Tables wrapped in `.table-wrap { overflow-x: auto }`; settings panes stack
      single-column under 900px.

**Demo-data removal:**

- [ ] Delete phantom `u1` option injection (`ui_assets_js_nodes.go:46-50`).
- [ ] Delete hardcoded `<option value="u1">u1</option>`
      (`ui_assets_markup_settings_b.go:173`) — terminal select = `local` +
      `/api/ssh/nodes` only.
- [ ] Delete the fake Antigravity card (`PID: 4757` block); instances panel renders
      from `/api/instances` with empty state
      "No instances detected — click Scan Instances."
- [ ] Sweep assets for leftover demo literals (`4757`, standalone `u1` options).

**Honesty fixes:**

- [ ] `ui_cmd.go`: cmdui-local FlagSet parses `-p/--port` (default 8080), `--host`
      (default `127.0.0.1`), `--no-browser`; positional arg stays the page name.
      `--port` MUST NOT become the page name anymore.
- [ ] `ui_help_menu.go`: pages section lists exactly the 10 real tabs with one
      honest line each — do NOT invent dashboard/cluster/templates/repos tabs.
- [ ] `ui_help_menu.go:59`: "Native Go JSON endpoints with SSE live telemetry" →
      "Native Go JSON endpoints (poll-based live refresh)". Do NOT implement SSE.
- [ ] `--host` non-loopback prints the network-exposure WARNING (no prompt).

**Settings save bar:**

- [ ] Delete the 7 per-pane Save buttons; add ONE sticky `.settings-savebar`
      ("Save All Settings" → existing `saveSettings()`, toast feedback).

## Done criteria

- [ ] `gitmap ui --port 9090 --no-browser` binds 9090, no browser opens;
      `gitmap ui settings --port 8080` opens the settings page (flag ≠ page).
- [ ] Help menu shows the 10 real tabs; no SSE claim anywhere.
- [ ] Zero `alert(`, zero emoji in chrome, zero phantom nodes, zero fake cards.
- [ ] Every touched file ≤300 lines; no file deletions; no alias changes.

## Hard rules

- Search ONLY via `gitmap aum search` / `gitmap find` / `gitmap cat`. TOTAL BAN on
  grep, rg, git grep, Select-String.
- NEVER run git commands. NEVER run `go build` / `go test` (lead verifies centrally).
- Relative paths only. NEVER propose deleting files.
- Security invariant: default bind stays 127.0.0.1; terminal exec surface unchanged.
