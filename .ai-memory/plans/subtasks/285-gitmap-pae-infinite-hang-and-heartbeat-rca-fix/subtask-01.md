# Subtask 01: PAE Bar Teardown, Bounded Ignore Collect & Non-Interactive Remediation

## Status: DONE
- **Parent Plan**: `.ai-memory/plans/285-gitmap-pae-infinite-hang-and-heartbeat-rca-fix.md`
- **Owner**: Worker 01
- **Owned Files**:
  - `cli/cmdpull/pull_efficient.go`
  - `cli/cmdpull/pull_efficient_types.go`
  - `cli/cmdpull/pull_ignore_scan.go`

## Objectives
1. In `cli/cmdpull/pull_efficient.go`:
   - In `executeActiveEfficientBatch`:
     ```go
     executePull(partition.ActiveRecords, bar, pullOptions{all: true, useSSH: opts.UseSSH, useHTTPS: opts.UseHTTPS})
     stopHeartbeat()
     if !opts.IsJSON {
         bar.Stop()
     }
     collectAndRemediateEfficientIgnore(ignoreHandle, opts.Args, opts.IsJSON)
     return handleEfficientBatchFinish(total, bar, partition, allRecords, opts, time.Since(startTime))
     ```
   - Keep deferred calls as safety guards against panics, but invoke immediately upon pull completion. (DONE)
2. In `cli/cmdpull/pull_ignore_scan.go`:
   - In `Collect()`:
     Add a 3-second timeout select:
     ```go
     func (h *IgnoreScanHandle) Collect() []IgnoreRepoIssue {
         if h == nil {
             return nil
         }
         select {
         case issues := <-h.done:
             return issues
         case <-time.After(3 * time.Second):
             return nil
         }
     }
     ```
     (DONE)
3. In `cli/cmdpull/pull_efficient_types.go`:
   - In `collectAndRemediateEfficientIgnore`:
     If issues are found, do not prompt interactive stdin. Auto-skip or print a clean guidance note unless explicit `-y` was provided. (DONE)
4. FastGate: Ensure all modified Go files stay strictly <= 300 lines with nesting depth <= 1.
   - `cli/cmdpull/pull_efficient.go`: 298 lines (depth <= 1)
   - `cli/cmdpull/pull_efficient_types.go`: 285 lines (depth <= 1)
   - `cli/cmdpull/pull_ignore_scan.go`: 169 lines (depth <= 1)
   (DONE)
