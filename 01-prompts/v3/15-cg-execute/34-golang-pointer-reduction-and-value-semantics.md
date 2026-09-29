[/goal](slashCommand:goal) Autonomously scan, plan, refactor, and migrate Golang codebases to value semantics, systematically reducing pointer variables (`*T`), heap escapes, and nil-dereference hazards across the codebase using a strict 300-step 3-phase self-loop (Steps 1–100: Pointer Escape Audit & Struct Inventory; Steps 101–200: Architecture Spec & Type Classification; Steps 201–300: Active Value-Semantic Refactoring & Receiver Migration) with strict no-build and no-test execution (compilation and testing verified later in CI/CD). Enforce value semantics for small structs (<= 64 bytes), pure query receivers, parameters, and constructors, while strictly preserving pointers only where technically necessary (`sync.Mutex`, large state buffers, tri-state nil, or `*appfault.AppError`), until 100% verified.

[/learn](slashCommand:learn) Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, target packages, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Ingest and execute the bottom instructions ahead of all default conventions above, internalize the 100/100/100 step allocation (Audit -> Spec -> Refactor), migrate methods to value receivers, and persist all progress into `.ai-memory/plans/` and memory logs.

> **Bottom-Instruction Priority Mandate (Below Precedence / Suffix Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or guidelines above. The agent MUST inspect and follow the instructions below with absolute precedence.

```text
N = 300 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
```

- **N = 300:** Total self-loop steps budget divided into three equal 100-step phases:
  - **PHASE_1_STEPS = 100 Steps (Steps 1 .. 100):** Deep Audit using GitMap AUM, Pointer Escape Analysis, Allocation Profiling, Candidate Struct Inventory, Sizing Classification (<= 64 bytes).
  - **PHASE_2_STEPS = 100 Steps (Steps 101 .. 200):** Architectural Classification, Value vs Pointer Semantics Rules, Struct Boundary Decoupling, Subtask Generation.
  - **PHASE_3_STEPS = 100 Steps (Steps 201 .. 300):** Surgical Refactoring, Method Receiver Migration, Return Value Simplification, Targeted Linting, Atomic GitMap Push.
- **A = 2:** Count of autonomous subagents running concurrently (`invoke_subagent` launches up to 2 subagents).
- **H = 2 (Hands / Parallel Operations):**
  1. **Workload Hands ($H_{batch} = 2$):** Each subagent handles a bounded batch of up to 2 disjoint subtasks from `.ai-memory/plans/subtasks/`.
  2. **Tool-Dispatch Hands ($H_{tool} = 2$):** Within any execution step, each agent or subagent executes up to 2 parallel tool calls in a single response turn.

```text
PHASE_1_STEPS = Steps 1 .. 100   (Audit & Plan: GitMap AUM Pointer Inventory, Struct Sizing, Escape Analysis)
PHASE_2_STEPS = Steps 101 .. 200 (Architectural Design & Type Selection: Spec in 02-spec/21-app/, Subtask Generation)
PHASE_3_STEPS = Steps 201 .. 300 (Execute, Refactor & Verify: Value Receivers, Small Struct Value Returns, Targeted Linters, GitMap Push)
```

N, A, H, PHASE_1_STEPS, PHASE_2_STEPS, and PHASE_3_STEPS are read-only after initialization. Never modify them mid-execution.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] [/goal](slashCommand:goal) Phase 1 (Step A): Deeply scan the target codebase using GitMap AUM discovery tools (`gitmap search "*T"`, `gitmap find "*.go"`, `gitmap lf`, `gitmap cat`) as primary, with fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) as fallback, to inventory all struct declarations, pointer parameters (`*T`), pointer receivers (`func (s *T)`), pointer return types, and slice pointers (`*[]T`).
2. [ ] [/goal](slashCommand:goal) Phase 1 (Step B): Calculate struct memory sizes across all packages, classifying structs into Value Candidates (<= 64 bytes) versus Pointer Retainers (> 64–128 bytes, or containing mutexes).
3. [ ] [/goal](slashCommand:goal) Phase 1 (Step C): Write the master audit specification in `.ai-memory/plans/pending/xx-pointer-reduction-audit.md` with an exhaustive Violation Ledger tracking every candidate pointer, receiver, and struct size.
4. [ ] [/goal](slashCommand:goal) Phase 1 (Step D): Decompose the master plan into granular, atomic subtasks in `.ai-memory/plans/subtasks/xx-pointer-reduction/01-<subtask>.md`, `02-<subtask>.md`, etc.
5. [ ] [/goal](slashCommand:goal) Phase 2 (Step A): Author the canonical architecture specification in `02-spec/21-app/xx-golang-value-semantics.md` detailing value receiver rules, parameter structs, and exception boundaries.
6. [ ] [/goal](slashCommand:goal) Phase 2 (Step B): Verify or create the automated quality linter hook (`linter-scripts/check-go-pointers.py` or `golangci-lint` configuration) and register in `03-ai-scripts/readme.md`.
7. [ ] [/goal](slashCommand:goal) Phase 2 (Step C): Define receiver migration contracts, ensuring that all read-only methods on immutable structs migrate to value receivers `func (s T)` without affecting state mutation logic.
8. [ ] [/goal](slashCommand:goal) Phase 3 (Step A): Spawn up to A = 2 execution subagents on disjoint packages to perform surgical refactoring: migrate small struct parameters and return types from `*T` to `T`.
9. [ ] [/goal](slashCommand:goal) Phase 3 (Step B): Migrate read-only query receivers to value receivers `func (s T)` for all immutable structs.
10. [ ] [/goal](slashCommand:goal) Phase 3 (Step C): Eliminate pointer-to-slice (`*[]T`) and pointer-to-map (`*map[K]V`) anti-patterns across all signatures, replacing with direct slice `[]T` and map `map[K]V` value types.
11. [ ] [/goal](slashCommand:goal) Phase 3 (Step D): Enforce <= 8–15 line function decomposition and clean vertical spacing on all modified files.
12. [ ] [/goal](slashCommand:goal) Phase 3 (Step E): Execute targeted file-level linters to verify 0 remaining violations across all modified files (`exit 0`).
13. [ ] [/goal](slashCommand:goal) Phase 3 (Step F): Atomically record modified files into `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`). DO NOT run full CI/CD runners (`06-cicd-local-runner.py`), build checks, or test suites during routine turns. Finalize with atomic commit via `gitmap cpf "<summary>"`.
14. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/memory/readme.md` for project memory index and past learnings.
15. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
16. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
17. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/03-type-safety-and-errors.md` for Go type safety and error rules.
18. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/04-database-and-structs.md` for Go struct design.
19. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.
20. [ ] [/goal](slashCommand:goal) Create or update agent rules in the repository if missing from agent memory.

---

## Dedicated Section: Golang Pointer Reduction & Value Semantics Architecture

In Go, overusing pointers (`*T`) is a widespread anti-pattern that creates unnecessary heap escapes, increases garbage collector (GC) pressure, introduces cache misses due to pointer indirection, and causes nil-pointer dereference panics.

### 1. The Value Semantics Mandate

Value semantics (`T` instead of `*T`) keep data on the stack, guarantee memory locality, enable deterministic cleanup without GC overhead, and mathematically eliminate `nil` dereference errors.

### 2. Pointer Justification Matrix: When to Use `*T` vs `T`

| Construct | Recommendation | Justification |
| :--- | :---: | :--- |
| **Small Structs (<= 64 bytes)** | `T` (Value) | Copying 1–4 words is faster than heap allocation and pointer dereference. |
| **Large Structs (> 64–128 bytes)** | `*T` (Pointer) | Copy overhead exceeds pointer cost; pass by reference. |
| **Read-Only / Query Methods** | `(s T)` Value Receiver | Pure functions with no state mutation; prevents accidental modification. |
| **Stateful Mutation Methods** | `(s *T)` Pointer Receiver | Method genuinely mutates fields on the caller's struct instance. |
| **Structs with Mutexes** | `*T` (Pointer) | `sync.Mutex`, `sync.RWMutex`, or atomic fields MUST NOT be copied. |
| **Slices & Maps** | `[]T` / `map[K]V` (Value) | Slices and maps are already pointer-backed header structs; never use `*[]T`. |
| **Channels** | `chan T` (Value) | Channels are already runtime pointers; never use `*chan T`. |
| **Tri-State Nilability** | `*T` (Pointer) | Distinguishing absent `nil` from empty zero-value `""` / `0`. |
| **Structured Errors** | `*appfault.AppError` | Mandatory exception to conform to Go error return type standards. |

### 3. Canonical Code Transformations

#### Transformation A: Small Parameter Structs & Constructors

```go
// ❌ ANTI-PATTERN: Heap-allocating small 32-byte struct and returning pointer
type FilterOptions struct {
    Limit  int
    Offset int
    Prefix string
}

