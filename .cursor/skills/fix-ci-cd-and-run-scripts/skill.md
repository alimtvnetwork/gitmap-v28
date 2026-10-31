---
name: fix-ci-cd-and-run-scripts
description: "Executes the Cross-Platform CI/CD & Run Scripts Fix prompt. First N/2 steps (Phase 1): Review the central CI/CD pipeline definitions (.github/workflows, .gitlab-ci.yml, etc.) and cross-reference them with the local Python runner (03-ai-scripts/06-cicd-local-runner.py). Use when the user asks to run fix-ci-cd-and-run-scripts, or the task is about CI/CD fixes, pipeline creation, or zero-storage Actions."
---

# Cross-Platform CI/CD & Run Scripts Fix

Source prompt: `01-prompts/16-ci-cd/05-fix-ci-cd-and-run-scripts.md`

## Instructions

1. Read `01-prompts/16-ci-cd/05-fix-ci-cd-and-run-scripts.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

First N/2 steps (Phase 1): Review the central CI/CD pipeline definitions (.github/workflows, .gitlab-ci.yml, etc.) and cross-reference them with the local Python runner (03-ai-scripts/06-cicd-local-runner.py).
