---
name: cg-boolean-and-naming
description: Autonomously audits, refactors, and validates repository-wide boolean conventions, positive prefixes, implicit checks, enum Type suffixes, and nested if flattening against spec/02-coding-guidelines/.
---

# Skill: Coding Guidelines — Booleans, Naming & Enums (`cg-boolean`)

This skill governs autonomous execution for boolean conventions, semantic naming, enum standardization, and conditional flattening across all source files.

## Mandatory Architectural Rules

1. **Implicit Boolean Checks Only:**
   - NEVER write `if isReady == true` or `if (isValid === true)`.
   - Positive booleans MUST ALWAYS be evaluated implicitly: `if isReady { ... }` or `if !isReady { ... }`.
   - Never compare against boolean literals (`== false`, `!= true`).

2. **Boolean Prefixes (`is`, `has`):**
   - `is`, `has` as prefix is only acceptable and nothing else acceptable including but not limited to `can`, `should`, `was`, `will`, `did`, `must`, etc.
   - Every boolean identifier must begin with `is` or `has` (e.g. `isValid`, `hasAccess`).
   - No negative boolean identifiers (`isNotValid`, `hasNoData` are banned).

3. **No Inverted Success Checks:**
   - Never invert positive success checks (e.g. `!response.isSuccess`).
   - Use explicit failure states (e.g. `response.isFail`, `isError`).

4. **Zero Tolerance for Nested `if` (Nesting Depth <= 1):**
   - No `if` statements inside another `if` block.
   - Flatten all conditionals using guard clauses and early returns.
   - Never combine mixed polarity (`if isA && !isB` -> split into separate guard clauses).

5. **Enum Suffix `Type`:**
   - All enum declarations across TypeScript, Go, PHP must end with `Type` (e.g. `UserRoleType`, `CommandStatusType`).

6. **Function and File Size Caps:**
   - Functions: <= 8 lines preferred, <= 15 lines maximum.
   - Files: <= 100 lines coding maximum (recommended <= 80 lines).
   - Zero line compression (no single-line `if/else`, no deleted blank lines).

## Validation Linters

- Linter: `python linter-scripts/check-enum-and-boolean.py`
- Local Runner: `python .lovable/ai-fix-scripts/03-cicd-local-runner.py`
