# Plan 360: Mandatory AI Name Enforcement, On-The-Fly Teaching & Settings UI Configuration

- **Plan Slug:** `mandatory-ai-name-and-settings-ui-configuration`
- **Spec Reference:** `02-spec/21-app/360-mandatory-ai-name-and-settings-ui-configuration/01-architecture-spec.md`
- **Status:** `completed`
- **Created At:** 2026-10-31
- **Completed At:** 2026-10-31

---

## 1. Goal

Enforce mandatory AI agent names for autonomous operations, provide on-the-fly teaching examples when names are missing, integrate `ai-auto` defaults into GitMap settings and the Web Settings UI, and synchronize all skills and execution prompts across repositories.

---

## 2. Subtasks Breakdown

1. `01-subtask-mandatory-ai-name-validation-and-teaching.md`: Update `cli/cmdlogin/` and `cli/cmdcommit/` to enforce mandatory AI agent name with on-the-fly teaching examples and settings fallback. [completed]
2. `02-subtask-settings-cli-and-ui-configuration.md`: Add `ai_auto.default_*` settings to `cli/cmdui/` (Web UI markup and server handlers) and `cli/cmdagy/` (`gitmap settings set`). [completed]
3. `03-subtask-skills-and-prompts-synchronization.md`: Update skills and prompts in `coding-guidelines` and `gitmap` making AI agent names strictly mandatory for autonomous runs and detailing on-the-fly error recovery. [completed]
4. `04-subtask-minor-bump-and-release.md`: Bump minor version, update changelog/release notes, and perform release ceremony. [completed]
