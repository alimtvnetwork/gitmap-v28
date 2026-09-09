# Subtask 03: Dirty Numbers Accuracy & Direct Remediation Commands

## Objective
Unify the porcelain status parser to guarantee accurate counting of staged, modified, untracked, and deleted files, and generate per-repository direct runnable shell commands in the remediation summary.

## Files to Touch
- `gitmap/gitutil/dirty_inspect.go`
- `gitmap/gitutil/dirty_inspect_test.go`
- `gitmap/gitutil/gitutil.go`
- `gitmap/gitutil/remediation_generator.go`
- `gitmap/cmd/remediation_box.go`

## Detailed Implementation Steps
1. In `gitmap/gitutil/dirty_inspect.go`:
   - Refactor `classifyDirtyFile`:
     - Inspect character `X` (index) and `Y` (worktree):
       - If `prefix == "??"`: `recordUntracked(diagnosis, filePath)`
       - If `prefix[0] != ' ' && prefix[0] != '?'`:
         - If `prefix[0] == 'D'`: `recordDeleted(diagnosis, filePath)`
         - Else: `recordStaged(diagnosis, filePath)`
       - If `len(prefix) > 1 && prefix[1] != ' ' && prefix[1] != '?'`:
         - If `prefix[1] == 'D'`: `recordDeleted(diagnosis, filePath)`
         - Else: `recordModified(diagnosis, filePath)`
   - In `collectReasonParts`:
     - Include `StagedCount`: `parts = append(parts, "+"+strconv.Itoa(diagnosis.StagedCount)+" staged")`
     - Order parts logically: staged, modified, untracked, deleted.
2. In `gitmap/cmd/remediation_box.go`:
   - In `printPendingReposList`, under each dirty repository, display the exact runnable commands:
     - Direct Git Command:
       `↳ Direct Git:    git -C "<repoPath>" stash -u`
     - Direct Gitmap Command:
       `↳ Gitmap Fix:    gitmap fix <repoName> 1`
3. Add unit tests in `gitmap/gitutil/dirty_inspect_test.go` and `gitmap/cmd/remediation_box_test.go`.
