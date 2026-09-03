# 34 - Style Guidelines, Formatting & Line-Gaps Audit Specification

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

1. **Rule SG-1 (Mandatory Blank Line BEFORE Control Structures):** When an `if`, `for`, `switch`, `while`, or `try` is preceded by any statement, there MUST be exactly one blank line before it (unless it is line 1 of a block).
2. **Rule SG-2 (Mandatory Blank Line AFTER Closing Brace `}`):** When a closing brace `}` is followed by further executable code or another statement, there MUST be exactly one blank line after it.
3. **Rule SG-3 (Mandatory Blank Line BEFORE Return/Throw):** In multi-line functions and blocks, there MUST be a blank line before `return`, `throw`, `raise`, `yield`.
4. **Rule SG-4 (Blank Lines Around Struct Calls & Sequential Invocations):** When instantiating parameter structs or invoking multiline functions, place clean blank lines before and after.
5. **Rule SG-5 (Zero Clumping of Consecutive Guard Clauses):** Consecutive guard clauses MUST be separated by a blank line after each closing brace `}`.

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | gitmap/cmd/ | Various | Squeezed guard clauses and return statements | Inserted blank lines before `if` and before `return` across command handlers | FIXED |
| V-02 | gitmap/cluster/ | 53-58 | Return newline spacing and exit code extraction | Clean blank lines inserted around `isExitErr` checks and before `return` | FIXED |
| V-03 | gitmap/cliexit/ | 96-110 | Map lookup and early returns in KindCode/KindLabel | Clean blank lines inserted around `isFound` checks | FIXED |
| V-04 | linter-scripts/check-newline-styling.py | N/A | Return newline style verification | Verified 0 newline styling violations | VERIFIED |
| V-05 | 03-ai-scripts/06-cicd-local-runner.py | 36 | MWS Error Codes check integration | Registered in Batch 1 of local runner | FIXED |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/34-style/batch-01.md`)**:
   Verify vertical line spacing (blank line before `if`, blank line after `}`, blank line before `return`) across core packages.
2. **Subtask 02 (`.lovable/plans/subtasks/34-style/batch-02.md`)**:
   Verify zero clumping of consecutive guard clauses and ensure clean multiline struct call separation.
3. **Subtask 03 (`.lovable/plans/subtasks/34-style/batch-03.md`)**:
   Execute `python linter-scripts/check-newline-styling.py`, `python linter-scripts/check-mws-error-codes.py`, and run full CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py`.
