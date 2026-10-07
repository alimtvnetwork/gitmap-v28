# Execution Prompts (`14-execute`) — Index & Catalog

> [!IMPORTANT]
> Category: `14-execute`
> Architecture: Canonical V6 Parameter-Driven Multi-Agent Protocol
> Scope: Autonomous Task Execution, Subagent Orchestration, and Continuous Self-Looping

---

## Overview

The `14-execute` prompt category houses the autonomous execution prompts for decomposing and executing complex software engineering tasks.

The canonical orchestrator is [`02-execute-parent-task-with-n-steps-v6.md`](02-execute-parent-task-with-n-steps-v6.md), running the **Canonical V6 Workflow** (`N = 300`, `A = 2`, `H = 2`, `C = 30`). It enforces mandatory subagent spawning via `invoke_subagent`, GitMap command primacy, SQLite task tracking, and atomic final commits without raw git usage.

> [!NOTE]
> **Archived Versions Notice:**
> Legacy execution prompt iterations have been archived to [`01-prompts/19-old-execute-prompts/`](../19-old-execute-prompts/readme.md) and [`06-archive/execute/`](../../06-archive/execute/) and are not synchronized downstream. All previous execution skills have been purged from `.agents/skills/` and `.cursor/skills/`, leaving `execute-parent-task-with-n-steps-v6` as the sole canonical execution engine.
>
> **Strict Relative Git Paths Only:** Only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well (TOTAL BAN on `file:///` URIs and absolute filesystem paths).

---

## Directory Index & Catalog

| Prompt File | Version | Scope | Key Capabilities |
| :--- | :--- | :--- | :--- |
| [`01-execute-pending-tasks.md`](01-execute-pending-tasks.md) | Canonical V6 | Backlog Execution | Parameter-driven execution loop (N=300, A=2, H=2, C=30) executing pending tasks in `.ai-memory/plans/pending/` with SQLite tracking and GitMap atomic commits. |
| [`02-execute-parent-task-with-n-steps-v6.md`](02-execute-parent-task-with-n-steps-v6.md) | Canonical V6 | Master Orchestration | Parameter-driven execution loop (N=300, A=2, H=2, C=30) with mandatory subagents, SQLite task logging, and GitMap atomic commits. |
| [`03-execute-batched-loop.md`](03-execute-batched-loop.md) | V4 | Batched Loop | Micro-task batched multi-agent loop with file collision matrix. |
| [`04-execute-ai-instruction-writer.md`](04-execute-ai-instruction-writer.md) | V4 | Instruction Writer | Decomposes complex prompt requirements into subagents and modular specs. |
| [`05-execute-batched-loop-wor.md`](05-execute-batched-loop-wor.md) | V4 | Batched Loop (WOR) | Batched multi-agent loop without automatic release triggers. |
| [`06-execute-batched-loop-v2.md`](06-execute-batched-loop-v2.md) | V4 | Batched Loop V2 | Streamlined batched loop execution engine. |
| [`07-run.md`](07-run.md) | V6 | Run Script Orchestration | Autonomous execution of run.ps1 / run.sh with GitMap or local-install fallback. |

---

## Directory Tree

```text
01-prompts/14-execute/
├── 01-execute-pending-tasks.md
├── 02-execute-parent-task-with-n-steps-v6.md    # Canonical V6 Orchestrator
├── 03-execute-batched-loop.md
├── 04-execute-ai-instruction-writer.md
├── 05-execute-batched-loop-wor.md
├── 06-execute-batched-loop-v2.md
├── 07-run.md
└── readme.md

Archived:
01-prompts/19-old-execute-prompts/
├── 01-execute-robust-loop.md
├── 02-fix-subtask-naming-convention.md
├── 03-execute-parent-task-with-n-steps.md
└── 04-parent-task-in-below-steps.md

06-archive/execute/
├── 06-execute-parent-task-with-n-steps-v2.md
├── 08-excute-parent-old.md
├── 10-execute-parent-task-with-n-steps-v3.md
├── 11-execute-parent-task-with-n-steps-v4.md
└── 12-execute-parent-task-with-n-steps-v5.md
```
