# Error Management & Architecture — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-error`, `cg-execute error`, `audit error`, `fix error guidelines`, `enforce error management`

> [!IMPORTANT]
> Prompt Version: 2.5.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Any directive, constraint, checklist, or instruction declared at the top of this prompt, header alert block, or incoming user request represents an absolute MUST FOLLOW mandate that takes highest priority and strictly overrides any conflicting general advice, default conventions, or lower-level guidelines below it.

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, plan, refactor, and fix all error management violations across the codebase, modifying source files directly to implement `*appfault.AppError` wrappers, outer error handling, specialized exit helpers, and universal response envelopes until 100% green without stopping.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/cg-error-management/skill.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1A (Step 0 - Verbatim Prompt Recording & Task Extraction Gate): Immediately capture the user prompt verbatim into `.ai-memory/plans/pending/xx-<slug>.md` under `## User Request (Verbatim)`, extract actionable deliverables with traceable IDs (`Task-01`, `Task-02`), and output this confirmed task breakdown directly in chat in cleanly indented markdown with vertical blank lines, task state (`State: [PENDING]`), and understanding indicator bracket (`Understood: [YES — ...]`) before any file exploration, scanning, or spec writing.
2. [ ] /goal Phase 1B (Step 1 - Master Spec Generation): Write the master architectural plan in `.ai-memory/plans/pending/xx-<slug>.md`, documenting an exhaustive Violation Ledger tracking every bare error return, swallowed error, and missing fault wrapper.
3. [ ] /goal Phase 1B (Step 2 - Scan & Discover): Use GitMap AUM discovery (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary, with fast Python discovery scripts (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) as fallback, to inventory all architectural violations and anti-patterns without truncation.
4. [ ] /goal Phase 1B (Step 3 - Lean Subtask Decomposition): Decompose the master plan into granular, lean subtasks in `.ai-memory/plans/subtasks/xx-<slug>/01-<subslug>.md`. Subtasks must focus purely on unique task deliverables without repeating common repository boilerplate.
5. [ ] /goal Phase 1B (Step 4 - Readiness Audit Gate): Confirm all `Task-xx` deliverables are mapped to subtasks and disjoint files before execution.
6. [ ] /goal Phase 1B (Zero-Stop Transition): Immediately upon completing Phase 1, self-loop and transition directly into Phase 2 execution mode without pausing or asking for permission.
7. [ ] /goal Phase 2 (Step A - Active Execution & Refactoring): Open each target file and perform surgical refactoring following authoritative guidelines: wrap all received Go errors in `*appfault.AppError`, enforce monadic `result.Wrap[T]`, eliminate bare returns, and remove all swallowed errors.
8. [ ] /goal Phase 2 (Step B - Size Tier & Formatting Enforcement): Enforce <= 8–15 line function decomposition, single return types, guard clause flattening, and clean formatting.
9. [ ] /goal Phase 2 (Step C - Failure Memory & Error Recovery): If a subagent fails, record the failure log in `.ai-memory/plan.md` and `.ai-memory/memory/issues/`; subsequent agents must read the failure log first to remediate root causes.
10. [ ] /goal Phase 2 (Step D - Change Recording & Quality Linting): Record all modified files into `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`) and run targeted file-level linters on specifically modified files (`exit 0`). DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`), unit tests, or build checks during routine turns.
11. [ ] /goal Phase 3 (Step A - Consolidation & Atomic Push): Consolidate completed subtasks into `.ai-memory/plans/completed/xx-<slug>.md`, delete granular subtasks and pending plan, stage all changes, and push in a single grouped commit.
12. [ ] /goal Phase 3 (Step B - Completion & Confidence Reporting): Emit the final Task Completion Summary with green check mark emojis, modified files summary, and implementation confidence score.
13. [ ] /learn Ingest `.ai-memory/memory/readme.md` for project memory index and past learnings.
14. [ ] /learn Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
15. [ ] /learn Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
16. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for hallucination prevention and micro-tasking.
17. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for strict relative path citation requirements.
18. [ ] /learn Ingest `02-spec/02-coding-guidelines/` for domain-specific architectural specifications.
19. [ ] /learn Ingest `02-spec/03-error-manage/` for error handling architectures and AppError.
20. [ ] /learn Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.
21. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Codebase, Write .ai-memory/plans/pending/ Spec, Create .ai-memory/plans/subtasks/, Verify/Create Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Actively Edit Code, AppError Refactoring, Linter Verification, Local CI Runner Verification, Plan Completion)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## 1. Zero Swallowed Errors Policy (TOTAL BAN — Non-Negotiable)

> [!CAUTION]
> **SWALLOWING ERRORS IS AN AUTO-REJECT FAILURE TIER 0 VIOLATION.**
> Under NO circumstances may an error be discarded, suppressed, silenced, or silently ignored.
> Every error encountered MUST either be completely resolved with structured context logging (operation name, input parameters) or embedded/wrapped into `*appfault.AppError` and returned to the caller.

### Strict Prohibitions

1. **NO Empty Catch/Except Blocks:**
   - ❌ **BANNED:** `try { ... } catch (e) {}` or `except Exception: pass`
   - ✅ **REQUIRED:** Catch blocks must log operation name, key inputs, and rethrow or return a wrapped `AppError`.
2. **NO Blank Identifier Error Discards:**
   - ❌ **BANNED:** `_ = err` or `val, _ := fn()` in Go.
   - ✅ **REQUIRED:** Check every error explicitly: `if err != nil { return appfault.Wrap(...) }`.
3. **NO Silent Fallback Defaults:**
   - ❌ **BANNED:** Returning dummy values (`return nil`, `return ""`, `return false`, `return 0`) to mask an underlying error without caller notification.
   - ✅ **REQUIRED:** Return failure status via `*appfault.AppError` or monadic `result.WrapFailure[T]`.
4. **NO Silent Nil Return in Dual Handling:**
   - ❌ **BANNED:** Calling an internal exit handler or printing an error, then returning `nil` to deceive the caller into believing execution succeeded.
   - ✅ **REQUIRED:** Leaf functions MUST return the error directly (`return err`).

---

## 2. Strict Golang Error Wrapping Mandate (`*appfault.AppError` & `appfault.Fault`)

> [!IMPORTANT]
> **ALL GOLANG ERRORS MUST BE EMBEDDED IN FAULT WRAPPERS.**
> Whenever ANY Go function encounters, intercepts, or receives an error (from the standard library `os`, `io`, `json`, `sql`, `net`, or downstream services), it MUST be immediately embedded and wrapped into `*appfault.AppError` (`appfault.Fault`).

### Core Rules for Go Error Handling

1. **Standard Error Return Type:** All domain functions returning failure metadata MUST use `*appfault.AppError` (or `appfault.Fault`).
2. **Deterministic Enum Taxonomy:** Classify errors using `errtype.Variation uint16` (`errtype.Validation`, `errtype.NotFound`, `errtype.Database`, `errtype.Network`, `errtype.Timeout`, `errtype.IO`, `errtype.Internal`). Redundant string error codes are banned.
3. **Always Wrap Standard Library Errors:**
   - Standard library `error` instances (`err != nil`) must NEVER be returned raw.
   - Use `appfault.Wrap(errType, err, "ContextMessage")` or `appfault.WrapFile(errType, err, relativePath, "ContextMessage")`.
4. **Monadic Result Wrapper Mandate (`pkg/result`):**
   - Functions returning a value along with possible failure MUST return `result.Wrap[T]` (`appfault.Result[T]`).
   - Bare tuples `(T, error)` across public domain boundaries are strictly prohibited.
   - Return success using `result.WrapSuccess(val)` and failure using `result.WrapFailure[T](fault)`.
5. **Context Enrichment:**
   - Chain contextual metadata: `.WithOp("Package.Function")`, `.WithVar("key", val)`, `.WithSiteId(siteId)`.
   - Use relative repository paths only (TOTAL BAN on absolute paths or `file:///` URIs).
