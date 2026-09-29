[/goal](slashCommand:goal) Autonomously scan, discover, plan, refactor, and format all function signatures across the codebase, enforcing argument reduction via dedicated value-based parameter Structs/DTOs for signatures with >2–3 parameters, affirmative boolean prefixing (is and has only (can, should, was, etc. are banned)) on all struct fields, mandatory `*appfault.AppError` returns (eliminating bare "void" functions in Go), wrapping external framework errors into `*AppError`, and single `Result[T]` return envelopes until 100% green without stopping with strict no-build and no-test execution (NEVER run build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine execution turns; all compilation and testing are strictly verified later in CI/CD). Spawn autonomous subagents (A = 2, H = 2) for parallel reading and modular spec generation, use GitMap high-speed commands as primary, establish a single-agent blueprint in Phase 1 (first 50% steps budget), and execute disjoint refactoring in Phase 2 (remaining 50% steps budget) with continuous self-looping until 100% complete and finalized with an atomic push.

[/learn](slashCommand:learn) Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Ingest and execute the bottom instructions ahead of all default conventions above, internalize the single-agent architectural blueprint, and persist all progress into `.ai-memory/plans/` and memory logs.

> **Bottom-Instruction Priority Mandate (Below Precedence / Suffix Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or guidelines above. The agent MUST inspect and follow the instructions below with absolute precedence.

```text
N = 200 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
```

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Planning, Detailed Spec, and Lean Subtask Generation)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Execution, Self-Looping, Targeted Quality Linting)
```

N, A, H, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

#### High-Speed GitMap Acceleration Options (Run Everything Faster)

Always prefer native GitMap commands over slow generic shell pipelines:
1. **Ultra-Fast File & Directory Discovery (AUM Index & Walk):**
   - **Wildcard / Glob Search:** `gitmap find "<wildcard*>" [-ext <ext>]` (alias `gitmap f`)
   - **Exact Filename Search:** `gitmap find-files <name> [-ext <ext>]` (alias `gitmap ff`)
   - **Substring Filename Search:** `gitmap find-files-any "<str>" [-ext <ext>]` (alias `gitmap ffa`)
   - **Prefix / Suffix Search:** `gitmap find-files-startswith <prefix>` (`gitmap ffs`) / `gitmap find-files-endswith <suffix>` (`gitmap ffe`)
   - **List Indexed Repo Files:** `gitmap list-files [pattern] [-ext <ext>]` (alias `gitmap lf`)
   - **Directory Tree & Scaffolding:** `gitmap folder-tree` (alias `gitmap ft`)
   - **Zero-Write File Stream:** `gitmap cat <filepath>`
   - **Instant Multi-Core Regex Search:** `gitmap search "<term>"` or `gitmap aum search "<query>" [dir] --ext <ext>`
2. **Fast Repository Hygiene, Lowercase & Symlink Repair:**
   - **Auto-Lowercase Files (Safe 2-Step `git mv`):** `gitmap lowercase` (alias `gitmap lcf [--dry-run]`)
   - **Lowercase Root Readme:** `gitmap lowercase-readme`
   - **Sync Curated `.gitignore` / `.gitattributes` / `.prettierignore`:** `gitmap commons` (alias `gitmap co` or `gitmap sync all`)
   - **Repair Broken Symlinks:** `gitmap fix-link` (alias `gitmap fixlink`)
   - **Clean Update Temp & Inspect Storage:** `gitmap update-cleanup`, `gitmap storage` (alias `gitmap stor`)
3. **Fast Git State, Execution & Atomic Commits:**
   - **Repo Status & Remote Check:** `gitmap status` (`gitmap st`), `gitmap has-any-updates` (`gitmap hau`), `gitmap latest-branch` (`gitmap lb`)
   - **Fast Cross-Platform Shell Runner:** `gitmap pwsh "<command>"` (`gitmap ps`), `gitmap bash "<command>"` (`gitmap sh`), `gitmap async <cmd>` (`gitmap asyn`)
   - **Semantic Atomic Commit & Push:** `gitmap cpf "<summary>"` (Feature), `gitmap cpb "<summary>"` (Bug), `gitmap cpr "<summary>"` (Release), `gitmap pcp "<summary>"` (Pull-Commit-Push)
   - **Smart CI/CD Pipeline Waiting:** `gitmap pe`, `gitmap pipeline-ai status --json` (`gitmap pl-ai status -t <etaSeconds>`)

---

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] [/goal](slashCommand:goal) Phase 1 (Step A): Deeply scan the target codebase using the GitMap AUM discovery tools (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary, with fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) as fallback, to inventory all architectural violations and anti-patterns without truncation.
2. [ ] [/goal](slashCommand:goal) Phase 1 (Step B): Write the master audit specification in `.ai-memory/plans/pending/XX-argument-reduction-audit.md` with an exhaustive Parameter & Return Ledger table.
3. [ ] [/goal](slashCommand:goal) Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.ai-memory/plans/subtasks/XX-argument-reduction/`.
4. [ ] [/goal](slashCommand:goal) Phase 1 (Step D): Verify or create the automated parameter linter and register in `03-ai-scripts/readme.md`.
5. [ ] [/goal](slashCommand:goal) Phase 2 (Step A): Refactor multi-argument functions (>2–3 params) by encapsulating parameters into dedicated value-based Structs (`TrackResultParams`, `CloneOptions`) or parameter objects.
6. [ ] [/goal](slashCommand:goal) Phase 2 (Step B): Enforce strict boolean prefixes (is and has only (can, should, was, etc. are banned)) on all struct fields and queued tasks (e.g. `safePull` -> `isSafePull`).
7. [ ] [/goal](slashCommand:goal) Phase 2 (Step C): Eliminate all bare "void" functions in Go domain/service logic by mandating `*appfault.AppError` returns for side-effect operations and `Result[T]` for data operations.
8. [ ] [/goal](slashCommand:goal) Phase 2 (Step D): Convert all external/framework standard `error` returns to `*appfault.AppError` context wrappers (`appfault.WrapSimple(err, caller)`).
9. [ ] [/goal](slashCommand:goal) Phase 2 (Step E): Execute local linters (`python linter-scripts/check-function-lengths.py`, `check-newline-styling.py`) to verify 0 remaining violations.
10. [ ] [/goal](slashCommand:goal) Phase 2 (Step F): Execute targeted file-level linters and verification on modified files ensuring 0 remaining violations (`exit 0`). DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine coding guideline execution turns.
11. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/memory/readme.md` for project memory index and past learnings.
12. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
13. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
14. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for Rule 9a/9b multi-line parameter and call formatting.
15. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/01-cross-language/10-function-naming.md` for semantic verb and predicate prefix standards.
16. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for hallucination prevention and micro-tasking.
17. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for strict relative path citation requirements.
18. [ ] [/learn](slashCommand:learn) Ingest `02-spec/03-error-manage/readme.md` for universal AppError wrapping and error envelopes.
19. [ ] [/learn](slashCommand:learn) Ingest `02-spec/03-error-manage/readme.md` for Result[T] and standardized API envelopes.
20. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.
21. [ ] [/goal](slashCommand:goal) Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Loose Parameters & Void Functions, Build Violation Ledger, Subtasks)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Struct-Based Refactoring, AppError Returns, Local CI Verification)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Dedicated Section: Function Argument Reduction & Return Architecture

