---
name: fix-spec-from-audit
description: "Executes the Fix Spec From Audit prompt. Autonomously ingest the latest specification audit file from 02-spec/25-app-spec-audit/, decompose every finding into an exhaustive 1:1 remediation checklist, spawn parallel subagents to fix the specifications, verify. Use when the user asks to run fix-spec-from-audit, or the task is about planning, spec steps, or an app-spec audit."
---

# Fix Spec From Audit

Source prompt: `01-prompts/13-plan-audit/04-fix-spec-from-audit.md`

## Instructions

1. Read `01-prompts/13-plan-audit/04-fix-spec-from-audit.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously ingest the latest specification audit file from 02-spec/25-app-spec-audit/, decompose every finding into an exhaustive 1:1 remediation checklist, spawn parallel subagents to fix the specifications, verify.
