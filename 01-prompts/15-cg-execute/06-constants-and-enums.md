# Constants & Enums Architecture — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-enums`, `cg-constants`, `cg-execute enums`, `audit constants`, `fix enums`, `eliminate magic strings`, `eliminate magic numbers`, `enforce enum suffix`, `constants and enums audit`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, plan, refactor, and fix all constants, enums, magic string literals, raw character/rune literals (`rune(10)`), and magic number violations across the codebase, modifying source files directly to enforce the `*Type` enum suffix, extract magic numbers/strings into dedicated constant files, and use typed enums and traits until 100% green without stopping.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase to inventory all architectural violations and anti-patterns.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.lovable/plans/pending/` with an exhaustive Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.lovable/plans/subtasks/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/01-index.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and perform surgical refactoring following authoritative guidelines.
6. [ ] /goal Phase 2 (Step B): Enforce <= 8–15 line function decomposition, single return types, and clean formatting.
7. [ ] /goal Phase 2 (Step C): Execute local linters to verify 0 remaining violations across all modified files.
8. [ ] /goal Phase 2 (Step D): Execute local CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` with exit code 0 (`exit 0`).
9. [ ] /learn Ingest `.lovable/memory/01-index.md` for project memory index and past learnings.
10. [ ] /learn Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
11. [ ] /learn Ingest `spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
12. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for hallucination prevention and micro-tasking.
13. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for strict relative path citation requirements.
14. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/27-types-folder-convention.md` for domain-specific architectural specifications.
15. [ ] /learn Ingest `spec/02-coding-guidelines/` for domain-specific architectural specifications.
16. [ ] /learn Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.
17. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Codebase for Enum/Constant Violations, Write .lovable/plans/pending/ Spec, Create .lovable/plans/subtasks/, Verify/Create Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Actively Edit Code, Rename Enums with Type Suffix, Extract Constants, Update Call Sites, Verify Local CI)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Dedicated Section: Constants, Enums & Zero Magic Literals Architecture

Hardcoded string literals, raw character casts (`rune(10)`), and magic numbers degrade code maintainability, defeat type checking, and invite subtle runtime errors. Enums and constants must be strictly centralized, expressive, and strongly typed.

---

### 1. Mandatory `*Type` Suffix on All Enums

Every Enum type definition MUST end with `Type` (e.g. `UserRoleType`, `ExitCodeType`, `SeverityLevelType`, `OrderStatusType`, `HttpMethodType`).

- Go: `type UserRoleType string` or `type OrderStatusType byte`
- TypeScript: `enum UserRoleType` or `const UserRoleType = { ... } as const`
- PHP: `enum UserRoleType: string`
- Python: `class UserRoleType(StrEnum):`
- C#: `public enum UserRoleType`

---

### 2. Elimination of Magic Strings, Numbers, Runes & Delimiters

#### ❌ The Raw Character / Rune Anti-Pattern (User Issue Example)

Never use raw integer character codes or inline conversions like `string(rune(10))` or hardcoded delimiter strings across codebase logic:

```go
// ❌ FORBIDDEN: Raw character conversions and inline string literals
lines := strings.Split(string(data), string(rune(10))) // ❌ Magic rune literal
header := strings.Join(fields, ",")                    // ❌ Magic delimiter string

// ✅ REQUIRED: Centralized expressive constants
const (
    NewLineUnix = "\n"
    DelimiterComma = ","
)

lines := strings.Split(string(data), NewLineUnix)
header := strings.Join(fields, DelimiterComma)
```

---

### 3. The Logging & Test Assertion Exemption (What Is Allowed)

To avoid useless boilerplate, the following strings are **EXEMPT** from being extracted to constants:

1. **Informational Log Messages & Format Strings:**
   - `logger.Info("User successfully authenticated", "userId", userId)`
   - `fmt.Sprintf("processing item %d of %d", current, total)`
2. **Test Assertions & Error Descriptions in Test Files (`*_test.go`, `*.test.ts`):**
   - `t.Errorf("expected user to be active, got inactive")`
   - `expect(result).toBe("custom-test-value")`

> [!IMPORTANT]
> **What MUST ALWAYS be constants:**
>
> - Error codes / Error types (`E_INTERNAL_ERROR`, `AUTH_FAILED`).
> - Business entity statuses (`PENDING`, `ACTIVE`, `SUSPENDED`).
> - Protocol / API headers (`Authorization`, `Content-Type`, `X-Request-Id`).
> - Timeouts, retry counts, port numbers, buffer sizes, and pagination limits.
> - Delimiters, line breaks (`NewLineUnix`), and special encoding markers.

---

### 4. Language-Specific Examples & Best Practices

#### 4a. Go: Typed Enums, `iota`, and String Constants

```go
// ❌ FORBIDDEN: Untyped magic strings and missing Type suffix
const (
    RoleAdmin = "admin" // ❌ Missing UserRoleType type and Type suffix
    RoleUser  = "user"
)

