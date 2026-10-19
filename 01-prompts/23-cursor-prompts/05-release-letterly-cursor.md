# Release Mode (Cursor) — Letterly Prompt Formatter

Format whatever input text is provided according to the exact minor release template below, following the execute N-steps structure, using the Cursor skill format. Do NOT add conversational filler or commentary.

1. Clean the input text verbatim while strictly capturing version scope, changelog notes, and release constraints.
2. Derive `<task title>` from the input text — it is the SUBJECT of the task (e.g. `Release Minor Version`). Structure the output starting immediately with `# <task title>: high priority instruction, non-negotiable task`.
3. Put `${Input Text Verbatim}` directly beneath the title line.
4. After the verbatim input, output both the task name and the slug as subheaders:
   - `## task: <task title>` — the human-readable subject/name of the task.
   - `## slug: <task-slug>` — the lowercase-hyphenated task slug (e.g. `## slug: release-minor-version`), NEVER `## slug: high-priority-instruction`.
   Both the task name and slug are required so AI agents can identify and register the task via GitMap (`gitmap agent task init --title "<task title>" --slug "<task-slug>"`).
5. Construct `# Actionable Items Must Follow Non-Negotiable`:
   - Item 1 is ALWAYS: `1. Git pull BEFORE starting any work — first try \`gitmap pull <repo path>\` (or run \`gitmap pull\` in repository folder); if that fails, fall back to regular \`git pull\` and resolve any merge conflicts so the repository is completely up to date and clean before starting.`
   - Item 2 is ALWAYS: `2. Write spec and plan first adhering to 02-spec/16-generic-release/ and 01-prompts/17-release-management/02-minor-bump.md`
   - Item 3 is ALWAYS: `3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String`
   - Item 4 is ALWAYS: `4. Strictly use relative Git paths only; only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page, in release notes, changelog.md, manifests, and documentation`
   - Item 5 is ALWAYS: `5. Enforce zero-storage GitHub Actions rules (zero routine artifact uploads)`
   - Item 6 is ALWAYS: `6. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"`
   - Item 7 is ALWAYS: `7. Consolidate and update release notes in root changelog.md and manifests`
   - Item 8 is ALWAYS: `8. Commit atomically via gitmap cpf "<module> - release minor version" and push release tag to remote tracking branch`
6. Append the mandatory release skill invocation suffix `[minor-bump](file;.cursor/skills/minor-bump)`.
7. Output ONLY the resulting formatted prompt as plain markdown text. STRICT — NEVER wrap it in a fenced code block: never emit ```markdown, ```plaintext, or any ``` fence at the start or end of the output. `#` headers and inline markdown are correct; the fence around them is forbidden. No conversational filler, no commentary before or after it.
8. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page, in release notes, changelog.md, manifests, and documentation.
9. Additional Instructions Mandate: Always append the relative paths directive under ## Additional Instructions: '- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.'
10. Step Budget Mandate: Execute within the 400-step lifecycle budget, dedicating Phase 3 (Steps 391..400) to the release ceremony, changelog synchronization, atomic commit, and upstream push.
11. Heatmap Mandate: Track in-flight file touches, prevent write collisions, and query the cognitive activity graph via GitMap Heatmap (gitmap heatmap discuss / log / live) with transparent 3-tier database location reporting.

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# <task title>: high priority instruction, non-negotiable task

${Input Text Verbatim}

## task: <task title>
## slug: <task-slug>

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec and plan first adhering to 02-spec/16-generic-release/ and 01-prompts/17-release-management/02-minor-bump.md
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths only; only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page, in release notes, changelog.md, manifests, and documentation
5. Enforce zero-storage GitHub Actions rules (zero routine artifact uploads)
6. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"
7. Consolidate and update release notes in root changelog.md and manifests
8. Commit atomically via gitmap cpf "<module> - release minor version" and push release tag to remote tracking branch

## Must follow and spawn agent using

[minor-bump](file;.cursor/skills/minor-bump)

## Additional Instructions

- Execute within the 400-step lifecycle budget: Phase 3 release ceremony (Steps 391..400) for minor version bump, changelog update, and atomic push.
- [/plan](slashCommand;plan) first before doing the work to reduce the credits.
- [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap Heatmap (`gitmap heatmap discuss` / `gitmap heatmap log` / `gitmap heatmap live`) with transparent 3-tier database location reporting.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
