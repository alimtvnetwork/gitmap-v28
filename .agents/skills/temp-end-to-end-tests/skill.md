---
name: temp-end-to-end-tests
description: >-
  Autonomously design, implement, and verify temporary end-to-end integration tests combining complete system flows locally, enforcing strict skip-by-default isolation so tests never run in CI/CD pipelines or standard local test suites.
---

# Instruction (must follow): Temporary End-to-End Tests & Isolated On-Demand Validation

> [!IMPORTANT]
> Prompt Version: 2.5.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Whatever directives, constraints, checklists, or instructions are given before this section or prompt (including in the prompt preamble, header blocks, or incoming user request) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

/goal Autonomously design, implement, and execute temporary end-to-end integration tests combining complete subsystem flows locally, strictly isolating them with skip-by-default tags and environment guards so they NEVER execute in automated CI/CD pipelines or standard local test suites.

---

## Sequential Single-Agent Execution Mandate (No Step Budgets, No Parallelism)

Unlike multi-agent parallel workflows, temporary end-to-end test development enforces strict linear determinism:

1. **Zero Step Budgeting:** Do not track or calculate artificial step budgets (no `N = 200`, no `PHASE_1_STEPS = N / 2`, no step decrement equations). Work proceeds linearly through distinct phases until all verification gates pass cleanly.
2. **Zero Parallel Dispatch (Strictly Sequential):** Never spawn parallel execution subagents, multi-threaded worker pools, or concurrent dispatch routines. End-to-end tests touch shared resources, local databases, ports, and subprocess state; concurrency introduces race conditions, port collisions, and fixture corruption. All actions are executed sequentially by a single dedicated agent.

