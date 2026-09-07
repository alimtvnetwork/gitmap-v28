# Subtask 01: Chrome Profile Polymorphic Discovery & Preflight `ls`

## Scope
- Update `cmd/chromeprofile_smart_import.go` and add `cmd/chromeprofile_smart_import_discovery.go`:
  - Discard/skip `manifest.json` from being treated as a single profile.
  - Implement recursive profile discovery across subdirectories (`Default/`, `Profile 1/`, `Profile 2/`, etc.).
  - Recognize profile folders containing `<folder>.json`, `Preferences`, `Bookmarks`, or `Extensions/`.
  - Add preflight inspection handlers for:
    - `gitmap chrome import ls [path]`
    - `gitmap chrome import-all ls [path]`
    - `gitmap chrome import-check [path]`
  - Render clear tabular inspection view (Profile Directory Name, Display Name, Email, Bookmarks Count, Extensions Count, Source Path) with `--json` support.

## Files Touched
- `gitmap/cmd/chromeprofile_smart_import.go`
- `gitmap/cmd/chromeprofile_smart_import_discovery.go`
- `gitmap/cmd/chromeprofile_import_handlers.go`
- `gitmap/cmd/chrome.go`
