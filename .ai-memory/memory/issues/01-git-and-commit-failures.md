# Issue Domain 01: Git and Commit Failures

- **Domain:** Git Operations, Rebase Conflicts, and Non-Git Directory Execution
- **Status:** Consolidated Problem & Resolution Matrix

## 1. Execution in Non-Git Directories
- **Symptoms:** `gitmap cm` crashed with fatal error E9000 when executed in root directories or non-git folders.
- **Root Cause:** Direct execution of `git rev-parse` without defensive pre-flight validation.
- **Resolution:** Added `isGitRepository()` guard returning structured `appfault.NewNotFoundError` with suggestions (`gitmap init` or navigation).

## 2. Detached HEAD and Diverged Tracking Branches
- **Symptoms:** Pull-all operations failed on detached HEAD or un-tracked local branches.
- **Root Cause:** Git pull defaults failed to resolve ambiguous upstream branches.
- **Resolution:** Implemented fast-forward verification and fallback rebase strategy in `cmdpas`.
