# CI/CD Fix & Release Mode — Letterly Prompt Formatter

Format whatever input text is provided according to the exact CI/CD fix and minor release template below, following the execute N-steps structure. Do NOT add conversational filler or commentary.

1. Clean the input text verbatim while capturing all failing workflow names, errors, and reproduction steps.
2. Structure the output starting immediately with `# High Priority Instruction`.
3. Put `${Input Text Verbatim}` directly beneath the high priority header.
4. Construct `# Actionable Items Must Follow Non-Negotiable`:
   - Item 1 is ALWAYS: `1. Write spec and plan first`
   - Item 2: `2. Inspect live CI/CD pipeline errors and execution timeline using gitmap pe -t (or gitmap pipeline fix)`
   - Item 3 is ALWAYS: `3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well`
   - Item 4: `4. Perform grounded 4-part Root Cause Analysis (RCA) on exact failing step and log`
   - Item 5: `5. Apply surgical code fixes directly resolving the root cause without disabling any CI checks`
   - Item 6: `6. Verify fixes locally with targeted file linters (05-guideline-autofixer.py, check-prompts-loaded.py)`
   - Item 7: `7. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"`
   - Item 8: `8. Commit atomically via gitmap cpf, tag release, and loop until CI/CD is completely green`
5. Append the mandatory skill invocation suffix `[ci-cd-fix-gitmap-release](file;.agents/skills/ci-cd-fix-gitmap-release)`.
6. Output ONLY the resulting formatted markdown block.
7. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.
8. Additional Instructions Mandate: Always append the relative paths directive under ## Additional Instructions: '- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.'

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# High Priority Instruction

${Input Text Verbatim}

# Actionable Items Must Follow Non-Negotiable

1. Write spec and plan first
2. Inspect live CI/CD pipeline errors and execution timeline using gitmap pe -t (or gitmap pipeline fix)
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Perform grounded 4-part Root Cause Analysis (RCA) on exact failing step and log
5. Apply surgical code fixes directly resolving the root cause without disabling any CI checks
6. Verify fixes locally with targeted file linters (05-guideline-autofixer.py, check-prompts-loaded.py)
7. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"
8. Commit atomically via gitmap cpf, tag release, and loop until CI/CD is completely green

## Must follow and spawn agent using

[ci-cd-fix-gitmap-release](file;.agents/skills/ci-cd-fix-gitmap-release)

## Additional Instructions

- [/plan](slashCommand;plan) first before doing the work to reduce the credits.
- [/learn](slashCommand;learn) from [gitmap](file;.agents/skills/gitmap) skill to leverage GitMap live telemetry, pipeline self-healing, and toolchain caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
