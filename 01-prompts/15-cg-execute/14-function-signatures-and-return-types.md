# Function Signatures, Invocations & Result Envelopes — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-functions`, `cg-signatures`, `cg-return-types`, `cg-execute functions`, `audit function naming`, `fix return types`, `enforce apperror`, `enforce result envelope`, `single return type audit`, `multi-line arguments`, `function call formatting`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, discover, plan, refactor, and format all function definitions, call-site invocations, parameter lists, boolean predicate prefixes, multi-value returns, raw generic errors, and result envelopes across the codebase, enforcing one-argument-per-line formatting for >2 arguments, semantic verb/predicate naming, universal `*AppError` wrapping, single `Result[T]` return envelopes with complete predicate methods (`IsSuccess`, `IsFailed`, `HasError`, `HasNoError`, `HasValidError`), and strict relative Git paths until 100% green without stopping.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase to inventory all function definitions (>2 params) on a single line, all function invocations (>2 args) on a single line, functions with generic error returns (`(T, error)`), unformatted boolean predicates missing `is`/`has` (`can` is not acceptable) prefixes, raw `fmt.Errorf()` or `errors.New()`, anti-garbage names, and functions lacking typed `Result[T]` envelopes.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.lovable/plans/pending/XX-function-signatures-audit.md` with an exhaustive Violation Ledger table.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.lovable/plans/subtasks/XX-function-signatures/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/01-index.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and format parameter declarations and call-site invocations to one line per argument with trailing commas.
6. [ ] /goal Phase 2 (Step B): Refactor signatures to use semantic verb/predicate names, convert multi-value returns to single `Result[T]` envelopes with `*AppError` wrappers.
7. [ ] /goal Phase 2 (Step C): Enforce <= 8–15 line function decomposition, single return types, and clean formatting.
8. [ ] /goal Phase 2 (Step D): Execute local linters (`python linter-scripts/check-function-lengths.py`, `check-newline-styling.py`) to verify 0 remaining violations.
9. [ ] /goal Phase 2 (Step E): Execute local CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` with exit code 0 (`exit 0`).
10. [ ] /learn Ingest `.lovable/memory/01-index.md` for project memory index and past learnings.
11. [ ] /learn Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
12. [ ] /learn Ingest `spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
13. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for Rule 9a/9b multi-line parameter and call formatting.
14. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/10-function-naming.md` for semantic verb and predicate prefix standards.
15. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for hallucination prevention and micro-tasking.
16. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for strict relative path citation requirements.
17. [ ] /learn Ingest `spec/03-error-manage/01-index.md` for universal AppError wrapping and error envelopes.
18. [ ] /learn Ingest `spec/03-error-manage/01-index.md` for Result[T] and standardized API envelopes.
19. [ ] /learn Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.
20. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Signatures & Calls, Build Violation Ledger in .lovable/plans/pending/, Subtasks, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Actively Edit Code, Format Multi-Line Params & Calls, Refactor Signatures, Verify Local CI)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Dedicated Section: Function Definitions, Invocations & Multi-Line Standards

A clean codebase relies on deterministic function contracts and clean, scannable call sites. Squeezing long argument lists onto a single line destroys git diff readability and introduces merge conflicts.

---

### 1. Multi-Line Parameter Declarations (>2 Parameters)

When a function, method, or constructor definition has **more than two parameters** (i.e. 3 or more), or exceeds 100 characters in length, **each parameter MUST be placed on its own line** with consistent indentation and a trailing comma (where syntax permits).

#### ❌ FORBIDDEN (>2 parameters on a single line):

```go
// Go
func SaveRecord(label string, path string, isSuccess bool, errMsg string) error { ... }

// TypeScript
function saveRecord(label: string, path: string, isSuccess: boolean, error?: string): void { ... }

// PHP
function saveRecord(string $label, string $path, bool $isSuccess, ?string $error): void { ... }

