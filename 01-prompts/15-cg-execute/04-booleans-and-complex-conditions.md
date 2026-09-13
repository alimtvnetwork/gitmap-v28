# Boolean Principles, Negatives & Complex Conditions — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-boolean`, `cg-execute boolean`, `audit boolean`, `fix boolean negatives`, `fix complex conditions`, `affirmative booleans`, `boolean-parameter-naming`, `affirmative-boolean-parameters`, `is-stopped`, `fix-v-bool`

> **Prompt Version:** 2.3.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, plan, refactor, and fix all boolean naming, double negatives, mixed polarity, single-letter boolean parameters (`v bool`, `b bool`), bare verb/noun identifiers (`stop bool` -> `isStopped bool`, `stopOnFail` -> `isStopOnFail`), awkward `isExists` identifiers (`isExists` -> `isDefined`), compound negative chains (`!a || !b || c` -> discrete checks), and complex condition violations across the codebase, modifying source files directly to enforce affirmative prefixes (is and has only (can, should, was, etc. are banned)), implicit evaluation (no `== true`), positive framing (no `!isSuccess`), and discrete condition decomposition until 100% green without stopping.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase using AST and ripgrep to inventory all architectural violations: explicit `== true`/`== false`, negative names (`isNot*`), inverted success (`!isSuccess`), mixed polarity (`&& !`), single-letter boolean parameters (`v bool`, `b bool`), bare un-prefixed boolean identifiers (`stop`, `pause`, `force`, `dryRun`), awkward `isExists` identifiers, and compound negative chains (`!a || !b || c`).
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.lovable/plans/pending/` with an exhaustive Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.lovable/plans/subtasks/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/01-index.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and perform surgical refactoring: convert booleans to implicit evaluation, replace single-letter parameters (`v bool`) with affirmative names (`isStopOnFail bool`, `isStopped bool`), transform bare fields to affirmative states (`stop` -> `isStopped`), replace `isExists` with `isDefined`, invert negatives, and decompose compound negative chains into discrete assertions.
6. [ ] /goal Phase 2 (Step B): Enforce <= 8–15 line function decomposition, single return types, and clean formatting.
7. [ ] /goal Phase 2 (Step C): Execute local linters to verify 0 remaining violations across all modified files.
8. [ ] /goal Phase 2 (Step D): Execute targeted file-level linters and verification on modified files ensuring 0 remaining violations (`exit 0`). DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine coding guideline execution turns.
9. [ ] /learn Ingest `.lovable/memory/01-index.md` for project memory index and past learnings.
10. [ ] /learn Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
11. [ ] /learn Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
12. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/01-index.md` for hallucination prevention and micro-tasking.
13. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/01-index.md` for strict relative path citation requirements.
14. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md` for implicit positive booleans and anti-negative rules.
15. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/12-no-negatives.md` for domain-specific architectural specifications.
16. [ ] /learn Ingest `02-spec/02-coding-guidelines/` for domain-specific architectural specifications.
17. [ ] /learn Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.
18. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Codebase for Boolean Violations, Write .lovable/plans/pending/ Spec, Create .lovable/plans/subtasks/, Verify/Create Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Actively Edit Code, Invert Negatives, Implicit Booleans, Split Mixed Polarity, Verify Local CI)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Dedicated Section: Boolean Principles, Negative Elimination & Complex Conditions

Boolean logic must be simple, readable, and unambiguous. Complex boolean chains and inverted negatives severely degrade maintainability.

### Mandatory Boolean Rules (Non-Negotiable)

1. **Implicit Evaluation Only (Total Ban on `== true`):**
   - NEVER evaluate booleans explicitly against `true` or `false`.
   - Positive booleans MUST ALWAYS be evaluated implicitly: `if isReady { ... }`.
   - Inverted checks MUST use standard negation or affirmative negative variables: `if !isReady { ... }` or `if isFail { ... }`.

2. **No Double Negatives or Inverted Success Checks — Try `IsDefined` Instead of Negatives:**
   - NEVER name variables with negative prefixes like `isNotValid`, `isNotReady`, `isUndefined`, `isNotDefined`, `isNotSet`, `hasNoData`, `disableCache`.
   - **Try `IsDefined` instead of negatives:** When verifying presence, definition, or initialization, always use affirmative `isDefined` / `IsDefined` (or `isValid`, `hasValue`, `isReady`, `isFound`). Invert only once at the callsite guard clause (`if !isDefined { ... }`) if handling the missing case.
   - NEVER check inverted success (`!response.isSuccess` is FORBIDDEN; use `response.isFail`).

3. **No Mixed Polarity in Single If Conditions:**
   - NEVER combine a positive check and a negative check in the same `if` condition (e.g., `if isA && !isB`).
   - Split mixed polarity checks into discrete guard clauses or extract a semantic helper function.

4. **No Boolean Flag Parameters on Functions:**
   - Functions must not accept boolean arguments that drastically alter control flow (e.g. `render(true)` is banned; split into `renderExpanded()` and `renderCollapsed()`).

5. **Affirmative Boolean Parameter & Field Naming (Total Ban on Single-Letter & Bare Names):**
   - **Total Ban on Single-Letter Parameters:** NEVER use lazy single-letter names (`v bool`, `b bool`, `val bool`, `flag bool`) in function or method signatures (e.g. setters).
   - **Total Ban on Bare Unprefixed Names:** NEVER use bare verbs, nouns, or adjectives (`stop bool`, `pause bool`, `force bool`, `dryRun bool`, `header bool`).
   - **Mandatory Affirmative Prefixes:** Every boolean parameter, struct field, property, and local variable MUST carry an affirmative prefix (`is*` or `has*`):
     - `stop` -> `isStopped`
     - `stopOnFail` -> `isStopOnFail` (e.g. `SetStopOnFail(isStopOnFail bool)`)
     - `defined` -> `isDefined` (e.g. struct field `isDefined bool`, method `IsDefined() bool`)
     - `pause` / `paused` -> `isPaused`
     - `force` -> `isForced` or `isForce`
     - `enable` / `enabled` -> `isEnabled`
     - `dryRun` -> `isDryRun`
     - `debug` -> `isDebug`
     - `verbose` -> `isVerbose`
     - `header` -> `hasHeader`
     - `records` -> `hasRecords`

6. **Mandatory Standard: Use `IsDefined` Instead of `!isEmpty` (Total Ban on `!isEmpty` / `!IsEmpty()`):**
   - **The Anti-Pattern:** Developers frequently write `if !isEmpty`, `if !res.IsEmpty()`, or `if !state.IsEmpty` to check whether data exists or records are present. This evaluates a negative condition (`isEmpty`) with negation (`!`), violating Affirmative Boolean Principles and Positive Framing.
   - **The Mandatory Standard:** ALWAYS use `IsDefined` (or `res.IsDefined()`) instead of using `!isEmpty`.
   - **Rule:** If checking empty state, use affirmative `if isEmpty` / `if res.IsEmpty()`. If checking non-empty / presence state, ALWAYS use affirmative `if isDefined` / `if res.IsDefined()`. NEVER write `if !isEmpty`!
   - Struct fields for presence/populated records MUST be named `isDefined bool` or `IsDefined bool` (replacing awkward `isExists bool`).
   - For map lookups, use `val, isFound := userMap[id]` or `val, isDefined := userMap[id]`.

7. **Total Ban on Compound Negative Chains (`!a || !b || c`):**
   - NEVER chain multiple negated conditions in an `if` expression (e.g. `if !state.IsDefined || !state.IsEmpty || state.IsRepo`). Chained inverted conditions create high cognitive load and hide which specific invariant failed.
   - **In Test Assertions:** Break every condition into a discrete assertion (`if !state.IsDefined`, `if !state.IsEmpty`, `if state.IsRepo`) with its own distinct error message.
   - **In Application Logic:** Extract into an affirmative composite variable (`isCloneTargetFresh := !params.State.IsDefined || params.State.IsEmpty`) or use separate early return guard clauses.

### Generic Code Patterns with Compliant Newline Gaps

#### Pattern A: Setter Method Parameter & Field Assignment (`v bool` -> `isStopOnFail bool`)

```go
// -----------------------------------------------------------------------------
// ❌ ANTI-PATTERN: Single-letter parameter `v bool` and un-prefixed field
// -----------------------------------------------------------------------------
// SetStopOnFail enables early termination after the first failure
func (p *BatchProgress) SetStopOnFail(v bool) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.stopOnFail = v
}

