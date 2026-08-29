---
name: Batched Loop Execution Without Release
description: Autonomous batched multi-agent execution engine for `.lovable/plans/pending/` tasks with strict 3-agent limit, chunked commits, locking matrix, and WOR workflow.
---

# Batched Loop Execution Without Release (WOR)

**Goal:** Autonomously execute pending tasks from `.lovable/plans/pending/` using a strictly batched multi-agent loop with maximum 3 concurrent sub-agents, file collision locking via `.lovable/temp/active-locks.json`, chunked commits, and without cutting releases (WOR).

## Non-Negotiable Core Rules

1. **Max 3 Sub-Agents:** Maximum 3 sub-agents running concurrently.
2. **Local Unit Tests Only:** No live API or network calls. Only run isolated local unit tests.
3. **Locking Matrix:** Register active target files in `.lovable/temp/active-locks.json` to prevent collisions.
4. **Temp State Tracking:** Log sub-agent tasks and progress in `.lovable/temp-agents/<task-name>.md`.
5. **3-Strike Rollback:** On 3 failures, revert dirty changes and log RCA to `.lovable/memory/last-failure.md`.
6. **Task Migration:** Move completed tasks from `.lovable/plans/pending/` to `.lovable/plans/completed/` and update `.lovable/plans/index.md`.
7. **Explicit Output Stats:** Always output tasks done, total completed, total pending, and remaining task list.
8. **Strict Relative Git Paths:** Total ban on absolute drive letters (`D:\...`, `C:\...`) and `file:///` URIs.
9. **Without Release (WOR):** Never bump versions, edit changelogs, or cut releases.
