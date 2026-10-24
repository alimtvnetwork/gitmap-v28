# Subtask 03: Magic Strings Elimination and V6 Error Management Integration

## Status: COMPLETED
- **Subtask Slug**: `subtask-03-magic-strings-and-error-management`
- **Parent Plan**: `.ai-memory/plans/pending/318-gitmap-command-and-coding-guideline-synchronization.md`
- **Architecture Spec**: `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/03-magic-strings-and-error-management-spec.md`
- **Owner**: Guideline & Error Architecture Subagent
- **Target Subsystems**: `cli/cmd`, `cli/constants`, `01-prompts/14-execute`, `.agents/skills/execute-parent-task-with-n-steps-v6`
- **Exclusive Owned Files**:
  - `cli/constants/os_targets.go`
  - `cli/cmd/cmd_explorer.go`
  - `cli/cmd/cmd_whoami.go`
  - `cli/cmd/hd_command_exec.go`
  - `cli/cmd/hd_terminal.go`
  - `cli/cmd/async_exec.go`
  - `cli/cmd/alias.go`
  - `cli/cmd/async_cmd.go`
  - `cli/cmd/cmd_clean.go`
  - `cli/cmd/clihelpers_argv.go`
  - `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`
  - `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`
  - `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/03-magic-strings-and-error-management-spec.md`
  - `.ai-memory/plans/subtasks/318-gitmap-command-and-coding-guideline-synchronization/subtask-03-magic-strings-and-error-management.md`

---

## 1. Context & Objectives

A comprehensive audit of the GitMap codebase and AI execution prompts revealed two interrelated compliance issues:
1. **Magic String Literals in `cli/cmd/`:** Runtime OS comparisons (`runtime.GOOS == "windows"`, `runtime.GOOS == "darwin"`) and subcommand evaluations (`sub == "ls"`, `sub == "dev"`) are written as raw string literals. In `cli/constants/os_targets.go`, constants like `OSTargetWin = "win"` lack the mandatory `*Type` suffix and missing enum validation methods (`IsValid()`, `Is{OS}()`, `String()`).
2. **Missing Error Management Integration in V6 Execute Prompt:** The canonical V6 prompt (`01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`) and companion skill (`.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`) lack explicit requirements for Leaf Error Sovereignty, no `os.Exit` in leaf methods, strongly-typed exit code enums (`ExitCodeType`), and database audit capture in `gitmap.db`. When subagents fail to bubble errors upward, command audit logging is bypassed.

This subtask executes the full remediation across four disjoint micro-batches:
- **Enum Standards Modernization:** Refactor `cli/constants/os_targets.go` to declare `OSTargetType`, `OSGroupType`, and `OrderModeType` with mandatory `*Type` suffixes, backwards-compatible aliases, and positive boolean methods (`IsValid()`, `IsWindows()`, `IsMac()`, `IsUnix()`).
- **Magic String Elimination:** Replace raw OS and subcommand strings in `cli/cmd/` with typed enum evaluations and helper predicates.
- **V6 Prompt & Skill Synchronization:** Embed Leaf Error Sovereignty, Rule R11, `*appfault.AppError` wrappers, typed exit codes, and database audit capture into `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md` and `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`.
- **Strict Coding Standards:** Maintain functions <= 15 lines max, zero nested if statements, blank lines before return and after `}`, positive booleans (`is*`, `has*`), and strict relative Git paths.

---

## 2. Target Files Manifest

| Target File | Purpose / Modifications | Action |
| :--- | :--- | :--- |
| `cli/constants/os_targets.go` | Enforce `OSTargetType`, `OSGroupType`, `OrderModeType`, aliases, and predicate methods | Edit |
| `cli/cmd/cmd_explorer.go` | Replace `runtime.GOOS == "darwin"` with typed platform predicate | Edit |
| `cli/cmd/cmd_whoami.go` | Replace `runtime.GOOS == "windows"` with typed platform predicate | Edit |
| `cli/cmd/hd_command_exec.go` | Replace `runtime.GOOS == "windows"` with typed platform predicate | Edit |
| `cli/cmd/hd_terminal.go` | Replace `runtime.GOOS == "windows"` with typed platform predicate | Edit |
| `cli/cmd/async_exec.go` | Replace `constants.OSWindows == "windows"` with typed platform predicate | Edit |
| `cli/cmd/alias.go` | Refactor raw action comparisons to `AliasActionType` enum | Edit |
| `cli/cmd/async_cmd.go` | Refactor raw action comparisons to `AsyncActionType` enum | Edit |
| `cli/cmd/cmd_clean.go` | Refactor raw action comparisons to `CleanCategoryType` enum | Edit |
| `cli/cmd/clihelpers_argv.go` | Refactor help flag checks to typed enum helpers | Edit |
| `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md` | Embed Leaf Error Sovereignty, Rule R11, and audit capture | Edit |
| `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` | Synchronize error management invariants from V6 prompt | Edit |
| `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/03-magic-strings-and-error-management-spec.md` | Master specification for magic strings and error management | Verified |
| `.ai-memory/plans/subtasks/318-gitmap-command-and-coding-guideline-synchronization/subtask-03-magic-strings-and-error-management.md` | Subtask tracking ledger and checklist | Verified |

