# Execution Plan: Agent AI Synchronization and Task Management

- **Slug:** `agent-ai-synchronization-and-task-management`
- **Spec Directory:** `02-spec/21-app/agent-ai-synchronization-and-task-management/`
- **Subtasks Directory:** `.ai-memory/plans/subtasks/agent-ai-synchronization-and-task-management/`
- **Database Path:** `.ai-memory/temp-agents/01-agent-ai-synchronization-and-task-management/agent-task.db`
- **Status:** ACTIVE
- **Total Steps Budget:** 300
- **Operational Waves:** 4

---

## User Request (Verbatim)

```text
# Agent AI Synchronization and Task Management: high priority instruction, non-negotiable task

Can you please confirm a few things that if you have updated the skills based on the heatmap and agent AI commands so that the agents can synchronize, enter their task, plan their task, what they're analyzing, what they're editing, everything needs to be done through the agent AI so that the heatmaps can be discussed and it needs to be done very fast. Please confirm that these are implemented and the commands are put into these skills. That's the first thing. Second is that would also be referred from the executed steps in V6 prompt in the coding guideline. Please confirm that. And also mention and make sure that using the agent AI, when we create a task, it actually gives a JSON or the terminal based on how we are dealing with it. Usually, it should give a JSON response because it's for the AI nicely, like what is the task, what is the machine Id, and also so that the subtask can be created. Also, when the subtask is created, subtask will also tell us what is the main Task Id, machine, and so on. So that based on the subtask slug or Id, it can determine anything. And every time that the AI is trying to modify the file, it would also mention what it is trying to do, modify, analyze, or read or remove stuff, so a refactor. So all kinds of flag needs to be there nicely using enum. Make sure that these are there. And based on that, these skills were created in the coding guideline and also in the gitmap for the gitmap skill for the agent and the cursor. Please confirm that.

## slug: agent-ai-synchronization-and-task-management

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/agent-ai-synchronization-and-task-management/ and enqueue plan task in .ai-memory/plans/agent-ai-synchronization-and-task-management.md (subtasks in .ai-memory/plans/subtasks/agent-ai-synchronization-and-task-management/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Use `gitmap` AI agents to enter data
5. Task completion includes committing and pushing to Git
6. Confirm that skills are updated based on heatmap and agent AI commands for synchronization and task management.
7. Ensure tasks created by agent AI provide JSON responses detailing task and machine Id.
8. Verify that subtasks include main Task Id and machine information.
9. Implement flags using enum for AI file modifications, including modify, analyze, read, remove, and refactor.
10. Confirm skills are integrated into coding guidelines and GitMap for agent and cursor.
```

---

## Actionable Deliverables & Wave Schedule

### Wave 1: Specification Authoring & Subtask Decomposition (Phase 1)
- **Task-01:** Author `02-spec/21-app/agent-ai-synchronization-and-task-management/01-architecture-spec.md` and `02-component-and-cli-spec.md`.
- **Task-02:** Author subtasks in `.ai-memory/plans/subtasks/agent-ai-synchronization-and-task-management/` (`01-spec-and-schema.md`, `02-task-manager-engine.md`, `03-skills-and-prompts-update.md`, `04-verification-and-audit.md`).

### Wave 2: Engine Implementation — Machine ID & ActionType Enum (Phase 2)
- **Task-03:** Update `03-ai-scripts/46-agent-sqlite-task-manager.py` with `MachineId` column, deterministic machine identifier detection, and JSON output response containing `taskId`, `taskSlug`, `machineId`, `runDirectory`, `databasePath`.
- **Task-04:** Update `add-subtasks` / `subtask create` in `03-ai-scripts/46-agent-sqlite-task-manager.py` to record `ParentTaskId` and `MachineId`, returning JSON response containing parent task ID, subtask IDs, codes, slugs, and machine context.
- **Task-05:** Implement `ActionType` enum in `03-ai-scripts/46-agent-sqlite-task-manager.py` with valid values: `modify`, `analyze`, `read`, `remove`, `refactor`, `create`. Validate `--action` CLI flag and record in action log.
- **Task-06:** Implement `gitmap agent-ai heatmap discuss [<slug>]` command logic in task manager and ensure high-speed evaluation (<5ms).

### Wave 3: Skills & Prompts Synchronization (Phase 2)
- **Task-07:** Update V6 prompts and guideline skills:
  - `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`
  - `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`
  - `.cursor/skills/execute-parent-task-with-n-steps-v6/skill.md`
- **Task-08:** Update GitMap agent skills:
  - `.agents/skills/gitmap/SKILL.md`
  - `.cursor/skills/gitmap/skill.md`

### Wave 4: Targeted Verification, Consolidation & Atomic Push (Phase 3)
- **Task-09:** Run `python 03-ai-scripts/50-fastgate.py` and `python 03-ai-scripts/14-version-sync-checker.py` ensuring clean relative paths and zero violations.
- **Task-10:** Consolidate subtask plans into `.ai-memory/plans/completed/agent-ai-synchronization-and-task-management.md`, update registers, and commit/push atomically via GitMap (`gitmap cpf "agent-ai - synchronization task management and action enum telemetry"`).
