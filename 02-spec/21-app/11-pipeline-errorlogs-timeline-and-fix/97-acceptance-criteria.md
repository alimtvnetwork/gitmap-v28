# Acceptance Criteria: Pipeline Errorlogs Timeline & CI/CD Fix Suite

## Scenarios

### Scenario 1: Watching an In-Progress Pipeline with `-t`

- **Given** a repository with an active GitHub Actions workflow run in progress
- **When** the user executes `gitmap pipeline errorlogs -t`
- **Then** GitMap displays adaptive countdown progress lines until the workflow concludes
- **And** upon completion, displays error logs if failed, or success confirmation if passed.

### Scenario 2: Displaying Historical Rerun ETA on Failure

- **Given** a failed pipeline execution
- **When** the user invokes `gitmap pipeline errorlogs` or `gitmap pipeline errorlogs -t`
- **Then** GitMap calculates the average duration of past successful runs for that workflow
- **And** displays `Estimated pipeline rerun duration (ETA): ~Xs` in terminal output
- **And** populates `rerunEtaSeconds` in `--json` output.

### Scenario 3: Automated Diagnostic & Auto-Fix Execution

- **Given** source files in the repository
- **When** the user invokes `gitmap pipeline errorlogs --fix`
- **Then** GitMap runs all internal diagnostic probes
- **And** automatically applies `gofmt -w .` formatting to unformatted files
- **And** displays status lines for each probe with `PASS`, `FIXED`, or `FAIL`.

### Scenario 4: Non-Interactive Flag Support

- **Given** automated scripts running GitMap
- **When** `--fix -y` or `--check --json` is passed
- **Then** GitMap executes without blocking on terminal stdin prompts.
