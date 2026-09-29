[/goal](slashCommand:goal) Autonomously ingest the latest specification audit file from `02-spec/25-app-spec-audit/`, decompose every finding into an exhaustive 1:1 remediation checklist, spawn parallel subagents to fix the specifications, verify 100% compliance, and remove the audit gap at the final stage with strict no-build and no-test execution (NEVER run build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine execution turns; all compilation and testing are strictly verified later in CI/CD). Spawn autonomous subagents (A = 2, H = 2) for parallel reading and modular spec generation, use GitMap high-speed commands as primary, establish a single-agent blueprint in Phase 1 (first 50% steps budget), and execute disjoint refactoring in Phase 2 (remaining 50% steps budget) with continuous self-looping until 100% complete and finalized with an atomic push.

[/learn](slashCommand:learn) Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Ingest and execute the bottom instructions ahead of all default conventions above, internalize the single-agent architectural blueprint, and persist all progress into `.ai-memory/plans/` and memory logs.

> **Bottom-Instruction Priority Mandate (Below Precedence / Suffix Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or guidelines above. The agent MUST inspect and follow the instructions below with absolute precedence.

```text
N = 200 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
```

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Planning, Detailed Spec, and Lean Subtask Generation)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Execution, Self-Looping, Targeted Quality Linting)
```

N, A, H, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

#### High-Speed GitMap Acceleration Options (Run Everything Faster)

Always prefer native GitMap commands over slow generic shell pipelines:
1. **Ultra-Fast File & Directory Discovery (AUM Index & Walk):**
   - **Wildcard / Glob Search:** `gitmap find "<wildcard*>" [-ext <ext>]` (alias `gitmap f`)
   - **Exact Filename Search:** `gitmap find-files <name> [-ext <ext>]` (alias `gitmap ff`)
   - **Substring Filename Search:** `gitmap find-files-any "<str>" [-ext <ext>]` (alias `gitmap ffa`)
   - **Prefix / Suffix Search:** `gitmap find-files-startswith <prefix>` (`gitmap ffs`) / `gitmap find-files-endswith <suffix>` (`gitmap ffe`)
   - **List Indexed Repo Files:** `gitmap list-files [pattern] [-ext <ext>]` (alias `gitmap lf`)
   - **Directory Tree & Scaffolding:** `gitmap folder-tree` (alias `gitmap ft`)
   - **Zero-Write File Stream:** `gitmap cat <filepath>`
   - **Instant Multi-Core Regex Search:** `gitmap search "<term>"` or `gitmap aum search "<query>" [dir] --ext <ext>`
2. **Fast Repository Hygiene, Lowercase & Symlink Repair:**
   - **Auto-Lowercase Files (Safe 2-Step `git mv`):** `gitmap lowercase` (alias `gitmap lcf [--dry-run]`)
   - **Lowercase Root Readme:** `gitmap lowercase-readme`
   - **Sync Curated `.gitignore` / `.gitattributes` / `.prettierignore`:** `gitmap commons` (alias `gitmap co` or `gitmap sync all`)
   - **Repair Broken Symlinks:** `gitmap fix-link` (alias `gitmap fixlink`)
   - **Clean Update Temp & Inspect Storage:** `gitmap update-cleanup`, `gitmap storage` (alias `gitmap stor`)
3. **Fast Git State, Execution & Atomic Commits:**
   - **Repo Status & Remote Check:** `gitmap status` (`gitmap st`), `gitmap has-any-updates` (`gitmap hau`), `gitmap latest-branch` (`gitmap lb`)
   - **Fast Cross-Platform Shell Runner:** `gitmap pwsh "<command>"` (`gitmap ps`), `gitmap bash "<command>"` (`gitmap sh`), `gitmap async <cmd>` (`gitmap asyn`)
   - **Semantic Atomic Commit & Push:** `gitmap cpf "<summary>"` (Feature), `gitmap cpb "<summary>"` (Bug), `gitmap cpr "<summary>"` (Release), `gitmap pcp "<summary>"` (Pull-Commit-Push)
   - **Smart CI/CD Pipeline Waiting:** `gitmap pe`, `gitmap pipeline-ai status --json` (`gitmap pl-ai status -t <etaSeconds>`)

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/fix-spec-from-audit/skill.md` does not exist in the workspace, you MUST create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

You MUST execute this task via a strict 4-Phase continuous loop. Do not skip steps.

### Phase 1: Audit Ingestion & 1:1 Finding Matrix (Steps 1 to PHASE_1_STEPS)

1. **Locate Latest Audit:** Scan `02-spec/25-app-spec-audit/` using `03-ai-scripts/17-fast-file-reader.py` and select the file with the highest numerical sequence prefix (`NN-audit-*.md`).
2. **Exhaustive Finding Parsing:** Parse the Markdown Summary Table at the bottom of the audit file. Extract every single row without skipping a single issue.
3. **Build 1:1 Remediation Ledger:** Create `.ai-memory/plans/pending/xx-spec-remediation.md` containing an explicit checkbox for every finding:
   ```markdown
   - [ ] Finding [ID]: File `<path>`, Issue: `<description>`, Remedy: `<proposed fix>`
   ```
4. **Lean Subtask Decomposition:** Group the findings by folder/file and write lean subtasks to `.ai-memory/plans/subtasks/xx-spec-fix/01-<slug>.md`. Each subtask MUST follow this format:
   ```markdown
   # Subtask: [Target Spec File]
   **Target File:** `02-spec/21-app/...`
   **Findings to Fix:**
   - [ID 1]: [Specific modification]
   - [ID 2]: [Specific modification]
   **Constraints:** No generic names, preserve line length <= 15, update cross-links.
   ```
