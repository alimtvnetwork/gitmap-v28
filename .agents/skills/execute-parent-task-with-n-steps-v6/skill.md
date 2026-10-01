---
name: execute-parent-task-with-n-steps-v6
description: Autonomously orchestrate and execute complex parent tasks using N-step continuous self-looping (N = 300), mandatory subagents (A = 2, H = 2), GitMap high-speed commands, strict no-build/no-test rules, and atomic commits.
---

# Parent Task N-Step Continuous Loop (V6)

```text
N = 300 (Total self-loop steps budget — default: 300)
A = 2   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
PHASE_1_BUDGET = N / 2   (Steps 1 .. 150: Planning, Parallel Discovery Subagents, Detailed Spec, and Lean Subtask Generation)
PHASE_2_BUDGET = N / 2   (Steps 151 .. 300: Mandatory Parallel Subagent Execution, Self-Looping, Targeted Quality Linting)
WAVES = ceil(subtasks / (A x H))
```

## Mandatory Subagent Spawning Gate (A = 2, H = 2)
- The lead agent MUST explicitly call `invoke_subagent` at each stage:
  1. Planning stage (2 research subagents).
  2. Spec stage (2 spec-authoring subagents).
  3. Execution stage (2 worker subagents per wave).
- Solo execution without calling `invoke_subagent` is an auto-reject failure.
- Zero builds and test suites during routine execution turns (Rule R1).
- Targeted verification only on modified files (Rule R2).
- Strict relative paths and lowercase hygiene (Rule R11).
