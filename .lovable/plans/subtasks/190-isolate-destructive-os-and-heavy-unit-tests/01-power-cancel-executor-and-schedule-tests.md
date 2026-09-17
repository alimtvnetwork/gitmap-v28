# Subtask 01: Power Cancel Executor & Schedule Power Tests

## Parent Plan
- Parent Plan: `.lovable/plans/pending/190-isolate-destructive-os-and-heavy-unit-tests.md`

## Objectives
1. Refactor `cli/cmdschedule/schedule_os.go`:
   - Add `OSActionCancel OSActionType = "cancel"`.
   - Update `executeNativeOSAction` to skip countdown delay when `params.Action == OSActionCancel`.
   - Update `CancelSchedulePowerCLI` to construct `OSActionParams` with `Action: OSActionCancel`, `Executable: exe`, `Args: cmdArgs`, and route through `DefaultOSActionExecutor(params)`.
2. Update `cli/cmdschedule/schedule_os_test.go`:
   - Add `TestCancelSchedulePowerCLI_Mocked` with `setupMockExecutor` and `defer cleanup()`.
   - Verify that `DefaultOSActionExecutor` receives `OSActionCancel` and cancel arguments without executing real OS abort.
   - Verify fast duration math with 1s/2s delays and mock tick verification.
3. Update `cli/cmdschedule/schedule_power_state_test.go`:
   - Ensure all state tests are fully hermetic.

## Coding Standards
- Functions <= 15 lines (target <= 8 lines).
- Affirmative booleans only (`is*`, `has*`).
- Strict Unix LF line endings.
- Universal AppError return wrapping.
