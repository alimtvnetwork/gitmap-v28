# Multi-Language Enums, Traits & Pattern Matching — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-enums-traits`, `cg-enums`, `cg-execute enums`, `audit enums`, `php enums traits`, `rust enums`, `golang enums`, `multi-language enums`, `pattern matching audit`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, plan, refactor, and fix all Enum, Trait, and Pattern Matching architectures across Go, TypeScript, PHP, Rust, and Python codebases, enforcing string-backed enums, `HasEnumHelpers` traits, exhaustive pattern matching, `*Type` suffixes, central `enums/` folder collocation, and strict relative Git paths until 100% green without stopping.

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
14. [ ] /learn Ingest `spec/17-consolidated-guidelines/07-enum-standards.md` for cross-language enum and constant architectures.
15. [ ] /learn Ingest `spec/02-coding-guidelines/04-php/` for PHP 8.1+ backed enums and helper traits.
16. [ ] /learn Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.
17. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Enums & Traits, Build Violation Ledger in .lovable/plans/pending/, Subtasks)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Refactor Enums & Traits, Match Expressions, Run Typecheckers, Verify Local CI)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Dedicated Section: Multi-Language Enum & Trait Architecture

Enums represent a finite, closed set of valid domain values. Scattering loose string literals or bare integers destroys type safety. Every language must enforce strongly-typed enums with helper methods and pattern matching.

---

### 1. PHP 8.1+ String-Backed Enums & `HasEnumHelpers` Trait

In PHP 8.1+, all enums must be string-backed (`enum StatusType: string`) and use a standard `HasEnumHelpers` trait to provide `values()`, `names()`, `isValid()`, and `tryFromOrThrow()`.

```php
<?php

declare(strict_types=1);

namespace App\Enums;

use App\Exceptions\AppException;
use App\Enums\Traits\HasEnumHelpers;

enum OrderStatusType: string
{
    use HasEnumHelpers;

    case Pending    = 'pending';
    case Processing = 'processing';
    case Completed  = 'completed';
    case Cancelled  = 'cancelled';

    public function label(): string
    {
        return match ($this) {
            self::Pending    => 'Pending Payment',
            self::Processing => 'Processing Shipment',
            self::Completed  => 'Order Completed',
            self::Cancelled  => 'Order Cancelled',
        };
    }

    public function isTerminal(): bool
    {
        return match ($this) {
            self::Completed, self::Cancelled => true,
            self::Pending, self::Processing  => false,
        };
    }
}
```

```php
<?php

declare(strict_types=1);

namespace App\Enums\Traits;

use App\Exceptions\AppException;

trait HasEnumHelpers
{
    /** @return list<string> */
    public static function values(): array
    {
        return array_column(self::cases(), 'value');
    }

    /** @return list<string> */
    public static function names(): array
    {
        return array_column(self::cases(), 'name');
    }

    public static function isValid(string $value): bool
    {
        return self::tryFrom($value) !== null;
    }

    public static function fromOrThrow(string $value): static
    {
        $case = self::tryFrom($value);

        if ($case === null) {
            throw AppException::validation("Invalid enum value '$value' for " . static::class);
        }

        return $case;
    }
}
```

---

### 2. Rust Enums, Algebraic Data Types & Exhaustive Matching

In Rust, enums are first-class Algebraic Data Types. Use typed variants with payload data and exhaustive `match` branches:

```rust
// ✅ REQUIRED: Rust ADT Enum with *Type suffix and custom helper methods
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum TaskStatusType {
    Pending,
    Running { progress_pct: u8 },
    Completed { duration_ms: u64 },
    Failed { error_code: String, message: String },
}

impl TaskStatusType {
    pub fn is_finished(&self) -> bool {
        match self {
            Self::Completed { .. } | Self::Failed { .. } => true,
            Self::Pending | Self::Running { .. } => false,
        }
    }

    pub fn status_name(&self) -> &'static str {
        match self {
            Self::Pending => "PENDING",
            Self::Running { .. } => "RUNNING",
            Self::Completed { .. } => "COMPLETED",
            Self::Failed { .. } => "FAILED",
        }
    }
}
```

---

### 3. Go Type-Safe Multi-File Enums & `BasicEnum` Integration

In Go, enums MUST be scaffolded using `03-ai-scripts/30-enum-generator.py` into dedicated, self-contained package folders (`04-code/golang/pkg/enum/{name.lower()}type/` or `internal/enums/{name.lower()}type/`). Each enum is decomposed into a strict **4-file architecture**:

