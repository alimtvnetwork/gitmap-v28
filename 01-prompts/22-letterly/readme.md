# Letterly Prompt Formatters (`22-letterly`)

This collection contains deterministic prompt transformation templates designed to ingest raw voice-dictated user transcripts from Letterly (mobile and desktop) and format them into rigorous, non-negotiable agent prompts.

## Directory Index

| # | Prompt File | Target Mode | Description | Companion Skill |
| :---: | :--- | :--- | :--- | :--- |
| **01** | [`01-mobile-letterly.md`](01-mobile-letterly.md) | Mobile Single-Line | Compact one-liner without line breaks starting with `# High Priority Instruction:` | `letterly-mobile` |
| **02** | [`02-desktop-letterly.md`](02-desktop-letterly.md) | Desktop Structured | Formats voice dictation into high-priority instructions, action items, and V6 suffix | `letterly-desktop` |
| **03** | [`03-execute-n-steps-letterly.md`](03-execute-n-steps-letterly.md) | Execute N-Steps | Verbatim input capture with discrete action items and V6 execution engine | `letterly-execute-n-steps` |
| **04** | [`04-plan-letterly.md`](04-plan-letterly.md) | Planning Spec | Formats raw input into architectural planning directives with no-build rules | `letterly-plan` |
| **05** | [`05-release-letterly.md`](05-release-letterly.md) | Minor Release (Letter U) | Formats input into minor release update ceremony directives | `letterly-release` |
| **06** | [`06-cicd-fix-release-letterly.md`](06-cicd-fix-release-letterly.md) | CI/CD Fix & Release | Formats input into `gitmap pe -t` telemetry diagnosis, 4-part RCA, and minor release | `letterly-cicd-fix-release` |
| **07** | [`07-mobile-cicd-fix-letterly.md`](07-mobile-cicd-fix-letterly.md) | Mobile CI/CD Fix | One-liner mobile prompt diagnosing via `gitmap pe -t` and triggering release | `letterly-mobile-cicd-fix` |
| **08** | [`08-run-letterly.md`](08-run-letterly.md) | Run Script | Formats input into run instructions via `./run.ps1` / `./run.sh` with auto dependency install | `letterly-run` |
| **09** | [`09-execute-with-verification-letterly.md`](09-execute-with-verification-letterly.md) | Execute + Retrospective Audit | Formats input into execute N-steps with mandatory retrospective AI verification at task end | `letterly-execute-with-verification` |
| **10** | [`10-execute-with-release-letterly.md`](10-execute-with-release-letterly.md) | Execute + Minor Release | Formats input into execute N-steps followed directly by green CI verification and minor release | `letterly-execute-with-release` |

### Cursor Prompt Suite (`23-cursor-prompts/`)

For Cursor environments requiring `.cursor/skills/` resolution, see [`../23-cursor-prompts/readme.md`](../23-cursor-prompts/readme.md) which houses the dedicated Cursor prompt suite (`01-mobile-letterly-cursor.md` through `10-execute-with-release-letterly-cursor.md`).

## Core Invariants

1. **Zero Conversational Filler:** Never output "Certainly! Here is your output:".
2. **Lossless Verbatim Capture:** Never drop specific technical flags, file paths, or commands from input.
3. **Mandatory Native Skill Triggers:** Every formatted prompt must conclude with its designated skill trigger: `[<skill-name>](file;.agents/skills/<skill-name>)`.
4. **Standard Letterly Suffix:** All prompt files in this directory strictly carry the `-letterly.md` file suffix.
5. **Strict Relative Git Paths Only:** Only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well.
