# Plan 139: Nested If Elimination & Guard Clauses Coding Guideline (Consolidated Milestone)

> **Originating Request:** `# Nested If Elimination & Guard Clauses — Coding Guideline (must follow)`
> **Trigger Keywords & Aliases:** `cg-nested-if`, `cg-execute nested-if`, `audit nested if`, `fix nested if`, `flatten conditionals`, `enforce guard clauses`
> **Prompt Version:** 2.1.0
> **Execution Budget & Loop Trace:** Completed across 2 phases and 3 subtasks within self-looping budget of N=200 steps.

---

## 1. Executive Summary of Implementation

1. **Comprehensive Polyglot AST Audit**:
   - **Go Source Code (`cli/`, `cli-updater/`)**: Audited all 2,860 Go files using `go/parser` AST inspection; confirmed **0** nested `if` statements and zero single-line collapsed `if` statements.
   - **TypeScript / React Frontend (`src/`)**: Audited all 205 `.ts`/`.tsx` files using TypeScript Compiler API AST inspection; identified 6 nested `if` statements across 3 files (`TabOrderMap.tsx`, `clipboard.ts`, `Settings.tsx`).

2. **Refactoring of Violations & Guard Clause Enforcements**:
   - `src/components/docs/TabOrderMap.tsx`:
     - Inverted outer form element check in `getFormElementLabel` into early return guard clause `if (!isFormEl) return null;` and `if (!formEl.labels || formEl.labels.length === 0) return null;`.
     - Inverted inner label text check `if (!text) return null;`, completely flattening to depth 0.
     - Extracted `getAriaLabelledByText(section)` helper function with early guard returns, flattening `getAriaSection` from depth 2 to depth 0.
     - Replaced single-line unbraced `if/else` statement in `groupEntriesBySection` with explicit multi-line braces and spacing.
   - `src/lib/clipboard.ts`:
     - Extracted `isModernClipboardSupported()` helper with early return guards.
     - Extracted `tryModernClipboardCopy(text)` helper wrapping async `queryWrapper` call.
     - Flattened `copyToClipboard(text)` into a single guard clause with early exit, removing nested `if (!res.isFail && res.data)`.
   - `src/pages/Settings.tsx`:
     - Inverted outer success check in `useEffect` with early return guard clause `if (res.isFail || !res.data) { return; }`.
     - Expanded inner settings state setters with explicit curly braces `{}`.

3. **Linter Infrastructure & Skill Registration**:
   - Created `.agents/skills/cg-nested-if/skill.md` defining autonomous guard clause flattening rules and anti-compression mandates.
   - Created `linter-scripts/check-nested-ifs.mjs` running TypeScript Compiler API AST checks over all TypeScript source files.
   - Upgraded `linter-scripts/check-nested-ifs.py` to invoke `check-nested-ifs.mjs`, ensuring polyglot coverage for both Go and TypeScript.

---

## 2. Granular Subtask Trace & Resolution

- **Subtask 01 (`01-audit-and-scan-nested-ifs.md`)**:
  - Codebase scanned for conditional nesting depth > 1 across Go and TypeScript.
  - Authored `.agents/skills/cg-nested-if/skill.md`.
  - Authored `linter-scripts/check-nested-ifs.mjs` and updated `linter-scripts/check-nested-ifs.py`.
- **Subtask 02 (`02-flatten-conditionals-and-guard-clauses.md`)**:
  - Refactored `src/components/docs/TabOrderMap.tsx`, `src/lib/clipboard.ts`, and `src/pages/Settings.tsx`.
  - Verified 0 nested `if` statements remaining in TypeScript AST.
- **Subtask 03 (`03-verify-linters-and-anti-compression.md`)**:
  - Ran `check-nested-ifs.py` (0 violations across 2,863 Go files and 205 TypeScript files).
  - Ran `npx tsc --noEmit` (clean compilation, exit code 0).
  - Ran `check-boolean-guidelines.py` (0 violations).
  - Ran `check-relative-paths.py` (0 absolute paths across 6,823 files).
  - Ran `35-result-wrapper-auditor.py` (0 violations).
  - Ran `go vet ./...` in `cli/` (clean).

---

## 3. Violation Resolution Ledger

| Id | File | Line | Violation Description | Resolution | Status |
|:---|:---|:---:|:---|:---|:---:|
| V-01 | `src/components/docs/TabOrderMap.tsx` | 144 | Nested `if` inside `isFormEl` check | Inverted to guard clauses `if (!isFormEl)`, `if (!formEl.labels)`, `if (!text)` | Verified |
| V-02 | `src/components/docs/TabOrderMap.tsx` | 215 | Nested `if` inside `labelId` check | Extracted `getAriaLabelledByText` helper with early returns; flattened `getAriaSection` | Verified |
| V-03 | `src/components/docs/TabOrderMap.tsx` | 306-307 | Single-line unbraced `if/else` | Expanded with curly braces and proper spacing | Verified |
| V-04 | `src/lib/clipboard.ts` | 22 | Nested `if` inside clipboard API check | Extracted `isModernClipboardSupported` and `tryModernClipboardCopy` helpers | Verified |
| V-05 | `src/pages/Settings.tsx` | 27 | Nested `if` inside success check | Inverted outer check `if (res.isFail \|\| !res.data) { return; }` | Verified |
| V-06 | `src/pages/Settings.tsx` | 28 | Nested `if` inside success check | Guard clause flattened; added braces | Verified |
| V-07 | `src/pages/Settings.tsx` | 29 | Nested `if` inside success check | Guard clause flattened; added braces | Verified |
| V-08 | `linter-scripts/check-nested-ifs.py` | 200 | TS/TSX files ignored in nested if linter | Added `check-nested-ifs.mjs` and integrated with `check-nested-ifs.py` | Verified |

---

## 4. Final Quality Gates & Verification Evidence

```text
=== Running Nested If Linter (check-nested-ifs.py) ===
✅ PASS: Zero nested if statements or single-line compression violations found across 2863 files in 1.48s.
▸ TypeScript AST check: ✅ PASS: Zero nested if statements found across 205 TypeScript file(s).

=== TypeScript Type Check ===
npx tsc --noEmit -> Exit code 0 (clean)

=== Boolean Guidelines Linter ===
✅ PASS: Zero boolean guideline violations found across 2863 files.

=== Relative Paths Linter ===
✅ PASS: No absolute filesystem paths or file:/// URIs found across 6823 files (2.57s).

=== Result Wrapper Auditor ===
✅ PASS: Zero legacy multi-value map tuple error returns found. ResultMap envelopes verified.

=== Go Vet ===
go vet ./... in cli/ -> Exit code 0 (clean)
```
