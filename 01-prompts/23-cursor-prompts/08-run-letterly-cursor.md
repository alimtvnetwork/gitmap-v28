# Run Mode (Cursor) — Letterly Prompt Formatter

Format whatever input text is provided according to the exact project run template below, following the execute N-steps structure, using the Cursor skill format. Do NOT add conversational filler or commentary.

1. Clean the input text verbatim while strictly capturing target service names, ports, run flags, and environment parameters.
2. Derive `<task title>` from the input text — it is the SUBJECT of the task (e.g. `Run Application Services`). Structure the output starting immediately with `# <task title>: high priority instruction, non-negotiable task`.
3. Put `${Input Text Verbatim}` directly beneath the title line.
4. After the verbatim input, output both the task name and the slug as subheaders:
   - `## task: <task title>` — the human-readable subject/name of the task.
   - `## slug: <task-slug>` — the lowercase-hyphenated task slug (e.g. `## slug: run-application-services`), NEVER `## slug: high-priority-instruction`.
   Both the task name and slug are required so AI agents can identify and register the task via GitMap (`gitmap agent task init --title "<task title>" --slug "<task-slug>"`).
5. Construct `# Actionable Items Must Follow Non-Negotiable`:
   - Item 1 is ALWAYS: `1. Git pull BEFORE starting any work — first try \`gitmap pull <repo path>\` (or run \`gitmap pull\` in repository folder); if that fails, fall back to regular \`git pull\` and resolve any merge conflicts so the repository is completely up to date and clean before starting.`
   - Item 2 is ALWAYS: `2. Inspect run.config.json for target service configuration and port mappings`
   - Item 3 is ALWAYS: `3. Execute run script (./run.ps1 or ./run.sh) with automatic dependency verification`
   - Item 4 is ALWAYS: `4. If dependencies or compilers are missing, fall back to gitmap aum install or ./local-install.ps1, then auto-rerun`
   - Item 5 is ALWAYS: `5. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well`
   - Item 6..N capture discrete run parameters, port bindings, or service profiles from the input.
6. Append the mandatory run skill invocation suffix `[run](file;.cursor/skills/run)`.
7. Output ONLY the resulting formatted prompt as plain markdown text. STRICT — NEVER wrap it in a fenced code block: never emit ```markdown, ```plaintext, or any ``` fence at the start or end of the output. `#` headers and inline markdown are correct; the fence around them is forbidden. No conversational filler, no commentary before or after it.
8. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.
9. Additional Instructions Mandate: Always append the relative paths directive under ## Additional Instructions: '- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.'
10. Step Budget Mandate: Execute within the 400-step lifecycle budget, dedicating appropriate phases to configuration inspection, environment installation, and runner execution.
11. Heatmap Mandate: Track in-flight file touches, prevent write collisions, and query the cognitive activity graph via GitMap Heatmap (gitmap heatmap discuss / log / live) with transparent 3-tier database location reporting.

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# <task title>: high priority instruction, non-negotiable task

${Input Text Verbatim}

## task: <task title>
## slug: <task-slug>

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Inspect run.config.json for target service configuration and port mappings
3. Execute run script (./run.ps1 or ./run.sh) with automatic dependency verification
4. If dependencies or compilers are missing, fall back to gitmap aum install or ./local-install.ps1, then auto-rerun
5. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
6. [Discrete run directive extracted from input]

## Must follow and spawn agent using

[run](file;.cursor/skills/run)

## Additional Instructions

- Execute within the 400-step lifecycle budget: allocate planning, installation, and run execution phases systematically.
- [/plan](slashCommand;plan) first before doing the work to reduce the credits.
- [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap Heatmap (`gitmap heatmap discuss` / `gitmap heatmap log` / `gitmap heatmap live`) with transparent 3-tier database location reporting.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
