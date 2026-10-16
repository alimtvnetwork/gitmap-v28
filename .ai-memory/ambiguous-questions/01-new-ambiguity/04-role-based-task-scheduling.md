# Ambiguity 04: Role-Based Task Scheduling and Distributed Fleet Delegation

- **Slug:** `role-based-task-scheduling`
- **Status:** `open`
- **Raised At:** 2026-10-16
- **Blocking:** No (Deferred to future distributed scheduling phase per user request)
- **Spec Reference:** [02-component-spec.md](../../../02-spec/21-app/gitmap-command-enhancements/02-component-spec.md)
- **Plan Reference:** [.ai-memory/plans/gitmap-command-enhancements.md](../../plans/gitmap-command-enhancements.md)

---

## 1. Question / Unresolved Decision

The user mandated creating multi-node task visibility and delegation commands for `gitmap` while explicitly deferring role-based task scheduling to a future phase:

> *"In the future, we will add role-based system. So in future, you can put a question mark there, which we will discuss later on. But I think you should have all kinds of things... if the `gitmap` is already delegated to a list of tasks, because every task needs to be enqueued in the task list. This way, `gitmap` knows how many tasks are actually requested. So based on that, `gitmap` can give this information like how many tasks, and tasks can be also identified as heavy duty tasks, and that can be revealed from one node to another for now."*

Specific open architectural questions:
1. **Node Role Taxonomy (?):** What formal roles should be assigned to fleet nodes?
   - Candidate roles: `Primary Orchestrator`, `Heavy Compute Worker`, `Lightweight Scraper`, `CI/CD Build Runner`, `Database Maintainer`, `Storage Mirror`.
2. **Work-Stealing & Task Migration Protocol (?):** When a node is identified as overloaded with heavy-duty tasks (e.g. CPU > 85%, RAM > 90%), how should pending tasks in its `agent-task.db` queue be reassigned or stolen by idle nodes?
   - Should migration occur via direct SSH RPC command (`gitmap nodes delegate <task-id> <target-node>`)?
   - Or should nodes poll a centralized task registry or broadcast via cluster gossip?
3. **Role Affinity & Task Matching Rules (?):** How should task tags (e.g. `heavy-duty`, `os-clean`, `compile-polyglot`, `ai-eval`) map to node capabilities (e.g. Windows Server vs Ubuntu GPU node vs macOS Darwin)?

---

## 2. Options Considered

- **Option A (Central Master Registry via `ai_agents.db`):** The primary master node schedules tasks and dispatches to worker nodes based on their declared roles (`role=build`, `role=eval`).
- **Option B (Decentralized Peer-to-Peer Lease-Claiming):** Nodes register their capabilities and claim unassigned tasks from shared Split-DB queues using optimistic lease locks.
- **Option C (Ad-hoc CLI Delegation for Phase 1):** Phase 1 implements transparent queue visibility via `gitmap nodes tasks`, tagging tasks with `isHeavyDuty` boolean flags. Explicit automated migration and role scheduling remain marked with `?` for interactive user review.

---

## 3. Impact if Guessed Wrong

Prematurely designing a rigid role-scheduling matrix before the user defines cluster topologies across Windows, Linux, and Windows Server risks overcomplicating task claiming and causing lock contention on SQLite databases. Storing this ambiguity ledger entry ensures lossless capture of the user's requirement while keeping Phase 1 delivery robust and unblocked.