Long parameter lists obscure function contracts, make call sites fragile, and increase cyclomatic complexity. A clean architecture reduces parameter counts using structured parameter objects and eliminates silent failures by enforcing mandatory error returns.

---

### 1. Function Argument Reduction via Parameter Structs

When a function requires **more than 2–3 parameters**, do NOT pass them as loose arguments. Group them into a dedicated, strongly-typed **Parameter Struct** (or Options Object).

#### ❌ FORBIDDEN (Loose multi-parameter signature):

```go
// Go: 5 loose parameters, unformatted boolean 'safePull', bare void return
func trackResult(
    p *Progress,
    result model.CloneResult,
    rec model.ScanRecord,
    targetDir string,
    safePull bool,
) {
    // ...
}
```

#### ✅ REQUIRED (Value-based parameter struct + mandatory *AppError return):

```go
package cloner

import (
    "pkg/appfault"
    "gitmap/model"
)

// -----------------------------------------------------------------------------
// In types.go:
// -----------------------------------------------------------------------------
// TrackResultParams encapsulates all inputs required for tracking a clone/pull result.
type TrackResultParams struct {
    Progress   *Progress
    Result     model.CloneResult
    ScanRecord model.ScanRecord
    TargetDir  string
    IsSafePull bool // Note: strict affirmative boolean prefix!
}

// TrackResult updates progress based on clone/pull outcome and returns any processing error.
func TrackResult(params TrackResultParams) *appfault.AppError {
    if params.Progress == nil {
        return appfault.New(
            appfault.ErrCodeValidationFailed,
            "progress tracker cannot be nil",
            "TrackResult",
        )
    }

    if params.Result.IsSuccess {
        pulled := params.IsSafePull && isGitRepo(params.TargetDir)
        params.Progress.Done(params.Result, pulled)
    }

    return nil
}
```

