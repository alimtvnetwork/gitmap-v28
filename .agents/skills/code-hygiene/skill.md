---
name: code-hygiene
description: "Executes the Code Hygiene prompt. Autonomously scan, plan, refactor, and fix all code hygiene, file size, parameter bloat, LF line ending. Use when the user asks to run code-hygiene, or the task is about coding-guideline execution for this rule family."
---

# Code Hygiene

Source prompt: `01-prompts/15-cg-execute/09-code-hygiene.md`

## Instructions

1. Read `01-prompts/15-cg-execute/09-code-hygiene.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously scan, plan, refactor, and fix all code hygiene, file size, parameter bloat, LF line ending.
