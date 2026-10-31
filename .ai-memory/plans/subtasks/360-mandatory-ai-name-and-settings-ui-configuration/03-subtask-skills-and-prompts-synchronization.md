# Subtask 03: Skills & Prompts Synchronization across Repositories

- **Parent Plan:** `360-mandatory-ai-name-and-settings-ui-configuration`
- **Slug:** `mandatory-ai-name-and-settings-ui-configuration - Task 03`
- **Status:** `completed`
- **Assigned Subagent:** Worker 03

---

## Deliverables
1. Update `.agents/skills/gitmap/skill.md` and `.cursor/skills/gitmap/skill.md`:
   - State that providing the AI agent name is **MANDATORY** for autonomous runs.
   - Include on-the-fly teaching examples (`muse`, `gemini`, `cursor-model`, `antigravity`, `claude`).
   - Document how defaults are changed via CLI settings and the Web Settings UI.
2. Update `.agents/skills/gitmap-automation-bot/skill.md` and `.cursor/skills/gitmap-automation-bot/skill.md`:
   - Emphasize mandatory agent names and on-the-fly teaching behavior.
3. Update `execute-parent-task-with-n-steps-v7` and `v6` skills and prompts:
   - Make agent name specification mandatory in Phase 1 identity configuration.
   - Clarify the on-the-fly error message and how autonomous agents must self-correct.
4. Synchronize all changes into `d:/work/coding-guidelines/`.
