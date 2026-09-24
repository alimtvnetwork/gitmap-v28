# Subtask [07]: Prompt Suite 21-temp-e2e-tests and Local AUM Search Benchmark
Traceability ID: Task-08, Task-09
Spec Reference: [02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md](../../../02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md)
Target Files: 01-prompts/21-temp-e2e-tests/01-temp-e2e-test.md, 01-prompts/21-temp-e2e-tests/02-temp-benchmark.md, benchmark.md
Action: Create folder `01-prompts/21-temp-e2e-tests/` with `01-temp-e2e-test.md` (configurable `N=300` on top, multi-stage structure) and `02-temp-benchmark.md`. Run local benchmark of AUM search vs standard scanner, outputting to local uncommitted `benchmark.md`.
Acceptance Criteria:
1. `01-prompts/21-temp-e2e-tests/01-temp-e2e-test.md` contains `N = 300` at top and follows multi-step pipeline.
2. `01-prompts/21-temp-e2e-tests/02-temp-benchmark.md` defines the local search benchmarking workflow.
3. `benchmark.md` saved in local workspace and ignored from CI/CD.
Targeted Verification: Test file existence and content validation.
