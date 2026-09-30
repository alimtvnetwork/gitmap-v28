---
name: fix-with-rca
description: "Executes the Bug Fix with 4-Part RCA & Regression Verification prompt. Autonomously fix the provided bug/issue, strictly enforcing coding guidelines, and document the complete RCA before pushing the code. Use when the user asks to run fix-with-rca, or the task is about bug fixes or root-cause analysis."
---

# Bug Fix with 4-Part RCA & Regression Verification

Source prompt: `01-prompts/07-bug-fix/01-fix-with-rca.md`

## Instructions

1. Read `01-prompts/07-bug-fix/01-fix-with-rca.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously fix the provided bug/issue, strictly enforcing coding guidelines, and document the complete RCA before pushing the code.
