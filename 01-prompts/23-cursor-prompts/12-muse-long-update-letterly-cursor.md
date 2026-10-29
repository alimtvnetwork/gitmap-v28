# Muse Long Update — Letterly Prompt Formatter

Format input text according to exact high-priority execution template below without conversational filler or commentary.

1. Clean input text verbatim, stripping verbal fillers (`um`, `ah`, `uh`, `like`, leading `Okay`, excess `so`) while preserving technical details, flags, and file paths.
2. Derive `<task title>` as task subject (e.g. `SEO Writing and Folder Structure Instructions`). Output line 1: `# <task title>: high priority instruction, non-negotiable task`.
3. Output `${Input Text Verbatim}` directly beneath title.
4. Output subheaders: `## task: <task title>` and `## slug: <task-slug>` (lowercase-hyphenated).
5. Under `# Actionable Items Must Follow Non-Negotiable`, ensure:
   - Item 1: `1. Always update GitMap FIRST via \`gitmap update\` before starting work.`
   - Item 2: `2. Run workspace pull via \`gitmap pa\`. If run recently (<15m–1h), skip and run \`gitmap pull\` on the working repository.`
   - Item 3: `3. Write spec under 02-spec/21-app/XX-<slug>/ and enqueue plan in .ai-memory/plans/pending/XX-<slug>.md first (NO-BUILD, NO-TEST during planning).`
   - Item 4: `4. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat); TOTAL BAN on rg, ripgrep, grep, Select-String.`
   - Item 5: `5. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); no absolute paths or file:/// URIs.`
   - Items 6..N: Technical directives extracted from input.
   - Penultimate/Final (if requested): `Verify CI via gitmap pe -t until green` and `Execute minor version bump release ceremony via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"`.
6. Append the Muse Overnight Autopilot Checklist, agent suffix link, and Additional Instructions.
7. Output plain markdown only. Never wrap in fenced code blocks (no ```markdown or ``` fences).
8. Relative Paths: Only add relative paths (`02-spec/...`, `.ai-memory/...`, `cmd/...`); never add absolute paths.

${Input Text Verbatim} = Cleaned input text without verbal filler words.

Output Format:

# <task title>: high priority instruction, non-negotiable task

${Input Text Verbatim}

## task: <task title>
## slug: <task-slug>

# Actionable Items Must Follow Non-Negotiable

1. Always update GitMap FIRST — run `gitmap update` before starting any work.
2. Workspace pull via `gitmap pa` — if run recently (<15m–1h), skip `gitmap pa` and run `gitmap pull` on the target repository.
3. Write spec under 02-spec/21-app/XX-<slug>/ and enqueue plan in .ai-memory/plans/pending/XX-<slug>.md first (NO-BUILD, NO-TEST during planning).
4. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py); TOTAL BAN on rg, ripgrep, grep, Select-String.
5. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); never add absolute paths.
6. ....discrete technical directives extracted from input direction....

### Muse Overnight Autopilot Checklist
- [ ] 1. GitMap Updated (`gitmap update`) & Fleet/Repo Pulled (`gitmap pa` or `gitmap pull`)
- [ ] 2. Autopilot Ambiguity Resolution: Pick closest sensible solution without halting; log to `.ai-memory/ambiguity/` for morning review
- [ ] 3. Batch Questions: Record any questions for morning review; do not stop execution
- [ ] 4. Pragmatic Commits via GitMap (`gitmap cpf` / `gitmap cpc` / `gitmap commit`): Group files logically; no single-file commits; no "perfect commit" paralysis; intermediate code does not need to build
- [ ] 5. Zero Intermediate Builds/Tests: Intermediate build/test checks banned; verify build ONLY at end if CI/CD or test explicitly requested
- [ ] 6. Conditional CI/CD & Minor Release: If CI fix asked, verify `gitmap pe -t` to green and bump minor release
- [ ] 7. Upstream Push Complete: All modified files committed and pushed (`gitmap push`) — No Push = Not Done

## Must follow and spawn agent using

[execute-parent-task-with-n-steps-v6](file;.cursor/skills/execute-parent-task-with-n-steps-v6)

## Additional Instructions

- [/plan](slashCommand;plan) first before doing work for improved work or less friction, faster work and more pure and effective work.
- [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap RC Heatmap (`gitmap rc heatmap status` / `gitmap rc heatmap live`).
- Only add relative paths, never add absolute paths during work; this must be respected on release page and in release notes as well.
