# Completed Plan: 234-codebase-review-remediation-and-consolidation

## Metadata
- **Plan Slug:** `234-codebase-review-remediation-and-consolidation`
- **Master Ledger:** [02-spec/21-app/234-codebase-review-remediation-and-consolidation/00-master-audit-ledger.md](../../02-spec/21-app/234-codebase-review-remediation-and-consolidation/00-master-audit-ledger.md)
- **Architecture Spec:** [02-spec/21-app/234-codebase-review-remediation-and-consolidation/01-architecture-spec.md](../../02-spec/21-app/234-codebase-review-remediation-and-consolidation/01-architecture-spec.md)
- **Component Spec:** [02-spec/21-app/234-codebase-review-remediation-and-consolidation/02-component-and-cli-spec.md](../../02-spec/21-app/234-codebase-review-remediation-and-consolidation/02-component-and-cli-spec.md)
- **Honest Counter-Review:** [02-spec/21-app/234-codebase-review-remediation-and-consolidation/03-honest-counter-review.md](../../02-spec/21-app/234-codebase-review-remediation-and-consolidation/03-honest-counter-review.md)
- **Status:** `COMPLETED`
- **Completed Date:** `2026-10-06`
- **Target Version:** `v6.497.0` (Preserved — zero unauthorized version churn per user directive)

---

## User Request (Verbatim)
```text
docs/review

# High Priority Instruction

There are some honest reviews written by News Spark. I do agree with some of the reviews, but not all of them. First, read that and put out all the actionable items it asks us to do. Some parts I agree with, for example, overlapping clone, we can map together or keep together. Keeping everything together in one file, I do not agree. That we need to avoid. The next part is banning negation is bad, but I don't want to negate the switch we wanted to have, and also we do not remove this guideline. Let's say eight line to 15 lines. So we don't do that. Remember to focus on that, and then version numbering. Every change on version number, it's not every change actually. Depends on the big feature. So that is a stupidity that it has put out. Duplication is the highest cost problem. I do agree. If we have duplicate code, we should reduce that. Ad hoc exit, we should remove, or we should use the app error plus the CLI exit everywhere. There should be nowhere the D drive or absolute path. Remember to remove from the review as well. Consolidate the README. I agree with this. Fix the version identity from the version.json. I agree with that. Remove the committed artifact from the repo secret. I don't know what kind of artifacts have been committed. You can ask me so that we can confirm it. Do not remove anything unless we confirm from the repo secrets. Shrink the README to an index. I appreciate that. We will follow this. Action checklist: Consolidate the README. Yes, checked. Fix the version.json slag. Do that. But also the coding guideline needs to be there if it is the coding guideline it is talking about. CLI contains exe. Then we have to remove it and also remove from all the git history as well. Audit the repo secrets if there is anything. Repo secrets will contain secrets. No worries on this. Verify security token part. Yes, that you should do from this code base. Also, delete the cluster files, fix update runner, etc. Benchmark we need to have. That needs to be in the docs file as a reference so that anyone can see it, and also that would be referred from the root README. Now, can you please work on these review feedbacks and then finally remove it, but also at the same time, finally you do your honest feedback on this.

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Map or keep overlapping clones together, avoid keeping everything in one file.
5. Do not negate the switch we wanted to have; do not remove this guideline.
6. Avoid unnecessary version number changes; only change for significant features.
7. Reduce duplicate code.
8. Remove ad hoc exits; use app error plus CLI exit.
9. Remove absolute paths like D drive from the review.
10. Consolidate the README.
11. Fix version identity from version.json.
12. Confirm before removing committed artifacts from repo secrets.
13. Shrink README to an index.
14. Ensure coding guidelines are present if applicable.
15. Remove CLI exe and its history from git.
16. Audit repo secrets.
17. Verify security token part.
18. Delete cluster files, fix update runner, etc.
19. Include benchmark in docs file and reference from root README.

## Must follow and spawn agent using

@[.agents/skills/execute-parent-task-with-n-steps-v6]

## Additional Instructions

- /plan first before doing the work to reduce the credits.
- /learn from @[.agents/skills/gitmap] skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
```

---

## Deliverables Summary

