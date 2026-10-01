---
name: isolate-destructive-os-and-heavy-unit-tests
description: "Executes the Isolate Destructive Os And Heavy Unit Tests prompt. Autonomously scan, audit, refactor, and verify repository-wide unit tests to ensure that tests NEVER trigger real OS shutdown, reboot, power-off, system modifications, or heavy unmocked system calls, enforcing. Use when the user asks to run isolate-destructive-os-and-heavy-unit-tests, or the task is about coding-guideline execution for this rule family."
---

# Isolate Destructive Os And Heavy Unit Tests

Source prompt: `01-prompts/15-cg-execute/24-isolate-destructive-os-and-heavy-unit-tests.md`

## Instructions

1. Read `01-prompts/15-cg-execute/24-isolate-destructive-os-and-heavy-unit-tests.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously scan, audit, refactor, and verify repository-wide unit tests to ensure that tests NEVER trigger real OS shutdown, reboot, power-off, system modifications, or heavy unmocked system calls, enforcing.
