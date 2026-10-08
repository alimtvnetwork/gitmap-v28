# Mobile CI/CD Fix Mode — Letterly Prompt Formatter

Combine the entire output into exactly ONE continuous single-line paragraph with zero newlines, zero line gaps, and zero conversational filler.

1. Clean the input text verbatim without conversational filler words (`um`, `ah`, `uh`).
2. Prefix with `[/goal](slashCommand;goal) [/learn](slashCommand;learn) Run gitmap pe -t to diagnose live CI/CD errors, fix all pipeline failures via 4-part RCA, verify locally, commit atomically via gitmap cpf, and release minor update. `.
3. Append `${Input Text Verbatim}`.
4. Suffix with ` - strictly use relative git paths only; only add the relative paths, never add the absolute path during your work (this must be respected on the release page and in release notes as well), must follow the skill [ci-cd-fix-gitmap-release](file;.agents/skills/ci-cd-fix-gitmap-release)`.
5. Output exactly that single line with no leading "Output" or markdown code fences.
6. Relative Paths Mandate: Only add the relative paths, never add the absolute path during your work; ensure this is respected on the release page and in release notes as well.

${Input Text Verbatim} = The cleaned input text as it is, without filler words.

Output Format:
[/goal](slashCommand;goal) [/learn](slashCommand;learn) Run gitmap pe -t to diagnose live CI/CD errors, fix all pipeline failures via 4-part RCA, verify locally, commit atomically via gitmap cpf, and release minor update. ${Input Text Verbatim} - strictly use relative git paths only; only add the relative paths, never add the absolute path during your work (this must be respected on the release page and in release notes as well), must follow the skill [ci-cd-fix-gitmap-release](file;.agents/skills/ci-cd-fix-gitmap-release)
