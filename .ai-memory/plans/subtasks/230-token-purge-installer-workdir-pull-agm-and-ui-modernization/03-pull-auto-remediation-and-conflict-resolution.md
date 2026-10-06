# Subtask 03: GitMap Pull Auto-Remediation and Conflict Resolution

> **Parent Plan:** [230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md](../../pending/230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md)  
> **Spec Reference:** [02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/01-architecture-spec.md](../../../../02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/01-architecture-spec.md)  
> **Status:** `QUEUED`  
> **Target Files:**  
> - `cli/cloner/safe_pull.go`  
> - `cli/cmdpull/pull_remediation.go`  
> - `cli/cmdpull/pull_remediation_hint.go`  
> - `cli/cmdpull/pull_fallback.go`  

---

## 1. Technical Objective

Harden repository synchronization in `gitmap pull` (`gitmap pa`) to automatically resolve common git pull failure modes (divergent local commits, uncommitted local edits via autostash, and fast-forward divergence) while ensuring absolute working tree safety by aborting merge conflicts and providing an interactive batch remediation menu.

---

## 2. Implementation Specifications

1. **4-Stage Safe Pull Pipeline (`cli/cloner/safe_pull.go`):**
   - **Stage 1 (Fast-Forward Only):**
     Execute `git pull --progress --ff-only --autostash`. If exit code is 0, return success.
   - **Stage 2 (Divergence Detection):**
     Check for divergence indicators in stderr/stdout (`not possible to fast-forward`, `diverged`, `need to specify how to reconcile`).
   - **Stage 3 (Auto-Merge Fallback):**
     Execute `git pull --progress --no-rebase --no-edit --autostash`.
     If clean, return auto-merge success.
   - **Stage 4 (Conflict Detection & Safe Abort):**
     If merge conflicts occur (`automatic merge failed`), immediately execute `git merge --abort`.
     Ensure working tree remains in a clean state with zero leftover conflict markers.
     Return structured `ErrMergeConflict`.

2. **Interactive Remediation Prompt (`cli/cmdpull/pull_remediation.go`):**
   - If one or more repositories fail during a bulk pull operation:
     - Check `isInteractiveTerminal()`. If non-interactive (CI/CD), skip prompt and return summary error code.
     - In interactive sessions, invoke `PromptInteractiveBatchFix(failures)`:
       ```text
       [?] Remediate failed repositories?
           [1/a] Resolve all failed repositories at once
           [2/s] Step through one by one
           [q/n] Skip / Exit
       ```
   - For step-by-step remediation, present each failed repo with actionable options:
     - Force checkout & pull
     - Stash changes and retry
     - Open working directory in editor/terminal

3. **Error Reporting & Split-DB Persistence:**
   - Log pull diagnostics to SQLite `pipeline.db` error tables via `store/pull_split_db_errors.go`.
   - Render clean tabular error output via `cli/cmdpull/pull_table_format.go`.

---

## 3. Verification Protocol

```bash
# 1. Verify safe pull execution on clean repository
gitmap pull --fast

# 2. Simulate diverged branch and verify auto-merge fallback
# 3. Simulate merge conflict and verify automatic git merge --abort
# 4. Check error table rendering
gitmap pull-error ls
```

---

## 4. Acceptance Criteria

- [ ] All pull executions apply `--autostash` to protect uncommitted developer work.
- [ ] Non-fast-forward diverged branches automatically attempt non-rebase auto-merge.
- [ ] Conflicted merges trigger instant `git merge --abort`, preventing working-tree corruption.
- [ ] Interactive terminals offer `[1/a] All | [2/s] Step | [q/n] Skip` remediation menu.
- [ ] Headless/CI environments skip interactive prompts with standard exit codes.
- [ ] Code strictly follows positive booleans (`isSuccess`, `isDiverged`, `isInteractive`) and guard inversion.
