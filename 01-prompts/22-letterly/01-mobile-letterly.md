# Mobile Mode — Letterly Prompt Formatter

Format whatever input text is provided according to the exact professional mobile single-line output template below.

1. Clean the input text verbatim without conversational filler words (`um`, `ah`, `uh`, `like`) while strictly preserving every technical detail, file path, command, and requirement.
2. Derive `<task title>` from the input text — it is the SUBJECT of the task (e.g. `SEO Writing and Folder Structure Instructions`). NEVER use the literal words `High Priority Instruction` as the title. Derive `<task-slug>` as the lowercase-hyphenated task title (e.g. `seo-writing-and-folder-structure-instructions`).
3. Extract discrete understandable tasks from the input text as an ordered list: `1. [task-1], 2. [task-2], 3. [task-3]`.
4. Prefix the single line with slash commands: `[/goal](slashCommand;goal) [/learn](slashCommand;learn) from [gitmap](file;.agents/skills/gitmap) [/plan](slashCommand;plan) `.
5. Follow immediately with the chained title header: `# <task title> - high priority instruction, non-negotiable task (task: <task title>, slug: <task-slug>): `.
6. Immediately follow the colon with `${Input Text Verbatim}; `.
7. Follow with core actionable directives: `do gitmap pull <repo path> (fallback git pull) and resolve conflicts before starting; write spec under 02-spec/21-app/<slug>/ and enqueue plan in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/); execute within 400-step lifecycle budget (100 plan, 150 intermediate/archival, 150 implementation/sync); search codebase exclusively via GitMap; track file touches and prevent write collisions using GitMap Heatmap (gitmap heatmap discuss / log / live); strictly use relative git paths only; `.
8. Append the extracted task list: `discrete understandable tasks to execute: <extracted ordered tasks>`.
9. Conclude the single line with: ` - strictly use relative git paths only; only add the relative paths, never add the absolute path during your work (this must be respected on the release page and in release notes as well); must follow the skill [execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)`.
10. Output exactly that single continuous paragraph with ZERO newlines, ZERO line breaks, and NO leading "Output" text.
11. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.

${Input Text Verbatim} = The cleaned input text as it is, without filler words.

Output Format:
[/goal](slashCommand;goal) [/learn](slashCommand;learn) from [gitmap](file;.agents/skills/gitmap) [/plan](slashCommand;plan) # <task title> - high priority instruction, non-negotiable task (task: <task title>, slug: <task-slug>): ${Input Text Verbatim}; do gitmap pull <repo path> (fallback git pull) and resolve conflicts before starting; write spec under 02-spec/21-app/<slug>/ and enqueue plan in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/); execute within 400-step lifecycle budget (100 plan, 150 intermediate/archival, 150 implementation/sync); search codebase exclusively via GitMap; track file touches and prevent write collisions using GitMap Heatmap (gitmap heatmap discuss / log / live); strictly use relative git paths only; discrete understandable tasks to execute: 1. [task-1], 2. [task-2], 3. [task-3] - strictly use relative git paths only; only add the relative paths, never add the absolute path during your work (this must be respected on the release page and in release notes as well); must follow the skill [execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)
