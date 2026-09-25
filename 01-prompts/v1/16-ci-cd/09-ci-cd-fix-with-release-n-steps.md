# [V2] Release-Triggered CI/CD Fix N-Step Continuous Loop & Automated Release Ceremony — Workflow (must follow)

Trigger Keywords & Aliases: `cicd fix release n steps`, `ci release n steps`, `fix and release v2`, `ci release loop v2`, `09-ci-cd-fix-with-release-n-steps`

> [!IMPORTANT]
> Prompt Version: 2.5.0
> Synchronization: Main Meta-Repo & Connected Workspaces
> 
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Whatever directives, constraints, checklists, or instructions are given before this section or prompt (including in the prompt preamble, header blocks, or incoming user request) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

/goal Autonomously diagnose, repair all CI/CD pipeline issues using local runner scripts, execute targeted quality gates to verify 100% green status, and perform full automated release publication with version bump and changelog synchronization using an N-step continuous self-loop.

```text
N = 200
```

N = total self-loop steps budget that the agent will perform.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Pipeline Discovery, 4-Part RCA, Diagnostic Subtasks, and Local Runner Setup)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Surgical Fix Execution, Targeted Verification, Automated Release Ceremony)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

> [!CAUTION]
> **AUTOMATED RELEASE WARNING:** This workflow culminates in an automated release publication. Once all targeted quality gates exit code 0, the release orchestrator WILL bump the SemVer version, create a dedicated release branch, tag the commit, merge to main, and push to origin.

---

## Fast File Discovery & Diagnostic Toolchain (Mandatory Acceleration)

To rapidly locate failing pipeline definitions, broken source files, test fixtures, and error logs without hitting 50-result tool caps, the AI agent MUST utilize the diagnostic toolchain:
- **Remote Pipeline AI Status (<50ms):** `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`)
- **Remote Dynamic Timeout Wait:** `gitmap pipeline-ai status -t <etaSeconds>` (or alias `gitmap pl-ai status -t <sec>`)
- **Extract Failing Step Error Logs:** `gitmap pipeline error-logs` (or alias `gitmap pe`, clear with `gitmap pe clear -y`)
- **Pipeline Runner Targets & Cache Table:** `gitmap pipeline details` (or alias `gitmap pd`)
- **Scan Source & Test Files:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts,py --limit 100 --stats`
- **Fast Cached Pattern Search (<15ms):** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<error-or-symbol>" --limit 50`
- **Sub-Millisecond Folder Listing & Reader:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder .github/workflows --limit 20`
- **Read Workflow or Log File:** `python 03-ai-scripts/17-fast-file-reader.py --read-file .github/workflows/ci.yml`
- **Codebase Topology Overview:** `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`
- **Record Modified Files Under Lock:** `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/ci-cd-fix-with-release/skill.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

Execute this task via a strict 4-Phase pipeline. Do not skip steps.

### Phase 1A: Verbatim Capture, Failure Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, spec writing, or code changes, you must execute Phase 1A:

1. **Top-Instruction Priority Verification:** Whatever instructions, constraints, or checklists are given before this section or in the incoming user prompt must be ingested first as highest priority and non-negotiable.
2. **Verbatim Prompt & Stack Trace Capture:** Capture the incoming user request and failure logs verbatim into `.ai-memory/cicd-issues/xx-<slug>.md` under a dedicated `## User Request & Error Log (Verbatim)` section.
3. **Screenshot & Base64 Failure Logs Ingestion Protocol:** If the user request contains a screenshot URL, terminal snapshot, or base64 data URI (e.g. `data:image/png;base64,...`):
   - Convert that base64 encoding or downloaded image to the file system immediately, saving it as a persistent file under `assets/screenshots/<slug>-<NN>.png` or `assets/ui/<slug>-<NN>.png`.
   - Never leave raw, massive base64 strings or external ephemeral URLs in the prompt text, RCA logs, or subtasks.
   - Refer back to this saved image file in the RCA document and subtasks strictly as a relative markdown link (e.g. `![Failure Screenshot](assets/screenshots/<slug>-<NN>.png)`).