func AssignRole(user *User, role string) { // ❌ Takes raw string instead of typed enum
    if role == "admin" {                   // ❌ Magic string comparison
        user.IsAdmin = true
    }
}

// ✅ REQUIRED: Typed enum with *Type suffix and dedicated constants package
package enums

type UserRoleType string

const (
    UserRoleTypeAdmin UserRoleType = "ADMIN"
    UserRoleTypeUser  UserRoleType = "USER"
    UserRoleTypeGuest UserRoleType = "GUEST"
)

func AssignRole(user *User, role UserRoleType) {
    if role == UserRoleTypeAdmin {
        user.IsAdmin = true
        return
    }

    user.IsAdmin = false
}
```

```go
// ✅ REQUIRED: Efficient Byte/Int iota enum with *Type suffix
package enums

type OrderStatusType byte

const (
    OrderStatusTypePending OrderStatusType = iota
    OrderStatusTypeProcessing
    OrderStatusTypeCompleted
    OrderStatusTypeFailed
)
```

---

#### 4b. TypeScript: Native Enums & `as const` Object Enums

String unions (`type Role = "admin" | "user"`) are **banned for enumerations** because they cannot be iterated, cannot be safely renamed with refactoring tools, and encourage magic strings at call sites.

```typescript
// ❌ FORBIDDEN: String union and magic string comparisons
type Role = "admin" | "editor" | "viewer"; // ❌ String union banned

function checkAccess(role: Role) {
    if (role === "admin") { // ❌ Magic string literal
        grantSuperuser();
    }
}

// ✅ REQUIRED Option 1: Native TypeScript Enum with *Type suffix
export enum UserRoleType {
    Admin = "ADMIN",
    Editor = "EDITOR",
    Viewer = "VIEWER",
}

function checkAccess(role: UserRoleType) {
    if (role === UserRoleType.Admin) {
        grantSuperuser();
    }
}

// ✅ REQUIRED Option 2: `as const` Object Enum (Bundle Size Optimized)
export const UserRoleType = {
    Admin: "ADMIN",
    Editor: "EDITOR",
    Viewer: "VIEWER",
} as const;

export type UserRoleType = (typeof UserRoleType)[keyof typeof UserRoleType];
```

---

#### 4c. PHP: Backed Enums & Trait Composition (PHP 8.1+)

```php
// ❌ FORBIDDEN: Magic strings and raw constants
class OrderService {
    public function process(string $status) {
        if ($status === 'completed') { // ❌ Magic string
            $this->notifyCustomer();
        }
    }
}

// ✅ REQUIRED: PHP Backed Enum with *Type suffix, methods, and reusable trait
namespace App\Traits;

trait HasEnumHelpers {
    public static function values(): array {
        return array_column(self::cases(), 'value');
    }

    public static function names(): array {
        return array_column(self::cases(), 'name');
    }

    public static function isValid(string $value): bool {
        return in_array($value, self::values(), true);
    }
}

namespace App\Enums;

use App\Traits\HasEnumHelpers;

enum OrderStatusType: string {
    use HasEnumHelpers;

    case Pending = 'pending';
    case Processing = 'processing';
    case Completed = 'completed';
    case Cancelled = 'cancelled';

    public function label(): string {
        return match($this) {
            self::Pending => 'Pending Review',
            self::Processing => 'In Processing',
            self::Completed => 'Order Completed',
            self::Cancelled => 'Order Cancelled',
        };
    }
}

class OrderService {
    public function process(OrderStatusType $status): void {
        if ($status === OrderStatusType::Completed) {
            $this->notifyCustomer();
        }
    }
}
```

---

#### 4d. Python: `StrEnum` & `IntEnum`

```python
# ❌ FORBIDDEN: Magic string comparisons
def process_task(priority: str) -> None:
    if priority == "high": # ❌ Magic string literal
        execute_urgent()

# ✅ REQUIRED: StrEnum with *Type suffix
from enum import StrEnum

class TaskPriorityType(StrEnum):
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"
    CRITICAL = "critical"

