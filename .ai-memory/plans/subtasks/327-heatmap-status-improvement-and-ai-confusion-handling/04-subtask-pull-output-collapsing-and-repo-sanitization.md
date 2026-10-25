# Subtask 04: Pull Output Collapsing and Repository Sanitization

- **Parent Task:** 327 (`heatmap-status-improvement-and-ai-confusion-handling`)
- **Subtask Sequence:** 04
- **Subtask Path:** `.ai-memory/plans/subtasks/327-heatmap-status-improvement-and-ai-confusion-handling/04-subtask-pull-output-collapsing-and-repo-sanitization.md`
- **Assigned Subagent:** Subagent 3
- **Created At:** 2026-10-25
- **Status:** `IN_PROGRESS`

---

## 1. Context & Objectives

During recent GitMap pull executions, commit message strings (e.g. `COMMIT - support cfr verb in commit formatter and eliminate magic strings across cli`) were treated as failed repositories (`D:\work\COMMIT - support cfr...`) and logged to the split pull error database. Furthermore, pull operations with dozens of unchanged repositories flooded the terminal output, causing clutter.

This subtask addresses these two problems:
1. **Repository Name Candidate Sanitization:**
   - Filter invalid repository names across `cli/cmdpull/pull_remediation_query.go`, `cli/cloner/pulldiag.go`, and `cli/cmdpull/pull_batch.go`.
   - Prevent ghost repositories (missing directory + missing remote URL) from being queued for remediation.
2. **Pull Output Collapsing & Guidance:**
   - Confirm and enforce collapsing of up-to-date repositories in batch/efficient pull outputs.
   - Provide a helpful one-liner hint (`gitmap pull status`) for granular review.

---

## 2. Detailed Technical Plan

### 2.1 Candidate Validation (`isValidRepoNameCandidate`)
Implement `isValidRepoNameCandidate(name string) bool`:
- Trim whitespace.
- Reject if length is 0 or greater than 64 characters.
- Reject if string contains newlines (`\n` or `\r`) or `" - "`.
- Reject if string has prefixes: `COMMIT `, `FEAT `, `FIX `, `Feature: `, `Bug: `, `Release: `, `Chore: `.

### 2.2 File Modifications
1. **`cli/cmdpull/pull_remediation_query.go`**:
   - Add `isValidRepoNameCandidate` and sub-helpers adhering to the 15-line limit and positive boolean conventions.
   - Refactor `convertRecordsToFailures` to skip any record where `!isValidRepoNameCandidate(r.RepoSlug)` or where the directory does not exist and has no remote URL recorded.
   - Guard `ExtractPullFailures` against invalid candidate repository names.
2. **`cli/cloner/pulldiag.go`**:
   - Add `IsValidRepoNameCandidate` and internal helpers.
   - In `ClassifyNonRepoFolder`: validate candidate before inspecting; if candidate fails validation, return empty `NonRepoDiagnosis{}` (`IsNonRepoFolder == false`).
   - Decompose `ClassifyNonRepoFolder` into sub-functions <= 15 lines.
3. **`cli/cmdpull/pull_batch.go`**:
   - Add `AppendPullErrorIfValid` and `AppendFailureIfValid`.
   - Guard `handlePullRemediationForRecords`, `handleBatchFailuresRemediation`, and `finalizePullBatchStandard` to verify `isValidRepoNameCandidate` before enqueueing or prompting.

---

## 3. Strict Coding Standards Compliance

- **Function Line Limit:** All functions <= 15 lines (<= 8 lines preferred).
- **Conditionals:** Zero nested `if` statements; early returns and guard clauses only.
- **Boolean Naming:** Positive prefixes only (`is*` / `has*`); no negative booleans.
- **Formatting:** Blank lines before `if`, after `}`, and before `return`.
- **Error Handling:** Return `error` or `*apperror.AppError`; no leaf-level exits.
- **Build Discipline:** Strict NO BUILD and NO TEST during editing.

---

## 4. Verification Criteria

- [ ] `isValidRepoNameCandidate("COMMIT - foo")` returns `false`.
- [ ] `isValidRepoNameCandidate("gitmap-v28")` returns `true`.
- [ ] `convertRecordsToFailures` filters out phantom records with non-existent paths and empty remotes.
- [ ] `ClassifyNonRepoFolder` ignores strings with commit message patterns.
- [ ] Pull batch failure handling cleanly ignores invalid repository slugs.
