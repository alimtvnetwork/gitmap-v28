---
name: commit-execute-pending-tasks
description: "Executes the Multi-Agent Task Dispatcher & Subtask Execution prompt. Autonomously orchestrate and execute ALL pending tasks in a continuous N-step self-loop until the entire queue is completely resolved without a single failure. Use when the user asks to run commit-execute-pending-tasks, or the task is about commits, boolean cleanups, or multi-agent code fixes."
---

# Multi-Agent Task Dispatcher & Subtask Execution

Source prompt: `01-prompts/09-commit-and-multi-agent-code-fix/02-execute-pending-tasks.md`

## Instructions

1. Read `01-prompts/09-commit-and-multi-agent-code-fix/02-execute-pending-tasks.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously orchestrate and execute ALL pending tasks in a continuous N-step self-loop until the entire queue is completely resolved without a single failure.
