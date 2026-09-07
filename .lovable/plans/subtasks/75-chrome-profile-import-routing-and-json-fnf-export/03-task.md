# Subtask 03: Helptext Documentation & Golden Tests

## Objective
Update helptext markdown documentation across relevant profile and import commands, adding examples for JSON export, `--file`, and `--fnf`, and ensuring all helptext golden tests pass.

## Affected Files
- `gitmap/helptext/import-check.md`
- `gitmap/helptext/import-all.md`
- `gitmap/helptext/profile.md`
- `gitmap/helptext/chrome-profile-import.md`

## Requirements
1. Document `--json`, `--file <path>`, `--fnf <path>`, and `--tempfile <filename>`.
2. Add explicit examples showing both `--json` and `--json --file <path>`.
3. Use fenced code blocks (` ``` `), strictly no 4-space indented blocks.
4. Pass `go test ./gitmap/helptext/... -run Golden -count=1`.