// Python
def save_record(label: str, path: str, is_success: bool, error: Optional[str] = None) -> None: ...
```

#### ✅ REQUIRED (One parameter per line with trailing comma):

```go
// Go
func SaveRecord(
    label string,
    path string,
    isSuccess bool,
    errMsg string,
) Result[bool] {
    // ...
}
```

```typescript
// TypeScript
function saveRecord(
    label: string,
    path: string,
    isSuccess: boolean,
    error?: string,
): Result<void> {
    // ...
}
```

```php
// PHP
function saveRecord(
    string $label,
    string $path,
    bool $isSuccess,
    ?string $error,
): Result {
    // ...
}
```

```python
# Python
def save_record(
    label: str,
    path: str,
    is_success: bool,
    error: Optional[str] = None,
) -> Result[bool]:
    # ...
```

```rust
// Rust
pub fn save_record(
    label: &str,
    path: &Path,
    is_success: bool,
    error: Option<&str>,
) -> Result<bool, AppError> {
    // ...
}
```

---

### 2. Multi-Line Function Invocations & Call Sites (>2 Arguments)

When calling any function, method, or constructor with **more than two arguments**, or when the call site exceeds 100 characters, **each argument MUST be placed on its own line** with consistent indentation and a trailing comma.

#### ❌ FORBIDDEN (>2 arguments on a single line):

```go
// Go
res := executeQuery(ctx, "SELECT * FROM users WHERE status = ?", statusVal, 50, 0)

// TypeScript
const res = logAction(userId, ActionType.Login, null, StatusType.Success, null, clientIp);

// PHP
$this->logAction($agentId, ActionType::AgentTest->value, null, StatusType::Failed->value, $error->getMessage());
```

#### ✅ REQUIRED (One argument per line with trailing comma):

```go
// Go
res := executeQuery(
    ctx,
    "SELECT * FROM users WHERE status = ?",
    statusVal,
    50,
    0,
)
```

```typescript
// TypeScript
const res = logAction(
    userId,
    ActionType.Login,
    null,
    StatusType.Success,
    null,
    clientIp,
);
```

```php
// PHP
$this->logAction(
    $agentId,
    ActionType::AgentTest->value,
    null,
    StatusType::Failed->value,
    $error->getMessage(),
);
```

```python
# Python
response = dispatch_event(
    event_name=EventNameType.USER_CREATED,
    payload=user_payload,
    retry_count=3,
    timeout_seconds=30,
)
```

---

### 3. No Boolean Flag Parameters (Split Intent-Driven Methods)

When a boolean parameter changes the fundamental **behavior or meaning** of an operation, **NEVER pass a bare boolean flag**. Split the behavior into two explicitly named, self-documenting functions.

```typescript
// ❌ FORBIDDEN: Boolean flag hides caller intent
function logMessage(message: string, isWithStack: boolean): void { ... }
logMessage("Payment failed", true); // What does 'true' do?

// ✅ REQUIRED: Explicitly named functions
function logMessage(message: string): void { ... }
function logMessageWithStack(message: string): void { ... }

logMessage("User saved");
logMessageWithStack("Payment failed");
```

---

### 4. The Single Return Type Principle & Universal `Result[T]` Envelope Architecture

In domain services, handlers, and internal business logic, functions MUST return a single encapsulated result envelope rather than raw multi-value tuples `(T, error)` or unhandled exceptions.

#### 4a. Production-Ready Go `Result[T]` Architecture

```go
package result

import "gitmap/apperror"

// Result encapsulates a computation outcome with typed value or *apperror.AppError.
type Result[T any] struct {
    Value    T
    Err      *apperror.AppError
    Data     T
    AppError error
}

// IsSuccess reports whether the result represents a successful operation.
func (r Result[T]) IsSuccess() bool {
    return r.Err == nil && r.AppError == nil
}

// IsFailed reports whether the result represents a failed operation.
func (r Result[T]) IsFailed() bool {
    return r.Err != nil || r.AppError != nil
}

// IsFailure reports whether the result represents a failed operation (alias).
func (r Result[T]) IsFailure() bool {
    return r.IsFailed()
}

// IsInvalid reports whether the result is invalid or failed.
func (r Result[T]) IsInvalid() bool {
    return r.IsFailed()
}

