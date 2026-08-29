# Subtask 01: Refactor Explicit Boolean Comparisons & Inverted Success Checks

> **Parent Plan:** `17-boolean-and-naming-audit.md`  
> **Scope:** All repository files with `== true`, `== false`, `=== true`, `=== false`, and `!isSuccess`

## Objective

Refactor all explicit boolean literal checks to implicit evaluation and replace inverted success checks with explicit positive failure variables.

## Action Steps

1. Scan for and replace all `== true` and `=== true` checks with direct evaluation (`if isFoo`).
2. Scan for and replace all `== false` and `=== false` checks with negation or explicit negative flags (`if !isFoo` or `if isFail`).
3. Replace all `!isSuccess` checks with explicit error states (`isFail`, `isError`).
4. Ensure zero compilation errors or logic inversions.
