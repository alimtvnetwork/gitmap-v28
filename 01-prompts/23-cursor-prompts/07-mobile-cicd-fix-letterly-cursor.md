# Mobile CI/CD Fix Mode (Cursor) — Letterly Prompt Formatter

Combine the entire output into exactly ONE continuous single-line paragraph with zero newlines, zero line gaps, and zero conversational filler, using the Cursor skill format.

1. Clean the input text verbatim without conversational filler words (`um`, `ah`, `uh`).
2. Derive `<task title>` from the input text — it is the SUBJECT of the task (e.g. `CI/CD Pipeline Fix`). NEVER use the literal words `High Priority Instruction` as the title. Derive `<task-slug>` as the lowercase-hyphenated task title (e.g. `cicd-pipeline-fix`).
3. Extract discrete understandable tasks from the input text as an ordered list: `1. diagnose CI/CD via gitmap pe -t, 2. apply 4-part RCA fix, 3. verify locally, 4. atomic commit via gitmap cpf, 5. minor release`.
4. Prefix the single line with slash commands: `[/goal](slashCommand;goal) [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) [/plan](slashCommand;plan) `.
5. Follow immediately with the chained title header: `# <task title> - high priority instruction, non-negotiable task (task: <task title>, slug: <task-slug>): `.
6. Immediately follow the colon with `${Input Text Verbatim}; `.
7. Follow with core actionable directives: `first try gitmap pull <repo-path> (fallback git pull) and resolve conflicts before starting; execute within 400-step lifecycle budget (100 plan, 150 intermediate/archival, 150 implementation/sync); run gitmap pe -t to diagnose live CI/CD errors, fix all pipeline failures via 4-part RCA, track file touches and prevent write collisions using GitMap Heatmap (gitmap heatmap discuss / log / live); verify locally, commit atomically via gitmap cpf, and release minor update. `.
8. Append the extracted task list: `discrete understandable tasks to execute: <extracted ordered tasks>`.
9. Conclude the single line with: ` - strictly use relative git paths only; only add the relative paths, never add the absolute path during your work (this must be respected on the release page and in release notes as well), must follow the skill [ci-cd-fix-gitmap-release](file;.cursor/skills/ci-cd-fix-gitmap-release)`.
10. Output exactly that single continuous paragraph with ZERO newlines, ZERO line breaks, and NO leading "Output" text.
11. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.

${Input Text Verbatim} = The cleaned input text as it is, without filler words.

Output Format:
[/goal](slashCommand;goal) [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) [/plan](slashCommand;plan) # <task title> - high priority instruction, non-negotiable task (task: <task title>, slug: <task-slug>): ${Input Text Verbatim}; first try gitmap pull <repo-path> (fallback git pull) and resolve conflicts before starting; execute within 400-step lifecycle budget (100 plan, 150 intermediate/archival, 150 implementation/sync); run gitmap pe -t to diagnose live CI/CD errors, fix all pipeline failures via 4-part RCA, track file touches and prevent write collisions using GitMap Heatmap (gitmap heatmap discuss / log / live); verify locally, commit atomically via gitmap cpf, and release minor update. discrete understandable tasks to execute: 1. diagnose CI/CD via gitmap pe -t, 2. apply 4-part RCA fix, 3. verify locally, 4. atomic commit via gitmap cpf, 5. minor release - strictly use relative git paths only; only add the relative paths, never add the absolute path during your work (this must be respected on the release page and in release notes as well), must follow the skill [ci-cd-fix-gitmap-release](file;.cursor/skills/ci-cd-fix-gitmap-release)
