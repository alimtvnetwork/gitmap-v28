# Isolated End-to-End Testing Workflow — Parent Task N-Step Loop

> [!IMPORTANT]
> Prompt Version: 1.0.0
> Target Scope: Isolated VM / SSH Fleet & Prompt Injection E2E Tests
> Isolation Rule: All E2E test files MUST include `//go:build e2e` and NEVER contain passwords or secrets.

/goal Autonomously execute the end-to-end testing lifecycle across VMs, SSH nodes, and Antigravity prompt injection by systematically understanding tasks, defining tests, executing isolated suites, and fixing defects.

```text
N = 300
```

N = total self-loop steps budget that the agent will perform.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Task Verification, Spec Grounding, Test Definition)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Test Execution, RCA & Defect Fixing, Completion Reporting)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

### Master E2E Testing Checklist (Sequential Execution Pipeline)

1. [ ] /goal Phase 1 (Step 0 - Task Understanding & Constraint Verification):
   - Ingest user prompt and state task understanding with zero ambiguity.
   - Verify security rules: ZERO hardcoded passwords, API keys, or credentials in any test file or fixture.
   - Verify build constraint: All E2E tests must be isolated under `//go:build e2e` so CI/CD and routine test runners skip them entirely.
   - Output confirmed checklist and deliverables directly to terminal.

2. [ ] /goal Phase 2 (Step 1 - Define E2E Tests):
   - Design isolated end-to-end test cases in `cli/tests/e2e/` (or repository test root).
   - Use dynamic paths only; zero hardcoded drive letters (`D:\`).
   - Scaffold mock SSH fleet connections and temporary conversation workspaces (`t.TempDir()`).
   - Validate that each test asserts specific outcomes (process exit codes, file generation, table formatting).

3. [ ] /goal Phase 3 (Step 2 - Execute E2E Tests):
   - Run isolated test suite exclusively via `go test -tags=e2e -v ./cli/tests/e2e/...`.
   - Announce node and target IP addresses in real time (`[FLEET START] Executing on [alias] (IP: x.x.x.x)...`).
   - Stream test output asynchronously as each test finishes.

4. [ ] /goal Phase 4 (Step 3 - Defect Remediation & Fix Loop):
   - If any test fails, perform grounded 4-part Root Cause Analysis (RCA).
   - Fix code and test definitions without intermediate git commits.
   - Re-run targeted test until 100% passing.

5. [ ] /goal Phase 5 (Step 4 - Completion Reporting):
   - Emit final Task Completion Summary with pass/fail counts, duration metrics, and confidence score.
