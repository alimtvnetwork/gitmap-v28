package cmdagy

const defaultCicdFixWithReleasePromptFallback = `# Release-Triggered CI/CD Fix Loop — Workflow (must follow)

N = 200
`

const defaultRcaFixPromptFallback = `# Bug Fix with 4-Part RCA & Regression Verification — Workflow (must follow)

/goal Autonomously fix the failing CI/CD pipeline errors, strictly enforcing coding guidelines, and document the complete RCA before pushing.

## The 4-Part RCA Requirement (Mandatory Memory File)
Before modifying code, document the issue in .ai-memory/memory/issues/xx-<slug>.md:
1. Why it happened: High-level architectural breakdown of the failure.
2. How it happened: Technical execution flow that triggered the error.
3. Root Cause: Exact file, line, and dependency responsible.
4. Code Fix: Exact code snippets showing the remediation.

## Non-Negotiable Rules
- Functions <= 8-15 lines.
- Affirmative booleans only (is*, has*).
- Zero naked panics or unhandled errors; wrap with *apperror.AppError.
- Group all changes into a single atomic commit at the final step.
`
