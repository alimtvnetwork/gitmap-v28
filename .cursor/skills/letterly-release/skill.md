---
name: letterly-release
description: >-
  Formats raw voice dictation into minor version release ceremony instructions, changelog updates, and minor-bump skill invocation suffix.
---

# Release Mode — Letterly Prompt Formatter

Format whatever input text is provided according to the exact minor release template below, following the execute N-steps structure. Do NOT add conversational filler or commentary.

1. Clean the input text verbatim while strictly capturing version scope, changelog notes, and release constraints.
2. Structure the output starting immediately with `# High Priority Instruction`.
3. Put `${Input Text Verbatim}` directly beneath the high priority header.
4. Construct `# Actionable Items Must Follow Non-Negotiable`:
   - Item 1 is ALWAYS: `1. Write spec and plan first`
   - Item 2: `2. Enforce zero-storage GitHub Actions rules (zero routine artifact uploads)`
   - Item 3: `3. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"`
   - Item 4: `4. Consolidate and update release notes in root changelog.md and manifests`
   - Item 5: `5. Commit atomically via gitmap cpf and push release tag to remote tracking branch`
5. Append the mandatory release skill invocation suffix [minor-bump](file;.cursor/skills/minor-bump).
6. Output ONLY the resulting formatted markdown block.

${Input Text Verbatim} = The cleaned input text as it is, without conversational filler words.

Output Format:

# High Priority Instruction

${Input Text Verbatim}

# Actionable Items Must Follow Non-Negotiable

1. Write spec and plan first
2. Enforce zero-storage GitHub Actions rules (zero routine artifact uploads)
3. Execute minor version bump via python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"
4. Consolidate and update release notes in root changelog.md and manifests
5. Commit atomically via gitmap cpf "<module> - release minor version"
6. Tag release version and push to remote tracking branch

Must follow and spawn agent using

[minor-bump](file;.cursor/skills/minor-bump)

## Additional Instructions

learn [/learn](slashCommand;learn) if you have to learn something and [/plan](slashCommand;plan) stuff before working please.
