# 12 — Documentation Standards

Inline comments, README structure, API docs, and spec conventions.

## Inline Comments

Comment the *why*, not the *what*. Use `//` on the line above. Markers: `TODO:`, `FIXME:`, `HACK:`, `NOTE:`.

No commented-out code in production.

## README Structure

Required sections in order: Title, Overview, Quick Start, Usage, Configuration, Development, License. Quick Start under 10 lines.

## Function Docs

Every exported function includes: what it does (one line), non-obvious parameters, return value and error conditions.

## Package Docs

Each package has a doc comment explaining purpose in 1–3 sentences.

## CLI Help Text

Every command includes: one-line description, usage line, flag descriptions with defaults, at least one example.

## Spec Documents

Numeric prefixes for ordering. Required structure: Title (single H1), Overview, Behavior, Implementation, Constraints. Use tables over prose.

## Constraints

- Comments explain *why*, not *what*
- READMEs under 300 lines
- Every exported symbol has a doc comment
- Help text includes at least one example

---

Source: `spec/05-coding-guidelines/12-documentation-standards.md`
