# Cross-OS Power Command Contracts & CLI Specifications

## 1. Executive Summary

This document specifies the exact command-line interface (CLI) contracts, syntax rules, OS-specific invocation mappings, error envelopes, and test double invariants for GitMap's system power scheduling subsystem (`gitmap schedule shutdown` and `gitmap schedule restart`).

All power commands adhere to GitMap's universal error management architecture, strict boolean naming standards (`IsSuccess`, `HasActiveSchedule`), structured exit codes, and mockable OS execution gates.

---

## 2. CLI Command Hierarchy & Grammar

The schedule power command interface is exposed under the root `gitmap schedule` command group:

```
gitmap schedule shutdown [duration | status | cancel] [--json]
gitmap schedule restart  [duration | status | cancel] [--json]
```

### 2.1 Grammar & Subcommand Dispatch Matrix

| Subcommand Form | Arguments / Flags | Target Operation | Description |
|---|---|---|---|
| **Arm Shutdown** | `gitmap schedule shutdown [duration]` | `OSActionShutdown` | Arm system shutdown after specified duration. |
| **Arm Restart** | `gitmap schedule restart [duration]` | `OSActionRestart` | Arm system reboot after specified duration. |
| **Status Query** | `gitmap schedule <action> status` | Query State | Display active schedule, trigger timestamp, and remaining time. Aliases: `ls`, `check`. |
| **Cancel Action** | `gitmap schedule <action> cancel` | `OSActionCancel` | Cancel active schedule and invoke OS-level abort. Aliases: `abort`, `stop`, `rm`, `remove`. |
| **Help Menu** | `gitmap schedule <action> --help` | Display Help | Print usage syntax, supported duration formats, and examples. Alias: `-h`. |

---

## 3. Parameter Specifications & Duration Resolution

### 3.1 Duration Argument (`[duration]`)
- **Default Value**: `0` (immediate execution).
- **Position**: First positional argument following the power verb.
- **Allowed Formats**:
  - **Colon Formats**: `1:45hr`, `1:45h`, `01:30:00`, `2:15m`.
  - **Day Units**: `1day`, `1d`, `2days`, `0.5d`.
  - **Hour Units**: `2h`, `2hours`, `1.5hr`.
  - **Minute Units**: `120m`, `30mins`, `15m`.
  - **Second Units**: `90s`, `300sec`, `10s`.
  - **Immediate**: `now`, `0`, `0s`.
- **Validation**:
  - Negative values (e.g., `-10m`) are rejected with `E_INVALID_DURATION`.
  - Malformed strings (e.g., `invalid`, `abc`) return an `apperror.AppError` and exit with `constants.ExitUsageError`.

---

## 4. Cross-OS Command Execution Specifications

### 4.1 Native Command Construction Matrix

When a power action is armed, GitMap translates the abstract `OSActionParams` into OS-native executables and arguments via `BuildOSActionCommand(action, seconds)`:

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                               OSActionParams                                 │
│  Action: OSActionShutdown | OSActionRestart                                  │
│  Seconds: int64 (calculated from parsed duration)                            │
└──────────────────────────────────────┬───────────────────────────────────────┘
                                       │
                ┌──────────────────────┼──────────────────────┐
                ▼                      ▼                      ▼
           [ Windows ]             [ Linux ]              [ Darwin ]
     shutdown /s /t <sec>        sudo poweroff          sudo shutdown -h now
     shutdown /r /t <sec>        sudo reboot            sudo shutdown -r now
