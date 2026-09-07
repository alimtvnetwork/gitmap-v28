# Subtask 03: Fix macOS Power Driver and Cross-Platform Driver Interaction Test

## Scope
- In `gitmap/power/driver_darwin.go`:
  - Replace hardcoded `IsNeverSleep: false` in `GetStatus()` with real `pmset -g` parsing.
  - Parse `displaysleep` and `sleep` timeout minutes from `pmset -g` output.
  - If both `displaysleep == 0` and `sleep == 0`, set `IsNeverSleep = true`.
- In `gitmap/power/manager_test.go`:
  - Update `TestMockRunner_DriverInteraction` to supply mock output matching the platform being tested:
    - Windows: `powercfg` output
    - Darwin: `pmset -g` output (`displaysleep 0`, `sleep 0`)
    - Linux: `gsettings` output
- Keep all functions <= 15 lines.
- Blank line before every return.
- Zero nested ifs.

## Files Touched
- `gitmap/power/driver_darwin.go`
- `gitmap/power/manager_test.go`
- `gitmap/power/parser_pmset.go`
