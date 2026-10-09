# Spec 04 — UI/UX Overhaul: Themes, Polish & Honesty Fixes (program 246)

## 1. Goal

Make `gitmap ui` genuinely premium — a real working theme system, coherent visual
language, honest help text — while **keeping the Go-string-constant asset model**
(zero-dependency single binary; no npm, no build step). Owner's words: "very
lucrative". This spec defines that as concrete, verifiable craft (see §2).

## 2. What "very attractive" means (concrete, verifiable)

- **Contrast:** body text ≥ 7:1 on surfaces (dark) and ≥ 4.5:1 (light); muted text
  ≥ 4.5:1 in both themes. No gray-on-gray telemetry.
- **Rhythm:** a 4pt spacing scale and a fixed type scale; zero magic-number
  padding/font-size values in the style block.
- **Hierarchy:** three card tiers (default / hero / flat), one primary action per
  view, section titles that scan.
- **Motion:** 150ms ease transitions on interactive elements; toast slide/fade;
  no layout shift on tab switch.
- **Craft details:** visible focus rings (keyboard a11y), tabular numerals for
  telemetry numbers, 8px radii, hairline borders instead of heavy boxes,
  dark code/terminal surfaces in BOTH themes (deliberate, premium feel).

Keep what works: the amber primary (`#f59e0b`), the 4-plane neutral depth
hierarchy (`--bg/--raised/--card/--popover`), the system font stack.

## 3. Ground truth defects (verified 2026-10-09, `cli/cmdui/`)

1. Theme select offers `dark`/`light` but **the light theme does not exist** — no
   light tokens, and no JS ever applies `data.theme` (dead control).
2. `var(--accent)` is referenced (`ui_assets_markup_ops.go:68`) but **never defined**.
3. Hardcoded hex mixed with vars: `#131d2e`, `#1e293b`, `#f8fafc`, `#e2e8f0`,
   `#cbd5e1` — these break any second theme.
4. Help menu advertises tabs `dashboard, cluster, templates, repos` — **none exist**.
   Real tabs: `settings, commitin, ssh, macro, installer, prompts, import-export,
   schedules, editor, help` (`ui_assets_markup_misc.go` shell).
5. `gitmap ui -p/--port/--host/--no-browser` are advertised but **not parsed** —
   `--port` becomes the page name (`ui_cmd.go`: `page = args[0]`).
6. "SSE live telemetry" is claimed (`ui_help_menu.go:59`); **no SSE endpoint exists**
   and no `EventSource` in the JS.
7. Save feedback is native `alert()` (19 call sites across fleet/nodes/settings JS).
8. Hardcoded demo data: phantom `u1` node injected into the terminal select
   (`ui_assets_js_nodes.go:46-50`) and hardcoded again (`ui_assets_markup_settings_b.go:173`);
   fake Antigravity card `PID: 4757` (`ui_assets_markup_settings_b.go:34`).
9. Emoji iconography in the sidebar chrome (⚙️📦🖥️…); flat 2015-admin look;
   fixed 240px sidebar, no responsive behavior; 7 redundant Save buttons on
   "All Settings".

## 4. Theme system

### 4.1 Token contract

- `:root` holds the **dark** theme (default, unchanged values, refined).
- `:root[data-theme="light"]` holds the **full light** token set.
- `--accent` is defined in both themes (fixes defect §3.2).
- New tokens: `--warning`, `--info`, `--radius-sm/md/lg`, `--font-sans`,
  `--font-mono`, `--shadow-card`, spacing scale `--sp-1…--sp-8`, type scale
  `--fs-xs…--fs-2xl`.

Dark (refined — keep current planes):

```css
--bg: hsl(230, 25%, 8%); --raised: hsl(230, 18%, 18%);
--card: hsl(230, 20%, 12%); --popover: hsl(230, 20%, 16%);
--border: hsl(230, 18%, 20%); --text: #f8fafc; --muted: #94a3b8;
--primary: #f59e0b; --primary-hover: #d97706; --primary-ink: #000000;
--accent: #fbbf24; --success: #10b981; --warning: #fbbf24; --info: #38bdf8;
--error: #ef4444; --code-bg: hsl(230, 25%, 6%);
```