```

#### Detailed Command Formats:
1. **Windows (`runtime.GOOS == "windows"`)**:
   - Shutdown: `shutdown /s /t <seconds> /c "GitMap scheduled shutdown"`
   - Restart: `shutdown /r /t <seconds> /c "GitMap scheduled restart"`
   - Cancel: `shutdown /a`
   - *Execution Note*: Windows kernel manages the countdown timer via RPC service. GitMap exits immediately upon command registration.

2. **Linux (`runtime.GOOS == "linux"`)**:
   - Shutdown: `sudo poweroff`
   - Restart: `sudo reboot`
   - Cancel: `shutdown -c`
   - *Execution Note*: GitMap executes an in-process delay (`performDelayCountdown`), ticking every second, and executes `sudo poweroff` / `sudo reboot` when `remaining <= 0`.

3. **macOS / Darwin (`runtime.GOOS == "darwin"`)**:
   - Shutdown: `sudo shutdown -h now`
   - Restart: `sudo shutdown -r now`
   - Cancel: `shutdown -c`
   - *Execution Note*: GitMap executes an in-process delay, and executes immediate shutdown when the timer elapses.

---

## 5. Universal Response Envelope & Status Schemas

When invoked programmatically or with `--json`, power commands emit GitMap's standard Universal JSON Response Envelope.

### 5.1 Active Schedule Status Schema (`PowerScheduleState`)

```json
{
  "isSuccess": true,
  "payload": {
    "action": "shutdown",
    "status": "ARMED",
    "scheduledAt": "2026-10-21T10:30:00Z",
    "triggerAt": "2026-10-21T12:15:00Z",
    "durationSeconds": 6300,
    "remainingDuration": "1h 45m 0s",
    "command": "shutdown /s /t 6300 /c GitMap scheduled shutdown",
    "hasActiveSchedule": true
  },
  "error": null
}
```

### 5.2 Error Envelope Schema

```json
{
  "isSuccess": false,
  "payload": null,
  "error": {
    "code": "E_INVALID_DURATION",
    "message": "invalid schedule duration \"invalid-timer\"",
    "component": "cmdschedule",
    "severity": "ERROR",
    "type": "VALIDATION",
    "details": {
      "raw": "invalid-timer"
    }
  }
}
```

---

## 6. Exit Codes & Error Sovereignty Architecture

Adhering to GitMap's strict error architecture guidelines:
1. **Error Return Sovereignty**: Functions such as `RunSchedulePowerCLI()`, `executeSchedulePowerCLI()`, and `ParseScheduleDuration()` return errors directly to the caller. Leaf functions **never** invoke `os.Exit` or `cliexit.HandleError`.
2. **Outer Command Dispatcher**: The root CLI dispatcher in `cli/root.go` captures the returned error, logs it to the database audit trail via `finishCommandAudit`, and invokes the appropriate exit handler.
3. **Structured Exit Codes**:
   - `constants.ExitSuccess` (`0`): Command completed successfully.
   - `constants.ExitGeneralError` (`1`): System power execution failure or filesystem state I/O error.
   - `constants.ExitUsageError` (`2`): Invalid duration syntax, unknown flags, or unrecognized subcommands.

---

## 7. Safety Invariants & Unit Test Doubles

Automated test execution must satisfy the `cg-isolate-os-tests` guardrails.

### 7.1 Mocking Protocol (`setupMockExecutor`)
In test files (e.g., `cli/cmdschedule/schedule_os_test.go`), native execution is intercepted:

```go
type mockExecutorState struct {
    calledCount int
    lastParams  OSActionParams
    returnErr   error
}

func (m *mockExecutorState) execute(params OSActionParams) error {
    m.calledCount++
    m.lastParams = params
    return m.returnErr
}
```

### 7.2 Safety Invariant Checklist
- [x] `DefaultOSActionExecutor` is replaced with `mock.execute` before any test logic runs.
- [x] A `defer cleanup()` hook restores the original executor when tests exit.
- [x] State file mutations are redirected to an isolated directory using `stateFileOverride = filepath.Join(t.TempDir(), ...)`.
- [x] Under no circumstances shall `exec.Command("shutdown", ...)` or `exec.Command("sudo", "poweroff")` be executed on the host during unit or CI test runs.

---

## 8. Docker Verification & Container Contracts

To guarantee Linux power command compatibility without requiring a dedicated Linux host, containerized verification is specified:

1. **Test Runner Container**:
   - Dockerfile: `Dockerfile.test`
   - Base Image: `golang:1.24-alpine`
   - Command: `go test -v -cover ./cli/cmdschedule/...`
2. **Linux Execution Contract Assertions**:
   - `buildLinuxActionCommand(OSActionShutdown)` yields `("sudo", ["poweroff"])`.
   - `buildLinuxActionCommand(OSActionRestart)` yields `("sudo", ["reboot"])`.
   - `BuildCancelOSActionCommand()` under Linux yields `("shutdown", ["-c"])`.
   - Countdown step logic (`resolveCountdownStep`) functions identically across OS boundaries.
