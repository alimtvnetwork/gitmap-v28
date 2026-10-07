---
name: commit-and-push-all-repos
description: Autonomously discover, stage, commit, and push all Git repositories across the target workspace (e.g. d:\work, ~/git-work), enforcing clean working trees, non-owned repository exclusion (omis, zsh), pre-commit pooling, and the mandatory 'No Push = Not Done' upstream completion invariant.
---

# Multi-Repository Commit & Push Orchestrator

Autonomously scans the entire workspace root directory (e.g. `D:\work` on Windows or `~/git-work` on POSIX/Linux), audits working tree status across dozens of Git repositories, excludes non-owned tools/shells, pools remote updates, stages dirty files, authors atomic conventional commits, and pushes to remote tracking branches with a strict "No Push = Not Done" completion invariant.

## Core Capabilities

1. **Cross-Platform Workspace Discovery & Non-Owned Exclusions:**
   - Auto-detects workspace roots across Windows (`D:\work`, `C:\work`), Linux/macOS (`~/git-work`, `~/work`), environment variables (`WORK_DIR`), or explicit `--dir` arguments.
   - Accurately filters out transient build directories (`target`, `node_modules`, `dist`, `vendor`, `.cache`, `tmp`).
   - **Non-Owned & Third-Party Filtering (TOTAL BAN):** Automatically ignores third-party tools and shell configs not owned by us (`omis`, `oh-my-zsh`, `ohmyzsh`, `zsh`, `oh-my-posh`, `dotfiles`, package runtimes), with `--exclude <pattern>` custom regex support.

2. **Full-Fleet Working Tree & Divergence Audit:**
   - Identifies dirty, modified, untracked, ahead, behind, diverged, and detached HEAD states in milliseconds.
   - Non-destructive and safe: skips detached HEAD states to prevent orphan commit branches.

3. **Pre-Commit Pooling / Pulling:**
   - Non-destructively pulls latest remote changes (`git pull origin <branch> --no-rebase`) prior to committing dirty changes to incorporate upstream updates cleanly and eliminate merge conflicts.
   - Configurable override via `--no-pull` for offline or isolated runs.

4. **Atomic Conventional Commits:**
   - Cleans and stages all modifications (`git add -A`).
   - Automatically authors structured conventional commits (`chore(sync): ...`) with file counts or user-specified messages.

5. **Upstream Remote Synchronization:**
   - Pushes cleanly to origin, automatically binding unbound branches (`git push -u origin <branch>`).
   - Gracefully handles read-only external repositories without breaking the fleet execution run.
   - Supports non-interactive SSH authentication (`git@github.com`).

6. **Mandatory Completion Invariant ("No Push = Not Done"):**
   - Verifies all dirty working trees are committed and all local commits are confirmed pushed upstream to GitHub.
   - If any repository has uncommitted dirty changes or unpushed commits remaining, the task is strictly considered **INCOMPLETE and NOT DONE**, exiting with code `1`.

7. **Dry-Run, JSON, & Audit Modes:**
   - Supports `--dry-run` for safe preview, `--check` for read-only fleet auditing, and `--json` for programmatic integration.

## CLI Execution

```bash
# 1. Run internal self-tests
python 03-ai-scripts/49-commit-and-push-all-repos.py --self-test

# 2. Audit all repositories across default workspace (read-only, non-owned excluded)
python 03-ai-scripts/49-commit-and-push-all-repos.py --check

# 3. Pull, commit dirty files, and push all ahead repos across workspace
python 03-ai-scripts/49-commit-and-push-all-repos.py

# 4. Target explicit directory (e.g. D:\work or ~/git-work)
python 03-ai-scripts/49-commit-and-push-all-repos.py --dir "d:/work"

# 5. Skip pre-commit pull/pooling phase
python 03-ai-scripts/49-commit-and-push-all-repos.py --no-pull

# 6. Add custom exclusion pattern
python 03-ai-scripts/49-commit-and-push-all-repos.py --exclude "external-tool"

# 7. Preview changes safely without mutating disk or remotes
python 03-ai-scripts/49-commit-and-push-all-repos.py --dry-run

# 8. Push pending commits only without committing dirty trees
python 03-ai-scripts/49-commit-and-push-all-repos.py --push-only

# 9. Commit dirty working trees only without pushing
python 03-ai-scripts/49-commit-and-push-all-repos.py --commit-only

# 10. Commit and push with custom conventional message
python 03-ai-scripts/49-commit-and-push-all-repos.py -m "feat(sync): batch multi-repo fleet updates"

# 11. Output machine-readable JSON summary
python 03-ai-scripts/49-commit-and-push-all-repos.py --json
```
