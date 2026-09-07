# Subtask 03: Windows Coding Guidelines Installer & Smoke Workflow

## Scope
- Update `gitmap/cmd/codingguidelines.go`:
  - Implement `writeCGCompatScriptWindows(url)` to download `error-manage-install.ps1`, patch invalid variable references (`\$([A-Za-z0-9_]+):` -> `${$1}:`), and invoke with `pwsh`/`powershell`.
  - Ensure temp files are cleaned up via defer/cleanup function.
- Update `.github/workflows/goreleaser-smoke.yml`:
  - In `Invoke-CfrCg`, pipe `Tee-Object` output to `Out-Null` so only `$LASTEXITCODE` is returned as the function result.

## Files Touched
- `gitmap/cmd/codingguidelines.go`
- `gitmap/cmd/codingguidelines_test.go`
- `.github/workflows/goreleaser-smoke.yml`
