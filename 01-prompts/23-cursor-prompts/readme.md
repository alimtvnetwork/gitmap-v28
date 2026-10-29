# Cursor Letterly Prompt Formatters (`23-cursor-prompts`)

This directory contains deterministic prompt transformation templates specifically structured for Cursor environments. These templates ingest raw voice-dictated user transcripts from Letterly (mobile and desktop) and format them into rigorous, non-negotiable agent prompts referencing the `.cursor/skills/` skill tree.

## Directory Index

| # | Prompt File | Target Mode | Description | Canonical Successor Cursor Skill |
| :---: | :--- | :--- | :--- | :--- |
| **01** | [`01-mobile-letterly-cursor.md`](01-mobile-letterly-cursor.md) | Mobile Single-Line | Compact one-liner without line breaks starting with `# High Priority Instruction:` | `execute-parent-task-with-n-steps-v6` |
| **02** | [`02-desktop-letterly-cursor.md`](02-desktop-letterly-cursor.md) | Desktop Structured | Formats voice dictation into high-priority instructions, action items, and V6 Cursor suffix | `execute-parent-task-with-n-steps-v6` |
| **03** | [`03-execute-n-steps-letterly-cursor.md`](03-execute-n-steps-letterly-cursor.md) | Execute N-Steps | Verbatim input capture with discrete action items and V6 Cursor execution engine | `execute-parent-task-with-n-steps-v6` |
| **04** | [`04-plan-letterly-cursor.md`](04-plan-letterly-cursor.md) | Planning Spec | Formats raw input into architectural planning directives with Cursor plan skill | `plan-spec-steps-v2` |
| **05** | [`05-release-letterly-cursor.md`](05-release-letterly-cursor.md) | Minor Release (Letter U) | Formats input into minor release update ceremony directives with Cursor release skill | `minor-bump` |
| **06** | [`06-cicd-fix-release-letterly-cursor.md`](06-cicd-fix-release-letterly-cursor.md) | CI/CD Fix & Release | Formats input into `gitmap pe -t` telemetry diagnosis, 4-part RCA, and Cursor CI/CD release | `ci-cd-fix-gitmap-release` |
| **07** | [`07-mobile-cicd-fix-letterly-cursor.md`](07-mobile-cicd-fix-letterly-cursor.md) | Mobile CI/CD Fix | One-liner mobile prompt diagnosing via `gitmap pe -t` and triggering Cursor release | `ci-cd-fix-gitmap-release` |
| **08** | [`08-run-letterly-cursor.md`](08-run-letterly-cursor.md) | Run Script | Formats input into run instructions via `./run.ps1` / `./run.sh` with auto dependency install | `run` |
| **09** | [`09-execute-with-verification-letterly-cursor.md`](09-execute-with-verification-letterly-cursor.md) | Execute + Retrospective Audit | Formats input into execute N-steps with mandatory retrospective AI verification at task end | `execute-parent-task-with-n-steps-v6` / `ai-verification` |
| **10** | [`10-execute-with-release-letterly-cursor.md`](10-execute-with-release-letterly-cursor.md) | Execute + Minor Release | Formats input into execute N-steps followed directly by green CI verification and minor release | `execute-parent-task-with-n-steps-v6` / `minor-bump` |
| **11** | [`11-git-pull-first-letterly-cursor.md`](11-git-pull-first-letterly-cursor.md) | Git Pull First | Formats input into pull-first execution directives ensuring gitmap pull runs before any work | `execute-parent-task-with-n-steps-v6` |
| **12** | [`12-muse-long-update-letterly-cursor.md`](12-muse-long-update-letterly-cursor.md) | Muse Long Update | Formats input into Muse long-running overnight execution with gitmap update first, rate-limited gitmap pa / repo pull, autopilot ambiguity logging, pragmatic commits, and conditional CI release | `execute-parent-task-with-n-steps-v6` |
| **13** | [`13-muse-sleepy-letterly-cursor.md`](13-muse-sleepy-letterly-cursor.md) | Muse Sleepy Mode | Formats input into "I am sleepy" autopilot execution with repo pull first, closest-choice ambiguity resolution, morning review log, no-build/no-test default, and pragmatic commits | `execute-parent-task-with-n-steps-v6` |
| **14** | [`14-execute-recommit-time-restricted-letterly-cursor.md`](14-execute-recommit-time-restricted-letterly-cursor.md) | Recommit & Time-Restricted | Formats voice dictation with on-call time restrictions (< 08:00 AM weekdays, anytime weekends in target timezone), recommit instructions, autobot identity, and V7 Cursor execution engine | `execute-parent-task-with-n-steps-v7` |

## Core Invariants

1. **Zero Conversational Filler:** Never output "Certainly! Here is your output:".
2. **Lossless Verbatim Capture:** Never drop specific technical flags, file paths, or commands from input.
3. **Mandatory Canonical Skill Triggers:** Every formatted prompt must conclude with its designated canonical skill trigger: `[<skill-name>](file;.cursor/skills/<skill-name>)`.
4. **Standard Cursor Suffix:** All prompt files in this directory strictly carry the `-letterly-cursor.md` file suffix.
5. **Strict Relative Git Paths Only:** Only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well.
6. **Decoupled Execution Engine:** Cursor prompts remain generic and decoupled from internal engine step budgeting; lifecycle execution parameters and concurrency are encapsulated directly within the invoked skills (e.g. `execute-parent-task-with-n-steps-v6`).
7. **Heatmap Telemetry & Collision Prevention:** All formatters enforce tracking in-flight file touches, avoiding write collisions across concurrent subagents, and querying cognitive graphs via GitMap Heatmap (`gitmap heatmap discuss` / `gitmap heatmap log` / `gitmap heatmap live`) with transparent 3-tier database location reporting.

## GitMap Heatmap One-Liner Integration

All 11 prompts integrate the standardized GitMap Heatmap one-liner to track in-flight file touches, avoid write collisions across concurrent subagents ($A = 3$, Capacity = 6), and query cognitive activity graphs with transparent 3-tier database location reporting:

- **Desktop Formatters (`02`, `03`, `04`, `05`, `06`, `08`, `09`, `10`, `11`):** Numbered Heatmap Mandate rule in the prompt specification plus an actionable directive bullet under `## Additional Instructions`.
- **Mobile Formatters (`01`, `07`):** Embedded single-line directive within Step 7 and continuous single-line output template.
- **Commands:** `gitmap heatmap discuss` / `gitmap heatmap log` / `gitmap heatmap live`.

