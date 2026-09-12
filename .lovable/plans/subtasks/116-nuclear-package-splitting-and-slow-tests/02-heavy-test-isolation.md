# Subtask 02: Heavy Test Isolation into Dedicated Package

## Objective
Extract the 19 heavy test files from `gitmap/cmd` into a dedicated package `gitmap/tests/heavy_test` (`package heavy_test`) to eliminate subprocess and compile latency from routine package tests.

## Steps
1. Create directory `gitmap/tests/heavy_test/`.
2. Move heavy test files from `gitmap/cmd` that invoke `exec.Command` or `time.Sleep` into `gitmap/tests/heavy_test/`.
3. Update package declarations to `package heavy_test` and import `github.com/alimtvnetwork/gitmap-v28/gitmap/cmd` as an external test package.
4. Export any required helper functions in `cmd` needed by the tests.
5. Verify `go vet ./tests/heavy_test` and `go vet ./cmd` pass with exit code 0.
