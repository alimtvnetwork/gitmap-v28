# Style Guidelines, Formatting & Line-Gaps — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-style`, `cg-execute style`, `audit style`, `fix formatting`, `enforce newline styling`, `flatten nested if`, `newline before if`, `return newline style`, `style guidelines audit`, `line gaps audit`, `fix line endings`, `enforce utf8 lf`, `fix function newlines`, `newline refactor`

> **Prompt Version:** 2.2.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously inventory, scan, partition, refactor, and fix all coding style, vertical newline spacing, blank line before `if`, blank line after `}`, blank line before `return`, blank lines around parameter struct instantiations and sequential function invocations, nested `if` elimination, function length (<= 8–15 lines), file size (<= 100 lines), LF line endings (`\n`), and UTF-8 (no BOM) encoding across ALL source files in the repository in bounded micro-batches of 5–8 files per subtask, running a continuous 2-agent unstoppable self-loop until 100% of codebase files are verified and refactored without stopping.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase using the fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py` with `--limit`) to inventory all architectural violations and anti-patterns without truncation.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.ai-memory/plans/pending/XX-style-guidelines-audit.md` with an exhaustive File Inventory Manifest and Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose ALL source files into granular, bounded subtask batches of **5–8 files each** in `.ai-memory/plans/subtasks/XX-style/batch-01.md`, `batch-02.md`, etc.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated style autofixer in `03-ai-scripts/05-guideline-autofixer.py` and register in `03-ai-scripts/01-index.md`.
5. [ ] /goal Phase 2 (Step A): Spawn 2 execution subagents (max 2 threads each) to process subtasks concurrently, opening and surgically editing each 5–8 file batch line-by-line.
6. [ ] /goal Phase 2 (Step B): Enforce Return New Line rules (R13-R16): blank line before `if`, blank line after `}`, blank line before `return`, blank lines around multiline struct calls, and zero clumped guard clauses.
7. [ ] /goal Phase 2 (Step C): Decompose functions exceeding 8–15 lines into focused single-responsibility helpers and flatten nested conditionals (depth 0).
8. [ ] /goal Phase 2 (Step D): Verify that actual source files (`*.go`, `*.ts`, etc.) have real modifications via `git diff --stat` (auto-reject if only `.ai-memory/` markdown files were changed).
9. [ ] /goal Phase 2 (Step E): Move completed batch subtasks to `.ai-memory/plans/completed/` and immediately self-loop to dispatch the next pending batches until 0 batches remain.
10. [ ] /goal Phase 2 (Step F): Execute local linters (`python linter-scripts/check-newline-styling.py`, `check-function-lengths.py`) to verify 0 remaining violations.
11. [ ] /goal Phase 2 (Step G): Execute targeted file-level linters and verification on modified files ensuring 0 remaining violations (`exit 0`). DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine coding guideline execution turns.
12. [ ] /learn Ingest `.ai-memory/memory/01-index.md` for project memory index and past learnings.
13. [ ] /learn Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
14. [ ] /learn Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
15. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/01-index.md` for hallucination prevention and micro-tasking.
16. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/01-index.md` for strict relative path citation requirements.
17. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/04-code-style/` for domain-specific architectural specifications.
18. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/21-newline-styling-examples.md` for newline styling examples.
19. [ ] /learn Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.
20. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Inventory ALL Source Files, Partition into 5-8 File Batches in .ai-memory/plans/subtasks/, Verify Autofixer)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Unstoppable 2-Agent Loop over All Batches, Surgical Function Edits, Real Git Diffs, Local CI Runner)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## AI Diagnostic Guide: Why Newline Violations are Missed in Functions & How to Fix Them

