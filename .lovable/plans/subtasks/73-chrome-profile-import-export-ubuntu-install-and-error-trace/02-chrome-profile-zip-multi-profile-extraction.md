# Subtask 02: Chrome Profile ZIP Multi-Profile Extraction & Batch Import

## Scope
- Update `cmd/chromeprofile_zip_import.go`:
  - Prevent ZIP archives from being treated as a single generic profile named after the archive file.
  - Safely extract `.zip` archive into temporary staging sandbox (`/tmp/gitmap-stage-chrome-*` or `os.TempDir()`).
  - Read `manifest.json` if present in root or probe top-level directories for profiles.
  - For each discovered profile in the archive:
    - Extract preferences, bookmarks, extensions.
    - Dispatch to `importChromeJSON` / `applyChromeExport`.
  - Clean up sandbox upon completion.

## Files Touched
- `gitmap/cmd/chromeprofile_zip_import.go`
- `gitmap/cmd/chromeprofile_smart_import.go`
