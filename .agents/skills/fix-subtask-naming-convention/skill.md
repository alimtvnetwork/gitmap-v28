---
name: fix-subtask-naming-convention
description: "Executes the Subtask Naming Normalization & Sequence Repair prompt. Your objective is to deeply audit the .ai-memory/plans/ directory for any subtask files that incorrectly use the SS- or SS-XX- prefix and fix them. Use when the user asks to run fix-subtask-naming-convention, or the task is about a legacy execute loop."
---

# Subtask Naming Normalization & Sequence Repair

Source prompt: `01-prompts/19-old-execute-prompts/02-fix-subtask-naming-convention.md`

## Instructions

1. Read `01-prompts/19-old-execute-prompts/02-fix-subtask-naming-convention.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Your objective is to deeply audit the .ai-memory/plans/ directory for any subtask files that incorrectly use the SS- or SS-XX- prefix and fix them.