> [!IMPORTANT]
> **CRITICAL FAILURE ROOT CAUSE (PREVENTING PREMATURE COMPLETION):**
>
> In past runs, AI agents failed to fix function newlines because they:
>
> 1. **Only ran surface linters** (e.g. checking file-level EOF newlines) without opening individual function bodies.
> 2. **Wrote a planning markdown file and immediately declared completion** without editing a single `.go` or `.ts` file.
> 3. **Attempted to process the entire codebase in one prompt**, causing context exhaustion and truncated file edits.
>
> **THE MANDATORY REMEDY:**
>
> - You MUST partition the full list of codebase files into **batches of 5–8 files each**.
> - Subagents MUST open and edit every single file in their assigned batch line-by-line.
> - The master orchestrator MUST continuously self-loop across all batches until every single batch in `.ai-memory/plans/subtasks/` is completed.

---

## Dedicated Section: Comprehensive Coding Style, Line-Gaps & Anti-Pattern Gallery (Zero Tolerance)

Proper vertical spacing and code hygiene are essential for readability and automated analysis. Dense, squeezed code without blank lines around control flow statements leads to missed edge cases, obscured invariants, and severe cognitive fatigue.

---

### Rule 1: Mandatory Blank Line BEFORE Control Structures (`if`, `for`, `switch`, `while`, `try`)

Whenever a control structure (`if`, `for`, `switch`, `while`, `try`) is preceded by **any statement** (variable declaration, assignment, method call, channel receive, or loop), there **MUST be exactly one blank line before the control structure**.

*Exception:* If the control structure is the **very first line** of a function body or immediately follows an opening brace `{`, no blank line is required before it.

#### 1a. Go: Blank Line Before `if`

```go
// ❌ WRONG: Squeezed variable declaration / map lookup directly against if
func ProcessUser(id string) error {
    user, isFound := userCache.Get(id)
    if !isFound {
        return ErrUserNotFound
    }
    config, isLoaded := loadConfig()
    if !isLoaded {
        return ErrConfigMissing
    }
    return executeUser(user, config)
}

// ✅ CORRECT: Clean blank line before each if statement, after each closing brace, and before return
func ProcessUser(id string) error {
    user, isFound := userCache.Get(id)

    if !isFound {
        return ErrUserNotFound
    }

    config, isLoaded := loadConfig()

    if !isLoaded {
        return ErrConfigMissing
    }

    return executeUser(user, config)
}
```

---

#### 1b. TypeScript / React: Blank Line Before `if`

```typescript
// ❌ WRONG: Squeezed variable declarations and function calls against if
function getFormattedPrice(item: Item): string {
    const rawPrice = calculateBasePrice(item);
    const hasDiscount = item.discountPercent > 0;
    if (hasDiscount) {
        return applyDiscount(rawPrice, item.discountPercent);
    }
    const formatted = formatCurrency(rawPrice);
    return formatted;
}

// ✅ CORRECT: Clean blank line before if and before final return
function getFormattedPrice(item: Item): string {
    const rawPrice = calculateBasePrice(item);
    const hasDiscount = item.discountPercent > 0;

    if (hasDiscount) {
        return applyDiscount(rawPrice, item.discountPercent);
    }

    const formatted = formatCurrency(rawPrice);

    return formatted;
}
```

---

#### 1c. Python: Blank Line Before `if`

```python
# ❌ WRONG: Assignment directly followed by if without blank line
def fetch_user_profile(user_id: str) -> Profile:
    auth_token = get_session_token()
    is_valid_token = verify_token(auth_token)
    if not is_valid_token:
        raise UnauthorizedError()
    user_record = db.find_user(user_id)
    if user_record is None:
        raise NotFoundError()
    return Profile.from_record(user_record)

# ✅ CORRECT: Clean blank line before if and before returns
def fetch_user_profile(user_id: str) -> Profile:
    auth_token = get_session_token()
    is_valid_token = verify_token(auth_token)

    if not is_valid_token:
        raise UnauthorizedError()

    user_record = db.find_user(user_id)

    if user_record is None:
        raise NotFoundError()

    return Profile.from_record(user_record)
```