6. **Zero Redundant Re-Wrapping:**
   - If a downstream function already returns `*appfault.AppError` or `result.Wrap[T]`, propagate the existing fault directly using `result.WrapFailureFromWrap[T](downstreamRes)` rather than wrapping it again.
7. **Mandatory Concrete Types in `types.go` (Total Ban on Leaking Raw Generics Across Signatures):**
   - **No Leaked Raw Generics:** NEVER leak raw generic Result wrappers (`result.Wrap[*Config]`, `result.Wrap[User]`, `result.ResultSlice[T]`) across function signatures, service boundaries, or public packages.
   - **Convert Reused Types to Concrete Named Types:** Rather than scattering raw generics everywhere, if a result type is used or reused across functions or layers, define a single reusable concrete type alias in `types.go` (e.g. `type ConfigResult = result.Wrap[*Config]`, `type UserResult = result.Wrap[User]`) for Golang (and equivalent leaf type definitions for other languages, e.g. `export type UserResult = Result<User>;`).
   - **Explanatory Code Comments:** Code examples and implementation files MUST include comments showing how the concrete type is declared in `types.go` and follows through into the function signatures.

---

## 3. Production Go Code Samples (Refer to `04-code/golang/examples/`)

> Real-world implementations are maintained in [`04-code/golang/examples/database_query.go`](04-code/golang/examples/database_query.go), [`04-code/golang/examples/workflow_service.go`](04-code/golang/examples/workflow_service.go), and [`04-code/golang/examples/types.go`](04-code/golang/examples/types.go).
> Always inspect those source files as the canonical ground truth.

