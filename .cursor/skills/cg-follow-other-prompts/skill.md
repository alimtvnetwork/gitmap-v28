---
name: cg-follow-other-prompts
description: "Executes the Cg Follow Other Prompts prompt. Autonomously ingest, follow, and execute the referenced external prompts, task instructions, and coding guideline directives across the target codebase in bounded 5-8 file micro-batches with strict no-build and no-test. Use when the user asks to run cg-follow-other-prompts, or the task is about coding-guideline execution for this rule family."
---

# Cg Follow Other Prompts

Source prompt: `01-prompts/15-cg-execute/32-cg-follow-other-prompts.md`

## Instructions

1. Read `01-prompts/15-cg-execute/32-cg-follow-other-prompts.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously ingest, follow, and execute the referenced external prompts, task instructions, and coding guideline directives across the target codebase in bounded 5-8 file micro-batches with strict no-build and no-test.
