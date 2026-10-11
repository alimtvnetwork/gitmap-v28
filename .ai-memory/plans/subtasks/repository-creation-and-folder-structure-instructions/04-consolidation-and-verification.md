# Subtask 06: Quality Gate, End-to-End Verification & Plan Consolidation

> **Document ID:** `.ai-memory/plans/subtasks/repository-creation-and-folder-structure-instructions/04-consolidation-and-verification.md`  
> **Parent Plan:** `.ai-memory/plans/repository-creation-and-folder-structure-instructions.md`  
> **Architecture Spec:** `02-spec/21-app/repository-creation-and-folder-structure-instructions/02-ide-ecosystem-and-execution-spec.md`  
> **Target Subsystems:** Quality Assurance, FastGate AST Linter, Path Sanitizer, Plan Archival  
> **Subtask Covered:** Task-06  
> **Status:** PENDING  

---

## 1. Objective & Scope

Subtask 06 represents the concluding quality gate and verification phase for milestone `repository-creation-and-folder-structure-instructions`. This phase systematically verifies that:
1. All repository creation, folder structuring, and multi-IDE synchronization components adhere strictly to architectural standards.
2. The entire codebase and documentation set satisfies 100% relative path hygiene with zero absolute path leaks.
3. Every schema, struct, and configuration variable adheres to positive affirmative boolean polarity.
4. All subtasks (Task-01 through Task-06) achieve verified completion and the parent execution plan is archived into `.ai-memory/plans/completed/`.
5. The workspace is prepared for clean atomic commit and upstream push synchronization under the mandatory "No Push = Not Done" completion invariant.

---

## 2. Concrete Verification Matrix & Execution Commands

