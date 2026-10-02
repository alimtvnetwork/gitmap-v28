---
name: execute-parent-task-with-n-steps-v4
description: "Executes the [V4] Parent Task N-Step Loop: Antigravity-Native Orchestrator prompt (N = 300 step ceiling, A = 2 worker subagents, evidence-gated checks, explicit-path staging, resumable ledger). Use when the user asks to run execute-parent-task-with-n-steps-v4, or wants a parent task executed in N steps with the V4 prompt."
---

# [V4] Parent Task N-Step Loop: Antigravity-Native Orchestrator (must follow)

Source prompt: `01-prompts/14-execute/11-execute-parent-task-with-n-steps-v4.md`

## Instructions

1. Read `01-prompts/14-execute/11-execute-parent-task-with-n-steps-v4.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. The prompt targets Google Antigravity. In Cursor, use the available subagent tool wherever the prompt says `invoke_subagent`, with the same worker brief and report block; if no subagent tool is available, follow the prompt's fallback rule (R5).
4. A direct instruction in the current user message overrides the prompt when they conflict.
5. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Execute a parent task end to end: capture it verbatim, plan it in the repo, run it through worker subagents that own disjoint files, back every claim with evidence, and finish with one commit that holds only this task's files.