### Sample 1: Standard Library File & JSON Handling (Concrete `ConfigResult` in `types.go`)

```go
// -----------------------------------------------------------------------------
// Step 1: Declare Concrete Types in `types.go` (Mandatory Rule)
// -----------------------------------------------------------------------------
// In types.go:
// type (
//     // Config contains application configuration fields.
//     Config struct {
//         Port int    `json:"port"`
//         Host string `json:"host"`
//     }
//
//     // ConfigResult is the single reusable concrete result envelope for *Config.
//     // RULE: Convert raw generic result.Wrap[*Config] to an explicit concrete type
//     // in types.go so all signatures and callers share the exact same definition!
//     ConfigResult = result.Wrap[*Config]
// )
// -----------------------------------------------------------------------------

// ❌ FORBIDDEN: Bare error returns, uninformative errors.New, swallowed errors, missing blank lines
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        _ = err // ❌ SWALLOWED ERROR
        return nil, err // ❌ Missing blank line before return, bare error return
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil { // ❌ Semicolon in if
        return nil, errors.New("invalid json") // ❌ BARE ERROR WITHOUT CAUSE
    }

    return &cfg, nil
}

// ✅ REQUIRED: Strict appfault wrapping with errtype, concrete ConfigResult from types.go, mandatory blank lines, flat ifs
func LoadConfig(path string) ConfigResult {
    if path == "" {
        fault := appfault.New(errtype.Validation, "config path cannot be empty").
            WithOp("config.LoadConfig")

        return result.WrapFailure[*Config](fault)
    }

    data, err := os.ReadFile(path)

    if err != nil {
        fault := appfault.WrapFile(errtype.IO, err, path, "failed to read configuration file").
            WithOp("config.LoadConfig")

        return result.WrapFailure[*Config](fault)
    }

    var cfg Config
    err = json.Unmarshal(data, &cfg)

    if err != nil {
        fault := appfault.Wrap(errtype.Validation, err, "failed to parse configuration json").
            WithOp("config.LoadConfig").
            WithVar("path", path)

        return result.WrapFailure[*Config](fault)
    }

    return result.WrapSuccess(&cfg)
}
```

### Sample 2: Database Query with Result Monad (Concrete `UserResult` in `types.go`)

```go
// -----------------------------------------------------------------------------
// Step 1: Declare Concrete Types in `types.go` (Mandatory Rule)
// -----------------------------------------------------------------------------
// In types.go:
// type (
//     // User represents the persistent user entity model.
//     User struct {
//         Id    int64  `json:"id"`
//         Name  string `json:"name"`
//         Email string `json:"email"`
//     }
//
//     // UserResult is the canonical single reusable result envelope for User.
//     // RULE: Convert raw generic result.Wrap[User] to a concrete named type in types.go.
//     // Never leak raw generic parameters across package boundaries and service signatures.
//     UserResult = result.Wrap[User]
// )
// -----------------------------------------------------------------------------

// ❌ FORBIDDEN: Combined if with semicolon, nested if, error-type branching, missing blank lines
func (r *UserRepository) FindUser(ctx context.Context, id int64) (*User, error) {
    row := r.db.QueryRowContext(ctx, "SELECT name, email FROM users WHERE id = ?", id)
    var user User

    if err := row.Scan(&user.Name, &user.Email); err != nil { // ❌ Combined semicolon if
        if errors.Is(err, sql.ErrNoRows) { // ❌ FORBIDDEN: Nested if and error-type branching
            return nil, nil // ❌ Deceptive swallowed error
        }
        return nil, err // ❌ Missing blank line before return
    }

    user.Id = id
    return &user, nil // ❌ Missing blank line before return
}

// ✅ REQUIRED: Flat if guard, direct error typing without branching, concrete UserResult from types.go, blank lines before return
func (r *UserRepository) FindUser(ctx context.Context, id int64) UserResult {
    if id <= 0 {
        fault := appfault.New(errtype.Validation, "user id must be positive").
            WithOp("UserRepository.FindUser").
            WithVar("id", id)

        return result.WrapFailure[User](fault)
    }

    row := r.db.QueryRowContext(ctx, "SELECT name, email FROM users WHERE id = ?", id)
    var user User

    err := row.Scan(&user.Name, &user.Email)

    if err != nil {
        // Direct error typing: select the error type reflecting this layer (errtype.Database)
        // Attach the ID and variables directly. Never branch on error types or nest ifs!
        fault := appfault.Wrap(errtype.Database, err, "failed to scan user row from database").
            WithOp("UserRepository.FindUser").
            WithVar("id", id)

        return result.WrapFailure[User](fault)
    }

    user.Id = id

    return result.WrapSuccess(user)
}
```

