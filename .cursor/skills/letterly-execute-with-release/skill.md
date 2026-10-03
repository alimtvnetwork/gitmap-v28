---
name: letterly-execute-with-release
description: >-
  Formats raw voice dictation into execute N-steps structure that executes the task, verifies CI/CD, and directly initiates minor release ceremony via minor-bump skill, using the cursor skill format.
---

# Execute with Release Mode (Cursor) — Letterly Prompt Formatter

Format whatever input text is provided according to the exact execution with minor release template below, following the execute N-steps structure and using the Cursor skill format. Do NOT add conversational filler or commentary.

1. Capture and clean the input text verbatim, stripping verbal filler words (`um`, `ah`, `uh`, `like`) while preserving every technical directive, parameter, flag, and file path.
2. Structure the output starting immediately with `[/plan](slashCommand;plan)` followed by `# High Priority Instruction`.
3. Output `${Input Text Verbatim}` directly beneath the header.
4. Construct `# Actionable Items Must Follow Non-Negotiable`:
   - Item 1 is ALWAYS: `1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first`
   - Item 2..N are sequential, discrete technical directives extracted from the input.
   - Penultimate Item is ALWAYS: `Verify live CI/CD pipeline health via gitmap pe -t until green`
   - Final Item is ALWAYS: `Execute minor version bump release ceremony via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>", update changelog.md, commit atomically via gitmap cpf, tag release, and push to remote tracking branch`
5. Append the mandatory release skill invocation suffix pointing to [minor-bump](file;.cursor/skills/minor-bump).
6. Output ONLY the resulting formatted markdown block.

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

[/plan](slashCommand;plan)

# High Priority Instruction

${Input Text Verbatim}

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. [Second actionable technical directive extracted from input]
3. [Third actionable technical directive extracted from input]
4. Verify live CI/CD pipeline health via gitmap pe -t until green
5. Execute minor version bump release ceremony via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>", update changelog.md, commit atomically via gitmap cpf, tag release, and push to remote tracking branch

Must follow and spawn agent using

[minor-bump](file;.cursor/skills/minor-bump)

## Additional Instructions

learn [/learn](slashCommand;learn) if you have to learn something and [/plan](slashCommand;plan) stuff before working please.
