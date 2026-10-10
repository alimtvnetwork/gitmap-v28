# Plan Mode (Cursor) — Letterly Prompt Formatter

Format whatever input text is provided according to the exact planning template below, following the execute N-steps structure and using the Cursor skill format. Do NOT add conversational filler or commentary.

1. Clean the input text verbatim while strictly capturing all architectural requirements, scope boundaries, and design constraints.
2. Structure the output starting immediately with `# High Priority Instruction`.
3. Put `${Input Text Verbatim}` directly beneath the high priority header.
4. Construct `# Actionable Items Must Follow Non-Negotiable`:
   - Item 1 is ALWAYS: `1. Write the plan and architectural spec first under 02-spec/21-app/<slug>/ adhering to 02-spec/01-spec-authoring-guide/`
   - Item 2 is ALWAYS: `2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String`
   - Item 3 is ALWAYS: `3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well`
   - Item 4 is ALWAYS: `4. Write the first step on task decomposition and bounded subtasks under .ai-memory/plans/subtasks/<slug>/`
   - Item 5 is ALWAYS: `5. Enforce strict no-build and no-test rules throughout the planning phase`
   - Item 6..N capture discrete architectural requirements from the input.
5. Append the mandatory planning skill invocation suffix `[plan-spec-steps-v2](file;.cursor/skills/plan-spec-steps-v2)`.
6. Output ONLY the resulting formatted prompt as plain markdown text. STRICT — NEVER wrap it in a fenced code block: never emit ```markdown, ```plaintext, or any ``` fence at the start or end of the output. `#` headers and inline markdown are correct; the fence around them is forbidden. No conversational filler, no commentary before or after it.
7. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.
8. Additional Instructions Mandate: Always append the relative paths directive under ## Additional Instructions: '- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.'

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# High Priority Instruction

${Input Text Verbatim}

# Actionable Items Must Follow Non-Negotiable

1. Write the plan and architectural spec first under 02-spec/21-app/<slug>/ adhering to 02-spec/01-spec-authoring-guide/
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Write the first step on task decomposition and bounded subtasks under .ai-memory/plans/subtasks/<slug>/
5. Enforce strict no-build and no-test rules throughout the planning phase
6. Define binary acceptance criteria for every subtask

## Must follow and spawn agent using

[plan-spec-steps-v2](file;.cursor/skills/plan-spec-steps-v2)

## Additional Instructions

- [/plan](slashCommand;plan) first before doing the work to reduce the credits.
- [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
