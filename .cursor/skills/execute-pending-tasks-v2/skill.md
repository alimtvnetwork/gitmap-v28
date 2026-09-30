---
name: execute-pending-tasks-v2
description: "Executes the Subtask Execution Engine (v2) prompt. Autonomously orchestrate and execute ALL pending tasks in a continuous N-step self-loop until the entire queue is completely resolved without a single failure. Use when the user asks to run execute-pending-tasks-v2, or the task is about commits, boolean cleanups, or multi-agent code fixes."
---

# Subtask Execution Engine (v2)

Source prompt: `01-prompts/09-commit-and-multi-agent-code-fix/04-execute-pending-tasks-v2.md`

## Instructions

1. Read `01-prompts/09-commit-and-multi-agent-code-fix/04-execute-pending-tasks-v2.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously orchestrate and execute ALL pending tasks in a continuous N-step self-loop until the entire queue is completely resolved without a single failure.
