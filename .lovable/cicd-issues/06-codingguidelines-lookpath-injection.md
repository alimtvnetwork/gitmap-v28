# CI/CD Issue 06: LookPath Injection in Coding Guidelines & Cross-Platform Fallback Assertion

- **Stage**: CI `go test -race` / Test Matrix
- **Status**: ✅ Resolved
- **Affected Packages**: `github.com/alimtvnetwork/gitmap-v28/gitmap/cmd`

## 1. Problem Description

During CI test runs on Linux/Windows runners, `TestRunCodingGuidelinesInstall_ShellMissing` in `cmd/codingguidelines_test.go` either caused flakiness across parallel tests or failed assertions:
1. `TestRunCodingGuidelinesInstall_ShellMissing` used `t.Setenv("PATH", "")` to simulate a missing shell, which wiped the process-wide `$PATH` while concurrent tests (`exec.Command("git", ...)`, `exec.Command("go", ...)`, `fakeTrueRunner`) were running, causing random `executable file not found in $PATH` errors across the package.
2. `assertCGMissingShell` strictly expected the Unix fallback recipe (`curl -fsSL`), failing on Windows where `irm` (PowerShell) is printed.

## 2. Root Cause Analysis

- **Process-Wide `$PATH` Zeroing**: `t.Setenv("PATH", "")` affects the entire OS process. In packages with parallel tests (`t.Parallel()`), zeroing `$PATH` breaks any test spawning child processes or checking binaries via `exec.LookPath`.
- **Platform-Specific Output Assumption**: Error messages and fallback recipes differ by OS (`curl -fsSL` on Unix vs `irm` on Windows).

## 3. Solution

1. **Injectable `LookPath` in `CodingGuidelinesOpts`**:
   - Added `LookPath func(file string) (string, error)` to `CodingGuidelinesOpts` (defaulting to `exec.LookPath`).
   - Wired `LookPath` through `dispatchCGUnix` and `dispatchCGWindows` (via `resolvePowerShellBinaryWithLookPath`).
   - Updated `TestRunCodingGuidelinesInstall_ShellMissing` to inject a stub `LookPath` returning `exec.ErrNotFound` without modifying the global `$PATH`, and marked it `t.Parallel()`.

2. **Cross-Platform Fallback Assertion**:
   - Updated `assertCGMissingShell` to check for either Unix (`curl -fsSL`) or Windows (`irm`) manual fallback hints.

3. **`t.Setenv` in `updatecleanup_handoff_test.go`**:
   - Cleaned up raw `os.Setenv` mutations in `TestBuildCleanupChildEnvForwardsDelayAndJSONPath` to use `t.Setenv`.

## 4. What NOT to Repeat

- **Never zero `$PATH` via `t.Setenv("PATH", "")` in packages with parallel tests**: Always use dependency injection (`LookPath func(...)`) for executable lookup testing.
- **Never hardcode platform-specific strings in cross-platform assertions**: Verify either platform variant (`irm` or `curl -fsSL`) when tests run across multiple operating systems.
