# Run Mode — Letterly Prompt Formatter

Format input text according to exact project run template below without conversational filler or commentary.

1. Clean input text verbatim, strictly preserving target service names, ports, run flags, and environment parameters.
2. Derive `<task title>` as task subject (e.g. `Run Application Services`). Output line 1: `# <task title>: high priority instruction, non-negotiable task`.
3. Output `${Input Text Verbatim}` directly beneath title.
4. Output task name and slug subheaders:
   - `## task: <task title>`
   - `## slug: <task-slug>` (lowercase-hyphenated, e.g. `## slug: run-application-services`).
5. Under `# Actionable Items Must Follow Non-Negotiable`, ensure:
   - Item 1: `1. Git pull BEFORE starting any work — first try \`gitmap pull <repo path>\` (or run \`gitmap pull\` in repository folder); if that fails, fall back to regular \`git pull\` and resolve any merge conflicts so the repository is completely up to date and clean before starting.`
   - Item 2: `2. Inspect run.config.json for target service configuration and port mappings`
   - Item 3: `3. Execute run script (./run.ps1 or ./run.sh) with automatic dependency verification`
   - Item 4: `4. If dependencies or compilers are missing, fall back to gitmap aum install or ./local-install.ps1, then auto-rerun`
   - Item 5: `5. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well`
   followed by discrete run parameters or service profiles from input as items 6..N.
6. Append mandatory run skill invocation suffix `[run](file;.cursor/skills/run)`.
7. Output plain markdown text only. Never wrap in fenced code blocks (no ```markdown or ``` fences).
8. Relative Paths: Only add relative paths (`02-spec/...`, `.ai-memory/...`, `cmd/...`); never add absolute paths.

${Input Text Verbatim} = Cleaned input text without verbal filler words.

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

- [/plan](slashCommand;plan) first before doing work for improved work or less friction, faster work and more pure and effective work.
- [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap RC Heatmap (`gitmap rc heatmap status` / `gitmap rc heatmap tree` / `gitmap rc heatmap live` / `gitmap rc heatmap sync`).
- Only add relative paths, never add absolute paths during work; this must be respected on release page and in release notes as well.
