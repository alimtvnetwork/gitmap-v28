# Specification: Stash Pop Untracked File Collision Fix

## Root Cause Analysis
When `gitmap fix` executes the "Stash & Re-apply" option, it generates the following recipe steps:
1. `git stash -u` (Stash all tracked and untracked changes)
2. `git pull` (Fetch and merge remote changes)
3. `git stash pop` (Restore stashed changes)

A collision occurs when a locally untracked file (e.g., `file.sh`) is stashed, and the subsequent `git pull` brings down the same file from the remote as a tracked file. When `git stash pop` attempts to restore the stashed untracked file, Git refuses to overwrite the newly tracked file. 
The command fails with the following specific error:
```
error: could not restore untracked files from stash
<file-path> already exists, no checkout
```

Despite this error, Git successfully restores all other non-colliding untracked files and tracked modifications from the stash to the working directory. The only remainder in the stash is the colliding file, which is kept because of the abort.

## Recovery Strategy
Because the file was successfully pulled from the upstream repository, the local working directory has the tracked, correctly versioned copy. 
The recovery strategy is to safely **drop the stash entry** since its contents are either identical to the pulled file or safely discarded in favor of the upstream version.

### Logic Definition
1. **Detect Collision**: Inside the error handler (`handleStepFailure` in `cli/cmd/fix_execute.go`), evaluate if the failed command is `git stash pop`.
2. **Match Error String**: Check if the command output contains both:
   - `"could not restore untracked files from stash"`
   - `"already exists, no checkout"`
3. **Extract Path**: Extract the repository path from the step arguments (e.g., `step.Args[1]` if `step.Args[0] == "-C"`).
4. **Execute Recovery**: Run `git -C <repoPath> stash drop`.
5. **Success Override**: If `stash drop` succeeds, output a recovery message (`recovered (dropped stashed untracked collision)`) and mark the failed step as successful (`return nil`). If it fails, fallback to the standard failure output.
