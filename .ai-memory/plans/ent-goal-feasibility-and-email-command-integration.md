# Master Plan: ENT Goal Feasibility and Email Command Integration

**Slug:** `ent-goal-feasibility-and-email-command-integration`  
**Parent Task ID:** `task-20261019-ent-goal-feasibility-and-email-command-integration`  
**Budget:** 400 steps  
**Specification:** `02-spec/21-app/ent-goal-feasibility-and-email-command-integration/spec.md`  

---

## User Request (Verbatim)

```text
# ENT Goal Feasibility and Email Command Integration: high priority instruction, non-negotiable task

In the `gitmap`, you already did some planning on the ENT goal, right? You did some feasibility analysis, things like that. My question is, if we use the ENT goal, is our system going to be slower or not? That's one question. Second is, can you do it in one shot without any hassle? Also, for the email section, we should be able to send commands using email, just like `Antigravity Manager or Antigravity Tools` package. You can check the `Antigravity Manager or Antigravity Tools`, Antigravity Tool, how we send and receive email with several formats. We should have something similar to the `Antigravity Manager or Antigravity Tools`, so plan and craft that and have similar features. Usually, the email should be checked. There should be a lemon CLI or service running on the OS, Windows, Ubuntu, macOS, that checks every five seconds, not less than that, but it can be configurable. I think eight seconds is absolutely fine. If we send any command, when it receives, it will reply back again that it received that command. Based on that, it can run commands, it could work with `Antigravity Manager or Antigravity Tools`, it could work with many different things. Work on that, and if the ENT goal is feasible, viable, I want you to migrate the whole code base to ENT goal. Make sure we have type-safe code. We shouldn't have any dynamic SQL code. Remember that. Make sure the commands that you have recently created, you write the end-to-end test cases, as many as possible, so that the edge cases are found. Take your time and do that.

## task: ENT Goal Feasibility and Email Command Integration
## slug: ent-goal-feasibility-and-email-command-integration

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/ent-goal-feasibility-and-email-command-integration/ and enqueue plan task in .ai-memory/plans/ent-goal-feasibility-and-email-command-integration.md (subtasks in .ai-memory/plans/subtasks/ent-goal-feasibility-and-email-command-integration/) first
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Conduct feasibility analysis on the ENT goal to determine if it will slow down the system.
6. Plan and implement email command integration similar to `Antigravity Manager or Antigravity Tools`.
7. Develop a lemon CLI or service for checking emails every five to eight seconds, configurable, on Windows, Ubuntu, macOS.
8. Ensure type-safe code and eliminate any dynamic SQL code.
9. Write comprehensive end-to-end test cases for newly created commands to cover edge cases.
10. Migrate the entire codebase to the ENT goal if feasible and viable.

## Must follow and spawn agent using

@[.agents/skills/execute-parent-task-with-n-steps-v6]

## Additional Instructions

- /plan first before doing the work to reduce the credits.
- /learn from @[.agents/skills/gitmap] skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
```

---

## Subtask Breakdown & Execution Roadmap

| Subtask ID | Title | Disjoint File Ownership | State |
|---|---|---|---|
| **Subtask-01** | Feasibility Analysis & Type-Safe SQL Audit | `02-spec/21-app/ent-goal-feasibility-and-email-command-integration/ent-feasibility.md`, `cli/dbengine/` | `COMPLETED` |
| **Subtask-02** | Email Command Parser & Multi-Format Envelope | `cli/cmdlemon/lemon_types.go`, `cli/cmdlemon/email_parser.go`, `cli/cmdlemon/email_parser_test.go` | `COMPLETED` |
| **Subtask-03** | Immediate Receipt Acknowledgment (ACK) & Command Replier | `cli/cmdlemon/email_ack.go`, `cli/cmdlemon/email_ack_test.go`, `cli/cmdlemon/email_reporter.go` | `COMPLETED` |
| **Subtask-04** | Command Execution Pipeline & AGM Alignment | `cli/cmdlemon/command_executor.go`, `cli/cmdlemon/command_executor_test.go` | `COMPLETED` |
| **Subtask-05** | GitMap Daemon (`gitmap-daemon`) Polling Loop, Settings Toggle & Cross-Platform Service Hooks | `cli/cmdlemon/lemon_daemon.go`, `cli/cmdlemon/service_manager.go`, `cli/cmdlemon/daemon_settings.go`, `cli/cmdlemon/lemon_cmd.go` | `COMPLETED` |
| **Subtask-06** | End-to-End Test Suite Covering Edge Cases, Settings Toggle & Verification | `cli/tests/e2e/lemon_email_e2e_test.go`, `cli/tests/e2e/lemon_daemon_e2e_test.go`, `cli/tests/e2e/lemon_edgecases_e2e_test.go` | `COMPLETED` |

---

## Concurrency Topology & Wave Allocations

- **Wave 1 (Research & Feasibility, A=3)**:
  - Agent 1: Research SQL patterns and Ent Go blockers.
  - Agent 2: Research Antigravity Manager / Antigravity Tools multi-format email and ACK lifecycle.
  - Agent 3: Research Lemon daemon loops, cross-OS service hooks, and interval guards.
- **Wave 2 (Core Implementation, A=3, H=2)**:
  - Worker 1: Subtask-02 (Multi-Format Email Parser) + Subtask-03 (Immediate ACK Replier).
  - Worker 2: Subtask-04 (Command Execution Pipeline with AGM & AGY integration).
  - Worker 3: Subtask-05 (Lemon Daemon Polling Loop with $\ge 5$s interval guard & OS service managers).
- **Wave 3 (Integration, E2E Testing & Release)**:
  - Worker 1: Subtask-06 (End-to-End Test Suite).
  - Lead: Audit ledger finalization, documentation sync, version verification.
