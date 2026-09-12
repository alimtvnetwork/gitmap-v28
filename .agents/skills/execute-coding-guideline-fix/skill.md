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
- **Targeted Batch Verification:** Run targeted linter / autofixer on the modified files in the batch (e.g. `python linter-scripts/check-nested-ifs.py <files>` or `08-naming-autofixer.py <files>`). DO NOT run `06-cicd-local-runner.py` or full builds (`npm run build`, `go build`) during routine batch fixes.
- **Consolidated Commits:** NEVER commit isolated 1-2 plan/doc files alone. Commit all modified source files, tests, and plans together as a single atomic batch.
- **Immediate Git Push:** Always push immediately to remote (`git push origin <branch>`) after every commit.
- **Atomic Change Cache Recording:** Record all modified files to `.lovable/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`).
