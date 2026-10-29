# Subtask 02: Letterly Recommit & Time-Restricted Prompt and V7 Execution Engine Implementation

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/pending/348-recommit-and-skill-creation-instructions.md`  
> **Target Path:** `.ai-memory/plans/subtasks/348-recommit-and-skill-creation-instructions/02-letterly-and-execution-prompts.md`  
> **Component Spec:** `02-spec/21-app/348-recommit-and-skill-creation-instructions/02-component-spec.md`  
> **Status:** READY  
> **Estimated Effort:** 3 Engineering Units  
> **Dependencies:** `02-spec/21-app/348-recommit-and-skill-creation-instructions/01-architecture-spec.md`, `02-spec/21-app/348-recommit-and-skill-creation-instructions/02-component-spec.md`

---

## 1. Objectives & Scope

This subtask governs the discrete execution plan for authoring the new Letterly prompt formatter (12), authoring the new V7 execution prompt (09), registering the new prompt in the Letterly index, and establishing on-call skill linkages.

### 1.1 Target Deliverables
1. **Letterly Recommit & Time-Restricted Formatter:**  
   `01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly.md`
2. **V7 Parent Task N-Step Execution Engine:**  
   `01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md`
3. **Letterly Readme Registration:**  
   Update `01-prompts/22-letterly/readme.md` to register entry 12 with canonical successor `execute-parent-task-with-n-steps-v7`.
4. **V6 Immutability Preservation:**  
   Strict verification that `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md` remains completely unmodified.

### 1.2 Non-Negotiable Constraints
- **Relative Git Paths Only:** Strictly use relative repository paths (`01-prompts/...`, `02-spec/...`, `.agents/skills/...`). Absolute paths are strictly banned.
- **Positive Booleans Only:** All parameters, structural flags, and status variables must use `is*` or `has*` prefixes (`isTimeRestricted`, `hasRecommitWindow`, `isAutobotEnabled`, `isValid`).
- **On-Call Skill Linkages:** Skills `gitmap-time-restriction` and `gitmap-automation-bot` must be wired on an "on-call basis" (activated conditionally when prompt directives or flags require them, never forced into general flows).
- **Leaf Error Sovereignty:** Any evaluation helpers or utility functions must return errors directly without invoking exit handlers or `os.Exit`.
- **Function Line Limits:** Any code snippets or helpers must remain <= 15 lines (<= 8 lines preferred).

---

## 2. Step-by-Step Execution Plan

```
+-------------------------------------------------------------------------------+
| Step 1: Pre-Execution Validation & V6 Immutability Baseline                   |
| - Verify git status is clean                                                  |
| - Record baseline hash of 01-prompts/14-execute/02-execute-...-v6.md          |
+-------------------------------------------------------------------------------+
                                        |
                                        v
+-------------------------------------------------------------------------------+
| Step 2: Author 01-prompts/22-letterly/12-execute-recommit-time-restricted-... |
| - Structure lossless verbatim capture template                                |
| - Implement standard Actionable Items 1-8                                     |
| - Implement Time-Restriction & Recommit Mandates 9-11                         |
| - Embed on-call skill links (.agents/skills/gitmap-time-restriction, etc.)    |
| - Set successor trigger: execute-parent-task-with-n-steps-v7                  |
+-------------------------------------------------------------------------------+
                                        |
                                        v
+-------------------------------------------------------------------------------+
| Step 3: Author 01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7   |
| - Adopt multi-agent swarm architecture (A = 3, H = 2, C = 30, Cap = 6)       |
| - Implement Phase 0: Antigravity Skill Bootstrap                              |
| - Implement Phase 1A: Chat Showcase Gate & Understood Mandate                 |
| - Implement Phase 1B: On-Call Skill Resolution Gate                           |
| - Implement Phase 1C: Brain Planning & Roadmap                                |
| - Implement Phase 2: Tree Swarm Parallel Execution (disjoint files)           |
| - Implement Phase 3: Time-Window Validation, Recommit Gate, & Atomic Commit    |
+-------------------------------------------------------------------------------+
                                        |
                                        v
