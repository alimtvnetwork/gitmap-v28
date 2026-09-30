---
name: write-memory
description: "Executes the Memory Persistence & Issue Logging prompt. Persist what happened this turn so the next AI knows everything without guessing. Use when the user asks to run write-memory, or the task is about reading memory, writing memory, proofreading, or conversation logs."
---

# Memory Persistence & Issue Logging

Source prompt: `01-prompts/03-read-write/04-write-memory.md`

## Instructions

1. Read `01-prompts/03-read-write/04-write-memory.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Persist what happened this turn so the next AI knows everything without guessing.
