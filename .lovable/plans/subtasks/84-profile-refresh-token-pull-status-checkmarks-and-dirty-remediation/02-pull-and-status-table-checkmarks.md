# Subtask 02: Pull and Status Table Checkmarks

## Objective
Update the `gitmap pull` and status table displays to show clean check marks (`✔ up-to-date` or `✔ active`) instead of raw uppercase strings or plain `active`.

## Files to Touch
- `gitmap/cmd/pull_table_style.go`
- `gitmap/cmd/pull_table_row.go`
- `gitmap/cmd/pull_table_layout.go`
- `gitmap/cmd/pull_table_test.go`

## Detailed Implementation Steps
1. In `gitmap/cmd/pull_table_style.go`:
   - Remove `//nolint:unused`.
   - Update `formatPullStatus(status string, isDirty bool) string`:
     - If `isDirty`: return `constants.ColorYellow + "● dirty" + constants.ColorReset`
     - Case `"UP_TO_DATE"`, `"up-to-date"`, `"synced"`: return `constants.ColorGreen + "✔ up-to-date" + constants.ColorReset`
     - Case `"active"`, `"ACTIVE"`: return `constants.ColorGreen + "✔ active" + constants.ColorReset`
     - Case `"SUCCESS"`, `"ok"`, `"updated"`: return `constants.ColorGreen + "✔ updated" + constants.ColorReset`
     - Case `"FAILED"`, `"fail"`, `"error"`: return `constants.ColorRed + "✖ failed" + constants.ColorReset`
2. In `gitmap/cmd/pull_table_row.go`:
   - In `printWideRow` and `printCompactRow`, replace `renderedStatus := statusStyle.Render(r.PullStatus)` with:
     ```go
     formattedStatus := formatPullStatus(r.PullStatus, r.IsDirty)
     padStatus := calcAnsiPadding(formattedStatus, l.MaxStatus)
     ```
   - Ensure the padding and column alignment correctly handles UTF-8 runes (using `calcAnsiPadding`).
3. In `gitmap/cmd/pull_table_test.go`:
   - Add/update unit test verifying `formatPullStatus` renders `✔` for `UP_TO_DATE`, `synced`, and `active`.