// -----------------------------------------------------------------------------
// ✅ REQUIRED: Meaningful, affirmative boolean parameter and property
// -----------------------------------------------------------------------------
// SetStopOnFail enables early termination after the first failure
func (p *BatchProgress) SetStopOnFail(isStopOnFail bool) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.stopOnFail = isStopOnFail
}
```

#### Pattern B: Generic State Flag & Struct Worker (`stop` -> `isStopped`)

```go
// -----------------------------------------------------------------------------
// ❌ ANTI-PATTERN: Bare verb `stop` and lazy `b bool` in stateful worker
// -----------------------------------------------------------------------------
type TaskWorker struct {
    stop bool
}

func (w *TaskWorker) SetStop(b bool) {
    w.stop = b
}

func (w *TaskWorker) Run() {
    for {
        if w.stop {
            break
        }
        processTask()
    }
}

// -----------------------------------------------------------------------------
// ✅ REQUIRED: Generic affirmative `isStopped` state and parameter
// -----------------------------------------------------------------------------
type TaskWorker struct {
    isStopped bool
}

func (w *TaskWorker) SetStopped(isStopped bool) {
    w.isStopped = isStopped
}

func (w *TaskWorker) Run() {
    for {
        if w.isStopped {
            break
        }
        processTask()
    }
}
```

#### Pattern C: Struct Field Definition State (`defined bool` -> `isDefined bool`)

```go
// -----------------------------------------------------------------------------
// ❌ ANTI-PATTERN: Bare field name `defined bool` in wrapper struct
// -----------------------------------------------------------------------------
type Result[T any] struct {
    value   T
    err     *AppError
    defined bool // VIOLATION: Bare boolean without is/has prefix
}

