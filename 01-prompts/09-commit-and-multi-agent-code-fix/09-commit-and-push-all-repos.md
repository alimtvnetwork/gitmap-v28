# Multi-Repository Commit & Push Automation — Workflow (must follow)

> **Prompt Version:** 1.0.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

[/goal](slashCommand;goal) Autonomously discover, verify, stage, commit, and push all Git repositories across the target workspace (e.g. `d:\work` on Windows, `~/git-work` on POSIX/Linux, or custom path), enforcing zero uncommitted data loss, branch immutability, atomic conventional commit formatting, upstream binding, and remote synchronization.

---

## Architectural Motivation & Why This Is Very Important

In distributed multi-agent workflows, developer workstations, and polyglot meta-repositories, dozens of Git repositories coexist in a single parent workspace root (such as `D:\work` on Windows or `~/git-work` on Linux). When working across projects:
1. **Uncommitted Work Stranding:** Developers or automated agents frequently create or edit files in multiple repositories without committing or pushing, leading to work loss, workspace drift, and divergence when switching machines.
2. **Cross-Platform Path Portability:** Scripts and prompts must seamlessly support Windows environments (`D:\work`, `C:\work`), POSIX/Linux paths (`~/git-work`, `/home/...`), and continuous integration environments without requiring hardcoded manual path changes.
3. **Build Artifact & Dependency Isolation:** Scanners must strictly distinguish genuine repository working trees from transient or nested build artifacts (`node_modules`, `target/`, `dist/`, `vendor/`, `.cache/`, `tmp/`).
4. **Upstream Safety & Monotonic History:** Commits must be clean, atomic, and structured with conventional semantics (`chore(sync): ...`). Pushes must respect upstream tracking branches (`git push -u origin <branch>` when unbound) and never force-push or rewrite published Git history.

---

## Step-by-Step Instructions for Multi-Repo Synchronization

Follow these structured steps sequentially to execute multi-repository audit, commit, and push operations:

### Step 1: Pre-Flight Workspace Root Auto-Discovery & Health Audit
- Identify the target workspace root directory:
  - If explicitly provided via `--dir` or parameter, validate existence.
  - Check environment variable `WORK_DIR`.
  - Check standard platform paths (`D:\work`, `C:\work`, `~/git-work`, `~/work`, or parent of current repository).
- Confirm Git executable availability and SSH / credential health:
  - Run `ssh -T git@github.com` or verify `git config url."git@github.com:".insteadOf "https://github.com/"` for non-interactive SSH authentication.

### Step 2: Multi-Repository Recursive Topology Discovery
- Recursively discover all top-level directories containing `.git`.
- **Exclusion Invariant:** Automatically ignore directories located within `node_modules`, `target`, `dist`, `build`, `vendor`, `.cache`, or `tmp`.
- Sort discovered repositories deterministically by relative path.

### Step 3: Git Status & Divergence Triage
For each discovered repository:
- Run `git status --porcelain` to identify untracked (`??`), modified (`M`), staged (`A`), or deleted (`D`) files.
- Determine the active branch: `git rev-parse --abbrev-ref HEAD`.
  - If in detached `HEAD` state, log a warning and skip automatic commit/push to prevent orphaned commits.
- Check upstream tracking ref: `git rev-parse --abbrev-ref @{u}`.
  - If upstream exists, calculate ahead count (`git rev-list @{u}..HEAD --count`) and behind count (`git rev-list HEAD..@{u} --count`).
  - If upstream does not exist, check if the repository has commits and locate default remote (e.g. `origin`).
- Categorize repository into `CLEAN`, `DIRTY`, `AHEAD`, `BEHIND`, `DIVERGED`, `DETACHED`, or `NO_COMMITS`.

### Step 4: Ephemeral & OS Cache Sanitization
- Prior to staging, ensure no unintended temporary files (e.g. misplaced `~/` directories, swap files, `.DS_Store`) exist in untracked lists.
- If unwanted artifacts exist, clean them or ensure they are added to `.gitignore`.

### Step 5: Atomic Staging & Meaningful Conventional Commit Creation
For each repository marked `DIRTY`:
- Stage all tracked and untracked changes:
  ```bash
  git add -A
  ```
- Generate a conventional commit message describing the changes:
  ```bash
  git commit -m "chore(sync): automated repository sync and working tree commit (<N> files)"
  ```
  *(Or use the user-supplied custom commit message if specified).*
- Verify commit creation and record new commit SHA (`git rev-parse --short HEAD`).

### Step 6: Upstream Reconciliation & Remote Tracking Push
For each repository marked `AHEAD` (or newly committed):
- Identify primary remote (`origin`).
- Push active branch cleanly:
  - If upstream tracking branch is already configured:
    ```bash
    git push origin <branch>
    ```
  - If branch has no upstream tracking set:
    ```bash
    git push -u origin <branch>
    ```
- If push is rejected due to remote divergence:
  - Fetch remote: `git fetch origin`.
  - Safely merge remote tracking branch non-destructively: `git merge --no-ff origin/<branch>`.
  - Re-verify and push.
- If push fails due to 403 authorization (e.g. read-only upstream forks), flag cleanly as `SKIPPED (External Repo)` without halting the remaining fleet.

### Step 7: Post-Push Verification Gate & Telemetry Summary
- Query all repositories to verify 0 dirty working trees and 0 ahead commits remaining.
- Render a concise structured report table displaying:
  - Repository relative path
  - Active branch
  - Initial status vs final outcome
  - Number of files committed and commits pushed
  - Any exceptions or skips

---

## Automated Execution via Python AI Script

Use the canonical Python automation script in `03-ai-scripts/49-commit-and-push-all-repos.py`:

```bash
# 1. Run internal self-tests to verify discovery, audit, and staging logic:
python 03-ai-scripts/49-commit-and-push-all-repos.py --self-test

# 2. Read-only status audit of all repositories in workspace:
python 03-ai-scripts/49-commit-and-push-all-repos.py --check

# 3. Commit and push all dirty and ahead repositories in workspace:
python 03-ai-scripts/49-commit-and-push-all-repos.py

# 4. Target explicit workspace directory (e.g. D:\work or ~/git-work):
python 03-ai-scripts/49-commit-and-push-all-repos.py --dir "d:/work"

# 5. Preview changes safely without modifying disk or remotes:
python 03-ai-scripts/49-commit-and-push-all-repos.py --dry-run

# 6. Push ahead commits only (skip committing dirty working trees):
python 03-ai-scripts/49-commit-and-push-all-repos.py --push-only

# 7. Provide custom commit message:
python 03-ai-scripts/49-commit-and-push-all-repos.py -m "feat(fleet): synchronize multi-repo core updates"

# 8. Output machine-readable JSON for automated CI/CD:
python 03-ai-scripts/49-commit-and-push-all-repos.py --json
```

---

## Action Items — Must Follow (Non-Negotiable)

- [ ] Discover all genuine repositories under workspace root, skipping build folders (`target`, `node_modules`, `dist`, `vendor`).
- [ ] Inspect git porcelain status, active branch, and remote divergence before mutating working trees.
- [ ] Strictly avoid force-pushing (`--force`, `-f`) or history rewrites on published branches.
- [ ] Stage all modifications (`git add -A`) and author atomic conventional commits (`chore(sync): ...`).
- [ ] Push to remote tracking branches, establishing `-u` upstream tracking if unbound.
- [ ] Handle read-only external forks and detached HEAD states gracefully without crashing.
- [ ] Execute post-flight verification verifying zero stranded dirty files or unpushed local commits.

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.
