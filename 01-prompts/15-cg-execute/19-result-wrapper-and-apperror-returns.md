# Result Wrapper Types, Collections & AppError Returns — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-result-wrapper`, `cg-apperror-returns`, `cg-execute result-wrapper`, `audit result wrapper`, `fix map return error`, `fix slice return error`, `single return object audit`, `enforce apperror returns`, `enforce result map`, `fix multi-value returns`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, discover, plan, refactor, and verify all Go functions returning multi-value error tuples (such as `(map[K]V, error)`, `([]T, error)`, or `(T, error)`), eliminating raw standard library error returns, replacing them with strongly-typed result wrappers (`ResultMap[K, V]`, `ResultSlice[T]`, `Result[T]`) and structured `*appfault.AppError` returns, guaranteeing a single return object, standardized outer-layer inspection predicates (`IsSuccess`, `IsFailure`, `HasError`, `IsEmptyError`, `IsEmpty`, `Data`, `AppError`, `Fault`, `Get`, `Has`, `Count`), and zero dual-handling across the entire codebase until 100% green without stopping.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase using ripgrep to inventory all functions returning multi-value tuples `(T, error)`, `(map[K]V, error)`, `([]T, error)`, and raw stdlib `error` returns.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.lovable/plans/pending/XX-result-wrapper-audit.md` with an exhaustive Violation Ledger table.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.lovable/plans/subtasks/XX-result-wrapper/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/01-index.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and refactor function signatures from multi-value returns to single `ResultMap[K, V]`, `ResultSlice[T]`, or `Result[T]` envelopes.
6. [ ] /goal Phase 2 (Step B): Replace raw stdlib `error` returns with structured `*appfault.AppError` instances using `appfault.New()` or `appfault.Wrap()`.
7. [ ] /goal Phase 2 (Step C): Modernize all caller call sites to utilize outer-layer inspection methods (`res.IsSuccess()`, `res.IsFailure()`, `res.IsEmpty()`, `res.Get()`, `res.AppError()`).
8. [ ] /goal Phase 2 (Step D): Enforce <= 8–15 line function decomposition and clean blank-line spacing.
9. [ ] /goal Phase 2 (Step E): Execute targeted file-level linters (`python linter-scripts/check-function-lengths.py`, `check-mws-error-codes.py`, `check-newline-styling.py`) to verify 0 remaining violations. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine coding guideline execution turns.
10. [ ] /learn Ingest `.lovable/memory/01-index.md` for project memory index and past learnings.
11. [ ] /learn Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
12. [ ] /learn Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
13. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/01-index.md` for single return type mandates and micro-tasking.
14. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/01-index.md` for strict relative path citation requirements.
15. [ ] /learn Ingest `02-spec/03-error-manage/01-index.md` for universal AppError wrapping and error envelopes.
16. [ ] /learn Ingest `02-spec/03-error-manage/02-error-architecture/02-error-handling-reference.md` for error handling architecture and Result wrappers.
17. [ ] /learn Ingest `02-spec/03-error-manage/03-error-code-registry/02-registry.md` for structured error code catalog.
18. [ ] /learn Ingest `02-spec/03-error-manage/02-error-architecture/05-response-envelope/05-response-envelope-reference.md` for response envelope schemas.
19. [ ] /learn Ingest `02-spec/03-error-manage/02-error-architecture/06-apperror-package/03-go-apperror-linter-spec.md` for Go AppError implementation specifications.
20. [ ] /learn Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.
21. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Multi-Value Returns, Build Violation Ledger in .lovable/plans/pending/, Subtasks, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Actively Edit Code, Refactor Signatures to ResultMap/Result, Modernize Call Sites, Verify Local Linters)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## The Code Pattern & Anti-Pattern Analysis

In legacy or standard Go codebases, developers frequently write functions that return multi-value tuples pairing collections with the standard library `error` interface.

### The Problematic Legacy Pattern

Consider this common database query implementation:

