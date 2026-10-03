---
name: plan-spec-steps-v2
description: "Executes the Specification Planning Engine (v2) prompt. Autonomously author comprehensive application specifications in 02-spec/21-app/ and lean subtask plans in .ai-memory/plans/ with strict no-build and no-test execution. Use when the user asks to run plan-spec-steps-v2, or the task is about planning, spec steps, or an app-spec audit."
---

# Plan Spec Steps (v2) — Specification Planning Engine

Source prompt: `01-prompts/13-plan-audit/02-plan-spec-steps-v2.md`

## Instructions

1. Read `01-prompts/13-plan-audit/02-plan-spec-steps-v2.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Reuse First: I have rigorously scanned and learned `03-ai-scripts/readme.md` to check if a helper script already exists before writing any new temporary code with strict no-build and no-test execution.