---

#### 1d. PHP: Blank Line Before `if` and `foreach`

```php
// ❌ WRONG: Squeezed statements before if and foreach
$result = $this->apiRequest($agentId, HttpMethodType::Post->value, $endpoint);
if (is_wp_error($result)) {
    return $result;
}
$items = $this->fetchItems();
foreach ($items as $item) {
    $this->process($item);
}

// ✅ CORRECT: Separated with blank lines before control structures
$result = $this->apiRequest($agentId, HttpMethodType::Post->value, $endpoint);

if (is_wp_error($result)) {
    return $result;
}

$items = $this->fetchItems();

foreach ($items as $item) {
    $this->process($item);
}
```

---

### Rule 2: Mandatory Blank Line AFTER Closing Brace `}` When Followed by Code

Whenever a closing brace `}` (from an `if`, `for`, `switch`, `while`, or `try/catch` block) is followed by further executable code or another statement, there **MUST be exactly one blank line after `}`**.

*Exception:* No blank line is needed when `}` is followed by another closing `}`, `else`, `catch`, `finally`, or the end of a function body.

#### 2a. Go: Blank Line After `}` Following Control Flow

```go
// ❌ WRONG: Closing brace squeezed against next statement
func ExecuteStep(step Step) error {
    if err := step.Validate(); err != nil {
        return err
    }
    result, err := step.Run()
    if err != nil {
        return err
    }
    return saveResult(result)
}

// ✅ CORRECT: Clean blank line after every closing brace
func ExecuteStep(step Step) error {
    if err := step.Validate(); err != nil {
        return err
    }

    result, err := step.Run()
    if err != nil {
        return err
    }

    return saveResult(result)
}
```

---

#### 2b. TypeScript: Blank Line After Loops & Try/Catch

```typescript
// ❌ WRONG: Loop and try/catch squeezed against subsequent logic
for (const item of items) {
    processed.push(transform(item));
}
const result = merge(processed);

try {
    saveToStorage(result);
} catch (error) {
    logger.error(error);
}
cleanup();

// ✅ CORRECT: Blank line after each closing brace
for (const item of items) {
    processed.push(transform(item));
}

const result = merge(processed);

try {
    saveToStorage(result);
} catch (error) {
    logger.error(error);
}

cleanup();
```

---

#### 2c. Go: Blank Line After Multiline Map/Slice Literals & Loops with Boolean Extraction

