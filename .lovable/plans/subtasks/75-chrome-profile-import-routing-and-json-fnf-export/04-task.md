# Subtask 04: End-to-End Tests & CI Verification

## Objective
Author and execute comprehensive unit and end-to-end tests validating profile subcommand routing, JSON export, file writing, and preflight inspect behavior.

## Affected Files
- `gitmap/cmd/profile_route_test.go` (new)
- `gitmap/cmd/chromeprofile_export_flags_test.go` (new)

## Requirements
1. Test `gitmap profile import`, `import-all`, `export`, `export-all`, `inspect`, `preview`, `check` routes correctly.
2. Test `import-check` with `--json`, `--file <path>`, `--fnf <path>`, and `--tempfile <filename>`.
3. Verify files are created with valid JSON on disk.
4. Run all local tests and CI quality gates.
