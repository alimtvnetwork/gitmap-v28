# File Size & Function Size Reduction — Coding Guideline Execution (must follow)

Trigger Keywords & Aliases: `cg-size-reduction`, `cg-file-reduction`, `cg-function-reduction`, `cg-execute size`, `reduce file size`, `split large files`, `decompose functions`, `audit file sizes`

> **Prompt Version:** 2.2.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform (configurable per run).

/goal Autonomously scan, plan, decompose, and refactor files exceeding the 100-line cap (recommended <= 80 lines) and functions exceeding 8–15 lines across the codebase. Enforce a two-part decomposition strategy (functions first, then files), preserve all formatting and whitespace (zero line-compression cheating), utilize wrapper objects for multi-value returns, enforce boolean conventions, and defer build verification strictly to the final step without running intermediate tests or builds.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

You MUST execute this task via a strict 3-Phase pipeline governed by the N-step budget. Do not skip steps.

```text
N = 200  (Total self-loop steps budget, read-only after initialization)
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Codebase, Map Violations, Plan Decompositions in .ai-memory/plans/pending/)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Function & File Extraction, Subtasks, Wrapper Objects, Final Build Fix)
```

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A - Discovery & Inventory): Deeply scan the target codebase using the fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py` with `--limit`) to inventory all architectural violations and anti-patterns without truncation.
2. [ ] /goal Phase 1 (Step B - Master Plan Generation): Write the master architectural specification into `.ai-memory/plans/pending/xx-size-reduction.md` with an exhaustive Violation Ledger table (File, Initial Lines, Functions to Extract, Planned Destination Files, Wrapper Structs Needed).
3. [ ] /goal Phase 1 (Step C - Subtask Decomposition): Decompose the master plan into lean, single-responsibility subtask files in `.ai-memory/plans/subtasks/xx-size-reduction/01-<subtask>.md`, `02-<subtask>.md`, etc.
4. [ ] /goal Phase 1 (Step D - Mandatory Auto-Loop): As soon as Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for permission**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.
5. [ ] /goal Phase 2 (Step A - Two-Part Decomposition): For each target file, execute Part 1 (decompose large functions to <= 8 lines, max 15 lines) and Part 2 (extract decomposed functions and helpers into separate sibling files to bring the file under 100 lines).
6. [ ] /goal Phase 2 (Step B - Wrapper Objects & Clean Signatures): If an extracted function returns multiple interrelated values or requires >2–3 parameters, encapsulate them into a dedicated wrapper object/struct.
7. [ ] /goal Phase 2 (Step C - Zero Line-Compression & Style Preservation): Strictly preserve all blank lines, block separation, and indentation. NEVER delete whitespace, collapse `if/else`, or merge statements to artificially lower line counts.
8. [ ] /goal Phase 2 (Step D - Boolean & Control Flow Concurrency): Enforce positive boolean naming (`is*`, `has*`), zero explicit `== true`, zero negative polarity in conditionals, and flatten nested `if` statements to depth <= 1 using guard clauses.
9. [ ] /goal Phase 2 (Step E - Banned Intermediate Verification): DO NOT run unit tests (`go test`, `pytest`, npm test) and DO NOT verify builds during intermediate micro-refactoring steps.
10. [ ] /goal Phase 2 (Step F - Final Step Build Verification): At the conclusion of all refactoring subtasks, run targeted syntax/build checks to resolve any compilation errors or import issues across all modified files.
11. [ ] /goal Phase 3 (Step A - Task Consolidation): Consolidate all completed subtasks into `.ai-memory/plans/completed/xx-size-reduction.md`, delete granular subtask files, and update `.ai-memory/plans/01-index.md`.
12. [ ] /goal Phase 3 (Step B - Final Step Git Commit & Push): Stage all modified files, consolidated plans, and memory records (`git add -A`), commit them in a single clean grouped atomic commit, and push to git. Never commit per-file.
13. [ ] /learn Ingest `.ai-memory/memory/01-index.md` for project memory index and past learnings.
14. [ ] /learn Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
15. [ ] /learn Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical size tiers.
16. [ ] /learn Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.

---

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
   python 03-ai-scripts/13-file-size-guard.py
   ```
Do not rely on standard search tools with 50-item truncation when discovering repository-wide violations.

---

## 1. Canonical Size Limits & Allowed Exceptions

### A. Strict File Size Caps
- **Target File Size:** Under 100 lines of code (recommended <= 80 lines).
- **Target Function Size:** <= 8 lines preferred (hard cap of <= 15 lines). In extreme cases (e.g. exhaustive codegen switch), maximum <= 25 lines with waiver, but always strive for <= 8 lines.

### B. Allowed File Exceptions (Exempt from 100-Line Limit)
The following file categories are explicitly permitted to exceed 100 lines and MUST NOT be decomposed artificially:
1. **Data & Schema Files:** JSON files (`*.json`), YAML/TOML configuration manifests.
2. **Type Definition Files:** Dedicated types files (`types.go`, `*.d.ts`, `types.ts`, `models.go`, `schema.ts`) that centralize exported domain models, structs, or discriminated unions.
3. **Lookup & Map Files:** Static lookup tables, dictionary mappings, registries, and transition matrices (`map.go`, `registry.go`, `lookup.ts`, `routes.go`).
4. **Variables & Constants Files:** Centralized constant packages, enum registrations, and configuration variable dictionaries (`vars.go`, `consts.go`, `constants.ts`, `vars.ts`, `tokens.ts`).