```go
// ❌ ANTI-PATTERN: Multi-value tuple return with raw standard library error
func (s *SQLiteStore) queryAllMacroSteps(db *sql.DB) (map[string][]MacroStep, error) {
    rows, err := db.Query("SELECT macro_id, step_name, action, payload FROM macro_steps ORDER BY macro_id, step_order")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    return scanMacroStepsMap(rows)
}

// ❌ ANTI-PATTERN: Secondary scanner returning raw map and error tuple
func scanMacroStepsMap(rows *sql.Rows) (map[string][]MacroStep, error) {
    stepsMap := make(map[string][]MacroStep)
    for rows.Next() {
        var macroId, name, action, payload string
        if err := rows.Scan(&macroId, &name, &action, &payload); err != nil {
            return nil, err
        }
        stepsMap[macroId] = append(stepsMap[macroId], MacroStep{
            Name:    name,
            Action:  action,
            Payload: payload,
        })
    }

    return stepsMap, rows.Err()
}
```

### Why This Pattern Violates Repository Guidelines

1. **Violates the Single Return Object Mandate:**
   - Multi-value returns like `(map[string][]MacroStep, error)` violate the core repository architectural standard requiring functions to return a single, strongly-typed envelope object.
   - Returning tuples leaks internal branching mechanics into the caller and forces dual assignment unpacking (`val, err := ...`).
2. **Raw Standard Library `error` Anti-Pattern:**
   - Returning the bare standard library `error` interface strips domain context, structured error codes, file/line tracing, and machine-readable metadata.
   - All errors in the repository MUST use structured `*appfault.AppError`.
3. **Dual-Handling Risk and Ambiguous Empty Returns:**
   - Does returning `(nil, nil)` represent an empty database table or an uninitialized store?
   - If an error occurs midway through row scanning, returning `(nil, err)` drops partially collected records, while returning `(stepsMap, err)` tempts callers into dual handling (processing data *and* logging error).
4. **Call-Site Boilerplate Multiplication:**
   - Callers are forced to write verbose, error-prone unpacking boilerplate:
     ```go
     stepsMap, err := store.queryAllMacroSteps(db)
     if err != nil {
         return fmt.Errorf("failed: %w", err)
     }
     if len(stepsMap) == 0 { ... }
     ```

---

## The Modern Refactored Architecture

Under the Prompt Architect coding guidelines, all multi-value returns are refactored into dedicated result container types from package `pkg/appfault`:
- Key-Value Maps: `appfault.ResultMap[K, V]`
- Lists & Slices: `appfault.ResultSlice[T]`
- Scalar Values: `appfault.Result[T]`
- Pure Side-Effects: `*appfault.AppError` (zero bare `void` / empty returns)

### Modern Refactored Implementation

```go
// ✅ MODERN PATTERN: Single ResultMap return envelope with structured AppError
func (s *SQLiteStore) queryAllMacroSteps(db *sql.DB) appfault.ResultMap[string, []MacroStep] {
    rows, err := db.Query("SELECT macro_id, step_name, action, payload FROM macro_steps ORDER BY macro_id, step_order")
    if err != nil {
        return appfault.FailMap[string, []MacroStep](
            appfault.New(appfault.ErrDatabaseQuery).
                WithCause(err).
                WithMessage("failed to query macro steps from database"),
        )
    }
    defer rows.Close()

    return scanMacroStepsMap(rows)
}

// ✅ MODERN PATTERN: Scanner returning strongly-typed ResultMap
func scanMacroStepsMap(rows *sql.Rows) appfault.ResultMap[string, []MacroStep] {
    stepsMap := make(map[string][]MacroStep)
    for rows.Next() {
        var macroId, name, action, payload string
        if err := rows.Scan(&macroId, &name, &action, &payload); err != nil {
            return appfault.FailMap[string, []MacroStep](
                appfault.New(appfault.ErrDatabaseScan).
                    WithCause(err).
                    WithMessage("failed to scan macro step row"),
            )
        }

        stepsMap[macroId] = append(stepsMap[macroId], MacroStep{
            Name:    name,
            Action:  action,
            Payload: payload,
        })
    }

    if err := rows.Err(); err != nil {
        return appfault.FailMap[string, []MacroStep](
            appfault.New(appfault.ErrDatabaseIteration).
                WithCause(err).
                WithMessage("row iteration failed for macro steps"),
        )
    }

    return appfault.OkMap(stepsMap)
}
```

