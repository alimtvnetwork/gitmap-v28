# Guideline Fix Execution & Linter Remediation — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Apply concrete, surgical fixes for all coding guideline violations listed in the pending tasks. Strictly adhere to all style rules, boolean principles, function size limits (< 8 lines), and type-safety standards without introducing regressions.

/learn Ingest and internalize all coding standards, boolean extraction patterns, and error management rules before applying fixes.

Context & References:

- /learn master cross-language coding guidelines: `02-spec/02-coding-guidelines/01-cross-language/15-master-coding-guidelines/`
- /learn braces, nesting, and conditions: `02-spec/02-coding-guidelines/01-cross-language/01-index.md` & `02-conditions-and-extraction.md`
- /learn function size caps (< 8 lines): `02-spec/02-coding-guidelines/01-cross-language/01-index.md`
- /learn boolean naming & inverse rules: `02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md` & `12-no-negatives.md`
- /learn error management & logging: `02-spec/03-error-manage/01-index.md`
- /learn language-specific rules in `02-spec/02-coding-guidelines/` (Go wrapped booleans: `02-spec/02-coding-guidelines/03-golang/09-wrapped-boolean-results.md`)

/goal - Apply concrete fixes for all violations listed in the pending tasks.

- Enforce non‑negotiable rules:
  - DRY – eliminate duplicated logic.
  - Use typed enums instead of magic strings or numbers.
  - Functions must be strictly < 8 lines. (NON-NEGOTIABLE)
  - Source files must be ≤ 80 lines.
  - No code mutation – only apply fixes, never introduce new bugs.
  - Positive boolean naming (`is` / `has`). No `isNot`, `hasNo`, `isUndefined`, or negative prefixes. Always use affirmative positive framing: try `isDefined` / `IsDefined` instead of negatives (e.g., use `isDefined` instead of `isUndefined` or `isNotDefined`, and invert with `!isDefined` in guard clauses; use `isValid` instead of `isNotValid`, `hasValue` instead of `hasNoValue`, `isReady` instead of `isNotReady`). No nested if statements, no magic values.
  - Style: Ensure a blank line before every `return` statement.
  - Mandatory File Path & Variable Context: Any error created or returned when operating on files, paths, or variables must embed the target path via `.WithPath(path)` / `WrapWithPath(..., path)` and variable context via `.WithVar(name, value)` (Rule R7).
  - Golang Single Return & Wrapped Booleans: Strictly return a single parameter (bundle multiple returns into a struct). No raw booleans returned in Go. Return a single Result struct (bundling Data, AppError, and Status together) with two flags (`IsSuccess` and `IsFailed`) managed by a constructor (`NewSuccess`/`NewFailure`).
  - Example usage (Note the explicit variable name `paymentStatus`, no short names like `res`):
    ```go
    paymentStatus := ProcessPayment(100)
    if paymentStatus.IsFailed {
        // handle error
    } else if paymentStatus.IsSuccess {
        // handle success
    }
    ```

#

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

## Checklist (execute phase)

0. Automated Pre-Pass: Run `python 03-ai-scripts/05-guideline-autofixer.py <target-dir>` to automatically clean up return new lines (R13-R16) and remove explicit `== true` checks.
1. Read pending tasks from the `.lovable/plans/subtasks/01-coding-guideline-fixes/` folder.
2. For each task:
   - Locate the affected source file.
   - Verify the file size ≤ 80 lines; if > 80, split into logical modules.
   - Refactor duplicated code into a shared helper function.
   - Replace magic strings/numbers with a newly defined enum in a dedicated `enums.dart` (or appropriate language file).
   - Ensure the refactored function body is ≤ 8 lines; extract sub‑logic to private helpers if needed.
   - Run the project's test suite and the Go race detector (`go test -race ./...`).
   - If tests pass, stage the changes.
3. Commit each fix using the CI/CD fix workflow:
   - Run `git add <modified files>`.
   - Commit with message `fix(coding-guidelines): resolve <issue‑id>`.
   - Push commits to remote branch (`git push origin main`). No automatic version bumps unless explicitly commanded by user.
4. Update the pending task file to mark it as completed.
5. If no pending tasks remain, output a summary of all fixes applied.

### Non‑Hallucination Policy

- Do not assume the existence of a file or enum that is not present; if uncertain, raise a question to the user.
- If a fix would require a large architectural change beyond the scope, create a new pending task instead of applying it directly.

### Execution Loop

- Process up to 50 tasks per run to avoid long‑running blocks.
- After each batch, report progress and await user confirmation before proceeding to the next batch.

---

/goal Apply coding‑guideline fixes safely and push a minor release.
/learn Ensure future prompts respect the same checklist and constraints.

## Actionable Items & Checklist

- [ ] Read and adhere to: `.lovable/coding-guidelines.md`
- [ ] Read and adhere to: `02-spec/02-coding-guidelines/01-cross-language/01-index.md`
- [ ] Read and adhere to: `02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md`
- [ ] Read and adhere to: `02-spec/02-coding-guidelines/03-golang/09-wrapped-boolean-results.md`
- [ ] Read and adhere to: `02-spec/02-coding-guidelines/06-ai-optimization/01-index.md`
- [ ] Read and adhere to: `02-spec/03-error-manage/01-index.md`
- [ ] Read the overarching main task plan.
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] **Boolean Standards & IsDefined (Rule R3):** Positive prefixes only (`is` and `has`). MANDATORY: Use `isDefined` (or `res.IsDefined()`) instead of inverted empty checks (`!isEmpty` / `!res.IsEmpty()`). Never write `if !isEmpty`. Map lookups use `val, isFound := userMap[id]` or `val, isUserExist := userMap[id]`.
- [ ] **GitHub Actions Zero Storage (Rule R18):** Never upload build binaries, logs, test artifacts, or reports in CI workflows (`actions/upload-artifact` is strictly banned in CI). Free tier accounts have an account-wide cap of 0.5 GB (500 MB). Releases belong exclusively in GitHub Releases (`release.yml`), never in Actions artifact storage.
- [ ] Group all completed work into a single logical commit.
- [ ] Anti-Hallucination Check: Before finalizing the plan or writing code, you MUST use read/search tools to verify that every file, function, enum, and variable mentioned actually exists in the codebase. Do not guess.
- [ ] Push the commit to the remote repository.
- [ ] Self-loop continuously until all the code issues are listed out in tasks and pending tasks.
- [ ] Describe all issues and files that need to be tested against for each file.
- [ ] Make a detailed plan/task for each file.
- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

---

slug: execute-coding-guideline-fix
status: active
---

## The 4-Part RCA Requirement (Mandatory Memory File)

Before you write any code to fix the problem, you MUST document the issue in `.lovable/memory/issues/xx-<slug>.md` (where XX is the next available sequential number). The file MUST contain these exact four sections:

1. **Why it happened:** The high-level business, logical, or architectural breakdown of the failure.
2. **How it happened:** The technical execution flow that triggered the bug.
3. **Root Cause:** The exact file, line, and dependency responsible for the failure.
4. **Code Fix:** The exact code snippets showing what needed to be changed to fix the root cause.