+-------------------------------------------------------------------------------+
| Step 4: Synchronize Letterly Directory Readme (01-prompts/22-letterly/readme) |
| - Append Row 12 to the directory index table                                  |
| - Verify rows 01-11 are untouched                                             |
+-------------------------------------------------------------------------------+
                                        |
                                        v
+-------------------------------------------------------------------------------+
| Step 5: Retrospective Quality Audit & Verification                            |
| - Verify zero diff on V6 prompt                                               |
| - Check relative path compliance across all modified files                    |
| - Confirm positive boolean naming standard                                    |
+-------------------------------------------------------------------------------+
```

---

### Step 1: Pre-Execution Validation & V6 Immutability Baseline
1. Check repository status via GitMap CLI to ensure working tree readiness:
   `gitmap ps` or `git status --porcelain`.
2. Inspect the current hash and state of `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`.
3. Establish pre-flight assertion: any edit touching `02-execute-parent-task-with-n-steps-v6.md` is an immediate test failure.

### Step 2: Authoring Letterly Formatter 12
- **File:** `01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly.md`
- **Actions:**
  1. Header and preamble definition:
     - Clear directive: format voice dictation into non-negotiable execution template.
     - Lossless capture rules: strip conversational filler (`um`, `uh`, `like`, leading `so`, `okay`) while strictly keeping parameters, flags, and file paths.
  2. Subject derivation and markdown structure:
     - Line 1: `# <task title>: high priority instruction, non-negotiable task`
     - Line 2: Blank
     - Line 3: `${Input Text Verbatim}`
     - Metadata: `## task: <task title>`, `## slug: <task-slug>`
  3. Actionable Items Specification:
     - Items 1–8: Standard invariants (pull first, spec/plan first, GitMap search primacy, relative paths, SQL repository pattern, enum action dispatch, enum authoring checkup, error management).
     - Item 9: Time-Restriction Enforcement (On-Call activation of `gitmap-time-restriction`).
     - Item 10: Recommit & History Adjustment (On-Call activation for re-evaluating commits and applying `gitmap` semantic commit verbs).
     - Item 11: AI Automation Bot Identity (On-Call activation of `gitmap-automation-bot` for author credential switches).
     - Item 12..N: Discrete user directives from input.
  4. On-Call Skills Section:
     - Relative link: `[gitmap-time-restriction](file;.agents/skills/gitmap-time-restriction)`
     - Relative link: `[gitmap-automation-bot](file;.agents/skills/gitmap-automation-bot)`
  5. Suffix Invocation:
     - `[execute-parent-task-with-n-steps-v7](file;.agents/skills/execute-parent-task-with-n-steps-v7)`
  6. Additional Instructions:
     - `/plan` and `/learn` recommendations.
     - GitMap Heatmap live sync command (`gitmap rc heatmap status` / `gitmap rc heatmap live`).
     - Relative path invariant enforcement.