---

## Standardized Outer-Layer Inspection Methods

Every result wrapper provides a unified suite of predicates and accessors. Outer calling layers inspect results cleanly without unpacking tuples:

| Method | Return Type | Purpose & Behavior |
|---|---|---|
| `res.IsSuccess()` | `bool` | Returns `true` if the operation succeeded and no error is present. |
| `res.IsFailure()` / `res.IsFailed()` | `bool` | Returns `true` if the operation encountered a failure or error. |
| `res.HasError()` | `bool` | Returns `true` if an active error is attached to the result. |
| `res.IsEmptyError()` / `res.HasNoError()` | `bool` | Returns `true` if no active error exists. |
| `res.IsEmpty()` | `bool` | Returns `true` if the underlying map/slice has 0 items or result is null/empty. |
| `res.Data` / `res.Value()` | `map[K]V` / `T` | Accesses the underlying data payload directly. |
| `res.AppError()` / `res.Fault()` | `*appfault.AppError` | Retrieves the structured error context for logging, propagation, or HTTP response mapping. |
| `res.Get(key)` | `(V, bool)` | Safely retrieves a map entry by key without nil-map panics. |
| `res.Has(key)` | `bool` | Checks whether a key exists within the result map. |
| `res.Count()` | `int` | Returns the number of entries in the collection (or `0` if failed). |
| `res.Keys()` | `[]K` | Returns a deterministically sorted slice of all map keys. |
| `res.Values()` | `[]V` | Returns a slice of map values ordered according to `Keys()`. |

### Caller Call-Site Modernization

#### ❌ Legacy Caller:

```go
steps, err := store.queryAllMacroSteps(db)
if err != nil {
    return fmt.Errorf("failed to load macro steps: %w", err)
}
if len(steps) == 0 {
    log.Println("No steps found")
    return nil
}
processSteps(steps["init"])
```

#### ✅ Modern Caller:

```go
res := store.queryAllMacroSteps(db)
if res.IsFailure() {
    return res.AppError().WithContext("caller", "executeWorkflow")
}

if res.IsEmpty() {
    log.Println("No macro steps found")
    return nil
}

if steps, exists := res.Get("init"); exists {
    processSteps(steps)
}
```

---

## Error Management Learning Checklist (`02-spec/03-error-manage/`)

Before refactoring error handling in any package, the agent must study and enforce the repository error management specifications:

- [ ] **Universal `*appfault.AppError` Standard (`02-spec/03-error-manage/01-index.md`):**
  - Never return bare `error` from domain services, repositories, or business logic.
  - Wrap third-party and standard library errors with `appfault.New()` or `appfault.Wrap()`.
- [ ] **Structured Error Codes (`02-spec/03-error-manage/03-error-code-registry/02-registry.md`):**
  - All errors must carry a typed `ErrorCode` string identifying the fault category (e.g. `ErrDatabaseQuery`, `ErrValidationFailed`, `ErrNotFound`).
- [ ] **Deterministic Error Handling & Envelopes (`02-spec/03-error-manage/02-error-architecture/02-error-handling-reference.md`):**
  - Use `appfault.Ok()`, `appfault.OkMap()`, and `appfault.OkSlice()` for successful results.
  - Use `appfault.Fail()`, `appfault.FailMap()`, and `appfault.FailSlice()` for failed results.
- [ ] **Universal Response Envelopes (`02-spec/03-error-manage/02-error-architecture/05-response-envelope/05-response-envelope-reference.md`):**
  - HTTP handlers and JSON serializers marshal `Result` and `ResultMap` into universal JSON response envelopes `{ "data": ..., "appError": ... }`.