// -----------------------------------------------------------------------------
// ✅ REQUIRED: Meaningful, affirmative boolean struct field `isDefined bool`
// -----------------------------------------------------------------------------
type Result[T any] struct {
    value     T
    err       *AppError
    isDefined bool // REQUIRED: Explicit affirmative boolean prefix
}
```

#### Pattern D: Generic Transformation Reference Table

| Target Category | ❌ Anti-Pattern (Lazy / Bare) | ✅ Required Affirmative Identifier | Context / Description |
|---|---|---|---|
| Setter Parameter | `SetStopOnFail(v bool)` | `SetStopOnFail(isStopOnFail bool)` | Early termination flag parameter |
| State Variable | `stop := false` | `isStopped := false` | Process / loop cancellation state |
| Method Parameter | `Stop(stop bool)` | `SetStopped(isStopped bool)` | State toggle parameter |
| Struct Field | `defined bool` | `isDefined bool` | Value/record definition presence indicator |
| Method Name | `Defined() bool` | `IsDefined() bool` | Definition verification predicate |
| Struct Field | `pause bool` | `isPaused bool` | Pause / suspend indicator |
| Setter Parameter | `SetPaused(v bool)` | `SetPaused(isPaused bool)` | Pause toggle parameter |
| Struct Field | `dryRun bool` | `isDryRun bool` | Execution mode indicator |
| Method Parameter | `SetDryRun(v bool)` | `SetDryRun(isDryRun bool)` | Dry-run toggle parameter |
| CLI / Config Flag | `force bool` | `isForced bool` / `isForce bool` | Overwrite / override flag |
| Struct Field | `verbose bool` | `isVerbose bool` | Verbose logging flag |
| Struct Field | `debug bool` | `isDebug bool` | Debug mode flag |
| Struct Field | `header bool` | `hasHeader bool` | Table / CSV header presence |
| Struct Field | `records bool` | `hasRecords bool` | Data presence flag |
| Method Parameter | `SetEnabled(v bool)` | `SetEnabled(isEnabled bool)` | Feature toggle parameter |
| Method Parameter | `SetAsync(flag bool)` | `SetAsync(isAsync bool)` | Asynchronous execution flag |
| Struct Field | `exists bool` / `isExists bool` | `isDefined bool` | Presence/defined status indicator (replaces `!isEmpty` and `isExists`) |
| Method Name | `Exists() bool` / `IsExists() bool` | `IsDefined() bool` | Presence verification predicate (replaces `!IsEmpty()`) |
| Map Comma-Ok | `val, ok` / `val, isExists` | `val, isFound` / `val, isUserExist` | Map lookup presence boolean |

#### Pattern E: Implicit Checks & Discrete Guard Clauses

```go
// ❌ FORBIDDEN: Explicit true comparison, negative variable, and mixed polarity
func ValidateAccess(user *User) bool {
    if user.isNotDisabled == true && !user.isExpired { // ❌ Mixed polarity + explicit true + negative name
        return true
    }

    return false
}