### Step 3: Authoring V7 Parent Task Execution Engine
- **File:** `01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md`
- **Actions:**
  1. Swarm Constants & Concurrency:
     - `N = 400`, `A = 3`, `H = 2`, `C = 30`, Capacity = 6.
     - Budgets: Phase 1 = 100, Phase 2 = 150, Phase 3 = 150.
  2. Phase 0: Antigravity Skill Auto-Bootstrap:
     - Ensure `.agents/skills/execute-parent-task-with-n-steps-v7/skill.md` is populated with YAML frontmatter.
  3. Phase 1A: Chat Showcase Gate & Understood Mandate:
     - Turn 1 visible chat showcase of Strategic Intent, Boundaries, Constraints, Blast Radius, Concurrency.
     - 6-part task verification schema (`Understood / AI Comprehension`, `Model Runtime Context`).
     - Lossless verbatim capture under `## User Request (Verbatim)`.
  4. Phase 1B: On-Call Skill Resolution Gate (V7 Innovation):
     - Inspect incoming request for time-restriction keywords (`time restrict`, `time window`). If present, bind `gitmap-time-restriction`.
     - Inspect request for automation bot keywords (`autobot`, `automation bot`). If present, bind `gitmap-automation-bot`.
     - Inspect request for recommit keywords (`recommit`, `amend`). If present, activate recommit audit workflow.
     - If neither is present, proceed with standard clean workflow without extraneous overhead.
  5. Phase 1C: Up-Front Planning & Roadmap:
     - Author unified Brain roadmap in `<appDataDir>/brain/<convo-id>/roadmap.md`.
     - Create lean subtask plans in `.ai-memory/plans/subtasks/XX-<slug>/`.
  6. Phase 2: Tree Swarm Parallel Execution:
     - Dispatch 3 concurrent worker subagents (`invoke_subagent`).
     - Strict disjoint file box isolation (zero cross-worker file collisions).
     - Immediate dynamic join on subtask completion (zero idle workers, 100% duty cycle).
  7. Phase 3: Verification, Pacing & Atomic GitMap Commit:
     - Time-window compliance check (if on-call time restriction is active).
     - Recommit history verification (if on-call recommit is active).
     - Final atomic GitMap commit (`gitmap cpf`, `gitmap cbf`, `gitmap cfr`, `gitmap cpc`, `gitmap cpr`).

### Step 4: Synchronizing Letterly Directory Readme
- **File:** `01-prompts/22-letterly/readme.md`
- **Actions:**
  1. Read `01-prompts/22-letterly/readme.md`.
  2. Locate Directory Index table.
  3. Append Row 12:
     `| **12** | [`12-execute-recommit-time-restricted-letterly.md`](12-execute-recommit-time-restricted-letterly.md) | Recommit & Time-Restricted | Formats voice dictation with on-call time restrictions, recommit instructions, autobot identity, and V7 execution engine | `execute-parent-task-with-n-steps-v7` |`
  4. Ensure all preceding rows (01–11) are preserved unchanged.

### Step 5: Retrospective Quality Audits & Verification
- **Actions:**
  1. Execute `git diff --name-only` to verify modified files match expectations.
  2. Confirm `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md` is NOT in the diff list.
  3. Scan authored files for illegal absolute path prefixes (`C:`, `D:`, `file:///C:`, `file:///D:`).
  4. Validate all boolean identifiers conform to `is*` or `has*`.

---

## 3. Risk Assessment & Mitigation Strategies

| Risk Description | Severity | Probability | Mitigation Strategy |
| :--- | :--- | :--- | :--- |
| **Accidental V6 Modification** | HIGH | LOW | Explicit pre-flight baseline hash check; post-authoring git status assertion. |
| **Absolute Path Leakage** | MEDIUM | LOW | Systematic relative path check across markdown links; ban `C:\` and `D:\`. |
| **Unconditional Skill Bloat** | MEDIUM | LOW | Specifying clear "on-call" conditional criteria so skills are only invoked when requested. |
| **Negative Boolean Ingestion** | LOW | LOW | Enforce strict naming rules (`isTimeRestricted`, `hasRecommitWindow`). |

---

## 4. Acceptance Criteria & Definition of Done

- [ ] `01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly.md` is authored with complete sections and on-call skill linkages.
- [ ] `01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md` is authored with multi-agent architecture ($A=3$, $H=2$, Cap=6) and On-Call Skill Gate.
- [ ] `01-prompts/22-letterly/readme.md` table includes entry 12 pointing to `execute-parent-task-with-n-steps-v7`.
- [ ] `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md` has zero git diff.
- [ ] All file references across newly authored documents use relative Git repository paths only.
- [ ] All boolean identifiers adhere to positive naming (`is*`, `has*`).