Light (full set — amber darkened for contrast on white):

```css
--bg: #eef2f7; --raised: #ffffff; --card: #ffffff; --popover: #ffffff;
--border: #e2e8f0; --text: #0f172a; --muted: #64748b;
--primary: #b45309; --primary-hover: #92400e; --primary-ink: #ffffff;
--accent: #b45309; --success: #047857; --warning: #b45309; --info: #0284c7;
--error: #dc2626; --code-bg: hsl(230, 25%, 8%);  /* terminal stays dark */
```

Terminal/editor surfaces (`--code-bg`, `#editor-area`) stay dark in both themes —
documented, deliberate.

### 4.2 Switcher plumbing

- `loadSettings()` (existing, `ui_assets_js_nodes.go`) calls new `applyTheme(name)`:
  sets `document.documentElement.dataset.theme`, syncs the `#setting-theme` select.
- Theme select `onchange`: apply **immediately** (optimistic), then POST
  `/api/settings`; on failure show an error toast and revert the select.
- Server: **no change needed** — `theme` is already persisted in `SettingsData`
  (`ui_server.go` `loadSettings`/`saveSettingsData`, default `"dark"`).
- `prefers-color-scheme`: ignored — the persisted setting is the single source of
  truth (deterministic across machines).

### 4.3 Per-tab application

Sweep every hardcoded hex to its var equivalent (`#131d2e`→`var(--raised)`,
`#1e293b`→`var(--popover)`, text hexes→`var(--text)`/`var(--muted)`). All 10 tabs
must render legibly in both themes. Terminal/editor keep dark code surfaces.

### 4.4 Asset file budget

Extract the `uiAssetsStyle` const out of `ui_assets_markup_misc.go` (237 lines)
into a new `cli/cmdui/ui_assets_theme.go` (style + light overrides + toast CSS +
icon helper). Register it in the `IndexHTML` concat in `ui_assets.go` at the same
position. No other concat-order change. Every asset file ≤300 lines.

## 5. Visual polish (within the asset model)

### 5.1 Scales

```css
--sp-1: 4px; --sp-2: 8px; --sp-3: 12px; --sp-4: 16px;
--sp-5: 20px; --sp-6: 24px; --sp-8: 32px;
--fs-xs: 0.75rem; --fs-sm: 0.85rem; --fs-md: 0.95rem;
--fs-lg: 1.1rem; --fs-xl: 1.35rem; --fs-2xl: 1.6rem;
--radius-sm: 6px; --radius-md: 8px; --radius-lg: 12px;
```

Replace all ad-hoc sizes in the style block with these vars.

### 5.2 Icons — inline SVG, no emoji in chrome

New JS helper `icon(name)` returning an inline `<svg>` string
(24×24 viewBox, `fill="none"`, `stroke="currentColor"`, `stroke-width="2"` —
Lucide-style; implementer copies path data from lucide.dev for each named icon).
Required set: `settings` (gear), `package`, `server`, `scroll-text`, `rocket`,
`message-square`, `arrow-left-right`, `clock`, `pen-line`, `book-open`,
`check`, `alert-triangle`, `info`, `x`.

Sidebar buttons become `icon + label`; labels hide in rail mode (`title` attr keeps
the name). Sub-tab buttons and toast kinds reuse the same set. No emoji remains
in nav chrome, buttons, or headers (emoji inside user content/help text is fine).

### 5.3 Toasts replace `alert()`

`showToast(message, kind)` — kinds `success|error|info`; fixed bottom-right
`#toast-stack`; CSS slide/fade-in; auto-dismiss 3.5s; click-to-dismiss; max 4
stacked. Replace **all 19** `alert(` call sites (fleet deploy/export/import,
nodes file read/save, settings save ×3, terminal). After the swap, zero `alert(`
remains in the assets.

### 5.4 Card tiers + header

- `.card` — default surface (unchanged role).
- `.card-hero` — elevated (`var(--popover)` + `--shadow-card` + 2px amber top
  accent) for the primary panel / key numbers of a tab.
