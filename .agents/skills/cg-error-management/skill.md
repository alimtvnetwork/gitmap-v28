---
name: cg-error-management
description: >-
  Autonomously audits, refactors, and validates repository-wide error management against spec/03-error-manage/ using AppError wrappers, universal response envelopes, and CI linters.
---

# Error Management & Architecture Coding Guidelines (`cg-error-management`)

This skill provides autonomous audit, refactoring, and validation of repository-wide error handling based on `spec/03-error-manage/` and `.lovable/coding-guidelines/coding-guidelines.md`.

## Core Invariants

1. **No Bare Panics or Bare Exits**: Zero calls to `panic("...")`, `panic(err)`, or `os.Exit(...)` outside the central dispatcher (`cliexit.HandleError`).
2. **Context-Rich `AppError` Wrappers**: All errors MUST be wrapped with `Op`, `Code`, `Type`, `Severity`, `Creator`, `Message`, `Ctx`, and `Cause`.
3. **Universal Response Envelope**: All API endpoints return `{ "data": ..., "errors": [...], "meta": ... }`.
4. **Never Swallow Errors**: Every catch block and error return must be recorded and handled explicitly.
5. **Targeted Verification**: Continuous verification via `python linter-scripts/check-error-management.py <files>`. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine fixes.


## Change Tracking & Test Avoidance
- Test execution disabled; pass `--no-tests` to local runner.
- Append modified files to `.lovable/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`).


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
