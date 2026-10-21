# Subtask 04: Docker Environment & Linux Containerized Verification

- **Task ID:** `Subtask-04`
- **Subtask Slug:** `04-subtask-docker-and-linux-verification`
- **Assigned Role:** Worker 02
- **Status:** PENDING
- **Parent Plan:** `.ai-memory/plans/gitmap-testing-and-schedule-management.md`
- **Target Specifications:**
  - `02-spec/21-app/gitmap-testing-and-schedule-management/03-docker-and-linux-container-strategy.md`
  - `02-spec/21-app/gitmap-testing-and-schedule-management/04-safety-and-os-test-isolation.md`

---

## 1. Scope & Exclusive File Ownership

### 1.1 Owned Files (Strictly Disjoint — Exclusive to Worker 02)
- `Dockerfile.test`
- `.github/workflows/docker-test.yml`
- `03-ai-scripts/51-verify-linux-cross-compile.py`

### 1.2 Read-Only Context References
- `02-spec/21-app/gitmap-testing-and-schedule-management/03-docker-and-linux-container-strategy.md`
- `go.mod`
- `cli/cmdschedule/schedule_os.go`
- `cli/cmdschedule/schedule_cmd.go`

---

## 2. Core Objectives

1. **Host Discovery Alignment:** Formally document that the local Windows host lacks Docker daemon and WSL distribution, establishing pure Go cross-compilation as the primary local gate.
2. **Hermetic `Dockerfile.test` Harness:** Author the reproducible test container targeting `golang:1.24-alpine` equipped with mock OS power binaries (`/usr/local/bin/shutdown`, `poweroff`, `reboot`) and non-root execution (`1001:gitmap`).
3. **Pure Go Linux Cross-Compilation:** Verify that all packages compile cleanly under `GOOS=linux GOARCH=amd64 CGO_ENABLED=0` without symbol drift or platform-dependent build tag errors.
4. **Automated Verification Runner:** Author a cross-platform Python runner (`03-ai-scripts/51-verify-linux-cross-compile.py`) to automate cross-compilation checks and container smoke tests.
5. **CI/CD Integration:** Wire container and cross-compilation gates into `.github/workflows/docker-test.yml`.

---

## 3. Step-by-Step Implementation Work Breakdown

### Phase 4.1: Authoring `Dockerfile.test`
- **Action 4.1.1:** Create `Dockerfile.test` adhering to specification in `02-spec/21-app/gitmap-testing-and-schedule-management/03-docker-and-linux-container-strategy.md`:
  - Base: `golang:1.24-alpine`.
  - Install `git`, `bash`, `ca-certificates`, `tzdata`.
  - Stub mock power binaries in `/usr/local/bin/` with executable permissions.
  - Create unprivileged user `gitmap` (`uid=1001`).
  - Pre-copy `go.mod` and `go.sum` for layer caching.
  - Set `CGO_ENABLED=0` and non-root user.
  - Set entrypoint: `["go", "test", "-v", "./cli/cmdschedule/..."]`.

### Phase 4.2: Pure Go Linux & Darwin Cross-Compilation Suite
- **Action 4.2.1:** Execute package test symbol check:
  ```bash
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -run=^$ ./cli/cmdschedule/...
  ```
- **Action 4.2.2:** Execute full repository Linux binary compilation check:
  ```bash
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /dev/null ./...
  ```
- **Action 4.2.3:** Execute full repository Darwin binary compilation check:
  ```bash
  GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o /dev/null ./...
  ```
- **Action 4.2.4:** Audit all file-handling and path-building code in `cli/cmdschedule` to ensure `filepath` and forward-slash normalization are strictly followed.

### Phase 4.3: Automation Verification Runner Script
- **Action 4.3.1:** Author `03-ai-scripts/51-verify-linux-cross-compile.py`:
  - Enforce Python 3 standard library only (`subprocess`, `sys`, `platform`, `os`).
  - Check availability of `go` binary and verify version $\ge 1.24$.
  - Probe `docker` CLI: if present, attempt container build and test run; if absent, print informational skip notice and execute host cross-compilation gates.
  - Run `GOOS=linux go test -run=^$` and `GOOS=linux go build`.
  - Output structured, color-coded terminal report with summary exit code 0 or 1.

### Phase 4.4: CI/CD Workflow Authoring (`.github/workflows/docker-test.yml`)
- **Action 4.4.1:** Author workflow `.github/workflows/docker-test.yml`:
  - Trigger on push to `main` and pull requests touching `cli/cmdschedule/**`, `Dockerfile.test`, or `.github/workflows/docker-test.yml`.
  - Job 1: `container-test` running on `ubuntu-latest` building and running `Dockerfile.test`.
  - Job 2: `cross-compile-gate` running on `ubuntu-latest` verifying Linux and Darwin cross-compilations with `CGO_ENABLED=0`.

### Phase 4.5: Quality Gate & Boundary Review
- **Action 4.5.1:** Verify that `Dockerfile.test` and workflow files contain zero hardcoded absolute paths.
- **Action 4.5.2:** Verify that all script output uses relative Git paths.
- **Action 4.5.3:** Validate zero negative booleans across authored scripts and workflows.

---

## 4. Verification Evidence & Quality Gates

```bash
# Gate 1: Pure Go Linux cross-compilation test check
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -run=^$ ./cli/cmdschedule/...

# Gate 2: Full repository Linux cross-compilation build check
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /dev/null ./...

# Gate 3: Python verification runner execution
python 03-ai-scripts/51-verify-linux-cross-compile.py
```

### Acceptance Checklist
- [ ] `Dockerfile.test` authored with `golang:1.24-alpine`, mock power binaries, and non-root user.
- [ ] Pure Go Linux cross-compilation passes cleanly with zero symbol drift.
- [ ] Darwin cross-compilation passes cleanly.
- [ ] Automation runner script executes on Windows host and validates cross-compilation cleanly.
- [ ] Zero host OS reboots or destructive actions during any validation step.
- [ ] All files strictly adhere to <= 300 lines limit.
