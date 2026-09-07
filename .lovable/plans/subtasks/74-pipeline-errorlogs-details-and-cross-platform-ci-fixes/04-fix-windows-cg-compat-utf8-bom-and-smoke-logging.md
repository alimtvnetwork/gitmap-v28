# Subtask 04: Fix Windows CG Compat Script UTF-8 BOM & GoReleaser Smoke Failure Logging

## Scope
- In `gitmap/cmd/codingguidelines_compat.go`:
  - In `patchCGWindowsScriptFile`, ensure the file is written with a UTF-8 BOM (`\xef\xbb\xbf`) so that Windows PowerShell 5.1 parses unicode characters (such as em-dash `—`) as UTF-8 rather than Windows-1252 ANSI, preventing `Unexpected token 'local' in expression or statement`.
- In `.github/workflows/goreleaser-smoke.yml`:
  - When `cfr cg` fails (`if ($rc -ne 0)`), dump `$logText` (`Get-Content $log`) before exiting so that failure reasons are visible in CI output and retrievable by `gitmap pipeline errorlogs`.
- Keep all functions <= 15 lines.
- Blank line before every return.
- Zero nested ifs.

## Files Touched
- `gitmap/cmd/codingguidelines_compat.go`
- `.github/workflows/goreleaser-smoke.yml`
