---
name: cg-isolate-os-tests
description: Autonomously audits, refactors, and validates repository-wide unit tests to ensure that tests NEVER trigger real OS shutdown, reboot, power-off, system modifications, or heavy unmocked system calls, enforcing injectable executors and mock duration verification across polyglot stacks.
---

# Skill: Unit Test Isolation & Destructive OS Prevention (`cg-isolate-os-tests`)

This skill governs autonomous auditing, refactoring, and validation of unit tests to strictly prevent execution of destructive host OS commands (shutdown, reboot, poweroff, system package modifications, user group alterations, unmocked filesystem purges) by enforcing injectable command executors and mock duration math.

---

## 1. Mandatory Architectural Invariants

### 1. Total Ban on Real OS Shutdown / Reboot in Unit Tests
The following command invocations MUST NEVER execute directly against the host OS during test runs:
- `shutdown` (`/s`, `/r`, `/a`, `-h`, `-r`, `-c`)
- `reboot`
- `poweroff` / `halt`
- `init 0` / `init 6`
- `systemctl poweroff` / `systemctl reboot`
- Direct package managers: `apt-get upgrade`, `Get-WindowsUpdate`, `softwareupdate`

### 2. Mandatory Injectable Executors
Every production subsystem that interacts with host power controls, system packages, or filesystem caches MUST expose an injectable executor or runner function:

```go
type OSActionExecutor func(params OSActionParams) error
var DefaultOSActionExecutor OSActionExecutor = executeNativeOSAction
```

```go
type FileRemover func(filePath string) (int, int64)
var defaultFileRemover FileRemover = removeSingleFileSafely
```

```go
type TempDirResolver func() []string
var defaultTempDirResolver TempDirResolver = resolveTempDirectories
```

```go
type OSCommandRunner func(cmd *exec.Cmd) error
var defaultOSCommandRunner OSCommandRunner = func(cmd *exec.Cmd) error { return cmd.Run() }
```

### 3. Hermetic Mock Testing with `defer` Restoration
Unit tests must ALWAYS substitute the default executor inside a `defer` restoration block so that test execution remains 100% isolated:

```go
mock, cleanup := setupMockExecutor()
defer cleanup()

// Execute CLI or unit logic
err := RunSchedulePowerCLI(OSActionShutdown, []string{"1s"})
// Assert mock captured parameters without invoking host commands
```

### 4. Fast Duration Math (1s/2s Checks)
Tests for countdown timers, delayed power commands, or recurring schedules must test:
- Duration string parsing (`ParseScheduleDuration`)
- Arithmetic target calculation (`TargetTime = ScheduledAt + Duration`)
- Fast simulated delays (`1s`, `2s`) with mock tick callbacks rather than sleeping or arming host OS timers.

### 5. Function Sizing & Affirmative Booleans
- Functions target **<= 8 lines** of body logic (hard cap: 15 lines).
- Affirmative booleans only (`is*`, `has*`). No negative booleans.
- Universal `*apperror.AppError` return wrapping.

---

## 2. Code Patterns: Before vs After

### Pattern 1: Power Cancellation Command (Go)

```go
// ❌ ANTI-PATTERN: Executes real shutdown /a or shutdown -c on developer workstation!
func CancelSchedulePowerCLI(action OSActionType) error {
    _ = CancelActivePowerSchedule()
    exe, cmdArgs := BuildCancelOSActionCommand()
    cmd := exec.Command(exe, cmdArgs...)
    _ = cmd.Run()
    return nil
}
```

```go
// ✅ COMPLIANT PATTERN: Routes through injectable DefaultOSActionExecutor
func CancelSchedulePowerCLI(action OSActionType) error {
    _ = CancelActivePowerSchedule()
    exe, cmdArgs := BuildCancelOSActionCommand()
    params := buildCancelActionParams(exe, cmdArgs)
    if err := DefaultOSActionExecutor(params); err != nil {
        return err
    }
    fmt.Printf("✓ Canceled active scheduled %s.\n", action)
    return nil
}

func buildCancelActionParams(exe string, cmdArgs []string) OSActionParams {
    return OSActionParams{
        Action:     OSActionCancel,
        Executable: exe,
        Args:       cmdArgs,
    }
}
```

### Pattern 2: Hermetic Test Substitution (Go)

```go
// ❌ ANTI-PATTERN: Tests that call real unmocked system commands
func TestCancelSchedulePower(t *testing.T) {
    err := CancelSchedulePowerCLI(OSActionShutdown) // Shuts down or errors on host!
}
```

```go
// ✅ COMPLIANT PATTERN: Hermetic mock with defer restoration
func TestCancelSchedulePowerCLI_Mocked(t *testing.T) {
    mock, cleanup := setupMockExecutor()
    defer cleanup()

    if err := CancelSchedulePowerCLI(OSActionShutdown); err != nil {
        t.Fatalf("unexpected cancel error: %v", err)
    }
    if mock.calledCount != 1 || mock.lastParams.Action != OSActionCancel {
        t.Errorf("expected 1 mock call with OSActionCancel")
    }
}
```

### Pattern 3: Ephemeral File Cleaner Mocking (Go)

```go
// ❌ ANTI-PATTERN: Tests delete real user cache directories
func TestAIClean(t *testing.T) {
    RunOSAICleanCLI([]string{"-y"}) // Deletes ~/.gemini or temp files!
}
```

```go
// ✅ COMPLIANT PATTERN: Mock file remover intercepts all file deletions
func TestRunOSAICleanCLI_Mocked(t *testing.T) {
    captured, cleanup := setupMockFileRemover()
    defer cleanup()

    if err := RunOSAICleanCLI([]string{"--dry-run"}); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}
```
