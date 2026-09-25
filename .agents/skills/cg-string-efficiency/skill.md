---
name: cg-string-efficiency
description: Autonomously audits, optimizes, and refactors string operations and comparisons across polyglot codebases for zero-allocation case folding, short-circuiting lazy evaluation, and loop hoisting while strictly preserving business logic.
---

# Skill: String Operations & Memory Efficiency (`cg-string-efficiency`)

This skill governs autonomous auditing, optimization, and refactoring of inefficient string operations, comparisons, case conversions, and loop allocations across polyglot codebases (Go, TypeScript, Python, PHP, Rust, C#).

---

## 1. Mandatory Architectural Invariants

### 1. Prime Directive: Business Logic Safety First (Zero-Regression Mandate)
Code efficiency optimizations MUST NEVER alter program behavior. Before modifying any string check, the agent MUST analyze whether the original code required:
- **Exact Case-Insensitive Equality:** Use zero-allocation folding (`strings.EqualFold` in Go).
- **Case-Insensitive Substring Containment:** DO NOT change to `EqualFold`! Keep substring search semantics (`Contains`), but optimize via lazy short-circuiting and loop hoisting.
- **Prefix / Suffix Matching:** Use `HasPrefix` / `HasSuffix` or `startsWith` / `endsWith`.

### 2. Zero-Allocation Case-Folding

- **Go:** Replace `strings.ToLower(a) == strings.ToLower(b)` or `strings.ToLower(s) == "val"` with `strings.EqualFold(a, b)` (0 allocations, ~10x faster).
- **TypeScript/JS:** Avoid lowering both sides when comparing with a constant literal (`status.toLowerCase() === "active"`).
- **Python:** Use `a.casefold() == b.casefold()` for full Unicode case-folding.
- **Rust:** Replace `a.to_lowercase() == b.to_lowercase()` with `a.eq_ignore_ascii_case(&b)` (0 heap allocations).
- **C#:** Replace `a.ToLower() == b.ToLower()` with `string.Equals(a, b, StringComparison.OrdinalIgnoreCase)`.
- **PHP:** Replace `strtolower($a) === strtolower($b)` with `strcasecmp($a, $b) === 0`.

### 3. Short-Circuiting Lazy Evaluation
In multi-field matching predicates, never compute all lowerings or searches eagerly:
- Return immediately on the first true branch using early guard returns.
- Avoid executing secondary case conversions when the first condition already matches.

### 4. Loop Hoisting of Case Conversions

- Never repeatedly convert the same search pattern or filter string inside an iteration.
- Lower the pattern once prior to the loop and pass the pre-lowered term to item matchers.

### 5. String Builders for Cumulative Concatenation

- In Go, replace `s += chunk` inside loops with `strings.Builder`.
- In C#, use `StringBuilder`. In Python, use `''.join(parts)`. In TS, use array accumulator `parts.join('')`.

### 6. Function Sizing & File Hygiene

- Functions target **<= 8 lines** of body logic (hard cap: 15 lines).
- Files target **<= 80 lines** (hard cap: 100 lines, excluding allowed exceptions).
- Whitespace preservation: exactly ONE blank line before return, exactly ONE blank line after closing brace `}`.

### 7. Execution & Build Policy

- **NO INTERMEDIATE TEST RUNNING:** NEVER run unit test suites (`go test ./...`, `npm test`, `pytest`) during routine refactoring turns.
- **NO INTERMEDIATE BUILD CHECKING:** DO NOT execute build commands (`go build`, `npm run build`) after individual file edits.
- **FINAL STEP BUILD VERIFICATION ONLY:** Verify compilation strictly at the final step after all file extractions and import adjustments are completed.

---

## 2. Code Patterns: Before vs After

### Pattern 1: Case-Insensitive Equality (Go)

```go
// ❌ ANTI-PATTERN: Heap-allocates 2 new strings on every comparison
if strings.ToLower(role) == strings.ToLower(targetRole) {
    return true
}

// ✅ REFACTORED: Zero heap allocations, Unicode case-folding
if strings.EqualFold(role, targetRole) {
    return true
}
```

### Pattern 2: Multi-Field Filter Matching with Short-Circuiting (Go)

```go
// ❌ ANTI-PATTERN: Always allocates lowerings for BOTH name and path, even if name matches
func matchesFilter(name, path, filter string) bool {
    term := strings.ToLower(filter)
    matchesName := strings.Contains(strings.ToLower(name), term)
    matchesPath := strings.Contains(strings.ToLower(path), term)

    return matchesName || matchesPath
}

// ✅ REFACTORED: Lazy short-circuit; path is never lowered if name matches; <= 8 lines
func matchesFilter(name, path, lowerFilter string) bool {
    if strings.Contains(strings.ToLower(name), lowerFilter) {
        return true
    }

    return strings.Contains(strings.ToLower(path), lowerFilter)
}
```

### Pattern 3: Loop Hoisting of Invariant Filter Lowering (Go)

```go
// ❌ ANTI-PATTERN: filter is converted to lowercase on every iteration
for _, item := range items {
    if matchesFilter(item.Name, item.Path, filter) {
        results = append(results, item)
    }
}

// ✅ REFACTORED: lowered once outside the loop
lowerFilter := strings.ToLower(filter)
for _, item := range items {
    if matchesFilter(item.Name, item.Path, lowerFilter) {
        results = append(results, item)
    }
}
```

---

## 3. Fast File Discovery & Reading via Python Toolchain

To avoid 50-result tool truncation limits, use the repository's dedicated Python scripts:

1. **Grep for Case Conversion Violations:**
   ```bash
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "strings.ToLower" --lang go --limit 100
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "toLowerCase()" --lang ts --limit 100
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "ToLower()" --limit 50
   ```
2. **Read File Content Efficiently:**
   ```bash
   python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000
   ```
3. **Verify Function & File Sizing:**
   ```bash
   gitmap find "*.go" -ext "go"
   python 03-ai-scripts/13-file-size-guard.py
   # Fallback: python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 50
   ```

---

## 4. Banned Operations Checklist

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts, Go (`go test ./...`), or any test runner during routine execution turns.
- [ ] **NO INTERMEDIATE BUILD CHECKING (TOTAL BAN):** NEVER run build commands after individual file edits. Build verification is checked ONLY at the final step.
- [ ] **NO SEMANTIC LOGIC ALTERATION (TOTAL BAN):** NEVER change substring search (`Contains`) to equality comparison (`EqualFold`), or alter filtering semantics.
- [ ] **NO LINE COMPRESSION (TOTAL BAN):** NEVER remove blank lines or compress statements to reduce line count.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions or trigger releases unless explicitly commanded.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** Accumulate all modifications and commit once at the final step.

---

## 5. Final Step Git Commit & Push Mandate

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS:** All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!
