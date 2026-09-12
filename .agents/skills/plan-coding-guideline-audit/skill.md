---
name: plan-coding-guideline-audit
description: Plan a structured coding guideline audit across repository codebases against spec/02-coding-guidelines/.
---

# Plan Coding Guideline Audit

Autonomously plans a comprehensive audit of repository codebases against the master coding guidelines.

## Audit Areas
1. **Boolean Conventions:** `is` and `has` prefixes only, implicit positive checks, no mixed polarity.
2. **Control Flow:** Maximum nesting depth 1, guard clauses, flattened conditionals.
3. **Naming & Types:** Enum `Type` suffixes, PascalCase types, no underscores in Go.
4. **Error Handling:** `*apperror.AppError` returns, no swallowed errors, typed error responses.
5. **Code Metrics:** Functions <= 15 lines, files <= 100 lines coding, blank line padding.

## Output
Generates structured audit logs and phased remediation plans in `.lovable/plans/pending/` with subtask micro-batches.
