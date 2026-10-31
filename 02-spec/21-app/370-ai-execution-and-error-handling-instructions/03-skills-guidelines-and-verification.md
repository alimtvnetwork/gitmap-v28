# Skills, Guidelines & Verification Specification

> **Task:** 370-ai-execution-and-error-handling-instructions  
> **Subsystem:** Skills, Anti-Pattern Guidelines, Test Suites & Verification  
> **Related Plans:** `.ai-memory/plans/pending/370-ai-execution-and-error-handling-instructions.md`, `.ai-memory/plans/subtasks/370-ai-execution-and-error-handling-instructions/`  
> **Related Specifications:** [`01-architecture-and-powershell-diagnostics.md`](01-architecture-and-powershell-diagnostics.md), [`readme.md`](readme.md)  
> **Target Release:** `v6.645.0`  
> **Date:** 31-Oct-26  
> **Status:** APPROVED  

---

## 1. Executive Overview

When AI agents navigate repositories or inspect filesystem artifacts, they occasionally mistake execution runners (`gitmap ps`, `gitmap pwsh`) for directory listing or document inspection utilities. This causes PowerShell to interpret directory and document paths as executable cmdlets, throwing `is not recognized as a name of a cmdlet, function, script file, or executable program` or syntax parser errors (e.g., `Unexpected token '-spec'`), terminating with Go execution stack traces (`exit status 1`).

This specification establishes:
1. **Canonical skill updates** to `.agents/skills/gitmap/skill.md` and `.cursor/skills/gitmap/skill.md`.
2. **Anti-pattern guideline additions** to `02-spec/02-coding-guidelines/06-ai-optimization/10-anti-pattern-replacements.md`.
3. **Automated test suite specifications** across `cli/cmddiagnostics/` and `cli/cmd/cmd_powershell_test.go`.
4. **End-to-end verification checklist** covering test scenarios `TC-PS-01` through `TC-PS-07` ensuring zero-stack-trace educational error handling across all repositories.

---

## 2. Skill Updates: `.agents/skills/gitmap/skill.md`

### 2.1 Non-Negotiable Command Replacement Matrix Addition

In Section 2 (*Non-Negotiable Command Replacement Matrix*), append the following mandatory replacement row:

| Prohibited Legacy Tool / Pattern | GitMap Mandatory Command | Why GitMap is Required |
| :--- | :--- | :--- |
| ❌ Passing directory or document paths to PowerShell runners: `gitmap ps <dir>`, `gitmap ps <file.md>`, `gitmap pwsh <file.json>` | ✅ Directory inspection: `gitmap lf <dir>` or `gitmap plans`<br>✅ File reading: `gitmap cat <filepath>` or `gitmap heatmap <filepath>`<br>✅ PowerShell command: `gitmap ps "<cmdlet> [args]"` | `gitmap ps` invokes `pwsh` / `powershell.exe` with `-NoProfile -Command`, executing arguments directly as code. Directories and non-executable documents cannot be run as cmdlets. Bypasses raw Go stack traces and guides agents to instant discovery tools. |

### 2.2 Section 3 Updates (*Script Execution & Runner Engines*)

Update the documentation for `gitmap ps` and `gitmap pwsh` in Section 3 to explicitly outline execution semantics, common agent mistakes, and the 2-to-3 step learning roadmap:

