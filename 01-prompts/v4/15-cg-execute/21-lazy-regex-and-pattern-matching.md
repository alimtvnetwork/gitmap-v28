[/goal](slashCommand:goal) Autonomously scan, plan, refactor, and verify all regular expression usage across the codebase, eliminating raw `regexp.MustCompile` and `regexp.Compile` outside package `lazyregex`, replacing inline compilation and blind `re.MatchString` test assertions with thread-safe `lazyregex.New(...)` and wrapped `lazyregex.MatchResult(...)` returning `*MatchResult` (`ResultGroup`) and structured `*appfault.AppError` diagnostics, enforcing affirmative evaluation (`rs.IsMatch()`, `rs.IsSuccess()`, `rs.IsFailed()`), and providing rich failure diagnostics showing pattern, compared text, character count, and failure cause until 100% green without stopping with strict no-build and no-test execution (NEVER run build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine execution turns; all compilation and testing are strictly verified later in CI/CD). Spawn autonomous subagents (A = 2, H = 2) for parallel reading and modular spec generation, use GitMap high-speed commands as primary, establish a single-agent blueprint in Phase 1 (first 50% steps budget), and execute disjoint refactoring in Phase 2 (remaining 50% steps budget) with continuous self-looping until 100% complete and finalized with an atomic push.

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
2. [ ] [/goal](slashCommand:goal) Phase 1 (Step B): Write the master audit specification in `.ai-memory/plans/pending/XX-lazy-regex-audit.md` with an exhaustive Violation Ledger.
3. [ ] [/goal](slashCommand:goal) Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.ai-memory/plans/subtasks/XX-lazy-regex/`.
4. [ ] [/goal](slashCommand:goal) Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/readme.md`.
5. [ ] [/goal](slashCommand:goal) Phase 2 (Step A): Open each target file and refactor raw `regexp.MustCompile` / `regexp.Compile` to centralized `lazyregex.New(...)`.
6. [ ] [/goal](slashCommand:goal) Phase 2 (Step B): Modernize all test assertions from blind `if !re.MatchString(s) { t.Error(...) }` to `rs := lazyRegex.MatchResult(s); if rs.IsFailed() { t.Error(rs.AppError()) }`.
7. [ ] [/goal](slashCommand:goal) Phase 2 (Step C): Refactor group and submatch extraction to fluent methods: `rs.Items()`, `rs.Map()`, `rs.First()`, `rs.Last()`, `rs.FirstOrDefault()`, `rs.At()`.
8. [ ] [/goal](slashCommand:goal) Phase 2 (Step D): Enforce <= 8–15 line function decomposition and clean blank-line spacing.
9. [ ] [/goal](slashCommand:goal) Phase 2 (Step E): Execute targeted package tests and linters (`golangci-lint run -c .golangci.yml`) to verify 0 remaining violations. DO NOT run full CI runner during routine turns.
10. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/memory/readme.md` for project memory index and past learnings.
11. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
12. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
13. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for hallucination prevention and micro-tasking.
14. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/01-cross-language/17-regex-usage-guidelines.md` for lazy regex and pattern caching rules.
15. [ ] [/learn](slashCommand:learn) Ingest `02-spec/03-error-manage/` for error handling architectures and AppError.
16. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.
17. [ ] [/goal](slashCommand:goal) Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Raw Regex & Blind Test Matches, Spec in .ai-memory/plans/pending/, Subtasks, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Actively Edit Code, Refactor to lazyregex.MatchResult, Modernize Tests, Verify Local Linters)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

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

---

## Dedicated Section: Lazy Regex Architecture & Rich Test Diagnostics

Regular expressions are expensive to compile and prone to backtracking. Uncached compilation wastes CPU cycles, while opaque boolean test assertions (`if !re.MatchString(output)`) hide failure context from developers and autonomous AI agents.

### Mandatory Rules (Non-Negotiable)

1. **Total Ban on Raw `regexp.MustCompile` & `regexp.Compile` Across Domain Code:**
   - NEVER call `regexp.MustCompile()` or `regexp.Compile()` inside package functions, methods, or loops.
   - ALWAYS use `lazyregex.New(pattern)` or `lazyregex.NewLock(pattern)` from `cli/lazyregex`.
   - Patterns are compiled on demand (lazy) at most once, and cached in a global deduplicated thread-safe registry.

