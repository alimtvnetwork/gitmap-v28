# Code Style Checklist & Cross-References (AI Execution Prompt)

> **/goal** Provide a consolidated, actionable verification checklist and cross-reference catalog covering all cross-language code style rules (braces, blank lines, function sizing, multi-line formatting, and dead code elimination).
> **/learn** Internalize the complete taxonomy of code style rules, enforce mandatory pull-request checklist items, and verify compliance across PHP, TypeScript, and Go implementations.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify pull requests and code modifications against the 19-point consolidated code style checklist.
- [ ] `/learn` Ensure zero tolerance for nested `if` statements, missing braces, or omitted blank line vertical separators.
- [ ] `/goal` Validate all signatures and calls with >2 arguments use multi-line formatting with trailing commas.
- [ ] `/learn` Review cross-references across Golang, TypeScript, and PHP specifications for language-specific nuances.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version:** 4.0.0
> **Updated:** 2026-03-31

---

## Checklist Summary (Copy for PRs)

```
[ ] No single-line `if (...) return;` — always use braces
[ ] No nested `if` — ZERO TOLERANCE — flatten with early returns or combined conditions
[ ] No inline multi-part `if` (2+ operators) — extract to named variable or method
[ ] Blank line before `return` or `throw` when preceded by other statements
[ ] Blank line after closing `}` when followed by more code
[ ] Functions max 15 lines — extract helpers for longer logic
[ ] Error handling lines exempt from 15-line count (guards, wrapping, context)
[ ] No deeply nested control flow — extract loop/condition bodies to helpers
[ ] No leading backslash on `Throwable` or other global types in catch/type hints
[ ] Functions/calls with >2 args — one arg per line with trailing comma (signatures AND calls)
[ ] PHP arrays with >2 items — each item on its own line with trailing comma
[ ] Blank line before control structures (`if`/`for`/`foreach`/`while`) when preceded by statements
[ ] Method chaining — each `.Method()` on its own line (>2 calls)
[ ] apperror.Wrap/New with 2+ args — always multi-line, each arg on its own line
[ ] Nested apperror.Fail[T](Wrap(...)) — always expanded to multi-line
[ ] No commented-out or dead code — delete it, use version control
[ ] Space after `//` in all line comments
[ ] Doc comments on all exported/public functions and methods
[ ] Struct/class files max 120 lines — split by concern when exceeded
```

---

## Core Rules & Code Comparisons

### Single-Line Statements & Braces

- ❌ FORBIDDEN:
  ```go
  if err != nil { return err }
  ```
- ✅ REQUIRED:
  ```go
  if err != nil {
      return err
  }
  ```

---

## Cross-References

- [No Raw Negations](../12-no-negatives.md) — Positive guard functions instead of `!` (all languages)
- [Function Naming](../10-function-naming.md) — No boolean flag parameters (all languages)
- [Strict Typing](../13-strict-typing.md) — Type declarations, max 3 parameters (all languages)
- [Boolean Principles](../02-boolean-principles/readme.md) — P1–P6 boolean naming rules (all languages)
- [Go Enum Specification](../../03-golang/01-enum-specification/readme.md) — Go enum pattern, required methods, folder structure
- [TypeScript Enums](../../02-typescript/readme.md) — TypeScript string enum definitions and usage patterns
- [PHP Enum Classes](../../04-php/02-enums.md) — PHP backed enum patterns
- [PHP Coding Standards](../../04-php/07-php-standards-reference/readme.md) — PHP-specific rules that reference this spec
- [PHP Forbidden Patterns](../../04-php/03-forbidden-patterns.md) — PHP checklist

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-STYLE-008: Comprehensive Code Style Checklist

**Given** Any pull request or modified source file within the repository across PHP, TypeScript, or Go.
**When** Static analysis, guideline autofixers, or CI linting checks run against the changes.
**Then** All 19 checklist items pass with zero violations, including braces, blank lines, multi-line arguments, line counts, and dead code removal.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.

---

*Part of [Code Style](./readme.md)*
