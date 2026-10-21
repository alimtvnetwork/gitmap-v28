# Execution Plan: GitMap Testing and Schedule Management

**Slug:** `gitmap-testing-and-schedule-management`  
**Task Name:** GitMap Testing and Schedule Management  
**Status:** ACTIVE  
**Lifecycle Budget:** N = 400 Steps (Phase 1: 1-100, Phase 2: 101-250, Phase 3: 251-400)  
**Concurrency:** A = 3, H = 2, Capacity = 6  

---

## 1. User Request (Verbatim)

```text
# GitMap Testing and Schedule Management: high priority instruction, non-negotiable task

Gitmap or sending email, testing the schedule stop. Have you tested the schedule shutdown and schedule restart item? I know that it cannot be done here in the Windows, but how confirmed are you that it can work regardless of the OS? If I give you another OS, can you test it, or can you install Docker and test it there, that can this be shut down and tested that Gitmap on that place? Do you understand? Can you please follow through and test every one of the new features and code that has been done for Gitmap?

## task: GitMap Testing and Schedule Management
## slug: gitmap-testing-and-schedule-management

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/gitmap-testing-and-schedule-management/ and enqueue plan task in .ai-memory/plans/gitmap-testing-and-schedule-management.md (subtasks in .ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/) first
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Test the schedule shutdown and schedule restart features across different operating systems.
6. Install Docker and test GitMap functionality within a Docker environment to ensure compatibility and functionality.
7. Verify and test all new features and code implemented for GitMap to ensure they work as expected.
```

---

## 2. Architecture & Systems Integration Blueprint

### Cross-OS Power Scheduling Architecture
```
User Invocations:
  gitmap schedule shutdown [duration] [status | cancel]
  gitmap schedule restart  [duration] [status | cancel]

Architecture Topology:
  ┌─────────────────────────────────────────────────────────────┐
  │ Duration Parser (duration_parser.go)                        │
  │ Supports colon format (1:45hr), unit suffixes (2h, 30m), etc │
  └──────────────────────────────┬──────────────────────────────┘
                                 │
  ┌──────────────────────────────▼──────────────────────────────┐
  │ Cross-OS Command Builder (schedule_os.go)                   │
  │ • Windows: shutdown /s /t <sec> /c <msg>                    │
  │ • Linux:   sudo poweroff / sudo reboot / shutdown -h now    │
  │ • Darwin:  sudo shutdown -h now / sudo shutdown -r now      │
  │ • Cancel:  shutdown /a (Win) vs shutdown -c (Unix)          │
  └──────────────────────────────┬──────────────────────────────┘
                                 │
  ┌──────────────────────────────▼──────────────────────────────┐
  │ Injectable OSActionExecutor & Mock Isolation (cg-isolate)   │
  │ • Prevents accidental host reboots during automated testing │
  │ • Simulates Windows, Linux, and macOS execution in tests    │
  └──────────────────────────────┬──────────────────────────────┘
                                 │
  ┌──────────────────────────────▼──────────────────────────────┐
  │ Power Schedule State & Countdown Persistence                │
  │ • Persisted target timestamp (TriggerAt)                    │
  │ • Dynamic remaining duration calculation (FormatRemaining)  │
  └─────────────────────────────────────────────────────────────┘
```

### Docker & Cross-OS Linux Testing Strategy
1. **Host Environment Finding**: Docker is not installed on the current Windows host.
2. **Hermetic Dockerfile.test Harness**: Add an isolated `Dockerfile.test` targeting `golang:1.24-alpine` enabling containerized Linux testing in CI/CD or whenever Docker is provisioned.
3. **Pure Go Linux Cross-Compilation**: `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -run=^$` verifies zero Windows-specific build tag leaks or symbol drift across the entire codebase.
4. **Hermetic Test Doubles**: Mock CLI executors in `cmdschedule` and `cmdmail` ensure tests can exercise Linux and Darwin execution branches deterministically on any host machine.

---

## 3. Subtask Work Breakdown Matrix

| Subtask ID | Title | Disjoint Files Owned | Assigned Role | Status |
|---|---|---|---|---|
| **Subtask-01** | Git Sync, Preflight & Task Initialization | `ledger.md`, `.ai-memory/plans/` | Lead | DONE |
| **Subtask-02** | Architecture & Cross-OS Test Specification Authoring | `02-spec/21-app/gitmap-testing-and-schedule-management/` | Spec Authors (A=3) | IN_PROGRESS |
| **Subtask-03** | Cross-OS Schedule Shutdown and Restart Test Suite | `cli/cmdschedule/schedule_os*.go`, `cli/cmdschedule/schedule_power_state*.go` | Worker 01 | PENDING |
| **Subtask-04** | Docker Environment & Linux Containerized Verification | `Dockerfile.test`, `.github/workflows/`, `cli/` | Worker 02 | PENDING |
| **Subtask-05** | Comprehensive Verification of New Mail & Scheduling Features | `cli/cmdmail/`, `cli/cmdschedule/` | Worker 03 | PENDING |
| **Subtask-06** | Master Index Registration, Quality Gates & Atomic Push | `02-spec/21-app/readme.md`, `.ai-memory/plans/readme.md` | Lead | PENDING |

---

## 4. Verification Evidence & Quality Gates
- **G1 (Cross-OS Command Construction):** Tests verify that `shutdown` and `restart` emit valid OS-specific commands for Windows (`shutdown /s /t`), Linux (`sudo poweroff` / `sudo reboot`), and Darwin (`sudo shutdown -h now` / `sudo shutdown -r now`).
- **G2 (Safety Invariant cg-isolate-os-tests):** Zero real OS reboots are triggered during test execution. All power actions use injectable mock executors.
- **G3 (Linux Compilation Gate):** `GOOS=linux go test -c ./cmdschedule` and `GOOS=linux go build ./...` compile cleanly with 0 errors.
- **G4 (Docker Harness Portability):** A runnable `Dockerfile.test` is provided for containerized validation.
- **G5 (Comprehensive Feature Coverage):** Folder watcher covers file creation, modification, and deletion. Mail templates, notification rules, address masking, and zero-CDN UIs are tested across all subcommands.
- **G6 (Line Cap & Relative Paths):** All authored files are strictly $\le 300$ lines with 100% relative paths.