```markdown
### 3. Script Execution & Runner Engines
- `gitmap py "<code-or-script>"` — High-performance cross-platform Python script execution.
- `gitmap py -c "<code-or-expression>"` — Direct inline Python command evaluation.
- `gitmap pwsh "<cmd>"` / `gitmap ps "<cmd>"` — Cross-platform PowerShell execution with `-NoProfile`.
- `gitmap ps -c "<cmd>"` — Direct inline PowerShell command evaluation.
- `gitmap bash "<cmd>"` / `gitmap sh "<cmd>"` — Cross-platform Bash execution.
- `gitmap bash -c "<cmd>"` — Direct inline POSIX Bash evaluation.

> ⚠️ **CRITICAL RUNNER RULE FOR AI AGENTS:**  
> `gitmap ps` and `gitmap pwsh` execute strings as **PowerShell code**, not as file viewers or directory navigators!  
> - **DO NOT PASS DIRECTORIES:** `gitmap ps .ai-memory/plans/subtasks/...` fails because folders are not executable cmdlets.  
> - **DO NOT PASS DOCUMENTS:** `gitmap ps readme.md` fails because markdown files are not executable scripts.  
> - **LEARNING ROADMAP:**  
>   1. To list files in a folder: run `gitmap lf <dir>` or `gitmap plans`.  
>   2. To inspect file contents: run `gitmap cat <filepath>` or `gitmap heatmap <filepath>`.  
>   3. To execute PowerShell: pass valid cmdlets in quotes: `gitmap ps "Get-ChildItem -Path '<dir>'"` or `gitmap ps "Test-Path '<path>'"`
```

---

## 3. Anti-Pattern Guidelines Addition: `10-anti-pattern-replacements.md`

Add a dedicated anti-pattern section to `02-spec/02-coding-guidelines/06-ai-optimization/10-anti-pattern-replacements.md`:

### 3.1 Substitution Matrix Row

| Prohibited Anti-Pattern | Mandatory GitMap Command | Core Rationale & Technical Advantage |
| :--- | :--- | :--- |
| ❌ Invoking `gitmap ps` on directories or document files (`gitmap ps .ai-memory/...`, `gitmap ps spec.md`) | ✅ `gitmap lf <dir>`<br>`gitmap cat <filepath>`<br>`gitmap plans`<br>`gitmap ps "Get-ChildItem '<dir>'"` | `gitmap ps` spawns a live PowerShell shell. Non-executable files and directories trigger cmdlet resolution failures and noisy Go stack traces. Native GitMap inspection tools return in <2ms with clean ANSI tables. |

### 3.2 Deep-Dive Specification: Anti-Pattern 6

```markdown
### Anti-Pattern 6: Directory / Document Execution via `gitmap ps` vs Native GitMap Discovery

#### Context & Forbidden Pattern
AI agents operating in multi-task workflows frequently confuse `gitmap ps` with:
- "Process status" (`ps` in UNIX)
- "Plan summary" (`plans`)
- "Path inspection" (`dir`, `ls`)

When an agent passes a directory path (e.g., `.ai-memory/plans/subtasks/240-...`) or a markdown document to `gitmap ps`, PowerShell attempts to execute the path as a binary or cmdlet, producing:
```text
.ai-memory/plans/... : The term '.ai-memory/plans/...' is not recognized as a name of a cmdlet...
gitmap ps: execute failed: [E9000:EXECUTION] execute powershell: exit status 1
```

```bash
<!-- ❌ FORBIDDEN: Passing directories or text documents directly to gitmap ps -->
gitmap ps .ai-memory/plans/subtasks/240-running-and-queued-prompt-detection-in-multi-instance
gitmap ps 02-spec/21-app/370-ai-execution-and-error-handling-instructions/readme.md
gitmap pwsh cli/cmd/
```

#### Required Replacement
Use dedicated GitMap inspection commands or wrap PowerShell syntax with proper cmdlets:

```bash
<!-- ✅ REQUIRED: Native GitMap tools for directories, files, and plans -->
# 1. Directory listing
gitmap lf .ai-memory/plans/subtasks/240-running-and-queued-prompt-detection-in-multi-instance
gitmap list-files cli/cmd/

# 2. Plan and subtask inspection
gitmap plans
gitmap agent subtask ls 240

# 3. File inspection with automatic cognitive heatmap touch
gitmap cat 02-spec/21-app/370-ai-execution-and-error-handling-instructions/readme.md
gitmap heatmap 02-spec/21-app/370-ai-execution-and-error-handling-instructions/readme.md

