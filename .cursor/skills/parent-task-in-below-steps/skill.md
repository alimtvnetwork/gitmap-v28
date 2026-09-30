---
name: parent-task-in-below-steps
description: "Executes the Subtask [01]: [Descriptive Subtask Name] prompt. Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous self-loop with strict no-build and no-test execution. Use when the user asks to run parent-task-in-below-steps, or the task is about executing pending tasks, a parent task, or a batched loop."
---

# Subtask [01]: [Descriptive Subtask Name]

Source prompt: `01-prompts/14-execute/09-parent-task-in-below-steps.md`

## Instructions

1. Read `01-prompts/14-execute/09-parent-task-in-below-steps.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous self-loop with strict no-build and no-test execution.
