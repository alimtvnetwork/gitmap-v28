---
name: letterly-run
description: >-
  Formats raw voice dictation into project run instructions using run.ps1 / run.sh driven by run.config.json and the cursor run skill suffix.
---

# Run Mode (Cursor) — Letterly Prompt Formatter

Format whatever input text is provided according to the exact project run template below, following the execute N-steps structure and using the Cursor skill format. Do NOT add conversational filler or commentary.

1. Clean the input text verbatim while strictly capturing target service names, ports, run flags, and environment parameters.
2. Structure the output starting immediately with `# High Priority Instruction`.
3. Put `${Input Text Verbatim}` directly beneath the high priority header.
4. Construct `# Actionable Items Must Follow Non-Negotiable`:
   - Item 1 is ALWAYS: `1. Inspect run.config.json for target service configuration and port mappings`
   - Item 2 is ALWAYS: `2. Execute run script (./run.ps1 or ./run.sh) with automatic dependency verification`
   - Item 3 is ALWAYS: `3. If dependencies or compilers are missing, fall back to gitmap aum install or ./local-install.ps1, then auto-rerun`
   - Item 4 is ALWAYS: `4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well`
   - Item 5..N capture discrete run parameters, port bindings, or service profiles from the input.
5. Append the mandatory run skill invocation suffix `[run](file;.cursor/skills/run)`.
6. Output ONLY the resulting formatted markdown block.
7. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.
8. Additional Instructions Mandate: Always append the relative paths directive under ## Additional Instructions: '- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.'

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# High Priority Instruction

${Input Text Verbatim}

# Actionable Items Must Follow Non-Negotiable

1. Inspect run.config.json for target service configuration and port mappings
2. Execute run script (./run.ps1 or ./run.sh) with automatic dependency verification
3. If dependencies or compilers are missing, fall back to gitmap aum install or ./local-install.ps1, then auto-rerun
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. [Discrete run directive extracted from input]

## Must follow and spawn agent using

[run](file;.cursor/skills/run)

## Additional Instructions

- [/plan](slashCommand;plan) first before doing the work to reduce the credits.
- [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
