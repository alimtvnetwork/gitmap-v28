# Subtask 01: Heavy Test Isolation & Duration Profiling

## Objective
Isolate tests executing subprocess commands, file system hierarchies, or timing delays out of unit packages into `cli/tests/heavy_test/` (`package heavy_test`).

## Execution Details
1. Audit `cli/tests/release_test`, `cli/tests/fixrepo_test`, and `cli/cmd/installctx_*_e2e_test.go`.
2. Migrate heavy subprocess tests into `cli/tests/heavy_test/`.
3. Ensure parent packages contain pure in-memory unit tests with <0.01s execution latency.
4. Verify with `go vet -C cli ./tests/heavy_test`.
