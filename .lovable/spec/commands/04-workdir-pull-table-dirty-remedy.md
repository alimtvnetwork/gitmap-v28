# Specification: WorkDir Management, Recursive Top-Level Pull, Rich Pull Table & Dirty Remediation

## 1. Overview
This specification details:
1. **WorkDir Command Suite (`gitmap workdir` / `gitmap wd`)**:
   - `gitmap workdir ls` / `list`: Displays registered work paths, labels, repo counts, and highlights the default work directory (first scanned path or explicitly set default).
   - `gitmap workdir set-default <path|id>`: Overrides the active default work directory.
   - `gitmap workdir add <path> [--label <label>]`: Registers a new work directory.
   - `gitmap workdir rm <path|id>`: Removes a registered work directory.

2. **Recursive Top-Level Git Repository Auto-Discovery (`gitmap pull` & `gitmap push`)**:
   - When executed in a directory that is not a Git repo:
     - Checks if current directory is a registered work directory or has child repositories.
     - Recursively traverses subdirectories to discover all Git repositories.
     - **Pruning Rule**: Once a directory contains `.git`, it is recorded as a repository and its subdirectories are pruned (never recursing into nested sub-repositories or submodules).
     - Automatically pulls/pushes all discovered top-level repositories.

3. **Rich Pull Table with Last SHA, Branch, PR Indicator**:
   - The pull progress / summary table displays:
     - **Repo Name**
     - **Current / Latest Branch**
     - **Last Commit SHA** (short 7-char hash)
     - **PR / Upstream Status** (e.g. `PR Open`, `Up to date`, `Ahead 2`, `Behind 1`)
     - **Pull Result** (SUCCESS / SKIPPED / CONFLICT)
     - **Elapsed Time**

4. **Dirty Status Diagnosis & Actionable Remediation Hints**:
   - When a repository is dirty or cannot be pulled cleanly:
     - Analyzes uncommitted/untracked status in detail (e.g., `+1 modified`, `+2 untracked`).
     - Outputs exact, copy-pasteable remediation command recipes:
       - *Option 1 (Stash)*: `cd <repo> && git stash && git pull && git stash pop`
       - *Option 2 (Commit)*: `cd <repo> && git add -A && git commit -m "wip" && git pull`
       - *Option 3 (Discard)*: `cd <repo> && git reset --hard HEAD && git clean -fd && git pull`

5. **Release & CI/CD Verification**:
   - Version bump and local test suite validation across packages.
