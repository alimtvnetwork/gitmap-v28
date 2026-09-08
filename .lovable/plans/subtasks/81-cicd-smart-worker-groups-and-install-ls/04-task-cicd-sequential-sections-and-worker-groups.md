# Subtask 04: Sequential Section Execution & Intra-Section Worker Groups

## 1. Description
Enforce strict section-by-section sequential execution in `03-ai-scripts/06-cicd-local-runner.py` (Section 1 -> Section 2 -> Section 3 ...), where within each section tests execute concurrently via a worker group (`ThreadPoolExecutor`). Integrate the smart Go test runner into the pipeline, and run full verification.

## 2. Files to Modify
- `03-ai-scripts/06-cicd-local-runner.py`:
  - Ensure sequential section execution: Section 1 (Linters & AST) -> Section 2 (Compile Gates) -> Section 3 (Packaging) -> Section 4 (E2E Smoke) -> Section 5 (Go Smart Incremental Tests) -> Section 6 (Coverage Verification) -> Section 7 (Race Detection).
  - Within each section, execute tasks as a concurrent worker group using `ThreadPoolExecutor`.
  - In Section 5 (Go Smart Incremental Tests): run only changed/affected tests via the worker group, parse individual test results and durations, and update `.lovable/test-inventory.json`.
  - Run linters: `check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-error-management.py`.
  - Rebuild `gitmap.exe` and sync to all 4 executable paths.

## 3. Invariants
- Sections must never run concurrently with each other.
- Concurrency occurs exclusively *within* each section via its worker pool.
- Pipeline returns exit code 0 when all active gates pass.
