# Memory: style/code-quality-improvement
Updated: 2026-03-29

Process name: **Code Quality Improvement**. All coding guidelines are documented in `spec/05-coding-guidelines/01-code-quality-improvement.md`. Universal rules (all languages including TypeScript and Go):

1. **No magic strings** — comparison groups use enums (TS) or const groups (Go); standalone values use named constants.
2. **Exported object constants** — PascalCase (e.g., `WsTierLabels`, not `ws_tier_labels`).
3. **No inline type definitions** — always extract named types/interfaces for reusability; never define types in-place.
4. **Function length** — 8–25 lines max (no cramming multiple statements per line).
5. **Simple conditionals** — no negation (`!`, `!=`), no complex compound logic inline; extract into well-named boolean functions.
6. **Boolean naming** — always prefix with `is` or `has` (variables, functions, constants).
7. **Meaningful variable names** — no single-char names like `s`, `x`, `d`.
8. **Blank line before return** — unless return is the only line in an `if`.
9. **Self-documenting code** — if a section needs a comment, extract it into a function.
10. **File length** — max 200 lines; split by responsibility.

Go-specific rules remain in `spec/03-general/06-code-style-rules.md` and `spec/04-generic-cli/08-code-style.md`.
