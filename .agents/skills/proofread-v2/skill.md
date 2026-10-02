---
name: proofread-v2
description: "Executes the Proofread V2 (On the Fly) Instruction prompt. Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion without a single failure with strict no-build and no-test execution. Use when the user asks to run proofread-v2, or the task is about reading memory, writing memory, proofreading, or conversation logs."
---

# Proofread V2 (On the Fly) Instruction

Source prompt: `01-prompts/03-read-write/09-proofread-v2.md`

## Instructions

1. Read `01-prompts/03-read-write/09-proofread-v2.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion without a single failure with strict no-build and no-test execution.
