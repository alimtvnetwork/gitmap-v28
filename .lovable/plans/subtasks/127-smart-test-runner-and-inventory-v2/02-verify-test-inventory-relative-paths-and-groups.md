# 02-verify-test-inventory-relative-paths-and-groups.md: Subtask 2 - Verify Test Inventory Relative Paths & Groups

**Status: completed**

## Objectives
1. Connect `06-cicd-local-runner.py` to `33-test-inventory-generator.py` so that all 3,534 repository tests across all packages are indexed without omitting non-cli packages.
2. Preserve all summary metrics in `.lovable/test-inventory.json`: `slow_tests`, `fast_tests`, `heavy_tests`, `unit_tests`, `slow_threshold_sec`, `estimated_slow_sec`, `estimated_fast_sec`.
3. Verify that 100% of tests have non-empty, repository-relative paths for `target_file` and `test_file` (0 absolute paths, 0 empty paths).
4. Verify configurable slow threshold (>4.0s default, configurable via `--slow-threshold` and `GITMAP_SLOW_TEST_THRESHOLD`).
5. Ensure baseline profiling with `--force-run-all` marks all tests dirty for complete execution.
