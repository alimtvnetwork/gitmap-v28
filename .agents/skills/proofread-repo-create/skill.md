---
name: proofread-repo-create
description: "Executes the Proofread Repo Create Instruction prompt. Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion without a single failure with strict no-build and no-test execution. Use when the user asks to run proofread-repo-create, or the task is about reading memory, writing memory, proofreading, or conversation logs."
---

# Proofread Repo Create Instruction

Source prompt: `01-prompts/03-read-write/10-proofread-repo-create.md`

## Instructions

1. Read `01-prompts/03-read-write/10-proofread-repo-create.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion without a single failure with strict no-build and no-test execution.
