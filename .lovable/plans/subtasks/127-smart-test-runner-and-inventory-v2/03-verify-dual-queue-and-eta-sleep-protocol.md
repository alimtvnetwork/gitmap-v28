# 03-verify-dual-queue-and-eta-sleep-protocol.md: Subtask 3 - Verify Dual Worker Queue & ETA Sleep Protocol

**Status: completed**

## Objectives
1. Verify dual worker queue in `run_smart_go_tests`:
   - Queue 1 (Slow tests): 4 workers, 2 tests per batch.
   - Queue 2 (Fast tests): 4 workers, 4 tests per batch, processing chunks of 100 tests from inventory.
2. Dynamic ETA calculation based on duration estimation from test inventory.
3. Verify live telemetry updates to `.lovable/temp/runner-eta.json`: status, total_eta_sec, elapsed_sec, remaining_eta_sec, completed, passed, failed.
4. Verify in-flight runner heartbeats are clean (>=25s).
5. Ensure AI sleep protocol: when checking background runner tasks, if active, AI agents must read remaining ETA and sleep for that duration instead of busy-polling.