def process_task(priority: TaskPriorityType) -> None:
    if priority == TaskPriorityType.HIGH:
        execute_urgent()
```

---

#### 4e. C#: Typed Enums & Strongly-Typed Constants

```csharp
// ❌ FORBIDDEN: Raw string comparison
if (user.Role == "Admin") { ... }

// ✅ REQUIRED: C# Enum with *Type suffix
public enum UserRoleType
{
    Admin,
    Editor,
    Viewer,
}

public static class AppConstants
{
    public const int DefaultTimeoutSeconds = 30;
    public const string DefaultContentType = "application/json";
}
```

---

## Dedicated Definition Files Structure

All enums and constants MUST live in dedicated directories:

- **Go:** `enums/` and `constants/` (e.g. `enums/user_role_type.go`, `constants/http_constants.go`).
- **TypeScript:** `src/types/` and `src/enums/` (e.g. `src/enums/user-role-type.ts`).
- **PHP:** `src/Enums/` and `src/Constants/` (e.g. `src/Enums/OrderStatusType.php`).
- **Python:** `enums/` or `constants.py` (e.g. `enums/user_role_type.py`).

Never define enums or constants inline inside controllers, services, or business components.

---

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

## Strict In-Repository Execution & `.lovable/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.lovable/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.lovable/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.lovable/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`, `07-relative-path-fixer.py`, `05-naming-autofixer.py`).
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
- [ ] **All Enums End with `Type`:** Verified `*Type` suffix on every enum definition.
- [ ] **Zero Magic Strings / Numbers / Runes:** All status literals, character codes (`rune(10)`), and numbers extracted to constants.
- [ ] **No String Unions for Enums:** TypeScript uses native enums or `as const` object enums.
- [ ] **Dedicated Definition Files:** All enums/constants housed in dedicated modules.
- [ ] Coding Guidelines & Master Consolidated File: I have fully read, checked, and strictly enforced every file in `spec/02-coding-guidelines/`, as well as the master consolidated coding guideline file at `.lovable/coding-guidelines.md`.

1. [ ] /learn and apply as a /goal `.lovable/coding-guidelines.md` and also make sure the agent rules are created in the repo to read in the future quickly.

- [ ] `python linter-scripts/check-enum-guidelines.py` and `node linter-scripts/check-enum-and-boolean.mjs` exited with code 0.
- [ ] Local CI runner `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` exited with code 0.

---

## Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Strict Relative Git Paths: All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths or `file:///` URIs.
- [ ] Master Guidelines: I have fully read and strictly enforced `spec/02-coding-guidelines/01-cross-language/27-types-folder-convention.md` and `.lovable/coding-guidelines.md`.
- [ ] Enum Suffix: Every enum ends with `Type`.
- [ ] Zero Magic Literals: Zero raw magic strings or magic numbers at business logic call sites.

1. [ ] /learn the section as a /goal [AI Fix Scripts Memory](#ai-fix-scripts-memory)

- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

---

## Mandatory Linter & CI/CD Integration

1. **Linter Scripts:** `linter-scripts/check-enum-guidelines.py`, `linter-scripts/check-enum-and-boolean.mjs`
2. **Local Run Command:** `node linter-scripts/check-enum-and-boolean.mjs`
3. **Autofixer Command:** `python 03-ai-scripts/05-guideline-autofixer.py <file>`
4. **CI/CD Integration (`.github/workflows/ci.yml`):**
   ```yaml
   - name: Validate Enums & Constants
     run: |
       python linter-scripts/check-enum-guidelines.py
       node linter-scripts/check-enum-and-boolean.mjs
   ```
5. **Runner Registration (`03-ai-scripts/06-cicd-local-runner.py`):**
   ```python
   JOBS = {
       "Enum Guidelines Check": [sys.executable, "linter-scripts/check-enum-guidelines.py"],
   }
   ```

---

## Strictly Avoid: No Automatic Releases & No Test Running (Strict Policy)

- **NO RELEASES (Strict Policy):** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
- **NO TEST RUNNING (Strict Policy):** Test execution is strictly disabled. You MUST NOT execute unit tests, integration tests, or test suites unless explicitly commanded by the repository owner.
- **CI/CD Quality Gates:** Execute local CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` ensuring all quality gates exit with code 0 (`exit 0`).
- **Test Inventory & Recent Changes Tracking:** Whenever any file is modified, append its repository-relative path to `.lovable/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`), cross-referencing `.lovable/test-inventory.json` so associated tests are known for future release verification.
