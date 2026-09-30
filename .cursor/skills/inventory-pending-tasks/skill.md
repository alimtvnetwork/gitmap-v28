---
name: inventory-pending-tasks
description: "Executes the Inventory Audit of Pending Tasks - Read-Only Proposal prompt. Perform a strictly read-only scan of the entire repository, 02-spec/, and .ai-memory/ directory to compile a comprehensive, deduplicated inventory of every pending task, subtask, unresolved issue, and open requirement. Use when the user asks to run inventory-pending-tasks, or the task is about planning, spec steps, or an app-spec audit."
---

# Inventory Audit of Pending Tasks - Read-Only Proposal

Source prompt: `01-prompts/13-plan-audit/01-inventory-pending-tasks.md`

## Instructions

1. Read `01-prompts/13-plan-audit/01-inventory-pending-tasks.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Perform a strictly read-only scan of the entire repository, 02-spec/, and .ai-memory/ directory to compile a comprehensive, deduplicated inventory of every pending task, subtask, unresolved issue, and open requirement.
