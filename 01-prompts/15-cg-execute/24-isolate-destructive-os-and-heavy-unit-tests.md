# CG-Execute: Isolate Destructive OS & Heavy Unit Tests — Coding Guideline 24

> **Prompt Version:** 2.1.0  
> **Synchronization:** Main Meta-Repo & Connected Workspaces  

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, audit, refactor, and verify repository-wide unit tests to ensure that tests NEVER trigger real OS shutdown, reboot, power-off, system modifications, or heavy unmocked system calls, enforcing injectable executors and mock duration verification across all test suites.

---

## 1. The Principle of Destructive OS Test Isolation

Unit tests must be fast, hermetic, and safe. Executing destructive OS commands or heavy operations directly inside unit tests causes machine instability, shuts down developer workstations, and breaks CI/CD runners.

### Core Mandates:
1. **Total Ban on Real OS Shutdown / Reboot**:
   - `exec.Command("shutdown", ...)`
   - `exec.Command("reboot", ...)`
   - `exec.Command("poweroff", ...)`
   - `exec.Command("init", "0")`
   - `exec.Command("systemctl", "poweroff"|"reboot")`
   Must NEVER be called directly in unit tests.
2. **Mandatory Injectable Executors**:
   All power actions, system cleanups, and heavy alterations must be decoupled behind an injectable function or interface:
   ```go
   type OSActionExecutor func(params OSActionParams) error
   var DefaultOSActionExecutor OSActionExecutor = executeNativeOSAction
   ```
3. **Hermetic Mock Testing**:
   Unit tests must substitute `DefaultOSActionExecutor` using a `defer` restoration block:
   ```go
   var captured OSActionParams
   oldExecutor := DefaultOSActionExecutor
   defer func() { DefaultOSActionExecutor = oldExecutor }()

   DefaultOSActionExecutor = func(params OSActionParams) error {
       captured = params
       return nil
   }
   ```
4. **Fast Duration Math (1s/2s Checks)**:
   Tests for countdown timers or delayed triggers must test duration parsing, math calculation, and short simulated delays (1s, 2s) with mock tick callbacks rather than blocking or arming the host OS power manager.

---

## 2. Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase using fast Python discovery tools to identify any unit tests invoking raw OS power commands, heavy shell commands, or unmocked filesystem cleaners.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.lovable/plans/pending/` with an exhaustive Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose the plan into granular subtasks in `.lovable/plans/subtasks/`.
4. [ ] /goal Phase 2 (Step A): Refactor target production code to introduce injectable executors (`DefaultOSActionExecutor`) and parameter structs.
5. [ ] /goal Phase 2 (Step B): Refactor unit tests to swap the executor with a mock and assert captured parameters.
6. [ ] /goal Phase 2 (Step C): Verify all modified files adhere to <= 15 line function limits and affirmative boolean rules.
7. [ ] /learn Ingest `spec/02-coding-guidelines/` for domain-specific architectural specifications.
8. [ ] /learn Ingest `spec/03-error-manage/` for AppError wrapping.
