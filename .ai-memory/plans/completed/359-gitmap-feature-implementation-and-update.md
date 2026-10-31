# Plan 359: GitMap Feature Implementation and Update

- **Plan Slug:** `gitmap-feature-implementation-and-update`
- **Spec Reference:** `02-spec/21-app/359-gitmap-feature-implementation-and-update/01-architecture-spec.md`
- **Status:** `completed`
- **Created At:** 2026-10-31
- **Completed At:** 2026-10-31

---

## 1. Goal

Implement autonomous GitMap information configuration (`gitmap login set ai-auto [agent-name]`, `gitmap login set global ai-auto`), commit `--ai-agent` flag with Split-DB heatmap telemetry attribution, and cross-repository synchronization of the canonical `gitmap` skill, `gitmap-automation-bot` skill, and V6/V7 execution prompts.

---

## 2. Subtasks Breakdown

1. `01-subtask-gitmap-login-set-ai-auto.md`: Implement `gitmap login set ai-auto [agent-name]` and `gitmap login set global ai-auto` in `cli/cmdlogin/` with unit tests and typo corrections. [completed]
2. `02-subtask-gitmap-commit-ai-agent-flag-and-heatmap.md`: Add `--ai-agent` flag to `cli/cmdcommit/`, auto-format bot author identity, and update heatmap logging with agent attribution. [completed]
3. `03-subtask-skills-prompts-and-v6-v7-checklist-sync.md`: Update `gitmap`, `gitmap-automation-bot`, `execute-parent-task-with-n-steps-v6`, and `execute-parent-task-with-n-steps-v7` across the `gitmap-v28` and `coding-guidelines-v24` repos. [completed]
4. `04-subtask-minor-bump-and-release.md`: Bump minor version, update manifests (`version.json`, `package.json`, `changelog.md`), and execute release ceremony. [completed]
