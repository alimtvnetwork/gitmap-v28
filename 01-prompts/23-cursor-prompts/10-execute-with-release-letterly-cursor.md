# Execute with Release Mode (Cursor) — Letterly Prompt Formatter

Format whatever input text is provided according to the exact execution with minor release template below, following the execute N-steps structure, using the Cursor skill format. Do NOT add conversational filler or commentary.

1. Capture and clean the input text verbatim, stripping verbal filler words (`um`, `ah`, `uh`, `like`) while preserving every technical directive, parameter, flag, and file path.
2. Derive `<task title>` from the input text — it is the SUBJECT of the task (e.g. `SEO Writing and Folder Structure Instructions`). Structure the output starting immediately with `# <task title>: high priority instruction, non-negotiable task`.
3. Put `${Input Text Verbatim}` directly beneath the title line.
4. After the verbatim input, output both the task name and the slug as subheaders:
   - `## task: <task title>` — the human-readable subject/name of the task.
   - `## slug: <task-slug>` — the lowercase-hyphenated task slug (e.g. `## slug: seo-writing-and-folder-structure-instructions`), NEVER `## slug: high-priority-instruction`.
   Both the task name and slug are required so AI agents can identify and register the task via GitMap (`gitmap agent task init --title "<task title>" --slug "<task-slug>"`).
5. Construct `# Actionable Items Must Follow Non-Negotiable`:
   - Item 1 is ALWAYS: `1. Git pull BEFORE starting any work — first try \`gitmap pull <repo path>\` (or run \`gitmap pull\` in repository folder); if that fails, fall back to regular \`git pull\` and resolve any merge conflicts so the repository is completely up to date and clean before starting.`
   - Item 2 is ALWAYS: `2. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first`
   - Item 3 is ALWAYS: `3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String`
   - Item 4 is ALWAYS: `4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well`
   - Item 5..N are sequential, discrete technical directives extracted from the input.
   - Penultimate Item is ALWAYS: `Verify live CI/CD pipeline health via gitmap pe -t until green`
   - Final Item is ALWAYS: `Execute minor version bump release ceremony via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>", update changelog.md, commit atomically via gitmap cpf, tag release, and push to remote tracking branch`
6. Append the mandatory release skill invocation suffix pointing to `[minor-bump](file;.cursor/skills/minor-bump)`.
7. Output ONLY the resulting formatted prompt as plain markdown text. STRICT — NEVER wrap it in a fenced code block: never emit ```markdown, ```plaintext, or any ``` fence at the start or end of the output. `#` headers and inline markdown are correct; the fence around them is forbidden. No conversational filler, no commentary before or after it.
8. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.
9. Additional Instructions Mandate: Always append the relative paths directive under ## Additional Instructions: '- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.'
10. Step Budget Mandate: Execute within the 400-step lifecycle budget (Phase 1: 100 plan, Phase 2: 150 intermediate/archival, Phase 3: 150 implementation/release).
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
5. [Fifth actionable technical directive extracted from input]
6. Verify live CI/CD pipeline health via gitmap pe -t until green
7. Execute minor version bump release ceremony via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>", update changelog.md, commit atomically via gitmap cpf, tag release, and push to remote tracking branch

## Must follow and spawn agent using

[minor-bump](file;.cursor/skills/minor-bump)

## Additional Instructions

- Execute within the 400-step lifecycle budget: Phase 1 (Steps 1..100) planning and spec authoring, Phase 2 (Steps 101..250) sanitization, archival, and prerequisites, Phase 3 (Steps 251..400) worker swarm implementation, green CI verification, and minor release.
- [/plan](slashCommand;plan) first before doing the work to reduce the credits.
- [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap Heatmap (`gitmap heatmap discuss` / `gitmap heatmap log` / `gitmap heatmap live`) with transparent 3-tier database location reporting.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
