# Mobile Mode — Letterly Prompt Formatter

Format input text according to exact mobile single-line template below without conversational filler or commentary.

1. Clean input text verbatim, stripping verbal fillers (`um`, `ah`, `uh`, `like`, leading `Okay`, excess `so`) while strictly preserving technical details, file paths, commands, and requirements.
2. Derive `<task title>` as task subject (e.g. `SEO Writing and Folder Structure Instructions`). Derive `<task-slug>` as lowercase-hyphenated task title (e.g. `seo-writing-and-folder-structure-instructions`).
3. Extract discrete understandable tasks from input text as ordered list: `1. [task-1], 2. [task-2], 3. [task-3]`.
4. Prefix single line with slash commands: `[/goal](slashCommand;goal) [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) [/plan](slashCommand;plan) `.
5. Follow immediately with chained title header: `# <task title> - high priority instruction, non-negotiable task (task: <task title>, slug: <task-slug>): `.
6. Immediately follow colon with `${Input Text Verbatim}; `.
7. Follow with core actionable directives: `do gitmap pull <repo path> (fallback git pull) and resolve conflicts before starting; write spec under 02-spec/21-app/XX-<slug>/ and enqueue plan in .ai-memory/plans/pending/XX-<slug>.md (subtasks in .ai-memory/plans/subtasks/XX-<slug>/) first (xx is sequence 01, 02); planning is strictly dedicated to architecture and decoupled from execution (no-build, no-test, no-code-changes); search codebase exclusively via GitMap; track file touches and prevent write collisions using GitMap RC Heatmap (gitmap rc heatmap status / tree / live / sync); strictly use relative git paths only; `.
8. Append extracted task list: `discrete understandable tasks to execute: <extracted ordered tasks>`.
9. Conclude single line with: ` - strictly use relative git paths only; only add relative paths, never add absolute paths during work (this must be respected on release page and in release notes as well); must follow the skill [execute-parent-task-with-n-steps-v6](file;.cursor/skills/execute-parent-task-with-n-steps-v6)`.
10. Output exactly that single continuous paragraph with ZERO newlines and ZERO line breaks.

${Input Text Verbatim} = Cleaned input text without verbal filler words.

Output Format:
[/goal](slashCommand;goal) [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) [/plan](slashCommand;plan) # <task title> - high priority instruction, non-negotiable task (task: <task title>, slug: <task-slug>): ${Input Text Verbatim}; do gitmap pull <repo path> (fallback git pull) and resolve conflicts before starting; write spec under 02-spec/21-app/XX-<slug>/ and enqueue plan in .ai-memory/plans/pending/XX-<slug>.md (subtasks in .ai-memory/plans/subtasks/XX-<slug>/) first (xx is sequence 01, 02); planning is strictly dedicated to architecture and decoupled from execution (no-build, no-test, no-code-changes); search codebase exclusively via GitMap; track file touches and prevent write collisions using GitMap RC Heatmap (gitmap rc heatmap status / tree / live / sync); strictly use relative git paths only; discrete understandable tasks to execute: 1. [task-1], 2. [task-2], 3. [task-3] - strictly use relative git paths only; only add relative paths, never add absolute paths during work (this must be respected on release page and in release notes as well); must follow the skill [execute-parent-task-with-n-steps-v6](file;.cursor/skills/execute-parent-task-with-n-steps-v6)
