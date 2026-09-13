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
- **Guard Clauses:** Invert early checks to return immediately and flatten nested blocks.
- **Go Errors & Result Containers:** Return `*appfault.AppError` and replace multi-value error tuples with `appfault.Result[T]`, `appfault.ResultSlice[T]`, `appfault.ResultMap[K, V]`.
- **Mandatory types.go Single Reusable Type Definition:** Centralize all domain payload structs and repeated Result aliases (e.g. `type ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]`) into a dedicated `types.go` file within the package as a single reusable named type. Never leave unexported structs or raw generic Result declarations scattered inline.
- **Pointer Null Safety & Method Composition:** All Result inspection methods must attach to pointer receivers (`*Result[T]`) with line-1 `if r == nil` guards, composing existing methods (`r.IsFailure()`, `r.IsSuccess()`, `r.Count()`) rather than repeating raw pointer/error checks.
- **Fluent Predicates:** Enforce `IsCountOtherThan(N)`, `IsEmpty()`, `HasRecord()`, `IsDefined()` at call sites instead of compound checks (`err != nil || len(...) != N`).
- **No Releases:** Strictly forbidden from bumping versions or cutting releases at the end of this task.
- **No Test Execution:** Test execution is disabled unless explicitly commanded by the repository owner.
- **Atomic Change Tracking:** Append all modified files to `.lovable/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`), mapping to associated tests in `.lovable/test-inventory.json`.
- **Targeted Batch Verification:** Run targeted linter / autofixer on the modified files in the batch (e.g. `python linter-scripts/check-nested-ifs.py <files>` or `08-naming-autofixer.py <files>`). DO NOT run `06-cicd-local-runner.py` or full builds (`npm run build`, `go build ./...`) during routine batch fixes.
- **Consolidated Commits:** NEVER commit isolated 1-2 plan/doc files alone. Commit all modified source files, tests, and plans together as a single atomic batch.
- **Immediate Git Push:** Always push immediately to remote (`git push origin <branch>`) after every commit.
- **Atomic Change Cache Recording:** Record all modified files to `.lovable/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`).
