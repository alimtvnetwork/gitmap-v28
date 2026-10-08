# Subtask 01: CI/CD Pipeline Architecture, Desktop Mock Adapters & Telemetry Cache

> **Subtask ID:** Subtask-01  
> **Parent Plan:** `.ai-memory/plans/241-cicd-enhancements-help-text-grouping-and-suggestion-engine-design.md`  
> **Target Path:** `.ai-memory/plans/subtasks/241-cicd-enhancements-help-text-grouping-and-suggestion-engine-design/01-cicd-pipeline-and-desktop-mock-enhancements.md`  
> **Status:** PENDING  

---

## 1. Objectives & Scope
1. Architect injectable mock adapters (`RegistryAccessor`, `LinuxDesktopRunner`) for `cli/cmdos/` to allow headless unit testing in GitHub Actions runners.
2. Design staged-only fast-gate script (`03-ai-scripts/50-fastgate.py`) to reduce pre-commit verification latency from 25s to <1.5s.
3. Design Split-DB telemetry caching (`.gitmap/repodb/pipeline.db`) across GitHub Actions workflow runs to monitor duration regressions.
4. Establish duration-weighted dynamic test matrix sharding driven by `.ai-memory/test-inventory.json`.
5. Integrate closed-loop AI telemetry streaming from `gitmap pe -t` step failures to `store.AiAnalysisSplitDB`.
6. Define Go test binary precompilation stage (`go test -c`) to eliminate duplicate compilation overhead.

---

## 2. Detailed Technical Plan

### 2.1 Headless Desktop Mock Adapters
- In `cli/cmdos/os_dock_types.go`:
  - Define `RegistryAccessor` and `LinuxDesktopRunner` interfaces.
  - Implement default live operators and headless test mock operators.
- In `cli/cmdos/os_dock_linux_test.go` and `os_dock_win_test.go`:
  - Verify dock positioning logic without requiring active display server or elevated privileges.

### 2.2 Staged-Only Fast-Gate Script (`03-ai-scripts/50-fastgate.py`)
- Discover staged files via `git diff --cached --name-only --diff-filter=ACMR`.
- Filter staged list against file extension handlers:
  * `.go` -> `linter-scripts/check-relative-paths.py`, `check-nested-ifs.py`, `check-boolean-guidelines.py`.
  * `.py` -> syntax checks and PEP 8 style validation.
- Execute linters exclusively on staged files with `--files` flag.
- Add `--full` fallback flag for complete repository audit.

### 2.3 Split-DB Telemetry Cache
- Configure `.github/workflows/ci.yml` with `actions/cache@v4` pointing to `.gitmap/repodb/pipeline.db`.
- Query historical execution times in `03-ai-scripts/06-cicd-local-runner.py`.
- Emit warning if current step duration exceeds historical baseline by >50%.

---

## 3. Verification Criteria
- [ ] Mock interfaces compile and decouple OS commands from physical hardware.
- [ ] Fast-gate script executes in <1.5 seconds on typical 3-5 file staged sets.
- [ ] Telemetry caching preserves `pipeline.db` records across CI workflow runs.
- [ ] Test matrix evenly balances execution load across parallel shards.