| Verification ID | Category | Command / Audit Probe | Expected Outcome | Failure Remediation |
|---|---|---|---|---|
| **V-06-01** | **FastGate AST Linter** | `python 03-ai-scripts/50-fastgate.py` | Exit code 0, 0 AST errors, execution time $<1.5\text{s}$. | Resolve syntax, type, or import errors flagged by AST analyzer. |
| **V-06-02** | **Relative Path Hygiene** | `python 03-ai-scripts/49-verify-privacy-and-relative-paths.py` | Exit code 0, zero drive letters (`C:\`, `D:\`), zero absolute paths, zero file URI schemes. | Replace absolute path references with Git root-relative paths. |
| **V-06-03** | **Positive Boolean Audit** | `python 03-ai-scripts/26-boolean-logic-modernizer.py --audit` | Exit code 0, 100% affirmative booleans across Go and specs. | Invert negative boolean names (`skipDesktop` &rarr; `isDesktopSyncEnabled`). |
| **V-06-04** | **Folder Structure Audit** | `python 03-ai-scripts/38-sync-prompts-skills-scripts.py --check` | 100% alignment with `01-canonical-folder-structure.md` and `.ai-memory/folder-structure.md`. | Realign mismatched folder paths or missing taxonomy files. |
| **V-06-05** | **VS Code Registration Test** | `go test -v ./cli/vscodepm/... ./cli/cmdvscode/...` | All tests pass, `projects.json` writes succeed in memory and disk fixtures. | Fix path canonicalization or JSON marshaling in `vscodepm`. |
| **V-06-06** | **Antigravity Sync Test** | `go test -v ./cli/workspacesync/... -run TestSyncAntigravity` | Project JSON generated with RFC 4122 UUID, valid folder URI, and `TURBO` mode. | Verify UUID generator and JSON template in `workspacesync`. |
| **V-06-07** | **Cursor Sync Test** | `go test -v ./cli/cmdcursor/...` | Tests pass, Cursor configuration resolution handles Windows/Linux paths correctly. | Fix path resolver in `cli/cmdcursor/cursor_sync.go`. |
| **V-06-08** | **GitHub Desktop Test** | `go test -v ./cli/cmddesktopsync/... ./cli/desktop/...` | Tests pass, CLI resolution and `AddRepos` summary calculations correct. | Fix CLI resolution or summary counters in `desktop`. |
| **V-06-09** | **Fallback Dry-Run Test** | `gitmap install vscode --dry-run`<br>`gitmap install antigravity --dry-run`<br>`gitmap cursor install --dry-run`<br>`gitmap gd --install --dry-run` | All 4 commands emit clean expected package manager invocations with exit code 0. | Fix argument parsing or template strings in `cmdinstall` and `cmdcursor`. |
| **V-06-10** | **Working Tree Hygiene** | `git status --porcelain` | Clean working tree with only tracked, intended changes. | Remove any untracked temporary artifacts or logs. |

---

## 3. Full Acceptance Criteria Traceability Matrix

| AC ID | Specification Source | Requirement Description | Verification ID | Status |
|---|---|---|---|---|
| **AC-01** | `02-ide-ecosystem-spec.md` | VS Code Project Manager registration via `projects.json` | V-06-05 | PENDING VERIFICATION |
| **AC-02** | `02-ide-ecosystem-spec.md` | Google Antigravity project descriptor generation (`TURBO` mode) | V-06-06 | PENDING VERIFICATION |
| **AC-03** | `02-ide-ecosystem-spec.md` | Cursor IDE workspace configuration mirroring and sync | V-06-07 | PENDING VERIFICATION |
| **AC-04** | `02-ide-ecosystem-spec.md` | GitHub Desktop repository addition via CLI endpoint | V-06-08 | PENDING VERIFICATION |
| **AC-05** | `02-ide-ecosystem-spec.md` | Sub-Point 3 Fallback commands for missing tools | V-06-09 | PENDING VERIFICATION |
| **AC-06** | `02-ide-ecosystem-spec.md` | Graceful degradation on missing IDE binaries (`[tool: skipped]`) | V-06-05, V-06-08 | PENDING VERIFICATION |
| **AC-07** | `01-architecture-spec.md` | Git initialization on default branch `main` | V-06-04 | PENDING VERIFICATION |
| **AC-08** | `01-architecture-spec.md` | Canonical `02-spec/` layered folder structure scaffolding | V-06-04 | PENDING VERIFICATION |
| **AC-09** | `01-architecture-spec.md` | Canonical `.ai-memory/` taxonomy (`plans/`, `memory/`) | V-06-04 | PENDING VERIFICATION |
| **AC-10** | Core Invariants | Strict relative Git path hygiene (zero drive letters, zero URIs) | V-06-02 | PENDING VERIFICATION |
| **AC-11** | Core Invariants | Positive affirmative boolean naming standards | V-06-03 | PENDING VERIFICATION |
| **AC-12** | Release Protocol | FastGate AST check $<1.5\text{s}$, plan consolidation, atomic commit | V-06-01, V-06-10 | PENDING VERIFICATION |

---

## 4. Plan Consolidation & Archival Protocol

Upon successful completion of all verification checks:
1. **Update Subtask Statuses:** Ensure Task-01 through Task-06 statuses in `.ai-memory/plans/subtasks/repository-creation-and-folder-structure-instructions/` are updated from `PENDING` to `COMPLETED`.
2. **Synchronize Master Plan:** Update `.ai-memory/plans/repository-creation-and-folder-structure-instructions.md` marking the milestone complete.
3. **Archive to Completed Directory:** Move the master plan file to `.ai-memory/plans/completed/repository-creation-and-folder-structure-instructions.md`.
4. **Update Master Registry:** Update `.ai-memory/plans/readme.md` recording the archived plan under completed milestones.

---

## 5. Atomic Commit & Upstream Push Protocol

Follow the standard GitMap atomic commit protocol adhering to the "No Push = Not Done" rule:
1. **Stage Files:** Stage all created and updated specifications, plans, and code files.
2. **Atomic Commit:**
   ```bash
   gitmap cpf "feat(repo-creation,ide-sync): canonical folder scaffolding and multi-ide registration engine"
   ```
3. **Upstream Remote Push:** Ensure commit is pushed to the remote tracking branch immediately.
