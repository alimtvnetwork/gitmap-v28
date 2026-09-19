---
name: batched-loop-orchestration
description: >-
  Autonomously orchestrate and execute pending tasks using batched loops, execution waves, subtask decomposition, and parallel sub-agents without build or test execution.
---

# [V2] Batched Loop & Execution Wave Orchestration — Workflow (must follow)

> **Prompt Version:** 2.2.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Autonomously orchestrate and execute the pending tasks by decomposing them into subtasks and running a continuous N-step self-loop until completion without a single failure.

## Non-Negotiable Rules (Auto-Reject on Violation)

1. Maximum 3 sub-agents may run concurrently at any time. Never exceed this limit.
2. TOTAL BAN on test running and build checking during routine execution: DO NOT run tests using Python scripts, Go (`go test`), or any test runner. DO NOT check builds (`go build`, compiler checks). All test execution and build verification is deferred to CI/CD.
3. At the end of every loop, output explicit task statistics (done, pending, remaining list).

## Anti-Hallucination & Relative Path Rules

> [!CAUTION]
> **STRICT RELATIVE GIT PATHS ONLY — NO ABSOLUTE PATHS / NO `file:///` URIs:**
> All file paths, markdown links, citations, and task targets MUST be relative paths starting from the repository root (e.g. `02-spec/03-error-manage/01-index.md`, `cmd/main.go`).
> NEVER write drive letters or absolute OS paths (`/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///...`) into ANY file.

## Execution Waves & Disjoint Locking Matrix

- Group pending tasks into Execution Waves.
- Ensure parallel tasks touch completely disjoint files to prevent git merge conflicts.
- Create isolated task directory `.ai-memory/temp-agents/xx-<task-name>/state.md`.
- Functions must be 8–15 lines max; files 100–200 lines max.
- Zero nested ifs (nesting depth > 1 is strictly forbidden).
- Affirmative booleans (`is*`, `has*`, `can*` prefixes; no negative inverts).
- Single return types using `result.Result[T]` or `*apperror.AppError`. No `(T, error)` tuples!

## Task Consolidation & File Reduction (End of Loop)

1. Combine all completed granular subtasks from `.ai-memory/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.ai-memory/plans/completed/xx-<slug>.md`.
2. Include a header explicitly referencing how the main task started and documenting step counts.
3. Delete the original granular `.md` files in `.ai-memory/plans/subtasks/xx-<slug>/`.
4. Delete the original parent plan `.ai-memory/plans/pending/xx-<slug>.md`.
5. Update `.ai-memory/plans/01-index.md` and `.ai-memory/plans/index.md`.
6. Final step: stage all files (`git add -A`), commit atomically, and push to git (`git push origin <branch>`).
