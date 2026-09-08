# Plan 81: Smart Incremental CI/CD Worker Groups, Code-to-Test Mapping & Install LS Enhancements

## 1. Overview & Context

This plan addresses two core developer experience and pipeline performance requirements:
1. **`gitmap install ls` Enhancements:**
   - Eliminate generic `"found"` placeholder and query live version strings (`v22.14.0`, `1.24.1`, `2.47.0`, etc.) with sub-second timeouts.
   - Include Gitmap CLI itself at the top of Core Tools with `constants.Version` (`v6.199.0`) and recent release tag history banner.
   - Register newly supported packages: `ToolGitmap`, `ToolComposer`, `ToolWpCli`, `ToolOpenVmTools`.
2. **CI/CD Pipeline Python Runner (`03-ai-scripts/06-cicd-local-runner.py`):**
   - **Section-by-Section Sequential Execution:** Run each section strictly in sequence (Section 1 -> Section 2 -> Section 3 -> ...), ensuring sections do not run parallelly with each other.
   - **Worker Group per Section:** Concurrently execute tests and gates within each section across a configurable worker group (`ThreadPoolExecutor`).
   - **Test Inventory Manifest (`.lovable/test-inventory.json`):**
     - Discover and catalog all existing tests (2,277+ Go unit tests across 633 test files, linter gates, and E2E suites).
     - Store exact test execution timings (duration in seconds/ms) and status in the JSON manifest.
   - **Code-to-Test Mapping & Impact-Driven Change Detection:**
     - Map each test to its target source file and target function.
     - Compute and record cryptographic hashes of both the target source code and the test function.
     - Only execute a test if its target function/code or the test itself changed since the last green run; otherwise skip with `[cached]`.

---

## 2. Task-Specific Rules & Invariants

1. **Strict Relative Paths:** Never use absolute paths or `file:///` in plans, code, or markdown.
2. **Coding Guidelines Adherence:** All Go functions must be <= 15 lines with a mandatory blank line before every return statement, affirmative booleans (`is*`, `has*`), and zero nested `if` blocks.
3. **Bounded Folders:** All logs, plans, and test caches are stored within `.lovable/`.
4. **Resilient Fallbacks:** Tool version detection in `gitmap install ls` must employ short context timeouts (<= 300ms) so command execution never hangs.

---

## 3. Subtask Decomposition

- [01-task-install-ls-enhancements.md](../subtasks/81-cicd-smart-worker-groups-and-install-ls/01-task-install-ls-enhancements.md): Add `gitmap`, `composer`, `wp-cli`, `open-vm-tools` to constants; implement live version detection in `installlist.go` with recent releases banner.
- [02-task-cicd-test-inventory-generator.md](../subtasks/81-cicd-smart-worker-groups-and-install-ls/02-task-cicd-test-inventory-generator.md): Implement test discovery in `06-cicd-local-runner.py` generating `.lovable/test-inventory.json` with execution timings.
- [03-task-cicd-code-to-test-change-detector.md](../subtasks/81-cicd-smart-worker-groups-and-install-ls/03-task-cicd-code-to-test-change-detector.md): Implement Go AST / function hashing, code-to-test mapping, and impact-based change skipping.
- [04-task-cicd-sequential-sections-and-worker-groups.md](../subtasks/81-cicd-smart-worker-groups-and-install-ls/04-task-cicd-sequential-sections-and-worker-groups.md): Enforce strict section-by-section sequential execution with intra-section worker groups and pass quality verification.

---

## 4. Verification Plan

1. Verify `gitmap install ls` displays actual versions for installed tools and `gitmap` at `v6.199.0` with releases banner.
2. Verify Go unit tests for install commands: `go test -v ./cmd -run "TestResolveToolStatus|TestInstaller"`.
3. Verify test inventory generation: run `python 03-ai-scripts/06-cicd-local-runner.py --inventory-only` and inspect `.lovable/test-inventory.json`.
4. Verify section sequential execution and worker group parallelism.
5. Verify incremental skip logic: running the runner twice skips unchanged tests and only runs tests whose functions changed.
6. Verify quality linters pass (`check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-error-management.py`).
