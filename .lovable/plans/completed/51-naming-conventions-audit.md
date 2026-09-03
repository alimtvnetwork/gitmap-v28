# 30 - Naming Conventions, Boolean Prefixes & Anti-Ok Variables Audit Specification

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

1. **Rule NC-1 (Mandatory `is`/`has` Boolean Prefixes):** Every boolean variable, parameter, struct field, or property MUST begin with `is` or `has` ONLY (`isValid`, `hasPermission`, `isReady`). All other prefixes (`can`, `should`, `was`, `will`, `did`, `must`) are strictly BANNED.
2. **Rule NC-2 (TOTAL BAN on Bare `ok` Identifiers):** In Go comma-ok idioms (type assertions, map lookups, channel receives), bare `ok` is strictly forbidden. Replace with semantic affirmative booleans (`isAppErr`, `isFound`, `hasValue`, `isChannelOpen`).
3. **Rule NC-3 (TOTAL BAN on Negative Boolean Identifiers):** Negative prefixes like `isNot*`, `hasNo*`, `disable*` are strictly banned. Use positive framing (`hasColors`, `hasPayload`, `isEnabled`).
4. **Rule NC-4 (Positive Framing with Inverted Guards):** Handle absence, empty states, or failure via positive boolean declaration and inverted `if` guard clauses (`hasColors := len > 0; if !hasColors { return }`).
5. **Rule NC-5 (Acronym & Filename Normalization):** All filenames must be lowercase. Acronyms must follow PascalCase (`UserId`, `ApiUrl`, `JsonData`, `IpResolver`).

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | gitmap/cliexit/kind.go | 96 | `code, ok := kindCodes[k]` | Replace bare `ok` with `isFound` | FIXED |
| V-02 | gitmap/cliexit/kind.go | 107 | `label, ok := kindLabels[k]` | Replace bare `ok` with `isFound` | FIXED |
| V-03 | gitmap/apperror/apperror.go | 87 | `_, file, line, ok := runtime.Caller(skip)` | Replace bare `ok` with `isCallerAvailable` | FIXED |
| V-04 | gitmap/cluster/exec_cmd.go | 56 | `exitErr, ok := err.(*exec.ExitError)` | Replace bare `ok` with `isExitErr` | FIXED |
| V-05 | gitmap/cluster/exec_lifecycle.go | 116 | `exitErr, ok := err.(*exec.ExitError)` | Replace bare `ok` with `isExitErr` | FIXED |
| V-06 | gitmap/cmd/agy_types.go | 86 | `t, ok := parseTimestampString(...)` | Replace bare `ok` with `isTimestampValid` | FIXED |
| V-07 | gitmap/cmd/ip_resolver.go | 27 | `skipLoopback bool` | Refactor to affirmative `isSkipLoopback bool` | FIXED |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/30-naming/01-bare-ok-refactoring.md`)**:
   Refactor bare `ok` variables across core packages (`cliexit`, `apperror`, `cluster`, `cmd/agy_types`) to semantic `is*` and `has*` booleans.
2. **Subtask 02 (`.lovable/plans/subtasks/30-naming/02-affirmative-boolean-prefixes.md`)**:
   Audit and enforce strict `is`/`has` prefixes and positive framing across command handlers.
3. **Subtask 03 (`.lovable/plans/subtasks/30-naming/03-ci-and-linter-verification.md`)**:
   Verify `python linter-scripts/check-enum-and-boolean.py`, `python linter-scripts/check-boolean-guidelines.py`, and run `python 03-ai-scripts/06-cicd-local-runner.py`.