---

### 2. Value-Based vs Pointer-Based Structs in Go

In Go, parameter structs MUST be passed as **value types** (`params TrackResultParams`) by default:

1. **Value-Based Structs (Default):**
   - Eliminates `nil` pointer panics at the call site.
   - Communicates immutability and data encapsulation.
   - Lightweight and cache-friendly for standard parameter objects (< 1KB).
2. **Pointer-Based Structs (Only when required):**
   - Use pointers (`params *TrackResultParams`) ONLY when the function explicitly needs to mutate the caller's struct state or when holding large non-copyable buffers (`sync.Mutex`, large byte arrays).

---

### 3. Boolean Prefix Enforcement on Struct Fields

When grouping parameters into a struct, all boolean fields **MUST adhere strictly to affirmative prefixes**:

- ❌ `safePull bool` ➔ ✅ `IsSafePull bool`
- ❌ `force bool` ➔ ✅ `IsForce bool`
- ❌ `dryRun bool` ➔ ✅ `IsDryRun bool`
- ❌ `verbose bool` ➔ ✅ `IsVerbose bool`
- ❌ `skipCache bool` ➔ ✅ `IsSkipCache bool` (or `HasSkipCache bool`)

#### Queued Task Protocol for Legacy Callers

If a parameter or struct field cannot be immediately refactored across the entire codebase in a single turn without breaking external packages:

1. Formulate a **Queued Task** in `.ai-memory/plans/pending/XX-boolean-naming-queue.md`.
2. Record the exact symbol, file path, line number, and required affirmative replacement.
3. Schedule the subtask for sequential execution in Phase 2.

---

### 4. Mandatory Return Architecture in Go (Zero Bare "Void" Functions)

In Go, **99.99% of functions MUST have a return type**. Bare "void" functions (`func DoWork()`) that return nothing are strictly prohibited in domain, business logic, service, and utility layers.

#### 4a. Side-Effect & Mutation Functions (Return `*appfault.AppError`)

If a function performs an action, I/O operation, or state mutation that produces no return data, it **MUST return `*appfault.AppError`**:

```go
// ❌ FORBIDDEN: Bare void function swallows or ignores potential execution failures
func SaveConfig(cfg *Config) {
    data, _ := json.Marshal(cfg)
    os.WriteFile("config.json", data, 0644)
}

// ✅ REQUIRED: Returns *appfault.AppError with complete contextual wrapping
func SaveConfig(cfg *Config) *appfault.AppError {
    if cfg == nil {
        return appfault.New(
            appfault.ErrCodeValidationFailed,
            "configuration cannot be nil",
            "SaveConfig",
        )
    }

    data, marshalErr := json.Marshal(cfg)

    if marshalErr != nil {
        return appfault.WrapSimple(marshalErr, "SaveConfig.Marshal")
    }

    if writeErr := os.WriteFile("config.json", data, 0644); writeErr != nil {
        return appfault.WrapSimple(writeErr, "SaveConfig.WriteFile")
    }

    return nil
}
```