4. **Actionable Failure Extraction:** Break down all reported pipeline errors, failing tests, and lint violations into discrete, actionable items with ordered traceable IDs (`Error-01`, `Error-02`, etc.) under `## Extracted Actionable Error List`.
5. **Mandatory Chat Output Gate & Same-Turn Tool Chaining (TOTAL BAN ON CLOSING CONVERSATION):**
   - You MUST output the confirmed deliverables list directly in chat, but you MUST NOT end your turn or pause after emitting it!
   - In the EXACT SAME RESPONSE turn as your markdown output, you MUST immediately invoke the first tool call (e.g. `run_command` to query `gitmap pipeline error-logs` or `write_to_file` to initialize the RCA log).
   - NEVER emit the breakdown text without invoking a tool call. Emitting visible text without a tool call terminates the AI turn, causing the engine to close or hang the conversation while errors remain unresolved.
   - TOTAL BAN ON WAITING FOR APPROVAL: Do not pause, do not wait for the user to reply, and do not ask "Should I proceed?". Override any planning mode stop directives. Proceed immediately into Phase 1B in the same turn.
   - Line-by-Line Output Format Structure:

```markdown
### 📋 Confirmed CI/CD Failure Breakdown & Diagnostic Ingestion

1. **Error-01: [Descriptive Failure Title / Job Name]**
   - **State:** `[IN PROGRESS — DIAGNOSING IMMEDIATELY]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of the root cause hypothesis, intent, and verified constraints]
   - **Actionable Scope:** [Precise package, function, or configuration requiring surgical remediation]
   - **Target Files / Area:** `[relative/path/to/failing_file.go]`

2. **Error-02: [Descriptive Failure Title / Job Name]**
   - **State:** `[QUEUED — EXECUTING SURGICAL FIX NEXT]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of failure signature and expected resolution]
   - **Actionable Scope:** [Precise package or linter rule requiring repair]
   - **Target Files / Area:** `[relative/path/to/target]`

Proceeding directly to Phase 1B: RCA Generation & Runner Verification (Active Tool Call Running Below).
```

MANDATORY SAME-TURN TOOL CHAIN: The breakdown text above and your first tool call MUST be emitted in the EXACT SAME TURN. Never end the turn with text alone.

---

### Phase 1B: Grounded 4-Part RCA & Diagnostic Subtasks (Steps 1 .. N/2)

#### Step 1: Grounded 4-Part Root Cause Analysis (RCA) & Issue Routing

Author the structured RCA post-mortem following the mandatory issue destination routing:
- **CI/CD Issues & Pipeline Failures:** Record the RCA in `.ai-memory/cicd-issues/xx-<slug>.md` and index it in `.ai-memory/cicd-index.md`.
- **Non-CI/CD Issues (Application Bugs, Feature Defects, Logic/Runtime Errors):** If during the diagnosis the failure is identified as an application bug or domain defect rather than a pipeline script/workflow issue, document it in `02-spec/22-app-issues/xx-<slug>.md` (indexed in `02-spec/22-app-issues/readme.md`, cross-referencing in `.ai-memory/memory/issues/`).

Follow the mandatory 4-part schema:
- **Part 1: The Failing Symptom & Exact Error Signature:** Verbatim stack trace snippet (strictly bounded: 5 preceding + 20 trailing lines).
- **Part 2: The Root Mechanism (Why it Broke):** Deep architectural diagnosis explaining why the failure occurred without speculation.
- **Part 3: The Surgical Remedy (Zero Side-Effects):** Specific code edits, types, or configuration lines to change.
- **Part 4: Quality Gate & Verification Proof:** Exact command to verify the fix (`python 03-ai-scripts/06-cicd-local-runner.py run-smart`, `--pkg <target>`, or `--changed-only`).

