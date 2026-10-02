---
name: execute-parent-task-with-n-steps-v4
description: >-
  Use this skill when the user asks you to execute a parent task with N steps using the V4 prompt (N = 300 step ceiling, A = 2 worker subagents via invoke_subagent, evidence-gated checks, explicit-path staging, and a resumable ledger).
---

# [V4] Parent Task N-Step Loop: Antigravity-Native Orchestrator (must follow)

Source prompt: `01-prompts/14-execute/11-execute-parent-task-with-n-steps-v4.md`

## Instructions

1. Read `01-prompts/14-execute/11-execute-parent-task-with-n-steps-v4.md` in full before doing anything else.
2. The parent task is the text the user gave with this skill (after `/execute-parent-task-with-n-steps-v4`), plus any instructions above it.
3. Execute that prompt exactly. It is the source of truth; this skill adds no rules of its own.
4. A direct instruction in the user's message overrides the prompt when they conflict.
5. Do not shorten, paraphrase, or skip any step or checklist item in the prompt.
