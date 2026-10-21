# GitMap Testing & Schedule Management Specification Suite

**Specification Slug:** `gitmap-testing-and-schedule-management`  
**Root Path:** `02-spec/21-app/gitmap-testing-and-schedule-management/`  
**Parent Plan:** `.ai-memory/plans/gitmap-testing-and-schedule-management.md`  
**Subtasks Path:** `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/`  
**Status:** APPROVED  

---

## 1. Executive Summary

This specification suite governs the architecture, CLI contracts, cross-platform power scheduling, Dockerized Linux testing strategy, safety test isolation (`cg-isolate-os-tests`), mail and schedule E2E verification matrices, and master acceptance criteria for GitMap.

The primary directive is ensuring that all newly implemented GitMap capabilities—specifically cross-OS system power scheduling (`shutdown` and `restart`), scheduled email dispatch, multi-mailbox inspection, folder watcher event triggers, privacy masking, and zero-CDN embedded web dashboards—are comprehensively verified across Windows, Linux, and macOS environments without endangering host workstations or destabilizing production workflows.

```
Subsystem Interoperability Architecture:
┌────────────────────────────────────────────────────────────────────────┐
│                        GitMap CLI Architecture                         │
├───────────────────────────────────┬────────────────────────────────────┤
│       System Power Scheduling     │        Mail & Notification Core    │
│  - Windows: shutdown /s /t        │  - Folder Watcher Event Lifecycle  │
│  - Linux:   sudo poweroff / reboot│  - Split-DB SQLite Storage & Logs  │
│  - Darwin:  sudo shutdown -h / -r │  - Privacy-Preserving Address Mask │
├───────────────────────────────────┼────────────────────────────────────┤
│       Safety & Test Isolation     │        Embedded Web Dashboards     │
│  - Injectable OSActionExecutor    │  - /mail/ui: Zero-CDN Composer     │
│  - Mock Doubles (Zero Host Boot)  │  - /schedule/ui: Settings Manager  │
│  - Docker Container Harness       │  - HTTP 405 Method Guard Security  │
└───────────────────────────────────┴────────────────────────────────────┘
```

---

## 2. Master Specification Document Catalog

This catalog indexes all six canonical specifications comprising the test and schedule management architecture:

| Spec # | Specification Document | Focus Area & Description |
|---|---|---|
| **01** | [01-architecture-and-cross-os-overview.md](01-architecture-and-cross-os-overview.md) | High-level architectural overview, cross-OS power scheduling topology, duration parsing (`1:45hr`, `2h`, `30m`), hybrid countdown lifecycles, and subagent swarm concurrency. |
| **02** | [02-cross-os-power-command-contracts.md](02-cross-os-power-command-contracts.md) | Concrete CLI command contracts for `gitmap schedule shutdown` and `gitmap schedule restart`, flag definitions (`--dry-run`, `--force`, `--json`), and OS command translation. |
| **03** | [03-docker-and-linux-container-strategy.md](03-docker-and-linux-container-strategy.md) | Docker environment discovery findings, pure Go cross-compilation (`GOOS=linux`), and hermetic containerized validation via `Dockerfile.test` targeting Alpine Linux. |
| **04** | [04-safety-and-os-test-isolation.md](04-safety-and-os-test-isolation.md) | Enforces the `cg-isolate-os-tests` standard, preventing accidental machine restarts during automated runs via injectable `OSActionExecutor` and deterministic test doubles. |
| **05** | [05-mail-and-schedule-e2e-matrix.md](05-mail-and-schedule-e2e-matrix.md) | Feature verification matrix covering folder watcher mutations (`EventCreate`, `EventModify`, `EventDelete`), Split-DB CRUD, `--file`/`--html` sending, and flag variations. |
| **06** | [06-verification-and-acceptance-criteria.md](06-verification-and-acceptance-criteria.md) | Numbered quality gates (`G1` through `G9`), pass invariants, automated execution command recipes, and release sign-off checklist. |

---

## 3. Core Architectural Principles & Invariants

1. **Relative Path Sovereignty:** All specification documents, subtask plans, and code references use strictly relative Git paths (`02-spec/...`, `.ai-memory/...`, `cli/...`). Absolute paths and file URI schemes are strictly prohibited.
2. **Safety Isolation Invariant (`cg-isolate-os-tests`):** Unit and integration test suites must NEVER invoke native OS shutdown or reboot binaries directly on the host machine. Mock executors capture parameters deterministically.
3. **Pure Go Cross-OS Portability:** Linux and Darwin builds must compile cleanly without CGO or Windows-specific symbol leaks (`GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./...`).
4. **Offline Zero-CDN Operation:** Embedded Web UI dashboards (`/mail/ui` and `/schedule/ui`) operate entirely offline without external CDN assets, scripts, or fonts.
5. **Privacy Governance:** Email addresses rendered in console output, logs, or JSON structures are sanitized via boundary masking (`u***r@d***n.com`) unless explicit unmasking is authorized.
6. **Structural Hygiene Budget:** All specification and plan files adhere to the strict $\le 300$ line count limit and lowercase hyphenated naming standard.

---

## 4. Associated Subtasks & Execution Roadmap

The implementation and verification of this specification suite is partitioned into six atomic subtasks:

- **Subtask-01:** Git Sync, Preflight & Task Initialization (`.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/01-subtask-git-sync-and-preflight.md`) — **DONE**
- **Subtask-02:** Architecture & Cross-OS Test Specification Authoring (`02-spec/21-app/gitmap-testing-and-schedule-management/`) — **DONE**
- **Subtask-03:** Cross-OS Schedule Shutdown and Restart Test Suite (`.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/03-subtask-cross-os-schedule-tests.md`) — **PENDING**
- **Subtask-04:** Docker Environment & Linux Containerized Verification (`.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/04-subtask-docker-linux-container-verification.md`) — **PENDING**
- **Subtask-05:** Comprehensive Verification of New Mail & Scheduling Features (`.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/05-subtask-new-features-e2e-verification.md`) — **PENDING**
- **Subtask-06:** Master Index Registration, Quality Gates & Atomic Push (`.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/06-subtask-index-quality-gate-and-push.md`) — **PENDING**
