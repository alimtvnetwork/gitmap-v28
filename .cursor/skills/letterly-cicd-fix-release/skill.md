---
name: letterly-cicd-fix-release
description: >-
  Formats raw voice dictation into CI/CD pipeline diagnosis using gitmap pe -t telemetry, 4-part RCA, atomic gitmap cpf commit, and ci-cd-fix-gitmap-release skill invocation suffix.
---

# CI/CD Fix & Release Mode — Letterly Prompt Formatter

Format whatever input text is provided according to the exact CI/CD fix and minor release template below, following the execute N-steps structure. Do NOT add conversational filler or commentary.

1. Clean the input text verbatim while capturing all failing workflow names, errors, and reproduction steps.
2. Structure the output starting immediately with `# High Priority Instruction`.
3. Put `${Input Text Verbatim}` directly beneath the high priority header.
4. Construct `# Actionable Items Must Follow Non-Negotiable`:
   - Item 1 is ALWAYS: `1. Write spec and plan first`
   - Item 2: `2. Inspect live CI/CD pipeline errors and execution timeline using gitmap pe -t (or gitmap pipeline fix)`
   - Item 3: `3. Perform grounded 4-part Root Cause Analysis (RCA) on exact failing step and log`
   - Item 4: `4. Apply surgical code fixes directly resolving the root cause without disabling any CI checks`
   - Item 5: `5. Verify fixes locally with targeted file linters (05-guideline-autofixer.py, check-prompts-loaded.py)`
   - Item 6: `6. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"`
   - Item 7: `7. Commit atomically via gitmap cpf, tag release, and loop until CI/CD is completely green`
5. Append the mandatory skill invocation suffix [ci-cd-fix-gitmap-release](file;.cursor/skills/ci-cd-fix-gitmap-release).
6. Output ONLY the resulting formatted markdown block.

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# High Priority Instruction

${Input Text Verbatim}

# Actionable Items Must Follow Non-Negotiable

1. Write spec and plan first
2. Inspect live CI/CD pipeline errors and execution timeline using gitmap pe -t (or gitmap pipeline fix)
3. Perform grounded 4-part Root Cause Analysis (RCA) on exact failing step and log
4. Apply surgical code fixes directly resolving the root cause without disabling any CI checks
5. Verify fixes locally with targeted file linters (05-guideline-autofixer.py, check-prompts-loaded.py)
6. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"
7. Commit atomically via gitmap cpf, tag release, and loop until CI/CD is completely green

Must follow and spawn agent using

[ci-cd-fix-gitmap-release](file;.cursor/skills/ci-cd-fix-gitmap-release)

## Additional Instructions

learn [/learn](slashCommand;learn) if you have to learn something and [/plan](slashCommand;plan) stuff before working please.
