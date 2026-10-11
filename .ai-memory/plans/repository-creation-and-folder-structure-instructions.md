# Master Plan: Repository Creation and Folder Structure Instructions

## User Request (Verbatim)
```text
# Repository Creation and Folder Structure Instructions: high priority instruction, non-negotiable task

Okay, I'm not sure if you have run those and created the final repo in the Git Hub. I do see that you didn't. The first observation would be remove these two folders. So that is absolutely wrong. Create the name as Icon Studio-2026 and then V1. Okay, that should be the last name. And you have lots of inconsistencies in the file and how it was put together. So there should be no replacing. That's the first thing I think you should know. And everything you do, you should maintain a JSON about the commit information and where you're getting from, so that the AI can fix the final tune. And also you create a final instruction for the AI, so that AI can take a decision how the folder structure should be final one. The final one should be smoother one, where the same folder, same file should be. If the same one, then it will consolidate the information, things like that. Or keep both if you find any differences, but that would be a different case. It could verify with human. Also, I think you did not respected that. So you need to make sure that the Merge AI works properly and it can work from the text file. That's the first thing. So make sure that you put the text file to something, and then you give the name of the Icon Studio, something like this. That should be in the first quotation. You have given it in the double quotation. It should be Icon Studio 2026 V1. That should be the name of the repo. So make sure that you respect that. Try to do the commits as well, and you make final commit and final push to the Git as well. And the repo needs to be added to VS Code, Antigravity if Antigravity exists, Cursor if Cursor exists, GitHub Desktop if GitHub Desktop exists, and VS Code as well if VS Code exists. If none of those exist, then it would also share the suggestion to install each one of them as a sub-point three with the `gitmap` command, how it can be installed and settings can be added. Do you understand?

## slug: repository-creation-and-folder-structure-instructions
```

---

## 1. Executive Summary & Architecture Boundary
This plan coordinates the end-to-end execution of repository creation, multi-repo non-destructive amalgamation, AI tuning metadata generation, and cross-IDE ecosystem integration for **`Icon Studio 2026 V1`** (canonical slug: **`icon-studio-2026-v1`**).

### Canonical Specifications:
- `02-spec/21-app/repository-creation-and-folder-structure-instructions/01-architecture-spec.md`
- `02-spec/21-app/repository-creation-and-folder-structure-instructions/02-ide-ecosystem-and-execution-spec.md`

### Invariants:
1. **Slug Invariant:** Target name `"Icon Studio 2026 V1"` maps strictly to lowercase kebab-case slug `icon-studio-2026-v1`.
2. **Zero Replacement Invariant:** When files collide across repositories, no file is replaced or destroyed. Files are sequenced with 2-digit chronological prefixes (`01_`, `02_`, `03_`).
3. **Commit Tracking Invariant:** Maintain complete provenance in `merge-ai-manifest.json` including commit SHAs, dates, author info, branch names, and per-file origins.
4. **AI Consolidation Guide Invariant:** Generate `instruction.md` detailing how subsequent AI iterations can consolidate or preserve variants.
5. **Multi-IDE Registration Invariant:** Register the newly staged repository into VS Code (`projects.json`), Google Antigravity (`.gemini/config/projects/`), Cursor IDE (`projects.json`), and GitHub Desktop (`github <path>`). If any IDE is not installed, output clear Sub-Point 3 fallback commands with `gitmap`.

---

## 2. Discrete Task Breakdown & Allocation Matrix

| Task ID | Task Title | Subtask Document | Assigned Owner | Scope / Target Area | Status |
|:---|:---|:---|:---|:---|:---|
| **Task-01** | Workspace Sanitization & Folder Removal | `.ai-memory/plans/subtasks/repository-creation-and-folder-structure-instructions/01-workspace-cleanup-and-manifest-spec.md` | Lead Orchestrator | Remove old `icon-studio-2026` and `icon-studio-v1` folders | **DONE** |
| **Task-02** | Specification & Schema Authoring | `.ai-memory/plans/subtasks/repository-creation-and-folder-structure-instructions/01-workspace-cleanup-and-manifest-spec.md` | Spec Authors 01 & 02 | Author 01-architecture and 02-ide specs | **DONE** |
| **Task-03** | Merge-AI Engine Enhancement | `.ai-memory/plans/subtasks/repository-creation-and-folder-structure-instructions/02-merge-ai-engine-and-manifest-tracking.md` | Worker 01 | Enhance `cli/cmdmergeai/` for rich commit tracking | **IN PROGRESS** |
| **Task-04** | Repository Creation & Amalgamation Execution | `.ai-memory/plans/subtasks/repository-creation-and-folder-structure-instructions/03-repo-creation-and-ide-registration.md` | Worker 01 | Run `gitmap merge-ai` for `icon-studio-2026-v1` | **QUEUED** |
| **Task-05** | IDE Ecosystem Registration & Sub-Point 3 Fallbacks | `.ai-memory/plans/subtasks/repository-creation-and-folder-structure-instructions/03-repo-creation-and-ide-registration.md` | Worker 02 | Register in VS Code, Antigravity, Cursor, GitHub Desktop | **QUEUED** |
| **Task-06** | Verification, Memory Consolidation & Atomic GitMap Commit | `.ai-memory/plans/subtasks/repository-creation-and-folder-structure-instructions/04-consolidation-and-verification.md` | Lead Orchestrator | Linters, ledger update, `gitmap cpf` atomic push | **QUEUED** |

---

## 3. Wave Execution Schedule (A = 2, H = 2)

- **Wave 1 (Execution Wave):**
  - **Worker 01:** Executes Task-03 & Task-04 (Enhances `cli/cmdmergeai/merge_ai_inspect.go` and `merge_ai_manifest.go` with enhanced commit metadata, then runs `gitmap merge-ai "Icon Studio 2026 V1" merge-ai-icon-studio.txt`, stages and makes initial commit in `icon-studio-2026-v1`).
  - **Worker 02:** Executes Task-05 (Inspects installed IDEs, registers `icon-studio-2026-v1` in VS Code, Antigravity, Cursor, and GitHub Desktop, and outputs the Sub-Point 3 installation guide).
- **Phase 3 (Consolidation & Push):**
  - Lead verifies targeted linters, completes SQLite task DB, consolidates plans into `.ai-memory/plans/completed/`, and performs atomic GitMap commit and push.
