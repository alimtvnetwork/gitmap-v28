---
name: commit-and-push-all-repos
description: Autonomously discover, stage, commit, and push all Git repositories across the target workspace (e.g. d:\work, ~/git-work), enforcing clean working trees, atomic conventional commits, and upstream remote synchronization.
---

# Multi-Repository Commit & Push Orchestrator

Autonomously scans the entire workspace root directory (e.g. `D:\work` on Windows or `~/git-work` on POSIX/Linux), audits working tree status across dozens of Git repositories, stages dirty files, authors atomic conventional commits, and pushes to remote tracking branches.

## Core Capabilities

1. **Cross-Platform Workspace Discovery:**
   - Auto-detects workspace roots across Windows (`D:\work`, `C:\work`), Linux/macOS (`~/git-work`, `~/work`), environment variables (`WORK_DIR`), or explicit `--dir` arguments.
   - Accurately filters out transient build directories (`target`, `node_modules`, `dist`, `vendor`, `.cache`, `tmp`).

2. **Full-Fleet Working Tree & Divergence Audit:**
   - Identifies dirty, modified, untracked, ahead, behind, diverged, and detached HEAD states in milliseconds.
   - Non-destructive and safe: skips detached HEAD states to prevent orphan commit branches.

3. **Atomic Conventional Commits:**
   - Cleans and stages all modifications (`git add -A`).
   - Automatically authors structured conventional commits (`chore(sync): ...`) with file counts or user-specified messages.

4. **Upstream Remote Synchronization:**
   - Pushes cleanly to origin, automatically binding unbound branches (`git push -u origin <branch>`).
   - Gracefully handles read-only external repositories without breaking the fleet execution run.
   - Supports non-interactive SSH authentication (`git@github.com`).

5. **Dry-Run, JSON, & Audit Modes:**
   - Supports `--dry-run` for safe preview, `--check` for read-only fleet auditing, and `--json` for programmatic integration.

## CLI Execution

```bash
# 1. Run internal self-tests
python 03-ai-scripts/49-commit-and-push-all-repos.py --self-test

# 2. Audit all repositories across default workspace (read-only)
python 03-ai-scripts/49-commit-and-push-all-repos.py --check

# 3. Commit dirty files and push all ahead repos across workspace
python 03-ai-scripts/49-commit-and-push-all-repos.py

# 4. Target explicit directory (e.g. D:\work or ~/git-work)
python 03-ai-scripts/49-commit-and-push-all-repos.py --dir "d:/work"

# 5. Preview changes safely without mutating disk or remotes
python 03-ai-scripts/49-commit-and-push-all-repos.py --dry-run

# 6. Push pending commits only without committing dirty trees
python 03-ai-scripts/49-commit-and-push-all-repos.py --push-only

# 7. Commit dirty working trees only without pushing
python 03-ai-scripts/49-commit-and-push-all-repos.py --commit-only

# 8. Commit and push with custom conventional message
python 03-ai-scripts/49-commit-and-push-all-repos.py -m "feat(sync): batch multi-repo fleet updates"

# 9. Output machine-readable JSON summary
python 03-ai-scripts/49-commit-and-push-all-repos.py --json
```