---

#### 4b. External & Framework Error Conversion

Whenever code calls standard library functions (`os.*`, `io.*`, `exec.*`, `json.*`) or third-party packages that return standard `error`:

1. **Never return standard `error` directly** from domain or service layers.
2. **Always convert and wrap immediately** into `*appfault.AppError` using `appfault.WrapSimple(err, caller)` or `appfault.New(ErrCode, msg, caller)`:

```go
// ✅ REQUIRED: Converting framework error to *appfault.AppError
cmd := exec.Command("git", "status")
output, cmdErr := cmd.CombinedOutput()

if cmdErr != nil {
    return appfault.WrapWithDetails(
        cmdErr,
        appfault.ErrCodeGitExecutionFailed,
        string(output),
        "ExecuteGitStatus",
    )
}
```

---

#### 4c. Data-Producing Functions (Return `Result[T]`)

If a function computes or retrieves data, return the single `Result[T]` envelope:

```go
// ✅ REQUIRED: Single Result[T] envelope return
func LoadConfig(path string) Result[*Config] {
    if path == "" {
        appErr := appfault.New(
            appfault.ErrCodeValidationFailed,
            "config path is required",
            "LoadConfig",
        )
        return FailureResult[*Config](appErr)
    }

    data, readErr := os.ReadFile(path)

    if readErr != nil {
        appErr := appfault.WrapSimple(readErr, "LoadConfig.ReadFile")
        return FailureResult[*Config](appErr)
    }

    var cfg Config
    if unmarshalErr := json.Unmarshal(data, &cfg); unmarshalErr != nil {
        appErr := appfault.WrapSimple(unmarshalErr, "LoadConfig.Unmarshal")
        return FailureResult[*Config](appErr)
    }

    return SuccessResult[*Config](&cfg)
}
```

---

### 5. Multi-Language Parameter Object Architecture

#### 5a. TypeScript Parameter Object (`interface *Options`)

```typescript
// ✅ REQUIRED: Options interface with readonly properties and affirmative booleans
export interface TrackResultOptions {
    readonly progress: ProgressTracker;
    readonly result: CloneResult;
    readonly scanRecord: ScanRecord;
    readonly targetDir: string;
    readonly isSafePull: boolean;
}

export function trackResult(options: TrackResultOptions): Result<void> {
    if (!options.progress) {
        return failureResult(new AppError(
            ErrorCodeType.ValidationFailed,
            "Progress tracker is required",
            "trackResult",
        ));
    }

    // ...
    return successResult(undefined);
}
```

---

#### 5b. PHP 8.1+ Readonly DTO Parameter Object

```php
<?php

declare(strict_types=1);

namespace App\Cloner;

use App\Common\Result;
use App\Common\Exceptions\AppException;

final readonly class TrackResultParams
{
    public function __construct(
        public ProgressTracker $progress,
        public CloneResult $result,
        public ScanRecord $scanRecord,
        public string $targetDir,
        public bool $isSafePull = false,
    ) {}
}

final class ClonerService
{
    public function trackResult(TrackResultParams $params): Result
    {
        // ...
        return Result::success(null);
    }
}
```

---

#### 5c. Rust Parameter Struct

```rust
pub struct TrackResultParams<'a> {
    pub progress: &'a mut ProgressTracker,
    pub result: CloneResult,
    pub scan_record: ScanRecord,
    pub target_dir: &'a Path,
    pub is_safe_pull: bool,
}

pub fn track_result(params: TrackResultParams) -> Result<(), AppError> {
    // ...
    Ok(())
}
```

---

#### 5d. Python Frozen Dataclass Parameter Object

```python
from dataclasses import dataclass
from pathlib import Path
from typing import Optional

@dataclass(frozen=True)
class TrackResultParams:
    progress: ProgressTracker
    result: CloneResult
    scan_record: ScanRecord
    target_dir: Path
    is_safe_pull: bool = False

def track_result(params: TrackResultParams) -> Result[None]:
    # ...
    return SuccessResult(None)
```

