---
name: execute-coding-guideline-fix
description: Execute coding guideline refactoring, fixing booleans, nesting, naming, and function sizes in 5-8 file micro-batches.
---

# Execute Coding Guideline Fix

Autonomously refactors code violations against `02-spec/02-coding-guidelines/` in strictly bounded 5-8 file micro-batches.

## Rules

- **No Line Compression:** Maintain mandatory blank lines before `if`, after `}`, before `return`.
- **Implicit Booleans:** Never write `== true`. Replace with implicit checks.
- **Affirmative Boolean Parameters & Fields:** Ban single-letter (`v bool`, `b bool`, `val bool`), bare identifiers (`stop bool`, `pause bool`, `defined bool`), and negative identifiers (`isNot*`, `isUndefined`, `hasNo*`). Rename to `isStopOnFail bool`, `isStopped bool`, `isPaused bool`, `isDefined bool`. Try `isDefined` / `IsDefined` instead of negatives (e.g., use `isDefined` instead of `isUndefined` or `isNotDefined`, and invert at callsite with `!isDefined` if testing for absence; use `isValid` instead of `isNotValid`, `hasValue` instead of `hasNoValue`).
- **IsDefined over `!isEmpty`:** MANDATORY: Use `isDefined` (or `res.IsDefined()`) instead of inverted empty checks (`!isEmpty` / `!res.IsEmpty()`). Never write `if !isEmpty`. Map lookups use `val, isFound := userMap[id]` or `val, isUserExist := userMap[id]`.
- **Guard Clauses:** Invert early checks to return immediately and flatten nested blocks.
- **Go Errors & Result Containers:** Return `*appfault.AppError` and replace multi-value error tuples with `appfault.Result[T]`, `appfault.ResultSlice[T]`, `appfault.ResultMap[K, V]`.
- **Mandatory types.go Single Reusable Type Definition:** Centralize all domain payload structs and repeated Result aliases (e.g. `type ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]`) into a dedicated `types.go` file within the package as a single reusable named type. Never leave unexported structs or raw generic Result declarations scattered inline.
- **Pointer Null Safety & Method Composition:** All Result inspection methods must attach to pointer receivers (`*Result[T]`) with line-1 `if r == nil` guards, composing existing methods (`r.IsFailure()`, `r.IsSuccess()`, `r.Count()`) rather than repeating raw pointer/error checks.
- **Fluent Predicates:** Enforce `IsCountOtherThan(N)`, `IsEmpty()`, `HasRecord()`, `IsDefined()` at call sites instead of compound checks (`err != nil || len(...) != N`).
- **No Releases (Strict Policy):** Strictly forbidden from bumping versions or cutting releases at the end of this task.
- **TOTAL BAN on Test Running & Build Checking:** Test execution (`go test`, `pytest`, python runners) and build checks (`go build`, compiler verification) are strictly banned. Testing and builds are verified later on in CI/CD.
- **Atomic Change Tracking:** Append all modified files to `.lovable/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`), mapping to associated tests in `.lovable/test-inventory.json` for subsequent CI/CD verification.
- **Targeted Batch Verification:** Run targeted linter / autofixer on the modified files in the batch (e.g. `python linter-scripts/check-nested-ifs.py <files>` or `08-naming-autofixer.py <files>`). DO NOT run `06-cicd-local-runner.py`, test runners, or full builds (`npm run build`, `go build ./...`) during routine batch fixes.
- **Consolidated Commits:** NEVER commit isolated 1-2 plan/doc files alone. Commit all modified source files, tests, and plans together as a single atomic batch.
- **Immediate Git Push:** Always push immediately to remote (`git push origin <branch>`) after every commit.

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