// HasError reports whether an error is present.
func (r Result[T]) HasError() bool {
    return r.Err != nil || r.AppError != nil
}

// HasNoError reports whether no error exists.
func (r Result[T]) HasNoError() bool {
    return r.Err == nil && r.AppError == nil
}

// HasValidError reports whether an AppError exists and is properly structured.
func (r Result[T]) HasValidError() bool {
    if r.Err != nil {
        return r.Err.IsValid()
    }

    return r.AppError != nil
}

// Unwrap returns the value and error tuple.
func (r Result[T]) Unwrap() (T, *apperror.AppError) {
    if r.Err != nil {
        return r.Value, r.Err
    }

    if r.AppError != nil {
        if appErr, isAppErr := r.AppError.(*apperror.AppError); isAppErr {
            return r.Value, appErr
        }

        return r.Value, apperror.WrapSimple(r.AppError, "result.Unwrap")
    }

    return r.Value, nil
}

// UnwrapOr returns the value if success, or defaultVal if failed.
func (r Result[T]) UnwrapOr(defaultVal T) T {
    if r.IsSuccess() {
        return r.Value
    }

    return defaultVal
}

// ValueOrPanic returns the value if success, or panics with the error.
func (r Result[T]) ValueOrPanic() T {
    if r.IsSuccess() {
        return r.Value
    }

    if r.Err != nil {
        panic(r.Err.Error())
    }

    panic(r.AppError.Error())
}

// SuccessResult constructs a successful Result envelope with Value and Data.
func SuccessResult[T any](val T) Result[T] {
    return Result[T]{
        Value: val,
        Data:  val,
    }
}

// FailureResult constructs a failed Result envelope with *apperror.AppError.
func FailureResult[T any](err *apperror.AppError) Result[T] {
    return Result[T]{
        Err:      err,
        AppError: err,
    }
}

// NewSuccess constructs a successful Result envelope with Data.
func NewSuccess[T any](data T) Result[T] {
    return SuccessResult(data)
}

// NewFailure constructs a failed Result envelope from any error.
func NewFailure[T any](err error) Result[T] {
    if appErr, isAppErr := err.(*apperror.AppError); isAppErr {
        return FailureResult[T](appErr)
    }

    if err == nil {
        return Result[T]{}
    }

    appErr := apperror.WrapSimple(err, "result.NewFailure")

    return FailureResult[T](appErr)
}

// NewFailureWithType constructs a typed failed Result with code, message, and caller.
func NewFailureWithType[T any](
    errCode apperror.ErrorCodeType,
    msg string,
    caller string,
) Result[T] {
    appErr := apperror.New(errCode, msg, caller)
    return FailureResult[T](appErr)
}
```

---

#### 4b. `AppError` Methods & Error Code Comparison

```go
package apperror

// HasError reports whether an error exists.
func (e *AppError) HasError() bool {
    return e != nil
}

// HasNoError reports whether no error exists.
func (e *AppError) HasNoError() bool {
    return e == nil
}

// HasValidError reports whether the AppError is non-nil and has a valid code.
func (e *AppError) HasValidError() bool {
    return e != nil && e.Code != ""
}

// IsErrorCode reports whether the AppError matches the specified ErrorCodeType.
func (e *AppError) IsErrorCode(code ErrorCodeType) bool {
    return e != nil && e.Code == code
}

// IsCode alias for IsErrorCode.
func (e *AppError) IsCode(code ErrorCodeType) bool {
    return e.IsErrorCode(code)
}
```

---

#### 4c. TypeScript `Result<T>` Envelope Architecture

```typescript
import { AppError, ErrorCodeType } from "./apperror";

export type Result<T> = {
    readonly isSuccess: boolean;
    readonly isFailed: boolean;
    readonly hasError: boolean;
    readonly hasNoError: boolean;
    readonly value: T | null;
    readonly data: T | null;
    readonly error: AppError | null;
    unwrap(): [T | null, AppError | null];
    unwrapOr(defaultVal: T): T;
};

