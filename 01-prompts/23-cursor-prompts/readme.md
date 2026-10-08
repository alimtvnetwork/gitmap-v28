# Cursor Letterly Prompt Formatters (`23-cursor-prompts`)

This directory contains deterministic prompt transformation templates specifically structured for Cursor environments. These templates ingest raw voice-dictated user transcripts from Letterly (mobile and desktop) and format them into rigorous, non-negotiable agent prompts referencing the `.cursor/skills/` skill tree.

## Directory Index

| # | Prompt File | Target Mode | Description | Companion Cursor Skill |
| :---: | :--- | :--- | :--- | :--- |
| **01** | [`01-mobile-letterly-cursor.md`](01-mobile-letterly-cursor.md) | Mobile Single-Line | Compact one-liner without line breaks starting with `# High Priority Instruction:` | `letterly-mobile` |
| **02** | [`02-desktop-letterly-cursor.md`](02-desktop-letterly-cursor.md) | Desktop Structured | Formats voice dictation into high-priority instructions, action items, and V6 Cursor suffix | `letterly-desktop` |
| **03** | [`03-execute-n-steps-letterly-cursor.md`](03-execute-n-steps-letterly-cursor.md) | Execute N-Steps | Verbatim input capture with discrete action items and V6 Cursor execution engine | `letterly-execute-n-steps` |
| **04** | [`04-plan-letterly-cursor.md`](04-plan-letterly-cursor.md) | Planning Spec | Formats raw input into architectural planning directives with Cursor plan skill | `letterly-plan` |
| **05** | [`05-release-letterly-cursor.md`](05-release-letterly-cursor.md) | Minor Release (Letter U) | Formats input into minor release update ceremony directives with Cursor release skill | `letterly-release` |
| **06** | [`06-cicd-fix-release-letterly-cursor.md`](06-cicd-fix-release-letterly-cursor.md) | CI/CD Fix & Release | Formats input into `gitmap pe -t` telemetry diagnosis, 4-part RCA, and Cursor CI/CD release | `letterly-cicd-fix-release` |
| **07** | [`07-mobile-cicd-fix-letterly-cursor.md`](07-mobile-cicd-fix-letterly-cursor.md) | Mobile CI/CD Fix | One-liner mobile prompt diagnosing via `gitmap pe -t` and triggering Cursor release | `letterly-mobile-cicd-fix` |
| **08** | [`08-run-letterly-cursor.md`](08-run-letterly-cursor.md) | Run Script | Formats input into run instructions via `./run.ps1` / `./run.sh` with auto dependency install | `letterly-run` |
| **09** | [`09-execute-with-verification-letterly-cursor.md`](09-execute-with-verification-letterly-cursor.md) | Execute + Retrospective Audit | Formats input into execute N-steps with mandatory retrospective AI verification at task end | `letterly-execute-with-verification` |
| **10** | [`10-execute-with-release-letterly-cursor.md`](10-execute-with-release-letterly-cursor.md) | Execute + Minor Release | Formats input into execute N-steps followed directly by green CI verification and minor release | `letterly-execute-with-release` |

## Core Invariants

1. **Zero Conversational Filler:** Never output "Certainly! Here is your output:".
2. **Lossless Verbatim Capture:** Never drop specific technical flags, file paths, or commands from input.
3. **Mandatory Native Skill Triggers:** Every formatted prompt must conclude with its designated skill trigger: `[<skill-name>](file;.cursor/skills/<skill-name>)`.
4. **Standard Cursor Suffix:** All prompt files in this directory strictly carry the `-letterly-cursor.md` file suffix.
5. **Strict Relative Git Paths Only:** Only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well.