All other source files (logic, services, controllers, handlers, utilities, hooks, CLI commands, scripts) **MUST** strictly adhere to the <= 100 line cap.

---

## 2. Zero Line-Compression Cheating (Non-Negotiable)

> [!CAUTION]
> **TOTAL BAN ON ARTIFICIAL LINE COMPRESSION & WHITESPACE REMOVAL:**
> Never attempt to meet the 100-line file cap or 8-line function cap by deleting blank lines, merging multiple statements onto a single line, or writing one-line `if/else` blocks. Doing so violates repository style guidelines and is an immediate auto-reject failure.

### Strict Formatting Preservation Invariants:
1. **Return New Line Concept (Mandatory):**
   - Always leave exactly **ONE blank line BEFORE** every `return`, `throw`, or `break` statement (unless it is the only statement in a block).
   - Always leave exactly **ONE blank line AFTER** every closing curly brace `}` of an `if`, `for`, `switch`, or helper block.
2. **No Multi-Statement Lines:** Every statement, assignment, condition, and return statement must reside on its own dedicated line.
3. **No Squished Functions:** Every function body starts immediately on line 1 after `{` (no blank line on line 1), but must maintain clean vertical breathing room between logical steps.
4. **Decomposition Over Compression:** If a file has 140 lines, you MUST extract cohesive logic into a sibling file (e.g., `parser.go` -> `parser_helpers.go` or `parser_validator.go`), NOT remove 40 blank lines!

### Total Ban on Compound Inline `if init; cond` Cramming (The Multi-Line Separation Mandate)

> [!CAUTION]
> **NEVER CRAM ASSIGNMENTS AND COMPOUND CONDITIONS INTO THE `if` HEADER TO FAKE FILE SIZE REDUCTION:**
> A frequent anti-pattern when attempting to meet the 100-line limit is compressing variable declarations, type assertions, and multiple boolean checks into a single compound line (e.g. `if v, isString := rawMap[key].(string); isString && len(v) > 0 {`).
>
> This is **strictly prohibited**. It severely obscures cognitive readability, mixes multiple concerns, and breaks the single-condition `if` rule. Statements MUST be placed on separate lines, booleans evaluated **BEFORE** the `if`, and the `if` condition kept simple with a single variable check.

#### The Canonical Example: What NOT to Do vs What to Do

```go
// ❌ BANNED ANTI-PATTERN:
// Cramming type assertion assignment and compound condition into one line to save vertical lines:
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
// 2. Boolean evaluation named affirmatively and computed BEFORE the if statement.
// 3. Clean vertical spacing (blank line before if).
// 4. if condition is dead-simple, checking exactly ONE variable.
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

File size and function size reduction MUST be achieved through **modular decomposition** (extracting cohesive helpers into sibling files), NEVER by squishing statements onto fewer lines.

---

## 3. Two-Part Decomposition Strategy

Refactoring a large file must be performed in two sequential phases:

```
┌───────────────────────────────────────────────────────────┐
│ PART 1: Function-Level Analysis & Decomposition           │
│ 1. Scan all functions in the target file.                 │
│ 2. Decompose functions > 15 lines down to <= 8 lines.     │
│ 3. Extract sub-steps, validation, formatting to helpers.  │
└─────────────────────────────┬─────────────────────────────┘
                              │
                              ▼
