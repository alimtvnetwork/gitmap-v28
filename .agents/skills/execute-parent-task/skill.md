---
name: execute-parent-task
description: >-
  Autonomously orchestrate and execute parent tasks by decomposing them into subtasks and running continuous N-step self-loops with strict coding guidelines and error management.
---

# Instruction (must follow): Execute Parent Task (N-Step Continuous Loop & Multi-Agent)

/goal Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion without a single failure.

- First N/2 steps will be given for spec writing and breaking down into subtasks.
- Second N/2 steps will be given to execute created tasks following coding guidelines and error management.
- Ingest memory, coding guidelines, error management, and strictly avoid rules before taking action.
- Strictly adhere to centralized error handling with AppError wrappers, never bare panics or bare exits.
