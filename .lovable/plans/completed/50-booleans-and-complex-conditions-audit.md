# 29 - Boolean Principles, Negatives & Complex Conditions Audit Specification

## 1. Verbatim Acceptance Criteria Echo (from spec/02-coding-guidelines/01-cross-language/97-acceptance-criteria.md)

### AC-01: Guideline Coverage
- [ ] Boolean principles define naming, evaluation, and composition patterns
- [ ] Casting elimination patterns cover type-safe alternatives to type assertions
- [ ] Code style defines formatting, naming, and structural conventions

### AC-02: Enforcement
- [ ] All guidelines include ❌ (forbidden) and ✅ (compliant) code examples
- [ ] ESLint/linter rules are documented for automated enforcement
- [ ] Master guidelines document consolidates all standards for AI reference

---

## 2. Task-Specific Rule Set (Domain Rules)

1. **Rule BP-1 (Implicit Evaluation Only):** Booleans MUST NEVER be evaluated against literal `true` or `false` (`if isReady == true` is strictly forbidden; use `if isReady`).
2. **Rule BP-2 (Affirmative Naming & Prefix Restrictions):** All boolean variables MUST begin with `is` or `has` ONLY. Banned prefixes include `can`, `should`, `was`, `will`, `did`, `must`. Negative prefixes (`isNot`, `hasNo`, `disableFeature`) are forbidden.
3. **Rule BP-3 (Zero Inverted Success Checks):** Checking inverted success (`!response.isSuccess`) is strictly banned; use affirmative failure states (`response.isFail`).
4. **Rule BP-4 (Zero Mixed Polarity in Single Condition):** Combining positive and negative conditions in the same `if` statement (`if isA && !isB`) is strictly forbidden. Split into discrete guard clauses or extract a single-purpose helper function.
5. **Rule BP-5 (Zero Boolean Flag Parameters):** Functions must not accept primitive boolean flag arguments that alter fundamental control flow; split into dedicated semantic methods.

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | gitmap/cmd/ip_resolver.go | 27 | `skipLoopback bool` parameter | Refactored to affirmative `isSkipLoopback bool` and discrete helper decomposition | FIXED |
| V-02 | gitmap/cmd/ip_resolver.go | 65 | Mixed polarity and nested loop checks | Split into `resolveIPv4FromAddr` with discrete guard clauses | FIXED |
| V-03 | linter-scripts/check-boolean-guidelines.py | 23-26 | Regex definitions | Scanned repository; 0 violations across 3095 source files | VERIFIED |
| V-04 | 03-ai-scripts/06-cicd-local-runner.py | 31 | Missing boolean linter registration | Added `Boolean Guidelines Linter` to Batch 1 | FIXED |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/29-booleans/01-implicit-booleans.md`)**:
   Verify repository-wide implicit boolean evaluation (`if isReady`) with 0 `== true` comparisons.
2. **Subtask 02 (`.lovable/plans/subtasks/29-booleans/02-negative-inversion.md`)**:
   Enforce affirmative `is` and `has` naming and eliminate negative prefixes (`isNot*`, `hasNo*`).
3. **Subtask 03 (`.lovable/plans/subtasks/29-booleans/03-split-mixed-polarity.md`)**:
   Eliminate mixed polarity conditions and verify local CI gates via `python 03-ai-scripts/06-cicd-local-runner.py`.