---

## 6. Phase 1 Violation Ledger Format

In Phase 1, you MUST generate `.ai-memory/plans/pending/XX-argument-reduction-audit.md` containing the master inventory table:

```markdown
| Symbol / Function | File Path | Line | Param Count | Current Signature | Violation | Target Refactoring | Status |
|---|---|:---:|:---:|---|---|---|:---:|
| `trackResult` | `gitmap/cloner/runners.go` | 118 | 5 | `(p, res, rec, dir, safePull)` | >3 loose params, bare void | Create `TrackResultParams`, return `*AppError` | PENDING |
| `dispatchTask` | `src/cluster/exec.go` | 64 | 4 | `(ctx, cmd, timeout, force)` | >3 loose params, `force` bool | Create `DispatchTaskParams`, `isForce` | PENDING |
| `cleanupTemp` | `src/storage/temp.go` | 210 | 1 | `(path string)` (void) | Bare void function | Return `*appfault.AppError` | PENDING |
```

---

## Continuous 2-Phase Self-Loop & 2-Agent Concurrency Architecture

To guarantee full execution without stopping after planning mode, the master orchestrator MUST enforce this continuous 2-phase loop:

### 1. 2-Agent Concurrency & Strict `.ai-memory/` Bounding

- **2-Agent Limit (Max 2 Threads Each):** When dispatching work, spawn **at most 2 sub-agents concurrently**, with **no more than 2 threads per agent**.
- **Strict Folder Bounding (`.ai-memory/`):** Subagents can ONLY write planning files, subtasks, status reports, and logs inside `.ai-memory/` (`.ai-memory/plans/`, `.ai-memory/readme.md`, `.ai-memory/memory/issues/`).
- **Context Diet:** Provide subagents with minimal instructions (e.g. "Read subtask file `.ai-memory/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md` and execute it"). Do not paste huge files into agent prompts.

### 2. Phase 1: Planning Mode & Subtask Generation (Steps 1 .. N/2)

### Fast File Discovery & Reading via Python Toolchain (Mandatory Acceleration)

To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, the AI agent MUST use the fast 2-tier discovery toolchain:

### Tier 1: GitMap AUM Acceleration (PRIMARY)
1. **Universal File Search:**
   ```bash
   gitmap find "<pattern>" [-ext <ext>]
   ```
2. **List Indexed Files & Substring Lookup:**
   ```bash
   gitmap list-files [pattern]
   gitmap find-files-any "<substring>"
   ```
3. **Stream File Content:**
   ```bash
   gitmap cat <filepath>
   ```
4. **Instant Code Walk Search:**
   ```bash
   gitmap search "<term>"
   ```

### Tier 2: Fast Cached Python Toolchain (FALLBACK)
1. **Inventory Target Files (with `--limit` option):**
   ```bash
   python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats
   ```
2. **Fast Cached Grep (<15ms, with `--limit` option):**
   ```bash
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "<search-pattern>" --lang go --limit 50
   ```
3. **Sub-Millisecond Folder & File Exploration (with `--limit` option):**
   ```bash
   python 03-ai-scripts/17-fast-file-reader.py --list-folder <folder-path> --ext .go --limit 50
   python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000
   python 03-ai-scripts/17-fast-file-reader.py --search-pattern "<pattern>" --path <folder-path> --limit 50
   ```
4. **Subsystem & Topology Overview:**
   ```bash
   python 03-ai-scripts/18-codebase-topology-discoverer.py --summary
   ```
Do not rely on standard search tools with 50-item truncation when discovering repository-wide violations.

- Spawn 2 planning subagents to scan the codebase for target guideline violations.
- Write the master architectural specification in `.ai-memory/plans/pending/xx-audit.md` with an exhaustive Violation Ledger table.
- Decompose the master plan into granular subtasks in `.ai-memory/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md`, `02-<subtask-title>.md`, etc.
- **MANDATORY AUTO-LOOP (DO NOT STOP):** Once Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for confirmation**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

### 3. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