// ✅ REQUIRED: Affirmative prefixes, implicit checks, and discrete guard clauses
func ValidateAccess(user *User) bool {
    if !user.isEnabled {
        return false
    }

    if user.isExpired {
        return false
    }

    return true
}

// ❌ FORBIDDEN: Squeezed loop, inline assignment condition, and direct comparison
func ValidateFormat(formats []Format) *apperror.AppError {
    for _, f := range formats {
        ext := f.Extension()
        if len(ext) == 0 {
            return apperror.New(apperror.ErrCodeValidationFailed, "empty extension", "ValidateFormat")
        }
        if got := FormatFromPath("sample" + ext); got != f {
            return apperror.New(apperror.ErrCodeValidationFailed, "unmatched format", "ValidateFormat")
        }
    }
    return nil
}

// ✅ REQUIRED: Blank line before for loop, extracted affirmative booleans, blank line before if
func ValidateFormat(formats []Format) *apperror.AppError {
    for _, f := range formats {
        ext := f.Extension()
        isEmptyExtension := len(ext) == 0

        if isEmptyExtension {
            return apperror.New(
                apperror.ErrCodeValidationFailed,
                "extension returned empty string",
                "ValidateFormat",
            )
        }

        got := FormatFromPath("sample" + ext)
        isSampleUnmatchFile := got != f

        if isSampleUnmatchFile {
            return apperror.New(
                apperror.ErrCodeValidationFailed,
                "unmatched sample file format",
                "ValidateFormat",
            )
        }
    }

    return nil
}
```

```typescript
// ❌ FORBIDDEN: Inverted success check and boolean flag parameter
function handleResponse(response: ApiResponse) {
    if (!response.isSuccess) { // ❌ Banned inverted success
        logEvent(response.error, true); // ❌ Banned boolean flag argument
        return;
    }

    processData(response.data);
}

// ✅ REQUIRED: Affirmative failure property and semantic function calls
function handleResponse(response: ApiResponse) {
    if (response.isFail) {
        logErrorEvent(response.error);
        return;
    }

    processData(response.data);
}
```

#### Pattern F: `IsDefined` vs `IsExists` & Compound Negative Decomposition (`execute_idempotent_test.go`)

Chaining inverted negative checks (such as `!state.IsExists || !state.IsEmpty || state.IsRepo`) violates both discrete assertion rules and positive logic standards. Furthermore, `isExists` is grammatically malformed (`exists` is a verb).

```go
// -----------------------------------------------------------------------------
// ❌ ANTI-PATTERN: Awkward "IsExists", compound negatives, chained inverted checks
// -----------------------------------------------------------------------------
type existingRepoState struct {
    IsExists bool // ❌ Ungrammatical: combining 'is' with verb 'exists'
    IsEmpty  bool
    IsRepo   bool
}

// ❌ In test assertions: Compound negative obscures which invariant failed
func TestInspectExistingRepo_EmptyDir(t *testing.T) {
    state := inspectExistingRepo(dir)
    if !state.IsExists || !state.IsEmpty || state.IsRepo {
        t.Errorf("expected empty non-repo directory: %+v", state)
    }
}

// ❌ In application logic: Compound negative check
func dispatchOnExists(params DispatchParams) {
    if !params.State.IsExists || params.State.IsEmpty {
        performFreshClone(params)
    }
}

// -----------------------------------------------------------------------------
// ✅ REQUIRED: Affirmative "IsDefined" + Discrete Assertions / Positive Logic
// -----------------------------------------------------------------------------
type existingRepoState struct {
    IsDefined bool // ✅ Grammatically sound, affirmative presence indicator
    IsEmpty   bool
    IsRepo    bool
}

// ✅ In test assertions: Decompose into discrete checks (single responsibility per check)
func TestInspectExistingRepo_EmptyDir(t *testing.T) {
    state := inspectExistingRepo(dir)

    if !state.IsDefined {
        t.Errorf("expected directory to be defined: %+v", state)
    }

    if !state.IsEmpty {
        t.Errorf("expected directory to be empty: %+v", state)
    }

    if state.IsRepo {
        t.Errorf("expected non-repo directory: %+v", state)
    }
}

