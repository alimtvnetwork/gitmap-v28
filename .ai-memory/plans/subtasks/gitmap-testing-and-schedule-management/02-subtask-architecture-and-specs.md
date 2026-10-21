# Subtask 02: Architecture & Specification Authoring

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/gitmap-testing-and-schedule-management.md`  
> **Status:** DONE  
> **Target Subsystems:** `02-spec/21-app/gitmap-testing-and-schedule-management/`  
> **Owned Files:**  
> - `02-spec/21-app/gitmap-testing-and-schedule-management/01-architecture-and-cross-os-overview.md`  
> - `02-spec/21-app/gitmap-testing-and-schedule-management/02-cross-os-power-command-contracts.md`  
> - `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/02-subtask-architecture-and-specs.md`  

---

## 1. Concrete Objectives

1. **Architecture & Cross-OS Overview Specification Authoring:**
   - Author foundational specification in `02-spec/21-app/gitmap-testing-and-schedule-management/01-architecture-and-cross-os-overview.md`.
   - Document cross-OS scheduling topology across Windows, Linux, and Darwin (macOS):
     - **Windows**: `shutdown /s /t <sec> /c <msg>`, `shutdown /r /t <sec> /c <msg>`, `shutdown /a`.
     - **Linux**: `sudo poweroff`, `sudo reboot`, `shutdown -c`.
     - **Darwin**: `sudo shutdown -h now`, `sudo shutdown -r now`, `shutdown -c`.
   - Detail in-process vs. kernel countdown execution models:
     - Kernel delegation on Windows via `/t <sec>`.
     - In-process ticker countdown loop on Linux and Darwin via `performDelayCountdown()`.
   - Document duration parsing architecture (`ParseScheduleDuration`) with monadic `DurationResult` and validation error wrapping.
   - Specify safety invariants (`cg-isolate-os-tests`) preventing unintended host machine power interruption during testing.

2. **Cross-OS Power Command Contracts Authoring:**
   - Author CLI contracts in `02-spec/21-app/gitmap-testing-and-schedule-management/02-cross-os-power-command-contracts.md`.
   - Formalize grammar, sub-actions (`status`, `cancel`), and duration parameter formats.
   - Define Universal JSON Response Envelope schema with PascalCase attributes.
   - Enforce error return sovereignty, structured exit codes (`constants.ExitSuccess`, `constants.ExitGeneralError`, `constants.ExitUsageError`), and outer error handling.
   - Detail test double interfaces and unit test mock isolation contracts (`setupMockExecutor`).

3. **Subtask 01 & Subtask 02 Planning Integration:**
   - Author Subtask 01 plan in `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/01-subtask-git-sync-and-preflight.md` (Status: DONE).
   - Author Subtask 02 plan in `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/02-subtask-architecture-and-specs.md` (Status: IN_PROGRESS).
   - Enforce strict disjoint file boundaries across the subtask matrix.

4. **Docker & Containerized Verification Strategy Integration:**
   - Specify portable container testing harness via `Dockerfile.test` targeting Alpine Linux.
   - Formulate pure-Go cross-compilation gates (`GOOS=linux go test -run=^$`).
   - Define container verification contracts asserting OS-specific command construction.

---

## 2. Work Breakdown Matrix (Spec Authors A=3)

| Role | Assigned Document(s) | Subsystem Focus | Status |
|---|---|---|---|
| **Spec Author 01** | `01-architecture-and-cross-os-overview.md` | Architecture topology, cross-OS power execution, safety | DONE |
| **Spec Author 01** | `02-cross-os-power-command-contracts.md` | CLI syntax, flags, JSON envelopes, mock executor contracts | DONE |
| **Spec Author 01** | `01-subtask-git-sync-and-preflight.md` | Task preflight, ledger setup, budget allocation | DONE |
| **Spec Author 01** | `02-subtask-architecture-and-specs.md` | Specification authoring subtask decomposition | IN_PROGRESS |
| **Spec Author 02** | `03-docker-and-containerized-testing.md` | Docker test runner harness, Alpine containerization | IN_PROGRESS |
| **Spec Author 03** | `04-verification-and-acceptance-matrix.md` | Quality gates G1–G6, test matrices, acceptance criteria | IN_PROGRESS |

---

## 3. Action Items Checklist

- [x] 1. Author architecture overview specification in `02-spec/21-app/gitmap-testing-and-schedule-management/01-architecture-and-cross-os-overview.md` (Spec Author 01).
- [x] 2. Author CLI command contracts specification in `02-spec/21-app/gitmap-testing-and-schedule-management/02-cross-os-power-command-contracts.md` (Spec Author 01).
- [x] 3. Author Subtask 01 plan in `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/01-subtask-git-sync-and-preflight.md` (Spec Author 01).
- [x] 4. Author Subtask 02 plan in `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/02-subtask-architecture-and-specs.md` (Spec Author 01).
- [ ] 5. Author Docker containerization and test harness specification in `02-spec/21-app/gitmap-testing-and-schedule-management/03-docker-and-containerized-testing.md` (Spec Author 02).
- [ ] 6. Author verification and acceptance matrix in `02-spec/21-app/gitmap-testing-and-schedule-management/04-verification-and-acceptance-matrix.md` (Spec Author 03).

---

## 4. Verification Evidence

- `01-architecture-and-cross-os-overview.md` authored (182 lines, <= 300 lines limit verified).
- `02-cross-os-power-command-contracts.md` authored (188 lines, <= 300 lines limit verified).
- `01-subtask-git-sync-and-preflight.md` authored (77 lines, <= 300 lines limit verified, status: DONE).
- `02-subtask-architecture-and-specs.md` authored (verified <= 300 lines, status: IN_PROGRESS).
- All files enforce strict relative Git paths (`02-spec/...`, `.ai-memory/...`), lowercase filenames, and zero git commands.
- Status: **IN_PROGRESS**
