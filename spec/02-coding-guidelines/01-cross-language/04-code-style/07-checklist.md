# Code Style Checklist & Cross-References

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

## Cross-References

- [No Raw Negations](../12-no-negatives.md) — Positive guard functions instead of `!` (all languages)
- [Function Naming](../10-function-naming.md) — No boolean flag parameters (all languages)
- [Strict Typing](../13-strict-typing.md) — Type declarations, max 3 parameters (all languages)
- [Boolean Principles](../02-boolean-principles/00-overview.md) — P1–P6 boolean naming rules (all languages)
- [Go Enum Specification](../../03-golang/01-enum-specification/00-overview.md) — Go enum pattern, required methods, folder structure
- [TypeScript Enums](../../02-typescript/00-overview.md) — TypeScript string enum definitions and usage patterns
- [PHP Enum Classes](../../04-php/01-enums.md) — PHP backed enum patterns
- [PHP Coding Standards](../../04-php/07-php-standards-reference/00-overview.md) — PHP-specific rules that reference this spec
- [PHP Forbidden Patterns](../../04-php/02-forbidden-patterns.md) — PHP checklist

---

*Part of [Code Style](./00-overview.md)*
