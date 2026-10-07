---
name: letterly-mobile
description: >-
  Formats raw voice dictation into mobile single-paragraph goal and learn execution directives followed by execute-parent-task-with-n-steps-v6 skill suffix without extra action items or conversational prefixes.
---

# Mobile Mode — Letterly Prompt Formatter

Format whatever input text is provided according to the exact mobile single-line output template below.

1. Clean the input text verbatim without conversational filler words (`um`, `ah`, `uh`, `like`).
2. Do NOT add `[/goal](slashCommand;goal)` or `[/learn](slashCommand;learn)` at the beginning.
3. Start the line immediately with `# High Priority Instruction: `.
4. Append `${Input Text Verbatim}`.
5. Conclude the single line with ` - strictly use relative git paths only; only add the relative paths, never add the absolute path during your work (this must be respected on the release page and in release notes as well), must follow the skill [execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)`.
6. Output exactly that single continuous paragraph with ZERO newlines, ZERO line breaks, and NO leading "Output" text.
7. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.

${Input Text Verbatim} = The cleaned input text as it is, without filler words.

Output Format:
# High Priority Instruction: ${Input Text Verbatim} - strictly use relative git paths only; only add the relative paths, never add the absolute path during your work (this must be respected on the release page and in release notes as well), must follow the skill [execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)
