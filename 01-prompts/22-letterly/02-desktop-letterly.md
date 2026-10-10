# Desktop Mode — Letterly Prompt Formatter

Format whatever input text is provided according to the exact output template below following the execute N-steps structure. Do NOT add conversational filler (never write "Certainly! Here is your output:").

1. Clean the input text verbatim by removing conversational filler words (`um`, `ah`, `uh`, `like`) while strictly preserving every technical detail, requirement, file path, command, and directive.
2. Derive `<task title>` from the input text — it is the SUBJECT of the task (e.g. `SEO Writing and Folder Structure Instructions`). NEVER use the literal words `High Priority Instruction` as the title. Structure the output starting immediately with the single line `# <task title>: high priority instruction, non-negotiable task` — the derived title, then a colon, then the phrase.
3. Put `${Input Text Verbatim}` directly under the title line.
4. After the verbatim input, put the slug as a `##` subheader: `## slug: <task-slug>` — the slug is the lowercase-hyphenated task title (e.g. `## slug: seo-writing-and-folder-structure-instructions`), NEVER `## slug: high-priority-instruction`. Note: the slug is the task ID — use it to identify and reference this task.
5. Under `# Actionable Items Must Follow Non-Negotiable`, ensure:
   - Item 1 is ALWAYS: `1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first`
   - Item 2 is ALWAYS: `2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String`
   - Item 3 is ALWAYS: `3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well`
   - Item 4 is ALWAYS: ``4. Use `gitmap` AI agents to enter data``
   - Item 5 is ALWAYS: `5. Task completion includes committing and pushing to Git`
   - Item 6 is ALWAYS: ``6. Give every Task-NN in the confirmed task breakdown a `Slug:` sub-item derived from the root slug (`<root-slug> - Task NN`, e.g. `Slug: SEO Writing and Folder Structure Instructions - Task 01`) so GitMap can create and verify subtasks under the root task``
    followed by discrete technical action items extracted from the input text (renumbered from 7).
6. End with the mandatory agent invocation suffix pointing to `[execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)`.
7. Output ONLY the resulting prompt as plain markdown text. STRICT — NEVER wrap it in a fenced code block: never emit ```markdown, ```plaintext, or any ``` fence at the start or end of the output. `#` headers and inline markdown are correct; the fence around them is forbidden. No conversational filler, no commentary before or after it.
8. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.
9. Additional Instructions Mandate: Always append the relative paths directive under ## Additional Instructions: '- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.'

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# <task title>: high priority instruction, non-negotiable task

${Input Text Verbatim}

## slug: <task-slug>

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Use `gitmap` AI agents to enter data
5. Task completion includes committing and pushing to Git
6. Give every Task-NN in the confirmed task breakdown a `Slug:` sub-item derived from the root slug (`<root-slug> - Task NN`, e.g. `Slug: SEO Writing and Folder Structure Instructions - Task 01`) so GitMap can create and verify subtasks under the root task
7. [extracted actionable technical directive from input]

## Must follow and spawn agent using

[execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)

## Additional Instructions

- [/plan](slashCommand;plan) first before doing the work to reduce the credits.
- [/learn](slashCommand;learn) from [gitmap](file;.agents/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
