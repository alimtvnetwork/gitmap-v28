# Subtask 03: Summary & Profile Import Sign-in Integration

## Objective
Wire the responsive layout into `gitmap/cmd/pull_table_summary.go` and verify profile import sign-in scrubbing integration in `chromeprofile_preferences.go` and `chromeprofile_smart_import_archive.go`.

## Target Files
- `gitmap/cmd/pull_table_summary.go`
- `gitmap/cmd/chromeprofile_preferences.go`
- `gitmap/cmd/chromeprofile_smart_import_archive.go`
- `gitmap/cmd/chromeprofile_export.go`

## Implementation Details
1. In `pull_table_summary.go`:
   - `RenderPullBatchTable(rows []model.PullTableRow)`:
     - Instantiate `layout := NewPullTableLayout(rows)` (which automatically detects terminal width).
     - Print header, render each row via `layout.PrintRow(r)`.
     - Print dynamic bottom divider matching `layout.DividerLen`.
2. In Chrome Profile files:
   - Verify `patchImportedChromeProfilePreferencesWithOptions` cleanly scrubs `account_info`, `sync`, `google`, `gaia_cookie`, and sets `signin.allowed = false`, `browser.has_seen_welcome_page = true`.
   - Verify `copyProfileDirectoryDiskFiles` preserves `Cookies` and `Network/Cookies` with automatic parent directory creation.
