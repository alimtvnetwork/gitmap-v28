# 40 - Multi-Language Enums, Traits & Pattern Matching Audit Specification

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

1. **Rule EN-1 (`*Type` Suffix):** All Enum and enum-like types in Go, TypeScript, PHP, Rust, and Python must end with `Type`.
2. **Rule EN-2 (String-Backed Enums in PHP):** All PHP enums must be string-backed with `HasEnumHelpers` trait.
3. **Rule EN-3 (Exhaustive Pattern Matching):** `match` and `switch` statements must explicitly handle all variants.
4. **Rule EN-4 (`as const` in TypeScript):** Object enums must use `as const` definitions.
5. **Rule EN-5 (Zero Loose Literals):** Zero raw string or numerical literals where a typed enum exists.

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | linter-scripts/check-enum-guidelines.py | N/A | Enum and constant conventions | Verified exit code 0 | VERIFIED |
| V-02 | linter-scripts/check-enum-and-boolean.py | N/A | Boolean, enum, and conditional compliance | Verified exit code 0 across 2,202 files | VERIFIED |
| V-03 | node linter-scripts/check-enum-and-boolean.mjs | N/A | TypeScript enum validation | Verified exit code 0 | VERIFIED |
| V-04 | 03-ai-scripts/06-cicd-local-runner.py | N/A | Local CI runner execution | Verified all gates exit 0 | VERIFIED |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/40-enums/01-enum-type-suffix-and-traits.md`)**:
   Verify `*Type` suffixes, traits, and exhaustive match handling.
2. **Subtask 02 (`.lovable/plans/subtasks/40-enums/02-linters-and-ci-runner.md`)**:
   Execute enum linters and full CI runner suite.