# 4. Valid PowerShell execution (if PowerShell is truly needed)
gitmap ps "Get-ChildItem -Path '.ai-memory/plans/subtasks'"
gitmap ps "Test-Path 'readme.md'"
```

#### AI Reflection Diagnostic Prompting
When intercepted by GitMap's pre-flight or post-execution engine, the terminal renders a structured Tokyo Night diagnostic card asking:
1. **"Why did you run this? What was your expectation?"**
2. **"Did you mean to list files, read a document, or execute PowerShell?"**
3. Concrete copy-pasteable replacement commands tailored to the exact path supplied.
```

---

## 4. Test Suite Additions: `cli/cmddiagnostics/` & `cli/cmd/cmd_powershell_test.go`

Unit and integration tests must be implemented in `cli/cmddiagnostics/` and `cli/cmd/cmd_powershell_test.go` to guarantee zero regressions.

### 4.1 Test Structure & Coverage Matrix

| Test Function | Target Logic | Test Input | Expected Outcome |
| :--- | :--- | :--- | :--- |
| `TestInterceptPowerShellPathArgument_Directory` | Pre-flight directory detection | Existing directories (`.ai-memory`, `cli`, `.`, `02-spec`) | Returns `isIntercepted=true, ArgumentType=TargetDirectory` |
| `TestInterceptPowerShellPathArgument_DocumentFiles` | Pre-flight document file detection | Document files (`readme.md`, `config.json`, `notes.txt`) | Returns `isIntercepted=true, ArgumentType=TargetDocument` |
| `TestInterceptPowerShellPathArgument_SubtaskKeyword` | Pre-flight subtask keyword detection | Arguments containing `subtask` or `plans` | Returns `isIntercepted=true, SuggestedCommand="gitmap agent subtask ls"` |
| `TestInterceptPowerShellPathArgument_ExecutableScripts` | Allow script execution | Executable files (`script.ps1`, `build.bat`, `tool.exe`) | Returns `isIntercepted=false` (pass to PowerShell) |
| `TestInterceptPowerShellPathArgument_ValidCommands` | Allow valid cmdlets | Cmdlets (`Get-ChildItem`, `git status`, `Write-Output "test"`) | Returns `isIntercepted=false` (pass to PowerShell) |
| `TestRenderPowerShellDiagnosticInquiry_Directory` | Tokyo Night advice formatting | Directory context | Contains diagnostic questions & suggests `gitmap lf`, `gitmap plans` |
| `TestRenderPowerShellDiagnosticInquiry_Document` | Tokyo Night advice formatting | `.md` file context | Contains diagnostic questions & suggests `gitmap cat`, `gitmap heatmap` |
| `TestPowerShellInterceptor_CmdletNotFound` | Post-execution stderr filter | Stderr with "is not recognized as a name of a cmdlet" | Intercepts error, suppresses raw Go stack, returns clean exit code 2 |

### 4.2 Test Scenario Mapping to Acceptance Criteria (`TC-PS-01` to `TC-PS-07`)

All test cases implement verification matching the scenarios defined in `01-architecture-and-powershell-diagnostics.md`:

```go
// Example test case contract in cli/cmddiagnostics/powershell_diagnostics_test.go
func TestTC_PS_01_DirectoryInterception(t *testing.T) {
    args := []string{"02-spec/21-app"}
    ctx, isIntercepted := InterceptPowerShellPathArgument(args)
    if !isIntercepted {
        t.Fatalf("expected interception for directory target: %v", args)
    }

    if ctx.ArgumentType != TargetDirectory {
        t.Errorf("expected TargetDirectory, got %s", ctx.ArgumentType)
    }

    if ctx.SuggestedCommand != "gitmap ls 02-spec/21-app" {
        t.Errorf("expected suggested command 'gitmap ls 02-spec/21-app', got '%s'", ctx.SuggestedCommand)
    }
}
```

