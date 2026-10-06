# Subtask 04: AGM Linux Update Quiet Mode and Failure Resolver

> **Parent Plan:** [230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md](../../pending/230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md)  
> **Spec Reference:** [02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/01-architecture-spec.md](../../../../02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/01-architecture-spec.md)  
> **Status:** `QUEUED`  
> **Target Files:**  
> - `cli/cmdinstall/agm_update.go`  
> - `cli/cmdinstall/installagmanager.go`  
> - `cli/cmdnodes/nodes_deploy_agm.go`  

---

## 1. Technical Objective

Refactor the Antigravity Manager (AGM) Linux update pipeline in `gitmap agm update` to eliminate stdout clutter on successful runs by executing in quiet mode, preserve full bounded stack traces and terminal error logs when an update fails, and provide an interactive failure resolution prompt when multiple updates fail across fleet machines.

---

## 2. Implementation Specifications

1. **Quiet Output Execution on Success (`cli/cmdinstall/installagmanager.go`):**
   - Intercept subprocess stdout and stderr in an in-memory buffer (`bytes.Buffer`) rather than streaming directly to `os.Stdout` / `os.Stderr` when `isUpdate = true`.
   - On successful exit (`err == nil`):
     - Discard verbose script trace.
     - Emit only a clean, minimal success badge:
       ```text
       ✓ Antigravity Manager updated successfully (latest).
       ```

2. **Bounded Stack Trace & Error Logging on Failure:**
   - On non-zero exit code (`err != nil`):
     - Print a prominent failure banner in red/yellow.
     - Dump the captured execution log buffer (bounded to the last 50 lines to prevent terminal overflow).
     - Include system context: architecture (`runtime.GOARCH`), OS (`runtime.GOOS`), target release version, and exit code.
     - Return structured `apperror.AppError`.

3. **Multi-Failure Interactive Prompt (`cli/cmdinstall/agm_update.go`):**
   - In batch update modes (`gitmap agm update-all` or remote fleet updates):
     - Collect all failed machine/node results into a `[]AgmUpdateFailure` slice.
     - If failures > 0 and `isInteractiveTerminal()` is true:
       - Present interactive choice prompt:
         ```text
         [?] Antigravity Manager update failed on 2 nodes (u1, w2).
             [1/a] Retry failed nodes
             [2/s] Show detailed failure logs for each node
             [q/n] Skip / Exit
         Choice [1/2/q]:
         ```
       - If non-interactive, print summarized error list and return non-zero exit code.

---

## 3. Verification Protocol

```bash
# 1. Simulate quiet update execution (dry-run)
gitmap agm update --dry-run

# 2. Verify command aliases and help documentation
gitmap agm update --help
gitmap agm update-all --help

# 3. Code review: verify stdout/stderr buffering in dispatchAgManagerUnixWithVersion
```

---

## 4. Acceptance Criteria

- [ ] Linux AGM update runs silently on success without streaming verbose curl or script traces.
- [ ] Successful update prints a single clean success line with green checkmark.
- [ ] Failed update dumps buffered execution log, exit code, and bounded stack trace.
- [ ] Fleet update failures trigger interactive `[1/a] Retry | [2/s] Inspect | [q/n] Skip` prompt when interactive.
- [ ] Positive booleans (`isUpdate`, `isSuccess`, `isInteractive`) and early return guard clauses enforced.