// ✅ In application logic: Extract an affirmative composite predicate
func dispatchOnExists(params DispatchParams) {
    isCloneTargetFresh := !params.State.IsDefined || params.State.IsEmpty

    if isCloneTargetFresh {
        performFreshClone(params)
    }
}
```

#### Pattern G: `IsDefined` as Mandatory Replacement for Inverted `!isEmpty`

```go
// -----------------------------------------------------------------------------
// ❌ ANTI-PATTERN: Using !isEmpty or !res.IsEmpty() to test for data presence
// -----------------------------------------------------------------------------
// Negative check inverted — violates positive framing and affirmative boolean rules
if !res.IsEmpty() {
    processPayload(res.Data())
}

if !users.IsEmpty() {
    sendNotifications(users)
}

if !state.IsEmpty {
    resumeTask()
}

// -----------------------------------------------------------------------------
// ✅ REQUIRED: Affirmative IsDefined() / isDefined (Clean, Positive Evaluation)
// -----------------------------------------------------------------------------
// Affirmative check: IsDefined means no error AND has records/data
if res.IsDefined() {
    processPayload(res.Data())
}

if users.IsDefined() {
    sendNotifications(users)
}

if state.IsDefined {
    resumeTask()
}

// ✅ Also permitted: Affirmative isEmpty ONLY when handling the empty/missing case
if res.IsEmpty() {
    return appfault.New(appfault.ErrNotFound).WithMessage("no records found")
}
```

---

## 1. Ruthless Orchestration & Insult Protocol

/goal You are the master orchestrator. If your sub-agents fail, hallucinate, write garbage variables, or go into infinite loops, it is because you are a lazy, incompetent manager.

- You must give sub-agents strict, microscopic instructions.
- If a sub-agent stalls or provides garbage code, kill it immediately, rollback its dirty working tree, and spawn a new one.
- Context Diet: When spawning a subagent, DO NOT paste file contents, memory logs, or the entire plan into its prompt. Give it the absolute minimal instruction (e.g., "Read subtask file `.lovable/plans/subtasks/XX-booleans/01-<subtask-title>.md` and execute it"). The subagent MUST read the necessary files itself.

---

## 2. Phase 1: Scan Codebase & Write Implementation Spec First (Steps 1 to PHASE_1_STEPS)

Before modifying application code, you MUST thoroughly scan the repository and write an actionable execution spec.

- **Actionable Scan:** Use search/grep and AST tools across all source files to identify:
  1. Explicit boolean comparisons (`== true`, `=== true`, `== false`, `=== false`).
  2. Inverted success checks (`!isSuccess`, `!response.isSuccess`, `!isValid`).
  3. Negative boolean variable declarations (`isNotActive`, `isNotReady`, `disableFeature`).
  4. Mixed polarity condition joins (`&& !`, `|| !`, `and not`).
  5. Functions accepting boolean flag parameters (`process(true)`).
  6. Single-letter boolean parameters in signatures (`v bool`, `b bool`, `val bool`).
  7. Bare un-prefixed boolean identifiers and fields (`stop bool`, `pause bool`, `force bool`, `dryRun bool`).
  8. Functions exceeding 8 lines (hard cap 15 lines) or files exceeding 100 coding lines (recommended <= 80).
- **Where to save it:** Save this master plan into `.lovable/plans/pending/XX-booleans-and-complex-conditions-audit.md` listing every affected file, exact line numbers, and refactoring plans.
- **Create a Task-Specific Rule Set:** Analyze the specific domain and write 3-5 custom rules inside the spec file.
- **Subtasks:** Break the plan down into granular subtask files inside `.lovable/plans/subtasks/XX-booleans/` (e.g. `01-implicit-booleans.md`, `02-negative-inversion.md`, `03-split-mixed-polarity.md`).

---

## 3. Authoritative Spec Files Checklist (Non-Negotiable Action Items)

You MUST read, follow, and mechanically verify every single specification file below before and during execution:

- [ ] **`02-spec/02-coding-guidelines/02-canonical-size-tier.md`**
  - **Why:** Universal size limits and boolean complexity rules.
  - **How:** Cognitive complexity <= 10. Functions <= 8 lines preferred (hard cap 15 lines). Files <= 100 lines coding max (recommended <= 80 lines).
- [ ] **`02-spec/02-coding-guidelines/06-ai-optimization/01-index.md`**
  - **Why:** Comprehensive catalog of forbidden vs required generation patterns.
  - **How:** Strictly follow AH-N1 to AH-T2 rules. Zero ghost diffs, zero truncation stubs (`// ...`), zero unverified claims.