### Sample 3: Propagating Errors Across Boundaries (Concrete `UserResult` Across Services)

```go
// ✅ REQUIRED: Propagate downstream Fault directly without redundant nested wrapping, using concrete UserResult
func (s *UserService) ActivateUser(ctx context.Context, userId int64) UserResult {
    userRes := s.repo.FindUser(ctx, userId)

    if userRes.IsFailed() {
        s.log.LogError(userRes.Fault())

        // Propagate existing Fault directly with zero re-wrapping
        return result.WrapFailureFromWrap[User](userRes)
    }

    user := userRes.Value()
    user.IsActive = true

    updateRes := s.repo.UpdateUser(ctx, user)

    if updateRes.IsFailed() {
        s.log.LogError(updateRes.Fault())

        return result.WrapFailureFromWrap[User](updateRes)
    }

    return result.WrapSuccess(user)
}
```

---

## 4. Dedicated Section: Error Return Contract & Outer Handling Principle (Zero Dual-Handling)

A function that declares an error or result return type MUST return the actual error instance directly to the caller. It MUST NEVER invoke an exit handler, terminate the process, or panic internally and then return `nil`.

### Why Dual-Handling & Internal Exit Is Forbidden

1. **Broken Caller Sovereignty:** When a leaf function handles its own exit internally and returns `nil`, the caller is deceived into believing the operation succeeded.
2. **Impossible Testability:** Unit tests cannot assert returned error types or values if the helper function kills the process or handles errors internally.
3. **Dual Execution Hazards:** Calling an exit handler inside a helper while returning a result creates race conditions, partial database mutations, and skipped resource cleanups.

### Mandatory Outer Handling Pattern

- **Leaf/Service Functions:** Construct or wrap `*appfault.AppError` and return it.
- **Top-Level Root Dispatcher / HTTP Router:** Only the outer controller handles the error, decides the exit code via `ExitCodeType` enum, and writes the Universal Response Envelope:

```go
// ✅ REQUIRED: Top-level caller handles the error and exit
func MainCommandDispatcher(args []string) {
    res := ExecuteOperation(args)
    if res.IsFailed() {
        exitHandler.HandleValidationError(res.Fault())
        return
    }

    exitHandler.HandleSuccess()
}
```

---

## 5. Phase 1: Scan Codebase & Write Implementation Spec First (Steps 1 to PHASE_1_STEPS)

Before modifying application code, you MUST thoroughly scan the repository and write an actionable execution spec.

### Step 2: Scan & Discover (Python Toolchain Acceleration)