- Spawn 2 execution subagents (max 2 threads each) to execute subtasks in parallel on disjoint files.
- Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, universal `*AppError` wrapping, Unix LF line endings).
- Move completed subtasks from `.ai-memory/plans/subtasks/` to `.ai-memory/plans/completed/` and update `.ai-memory/plans/readme.md`.
- **Failure Memory & Feedback Loop:** If a subagent fails:
  - Rollback dirty working tree and log error details to `.ai-memory/plan.md` and `.ai-memory/memory/issues/xx-failure.md`.
  - The next subagent spawned MUST read the previous failure log first, record it as a pending memory task, and implement the necessary fix.
- Execute targeted local linters on modified files ensuring `exit 0` before concluding. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine loops.

---

## Strict In-Repository Execution & `.ai-memory/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.ai-memory/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.ai-memory/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.ai-memory/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`, `07-relative-path-fixer.py`, `05-naming-autofixer.py`, `09-cli-help-auditor.py`).
>    - RCA & Issue Logs: `.ai-memory/memory/issues/` and `.ai-memory/cicd-issues/`.
>    - Execution Plans & Subtasks: `.ai-memory/plans/pending/`, `.ai-memory/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.ai-memory/coding-guidelines.md`.
> 3. **Worker Pool & Log Aggregation Architecture:** All local runners and test orchestrators must use a concurrent worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without canceling sibling workers, and print a consolidated final summary with full stdout/stderr error logs for failed jobs.
> 4. **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script from scratch**, regardless of whether the file already exists on disk.
> 5. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.

---

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] [/goal](slashCommand:goal) **Reuse First:** I have rigorously scanned and [/learn](slashCommand:learn)ed `03-ai-scripts/readme.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict 03-ai-scripts/ Tooling Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Automated Naming & Style Fixers:** Use `python 03-ai-scripts/08-naming-autofixer.py` and `05-guideline-autofixer.py` to audit boolean prefixes and newlines.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/readme.md` using sequential script naming. For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

---

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Sub-agents are actively assigned disjoint files verified against `.ai-memory/readme.md`.
- [ ] Completed tasks were `mv`'d to `plans/completed/` and `.ai-memory/plans/readme.md` was updated.
- [ ] 3-strike rule respected: failed tasks cleanly rolled back and logged to `last-failure.md`.
- [ ] **Strict Relative Git Paths:** All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths (`/absolute/path/to/...`, `/absolute/path/to/...`) or `file:///` URIs.
- [ ] **Argument Reduction via Structs:** All functions with >2–3 parameters encapsulated into value-based parameter structs (`*Params`).
- [ ] **Boolean Prefix Compliance:** All struct fields and boolean parameters use affirmative prefixes (is and has only (can, should, was, etc. are banned)).
- [ ] **Mandatory AppError Returns:** Zero bare "void" functions in Go domain/service logic; all side-effect functions return `*appfault.AppError`.
- [ ] **Framework Error Conversion:** All standard library / framework errors converted and wrapped into `*appfault.AppError`.
- [ ] **Single Return Types:** Multi-value `(T, error)` returns refactored to single `Result[T]` envelopes.
- [ ] **Multi-Line Formatting (Rule 9a/9b):** All definitions and call sites with >2 arguments formatted one argument per line with trailing commas.
- [ ] **LF Line Endings (`\n`):** All files use Unix LF line endings. Zero CRLF (`\r\n`).
- [ ] **UTF-8 Encoding (No BOM):** All files encoded in UTF-8 without BOM.
- [ ] **Single Trailing Newline:** Every file ends with exactly one terminating newline (`\n`).
- [ ] **Blank Line Before `if`:** Exactly one blank line precedes every `if` statement (unless at the top of a block).
- [ ] **Blank Line After `}`:** Exactly one blank line follows every closing brace `}` (unless closing the enclosing block).
- [ ] **Blank Line Before `return`:** Exactly one blank line precedes `return` / `throw` in multi-line blocks.
- [ ] **Zero Nested `if`:** All conditionals flattened to depth 0 using guard clauses and early returns.
- [ ] **Function Sizing:** All functions <= 8 lines preferred (hard cap 15 lines).
- [ ] `python linter-scripts/check-newline-styling.py` and `python linter-scripts/check-function-lengths.py` exited with code 0.
- [ ] **Atomic File Recording:** Modified files recorded to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`). Builds, tests, and CI runner are deferred to CI/CD fix.

---

## Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

[/goal](slashCommand:goal) You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Strict Relative Git Paths: All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths or `file:///` URIs.
- [ ] Master Guidelines: I have fully read and strictly enforced `02-spec/02-coding-guidelines/01-cross-language/readme.md` and `.ai-memory/coding-guidelines.md`.
- [ ] Parameter Structs: Enforced parameter structs (`*Params`) for all functions with >2–3 parameters.
- [ ] Zero Void in Go: Enforced `*appfault.AppError` returns for all side-effect functions in Go.
- [ ] Framework Error Conversion: Converted all stdlib/framework `error` instances to `*appfault.AppError`.
- [ ] Result Envelope: Enforced `Result[T]` with complete predicate methods across domain services.
- [ ] LF Line Endings & UTF-8 (No BOM): Verified Unix LF and UTF-8 across all files.
- [ ] Blank Line Before `if`: Verified blank line before every `if` statement across all modified files.
- [ ] Blank Line After `}`: Verified blank line after every closing brace `}` followed by code.
- [ ] Blank Line Before `return`: Verified blank line before every `return`/`throw` in multi-line blocks.
- [ ] Zero Nested `if`: Zero nested `if` statements (depth > 1).
- [ ] [/learn](slashCommand:learn) the section as a [/goal](slashCommand:goal) [AI Fix Scripts Catalog](03-ai-scripts/readme.md)
- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

