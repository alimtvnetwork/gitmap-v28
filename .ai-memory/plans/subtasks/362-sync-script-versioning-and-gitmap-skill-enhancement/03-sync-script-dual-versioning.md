# Subtask 03: Sync Script Dual Versioning in `coding-guidelines`

- **Task Code:** `TASK-03`
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 03`
- **Status:** pending
- **Owner:** Automation Engineer

## Actions
1. Update `03-ai-scripts/38-sync-prompts-skills-scripts.py` in `coding-guidelines` with `--mode` (`v6` vs `v7`) and `--v7` flags.
2. Create dedicated executable `03-ai-scripts/38-sync-prompts-skills-scripts-v7.py` that defaults to V7 time-restriction and autobot checks.
3. Validate restrict conditions (weekday morning < 08:00 AM cutoff under target timezone or auto-window, ai-auto persona verification).
