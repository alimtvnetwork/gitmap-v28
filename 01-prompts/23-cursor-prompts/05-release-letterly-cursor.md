# Release Mode — Letterly Prompt Formatter

Format input text according to exact minor release template below without conversational filler or commentary.

1. Clean input text verbatim, stripping verbal fillers while strictly preserving version scope, changelog notes, and release constraints.
2. Derive `<task title>` as task subject (e.g. `Release Minor Version`). Output line 1: `# <task title>: high priority instruction, non-negotiable task`.
3. Output `${Input Text Verbatim}` directly beneath title.
4. Output task name and slug subheaders:
   - `## task: <task title>`
   - `## slug: <task-slug>` (lowercase-hyphenated, e.g. `## slug: release-minor-version`).
5. Under `# Actionable Items Must Follow Non-Negotiable`, ensure:
   - Item 1: `1. Git pull BEFORE starting any work — first try \`gitmap pull <repo path>\` (or run \`gitmap pull\` in repository folder); if that fails, fall back to regular \`git pull\` and resolve any merge conflicts so the repository is completely up to date and clean before starting.`
   - Item 2: `2. Write spec and plan first adhering to 02-spec/16-generic-release/ and 01-prompts/17-release-management/02-minor-bump.md`
   - Item 3: `3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String`
   - Item 4: `4. Strictly use relative Git paths only; only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page, in release notes, changelog.md, manifests, and documentation`
   - Item 5: `5. Enforce zero-storage GitHub Actions rules (zero routine artifact uploads)`
   - Item 6: `6. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"`
   - Item 7: `7. Consolidate and update release notes in root changelog.md and manifests`
   - Item 8: `8. Commit atomically via gitmap cpr "v<version> <summary>" and push release tag to remote tracking branch`
6. Append mandatory release skill invocation suffix `[minor-bump](file;.cursor/skills/minor-bump)`.
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
2. Write spec and plan first adhering to 02-spec/16-generic-release/ and 01-prompts/17-release-management/02-minor-bump.md
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths only; only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page, in release notes, changelog.md, manifests, and documentation
5. Enforce zero-storage GitHub Actions rules (zero routine artifact uploads)
6. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"
7. Consolidate and update release notes in root changelog.md and manifests
8. Commit atomically via gitmap cpr "v<version> <summary>" and push release tag to remote tracking branch

## Must follow and spawn agent using

[minor-bump](file;.cursor/skills/minor-bump)

## Additional Instructions

- [/plan](slashCommand;plan) first before doing work for improved work or less friction, faster work and more pure and effective work.
- [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap RC Heatmap (`gitmap rc heatmap status` / `gitmap rc heatmap tree` / `gitmap rc heatmap live` / `gitmap rc heatmap sync`).
- Only add relative paths, never add absolute paths during work; this must be respected on release page and in release notes as well.
