# Master Coding Guidelines — Lint, enum sync, tests, lazy eval, regex, mutation, null safety, nesting, newlines, defer

> **Parent:** [Master Coding Guidelines](./00-overview.md)  
> **Version:** 2.1.0  
> **Updated:** 2026-03-31

---

## 11. Lint Scripts (Go)

| Script | Rule | Status |
|--------|------|--------|
| `linter-scripts/lint-file-size.sh` | No `.go` file > 300 lines | ✅ Active |
| `linter-scripts/lint-func-size.sh` | No function body > 15 lines | ✅ Active |
| `linter-scripts/lint-negative.sh` | No `IsNot*`, `HasNo*` function names | ✅ Active |

---


---

## 12. Cross-Language Enum Synchronization

Any modification to an enum must follow the enum-consumer-checklist <!-- external: spec/02-spec-management-software/18-enum-consumer-checklist.md -->:
1. Update PHP enum file
2. Update Go enum file (if mirrored)
3. Update TypeScript constants/types
4. Update database migration (if stored values change)
5. Update API documentation
6. Update admin templates referencing the enum

---


---

## 13. Test Naming & Structure

> Full reference: [14-test-naming-and-structure.md](../14-test-naming-and-structure.md)

| Rule | Description |
|------|-------------|
| R1 | Test file mirrors source file name with test suffix |
| R2 | Three-part naming: `Test{Unit}_{Scenario}_{ExpectedOutcome}` |
| R3 | 3+ similar scenarios → table-driven tests with named cases |
| R4 | Test helpers use descriptive verb prefixes; `t.Helper()` required (Go) |
| R5 | Test body follows Arrange / Act / Assert with blank line separators |
| R6 | Each test independently runnable — no shared mutable state |
| R7 | Integration tests skippable in short/fast mode |

---


---

## 14. Lazy Evaluation

> Full reference: [16-lazy-evaluation-patterns.md](../16-lazy-evaluation-patterns.md)

| Rule | Description |
|------|-------------|
| L1 | Use lazy init for expensive fields (DB queries, file parsing, computed constants) |
| L2 | Make lazy fields non-exported; expose via getter method |
| L3 | Always access lazy fields through getter — never direct field access |
| L4 | Use mutex lock for concurrent access (`MembersLock()` pattern) |
| L5 | If a required field is lazy, dependent fields MUST also be lazy |
| L6 | Do NOT apply lazy evaluation when data changes per request or is cheap to compute |

---


---

## 15. Regex Usage

> Full reference: [17-regex-usage-guidelines.md](../17-regex-usage-guidelines.md)

| Rule | Description |
|------|-------------|
| RX1 | Regex is **last resort** — prefer `strings.Contains()`, `HasPrefix()`, `Split()` for simple patterns |
| RX2 | Go: compile regex at package level (`var re = regexp.MustCompile(...)`) — never inside functions |
| RX3 | Add sample data as comments above regex declarations |
| RX4 | Never use regex in loops without reviewer approval |
| RX5 | TS: use regex literals for static patterns; `new RegExp()` only for dynamic |

---


---

## 16. Code Mutation Avoidance

> Full reference: [18-code-mutation-avoidance.md](../18-code-mutation-avoidance.md)

| Rule | Description |
|------|-------------|
| M1 | Every variable should be assigned **once** — prefer immutable values |
| M2 | No post-construction mutation — pass all values to constructor or use struct literals |
| M3 | When mutation is unavoidable, use mutex locks; provide both locked and unlocked versions |
| M4 | Never modify a variable across multiple methods — keep mutation in one method |

**Exemptions:** lazy evaluation/caching, loop accumulation (`append`), builder pattern, design pattern implementations.

---


---

## 17. Null Pointer Safety

> Full reference: [19-null-pointer-safety.md](../19-null-pointer-safety.md)

| Rule | Description |
|------|-------------|
| N1 | Always check `err` before using the returned value |
| N2 | Never call methods on unchecked return values — separate creation from execution |
| N3 | Check pointer for nil before dereference |
| N4 | Check array/slice for nil AND empty before index access |
| N5 | Nil-check before creating pointer returns |

See also §3.1 for the `isDefined`/`isDefinedAndValid` guard pattern.

---


---

## 18. Nesting Resolution

> Full reference: [20-nesting-resolution-patterns.md](../20-nesting-resolution-patterns.md)

| Rule | Description |
|------|-------------|
| NR1 | **Zero nested `if` is the target** — reinforces R2/R7 |
| NR2 | Method 1: Extract nested condition to named function |
| NR3 | Method 2: Inverse logic with early return — eliminate `else` |
| NR4 | Method 3: Decompose complex conditions into named boolean variables |
| NR5 | Exception: one loop with a single `if` inside is acceptable |

---


---

## 19. Newline Styling

> Full reference: [21-newline-styling-examples.md](../21-newline-styling-examples.md)

Supplements Code Style rules R4, R5, R10, R12, R13 with detailed before/after examples.

| Rule | Description |
|------|-------------|
| NL1 | Single-line function: no blank line before `return` needed |
| NL2 | Multi-line function: blank line before `return` **required** (R4) |
| NL3 | No empty lines between consecutive closing braces (R5 exception) |
| NL4 | No empty line at start of function (R12) |
| NL5 | No double empty lines anywhere (R13) |
| NL6 | Use `constants.NewLineUnix` (`"\n"`) in 90% of cases — `constants.NewLine` only for OS-specific needs |

---


---

## 20. Defer Rules (Go)

> Full reference: [05-defer-rules.md](../../03-golang/05-defer-rules.md)

| Rule | Description |
|------|-------------|
| D1 | Max **one** `defer` per function — defers are LIFO stack, multiple defers are confusing |
| D2 | Place `defer` at **top** or **bottom** of function — never buried in middle logic |
| D3 | Need multiple defers? Extract into separate functions, each with its own single defer |

---
