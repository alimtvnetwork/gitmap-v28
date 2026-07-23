# 01 — Code Quality

Universal rules enforced across all languages (TypeScript, Go).

## Rules

| # | Rule | Scope |
|---|------|-------|
| 1 | No magic strings — enums or constants | All |
| 2 | Exported object constants — PascalCase | All |
| 3 | No inline type definitions — extract named types | TS |
| 4 | Function length — 8 to 15 lines (Go), 8 to 25 lines (TS) | All |
| 5 | Simple conditionals — no negation (`!`, `!=`) | All |
| 6 | Boolean naming — `is` / `has` prefix | All |
| 7 | Meaningful variable names — no single-char | All |
| 8 | Blank line before `return` | All |
| 9 | Self-documenting code — no restating comments | All |
| 10 | File length — max 200 lines | All |

## No Magic Strings

Every string literal used for comparison, defaults, labels, or keys must live in a dedicated constants file. Comparison groups use enums (TS) or const groups (Go). Standalone values use named constants.

## Simple Conditionals

No `!`, no `!=`, no negative function names like `isNotValid`. Extract complex compound conditions into well-named boolean functions. Even `=== false` is preferred over `!` for positive-logic readability.

## Function Length

Target 8–15 lines in Go, 8–25 lines in TypeScript. Extract helpers when exceeding limits. Each function does one thing. Do not cram multiple statements per line.

## File Length

Max 200 lines per source file. Split by responsibility when exceeded.

---

Source: `spec/05-coding-guidelines/01-code-quality-improvement.md`
