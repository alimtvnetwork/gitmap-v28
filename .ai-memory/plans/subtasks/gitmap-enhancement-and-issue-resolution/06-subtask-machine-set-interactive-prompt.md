# Subtask 06: Interactive Machine Set Prompt & Input Validation

- **Subtask ID**: `06-subtask-machine-set-interactive-prompt`
- **Parent Task**: `gitmap-enhancement-and-issue-resolution`
- **Component Spec**: `02-spec/21-app/gitmap-enhancement-and-issue-resolution/02-component-spec.md`
- **Owned Files**:
  - `cli/cmdos/os_machine_alias.go`
  - `cli/cmdos/os_machine_alias_prompt.go`
  - `cli/cmdos/os_machine_alias_test.go`
- **Status**: COMPLETED

---

## 1. Objectives & Context

When executing `gitmap machine set` or `gitmap alias set` without arguments in an interactive terminal, GitMap currently returns an immediate error:
`missing target name/alias (e.g. gitmap machine set dev-win-01 -y)` with error code `E_MISSING_MACHINE_NAME`.

While throwing an error is correct for automated non-interactive scripts and machine-readable JSON invocations (`--json`), interactive terminal users should be prompted with a sensible default (the detected system hostname) and guided to input or confirm their target machine name/alias.

This subtask delivers:
1. Interactive TTY detection via `!isJSON && term.IsTerminal(int(os.Stdin.Fd()))`.
2. Interactive terminal prompt with default fallback to system hostname.
3. Token sanitization and validation (alphanumeric, hyphens, no special characters).
4. Direct application via `applyLocalMachineOrAliasSet` without failing.
5. Preservation of existing non-interactive / JSON error semantics.

---

## 2. Step-by-Step Implementation Plan

### Step 2.1: Modular Prompt Logic (`cli/cmdos/os_machine_alias_prompt.go`)
- Create `cli/cmdos/os_machine_alias_prompt.go` to keep `os_machine_alias.go` under the 300-line limit.
- Implement `isInteractiveTerminal() bool` (<= 5 lines):
  - Return `term.IsTerminal(int(os.Stdin.Fd()))`.
- Implement `resolveDefaultMachineName() string` (<= 10 lines):
  - Hostname lookup via `os.Hostname()`.
  - Fallback to current configured name or `"gitmap-node"`.
- Implement `promptInteractiveMachineName(mode string, defaultName string) (string, error)` (<= 15 lines):
  - Display styled prompt: `? Enter new %s name [default: %s]: `.
  - Read input via `bufio.NewReader(os.Stdin).ReadString('\n')`.
  - Trim whitespace and newlines.
  - If input is empty, return `defaultName, nil`.
  - Return sanitized input token.

### Step 2.2: Refactor `executeMachineAliasSet` (`cli/cmdos/os_machine_alias.go`)
- In `executeMachineAliasSet(mode string, args []string, isSSH, isJSON bool)`:
  - If `len(args) == 0`:
    - Check guard: `if !isJSON && isInteractiveTerminal()`:
      - Call `promptInteractiveMachineName(mode, resolveDefaultMachineName())`.
      - If prompt returns valid token, set `args = []string{token}`.
    - If `len(args) == 0`:
      - Return `apperror.NewSimple("missing target name/alias (e.g. gitmap "+mode+" set dev-win-01 -y)", "E_MISSING_MACHINE_NAME")`.
  - Continue standard dispatch to `executeSSHMachineAliasSet` or `applyLocalMachineOrAliasSet`.

### Step 2.3: Input Validation & Token Sanitization
- Implement `isValidMachineToken(token string) bool` (<= 10 lines):
  - Ensure token length > 0 and <= 63 characters.
  - Ensure token contains only alphanumeric characters, hyphens, and dots.
- If input fails validation:
  - Print user-friendly formatting guidance.
  - Return `apperror.NewSimple("invalid machine name format: alphanumeric and hyphens only", "E_INVALID_MACHINE_TOKEN")`.

### Step 2.4: Unit & Interactive Mock Tests (`cli/cmdos/os_machine_alias_prompt_test.go`)
- Test non-interactive mode: verify `executeMachineAliasSet` returns `E_MISSING_MACHINE_NAME` when `isJSON = true`.
- Test mocked stdin input:
  - Empty string simulates pressing ENTER -> asserts default hostname is chosen.
  - Non-empty string `dev-win-01` -> asserts `dev-win-01` is applied.
  - Whitespace string `   ` -> asserts default hostname is chosen.
- Test token sanitization: spaces converted to hyphens.

---

## 3. Coding Guidelines Invariants Checklist

- [x] Error Sovereignty: All functions return error directly; leaf functions do not call exit handlers.
- [x] Boolean Naming: All booleans use `is*` or `has*` only (`isInteractiveTerminal`, `isJSON`, `isAppErr`, `isTerminalFn`).
- [x] Function Size: All functions <= 15 lines (<= 8 lines preferred).
- [x] Conditional Flattening: Zero nested `if` statements; early returns used throughout.
- [x] Relative Path Hygiene: Zero absolute paths; strictly relative paths in code and tests.

---

## 4. Verification & Validation Commands

1. **Unit Tests Execution**:
   ```bash
   go test -v ./cli/cmdos/ -run TestPromptForMachineName
   go test -v ./cli/cmdos/ -run TestExecuteMachineAliasSet
   ```
2. **Interactive Terminal Verification**:
   ```bash
   gitmap machine set
   ```
   Verify prompt appears, allows typing or pressing ENTER, and updates successfully.
3. **Automated / JSON Non-Interactive Verification**:
   ```bash
   gitmap machine set --json
   ```
   Verify exit code is non-zero and outputs structured JSON error with `E_MISSING_MACHINE_NAME`.

---

## 5. Subtask Execution & Verification Summary

- **Completed**: 2026-10-21
- **Subtask**: `06-subtask-machine-set-interactive-prompt`
- **Component**: Component C (Interactive Machine Name & Alias Configuration Prompt)
- **Delivered**:
  - `cli/cmdos/os_machine_alias_prompt.go`: `isInteractiveTerminal()`, `resolveDefaultMachineName()`, `promptForMachineName()`, `resolveMachineSetArgs()`, `promptInteractiveMachineArg()` with zero nested ifs and <= 15 lines.
  - `cli/cmdos/os_machine_alias.go`: refactored `executeMachineAliasSet` to dispatch through `resolveMachineSetArgs`, preserving non-interactive / `--json` `E_MISSING_MACHINE_NAME` error semantics while enabling interactive fallback.
  - `cli/cmdos/os_machine_alias_test.go`: unit tests verifying default hostname resolution, non-interactive JSON error, non-interactive terminal error, interactive prompt reading, default fallback on ENTER/whitespace, and token sanitization.