- [ ] **`02-spec/02-coding-guidelines/06-ai-optimization/06-citation-requirement.md`**
  - **Why:** Grounded rule enforcement and traceability.
  - **How:** Cite authoritative spec files for every code modification made.
- [ ] **`02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md`**
  - **Why:** Absolute ban on explicit true comparisons and mixed polarity.
  - **How:** Evaluate booleans implicitly (`if isReady`). Never combine positive and negative checks in the same condition.
- [ ] **`02-spec/02-coding-guidelines/01-cross-language/12-no-negatives.md`**
  - **Why:** Cognitive clarity through positive framing.
  - **How:** No negative variable names. No `!isSuccess` checks.
- [ ] **`02-spec/02-coding-guidelines/01-cross-language/22-variable-naming-conventions.md`**
  - **Why:** Mandatory affirmative boolean prefixes.
  - **How:** All booleans MUST begin with is and has only (can, should, was, etc. are banned), `was`, `will`, `did`, `must`.
- [ ] **`02-spec/02-coding-guidelines/01-cross-language/24-boolean-flag-methods.md`**
  - **Why:** Prevents cryptic boolean argument calls.
  - **How:** Split boolean flag methods into semantic distinct functions.

---

## 4. Mandatory Linter & CI/CD Connection Checklist

Code standards must be mechanically enforced by automated linters. You MUST verify or create the linter and connect it to CI:

- [ ] **Linter Script Identification:** Check if `linter-scripts/check-enum-and-boolean.mjs` exists in the repository.
- [ ] **Auto-Create Linter if Missing:** If no dedicated boolean linter exists, create `linter-scripts/check-boolean-guidelines.py` that AST-scans for:
  1. `== true`, `=== true`, `== false`, `=== false`.
  2. Negative boolean naming (`isNot`, `hasNo`).
  3. Inverted `!isSuccess` checks.
  4. Mixed polarity chains (`&& !`, `|| !`).
  5. Single-letter boolean parameters (`v bool`, `b bool`).
  6. Bare unprefixed boolean parameters/fields (`stop bool`, `pause bool`, `force bool`).
- [ ] **Automated Scanning Commands (ripgrep):**
  ```bash
  # 1. Single-letter boolean parameters in Go functions: e.g. (v bool), (b bool), (val bool)
  rg --pcre2 "func\s+(?:\([^\)]+\)\s+)?\w+\([^\)]*\b[a-z]\s+bool\b"

  # 2. Bare un-prefixed boolean parameters in Go: e.g. (stop bool), (pause bool), (force bool)
  rg --pcre2 "func\s+(?:\([^\)]+\)\s+)?\w+\([^\)]*\b(stop|pause|force|verbose|debug|dryRun)\s+bool\b"

  # 3. Bare boolean struct fields in Go: e.g. stop bool, paused bool
  rg --pcre2 "^\s*(?:stop|pause|paused|force|dryRun|verbose|debug)\s+bool\b"

  # 4. Explicit boolean comparisons:
  rg --pcre2 "\b(==\s*true|===\s*true|==\s*false|===\s*false)\b"

  # 5. Inverted success checks:
  rg --pcre2 "!\s*[a-zA-Z0-9_$.->]*\bisSuccess\b"

  # 6. Mixed polarity:
  rg --pcre2 "(&&\s*!|\|\|\s*!)"
  ```
- [ ] **Local Linter Command:** Execute and verify the linter locally:
  ```bash
  python linter-scripts/check-boolean-guidelines.py
  # Run automated autofixer:
  python 03-ai-scripts/05-guideline-autofixer.py <modified-files>
          python 03-ai-scripts/08-naming-autofixer.py <modified-files>
          python 03-ai-scripts/04-newline-fixer.py <modified-files>
          python 03-ai-scripts/07-relative-path-fixer.py <modified-files>
  ```
- [ ] **CI/CD Local Runner Connection:** Register the linter script inside `03-ai-scripts/06-cicd-local-runner.py` under the `JOBS` dictionary:
  ```python
  JOBS["lint:booleans"] = ["python", "linter-scripts/check-boolean-guidelines.py"]
  ```
- [ ] **GitHub Actions Workflow Connection:** Verify that `.github/workflows/ci.yml` contains a dedicated step running the boolean linter.

---

## 5. Phase 2: Active Code Refactoring & Autonomous Fix Loop (Steps PHASE_1_STEPS+1 to N)