---

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Preamble Precedence Verification: Whatever is given before this section or prompt (user preamble, header constraints, prior instructions) has been verified as highest priority and non-negotiable, and is strictly incorporated into the task scope ahead of all other guidelines.
2. [ ] /goal Phase 1A (Step 0 - Verbatim Prompt Recording & Task Extraction Gate): Immediately capture the user prompt verbatim into `.ai-memory/plans/pending/xx-<slug>.md` under `## User Request (Verbatim)`. Extract actionable deliverables with traceable IDs (`Task-01`, `Task-02`), and output this confirmed task breakdown directly in chat in cleanly indented markdown with vertical blank lines, task state (`State: [IN PROGRESS — EXECUTING IMMEDIATELY]`), and understanding indicator bracket (`Understood: [YES — ...]`) before any file exploration, scanning, or spec writing. Same-turn tool chaining is mandatory (never terminate turn with text alone).
3. [ ] /goal Phase 1B (Step 1 - Canonical Spec Generation in Folder 21): Write the canonical specification in `02-spec/21-app/xx-<slug>.md` (or directory `02-spec/21-app/xx-<slug>/` for complex features) with lossless verbatim prompt capture, register it in `02-spec/21-app/readme.md`, and initialize the execution plan in `.ai-memory/plans/pending/xx-<slug>.md` linking back to the spec.
4. [ ] /goal Phase 1B (Step 2 - Scan & Discover): Use GitMap AUM discovery (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary, with fast Python discovery scripts (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) as fallback, to inventory target integration files and fixtures without tool truncation limits.
5. [ ] /goal Phase 1B (Step 3 - Lean Subtask Decomposition): Decompose the plan into lean subtask files in `.ai-memory/plans/subtasks/xx-<slug>/01-<subslug>.md` cross-referencing the canonical spec. Subtasks must focus purely on unique task deliverables without repeating common repository boilerplate.
6. [ ] /goal Phase 1B (Step 4 - Readiness Audit Gate): Confirm canonical spec is registered in `02-spec/21-app/readme.md`, and all `Task-xx` deliverables are mapped to subtasks linking back to the spec before execution.
7. [ ] /goal Phase 1B (Step 5 - Unconditional Zero-Question Execution Mandate): Immediately upon completing Phase 1, transition directly into Phase 2 execution mode without pausing, asking questions, or seeking user confirmation.
8. [ ] /goal Phase 2 (Execution & Sequential Implementation): Unconditionally implement code refactoring and test fixtures sequentially following all repository coding guidelines (boolean standards, concrete types in `types.go`, `*appfault.AppError`, function lengths <= 8-15 lines, Unix LF).
9. [ ] /goal Phase 2 (Strict Skip-by-Default Isolation): Tag all temporary E2E test files with mandatory build tags (`//go:build tempe2e`), pytest markers (`@pytest.mark.temp_e2e`), or Vitest/Jest skip guards (`describe.skipIf(!isTempE2EActive)`), and apply runtime guards (`RUN_TEMP_E2E`).
10. [ ] /goal Phase 2 (Targeted On-Demand Verification ONLY): Run ONLY the specific isolated temporary E2E test using the explicit on-demand command (e.g. `RUN_TEMP_E2E=1 go test -tags=tempe2e -v ...`). Do not run `06-cicd-local-runner.py` or standard unflagged test commands.
11. [ ] /goal Phase 3 (Consolidation & Atomic Push): Consolidate completed subtasks into `.ai-memory/plans/completed/xx-<slug>.md` preserving the canonical spec reference (canonical spec in `02-spec/21-app/` remains permanently intact), delete granular subtasks and pending plan, stage all changes, and push in a single grouped commit.
12. [ ] /goal Phase 3 (Completion & Confidence Reporting): Emit the final Task Completion Summary with strict line-by-line bullet format (- ✅), isolation verification summary, modified files summary, and implementation confidence score.
13. [ ] /learn Ingest `.ai-memory/memory/readme.md` for project memory index and past learnings.
14. [ ] /learn Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
15. [ ] /learn Ingest `02-spec/02-coding-guidelines/` for domain-specific architectural specifications.
16. [ ] /learn Ingest `02-spec/03-error-manage/` for error handling architectures and *appfault.AppError.
17. [ ] /learn Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.

---

## Temporary E2E Test Architecture & Strict Skip-by-Default Isolation Standard

Temporary end-to-end integration tests assemble real database connections, subprocess lifecycles, configuration states, CLI binaries, or network sockets to verify complex multi-component workflows locally.

Because these tests are heavy, stateful, and non-deterministic on virtualized cloud runners:

### 1. The Zero-CI Quarantine Standard (TOTAL BAN on Running in CI/CD)

- **Absolute Prohibition in CI/CD:** Temporary E2E tests MUST NEVER execute in GitHub Actions workflows (`.github/workflows/*.yml`), automated pull request checks, or deployment pipelines.
- **Zero Storage & Speed Guarantee:** CI/CD runners must stay lightweight, ultra-fast (<2 minutes), and strictly zero-storage. Heavy E2E tests cause pipeline timeouts, flaky runner hangs, and exhausted free-tier quotas.
- **CI Configuration Invariant:** Never modify any `.github/workflows/*.yml` or CI scripts to enable, trigger, or pass temporary E2E test flags.

### 2. The Routine Local Run Exclusion Standard (Skip-by-Default Everywhere)

- **Default Test Suite Protection:** Developers and automated local runners running `go test ./...`, `pytest`, `npm test`, or `python 03-ai-scripts/06-cicd-local-runner.py` must NEVER inadvertently run temporary E2E tests.
- **Skip-by-Default Guarantee:** If any test runner executes without explicit, specialized opt-in flags, temporary E2E tests MUST be skipped automatically (0 failures, 0 runtime side effects).

### 3. Polyglot Tagging & Isolation Specification Matrix

Every temporary end-to-end test file MUST implement the mandatory isolation controls for its language:

| Ecosystem | Mandatory Build Tag / File Pattern | Mandatory Runtime Skip Guard | Default Execution State |
|---|---|---|---|
| **Go** | `//go:build tempe2e`<br>File: `*_tempe2e_test.go` | `if os.Getenv("RUN_TEMP_E2E") != "1" { t.Skip(...) }` | **Skipped** (`go test ./...` ignores file unless `-tags=tempe2e` passed) |
| **Python** | `@pytest.mark.temp_e2e`<br>File: `test_*_tempe2e.py` | `if os.getenv("RUN_TEMP_E2E") != "1": pytest.skip(...)` | **Skipped** (`pytest` excludes marker via `-m "not temp_e2e"` by default) |
| **TypeScript / Vitest** | Pattern: `*.tempe2e.test.ts`<br>Dir: `tests/tempe2e/` | `describe.skipIf(process.env.RUN_TEMP_E2E !== '1')(...)` | **Skipped** (`vitest.config.ts` excludes `**/*.tempe2e.test.ts` by default) |
| **Rust** | `#[cfg(feature = "temp_e2e")]`<br>or `#[ignore = "tempe2e"]` | Checked at compile or harness level | **Skipped** (`cargo test` ignores without `--features temp_e2e` or `-- --ignored`) |
| **PHP** | `@group temp_e2e` docblock annotation | Group excluded in `phpunit.xml` default testsuite | **Skipped** (`phpunit` ignores group by default) |

#### Concrete Go Isolation Pattern:

```go
//go:build tempe2e

package e2e_test

import (
	"os"
	"testing"
)

func TestCustomerCheckoutWorkflow_TempE2E(t *testing.T) {
	if os.Getenv("RUN_TEMP_E2E") != "1" {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	// End-to-end orchestration logic combining subsystems
}
```

#### Concrete Python Isolation Pattern:

```python
import os
import pytest

# Module-level skip-by-default guard
if os.getenv("RUN_TEMP_E2E") != "1":
    pytest.skip("skipping temporary e2e test; enable on-demand with RUN_TEMP_E2E=1", allow_module_level=True)

@pytest.mark.temp_e2e
def test_full_pipeline_orchestration_tempe2e():
    """Verify multi-component workflow on-demand."""
    ...
```

#### Concrete TypeScript Isolation Pattern:

```typescript
import { describe, it, expect } from 'vitest';

const isTempE2EActive = process.env.RUN_TEMP_E2E === '1';

describe.skipIf(!isTempE2EActive)('Temporary E2E: Full Integration Flow', () => {
  it('combines database, server, and worker processes', async () => {
    // End-to-end flow execution
  });
});
```

### 4. On-Demand Local Execution Command Protocol

Temporary E2E tests are executed **ONLY and EXCLUSIVELY** when explicitly commanded by the user or on-demand workflow via targeted flags:

- **Go On-Demand:**
  ```bash
  RUN_TEMP_E2E=1 go test -tags=tempe2e -v ./tests/tempe2e/... -run TestCustomerCheckoutWorkflow_TempE2E
  ```
- **Python On-Demand:**
  ```bash
  RUN_TEMP_E2E=1 pytest -v -m temp_e2e tests/tempe2e/test_workflow.py
  ```
- **TypeScript On-Demand:**
  ```bash
  RUN_TEMP_E2E=1 npx vitest run tests/tempe2e/workflow.tempe2e.test.ts
  ```

---

## Phase 2: Sequential Execution & Isolated Local Verification

1. **Sequential Single-Agent Execution:** Unconditionally execute code refactoring and test construction sequentially. Never dispatch multiple agents or concurrent threads.
2. **Strict Coding Guidelines:** All test code, mocks, test harnesses, and fixtures MUST strictly comply with repository coding guidelines (boolean standards with `is` and `has`, concrete types in `types.go`, `*appfault.AppError`, function lengths <= 8-15 lines, Unix LF).
3. **Local Teardown & Hygiene:** Temporary E2E tests must clean up their own SQLite files, temporary sockets, port listeners, and scratch files upon completion (using `t.Cleanup()`, `finally`, or teardown fixtures).
4. **Targeted On-Demand Verification ONLY:** Run ONLY the specific isolated temporary E2E test using the explicit on-demand command (e.g. `RUN_TEMP_E2E=1 go test -tags=tempe2e -v ...`).
5. **Total Ban on Routine Test Suites:** Do NOT run `06-cicd-local-runner.py` or standard untagged test suites (`go test ./...` without tags). Routine test suites must remain completely untouched.

---

## Phase 3: Task Consolidation & File Reduction (End of Loop)

To reduce markdown file count and bloat, consolidate subtasks when all deliverables are verified:

1. Combine all completed granular subtasks from `.ai-memory/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.ai-memory/plans/completed/xx-<slug>.md`.
2. Explicitly reference the canonical spec `[02-spec/21-app/xx-<slug>.md](../../../02-spec/21-app/xx-<slug>.md)` in the consolidated header. The canonical spec remains permanently intact.
3. Delete the original granular `.md` files in `.ai-memory/plans/subtasks/xx-<slug>/`.
4. Delete the original parent plan `.ai-memory/plans/pending/xx-<slug>.md`.
5. Update `.ai-memory/plans/readme.md` to point to the newly consolidated completed file.
6. **Final Step Git Commit & Push (Mandatory):** Stage all modified files, consolidated plans, and memory records (`git add -A`), commit them in a single clean grouped atomic commit (`git commit -m "<type>(<scope>): <summary>"`), and push to git (`git push origin <branch>`). Under no circumstances commit each file individually.

---

## End-of-Turn Verification & Confidence Reporting (Mandatory Output)

At the completion of all tasks and before concluding the turn, emit this structured verification summary in the chat response:

> [!CRITICAL]
> **STRICT LINE-BY-LINE OUTPUT MANDATE (TOTAL BAN ON HORIZONTAL CONCATENATION):**
> Every single completed task in the `Task Completion Summary` MUST be rendered on its OWN SEPARATE LINE starting with an individual markdown list bullet (`- ✅`).
> NEVER concatenate multiple tasks horizontally into a single run-on paragraph or single wrapped line.
> In Markdown, consecutive lines without bullet markers (`- `) collapse into a single run-on horizontal sentence. You MUST format each task as a discrete bullet list item (`- ✅`) followed by an explicit newline!
>
> ❌ **BANNED (Horizontal Run-on Concat):**
> `✅ #1. Task-01: [Title] — Completed ✅ #2. Task-02: [Title] — Completed ✅ #3. Task-03: [Title] — Completed`
>
> ✅ **MANDATORY (Strict Line-by-Line Vertical Markdown List):**
> ```markdown
> ### Task Completion Summary
>
> - ✅ **Task-01: [Descriptive Task Title]** — `[Completed]`
> - ✅ **Task-02: [Descriptive Task Title]** — `[Completed]`
> - ✅ **Task-03: [Descriptive Task Title]** — `[Completed]`
> ```

```markdown
### Task Completion Summary

- ✅ **Task-01: [Descriptive Task Title]** — `[Completed]`
- ✅ **Task-02: [Descriptive Task Title]** — `[Completed]`
(If any task failed or was deferred, mark with `- ❌` or `- ⏳` on its own separate line and explain why)

### Temporary E2E Test Isolation Verification

- **Target Test File:** `[relative/path/to/test]`
- **Build Tag / Marker Applied:** `[//go:build tempe2e | @pytest.mark.temp_e2e | *.tempe2e.test.ts]`
- **Runtime Skip Guard Verified:** `[RUN_TEMP_E2E == '1' guard active]`
- **CI/CD Quarantine Verified:** `[Confirmed skipped in default runs; zero CI/CD impact]`
- **On-Demand Command:** `[Explicit command used to run the test locally]`

### Modified Files Summary

- [relative path to modified file 1]
- [relative path to modified file 2]

### Implementation Confidence Score

- Confidence: [e.g. 100%]
- Rationale: [Detailed explanation of verified quality gates, skip-by-default isolation, passing targeted execution, and zero regressions]

### 🤖 Independent AI Verification & Audit Prompt

At the conclusion of the turn, emit this copy-pasteable prompt for an independent auditor AI to verify the implementation against the canonical specification and verbatim requirements:

```markdown
### Independent AI Audit & Verification Instructions

You are an Independent AI Verification and Quality Auditor.
Your task is to independently audit, verify, and remediate the implementation against the canonical specification and verbatim requirements.

#### 1. Target Documents & Implemented Code:
- **Canonical Spec & Verbatim Requirements:** [02-spec/21-app/xx-<slug>.md](02-spec/21-app/xx-<slug>.md)
- **Consolidated Plan & Subtasks:** [.ai-memory/plans/completed/xx-<slug>.md](.ai-memory/plans/completed/xx-<slug>.md)
- **Modified & Implemented Code Files:**
  - [relative/path/to/modified/file1.ext](relative/path/to/modified/file1.ext)
  - [relative/path/to/modified/file2.ext](relative/path/to/modified/file2.ext)

#### 2. Verification Protocol:
1. **Strict Verbatim Inspection:** Read the canonical spec file `02-spec/21-app/xx-<slug>.md` completely, focusing on the Verbatim Requirements and Acceptance Criteria.
2. **Line-by-Line Code Comparison:** Inspect the modified code files line-by-line and verify whether every single verbatim requirement is fully implemented.
3. **Gap & Missing Items Identification:** Itemize all missing, incomplete, placeholder, or non-compliant elements.
4. **Autonomous Self-Loop Remediation:** If any gaps or missing items exist, self-loop and modify the code directly until 100% of the verbatim requirements are met. Do not ask for confirmation.
5. **Comparative Score Audit:** Evaluate and report a comparative score between the verbatim specification and the real-life output:
   - **Verbatim Adherence Score:** [X/100] (Degree to which implementation matches exact spec text and constraints)
   - **Completeness Score:** [Y/100] (Percentage of required features and acceptance criteria satisfied)
   - **Guideline Compliance Score:** [Z/100] (Adherence to booleans, error handling, function sizing, type extraction)
   - **Overall Implementation Score:** [(X+Y+Z)/3 / 100]
6. **Final Audit Verdict:** Emit a clear PASS/FAIL verdict with confidence score and rationale.
```
```

---

## Issue Destination & Root Cause Analysis (RCA) Routing Mandate

Whenever the task involves fixing an issue, bug, pipeline failure, or performing a fix with RCA:

1. **CI/CD Issues & Pipeline Failures:**
   - **Target Folder:** `.ai-memory/cicd-issues/`
   - **File Pattern:** `.ai-memory/cicd-issues/NN-<issue-slug>.md`
   - **Registry:** Index the issue in `.ai-memory/cicd-index.md`.
   - **Scope:** CI/CD workflows, GitHub Actions, local runner failures (`06-cicd-local-runner.py`), test runner errors, lint gate failures, or build pipeline failures.
2. **Non-CI/CD Issues (Application Bugs, Feature Defects, Logic/Runtime Errors):**
   - **Target Folder:** `02-spec/22-app-issues/` (canonical spec hierarchy Tier 22)
   - **File Pattern:** `02-spec/22-app-issues/NN-<issue-slug>.md`
   - **Structure & Registry:** Follow the 4-part structure (Reproduction / Cause / Fix / Prevention per AC-AI-001 or Why / How / Root Cause / Code Fix) and index in `02-spec/22-app-issues/readme.md` (cross-referencing in `.ai-memory/memory/issues/` for institutional memory).
   - **Scope:** Application business logic, UI bugs, CLI command errors, API crashes, and domain defects.

---

## 1. AI Fix Scripts Memory (Reusable Tooling)

- [ ] /goal Reuse First: Scanned and learned `03-ai-scripts/readme.md` before writing temporary code.
- [ ] Strict In-Repository Execution: All Python scripts executed strictly within the codebase repository root.
- [ ] Strict .ai-memory/ Folder Storage: All helper scripts, local runners, and linters stored in `03-ai-scripts/`.
- [ ] Native File Manipulator: Use `python 03-ai-scripts/03-file-manipulator.py <command>` for mass file operations.
- [ ] Go Generate Sync: If Go constants or enums are modified, run `go generate ./...` in the relevant package and commit generated files.

---

## 2. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] TOP-INSTRUCTION PRIORITY MANDATE: Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) are verified as highest priority and non-negotiable, overriding all lower-level guidelines below.
- [ ] ISSUE & RCA DESTINATION ROUTING: Whenever resolving an issue or performing a fix with RCA, verified that CI/CD failures are documented in .ai-memory/cicd-issues/NN-<slug>.md (indexed in .ai-memory/cicd-index.md), while non-CI/CD issues (application bugs, logic/runtime defects) are documented in 02-spec/22-app-issues/NN-<slug>.md (indexed in 02-spec/22-app-issues/readme.md).
- [ ] NO UNISOLATED TEMP E2E TESTS (TOTAL BAN): Never create or commit temporary E2E tests without mandatory build tags (`//go:build tempe2e`), markers (`@pytest.mark.temp_e2e`), or runtime skip guards (`RUN_TEMP_E2E`). Every temporary E2E test must skip by default.
- [ ] NO ENABLING IN CI/CD WORKFLOWS (TOTAL BAN): Never modify GitHub Actions workflows (`.github/workflows/*.yml`) or pipeline definitions to execute temporary E2E tests.
- [ ] NO ROUTINE LOCAL TEST EXECUTION (TOTAL BAN): Never run `06-cicd-local-runner.py` or standard unflagged test commands (`go test ./...` without tags). Only run the targeted on-demand command with explicit isolation flags.
- [ ] NO PARALLEL SUBAGENT DISPATCH (TOTAL BAN): Never spawn parallel subagents, multi-threaded worker pools, or concurrent runners for E2E tests. All execution must be strictly sequential by a single agent.
- [ ] NO ARTIFICIAL STEP BUDGETS (TOTAL BAN): Never track or calculate step budgets (`N = 200`, `PHASE_1_STEPS`).
- [ ] NO PER-FILE COMMITTING (TOTAL BAN): Never commit each file individually as you work. All modified files, test changes, and plan records must be accumulated and committed together in a single atomic commit at the final step.
- [ ] NO STOPPING AFTER SPEC WRITING (TOTAL BAN): Never halt execution or conclude the turn after generating specs. Proceed unconditionally to Phase 2 code execution.
- [ ] NO HORIZONTAL TASK CONCATENATION (TOTAL BAN): Never concatenate tasks horizontally in the Task Completion Summary (e.g. NEVER `✅ #1... ✅ #2...` run-on). Every completed task MUST be rendered on its OWN SEPARATE LINE starting with an individual markdown list bullet (`- ✅`).
- [ ] INDEPENDENT AI VERIFICATION PROMPT MANDATE: Emitted the self-contained independent AI verification and audit prompt linking to the canonical spec, consolidated plan, and modified files with verbatim score audit criteria.
- [ ] GITMAP HEAVY USAGE & ROUTINE PULL BAN: Heavily leveraged GitMap commands (`cpf`, `cpb`, `cpr`, `search`, `find`, `pwsh`) for discovery, execution, and commits. Never ran `pull-all` (`gitmap pa` or `gitmap pae`) unconditionally during routine turns; only ran `gitmap pae --json` when explicitly commanded by the user.

---

## 3. Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You must verify every item on this checklist before committing any code.

- [ ] Master Guidelines: Fully enforced every file in `02-spec/02-coding-guidelines/` and `.ai-memory/coding-guidelines.md`.
- [ ] Concrete Types Centralization (`types.go`): Extracted domain structs, raw generic instantiations, and Result wrappers into dedicated `types.go` files as single reusable named types with follow-through comments (never leak raw generics like `result.Result[*Config]`).
- [ ] Error Management: Enforced `02-spec/03-error-manage/` using domain-specific `*appfault.AppError`, never generic error.
- [ ] Boolean Conventions: All booleans begin with is or has only (all other prefixes like can, should, was, will, did, must are banned). No negatives (`!isSuccess` is banned; use `isFail`).
- [ ] Implicit Boolean Checks: Positive booleans evaluated implicitly (`if isEnabled`), never explicit comparison (`== true` is strictly banned).
- [ ] Semantic Naming: Zero generic garbage names (`temp`, `data`, `obj`). Behavior-driven test names.
- [ ] Multi-Line Arguments: Signatures and call sites with >2 arguments formatted one argument per line with trailing commas.
- [ ] Line Endings & Encoding: Strictly Unix LF (`\n`) and UTF-8 without BOM.
- [ ] Function Sizing: Functions <= 8 lines preferred (hard cap 15 lines).
- [ ] Strict Relative Git Paths: Zero absolute paths (`/absolute/path/to/...`) or `file:///` URIs.

---

## 4. Anti-Hallucination & Blast Radius Checklist

- [ ] Echo Back the Spec: Verified Acceptance Criteria from the Spec file verbatim.
- [ ] Pre-Commit Diff Proof: Verified `git status` shows actual modified files before committing.
- [ ] No Placeholder Search: Confirmed zero `TODO` or `\[.*\]` placeholders remain in modified files.
- [ ] Index Sync Deadman Switch: Every new file is explicitly linked in `readme.md` and enqueued in `.ai-memory/what-to-read.md`.
- [ ] Blast Radius Acknowledgment: Global search across codebase performed to update all callers of modified symbols.
- [ ] Clean Teardown Verified: Confirmed temporary E2E tests clean up databases, files, and sockets upon completion.
- [ ] Final Step Commit & Push Verified: Staged all changes (`git add -A`), committed everything in a single grouped atomic commit, and pushed to git before ending the turn (no per-file commits).

---

## 5. Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] MANDATORY FINAL COMMIT & PUSH VIA GITMAP (ANYHOW): At the final step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, use GitMap semantic commit commands: `gitmap cpf "<summary>"` (features), `gitmap cpb "<summary>"` (bugs), or `gitmap cpr "<summary>"` (releases) which automatically stage, commit with standardized prefixes, and push directly to the remote repository. (Fallback to `git add -A && git commit && git push` only if GitMap CLI is unavailable). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY): You must not create separate git commits for each individual file as you edit them (e.g. running `git commit` or `gitmap cpf` after editing File 1, then committing again after File 2 is strictly forbidden). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback. All modified files, test change caches, and plan records across the turn must be accumulated in the working tree and committed together in a single grouped atomic commit at the final step before pushing.