```go
// ❌ FORBIDDEN (Unacceptable): Squeezed loops against map literals, inline conditional assignments, and missing blank lines between if blocks
func ValidateDoubleExtensionFormats(targetPath string) *apperror.AppError {
    cases := map[string]Format{
        "archive.tar.gz":  FormatTarGz,
        "archive.tgz":     FormatTarGz,
        "archive.tar.bz2": FormatTarBz2,
        "archive.tbz2":    FormatTarBz2,
        "archive.tar.xz":  FormatTarXz,
        "archive.txz":     FormatTarXz,
        "archive.tar.zst": FormatTarZst,
        "archive.tzst":    FormatTarZst,
    }
    for path, expectedFormat := range cases {
        if got := FormatFromPath(path); got != expectedFormat {
            return apperror.New(apperror.ErrCodeValidationFailed, "mismatch", "ValidateDoubleExtensionFormats")
        }
    }
    return nil
}

func ValidateExtensionRoundTrip(formats []Format) *apperror.AppError {
    for _, f := range formats {
        ext := f.Extension()
        if len(ext) == 0 {
            return apperror.New(apperror.ErrCodeValidationFailed, "empty extension", "ValidateExtensionRoundTrip")
        }
        if got := FormatFromPath("sample" + ext); got != f {
            return apperror.New(apperror.ErrCodeValidationFailed, "unmatched format", "ValidateExtensionRoundTrip")
        }
    }
    return nil
}

// ✅ REQUIRED (Right Practice): Blank line after map literal closing brace, blank line before loops, blank line before if, blank line after closing brace, and extracted affirmative booleans
func ValidateDoubleExtensionFormats(targetPath string) *apperror.AppError {
    cases := map[string]Format{
        "archive.tar.gz":  FormatTarGz,
        "archive.tgz":     FormatTarGz,
        "archive.tar.bz2": FormatTarBz2,
        "archive.tbz2":    FormatTarBz2,
        "archive.tar.xz":  FormatTarXz,
        "archive.txz":     FormatTarXz,
        "archive.tar.zst": FormatTarZst,
        "archive.tzst":    FormatTarZst,
    }

    for path, expectedFormat := range cases {
        resolvedFormat := FormatFromPath(path)
        isFormatMismatch := resolvedFormat != expectedFormat

        if isFormatMismatch {
            return apperror.New(
                apperror.ErrCodeValidationFailed,
                "format mismatch detected",
                "ValidateDoubleExtensionFormats",
            )
        }
    }

    return nil
}

func ValidateExtensionRoundTrip(formats []Format) *apperror.AppError {
    for _, f := range formats {
        ext := f.Extension()
        isEmptyExtension := len(ext) == 0

        if isEmptyExtension {
            return apperror.New(
                apperror.ErrCodeValidationFailed,
                "extension returned empty string",
                "ValidateExtensionRoundTrip",
            )
        }

        got := FormatFromPath("sample" + ext)
        isSampleUnmatchFile := got != f

        if isSampleUnmatchFile {
            return apperror.New(
                apperror.ErrCodeValidationFailed,
                "unmatched sample file format",
                "ValidateExtensionRoundTrip",
            )
        }
    }

    return nil
}
```

---

#### 2d. Go: Blank Lines Around Struct Instantiations & Sequential Function Invocations

When instantiating a parameter struct or invoking a multi-line function, there **MUST be a blank line before the invocation** (if preceded by assignments or statements) and **MUST be a blank line after the invocation closing brace `}`** before subsequent statements, `if` conditions, or other function calls.

```go
// ❌ FORBIDDEN (Unacceptable): Squeezing variable assignments, multiline struct invocations, and following if statements without vertical line gaps
func PrintIdentityBlock(cwd string) {
    fmt.Println(" " + constants.ColorCyan + "Identity Block" + constants.ColorReset)
    src := getSourceDirectory()
    emitIdentityRows(IdentityRowParams{
        Dir:            src,
        RepoOverride:   buildRepo,
        BranchOverride: buildBranch,
        ShaOverride:    buildCommit,
    })
    if len(buildDate) > 0 {
        fmt.Printf(" Built: %s\n", buildDate)
    }
    emitIdentityRows(IdentityRowParams{
        Dir: cwd,
    })
    fmt.Println()
}

// ✅ REQUIRED (Right Practice): Clean blank lines before and after multiline struct calls, separating discrete execution stages
func PrintIdentityBlock(cwd string) {
    fmt.Println(" " + constants.ColorCyan + "Identity Block" + constants.ColorReset)

    src := getSourceDirectory()

    emitIdentityRows(IdentityRowParams{
        Dir:            src,
        RepoOverride:   buildRepo,
        BranchOverride: buildBranch,
        ShaOverride:    buildCommit,
    })

    if len(buildDate) > 0 {
        fmt.Printf(" Built: %s\n", buildDate)
    }

    emitIdentityRows(IdentityRowParams{
        Dir: cwd,
    })

    fmt.Println()
}
```

---

### Rule 3: Mandatory Blank Line BEFORE `return`, `throw`, `raise`, `yield`

In multi-line functions and blocks, there **MUST be a blank line before `return` / `throw` / `raise` / `yield`**.

*Exception:* Single-statement function body (`func GetId() string { return c.Id }`) or when `return` is the immediate first statement inside a block.

