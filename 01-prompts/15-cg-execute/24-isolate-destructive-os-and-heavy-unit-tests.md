# Unit Test Isolation & Destructive OS Prevention — Coding Guideline Execution (must follow)

Trigger Keywords & Aliases: `cg-isolate-os-tests`, `cg-mock-destructive`, `cg-execute os-tests`, `isolate destructive tests`, `mock os shutdown`, `hermetic test isolation`

> **Prompt Version:** 2.2.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, audit, refactor, and verify repository-wide unit tests to ensure that tests NEVER trigger real OS shutdown, reboot, power-off, system modifications, or heavy unmocked system calls, enforcing injectable executors and mock duration verification across all test suites. Keep functions <= 8–15 lines, enforce affirmative booleans, and defer build verification strictly to the final step without running intermediate tests or builds.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

You MUST execute this task via a strict 3-Phase pipeline governed by the N-step budget. Do not skip steps.

```text
N = 200  (Total self-loop steps budget, read-only after initialization)
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Deep Scan Codebase, Map Violations in .lovable/plans/pending/, Subtasks)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Subtasks, Injectable Executors, Hermetic Mocks, Defer Restorations)
```

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A - Discovery & Inventory): Deeply scan the target codebase using fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) to inventory all direct OS commands, system altering functions, or unmocked filesystem cleaners.
2. [ ] /goal Phase 1 (Step B - Violation Ledger Mapping): Identify every site where `exec.Command` or raw system calls (`shutdown`, `reboot`, `poweroff`, `apt-get`, `Get-WindowsUpdate`, `os.Remove`) could be reached from test suites.
3. [ ] /goal Phase 1 (Step C - Master Plan Generation): Write the master architectural specification into `.lovable/plans/pending/xx-isolate-destructive-os-and-heavy-unit-tests.md` with an exhaustive Violation Ledger table.
4. [ ] /goal Phase 1 (Step D - Subtask Decomposition): Decompose the plan into granular subtasks in `.lovable/plans/subtasks/xx-isolate-destructive-os-and-heavy-unit-tests/`.
5. [ ] /goal Phase 2 (Step A - Injectable Executor Introduction): Refactor target production code to introduce injectable executors (`DefaultOSActionExecutor`, `FileRemover`, `OSCommandRunner`, `TempDirResolver`) and parameter structs.
6. [ ] /goal Phase 2 (Step B - Hermetic Mock Testing): Refactor unit tests to swap the executor with a mock and assert captured parameters using `defer` restoration blocks.
7. [ ] /goal Phase 2 (Step C - Fast Duration Math): Verify countdown timers and delays use duration arithmetic and short simulated ticks (1s, 2s) without sleeping or arming host OS power timers.
8. [ ] /goal Phase 2 (Step D - Function & File Size Compliance): Ensure all refactored functions remain <= 8 lines of body logic (hard cap of <= 15 lines) and files remain under 100 lines.
9. [ ] /goal Phase 2 (Step E - Boolean & Style Conventions): Enforce affirmative boolean naming (`is*`, `has*`), zero explicit `== true`, and zero negative booleans.
10. [ ] /goal Phase 2 (Step F - Banned Intermediate Verification): DO NOT run unit tests (`go test`, `pytest`, `npm test`) during intermediate micro-refactoring steps.
11. [ ] /goal Phase 2 (Step G - Final Step Build Verification): At the conclusion of all refactoring subtasks, run targeted syntax and quality gate checks to verify clean compilation.
12. [ ] /goal Phase 3 (Step A - Task Consolidation): Consolidate completed subtasks into `.lovable/plans/completed/`, delete granular subtask files, and update `.lovable/plans/01-index.md`.
13. [ ] /goal Phase 3 (Step B - Atomic Git Commit & Push): Stage all modified files (`git add -A`), commit them in a single clean grouped atomic commit, and push to git.
14. [ ] /learn Ingest `spec/02-coding-guidelines/` for domain-specific architectural specifications.
15. [ ] /learn Ingest `spec/03-error-manage/` for AppError wrapping.
16. [ ] /learn Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.

---

## 1. Core Principles of Destructive OS Test Isolation

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

## 2. Polyglot Isolation Patterns

### Pattern 1: Power Cancellation & Abort Command (Go)

```go
// ❌ ANTI-PATTERN: Executes real shutdown /a or shutdown -c on host system
func CancelSchedulePowerCLI(action OSActionType) error {
    _ = CancelActivePowerSchedule()
    exe, cmdArgs := BuildCancelOSActionCommand()
    cmd := exec.Command(exe, cmdArgs...)
    _ = cmd.Run()
    return nil
}
```

```go
// ✅ COMPLIANT PATTERN: Routes through injectable DefaultOSActionExecutor
func CancelSchedulePowerCLI(action OSActionType) error {
    _ = CancelActivePowerSchedule()
    exe, cmdArgs := BuildCancelOSActionCommand()
    params := buildCancelActionParams(exe, cmdArgs)
    if err := DefaultOSActionExecutor(params); err != nil {
        return err
    }
    fmt.Printf("✓ Canceled active scheduled %s.\n", action)
    return nil
}

func buildCancelActionParams(exe string, cmdArgs []string) OSActionParams {
    return OSActionParams{
        Action:     OSActionCancel,
        Executable: exe,
        Args:       cmdArgs,
    }
}
```

### Pattern 2: Ephemeral Cache Purge Mocking (Go)

```go
// ❌ ANTI-PATTERN: Direct unmocked os.Remove in production logic
func purgeCategoryFiles(paths []string) (int, int64) {
    for _, p := range paths {
        _ = os.Remove(p)
    }
}
```

```go
// ✅ COMPLIANT PATTERN: Injectable FileRemover allows hermetic testing
type FileRemover func(filePath string) (int, int64)
var defaultFileRemover FileRemover = removeSingleFileSafely

func purgeCategoryFiles(paths []string) (int, int64) {
    for _, p := range paths {
        count, bytes := defaultFileRemover(p)
    }
}
```

### Pattern 3: Temporary Directory Sweep Isolation (Go)

```go
// ❌ ANTI-PATTERN: Sweeps host system /tmp or C:\Users\...\AppData\Local\Temp directly
func CleanTempDirectories(opts CleanOptions) CleanResult {
    dirs := resolveTempDirectories() // Returns real OS temp dirs
    for _, dir := range dirs {
        os.RemoveAll(dir)
    }
}
```

```go
// ✅ COMPLIANT PATTERN: Injectable TempDirResolver defaults to OS dirs, overridden with t.TempDir() in tests
type TempDirResolver func() []string
var defaultTempDirResolver TempDirResolver = resolveTempDirectories

func CleanTempDirectories(opts CleanOptions) CleanResult {
    dirs := defaultTempDirResolver()
    for _, dir := range dirs {
        // Safe execution against resolved dirs
    }
}
```