> [!IMPORTANT]
> **AUTONOMOUS EXECUTION MANDATE — DO NOT STOP.**
> Open the offending source code files and directly rewrite the code to eliminate violations. Maintain continuous self-looping until all checks pass 100% green.

```text
STEP = 0
WHILE (STEP < PHASE_2_STEPS):
    STEP += 1

    1. Read the next subtask from .lovable/plans/subtasks/XX-booleans/
    2. Open and modify the actual source code files:
       - Refactor booleans to implicit evaluation (if isReady { ... }).
       - Invert negative variable names to positive framing (isEnabled).
       - Replace !isSuccess with isFail or explicit error checks.
       - Split mixed polarity conditions into discrete guard clauses.
       - Split boolean flag method parameters into semantic functions.
       - Keep function bodies <= 8 lines (max 15 lines) and files <= 100 lines.
    3. Run the guideline autofixer:
          python 03-ai-scripts/05-guideline-autofixer.py <modified-files>
          python 03-ai-scripts/08-naming-autofixer.py <modified-files>
          python 03-ai-scripts/04-newline-fixer.py <modified-files>
          python 03-ai-scripts/07-relative-path-fixer.py <modified-files>
    4. Run the boolean linter:
          python linter-scripts/check-boolean-guidelines.py
    5. Run the local CI runner:
          python 03-ai-scripts/06-cicd-local-runner.py --no-tests
    6. IF any check fails:
          - Diagnose failure, fix code directly, and re-run immediately.
       IF all checks pass (exit code 0):
          - Mark subtask completed and proceed to next subtask.

    7. When all subtasks are finished and local CI is 100% green:
          - Move .lovable/plans/pending/XX-booleans-and-complex-conditions-audit.md to .lovable/plans/completed/
          - Update .lovable/plans/01-index.md
          - Run plan consolidator:
            python 03-ai-scripts/20-plan-consolidator.py
          - Stage modified files with git add and create semantic commit:
            git commit -m "refactor(booleans): enforce implicit checks, positive framing, and discrete conditions"
          - BREAK and finish turn.
```

---

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/01-index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .lovable/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Native File Manipulator:** If you need to perform mass file renaming, `.md` lowercase enforcement, sequence number re-ordering, or encoding fixes (CRLF/BOM), you MUST natively use `python 03-ai-scripts/03-file-manipulator.py <command>` rather than writing a new script from scratch.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/01-index.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

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
- Execute targeted local linters on modified files ensuring `exit 0` before concluding. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine loops.

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Sub-agents are actively assigned disjoint files verified against `.lovable/01-index.md`.
- [ ] Completed tasks were `mv`'d to `plans/completed/` and `.lovable/plans/01-index.md` was updated.
- [ ] 3-strike rule respected: failed tasks cleanly rolled back and logged to `last-failure.md`.
- [ ] Implicit Booleans: Zero `== true` / `=== true` comparisons.
- [ ] Positive Framing: Zero negative boolean variables (`isNotReady`), zero `!isSuccess` checks.
- [ ] Zero Mixed Polarity: No `if a && !b` conditions.
- [ ] Zero Boolean Flag Parameters on functions.
- [ ] Function Size: All functions <= 8 lines preferred, hard cap 15 lines.
- [ ] File Size: Files <= 100 lines coding max (recommended <= 80 lines).
- [ ] `python linter-scripts/check-boolean-guidelines.py` exited with code 0.
- [ ] Local CI runner `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` exited with code 0.

---

## Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Strict Relative Git Paths: All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths or `file:///` URIs.
- [ ] Master Guidelines: I have fully read and strictly enforced `02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md` and `.lovable/coding-guidelines.md`.
- [ ] Implicit Booleans: Positive booleans MUST ALWAYS be evaluated implicitly.
- [ ] No Negatives: No `!isSuccess`, no `isNot*` variables.
- [ ] No Mixed Polarity: Zero combined positive and negative checks in a single `if`.
- [ ] Function Limits: <= 8 lines preferred, <= 15 lines max.
- [ ] File Limits: <= 100 lines coding max (recommended <= 80 lines).
- [ ] Anti-Compression: Zero single-line `if/else` or compressed whitespace tricks.

