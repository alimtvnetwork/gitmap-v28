---
name: spec-write-from-reverse-engineering
description: "Executes the Reverse Engineering Spec Writer & System Architecture Discovery prompt. Autonomously scan, reverse-engineer, and synthesize an exhaustive, multi-file architectural specification of any target codebase into 02-spec/21-app/. Use when the user asks to run spec-write-from-reverse-engineering, or the task is about reading memory, writing memory, proofreading, or conversation logs."
---

# Reverse Engineering Spec Writer & System Architecture Discovery

Source prompt: `01-prompts/03-read-write/08-spec-write-from-reverse-engineering.md`

## Instructions

1. Read `01-prompts/03-read-write/08-spec-write-from-reverse-engineering.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously scan, reverse-engineer, and synthesize an exhaustive, multi-file architectural specification of any target codebase into 02-spec/21-app/.