2. **Mandatory `MatchResult` Envelope for Test Matching:**
   - In unit tests and integration tests comparing output or parsing tokens, NEVER write opaque boolean checks:
     ```go
     // ❌ BANNED: Blind regex test assertion with no diagnostic context
     if !re.MatchString(content) {
         t.Error("expected output to match pattern")
     }
     ```
   - ALWAYS execute matching via `lazyregex.MatchResult(content)`:
     ```go
     // ✅ REQUIRED: Wrapped Result with affirmative evaluation and rich diagnostics
     var rsLazyRegex = lazyRegex.MatchResult(comparing)
     if rsLazyRegex.IsFailed() {
         t.Error(rsLazyRegex.AppError())
     }
     ```

3. **Rich Diagnostic Error Contract:**
   - On matching failure, `rsLazyRegex.AppError()` generates a fully contextualized diagnostic report containing:
     1. The regular expression pattern that failed to match.
     2. The actual content being compared (with safe preview truncation for very large buffers).
     3. The exact length of the compared string in bytes.
     4. Context fields: `pattern`, `comparing_preview`, `comparing_len`.
     5. Structured AppError code `[E1000:VALIDATION]` with file and line origin.

4. **Fluent MatchGroup (`ResultGroup`) Accessors:**
   - When a match succeeds, extract submatches and named capture groups using fluent helper methods:
     - `rs.Group()`: Accesses the underlying `*MatchGroup` (alias `ResultGroup`).
     - `rs.Items()`: Returns `[]string` of all captured submatches (`[0]` = full match, `[1..N]` = groups).
     - `rs.Map()`: Returns `GroupMap` containing named capture groups (`(?P<name>...)`).
     - `rs.First()`: Returns the full match string or empty string.
     - `rs.Last()`: Returns the last captured group.
     - `rs.FirstOrDefault(defaultValue)`: Returns the first match or fallback default if empty.
     - `rs.At(index)`: Safely accesses submatch at index without slice out-of-bounds panic.
     - `rs.Count()` / `rs.Len()`: Returns the number of captured submatches.

5. **Multi-Match Evaluation (`MatchAllResults`):**
   - For extracting multiple non-overlapping occurrences across text:
     ```go
     all := lazyRegex.MatchAllResults(content)
     if all.IsFailed() {
         t.Error(all.AppError())
     }
     for _, group := range all.Groups() {
         role := group.GetNamed("role")
         ip := group.GetNamed("ip")
     }
     ```

---

## Code Pattern Comparison

### Pattern A: Test Pattern Matching

```go
// -----------------------------------------------------------------------------
// ❌ ANTI-PATTERN: Opaque regex boolean check hiding pattern and compared text
// -----------------------------------------------------------------------------
func TestLegacySubcommandMatching(t *testing.T) {
    re := regexp.MustCompile(`^(add|join|nodes)$`)
    output := getCommandOutput()

    if !re.MatchString(output) {
        t.Error("expected valid subcommand output") // Tells AI nothing about what failed!
    }
}

// -----------------------------------------------------------------------------
// ✅ REQUIRED: Wrapped Result with rich diagnostic error reporting
// -----------------------------------------------------------------------------
func TestModernSubcommandMatching(t *testing.T) {
    re := lazyregex.New(`^(?P<action>add|join|nodes)$`)
    output := getCommandOutput()

    var rsLazyRegex = re.MatchResult(output)
    if rsLazyRegex.IsFailed() {
        t.Error(rsLazyRegex.AppError())
    }

    if rsLazyRegex.Map().Get("action") != "nodes" {
        t.Errorf("expected action 'nodes', got: %q", rsLazyRegex.Map().Get("action"))
    }
}
```

### Pattern B: What the Diagnostic Error Looks Like on Failure

When a test using `t.Error(rsLazyRegex.AppError())` fails, it outputs:

```text
[E1000:VALIDATION] validation: regex pattern does not match content
  Pattern:   "^(?P<action>add|join|nodes)$"
  Comparing: "gitmap cluster unknown-command"
  Length:    29 bytes (at=lazyregex/match_error.go:32) (ctx=map[comparing_len:29 comparing_preview:gitmap cluster unknown-command pattern:^(?P<action>add|join|nodes)$])
```

Autonomous AI agents and engineers can immediately diagnose:
1. **What was expected**: `^(?P<action>add|join|nodes)$`
2. **What was received**: `"gitmap cluster unknown-command"`
3. **Exact failure point**: File and line number in `match_error.go`
4. **Error classification**: `[E1000:VALIDATION]`

--

## 🚨 Highest Priority Instructions (Appended User Tasks & Instructions Below)

[PASTE USER REQUEST / TASK INSTRUCTIONS HERE — THE AGENT MUST EXECUTE WHATEVER IS WRITTEN BELOW WITH ABSOLUTE PRIORITY AND PRECEDENCE OVER ALL GENERAL GUIDELINES ABOVE]