┌───────────────────────────────────────────────────────────┐
│ PART 2: File-Level Extraction & Sibling Modularization    │
│ 1. Count remaining total lines in the file.               │
│ 2. If > 100 lines, group extracted helpers by capability. │
│ 3. Move cohesive groups to dedicated sibling files.       │
│ 4. Verify parent file is now under 100 lines.             │
└───────────────────────────────────────────────────────────┘
```

### Part 1: Function Decomposition (Target: <= 8 lines, Max: 15 lines)
- Identify long functions and extract sub-operations into dedicated, private helper functions.
- Replace complex multi-stage loops or deep logic with readable pipeline calls:
  ```go
  // ❌ BEFORE (35 lines: monolithic function doing parsing, validation, and storage)
  func ProcessOrder(order Order) error { ... }

  // ✅ AFTER (8 lines: orchestrates small decomposed functions)
  func ProcessOrder(order Order) *appfault.AppError {
      if isInvalid := validateOrder(order); isInvalid {
          return appfault.New("invalid order")
      }

      payload := prepareOrderPayload(order)
      return saveOrder(payload)
  }
  ```

### Part 2: File Modularization (Target: < 100 lines)
- When a file has multiple decomposed functions that keep the total lines above 100, extract cohesive clusters into new sibling files:
  - Validation logic -> `<name>_validator.go` (or `.ts`)
  - Transformation/mapping -> `<name>_converter.go` (or `.ts`)
  - State/helpers -> `<name>_helpers.go` (or `.ts`)
- Maintain package cohesion: extracted files reside in the same package/directory to avoid unnecessary export friction.

---

## 4. Wrapper Objects & Struct Encapsulation Pattern

When extracting functions, passing or returning multiple distinct values leads to argument bloat and fragile signatures.

### The Wrapper Object Mandate:
1. **Return Wrapper Struct:** If an extracted helper returns 2 or more related values (beyond standard error/Result), encapsulate them into a dedicated named struct or type:
   ```go
   // ❌ BAD: returning multiple raw values
   func parseHeader(raw string) (string, int, bool, *appfault.AppError)

   // ✅ GOOD: encapsulate in a clean wrapper struct
   type HeaderMeta struct {
       ContentType   string
       ContentLength int
       IsChunked     bool
   }

   func parseHeader(raw string) result.Result[HeaderMeta]
   ```
2. **Parameter Structs for >2 Arguments:** If an extracted function requires more than 2 parameters, combine them into an options or params struct (`*Params` / `*Options`).

---

## 5. Boolean Principles During Extraction (TOTAL BAN on Negatives)

When writing new helper functions or guard clauses during extraction:

- **Implicit Positive Checks Only:** NEVER write `if isReady == true`. ALWAYS write `if isReady`.
- **No Negative Polarity:** NEVER combine a positive and negative condition (`if isA && !isB`). Split into separate guard clauses!
- **Positive Naming Only:** All booleans MUST use affirmative prefixes: `is*` or `has*` only.
  - ❌ `!isSuccess` is BANNED -> Use `isFail` or `if isFail { ... }`
  - ❌ `isNotValid` is BANNED -> Use `isInvalid`
  - ❌ `isDisabled` -> Use `isPaused` or `isInactive`
- **Zero Nested If Statements:** Every extracted function MUST have a nesting depth of <= 1. Flatten all conditional branches using early guard returns!

---

## 6. Execution & Build Policy (Strict Ban on Intermediate Runs)

To maintain maximum focus and velocity during micro-tasking:

1. **NO INTERMEDIATE TEST RUNNING (TOTAL BAN):**
   - NEVER run unit test suites (`go test ./...`, `npm test`, `pytest`, runner scripts) during routine file refactoring turns.
   - All unit test verification is strictly deferred to dedicated QA and CI/CD pipelines.
2. **NO INTERMEDIATE BUILD CHECKING (TOTAL BAN):**
   - DO NOT execute build commands (`go build`, `npm run build`, compiler invocations) after editing individual files.
3. **FINAL STEP BUILD VERIFICATION (END OF STEPS ONLY):**
   - Syntax and compilation verification is performed strictly at the **final step** of the refactoring run, after all file extractions and import adjustments are complete, ensuring any dangling references or unexported symbols are resolved before concluding.

---

## 7. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO INTERMEDIATE BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) after individual file edits. Build verification is checked ONLY at the final step.
- [ ] **NO LINE COMPRESSION (TOTAL BAN):** NEVER remove blank lines, merge statements, or compress `if/else` blocks to reduce line count.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.

---

## 8. Non-Negotiable Coding Guidelines Checklist

- [ ] Canonical Size Limits: Every refactored file is under 100 lines (recommended <= 80 lines), excluding allowed exceptions (JSON, types, maps, vars, consts).
- [ ] Function Sizing: Functions <= 8 lines preferred (hard cap 15 lines, absolute waiver cap 25 lines).
- [ ] Whitespace & Style Preserved: Full vertical breathing room maintained; blank line before return, blank line after closing brace `}`.
- [ ] Wrapper Objects Used: Complex multi-value returns or extracted parameters encapsulated in clean structs/types.
- [ ] Boolean Conventions: All booleans prefixed with `is` or `has`. Zero negative checks (`!isSuccess` banned; use `isFail`). Implicit positive checks only.
- [ ] Control Flow Flattened: Zero nested `if` statements (nesting depth <= 1). Early guard returns used throughout.
- [ ] Line Endings & Encoding: Strictly Unix LF (`\n`) and UTF-8 without BOM.
- [ ] Strict Relative Git Paths: All citations and links use relative paths from repository root.

---

## 9. Task Consolidation & Final Step Git Commit & Push Mandate

### Task Consolidation & File Reduction (End of Loop)
When all subtasks for the parent task (`.ai-memory/plans/pending/xx-size-reduction.md`) are finished:
1. Combine all completed granular subtasks from `.ai-memory/plans/subtasks/xx-size-reduction/*.md` into `.ai-memory/plans/completed/xx-size-reduction.md`.
2. Include a header explicitly documenting the initial file sizes, final file sizes, and loop step metrics.
3. Delete the original granular `.md` files in `.ai-memory/plans/subtasks/xx-size-reduction/`.
4. Delete the original parent plan `.ai-memory/plans/pending/xx-size-reduction.md`.
5. Update `.ai-memory/plans/01-index.md` to point to the newly consolidated completed file.

### Final Step Git Commit & Push Mandate (Strict Checklist)
- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!
