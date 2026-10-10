# Muse Prompts (`27-muse-prompts`) — Index & Catalog

> [!IMPORTANT]
> Category: `27-muse-prompts`
> Architecture: Self-Bootstrapping Autonomous Onboarding Directive
> Scope: Muse AI Session Bootstrapping, Operating Contract, and Multi-Agent Execution

---

## Overview

The `27-muse-prompts` prompt category houses the master onboarding prompt for
Muse AI sessions. Pasting the master prompt into a fresh Muse chat boots a
fully-onboarded autonomous agent: repository intake, terminal + GitHub auth
bootstrap, GitMap clone/build/learn, a memory operating-contract checklist
(turbo mode, strict no-delete guards), coding-guideline and design-system
intake — and then the standing **Task Confirmation & Multi-Agent Execution
Protocol**: every user message that is work gets a confirmed task breakdown
first (`Understood: [YES]` per task), then multiple concurrent agents complete it.

The canonical execution engine it delegates to is
[`02-execute-parent-task-with-n-steps.md`](../14-execute/02-execute-parent-task-with-n-steps.md)
(Canonical V6: `N = 300`, `A = 2`, `H = 2`, `C = 30`).

Its installable skill is `muse-master-prompt`
(`.agents/skills/muse-master-prompt/skill.md`, mirrored in `.cursor/skills/`).

---

## Directory Index & Catalog

| Prompt File | Version | Scope | Key Capabilities |
| :--- | :--- | :--- | :--- |
| [`01-muse-master-prompt.md`](01-muse-master-prompt.md) | 6.0.0 | Master Onboarding | Self-loop Phases 0–5, auth bootstrap, GitMap intake, memory contract checklist, task-confirm + multi-agent protocol (A=2/H=2), verification gates, Top-Instruction Priority Mandate. |
| [`02-muse-execute-in-a-step.md`](02-muse-execute-in-a-step.md) | 1.0.0 | Execute-in-a-Step | Breakdown-first listing, RUNNING + ETA declaration ("Are you running or not?"), 5-minute status pings, commit+push completion, both prompts as MD code blocks for Literally. |

---

## Directory Tree

```text
01-prompts/27-muse-prompts/
├── 01-muse-master-prompt.md   # the master prompt (paste into a fresh Muse session)
├── 02-muse-execute-in-a-step.md   # one-shot execution prompt (breakdown → RUNNING+ETA → 5-min pings → commit+push)
└── readme.md                   # this index
```
