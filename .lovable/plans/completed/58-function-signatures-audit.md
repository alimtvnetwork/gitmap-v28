# 37 - Function Signatures, Invocations & Result Envelopes Audit Specification

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

1. **Rule FS-1 (Multi-Line Parameter Declarations):** Definitions with >2 parameters must format with exactly one parameter per line with trailing commas.
2. **Rule FS-2 (Multi-Line Function Invocations):** Call sites with >2 arguments must format with exactly one argument per line with trailing commas.
3. **Rule FS-3 (Single Result Envelope):** Domain services encapsulate execution returns in `Result[T]` envelopes containing typed value and `*apperror.AppError`.
4. **Rule FS-4 (Complete Predicate Methods):** Result envelopes provide `IsSuccess()`, `IsFailed()`, `HasError()`, `HasNoError()`, and `HasValidError()`.
5. **Rule FS-5 (Semantic Verb & Predicate Naming):** Action functions start with active verbs; boolean predicates start with `is` and `has` ONLY (`can` is banned).

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | gitmap/result/result.go | 1-132 | `Result[T]` generic envelope implementation | Implemented with all predicates (`IsSuccess`, `IsFailed`, `HasError`, `HasNoError`, `HasValidError`) | VERIFIED |
| V-02 | gitmap/apperror/apperror.go | 239-267 | `AppError` predicate methods | Implemented `HasError`, `HasNoError`, `HasValidError`, `IsErrorCode` | VERIFIED |
| V-03 | gitmap/result/result_test.go | 1-90 | Unit tests for `Result[T]` envelope | Verified with `go test -C gitmap ./result/...` | VERIFIED |
| V-04 | gitmap/cmd/ | Various | Multi-line call site formatting | Verified one argument per line for multi-parameter calls | VERIFIED |
| V-05 | 03-ai-scripts/06-cicd-local-runner.py | N/A | Full CI quality gate execution | Verified exit code 0 | VERIFIED |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/37-function-signatures/01-multi-line-formatting-and-signatures.md`)**:
   Format function definitions and invocations with >2 parameters/arguments to one per line.
2. **Subtask 02 (`.lovable/plans/subtasks/37-function-signatures/02-result-envelope-and-apperror.md`)**:
   Verify `Result[T]` envelope and `AppError` predicate methods.
3. **Subtask 03 (`.lovable/plans/subtasks/37-function-signatures/03-ci-runner-verification.md`)**:
   Execute `go test -C gitmap ./result/...` and full CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py`.
