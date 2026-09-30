---
name: execute-parent-task-with-n-steps-v3
description: "Executes the [V3] Parent Task N-Step Continuous Loop & Mandatory Multi-Agent Subagent Orchestration prompt (N = 300 top header, A = 2, H = 2 mandatory invoke_subagent gate). Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion without a single failure with strict no-build and no-test execution."
---

# [V3] Parent Task N-Step Continuous Loop & Mandatory Multi-Agent Subagent Orchestration — Workflow (must follow)

```text
N = 300 (Total self-loop steps budget — editable top-header parameter, default: 300)
A = 2   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
PHASE_1_STEPS = N / 2   (Steps 1 .. 150: Planning, Parallel Discovery Subagents, Detailed Spec, and Lean Subtask Generation)
PHASE_2_STEPS = N / 2   (Steps 151 .. 300: Mandatory Parallel Subagent Execution, Self-Looping, Targeted Quality Linting)
```

Source prompt: `01-prompts/14-execute/10-execute-parent-task-with-n-steps-v3.md`

## Instructions

1. Read `01-prompts/14-execute/10-execute-parent-task-with-n-steps-v3.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. Spawning subagents via `invoke_subagent` (`A = 2`, `H = 2`, `TypeName: "self"`) in Phase 1 and Phase 2 is strictly mandatory (zero solo execution allowed).
4. A direct instruction in the current user message overrides the prompt when they conflict.
5. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous 300-step (`N = 300`) self-loop with mandatory `invoke_subagent` (`A = 2, H = 2`) orchestration until completion without a single failure with strict no-build and no-test execution.
