# Learned Memory: Plan 129 — Smart Test Runner, Dual-Queue Worker Pools, and Dynamic ETA Sleep Protocol Synchronization

> **Version:** 1.0.0
> **Date:** 2026-09-12
> **Context:** Smart incremental test execution, dual-queue concurrency pools, temp/failure isolation, and cross-repo synchronization with coding-guidelines.

## Key Architectural Decisions

1. **Dual-Queue Worker Pools**:
   - **Slow Tests Queue**: 4 worker threads, max 2 tests running per batch. Mitigates disk I/O and process contention for heavy integration and end-to-end tests.
   - **Fast Tests Queue**: 4 worker threads, max 4 tests running per batch, dispatched in chunks of 100 tests from the test inventory queue.
2. **Directory & Failure Isolation**:
   - Temporary artifacts, caches, and ETA telemetry reside strictly in `.lovable/temp/`. Root-level `temp/` is prohibited.
   - Failures are captured individually in `.lovable/temp/failures/<test-id>.log`.
   - Passing tests are 100% silent in terminal output and the filesystem (no logs created).
3. **Relative Path Manifests**:
   - `.lovable/test-inventory.json` stores all paths (`target_file` and `test_file`) as repo-relative forward-slash paths (e.g. `cli/cmd/root.go`, `cli/cmd/root_test.go`).
   - Incremental test execution triggers solely on dirty SHA-256 hash changes or explicit CLI filters.
4. **Dynamic ETA Sleep Protocol**:
   - Runner writes progress and remaining duration estimates to `.lovable/temp/runner-eta.json`.
   - AI agents observe runner ETA, sleep for the estimated duration or 60s, rather than polling actively, conserving tokens and execution bandwidth.
5. **Cross-Repository Synchronization**:
   - Parity maintained across `gitmap` and `coding-guidelines` (`03-ai-scripts/06-cicd-local-runner.py`, `03-ai-scripts/33-test-inventory-generator.py`, `.agents/skills/smart-test-runner-and-inventory/skill.md`, prompts).
