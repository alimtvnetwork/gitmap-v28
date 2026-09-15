---
name: cg-result-wrapper
description: Autonomously scan, audit, refactor, and verify Go functions returning multi-value error tuples, replacing them with strongly-typed Result, ResultMap, and ResultSlice wrappers and structured AppError returns.
---

# Skill: Result Wrapper Types, Collections & AppError Returns (`cg-result-wrapper`)

This skill governs autonomous scanning, auditing, refactoring, and verification of Go functions returning multi-value tuples pairing return values with standard library errors (such as `(map[K]V, error)`, `([]T, error)`, or `(T, error)`).

## Core Architectural Directives

1. **Single Return Object Mandate:**
   - Multi-value error returns (e.g. `(map[string]int, error)`, `([]string, error)`, `(MyType, error)`) are prohibited on domain, store, and service functions.
   - All operations must return a single, strongly-typed envelope:
     - Key-Value Maps: `result.ResultMap[K, V]`
     - Slices / Lists: `result.ResultSlice[T]`
     - Scalar Values: `result.Result[T]`
     - Pure Side-Effects: `*apperror.AppError` (zero bare `error` returns)

2. **Structured Domain Errors:**
   - Raw standard library `error` returns are strictly banned in domain, store, and service layers.
   - All errors must use structured `*apperror.AppError` created with `apperror.New()` or `apperror.Wrap()`.

3. **Standardized Outer-Layer Inspection Predicates:**
   - `res.IsSuccess()`: `bool` - operation completed without error.
   - `res.IsFailure()` / `res.IsFailed()`: `bool` - operation encountered an error.
   - `res.HasError()`: `bool` - active error attached.
   - `res.IsEmptyError()` / `res.HasNoError()`: `bool` - no active error attached.
   - `res.IsEmpty()`: `bool` - collection has 0 items or data is empty.
   - `res.Data` / `res.Value`: Direct access to underlying payload.
   - `res.AppError()` / `res.Fault()`: Retrieves structured `*apperror.AppError`.
   - `res.Get(key)`: Safely retrieves map element `(V, bool)`.
   - `res.Has(key)`: Reports whether key exists in map.
   - `res.Count()`: Returns length of collection.
   - `res.Keys()`: Returns deterministically sorted slice of keys.
   - `res.Values()`: Returns slice of values ordered by keys.

4. **Zero Inversion & Affirmative Predicates:**
   - Never write `if !res.IsSuccess()`.
   - Always write `if res.IsFailure() { ... }` or `if res.IsFailed() { ... }`.

5. **Guard Clauses & Zero Nested Ifs:**
   - Maximum nesting depth is 1 (`if` within `if` is prohibited).
   - Use early returns and guard clauses.
   - Separate consecutive guard clauses with blank lines.
   - Always insert a blank line before `return`.

6. **Targeted Verification & CI Protection:**
   - DO NOT run `06-cicd-local-runner.py` during routine refactoring turns.
   - Run targeted linters (`golangci-lint run <pkg>`, `go vet <pkg>`, `python linter-scripts/...`).
   - Append modified files to `.lovable/temp/recent-file-changes.json` under lock.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
