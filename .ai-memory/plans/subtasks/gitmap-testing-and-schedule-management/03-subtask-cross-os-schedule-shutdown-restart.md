# Subtask 03: Cross-OS Schedule Shutdown and Restart Test Suite

- **Task ID:** `Subtask-03`
- **Subtask Slug:** `03-subtask-cross-os-schedule-shutdown-restart`
- **Assigned Role:** Worker 01
- **Status:** PENDING
- **Parent Plan:** `.ai-memory/plans/gitmap-testing-and-schedule-management.md`
- **Target Specifications:**
  - `02-spec/21-app/gitmap-testing-and-schedule-management/04-safety-and-os-test-isolation.md`
  - `02-spec/21-app/gitmap-testing-and-schedule-management/03-docker-and-linux-container-strategy.md`
  - `.agents/skills/cg-isolate-os-tests/skill.md`

---

## 1. Scope & Exclusive File Ownership

### 1.1 Owned Files (Strictly Disjoint — Exclusive to Worker 01)
- `cli/cmdschedule/schedule_os.go`
- `cli/cmdschedule/schedule_os_test.go`
- `cli/cmdschedule/schedule_power_state.go`
- `cli/cmdschedule/schedule_power_state_test.go`

### 1.2 Read-Only Context References
- `02-spec/21-app/gitmap-testing-and-schedule-management/04-safety-and-os-test-isolation.md`
- `cli/cmdschedule/duration_parser.go`
- `cli/cmdschedule/duration_parser_test.go`
- `cli/apperror/apperror.go`
- `.agents/skills/cg-isolate-os-tests/skill.md`

---

## 2. Core Objectives

1. **Safety & Zero-Reboot Guarantee:** Enforce that zero tests execute native host power commands (`shutdown`, `reboot`, `poweroff`). All tests must use `setupMockExecutor()` with strict `defer cleanup()`.
2. **Cross-OS Command Construction Matrix:** Verify deterministic command emission for Windows (`shutdown /s /t`), Linux (`sudo poweroff` / `sudo reboot`), and Darwin (`sudo shutdown -h now` / `-r now`), along with cancellation commands.
3. **Power State Lifecycle & Persistence:** Test state transitions (`ARMED` -> `COMPLETED` / `CANCELED`), target timestamp calculation, and remaining duration formatting with isolated temporary state files (`t.TempDir()`).
4. **Mock Countdown & Tick Verification:** Test countdown step resolution and `OnTick` callback progression without real wall-clock sleeps.
5. **CLI Subcommand Dispatch:** Validate `status`, `cancel`, `help`, and error handling for invalid duration arguments.

---

## 3. Step-by-Step Implementation Work Breakdown

### Phase 3.1: Injectable Mock Harness Hardening (`schedule_os.go`, `schedule_os_test.go`)
- **Action 3.1.1:** Review `DefaultOSActionExecutor` and ensure all power CLI execution paths route strictly through this function pointer.
- **Action 3.1.2:** Verify `setupMockExecutor()` restores the original executor via `defer cleanup()` in every test function.
- **Action 3.1.3:** Author assertion helper `assertMockExecuted(t, mock, expectedAction, expectedSeconds)` to check call count and parameter fidelity.

### Phase 3.2: Multi-OS Command Generation Test Suite (`schedule_os_test.go`)
- **Action 3.2.1:** Implement unit tests for `buildWindowsActionCommand`:
  - Assert `OSActionShutdown` produces `shutdown [/s /t <sec> /c <msg>]`.
  - Assert `OSActionRestart` produces `shutdown [/r /t <sec> /c <msg>]`.
- **Action 3.2.2:** Implement unit tests for `buildLinuxActionCommand`:
  - Assert `OSActionShutdown` produces `sudo [poweroff]`.
  - Assert `OSActionRestart` produces `sudo [reboot]`.
- **Action 3.2.3:** Implement unit tests for `buildDarwinActionCommand`:
  - Assert `OSActionShutdown` produces `sudo [shutdown -h now]`.
  - Assert `OSActionRestart` produces `sudo [shutdown -r now]`.
- **Action 3.2.4:** Implement unit tests for `BuildCancelOSActionCommand`:
  - Verify Windows produces `shutdown [/a]`.
  - Verify Unix-like platforms produce `shutdown [-c]`.

### Phase 3.3: Power Schedule State Persistence (`schedule_power_state_test.go`)
- **Action 3.3.1:** Author test helper to redirect `stateFileOverride` to `filepath.Join(t.TempDir(), "test_power_state.json")`.
- **Action 3.3.2:** Test `SavePowerScheduleState` and `readPowerScheduleState` round-trip JSON serialization.
- **Action 3.3.3:** Test `GetActivePowerSchedule`:
  - Returns active state when `TriggerAt` is in the future.
  - Automatically transitions state to `COMPLETED` and returns nil when `TriggerAt` has elapsed.
  - Returns nil when status is `CANCELED`.
- **Action 3.3.4:** Test `CancelActivePowerSchedule` marks state as `CANCELED` and saves to disk.
- **Action 3.3.5:** Test `FormatRemainingDuration` across multiple duration boundaries (hours, minutes, seconds, expired).

### Phase 3.4: Countdown Logic & Fast Tick Emulation (`schedule_os_test.go`)
- **Action 3.4.1:** Test `resolveCountdownStep`:
  - Duration $> 1s$ returns `1s`.
  - Duration $< 1s$ returns the exact sub-second duration.
- **Action 3.4.2:** Test `performDelayCountdown` using a simulated fast-ticker callback:
  - Verify tick counter increments properly without invoking `time.Sleep`.
  - Verify that Windows or `OSActionCancel` immediately bypasses the countdown loop.

### Phase 3.5: CLI Subcommand Integration & Error Gates (`schedule_os_test.go`)
- **Action 3.5.1:** Test `RunSchedulePowerCLI` with status subcommands (`status`, `ls`, `check`):
  - Correctly reports active schedule or empty state message.
  - Zero executor invocations.
- **Action 3.5.2:** Test `RunSchedulePowerCLI` with cancel subcommands (`cancel`, `abort`, `stop`, `rm`):
  - Invokes `BuildCancelOSActionCommand` via mock executor.
  - Marks state canceled in state store.
- **Action 3.5.3:** Test invalid duration input (e.g. `invalid`, `-5m`):
  - Returns validation error wrapped with `apperror`.
  - Zero executor invocations.
- **Action 3.5.4:** Test executor failure propagation:
  - When mock returns an error, `RunSchedulePowerCLI` returns the wrapped error.

---

## 4. Verification Evidence & Quality Gates

```bash
# Gate 1: Run isolated power scheduling tests
go test -v -run TestRunSchedulePower ./cli/cmdschedule/...

# Gate 2: Run power state persistence tests
go test -v -run TestPowerScheduleState ./cli/cmdschedule/...

# Gate 3: Fast test duration gate (entire suite < 500ms)
go test -v ./cli/cmdschedule -run "Test.*SchedulePower.*"
```

### Acceptance Checklist
- [ ] Zero real OS shutdown or restart commands executed against the host.
- [ ] 100% of power action tests use `setupMockExecutor` with `defer cleanup()`.
- [ ] Windows, Linux, and Darwin command generation verified via explicit unit tests.
- [ ] State persistence tests execute against `t.TempDir()` with zero side effects on `~/.gitmap/`.
- [ ] All functions adhere to the <= 15 lines rule with clean guard clauses.
