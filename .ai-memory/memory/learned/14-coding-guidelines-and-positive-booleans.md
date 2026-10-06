# 14 — Coding Guidelines and Positive Boolean Conventions

- **Subsystem:** Quality Governance & Linter Standards
- **Status:** Authoritative Reference

## 1. Affirmative Boolean Prefixing
- All boolean variables, struct fields, and functions must use positive prefixes: `is*` or `has*` (e.g. `isConsolidated`, `hasVerifiedOutcome`).
- Total ban on double negatives or negative polarity flags (`isNotValid`, `hasNoErrors`).

## 2. Function Sizing & Inverted Guard Clauses
- Functions strictly bounded to <= 8–15 lines.
- Nested `if/else` ladders refactored into early return guard clauses.
- Mandatory vertical whitespace: blank lines before `if` statements and return statements.
