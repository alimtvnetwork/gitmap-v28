# GitMap v6.573.0

## What's Changed in v6.573.0

### CI/CD Pipeline Defect Resolutions

1. **Dead Code Elimination (`cli/cmddb`)**:
   - Removed uncalled legacy `cmddb_help.go` functions (`runDbHelp`, `printDBHelpHeader`, `printDBHelpSubcommands`, `printDBHelpFlags`, `printDBHelpExamples`), eliminating golangci-lint `[unused]` linter failures. Database help is centrally dispatched through `cmddb_dispatch.go`.

2. **Static Analysis Nil Pointer Dereference Guard (`cli/cmdupdate`)**:
   - Added immediate nil candidate guard in `finalizeUpdate` (`cli/cmdupdate/update_binary_swap.go`), resolving staticcheck `SA5011` warning where candidate pointers were dereferenced prior to nil checks.

3. **Async Poller Interval Clamping Test Alignment (`cli/cmd`)**:
   - Aligned `TestParseAsyncArgs` in `cli/cmd/async_cmd_test.go` with the system-wide `>= 30s` poller floor invariant, resolving race-detector test failures.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.573.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.573.0/install.sh | sh
```
