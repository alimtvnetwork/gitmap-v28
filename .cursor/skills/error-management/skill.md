---
name: error-management
description: "Executes the Error Management prompt. Autonomously scan, plan, refactor, and fix all error management violations across the codebase, modifying source files directly to implement appfault.AppError wrappers, outer error handling, specialized exit helpers. Use when the user asks to run error-management, or the task is about coding-guideline execution for this rule family."
---

# Error Management

Source prompt: `01-prompts/15-cg-execute/02-error-management.md`

## Instructions

1. Read `01-prompts/15-cg-execute/02-error-management.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously scan, plan, refactor, and fix all error management violations across the codebase, modifying source files directly to implement appfault.AppError wrappers, outer error handling, specialized exit helpers.
