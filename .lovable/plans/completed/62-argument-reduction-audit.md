# 42 - Argument Reduction, Parameter Structs & Return Architecture Audit Specification

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

1. **Rule AR-1 (Parameter Reduction via Structs):** Multi-argument signatures with >2–3 parameters must be grouped into dedicated value-based parameter Structs (`*Params`) or parameter objects.
2. **Rule AR-2 (Affirmative Boolean Fields):** All boolean fields in parameter structs must use affirmative prefixes (`is` and `has` only).
3. **Rule AR-3 (Mandatory `*apperror.AppError` Returns):** Bare "void" functions in domain and service logic must return `*apperror.AppError` for side-effect operations.
4. **Rule AR-4 (Framework Error Wrapping):** External and stdlib `error` returns must be converted and wrapped into `*apperror.AppError` context wrappers.
5. **Rule AR-5 (Single `Result[T]` Envelopes):** Functions that produce data must return single `Result[T]` envelopes with complete predicate methods.

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | gitmap/cmd/rootusagefooter.go | 100-109 | `IdentityRowParams` parameter struct | Value-based parameter struct implementation | VERIFIED |
| V-02 | gitmap/result/result.go | 1-132 | `Result[T]` generic envelope | Full generic result envelope with predicates | VERIFIED |
| V-03 | gitmap/apperror/apperror.go | 1-268 | `*AppError` wrapping and constructors | Full contextual error wrapping and predicates | VERIFIED |
| V-04 | 03-ai-scripts/06-cicd-local-runner.py | N/A | CI quality gates verification | Verified all 16 quality gates exit with code 0 | VERIFIED |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/42-argument-reduction/01-parameter-struct-reduction.md`)**:
   Enforce value-based parameter structs (`*Params`) for functions with >2–3 parameters and affirmative boolean field prefixes.
2. **Subtask 02 (`.lovable/plans/subtasks/42-argument-reduction/02-apperror-and-result-returns.md`)**:
   Mandate `*apperror.AppError` returns for side-effect operations and `Result[T]` for data-producing functions.
3. **Subtask 03 (`.lovable/plans/subtasks/42-argument-reduction/03-ci-runner-verification.md`)**:
   Execute full CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py`.
