# Subtask 05: GitMap Skill Updates & Cross-Repository Fleet Synchronization

- **Parent Plan:** [368-sql-query-execution-and-gitmap-skill-enhancement.md](../../pending/368-sql-query-execution-and-gitmap-skill-enhancement.md)
- **Slug:** `sql-query-execution-and-gitmap-skill-enhancement - Task 05`
- **Status:** `pending`
- **Assignee:** V7 Master Autonomous Execution Engine

## Action Items
1. Update `.agents/skills/gitmap/skill.md` and `.cursor/skills/gitmap/skill.md` in both `gitmap` and `coding-guidelines`:
   - Enshrine `gitmap agy summaries` and `gitmap db-read agy "<sql>"` as the mandatory replacement for Antigravity database queries.
   - Enshrine `gitmap heatmap touch <path> --reason "<why>"` before analyzing files.
   - Add reference to `.agents/rules/ai_code_execution_restrictions.md`.
2. Run `python 03-ai-scripts/38-sync-prompts-skills-scripts-v7.py --force --auto-window --schedule` to synchronize all 43 repositories.
