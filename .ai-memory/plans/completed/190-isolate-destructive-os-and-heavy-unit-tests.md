# Plan 190: Isolate Destructive OS & Heavy Unit Tests — Coding Guideline 24

## Executive Summary
This plan establishes strict repository-wide compliance with **Coding Guideline 24** (`02-spec/02-coding-guidelines/24-isolate-destructive-os-and-heavy-unit-tests.md`), ensuring that unit tests NEVER trigger real OS shutdown, reboot, power-off, system package upgrades, user modifications, or unmocked filesystem purges. All power actions, system cleanups, and heavy alterations are decoupled behind injectable executors (`OSActionExecutor`, `FileRemover`, `OSCommandRunner`, `TempDirResolver`) with hermetic mock testing and fast duration verification (1s/2s checks).

---

## Violation Ledger

| Violation ID | Component / File | Severity | Nature of Violation | Remediation Architecture |
|---|---|---|---|---|
| **V-01** | `cli/cmdschedule/schedule_os.go:190` | **HIGH** | `CancelSchedulePowerCLI` invokes raw `exec.Command("shutdown", ...)` directly without going through `DefaultOSActionExecutor`. | Introduce `OSActionCancel` into `OSActionType` and route cancel execution through `DefaultOSActionExecutor(OSActionParams{Action: OSActionCancel, ...})`. |
| **V-02** | `cli/cluster/exec_lifecycle.go:42, 65` | **HIGH** | `ExecShutdown` and `ExecLogoff` construct native `exec.CommandContext("shutdown", ...)` and execute via `runCmdFunc` with no dedicated unit test coverage. | Add hermetic mock tests in `cli/cluster/exec_lifecycle_test.go` covering `ExecShutdown`, `ExecRestart`, and `ExecLogoff` with `defer` restoration blocks. |
| **V-03** | `cli/cmdos/os_ai_clean.go:226` | **MEDIUM** | `removeSingleFileSafely` executes direct `os.Remove(filePath)` on real system paths discovered during AI cache sweeps. | Decouple removal behind injectable `FileRemover` (`var defaultFileRemover FileRemover = removeSingleFileSafely`) and assert captured mock paths in tests. |
| **V-04** | `cli/osclean/clean.go:13` | **MEDIUM** | `CleanTempDirectories` resolves host system temp roots (`os.TempDir()`, `TEMP`, `/tmp`) with direct `os.RemoveAll` on non-dry-run without mock hook. | Introduce injectable `TempDirResolver` (`var defaultTempDirResolver TempDirResolver = resolveTempDirectories`) and add hermetic unit tests in `clean_test.go` targeting `t.TempDir()`. |
| **V-05** | `cli/cmdos/os_update.go`, `os_full_upgrade.go` | **MEDIUM** | `ExecuteOSUpdate` and `ExecuteOSFullUpgrade` execute direct `cmd.Run()` (`Get-WindowsUpdate`, `apt-get upgrade`, `do-release-upgrade`). | Introduce injectable `OSCommandRunner` and unit tests in `os_update_test.go` asserting command lines without spawning native updaters. |
| **V-06** | `cli/cmdservice/driver.go`, `service_cmd.go` | **HIGH** | `runService*` invoked native `ResolveServiceDriver()` directly, interacting with real `sc.exe`, `systemctl`, `launchctl`. | Introduce injectable `DefaultServiceDriverResolver ServiceDriverResolver` and comprehensive mock-driven unit tests in `cli/cmdservice/service_cmd_test.go`. |
| **V-07** | `cli/osuser/create_root.go`, `remove.go`, `kill.go`, `sudoers.go` | **HIGH** | `CreateRootUser`, `RemoveEnhancedUser`, and `KillUserProcesses` executed unmocked `net user`, `useradd`, `deluser`, `pkill`, `taskkill`. | Introduce injectable `defaultOSCommandRunner OSCommandRunner` and hermetic mock unit tests in `cli/osuser/osuser_test.go`. |
| **V-08** | `cli/cmdos/os_ip.go` | **MEDIUM** | `runOSIP` created network IP manager directly without factory hook for mock drivers. | Introduce injectable `defaultNetIPManagerFactory NetIPManagerFactory` and hermetic unit tests in `cli/cmdos/os_ip_test.go`. |

---

## Architecture Specification

### 1. Power Action Injectable Executor Decoupling (`cli/cmdschedule/schedule_os.go`)
```go
const (
    OSActionShutdown OSActionType = "shutdown"
    OSActionRestart  OSActionType = "restart"
    OSActionCancel   OSActionType = "cancel"
)

type OSActionParams struct {
    Action     OSActionType
    Delay      time.Duration
    Seconds    int64
    Executable string
    Args       []string
    OnTick     func(remaining time.Duration)
}

type OSActionExecutor func(params OSActionParams) error
var DefaultOSActionExecutor OSActionExecutor = executeNativeOSAction
```
In `CancelSchedulePowerCLI`:
```go
func CancelSchedulePowerCLI(action OSActionType) error {
    _ = CancelActivePowerSchedule()
    exe, cmdArgs := BuildCancelOSActionCommand()
    params := OSActionParams{
        Action:     OSActionCancel,
        Executable: exe,
        Args:       cmdArgs,
    }
    if err := DefaultOSActionExecutor(params); err != nil {
        return err
    }
    fmt.Printf("✓ Canceled active scheduled %s.\n", action)
    return nil
}
```

