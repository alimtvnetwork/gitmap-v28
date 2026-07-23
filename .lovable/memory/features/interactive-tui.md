# Memory: features/interactive-tui
Updated: now

The Terminal User Interface (TUI), built with Bubble Tea and Lipgloss, provides a multi-view management experience with 9 primary views: 'Repos' (fuzzy-search), 'Actions' (batch operations), 'Groups' (collection management), 'Status' (live git dashboard), 'Releases' (history and metadata browser), 'Temp Releases' (lightweight branch visualization), 'Zip Groups' (archive management), 'Aliases' (shortcut mapping), and 'Logs' (searchable command history and detail viewer).

## Views

1. **Repos** — Browse tracked repositories with fuzzy search, multi-select (Space), detail panel (Enter).
2. **Actions** — Batch git operations (pull, exec, status, group-add) on selected repos. Supports `--stop-on-fail` to halt after first failure.
3. **Groups** — Manage repository groups: list, create inline, delete with confirmation.
4. **Status** — Live-refreshing dashboard with dirty/clean indicators, ahead/behind counts, stash counts. Configurable auto-refresh (CLI flag, config, or 30s default). Manual refresh via 'r'.
5. **Releases** — Browse release history with version, tag, branch, date columns. Detail view with changelog and notes. Draft/pre-release indicators.
6. **Temp Releases** — Interactive visualization of temp-release branches with three display modes: Flat List, Detail Panel (toggled via 'enter'), and Grouped View (toggled via 'g') which aggregates branches by prefix and calculates sequence ranges.
7. **Zip Groups** — Archive management view.
8. **Aliases** — Shortcut mapping view.
9. **Logs** — Searchable command history with detail viewer. Filter by command, alias, args, or exit code.

## Key Bindings

Global: `q`/`Esc` quit, `Tab` switch views. View-specific keys documented in the TUI help bar.

## Architecture

Core logic split across `tui.go` (root model), `tuiview.go` (rendering), and per-view files (`browser.go`, `actions.go`, `groups.go`, `dashboard.go`, `releases.go`, `tempreleases.go`, `trformat.go`, `zipgroups.go`, `aliases.go`, `logs.go`). All files respect the 200-line limit.