#### Step 2: Local Runner Synchronization (`03-ai-scripts/06-cicd-local-runner.py`)

- Cross-reference central CI/CD workflows with `03-ai-scripts/06-cicd-local-runner.py`. Ensure all jobs can be executed locally on the native host.

#### Step 3: Lean Diagnostic Subtask Decomposition

- Break down the fix plan into lean subtasks in `.ai-memory/plans/subtasks/xx-<slug>/01-<subtask>.md`.
- All file paths in subtasks MUST use strict relative Git paths (zero absolute paths or `file:///` URIs).

#### Step 4: Unconditional Zero-Question Execution Mandate (Total Ban on Stopping After RCA)

- **Strict 50/50 Time Budget:** RCA generation and subtask setup MUST complete within the first 50% of the budget (`PHASE_1_STEPS = N / 2`).
- As soon as Phase 1 completes, transition unconditionally into Phase 2 code execution without pausing or asking confirmation questions.

---

### Phase 2: Execution Mode & Singly-Done Self-Loop Fixing (Steps N/2+1 .. N)

1. **Singly-Done Autonomous Loop:** Zero in on ONE failing error per iteration. Fix root causes surgically, apply coding guidelines, and verify that specific gate before moving to the next.
2. **Smart Targeted Testing & Release Authority:**
   - **Priority Incremental Runner:** Run `python 03-ai-scripts/06-cicd-local-runner.py run-smart` (or alias `--smart`, `-s`) to build ONLY changed packages into OS temp and run Quad Runner.
   - **Specific Package Targeting:** Run/build ONLY packages and test functions directly cited in the failure stack trace: `python 03-ai-scripts/06-cicd-local-runner.py --pkg <target_package_or_file>`.
   - **Heatmap & Fast-Path Testing:** Use `--fast` to run only hot and warm tests based on `.ai-memory/test-heatmap.json`.
   - **Changed Packages from Last Git Hash:** Compare against `HEAD~1` (`git diff --name-only HEAD~1`) via `python 03-ai-scripts/06-cicd-local-runner.py --changed-only`.
   - **Strict Ban on Extraneous Runs:** NEVER run full test suites, spellcheckers, or unrelated checks.
3. **Remote CI/CD Pipeline Monitoring & Dynamic Waiting Protocol (GitMap Pipeline-AI):**
   - Query remote pipeline state via `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`).
   - Dynamic timeout wait: `gitmap pipeline-ai status -t <etaSeconds>` to eliminate credit waste.
   - Tight-loop polling (`gh run view` loops) is STRICTLY BANNED. Adaptive sleep intervals: ETA > 120s wait 20s-30s; 60s < ETA <= 120s wait 10s-20s; ETA <= 60s wait 5s-10s.
4. **Change Recording Under Lock:** Record all modified files to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`).
5. **Per-Task Agent Isolation Subfolders (`.ai-memory/temp-agents/xx-<task-name>/`):**
   - Maintain `state.md` to track error triage, assigned files, and intermediate diagnostic outputs.
6. **Total Ban on Interim Per-File Commits:**
   - NEVER commit individual files during the fix phase. Accumulate all fixes cleanly in the working tree. All changes will be committed atomically by the release orchestrator or in the final release ceremony.
7. **Strict Avoidance: Never Disable CI/CD or Linting Checks:**
   - NEVER disable, comment out, delete, or skip any CLI linting command, build step, or test suite.
   - NEVER add `|| true`, `continue-on-error: true`, `# nolint`, or `// eslint-disable`.

---

### Phase 3: Automated Release Ceremony & Branch Lifecycle

Once all targeted quality gates exit code 0 (`exit 0`), proceed immediately and autonomously to the release ceremony:

1. **Verify Clean Working Tree Pre-Flight:** Confirm all outstanding edits are staged or saved before branching.
2. **Execute Heavy-Lifting Release Orchestrator:**
   ```bash
   python 03-ai-scripts/29-release-orchestrator.py --tier <minor|patch|major> --scope "<Fix and Release summary>"
   ```
