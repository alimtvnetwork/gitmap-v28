# Subtask 04: Verification, Quality Gates & Audit Protocol — Agent AI Synchronization

> **Subtask ID:** Task-04  
> **Parent Plan:** `.ai-memory/plans/agent-ai-synchronization-and-task-management.md`  
> **Spec Reference:** `02-spec/21-app/agent-ai-synchronization-and-task-management/02-component-and-cli-spec.md`  
> **Assigned Role:** Lead / Verification Worker  
> **Status:** PENDING  
> **Owned Files:**  
> - `.ai-memory/plans/subtasks/agent-ai-synchronization-and-task-management/04-verification-and-audit.md`  

---

## 1. Deliverables & Verification Objectives

This subtask defines the comprehensive quality gate, verification test suite, and retrospective audit protocols for the **Agent AI Synchronization and Task Management** implementation. All code changes across Go binaries, Python task manager scripts, specification documents, and skills must be rigorously verified before final consolidation and release tagging.

### Key Verification Focus Areas:
1. **JSON Contract Fidelity:** Confirm all CLI commands (`create-task`, `subtask create`, `claim-files`, `editing`, `heatmap discuss`) emit valid, structured JSON containing the required attributes when passed `--json`.
2. **Machine Context Lineage:** Verify deterministic detection and propagation of `machineId` across parent tasks, subtasks, and JSON payloads.
3. **ActionType Enum Enforcement:** Verify that only canonical enum values (`modify`, `analyze`, `read`, `remove`, `refactor`, `create`) are accepted by `--action`, and illegal values are rejected with exit code 1.
4. **Performance Thresholds:** Confirm that pre-touch logging executes in $<2\text{ms}$ and heatmap discussion summaries aggregate in $<5\text{ms}$.
5. **Quality Linters:** Execute FastGate (`50-fastgate.py`), version sync (`14-version-sync-checker.py`), and privacy/relative path guards (`49-verify-privacy-and-relative-paths.py`).
6. **Zero Absolute Paths:** Guarantee 100% adherence to strict relative Git paths across all written code, specs, and plans.

---

## 2. Concrete Verification Commands & Test Matrix

