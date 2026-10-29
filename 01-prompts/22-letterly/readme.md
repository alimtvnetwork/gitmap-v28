# Letterly Prompt Formatters (`22-letterly`)

This collection contains deterministic prompt transformation templates designed to ingest raw voice-dictated user transcripts from Letterly (mobile and desktop) and format them into rigorous, non-negotiable agent prompts.

## Directory Index

| # | Prompt File | Target Mode | Description | Canonical Successor Skill |
| :---: | :--- | :--- | :--- | :--- |
| **01** | [`01-mobile-letterly.md`](01-mobile-letterly.md) | Mobile Single-Line | Compact one-liner without line breaks starting with `# High Priority Instruction:` | `execute-parent-task-with-n-steps-v6` |
| **02** | [`02-desktop-letterly.md`](02-desktop-letterly.md) | Desktop Structured | Formats voice dictation into high-priority instructions, action items, and V6 suffix | `execute-parent-task-with-n-steps-v6` |
| **03** | [`03-execute-n-steps-letterly.md`](03-execute-n-steps-letterly.md) | Execute N-Steps | Verbatim input capture with discrete action items and V6 execution engine | `execute-parent-task-with-n-steps-v6` |
| **04** | [`04-plan-letterly.md`](04-plan-letterly.md) | Planning Spec | Formats raw input into architectural planning directives with no-build rules | `plan-spec-steps-v2` |
| **05** | [`05-release-letterly.md`](05-release-letterly.md) | Minor Release (Letter U) | Formats input into minor release update ceremony directives | `minor-bump` |
| **06** | [`06-cicd-fix-release-letterly.md`](06-cicd-fix-release-letterly.md) | CI/CD Fix & Release | Formats input into `gitmap pe -t` telemetry diagnosis, 4-part RCA, and minor release | `ci-cd-fix-gitmap-release` |
| **07** | [`07-mobile-cicd-fix-letterly.md`](07-mobile-cicd-fix-letterly.md) | Mobile CI/CD Fix | One-liner mobile prompt diagnosing via `gitmap pe -t` and triggering release | `ci-cd-fix-gitmap-release` |
| **08** | [`08-run-letterly.md`](08-run-letterly.md) | Run Script | Formats input into run instructions via `./run.ps1` / `./run.sh` with auto dependency install | `run` |
| **09** | [`09-execute-with-verification-letterly.md`](09-execute-with-verification-letterly.md) | Execute + Retrospective Audit | Formats input into execute N-steps with mandatory retrospective AI verification at task end | `execute-parent-task-with-n-steps-v6` / `ai-verification` |
| **10** | [`10-execute-with-release-letterly.md`](10-execute-with-release-letterly.md) | Execute + Minor Release | Formats input into execute N-steps followed directly by green CI verification and minor release | `execute-parent-task-with-n-steps-v6` / `minor-bump` |
| **11** | [`11-git-pull-first-letterly.md`](11-git-pull-first-letterly.md) | Git Pull First | Formats input into pull-first execution directives ensuring gitmap pull runs before any work | `execute-parent-task-with-n-steps-v6` |
| **12** | [`12-execute-recommit-time-restricted-letterly.md`](12-execute-recommit-time-restricted-letterly.md) | Recommit & Time-Restricted | Formats voice dictation with on-call time restrictions, recommit instructions, autobot identity, and V7 execution engine | `execute-parent-task-with-n-steps-v7` |

### Cursor Prompt Suite (`23-cursor-prompts/`)

For Cursor environments requiring `.cursor/skills/` resolution, see [`../23-cursor-prompts/readme.md`](../23-cursor-prompts/readme.md) which houses the dedicated Cursor prompt suite (`01-mobile-letterly-cursor.md` through `11-git-pull-first-letterly-cursor.md`).

## Core Invariants

1. **Zero Conversational Filler:** Never output "Certainly! Here is your output:".
2. **Lossless Verbatim Capture:** Never drop specific technical flags, file paths, or commands from input.
3. **Mandatory Canonical Skill Triggers:** Every formatted prompt must conclude with its designated canonical skill trigger: `[<skill-name>](file;.agents/skills/<skill-name>)`.
4. **Standard Letterly Suffix:** All prompt files in this directory strictly carry the `-letterly.md` file suffix.
5. **Strict Relative Git Paths Only:** Only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well.
6. **Decoupled Execution Engine:** Letterly prompts remain generic and decoupled from internal engine step budgeting; lifecycle execution parameters and concurrency are encapsulated directly within the invoked skills (e.g. `execute-parent-task-with-n-steps-v6`).
7. **Heatmap Telemetry & Write Collision Mandate:** All Letterly formatters embed the GitMap Heatmap one-liner (`gitmap heatmap discuss` / `gitmap heatmap log` / `gitmap heatmap live`) to track in-flight file touches, prevent write collisions across parallel agent swarms ($A = 3$, Capacity = 6), and transparently report 3-tier database locations.

## GitMap Heatmap One-Liner Integration

All 11 Letterly formatters enforce GitMap Heatmap integration:
- **Desktop & Multi-Line Formatters (`02`–`06`, `08`–`11`):** Include the `Heatmap Mandate` in the numbered instructions and embed `- Track in-flight file touches and prevent write collisions using GitMap Heatmap (\`gitmap heatmap discuss\` / \`gitmap heatmap log\` / \`gitmap heatmap live\`) with transparent 3-tier database location reporting.` under `## Additional Instructions`.
- **Mobile Single-Line Formatters (`01`, `07`):** Embed `track file touches and prevent write collisions using GitMap Heatmap (gitmap heatmap discuss / log / live); ` within the continuous single-line directives.
- **3-Tier Database Reporting:** Enforces transparency across Host Root (`~/.gitmap/agent-root.db`), Repo-Cache (`repo-cache/.heatmap/<ai-slug>/<repo-slug>/heatmap.db`), and Local Task DB (`.ai-memory/temp-agents/<slug>/agent-task.db`).

