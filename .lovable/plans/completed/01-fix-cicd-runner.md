# Consolidated Plan: Fix CI/CD Runner Output, ETA, and Targeted Package/File Testing
Completed in 2 Steps/Loops.

## Goals Achieved
1. **OS Temp Directory Migration**: Updated `03-ai-scripts/06-cicd-local-runner.py` to store test execution logs, states, and errors directly in `Path(tempfile.gettempdir()) / ".lovable" / "cicd"`.
2. **Multiple Custom Output Paths**: Added support for `--output-paths` taking multiple file destinations for JSON reports.
3. **Accurate Parallel ETA Calculation**: Refactored `calculate_total_eta` to compute realistic execution estimates by accounting for parallel workers per batch.
4. **Hashed Run Directories & Isolated Failure Logs**:
   - Every run now creates a unique directory tagged with an 8-character hash and timestamp (`<run_hash>-<timestamp>`).
   - Any failing test creates a dedicated file inside `<session_dir>/failed_tests/<sanitized_name>.log` containing the test name, command, exit code, and full stack trace.
   - The AI Remediation Banner displays exact absolute paths to these failed test logs.
5. **Targeted Package, Code File Path, and Code File Name Testing**:
   - Added `--pkg`, `--package`, `-p`, `--target-file`, `--file`, and positional argument support.
   - Allows running tests by:
     - Code file path: `python 03-ai-scripts/06-cicd-local-runner.py gitmap/cmd/mkdir.go`
     - Code file name: `python 03-ai-scripts/06-cicd-local-runner.py -p mkdir.go`
     - Short package name: `python 03-ai-scripts/06-cicd-local-runner.py -p store`
     - Full package module path: `python 03-ai-scripts/06-cicd-local-runner.py --pkg github.com/alimtvnetwork/gitmap-v28/gitmap/cmd`
   - Automatically isolates execution to `Go Smart Incremental Tests` and maps code files/functions to their test suite.
6. **Strict Coding Guidelines Compliance (Nested-If Flattening)**:
   - Completely flattened nested if statements in `gitmap/cmd/mkdir.go` (`makeDirSingle`, `makeDirDeep`, `touchFilePrepareParent`).
   - Completely flattened nested if statements in `gitmap/store/store.go` (`lockDBIfNotMem`, `releaseLockIfNotMem`).
   - Verified `check-nested-ifs.py`, `check-enum-and-boolean.py`, and `check-error-management.py` all pass with 0 violations.
