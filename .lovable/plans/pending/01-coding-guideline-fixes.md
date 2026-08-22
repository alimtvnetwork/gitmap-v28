# Master Plan: Coding Guideline Fixes

This master plan details the codebase-wide audit for coding guideline violations across `gitmap` and `src` directories.

## Summary of Audit Findings
- **Inverted Booleans:** Found and queued 50 instances of inverted booleans (`isNot* := !is*`).
- **Nested Ifs:** Found and queued 50 instances of nested `if` statements requiring flattening.
- **Monolithic Functions:** Identified 2,200+ monolithic functions; queued 50 for refactoring.
- **Bare Bools:** Found 800+ instances of returning bare `true`/`false`; queued 50 for wrapping in `Result` objects.

## Subtasks
- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/01-inverted-booleans.md`
- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/02-nested-ifs.md`
- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/03-monolithic-functions.md`
- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/04-bare-bools.md`


- **Single-Character Vars:** Found 2,300+ instances of single character variables; queued 50 for renaming.
- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/05-single-character-vars.md`

Total queued exact file/line changes: 450.

- **Magic Strings and Numbers:** Found and queued 50 instances.
- **Swallowed Errors:** Found and queued 50 instances of missing error context or empty catch blocks.
- **Banned Short Identifiers:** Found and queued 50 instances of banned short variable names.
- **Missing Enum/Type Suffixes:** Found and queued 50 instances of types missing suffixes.

- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/06-magic-strings.md`
- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/07-swallowed-errors.md`
- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/08-banned-short-identifiers.md`
- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/09-missing-type-suffixes.md`
