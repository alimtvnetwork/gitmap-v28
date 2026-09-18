---
name: parent-task-n-step-loop
description: Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion.
---
# Parent Task N-Step Continuous Loop & Multi-Agent Orchestration

## Phase 1: Planning Mode
1. Verbatim Prompt Capture: Write user prompt into `.ai-memory/plans/pending/xx-<slug>.md`. Extract tasks.
2. Scan & Discover: Spawn 2 planning subagents to scan codebase.
3. Master Spec Generation: Save architectural plan. Add custom rules.
4. Lean Subtask Decomposition: Break down into focused subtasks in `.ai-memory/plans/subtasks/xx-<slug>/`.
5. Strict Relative Paths: Zero absolute paths.
6. Mandatory Auto-Loop: Transition directly to Phase 2 without stopping.

## Phase 2: Execution Mode
1. Parallel Dispatch: Spawn 2 execution subagents (max 2 threads each).
2. Coding Guidelines: Enforce max 15-line functions, single return types, Unix LF.
3. Failure Memory: On crash, rollback and log to `.ai-memory/memory/issues/`.
4. Targeted Linting: Run file-level linters ONLY. NO local test runners or routine builds.
5. Atomic Change Tracking: Record modified files using `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.
6. Remote CI/CD Pipeline Monitoring (GitMap Pipeline-AI): When checking remote CI/CD pipelines, agents MUST use GitMap: `gitmap pipeline-ai status --json` or `gitmap pl-ai status -t <sec>`. Strictly wait based on `etaSeconds` (`-t <sec>`) using adaptive intervals (ETA > 120s: 20s–30s; 60s < ETA <= 120s: 10s–20s; ETA <= 60s: 5s–10s). BANNED: Never loop rapidly or busy-poll (`gh run view` in tight loops).
7. Consolidate & Final Commit: Group all completed subtasks into a single `.ai-memory/plans/completed/` file. Single atomic git commit.

## Banned Operations
- **NO TEST RUNNING:** Never run full test suites or local runners during routine turns.
- **NO BUILD CHECKING:** Never run build commands (`go build`, `npm run build`) during routine turns.
- **NO RAPID CI/CD POLLING:** Never tight-loop on remote status; use GitMap Pipeline-AI with dynamic ETA waiting.
- **NO AUTOMATIC RELEASES:** Never bump versions or cut releases without explicit user command.