To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, use the repository's dedicated Python discovery scripts:
- **Inventory Target Files:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats`
- **Fast Cached Grep (<15ms):** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<search-pattern>" --limit 50`
- **Sub-Millisecond Folder & File Exploration:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder <folder-path> --limit 50`
- **Read Target File:** `python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000`
- **Fast Pattern Search:** `python 03-ai-scripts/17-fast-file-reader.py --search-pattern "<pattern>" --limit 50`
- **Subsystem & Topology Overview:** `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`

- **Actionable Scan:** Search across all Go, TypeScript, PHP, and Python files to identify:
  1. Internal exit handler calls in leaf functions followed by `return nil`.
  2. Empty `catch` / `except` blocks or swallowed errors (`_ = err`, unused errors).
  3. Bare error returns without contextual wrapping (`return err` instead of `appfault.Wrap`).
  4. Raw panic / exit invocations (`panic()`, `process.exit()`, `os.Exit()`).
  5. Magic integer exit codes (`HandleError(err, 1)` instead of enums).
  6. API endpoints returning raw text or unformatted error payloads instead of the `{ data, errors, meta }` envelope.
  7. Functions exceeding 8 lines (hard cap 15 lines) or files exceeding 100 coding lines (recommended <= 80).
  8. Nested `if` conditionals.
- **Where to save it:** Save this master plan into `.ai-memory/plans/pending/xx-error-management-audit.md` listing every affected file and exact line number.
- **Subtasks:** Break the plan down into granular subtask files inside `.ai-memory/plans/subtasks/xx-error-management/` (e.g. `01-leaf-error-returns.md`, `02-specialized-exit-helpers.md`, `03-api-response-envelopes.md`).

---

## 6. Authoritative Spec Files Checklist (Non-Negotiable Action Items)

You MUST read, follow, and mechanically verify every single specification file below before and during execution:

- [ ] **`02-spec/02-coding-guidelines/02-canonical-size-tier.md`**
  - **Why:** Universal size limits across all languages.
  - **How:** Functions <= 8 lines preferred (hard cap 15 lines). Files <= 100 lines coding max (recommended <= 80 lines). Zero line-compression cheating.
- [ ] **`02-spec/02-coding-guidelines/06-ai-optimization/readme.md`**
  - **Why:** Comprehensive catalog of forbidden vs required generation patterns.
  - **How:** Strictly follow AH-N1 to AH-T2 rules. Zero ghost diffs, zero truncation stubs (`// ...`), zero unverified claims.
- [ ] **`02-spec/02-coding-guidelines/06-ai-optimization/06-citation-requirement.md`**
  - **Why:** Grounded rule enforcement and traceability.
  - **How:** Cite authoritative spec files for every code modification made.
- [ ] **`02-spec/02-coding-guidelines/01-cross-language/04-code-style/02-braces-and-nesting.md`**
  - **Why:** Absolute zero tolerance for nested conditionals.
  - **How:** Flatten all nested `if` statements with guard clauses and early returns.
- [ ] **`02-spec/03-error-manage/readme.md`**
  - **Why:** Authoritative error management foundation across all services.
  - **How:** Never swallow errors; every `catch` logs with operation name and key inputs, then rethrows or returns a typed error.
- [ ] **`02-spec/03-error-manage/02-error-architecture/02-error-handling-reference.md`**
  - **Why:** Universal cross-stack `AppError` and `AppException` structure.
  - **How:** Implement `appfault.Wrap(errType, err, "OpName")` in Go, `throw new AppError(cause, { op, ctx })` in TS, and `AppException` in C#/PHP; preserve the root cause and causal stack.
- [ ] **`02-spec/03-error-manage/02-error-architecture/03-go-delegation-fix.md`**
  - **Why:** Prevents nil pointer panics and raw error leaks in Go routines.
  - **How:** Never delegate errors to uninitialized handlers; use explicit, typed error delegation channels with mutex guards.
- [ ] **`02-spec/03-error-manage/02-error-architecture/readme.md`**
  - **Why:** Standardized error severity and UI feedback mapping.
  - **How:** Map log levels strictly: `debug` (trace), `info` (lifecycle), `warn` (recoverable/amber), `error` (user-visible failure/red), `fatal` (process exit).
- [ ] **`02-spec/03-error-manage/02-error-architecture/05-response-envelope/readme.md`**
  - **Why:** Universal API response contract across all endpoints.
  - **How:** Every HTTP/RPC response MUST return the standard envelope: `{ "data": T, "errors": [AppError], "meta": Meta }`. Never return raw un-enveloped error text.
- [ ] **`02-spec/03-error-manage/03-error-code-registry/readme.md`**
  - **Why:** Stable error code registry and catalog.
  - **How:** All error codes must be registered constants (`errtype.Variation`). No ad-hoc string literals invented at the throw site.

---

## 7. Mandatory Linter & Targeted Verification Checklist

Code standards must be mechanically enforced by automated linters. You MUST verify or create the linter and connect it to CI:

- [ ] **Linter Script Identification:** Check if `linter-scripts/check-error-management.py` exists in the repository.
- [ ] **Auto-Create Linter if Missing:** If no dedicated error linter exists, create `linter-scripts/check-error-management.py` that AST-scans for:
  1. Internal exit handler invocations in non-main functions.
  2. Empty `catch` or `except` blocks (swallowed errors).
  3. Bare un-wrapped error returns (`return err` instead of `appfault.Wrap`).
  4. Bare panics/hard exits (`panic()`, `process.exit()`, `os.Exit()`).
  5. Non-standard API responses lacking the `{ data, errors, meta }` envelope.