3. **Mandatory 5-Step Release Branching Lifecycle:**
   - **Step 1:** Create and switch to release branch: `git checkout -b release/vX.Y.Z`.
   - **Step 2:** Bump SemVer using the repository-adapted Python bump script (`03-ai-scripts/37-bump-version.py`). Update `version.json`, `package.json`, `readme.md`, `changelog.md`, and sync.
   - **Step 3:** Commit on release branch: `release: vX.Y.Z <scope>`.
   - **Step 4:** Create annotated git tag on release commit: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`.
   - **Step 5:** Merge back to `main`, push `main`, `release/vX.Y.Z`, and tag `vX.Y.Z` to `origin`, then safely restore the starting branch.
4. **Clean up & Consolidate:** Consolidate subtasks into `.ai-memory/plans/completed/xx-<slug>.md` and mark RCA `RESOLVED`.

---

### End-of-Turn Verification & Confidence Reporting (Mandatory Output)

At the completion of all tasks and before concluding the turn, you must emit this structured verification summary in the chat response:

```markdown
### CI/CD Fix & Release Completion Summary

✅ #1. Error-01: [Error description] — Fixed & Verified
✅ #2. Error-02: [Error description] — Fixed & Verified
🚀 #3. Release Ceremony: vX.Y.Z — Published & Pushed to Origin

### Release Details

- Release Version: vX.Y.Z (Tier: MINOR / PATCH)
- Release Branch: `release/vX.Y.Z`
- Release Tag: `vX.Y.Z`
- Target Branch Restored: `main` (verified)

### Modified Files Summary

- [relative path to modified file 1]
- [relative path to modified file 2]

### Verification & Quality Gate Proof

- Gate Command: `python 03-ai-scripts/06-cicd-local-runner.py run-smart`
- Exit Code: 0 (PASS)
- Implementation Confidence Score: 100%
- Rationale: [Detailed explanation of verified quality gates, passing targeted tests, contract adherence, clean release publication, and zero regressions]
```

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] TOP-INSTRUCTION PRIORITY MANDATE: Whatever is given before this section or in the user prompt is verified as highest priority and non-negotiable.
- [ ] ISSUE & RCA DESTINATION ROUTING: Whenever resolving an issue or performing a fix with RCA, verified that CI/CD failures are documented in .ai-memory/cicd-issues/NN-<slug>.md (indexed in .ai-memory/cicd-index.md), while non-CI/CD issues (application bugs, logic/runtime defects) are documented in 02-spec/22-app-issues/NN-<slug>.md (indexed in 02-spec/22-app-issues/readme.md).
- [ ] NO INTERIM PER-FILE COMMITTING (TOTAL BAN): Never commit each file individually during the fix loop. All fixes must be accumulated and committed atomically by the release orchestrator.
- [ ] NO RAPID CI/CD POLLING (TOTAL BAN): Never query or loop rapidly (`gh run view` in tight loops). Use GitMap Pipeline-AI (`gitmap pl-ai status -t <sec>`) and sleep based on `etaSeconds`.
- [ ] NO FULL TEST SUITE RUNS (TOTAL BAN IN ROUTINE FIXES): Never run entire heavy test suites. Run tests strictly on affected packages cited in stack traces or changed from git hash (`HEAD~1`).
- [ ] NO DISABLING CI/CD GATES (TOTAL BAN): Never bypass, comment out, or add `|| true` to force a pipeline to pass.
- [ ] NO STOPPING AFTER RCA (TOTAL BAN): Never halt execution or ask user permission after writing the RCA. Proceed unconditionally to Phase 2 code execution.
- [ ] MANDATORY RELEASE PUSH TO GIT (ANYHOW): At the completion of the release, verify that the release branch, tag, and merged main branch are all pushed to origin.