### 2. Cluster Lifecycle Hermetic Testing (`cli/cluster/exec_lifecycle_test.go`)
All lifecycle commands (`ExecRestart`, `ExecShutdown`, `ExecLogoff`) MUST be tested using:
```go
origRunCmd := runCmdFunc
defer func() { runCmdFunc = origRunCmd }()

var capturedCmd *exec.Cmd
runCmdFunc = func(cmd *exec.Cmd) error {
    capturedCmd = cmd
    return nil
}
```
Assert that `capturedCmd.Path` and `capturedCmd.Args` contain the correct OS-specific flags (`/s`, `/r`, `-h`, `logoff`) without executing native binaries.

### 3. File Removal Mocking (`cli/cmdos/os_ai_clean.go`)
```go
type FileRemover func(filePath string) (int, int64)
var defaultFileRemover FileRemover = removeSingleFileSafely
```
Tests in `os_ai_clean_test.go` swap `defaultFileRemover` to record paths and simulated byte sizes.

### 4. Temp Directory Sweep Mocking (`cli/osclean/clean.go`)
```go
type TempDirResolver func() []string
var defaultTempDirResolver TempDirResolver = resolveTempDirectories
```
In `clean_test.go`:
```go
orig := defaultTempDirResolver
defer func() { defaultTempDirResolver = orig }()
dir := t.TempDir()
defaultTempDirResolver = func() []string { return []string{dir} }
```

### 5. OS Command Runner Injection (`cli/cmdos/os_update.go`, `os_full_upgrade.go`)
```go
type OSCommandRunner func(cmd *exec.Cmd) error
var defaultOSCommandRunner OSCommandRunner = func(cmd *exec.Cmd) error { return cmd.Run() }
```

---

## Subtask Decomposition

- **Subtask 01**: `.ai-memory/plans/subtasks/190-isolate-destructive-os-and-heavy-unit-tests/01-power-cancel-executor-and-schedule-tests.md`
  - Refactor `CancelSchedulePowerCLI` to route through `DefaultOSActionExecutor`.
  - Add `OSActionCancel` constant.
  - Add `TestCancelSchedulePowerCLI_Mocked` in `schedule_os_test.go`.
  - Add 1s/2s duration mock tests.
- **Subtask 02**: `.ai-memory/plans/subtasks/190-isolate-destructive-os-and-heavy-unit-tests/02-cluster-lifecycle-shutdown-mock-tests.md`
  - Add tests for `ExecShutdown` and `ExecLogoff` in `cli/cluster/exec_lifecycle_test.go`.
  - Ensure 100% of cluster power commands are hermetically mocked with `defer` restoration.
- **Subtask 03**: `.ai-memory/plans/subtasks/190-isolate-destructive-os-and-heavy-unit-tests/03-os-ai-clean-and-osclean-injectable-removers.md`
  - Refactor `os_ai_clean.go` to introduce `defaultFileRemover`.
  - Refactor `osclean/clean.go` to introduce `defaultTempDirResolver`.
  - Create unit tests in `cli/osclean/clean_test.go` and `cli/cmdos/os_ai_clean_test.go`.
- **Subtask 04**: `.ai-memory/plans/subtasks/190-isolate-destructive-os-and-heavy-unit-tests/04-os-update-full-upgrade-runner-injection.md`
  - Refactor `os_update.go` and `os_full_upgrade.go` to use `defaultOSCommandRunner`.
  - Create unit tests in `cli/cmdos/os_update_test.go` asserting platform command construction.
- **Subtask 05**: `.ai-memory/plans/subtasks/190-isolate-destructive-os-and-heavy-unit-tests/05-service-driver-resolution-decoupling.md`
  - Decouple `DefaultServiceDriverResolver` in `cli/cmdservice/driver.go` and `service_cmd.go`.
  - Decompose functions in `service_cmd.go` to stay <= 8-15 lines.
  - Create `cli/cmdservice/service_cmd_test.go` with mock `ServiceDriver` covering all service subcommands.
- **Subtask 06**: `.ai-memory/plans/subtasks/190-isolate-destructive-os-and-heavy-unit-tests/06-osuser-command-runner-decoupling.md`
  - Introduce `defaultOSCommandRunner OSCommandRunner` in `cli/osuser/runner.go`.
  - Wire `defaultOSCommandRunner` across `create_root.go`, `remove.go`, `kill.go`, and `sudoers.go`.
  - Add hermetic mock unit tests in `cli/osuser/osuser_test.go`.
- **Subtask 07**: `.ai-memory/plans/subtasks/190-isolate-destructive-os-and-heavy-unit-tests/07-netip-manager-factory-decoupling.md`
  - Introduce `defaultNetIPManagerFactory NetIPManagerFactory` in `cli/cmdos/os_ip.go`.
  - Add hermetic unit tests in `cli/cmdos/os_ip_test.go` verifying help, fallback, and show without touching network devices.

---

## Quality Gate & Verification Contract
- Functions <= 15 lines (target <= 8 lines).
- Affirmative booleans only (`is*`, `has*`).
- Zero nested ifs (nesting depth <= 1).
- 100% hermetic unit tests with zero real system alterations.
- Strict Unix LF line endings.
