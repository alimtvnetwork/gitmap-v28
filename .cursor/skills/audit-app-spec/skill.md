---
name: audit-app-spec
description: "Executes the Application Specification Audit prompt. Autonomously perform blind-AI readiness audits on application specs, generate scorecards and remediation matrices in 02-spec/25-app-spec-audit/ with strict no-build and no-test execution. Use when the user asks to run audit-app-spec, or the task is about planning, spec steps, or an app-spec audit."
---

# Application Specification Audit

Source prompt: `01-prompts/13-plan-audit/03-audit-app-spec.md`

## Instructions

1. Read `01-prompts/13-plan-audit/03-audit-app-spec.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Reuse First: I have rigorously scanned and learned `03-ai-scripts/readme.md` to check if a helper script already exists before writing any new temporary code with strict no-build and no-test execution.