```go
// ❌ WRONG: Return squeezed directly under statements
func CalculateTotal(items []Item, taxRate float64) float64 {
    subtotal := computeSubtotal(items)
    tax := subtotal * taxRate
    total := subtotal + tax
    return total
}

// ✅ CORRECT: Blank line before final return
func CalculateTotal(items []Item, taxRate float64) float64 {
    subtotal := computeSubtotal(items)
    tax := subtotal * taxRate
    total := subtotal + tax

    return total
}
```

---

### Rule 4: Zero Clumping of Consecutive Guard Clauses

When multiple guard clauses follow one another sequentially, **each guard clause MUST be separated by a blank line** after its closing brace `}`. Never clump or stack guard clauses together without vertical breathing room.

```go
// ❌ WRONG: Clumped guard clauses with zero spacing
func ValidateOrder(order *Order) error {
    if order == nil {
        return ErrNilOrder
    }
    if !order.HasItems() {
        return ErrEmptyOrder
    }
    if order.TotalAmount <= 0 {
        return ErrInvalidAmount
    }
    return nil
}

// ✅ CORRECT: Vertical breathing room between discrete guard clauses
func ValidateOrder(order *Order) error {
    if order == nil {
        return ErrNilOrder
    }

    if !order.HasItems() {
        return ErrEmptyOrder
    }

    if order.TotalAmount <= 0 {
        return ErrInvalidAmount
    }

    return nil
}
```

---

### Rule 5: Zero Nested `if` Statements (Mandatory Flattening to Depth 0)

Conditionals MUST NEVER exceed depth 1 (i.e., **no nested `if` statements inside another `if` block**). Flatten all branching logic using early guard returns or discrete helper functions.

---

### Rule 6: No Multi-Statement Lines & No Inline Compound Condition Cramming (Multi-Line Separation)

1. **No Semicolon Packing:** Never compress multiple statements onto a single line using semicolons (`a = 1; b = 2; return a + b`). Each statement MUST occupy its own line.
2. **Total Ban on Inline Compound Assignments (`if init; cond`):** NEVER cram variable declarations, type assertions, or multi-part boolean checks into the `if` header (e.g. `if v, isString := rawMap[key].(string); isString && len(v) > 0 {`).
3. **Multi-Line Statement Separation:**
   - Execute variable assignments / lookups on their own dedicated line.
   - Evaluate and assign the boolean condition to an affirmative variable (`is*` or `has*`) on its own dedicated line *before* the `if` statement.
   - Maintain vertical breathing room (blank line before `if`).
   - The `if` condition itself must be dead simple, checking **one single variable**.

#### Canonical Example: What NOT to Do vs What to Do

```go
// ❌ BANNED ANTI-PATTERN:
// Squeezing type assertion assignment and compound condition into one line to save vertical space:
func extractVersionValue(rawMap map[string]interface{}) string {
    for _, key := range []string{"Version", "version"} {
        if v, isString := rawMap[key].(string); isString && len(v) > 0 {
            return v
        }
    }

    return ""
}

// ✅ MANDATORY CLEAN PATTERN:
// 1. Assignment on its own dedicated line.
// 2. Boolean evaluated and named affirmatively BEFORE the if statement.
// 3. Clean vertical spacing (blank line before if).
// 4. if condition is dead-simple, evaluating exactly ONE variable.
func extractVersionValue(rawMap map[string]interface{}) string {
    for _, key := range []string{"Version", "version"} {
        v, isString := rawMap[key].(string)
        hasContent := isString && len(v) > 0

        if hasContent {
            return v
        }
    }

    return ""
}
```

---

### Rule 7: Universal File Hygiene, Line Endings (LF `\n` Only) & Encoding (UTF-8 No BOM)

