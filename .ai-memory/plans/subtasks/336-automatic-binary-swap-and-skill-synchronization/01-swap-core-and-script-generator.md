# Subtask 336.01 — Swap Core Engine and Platform Script Generator

> **Sequence ID:** 336  
> **Parent Plan:** [336-automatic-binary-swap-and-skill-synchronization.md](file:///d:/work/gitmap/.ai-memory/plans/pending/336-automatic-binary-swap-and-skill-synchronization.md)  
> **Target Packages:** `cli/cmdswap/`, `cli/cmdupdate/`  
> **Status:** Pending  

---

## 1. Objectives

1. Create a dedicated package `cli/cmdswap` providing data structures and core services for process discovery, script generation, and atomic replacement.
2. Implement data structures with positive boolean properties:
   - `SwapOptions`: `FromPath`, `ToPath`, `TargetPID`, `hasKill`, `hasRestart`, `hasForce`, `isScriptOnly`, `isJSONOutput`, `maxWaitSeconds`.
   - `ScriptPlan`: `Platform`, `ScriptPath`, `ScriptContent`, `TargetPIDList`, `isDetached`, `ExecutionArgs`, `EstimatedWait`.
   - `SwapResult`: `isSuccess`, `FromPath`, `ToPath`, `TerminatedPIDs`, `ScriptPlan`, `ErrorMessage`, `OutputMessage`, `DurationMs`.
   - `ProcessLockInfo`: `PID`, `ProcessName`, `Executable`, `isParent`.
3. Implement candidate source binary resolution (`./cli/gitmap.exe`, `./bin/gitmap.exe`, `./gitmap.exe`) and target resolution via `cmdinstall.ResolveDeployedAndConfigPaths()`.
4. Implement PID discovery and locking process detection across Windows (`tasklist`, WMI process snapshot) and Unix (`pgrep`, `/proc`).
5. Implement detached script generators:
   - Windows PowerShell (`swap_handoff.ps1`): Wait for parent PID, graceful stop, forceful termination fallback, `Move-Item -Force` with exponential backoff retry loop, post-swap verification (`gitmap --version`), and self-cleanup.
   - Unix POSIX Shell (`swap_handoff.sh`): Wait for parent PID, `SIGTERM`/`kill -9` fallback, atomic `mv`, verification, and self-cleanup.
6. Provide unit tests covering option resolution, script template interpolation, and lock detection guards.

---

## 2. File Manifest

- `cli/cmdswap/types.go` — Data contracts and option structs with positive boolean naming.
- `cli/cmdswap/resolver.go` — Source and target path discovery and binary checksum verification.
- `cli/cmdswap/discovery.go` — Sibling and locking process enumeration (`ProcessLockInfo`).
- `cli/cmdswap/generator_windows.go` — Detached PowerShell script template renderer.
- `cli/cmdswap/generator_unix.go` — Detached POSIX shell script template renderer.
- `cli/cmdswap/executor.go` — Detached process invocation and clean parent exit handler.
- `cli/cmdswap/swap_test.go` — Unit tests for options validation and script generation.

---

## 3. Verification Criteria

- [ ] `go test -C cli ./cmdswap` passes 100% with zero race conditions.
- [ ] Detached script contains parameterized retries, error traps, and post-swap verification checks.
- [ ] Options with `isScriptOnly = true` return `ScriptPlan` content without killing processes or modifying binary files.
