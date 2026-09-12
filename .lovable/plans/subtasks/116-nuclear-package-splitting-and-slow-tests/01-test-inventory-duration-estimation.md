# Subtask 01: Test Inventory Duration Estimation

## Objective
Analyze all 2,498 tests in `.lovable/test-inventory.json` and compute realistic duration estimates and tier classifications based on AST inspection.

## Steps
1. Inspect AST of all test files in `.lovable/test-inventory.json`.
2. Tests using `exec.Command`, `git` subprocesses, `time.Sleep`, or network sockets are classified as heavy (`duration_sec: 1.5 - 15.0`).
3. Pure in-memory unit tests are classified as fast (`duration_sec: 0.005`).
4. Update `.lovable/test-inventory.json` with the computed durations and summary statistics.
5. Verify JSON validity and structure.
