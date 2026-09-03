# 32 - React & Frontend Architecture Audit Specification

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

1. **Rule RF-1 (Component Modular Sizing):** Components should ideally target <= 80 lines (standard max <= 100 lines). Decompose large UI blocks into single-responsibility child components.
2. **Rule RF-2 (Hook Return Signatures):** Custom React hooks MUST NOT return array tuples (`[val, setVal]` is banned). Custom hooks MUST return named property objects (`{ val, onUpdate }`).
3. **Rule RF-3 (Zero Redundant useEffect):** Eliminate redundant `useEffect` hooks used for derived calculations or local state synchronization. Derive state inline during rendering.
4. **Rule RF-4 (Immutable State Updates):** Direct state mutations are strictly banned. All React state transitions must use immutable spreads or pure functional updater functions.
5. **Rule RF-5 (Zero Nested Conditionals):** Zero tolerance for nested `if` statements inside handlers and hooks. Flatten with early returns and guard clauses.

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | src/hooks/ | N/A | Custom hook return signatures | Verified 0 tuple returns across all custom hooks in `src/hooks/` | VERIFIED |
| V-02 | src/ | N/A | React component compilation | Verified `npm run build` succeeds cleanly with 0 errors | VERIFIED |
| V-03 | linter-scripts/check-enum-and-boolean.mjs | N/A | Frontend enum & boolean compliance | Verified exit code 0 | VERIFIED |
| V-04 | 03-ai-scripts/06-cicd-local-runner.py | 43 | Frontend build quality gate | Integrated Web App Build gate in Batch 2 | VERIFIED |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/32-react-frontend/01-custom-hook-signatures.md`)**:
   Verify that all custom hooks in `src/hooks/` return named property objects rather than tuple arrays.
2. **Subtask 02 (`.lovable/plans/subtasks/32-react-frontend/02-immutable-state-and-useeffect.md`)**:
   Verify absence of direct state mutations and redundant `useEffect` derived state syncing.
3. **Subtask 03 (`.lovable/plans/subtasks/32-react-frontend/03-ci-frontend-build-verification.md`)**:
   Execute `node linter-scripts/check-enum-and-boolean.mjs`, `npm run build`, and `python 03-ai-scripts/06-cicd-local-runner.py`.
