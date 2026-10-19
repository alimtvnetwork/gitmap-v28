# CI/CD Fix & Release Mode (Cursor) — Letterly Prompt Formatter

Format whatever input text is provided according to the exact CI/CD fix and minor release template below, following the execute N-steps structure, using the Cursor skill format. Do NOT add conversational filler or commentary.

1. Clean the input text verbatim while capturing all failing workflow names, errors, and reproduction steps.
2. Derive `<task title>` from the input text — it is the SUBJECT of the task (e.g. `CI/CD Pipeline Fix and Minor Release`). Structure the output starting immediately with `# <task title>: high priority instruction, non-negotiable task`.
3. Put `${Input Text Verbatim}` directly beneath the title line.
4. After the verbatim input, output both the task name and the slug as subheaders:
   - `## task: <task title>` — the human-readable subject/name of the task.
   - `## slug: <task-slug>` — the lowercase-hyphenated task slug (e.g. `## slug: cicd-pipeline-fix-and-minor-release`), NEVER `## slug: high-priority-instruction`.
   Both the task name and slug are required so AI agents can identify and register the task via GitMap (`gitmap agent task init --title "<task title>" --slug "<task-slug>"`).
5. Construct `# Actionable Items Must Follow Non-Negotiable`:
   - Item 1 is ALWAYS: `1. Git pull BEFORE starting any work — first try \`gitmap pull <repo path>\` (or run \`gitmap pull\` in repository folder); if that fails, fall back to regular \`git pull\` and resolve any merge conflicts so the repository is completely up to date and clean before starting.`
   - Item 2 is ALWAYS: `2. Write spec and plan first under 02-spec/21-app/<slug>/ and .ai-memory/plans/<slug>.md`
   - Item 3 is ALWAYS: `3. Inspect live CI/CD pipeline errors and execution timeline using gitmap pe -t (or gitmap pipeline fix)`
   - Item 4 is ALWAYS: `4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well`
   - Item 5: `5. Perform grounded 4-part Root Cause Analysis (RCA) on exact failing step and log`
   - Item 6: `6. Apply surgical code fixes directly resolving the root cause without disabling any CI checks`
   - Item 7: `7. Verify fixes locally with targeted file linters (05-guideline-autofixer.py, check-prompts-loaded.py)`
   - Item 8: `8. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"`
   - Item 9: `9. Commit atomically via gitmap cpf, tag release, and loop until CI/CD is completely green`
6. Append the mandatory skill invocation suffix `[ci-cd-fix-gitmap-release](file;.cursor/skills/ci-cd-fix-gitmap-release)`.
7. Output ONLY the resulting formatted prompt as plain markdown text. STRICT — NEVER wrap it in a fenced code block: never emit ```markdown, ```plaintext, or any ``` fence at the start or end of the output. `#` headers and inline markdown are correct; the fence around them is forbidden. No conversational filler, no commentary before or after it.
8. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.
9. Additional Instructions Mandate: Always append the relative paths directive under ## Additional Instructions: '- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.'
10. Step Budget Mandate: Execute within the 400-step lifecycle budget (Phase 1: 100 plan/diagnosis, Phase 2: 150 intermediate fixes/rca, Phase 3: 150 verification/bump/push).
11. Heatmap Mandate: Track in-flight file touches, prevent write collisions, and query the cognitive activity graph via GitMap Heatmap (gitmap heatmap discuss / log / live) with transparent 3-tier database location reporting.

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# <task title>: high priority instruction, non-negotiable task

${Input Text Verbatim}

## task: <task title>
## slug: <task-slug>

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec and plan first under 02-spec/21-app/<slug>/ and .ai-memory/plans/<slug>.md
3. Inspect live CI/CD pipeline errors and execution timeline using gitmap pe -t (or gitmap pipeline fix)
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Perform grounded 4-part Root Cause Analysis (RCA) on exact failing step and log
6. Apply surgical code fixes directly resolving the root cause without disabling any CI checks
7. Verify fixes locally with targeted file linters (05-guideline-autofixer.py, check-prompts-loaded.py)
8. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"
9. Commit atomically via gitmap cpf, tag release, and loop until CI/CD is completely green

## Must follow and spawn agent using

[ci-cd-fix-gitmap-release](file;.cursor/skills/ci-cd-fix-gitmap-release)

## Additional Instructions

- Execute within the 400-step lifecycle budget: 100 plan, 150 intermediate/archival, 150 implementation/sync.
- [/plan](slashCommand;plan) first before doing the work to reduce the credits.
- [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap live telemetry, pipeline self-healing, and toolchain caching.
- Track in-flight file touches and prevent write collisions using GitMap Heatmap (`gitmap heatmap discuss` / `gitmap heatmap log` / `gitmap heatmap live`) with transparent 3-tier database location reporting.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