- [ ] **Local Linter Command:** Execute and verify the linter locally on modified files:
  ```bash
  python linter-scripts/check-error-management.py
  ```

---

## 8. Phase 2: Active Code Refactoring & Autonomous Fix Loop (Steps PHASE_1_STEPS+1 to N)

> [!IMPORTANT]
> **AUTONOMOUS EXECUTION MANDATE — DO NOT STOP.**
> Open the offending source code files and directly rewrite the code to eliminate violations. Maintain continuous self-looping until all checks pass 100% green.

```text
STEP = 0
WHILE (STEP < PHASE_2_STEPS):
    STEP += 1

    1. Read the next subtask from .ai-memory/plans/subtasks/xx-error-management/
    2. Open and modify the actual source code files:
       - Ensure leaf functions return errors (return err) and do not call exit handlers internally.
       - Wrap all received Go errors using appfault.Wrap(errType, err, "msg") or result.WrapFailure[T].
       - Enforce zero swallowed errors (no empty catch/except, no _ = err, no silent nil returns).
       - Use typed Enums (ExitCodeType) for exit codes instead of magic numbers.
       - Create specialized helper functions to reduce repeated handler parameters.
       - Inject operation name and parameter context into all error logs.
       - Enforce the universal { data, errors, meta } API envelope.
       - Keep function bodies <= 8 lines (max 15 lines) and files <= 100 coding lines.
       - Flatten any nested ifs with guard clauses.
    3. Run the error management linter:
          python linter-scripts/check-error-management.py
    4. Run targeted guideline autofixers:
          python 03-ai-scripts/05-guideline-autofixer.py <modified-files>
          python 03-ai-scripts/08-naming-autofixer.py <modified-files>
          python 03-ai-scripts/04-newline-fixer.py <modified-files>
          python 03-ai-scripts/07-relative-path-fixer.py <modified-files>
    5. Record modified files (DO NOT run builds, test runners, or full CI runner):
          python 03-ai-scripts/33-test-inventory-generator.py --record <modified-files>
    6. IF any check fails:
          - Diagnose failure, fix code directly, and re-run immediately.
       IF all checks pass (exit code 0):
          - Mark subtask completed and proceed to next subtask.

    7. When all subtasks are finished and local checks are 100% green:
          - Move .ai-memory/plans/pending/xx-error-management-audit.md to .ai-memory/plans/completed/
          - Update .ai-memory/plans/readme.md
          - Run plan consolidator:
            python 03-ai-scripts/20-plan-consolidator.py
          - Stage modified files with git add and create semantic commit:
            git commit -m "refactor(errors): enforce AppError fault wrapping, outer error handling, and zero swallowed errors"
          - Push to git repository and finish turn.
```

---

## 9. Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] **Top-Instruction Priority Respected:** Highest precedence given to preamble directives and header mandates.
- [ ] **Git Working Tree Clean:** Working tree verified before new code changes.
- [ ] **Sub-agents Assigned Disjoint Files:** Verified against `.ai-memory/readme.md`.
- [ ] **Zero Swallowed Errors:** Verified no empty catches, no blank identifiers (`_ = err`), and no silent fallback returns.
- [ ] **Fault Wrapper Enforcement:** Every Go error encountered is strictly embedded in `*appfault.AppError` / `result.Wrap[T]`.
- [ ] **Error Return Sovereignty:** Leaf functions return `error` (`return err`); NO leaf functions call exit handlers or panics and return `nil`.
- [ ] **Typed Exit Codes:** Exit codes use strongly typed enums (`ExitCodeType`), zero magic integers.
- [ ] **Parameter Reduction:** Repeated handler parameters extracted into specialized helper functions.
- [ ] **Zero Nested Ifs:** NO nested `if` blocks exist; all flattened with guard clauses.
- [ ] **Function Size Limits:** All functions <= 8 lines preferred, hard cap 15 lines.
- [ ] **File Size Limits:** Files <= 100 lines coding max (recommended <= 80 lines).
- [ ] **Strict Relative Git Paths:** All file paths, markdown links, citations, and subtask references are strictly relative to git root. Zero absolute paths or `file:///` URIs.
- [ ] **Banned Operations Honored:** Zero unit tests run, zero build commands executed, zero uncommanded releases triggered, zero per-file commits.
- [ ] **Final Step Commit & Push Verified:** Staged all changes (`git add -A`), committed in a single grouped atomic commit, and pushed to git before ending the turn.

---

## 10. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.

---

## Metadata

- slug: cg-error-management
- priority: high
- status: active
