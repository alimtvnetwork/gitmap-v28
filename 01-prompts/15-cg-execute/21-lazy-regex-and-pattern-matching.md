# Lazy Regex, Centralized Pattern Caching & Match Result Diagnostics — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-lazyregex`, `cg-regex`, `cg-execute regex`, `audit regex`, `fix raw regex`, `enforce lazy regex`, `regex-match-result`, `lazyregex-match`, `ban-regexp-mustcompile`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, plan, refactor, and verify all regular expression usage across the codebase, eliminating raw `regexp.MustCompile` and `regexp.Compile` outside package `lazyregex`, replacing inline compilation and blind `re.MatchString` test assertions with thread-safe `lazyregex.New(...)` and wrapped `lazyregex.MatchResult(...)` returning `*MatchResult` (`ResultGroup`) and structured `*apperror.AppError` diagnostics, enforcing affirmative evaluation (`rs.IsMatch()`, `rs.IsSuccess()`, `rs.IsFailed()`), and providing rich failure diagnostics showing pattern, compared text, character count, and failure cause until 100% green without stopping.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase using the fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py` with `--limit`) to inventory all architectural violations and anti-patterns without truncation.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.lovable/plans/pending/XX-lazy-regex-audit.md` with an exhaustive Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.lovable/plans/subtasks/XX-lazy-regex/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/01-index.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and refactor raw `regexp.MustCompile` / `regexp.Compile` to centralized `lazyregex.New(...)`.
6. [ ] /goal Phase 2 (Step B): Modernize all test assertions from blind `if !re.MatchString(s) { t.Error(...) }` to `rs := lazyRegex.MatchResult(s); if rs.IsFailed() { t.Error(rs.AppError()) }`.
7. [ ] /goal Phase 2 (Step C): Refactor group and submatch extraction to fluent methods: `rs.Items()`, `rs.Map()`, `rs.First()`, `rs.Last()`, `rs.FirstOrDefault()`, `rs.At()`.
8. [ ] /goal Phase 2 (Step D): Enforce <= 8–15 line function decomposition and clean blank-line spacing.
9. [ ] /goal Phase 2 (Step E): Execute targeted package tests and linters (`golangci-lint run -c .golangci.yml`) to verify 0 remaining violations. DO NOT run full CI runner during routine turns.
10. [ ] /learn Ingest `.lovable/memory/01-index.md` for project memory index and past learnings.
11. [ ] /learn Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
12. [ ] /learn Ingest `spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
13. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for hallucination prevention and micro-tasking.
14. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/17-regex-usage-guidelines.md` for lazy regex and pattern caching rules.
15. [ ] /learn Ingest `spec/03-error-manage/` for error handling architectures and AppError.
16. [ ] /learn Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.
17. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Raw Regex & Blind Test Matches, Spec in .lovable/plans/pending/, Subtasks, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Actively Edit Code, Refactor to lazyregex.MatchResult, Modernize Tests, Verify Local Linters)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

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