1. **`variant.go`**: Core type declaration (`Variant byte`), type alias (`type OrderStatusType = Variant`), `iota` constants starting with `Invalid Variant = iota`, compile-time interface assertions, value/code methods, affirmative item predicates (`IsPending()`, `IsCompleted()`), formatting methods, and DRY JSON serialization.
2. **`vars.go`**: Canonical `Result` alias (`type Result = result.Wrap[Variant]`), `variantLabels` array, `baseenumer.NewBasicInteger` (or `NewBasicString` / `NewBasicSparseInteger`) integration, `All()`, `Values()`, and monadic `Parse(s string) Result`.
3. **`variant_test.go`**: 100% test coverage validating interfaces, properties, item predicates, parse lookups, and JSON roundtrips.
4. **`readme.md`**: Package documentation detailing zero circular dependencies, direct Result returns, high-speed lookups, and DRY JSON marshaling.

#### Automated Scaffolding Command:
```bash
python 03-ai-scripts/30-enum-generator.py --name OrderStatus --type byte --items Pending,Processing,Completed,Cancelled --zero-value Invalid
```

#### `variant.go` (Core Type, Constants, & Methods):
```go
package orderstatustype

import (
    "encoding/json"
    "fmt"

    "coding-guidelines/common/pkg/baseenumer"
)

type (
    Variant byte

    OrderStatusType = Variant

    VariantPredicate func(v Variant) bool
)

const (
    Invalid Variant = iota // zero value is always Invalid
    Pending
    Processing
    Completed
    Cancelled
)

var (
    _ baseenumer.BaseEnumer   = Variant(0)
    _ baseenumer.ByteEnumer   = Variant(0)
    _ baseenumer.NumberEnumer = Variant(0)
    _ json.Marshaler          = Variant(0)
    _ json.Unmarshaler        = (*Variant)(nil)
)

func (v Variant) Byte() byte       { return byte(v) }
func (v Variant) ValueByte() byte  { return byte(v) }
func (v Variant) Bytes() []byte    { return []byte{byte(v)} }
func (v Variant) Int() int         { return int(v) }
func (v Variant) Code() uint16     { return uint16(v) }

func (v Variant) IsValid() bool   { return baseenumer.IsBetween(v, Pending, Cancelled) }
func (v Variant) IsInvalid() bool { return baseenumer.IsNotBetween(v, Pending, Cancelled) }
func (v Variant) IsEnum() bool    { return v.IsValid() }

func (v Variant) IsPending() bool    { return v == Pending }
func (v Variant) IsProcessing() bool { return v == Processing }
func (v Variant) IsCompleted() bool  { return v == Completed }
func (v Variant) IsCancelled() bool  { return v == Cancelled }

func (v Variant) Name() string {
    if int(v) < len(variantLabels) {
        return variantLabels[v]
    }
    return fmt.Sprintf("OrderStatus(%d)", byte(v))
}

func (v Variant) Label() string       { return v.Name() }
func (v Variant) String() string      { return v.Name() }
func (v Variant) ValueString() string { return baseenumer.FormatNameValue(v.Name(), byte(v)) }

func (v Variant) MarshalJSON() ([]byte, error) {
    return baseenumer.MarshalJSON(v.Name())
}

func (v *Variant) UnmarshalJSON(data []byte) error {
    return basicEnum.UnmarshalJSON(data, v)
}
```

#### `vars.go` (`BasicEnum` Engine & Monadic Parser):
```go
package orderstatustype

import (
    "coding-guidelines/common/pkg/baseenumer"
    "coding-guidelines/common/pkg/errtype"
    "coding-guidelines/common/pkg/result"
)

type Result = result.Wrap[Variant]

var (
    variantLabels = [...]string{
        Invalid:    "Invalid",
        Pending:    "Pending",
        Processing: "Processing",
        Completed:  "Completed",
        Cancelled:  "Cancelled",
    }

    basicEnum  = baseenumer.NewBasicInteger(variantLabels[:], Invalid)
    variantMap = basicEnum.Map()
)

func All() []Variant {
    return basicEnum.All()
}

func Values() []string {
    return basicEnum.Values()
}

func Parse(s string) Result {
    v, err := basicEnum.Parse(s)
    if err != nil {
        return result.WrapFailureWithId[Variant](errtype.Validation, err.Error())
    }

    return result.WrapSuccess(v)
}
```

