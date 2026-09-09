---
name: execute-parent-task
description: Autonomously orchestrate and execute parent tasks by decomposing them into subtasks and running continuous N-step self-loops with strict coding guidelines and error management.
---

# Parent Task N-Step Continuous Loop & Multi-Agent Orchestration — Workflow

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

## Master Task Checklist

1. Phase 1 (Planning & Spec Generation, Steps 1..N/2): Spawn exactly 2 planning subagents (max 2 threads each) to scan the codebase and draft `.lovable/plans/pending/XX-<slug>.md`.
2. Phase 1 (Subtask Decomposition): Decompose the master plan into microscopic, actionable subtasks in `.lovable/plans/subtasks/XX-<slug>/*.md`.
3. Phase 1 (Strict Folder Bounding): Restrict all planning logs, active locks, and status reports strictly within `.lovable/` (`.lovable/plans/`, `.lovable/01-index.md`).
4. Phase 1 (Zero-Stop Transition): Immediately upon completing Phase 1, self-loop and transition directly into Phase 2 execution mode without pausing or stopping.
5. Phase 2 (Execution & Code Refactoring, Steps N/2+1..N): Spawn exactly 2 execution subagents (max 2 threads each) to execute subtasks on disjoint files in parallel.
6. Phase 2 (Failure Memory & Error Recovery): If a subagent fails, record the failure log in `.lovable/plan.md` and `.lovable/memory/issues/`; subsequent agents MUST read the failure log first to remediate root causes.
7. Phase 2 (Quality Gate Verification): Execute local linters and `python 03-ai-scripts/06-cicd-local-runner.py` ensuring `exit 0` before finishing.
8. Ingest `.lovable/memory/01-index.md` for project memory index and past learnings.
9. Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
10. Ingest `02-spec/02-coding-guidelines/` for domain-specific architectural specifications.
11. Ingest `02-spec/03-error-manage/` for error handling architectures and AppError.
12. Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.

## 1. 2-Agent Concurrency & Ruthless Orchestration
- Strict 2-Agent Limit (Max 2 Threads Each)
- Strict Folder Bounding (`.lovable/`)
- Context Diet: Minimal instruction, subagents read files themselves
- Fail Fast & Kill Stalls

## 2. Phase 1: Planning Mode & Subtask Generation FIRST (Steps 1 .. N/2)
- Scan & Discover: Spawn 2 planning subagents
- Master Spec in `.lovable/plans/pending/XX-<slug>.md`
- Subtask Decomposition in `.lovable/plans/subtasks/XX-<slug>/*.md`
- Relative Git Paths only
- Mandatory Auto-Loop: Immediate transition to Phase 2

## 3. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)
- Parallel Dispatch: 2 execution subagents on disjoint subtasks
- File Locking via `.lovable/01-index.md`
- Quality Gate Verification: Local linters & `06-cicd-local-runner.py`
