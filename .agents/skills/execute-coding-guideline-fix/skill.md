---
name: execute-coding-guideline-fix
description: Execute coding guideline refactoring, fixing booleans, nesting, naming, and function sizes in 5-8 file micro-batches.
---

# Execute Coding Guideline Fix

Autonomously refactors code violations against `spec/02-coding-guidelines/` in strictly bounded 5-8 file micro-batches.

## Rules
- **No Line Compression:** Maintain mandatory blank lines before `if`, after `}`, before `return`.
- **Implicit Booleans:** Never write `== true`. Replace with implicit checks.
- **Guard Clauses:** Invert early checks to return immediately and flatten nested blocks.
- **Go Errors:** Return `*apperror.AppError` and preserve full error causal chains.
- **Verification:** Run `python 03-ai-scripts/06-cicd-local-runner.py --no-tests --no-tests` after each batch (test execution is disabled unless explicitly commanded by the repository owner).
