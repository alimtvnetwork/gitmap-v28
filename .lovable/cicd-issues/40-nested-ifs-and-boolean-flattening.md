# CI/CD Issue 40: Nested If & Boolean Conditional Flattening

- Job: Nested If Linter & Boolean / Enum Linter
- Type: FAIL
- Detected: 2026-09-07T19:59:58Z
- Status: resolved

## Error

```text
❌ FAIL: Found 7 nested-if / anti-compression violation(s) across 6 file(s):
  gitmap/cmd/install_buildessential.go:103: Nested if statement found (depth 2 inside conditional block)
  gitmap/cmd/power_ops.go:93: Nested if statement found (depth 2 inside conditional block)
  gitmap/cmd/power_ops.go:145: Nested if statement found (depth 2 inside conditional block)
  gitmap/cmd/vmware_shared.go:87: Nested if statement found (depth 2 inside conditional block)
  gitmap/cmd/vmware_test.go:38: Nested if statement found (depth 2 inside conditional block)
  gitmap/downloaderconfig/staging.go:37: Nested if statement found (depth 2 inside conditional block)
  gitmap/store/installedtool.go:27: Nested if statement found (depth 2 inside conditional block)
```

## Root Cause

During the implementation of OS power management, VMware tools automation, and the installation split database, nested conditionals (`if` inside `if` or inside `else`) were introduced, exceeding the strict depth-1 conditional complexity limit enforced by `check-nested-ifs.py` and `check-enum-and-boolean.py`.

## Fix Applied

1. **`gitmap/cmd/install_buildessential.go`**: Flattened `probeToolVersion` by executing output parsing and checking error status with a single compound conditional.
2. **`gitmap/cmd/power_ops.go`**: Extracted `parseFirstPositional(fs)` and `loadPreviousOrDefault(target)` helpers to eliminate nested conditionals in CLI flag and reset handlers.
3. **`gitmap/cmd/vmware_shared.go`**: Refactored `ensureMountDirectory` using an early-return guard clause (`if !os.IsNotExist(err)`).
4. **`gitmap/cmd/vmware_test.go`**: Refactored `TestVmwareOSConstraint` with an early return on Linux.
5. **`gitmap/downloaderconfig/staging.go`**: Refactored `ResolvePersistentKeepDir` with an early return for `SUDO_USER`.
6. **`gitmap/store/installedtool.go`**: Extracted `migrateIfConnected` helper with guard clause in `openSplitOrFallback`.
7. **Verification**: Executed `check-nested-ifs.py` across 2,487 files (0 violations) and `check-enum-and-boolean.py` across 1,869 files (0 violations), followed by the complete 33-gate local CI runner (`06-cicd-local-runner.py`) which passed 100% with exit code 0.