| Step | Target Subsystem | Verification Command | Expected Outcome | Failure Remediation |
|---|---|---|---|---|
| **V-01** | **FastGate AST Linter** | `python 03-ai-scripts/50-fastgate.py` | Exit code 0, 0 AST errors, 0 path violations | Revert or fix reported file paths or syntax |
| **V-02** | **Version Synchronization** | `python 03-ai-scripts/14-version-sync-checker.py` | Exit code 0, versions in `version.json`, `package.json`, and `changelog.md` match | Align SemVer versions |
| **V-03** | **Relative Path Hygiene** | `python 03-ai-scripts/49-verify-privacy-and-relative-paths.py` | Exit code 0, zero drive letters (`C:\`, `D:\`), zero `file://` URIs | Replace with forward-slash relative Git paths |
| **V-04** | **Task DB Schema Parity** | `python 03-ai-scripts/46-agent-sqlite-task-manager.py schema --json` | Valid JSON schema with `MachineId` in `ParentTask` & `Subtask`, and `ActionType` CHECK constraint | Update `SCHEMA_DDL` and metadata in script |
| **V-05** | **Parent Task Creation (CLI)** | `gitmap agent-ai create-task "Verification Suite" --budget 300 --json` | Emits JSON with `isSuccess=true`, `taskSlug="verification-suite"`, non-empty `machineId`, status `ACTIVE` | Fix task registration logic in CLI router |
| **V-06** | **Subtask Lineage & Machine ID** | `gitmap agent-ai subtask create verification-suite "Verify CLI" --code Task-01 --files "cli/cmdagent/agent_cmd.go" --role "Worker 01" --json` | Emits JSON with `parentTaskId`, `machineId`, `taskCode="Task-01"`, `ownedFiles=["cli/cmdagent/agent_cmd.go"]`, status `PENDING` | Check parent task lookup and foreign key propagation |
| **V-07** | **Look-Ahead File Claim** | `gitmap agent subtask claim-files -t verification-suite -s Task-01 -f "cli/cmdagent/agent_cmd.go" -a "Worker 01" --json` | Emits JSON with `hasCollision=false`, `status="CLAIMED"` | Check `FileClaim` table insertion and querying |
| **V-08** | **Collision Alerting (Negative)** | `gitmap agent subtask claim-files -t verification-suite -s Task-02 -f "cli/cmdagent/agent_cmd.go" -a "Worker 02" --json` | Emits JSON with `hasCollision=true`, conflicting agent `"Worker 01"`, and prints collision banner | Verify cross-subtask claim matching in Tier 1 registry |
| **V-09** | **ActionType Enum (Negative)** | `gitmap agent-ai editing Task-01 "cli/cmdagent/agent_cmd.go" --action "tamper" --lines 1,10 --reasoning "Negative test"` | Rejects with error message detailing allowed enum values, exit code 1 | Validate action against canonical enum set |
| **V-10** | **ActionType Enum (Positive)** | `gitmap agent-ai editing Task-01 "cli/cmdagent/agent_cmd.go" --action modify --lines 10,25 --reasoning "Verification update" --json` | Emits JSON with `isSuccess=true`, `actionType="modify"`, latency $<2\text{ms}$ | Optimize SQLite WAL write and indexing |
| **V-11** | **Heatmap Discuss Speed** | `gitmap agent-ai heatmap discuss verification-suite --fast --json` | Emits JSON summary with heated files, contention scores, execution time $<5\text{ms}$ | Enable SQLite index scans and memory caching |
| **V-12** | **Go Package Tests** | `go test -v ./cli/cmdagent/...` | 100% passing tests, 0 failures, 0 race conditions | Fix failing unit test assertions |

---

## 3. Acceptance Criteria Traceability Matrix

| Spec Criterion | Description | Target Test Step | Status |
|---|---|---|---|
| **AC-01** | Parent Task Creation with `MachineId` | V-05 | Pending Execution |
| **AC-02** | Subtask Creation with `MachineId` & `ParentTaskId` | V-06 | Pending Execution |
| **AC-03** | Look-Ahead File Box Claiming | V-07 | Pending Execution |
| **AC-04** | Collision Detection & Warning Banner | V-08 | Pending Execution |
| **AC-05** | Canonical `ActionType` Enum Enforcement | V-09 | Pending Execution |
| **AC-06** | Pre-Touch Action Logging Speed ($<2\text{ms}$) | V-10 | Pending Execution |
| **AC-07** | Heatmap Discussion Speed ($<5\text{ms}$) & Scoring | V-11 | Pending Execution |
| **AC-08** | Structured JSON Schema Output Parity | V-04, V-05, V-06, V-07, V-10, V-11 | Pending Execution |
| **AC-09** | Human ANSI Terminal Formatting | Manual CLI invocation without `--json` | Pending Execution |
| **AC-10** | Strict Relative Git Paths Hygiene | V-01, V-03 | Pending Execution |
| **AC-11** | Positive Boolean Conventions | V-04 (Schema inspection: `isSuccess`, `isActive`, `hasCompleted`) | Pending Execution |
| **AC-12** | Python Engine & Go CLI Parity | Run test steps V-05 through V-11 against both engines | Pending Execution |

---

## 4. Evidence Checklist & Gatekeeper Protocol

Before the orchestrator marks this task complete, all items in this checklist must be satisfied with recorded execution output:

### 4.1 Linters & Code Hygiene
- [ ] `python 03-ai-scripts/50-fastgate.py` executed: Exit code 0, 0 AST violations.
- [ ] `python 03-ai-scripts/14-version-sync-checker.py` executed: Exit code 0, versions synchronized.
- [ ] `python 03-ai-scripts/49-verify-privacy-and-relative-paths.py` executed: Exit code 0, zero IP or path leaks.
- [ ] No trailing whitespaces or non-LF line endings in newly authored files.

### 4.2 CLI & Functional Evidence
- [ ] `gitmap agent-ai create-task` JSON response recorded proving `machineId`, `runDirectory`, `dbPath`.
- [ ] `gitmap agent-ai subtask create` JSON response recorded proving `parentTaskId` and `machineId` propagation.
- [ ] `gitmap agent subtask claim-files` collision report recorded showing high-contrast alert.
- [ ] Negative test on invalid `--action` recorded showing rejection of unapproved action strings.
- [ ] Positive test on `--action modify` recorded showing $<2\text{ms}$ write confirmation.
- [ ] `gitmap agent-ai heatmap discuss` performance verified under $5\text{ms}$.

### 4.3 Documentation & Memory Consolidation
- [ ] Specification `02-spec/21-app/agent-ai-synchronization-and-task-management/02-component-and-cli-spec.md` complete and indexed.
- [ ] Subtask plans in `.ai-memory/plans/subtasks/agent-ai-synchronization-and-task-management/` updated with completion marks and evidence logs.
- [ ] Parent plan `.ai-memory/plans/agent-ai-synchronization-and-task-management.md` updated with final status.
- [ ] Consolidated plan archived in `.ai-memory/plans/completed/agent-ai-synchronization-and-task-management.md`.

---

## 5. Non-Negotiable Boundaries Reminder
- **TOTAL BAN ON GIT COMMANDS FOR SUBAGENTS:** Subagents must NEVER execute `git` commands (`git add`, `git commit`, `git push`, `git status`). The lead orchestrator alone handles staging, committing, and pushing.
- **Strict Relative Git Paths:** Zero absolute filesystem paths and zero `file://` URIs in any document or command payload.
