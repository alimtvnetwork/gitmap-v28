# Dashboard Recent Flag and Terminal UI Enhancements

Slug: dashboard-recent-and-terminal-ui
Status: pending
Created: 2026-08-20

## Intent

Improve the visualization of the `gitmap dashboard` and `gitmap hd` terminal outputs by introducing colors and better layout formatting. Introduce a `--recent` flag to the dashboard command to specifically filter and highlight recent items. Update the HTML UI app (the dashboard output) to include a "Recent Items" view to prominently display these recent updates. Finally, ensure all corresponding help texts are thoroughly documented to cover these additions.

## Scope

- `gitmap/constants/constants_dashboard.go` and `constants_terminal.go` for new flags and color definitions.
- `gitmap/cmd/dashboard.go` to parse the `--recent` flag and use terminal coloring.
- `gitmap/dashboard/collector.go` to handle recent items logic.
- `gitmap/dashboard/templates/dashboard.html` to add the "Recent Items" UI component.
- `gitmap/helptext/dashboard.md` and `gitmap/helptext/hd.md` to document the changes.

## Acceptance Criteria

1. The `--recent` flag is successfully added to `gitmap dashboard`.
2. Terminal output for `dashboard` and `hd` commands shows improved visualization and coloring using standard ANSI codes.
3. The generated HTML dashboard includes a "Recent Items" section that defaults to active when `--recent` is passed, or is available as a view tab.
4. Help texts (`dashboard.md`, `hd.md`) are completely updated.
5. All codebase changes adhere to `spec/12-consolidated-guidelines/02-go-code-style.md`.

## Affected Files

- `gitmap/constants/constants_dashboard.go`
- `gitmap/constants/constants_terminal.go`
- `gitmap/cmd/dashboard.go`
- `gitmap/dashboard/collector.go`
- `gitmap/dashboard/templates/dashboard.html`
- `gitmap/helptext/dashboard.md`
- `gitmap/helptext/hd.md`

## Linked Issues/Commands

- Captured user instruction to improve terminal text, visualization, add `--recent` flag, and update the HTML UI App.