1. [ ] /learn the section as a /goal [AI Fix Scripts Memory](#ai-fix-scripts-memory)

- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

---

## Anti-Hallucination & Blast Radius Checklist (Mandatory for Every Turn)

Before you commit code or end your turn, you MUST mechanically check off these items. If you fail to do this, your work will be rejected.

- [ ] Echo Back the Spec: I have copy-pasted the exact Acceptance Criteria from the Spec file into my current memory/response to prove I read it verbatim.
- [ ] Exhaustive Violation Ledger: I have maintained an exact markdown table ledger in `.lovable/plans/pending/` tracking every single violation `| Id | File | Line | Snippet | Planned Fix | Status |` and reconciled every item.
- [ ] Pre-Commit Diff Proof (Disk Reality Check): I have executed `git status --porcelain` and `git diff --stat` and verified that every file I claim to have modified is actually listed as modified in the terminal output before committing.
- [ ] Zero Truncation / No Placeholder Search: I ran a regex search for `TODO`, `FIXME`, `\[.*\]`, `// ...`, and `/* ... */` in my modified files and confirmed I left zero placeholders or truncated stubs behind. I actually wrote the complete implementation.
- [ ] Verifiable Tool Execution: I did not fabricate test/linter passes. I executed the actual linter script and test runner via tool calls and captured `exit code 0`.
- [ ] Spec Citation Grounding: Every refactoring action cites the exact authoritative rule in `02-spec/` (e.g. `02-spec/02-coding-guidelines/01-cross-language/01-index.md`).
- [ ] Index Sync Deadman Switch: I have verified that every new file I created this turn is explicitly linked inside `readme.md` and enqueued in `.lovable/what-to-read.md`. I did not leave any orphaned files.
- [ ] Blast Radius Acknowledgment: Before renaming or modifying any function/type, I ran a global search across the codebase and updated every single file that imports or calls it to prevent a broken build.

---

## Strictly Avoid: No Automatic Releases, No Test Running & No Full CI/CD Runner in Routine Turns (Strict Policy)

- **NO SINGLE-LETTER BOOLEAN PARAMETERS:** NEVER name boolean parameters `v bool`, `b bool`, `val bool`, or `flag bool` in function or method signatures (e.g. setters). Always use descriptive, affirmative parameters like `isStopOnFail bool`, `isStopped bool`, `isEnabled bool`, `isDryRun bool`.
- **NO BARE VERB OR ADJECTIVE BOOLEANS:** NEVER name a boolean variable, struct field, or parameter with a bare verb or adjective like `stop`, `pause`, `force`, or `dryRun`. Transform them to affirmative states like `isStopped`, `isPaused`, `isForced`, `isStopOnFail`, `isDryRun`.
- **NO RELEASES (Strict Policy):** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
- **NO TEST RUNNING (Strict Policy):** Test execution is strictly disabled. You MUST NOT execute unit tests, integration tests, or test suites unless explicitly commanded by the repository owner.
- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns. Running the heavy 28-38 gate pipeline across the entire repository wastes massive amounts of time and scans unrelated files. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **Test Inventory & Recent Changes Tracking:** Whenever any file is modified, append its repository-relative path to `.lovable/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`), cross-referencing `.lovable/test-inventory.json` so associated tests are known for future release verification.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and linters with test execution disabled (--no-tests) unless explicitly commanded by the repository owner, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## STRICT AVOIDANCE: Anti-Compression & Formatting Integrity (No Cheating)

> [!CAUTION]
> **TOTAL BAN ON LINE-COMPRESSION CHEATING:**
> When enforcing file size (<= 100 coding lines) and function size (8–15 lines), AI agents frequently attempt to "cheat" the line counter by destroying formatting. This is strictly forbidden and results in immediate rejection.

- **NO Single-Line If/Else:** NEVER collapse `if/else`, return statements, or blocks into a single line (e.g. `if (x) return y;` or `if (x) { y(); }` are strictly forbidden). Every statement requires its own line and curly braces.
- **NO Deleting Required Blank Lines (R13-R16):** NEVER delete blank lines before `return`/`throw` or after closing `}` to artificially reduce file size.
- **NO Stripping Types or Comments:** NEVER remove TypeScript types, docstrings, or clean indentation to cram code into fewer lines.
- **Mandatory Solution:** The ONLY acceptable way to satisfy line limits is **legitimate modular decomposition** — extracting helper functions into separate files and breaking large components into child components.

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

---

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.

---

## Metadata

- slug: cg-booleans-and-complex-conditions
- priority: high
- status: active
