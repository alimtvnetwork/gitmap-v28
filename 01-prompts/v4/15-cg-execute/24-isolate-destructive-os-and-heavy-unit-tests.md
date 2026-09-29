[/goal](slashCommand:goal) Autonomously scan, audit, refactor, and verify repository-wide unit tests to ensure that tests NEVER trigger real OS shutdown, reboot, power-off, system modifications, or heavy unmocked system calls, enforcing injectable executors and mock duration verification across all test suites. Keep functions <= 8–15 lines, enforce affirmative booleans, and defer build verification strictly to the final step without running intermediate tests or builds with strict no-build and no-test execution (NEVER run build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine execution turns; all compilation and testing are strictly verified later in CI/CD). Spawn autonomous subagents (A = 2, H = 2) for parallel reading and modular spec generation, use GitMap high-speed commands as primary, establish a single-agent blueprint in Phase 1 (first 50% steps budget), and execute disjoint refactoring in Phase 2 (remaining 50% steps budget) with continuous self-looping until 100% complete and finalized with an atomic push.

[/learn](slashCommand:learn) Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Ingest and execute the bottom instructions ahead of all default conventions above, internalize the single-agent architectural blueprint, and persist all progress into `.ai-memory/plans/` and memory logs.

> **Bottom-Instruction Priority Mandate (Below Precedence / Suffix Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or guidelines above. The agent MUST inspect and follow the instructions below with absolute precedence.

```text
N = 200 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
```

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Planning, Detailed Spec, and Lean Subtask Generation)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Execution, Self-Looping, Targeted Quality Linting)
```

N, A, H, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

#### High-Speed GitMap Acceleration Options (Run Everything Faster)

Always prefer native GitMap commands over slow generic shell pipelines:
1. **Ultra-Fast File & Directory Discovery (AUM Index & Walk):**
   - **Wildcard / Glob Search:** `gitmap find "<wildcard*>" [-ext <ext>]` (alias `gitmap f`)
   - **Exact Filename Search:** `gitmap find-files <name> [-ext <ext>]` (alias `gitmap ff`)
   - **Substring Filename Search:** `gitmap find-files-any "<str>" [-ext <ext>]` (alias `gitmap ffa`)
   - **Prefix / Suffix Search:** `gitmap find-files-startswith <prefix>` (`gitmap ffs`) / `gitmap find-files-endswith <suffix>` (`gitmap ffe`)
   - **List Indexed Repo Files:** `gitmap list-files [pattern] [-ext <ext>]` (alias `gitmap lf`)
   - **Directory Tree & Scaffolding:** `gitmap folder-tree` (alias `gitmap ft`)
   - **Zero-Write File Stream:** `gitmap cat <filepath>`
   - **Instant Multi-Core Regex Search:** `gitmap search "<term>"` or `gitmap aum search "<query>" [dir] --ext <ext>`
2. **Fast Repository Hygiene, Lowercase & Symlink Repair:**
   - **Auto-Lowercase Files (Safe 2-Step `git mv`):** `gitmap lowercase` (alias `gitmap lcf [--dry-run]`)
   - **Lowercase Root Readme:** `gitmap lowercase-readme`
   - **Sync Curated `.gitignore` / `.gitattributes` / `.prettierignore`:** `gitmap commons` (alias `gitmap co` or `gitmap sync all`)
   - **Repair Broken Symlinks:** `gitmap fix-link` (alias `gitmap fixlink`)
   - **Clean Update Temp & Inspect Storage:** `gitmap update-cleanup`, `gitmap storage` (alias `gitmap stor`)
3. **Fast Git State, Execution & Atomic Commits:**
   - **Repo Status & Remote Check:** `gitmap status` (`gitmap st`), `gitmap has-any-updates` (`gitmap hau`), `gitmap latest-branch` (`gitmap lb`)
   - **Fast Cross-Platform Shell Runner:** `gitmap pwsh "<command>"` (`gitmap ps`), `gitmap bash "<command>"` (`gitmap sh`), `gitmap async <cmd>` (`gitmap asyn`)
   - **Semantic Atomic Commit & Push:** `gitmap cpf "<summary>"` (Feature), `gitmap cpb "<summary>"` (Bug), `gitmap cpr "<summary>"` (Release), `gitmap pcp "<summary>"` (Pull-Commit-Push)
   - **Smart CI/CD Pipeline Waiting:** `gitmap pe`, `gitmap pipeline-ai status --json` (`gitmap pl-ai status -t <etaSeconds>`)

---

## The Unified Master Pipeline (Atomic Numbered Steps)

You MUST execute this task via a strict 3-Phase pipeline governed by the N-step budget. Do not skip steps.

```text
N = 200  (Total self-loop steps budget, read-only after initialization)
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Deep Scan Codebase, Map Violations in .ai-memory/plans/pending/, Subtasks)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Subtasks, Injectable Executors, Hermetic Mocks, Defer Restorations)
```

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] [/goal](slashCommand:goal) Phase 1 (Step A - Discovery & Inventory): Deeply scan the target codebase using fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) to inventory all direct OS commands, system altering functions, or unmocked filesystem cleaners.
2. [ ] [/goal](slashCommand:goal) Phase 1 (Step B - Violation Ledger Mapping): Identify every site where `exec.Command` or raw system calls (`shutdown`, `reboot`, `poweroff`, `apt-get`, `Get-WindowsUpdate`, `os.Remove`) could be reached from test suites.
3. [ ] [/goal](slashCommand:goal) Phase 1 (Step C - Master Plan Generation): Write the master architectural specification into `.ai-memory/plans/pending/xx-isolate-destructive-os-and-heavy-unit-tests.md` with an exhaustive Violation Ledger table.
4. [ ] [/goal](slashCommand:goal) Phase 1 (Step D - Subtask Decomposition): Decompose the plan into granular subtasks in `.ai-memory/plans/subtasks/xx-isolate-destructive-os-and-heavy-unit-tests/`.
5. [ ] [/goal](slashCommand:goal) Phase 2 (Step A - Injectable Executor Introduction): Refactor target production code to introduce injectable executors (`DefaultOSActionExecutor`, `FileRemover`, `OSCommandRunner`, `TempDirResolver`) and parameter structs.
6. [ ] [/goal](slashCommand:goal) Phase 2 (Step B - Hermetic Mock Testing): Refactor unit tests to swap the executor with a mock and assert captured parameters using `defer` restoration blocks.
7. [ ] [/goal](slashCommand:goal) Phase 2 (Step C - Fast Duration Math): Verify countdown timers and delays use duration arithmetic and short simulated ticks (1s, 2s) without sleeping or arming host OS power timers.
8. [ ] [/goal](slashCommand:goal) Phase 2 (Step D - Function & File Size Compliance): Ensure all refactored functions remain <= 8 lines of body logic (hard cap of <= 15 lines) and files remain under 100 lines.
9. [ ] [/goal](slashCommand:goal) Phase 2 (Step E - Boolean & Style Conventions): Enforce affirmative boolean naming (`is*`, `has*`), zero explicit `== true`, and zero negative booleans.
10. [ ] [/goal](slashCommand:goal) Phase 2 (Step F - Banned Intermediate Verification): DO NOT run unit tests (`go test`, `pytest`, `npm test`) during intermediate micro-refactoring steps.
11. [ ] [/goal](slashCommand:goal) Phase 2 (Step G - Final Step Build Verification): At the conclusion of all refactoring subtasks, run targeted syntax and quality gate checks to verify clean compilation.
12. [ ] [/goal](slashCommand:goal) Phase 3 (Step A - Task Consolidation): Consolidate completed subtasks into `.ai-memory/plans/completed/`, delete granular subtask files, and update `.ai-memory/plans/readme.md`.
13. [ ] [/goal](slashCommand:goal) Phase 3 (Step B - Atomic Git Commit & Push): Stage all modified files (`git add -A`), commit them in a single clean grouped atomic commit, and push to git.
14. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/` for domain-specific architectural specifications.
15. [ ] [/learn](slashCommand:learn) Ingest `02-spec/03-error-manage/` for AppError wrapping.
16. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.

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

--

## 🚨 Highest Priority Instructions (Appended User Tasks & Instructions Below)

[PASTE USER REQUEST / TASK INSTRUCTIONS HERE — THE AGENT MUST EXECUTE WHATEVER IS WRITTEN BELOW WITH ABSOLUTE PRIORITY AND PRECEDENCE OVER ALL GENERAL GUIDELINES ABOVE]