---

## 3. Step-by-Step Atomic Checklist

### Phase 1: OS Constants & Enum Methods Modernization (`cli/constants/os_targets.go`)
- [x] **1.1 Define Strongly-Typed Enums with `*Type` Suffix:**
  - Define `type OSTargetType string`.
  - Define `type OSGroupType string`.
  - Define `type OrderModeType string`.
- [x] **1.2 Maintain Backwards-Compatible Aliases:**
  - Define `type OSTarget = OSTargetType`.
  - Define `type OSGroup = OSGroupType`.
  - Define `type OrderMode = OrderModeType`.
- [x] **1.3 Typed Constants & Zero Values:**
  - Declare `OSTargetInvalid OSTargetType = ""`, `OSTargetWin`, `OSTargetUbuntu`, `OSTargetCentOS`, `OSTargetDebian`, `OSTargetFedora`, `OSTargetArch`, `OSTargetMac`, `OSTargetUnix`, `OSTargetAll`.
  - Declare `OSGroupInvalid OSGroupType = ""`, `OSGroupWindows`, `OSGroupUnix`, `OSGroupMac`, `OSGroupPOSIX`.
  - Declare `OrderModeInvalid OrderModeType = ""`, `OrderModeUnixFirst`, `OrderModeOSFirst`, `OrderModeOSOnly`, `OrderModeFallback`.
- [x] **1.4 Implement Enum Predicate Methods (Positive Booleans Only):**
  - Implement `(t OSTargetType) IsValid() bool`.
  - Implement `(t OSTargetType) IsWindows() bool`.
  - Implement `(t OSTargetType) IsMac() bool`.
  - Implement `(t OSTargetType) IsUnix() bool`.
  - Implement `(t OSTargetType) String() string` and `(t OSTargetType) Label() string`.
  - Implement `(g OSGroupType) IsValid() bool` and `(g OSGroupType) IsWindows() bool`.

### Phase 2: Magic String Elimination across `cli/cmd/`
- [x] **2.1 OS Check Refactoring:**
  - In `cli/cmd/cmd_explorer.go`, replace `runtime.GOOS == "darwin"` with typed platform check.
  - In `cli/cmd/cmd_whoami.go`, replace `runtime.GOOS == "windows"` with typed platform check.
  - In `cli/cmd/hd_command_exec.go`, replace `runtime.GOOS == "windows"` with typed platform check.
  - In `cli/cmd/hd_terminal.go`, replace `runtime.GOOS == "windows"` with typed platform check.
  - In `cli/cmd/async_exec.go`, replace `constants.OSWindows == "windows"` with typed platform check.
- [x] **2.2 Subcommand & Argument Routing Refactoring:**
  - In `cli/cmd/alias.go`, replace raw string checks (`sub == "ls"`, `sub == "set"`) with typed `AliasActionType` enum comparisons.
  - In `cli/cmd/async_cmd.go`, replace raw string checks (`sub == "ls"`, `sub == "stop"`) with typed `AsyncActionType` enum comparisons.
  - In `cli/cmd/cmd_clean.go`, replace chained string branches (`sub == "terminal"`, `sub == "dev"`, `sub == "cache"`) with `CleanCategoryType` enum checks using `IsAnyOf(...)`.
  - In `cli/cmd/clihelpers_argv.go`, refactor `arg == "help" || arg == "--help"` to use centralized typed constants.