export function successResult<T>(val: T): Result<T> {
    return {
        isSuccess: true,
        isFailed: false,
        hasError: false,
        hasNoError: true,
        value: val,
        data: val,
        error: null,
        unwrap: () => [val, null],
        unwrapOr: () => val,
    };
}

export function failureResult<T>(err: AppError): Result<T> {
    return {
        isSuccess: false,
        isFailed: true,
        hasError: true,
        hasNoError: false,
        value: null,
        data: null,
        error: err,
        unwrap: () => [null, err],
        unwrapOr: (defaultVal: T) => defaultVal,
    };
}

export function newFailure<T>(
    code: ErrorCodeType,
    message: string,
    caller: string,
): Result<T> {
    const appErr = new AppError(code, message, caller);
    return failureResult<T>(appErr);
}
```

---

#### 4d. PHP 8.1+ `Result<T>` Class Architecture

```php
<?php

declare(strict_types=1);

namespace App\Common;

use App\Common\Exceptions\AppException;

final class Result
{
    public function __construct(
        public readonly mixed $value = null,
        public readonly ?AppException $error = null,
    ) {}

    public function isSuccess(): bool
    {
        return $this->error === null;
    }

    public function isFailed(): bool
    {
        return $this->error !== null;
    }

    public function isInvalid(): bool
    {
        return $this->isFailed();
    }

    public function hasError(): bool
    {
        return $this->error !== null;
    }

    public function hasNoError(): bool
    {
        return $this->error === null;
    }

    public function unwrap(): mixed
    {
        if ($this->error !== null) {
            throw $this->error;
        }

        return $this->value;
    }

    public function unwrapOr(mixed $defaultVal): mixed
    {
        if ($this->isSuccess()) {
            return $this->value;
        }

        return $defaultVal;
    }

    public static function success(mixed $value): self
    {
        return new self(value: $value);
    }

