# Dashboard Recent Flag and Terminal UI Enhancements

Slug: dashboard-recent-and-terminal-ui
Steps: 30
Status: pending
Created: 2026-08-20

## Context
The user requested better terminal visualization with more coloring for dashboard commands, a new `--recent` flag to filter recent items, and an update to the interactive HTML UI App (`dashboard.html`) to prominently display these recent items. All new flags and features must be documented in the corresponding help texts. Affected files span CLI parsing, data collection, HTML templates, and Markdown documentation.
Links:
- Spec: [.lovable/spec/tasks/08-dashboard-recent-and-terminal-ui.md](file:///d:/work/gitmap/.lovable/spec/tasks/08-dashboard-recent-and-terminal-ui.md)
- Adheres to: `.lovable/coding-guidelines.md`

## Steps
1. Group 1: Constants and Configuration Update (Execute via standalone agent "Config Updater"). Read `gitmap/constants/constants_dashboard.go`.
2. In `gitmap/constants/constants_dashboard.go`, define `FlagRecent` string constant with the value `"recent"`.
3. In `gitmap/constants/constants_dashboard.go`, define `FlagDescDashRecent` string constant with the value `"Show and highlight recent items only (e.g., last 7 days)"`.
4. Read `gitmap/constants/constants_terminal.go` to identify available ANSI color constants.
5. In `gitmap/constants/constants_dashboard.go`, define ANSI-colored message variants for `MsgDashCollecting`, `MsgDashWriteJSON`, `MsgDashWriteHTML`, and `MsgDashGenerated`.
6. Commit Group 1: Verify absolutely NO artifacts are staged. Ensure `.gitignore` explicitly excludes them. Fix any git issues and push the commit to the repository immediately.
7. Group 2: CLI Parsing and Terminal UI (Execute via standalone agent "Terminal UI Enhancer"). Read `gitmap/cmd/dashboard.go` to analyze `parseDashboardFlags` and `runDashboard`.
8. In `parseDashboardFlags` within `gitmap/cmd/dashboard.go`, add a boolean flag definition for `--recent` using `FlagRecent` and `FlagDescDashRecent`.
9. Modify `parseDashboardFlags` to return the new `recent` boolean alongside existing options.
10. Update the `dashboard.CollectOptions` struct initialization in `parseDashboardFlags` to include the `Recent` field set to the parsed boolean value.
11. Update `runDashboard` to replace standard `fmt.Println` and `fmt.Printf` calls with the newly defined ANSI-colored message constants from Step 5.
12. Commit Group 2: Verify absolutely NO artifacts are staged. Ensure `.gitignore` explicitly excludes them. Fix any git issues and push the commit to the repository immediately.
13. Group 3: Data Collection Logic (Execute via standalone agent "Data Collector"). Read `gitmap/dashboard/collector.go` to understand how commits and data are aggregated.
14. Ensure `dashboard.CollectOptions` struct definition in `gitmap/dashboard/collector.go` or related models file has a `Recent bool` field added.
15. In `gitmap/dashboard/collector.go`, implement logic to automatically set a dynamic `Since` date boundary (e.g., 7 days ago) if `opts.Recent` is true and `opts.Since` is empty.
16. Ensure the exported JSON model (`dashboard.WriteJSON`) passes down the `Recent` boolean flag in its metadata payload so the frontend knows if the view is recent-focused.
17. Commit Group 3: Verify absolutely NO artifacts are staged. Ensure `.gitignore` explicitly excludes them. Fix any git issues and push the commit to the repository immediately.
18. Group 4: HTML UI App Update (Execute via standalone agent "Frontend Developer"). Read `gitmap/dashboard/templates/dashboard.html` to analyze the HTML layout and JS logic.
19. In `gitmap/dashboard/templates/dashboard.html`, add a new navigation tab button labeled `Recent Items` in the main navigation bar.
20. In `gitmap/dashboard/templates/dashboard.html`, add a `<div id="recentItemsView" class="tab-content">` section to hold the recent items table and visualization.
21. In `gitmap/dashboard/templates/dashboard.html`, update the JavaScript logic to parse the `Recent` boolean flag from the JSON payload.
22. In `gitmap/dashboard/templates/dashboard.html`, write a `renderRecentItems()` JavaScript function that populates the `recentItemsView` with a summarized list of recent commits.
23. In `gitmap/dashboard/templates/dashboard.html`, modify the `init()` JavaScript function to automatically switch to the `recentItemsView` tab if the payload's `Recent` flag is true.
24. Commit Group 4: Verify absolutely NO artifacts are staged. Ensure `.gitignore` explicitly excludes them. Fix any git issues and push the commit to the repository immediately.
25. Group 5: Help Text Documentation (Execute via standalone agent "Documentation Writer"). Read `gitmap/helptext/dashboard.md` and `gitmap/helptext/hd.md`.
26. In `gitmap/helptext/dashboard.md`, document the new `--recent` flag under the Flags section and add a usage example.
27. In `gitmap/helptext/hd.md`, ensure any relevant references to recent visualization improvements or synergies are documented.
28. Commit Group 5: Verify absolutely NO artifacts are staged. Ensure `.gitignore` explicitly excludes them. Fix any git issues and push the commit to the repository immediately.
29. Group 6: Verification and Release (Execute via standalone agent "Release Manager"). Run `go test ./gitmap/cmd/... ./gitmap/dashboard/...` to verify the new flag parsing and collection logic does not break existing tests.
30. Execute the release ceremony per `11-release.md` (bump MINOR version, update changelog with the whole plan summary, pin new version in root README), and immediately commit and push the final release.

## Verification
- Run `gitmap dashboard --recent` and ensure terminal output utilizes colors.
- Open the resulting HTML file and verify the "Recent Items" tab is active and populated.
- Review `gitmap help dashboard` to verify `--recent` is documented.

## Appended from prior pending tasks
none
