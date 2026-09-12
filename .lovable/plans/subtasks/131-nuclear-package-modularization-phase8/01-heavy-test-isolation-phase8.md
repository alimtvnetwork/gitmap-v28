# Subtask 01: Heavy Test Isolation (Phase 8)

## Objective
Segregate tests taking >= 1.0s or executing external subprocesses into `cli/tests/heavy_test/` (`package heavy_test`).

## Execution Details
1. Audit and move tests: `cli/cluster.TestExecGit_Commands`, `cli/cmdinstaller.TestInstallerExportGitCmd`, `cli/macro.TestExecute_WithGitmapCdAndRelativeCd`, `cli/release.TestTrimGitOutputFallback`, `cli/probe.TestBackgroundRunner_HonorsWorkerCap`.
2. Ensure subpackage unit tests remain 100% fast, pure in-memory tests (<0.01s).
3. Verify compilation via `go vet -C cli ./tests/heavy_test`.
