---
name: execute-pending-tasks
description: "Executes the Execute Pending Tasks prompt. Autonomously orchestrate and execute ALL pending tasks in a continuous N-step self-loop until the entire queue is completely resolved without a single failure with strict no-build and no-test execution. Use when the user asks to run execute-pending-tasks, or the task is about executing pending tasks, a parent task, or a batched loop."
---

# Execute Pending Tasks

Source prompt: `01-prompts/14-execute/01-execute-pending-tasks.md`

## Instructions

1. Read `01-prompts/14-execute/01-execute-pending-tasks.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously orchestrate and execute ALL pending tasks in a continuous N-step self-loop until the entire queue is completely resolved without a single failure with strict no-build and no-test execution.