### Phase 3: Error Management Integration in V6 Execute Prompt & Skill
- [x] **3.1 Update Section 1.1.E in `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`:**
  - Mandate Leaf Error Sovereignty: leaf functions returning `error` or `result.Wrap[T]` must return errors directly upward.
  - Ban `os.Exit` and `cliexit.HandleError` inside leaf functions.
  - Mandate bubbling to root command dispatcher to ensure SQLite audit capture via `finishCommandAudit` and `persistLastError` in `gitmap.db`.
  - Enforce strongly-typed `ExitCodeType` enums and specialized exit helpers (`HandleGeneralError`, `HandleUsageError`, `HandleValidationError`).
- [x] **3.2 Embed Rule R11 in V6 Rules Taxonomy:**
  - Add Rule R11 explicitly detailing error return sovereignty and audit log preservation.
- [x] **3.3 Update Phase 2 Constraints & Section 13 Quality Checklist:**
  - Add verification checklist items for Leaf Error Sovereignty, zero leaf `os.Exit`, and typed exit codes.
- [x] **3.4 Synchronize Antigravity Skill `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`:**
  - Ensure skill instructions exactly match the updated canonical V6 prompt text.

### Phase 4: Verification & Quality Audit
- [x] **4.1 Run FastGate AST Linter:** Execute `python 03-ai-scripts/50-fastgate.py` to confirm zero syntax errors.
- [x] **4.2 Run Forbidden Strings Check:** Execute `python linter-scripts/check-forbidden-strings.py` to verify zero raw integer exit codes.
- [x] **4.3 Run Enum Guidelines Linter:** Execute `python linter-scripts/check-enum-guidelines.py` to ensure all enum types bear `Type` suffix.
- [x] **4.4 Run Relative Path Linter:** Execute `python linter-scripts/check-relative-paths.py` to confirm 100% relative Git paths.

---

## 4. Technical Invariants & Negative Guardrails

1. **Leaf Error Return Sovereignty:** Leaf functions in `cli/cmd/`, `cli/store/`, and `pkg/` that declare an `error` return type MUST return the error directly. They MUST NEVER call `cliexit.HandleError`, `cliexit.Exit`, or `os.Exit` and return `nil`.
2. **Database Audit Capture:** Exits belong exclusively to the top-level command dispatch orchestrator in `cli/cmd/root.go` (`runDispatch` / `dispatch`), ensuring that all errors are recorded in `gitmap.db` by `finishCommandAudit`.
3. **No Magic Integer Exit Codes:** Never pass un-typed integers (`1`, `2`) to `cliexit.HandleError`. Always use `cliexit.ExitCodeType` constants or specialized helpers (`HandleGeneralError`, `HandleUsageError`, `HandleValidationError`).
4. **Mandatory `*Type` Enum Suffix:** Every enum type definition must end in `Type`. Aliases (`type Enum = EnumType`) must be provided for backwards compatibility.
5. **Positive Boolean Naming:** All booleans must begin with `is` or `has` only (`isValid`, `hasMatch`). Negative booleans (`isNotValid`, `disabled`) and banned prefixes (`can`, `should`, `was`, `will`) are strictly forbidden.
6. **Strict Relative Git Paths:** Zero drive letters or absolute file URIs (`d:\`, `C:\`, `file:///`) are permitted in code, documentation, or commit messages.
7. **Function Length Limit:** Functions must stay <= 15 lines max (<= 8 lines preferred).
8. **Zero Nested Ifs:** All nested conditionals must be flattened with guard clauses and early returns.
9. **Return Styling:** Mandatory blank line after `}` and blank line before `return`.
10. **Zero Builds and Zero Tests (Rule R1):** Do not run full `go build` or `go test ./...`. FastGate AST verification only.

---

## 5. Verification & Acceptance Evidence Contract

### Verification Gates

```bash
# Gate 1: FastGate AST syntax and hygiene check
python 03-ai-scripts/50-fastgate.py

# Gate 2: Verify zero magic integer exit codes and forbidden string patterns
python linter-scripts/check-forbidden-strings.py

# Gate 3: Verify all enum types adhere to the *Type suffix convention
python linter-scripts/check-enum-guidelines.py

# Gate 4: Verify strict relative Git path compliance across all modified files
python linter-scripts/check-relative-paths.py
```

### Final Atomic Commit Verification

```bash
gitmap cpf "SPEC - add magic strings elimination and v6 error management spec and subtask"
```