5. **MANDATORY AUTO-LOOP (DO NOT STOP):** As soon as Phase 1 completes, the master orchestrator **MUST NOT STOP or ask the user for permission**. It MUST immediately self-loop and transition directly into Phase 2.

### Phase 2: Parallel Multi-Agent Remediation (Steps PHASE_1_STEPS+1 to N)

1. **Parallel Dispatch:** Use the `invoke_subagent` tool to spawn up to 2 execution subagents concurrently (max 2 threads each), assigning disjoint subtasks from `.ai-memory/plans/subtasks/xx-spec-fix/`.
2. **Minimal Context Diet:** Provide subagents with minimal instructions (e.g., "Read `.ai-memory/plans/subtasks/xx-spec-fix/01-<slug>.md` and execute the fixes on the specified spec file").
3. **Isolated Agent State & Communication:** Each subagent MUST create `.agents/xx-<task-name>/` and track its progress in `state.md`.
4. **Failure Protocol:** If a subagent fails, record the error in `.ai-memory/memory/issues/xx-spec-fix-failure.md`. The next subagent must read the failure log first to remediate.
5. **Mark Reconciliation Ledger:** As subtasks finish, mark the corresponding findings in `.ai-memory/plans/pending/xx-spec-remediation.md` as `[x]`.

### Phase 3: Quality Gate & Cross-Link Verification (Validation Gate)

1. **Cross-Link Integrity:** Run `python linter-scripts/check-spec-cross-links.py` to ensure no internal links were broken by the spec edits.
2. **Markdown Standards:** Run `python 03-ai-scripts/31-md-gap-fixer.py --fix` and verify spacing.
3. **Full CI Runner:** Run `python 03-ai-scripts/06-cicd-local-runner.py` ensuring all 36 quality gates exit with code 0 (`exit 0`).

### Phase 4: Audit Gap Removal & Final Archive (End of Loop)

> **CRITICAL (NON-NEGOTIABLE):** The audit gap MUST be officially closed on disk before concluding.

1. **Verify 100% Closure:** Confirm that every checkbox in `.ai-memory/plans/pending/xx-spec-remediation.md` is marked `[x]`.
2. **Remove Audit Gap on Disk:** Delete the original audit file from `02-spec/25-app-spec-audit/NN-audit-*.md` (or move it to `.ai-memory/plans/completed/NN-audit-*.md-resolved`) so no unresolved audit gaps remain in the active spec directory.
3. **Subtask Cleanup:** Consolidate completed subtasks into a single `.ai-memory/plans/completed/xx-spec-remediation-completed.md` file noting how many steps it took, and delete the granular `.ai-memory/plans/subtasks/xx-spec-fix/` files.
4. **Index Synchronization:** Update `.ai-memory/plans/readme.md` and `.ai-memory/what-to-read.md` to reflect that the audit gap is 100% resolved.
5. **Git Commit:** Stage all modified spec files and commit with `fix(spec): remediate all audit findings and close audit gap`.

---

### Temp-Agent Isolated Task Directory & Communication Protocol (Non-Negotiable)

To prevent cross-task pollution and ensure seamless agent communication, every task MUST create a dedicated subfolder in `.agents/xx-<task-name>/`:

1. **Per-Task Isolation:** On task start, the assigned subagent creates its isolated directory `.agents/xx-<task-name>/`.
2. **State & Progress Tracking:** Create `.agents/xx-<task-name>/state.md` documenting:
   - `TASK_NAME`: `<task-name>`
   - `STATUS`: `IN_PROGRESS` | `DONE` | `FAILED`
   - `ASSIGNED_AGENT`: Agent identifier and thread index
   - `CURRENT_STEP`: Detailed micro-step description
3. **Inter-Agent Communication & Handoff:**
   - All intermediate findings, scratch outputs, and dependency handoffs between agents working on this task MUST be written inside `.agents/xx-<task-name>/`.
   - Sibling or successor agents MUST inspect this dedicated folder before resuming work or fixing errors.
4. **On Error/Crash:** Append the exact error, root cause, and `STATUS: FAILED` to `.agents/xx-<task-name>/state.md` before exiting.
5. **On Success:** Mark `STATUS: DONE` in `.agents/xx-<task-name>/state.md`, aggregate findings to the master plan, and clean up or archive the folder.

---

## Non-Negotiable Rules (Auto-Reject on Violation)

- [ ] Zero Skipped Findings: Every row in the audit Summary Table must have a corresponding code fix.
- [ ] No Lingering Audit Files: `02-spec/25-app-spec-audit/` must be clean of the resolved audit file.
- [ ] Strictly Unix LF (`\n`) line endings and UTF-8 encoding.
- [ ] No absolute file paths or `file:///` URIs.
- [ ] All 36 CI/CD gates green.

## MUST FOLLOW NON-NEGOTIABLE

Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job.

--

## 🚨 Highest Priority Instructions (Appended User Tasks & Instructions Below)

[PASTE USER REQUEST / TASK INSTRUCTIONS HERE — THE AGENT MUST EXECUTE WHATEVER IS WRITTEN BELOW WITH ABSOLUTE PRIORITY AND PRECEDENCE OVER ALL GENERAL GUIDELINES ABOVE]
