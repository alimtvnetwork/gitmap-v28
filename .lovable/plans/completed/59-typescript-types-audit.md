# 38 - TypeScript Strict Typing & Discriminated Unions Audit Specification

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

1. **Rule TS-1 (Total Ban on `any`):** The `any` type is strictly forbidden. Use `unknown` with validation guards, generics `<T>`, or Discriminated Unions.
2. **Rule TS-2 (Discriminated Unions for State):** Multi-state models must use a common literal discriminator property (`status`, `kind`, or `type`) with exhaustive `assertNever` checks.
3. **Rule TS-3 (`as const` Enums with `*Type` Suffix):** Define enums as `as const` objects with `*Type` suffix and export the matching union type.
4. **Rule TS-4 (Immutability):** Mark configuration arrays, objects, and props as `readonly`.
5. **Rule TS-5 (Parameter Reduction):** Functions with >3 parameters must bundle arguments into a single typed interface.

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | src/ | N/A | TypeScript type checking | Verified `npx tsc --noEmit` exits with code 0 | VERIFIED |
| V-02 | src/ | N/A | Enum and boolean linting | Verified `node linter-scripts/check-enum-and-boolean.mjs` exits with code 0 | VERIFIED |
| V-03 | src/ | N/A | Production build compilation | Verified `npm run build` compiles cleanly in 9.85s | VERIFIED |
| V-04 | 03-ai-scripts/06-cicd-local-runner.py | N/A | Full CI quality gates | Verified all 16 gates exit with code 0 | VERIFIED |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/38-typescript/01-strict-types-and-discriminated-unions.md`)**:
   Verify total absence of `any`, validate discriminated unions, and check `as const` enums.
2. **Subtask 02 (`.lovable/plans/subtasks/38-typescript/02-tsc-and-build-verification.md`)**:
   Execute `npx tsc --noEmit` and `npm run build` to guarantee type safety and bundling.
3. **Subtask 03 (`.lovable/plans/subtasks/38-typescript/03-ci-runner-verification.md`)**:
   Execute full CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py`.
