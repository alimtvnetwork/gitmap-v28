---
name: execute-robust-loop
description: "Executes the Resilient Multi-Agent Loop Execution prompt. Read all pending tasks from .ai-memory/, allocate small micro-portions of work to sub-agents, and execute them in a continuous self-loop. Use when the user asks to run execute-robust-loop, or the task is about a legacy execute loop."
---

# Resilient Multi-Agent Loop Execution

Source prompt: `01-prompts/19-old-execute-prompts/01-execute-robust-loop.md`

## Instructions

1. Read `01-prompts/19-old-execute-prompts/01-execute-robust-loop.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Read all pending tasks from .ai-memory/, allocate small micro-portions of work to sub-agents, and execute them in a continuous self-loop.
