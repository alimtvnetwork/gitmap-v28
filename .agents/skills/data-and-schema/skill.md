---
name: data-and-schema
description: "Executes the Data And Schema prompt. Autonomously scan, plan, refactor, and fix all database schema, model, and query violations across the codebase, modifying migration scripts and ORM entities directly to enforce PascalCase tables, camelCase columns. Use when the user asks to run data-and-schema, or the task is about coding-guideline execution for this rule family."
---

# Data And Schema

Source prompt: `01-prompts/15-cg-execute/07-data-and-schema.md`

## Instructions

1. Read `01-prompts/15-cg-execute/07-data-and-schema.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously scan, plan, refactor, and fix all database schema, model, and query violations across the codebase, modifying migration scripts and ORM entities directly to enforce PascalCase tables, camelCase columns.
