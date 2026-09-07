# Subtask 02: JSON, File, and FNF Export Options for Inspection Preview

## Objective
Implement `--json`, `--file <path>`, `--fnf <path>`, and `--tempfile <filename>` export options for Chrome profile inspection preflight checks.

## Affected Files
- `gitmap/cmd/chromeprofile_transfer_opts.go`
- `gitmap/cmd/chromeprofile_smart_import_export.go` (new)
- `gitmap/cmd/chromeprofile_smart_import.go`
- `gitmap/cmd/chromeprofile_smart_import_preview.go`

## Requirements
1. In `gitmap/cmd/chromeprofile_transfer_opts.go`:
   - Extend `chromeTransferOptions` with `IsJSON`, `FilePath`, `TempFile`, `Fnf`.
   - Update `parseChromeTransferOptions` to safely parse `--file`, `-f`, `-o`, `--fnf`, `--tempfile`, `--json`.
2. In `gitmap/cmd/chromeprofile_smart_import_export.go`:
   - Create `dispatchPreviewOutput` to handle console table, raw JSON, or writing to disk (`--file`, `--fnf`, `--tempfile`).
   - Create parent directories automatically and print confirmation on write.
3. In `gitmap/cmd/chromeprofile_smart_import.go`:
   - Update `resolveCheckTarget` to skip valued flags so `--file <path>` does not pollute the target path.
   - Wire `runChromeProfileImportCheck` to use `dispatchPreviewOutput`.
4. In `gitmap/cmd/chromeprofile_smart_import_preview.go`:
   - Update `renderProfileCandidatesTable` usage hints to include `--json`, `--file`, and `--fnf` examples.
