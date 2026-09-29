---
name: cg-golang-pointer-reduction-and-value-semantics
description: Autonomously audits, plans, and refactors Golang codebases to reduce pointer usage (*T) and enforce value semantics across small structs, parameters, and query receivers using a 300-step 3-phase self-loop.
---

# Skill: Golang Pointer Reduction & Value Semantics (`cg-pointer-reduction`)

This skill governs autonomous execution for migrating Golang codebases to value semantics, reducing heap allocations, and eliminating pointer indirection and nil-dereference panics across all Go packages.

## The 300-Step 3-Phase Architecture

- **Phase 1: Deep Audit & Struct Sizing (Steps 1–100):**
  - Inventory all struct declarations, pointer fields, pointer parameters (`*T`), pointer receivers (`func (s *T)`), and pointer return types.
  - Calculate memory footprints: classify structs into Value Candidates (<= 64 bytes) vs. Pointer Retainers (> 64–128 bytes, or containing mutexes).
  - Write master audit spec in `.ai-memory/plans/pending/xx-pointer-reduction-audit.md` with an exhaustive Violation Ledger.
- **Phase 2: Architectural Specification & Type Selection (Steps 101–200):**
  - Author canonical architecture spec in `02-spec/21-app/xx-golang-value-semantics.md`.
  - Define clear migration rules: query methods to value receivers, small parameter structs by value, slice/map pointer elimination.
  - Decompose into granular subtasks in `.ai-memory/plans/subtasks/`.
- **Phase 3: Active Refactoring & Quality Verification (Steps 201–300):**
  - Execute surgical refactoring in micro-batches (5–8 files).
  - Migrate small structs and read-only receivers to value semantics.
  - Verify with targeted file-level linters (`exit 0`).
  - Consolidate plans and finalize with an atomic push.

## Pointer Justification Matrix

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

## Canonical Code Transformation

### ❌ Anti-Pattern: Unnecessary Pointers

```go
// ❌ WRONG — heap allocation for small struct, pointer receiver on query, pointer to slice
type QueryOptions struct {
    Limit  int
    Offset int
}

func NewQueryOptions() *QueryOptions {
    return &QueryOptions{Limit: 20}
}

func (q *QueryOptions) IsPaged() bool {
    return q.Limit > 0
}

func FetchRows(opts *QueryOptions, ids *[]string) { ... }
```

### ✅ Canonical Pattern: Value Semantics

```go
// ✅ REQUIRED — value semantics, stack allocated, value receiver, direct slice
type QueryOptions struct {
    Limit  int
    Offset int
}

func NewQueryOptions() QueryOptions {
    return QueryOptions{Limit: 20}
}

func (q QueryOptions) IsPaged() bool {
    return q.Limit > 0
}

func FetchRows(opts QueryOptions, ids []string) { ... }
```

## Cross-References

- [03-type-safety-and-errors.md](02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/03-type-safety-and-errors.md)
- [04-database-and-structs.md](02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/04-database-and-structs.md)
- [02-canonical-size-tier.md](02-spec/02-coding-guidelines/02-canonical-size-tier.md)
