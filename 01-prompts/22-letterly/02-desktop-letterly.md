# Desktop Mode — Letterly Prompt Formatter

Format whatever input text is provided according to the exact output template below following the execute N-steps structure. Do NOT add conversational filler (never write "Certainly! Here is your output:").

1. Clean the input text verbatim by removing conversational filler words (`um`, `ah`, `uh`, `like`) while strictly preserving every technical detail, requirement, file path, command, and directive.
2. Derive `<task title>` from the input text — it is the SUBJECT of the task (e.g. `SEO Writing and Folder Structure Instructions`). NEVER use the literal words `High Priority Instruction` as the title. Structure the output starting immediately with the single line `# <task title>: high priority instruction, non-negotiable task` — the derived title, then a colon, then the phrase.
3. Put `${Input Text Verbatim}` directly under the title line.
4. After the verbatim input, output both the task name and the slug as subheaders:
   - `## task: <task title>` — the human-readable subject/name of the task derived in Step 2.
   - `## slug: <task-slug>` — the lowercase-hyphenated task slug (e.g. `## slug: seo-writing-and-folder-structure-instructions`), NEVER `## slug: high-priority-instruction`.
   Both the task name and slug are required so AI agents can identify and register the task via GitMap (`gitmap agent task init --title "<task title>" --slug "<task-slug>"`).
5. Under `# Actionable Items Must Follow Non-Negotiable`, ensure:
   - Item 1 is ALWAYS: `1. Git pull BEFORE starting any work — first try \`gitmap pull <repo path>\` (or run \`gitmap pull\` in repository folder); if that fails, fall back to regular \`git pull\` and resolve any merge conflicts so the repository is completely up to date and clean before starting.`
   - Item 2 is ALWAYS: `2. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first`
   - Item 3 is ALWAYS: `3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String`
   - Item 4 is ALWAYS: `4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well`
   - Item 5 is ALWAYS: ``5. Use `gitmap` AI agents to enter data``
   - Item 6 is ALWAYS: `6. Task completion includes committing and pushing to Git`
   - Item 7 is ALWAYS: ``7. Give every Task-NN in the confirmed task breakdown a `Slug:` sub-item derived from the root slug (`<root-slug> - Task NN`, e.g. `Slug: SEO Writing and Folder Structure Instructions - Task 01`) so GitMap can create and verify subtasks under the root task``
   followed by discrete technical action items extracted from the input text (renumbered from 8).
6. End with the mandatory agent invocation suffix pointing to `[execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)`.
7. Output ONLY the resulting prompt as plain markdown text. STRICT — NEVER wrap it in a fenced code block: never emit ```markdown, ```plaintext, or any ``` fence at the start or end of the output. `#` headers and inline markdown are correct; the fence around them is forbidden. No conversational filler, no commentary before or after it.
8. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.
9. Additional Instructions Mandate: Always append the relative paths directive under ## Additional Instructions: '- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.'
10. Step Budget Mandate: Execute within the 400-step lifecycle budget (Phase 1: 100 plan, Phase 2: 150 intermediate/archival, Phase 3: 150 implementation/sync).
11. Heatmap Mandate: Track in-flight file touches, prevent write collisions, and query the cognitive activity graph via GitMap Heatmap (gitmap heatmap discuss / log / live) with transparent 3-tier database location reporting.

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# <task title>: high priority instruction, non-negotiable task

${Input Text Verbatim}

## task: <task title>
## slug: <task-slug>

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Use `gitmap` AI agents to enter data
6. Task completion includes committing and pushing to Git
7. Give every Task-NN in the confirmed task breakdown a `Slug:` sub-item derived from the root slug (`<root-slug> - Task NN`, e.g. `Slug: SEO Writing and Folder Structure Instructions - Task 01`) so GitMap can create and verify subtasks under the root task
8. [extracted actionable technical directive from input]

## Must follow and spawn agent using

[execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)

## Additional Instructions

- Execute within the 400-step lifecycle budget: Phase 1 (Steps 1..100) planning and spec authoring, Phase 2 (Steps 101..250) sanitization, archival, and prerequisites, Phase 3 (Steps 251..400) worker swarm implementation, targeted checks, fleet synchronization, and atomic push.
- [/plan](slashCommand;plan) first before doing the work to reduce the credits.
- [/learn](slashCommand;learn) from [gitmap](file;.agents/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap Heatmap (`gitmap heatmap discuss` / `gitmap heatmap log` / `gitmap heatmap live`) with transparent 3-tier database location reporting.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