1. **Unix LF (`\n`) Line Endings Only:** Every file MUST use Unix LF (`\n`). Total ban on Windows CRLF (`\r\n`).
2. **Strict UTF-8 Encoding (NO BOM):** Save all files in UTF-8 without BOM.
3. **Mandatory Single Trailing Newline at EOF:** Exactly one newline at the end of every file.

---

## Continuous 2-Phase Self-Loop & 2-Agent Concurrency Architecture

To guarantee full execution without stopping after planning mode, the master orchestrator MUST enforce this continuous 2-phase loop:

### 1. 2-Agent Concurrency & Strict `.ai-memory/` Bounding

- **2-Agent Limit (Max 2 Threads Each):** When dispatching work, spawn **at most 2 sub-agents concurrently**, with **no more than 2 threads per agent**.
- **Strict Folder Bounding (`.ai-memory/`):** Subagents can ONLY write planning files, subtasks, status reports, and logs inside `.ai-memory/` (`.ai-memory/plans/`, `.ai-memory/01-index.md`, `.ai-memory/memory/issues/`).
- **Context Diet:** Provide subagents with minimal instructions (e.g. "Read subtask file `.ai-memory/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md` and execute it"). Do not paste huge files into agent prompts.

### 2. Phase 1: Planning Mode & Micro-Batch Subtask Partitioning (Steps 1 .. N/2)

### Fast File Discovery & Reading via Python Toolchain (Mandatory Acceleration)

To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, the AI agent MUST use the repository's dedicated Python discovery scripts first:

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

1. **Comprehensive File Inventory:** Scan and list EVERY single source code file in the repository (`*.go`, `*.ts`, `*.py`, `*.php`).
2. **Partition into 5–8 File Batches:** Group the file list into numbered subtasks:
   - `.ai-memory/plans/subtasks/XX-style/batch-01.md`: Files 1–8
   - `.ai-memory/plans/subtasks/XX-style/batch-02.md`: Files 9–16
   - `.ai-memory/plans/subtasks/XX-style/batch-03.md`: Files 17–24
   - ... (continue until all files in the codebase are assigned to a batch).
3. **MANDATORY AUTO-LOOP (DO NOT STOP):** Once Phase 1 subtasks are written, the master orchestrator **MUST NOT STOP or ask the user for confirmation**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

### 3. Phase 2: Unstoppable Execution Mode & Parallel Batch Refactoring (Steps N/2+1 .. N)

1. **Parallel 2-Agent Dispatch:**
   - Spawn Subagent 1 on `batch-01.md` (Files 1–8).
   - Spawn Subagent 2 on `batch-02.md` (Files 9–16).
2. **Surgical Line-by-Line Refactoring:**
   - Subagents open each file in their batch, examine all function bodies, and apply vertical line gaps before `if`, after `}`, before `for`, before `return`, and around struct calls.
   - Run `python 03-ai-scripts/05-guideline-autofixer.py <file>` to automate newline insertions.
3. **Anti-Cheating Reality Check:**
   - Subagents MUST verify that actual source files (`*.go`, `*.ts`, etc.) were edited via `git diff --stat`.
   - If 0 source code files were modified, the batch is rejected as a hallucination.
4. **Continuous Self-Looping:**
   - Move completed batch subtasks from `.ai-memory/plans/subtasks/` to `.ai-memory/plans/completed/`.
   - Orchestrator checks for remaining pending batches. If any exist, immediately self-loop and dispatch the next 2 batches (`batch-03`, `batch-04`).
   - **DO NOT STOP until ALL batches are in `plans/completed/` and 0 pending batches remain.**
5. **Quality Gate Verification:**
   - Execute targeted local linters on modified files ensuring `exit 0` before concluding. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine loops.

---