- [ ] **Go AppError Architecture (`02-spec/03-error-manage/02-error-architecture/06-apperror-package/03-go-apperror-linter-spec.md`):**
  - Strict enforcement of `*appfault.AppError` return types and monadic helper methods across Go packages.

---

## Automated Codebase Scanning Guide

Use these exact `ripgrep` regex commands to discover legacy multi-value return patterns across the codebase:

```bash
# 1. Find functions returning multi-value map tuples: (map[...], error)
rg --pcre2 "func\s+\w+\([^\)]*\)\s*\(\s*map\[[^\]]+\][^,]+,\s*error\)"

# 2. Find functions returning multi-value slice tuples: ([]..., error)
rg --pcre2 "func\s+\w+\([^\)]*\)\s*\(\s*\[\][^,]+,\s*error\)"

# 3. Find any function returning a tuple ending in standard library error
rg --pcre2 "func\s+\w+\([^\)]*\)\s*\([^\)]*,\s*error\)"

# 4. Find functions returning bare standard library error
rg --pcre2 "func\s+\w+\([^\)]*\)\s+error\s*\{"

# 5. Find dual-assignment caller unpacking: val, err := ...
rg --pcre2 "\b(\w+),\s*err\s*:=\s*"
```

---

## 2-Agent Parallel Orchestration

To survive large codebases without hitting step limits or context loss, execute this prompt using a strict 2-agent parallel split:

```text
+-------------------------------------------------------------------------+
| MASTER ORCHESTRATOR (Budget: N = 200)                                   |
|                                                                         |
| Phase 1 (Steps 1..100): DISCOVERY & PLANNING                            |
| +---------------------------------------------------------------------+ |
| | Sub-Agent 1: Codebase Scanner & Spec Architect                      | |
| | - Runs ripgrep queries to catalog all multi-value error returns     | |
| | - Authors master audit plan in .lovable/plans/pending/             | |
| | - Generates granular subtasks in .lovable/plans/subtasks/           | |
| +---------------------------------------------------------------------+ |
|                                                                         |
| Phase 2 (Steps 101..200): SURGICAL REFACTORING                          |
| +---------------------------------------------------------------------+ |
| | Sub-Agent 2: Code Refactorer & Outer-Layer Modernizer                | |
| | - Refactors store and repository signatures to ResultMap/ResultSlice| |
| | - Updates scanner functions to use appfault.OkMap / FailMap         | |
| | - Modernizes caller call-sites with IsSuccess / IsFailure / Get     | |
| | - Verifies zero regressions with targeted file linters              | |
| +---------------------------------------------------------------------+ |
+-------------------------------------------------------------------------+
```

---

## Strictly Avoid: Anti-Patterns & Prohibitions

- **NO PIECEMEAL COMMITS:** NEVER commit 1 or 2 files in isolation. Consolidate all related changes across specs, code, and indices into a single atomic commit followed immediately by `git push origin main`.
- **NO ROUTINE FULL CI/CD RUNS:** DO NOT run `06-cicd-local-runner.py` during normal turns. It executes 28-38 heavy validation gates across unrelated packages and wastes minutes. Run targeted linters only on modified files.
- **NO UNIT TEST EXECUTION:** Test execution is strictly disabled unless explicitly commanded by the repository owner.
- **NO RAW `error` RETURNS:** Never leave bare `error` as a return type on domain or store functions; always use `*appfault.AppError` or `Result[T]`.
- **NO ABSOLUTE PATHS:** Never write absolute filesystem paths (`C:\...`, `/home/...`) or `file:///` URIs. Use strict relative Git paths starting from the repository root.
- **NO UPPERCASE FILENAMES:** Every file created or edited must be strictly lowercase.
- **NO MULTI-VALUE TUPLES:** Eliminate `(T, error)` in favor of `Result[T]`, `ResultMap[K, V]`, or `ResultSlice[T]`.
