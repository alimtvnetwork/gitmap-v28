---
name: temp-end-to-end-tests
description: Autonomously design, implement, and execute temporary end-to-end integration tests combining complete subsystem flows locally, strictly isolating them with skip-by-default tags (//go:build tempe2e) and environment guards (RUN_TEMP_E2E=1) so they never execute in CI/CD pipelines or standard test suites.
---

# Temporary End-to-End Tests & Isolated On-Demand Validation

## Core Mandate
1. **Skip-by-Default Everywhere**: All temporary E2E tests in Go MUST include `//go:build tempe2e`, be named `*_tempe2e_test.go`, and check `if os.Getenv("RUN_TEMP_E2E") != "1" { t.Skip(...) }`.
2. **Zero CI/CD Impact**: Never enable `tempe2e` in `.github/workflows/*.yml` or routine local runners.
3. **On-Demand Execution**: Execute only via `RUN_TEMP_E2E=1 go test -tags=tempe2e -v ./cli/tests/e2e/... -run <TestName>`.
4. **Clean Teardown**: Always clean up temporary SQLite databases, files, and sockets via `t.Cleanup()`.
