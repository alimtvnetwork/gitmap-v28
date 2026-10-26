# Subtask 333.1 — Specification Authoring & Architectural Design

> **Task Code:** `Task-01`  
> **Title:** Specification Authoring & Architectural Design  
> **Parent Plan:** `.ai-memory/plans/pending/333-clone-and-fix-git-repositories.md`  
> **Specification:** `02-spec/21-app/333-clone-and-fix-git-repositories/01-architecture-spec.md`  
> **Status:** COMPLETED  
> **Owned Files:**  
> - `02-spec/21-app/333-clone-and-fix-git-repositories/01-architecture-spec.md`  
> - `.ai-memory/plans/subtasks/333-clone-and-fix-git-repositories/01-subtask-architecture-and-specs.md`  

---

## 1. Objectives & Execution Summary

1. **Grounded Spec Authoring:**
   - Authored `02-spec/21-app/333-clone-and-fix-git-repositories/01-architecture-spec.md` incorporating the full verbatim User Request.
   - Diagnosed and detailed Root Cause 1: Bare repository slugs routed blindly to `executeManifestCloneFlow` which triggered `cloner.LoadRecords` and crashed with `open clone source "yt-manage-v1": The system cannot find the file specified` at `cli/cloner/cloner.go:94`.
   - Diagnosed and detailed Root Cause 2: `gitmap cfr <name>` collided with `cmdcommit.RunCommitPushRefactor` in `runCfrDisambiguated` (`cli/cmd/roottooling_util.go:177`), failing when run outside a Git repo (`d:\work`).
   - Diagnosed and detailed Root Cause 3: Omission of Antigravity IDE registration in `cli/cmdclone/clone_direct.go:207-221`.
2. **Unified Architecture Design:**
   - Designed non-file bare argument detection vs manifest detection (`fsutil.FileExists` guard, manifest extension check).
   - Designed 4-tier resolution engine: Direct URL -> Local File/Manifest -> SQLite `store.DB` (`Repo`, `OwnerRepoNameIndex`) -> GitHub CLI/API (`gh repo list`).
   - Designed interactive suggestion display and fuzzy did-you-mean matching across `clone`, `reclone`, `cfr`, `crfb`, and `crf`.
   - Designed post-clone automated Antigravity IDE synchronization via `workspacesync.SyncAntigravity` and `cmdide.SyncAntigravityPinned`.
3. **Coding Guidelines Governance:**
   - Enforced Leaf Error Sovereignty, zero nested conditionals, 15-line function limit, positive booleans (`is*`/`has*`), and strictly relative paths.

---

## 2. Acceptance Criteria & Verification Evidence

| Criteria ID | Requirement | Verification Evidence | Status |
|---|---|---|---|
| **AC-01** | Full verbatim User Request captured | Preserved verbatim in `01-architecture-spec.md` Section 1 and parent plan. | **PASS** |
| **AC-02** | Root Cause Analysis grounded in codebase | Explicitly traced to `cloner/cloner.go:94`, `roottooling_util.go:177`, and `clone_direct.go:207`. | **PASS** |
| **AC-03** | Multi-tier resolution & suggestion design | Defined 4 tiers in Section 3 with Mermaid flowchart and sequence diagram. | **PASS** |
| **AC-04** | Disambiguation & alias routing | Outlined `cfr`, `crfb`, and `crf` alias dispatch in Section 3.4. | **PASS** |
| **AC-05** | Antigravity IDE integration contract | Detailed three-tier registration protocol in Section 3.5. | **PASS** |
| **AC-06** | Coding guidelines compliance | Full compliance matrix documented with leaf error sovereignty and positive booleans. | **PASS** |

---

## 3. Strict Coding Guidelines Compliance
- **Error Sovereignty:** Leaf functions return structured `error` (`*apperror.AppError`); outer dispatchers handle process exits.
- **Function Sizing:** Max <= 15 lines per function.
- **Boolean Naming:** `is*` and `has*` only; no negative booleans.
- **Relative Path Hygiene:** Pure relative paths (`02-spec/...`, `.ai-memory/...`, `cli/...`) throughout all specifications and plans.
