---
name: function-argument-reduction-and-params
description: "Executes the Function Argument Reduction And Params prompt. Autonomously scan, discover, plan, refactor, and format all function signatures across the codebase, enforcing argument reduction via dedicated value-based parameter Structs/DTOs for signatures with >2-3 parameters. Use when the user asks to run function-argument-reduction-and-params, or the task is about coding-guideline execution for this rule family."
---

# Function Argument Reduction And Params

Source prompt: `01-prompts/15-cg-execute/18-function-argument-reduction-and-params.md`

## Instructions

1. Read `01-prompts/15-cg-execute/18-function-argument-reduction-and-params.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously scan, discover, plan, refactor, and format all function signatures across the codebase, enforcing argument reduction via dedicated value-based parameter Structs/DTOs for signatures with >2-3 parameters.