## Strict In-Repository Execution & `.ai-memory/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.ai-memory/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.ai-memory/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.ai-memory/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`, `07-relative-path-fixer.py`, `05-naming-autofixer.py`).
>    - RCA & Issue Logs: `.ai-memory/memory/issues/` and `.ai-memory/cicd-issues/`.
>    - Execution Plans & Subtasks: `.ai-memory/plans/pending/`, `.ai-memory/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.ai-memory/coding-guidelines.md`.
> 3. **Worker Pool & Log Aggregation Architecture:** All local runners and test orchestrators must use a concurrent worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without canceling sibling workers, and print a consolidated final summary with full stdout/stderr error logs for failed jobs.
> 4. **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script from scratch**, regardless of whether the file already exists on disk.
> 5. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.

---

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Sub-agents are actively assigned disjoint batches verified against `.ai-memory/01-index.md`.
- [ ] **Micro-Batch Sizing:** Each subtask is bounded to exactly 5–8 files.
- [ ] **Real Source Edits:** Verified with `git diff --stat` that actual source code files (`*.go`, `*.ts`, etc.) have newline insertions.
- [ ] Completed batch tasks were `mv`'d to `plans/completed/` and `.ai-memory/plans/01-index.md` was updated.
- [ ] **Blank Line Before `if`:** Verified blank line before every `if` statement across all modified files.
- [ ] **Blank Line After `}`:** Verified blank line after every closing brace `}` followed by code.
- [ ] **Blank Line Before `return`:** Verified blank line before every `return`/`throw` in multi-line blocks.
- [ ] **Blank Lines Around Struct Invocations:** Verified blank lines before and after multiline struct calls and loops.
- [ ] **Zero Clumped Guard Clauses:** Consecutive `if` statements separated by blank lines.
- [ ] **Zero Nested `if`:** All conditionals flattened to depth 0.
- [ ] **LF Line Endings (`\n`):** All files use Unix LF line endings.
- [ ] **UTF-8 Encoding (No BOM):** All files encoded in UTF-8 without BOM.
- [ ] `python linter-scripts/check-newline-styling.py` and `python linter-scripts/check-function-lengths.py` exited with code 0.
- [ ] **Atomic File Recording:** Modified files recorded to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`). Builds, tests, and CI runner are deferred to CI/CD fix.

---

## Mandatory Linter & CI/CD Integration

1. **Linter Scripts:** `linter-scripts/check-function-lengths.py`, `linter-scripts/check-mws-error-codes.py`, `linter-scripts/check-newline-styling.py`
2. **Local Run Command:** `python linter-scripts/check-function-lengths.py`
3. **Autofixer Command:** `python 03-ai-scripts/05-guideline-autofixer.py <file>`
4. **CI/CD Integration (`.github/workflows/ci.yml`):**
   ```yaml
   - name: Validate Newline Styling & Function Lengths
     run: |
       python linter-scripts/check-function-lengths.py
       python linter-scripts/check-mws-error-codes.py
       python linter-scripts/check-newline-styling.py
   ```
5. **Runner Registration (`03-ai-scripts/06-cicd-local-runner.py`):**
   ```python
   JOBS = {
       "Newline Styling Check": [sys.executable, "linter-scripts/check-newline-styling.py"],
       "Function Lengths Check": [sys.executable, "linter-scripts/check-function-lengths.py"],
       "Error Codes Check": [sys.executable, "linter-scripts/check-mws-error-codes.py"],
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

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

---

## Strictly Avoid: No Automatic Releases, No Test Running & No Full CI/CD Runner in Routine Turns (Strict Policy)

- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- **NO RELEASES (Strict Policy):** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
- **NO TEST RUNNING (Strict Policy):** Test execution is strictly disabled. You MUST NOT execute unit tests, integration tests, or test suites unless explicitly commanded by the repository owner.
- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns. Running the heavy 28-38 gate pipeline across the entire repository wastes massive amounts of time and scans unrelated files. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **Test Inventory & Recent Changes Tracking:** Whenever any file is modified, append its repository-relative path to `.ai-memory/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`), cross-referencing `.ai-memory/test-inventory.json` so associated tests are known for future release verification.
