# Plan Mode (Cursor) — Letterly Prompt Formatter

Format whatever input text is provided according to the exact planning template below, following the execute N-steps structure, using the Cursor skill format. Do NOT add conversational filler or commentary.

1. Clean the input text verbatim while strictly capturing all architectural requirements, scope boundaries, and design constraints.
2. Derive `<task title>` from the input text — it is the SUBJECT of the task (e.g. `SEO Writing and Folder Structure Instructions`). Structure the output starting immediately with `# <task title>: high priority instruction, non-negotiable task`.
3. Put `${Input Text Verbatim}` directly beneath the title line.
4. After the verbatim input, output both the task name and the slug as subheaders:
   - `## task: <task title>` — the human-readable subject/name of the task.
   - `## slug: <task-slug>` — the lowercase-hyphenated task slug (e.g. `## slug: seo-writing-and-folder-structure-instructions`), NEVER `## slug: high-priority-instruction`.
   Both the task name and slug are required so AI agents can identify and register the task via GitMap (`gitmap agent task init --title "<task title>" --slug "<task-slug>"`).
5. Construct `# Actionable Items Must Follow Non-Negotiable`:
   - Item 1 is ALWAYS: `1. Git pull BEFORE starting any work — first try \`gitmap pull <repo path>\` (or run \`gitmap pull\` in repository folder); if that fails, fall back to regular \`git pull\` and resolve any merge conflicts so the repository is completely up to date and clean before starting.`
   - Item 2 is ALWAYS: `2. Write the plan and architectural spec first under 02-spec/21-app/<slug>/ adhering to 02-spec/01-spec-authoring-guide/`
   - Item 3 is ALWAYS: `3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String`
   - Item 4 is ALWAYS: `4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well`
   - Item 5 is ALWAYS: `5. Write the first step on task decomposition and bounded subtasks under .ai-memory/plans/subtasks/<slug>/`
   - Item 6 is ALWAYS: `6. Enforce strict no-build and no-test rules throughout the planning phase`
   - Item 7..N capture discrete architectural requirements from the input.
6. Append the mandatory planning skill invocation suffix `[plan-spec-steps-v2](file;.cursor/skills/plan-spec-steps-v2)`.
7. Output ONLY the resulting formatted prompt as plain markdown text. STRICT — NEVER wrap it in a fenced code block: never emit ```markdown, ```plaintext, or any ``` fence at the start or end of the output. `#` headers and inline markdown are correct; the fence around them is forbidden. No conversational filler, no commentary before or after it.
8. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.
9. Additional Instructions Mandate: Always append the relative paths directive under ## Additional Instructions: '- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.'
10. Step Budget Mandate: Execute within the 400-step lifecycle budget, dedicating Phase 1 (Steps 1..100) to planning, spec authoring under 02-spec/21-app/<slug>/, and subtask decomposition under .ai-memory/plans/subtasks/<slug>/.
11. Heatmap Mandate: Track in-flight file touches, prevent write collisions, and query the cognitive activity graph via GitMap Heatmap (gitmap heatmap discuss / log / live) with transparent 3-tier database location reporting.

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# <task title>: high priority instruction, non-negotiable task

${Input Text Verbatim}

## task: <task title>
## slug: <task-slug>

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write the plan and architectural spec first under 02-spec/21-app/<slug>/ adhering to 02-spec/01-spec-authoring-guide/
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Write the first step on task decomposition and bounded subtasks under .ai-memory/plans/subtasks/<slug>/
6. Enforce strict no-build and no-test rules throughout the planning phase
7. Define binary acceptance criteria for every subtask

## Must follow and spawn agent using

[plan-spec-steps-v2](file;.cursor/skills/plan-spec-steps-v2)

## Additional Instructions

- Execute within the 400-step lifecycle budget: allocate Phase 1 (Steps 1..100) for planning, spec authoring, and subtask decomposition.
- [/plan](slashCommand;plan) first before doing the work to reduce the credits.
- [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap Heatmap (`gitmap heatmap discuss` / `gitmap heatmap log` / `gitmap heatmap live`) with transparent 3-tier database location reporting.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
