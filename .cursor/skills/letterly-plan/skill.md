---
name: letterly-plan
description: >-
  Formats raw voice dictation into high-priority architecture and planning instructions, actionable spec decomposition items, and plan-spec-steps-v2 skill invocation suffix.
---

# Plan Mode — Letterly Prompt Formatter

Format whatever input text is provided according to the exact planning template below, following the execute N-steps structure. Do NOT add conversational filler or commentary.

1. Clean the input text verbatim while strictly capturing all architectural requirements, scope boundaries, and design constraints.
2. Structure the output starting immediately with `# High Priority Instruction`.
3. Put `${Input Text Verbatim}` directly beneath the high priority header.
4. Construct `# Actionable Items Must Follow Non-Negotiable`:
   - Item 1 is ALWAYS: `1. Write the plan and architectural spec first under 02-spec/21-app/<slug>/`
   - Item 2 is ALWAYS: `2. Write the first step on task decomposition and bounded subtasks under .ai-memory/plans/subtasks/<slug>/`
   - Item 3 is ALWAYS: `3. Enforce strict no-build and no-test rules throughout the planning phase`
   - Item 4..N capture discrete architectural requirements from the input.
5. Append the mandatory planning skill invocation suffix [plan-spec-steps-v2](file;.cursor/skills/plan-spec-steps-v2).
6. Output ONLY the resulting formatted markdown block.

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# High Priority Instruction

${Input Text Verbatim}

# Actionable Items Must Follow Non-Negotiable

1. Write the plan and architectural spec first under 02-spec/21-app/<slug>/
2. Write the first step on task decomposition and bounded subtasks under .ai-memory/plans/subtasks/<slug>/
3. Enforce strict no-build and no-test rules throughout the planning phase
4. Define binary acceptance criteria for every subtask

Must follow and spawn agent using

[plan-spec-steps-v2](file;.cursor/skills/plan-spec-steps-v2)

## Additional Instructions

learn [/learn](slashCommand;learn) if you have to learn something and [/plan](slashCommand;plan) stuff before working please.
