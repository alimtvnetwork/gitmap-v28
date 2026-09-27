# Subtask 06: Temporary E2E Tests, Minor Release & CI/CD Verification
Traceability ID: Task-07
Spec Reference: [02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md](../../../02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md)
Target Files: cli/tests/e2e/test_gitmap_migration_tempe2e_test.go, cli/constants/constants.go
Action:
- Update `cli/tests/e2e/test_gitmap_migration_tempe2e_test.go` to test Array Async Pool probing, bootstrap config creation, short input parsing, and self-contained SEO template variable resolution.
- Verify with on-demand `RUN_TEMP_E2E=1 go test -tags=tempe2e -v ...`.
- Execute minor release bump `v6.348.0` via `python 03-ai-scripts/29-release-orchestrator.py --tier minor`.
- Build and install the global binary at `C:\Users\Administrator\AppData\Local\gitmap-cli\gitmap.exe`.
- Monitor remote GitHub Actions CI/CD via `gitmap pl-ai status` until all workflows pass 100% green.
Acceptance Criteria:
- Temporary E2E tests pass.
- Release `v6.348.0` published.
- Remote CI/CD is green.
Targeted Verification: gitmap pl-ai status -t 30