    public static function failure(AppException $error): self
    {
        return new self(error: $error);
    }
}
```

---

### 5. Semantic Function Naming & Predicate Prefixes

1. **Action Functions (Verb + Noun):**
   - Every function performing an action MUST start with a clear, active verb: `fetchUser()`, `calculateTax()`, `renderHelpRow()`, `validatePayload()`.
   - Ban vague garbage names: `handle()`, `process()`, `doStuff()`, `manage()`, `temp()`.
2. **Boolean Predicate Functions (`is` and `has` ONLY):**
   - Every function returning a boolean MUST begin with `is` or `has` ONLY: `isValid()`, `hasPermissions()`, `isExecutable()`, `isRetryRequired()`. All other prefixes (`can...`, `should...`, `was...`, `will...`, `did...`, `must...`) are **strictly BANNED**.
   - Negative prefixes (`isNotReady()`, `hasNoData()`) are **strictly prohibited**. Frame positively (`isReady()`, `hasData()`) and invert at the call site (`if !isReady { ... }`).

---

## 6. Phase 1 Violation Ledger Format

In Phase 1, you MUST generate `.lovable/plans/pending/XX-function-signatures-audit.md` containing the master inventory table:

```markdown
| Symbol / Call Site | File Path | Line | Category | Current Layout | Violation | Target Refactoring | Status |
|---|---|:---:|---|---|---|---|:---:|
| `SaveRecord` | `src/storage/db.go` | 42 | Definition | Single line (>2 params) | Rule 9a violation | Split to 1 parameter per line | PENDING |
| `logAction(...)` | `src/logger/log.ts` | 88 | Call Site | Single line (>2 args) | Rule 9b violation | Split to 1 argument per line | PENDING |
| `GetUser` | `src/services/user.go` | 104 | Return Type | `(*User, error)` | Multi-value return | Wrap in `Result[*User]` with `*AppError` | PENDING |
```

---

## Continuous 2-Phase Self-Loop & 2-Agent Concurrency Architecture

To guarantee full execution without stopping after planning mode, the master orchestrator MUST enforce this continuous 2-phase loop:

### 1. 2-Agent Concurrency & Strict `.lovable/` Bounding

- **2-Agent Limit (Max 2 Threads Each):** When dispatching work, spawn **at most 2 sub-agents concurrently**, with **no more than 2 threads per agent**.
- **Strict Folder Bounding (`.lovable/`):** Subagents can ONLY write planning files, subtasks, status reports, and logs inside `.lovable/` (`.lovable/plans/`, `.lovable/01-index.md`, `.lovable/memory/issues/`).
- **Context Diet:** Provide subagents with minimal instructions (e.g. "Read subtask file `.lovable/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md` and execute it"). Do not paste huge files into agent prompts.

### 2. Phase 1: Planning Mode & Subtask Generation (Steps 1 .. N/2)

- Spawn 2 planning subagents to scan the codebase for target guideline violations.
- Write the master architectural specification in `.lovable/plans/pending/xx-audit.md` with an exhaustive Violation Ledger table.
- Decompose the master plan into granular subtasks in `.lovable/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md`, `02-<subtask-title>.md`, etc.
- **MANDATORY AUTO-LOOP (DO NOT STOP):** Once Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for confirmation**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

### 3. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

- Spawn 2 execution subagents (max 2 threads each) to execute subtasks in parallel on disjoint files.
- Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, universal `*AppError` wrapping, Unix LF line endings).
- Move completed subtasks from `.lovable/plans/subtasks/` to `.lovable/plans/completed/` and update `.lovable/plans/01-index.md`.
- **Failure Memory & Feedback Loop:** If a subagent fails:
  - Rollback dirty working tree and log error details to `.lovable/plan.md` and `.lovable/memory/issues/xx-failure.md`.
  - The next subagent spawned MUST read the previous failure log first, record it as a pending memory task, and implement the necessary fix.
- Execute local linters and `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` ensuring `exit 0` before concluding.

---

## Strict In-Repository Execution & `.lovable/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.lovable/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.lovable/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.lovable/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`, `07-relative-path-fixer.py`, `05-naming-autofixer.py`, `09-cli-help-auditor.py`).
>    - RCA & Issue Logs: `.lovable/memory/issues/` and `.lovable/cicd-issues/`.
>    - Execution Plans & Subtasks: `.lovable/plans/pending/`, `.lovable/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.lovable/coding-guidelines.md`.
> 3. **Worker Pool & Log Aggregation Architecture:** All local runners and test orchestrators must use a concurrent worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without cancelling sibling workers, and print a consolidated final summary with full stdout/stderr error logs for failed jobs.
> 4. **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script from scratch**, regardless of whether the file already exists on disk.
> 5. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.

---

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/01-index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .lovable/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Automated Naming & Style Fixers:** Use `python 03-ai-scripts/08-naming-autofixer.py` and `05-guideline-autofixer.py` to audit boolean prefixes and newlines.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/01-index.md` using sequential script naming. For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

---

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Sub-agents are actively assigned disjoint files verified against `.lovable/01-index.md`.
- [ ] Completed tasks were `mv`'d to `plans/completed/` and `.lovable/plans/01-index.md` was updated.
- [ ] 3-strike rule respected: failed tasks cleanly rolled back and logged to `last-failure.md`.
- [ ] **Strict Relative Git Paths:** All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths (`/absolute/path/to/...`, `/absolute/path/to/...`) or `file:///` URIs.
- [ ] **Multi-Line Definitions (Rule 9a):** All function/method definitions with >2 parameters are formatted with exactly one parameter per line and trailing commas.
- [ ] **Multi-Line Invocations (Rule 9b):** All function/method call sites with >2 arguments are formatted with exactly one argument per line and trailing commas.
- [ ] **No Boolean Flag Parameters:** No boolean parameters used to switch behavior; split into distinct methods.
- [ ] **Semantic Naming:** All functions start with active verbs; all boolean functions start with is and has only (can, should, was, etc. are banned).
- [ ] **Single Return Types:** Multi-value `(T, error)` returns refactored to single `Result[T]` envelopes in services with complete predicate methods (`IsSuccess()`, `IsFailed()`, `HasError()`, `HasNoError()`, `HasValidError()`).
- [ ] **Universal `AppError`:** Zero generic `error` or `fmt.Errorf()` returns in domain logic.
- [ ] **LF Line Endings (`\n`):** All files use Unix LF line endings. Zero CRLF (`\r\n`).
- [ ] **UTF-8 Encoding (No BOM):** All files encoded in UTF-8 without BOM.
- [ ] **Single Trailing Newline:** Every file ends with exactly one terminating newline (`\n`).
- [ ] **Blank Line Before `if`:** Exactly one blank line precedes every `if` statement (unless at the top of a block).
- [ ] **Blank Line After `}`:** Exactly one blank line follows every closing brace `}` (unless closing the enclosing block).
- [ ] **Blank Line Before `return`:** Exactly one blank line precedes `return` / `throw` in multi-line blocks.
- [ ] **Zero Nested `if`:** All conditionals flattened to depth 0 using guard clauses and early returns.
- [ ] **Function Sizing:** All functions <= 8 lines preferred (hard cap 15 lines).
- [ ] `python linter-scripts/check-newline-styling.py` and `python linter-scripts/check-function-lengths.py` exited with code 0.
- [ ] Local CI runner `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` exited with code 0.