---

### 4. TypeScript `as const` Object Enums

```typescript
export const OrderStatusType = {
    Pending: 'pending',
    Processing: 'processing',
    Completed: 'completed',
    Cancelled: 'cancelled',
} as const;

export type OrderStatusType = (typeof OrderStatusType)[keyof typeof OrderStatusType];

export function isTerminalStatus(status: OrderStatusType): boolean {
    return status === OrderStatusType.Completed || status === OrderStatusType.Cancelled;
}
```

---

## 5. Phase 1 Violation Ledger Format

In Phase 1, you MUST generate `.lovable/plans/pending/XX-enums-and-traits-audit.md` containing the master inventory table:

```markdown
| Target File | Line | Identifier | Current Pattern | Language | Planned Refactoring | Status |
|---|:---:|---|---|---|---|:---:|
| `app/Models/Order.php` | 24 | `$status` | Loose string literal `'pending'` | PHP | Backed Enum `OrderStatusType` + `HasEnumHelpers` | PENDING |
| `src/task.rs` | 52 | `status_code: u8` | Numeric status code `0, 1, 2` | Rust | ADT Enum `TaskStatusType` with payload | PENDING |
| `pkg/api/order.go` | 18 | `Status string` | Raw unvalidated string | Go | Dedicated `orderstatustype.Variant` + `BasicEnum` | PENDING |
```

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
- [ ] **Strict 03-ai-scripts/ Tooling Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
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
- [ ] **`*Type` Suffix:** All enum type definitions end with `Type` (e.g., `OrderStatusType`).
- [ ] **Exhaustive Matching:** All `match`/`switch` expressions cover 100% of enum cases without unhandled branches.
- [ ] **PHP Backed Enums:** All PHP enums use string-backing and `HasEnumHelpers` trait.
- [ ] **Rust ADT Enums:** All Rust variants implement proper `match` arms and derive macros.
- [ ] **LF Line Endings (`\n`):** All files use Unix LF line endings. Zero CRLF (`\r\n`).
- [ ] **UTF-8 Encoding (No BOM):** All files encoded in UTF-8 without BOM.
- [ ] **Single Trailing Newline:** Every file ends with exactly one terminating newline (`\n`).
- [ ] **Blank Line Before `if`:** Exactly one blank line precedes every `if` statement (unless at the top of a block).
- [ ] **Blank Line After `}`:** Exactly one blank line follows every closing brace `}` (unless closing the enclosing block).
- [ ] **Blank Line Before `return`:** Exactly one blank line precedes `return` / `throw` in multi-line blocks.
- [ ] **Zero Nested `if`:** All conditionals flattened to depth 0 using guard clauses and early returns.
- [ ] **Function Sizing:** All functions <= 8 lines preferred (hard cap 15 lines).
- [ ] Local CI runner `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` exited with code 0.

---

## Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Strict Relative Git Paths: All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths or `file:///` URIs.
- [ ] Master Guidelines: I have fully read and strictly enforced `spec/17-consolidated-guidelines/07-enum-standards.md` and `.lovable/coding-guidelines.md`.
- [ ] Enum Suffix: Enforced `*Type` naming across all languages.
- [ ] LF Line Endings & UTF-8 (No BOM): Verified Unix LF and UTF-8 across all files.
- [ ] Blank Line Before `if`: Verified blank line before every `if` statement across all modified files.
- [ ] Blank Line After `}`: Verified blank line after every closing brace `}` followed by code.
- [ ] Blank Line Before `return`: Verified blank line before every `return`/`throw` in multi-line blocks.
- [ ] Zero Nested `if`: Zero nested `if` statements (depth > 1).

1. [ ] /learn the section as a /goal [AI Fix Scripts Catalog](03-ai-scripts/01-index.md)

- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

---

## Strictly Avoid: No Automatic Releases & No Test Running (Strict Policy)

- **NO RELEASES (Strict Policy):** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
- **NO TEST RUNNING (Strict Policy):** Test execution is strictly disabled. You MUST NOT execute unit tests, integration tests, or test suites unless explicitly commanded by the repository owner.
- **CI/CD Quality Gates:** Execute local CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` ensuring all quality gates exit with code 0 (`exit 0`).
- **Test Inventory & Recent Changes Tracking:** Whenever any file is modified, append its repository-relative path to `.lovable/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`), cross-referencing `.lovable/test-inventory.json` so associated tests are known for future release verification.
