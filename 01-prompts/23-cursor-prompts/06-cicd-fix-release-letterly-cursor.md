# CI/CD Fix & Release Mode — Letterly Prompt Formatter

Format input text according to exact CI/CD fix and minor release template below without conversational filler or commentary.

1. Clean input text verbatim, capturing failing workflow names, errors, and reproduction steps.
2. Derive `<task title>` as task subject (e.g. `CI/CD Pipeline Fix and Minor Release`). Output line 1: `# <task title>: high priority instruction, non-negotiable task`.
3. Output `${Input Text Verbatim}` directly beneath title.
4. Output task name and slug subheaders:
   - `## task: <task title>`
   - `## slug: <task-slug>` (lowercase-hyphenated; when authoring directories and plan files, format as XX-<task-slug> where XX is the zero-padded sequence of that folder, e.g. `## slug: 24-cicd-pipeline-fix-and-minor-release`).
5. Under `# Actionable Items Must Follow Non-Negotiable`, ensure:
   - Item 1: `1. Git pull BEFORE starting any work — first try \`gitmap pull <repo path>\` (or run \`gitmap pull\` in repository folder); if that fails, fall back to regular \`git pull\` and resolve any merge conflicts so the repository is completely up to date and clean before starting.`
   - Item 2: `2. Write spec and plan first under 02-spec/21-app/XX-<slug>/ and .ai-memory/plans/pending/XX-<slug>.md (where XX is zero-padded folder sequence). When writing about planning, there is no need to include execute steps — planning is strictly dedicated to architecture, specification, boundaries, and task decomposition (NO-BUILD, NO-TEST, NO-CODE-CHANGES) with execution decoupled to subsequent phases.`
   - Item 3: `3. Inspect live CI/CD pipeline errors and execution timeline using gitmap pe -t (or gitmap pipeline fix)`
   - Item 4: `4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well`
   - Item 5: `5. Perform grounded 4-part Root Cause Analysis (RCA) on exact failing step and log`
   - Item 6: `6. Apply surgical code fixes directly resolving the root cause without disabling any CI checks`
   - Item 7: `7. Verify fixes locally with targeted file linters (05-guideline-autofixer.py, check-prompts-loaded.py)`
   - Item 8: `8. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"`
   - Item 9: `9. Commit atomically via gitmap cpr, tag release, and loop until CI/CD is completely green`
6. Append mandatory skill invocation suffix `[ci-cd-fix-gitmap-release](file;.cursor/skills/ci-cd-fix-gitmap-release)`.
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
2. Write spec under 02-spec/21-app/XX-<slug>/ and enqueue plan task in .ai-memory/plans/pending/XX-<slug>.md (subtasks in .ai-memory/plans/subtasks/XX-<slug>/) first, (xx should be figured out by the folder structure and running AI tools, xx  is sequence 01, 02)
3. Inspect live CI/CD pipeline errors and execution timeline using gitmap pe -t (or gitmap pipeline fix)
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Perform grounded 4-part Root Cause Analysis (RCA) on exact failing step and log
6. Apply surgical code fixes directly resolving the root cause without disabling any CI checks
7. Verify fixes locally with targeted file linters (05-guideline-autofixer.py, check-prompts-loaded.py)
8. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"
9. Commit atomically via gitmap cpr, tag release, and loop until CI/CD is completely green

## Must follow and spawn agent using

[ci-cd-fix-gitmap-release](file;.cursor/skills/ci-cd-fix-gitmap-release)

## Additional Instructions

- [/plan](slashCommand;plan) first before doing work for improved work or less friction, faster work and more pure and effective work.
- [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap live telemetry, pipeline self-healing, and toolchain caching.
- Track in-flight file touches and prevent write collisions using GitMap RC Heatmap (`gitmap rc heatmap status` / `gitmap rc heatmap tree` / `gitmap rc heatmap live` / `gitmap rc heatmap sync`).
- Only add relative paths, never add absolute paths during work; this must be respected on release page and in release notes as well.
