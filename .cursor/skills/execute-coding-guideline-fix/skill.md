---
name: execute-coding-guideline-fix
description: "Executes the Execute Coding Guideline Fix prompt. Autonomously orchestrate and apply concrete, surgical refactoring fixes for all coding guideline violations across the target codebase in bounded 5-8 file micro-batches until 100% green without stopping with strict. Use when the user asks to run execute-coding-guideline-fix, or the task is about coding-guideline execution for this rule family."
---

# Execute Coding Guideline Fix

Source prompt: `01-prompts/15-cg-execute/01-execute-coding-guideline-fix.md`

## Instructions

1. Read `01-prompts/15-cg-execute/01-execute-coding-guideline-fix.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously orchestrate and apply concrete, surgical refactoring fixes for all coding guideline violations across the target codebase in bounded 5-8 file micro-batches until 100% green without stopping with strict.