- `.card-flat` — transparent background, border-only, for nested groups.
- `showTab` syncs `#header-title` and `document.title` (`"GitMap — <Tab>"`).
- Status indicator becomes a live dot + theme-aware badge (replaces the static
  "Connected: localhost" text).

### 5.5 Responsive

- ≥1200px: full 240px sidebar.
- 900–1200px: icon rail 68px (labels hidden).
- <900px: sidebar off-canvas; hamburger button in the header toggles
  `.sidebar-open`; tables wrapped in `.table-wrap { overflow-x: auto }`;
  settings panes stack single-column.

### 5.6 Demo-data removal

- Delete the phantom `u1` injection (`ui_assets_js_nodes.go:46-50`) and the
  hardcoded `<option value="u1">u1</option>` (`ui_assets_markup_settings_b.go:173`).
  Terminal node select = `local` + `/api/ssh/nodes` results only.
- Delete the fake Antigravity card (`PID: 4757`, `ui_assets_markup_settings_b.go`).
  The instances panel renders from `/api/instances`; empty state reads
  "No instances detected — click Scan Instances."
- Implementer sweeps assets for remaining hardcoded demo literals (`4757`,
  standalone `u1` options).

## 6. Honesty fixes — decisions

### 6.1 Flag parsing — IMPLEMENT

`RunUICmd` (`ui_cmd.go`) gains a cmdui-local `flag.FlagSet`:
`-p`/`--port` (default 8080), `--host` (default `127.0.0.1`), `--no-browser`.
Remaining positional arg = page name. This fixes the `--port`-becomes-page bug.
Help text (`ui_help_menu.go` usage lines) updated to match reality.
Rationale: advertised, user-facing, zero behavior risk.

### 6.2 Help menu pages — ALIGN TO THE 10 REAL TABS

Rewrite `buildUIPagesSection()` to describe exactly the existing tabs:
`settings, commitin, ssh, macro, installer, prompts, import-export, schedules,
editor, help` — one honest line each.
**DECISION: do NOT add `dashboard`/`cluster`/`templates`/`repos` tabs.**
Rationale: those are new product surfaces (commit graphs, template managers),
not polish — inventing tabs to match a wrong help menu is backwards. Honesty
without scope explosion.

### 6.3 SSE — REMOVE THE CLAIM, DO NOT IMPLEMENT

`ui_help_menu.go:59`: "Native Go JSON endpoints with SSE live telemetry" →
"Native Go JSON endpoints (poll-based live refresh)".
**DECISION: do not build SSE.** Rationale: no consumer, no user ask; an SSE
broadcast hub + goroutine lifecycle for zero benefit. Noted as future work.

## 7. Settings tab specifically

- **One sticky save bar** replaces the 7 per-pane Save buttons:
  `.settings-savebar { position: sticky; bottom: 0; }` with a single
  "Save All Settings" button calling the existing `saveSettings()` (it already
  serializes all panes). Rationale: fixes the "7 redundant Save buttons on All
  Settings" complaint; one mental model.
- Keep the 8 sub-tabs, now icon+label pills with an active state.
- Theme select options become "Dark" / "Light" (drop "(Standard)"); it now works (§4.2).
- The "UI/UX Repo Index" catalog pane stays (real curated content, not demo data).

## 8. Security invariant (preserved)

- Default bind stays `127.0.0.1` (`bindAvailablePort` unchanged).
- `--host` may override, but a non-loopback host prints:
  `WARNING: binding a non-loopback address exposes the embedded terminal to the network.`
  (warning only — no interactive prompt, turbo-mode friendly). Document the risk in
  `--help`.

## 9. Acceptance

- Toggling Dark/Light re-themes every tab instantly; choice persists across reload.
- Zero `alert(`, zero emoji in chrome, zero undefined CSS vars, zero phantom nodes.
- `gitmap ui --port 9090 --no-browser` binds 9090 without opening a browser;
  `gitmap ui settings --port 8080` still works (flag no longer becomes the page).
- Help menu lists the 10 real tabs; no SSE claim; `--host` warns on non-loopback.
- Every touched asset file ≤300 lines; no file deletions; no alias changes.
