---
name: letterly-desktop
description: >-
  Formats raw voice dictation into desktop high-priority instructions, non-negotiable action items starting with write spec and plan, and execute-parent-task-with-n-steps-v6 skill suffix.
---

# Desktop Mode — Letterly Prompt Formatter

Format whatever input text is provided according to the exact output template below following the execute N-steps structure. Do NOT add conversational filler (never write "Certainly! Here is your output:").

1. Clean the input text verbatim by removing conversational filler words (`um`, `ah`, `uh`, `like`) while strictly preserving every technical detail, requirement, file path, command, and directive.
2. Structure the output starting immediately with `# High Priority Instruction`.
3. Put `${Input Text Verbatim}` directly under the high priority header.
4. Under `# Actionable Items Must Follow Non-Negotiable`, ensure the first item is ALWAYS:
   `1. Write spec and plan first`
   followed by discrete technical action items extracted from the input text.
5. End with the mandatory agent invocation suffix pointing to `[execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)`.
6. Output ONLY the resulting markdown block.

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# High Priority Instruction

${Input Text Verbatim}

# Actionable Items Must Follow Non-Negotiable

1. Write spec and plan first
2. [Second actionable technical directive extracted from input]
3. [Third actionable technical directive extracted from input]

Must follow and spawn agent using

[execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)

## Additional Instructions

learn [/learn](slashCommand;learn) if you have to learn something and [/plan](slashCommand;plan) stuff before working please.
