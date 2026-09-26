---
name: cg-error-management
description: >-
  Autonomously audits, refactors, and validates repository-wide error management against 02-spec/03-error-manage/ using AppError wrappers, universal response envelopes, and CI linters.
---

# Error Management & Architecture Coding Guidelines (`cg-error-management`)

This skill provides autonomous audit, refactoring, and validation of repository-wide error handling based on `02-spec/03-error-manage/` and `.ai-memory/coding-guidelines.md`.

## Core Invariants

1. **Top-Instruction Priority Mandate (Preamble Precedence):** Any directive, constraint, or checklist placed at the top of an incoming prompt, header alert block, or user request represents an absolute MUST FOLLOW mandate that takes highest priority and strictly overrides generic guidelines below it.
2. **Strict Golang Error Wrapping Mandate (`*appfault.AppError` & `appfault.Fault`):**
   - Whenever ANY Go function encounters, intercepts, or receives an error (from stdlib `os`, `io`, `json`, `sql`, `net`, or downstream packages), it MUST be immediately embedded and wrapped using `appfault.Wrap(errType, err, "opName")` or `result.WrapFailure[T]`.
   - Raw standard library `error` returns (`return err`) in domain packages are strictly prohibited.
   - Functions with side-effects only MUST return `*appfault.AppError`. Functions returning data and possible error MUST return `result.Wrap[T]`.
3. **Zero Swallowed Errors Policy (TOTAL BAN):**
   - NEVER swallow, suppress, or silently ignore errors under any circumstances.
   - NO empty `catch` or `except:` blocks.
   - NO blank identifier error discards (`_ = err` or `val, _ := fn()`).
   - NO returning fallback default values (`return nil`, `return ""`, `return false`) to mask an underlying error without caller notification.
   - Every caught or received error MUST either be completely resolved with structured context logging (operation name, input parameters) OR embedded/wrapped in `*appfault.AppError` and returned to the caller.
4. **No Bare Panics or Bare Exits**: Zero calls to `panic("...")`, `panic(err)`, or `os.Exit(...)` outside the central dispatcher (`cliexit.HandleError`).
5. **Context-Rich Metadata**: All errors MUST use structured `*appfault.AppError` from `04-code/golang/pkg/appfault` and be wrapped with `Op`, `Code`, `Type`, `Severity`, `Creator`, `Message`, `Ctx`, and `Cause`.
6. **Single Result Container Return Types (`pkg/appfault`) & `types.go` Mandate**:
   - Multi-value returns returning errors (`(map[K]V, error)`, `([]T, error)`, `(T, error)`) are strictly banned.
   - Functions returning maps MUST return `appfault.ResultMap[K, V]` (or domain alias).
   - Functions returning slices MUST return `appfault.ResultSlice[T]` (or domain alias).
   - Functions returning single values MUST return `appfault.Result[T]` (or domain alias).
   - **Mandatory `types.go` Single Reusable Type Definition:** All domain payload structs and repeated generic Result envelopes MUST be defined in a dedicated `types.go` file within the package as a single reusable named type everywhere.
7. **Universal Response Envelope**: All API endpoints return `{ "data": ..., "errors": [...], "meta": ... }`.
8. **Targeted Verification**: Continuous verification via `python linter-scripts/check-error-management.py <files>`. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine fixes.
9. **No Releases**: Strictly forbidden from bumping versions or cutting releases at the end of this task.
10. **Atomic Change Tracking**: Append all modified files to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`), mapping to associated tests in `.ai-memory/test-inventory.json`.

---

## Fast File Discovery & Reading Toolchain (GitMap AUM Primary, Python Fallback)

To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, the AI agent MUST use the fast 2-tier discovery toolchain:

### Tier 1: GitMap AUM Acceleration (PRIMARY)
- **Universal File Search:** `gitmap find "<pattern>" [-ext <ext>]` (e.g. `gitmap find "*.go" -ext "go"`, `gitmap find "01*"`)
- **List Indexed Files:** `gitmap list-files [pattern]` (alias `gitmap lf [pattern] [-ext <ext>]`)
- **Substring Match:** `gitmap find-files-any "<substring>"` (alias `gitmap ffa "<str>"`)
- **Stream File Content:** `gitmap cat <filepath>` (streams to stdout with zero disk writes)
- **Instant Code Search:** `gitmap search "<term>"` (immediate multi-core filesystem walk)

### Tier 2: Fast Cached Python Toolchain (FALLBACK)
- **Inventory Target Files:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats`
- **Fast Cached Grep (<15ms):** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<search-pattern>" --limit 50`
- **Sub-Millisecond Folder & File Exploration:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder <folder-path> --limit 50`
- **Read Target File:** `python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000`
- **Subsystem & Topology Overview:** `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`

---

## Routine Execution Policy

- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns or micro-batch loops. Running the heavy 28-38 gate pipeline across the entire repository wastes massive amounts of time. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them. Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
