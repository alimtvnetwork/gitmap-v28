# Release Mode (Cursor) — Letterly Prompt Formatter

Format whatever input text is provided according to the exact minor release template below, following the execute N-steps structure and using the Cursor skill format. Do NOT add conversational filler or commentary.

1. Clean the input text verbatim while strictly capturing version scope, changelog notes, and release constraints.
2. Structure the output starting immediately with `# High Priority Instruction`.
3. Put `${Input Text Verbatim}` directly beneath the high priority header.
4. Construct `# Actionable Items Must Follow Non-Negotiable`:
   - Item 1 is ALWAYS: `1. Write spec and plan first adhering to 02-spec/16-generic-release/ and 01-prompts/17-release-management/02-minor-bump.md`
   - Item 2 is ALWAYS: `2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String`
   - Item 3 is ALWAYS: `3. Strictly use relative Git paths only; only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page, in release notes, changelog.md, manifests, and documentation`
   - Item 4 is ALWAYS: `4. Enforce zero-storage GitHub Actions rules (zero routine artifact uploads)`
   - Item 5 is ALWAYS: `5. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"`
   - Item 6 is ALWAYS: `6. Consolidate and update release notes in root changelog.md and manifests`
   - Item 7 is ALWAYS: `7. Commit atomically via gitmap cpf "<module> - release minor version" and push release tag to remote tracking branch`
5. Append the mandatory release skill invocation suffix `[minor-bump](file;.cursor/skills/minor-bump)`.
6. Output ONLY the resulting formatted markdown block.
7. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page, in release notes, changelog.md, manifests, and documentation.
8. Additional Instructions Mandate: Always append the relative paths directive under ## Additional Instructions: '- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.'

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# High Priority Instruction

${Input Text Verbatim}

# Actionable Items Must Follow Non-Negotiable

1. Write spec and plan first adhering to 02-spec/16-generic-release/ and 01-prompts/17-release-management/02-minor-bump.md
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths only; only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page, in release notes, changelog.md, manifests, and documentation
4. Enforce zero-storage GitHub Actions rules (zero routine artifact uploads)
5. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"
6. Consolidate and update release notes in root changelog.md and manifests
7. Commit atomically via gitmap cpf "<module> - release minor version" and push release tag to remote tracking branch

## Must follow and spawn agent using

[minor-bump](file;.cursor/skills/minor-bump)

## Additional Instructions

- [/plan](slashCommand;plan) first before doing the work to reduce the credits.
- [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
