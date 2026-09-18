---
name: cg-error-management
description: >-
  Autonomously audits, refactors, and validates repository-wide error management against 02-spec/03-error-manage/ using AppError wrappers, universal response envelopes, and CI linters.
---

# Error Management & Architecture Coding Guidelines (`cg-error-management`)

This skill provides autonomous audit, refactoring, and validation of repository-wide error handling based on `02-spec/03-error-manage/` and `.ai-memory/coding-guidelines.md`.

## Core Invariants

1. **No Bare Panics or Bare Exits**: Zero calls to `panic("...")`, `panic(err)`, or `os.Exit(...)` outside the central dispatcher (`cliexit.HandleError`).
2. **Context-Rich `*appfault.AppError` Wrappers**: All errors MUST use structured `*appfault.AppError` from `04-code/golang/pkg/appfault` and be wrapped with `Op`, `Code`, `Type`, `Severity`, `Creator`, `Message`, `Ctx`, and `Cause`. Raw standard library `error` returns are strictly prohibited.
3. **Single Result Container Return Types (`pkg/appfault`) & `types.go` Mandate**:
   - Multi-value returns returning errors (`(map[K]V, error)`, `([]T, error)`, `(T, error)`) are strictly banned.
   - Functions returning maps MUST return `appfault.ResultMap[K, V]` (or domain alias).
   - Functions returning slices MUST return `appfault.ResultSlice[T]` (or domain alias).
   - Functions returning single values MUST return `appfault.Result[T]` (or domain alias).
   - Functions with side-effects only MUST return `*appfault.AppError`.
   - **Mandatory `types.go` Single Reusable Type Definition:** All domain payload structs (e.g. `ScheduleExportBundle`) and repeated generic Result envelopes (`type ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]`) MUST be defined in a dedicated `types.go` file within the package as a single reusable named type everywhere. Never declare ad-hoc unexported structs or raw generic Result envelopes inline in implementation files.
   - Result struct fields MUST use affirmative prefixes (e.g. `isDefined bool`, TOTAL BAN on bare `defined bool`).
4. **Pointer-Attached Null Safety & Method Composition**:
   - All Result inspection methods MUST be attached to pointer receivers (`(r *Result[T])`, `(rs *ResultSlice[T])`, `(rm *ResultMap[K, V])`).
   - Total ban on value receivers on Result checking methods to eliminate nil pointer dereference panics.
   - **Method Composition Mandate:** Inspection methods MUST delegate to and compose existing methods (`r.IsFailure()`, `r.IsSuccess()`, `r.Count()`) rather than repeating raw pointer/error checks (`r == nil || r.err != nil`).
   - Line-1 `if r == nil` guards MUST return safe canonical defaults without crashing:
     - `IsFailure()` -> `true`
     - `IsSuccess()` -> `false`
     - `Count()` -> `0`
     - `IsEmpty()` -> `true`
     - `HasRecord()` / `HasRecords()` -> `false`
     - `IsDefined()` -> `false`
     - `IsCountOtherThan(n)` -> `true`
     - `AppError()` / `Fault()` -> `nil`
5. **The 4 Core Predicate Methods**:
   - `res.IsCountOtherThan(number int) bool`: Returns `true` if operation failed (or nil receiver) OR `Count() != number`. Replaces compound checks like `err != nil || len(...) != N` or `IsFailure() || Count() != N`.
   - `res.IsEmpty() bool`: Returns `true` if collection has 0 elements, payload data is empty/null/zero, or receiver is nil.
   - `res.HasRecord() bool` (and alias `res.HasRecords() bool`): Returns `true` if operation succeeded (no error) AND has **more than 0 records** (`Count() > 0 && !IsFailure()`).
   - `res.IsDefined() bool`: Returns `true` if operation succeeded (no error) AND `recordCount > 0` (or non-null/non-empty data `T`). Delegates error validation to `IsSuccess()`/`IsFailure()`.
6. **Universal Response Envelope**: All API endpoints return `{ "data": ..., "errors": [...], "meta": ... }`.
7. **Never Swallow Errors**: Every catch block and error return must be recorded and handled explicitly.
8. **Targeted Verification**: Continuous verification via `python linter-scripts/check-error-management.py <files>`. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine fixes.
9. **No Releases**: Strictly forbidden from bumping versions or cutting releases at the end of this task.
10. **Atomic Change Tracking**: Append all modified files to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`), mapping to associated tests in `.ai-memory/test-inventory.json`.

## Fast File Discovery & Reading via Python Toolchain (Mandatory Acceleration)

To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, the AI agent MUST use the repository's dedicated Python discovery scripts first:
- **Inventory Target Files:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats`
- **Fast Cached Grep (<15ms):** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<pattern>" --lang go --limit 50`
- **Sub-Millisecond Folder Explorer & Reader:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder <dir> --ext .go --limit 50`
- **Read Target File:** `python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000`
- **Fast Pattern Search:** `python 03-ai-scripts/17-fast-file-reader.py --search-pattern "<pattern>" --limit 50`
- **Codebase Topology:** `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`

---

## Routine Execution Policy

- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns or micro-batch loops. Running the heavy 28-38 gate pipeline across the entire repository wastes massive amounts of time. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.

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
