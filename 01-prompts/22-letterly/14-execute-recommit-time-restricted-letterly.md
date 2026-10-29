# Recommit & Time-Restricted Execution — Letterly Prompt Formatter

Format input text according to exact output template below without conversational filler or commentary. This formatter directs the AI agent to follow the unified V7 master execution skill (`execute-parent-task-with-n-steps-v7`), which encapsulates both Option 1 (Time-Restricted with timezone & weekday < 08:00 AM cutoff) and Option 2 (Unrestricted Time) alongside AI automation bot integration.

1. Clean input text verbatim, removing verbal fillers (`um`, `ah`, `uh`, `like`, leading `Okay`, excess `so`) while strictly preserving technical details, requirements, file paths, commands, flags, and directives.
2. Derive `<task title>` as task subject (e.g. `Time-Restricted Recommit and Bot Workflow`). Output line 1: `# <task title>: high priority instruction, non-negotiable task`.
3. Output `${Input Text Verbatim}` directly beneath title.
4. Output task name and slug subheaders:
   - `## task: <task title>`
   - `## slug: <task-slug>` (lowercase-hyphenated; when authoring directories and plan files, format as XX-<task-slug> where XX is the zero-padded sequence of that folder, e.g. `## slug: 25-time-restricted-recommit-and-bot-workflow`).
5. Under `# Actionable Items Must Follow Non-Negotiable`, ensure:
   - Item 1: `1. Git pull BEFORE starting any work — first try \`gitmap pull <repo path>\` (or run \`gitmap pull\` in repository folder); if that fails, fall back to regular \`git pull\` and resolve any merge conflicts so the repository is completely up to date and clean before starting.`
   - Item 2: `2. Write spec under 02-spec/21-app/XX-<slug>/ and enqueue plan task in .ai-memory/plans/pending/XX-<slug>.md (subtasks in .ai-memory/plans/subtasks/XX-<slug>/) first, (xx should be figured out by the folder structure and running AI tools, xx is sequence 01, 02)`
   - Item 3: `3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String`
   - Item 4: `4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well`
   - Item 5: ``5. Use `gitmap` AI agents to enter data``
   - Item 6: `6. Task completion includes committing and pushing to Git`
   - Item 7: ``7. Give every Task-NN in the confirmed task breakdown a `Slug:` sub-item derived from the root slug (`<root-slug> - Task NN`) so GitMap can create and verify subtasks under the root task``
   - Item 8: `8. Execute via V7 Master Engine with Dual-Option Support: Follow [execute-parent-task-with-n-steps-v7](file;.agents/skills/execute-parent-task-with-n-steps-v7) which provides two explicit execution options: Option 1 (Restricted Time: enforce commits strictly BEFORE 08:00 AM on weekdays and anytime on weekends in target timezone via \`--tz\`, \`--schedule\`, \`--auto-window\`, \`--non-interactive\`) or Option 2 (Unrestricted Time: immediate commits anytime). When automation bot persona is commanded, configure identity via \`gitmap autobot\` before modifying code.`
   - Item 9: ``9. Recommit & History Adjustment (On-Call): When recommit is requested, inspect git log, verify previous commit hash, and execute atomic recommits via `gitmap` semantic commit verbs (`gitmap cpf`, `gitmap cbf`, `gitmap cfr`, `gitmap cpc`, `gitmap cpr`) without polluting git history.``
   followed by discrete technical action items extracted from input text as items 10..N.
6. End with mandatory agent invocation suffix pointing to `@[.agents/skills/execute-parent-task-with-n-steps-v7]`.
7. Output plain markdown text only. Never wrap in fenced code blocks (no ```markdown or ``` fences).
8. Relative Paths: Only add relative paths (`02-spec/...`, `.ai-memory/...`, `cmd/...`); never add absolute paths.

${Input Text Verbatim} = Cleaned input text without verbal filler words.

Output Format:

# <task title>: high priority instruction, non-negotiable task

${Input Text Verbatim}

## task: <task title>
## slug: <task-slug>  (should not contain `xx-` here. )

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/XX-<slug>/ and enqueue plan task in .ai-memory/plans/pending/XX-<slug>.md (subtasks in .ai-memory/plans/subtasks/XX-<slug>/) first, (xx should be figured out by the folder structure and running AI tools, xx is sequence 01, 02)
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Use `gitmap` AI agents to enter data
6. Task completion includes committing and pushing to Git
7. Give every Task-NN in the confirmed task breakdown a `Slug:` sub-item derived from the root slug (`<root-slug> - Task NN`) so GitMap can create and verify subtasks under the root task
8. Execute via V7 Master Engine with Dual-Option Support: Follow [execute-parent-task-with-n-steps-v7](file;.agents/skills/execute-parent-task-with-n-steps-v7) which provides two explicit execution options: Option 1 (Restricted Time: enforce commits strictly BEFORE 08:00 AM on weekdays and anytime on weekends in target timezone via `--tz`, `--schedule`, `--auto-window`, `--non-interactive`) or Option 2 (Unrestricted Time: immediate commits anytime). When automation bot persona is commanded, configure identity via `gitmap autobot` before modifying code.
9. Recommit & History Adjustment (On-Call): When recommit is requested, inspect git log, verify previous commit hash, and execute atomic recommits via `gitmap` semantic commit verbs (`gitmap cpf`, `gitmap cbf`, `gitmap cfr`, `gitmap cpc`, `gitmap cpr`) without polluting git history.
10. [extracted actionable technical directive from input]

## Must follow and spawn agent using

@[.agents/skills/execute-parent-task-with-n-steps-v7]

## Additional Instructions

- [/plan](slashCommand;plan) first before doing work for improved work or less friction, faster execution for AI and coding without errors and following Coding Guidelines.
- [/learn](slashCommand;learn) from [gitmap](file;.agents/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap RC Heatmap (`gitmap rc heatmap status` / `gitmap rc heatmap tree` / `gitmap rc heatmap live` / `gitmap rc heatmap sync`).
- Only add relative paths, never add absolute paths during work; this must be respected on release page and in release notes as well.
