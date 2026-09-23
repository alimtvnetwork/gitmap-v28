# String Operations & Memory Efficiency — Coding Guideline Execution (must follow)

Trigger Keywords & Aliases: `cg-string-efficiency`, `cg-strings`, `cg-string-comparison`, `cg-equalfold`, `cg-execute strings`, `audit string efficiency`, `optimize string operations`, `string allocation reduction`

> **Prompt Version:** 2.2.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform (configurable per run).

/goal Autonomously scan, audit, plan, and refactor inefficient string operations, comparisons, case conversions, and loop allocations across the codebase. Enforce zero-allocation case-folding (`strings.EqualFold` in Go, `StringComparison.OrdinalIgnoreCase` in C#, `eq_ignore_ascii_case` in Rust, `strcasecmp` in PHP), short-circuiting lazy evaluation, loop hoisting, and string builder patterns without altering underlying business logic. Preserve semantic behavior as the absolute first priority, keep functions <= 8–15 lines, enforce boolean conventions, and defer build verification strictly to the final step without running intermediate tests or builds.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

You MUST execute this task via a strict 3-Phase pipeline governed by the N-step budget. Do not skip steps.

```text
N = 200  (Total self-loop steps budget, read-only after initialization)
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Codebase, Business Logic Analysis, Map Violations in .ai-memory/plans/pending/)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Subtasks, Logic-Safe String Optimizations, Guard Short-Circuiting)
```

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A - Discovery & Inventory): Deeply scan the target codebase using the fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py` with `--limit`) to inventory all string comparison, case conversion, and allocation anti-patterns without truncation.
2. [ ] /goal Phase 1 (Step B - Business Logic Intent Analysis): For every discovered violation, inspect the surrounding context to understand the exact business logic intent: Was exact case-insensitive equality intended? Was substring containment intended? Was prefix/suffix matching intended?
3. [ ] /goal Phase 1 (Step C - Master Plan Generation): Write the master architectural specification into `.ai-memory/plans/pending/xx-string-efficiency.md` with an exhaustive Violation Ledger table (File, Line, Current Pattern, Logic Intent, Safe Optimized Pattern, Status).
4. [ ] /goal Phase 1 (Step D - Subtask Decomposition): Decompose the master plan into lean, bounded subtask files in `.ai-memory/plans/subtasks/xx-string-efficiency/01-<subtask>.md`, `02-<subtask>.md`, etc.
5. [ ] /goal Phase 1 (Step E - Mandatory Auto-Loop): As soon as Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for permission**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.
6. [ ] /goal Phase 2 (Step A - Logic-Safe Refactoring): Execute each subtask, applying the safe optimized string pattern. Under NO circumstances change substring matching to equality checking or alter filtering behavior.
7. [ ] /goal Phase 2 (Step B - Lazy Evaluation & Short-Circuiting): Replace eager boolean assignments with early-returning guard clauses to avoid unnecessary secondary case conversions and string evaluations.
8. [ ] /goal Phase 2 (Step C - Loop Hoisting & Zero-Allocation Builders): Hoist invariant case conversions (e.g. lowering search terms or filter strings) outside iterations. Replace repeated `+` string concatenations in loops with dedicated builders (`strings.Builder`, `StringBuilder`, array join).
9. [ ] /goal Phase 2 (Step D - Function & File Size Compliance): Ensure all refactored functions remain <= 8 lines of body logic (hard cap of <= 15 lines) and files remain under 100 lines.
10. [ ] /goal Phase 2 (Step E - Boolean & Style Conventions): Enforce affirmative boolean naming (`is*`, `has*`), zero explicit `== true`, zero negative polarity in conditionals, and flatten nested `if` statements to depth <= 1 using guard clauses.
11. [ ] /goal Phase 2 (Step F - Banned Intermediate Verification): DO NOT run unit tests (`go test`, `pytest`, `npm test`) and DO NOT verify builds during intermediate micro-refactoring steps.
12. [ ] /goal Phase 2 (Step G - Final Step Build Verification): At the conclusion of all refactoring subtasks, run targeted syntax/build checks to resolve any compilation errors or import issues across all modified files.
13. [ ] /goal Phase 3 (Step A - Task Consolidation): Consolidate all completed subtasks into `.ai-memory/plans/completed/xx-string-efficiency.md`, delete granular subtask files, and update `.ai-memory/plans/01-index.md`.
14. [ ] /goal Phase 3 (Step B - Final Step Git Commit & Push): Stage all modified files, consolidated plans, and memory records (`git add -A`), commit them in a single clean grouped atomic commit, and push to git. Never commit per-file.
15. [ ] /learn Ingest `.ai-memory/memory/01-index.md` for project memory index and past learnings.
16. [ ] /learn Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
17. [ ] /learn Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical size tiers.
18. [ ] /learn Ingest `02-spec/02-coding-guidelines/03-golang/06-string-slice-internals.md` for string allocation internals.
19. [ ] /learn Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.

---

### Fast File Discovery & Reading via Python Toolchain (Mandatory Acceleration)

To avoid 50-result tool truncation limits and eliminate exploratory roundtrips, the AI agent MUST use the repository's dedicated Python discovery scripts first:

1. **Inventory String Operation Anti-Patterns (with `--limit` option):**
   ```bash
   # Find ToLower / ToUpper case conversions in Go
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "strings.ToLower" --lang go --limit 100
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "strings.ToUpper" --lang go --limit 100

   # Find toLowerCase / toUpperCase in TypeScript/JavaScript
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "toLowerCase()" --lang ts --limit 100

   # Find case conversions in Python / C# / PHP / Rust
   python 03-ai-scripts/12-fast-cached-grep.py --pattern ".lower()" --lang py --limit 50
   python 03-ai-scripts/12-fast-cached-grep.py --pattern ".ToLower()" --limit 50
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "strtolower(" --limit 50
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "to_lowercase()" --lang rs --limit 50
   ```

2. **Inspect Target Files Efficiently:**
   ```bash
   python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000
   python 03-ai-scripts/17-fast-file-reader.py --search-pattern "ToLower" --path <folder> --limit 50
   ```

3. **Verify Function & File Sizing:**
   ```bash
   python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 50
   python 03-ai-scripts/13-file-size-guard.py
   ```

Do not rely on standard search tools with 50-item truncation when discovering repository-wide string anti-patterns.

---

## 1. Prime Directive: Business Logic Safety First (Zero-Regression Mandate)

> [!IMPORTANT]
> **UNDERSTAND THE LOGIC FIRST — NEVER BREAK BUSINESS SEMANTICS:**
> Code efficiency optimizations MUST NEVER alter program behavior. Before modifying any string check, the AI agent MUST analyze whether the original code required:
> 1. **Exact Case-Insensitive Equality** (`full match`)
> 2. **Case-Insensitive Substring Containment** (`partial match / search filter`)
> 3. **Case-Insensitive Prefix or Suffix Matching** (`starts with / ends with`)
>
> Changing a substring search (`Contains`) to an equality check (`EqualFold`) is a CATASTROPHIC BUG that breaks filters and searches!

### Semantic Intent Mapping Matrix

| Original Expression | Intended Semantic | Incorrect (Banned) Fix | Correct Logic-Safe Optimization |
| :--- | :--- | :--- | :--- |
| `strings.ToLower(a) == strings.ToLower(b)` | Exact Equality | `strings.Contains(...)` | `strings.EqualFold(a, b)` (Zero allocations) |
| `strings.ToLower(a) == "admin"` | Exact Equality | `strings.Contains(...)` | `strings.EqualFold(a, "admin")` (Zero allocations) |
| `strings.Contains(strings.ToLower(name), term)` | Substring Match | ❌ `strings.EqualFold(name, term)` (Breaks substring search!) | ✅ Hoist `term = strings.ToLower(term)` outside caller loop; short-circuit evaluations |
| `strings.HasPrefix(strings.ToLower(s), "test")` | Prefix Match | `strings.EqualFold(s, "test")` | `strings.HasPrefix(strings.ToLower(s), "test")` with hoisted prefix |
| `s[:4] == "test"` | Prefix Match | Slicing without length check (panics) | `strings.HasPrefix(s, "test")` |

---

## 2. Cross-Language Efficiency Architectures

### A. Go String Optimization Rules

1. **Zero-Allocation Equality with `strings.EqualFold`:**
   Go's `strings.EqualFold(s, t)` compares two UTF-8 strings under Unicode case-folding without heap-allocating new strings:
   ```go
   // ❌ BAD (Allocates 2 new strings on the heap):
   if strings.ToLower(userRole) == strings.ToLower(targetRole) { ... }
   if strings.ToLower(command) == "status" { ... }

   // ✅ GOOD (0 heap allocations, ~10x faster):
   if strings.EqualFold(userRole, targetRole) { ... }
   if strings.EqualFold(command, "status") { ... }
   ```

2. **Lazy Evaluation & Short-Circuiting in Multi-Field Filters:**
   When matching multiple fields against a filter, NEVER compute all lowerings eagerly:
   ```go
   // ❌ BAD (Eager allocation: always lowers both name AND path even if name matches):
   func matchesFilter(name, path, filter string) bool {
       term := strings.ToLower(filter)
       matchesName := strings.Contains(strings.ToLower(name), term)
       matchesPath := strings.Contains(strings.ToLower(path), term)

       return matchesName || matchesPath
   }

   // ✅ GOOD (Short-circuited: path is NEVER lowered if name matches; <= 8 lines):
   func matchesFilter(name, path, lowerFilter string) bool {
       if strings.Contains(strings.ToLower(name), lowerFilter) {
           return true
       }

       return strings.Contains(strings.ToLower(path), lowerFilter)
   }
   ```

3. **Loop Hoisting of Case Conversions:**
   Never repeatedly convert the same search term or pattern inside a loop:
   ```go
   // ❌ BAD (Lowers filter string N times inside the loop):
   for _, item := range items {
       if matchesFilter(item.Name, item.Path, filter) {
           results = append(results, item)
       }
   }

   // ✅ GOOD (Lowers filter string once before the loop):
   lowerFilter := strings.ToLower(filter)
   for _, item := range items {
       if matchesFilter(item.Name, item.Path, lowerFilter) {
           results = append(results, item)
       }
   }
   ```

4. **String Concatenation in Loops (`strings.Builder`):**
   ```go
   // ❌ BAD (O(N^2) allocations due to string immutability):
   var result string
   for _, chunk := range chunks {
       result += chunk
   }

   // ✅ GOOD (O(N) with single pre-allocated or dynamic buffer):
   var builder strings.Builder
   for _, chunk := range chunks {
       builder.WriteString(chunk)
   }
   result := builder.String()
   ```

5. **Safe Prefix & Suffix Checking:**
   Avoid slicing syntax `s[:len(prefix)]` which causes out-of-bounds panics if `len(s) < len(prefix)`. Always use `strings.HasPrefix(s, prefix)` and `strings.HasSuffix(s, suffix)`.

---

### B. TypeScript & JavaScript Optimization Rules

1. **Avoid Double Case Conversion When Comparing to Constants:**
   ```ts
   // ❌ BAD (Redundant lowercasing of literal):
   if (status.toLowerCase() === "ACTIVE".toLowerCase()) { ... }

   // ✅ GOOD (Literal is already lowercase; 1 allocation instead of 2):
   if (status.toLowerCase() === "active") { ... }
   ```

2. **Short-Circuiting Filter Predicates:**
   ```ts
   // ❌ BAD (Evaluates both sides eagerly):
   const matches = name.toLowerCase().includes(term) || path.toLowerCase().includes(term);

   // ✅ GOOD (Short-circuits before lowering path):
   const isMatch = (name: string, path: string, lowerTerm: string): boolean => {
     if (name.toLowerCase().includes(lowerTerm)) {
       return true;
     }
     return path.toLowerCase().includes(lowerTerm);
   };
   ```

3. **Loop Aggregation via Arrays:**
   Avoid repeated `str += val` in large loops. Use an array accumulator `parts.push(val)` and finalize with `parts.join('')`.

---

### C. Python Optimization Rules

1. **Proper Unicode Case-Folding (`casefold`):**
   ```python
   # ❌ SUBOPTIMAL: lower() misses certain Unicode foldings (e.g. German ß -> ss)
   if a.lower() == b.lower(): ...

   # ✅ RECOMMENDED for robust case-insensitive comparison:
   if a.casefold() == b.casefold(): ...
   ```

2. **Loop Hoisting:**
   ```python
   # ❌ BAD: lowering search term in comprehension
   matched = [item for item in items if term.lower() in item.name.lower()]

   # ✅ GOOD: lower once
   term_lower = term.casefold()
   matched = [item for item in items if term_lower in item.name.casefold()]
   ```

3. **Concatenation in Loops:**
   Always use `''.join(chunks)` instead of cumulative string `+` in loops.

---

### D. Rust, C#, and PHP Optimization Rules

- **Rust:**
  - Equality: Replace `a.to_lowercase() == b.to_lowercase()` (allocates 2 heap `String`s) with `a.eq_ignore_ascii_case(&b)` (0 allocations for ASCII/UTF-8).
  - Pre-allocate strings via `String::with_capacity(expected_len)`.
- **C#:**
  - Equality: Replace `a.ToLower() == b.ToLower()` with `string.Equals(a, b, StringComparison.OrdinalIgnoreCase)`.
  - Substring: Use `a.Contains(b, StringComparison.OrdinalIgnoreCase)`.
  - Builders: Use `StringBuilder` with initial capacity for loops.
- **PHP:**
  - Equality: Replace `strtolower($a) === strtolower($b)` with `strcasecmp($a, $b) === 0`.
  - Substring: Use `stripos($haystack, $needle) !== false` instead of `strpos(strtolower(...), strtolower(...))`.

---

## 3. Function & File Size Rules During String Refactoring

- **Function Size Cap:** Target <= 8 lines of body logic; hard maximum of <= 15 lines.
- **File Size Cap:** Target <= 80 lines; hard maximum of <= 100 lines (excluding allowed JSON, types, vars, consts, maps exceptions).
- **Control Flow Flattening:** Nesting depth MUST remain <= 1. Use early guard returns.
- **Zero Line Compression:** Strictly preserve vertical whitespace, blank line before returns, and blank line after closing curly braces `}`.

---

## 4. Execution & Build Policy (Strict Ban on Intermediate Runs)

To maintain high throughput and avoid distracting CI noise during micro-refactoring:

1. **NO INTERMEDIATE TEST RUNNING (TOTAL BAN):**
   - NEVER run unit test suites (`go test ./...`, `npm test`, `pytest`, runner scripts) during routine file editing turns.
   - All unit test verification is strictly deferred to dedicated QA and CI/CD pipelines.
2. **NO INTERMEDIATE BUILD CHECKING (TOTAL BAN):**
   - DO NOT execute build commands (`go build`, `npm run build`, compiler invocations) after editing individual files.
3. **FINAL STEP BUILD VERIFICATION (END OF STEPS ONLY):**
   - Syntax and compilation verification is performed strictly at the **final step** of the refactoring run, after all file modifications are complete, ensuring any dangling references or import errors are cleanly resolved before concluding.

---

## 5. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO INTERMEDIATE BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) after individual file edits. Build verification is checked ONLY at the final step.
- [ ] **NO SEMANTIC LOGIC ALTERATION (TOTAL BAN):** NEVER change substring search (`Contains`) to equality comparison (`EqualFold`), or alter filtering semantics. Logic preservation is the #1 priority.
- [ ] **NO LINE COMPRESSION (TOTAL BAN):** NEVER remove blank lines, merge statements, or compress `if/else` blocks to reduce line count.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.

---

## 6. Non-Negotiable Coding Guidelines Checklist

- [ ] Business Logic Preserved: Exact semantic behavior verified; substring checks remain substring checks, equality checks remain equality checks.
- [ ] Zero-Allocation Equality: Used `strings.EqualFold` (Go), `string.Equals(..., OrdinalIgnoreCase)` (C#), `eq_ignore_ascii_case` (Rust), or `strcasecmp` (PHP) for case-insensitive comparisons.
- [ ] Lazy Evaluation & Short-Circuiting: Replaced eager multi-check assignments with early-returning guard clauses.
- [ ] Loop Hoisting: Invariant case conversions and search terms hoisted outside loops.
- [ ] Function Sizing: Functions <= 8 lines preferred (hard cap 15 lines).
- [ ] File Sizing: Files <= 80 lines preferred (hard cap 100 lines).
- [ ] Whitespace & Style Preserved: Blank line before return, blank line after closing brace `}`.
- [ ] Boolean Conventions: All booleans prefixed with `is` or `has`. Zero negative checks (`!isSuccess` banned; use `isFail`). Implicit positive checks only.
- [ ] Control Flow Flattened: Zero nested `if` statements (nesting depth <= 1). Early guard returns used throughout.
- [ ] Line Endings & Encoding: Strictly Unix LF (`\n`) and UTF-8 without BOM.
- [ ] Strict Relative Git Paths: All citations and links use relative paths from repository root.

---

## 7. Task Consolidation & Final Step Git Commit & Push Mandate

### Task Consolidation & File Reduction (End of Loop)
When all subtasks for the parent task (`.ai-memory/plans/pending/xx-string-efficiency.md`) are finished:
1. Combine all completed granular subtasks from `.ai-memory/plans/subtasks/xx-string-efficiency/*.md` into `.ai-memory/plans/completed/xx-string-efficiency.md`.
2. Include a header explicitly documenting initial and optimized patterns, allocation reductions, and loop step metrics.
3. Delete the original granular `.md` files in `.ai-memory/plans/subtasks/xx-string-efficiency/`.
4. Delete the original parent plan `.ai-memory/plans/pending/xx-string-efficiency.md`.
5. Update `.ai-memory/plans/01-index.md` to point to the newly consolidated completed file.

### Final Step Git Commit & Push Mandate (Strict Checklist)
- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

---

## Metadata

- slug: cg-string-efficiency
- version: 2.2.0
- category: performance-and-memory
- priority: high
- autoloop: true
- languages: [golang, typescript, python, rust, csharp, php]
