# 002: Interactive Macro Builder Verification, Git Hygiene & Atomic Release Sync

## 1. Objectives & Context

Execute final review, git hygiene, repository status verification, and single-commit synchronization for the Interactive Macro Builder features delivered under Plan 83 ([.lovable/plans/completed/83-interactive-macro-builder-pwd-ls-search.md](../.lovable/plans/completed/83-interactive-macro-builder-pwd-ls-search.md)).

### Mandatory Reference Documents
- `mem://01-index.md` / `.lovable/memory/01-index.md`: Core memory, CODE RED constraints, and zero-error-swallowing policy.
- `.lovable/coding-guidelines.md`: Strict function length (≤15 lines), blank line before return, and positive boolean naming.
- `.lovable/plans/01-index.md`: Active plans roadmap and completed milestone index.
- `.lovable/strictly-avoid.md`: Prohibited patterns, CI/CD bypass bans, and absolute path prohibitions.

---

## 2. Scope & Acceptance Criteria

### A. Code & Test Integrity Verification
- [ ] Re-verify unit tests for interactive macro builder:
  - `go test -v ./uipref` (all tests passing).
  - `go test -v ./cmd -run TestMacro` (all tests passing).
  - `go test -v ./cmd -run TestProcessInteractiveStepLine` (all tests passing).
- [ ] Verify repository quality linters:
  - `python linter-scripts/check-nested-ifs.py` (0 violations).
  - `python linter-scripts/check-boolean-guidelines.py` (0 violations).
  - `python linter-scripts/check-relative-paths.py` (0 violations).
  - `python linter-scripts/check-error-management.py` (0 violations).

### B. Git Hygiene & Artifact Exclusion
- [ ] Confirm `.gitignore` explicitly ignores test reports, coverage output, `.test-report.*`, and temporary test databases.
- [ ] Verify `git status` contains zero untracked test output artifacts, temporary data dumps, or binary executables.
- [ ] Confirm root `readme.md` is strictly lowercase.

### C. Atomic Single Commit & Remote Synchronization
- [ ] Stage all modified and added files adhering to conventional commit format.
- [ ] Commit with message: `feat(macro): add PWD header, in-builder ls/find/search/replace helpers, and uipref toggle`
- [ ] Complete work strictly on the current branch.
- [ ] Push the commit to the remote repository.

### D. File Change Summary Deliverable
- [ ] Deliver a comprehensive, file-by-file summary in the response listing:
  - File path.
  - Specific changes implemented.
  - Rationale and user requirement alignment.

---

## 3. Definition of Done
1. All Go unit tests pass cleanly.
2. All 4 quality linters report zero violations.
3. Git working tree is clean with all changes committed in a single logical commit on the current branch.
4. Remote branch is up-to-date with local branch.
5. Detailed file-by-file summary delivered in chat.

---

## 4. Self-Instruction & Post-Execution Feedback
- **Before acting:** Re-read `mem://01-index.md` and `.lovable/coding-guidelines.md`; restate which rules apply before executing.
- **After execution:** Suggest any further refinements or automated assertions that could be added to this instruction for future runs.