### 4.3 Coding Standards for Test Suite

In accordance with `.agents/rules/gitmap_guidelines.md`:
- Every helper function must be $\le 15$ lines.
- Zero nested `if` statements (enforce early returns and guard clauses).
- All booleans must begin strictly with `is` or `has` (`isDir`, `hasError`, `isIntercepted`).
- Blank lines must follow `}` and precede `return`.
- Typed exit codes (`constants.ExitUsageError` or `constants.ExitGeneralError`) must be asserted rather than magic numbers.

---

## 5. End-to-End Verification Checklist

Before releasing `v6.645.0`, all verification gates below must be executed and confirmed green:

### 5.1 Pre-Flight Path Interception Verification

- [ ] **Directory Target Interception (`TC-PS-01` & `TC-PS-02`):**
  ```bash
  gitmap ps 02-spec/21-app
  gitmap ps .ai-memory/plans/subtasks
  ```
  - **Expected:** Tokyo Night card displayed; asks reflection inquiry; recommends `gitmap ls 02-spec/21-app` or `gitmap lf .ai-memory/plans/subtasks` and `gitmap plans`; exits cleanly with code `1` or `2`.
  - **Forbidden:** No raw Go stack trace (`[E9000:EXECUTION] execute powershell: exit status 1`) and no `ParserError: Unexpected token '-spec'`.

- [ ] **Subtask Keyword Interception (`TC-PS-03`):**
  ```bash
  gitmap ps subtask
  ```
  - **Expected:** Tokyo Night card displayed; recommends `gitmap agent subtask ls`.

- [ ] **Document Target Interception (`TC-PS-04`):**
  ```bash
  gitmap ps cli/cmd/root.go
  gitmap ps readme.md
  ```
  - **Expected:** Tokyo Night card displayed; recommends `gitmap cat <path>` and `gitmap heatmap <path>`; exits cleanly without launching PowerShell.

- [ ] **Process List Without Arguments Pass-Through (`TC-PS-05`):**
  ```bash
  gitmap ps
  ```
  - **Expected:** Runs process list command, displays native system process table, exits cleanly with code `0`.

- [ ] **Valid Expression Pass-Through (`TC-PS-06`):**
  ```bash
  gitmap ps -Command "Get-Date"
  gitmap ps "Write-Output 'GitMap PS verification successful'"
  ```
  - **Expected:** Executes in PowerShell, outputs result, exits cleanly with code `0`.

- [ ] **PowerShell Script File Pass-Through (`TC-PS-07`):**
  ```bash
  gitmap ps scripts/test.ps1
  ```
  - **Expected:** Does not intercept; forwards script file to PowerShell runtime.

### 5.2 Post-Execution Fallback Interceptor Verification

- [ ] **Non-Existent Cmdlet Interception:**
  ```bash
  gitmap ps "InvalidCmdletNameForTesting"
  ```
  - **Expected:** Intercepts "is not recognized as a name of a cmdlet", surfaces educational advice, avoids panic or unhandled stack trace.

### 5.3 Automated Unit Test Verification

- [ ] Run the dedicated PowerShell test suite:
  ```bash
  go test -v ./cli/cmddiagnostics/... ./cli/cmd/... -run "TestPowerShell|TestTC_PS"
  ```
  - **Expected:** All unit tests pass with `PASS`.

### 5.4 Fleet Synchronization & Release Verification

- [ ] **Sync Across Connected Repositories:**
  ```bash
  gitmap cg sync
  ```
  - **Expected:** Skill updates in `.agents/skills/gitmap/skill.md` and guidelines in `02-spec/02-coding-guidelines/` synchronized across fleet nodes without merge conflicts.

- [ ] **Release Ceremony:**
  ```bash
  gitmap cpr "v6.645.0 AI execution instructions and powershell pre-flight diagnostic interceptor"
  ```
  - **Expected:** Minor version bumped, changelog updated, git tag generated, and clean push.
