# Subtask 05 — Theme System Implementation (program 246)

## Objective

Implement spec `02-spec/21-app/246-parallel-fix-command/04-ui-ux-overhaul.md` §4:
a REAL working theme system — refined dark theme + a full light theme — with a
functional switcher persisted to settings, and the undefined `--accent` fixed.
Keep the Go-string-constant asset model (no npm, no build step).

## Read first

- `02-spec/21-app/246-parallel-fix-command/04-ui-ux-overhaul.md` (§2–§4, token tables)
- `.ai-memory/what-to-read.md` — "The Mindset" (§5: ~300 lines max per file)

## Owned files (ONLY these — do not touch any other file)

New:

- `cli/cmdui/ui_assets_theme.go` (≤300 lines) — `uiAssetsStyle` const moved here,
  extended with light-theme overrides, toast CSS, and the `icon()` SVG helper.

Edit:

- `cli/cmdui/ui_assets.go` — concat order: swap the style const reference to the
  new file at the same position; no other order change.
- `cli/cmdui/ui_assets_markup_misc.go` — REMOVE the `uiAssetsStyle` const
  (moved, not deleted); add responsive shell bits + hamburger per spec §5.5.
- `cli/cmdui/ui_assets_js_nodes.go` — add `applyTheme(name)`; wire it into
  `loadSettings()` and the `#setting-theme` change handler (optimistic apply +
  POST, revert-on-failure with error toast).
- `cli/cmdui/ui_assets_markup_ops.go` — theme select options become
  "Dark" / "Light" (drop "(Standard)").

## Implementation checklist

- [ ] `uiAssetsStyle` moved byte-identical into `ui_assets_theme.go`, then extended:
      `:root` = refined dark tokens, `:root[data-theme="light"]` = full light set
      (exact values in spec §4.1).
- [ ] `--accent` defined in BOTH themes (fixes the undefined var at
      `ui_assets_markup_ops.go:68`).
- [ ] New tokens added: `--warning`, `--info`, `--radius-sm/md/lg`,
      `--font-sans`, `--font-mono`, `--shadow-card`, `--sp-1…--sp-8`,
      `--fs-xs…--fs-2xl` (spec §5.1).
- [ ] Terminal/editor surfaces (`--code-bg`, `#editor-area`) stay dark in both
      themes — deliberate, documented in a CSS comment.
- [ ] Hardcoded-hex sweep → vars across ALL asset files: `#131d2e`→`var(--raised)`,
      `#1e293b`→`var(--popover)`, `#f8fafc`/`#e2e8f0`→`var(--text)`,
      `#cbd5e1`→`var(--muted)`. No raw hex left outside the token definitions.
- [ ] `applyTheme`: sets `document.documentElement.dataset.theme`, syncs the select.
      Called on settings load AND on select change (optimistic apply, POST
      `/api/settings`, revert + error toast on failure).
- [ ] Server: NO CHANGE — `theme` already persists via `SettingsData`
      (`ui_server.go`); default stays `"dark"`.
- [ ] Every file ≤300 lines (`wc -l` each). No file deletions. No alias changes.

## Per-tab verification checklist (both themes, after implementation)

- [ ] settings (all 8 sub-tabs) — forms legible, no invisible text
- [ ] commitin, ssh, macro, installer, prompts, import-export, schedules, editor, help
- [ ] Theme choice survives a full page reload (persisted to `~/.gitmap/ui_settings.json`)
- [ ] Contrast spot-check: body text and muted text meet spec §2 ratios

## Done criteria

- [ ] Dark/Light toggle re-themes every tab instantly; persists across reload.
- [ ] Zero undefined CSS vars; zero hardcoded hex outside token blocks.
- [ ] Report lists each touched file with its line count (all ≤300).

## Hard rules

- Search ONLY via `gitmap aum search` / `gitmap find` / `gitmap cat`. TOTAL BAN on
  grep, rg, git grep, Select-String.
- NEVER run git commands. NEVER run `go build` / `go test` (lead verifies centrally).
- Relative paths only. NEVER propose deleting files.
- Security invariant: do not touch `bindAvailablePort` (127.0.0.1 default).
