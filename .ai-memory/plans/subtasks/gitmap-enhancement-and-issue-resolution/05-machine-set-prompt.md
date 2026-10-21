# Subtask 05: Machine Set Interactive Prompt for Missing Machine Name

- **Subtask Path**: `.ai-memory/plans/subtasks/gitmap-enhancement-and-issue-resolution/05-machine-set-prompt.md`
- **Parent Plan**: [.ai-memory/plans/gitmap-enhancement-and-issue-resolution.md](../../gitmap-enhancement-and-issue-resolution.md)
- **Specification Reference**: [02-spec/21-app/gitmap-enhancement-and-issue-resolution/02-component-spec.md](../../../02-spec/21-app/gitmap-enhancement-and-issue-resolution/02-component-spec.md)
- **Primary Owner**: Worker 02 (CLI UX & Configuration Specialist)
- **Target Files (Exclusive)**:
  - `cli/cmdos/os_machine_alias.go` (routing and dispatch update)
  - `cli/cmdos/os_machine_alias_prompt.go` (new interactive prompt module)
  - `cli/cmdos/os_machine_alias_prompt_test.go` (new unit and mock prompt tests)
- **Status**: READY
- **Target Release**: GitMap v6.590.0

---

## 1. Context & Objectives

Currently, when a user executes `gitmap machine set` or `gitmap alias set` without arguments:
- In `cli/cmdos/os_machine_alias.go`, `executeMachineAliasSet` checks `if len(args) == 0` and immediately terminates with an error:
  `missing target name/alias (e.g. gitmap machine set dev-win-01 -y)` with error code `E_MISSING_MACHINE_NAME`.
- While returning an error is required for non-interactive automation (CI/CD pipelines, headless scripts, and `--json` invocations), human users in terminal sessions should be prompted to supply a machine name, with the system hostname offered as an intuitive default.

### Key Objectives:
1. Detect whether the current session is an interactive terminal (`isInteractiveStdin()`).
2. If the machine name is omitted in an interactive terminal, prompt the user via `bufio.NewReader(os.Stdin)`.
3. Provide the detected OS hostname (or existing machine name) as a default candidate in brackets `[default: <hostname>]`.
4. If the user presses ENTER without typing, automatically adopt the default candidate.
5. Sanitize and validate input tokens (alphanumeric, hyphens, no invalid symbols).
6. Preserve non-interactive failure behavior (`E_MISSING_MACHINE_NAME`) for automated scripts, CI, and `--json` flags.

---

## 2. Step-by-Step Implementation Plan

### Step 2.1: Create `cli/cmdos/os_machine_alias_prompt.go`

To maintain file size discipline (keeping `os_machine_alias.go` compact and <= 250 lines):
1. **Interactive TTY Detection**:
   ```go
   func isInteractiveStdin() bool {
       if os.Getenv("CI") != "" || os.Getenv("GITMAP_NON_INTERACTIVE") == "1" {
           return false
       }
       fi, err := os.Stdin.Stat()
       if err != nil {
           return false
       }
       return (fi.Mode() & os.ModeCharDevice) != 0
   }
   ```
2. **Resolve Default Candidate**:
   - Query `os.Hostname()`. If empty or loopback, fallback to current configured `machine.name`. If still empty, fallback to `"gitmap-node"`.
3. **Interactive Read Routine (`promptMachineNameInput`)**:
   - Accept `mode string` ("machine" or "alias") and `defaultCandidate string`.
   - Print formatted prompt: `  ? Enter new %s name [default: %s]: `.
   - Read line using `bufio.NewReader(os.Stdin).ReadString('\n')`.
   - Trim `\r`, `\n`, and whitespace.
   - If trimmed length is 0, return `defaultCandidate, nil`.
   - Sanitize via `sanitizeMachineToken` and validate charset via `isValidMachineToken`.
   - Return sanitized token.
4. **Token Validation (`isValidMachineToken`)**:
   - Ensure length between 1 and 63 characters.
   - Disallow slashes, backslashes, colons, spaces, and control characters.
   - Return positive boolean `isValid`.

---

### Step 2.2: Refactor `executeMachineAliasSet` in `cli/cmdos/os_machine_alias.go`

1. **Update Guard Clause**:
   ```go
   func executeMachineAliasSet(mode string, args []string, isSSH, isJSON bool) error {
       fmt.Printf("  %sSample format: dev-win-01, ubuntu-node-02, mac-studio-01%s\n", constants.ColorYellow, constants.ColorReset)
       
       if len(args) == 0 {
           if !isJSON && isInteractiveStdin() {
               defaultName := resolveDefaultMachineCandidate()
               promptedName, promptErr := promptMachineNameInput(mode, defaultName)
               if promptErr == nil && len(promptedName) > 0 {
                   args = []string{promptedName}
               }
           }
       }
       
       if len(args) == 0 {
           return apperror.NewSimple("missing target name/alias (e.g. gitmap "+mode+" set dev-win-01 -y)", "E_MISSING_MACHINE_NAME")
       }
       
       if isSSH {
           return executeSSHMachineAliasSet(mode, args, isJSON)
       }
       
       newVal := sanitizeMachineToken(args[0])
       return applyLocalMachineOrAliasSet(mode, newVal, isJSON)
   }
   ```
2. Ensure `executeMachineAliasSet` remains <= 15 lines by decomposing the prompt check into a small helper `resolveTargetArgs(mode string, args []string, isJSON bool) ([]string, error)` if needed.

---

### Step 2.3: Implement Unit & Mock Tests in `cli/cmdos/os_machine_alias_prompt_test.go`

1. `TestPromptMachineNameDefaultOnEmptyEnter`: Pass simulated `\n` to `promptMachineNameInputFromReader`, verify it returns `defaultCandidate`.
2. `TestPromptMachineNameCustomInput`: Pass simulated `my-new-node\n`, verify it returns `my-new-node`.
3. `TestPromptMachineNameWhitespaceOnly`: Pass simulated `   \n`, verify it returns `defaultCandidate`.
4. `TestExecuteMachineAliasSetNonInteractiveJson`: Verify `executeMachineAliasSet` returns `E_MISSING_MACHINE_NAME` error when `isJSON = true`.
5. `TestIsValidMachineToken`: Test valid tokens (`dev-win-01`, `node.02`, `mac-mini`) vs invalid tokens (`node/name`, `bad:token`, `a b`).

---

## 3. Coding Guidelines & Invariant Verification

- [ ] **Error Sovereignty**: Leaf functions return `error` or `*apperror.AppError` directly. Zero calls to `cliexit.HandleError` or `os.Exit`.
- [ ] **Boolean Naming**: All booleans prefixed with `is*` or `has*` only (`isInteractive`, `isJSON`, `isValidToken`).
- [ ] **Function Sizing**: <= 15 lines max per function (<= 8 lines preferred).
- [ ] **Conditional Flattening**: Zero nested `if` statements. Guard clauses and early returns throughout.
- [ ] **Path Hygiene**: Strictly relative paths in code, comments, specs, and tests.

---

## 4. Verification & Validation Commands

1. **Unit Test Suite**:
   ```bash
   go test -v ./cli/cmdos/ -run TestPromptMachine
   ```
2. **Interactive Terminal Manual Test**:
   ```bash
   gitmap machine set
   ```
   Verify prompt appears, displays `[default: <hostname>]`, accepts ENTER to apply default, and prints green confirmation.
3. **Non-Interactive Verification**:
   ```bash
   gitmap machine set --json
   ```
   Verify non-zero exit code and structured error output with `E_MISSING_MACHINE_NAME`.
