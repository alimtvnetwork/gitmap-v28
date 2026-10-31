# Subtask 03: Coding-Guidelines Sync Script Dual Versioning (V6 vs V7)

- **Task Code:** `TASK-03`
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 03`
- **Status:** completed
- **Owner:** Fleet Sync Engineer
- **Specification:** [03-sync-script-dual-versioning-spec.md](../../../02-spec/21-app/367-sync-script-versioning-and-gitmap-skill-enhancement/03-sync-script-dual-versioning-spec.md)

## Actions
1. Verify `03-ai-scripts/38-sync-prompts-skills-scripts.py` and `38-sync-prompts-skills-scripts-v7.py` in `gitmap` and `coding-guidelines`.
2. Ensure V7 execution path strictly validates calendar morning window (< 08:00 AM on weekdays) and AI bot identity.
3. Ensure V6 execution path runs as is (unrestricted).
4. Verify parity between `.agents/scripts/` and `03-ai-scripts/`.