---

## Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Strict Relative Git Paths: All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths or `file:///` URIs.
- [ ] Master Guidelines: I have fully read and strictly enforced `spec/02-coding-guidelines/01-cross-language/01-index.md` and `.lovable/coding-guidelines.md`.
- [ ] Rule 9a/9b Multi-Line Formatting: Verified one parameter/argument per line for all definitions and call sites with >2 arguments.
- [ ] Result Envelope: Enforced `Result[T]` with `IsSuccess()`, `IsFailed()`, `HasError()`, `HasNoError()`, and `*AppError` across domain services.
- [ ] LF Line Endings & UTF-8 (No BOM): Verified Unix LF and UTF-8 across all files.
- [ ] Blank Line Before `if`: Verified blank line before every `if` statement across all modified files.
- [ ] Blank Line After `}`: Verified blank line after every closing brace `}` followed by code.
- [ ] Blank Line Before `return`: Verified blank line before every `return`/`throw` in multi-line blocks.
- [ ] Zero Nested `if`: Zero nested `if` statements (depth > 1).
- [ ] /learn the section as a /goal [AI Fix Scripts Memory](#ai-fix-scripts-memory)
- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

---

## Mandatory Linter & CI/CD Integration

1. **Linter Scripts:** `linter-scripts/check-function-lengths.py`, `linter-scripts/check-mws-error-codes.py`, `linter-scripts/check-newline-styling.py`
2. **Local Run Command:** `python linter-scripts/check-function-lengths.py`
3. **Autofixer Command:** `python 03-ai-scripts/05-guideline-autofixer.py <file>`
4. **CI/CD Integration (`.github/workflows/ci.yml`):**
   ```yaml
   - name: Validate Function Signatures & Error Envelopes
     run: |
       python linter-scripts/check-function-lengths.py
       python linter-scripts/check-mws-error-codes.py
       python linter-scripts/check-newline-styling.py
   ```
5. **Runner Registration (`03-ai-scripts/06-cicd-local-runner.py`):**
   ```python
   JOBS = {
       "Function Signatures Check": [sys.executable, "linter-scripts/check-function-lengths.py"],
       "Error Codes Check": [sys.executable, "linter-scripts/check-mws-error-codes.py"],
       "Newline Styling Check": [sys.executable, "linter-scripts/check-newline-styling.py"],
   }
   ```

---

## Strictly Avoid: No Automatic Releases & No Test Running (Strict Policy)

- **NO RELEASES (Strict Policy):** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
- **NO TEST RUNNING (Strict Policy):** Test execution is strictly disabled. You MUST NOT execute unit tests, integration tests, or test suites unless explicitly commanded by the repository owner.
- **CI/CD Quality Gates:** Execute local CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` ensuring all quality gates exit with code 0 (`exit 0`).
- **Test Inventory & Recent Changes Tracking:** Whenever any file is modified, append its repository-relative path to `.lovable/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`), cross-referencing `.lovable/test-inventory.json` so associated tests are known for future release verification.
