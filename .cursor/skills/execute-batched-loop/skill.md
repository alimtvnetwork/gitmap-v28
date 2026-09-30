---
name: execute-batched-loop
description: "Executes the Execute Batched Loop prompt. Execute pending tasks from .ai-memory/plans/pending/ using a strictly batched multi-agent loop. Use when the user asks to run execute-batched-loop, or the task is about executing pending tasks, a parent task, or a batched loop."
---

# Execute Batched Loop

Source prompt: `01-prompts/14-execute/03-execute-batched-loop.md`

## Instructions

1. Read `01-prompts/14-execute/03-execute-batched-loop.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Execute pending tasks from .ai-memory/plans/pending/ using a strictly batched multi-agent loop.
