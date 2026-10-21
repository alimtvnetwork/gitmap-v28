# Subtask 06: Master Index Registration, Quality Gates & Atomic Push

- **Subtask ID:** `06-subtask-index-quality-gate-and-push`
- **Parent Task:** `gitmap-testing-and-schedule-management`
- **Specification Reference:** `02-spec/21-app/gitmap-testing-and-schedule-management/06-verification-and-acceptance-criteria.md`
- **Catalog Reference:** `02-spec/21-app/gitmap-testing-and-schedule-management/readme.md`
- **Assigned Worker:** Lead
- **Status:** DONE

---

## 1. Owned Files (Exclusive)

The following files are strictly owned by Lead for final consolidation and registration:
- `02-spec/21-app/readme.md`
- `.ai-memory/plans/readme.md`
- `ledger.md`

---

## 2. Read-Only Context Files

- `02-spec/21-app/gitmap-testing-and-schedule-management/readme.md`
- `02-spec/21-app/gitmap-testing-and-schedule-management/01-architecture-and-cross-os-overview.md`
- `02-spec/21-app/gitmap-testing-and-schedule-management/02-cross-os-power-command-contracts.md`
- `02-spec/21-app/gitmap-testing-and-schedule-management/03-docker-and-linux-container-strategy.md`
- `02-spec/21-app/gitmap-testing-and-schedule-management/04-safety-and-os-test-isolation.md`
- `02-spec/21-app/gitmap-testing-and-schedule-management/05-mail-and-schedule-e2e-matrix.md`
- `02-spec/21-app/gitmap-testing-and-schedule-management/06-verification-and-acceptance-criteria.md`
- `.ai-memory/plans/gitmap-testing-and-schedule-management.md`
- `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/01-subtask-git-sync-and-preflight.md`
- `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/02-subtask-architecture-and-specs.md`
- `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/03-subtask-cross-os-schedule-tests.md`
- `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/04-subtask-docker-linux-container-verification.md`
- `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/05-subtask-new-features-e2e-verification.md`

---

## 3. Subtask Objectives & Step-by-Step Execution Plan

### Step 1: Comprehensive Quality Gate Audit (G1 through G9)
1. Verify cross-OS power scheduling command construction (`G1`):
   - Assert `schedule_os_test.go` exercises Windows, Linux, and Darwin power action mappings.
2. Verify safety test isolation (`G2`):
   - Confirm zero real host system calls; all power commands use mock `OSActionExecutor`.
3. Verify pure Go Linux cross-compilation (`G3`):
   - Run `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -run=^$ ./...`.
   - Run `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./...`.
4. Verify Docker container test harness (`G4`):
   - Assert `Dockerfile.test` exists, targets `golang:1.24-alpine`, and includes hermetic build steps.
5. Verify folder watcher, Split-DB, privacy masking, and zero-CDN UI gates (`G5`-`G8`):
   - Assert `go test -v -race ./cli/cmdmail/... ./cli/cmdschedule/...` passes with zero failures.
6. Verify structural hygiene budget (`G9`):
   - Scan all authored markdown and Go files: assert every file is strictly $\le 300$ lines.
   - Scan for absolute path leaks or URI schemes: assert 100% relative Git paths.

### Step 2: Specification Catalog & Plan Index Updates
1. Update `02-spec/21-app/readme.md`:
   - Append entry for `gitmap-testing-and-schedule-management` linking `02-spec/21-app/gitmap-testing-and-schedule-management/readme.md`.
   - Document scope: cross-OS power scheduling, Dockerized Linux testing harness, mail and schedule E2E verification matrix.
2. Update `.ai-memory/plans/readme.md`:
   - Register `gitmap-testing-and-schedule-management.md` under the active/completed plan index.
3. Update `ledger.md`:
   - Record transition of `gitmap-testing-and-schedule-management` from `IN_PROGRESS` to `COMPLETED`.
   - Log completed subtasks 01 through 06 with timestamps.

### Step 3: Working Tree Sanitization & Pre-Push Preflight
1. Clean temporary files, build caches, and test artifacts.
2. Execute full repository test suite:
   ```bash
   go test -short ./...
   ```
3. Confirm zero uncommitted or untracked temporary files outside target scope.

### Step 4: Atomic Commit & Upstream Delivery Protocol
1. Stage all modified and authored files:
   - Specification suite: `02-spec/21-app/gitmap-testing-and-schedule-management/`
   - Plan and subtasks: `.ai-memory/plans/` and `.ai-memory/plans/subtasks/`
   - Test suites and implementations: `cli/cmdmail/`, `cli/cmdschedule/`, `Dockerfile.test`
   - Catalogs: `02-spec/21-app/readme.md`, `.ai-memory/plans/readme.md`, `ledger.md`
2. Construct atomic Git commit message following GitMap standards:
   ```text
   feat(schedule): cross-os power testing, docker harness and e2e verification
   ```
3. Push to upstream tracking branch and verify clean remote synchronization:
   - "No Push = Not Done" completion invariant enforced.

---

## 4. Final Quality Gates & Completion Sign-Off

- [ ] All nine quality gates (`G1` through `G9`) verified and signed off.
- [ ] Master specification catalog (`02-spec/21-app/readme.md`) updated.
- [ ] Plan index (`.ai-memory/plans/readme.md`) and ledger (`ledger.md`) updated.
- [ ] Full test suite passes without regressions.
- [ ] All authored files strictly $\le 300$ lines with 100% relative Git paths.
- [ ] Atomic Git commit pushed to remote upstream repository.