---

## Mandatory Linter & CI/CD Integration

1. **Linter Scripts:** `linter-scripts/check-function-lengths.py`, `linter-scripts/check-mws-error-codes.py`, `linter-scripts/check-newline-styling.py`
2. **Local Run Command:** `python linter-scripts/check-function-lengths.py`
3. **Autofixer Command:** `python 03-ai-scripts/05-guideline-autofixer.py <file>`
4. **CI/CD Integration (`.github/workflows/ci.yml`):**
   ```yaml
   - name: Validate Function Arguments & AppError Envelopes
     run: |
       python linter-scripts/check-function-lengths.py
       python linter-scripts/check-mws-error-codes.py
       python linter-scripts/check-newline-styling.py
   ```
5. **Runner Registration (`03-ai-scripts/06-cicd-local-runner.py`):**
   ```python
   JOBS = {
       "Function Signatures & Params Check": [sys.executable, "linter-scripts/check-function-lengths.py"],
       "Error Codes Check": [sys.executable, "linter-scripts/check-mws-error-codes.py"],
       "Newline Styling Check": [sys.executable, "linter-scripts/check-newline-styling.py"],
   }
   ```

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

---

## Strictly Avoid: No Automatic Releases, No Test Running & No Full CI/CD Runner in Routine Turns (Strict Policy)

- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- **NO RELEASES (Strict Policy):** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
- **NO TEST RUNNING (Strict Policy):** Test execution is strictly disabled. You MUST NOT execute unit tests, integration tests, or test suites unless explicitly commanded by the repository owner.
- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns. Running the heavy 28-38 gate pipeline across the entire repository wastes massive amounts of time and scans unrelated files. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **Test Inventory & Recent Changes Tracking:** Whenever any file is modified, append its repository-relative path to `.ai-memory/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`), cross-referencing `.ai-memory/test-inventory.json` so associated tests are known for future release verification.

--

## 🚨 Highest Priority Instructions (Appended User Tasks & Instructions Below)

[PASTE USER REQUEST / TASK INSTRUCTIONS HERE — THE AGENT MUST EXECUTE WHATEVER IS WRITTEN BELOW WITH ABSOLUTE PRIORITY AND PRECEDENCE OVER ALL GENERAL GUIDELINES ABOVE]
