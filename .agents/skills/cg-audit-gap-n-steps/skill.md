---
name: cg-audit-gap-n-steps
description: "Executes the Gap Analysis & N-Step Guideline Audit - Planning Spec prompt. First N/2 steps (Phase 1): Deeply scan the entire codebase file-by-file, dividing N steps across files with 30-50 nested atomic checks per file, scoring guideline compliance from 0 to 100, and writing the master audit. Use when the user asks to run cg-audit-gap-n-steps, or the task is about planning or executing a coding-guideline audit."
---

# Gap Analysis & N-Step Guideline Audit - Planning Spec

Source prompt: `01-prompts/05-coding-guidelines/03-cg-audit-gap-n-steps.md`

## Instructions

1. Read `01-prompts/05-coding-guidelines/03-cg-audit-gap-n-steps.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

First N/2 steps (Phase 1): Deeply scan the entire codebase file-by-file, dividing N steps across files with 30-50 nested atomic checks per file, scoring guideline compliance from 0 to 100, and writing the master audit.