### 1. Root Clutter Purge & Binary Hygiene (Subtask 01)
- Removed 22 tracked clutter files from repository root:
  * 16 scratch scripts: `fix_agy_conv_ls.py`, `fix_agy_ls.py`, `fix_agy_most_conv.py`, `fix_cd_alias.py`, `fix_cg_path.py`, `fix_common_sync.py`, `fix_cr_alias.py`, `fix_cr_rootcore.py`, `fix_go_compile.py`, `fix_ip_aliases.py`, `fix_ip_conflict.py`, `fix_most_conv.py`, `fix_repo_create_init.py`, `fix_repo_create_init2.py`, `fix_watch.py`, `fix_words_flag.py`
  * 2 runner scripts: `update_runner.py`, `update_runner2.py`
  * 3 audit CSV spreadsheets: `audit_results_extra.csv`, `audit_results_go.csv`, `audit_results_ts.csv`
  * 1 prompt session dump: `user_prompt_full.txt`
  * 1 duplicate non-dot `gitignore`
- Updated `.gitignore` adding `*.syso` and `*.syso~`.
- Untracked `cli/rsrc_windows_386.syso` and `cli/rsrc_windows_amd64.syso` via `git rm --cached`. Verified zero `*.exe` ever committed in git history.

### 2. README Consolidation & Benchmark Elevation (Subtask 02)
- Consolidated byte-identical 3,503-line `readme.md` / `what-to-read.md` pair:
  * Root `readme.md` transformed into a clean, high-speed 155-line index linking to CLI commands, installation scripts, benchmarks, specs, and AI memory.
  * `what-to-read.md` reduced to a 32-line pointer document referencing `readme.md` and `.ai-memory/what-to-read.md`.
- Elevated root `benchmark.md` into `docs/benchmarks/benchmark.md` with verified sub-millisecond search numbers (40 µs hot, 0.82 ms cold) and fixed script pointer (`03-ai-scripts/40-run-search-benchmarks.py`).

### 3. Version Identity Alignment & Clone Architecture Spec (Subtask 03)
- Updated `version.json` identity to GitMap:
  * `Title`: `"GitMap"`
  * `RepoSlug`: `"gitmap-v28"`
  * `RepoUrl`: `"https://github.com/alimtvnetwork/gitmap-v28"`
  * `name`: `"gitmap"`
  * `Description`: Clarified GitMap CLI capabilities while preserving author attribution (MD ALIM UL KARIM / RISEUP ASIA LLC) and coding guidelines metadata.
- Preserved version at `6.497.0` (zero unnecessary version churn).
- Authored `docs/commands/cloning-architecture.md` comprehensively mapping the 5 specialized clone engines (`cloner`, `clonefrom`, `clonenow`, `clonepick`, `clonenext`) and `cloneconcurrency` while explicitly preserving modular separate-package architecture.

### 4. Error Exit Standardization & Security Token Verification (Subtask 04)
- Standardized CLI exits in `cli/cliexit/handle.go` with typed helpers: `HandleSuccess()`, `HandleUsageError()`, `HandleValidationError()`, `HandleGeneralError()`, `HandleNotFound()`, `HandleAppError()`.
- Refactored 27 ad-hoc `cliexit.HandleError(nil, 0)` / `(nil, 2)` calls across 23 files in `cli/cmd/` to typed helpers.
- Added comprehensive unit tests in `cli/cliexit/handle_test.go` (100% pass).
- Automated security token scan across all files: verified zero exposed OAuth refresh tokens, PATs, or private keys.

### 5. repo-secrets Non-Destructive Audit & Honest Counter-Review (Subtask 05)
- Completed full audit of `repo-secrets/`: verified all 21 files (4 documentation runbooks, 5 sanitized configuration manifests, 12 migration scripts) with zero live secrets. Preserved all files without any deletions per user mandate.
- Cleanly deleted temporary review directory `docs/review/` (`readme.md`, `honest-feedback.md`, `action-checklist.md`).
- Authored [02-spec/21-app/234-codebase-review-remediation-and-consolidation/03-honest-counter-review.md](../../02-spec/21-app/234-codebase-review-remediation-and-consolidation/03-honest-counter-review.md) presenting our objective engineering analysis.

---

## Verification Evidence
- `go vet ./cliexit ./cmd`: PASS (exit code 0).
- `go test -v ./cliexit`: PASS (100% pass rate).
- `python3 linter-scripts/check-forbidden-strings.py`: PASS (0 violations).
- `python3 linter-scripts/check-relative-paths.py`: PASS (0 violations across 8,871 files).
- `git status`: Verified clean working tree, clutter files deleted, `.syso` untracked, zero `.exe` files.
