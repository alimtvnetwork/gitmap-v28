# Recommit & Time-Restricted Execution — Letterly Prompt Formatter

Format input text according to exact output template below without conversational filler or commentary. Directs the AI agent to follow unified V7 master execution per @[execute-parent-task-with-n-steps-v7](file;.cursor/skills/execute-parent-task-with-n-steps-v7) (Option 1: Time-Restricted; Option 2: Unrestricted) and bot integration.

1. Clean input text verbatim, removing verbal fillers (`um`, `ah`, `uh`, `like`, leading `Okay`, excess `so`) while strictly preserving technical details, requirements, file paths, commands, flags, and directives.
2. Derive `<task title>` as task subject (e.g. `Time-Restricted Recommit and Bot Workflow`). Output line 1: `# <task title>: high priority instruction, non-negotiable task`.
3. Output `${Input Text Verbatim}` directly beneath title.
4. Output task name and slug subheaders:
   - `## task: <task title>`
   - `## slug: <task-slug>` (lowercase-hyphenated, formatted as XX-<task-slug> for directories and plan files, e.g. `## slug: 25-time-restricted-recommit-and-bot-workflow`).
5. Construct `# Actionable Items Must Follow Non-Negotiable`: include fixed items 1-9 (pull first, spec/plan enqueue with planning decoupled from execution, GitMap search primacy, relative Git paths, agent data entry, commit and push, task slug naming, V7 dual-option execution, and time-restricted recommits), followed by discrete technical directives from input as items 10..N.
6. End with mandatory agent invocation suffix pointing to @[execute-parent-task-with-n-steps-v7](file;.cursor/skills/execute-parent-task-with-n-steps-v7).
7. Output plain markdown text only. Never wrap in fenced code blocks (no ```markdown or ``` fences).
8. Relative Paths: Only add relative paths (`02-spec/...`, `.ai-memory/...`, `cmd/...`); never add absolute paths.

${Input Text Verbatim} = Cleaned input text without verbal filler words.

Output Format:

# <task title>: high priority instruction, non-negotiable task

${Input Text Verbatim}

## task: <task title>
## slug: <task-slug>  (should not contain `xx-` here. )

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve merge conflicts so the repository is up to date and clean before starting.
2. Write spec under 02-spec/21-app/XX-<slug>/ and enqueue plan task in .ai-memory/plans/pending/XX-<slug>.md (subtasks in .ai-memory/plans/subtasks/XX-<slug>/) first (XX is sequence 01, 02)
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add relative paths, never add absolute paths during work or on the release page and in release notes
5. Use `gitmap` AI agents to enter data
6. Task completion includes committing and pushing to Git
7. Give every Task-NN in confirmed task breakdown a `Slug:` sub-item derived from root slug (`<root-slug> - Task NN`) so GitMap can create and verify subtasks
8. Execute via V7 Master Engine per @[execute-parent-task-with-n-steps-v7](file;.cursor/skills/execute-parent-task-with-n-steps-v7): Option 1 (Time-Restricted commits before 08:00 AM weekdays/anytime weekends via `--tz`, `--schedule`, `--auto-window`) or Option 2 (Unrestricted commits); configure bot identity via `gitmap autobot` when commanded.
9. Enforce time-restricted recommits per @[gitmap-time-restriction](file;.cursor/skills/gitmap-time-restriction) using `gitmap` semantic commit verbs without polluting git history.
10. [extracted actionable technical directive from input]

## Must follow and spawn agent using

@[execute-parent-task-with-n-steps-v7](file;.cursor/skills/execute-parent-task-with-n-steps-v7)

## Additional Instructions

- [/plan](slashCommand;plan) first before doing work for improved work or less friction, faster work and more pure and effective work.
- [/learn](slashCommand;learn) from @[gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap RC Heatmap (`gitmap rc heatmap status` / `gitmap rc heatmap tree` / `gitmap rc heatmap live` / `gitmap rc heatmap sync`).
- Only add relative paths, never add absolute paths during work; this must be respected on release page and in release notes as well.
