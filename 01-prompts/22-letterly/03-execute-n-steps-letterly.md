# Execute N-Steps — Letterly Prompt Formatter

Format whatever input text is provided according to the exact high-priority execution template below. Do NOT add conversational filler or commentary (never write "Certainly! Here is your output:").

1. Capture and clean the input text verbatim, stripping verbal filler words (`um`, `ah`, `uh`, `like`, don't start with `Okay`, reduce if there are too many `so`) while preserving every technical directive, parameter, flag, and file path.
2. Derive `<task title>` from the input text — it is the SUBJECT of the task (e.g. `SEO Writing and Folder Structure Instructions`). Structure the output starting immediately with `# <task title>: high priority instruction, non-negotiable task`.
3. Output `${Input Text Verbatim}` directly beneath the title line.
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
6. Append the mandatory agent invocation suffix pointing to `[execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)`.
7. Make sure all the action items are listed and nothing pending.
8. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.
9. Additional Instructions Mandate: Always append the relative paths directive under ## Additional Instructions: '- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.'
10. Step Budget Mandate: Execute within the 400-step lifecycle budget (Phase 1: 100 plan, Phase 2: 150 intermediate/archival, Phase 3: 150 implementation/sync).
11. Don't output "```plaintext or ```markdown" ever please.
12. Heatmap & Agent Task Mandate: Track in-flight file touches, prevent write collisions, and query the cognitive activity graph via GitMap RC Heatmap (gitmap rc heatmap status / tables / query / tree / live / sync). Ensure automatic RC vault creation if not present, link the execution roadmap to heatmap telemetry, and verify task existence and prompt MD5 hash to prevent redundant reads and enable instant dynamic subagent join.

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
5. ....other steps and more steps sequentially from the input direction. Create more steps in between.

## Must follow and spawn agent using

[execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)

## Additional Instructions

- Execute within the 400-step lifecycle budget: Phase 1 (Steps 1..100) planning and spec authoring, Phase 2 (Steps 101..250) sanitization, archival, and prerequisites, Phase 3 (Steps 251..400) worker swarm implementation, targeted checks, fleet synchronization, and atomic push.
- [/plan](slashCommand;plan) first before doing the work to reduce the credits.
- [/learn](slashCommand;learn) from [gitmap](file;.agents/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches, prevent write collisions, and query the cognitive graph via GitMap RC Heatmap (`gitmap rc heatmap status` / `gitmap rc heatmap tree` / `gitmap rc heatmap live` / `gitmap rc heatmap sync`). If RC does not exist, auto-provision it linked to the mission folder; connect the roadmap to the heatmap command so newly spawned subagents can join dynamically without redundant reads; verify prior task existence and prompt MD5 hash code before initialization.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