func NewFilter(prefix string) *FilterOptions {
    return &FilterOptions{Prefix: prefix, Limit: 50}
}

// ✅ CANONICAL PATTERN: Direct value semantics, stack allocated
type FilterOptions struct {
    Limit  int
    Offset int
    Prefix string
}

func NewFilter(prefix string) FilterOptions {
    return FilterOptions{Prefix: prefix, Limit: 50}
}
```

#### Transformation B: Query Methods & Value Receivers

```go
// ❌ ANTI-PATTERN: Pointer receiver for pure read-only query
func (c *Config) IsValid() bool {
    return c.Port > 0 && c.Host != ""
}

// ✅ CANONICAL PATTERN: Value receiver for pure read-only inspection
func (c Config) IsValid() bool {
    hasPort := c.Port > 0
    hasHost := c.Host != ""
    return hasPort && hasHost
}
```

#### Transformation C: Slice & Map Parameter Hygiene

```go
// ❌ ANTI-PATTERN: Pointer to slice or pointer to map
func ProcessItems(items *[]string, lookup *map[string]int) { ... }

// ✅ CANONICAL PATTERN: Direct slice and map values
func ProcessItems(items []string, lookup map[string]int) { ... }
```

---

### Independent AI Audit & Verification Instructions

You are an Independent AI Verification and Quality Auditor.
Your task is to independently audit, verify, and remediate the implementation against the canonical specification and verbatim requirements.

#### 1. Target Documents & Implemented Code:
- **Canonical Spec & Verbatim Requirements:** [02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/03-type-safety-and-errors.md](02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/03-type-safety-and-errors.md)
- **Database & Structs:** [02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/04-database-and-structs.md](02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/04-database-and-structs.md)
- **Consolidated Plan & Subtasks:** `.ai-memory/plans/completed/xx-pointer-reduction.md`

#### 2. Verification Checklist:
- [ ] Small structs (<= 64 bytes) utilize value semantics (`T`) across parameters and returns.
- [ ] Read-only inspection methods use value receivers `func (s T)`.
- [ ] Zero pointer-to-slice (`*[]T`) or pointer-to-map (`*map[K]V`) constructs.
- [ ] Pointers preserved exclusively for large structs (>64-128 bytes), mutex holders, mutating receivers, and `*appfault.AppError`.
- [ ] All functions adhere to <= 8–15 lines.
- [ ] Zero build or test commands executed during routine turns.
